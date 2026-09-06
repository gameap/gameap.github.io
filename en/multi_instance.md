---
title: Multiple Panel Instances
layout: default
lang: en
category: Administration
order: 338
---

The panel can run as several instances behind a load balancer — for fault tolerance or to
spread the load. For that, the instances must share state: the database, the cache, files, and
event exchange.

> Most installations do not need this. A single panel handles hundreds of game servers — the
> heavy lifting is done by the daemons on the dedicated servers, not the panel.

## What Must Be Shared

| Component         | Setting                                    | What happens otherwise                                                                     |
|-------------------|--------------------------------------------|--------------------------------------------------------------------------------------------|
| Database          | PostgreSQL or MySQL                        | SQLite is not designed for several connected instances                                     |
| Cache             | `CACHE_DRIVER=redis`                       | With `memory`, sessions, setup keys and SSO login tickets are visible to one instance only |
| Event exchange    | `PUBSUB_DRIVER=redis` or `postgres`        | With `memory`, instances never learn about each other's events                             |
| Files             | `FILES_DRIVER=s3`                          | With `local`, each instance has its own files, certificates and plugin files               |

SSO login tickets show why the cache must be shared. `POST /api/auth/sso/tickets` stores a
single-use ticket in the cache, and `POST /api/auth/sso/exchange` takes it out again — usually
on another instance, whichever one the balancer sent the browser to. With `CACHE_DRIVER=memory`
the exchange succeeds only when it lands on the instance that issued the ticket, so SSO works
intermittently. The ticket lives `AUTH_SSO_TICKET_TTL` (60 seconds by default, capped at 120),
too short to compensate with retries. See [Security](/en/security.html).

Example configuration:

```dotenv
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://gameap:password@db.example.com:5432/gameap?sslmode=require

CACHE_DRIVER=redis
CACHE_REDIS_ADDR=redis.example.com:6379

PUBSUB_DRIVER=redis
PUBSUB_REDIS_ADDR=redis.example.com:6379
PUBSUB_INSTANCE_ID=panel-1

FILES_DRIVER=s3
FILES_S3_ENDPOINT=https://s3.example.com
FILES_S3_BUCKET=gameap
FILES_S3_ACCESS_KEY_ID=...
FILES_S3_SECRET_ACCESS_KEY=...
```

`AUTH_SECRET` and `ENCRYPTION_KEY` must be **identical** on all instances: otherwise a token
issued by one will not be accepted by another, and encrypted data will be unreadable. Since 4.5
this includes plugin secrets: the `plugin_secrets` table is encrypted with `ENCRYPTION_KEY`
(AES-256-GCM), so an instance with a different key cannot read a secret another instance wrote,
and with `PLUGINS_SECRETS_REQUIRE_ENCRYPTION=true` (the default) an instance without the key
refuses to write secrets at all.

## Instance Identifier

```dotenv
PUBSUB_INSTANCE_ID=panel-1
```

The value must be **unique for each instance**. Instances use it to tell their own messages
from others' and to address replies to each other.

If the variable is not set, the value `default` is used — the same on every instance, and
inter-instance exchange will not work. Set it explicitly.

## How It Works

A daemon connects to one of the instances — the one the load balancer routed it to — and keeps
a persistent connection with it. No other instance can reach that daemon directly.

So a command arriving at an instance that does not own the connection to the needed daemon is
forwarded through the shared event exchange to the instance that does; the response comes back
the same way. File operations, the console, and metrics requests work the same way.

It follows that shared event exchange is mandatory: without it, the panel will only control the
daemons that happened to connect to that particular instance.

Tasks that appeared in the database while a daemon was disconnected, or while its connection
was owned by another instance, are delivered immediately on connection rather than at the next
reconnect.

## Load Balancer and Ports

Separate the two streams:

* **HTTP (8025)** — can be balanced the usual way; session affinity is not needed: tokens are
  verified with the shared key, and session state lives in the shared cache.
* **gRPC (31718)** — connections are long-lived. The load balancer must be able to hold them
  and not cut them on an idle timeout: the channel has no keepalive mechanism of its own, and
  liveness rests on a heartbeat every 30 seconds.

Set `GRPC_EXTERNAL_HOST` to the address at which daemons see the panel through the balancer,
and `GRPC_EXTERNAL_PORT` if the port is published externally under a different number.

## Certificates

**The gRPC certificates must be shared.** They are stored in the panel's file storage, so with
`FILES_DRIVER=s3` they are automatically shared by all instances.

With `FILES_DRIVER=local`, each instance will generate its own certificate authority, and a
daemon registered through one instance will not be able to connect to another.

### HTTPS

If the panel certificate is issued via ACME, there are two things to know.

The ACME storage also lives in the file storage — with S3 it is shared.

The lock that keeps several instances from requesting a certificate at the same time works
through Redis. It is active only with `CACHE_DRIVER=redis`; with the in-memory cache the lock
is local, and the instances will get in each other's way.

The `http-01` challenge requires the certificate authority's request to land on the same
instance that started the issuance. It is simpler to use `dns-01` or terminate TLS at the
balancer. See [HTTPS and Certificates](/en/https.html).

## Limitations

**Login attempt counters** live in the cache. With `CACHE_DRIVER=redis` they are shared; with
`memory` each instance has its own, so the effective limit is multiplied by the number of
instances.

**The audit log** is written to each instance's output separately. Aggregating them requires an
external log collector, see [Security](/en/security.html).

**Plugin tasks** are deduplicated across instances with distributed locks, so a plugin's
periodic task runs once, not on every instance.

**Plugin state** is synchronized. The `plugins` table is the desired state, and every instance
reconciles against it: an install, update, uninstall, permission change or reload performed on
one instance is announced to the others through the shared event exchange (a `gameap:plugin:sync`
message) and picked up at once; the periodic pass every `PLUGINS_SYNC_REFRESH_INTERVAL`
(60 seconds by default) is the safety net for a lost message. A reload bumps the record's
`generation` counter, so **Reload** pressed on one instance restarts the plugin everywhere.

A plugin whose file is missing on an instance is downloaded from the store again and verified
against the recorded checksum. A plugin installed from a local file cannot be recovered this
way — no other instance can obtain the file — so such plugins need shared file storage
(`FILES_DRIVER=s3`). A plugin that fails to load on an instance is retried there with a backoff
growing from `PLUGINS_SYNC_MIN_BACKOFF` (15 seconds) to `PLUGINS_SYNC_MAX_BACKOFF` (15 minutes);
a plugin the runtime has disabled is restarted according to `PLUGINS_RECOVERY_*`.
`PLUGINS_SYNC_DISABLED=true` turns synchronization off — changes then reach an instance only
when it restarts. See [config.env Reference](/en/config.html).

Two things stay local to an instance:

* **SSH connections** opened by plugins through the `gameap-ssh` host library live in the memory
  of the instance that opened them and are not visible to the others (`PLUGINS_SSH_ENABLED` is
  `false` by default).
* **The plugin runtime cache** of compiled WebAssembly modules is kept in memory by default;
  when `PLUGINS_RUNTIME_CACHE_DIR` is set, it is a local path on that instance, not shared state.
  Give each instance its own directory.

## Checking

After starting, make sure that:

* a daemon registered through one instance can be controlled through another;
* a login performed on one instance is valid on the rest;
* the file manager opens regardless of which instance the request landed on;
* a plugin installed through one instance is loaded on the rest.

If any of this does not work, the cause is almost always that one of the four shared components
has remained local.
