---
title: Upgrade from v3 to v4
layout: default
lang: en
category: Install GameAP
order: 191
---

Upgrading from GameAP 3 while keeping the existing database is possible only to
**versions 4.0 and 4.1**. They recognize a third-version database, extend its schema with their
own tables, and keep working with the same users, servers, and games.

> **Starting with version 4.2, an in-place upgrade is not supported.** Later versions expect the
> fourth-version schema and will not work with a GameAP 3 database.
>
> If you need the latest version of the panel, install it on a clean database and migrate the
> data yourself — see [If you need the latest version](#if-you-need-the-latest-version).

The latest version in the supported line is **4.1.2**.

## Before upgrading

> **Back up the database manually.** `gameapctl` copies the v3 installation directory and the
> web server configuration, but **does not save the database**, and the upgrade changes its
> schema irreversibly.

```shell
# MySQL or MariaDB
mysqldump -u root -p gameap > gameap-v3-backup.sql

# PostgreSQL
pg_dump -U gameap gameap > gameap-v3-backup.sql

# SQLite
sqlite3 /var/www/gameap/database.sqlite ".backup 'gameap-v3-backup.sqlite'"
```

Check that the backup is not empty, and only then continue.

> For SQLite use `.backup`, not `cp`. Simply copying the file while the panel is running can
> produce a corrupted copy: at that moment part of the data is still in the WAL journal and has
> not yet been written to the main file.

### Which databases can be upgraded in place

| DBMS in GameAP 3 | In-place upgrade |
|------------------|------------------|
| MySQL, MariaDB   | yes              |
| SQLite           | yes              |
| PostgreSQL       | **no**           |

For MySQL and SQLite, the panel recognizes an existing GameAP 3 database and does not try to
create the tables anew. For PostgreSQL there is no such recognition: the very first migration
will try to create tables that already exist and fail.

If GameAP 3 runs on PostgreSQL, install GameAP 4 on a clean database and migrate the data
yourself.

## Upgrading to 4.1

`gameapctl` is installed together with the panel and can perform the transition.

```shell
gameapctl self-update
gameapctl panel upgrade
```

The utility will detect that the third version is installed and perform the upgrade.

> **Check which version `gameapctl` installs.** When moving from the third version, it ignores
> the `--version` flag and installs the latest stable version from the 4.x line. If that turns
> out to be 4.2 or newer, the production database must not be upgraded.
>
> First test the transition on a copy of the database — the procedure is described in
> [Testing on a database copy](#testing-on-a-database-copy). Make sure 4.0 or 4.1 got installed,
> and only then upgrade the production installation.

What happens during the upgrade:

1. The `.env` file of the GameAP 3 installation is read, and the database connection parameters
   are taken from it.
2. The v3 installation directory and the web server configuration are copied to a temporary
   directory. The path to the copy is printed to the console — write it down.
3. The GameAP 4 configuration `config.env` is created, with the same database plugged in and
   new `AUTH_SECRET` and `ENCRYPTION_KEY` keys generated.
4. GameAP 4 is installed and started.

On first start the panel applies migrations: it detects that the database belongs to the third
version, skips table creation, and adds only what the fourth version is missing.

### gameapctl is not installed

Download it from the [releases page](https://github.com/gameap/gameapctl/releases) or with a
command that substitutes the right version:

```shell
VERSION=$(curl -sL https://api.github.com/repos/gameap/gameapctl/releases/latest \
  | grep -m1 '"tag_name"' | cut -d'"' -f4)
curl -OL "https://github.com/gameap/gameapctl/releases/download/${VERSION}/gameapctl-${VERSION}-linux-amd64.tar.gz"
tar xvfz "gameapctl-${VERSION}-linux-amd64.tar.gz" -C /usr/local/bin
```

## What changes with the transition

**A web server and PHP are no longer needed.** GameAP 4 is a single executable with a built-in
web interface. After the upgrade, nginx or Apache can be kept as a reverse proxy or removed
entirely.

**The default port is 8025**, not 80.

**User passwords keep working.** Versions 4.0 and 4.1 verify passwords the same way the third
version did. There is no need to ask users to change their passwords.

**Daemons keep working over the old protocol.** In versions 4.0 and 4.1 the panel exchanges
data with GameAP Daemon the same way the third version did, and connects to the daemon itself.
Dedicated server settings need no changes, and the daemon port (`31717` by default) must stay
open.

> Do not run `gameapctl daemon upgrade --switch-to-grpc` after upgrading to 4.1. gRPC data
> exchange only appeared in version 4.2, and a 4.1 panel will not be able to connect to such a
> daemon.

## Rollback

If something does not work after the upgrade:

1. Stop GameAP 4: `gameapctl panel stop`.
2. Restore the database from the backup made before the upgrade.
3. Bring back the v3 installation directory from the temporary copy whose path `gameapctl`
   printed.
4. Restore the web server configuration and start it.

> A rollback without a database backup is impossible: the schema has been changed, and GameAP 3
> will no longer work with it.

## Testing on a database copy

Before upgrading the production installation, it is worth running the transition on a copy of
the database. It also shows which exact version `gameapctl` will install.

> Do not point the trial installation at the production database: on first start the panel will
> apply migrations to it, and GameAP 3 will stop working.

Create a copy of the database:

```shell
mysqldump -u root -p gameap > /tmp/gameap.sql
mysql -u root -p -e "CREATE DATABASE gameap_v4_test"
mysql -u root -p gameap_v4_test < /tmp/gameap.sql
```

Download a version from the 4.1 line from the
[releases page](https://github.com/gameap/gameap/releases) and unpack it:

```shell
curl -OL https://github.com/gameap/gameap/releases/download/v4.1.2/gameap-v4.1.2-linux-amd64.tar.gz
tar xvfz gameap-v4.1.2-linux-amd64.tar.gz -C /usr/bin
```

Create a separate configuration file `/etc/gameap-v4-test/config.env`:

```dotenv
DATABASE_DRIVER=mysql
DATABASE_URL=gameap:password@tcp(127.0.0.1:3306)/gameap_v4_test?parseTime=true

# generate the values with: openssl rand -base64 24
AUTH_SECRET=replace_with_32_random_bytes
ENCRYPTION_KEY=replace_with_32_random_bytes

HTTP_PORT=8125
GRPC_PORT=31818
```

> The paths, ports and service name here deliberately differ from the production ones. If the
> trial installation is deployed on the same `/etc/gameap`, `8025` and `31718`, it will overwrite
> the configuration of the production panel and take over its ports.

Run:

```shell
gameap --env /etc/gameap-v4-test/config.env
```

The panel will be available on port 8125. The full list of configuration parameters is in the
[config.env Reference](/en/config.html).

Make sure the data is intact: users can log in, game servers and games are displayed. After
that you can upgrade the production installation.

### Setting up a systemd service

For convenient management of the trial installation, create a separate service.

Create a user and a directory — with names that differ from the production installation:

```shell
useradd -r -s /usr/sbin/nologin -d /var/lib/gameap-v4-test gameap-test
mkdir -p /var/lib/gameap-v4-test
chown gameap-test:gameap-test /var/lib/gameap-v4-test
```

Then the file `/etc/systemd/system/gameap-v4-test.service`:

```ini
[Unit]
Description=GameAP 4 (trial installation)
Documentation=https://docs.gameap.com
After=network.target
Wants=network-online.target
Requires=network.target

[Service]
Type=simple
User=gameap-test
Group=gameap-test

WorkingDirectory=/var/lib/gameap-v4-test

ExecStart=/usr/bin/gameap --env /etc/gameap-v4-test/config.env

# Graceful stop
ExecStop=/bin/kill -TERM $MAINPID
KillMode=mixed
KillSignal=SIGTERM
TimeoutStopSec=30

# Restart policy
Restart=always
RestartSec=5
StartLimitInterval=60
StartLimitBurst=3

RuntimeDirectory=gameap-v4-test
PIDFile=/run/gameap-v4-test/gameap.pid

# Filesystem permissions
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true

ReadWritePaths=/var/lib/gameap-v4-test

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=gameap-v4-test

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```shell
systemctl daemon-reload
systemctl enable gameap-v4-test
systemctl start gameap-v4-test
```

When the trial installation is no longer needed, remove it entirely so that it does not start
with the system:

```shell
systemctl disable --now gameap-v4-test
rm /etc/systemd/system/gameap-v4-test.service
systemctl daemon-reload
```

## If you need the latest version

You cannot upgrade from GameAP 3 straight to 4.2 or newer. The procedure is:

1. Install the current version of GameAP 4 on a **clean** database, leaving the production
   installation untouched — see [Install on Linux](/en/install/install_on_linux.html).
2. Migrate the data from GameAP 3: recreate the users, dedicated servers, games, and game
   servers — manually or via the [API](https://openapi.gameap.io/).
3. Make sure everything works, and only then decommission the old installation.

Game server files do not need to be moved anywhere: they stay on the dedicated server; it is
enough to describe the servers in the new panel with the same directories and ports.

**Daemons, however, have to be registered again.** A clean panel has its own certificate
authority, its own dedicated server ids and its own access keys, while the configuration of a
running daemon still holds the `ds_id`, `api_key` and certificates of the old panel — pointing it
at the same paths and ports is not enough.

For each dedicated server:

1. In the new panel open **"Administration"** → **"Dedicated servers"** → **"Create"** and take
   the installation command or the connect URL from there.
2. On the dedicated server run the enrollment:
   `gameap-daemon enroll --connect=grpc://new-panel:31718/key`
3. Restart the daemon: `systemctl restart gameap-daemon`.
4. Make sure it has connected: the dedicated server appears in the panel, and the daemon log
   contains no `gRPC connection failed`.

> The `enroll` command overwrites the daemon configuration file entirely, without merging and
> without a backup. Save the old file if it contained manual settings — the process manager, the
> Steam account, repository address replacements.
