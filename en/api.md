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

## Personal Access Tokens

This is the main method for automation: the token is not tied to a session, does not expire,
and has its own set of permissions.

### Creating

A token can be created in the profile or with a request:

```bash
curl -X POST https://panel.example.com:8025/api/tokens \
  -H "Authorization: Bearer <session token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "ci-deploy", "abilities": ["server:list", "server:restart"]}'
```

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

| Ability                   | What it allows                          |
|---------------------------|------------------------------------------|
| `server:list`             | Viewing the server list                  |
| `server:start`            | Starting a server                        |
| `server:stop`             | Stopping a server                        |
| `server:restart`          | Restarting a server                      |
| `server:update`           | Updating a server                        |
| `server:console`          | Reading and writing to the console       |
| `server:rcon-console`     | RCON console                             |
| `server:rcon-players`     | Managing players via RCON                |
| `server:tasks-manage`     | Managing server tasks                    |
| `server:settings-manage`  | Managing server settings                 |
| `admin:server:create`     | Creating servers                         |
| `admin:gdaemon-task:read` | Reading daemon tasks                     |

Abilities with the `admin:` prefix can be granted only by an administrator — a regular user's
attempt to add them fails.

The current list is available with:

```http
GET /api/tokens/abilities
```

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

## Limits and Response Codes

| Code  | Reason                                                                  |
|-------|--------------------------------------------------------------------------|
| `401` | The token is missing, invalid, or revoked                                |
| `403` | The token or the user lacks permissions                                  |
| `422` | Request validation error                                                 |
| `429` | Login attempt limit exceeded                                             |

Rate limiting applies only to login and second-factor verification: 20 failed attempts per
address and 5 per login within 15 minutes. Requests with a personal token are not rate-limited.

## CORS

If the API is called from a browser on a different origin, list the allowed origins in
`HTTP_ALLOWED_ORIGINS` — in full, including the scheme. The `*` wildcard is not supported.
See the [config.env Reference](/en/config.html).
