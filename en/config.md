---
title: config.env Reference
layout: default
lang: en
category: Administration
order: 334
---

The panel is configured with environment variables. Values are read from the `config.env` file:

* Linux — `/etc/gameap/config.env`
* Windows — `C:\gameap\web\config.env`

The file format is one `NAME=value` pair per line, without quotes and without `export`. Lines
starting with `#` are ignored. After changing the file, the panel has to be restarted:

```bash
gameapctl panel restart
```

Variables set in the process environment take precedence over the file.

Only two variables are required: `DATABASE_URL` and `AUTH_SECRET`. The panel will not start without
them. Everything else has defaults suitable for a typical single-server installation.

Notation used in the tables: *duration* — a number with an `ms`, `s`, `m`, or `h` suffix (`100ms`,
`30s`, `5m`, `12h`); *list* — comma-separated values; *size* — a number with a binary suffix `K`,
`M`, `G`, `T`, `P` (`8M`, `100G`), a bare number is treated as bytes. The `KB` and `KiB` spellings
are accepted as well; every suffix is a power of 1024, and fractional values such as `1.5G` are
allowed.

## HTTP

| Variable               | Type   | Default   | Purpose                                                                          |
|------------------------|--------|-----------|-----------------------------------------------------------------------------------|
| `HTTP_HOST`            | string | `0.0.0.0` | Panel host name. Used for CORS and as the name in self-signed certificates        |
| `HTTP_BIND_IP`         | string | `""`      | IP address to listen on. Empty — all interfaces                                   |
| `HTTP_PORT`            | number | `8025`    | Web interface and API port                                                        |
| `HTTPS_PORT`           | number | `443`     | HTTPS port. Listened on only when a certificate is configured                     |
| `HTTP_ALLOWED_ORIGINS` | list   | `""`      | Origins allowed to make requests from the browser. Empty — a single origin derived from `HTTP_HOST` |

`HTTP_ALLOWED_ORIGINS` takes full origins including the scheme: `https://panel.example.com`.
The `*` wildcard is not supported.

## Database

| Variable                   | Type     | Default | Purpose                                                     |
|----------------------------|----------|---------|--------------------------------------------------------------|
| `DATABASE_DRIVER`          | string   | `mysql` | `mysql`, `postgres`, `sqlite`, `inmemory`                    |
| `DATABASE_URL`             | string   | —       | **Required.** Connection string                              |
| `DATABASE_CONNECT_TIMEOUT` | duration | `30s`   | How long to keep retrying the initial connection at startup  |

The PostgreSQL driver names are interchangeable: `postgres`, `postgresql`, `pgsql`, `pg`, `pgx`.

Within `DATABASE_CONNECT_TIMEOUT` a database that is not reachable yet — restarting during an
upgrade, starting later than the panel — is retried with backoff instead of the panel exiting at
once.

Connection string formats:

```text
# PostgreSQL
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://gameap:password@localhost:5432/gameap?sslmode=disable

# MySQL / MariaDB
DATABASE_DRIVER=mysql
DATABASE_URL=gameap:password@tcp(localhost:3306)/gameap?parseTime=true

# SQLite
DATABASE_DRIVER=sqlite
DATABASE_URL=file:/var/lib/gameap/db.sqlite?_busy_timeout=5000&_journal_mode=WAL&cache=shared
```

The `inmemory` driver keeps data in RAM only and is meant for tests: everything is lost on
restart.

## TLS and certificates

| Variable          | Type   | Default | Purpose                                            |
|-------------------|--------|---------|-----------------------------------------------------|
| `TLS_CERT_FILE`   | string | `""`    | Path to the certificate file                        |
| `TLS_KEY_FILE`    | string | `""`    | Path to the private key file                        |
| `TLS_CERT`        | string | `""`    | Certificate directly in the variable, PEM or base64 |
| `TLS_KEY`         | string | `""`    | Private key directly in the variable, PEM or base64 |
| `TLS_FORCE_HTTPS` | bool   | `false` | Redirect HTTP to HTTPS                              |

HTTPS is enabled when at least one certificate source is configured: files, values in variables,
or ACME. If none is configured, the panel serves HTTP only.

With `TLS_FORCE_HTTPS=true`, every HTTP request gets a `301` redirect, except
`/.well-known/acme-challenge/` — otherwise issuing a Let's Encrypt certificate would be impossible.

### ACME (Let's Encrypt)

| Variable                      | Type     | Default         | Purpose                                                        |
|-------------------------------|----------|-----------------|-----------------------------------------------------------------|
| `ACME_ENABLED`                | bool     | `false`         | Automatic certificate issuance                                  |
| `ACME_EMAIL`                  | string   | `""`            | Notification address. **Required** when ACME is enabled         |
| `ACME_DOMAINS`                | list     | `""`            | Certificate domains. **Required** when ACME is enabled          |
| `ACME_CHALLENGE_TYPE`         | string   | `http-01`       | `http-01` or `dns-01`                                           |
| `ACME_DNS_PROVIDER`           | string   | `""`            | DNS provider for `dns-01`, for example `cloudflare`             |
| `ACME_DIRECTORY_URL`          | string   | production ACME | ACME directory URL. For testing, point it at Let's Encrypt staging |
| `ACME_RENEWAL_THRESHOLD`      | duration | `720h`          | How long before expiry to renew. 30 days by default             |
| `ACME_RENEWAL_CHECK_INTERVAL` | duration | `12h`           | How often to check the expiry date                              |
| `ACME_PROPAGATION_TIMEOUT`    | duration | `180s`          | How long to wait for DNS record propagation with `dns-01`       |
| `ACME_STORAGE_PATH`           | string   | `acme`          | Directory for storing certificates and the ACME account key     |

ACME is enabled only when `ACME_ENABLED=true`, `ACME_EMAIL`, and `ACME_DOMAINS` are all set.
If anything is missing, the panel silently continues without ACME.

The default is the production Let's Encrypt directory with its strict limits on the number of
attempts. While debugging your setup, use staging:
`ACME_DIRECTORY_URL=https://acme-staging-v02.api.letsencrypt.org/directory`.

## Authentication

| Variable                        | Type     | Default  | Purpose                                                                  |
|---------------------------------|----------|----------|---------------------------------------------------------------------------|
| `AUTH_SECRET`                   | string   | —        | **Required.** Token signing key, exactly 32 random bytes                  |
| `ENCRYPTION_KEY`                | string   | `""`     | Encryption key for secrets in the database, 32 random bytes recommended   |
| `AUTH_SERVICE`                  | string   | `paseto` | Token format: `paseto` or `jwt`. Any other value — the panel will not start |
| `AUTH_BCRYPT_COST`              | number   | `13`     | bcrypt cost, 10 to 14                                                     |
| `AUTH_ALLOW_WEAK_PASSWORDS`     | bool     | `false`  | Disables the check against the compromised password list                  |
| `AUTH_REQUIRE_MFA_FOR_ADMINS`   | bool     | `true`   | Require 2FA from administrators                                           |
| `AUTH_MFA_HARD_FAIL_DAYS`       | number   | `30`     | Days until lockout. `0` — reminder only                                   |
| `AUTH_MFA_ENROLLMENT_TOKEN_TTL` | duration | `15m`    | Lifetime of the restricted session issued after the deadline              |
| `AUTH_SHORT_LIVED_TOKEN_TTL`    | duration | `10s`    | Lifetime of one-time `glst_` tokens. Effectively capped at 10 seconds     |
| `AUTH_SSO_TICKET_TTL`           | duration | `60s`    | Lifetime of a single-use SSO login ticket. Capped at 120 seconds          |

See the [Security](/en/security.html) page for details.

An SSO ticket is issued for another user and exchanged by an external system — a billing panel,
for example — for a logged-in session. The issuing handler caps its lifetime at 120 seconds; a zero
or negative value falls back to 60 seconds. In a multi-instance installation the ticket is issued on
one instance and redeemed on another, so it needs a shared `CACHE_DRIVER`: `redis`, `mysql`, or
`postgres`, not `memory`.

> `AUTH_SECRET` is silently coerced to 32 bytes: a short value is padded, a long one is truncated.
> Set exactly 32 characters, for example `openssl rand -base64 24`.

### First administrator

These three variables are read only during the initial seeding of an empty database:

| Variable         | Purpose                                                                        |
|------------------|---------------------------------------------------------------------------------|
| `ADMIN_LOGIN`    | Login of the first administrator                                                |
| `ADMIN_EMAIL`    | Email address                                                                   |
| `ADMIN_PASSWORD` | Password. If not set, a random one is generated and printed to the log on first start |

> The password from `ADMIN_PASSWORD` is **not** checked against the password policy: neither the
> length nor the compromised password list. Choose it deliberately.

## Access control and cache

| Variable                    | Type     | Default          | Purpose                                       |
|-----------------------------|----------|------------------|------------------------------------------------|
| `RBAC_CACHE_TTL`            | duration | `30s`            | Permission check cache lifetime                |
| `CACHE_DRIVER`              | string   | `memory`         | `memory`, `redis`, `postgres`, `mysql`. See below |
| `CACHE_REDIS_ADDR`          | string   | `localhost:6379` | Redis address                                  |
| `CACHE_REDIS_PASSWORD`      | string   | `""`             | Redis password                                 |
| `CACHE_REDIS_DB`            | number   | `0`              | Redis database number                          |
| `CACHE_TTL_RBAC`            | duration | `24h`            | Permissions cache lifetime                     |
| `CACHE_TTL_GAMES`           | duration | `48h`            | Games cache lifetime                           |
| `CACHE_TTL_NODES`           | duration | `24h`            | Dedicated servers cache lifetime               |
| `CACHE_TTL_USERS`           | duration | `6h`             | Users cache lifetime                           |
| `CACHE_TTL_PERSONAL_TOKENS` | duration | `24h`            | Personal tokens cache lifetime                 |
| `CACHE_TTL_SERVER_SETTINGS` | duration | `12h`            | Server settings cache lifetime                 |

`CACHE_DRIVER` accepts `memory` (alias `inmemory`), `redis`, `mysql` (alias `database`), and
`postgres` (aliases `postgresql`, `pgsql`, `pg`). An unknown value causes the panel to fail at
startup.

The cache holds more than reference data: the daemon setup key, the revoked token list, and the
login attempt counters live there. With `CACHE_DRIVER=memory` all of that is lost on restart and
is not shared between multiple panel instances. For a multi-instance installation, use `redis`.

## Files

| Variable                     | Type   | Default | Purpose                                |
|------------------------------|--------|---------|-----------------------------------------|
| `FILES_DRIVER`               | string | `local` | `local` or `s3`                         |
| `FILES_LOCAL_BASE_PATH`      | string | `""`    | Base directory for the `local` driver   |
| `FILES_S3_ENDPOINT`          | string | `""`    | Address of the S3-compatible storage    |
| `FILES_S3_USE_SSL`           | bool   | `true`  | Access the storage over HTTPS           |
| `FILES_S3_ACCESS_KEY_ID`     | string | `""`    | Access key ID                           |
| `FILES_S3_SECRET_ACCESS_KEY` | string | `""`    | Secret access key                       |
| `FILES_S3_BUCKET`            | string | `""`    | Bucket name                             |

### File uploads

| Variable                        | Type     | Default  | Purpose                                             |
|---------------------------------|----------|----------|------------------------------------------------------|
| `FILES_UPLOAD_CHUNK_SIZE`       | size     | `8M`     | Chunk size for chunked uploads                       |
| `FILES_UPLOAD_MAX_CHUNKS`       | number   | `100000` | Maximum number of chunks per file                    |
| `FILES_UPLOAD_SESSION_TTL`      | duration | `24h`    | How long an unfinished upload lives                  |
| `FILES_UPLOAD_DISPATCH_TIMEOUT` | duration | `2m`     | Timeout for handing the assembled file to the daemon |
| `FILES_UPLOAD_JANITOR_INTERVAL` | duration | `12h`    | How often to clean up expired uploads                |
| `FILES_UPLOAD_ALLOWED_MIMES`    | list     | `""`     | Extends the list of allowed types, does not replace it |
| `FILES_UPLOAD_ALLOW_ARCHIVES`   | bool     | `false`  | Allow archives: zip, tar, gzip, bzip2, 7z, xz        |
| `FILES_UPLOAD_ALLOW_BINARY`     | bool     | `false`  | Allow arbitrary binary files                         |

The maximum file size is `FILES_UPLOAD_CHUNK_SIZE` multiplied by `FILES_UPLOAD_MAX_CHUNKS`. With
the defaults that is about 780 GB. A separate 100 MB limit applies to a regular single-request
upload; it is not configurable.

### Archives

| Variable                              | Type   | Default  | Purpose                                                                                  |
|---------------------------------------|--------|----------|-------------------------------------------------------------------------------------------|
| `FILES_ARCHIVE_MAX_BYTES`             | size   | `100G`   | Maximum uncompressed size of one archive operation                                        |
| `FILES_ARCHIVE_MAX_FILES`             | number | `500000` | Maximum number of entries in one archive operation                                        |
| `FILES_ARCHIVE_CONCURRENT_PER_SERVER` | number | `2`      | Concurrent ZIP downloads of a directory per game server. Exceeding it returns HTTP `429` |

`FILES_ARCHIVE_MAX_BYTES` and `FILES_ARCHIVE_MAX_FILES` bound the streamed ZIP download of a
directory and, since 4.5.0, also the archive creation and extraction that the
[file manager](/en/gameap_configure/file_manager.html) runs on the dedicated server; on extraction
the size limit guards against decompression bombs. A value of `0` removes the limit from the
streamed download; for an operation on the dedicated server it leaves the bound to the daemon,
which then applies its own defaults: 10 GiB and 100 000 entries. `FILES_ARCHIVE_CONCURRENT_PER_SERVER`
applies to downloads only and does not limit creation or extraction on the dedicated server.

## gRPC

| Variable                      | Type   | Default    | Purpose                                            |
|-------------------------------|--------|------------|-----------------------------------------------------|
| `GRPC_PORT`                   | number | `31718`    | gRPC server port for daemons                        |
| `GRPC_TLS_ENABLED`            | bool   | `true`     | Encryption of daemon connections                    |
| `GRPC_REQUIRE_MTLS`           | bool   | `false`    | Require a client certificate from the daemon        |
| `GRPC_EXTERNAL_HOST`          | string | `""`       | Panel address reported to the daemon                |
| `GRPC_EXTERNAL_PORT`          | number | `0`        | Port reported to the daemon. `0` — `GRPC_PORT` is used |
| `GRPC_MAX_RECV_MSG_SIZE`      | number | `10485760` | Maximum incoming message size, bytes                |
| `GRPC_MAX_SEND_MSG_SIZE`      | number | `10485760` | Maximum outgoing message size, bytes                |
| `GRPC_MAX_CONCURRENT_STREAMS` | number | `100`      | Concurrent streams per connection                   |
| `GRPC_ENABLE_REFLECTION`      | bool   | `false`    | Schema reflection for debugging tools               |
| `DAEMON_SETUP_KEY`            | string | `""`       | Permanent daemon setup key instead of a temporary one |

See the [GRPC API](/en/daemon/grpc.html) page for details.

> `DAEMON_SETUP_KEY` sets a key that never expires. That is convenient for automated deployments,
> but such a key lets anyone who learns it register a new dedicated server in the panel. For a
> regular installation, leave this variable unset — the panel will issue a temporary key valid
> for one hour.

## Security

The full description is on the [Security](/en/security.html) page.

| Variable                           | Type   | Default                           | Purpose                                   |
|------------------------------------|--------|-----------------------------------|--------------------------------------------|
| `SECURITY_HEADERS_ENABLED`         | bool   | `true`                            | Master switch for security headers         |
| `SECURITY_CONTENT_TYPE_OPTIONS`    | bool   | `true`                            | `X-Content-Type-Options: nosniff`          |
| `SECURITY_FRAME_OPTIONS`           | string | `SAMEORIGIN`                      | `X-Frame-Options`                          |
| `SECURITY_REFERRER_POLICY`         | string | `strict-origin-when-cross-origin` | `Referrer-Policy`                          |
| `SECURITY_HSTS_ENABLED`            | bool   | `true`                            | HSTS. Sent over HTTPS only                 |
| `SECURITY_HSTS_MAX_AGE`            | number | `31536000`                        | HSTS lifetime in seconds                   |
| `SECURITY_HSTS_INCLUDE_SUBDOMAINS` | bool   | `false`                           | Extend HSTS to subdomains                  |
| `SECURITY_HSTS_PRELOAD`            | bool   | `false`                           | Add `preload`                              |
| `SECURITY_CSP_ENABLED`             | bool   | `true`                            | Content Security Policy                    |
| `SECURITY_CSP_REPORT_ONLY`         | bool   | `false`                           | Reports only, no blocking                  |
| `SECURITY_CSP_POLICY`              | string | `""`                              | Completely replaces the generated policy   |
| `SECURITY_CSP_REPORT_URI`          | string | `""`                              | Address for CSP reports                    |
| `SECURITY_CSP_EXTRA_SCRIPT_SRC`    | list   | `""`                              | Extends `script-src`                       |
| `SECURITY_CSP_EXTRA_STYLE_SRC`     | list   | `""`                              | Extends `style-src`                        |
| `SECURITY_CSP_EXTRA_CONNECT_SRC`   | list   | `""`                              | Extends `connect-src`                      |
| `SECURITY_CSP_EXTRA_IMG_SRC`       | list   | `""`                              | Extends `img-src`                          |
| `SECURITY_CSP_EXTRA_FRAME_SRC`     | list   | `""`                              | Extends `frame-src`                        |
| `SECURITY_CSP_EXTRA_FONT_SRC`      | list   | `""`                              | Extends `font-src`                         |
| `SECURITY_SENSITIVE_PATH_PREFIXES` | list   | see below                         | Paths whose responses must never be cached |

The default value of `SECURITY_SENSITIVE_PATH_PREFIXES`:
`/api/auth/,/api/profile/,/api/users/,/api/tokens/`.

### CAPTCHA

| Variable             | Type   | Default | Purpose                                                        |
|----------------------|--------|---------|-----------------------------------------------------------------|
| `CAPTCHA_PROVIDER`   | string | `""`    | `recaptcha_v2`, `recaptcha_v3`, or `turnstile`. Empty — disabled |
| `CAPTCHA_SITE_KEY`   | string | `""`    | Public key                                                      |
| `CAPTCHA_SECRET_KEY` | string | `""`    | Secret key                                                      |
| `CAPTCHA_MIN_SCORE`  | float  | `0.5`   | Threshold, reCAPTCHA v3 only                                    |
| `CAPTCHA_FAIL_OPEN`  | bool   | `false` | Let logins through when the verification service is unavailable |
| `CAPTCHA_VERIFY_URL` | string | `""`    | Custom verification address                                     |

### Audit log

| Variable                 | Type   | Default | Purpose                              |
|--------------------------|--------|---------|---------------------------------------|
| `AUDIT_ENABLED`          | bool   | `true`  | Recording of security events          |
| `AUDIT_CLIENT_IP_HEADER` | string | `""`    | Header carrying the real client IP    |

## Plugins

> Since 4.5.0 every plugin setting is spelled `PLUGINS_*`. The former `PLUGIN_*` names, as well as
> `PLUGINS_CACHE_ENABLED` and `PLUGINS_CACHE_DIR`, keep working for one release: the panel applies
> the value and logs `environment variable is deprecated and will be removed in a future release`
> with the name of the replacement. When both the old and the new name are set, the old one is
> ignored. The exception is three names that carried the unit: `PLUGIN_HTTP_MAX_TIMEOUT_SECONDS`,
> `PLUGIN_NET_MAX_TIMEOUT_SECONDS`, and `PLUGIN_NET_READ_BUFFER_BYTES` are no longer recognised and
> are silently ignored — replace them with `PLUGINS_HTTP_MAX_TIMEOUT=30s`, `PLUGINS_NET_MAX_TIMEOUT=10s`,
> and `PLUGINS_NET_READ_BUFFER=64K`. `gameapctl panel upgrade` rewrites all of these names in
> `config.env` automatically.

| Variable                    | Type   | Default                          | Purpose                                                                                             |
|-----------------------------|--------|----------------------------------|------------------------------------------------------------------------------------------------------|
| `PLUGINS_DISABLED`          | bool   | `false`                          | Disable plugins entirely                                                                             |
| `PLUGINS_AUTOLOAD`          | list   | `""`                             | Plugins loaded at startup                                                                            |
| `PLUGINS_STRICT_LOAD`       | bool   | `false`                          | Refuse to start when any plugin fails to load. Off: the plugin gets status `error` and is skipped   |
| `PLUGINS_STORE_URL`         | string | `https://plugins.gameap.dev/api` | Plugin catalog address                                                                               |
| `PLUGINS_STORE_LICENSE_KEY` | string | `""`                             | License key for paid plugins                                                                         |

### Plugin runtime

| Variable                          | Type   | Default | Purpose                                                              |
|-----------------------------------|--------|---------|-----------------------------------------------------------------------|
| `PLUGINS_RUNTIME_MAX_MEMORY`      | size   | `256M`  | Memory cap per plugin module. `0` — the wazero default of 4 GiB      |
| `PLUGINS_RUNTIME_MAX_MODULE_SIZE` | size   | `128M`  | Largest `.wasm` file accepted for install and load. `0` — unlimited  |
| `PLUGINS_RUNTIME_CACHE_ENABLED`   | bool   | `true`  | Cache of compiled WebAssembly modules                                 |
| `PLUGINS_RUNTIME_CACHE_DIR`       | string | `""`    | Directory of that cache. When set, compiled code survives panel restarts |

A module that declares a larger maximum memory is clamped to `PLUGINS_RUNTIME_MAX_MEMORY`; only a
module whose initial memory already exceeds the cap fails to load. `PLUGINS_RUNTIME_MAX_MODULE_SIZE`
applies to uploads, catalog installs, and autoload alike; a plugin uploaded through the web interface
is additionally limited to 100 MB regardless of this value.

### Plugin network requests

| Variable                                 | Type     | Default | Purpose                                          |
|------------------------------------------|----------|---------|---------------------------------------------------|
| `PLUGINS_HTTP_BLOCK_PRIVATE_IPS`         | bool     | `true`  | Block requests to private network addresses       |
| `PLUGINS_HTTP_ALLOWED_SCHEMES`           | list     | `https` | Allowed schemes                                   |
| `PLUGINS_HTTP_ALLOWED_HOSTS`             | list     | `""`    | Hosts exempt from the private address block       |
| `PLUGINS_HTTP_MAX_TIMEOUT`               | duration | `30s`   | Maximum request time                              |
| `PLUGINS_HTTP_MAX_REDIRECTS`             | number   | `5`     | Maximum number of redirects                       |
| `PLUGINS_HTTP_RESPONSE_HEADER_ALLOWLIST` | list     | `""`    | Extends the list of headers passed to the plugin  |

Cloud provider metadata service addresses are always blocked; `PLUGINS_HTTP_ALLOWED_HOSTS` has no
effect on them.

A separate group of settings governs custom RCON and Query protocols implemented by plugins:

| Variable                        | Type     | Default | Purpose                                                 |
|---------------------------------|----------|---------|----------------------------------------------------------|
| `PLUGINS_NET_ENABLED`           | bool     | `true`  | Allow plugins to talk to game servers over the network   |
| `PLUGINS_NET_BLOCK_PRIVATE_IPS` | bool     | `false` | Block private addresses. Off: game servers often live on the internal network |
| `PLUGINS_NET_ALLOWED_HOSTS`     | list     | `""`    | Exempt hosts                                             |
| `PLUGINS_NET_MAX_TIMEOUT`       | duration | `10s`   | Maximum time for a single operation                      |
| `PLUGINS_NET_READ_BUFFER`       | size     | `64K`   | Maximum size of a single read                            |
| `PLUGINS_NET_MAX_CONNECTIONS`   | number   | `8`     | Concurrent connections per plugin                        |

### Plugin scheduler

| Variable                                 | Type     | Default | Purpose                                       |
|------------------------------------------|----------|---------|------------------------------------------------|
| `PLUGINS_SCHEDULER_MIN_INTERVAL`         | duration | `1s`    | Smallest allowed task interval                 |
| `PLUGINS_SCHEDULER_MAX_TASKS_PER_PLUGIN` | number   | `32`    | Maximum number of tasks per plugin             |
| `PLUGINS_SCHEDULER_CALL_TIMEOUT`         | duration | `60s`   | Default handler call timeout                   |
| `PLUGINS_SCHEDULER_MAX_CALL_TIMEOUT`     | duration | `5m`    | Maximum timeout a plugin may request           |
| `PLUGINS_SCHEDULER_MAX_RETRIES`          | number   | `10`    | Maximum number of retries                      |
| `PLUGINS_SCHEDULER_MAX_RETRY_DELAY`      | duration | `10m`   | Maximum delay between retries                  |
| `PLUGINS_SCHEDULER_MAX_JITTER`           | duration | `30s`   | Maximum random start jitter                    |
| `PLUGINS_SCHEDULER_REFRESH_INTERVAL`     | duration | `30s`   | How often to re-read tasks from the database   |

### Plugin permissions

| Variable                        | Type     | Default | Purpose                                                                                     |
|---------------------------------|----------|---------|----------------------------------------------------------------------------------------------|
| `PLUGINS_PERMISSIONS_ENFORCE`   | bool     | `false` | Enforce the recorded permission grants                                                       |
| `PLUGINS_PERMISSIONS_CACHE_TTL` | duration | `30s`   | How long a plugin's grants stay cached in memory. `0` — read the plugin record on every check |

Enforcement is off in this release so that plugins written before grants existed keep working. With
`false` no grant check blocks anything; the grants are still recorded, shown, and editable in the
panel. Grant changes are pushed to every instance over pub/sub, so the cache TTL only bounds the
drift while the broker is unreachable.

### Plugin recovery

| Variable                         | Type     | Default | Purpose                                              |
|----------------------------------|----------|---------|-------------------------------------------------------|
| `PLUGINS_RECOVERY_ENABLED`       | bool     | `true`  | Reload plugins the runtime has disabled automatically |
| `PLUGINS_RECOVERY_INITIAL_DELAY` | duration | `30s`   | Delay before the first attempt                        |
| `PLUGINS_RECOVERY_MAX_DELAY`     | duration | `10m`   | Upper bound of the exponentially growing delay        |
| `PLUGINS_RECOVERY_MAX_ATTEMPTS`  | number   | `5`     | Number of attempts                                    |

A plugin the runtime disabled — a call overran its deadline or the module terminated itself — is
reloaded with exponential backoff. After the last attempt it stays in status `error` until an
operator reloads it or the panel restarts. With `PLUGINS_RECOVERY_ENABLED=false` the disable is still
recorded (status, reason, audit log), but no automatic reload happens.

### Plugin synchronization

| Variable                        | Type     | Default | Purpose                                                                            |
|---------------------------------|----------|---------|-------------------------------------------------------------------------------------|
| `PLUGINS_SYNC_DISABLED`         | bool     | `false` | Turn off syncing plugin changes between instances                                   |
| `PLUGINS_SYNC_REFRESH_INTERVAL` | duration | `60s`   | Period of the reconcile pass; bounds how long a lost hint leaves an instance stale  |
| `PLUGINS_SYNC_MIN_BACKOFF`      | duration | `15s`   | Minimum retry delay after a plugin failed to load on this instance                  |
| `PLUGINS_SYNC_MAX_BACKOFF`      | duration | `15m`   | Maximum retry delay                                                                 |

Relevant for [multi-instance](/en/multi_instance.html) installations: install, update, uninstall,
reload, and permission changes made on one instance reach the others through a pub/sub hint and a
periodic reconcile pass. With sync disabled, changes reach an instance only when it restarts.

### Plugin node file access

| Variable                       | Type   | Default        | Purpose                                                                                       |
|--------------------------------|--------|----------------|------------------------------------------------------------------------------------------------|
| `PLUGINS_NODEFS_MAX_INLINE`    | size   | `32M`          | Largest file a plugin may download or upload in one call. `0` — unlimited                     |
| `PLUGINS_NODEFS_PATH_POLICY`   | string | `unrestricted` | Where plugins may point node file operations: `unrestricted`, `node_workpath`, `server_dirs`  |
| `PLUGINS_NODEFS_ALLOWED_PATHS` | list   | `""`           | Extra absolute roots allowed on every node in the restricted modes                            |

`node_workpath` confines paths to the work path of the dedicated server, `server_dirs` — to the
directories of the game servers on that node. The policy also applies to the working directory of
commands plugins run on nodes. `..` segments are refused in every mode.

### Plugin storage and cache

| Variable                              | Type   | Default | Purpose                                             |
|---------------------------------------|--------|---------|------------------------------------------------------|
| `PLUGINS_STORAGE_MAX_KEYS_PER_PLUGIN` | number | `10000` | Number of entries one plugin may keep                |
| `PLUGINS_STORAGE_MAX_VALUE`           | size   | `1M`    | Largest single stored value                          |
| `PLUGINS_STORAGE_MAX_TOTAL`           | size   | `64M`   | Total size of all values of one plugin               |
| `PLUGINS_CACHE_MAX_VALUE`             | size   | `1M`    | Largest value in the plugin cache. `0` — unlimited   |

The storage quotas cannot be switched off: a zero or negative value falls back to the built-in
default. The plugin cache lives in the panel cache (`CACHE_DRIVER`), every plugin has its own key
namespace, and it is not cleaned up when a plugin is uninstalled.

### Plugin secrets

| Variable                              | Type   | Default | Purpose                                        |
|---------------------------------------|--------|---------|-------------------------------------------------|
| `PLUGINS_SECRETS_MAX_KEYS_PER_PLUGIN` | number | `64`    | Number of secrets one plugin may keep           |
| `PLUGINS_SECRETS_MAX_VALUE`           | size   | `8K`    | Largest plaintext size of a single secret       |
| `PLUGINS_SECRETS_REQUIRE_ENCRYPTION`  | bool   | `true`  | Refuse writes while `ENCRYPTION_KEY` is unset   |

Plugin secrets are encrypted at rest with `ENCRYPTION_KEY`. With
`PLUGINS_SECRETS_REQUIRE_ENCRYPTION=false` and no key set, secrets are stored in plaintext — no safer
than ordinary plugin storage.

### Plugin SSH

| Variable                                | Type     | Default   | Purpose                                                                                  |
|-----------------------------------------|----------|-----------|-------------------------------------------------------------------------------------------|
| `PLUGINS_SSH_ENABLED`                   | bool     | `false`   | Allow plugins to open SSH connections from the panel                                      |
| `PLUGINS_SSH_BLOCK_PRIVATE_IPS`         | bool     | `true`    | Block connections to private network addresses                                            |
| `PLUGINS_SSH_ALLOWED_HOSTS`             | list     | `""`      | Hosts exempt from the private address block                                               |
| `PLUGINS_SSH_ALLOW_ACCEPT_ANY_HOST_KEY` | bool     | `true`    | Allow the `accept_any` host key policy (trust on first use)                               |
| `PLUGINS_SSH_MAX_CONNECTIONS`           | number   | `8`       | Concurrent connections per plugin                                                         |
| `PLUGINS_SSH_MAX_OPERATIONS`            | number   | `16`      | Concurrently running commands per plugin                                                  |
| `PLUGINS_SSH_CONNECT_TIMEOUT`           | duration | `30s`     | Dial, handshake, and authentication together; also the ceiling for a plugin's own value   |
| `PLUGINS_SSH_MAX_EXEC_TIMEOUT`          | duration | `30m`     | Ceiling for a single remote command                                                       |
| `PLUGINS_SSH_IDLE_TIMEOUT`              | duration | `10m`     | Close a connection nothing has run on for this long                                       |
| `PLUGINS_SSH_MAX_OUTPUT_BYTES`          | number   | `1048576` | Captured stdout and stderr per command, bytes. The head is kept, the rest is truncated    |
| `PLUGINS_SSH_MAX_STDIN_BYTES`           | number   | `1048576` | Data a plugin may pipe into a command, bytes                                              |
| `PLUGINS_SSH_OPERATION_RETENTION`       | duration | `10m`     | How long a finished operation with its output stays readable                              |
| `PLUGINS_SSH_MAX_RETAINED_OPERATIONS`   | number   | `64`      | Finished operations kept per plugin; the oldest are evicted first                         |
| `PLUGINS_SSH_KEEPALIVE_INTERVAL`        | duration | `30s`     | Interval of liveness probes on open connections                                           |
| `PLUGINS_SSH_COMPLETION_CALL_TIMEOUT`   | duration | `30s`     | Timeout of one completion callback into the plugin. A plugin that does not answer in time is disabled |
| `PLUGINS_SSH_BUSY_RETRY_DELAY`          | duration | `2s`      | Pause between completion callback retries while the plugin is busy with another call      |
| `PLUGINS_SSH_BUSY_RETRIES`              | number   | `5`       | Number of retries of a busy completion callback                                           |

The `gameap-ssh` host library lets plugins connect to hosts over SSH and run commands there — for
example, to install GameAP Daemon on a fresh machine. It is off by default, and `PLUGINS_SSH_ENABLED`
is the operator's deliberate consent to outbound SSH from the panel. Cloud provider metadata
addresses are blocked even with `PLUGINS_SSH_BLOCK_PRIVATE_IPS=false`. Connections and operations
live in the memory of the panel instance that opened them.

### Plugin rate limits

Token buckets per plugin and per panel instance on the expensive host libraries: sustained calls per
second (`_RPS`) plus a burst (`_BURST`). A refused call gets a "rate limited" error in the response;
the plugin is not disabled for it. `_RPS` of `0` removes the limit for that class.

| Variable                                | Type   | Default | Purpose                                                                                |
|-----------------------------------------|--------|---------|-----------------------------------------------------------------------------------------|
| `PLUGINS_RATELIMIT_NODECMD_RPS`         | float  | `5`     | Commands on dedicated servers (`gameap-nodecmd`)                                        |
| `PLUGINS_RATELIMIT_NODECMD_BURST`       | number | `20`    | Burst for the same class                                                                |
| `PLUGINS_RATELIMIT_SERVERCONTROL_RPS`   | float  | `5`     | Game server control, daemon tasks, server and setting writes (`gameap-servercontrol`)  |
| `PLUGINS_RATELIMIT_SERVERCONTROL_BURST` | number | `20`    | Burst for the same class                                                                |
| `PLUGINS_RATELIMIT_NODEFS_RPS`          | float  | `50`    | Node file operations (`gameap-nodefs`)                                                 |
| `PLUGINS_RATELIMIT_NODEFS_BURST`        | number | `200`   | Burst for the same class                                                                |
| `PLUGINS_RATELIMIT_HTTP_RPS`            | float  | `20`    | Outbound HTTP requests (`gameap-http`)                                                  |
| `PLUGINS_RATELIMIT_HTTP_BURST`          | number | `50`    | Burst for the same class                                                                |
| `PLUGINS_RATELIMIT_RBAC_RPS`            | float  | `10`    | Role and ability writes (`gameap-rbac`)                                                 |
| `PLUGINS_RATELIMIT_RBAC_BURST`          | number | `50`    | Burst for the same class                                                                |
| `PLUGINS_RATELIMIT_SSH_RPS`             | float  | `20`    | Every SSH call, including polling a running command (`gameap-ssh`)                      |
| `PLUGINS_RATELIMIT_SSH_BURST`           | number | `60`    | Burst for the same class                                                                |

## Event exchange between instances

Needed only when running several panel instances. With a single instance the defaults are
sufficient. An unknown `PUBSUB_DRIVER` value causes the panel to fail at startup.

| Variable                     | Type     | Default  | Purpose                                            |
|------------------------------|----------|----------|-----------------------------------------------------|
| `PUBSUB_DRIVER`              | string   | `memory` | `memory`, `redis`, `postgres`                       |
| `PUBSUB_INSTANCE_ID`         | string   | `""`     | Instance identifier, must be unique                 |
| `PUBSUB_REDIS_ADDR`          | string   | `""`     | Redis address                                       |
| `PUBSUB_REDIS_PASSWORD`      | string   | `""`     | Redis password                                      |
| `PUBSUB_REDIS_DB`            | number   | `1`      | Redis database number                               |
| `PUBSUB_RETRY_ENABLED`       | bool     | `true`   | Retry delivery on error                             |
| `PUBSUB_RETRY_MAX_RETRIES`   | number   | `3`      | Maximum number of retries                           |
| `PUBSUB_RETRY_INITIAL_DELAY` | duration | `100ms`  | Initial delay before a retry                        |
| `PUBSUB_RETRY_MAX_DELAY`     | duration | `5s`     | Maximum delay before a retry                        |
| `PUBSUB_RETRY_MULTIPLIER`    | float    | `2.0`    | Factor by which the delay grows                     |
| `PUBSUB_DLQ_ENABLED`         | bool     | `false`  | Put undelivered events into a separate queue        |
| `PUBSUB_DLQ_DRIVER`          | string   | `memory` | Storage for that queue                              |
| `PUBSUB_DLQ_MAX_SIZE`        | number   | `1000`   | Maximum queue size                                  |

## Miscellaneous

| Variable                      | Type     | Default                  | Purpose                                        |
|-------------------------------|----------|--------------------------|-------------------------------------------------|
| `LOGGER_LEVEL`                | string   | `info`                   | `debug`, `info`, `warn`, `error`                |
| `LOGGER_LOG_DB_QUERIES`       | bool     | `false`                  | Log database queries. Debugging only            |
| `DEFAULT_LANGUAGE`            | string   | `""`                     | Default interface language: `en`, `ru`, `es`, `de`. Empty — the browser language |
| `GLOBAL_API_URL`              | string   | `https://api.gameap.com` | Global API address — bug reports                |
| `GAMES_CDN_URLS`              | list     | see below                | Game catalog sources, tried in order            |
| `TASK_REAPER_INTERVAL`        | duration | `1m`                     | How often to look for stuck tasks               |
| `TASK_REAPER_STALE_THRESHOLD` | duration | `10m`                    | Idle time after which a task is considered stuck |

The default value of `GAMES_CDN_URLS`:
`https://cdn.gameap.ru/games.json,https://cdn.gameap.com/games.json`.

`DEFAULT_LANGUAGE` is lowercased when read. Spanish (`es`) and German (`de`) were added in 4.5.0.
The value applies only until the user picks a language in the interface: a saved choice comes
first, then `DEFAULT_LANGUAGE`, then the browser language.

## Update check

Since 4.5.0 the panel looks up the latest GameAP and GameAP Daemon releases and shows them to
administrators in the versions block on the home page.

| Variable               | Type     | Default   | Purpose                                            |
|------------------------|----------|-----------|-----------------------------------------------------|
| `UPDATE_CHECK_ENABLED` | bool     | `true`    | Check for new GameAP and GameAP Daemon releases    |
| `UPDATE_CHECK_URLS`    | list     | see below | Release sources, tried in order until one answers  |
| `UPDATE_CHECK_TTL`     | duration | `6h`      | How long a successful lookup is cached             |

The default value of `UPDATE_CHECK_URLS`:
`https://cdn.gameap.com/{component}/releases.json,https://cdn.gameap.ru/{component}/releases.json,https://api.github.com/repos/gameap/{repo}/releases`.
`{component}` is replaced with `gameap` or `gameap-daemon`, `{repo}` — with the GitHub repository
name, `gameap` or `daemon`. A request times out after 15 seconds; a failed lookup is cached for
15 minutes. With `UPDATE_CHECK_ENABLED=false` the panel makes no outbound requests and shows the
installed versions only. See the [Upgrade](/en/upgrade.html) page.

## Variables that have no effect

These variables are parsed by the panel but do not affect its behavior. Do not rely on them:

| Variable                        | What is wrong                                                                          |
|---------------------------------|-----------------------------------------------------------------------------------------|
| `AUTH_SESSION_IDLE_TIMEOUT`     | There is no idle session termination. A session lives its full lifetime — 24 hours or 7 days |
| `AUTH_SESSION_IDLE_UPDATE_FREQ` | Same as above                                                                           |
| `GRPC_FILE_TRANSFER_BASE_PATH`  | The value reaches the gRPC server configuration, but no handler reads it                |
| `GRPC_ENABLED`                  | No such variable exists at all. gRPC always runs and cannot be turned off. Old versions of `gameapctl` append this line to `config.env` |

Login rate limits are not configurable via environment variables either — they are set as
constants in the code. The values are listed on the [Security](/en/security.html) page.

## Minimal configuration example

```dotenv
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://gameap:password@localhost:5432/gameap?sslmode=disable

AUTH_SECRET=replace_with_32_random_bytes
ENCRYPTION_KEY=replace_with_32_random_bytes

HTTP_HOST=panel.example.com
HTTP_PORT=8025
```

Everything else falls back to the defaults. `ENCRYPTION_KEY` is not formally required, but without
it plugins cannot store secrets: with the default `PLUGINS_SECRETS_REQUIRE_ENCRYPTION=true` such
writes are refused. Its length is not checked: the value is hashed in full, so 32 random bytes is a
recommendation rather than a requirement.
