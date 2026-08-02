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

| Component         | Setting                                    | What happens otherwise                                      |
|-------------------|--------------------------------------------|--------------------------------------------------------------|
| Database          | PostgreSQL or MySQL                        | SQLite is not designed for several connected instances       |
| Cache             | `CACHE_DRIVER=redis`                       | With `memory`, sessions and setup keys are visible to one instance only |
| Event exchange    | `PUBSUB_DRIVER=redis` or `postgres`        | With `memory`, instances never learn about each other's events |
| Files             | `FILES_DRIVER=s3`                          | With `local`, each instance has its own files and its own certificates |

Example configuration:

```
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
issued by one will not be accepted by another, and encrypted data will be unreadable.

## Instance Identifier

```
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

## Checking

After starting, make sure that:

* a daemon registered through one instance can be controlled through another;
* a login performed on one instance is valid on the rest;
* the file manager opens regardless of which instance the request landed on.

If any of this does not work, the cause is almost always that one of the four shared components
has remained local.
