---
title: Database
layout: default
lang: en
category: Administration
order: 335
---

The panel works with PostgreSQL, MySQL or MariaDB, and SQLite. The choice is set with two
variables in `config.env`:

```
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

```
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://user:password@host:5432/database?sslmode=disable
```

The driver names are interchangeable: `postgres`, `postgresql`, `pgsql`, `pg`, `pgx`.

The `sslmode` parameter sets the connection encryption mode: `disable` for a local database,
`require` or `verify-full` for a remote one.

### MySQL and MariaDB

```
DATABASE_DRIVER=mysql
DATABASE_URL=user:password@tcp(host:3306)/database?parseTime=true
```

The format differs from the usual URL with a scheme — it is the Go driver format. The
`parseTime=true` parameter is required.

Connecting through a socket:

```
DATABASE_URL=user:password@unix(/var/run/mysqld/mysqld.sock)/database?parseTime=true
```

### SQLite

```
DATABASE_DRIVER=sqlite
DATABASE_URL=file:/var/lib/gameap/db.sqlite?_busy_timeout=5000&_journal_mode=WAL&cache=shared
```

The file is created automatically. The parameters in the example enable WAL journaling and lock
waiting — keep them, they noticeably improve behavior under concurrent requests.

The directory with the database file must be writable by the user the panel runs as.

### inmemory

```
DATABASE_DRIVER=inmemory
```

Data is kept in RAM only and is lost on restart. Meant for tests; not suitable for a production
installation.

## Migrations

The panel applies migrations itself at startup — there is no separate command. The schema
version is stored in a service table inside the same database.

Two practical rules follow from this:

* **A panel upgrade changes the schema on the very first start.** Make the backup before it.
* **Rolling back to a previous panel version without restoring the database will not work** —
  the schema has already changed.

Migrations are applied even with gaps in the numbering, so skipping an intermediate panel
version during an upgrade is fine.

## Backup

The panel does not make backups — neither by itself nor via `gameapctl`. Set up backups with
the DBMS tools.

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

* **`config.env`** — it holds `AUTH_SECRET` and `ENCRYPTION_KEY`. Without `ENCRYPTION_KEY`,
  some data from the backup cannot be restored, and two-factor authentication will stop working
  for every user;
* **the panel files directory** — it holds the gRPC certificates daemons connect with and the
  ACME data. The path is set by `FILES_LOCAL_BASE_PATH`.

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
