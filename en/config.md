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

Notation used in the tables: *duration* — a number with an `s`, `m`, or `h` suffix (`30s`, `5m`,
`12h`); *list* — comma-separated values; *size* — a number with a binary suffix `K`, `M`, `G`, `T`
(`8M`, `100G`), a bare number is treated as bytes.

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

| Variable          | Type   | Default | Purpose                                   |
|-------------------|--------|---------|--------------------------------------------|
| `DATABASE_DRIVER` | string | `mysql` | `mysql`, `postgres`, `sqlite`, `inmemory`  |
| `DATABASE_URL`    | string | —       | **Required.** Connection string            |

The PostgreSQL driver names are interchangeable: `postgres`, `postgresql`, `pgsql`, `pg`, `pgx`.

Connection string formats:

```
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
| `ENCRYPTION_KEY`                | string   | `""`     | Encryption key for secrets in the database, exactly 32 random bytes       |
| `AUTH_SERVICE`                  | string   | `paseto` | Token format: `paseto` or `jwt`. Any other value — the panel will not start |
| `AUTH_BCRYPT_COST`              | number   | `13`     | bcrypt cost, 10 to 14                                                     |
| `AUTH_ALLOW_WEAK_PASSWORDS`     | bool     | `false`  | Disables the check against the compromised password list                  |
| `AUTH_REQUIRE_MFA_FOR_ADMINS`   | bool     | `true`   | Require 2FA from administrators                                           |
| `AUTH_MFA_HARD_FAIL_DAYS`       | number   | `30`     | Days until lockout. `0` — reminder only                                   |
| `AUTH_MFA_ENROLLMENT_TOKEN_TTL` | duration | `15m`    | Lifetime of the restricted session issued after the deadline              |
| `AUTH_SHORT_LIVED_TOKEN_TTL`    | duration | `10s`    | Lifetime of one-time `glst_` tokens. Effectively capped at 10 seconds     |

See the [Security](/en/security.html) page for details.

> `AUTH_SECRET` is silently coerced to 32 bytes: a short value is padded, a long one is truncated.
> Set exactly 32 random bytes, for example `openssl rand -hex 16`.

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

| Variable                               | Type   | Default  | Purpose                                        |
|----------------------------------------|--------|----------|-------------------------------------------------|
| `FILES_ARCHIVE_MAX_BYTES`              | size   | `100G`   | Maximum size of a created archive               |
| `FILES_ARCHIVE_MAX_FILES`              | number | `500000` | Maximum number of files in an archive           |
| `FILES_ARCHIVE_CONCURRENT_PER_SERVER`  | number | `2`      | Concurrent archive operations per server        |

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

| Variable                   | Type   | Default                          | Purpose                                  |
|----------------------------|--------|----------------------------------|-------------------------------------------|
| `PLUGINS_DISABLED`         | bool   | `false`                          | Disable plugins entirely                  |
| `PLUGINS_AUTOLOAD`         | list   | `""`                             | Plugins loaded at startup                 |
| `PLUGINS_CACHE_ENABLED`    | bool   | `true`                           | Cache of compiled WebAssembly modules     |
| `PLUGINS_CACHE_DIR`        | string | `""`                             | Directory of that cache                   |
| `PLUGIN_STORE_URL`         | string | `https://plugins.gameap.dev/api` | Plugin catalog address                    |
| `PLUGIN_STORE_LICENSE_KEY` | string | `""`                             | License key for paid plugins              |

### Plugin network requests

| Variable                                | Type   | Default | Purpose                                          |
|-----------------------------------------|--------|---------|---------------------------------------------------|
| `PLUGIN_HTTP_BLOCK_PRIVATE_IPS`         | bool   | `true`  | Block requests to private network addresses       |
| `PLUGIN_HTTP_ALLOWED_SCHEMES`           | list   | `https` | Allowed schemes                                   |
| `PLUGIN_HTTP_ALLOWED_HOSTS`             | list   | `""`    | Hosts exempt from the private address block       |
| `PLUGIN_HTTP_MAX_TIMEOUT_SECONDS`       | number | `30`    | Maximum request time                              |
| `PLUGIN_HTTP_MAX_REDIRECTS`             | number | `5`     | Maximum number of redirects                       |
| `PLUGIN_HTTP_RESPONSE_HEADER_ALLOWLIST` | list   | `""`    | Extends the list of headers passed to the plugin  |

Cloud provider metadata service addresses are always blocked; `PLUGIN_HTTP_ALLOWED_HOSTS` has no
effect on them.

A separate group of settings governs custom RCON and Query protocols implemented by plugins:

| Variable                         | Type   | Default | Purpose                                                 |
|----------------------------------|--------|---------|----------------------------------------------------------|
| `PLUGIN_NET_ENABLED`             | bool   | `true`  | Allow plugins to talk to game servers over the network   |
| `PLUGIN_NET_BLOCK_PRIVATE_IPS`   | bool   | `false` | Block private addresses. Off: game servers often live on the internal network |
| `PLUGIN_NET_ALLOWED_HOSTS`       | list   | `""`    | Exempt hosts                                             |
| `PLUGIN_NET_MAX_TIMEOUT_SECONDS` | number | `10`    | Maximum time for a single operation                      |
| `PLUGIN_NET_READ_BUFFER_BYTES`   | number | `65536` | Maximum size of a single read                            |
| `PLUGIN_NET_MAX_CONNECTIONS`     | number | `8`     | Concurrent connections per plugin                        |

### Plugin scheduler

| Variable                                | Type     | Default | Purpose                                       |
|-----------------------------------------|----------|---------|------------------------------------------------|
| `PLUGIN_SCHEDULER_MIN_INTERVAL`         | duration | `1s`    | Smallest allowed task interval                 |
| `PLUGIN_SCHEDULER_MAX_TASKS_PER_PLUGIN` | number   | `32`    | Maximum number of tasks per plugin             |
| `PLUGIN_SCHEDULER_CALL_TIMEOUT`         | duration | `60s`   | Default handler call timeout                   |
| `PLUGIN_SCHEDULER_MAX_CALL_TIMEOUT`     | duration | `5m`    | Maximum timeout a plugin may request           |
| `PLUGIN_SCHEDULER_MAX_RETRIES`          | number   | `10`    | Maximum number of retries                      |
| `PLUGIN_SCHEDULER_MAX_RETRY_DELAY`      | duration | `10m`   | Maximum delay between retries                  |
| `PLUGIN_SCHEDULER_MAX_JITTER`           | duration | `30s`   | Maximum random start jitter                    |
| `PLUGIN_SCHEDULER_REFRESH_INTERVAL`     | duration | `30s`   | How often to re-read tasks from the database   |

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
| `DEFAULT_LANGUAGE`            | string   | `""`                     | Default interface language, for example `ru`    |
| `GLOBAL_API_URL`              | string   | `https://api.gameap.com` | Global API address — game updates               |
| `GAMES_CDN_URLS`              | list     | see below                | Game catalog sources, tried in order            |
| `TASK_REAPER_INTERVAL`        | duration | `1m`                     | How often to look for stuck tasks               |
| `TASK_REAPER_STALE_THRESHOLD` | duration | `10m`                    | Idle time after which a task is considered stuck |

The default value of `GAMES_CDN_URLS`:
`https://cdn.gameap.ru/games.json,https://cdn.gameap.com/games.json`.

## Variables that have no effect

These variables are parsed by the panel but do not affect its behavior. Do not rely on them:

| Variable                        | What is wrong                                                                          |
|---------------------------------|-----------------------------------------------------------------------------------------|
| `AUTH_SESSION_IDLE_TIMEOUT`     | There is no idle session termination. A session lives its full lifetime — 24 hours or 7 days |
| `AUTH_SESSION_IDLE_UPDATE_FREQ` | Same as above                                                                           |
| `GRPC_FILE_TRANSFER_BASE_PATH`  | The value is not passed anywhere                                                        |
| `GRPC_ENABLED`                  | No such variable exists at all. gRPC always runs and cannot be turned off. Old versions of `gameapctl` append this line to `config.env` |

Login rate limits are not configurable via environment variables either — they are set as
constants in the code. The values are listed on the [Security](/en/security.html) page.

## Minimal configuration example

```
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://gameap:password@localhost:5432/gameap?sslmode=disable

AUTH_SECRET=replace_with_32_random_bytes
ENCRYPTION_KEY=replace_with_32_random_bytes

HTTP_HOST=panel.example.com
HTTP_PORT=8025
```

Everything else falls back to the defaults.
