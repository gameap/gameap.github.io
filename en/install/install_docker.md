---
title: Install in Docker
layout: default
lang: en
category: Install GameAP
order: 102
---

The panel is distributed as a ready-made `gameap/gameap` image. The daemon is not run in Docker:
it has to manage game server processes on the dedicated server directly.

## Quick Start

The simplest option is SQLite, with no separate database:

```bash
docker run -d \
  --name gameap \
  -p 8025:8025 \
  -p 31718:31718 \
  -e DATABASE_DRIVER=sqlite \
  -e 'DATABASE_URL=file:/var/lib/gameap/db.sqlite?_busy_timeout=5000&_journal_mode=WAL&cache=shared' \
  -e AUTH_SECRET=$(openssl rand -base64 24) \
  -e ENCRYPTION_KEY=$(openssl rand -hex 32) \
  -e GRPC_EXTERNAL_HOST=panel.example.com \
  -v gameap-data:/var/lib/gameap \
  gameap/gameap:4.5
```

The panel will be available at `http://localhost:8025`. The first administrator's login and
password will be in the container log:

```bash
docker logs gameap
```

> Two parameters from the example are essential, and both are easy to miss because the
> ready-made `docker-compose.yml` files lack them. See [Port 31718](#port-31718) and
> [Panel address for daemons](#panel-address-for-daemons).

## Port 31718

Daemons connect to the panel through this port. The `Dockerfile` declares only `8025`, and the
example `docker-compose.yml` from the repository publishes only that one as well.

**Port 31718 has to be published yourself** — otherwise the panel will work, but not a single
dedicated server will connect to it.

```yaml
ports:
  - "8025:8025"
  - "31718:31718"
```

## Panel address for daemons

The panel substitutes its own address into the daemon installation command. By default it takes
it from the request header, and inside a container that is usually `localhost` or the service
name — such a command will not work on a remote dedicated server.

Set the address at which the panel is reachable from outside:

```yaml
environment:
  GRPC_EXTERNAL_HOST: panel.example.com
```

If port 31718 is published externally under a different number, set it too:

```yaml
environment:
  GRPC_EXTERNAL_PORT: 41718
```

> `GRPC_EXTERNAL_HOST` goes into the name list of the self-signed gRPC certificate. Set the
> variable before the first start if you can. If you set it later, just restart the container:
> the panel notices that the required name is missing from the existing certificate and
> re-issues it automatically — nothing needs to be deleted by hand. See
> [GRPC API](/en/daemon/grpc.html) for details.

## Image Tags

Every release is published under several tags:

| Tag                                             | What it points to                          |
|-------------------------------------------------|--------------------------------------------|
| `gameap/gameap:4.5.0`                           | The exact release                          |
| `gameap/gameap:4.5`                             | The latest patch release of the 4.5 line   |
| `gameap/gameap:4`                               | The latest release of the 4.x line         |
| `gameap/gameap:latest`                          | The latest release                         |
| `gameap/gameap:sha-<sha>`, `gameap/gameap:main` | Development builds from the `main` branch  |

Migrations run automatically at container start and are irreversible, so pin at least the minor
line (`4.5`) in `docker-compose.yml` and bump it deliberately after reading the notes for the new
version on the [Upgrade](/en/upgrade.html) page. Do not use `latest` or `main` in production.

## Docker Compose

A ready-made `docker-compose.yml` is available in the
[panel repository](https://github.com/gameap/gameap). It brings up the panel together with
PostgreSQL and Redis:

```bash
git clone https://github.com/gameap/gameap.git
cd gameap
docker compose up -d
```

You need to add the port 31718 publication and `GRPC_EXTERNAL_HOST` to it, as described above.

Passwords and keys are set via a `.env` file next to `docker-compose.yml`:

```dotenv
AUTH_SECRET=replace_with_32_random_bytes
ENCRYPTION_KEY=replace_with_32_random_bytes
POSTGRES_PASSWORD=replace_with_a_strong_password
LOGGER_LEVEL=info
```

> In the example `docker-compose.yml`, the default `AUTH_SECRET` and `ENCRYPTION_KEY` values
> are `change-me-in-production`. With such keys, session tokens are trivial to forge.
> **Replace both values before the first `docker compose up`** — before, not after.
>
> Changing `ENCRYPTION_KEY` later is not painless: the stored TOTP secrets become unreadable,
> and every user has to set up two-factor authentication again. See
> [Security](/en/security.html).

The two keys are treated differently: `AUTH_SECRET` is reduced to exactly 32 bytes, while
`ENCRYPTION_KEY` is hashed whole. That is why the generation commands in the example above
differ: `openssl rand -base64 24` yields exactly 32 characters, `openssl rand -hex 32` yields 64.

## Data and Volumes

All panel state is stored in `/var/lib/gameap`:

| What               | Where                                                      |
|--------------------|------------------------------------------------------------|
| Panel files        | The directory from `FILES_LOCAL_BASE_PATH`, `/var/lib/gameap/files` in the compose example |
| gRPC certificates  | The `certs/` subdirectory inside it                        |
| ACME certificates  | The subdirectory from `ACME_STORAGE_PATH`, `acme/` by default |
| SQLite database    | Wherever `DATABASE_URL` points                             |

Certificates and ACME data are stored not in the volume root but **inside the panel's files
directory**. With the settings from the example `docker-compose.yml`, the full path to the gRPC
certificate is `/var/lib/gameap/files/certs/server/api-server.crt`.

Mount this directory as a volume, otherwise recreating the container will lose the
certificates — and all daemons will stop connecting.

The container runs as the unprivileged `gameap` user, so the volume must have suitable
permissions.

## Environment Variables

The image is configured with the same variables as a regular installation — the full list is in
the [config.env Reference](/en/config.html). No `config.env` file is needed inside the
container.

The main variables:

| Variable             | Purpose                                                  |
|----------------------|----------------------------------------------------------|
| `DATABASE_DRIVER`    | `postgres`, `mysql`, or `sqlite`                         |
| `DATABASE_URL`       | Database connection string                               |
| `AUTH_SECRET`        | Token signing key, 32 random bytes                       |
| `ENCRYPTION_KEY`     | Secrets encryption key, 32 random bytes                  |
| `GRPC_EXTERNAL_HOST` | Panel address daemons use to connect                     |

`DATABASE_URL` and `AUTH_SECRET` are the ones the panel really needs: `AUTH_SECRET` has no default
at all, and the only default for `DATABASE_URL` is the `file:/db.sqlite` set in the image itself —
at the container root, outside the mounted volume. `DATABASE_DRIVER` is defaulted too (`sqlite` in
the image, `mysql` in the panel), so it has to be set for any other database, as the compose
example with PostgreSQL does. Set all three explicitly.

`ENCRYPTION_KEY` is formally optional, but without it the daemon connection password is stored in
plain text and plugins cannot store secrets: with the default
`PLUGINS_SECRETS_REQUIRE_ENCRYPTION=true` such writes are refused, and with
`PLUGINS_SECRETS_REQUIRE_ENCRYPTION=false` they are kept in plain text — see
[Security](/en/security.html). `GRPC_EXTERNAL_HOST` is needed when the address
the panel takes from the request is not reachable from outside; inside a container that is the
usual case — see [Panel address for daemons](#panel-address-for-daemons).

Often useful as well:

| Variable               | Default | Purpose                                                                                                           |
|------------------------|---------|-------------------------------------------------------------------------------------------------------------------|
| `DEFAULT_LANGUAGE`     | `""`    | Interface language for users who have not chosen one: `en`, `ru`, `es`, `de`. Empty — follow the browser language |
| `UPDATE_CHECK_ENABLED` | `true`  | Check for new panel and daemon versions. Set `false` for installations without outbound internet access           |

The update check queries the addresses from `UPDATE_CHECK_URLS` and caches the result for
`UPDATE_CHECK_TTL` (`6h` by default).

Plugin variables use the `PLUGINS_*` prefix since 4.5.0. If an older compose file sets
`PLUGIN_*` variables, rename them — see [Upgrading to 4.5.0](#upgrading-to-450).

## Health Check

The image includes a health check — a request to `/api/health`:

```bash
{% raw %}docker inspect --format='{{.State.Health.Status}}' gameap{% endraw %}
```

The panel log:

```bash
docker logs -f gameap
```

## Behind a Reverse Proxy

If HTTPS is served by nginx or Traefik, there is no need to configure certificates in the
panel. The proxy must pass the `X-Forwarded-Proto: https` header and **overwrite** it, not
append to it.

Port 31718 usually does not go through a reverse proxy: daemons must connect to the panel
directly. See [HTTPS and Certificates](/en/https.html) for details.

## Upgrading

First back up the database and the data volume: migrations are applied at startup and are
irreversible.

The command depends on which database the installation uses. PostgreSQL from the example
`docker-compose.yml`:

```bash
docker compose exec postgres pg_dump -U gameap gameap > gameap-backup.sql
```

MySQL or MariaDB in a neighbouring container, under its own service name:

```bash
docker compose exec mysql mysqldump -u gameap -p gameap > gameap-backup.sql
```

An external database is backed up with its own tools, on the host where it runs. A SQLite database
lies inside the volume and is covered by the copy below — provided the container is stopped, so
that the WAL journal does not leave the copy inconsistent.

The volume holds the gRPC certificates as well, and without them daemons stop connecting:

```bash
docker stop gameap
docker run --rm -v gameap-data:/data -v "$PWD:/backup" alpine \
  tar czf /backup/gameap-data.tar.gz -C /data .
```

Take the volume name from your own configuration: `gameap-data` in the quick start above,
`gameap-storage` in the example `docker-compose.yml`. What else has to be saved, and how to
restore it, is on the [Database](/en/database.html) page.

Then bump the image tag in `docker-compose.yml` rather than relying on `latest` — see
[Image Tags](#image-tags) — and pull the new image:

```bash
docker compose pull
docker compose up -d
```

General upgrade notes are on the [Upgrade](/en/upgrade.html) page.

### Upgrading to 4.5.0

There is no `gameapctl` inside the container to rewrite the configuration, so the changes below
are made by hand in `docker-compose.yml` or `.env`.

**Plugin variables were renamed.** Every `PLUGIN_*` variable became `PLUGINS_*`, and
`PLUGINS_CACHE_ENABLED` / `PLUGINS_CACHE_DIR` became `PLUGINS_RUNTIME_CACHE_ENABLED` /
`PLUGINS_RUNTIME_CACHE_DIR`. The old names keep working for one release and produce a deprecation
warning in the log; if both the old and the new name are set, the new one wins. Rename them
anyway — the old names will be removed in a future release.

Three variables also changed their value format and have **no** compatibility fallback: the old
names are silently ignored, and the panel runs with the defaults.

| Old name                          | New name                   | Value                                            |
|-----------------------------------|----------------------------|--------------------------------------------------|
| `PLUGIN_HTTP_MAX_TIMEOUT_SECONDS` | `PLUGINS_HTTP_MAX_TIMEOUT` | Duration with a unit: `30s`                      |
| `PLUGIN_NET_MAX_TIMEOUT_SECONDS`  | `PLUGINS_NET_MAX_TIMEOUT`  | Duration with a unit: `10s`                      |
| `PLUGIN_NET_READ_BUFFER_BYTES`    | `PLUGINS_NET_READ_BUFFER`  | Size, a plain byte count or with a suffix: `64K` |

The full list of variables is in the [config.env Reference](/en/config.html).

**Logins and e-mail addresses are lowercased.** On the first start of 4.5.0 a migration folds
every stored login and e-mail to lower case; users who signed in with capital letters continue
to sign in, because the panel lowercases what is typed on the login form as well. If two
accounts fold to the same login or e-mail, only one of them keeps it and the panel log names both
accounts by id, in the `user_id` and `kept_by_user_id` fields; the identifiers themselves are not
logged — see [Upgrade](/en/upgrade.html).
