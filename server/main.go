package main

import (
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

//go:embed site
var embeddedFS embed.FS

var siteFS fs.FS

// Config holds server configuration from environment variables
type Config struct {
	Lang      string            // LANG - current language
	Redirects map[string]string // REDIRECTS - map of lang -> domain
	Port      string            // PORT - server port
}

var cfg Config

func init() {
	// Create sub-filesystem rooted at "site"
	var err error
	siteFS, err = fs.Sub(embeddedFS, "site")
	if err != nil {
		log.Fatal(err)
	}

	cfg = Config{
		Lang:      getEnv("LANG", "en"),
		Port:      getEnv("PORT", "80"),
		Redirects: parseRedirects(os.Getenv("REDIRECTS")),
	}
	log.Printf("Config: LANG=%s, REDIRECTS=%v", cfg.Lang, cfg.Redirects)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// parseRedirects parses REDIRECTS=ru:domain-ru.com,de:domain-de.com
func parseRedirects(s string) map[string]string {
	m := make(map[string]string)
	if s == "" {
		return m
	}
	for _, pair := range strings.Split(s, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func main() {
	handler := gzipMiddleware(http.HandlerFunc(handleRequest))
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Printf("Listening on :%s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}

// gzip middleware
type gzipWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Pre-set Content-Type for known file extensions before gzip wrapping
		urlPath := r.URL.Path
		if ext := path.Ext(urlPath); ext != "" {
			if ct := mime.TypeByExtension(ext); ct != "" {
				w.Header().Set("Content-Type", ct)
			}
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		next.ServeHTTP(gzipWriter{Writer: gz, ResponseWriter: w}, r)
	})
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	// Set cache headers
	if isStaticAsset(urlPath) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=3600")
	}

	// Root path
	if urlPath == "/" {
		serveFile(w, r, path.Join(cfg.Lang, "index.html"))
		return
	}

	// Static assets (css, images) - check BEFORE language prefix
	// because paths like /css/ look like language codes
	if isStaticAsset(urlPath) {
		serveFile(w, r, urlPath[1:]) // Remove leading slash
		return
	}

	// Check if path starts with /$lang/
	if lang, suffix, ok := extractLangPrefix(urlPath); ok {
		if lang == cfg.Lang {
			// Current language - serve locally
			serveFile(w, r, path.Join(lang, suffix))
		} else if domain, exists := cfg.Redirects[lang]; exists {
			// Other language from REDIRECTS - redirect
			http.Redirect(w, r, fmt.Sprintf("https://%s/%s", domain, suffix), http.StatusMovedPermanently)
		} else {
			// Language not in REDIRECTS - 404
			http.NotFound(w, r)
		}
		return
	}

	// Default - current language
	serveFile(w, r, cfg.Lang+urlPath)
}

// extractLangPrefix extracts language from path /en/page -> ("en", "page", true)
func extractLangPrefix(urlPath string) (lang, suffix string, ok bool) {
	if len(urlPath) < 2 || urlPath[0] != '/' {
		return "", "", false
	}
	rest := urlPath[1:]
	slashIdx := strings.Index(rest, "/")
	if slashIdx == -1 {
		return "", "", false
	}
	lang = rest[:slashIdx]
	suffix = rest[slashIdx+1:]

	// Check if it looks like a language code (2-3 lowercase letters)
	if len(lang) < 2 || len(lang) > 3 {
		return "", "", false
	}
	for _, c := range lang {
		if c < 'a' || c > 'z' {
			return "", "", false
		}
	}
	return lang, suffix, true
}

func isStaticAsset(urlPath string) bool {
	prefixes := []string{"/css/", "/images/", "/assets/", "/media/"}
	for _, p := range prefixes {
		if strings.HasPrefix(urlPath, p) {
			return true
		}
	}
	return urlPath == "/favicon.ico"
}

func serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	// Clean the path and ensure no leading slash (fs.FS uses relative paths)
	filePath = path.Clean(filePath)
	if strings.HasPrefix(filePath, "/") {
		filePath = filePath[1:]
	}

	// Try to stat the file
	info, err := fs.Stat(siteFS, filePath)
	if err != nil {
		// Try with .html extension
		htmlPath := filePath + ".html"
		if _, err := fs.Stat(siteFS, htmlPath); err == nil {
			filePath = htmlPath
		} else {
			// Try as directory with index.html
			indexPath := path.Join(filePath, "index.html")
			if _, err := fs.Stat(siteFS, indexPath); err == nil {
				filePath = indexPath
			} else {
				http.NotFound(w, r)
				return
			}
		}
		// Re-stat the resolved path
		info, _ = fs.Stat(siteFS, filePath)
	}

	// If it's a directory, try to serve index.html
	if info != nil && info.IsDir() {
		indexPath := path.Join(filePath, "index.html")
		if _, err := fs.Stat(siteFS, indexPath); err == nil {
			filePath = indexPath
		} else {
			http.NotFound(w, r)
			return
		}
	}

	// Read file content from embedded filesystem
	content, err := fs.ReadFile(siteFS, filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Set Content-Type based on file extension
	ext := path.Ext(filePath)
	if ct := mime.TypeByExtension(ext); ct != "" {
		w.Header().Set("Content-Type", ct)
	}

	w.Write(content)
}
