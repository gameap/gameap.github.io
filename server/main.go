package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config holds server configuration from environment variables
type Config struct {
	Lang      string            // LANG - current language
	Redirects map[string]string // REDIRECTS - map of lang -> domain
	SiteRoot  string            // SITE_ROOT - site root directory
	Port      string            // PORT - server port
}

var cfg Config

func init() {
	cfg = Config{
		Lang:      getEnv("LANG", "en"),
		SiteRoot:  getEnv("SITE_ROOT", "/srv"),
		Port:      getEnv("PORT", "80"),
		Redirects: parseRedirects(os.Getenv("REDIRECTS")),
	}
	log.Printf("Config: LANG=%s, SITE_ROOT=%s, REDIRECTS=%v", cfg.Lang, cfg.SiteRoot, cfg.Redirects)
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
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		next.ServeHTTP(gzipWriter{Writer: gz, ResponseWriter: w}, r)
	})
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Set cache headers
	if isStaticAsset(path) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=3600")
	}

	// Check if path starts with /$lang/
	if lang, suffix, ok := extractLangPrefix(path); ok {
		if lang == cfg.Lang {
			// Current language - serve locally
			serveFile(w, r, filepath.Join(cfg.SiteRoot, lang, suffix))
		} else if domain, exists := cfg.Redirects[lang]; exists {
			// Other language from REDIRECTS - redirect
			http.Redirect(w, r, fmt.Sprintf("https://%s/%s", domain, suffix), http.StatusMovedPermanently)
		} else {
			// Language not in REDIRECTS - 404
			http.NotFound(w, r)
		}
		return
	}

	// Root path
	if path == "/" {
		serveFile(w, r, filepath.Join(cfg.SiteRoot, cfg.Lang, "index.html"))
		return
	}

	// Static assets (css, images)
	if isStaticAsset(path) {
		serveFile(w, r, filepath.Join(cfg.SiteRoot, path))
		return
	}

	// Default - current language
	serveFile(w, r, filepath.Join(cfg.SiteRoot, cfg.Lang, path))
}

// extractLangPrefix extracts language from path /en/page -> ("en", "page", true)
func extractLangPrefix(path string) (lang, suffix string, ok bool) {
	if len(path) < 2 || path[0] != '/' {
		return "", "", false
	}
	rest := path[1:]
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

func isStaticAsset(path string) bool {
	prefixes := []string{"/css/", "/images/", "/assets/", "/media/"}
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return path == "/favicon.ico"
}

func serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	filePath = filepath.Clean(filePath)

	if _, err := os.Stat(filePath); err != nil {
		// Try with .html extension
		if _, err := os.Stat(filePath + ".html"); err == nil {
			filePath += ".html"
		} else if _, err := os.Stat(filepath.Join(filePath, "index.html")); err == nil {
			filePath = filepath.Join(filePath, "index.html")
		} else {
			http.NotFound(w, r)
			return
		}
	}

	info, _ := os.Stat(filePath)
	if info != nil && info.IsDir() {
		indexPath := filepath.Join(filePath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			filePath = indexPath
		} else {
			http.NotFound(w, r)
			return
		}
	}

	http.ServeFile(w, r, filePath)
}
