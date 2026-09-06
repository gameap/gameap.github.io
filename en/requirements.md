---
title: Requirements
description: "GameAP 4 requirements: memory and disk for the panel and the daemon, supported systems and architectures, the ports needed and the choice of database."
layout: default
lang: en
category: Main
order: 10
---

GameAP 4 is a single executable with a built-in web interface. No web server, PHP, composer, or
Node.js is needed to run the panel.

## System requirements

### Panel

* RAM: 512 MB or more
* Disk: 200 MB or more
* The panel is not demanding on the CPU, one core is enough

### GameAP Daemon

The values below are for the daemon itself, **not counting game servers**. Budget resources for
game servers separately: some games have high requirements.

* RAM: 128 MB
* Disk: 1 GB or more, plus space for game server files
* The daemon is not demanding on the CPU, one core is enough

The panel and the daemon can be installed on the same server.

## Operating systems

Builds are released for:

| System  | Architectures           |
|---------|-------------------------|
| Linux   | `amd64`, `arm64`, `386` |
| Windows | `amd64`, `arm64`, `386` |
| macOS   | `amd64`, `arm64`        |

Installation via `gameapctl` is described on the
[Install on Linux](/en/install/install_on_linux.html) and
[Install on Windows](/en/install/install_on_windows.html) pages.

Supported Windows versions are listed on the installation page.

## Network

| Port    | Who needs it                                                 | Required                                   |
|---------|--------------------------------------------------------------|--------------------------------------------|
| `8025`  | Web interface and API. Needed by the administrator's browser | yes                                        |
| `31718` | gRPC. Daemons connect to the panel through it                | yes, if there are remote dedicated servers |
| `443`   | HTTPS, when the panel serves it itself                       | no                                         |
| `80`    | Let's Encrypt `http-01` challenge                            | yes, if this method is used                |

Ports 8025 and 443 are changed with the `HTTP_PORT` and `HTTPS_PORT` variables, the gRPC port —
with `GRPC_PORT`. See the [config.env Reference](/en/config.html).

The daemon needs outbound access to the panel on port 31718 — and that is the only panel port it
requires. The panel does not connect to the daemon; no inbound ports need to be opened on the
dedicated server.

> This is how versions **4.2 and newer** work. In 4.0 and 4.1 the exchange uses the old protocol:
> the panel connects, and **inbound port 31717** must be open on the dedicated server. This matters
> if you are upgrading from GameAP 3 — that upgrade lands exactly on 4.1,
> see [Upgrade from v3 to v4](/en/upgrade_from_v3_to_v4.html).

### Outbound connections from the panel

The panel itself makes a few outbound HTTPS requests. None of them is required to run the panel:
without outbound access it only cannot check for updates, refresh the games catalog or open the
plugin catalog.

| Address                                             | Purpose                                            | Variable                                                        |
|-----------------------------------------------------|----------------------------------------------------|-----------------------------------------------------------------|
| `cdn.gameap.com`, `cdn.gameap.ru`, `api.github.com` | Checking for new GameAP and GameAP Daemon releases | `UPDATE_CHECK_URLS`; disabled with `UPDATE_CHECK_ENABLED=false` |
| `cdn.gameap.ru`, `cdn.gameap.com`                   | Games catalog                                      | `GAMES_CDN_URLS`                                                |
| `api.gameap.com`                                    | Global API — game updates                          | `GLOBAL_API_URL`                                                |
| `plugins.gameap.dev`                                | Plugin catalog                                     | `PLUGINS_STORE_URL`                                             |

## Database

One of the following is needed:

| DBMS            | Note                                                                  |
|-----------------|-----------------------------------------------------------------------|
| PostgreSQL      | Recommended for installations with several panel instances            |
| MySQL / MariaDB |                                                                       |
| SQLite          | No separate database server needed, the file is created automatically |

For a small installation SQLite is enough: it requires neither a separate service nor any setup.

The connection string is set with the `DATABASE_DRIVER` and `DATABASE_URL` variables; the formats
are given in the [config.env Reference](/en/config.html).

Upgrading from GameAP 3 while keeping the existing database is possible **only to versions 4.0
and 4.1** and only with MySQL, MariaDB or SQLite. There is no such path for PostgreSQL, nor for
versions 4.2 and later. Details and the procedure — [Upgrade from v3 to v4](/en/upgrade_from_v3_to_v4.html).

## Optional components

Needed only in particular scenarios, not for a regular installation.

| Component     | When it is needed                                                           |
|---------------|-----------------------------------------------------------------------------|
| Redis         | Shared cache and event exchange between several panel instances             |
| S3 storage    | Shared storage for files and ACME certificates with several panel instances |
| Reverse proxy | When HTTPS is served by nginx or Traefik rather than the panel itself       |

## Building from source (optional)

Regular installations use prebuilt binaries and need none of this. A toolchain is required only
when the panel or the daemon is built from source — the `--github` option of
`gameapctl panel install` and `gameapctl daemon install`, and of the matching `upgrade` commands.
See [Install on Linux](/en/install/install_on_linux.html).

* GameAP 4.5.0: Go 1.27, Node.js 24 and git;
* GameAP Daemon 4.1.2: Go 1.26.5 or newer, and git (Node.js is not needed);
* outbound access to `github.com` and the Go module proxy (`proxy.golang.org`).

When run as root, `gameapctl` installs these itself; in user scope (`--scope user`) they must
already be present.

## Dedicated server requirements

Automatic daemon installation needs:

* root (Linux) or administrator (Windows) privileges;
* the `curl`, `tar`, and `install` utilities — for the Linux installation script;
* outbound access to `github.com` and `api.github.com` — gameapctl and the daemon are downloaded
  from there;
* outbound access to the panel on port 31718.

The packages required by the daemon and game servers — including SteamCMD and a process manager —
are installed by `gameapctl` itself. See
[Dedicated Servers](/en/gameap_configure/dedicated_servers.html) for details.

### Curl

If `curl` is not present on the system, install it:

```shell
# Debian, Ubuntu
apt install curl

# CentOS, RHEL, Fedora
yum install curl
```
