---
title: API and Tokens
layout: default
lang: en
category: Administration
order: 336
---

The panel is fully controlled through the HTTP API — the interface works through the same API.
A complete description of the methods with request and response schemas is available at
[openapi.gameap.io](https://openapi.gameap.io/).

This page is about how to authenticate with the API.

## Authentication Methods

| Method                          | What for                                            | Lifetime            |
|---------------------------------|------------------------------------------------------|---------------------|
| Session token                   | The interface                                        | 24 hours or 7 days  |
| Personal access token (PAT)     | Scripts, integrations, automation                    | Indefinite          |
| Short-lived token               | WebSocket and file downloads                         | 10 seconds          |

The token is passed in a header:

```http
Authorization: Bearer <token>
```

An external system can also sign a user in without a password — with a single-use
[SSO login ticket](#sso-login-tickets). The ticket itself is not a bearer token: it is exchanged
for an ordinary session token.

## Personal Access Tokens

This is the main method for automation: the token is not tied to a session, does not expire,
and has its own set of permissions.

### Creating

A token can be created on the **Tokens** page (the user menu in the top bar) or with a request:

```bash
curl -X POST https://panel.example.com:8025/api/tokens \
  -H "Authorization: Bearer <session token>" \
  -H "Content-Type: application/json" \
  -d '{"token_name": "ci-deploy", "abilities": ["server:list", "server:restart"]}'
```

| Field        | Constraints                                                          |
|--------------|----------------------------------------------------------------------|
| `token_name` | Required; up to 255 characters                                       |
| `abilities`  | Required; from 1 to 100 abilities from the list below, no duplicates |

A token can be created only from a session: a request authenticated with a personal token is
refused with `403`. Abilities with the `admin:` prefix can be requested only by an administrator —
for any other user the request fails validation.

The response contains the full token:

```json
{"token": "12|kJ3n8sQm..."}
```

> The token is shown **only once**. Only its SHA-256 hash is stored in the database, and the
> value cannot be recovered — if you lose it, issue a new one.

The token format is `{id}|{secret}`. The `|` separator is required: pass the value in full,
exactly as issued.

### Token Abilities

Abilities are specified at creation time and restrict the token independently of the user's
permissions: the token can do no more than it is allowed, and no more than its owner is allowed.

| Ability                   | What it allows                                                        |
|---------------------------|-----------------------------------------------------------------------|
| `server:list`             | Viewing the server list                                               |
| `server:start`            | Starting a server                                                     |
| `server:stop`             | Stopping a server                                                     |
| `server:restart`          | Restarting a server                                                   |
| `server:update`           | Updating a server                                                     |
| `server:console`          | Reading and writing to the console                                    |
| `server:rcon-console`     | RCON console                                                          |
| `server:rcon-players`     | Managing players via RCON                                             |
| `server:tasks-manage`     | Managing server tasks                                                 |
| `server:settings-manage`  | Managing server settings                                              |
| `admin:server:create`     | Creating servers                                                      |
| `admin:gdaemon-task:read` | Reading daemon tasks                                                  |
| `admin:user:read`         | Reading users and their server assignments                            |
| `admin:user:manage`       | Creating and updating users, assigning servers and server permissions |
| `admin:node:read`         | Reading nodes, their IP addresses and busy ports                      |
| `admin:game:read`         | Reading games and game mods                                           |
| `admin:user:sso`          | Issuing SSO login tickets for other users                             |

Abilities with the `admin:` prefix can be granted only by an administrator — a regular user's
attempt to add them fails.

Endpoints declare the abilities they require; if an endpoint declares several, the token must
carry all of them. A token without the required ability receives `403`.
Administrative endpoints are open to a personal token only when they declare an ability:

| Ability                   | Endpoints                                                                                                                                                               |
|---------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `admin:server:create`     | `POST /api/servers`, `PUT` and `DELETE /api/servers/{id}`                                                                                                               |
| `admin:gdaemon-task:read` | `GET /api/gdaemon_tasks/{id}`                                                                                                                                           |
| `admin:user:read`         | `GET /api/users`, `GET /api/users/{id}`, `GET /api/users/{id}/servers`, `GET /api/users/{id}/servers/{server}/permissions`                                              |
| `admin:user:manage`       | `POST /api/users`, `PUT /api/users/{id}`, `PUT /api/users/{id}/servers/{server}/permissions`, `PUT` and `DELETE /api/users/{id}/servers/{server}`                       |
| `admin:node:read`         | `GET /api/nodes`, `GET /api/nodes/{id}`, `GET /api/nodes/{node}/busy_ports`, `GET /api/nodes/{node}/ip_list`                                                            |
| `admin:game:read`         | `GET /api/games`, `GET /api/games/{code}`, `GET /api/games/{code}/mods`, `GET /api/game_mods`, `GET /api/game_mods/{id}`, `GET /api/game_mods/get_list_for_game/{game}` |
| `admin:user:sso`          | `POST /api/auth/sso/tickets`                                                                                                                                            |

Every other administrative endpoint — deleting users, creating and editing nodes, games and game
mods, managing plugins, `GET /api/version` — is closed to personal tokens regardless of their
abilities: the response is `403` "personal access tokens cannot access this administrative
endpoint".

Regardless of abilities, a personal token also cannot (`403`):

- assign an administrative role when creating or updating a user;
- modify a user who is an administrator, including attaching or detaching servers;
- change a user's password via `PUT /api/users/{id}`;
- manage two-factor authentication — setup, confirmation, disabling, recovery codes;
- create another token.

See [Security](/en/security.html) for the reasoning behind these limits.

The current list of abilities is available with:

```http
GET /api/tokens/abilities
```

The response is grouped (`server`, `gdaemon-task`, `user`, `node`, `game`). The `admin:`
abilities and their groups are returned only to administrators.

### Listing and Revoking

```http
GET    /api/tokens        — list your tokens
DELETE /api/tokens/{id}   — revoke a token
```

The list shows the name, abilities, and last-used time — handy for finding unused tokens.

> **Changing the password revokes tokens.** All personal tokens created before the password
> change stop working. After changing the password, issue the tokens again.

### Usage Example

```bash
TOKEN='12|kJ3n8sQm...'

# server list
curl -H "Authorization: Bearer $TOKEN" \
  https://panel.example.com:8025/api/servers

# restart a server
curl -X POST -H "Authorization: Bearer $TOKEN" \
  https://panel.example.com:8025/api/servers/1/restart
```

## Session Token

Issued on login with a username and password:

```bash
curl -X POST https://panel.example.com:8025/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login": "admin", "password": "..."}'
```

The token format is PASETO v4.local. A regular session lives 24 hours, or 7 days with
"remember me" checked. Logging out (`POST /api/auth/logout`) puts the token on the revocation
list, which is checked on every request.

If the user has two-factor authentication enabled, login returns `two_factor_required` together
with a `challenge_token` instead of a token; the second factor is confirmed with
`POST /api/auth/2fa/verify`. See [Security](/en/security.html) for details.

Session tokens are inconvenient for automation: they expire, and login is protected by rate
limiting and possibly a CAPTCHA. Use personal tokens.

## Short-Lived Tokens

Needed where the token has to be passed in the page address — WebSocket connections and file
downloads. They are issued with `POST /api/auth/short-lived-token`, carry the `glst_` prefix,
are single-use, and live no longer than 10 seconds regardless of settings.

Only these tokens are accepted in the `?token=` query parameter; a personal token cannot be
passed there — this keeps it out of web server logs and browser history.

## SSO Login Tickets

A single-use ticket lets an external system that already knows the user — for example a billing
panel with an "open the game panel" button — sign that user into GameAP without a password.
There is no interface for issuing tickets: this is an API-only flow.

**1. Issue a ticket** with `POST /api/auth/sso/tickets` — using an administrator's session token
or a personal token with the `admin:user:sso` ability:

```bash
curl -X POST https://panel.example.com:8025/api/auth/sso/tickets \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 42, "redirect_to": "/servers/6"}'
```

| Field         | Meaning                                                                                                                                                                                                              |
|---------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `user_id`     | Required. The user the ticket signs in                                                                                                                                                                               |
| `redirect_to` | Optional. Path inside the panel to open after sign-in; must be a relative path starting with `/`, otherwise `422`                                                                                                    |
| `client_ip`   | Optional. Binds the ticket to the browser's IP address — the exchange from any other address is refused. Must be a literal IPv4 or IPv6 address; only useful behind a proxy with `AUDIT_CLIENT_IP_HEADER` configured |

The response:

```json
{"ticket": "glsso_aB3dE5fG7hJ9kL1mN3pQ5rS7tU9vW1xY3zA5bC7dE9f", "expires_in": 60, "redirect_to": "/servers/6"}
```

The ticket carries the `glsso_` prefix, is single-use, and lives for `AUTH_SSO_TICKET_TTL`
(60 seconds by default, at most 120). A ticket for an administrator is issued only for the account
the request itself is authenticated as; for any other administrator the response is `403`.
An unknown `user_id` returns `404`.

**2. Deliver the ticket to the browser** in the URL fragment:

```
https://panel.example.com:8025/sso#t=<ticket>
```

The fragment is not sent to the server, so the ticket does not end up in web server or proxy
logs. The `/sso` page reads it and performs the exchange itself.

**3. Exchange** with `POST /api/auth/sso/exchange` and the body `{"ticket": "glsso_..."}` — no
authentication. The ticket is consumed atomically, so a repeated request loses. The response is
the same as for a password login:

- an ordinary session — `token`, `expires_in`, `user`, `redirect_to`;
- with two-factor authentication enabled — `two_factor_required: true` and a `challenge_token`;
  finish with `POST /api/auth/2fa/verify`. SSO does not bypass the second factor;
- for an administrator without a second factor — the admin MFA policy applies exactly as at
  password login: `mfa_nudge` during the grace period, then `mfa_enrollment_required: true` with a
  token scoped to the 2FA enrolment endpoints.

An invalid, expired, already used or IP-mismatched ticket is refused with `401`, as is a ticket
for an account that turns out to be an administrator other than the issuer (for example, promoted
after the ticket was issued). Failed exchanges are rate-limited per address under a counter of
their own — 60 within 15 minutes.

Everywhere except the exchange endpoint the ticket is worthless: presented in the
`Authorization` header, the query string or a cookie it is rejected with `401`.

> Behind a load balancer the ticket may be issued on one panel instance and redeemed on another —
> a shared cache (`CACHE_DRIVER=redis`, `mysql` or `postgres`) is required.
> See [Multiple Panel Instances](/en/multi_instance.html).

The threat model and the restrictions on administrators are described in
[Security](/en/security.html).

## Version Endpoint

`GET /api/version` returns the running panel version and, when the update check is enabled, the
latest available releases of the panel and GameAP Daemon. Administrators only, and only with a
session token: no ability exists for this endpoint, so a personal token is refused with `403`.

```json
{
  "panel": {
    "current": "4.4.1",
    "build_date": "2026-08-06T14:42:13Z",
    "is_release": true,
    "latest_stable": "4.4.2",
    "latest_stable_url": "https://github.com/gameap/gameap/releases/tag/v4.4.2",
    "update_available": true
  },
  "daemon": {
    "latest_stable": "4.1.2",
    "latest_stable_url": "https://github.com/gameap/daemon/releases/tag/v4.1.2"
  },
  "update_check_enabled": true
}
```

| Field                                              | Meaning                                                                                                                |
|----------------------------------------------------|------------------------------------------------------------------------------------------------------------------------|
| `panel.current`, `panel.build_date`                | Version and build date of the running panel                                                                            |
| `panel.is_release`                                 | `false` for development builds whose version is not a valid semantic version; they are never compared against releases |
| `panel.latest_stable`, `panel.latest_stable_url`   | Latest stable release and its release page                                                                             |
| `panel.latest_beta`, `panel.latest_beta_url`       | Latest pre-release; present only when it is newer than the latest stable release                                       |
| `panel.update_available`                           | Whether the latest stable release is newer than the running version                                                    |
| `daemon.latest_stable`, `daemon.latest_stable_url` | Latest stable GameAP Daemon release and its release page                                                               |
| `daemon.latest_beta`, `daemon.latest_beta_url`     | Latest GameAP Daemon pre-release; present only when newer than the stable one                                          |
| `update_check_enabled`                             | Whether the panel is allowed to check for new releases                                                                 |

Network failures never fail the request: when the release source is unreachable, or the update
check is disabled, the `latest_*` fields are simply omitted. The update check is configured with
the `UPDATE_CHECK_*` variables — see the [config.env Reference](/en/config.html).

## Limits and Response Codes

| Code  | Meaning                                                                                                      |
|-------|--------------------------------------------------------------------------------------------------------------|
| `204` | Success without a response body — for example, revoking a token, attaching or detaching a server from a user |
| `400` | Malformed request body                                                                                       |
| `401` | The token is missing, invalid, or revoked; an invalid or expired SSO ticket                                  |
| `403` | The token or the user lacks permissions                                                                      |
| `404` | Object not found                                                                                             |
| `409` | Conflict — for example, a user with this login or e-mail already exists                                      |
| `422` | Request validation error                                                                                     |
| `429` | Attempt limit exceeded — login, second-factor verification, SSO ticket exchange                              |

Rate limiting applies to login and second-factor verification — 20 failed attempts per address
and 5 per login within 15 minutes — and to the SSO ticket exchange, which has a counter of its
own: 60 failed attempts per address within 15 minutes. Requests with a personal token are not
rate-limited.

## CORS

If the API is called from a browser on a different origin, list the allowed origins in
`HTTP_ALLOWED_ORIGINS` — in full, including the scheme. The `*` wildcard is not supported.
See the [config.env Reference](/en/config.html).
