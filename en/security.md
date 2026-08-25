---
title: Security
layout: default
lang: en
category: Administration
order: 331
---

Security settings are defined by environment variables in the panel configuration file —
`/etc/gameap/config.env` on Linux, `C:\gameap\web\config.env` on Windows. After changing the file,
the panel has to be restarted: `gameapctl panel restart`.

## Two-factor authentication

### Mandatory 2FA for administrators

**Enabled by default.** An administrator without two-factor authentication first sees a reminder, and
after 30 days a requirement to enable it, without which the panel cannot be used.

| Variable                        | Default | Purpose                                                              |
|---------------------------------|---------|----------------------------------------------------------------------|
| `AUTH_REQUIRE_MFA_FOR_ADMINS`   | `true`  | Require 2FA from administrators. `false` disables the mechanism entirely |
| `AUTH_MFA_HARD_FAIL_DAYS`       | `30`    | How many days are given to enable it. `0` — reminder only, no lockout |
| `AUTH_MFA_ENROLLMENT_TOKEN_TTL` | `15m`   | Lifetime of the restricted session issued after the deadline         |

An administrator is a user who has the **`admin roles & permissions`** permission granted globally —
whether through a role or directly. There is no separate "administrator role" in the panel, and the
user with `id = 1` gets no privileges by default.

### How the 30 days are counted

The countdown starts **at the administrator's first successful login without 2FA** — not at panel
installation, not at the upgrade date, and not when the setting was enabled.

The date of the first reminder is stored in the `metadata` field of the `users` table under the
`mfa_first_shown_at` key. If an administrator has never logged in to the panel, their countdown has
not started yet.

It follows that after upgrading an existing installation, every administrator gets the full 30 days
starting from their next login.

The "Remind me later" button postpones the dialog for 24 hours, but **does not move the deadline**.
Closing the dialog with the cross or the Esc key works the same way.

### What happens after the deadline

Login still succeeds, but instead of a normal session a restricted token valid for 15 minutes is
issued. It gives access to only five routes, the ones needed to enable 2FA:

* `GET /api/config/public`
* `POST /api/auth/logout`
* `GET /api/profile`
* `POST /api/profile/2fa/setup`
* `POST /api/profile/2fa/confirm`

All other requests return `403` with the message `session is restricted to two-factor enrollment`.
The interface shows a modal dialog with no close button.

> This is not an account lockout. Enable 2FA, log in again — and work continues.
> A restricted token cannot be "upgraded" to a full one, so once 2FA is enabled the panel logs you
> out itself and asks you to log in again.

Two exceptions worth knowing about in advance:

* **A full session issued before the deadline keeps working** until it expires. A "remember me" token
  obtained on day 29 will remain valid for another 7 days.
* **Personal access tokens (PATs) are not restricted at all.** A token issued in advance will keep
  working with the API after the deadline — a sensible fallback for automation.

### Enabling 2FA

Profile → two-factor authentication. The panel shows a QR code and a secret for manual entry, after
which you enter the code from your authenticator app.

The parameters are fixed and not configurable: **TOTP per RFC 6238, HMAC-SHA1, 6 digits, 30-second
period**, tolerance ±1 step (that is, roughly ±30 seconds of clock drift). This combination is
supported by all common apps — Google Authenticator, Authy, 1Password and others. In the app the
entry will be named `GameAP`, and the account name will be your login.

A code that has been used cannot be used again, even within its 30-second window.

If the clock on the server or on the phone drifts by more than half a minute, codes stop matching —
check time synchronization on both sides.

### Recovery codes

When 2FA is enabled, **10 recovery codes** of the form `abcde-fghjk` are issued. The alphabet
contains no vowels and no easily confused characters (`0`, `o`, `1`, `l`, `i`).

* Each code is single-use.
* The codes are shown **exactly once** — when 2FA is enabled. Only their hashes are stored in the
  database, so there is no way to view the codes again.
* A recovery code can be used both to log in and to disable 2FA.
* Regenerating: profile → regenerate recovery codes, your password will be required. All previous
  codes are invalidated.

Save the codes right away, and not in the same password manager that holds your panel password.

### Lost access

There is **no** built-in CLI command to reset 2FA: neither `gameapctl` nor the panel itself can
disable two-factor authentication for another user. The options, in order:

**1. A recovery code.** Enter it instead of the code from the app.

**2. The deadline has passed, but 2FA is not enabled yet.** This is not lost access: log in as usual
and finish enabling it — the pages needed for that are available.

**3. Remove the requirement entirely.** In `config.env`:

```dotenv
AUTH_REQUIRE_MFA_FOR_ADMINS=false
```

and `gameapctl panel restart`. The restriction is lifted, and already-issued restricted tokens become
full ones.

> This option **will not help someone who has already enabled TOTP and lost their device**: the second
> factor is verified before the requirement is checked, so the panel will keep asking for a code.

**4. Keep the reminder, but remove the lockout.**

```dotenv
AUTH_MFA_HARD_FAIL_DAYS=0
```

**5. Editing the database.** The only way left when both the device and the recovery codes are lost.
Stop the panel, back up the database and run the query.

PostgreSQL:

```sql
UPDATE users
   SET two_factor_enabled = false,
       two_factor_secret = NULL,
       two_factor_recovery_codes = NULL,
       two_factor_last_used_step = NULL
 WHERE login = 'admin';
```

MySQL and SQLite — the same, but with `two_factor_enabled = 0`.

To reset the 30-day countdown as well, remove the `mfa_first_shown_at` key from the `metadata` field.
In PostgreSQL, where this field is of type `JSONB`:

```sql
UPDATE users SET metadata = metadata - 'mfa_first_shown_at' WHERE login = 'admin';
```

In MySQL the field is stored as JSON text, and the key is removed like this:

```sql
UPDATE users SET metadata = JSON_REMOVE(metadata, '$.mfa_first_shown_at') WHERE login = 'admin';
```

In SQLite — starting with version 3.38:

```sql
UPDATE users SET metadata = json_remove(metadata, '$.mfa_first_shown_at') WHERE login = 'admin';
```

> Do not clear the `metadata` field entirely (`SET metadata = NULL`): besides the 2FA countdown it
> may hold other information about the user, and that would be lost.

After that, start the panel and enable 2FA again.

## Passwords

Password requirements: **no fewer than 12 and no more than 128 bytes**. There are deliberately no
composition requirements (uppercase letters, digits, special characters) — a check against a list of
compromised passwords is used instead.

> The limit is counted in bytes, not characters. A password of 12 Cyrillic letters takes 24 bytes and
> passes the check with room to spare.

The list of common passwords is taken from [SecLists](https://github.com/danielmiessler/SecLists)
(`xato-net-10-million-passwords`), filtered by length, and contains about 46,000 entries. It is
**embedded in the binary** — the panel contacts nothing when checking a password and transmits
nothing about it.

| Variable                    | Default | Purpose                                                                |
|-----------------------------|---------|------------------------------------------------------------------------|
| `AUTH_ALLOW_WEAK_PASSWORDS` | `false` | Disables only the list check. The length limits remain                 |
| `AUTH_BCRYPT_COST`          | `13`    | bcrypt cost. Allowed range is 10 to 14, otherwise the panel will not start |

Passwords are hashed with bcrypt on top of a preliminary SHA-256, so bcrypt's 72-byte limit does not
truncate long passwords. On login, a hash with a cost lower than the current one is rehashed
automatically; the cost of already-stored hashes cannot be lowered, even if you decrease the value of
the variable.

The check applies when an administrator creates a user, when a user is modified, and when you change
your own password. It does not apply on login — otherwise users with old weak passwords would lose
access.

> The first administrator's password, set with the `ADMIN_PASSWORD` variable during initial database
> seeding, is **not** checked. Choose it deliberately.

## CAPTCHA

Disabled by default. It protects **only** the login form (`POST /api/auth/login`); second-factor
verification is not covered by the CAPTCHA.

| Variable              | Default | Purpose                                                                |
|-----------------------|---------|------------------------------------------------------------------------|
| `CAPTCHA_PROVIDER`    | `""`    | `recaptcha_v2`, `recaptcha_v3` or `turnstile`. Empty — disabled        |
| `CAPTCHA_SITE_KEY`    | `""`    | Public key, sent to the browser                                        |
| `CAPTCHA_SECRET_KEY`  | `""`    | Secret key, never sent outside                                         |
| `CAPTCHA_MIN_SCORE`   | `0.5`   | Threshold for reCAPTCHA v3 only, ignored by the other providers        |
| `CAPTCHA_FAIL_OPEN`   | `false` | Whether to allow login if the verification service is unavailable      |
| `CAPTCHA_VERIFY_URL`  | `""`    | Custom verification address — for proxying outbound traffic            |

With `CAPTCHA_FAIL_OPEN=false`, an unavailable verification service means a `503` on panel login.

> If you set `CAPTCHA_PROVIDER` but not `CAPTCHA_SECRET_KEY`, the widget will appear in the login
> form, but server-side verification will silently stay off. Set both variables together.

The panel adds the selected provider's domains to the CSP policy itself — there is no need to
configure `SECURITY_CSP_EXTRA_SCRIPT_SRC` additionally.

### reCAPTCHA v3

Keys are issued in the [reCAPTCHA console](https://www.google.com/recaptcha/admin): register the site,
choose the **reCAPTCHA v3** type and specify the panel domain.

```dotenv
CAPTCHA_PROVIDER=recaptcha_v3
CAPTCHA_SITE_KEY=6LcExampleSiteKeyExampleSiteKeyExam
CAPTCHA_SECRET_KEY=6LcExampleSecretKeyExampleSecretKeyEx
CAPTCHA_MIN_SCORE=0.5
```

reCAPTCHA v3 asks the user nothing: it returns a score from `0.0` to `1.0`, where one means almost
certainly a human. Login is rejected if the score is below `CAPTCHA_MIN_SCORE`.

Start with the default value of `0.5` and adjust it as circumstances require: if people complain they
cannot log in — lower it, if password guessing continues — raise it. A value above `0.7` noticeably
gets in the way of users with browser blockers and in incognito mode.

The panel must be opened at the domain specified in the key settings, otherwise verification will
fail. For multiple domains, list them all in the reCAPTCHA console.

Verification is performed with a request to `https://www.google.com/recaptcha/api/siteverify` — this
address must be reachable from the panel server.

### Turnstile

Keys are issued in the [Cloudflare](https://dash.cloudflare.com/) dashboard, in the **Turnstile**
section. A free account is enough, and the domain does not have to be delegated to Cloudflare.

```dotenv
CAPTCHA_PROVIDER=turnstile
CAPTCHA_SITE_KEY=0x4AAAAAAAExampleSiteKey
CAPTCHA_SECRET_KEY=0x4AAAAAAAExampleSecretKey
```

`CAPTCHA_MIN_SCORE` does not apply to Turnstile — the provider returns only "pass" or "fail", so there
is no need to set this variable.

In most cases Turnstile passes without the user noticing, and shows a short challenge when something
looks suspicious. The widget mode (**Managed**, **Non-interactive** or **Invisible**) is chosen on the
Cloudflare side when creating the key; it is not configurable from the panel.

Verification is performed with a request to
`https://challenges.cloudflare.com/turnstile/v0/siteverify`.

reCAPTCHA v2 is configured the same way as v3, but without `CAPTCHA_MIN_SCORE`.

## Login rate limiting

The limits are **hard-coded**, there are no environment variables for them:

* window — 15 minutes;
* no more than 20 failed attempts from a single IP address;
* no more than 5 failed attempts per login.

The limit is applied on two routes: `POST /api/auth/login` and `POST /api/auth/2fa/verify`. The client
gets a `429` and a `Retry-After: 900` header. Only `401` responses increment the counter; a successful
login resets the counter for the login, but not for the IP address.

The counters are stored in the panel cache. With `CACHE_DRIVER=memory` (the default value) they are
reset on restart and are not shared between multiple panel instances — for a fault-tolerant
installation use `CACHE_DRIVER=redis`.

Separately from this, second-factor verification allows no more than 5 code entry attempts per login
attempt, after which the password has to be entered again.

The panel has no account lockout — brute force is limited only by what is described above.

## HTTP headers and Content Security Policy

Security headers are enabled by default.

| Variable                            | Default                           | Header                                          |
|-------------------------------------|-----------------------------------|--------------------------------------------------|
| `SECURITY_HEADERS_ENABLED`          | `true`                            | Master switch                                    |
| `SECURITY_CONTENT_TYPE_OPTIONS`     | `true`                            | `X-Content-Type-Options: nosniff`                |
| `SECURITY_FRAME_OPTIONS`            | `SAMEORIGIN`                      | `X-Frame-Options`, an empty value removes the header |
| `SECURITY_REFERRER_POLICY`          | `strict-origin-when-cross-origin` | `Referrer-Policy`                                |
| `SECURITY_HSTS_ENABLED`             | `true`                            | `Strict-Transport-Security`                      |
| `SECURITY_HSTS_MAX_AGE`             | `31536000`                        | Duration in seconds, one year by default         |
| `SECURITY_HSTS_INCLUDE_SUBDOMAINS`  | `false`                           | Adds `includeSubDomains`                         |
| `SECURITY_HSTS_PRELOAD`             | `false`                           | Adds `preload`                                   |

HSTS is only sent when the panel is accessed over HTTPS — by an actual TLS connection, by the
`X-Forwarded-Proto: https` header or with `TLS_FORCE_HTTPS=true`. Working over plain HTTP during
development will not get the browser "stuck".

### CSP policy

| Variable                            | Default | Purpose                                                                   |
|-------------------------------------|---------|---------------------------------------------------------------------------|
| `SECURITY_CSP_ENABLED`              | `true`  | Enables the policy                                                        |
| `SECURITY_CSP_REPORT_ONLY`          | `false` | Send `Content-Security-Policy-Report-Only` instead of the blocking policy  |
| `SECURITY_CSP_POLICY`               | `""`    | Completely replaces the generated policy                                  |
| `SECURITY_CSP_REPORT_URI`           | `""`    | Adds `report-uri`                                                         |
| `SECURITY_CSP_EXTRA_SCRIPT_SRC`     | `""`    | Extends `script-src`, comma-separated values                              |
| `SECURITY_CSP_EXTRA_STYLE_SRC`      | `""`    | Extends `style-src`                                                       |
| `SECURITY_CSP_EXTRA_CONNECT_SRC`    | `""`    | Extends `connect-src`                                                     |
| `SECURITY_CSP_EXTRA_IMG_SRC`        | `""`    | Extends `img-src`                                                         |
| `SECURITY_CSP_EXTRA_FRAME_SRC`      | `""`    | Extends `frame-src`                                                       |
| `SECURITY_CSP_EXTRA_FONT_SRC`       | `""`    | Extends `font-src`                                                        |

The generated policy:

```text
default-src 'self'; base-uri 'self'; object-src 'none'; frame-ancestors 'self'; form-action 'self';
script-src 'self' blob: 'wasm-unsafe-eval' <hashes of inline scripts>;
style-src 'self' 'unsafe-inline';
img-src 'self' data: blob:;
font-src 'self';
connect-src 'self';
frame-src 'self';
worker-src 'self' blob:
```

Why it contains relaxations:

* `'wasm-unsafe-eval'` — when uploading files, the SHA-256 checksum is calculated in the browser via
  WebAssembly. Without this permission file uploads will stop working.
* `blob:` in `script-src` — this is how the frontend parts of plugins are loaded.
* `'unsafe-inline'` in `style-src` — plugin and Vue styles are added to the page inline.

The addresses of the selected CAPTCHA provider are added to `script-src` and `frame-src`
automatically.

> **A plugin that loads scripts from a third-party CDN will be blocked by the policy.** Add the
> required domain to `SECURITY_CSP_EXTRA_SCRIPT_SRC`. It is better not to use `SECURITY_CSP_POLICY`
> for this: it replaces the whole policy, together with the automatic allowance of CAPTCHA domains
> and the hashes of the panel's inline scripts.

## Secret encryption

| Variable         | Required | Purpose                                                              |
|------------------|----------|----------------------------------------------------------------------|
| `AUTH_SECRET`    | yes      | Signing key for session tokens. Without it the panel will not start   |
| `ENCRYPTION_KEY` | no       | Encryption key for secrets in the database                            |

Both values must be random rather than a passphrase. Their length requirements, however, are
**different** — the panel handles them differently.

**`AUTH_SECRET` is used as is and coerced to exactly 32 bytes:** a shorter value is padded, a longer
one is **truncated**, and a warning goes to the log. So give it exactly 32 characters:

```bash
openssl rand -base64 24
```

Do not use `openssl rand -hex 32` here: it produces 64 characters, the panel discards half of them,
and the strength stays the same.

**`ENCRYPTION_KEY` is hashed in full with SHA-256**, its length is not limited and nothing is lost.
A longer value can be used here:

```bash
openssl rand -hex 32
```

Hashing preserves the entropy of the original value but does not increase it, so the key still has
to be random. A passphrase is unsafe here: it can be brute-forced if the encrypted value leaks.

> `AUTH_SECRET` is silently coerced to 32 bytes: a shorter value is padded, a longer one is
> truncated, and only a warning goes to the log. A short or predictable `AUTH_SECRET` means session
> tokens can be forged.

`ENCRYPTION_KEY` is used to encrypt TOTP secrets and the daemon connection password in the database.
If it is not set, TOTP secrets are encrypted with a key derived from `AUTH_SECRET`, and the daemon
password is stored in plain text — the panel warns about this at startup.

> **Do not set `ENCRYPTION_KEY` for the first time on a running installation that already has 2FA
> enabled.** The encryption key for TOTP secrets will switch from `AUTH_SECRET` to `ENCRYPTION_KEY`,
> previously stored secrets will become unreadable, and **every user will have to enable 2FA again**.
> Recovery codes will keep working — they are stored separately.
>
> The same will happen if `ENCRYPTION_KEY` is lost or changed. Store it together with the database
> backup: without it, part of the data in the backup cannot be restored.

User passwords are hashed with bcrypt, personal tokens and the daemon key are stored as SHA-256 —
these are irreversible transformations, and `ENCRYPTION_KEY` has nothing to do with them.

## Audit log

| Variable                 | Default | Purpose                                              |
|--------------------------|---------|------------------------------------------------------|
| `AUDIT_ENABLED`          | `true`  | Recording of security events                         |
| `AUDIT_CLIENT_IP_HEADER` | `""`    | Header with the real client IP, for example `X-Real-IP` |

What is recorded: successful and failed logins, rate limit hits, access denials, enabling and
disabling 2FA, regenerating recovery codes, user changes and role assignments, token creation and
revocation, changes to and deletion of dedicated servers, file operations, plugin installation and
removal.

Each record contains: event type, category, outcome, the acting user's identifier and login,
authentication method, IP address, User-Agent, request method and path, request identifier.

> The audit log is a set of structured application log lines with the `component=audit` field. There
> is **no** separate database table, separate file, rotation, viewing interface or read API. If the
> records need to be stored and searched, set up collection of the panel log with the standard tools
> of your system — for example, through `journald` and an external log collector.
>
> `AUDIT_CLIENT_IP_HEADER` trusts the specified header from **any** sender — the panel has no trusted
> proxy list. Enable this variable only if the reverse proxy is guaranteed to overwrite the header in
> incoming requests. Otherwise the IP address can be spoofed, and with it the per-IP login rate limit
> can be bypassed.

## Sessions and tokens

Sessions are issued in the PASETO v4.local format (`AUTH_SERVICE=paseto`, the alternative is `jwt`).
A normal session lives for 24 hours, with the "remember me" option — 7 days. On logout the token goes
into a revocation list, which is checked on every request.

For cases where the token has to be passed in the page address — WebSocket connections, file
downloads — single-use short-lived tokens with the `glst_` prefix are issued. Their lifetime is
limited to 10 seconds regardless of the value of `AUTH_SHORT_LIVED_TOKEN_TTL`.

The panel has no CSRF protection and does not need it: authentication goes through the
`Authorization` header, not cookies.

The list of origins allowed to access the API from a browser is set with the `HTTP_ALLOWED_ORIGINS`
variable (comma-separated values). If it is empty, a single origin computed from `HTTP_HOST` is
allowed. The `*` character is not supported.

## File uploads

The type of an uploaded file is determined by its content, not by the extension and not by the header
sent by the client. By default images, text files, JSON, XML, CSV, YAML and PDF are allowed. SVG and
HTML are deliberately forbidden — they can contain scripts.

| Variable                      | Default | Purpose                                                     |
|-------------------------------|---------|-------------------------------------------------------------|
| `FILES_UPLOAD_ALLOWED_MIMES`  | `""`    | Extends the list of allowed types, does not replace it      |
| `FILES_UPLOAD_ALLOW_ARCHIVES` | `false` | Allow archives: zip, tar, gzip, bzip2, 7z, xz               |
| `FILES_UPLOAD_ALLOW_BINARY`   | `false` | Allow arbitrary binary files                                |

> The ban on archives and binary files is the most common reason for the question "why won't the file
> upload". An archive can contain executables that will be unpacked on the dedicated server, which is
> why uploading is forbidden by default. Enable these settings deliberately.

Rejected uploads go into the audit log with the detected file type and the reason for the rejection.
The size limit for a single file is 100 MB and is not configurable.

## Plugins

Plugins run in a WebAssembly sandbox and have no direct access to the system. Plugin network requests
are restricted separately:

| Variable                          | Default | Purpose                                                       |
|-----------------------------------|---------|---------------------------------------------------------------|
| `PLUGINS_DISABLED`                | `false` | Disable the plugin mechanism entirely                         |
| `PLUGIN_HTTP_BLOCK_PRIVATE_IPS`   | `true`  | Forbid requests to internal network addresses                 |
| `PLUGIN_HTTP_ALLOWED_SCHEMES`     | `https` | Allowed schemes                                               |
| `PLUGIN_HTTP_ALLOWED_HOSTS`       | `""`    | List of allowed hosts, empty — no host restrictions           |
| `PLUGIN_HTTP_MAX_TIMEOUT_SECONDS` | `30`    | Request time limit                                            |
| `PLUGIN_HTTP_MAX_REDIRECTS`       | `5`     | Redirect limit, each one is checked anew                      |

The addresses of cloud provider metadata services are always blocked and cannot be unblocked. The
`Set-Cookie`, `Authorization`, `WWW-Authenticate` and `Clear-Site-Data` headers are not passed to the
plugin.

See details on the [Plugins](/en/plugins/index.html) page.

## What is hard-coded

Some variables are present in the configuration but have no effect on the panel's behaviour. Do not
rely on them:

* **`AUTH_SESSION_IDLE_TIMEOUT` and `AUTH_SESSION_IDLE_UPDATE_FREQ`** — there is currently no session
  termination on inactivity. A session lives exactly its full term: 24 hours or 7 days.
* **`GRPC_ENABLED`** — no such setting exists, the gRPC server always runs. Older versions of
  `gameapctl` append this line to `config.env`; it is harmless, but it does nothing.
* Login rate limits are set by constants in the code, there are no `RATE_LIMIT*` variables.

The OWASP ASVS compliance percentages in the `docs/security/ASVS.md` and `ASVS_L2.md` files in the
panel repository are out of date — refer to this page instead.

## Reporting a vulnerability

Vulnerabilities are accepted through a
[GitHub Security Advisory](https://github.com/gameap/gameap/security/advisories/new) — this is the
preferred channel — or by email to `security@gameap.com`.

Timelines: acknowledgement of receipt — 72 hours, assessment — 14 days, public disclosure — 90 days.
Fixes: critical — 14 days, high — 30 days, medium — 60 days, low — in the next release.

An unintentionally insecure configuration set by the administrator themselves (for example,
`AUTH_ALLOW_WEAK_PASSWORDS=true` or `SECURITY_HEADERS_ENABLED=false`) is not considered a
vulnerability.
