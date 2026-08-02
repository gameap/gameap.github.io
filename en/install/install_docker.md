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
  -e AUTH_SECRET=$(openssl rand -hex 16) \
  -e ENCRYPTION_KEY=$(openssl rand -hex 16) \
  -e GRPC_EXTERNAL_HOST=panel.example.com \
  -v gameap-data:/var/lib/gameap \
  gameap/gameap:latest
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

> `GRPC_EXTERNAL_HOST` goes into the name list of the self-signed gRPC certificate, and the
> certificate is created on first start. Set the variable **before** the first start. If you
> set it later, delete `certs/server/api-server.crt` and `certs/server/api-server.key` in the
> panel's files directory and restart the container. See [GRPC API](/en/daemon/grpc.html) for
> details.

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

```
AUTH_SECRET=replace_with_32_random_bytes
ENCRYPTION_KEY=replace_with_32_random_bytes
POSTGRES_PASSWORD=replace_with_a_strong_password
LOGGER_LEVEL=info
```

> In the example `docker-compose.yml`, the default `AUTH_SECRET` and `ENCRYPTION_KEY` values
> are `change-me-in-production`. With such keys, session tokens are trivial to forge. Set your
> own values before the first start.
>
> Keep in mind: changing `ENCRYPTION_KEY` on a running installation breaks two-factor
> authentication for every user. See [Security](/en/security.html).

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

The bare minimum:

| Variable             | Purpose                                                 |
|----------------------|----------------------------------------------------------|
| `DATABASE_DRIVER`    | `postgres`, `mysql`, or `sqlite`                        |
| `DATABASE_URL`       | Database connection string                               |
| `AUTH_SECRET`        | Token signing key, 32 random bytes                       |
| `ENCRYPTION_KEY`     | Secrets encryption key, 32 random bytes                  |
| `GRPC_EXTERNAL_HOST` | Panel address daemons use to connect                     |

## Health Check

The image includes a health check — a request to `/api/health`:

```bash
docker inspect --format='{{.State.Health.Status}}' gameap
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

```bash
docker compose pull
docker compose up -d
```

Before upgrading, back up the database and the data volume: migrations are applied at startup
and are irreversible.

```bash
docker compose exec postgres pg_dump -U gameap gameap > gameap-backup.sql
```
