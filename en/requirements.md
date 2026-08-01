---
title: Requirements
layout: default
lang: en
category: Main
order: 10
---

* This will become a table of contents (this text will be scraped).
{:toc}

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

| System  | Architectures            |
|---------|--------------------------|
| Linux   | `amd64`, `arm64`, `386`  |
| Windows | `amd64`, `arm64`, `386`  |
| macOS   | `amd64`, `arm64`         |

Installation via `gameapctl` is described on the
[Install on Linux](/en/install/install_on_linux.html) and
[Install on Windows](/en/install/install_on_windows.html) pages.

Supported Windows versions are listed on the installation page.

## Network

| Port    | Who needs it                                               | Required |
|---------|------------------------------------------------------------|----------|
| `8025`  | Web interface and API. Needed by the administrator's browser | yes      |
| `31718` | gRPC. Daemons connect to the panel through it              | yes, if there are remote dedicated servers |
| `443`   | HTTPS, when the panel serves it itself                     | no       |
| `80`    | Let's Encrypt `http-01` challenge                          | no       |

Ports 8025 and 443 are changed with the `HTTP_PORT` and `HTTPS_PORT` variables, the gRPC port —
with `GRPC_PORT`. See the [config.env Reference](/en/config.html).

The daemon needs outbound access to the panel on port 31718 — and that is the only panel port it
requires. The panel does not connect to the daemon; no inbound ports need to be opened on the
dedicated server.

## Database

One of the following is needed:

| DBMS             | Note                                                             |
|------------------|-------------------------------------------------------------------|
| PostgreSQL       | Recommended for installations with several panel instances        |
| MySQL / MariaDB  |                                                                   |
| SQLite           | No separate database server needed, the file is created automatically |

For a small installation SQLite is enough: it requires neither a separate service nor any setup.

The connection string is set with the `DATABASE_DRIVER` and `DATABASE_URL` variables; the formats
are given in the [config.env Reference](/en/config.html).

When upgrading from GameAP 3, the panel connects to the existing MySQL database and migrates it
in place — see [Upgrade from v3 to v4](/en/upgrade_from_v3_to_v4.html).

## Optional components

Needed only in particular scenarios, not for a regular installation.

| Component        | When it is needed                                                            |
|------------------|-------------------------------------------------------------------------------|
| Redis            | Shared cache and event exchange between several panel instances               |
| S3 storage       | Shared storage for files and ACME certificates with several panel instances   |
| Reverse proxy    | When HTTPS is served by nginx or Traefik rather than the panel itself         |

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
