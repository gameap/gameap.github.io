---
title: GRPC API
layout: default
lang: en
category: GameAP Daemon
order: 405
---

Starting with GameAP 4.2 and GameAP Daemon 4.0, support for the GRPC Bidirectional Streaming API has
been added for exchanging data between the panel and the daemon.
The new way of exchanging data has replaced the old one that used the BINN and REST API protocols.

The GRPC Bidi API provides more efficient and reliable realtime communication between the panel and
the daemon.

## Basic configuration

### GameAP

To use the GRPC API, the following parameters must be specified in the GameAP configuration (/etc/gameap/config.env):

```
GRPC_ENABLED=true
GRPC_PORT=31718
GRPC_TLS_ENABLED=true
```

### GameAP Daemon

