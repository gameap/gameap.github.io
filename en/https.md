---
title: HTTPS and Certificates
layout: default
lang: en
category: Administration
order: 333
---

The panel can serve HTTPS on its own, without a reverse proxy. A certificate can be taken from
files, set directly in the configuration, or obtained automatically via Let's Encrypt.

Settings are defined in `config.env` — `/etc/gameap/config.env` on Linux, `C:\gameap\web\config.env`
on Windows. After a change, a restart is required: `gameapctl panel restart`.

## Certificate sources

The panel picks the source itself, in this order:

1. **ACME** — if `ACME_ENABLED=true`, `ACME_EMAIL`, and `ACME_DOMAINS` are set.
2. **Files** — if `TLS_CERT_FILE` and `TLS_KEY_FILE` are set.
3. **Values in the configuration** — if `TLS_CERT` and `TLS_KEY` are set.
4. **No certificate** — the panel serves HTTP only.

It is **pairs** that are checked: a `TLS_CERT_FILE` alone, without `TLS_KEY_FILE`, does not count
as a configured source and is silently ignored.

HTTPS listens on the port from `HTTPS_PORT` (`443` by default) and only when a certificate is
available. HTTP on `HTTP_PORT` (`8025` by default) is always on.

> The panel certificate is unrelated to the gRPC certificates the panel uses to talk to daemons.
> Those are issued automatically by an internal certificate authority; ACME does not apply to them.
> See [GRPC API](/en/daemon/grpc.html) for details.

## Certificate from files

```
TLS_CERT_FILE=/etc/gameap/certs/panel.crt
TLS_KEY_FILE=/etc/gameap/certs/panel.key
HTTPS_PORT=443
```

The certificate file must contain the full chain: the certificate itself, then the intermediates.
Without the intermediates, some clients will not be able to verify the signature.

The files are read at startup. After replacing the certificate, restart the panel — it does not
watch the files for changes on its own.

## Certificate directly in the configuration

Convenient when the configuration is rolled out by a secret management system and extra files on
disk are undesirable.

```
TLS_CERT=LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t...
TLS_KEY=LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0t...
```

Both plain PEM and base64-encoded PEM are accepted — the panel detects the format itself. Since
the `config.env` format does not support multi-line values, in practice base64 is used:

```bash
base64 -w0 panel.crt
base64 -w0 panel.key
```

## Let's Encrypt

The panel has a built-in ACME client: the certificate is issued and renewed without external
tools like certbot.

| Variable                      | Default         | Purpose                                                    |
|-------------------------------|-----------------|-------------------------------------------------------------|
| `ACME_ENABLED`                | `false`         | Enables automatic issuance                                  |
| `ACME_EMAIL`                  | `""`            | Address for expiry notifications. Required                  |
| `ACME_DOMAINS`                | `""`            | Comma-separated domains. Required                           |
| `ACME_CHALLENGE_TYPE`         | `http-01`       | Challenge type: `http-01` or `dns-01`                       |
| `ACME_DNS_PROVIDER`           | `""`            | DNS provider, `dns-01` only                                 |
| `ACME_DIRECTORY_URL`          | production ACME | ACME directory URL                                          |
| `ACME_RENEWAL_THRESHOLD`      | `720h`          | How long before expiry to renew. 30 days by default         |
| `ACME_RENEWAL_CHECK_INTERVAL` | `12h`           | How often to check the expiry date                          |
| `ACME_PROPAGATION_TIMEOUT`    | `180s`          | How long to wait for DNS record propagation with `dns-01`   |
| `ACME_STORAGE_PATH`           | `acme`          | Directory for certificates and the ACME account key         |

> ACME is enabled **only** when `ACME_ENABLED=true`, `ACME_EMAIL`, and `ACME_DOMAINS` are all
> set, and for `dns-01` also `ACME_DNS_PROVIDER`. If anything is missing, the panel starts
> without ACME and without an error. Check the outcome via the status in the admin area.

### http-01 challenge

The default method. It requires nothing except the panel being reachable from the internet.

```
ACME_ENABLED=true
ACME_EMAIL=admin@example.com
ACME_DOMAINS=panel.example.com
ACME_CHALLENGE_TYPE=http-01
TLS_FORCE_HTTPS=true
```

What is needed:

* the domains from `ACME_DOMAINS` must resolve to this server's address;
* **requests to port 80 must reach the panel** — that is the port the certificate authority
  connects to;
* the panel serves the `/.well-known/acme-challenge/` path itself, on its HTTP port.

> By default the panel listens on port **8025**, while the certificate authority always connects
> to port **80**. On its own, these do not match. Either set `HTTP_PORT=80` or forward port 80
> to the panel port using system tools:
>
> ```bash
> iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 8025
> ```
>
> This is the most common reason `http-01` issuance fails.

Port 80 is needed not only for the first issuance: the challenge is repeated on every renewal,
so it must not be closed afterwards.

The `http-01` method **does not issue wildcard certificates** (`*.example.com`) — those require
`dns-01`.

### dns-01 challenge

Needed when the panel is not reachable from the internet on port 80 or a wildcard certificate
is required.

Of the built-in providers, only **Cloudflare** is supported. Others are added via plugins: in
that case `ACME_DNS_PROVIDER` is set to `<plugin-id>:<provider-name>`.

```
ACME_ENABLED=true
ACME_EMAIL=admin@example.com
ACME_DOMAINS=panel.example.com,*.example.com
ACME_CHALLENGE_TYPE=dns-01
ACME_DNS_PROVIDER=cloudflare
CLOUDFLARE_DNS_API_TOKEN=token_from_the_cloudflare_dashboard
```

The token is created in Cloudflare with the **Zone → DNS → Edit** permission for the zone in
question. Besides `CLOUDFLARE_DNS_API_TOKEN`, the variables `CF_DNS_API_TOKEN`,
`CLOUDFLARE_API_TOKEN`, and `CF_API_TOKEN` are accepted, as well as the legacy
"global key + email" pair: `CLOUDFLARE_API_KEY` together with `CLOUDFLARE_EMAIL`. A token with
restricted permissions is preferable.

If DNS records propagate slowly, increase `ACME_PROPAGATION_TIMEOUT`.

### Setup via gameapctl

Instead of editing `config.env` by hand, you can use the wizard:

```bash
gameapctl panel letsencrypt setup
```

It asks for the domains, the email address, and the challenge type, writes the settings to
`config.env`, and restarts the panel.

The same call without questions:

```bash
gameapctl panel letsencrypt setup --non-interactive \
  --domains=panel.example.com \
  --email=admin@example.com \
  --challenge=http-01
```

Useful flags:

| Flag                | Purpose                                                                      |
|---------------------|-------------------------------------------------------------------------------|
| `--challenge`       | `http-01` or `dns-01`                                                         |
| `--domains`         | Comma-separated domains                                                       |
| `--email`           | ACME account address                                                          |
| `--dns-provider`    | DNS provider for `dns-01`                                                     |
| `--env`             | Extra `KEY=VALUE` lines for `config.env` — for DNS credentials                |
| `--staging`         | Let's Encrypt staging directory                                               |
| `--non-interactive` | Ask no questions; fail with an error when parameters are missing              |

Disabling:

```bash
gameapctl panel letsencrypt disable
```

The command removes the `ACME_*` variables from `config.env` and restarts the panel. The
`--purge-certs` flag is declared but not implemented yet — issued certificates remain on disk.

### Debugging issuance

The production Let's Encrypt directory has strict limits on the number of attempts per domain,
and it is easy to exhaust them while setting things up. Until your setup works, use the staging
directory:

```
ACME_DIRECTORY_URL=https://acme-staging-v02.api.letsencrypt.org/directory
```

Browsers treat its certificates as untrusted, but the limits are far more relaxed. Once issuance
succeeds, remove this variable, delete the contents of the `ACME_STORAGE_PATH` directory, and
restart the panel to obtain a production certificate.

### Renewal

The panel checks the expiry date every `ACME_RENEWAL_CHECK_INTERVAL` (every 12 hours by default)
and renews the certificate when less than `ACME_RENEWAL_THRESHOLD` remains until expiry (30 days
by default). No separate scheduler or cron job is needed.

Certificates, the ACME account key, and housekeeping data are stored in the `ACME_STORAGE_PATH`
directory inside the panel's file storage. With `FILES_DRIVER=s3` they end up in S3 — this is
what lets several panel instances share one certificate.

## Certificate status

The current status is available to an administrator at `GET /api/admin/letsencrypt/status`:

```json
{
  "enabled": true,
  "state": "active",
  "challenge_type": "http-01",
  "domains": ["panel.example.com"],
  "not_after": "2026-10-30T12:00:00Z",
  "last_renewal_at": "2026-08-01T12:00:00Z",
  "next_renewal_check_at": "2026-08-02T00:00:00Z"
}
```

Possible `state` values:

| Value      | Meaning                                              |
|------------|-------------------------------------------------------|
| `disabled` | ACME is disabled                                      |
| `pending`  | The certificate has not been issued yet               |
| `active`   | The certificate is issued and valid                   |
| `renewing` | Renewal is in progress                                |
| `failed`   | The last attempt failed; the reason is in `last_error` |

## Redirecting to HTTPS

```
TLS_FORCE_HTTPS=true
```

All HTTP requests get a `301` redirect, except `/.well-known/acme-challenge/` — otherwise the
`http-01` challenge would stop working.

The same variable affects two other mechanisms: the HSTS header starts being sent even when TLS
terminates at a reverse proxy, and the CORS origin is computed with the `https` scheme.

## Panel behind a reverse proxy

If TLS terminates at nginx, Traefik, or another proxy, there is no need to configure certificates
in the panel — leave `ACME_ENABLED=false` and do not set `TLS_*`. The panel will serve HTTP on
`8025`, and the proxy will handle HTTPS.

What matters in this setup:

* The proxy must pass the `X-Forwarded-Proto: https` header, otherwise the panel will not know
  the connection is secure and will not send HSTS.
* The proxy must **overwrite** the `X-Forwarded-Proto` header and the header from
  `AUDIT_CLIENT_IP_HEADER`, not append to them: the panel trusts them without verifying the
  sender.
* Port **31718** usually does not go through the proxy — daemons must connect to the panel
  directly. Set `GRPC_EXTERNAL_HOST` to the address at which the panel is reachable by daemons.
* If the public address differs from `HTTP_HOST`, list it in `HTTP_ALLOWED_ORIGINS`.

## Common problems

| Symptom                                              | Cause                                                                                      |
|------------------------------------------------------|---------------------------------------------------------------------------------------------|
| The panel started, but HTTPS is not listening        | A variable pair is not set in full, or one of the required `ACME_*` variables is missing     |
| `state: failed` with `http-01`                       | Port 80 is not reachable from outside, the domain does not resolve to this server, or another service intercepts it |
| `state: failed` with `dns-01`                        | The DNS token lacks permissions, or the record did not propagate in time — increase `ACME_PROPAGATION_TIMEOUT` |
| Issuance stopped working after several attempts      | The production Let's Encrypt rate limit is exhausted. Switch to the staging directory and finish the setup there |
| The browser complains about the chain                | `TLS_CERT_FILE` contains only the certificate, without the intermediates                     |
| The certificate was replaced, but the old one is served | The files are read at startup — `gameapctl panel restart` is needed                       |
