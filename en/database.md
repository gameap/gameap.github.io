---
title: Database
layout: default
lang: en
category: Administration
order: 335
---

The panel works with PostgreSQL, MySQL or MariaDB, and SQLite. The choice is set with two
variables in `config.env`:

```dotenv
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://gameap:password@localhost:5432/gameap?sslmode=disable
```

## Which One to Choose

| DBMS               | When it fits                                                                                                                                               |
|--------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **SQLite**         | A single panel, up to a few dozen game servers. No separate service needed                                                                                 |
| **PostgreSQL**     | When high performance is needed. Recommended for running multiple panel instances                                                                          |
| **MySQL, MariaDB** | The familiar option, including when upgrading from GameAP 3. Slightly slower than PostgreSQL, but still suitable for multiple instances                    |

For a typical single-server installation SQLite is enough: it requires no separate service, no
setup, and no backup more complex than copying a file.

## Connection String

### PostgreSQL

```dotenv
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://user:password@host:5432/database?sslmode=disable
```

The driver names are interchangeable: `postgres`, `postgresql`, `pgsql`, `pg`, `pgx`.

The `sslmode` parameter sets the connection encryption mode: `disable` for a local database,
`require` or `verify-full` for a remote one.

### MySQL and MariaDB

```dotenv
DATABASE_DRIVER=mysql
DATABASE_URL=user:password@tcp(host:3306)/database?parseTime=true
```

The format differs from the usual URL with a scheme — it is the Go driver format. The
`parseTime=true` parameter is required.

Connecting through a socket:

```dotenv
DATABASE_URL=user:password@unix(/var/run/mysqld/mysqld.sock)/database?parseTime=true
```

### SQLite

```dotenv
DATABASE_DRIVER=sqlite
DATABASE_URL=file:/var/lib/gameap/db.sqlite?_busy_timeout=5000&_journal_mode=WAL&cache=shared
```

The file is created automatically. The parameters in the example enable WAL journaling and lock
waiting — keep them, they noticeably improve behavior under concurrent requests.

The directory with the database file must be writable by the user the panel runs as.

### inmemory

```dotenv
DATABASE_DRIVER=inmemory
DATABASE_URL=inmemory
```

Data is kept in RAM only and is lost on restart. Meant for tests; not suitable for a production
installation.

`DATABASE_URL` must be set even here: the panel checks that it is non-empty before it looks at
the driver and will not start without it. The value itself is not used.

## Migrations

The panel applies migrations itself at startup — there is no separate command. The schema
version is stored in a service table inside the same database.

Two practical rules follow from this:

* **A panel upgrade changes the schema on the very first start.** Make the backup before it.
* **Rolling back to a previous panel version without restoring the database will not work** —
  the schema has already changed. Some migrations cannot be undone even in principle: migration
  022 (4.5.0) lowercases every login and email, and case folding is irreversible — its rollback
  step does nothing; migrations 016 and 023 (4.5.0) delete duplicate rows from `plugin_storage`
  and `server_user`. A database dump is the only way back.

Migrations are applied even with gaps in the numbering, so skipping an intermediate panel
version during an upgrade is safe for the schema. It is not safe for `config.env`: when a
variable is renamed, the panel keeps reading the old name for exactly one release and then
drops it, so an operator who skips a release silently loses that setting.
`gameapctl panel upgrade` rewrites `config.env` on every upgrade and remembers the whole chain
of renames; if you upgrade by other means, check `config.env` against the
[config.env Reference](/en/config.html) afterwards.

Some migrations rewrite tables rather than just add them. Upgrading from 4.4.1 or earlier
applies `014_widen_port_columns`, which on PostgreSQL runs `ALTER TABLE … TYPE INTEGER` on the
port columns of `dedicated_servers` and `servers` (they were `SMALLINT`, so ports above 32767
did not fit). PostgreSQL rewrites those tables and holds an exclusive lock while it does, so
the first start after the upgrade can take noticeably longer on a large installation. On
MySQL/MariaDB and SQLite the same migration is a no-op.

4.5.0 also adds a unique index on `server_user (user_id, server_id)` (migration 023) after
removing duplicates. A full dump made with the commands below restores the old schema together
with the migration-version table, so the panel simply re-applies the migration on the next
start. Only a data-only restore (`--no-create-info`, copying selected tables) into an already
migrated 4.5 database can fail on duplicate pairs — deduplicate them first or restore the full
dump.

## Backup

The panel does not make backups — neither by itself nor via `gameapctl`. Set up backups with
the DBMS tools.

Restore a full dump into an empty database — create it fresh, or drop and recreate the existing
one, before running the commands below. A dump does not clear what is already in the target, so
loading it over an existing schema fails on the tables that are already there.

### PostgreSQL

```bash
pg_dump -U gameap gameap > gameap-$(date +%F).sql
```

Restoring:

```bash
psql -U gameap gameap < gameap-2026-08-02.sql
```

### MySQL and MariaDB

```bash
mysqldump -u gameap -p gameap > gameap-$(date +%F).sql
```

Restoring:

```bash
mysql -u gameap -p gameap < gameap-2026-08-02.sql
```

### SQLite

Stopping the panel is not necessary if you use the built-in command:

```bash
sqlite3 /var/lib/gameap/db.sqlite ".backup '/backup/gameap-$(date +%F).sqlite'"
```

Simply copying the file while the panel is running can produce a corrupted copy because of the
WAL journal.

### What Else to Save

The database is not enough. Along with it, save:

* **`config.env`** — it holds `AUTH_SECRET` and `ENCRYPTION_KEY`. Restore the **same values**
  that were in place when the backup was made: with a different `AUTH_SECRET`, issued tokens
  stop being accepted. If `ENCRYPTION_KEY` was set on the installation, encrypted data cannot
  be restored without it and two-factor authentication stops working for every user; since 4.5
  this also covers plugin secrets — the `plugin_secrets` table is encrypted with this key
  (AES-256-GCM) and is unrecoverable without it. If `ENCRYPTION_KEY` was never set, the TOTP
  secrets are encrypted with a key derived from `AUTH_SECRET` — then saving `AUTH_SECRET` is
  enough, and `ENCRYPTION_KEY` must **not** be added during the restore: it breaks 2FA for
  every user;
* **the panel files directory** — it holds the gRPC certificates daemons connect with, the ACME
  data and the files of installed plugins. The path is set by `FILES_LOCAL_BASE_PATH`.

Game server files live on the dedicated servers and are not part of the panel backup.

## Switching to Another DBMS

There is no built-in data transfer between DBMSes: the PostgreSQL, MySQL, and SQLite schemas
are created independently, and the panel provides no migration tool.

The procedure:

1. Back up the current database.
2. Prepare the new database and point `DATABASE_DRIVER` and `DATABASE_URL` at it.
3. Start the panel — it will create the schema from scratch.
4. Migrate the data: recreate the users, dedicated servers, games, and game servers manually or
   via the [API](/en/api.html).

Game servers themselves are untouched: their files stay on the dedicated servers; it is enough
to describe the servers in the panel with the same directories and ports.

Loading a dump from one DBMS directly into another will not work — column types and syntax
differ.

## Multiple Panel Instances

For a multi-instance installation, use PostgreSQL or MySQL. SQLite will not do: it is not
designed for several clients working with the database file at the same time.

Besides the database, you will need a shared cache, shared event exchange, and shared file
storage. See [Multiple Panel Instances](/en/multi_instance.html) for details.

## Checking

If the panel fails to start because of the database, the log will contain a connection or
migration message:

```bash
journalctl -u gameap -n 50 --no-pager
```

Common causes: a wrong `DATABASE_URL` format (especially for MySQL — it needs the Go driver
format, not a URL with a scheme), no permissions on the SQLite file directory, the database not
created, or the user not granted access to it.
