---
title: GRPC API
layout: default
lang: en
category: GameAP Daemon
order: 405
---

Starting with GameAP 4.2 and GameAP Daemon 4.0, the panel and the daemon exchange data over gRPC
using a bidirectional stream (bidirectional streaming). This method has replaced the old exchange
over BINN and the REST API.

The connection is established by the daemon: it connects to the panel itself and holds a single
long-lived stream that carries everything — registration, heartbeat, metrics, tasks, commands, the
console, file operations. The panel does not connect to the daemon and requires no inbound ports on
the dedicated server.

## Panel configuration

There is no separate setting to enable gRPC: the server is always started. Only the address,
encryption and limits are configurable.

| Variable                      | Default    | Purpose                                                                    |
|-------------------------------|------------|----------------------------------------------------------------------------|
| `GRPC_PORT`                   | `31718`    | Port the panel's gRPC server listens on                                    |
| `GRPC_TLS_ENABLED`            | `true`     | Connection encryption                                                      |
| `GRPC_REQUIRE_MTLS`           | `false`    | Require a client certificate from the daemon                               |
| `GRPC_EXTERNAL_HOST`          | `""`       | Panel address reported to the daemon. Empty — determined from the request  |
| `GRPC_EXTERNAL_PORT`          | `0`        | Port reported to the daemon. `0` — `GRPC_PORT` is used                     |
| `GRPC_MAX_RECV_MSG_SIZE`      | `10485760` | Maximum size of an incoming message, bytes                                 |
| `GRPC_MAX_SEND_MSG_SIZE`      | `10485760` | Maximum size of an outgoing message, bytes                                 |
| `GRPC_MAX_CONCURRENT_STREAMS` | `100`      | Number of simultaneous streams per connection                              |
| `GRPC_ENABLE_REFLECTION`      | `false`    | Schema reflection for debugging tools such as `grpcurl`. Do not enable it in production |

> The `GRPC_ENABLED` variable **does not exist** — the panel does not read it. Older versions of
> `gameapctl` append a `GRPC_ENABLED=true` line to `config.env`: it is harmless, but has no effect
> whatsoever. The gRPC server cannot be turned off.

### Ports

gRPC runs on a **separate port, 31718**, while the web interface and the API run on `HTTP_PORT`
(`8025` by default). These are two different listeners, not a single port with protocol detection.

Port 31718 must be reachable from every dedicated server. It is the only panel port a working daemon
needs.

> Only port 8025 is published in the panel's `Dockerfile` and `docker-compose.yml`. When deploying
> with Docker, you have to forward port 31718 yourself.

### Panel address for the daemon

The panel substitutes its own address into the daemon installation command and into a connect URL of
the form `grpc://host:port/key`. By default the host is taken from the header of the request with
which the administrator opened the dedicated server creation page — that is, from the address in the
browser's address bar.

Set `GRPC_EXTERNAL_HOST` if that address differs from the one the daemons are supposed to connect to:

* the panel is behind a reverse proxy, gRPC does not pass through the proxy and the daemons have to
  go directly;
* the panel is behind NAT and has different addresses inside and outside;
* the panel runs in Docker, where the header ends up containing `localhost` or the container name.

`GRPC_EXTERNAL_PORT` is needed when port 31718 is published externally under a different number.

If the variables are not set, the gRPC server itself works fine — only the address in the generated
installation command comes out wrong, and the daemon will not be able to connect.

> `GRPC_EXTERNAL_HOST` goes into the list of subject alternative names (SAN) of the self-signed gRPC
> certificate. The certificate is created once, at the first start, so the variable has to be set
> **before** the panel is first started. If you set it later, the daemon will reject the connection
> because of the name mismatch in the certificate: delete `certs/server/api-server.crt` and
> `certs/server/api-server.key` and restart the panel so the certificate is generated again.

### Encryption and certificates

With `GRPC_TLS_ENABLED=true` (the default value) the panel uses a self-signed certificate issued by
its own internal certificate authority: `certs/root.crt` and `certs/root.key`. The server certificate
is `certs/server/api-server.crt`. The keys are RSA 2048 bit, valid for 10 years, and everything is
created automatically on first use.

These certificates have nothing to do with the panel's own HTTPS certificate: ACME and Let's Encrypt
do not apply to gRPC, and there is no need to configure them separately.

During registration the daemon receives `ca.crt`, `server.crt` and `server.key` from the panel and
puts them into its certificate directory. From then on it verifies the panel's certificate against
the received CA and presents its own client certificate.

In addition, every request is authenticated with the node's API key, which is stored in the panel's
database as a SHA-256 hash and compared in constant time.

### Mutual authentication (mTLS)

`GRPC_REQUIRE_MTLS=true` makes the panel require a certificate issued by its own certificate
authority from the client, and reject requests without one.

There is no need to configure the daemon specially for this: a registered daemon always presents its
certificate anyway.

> **Enable mTLS only after all daemons have been registered.** Registration itself happens without a
> client certificate — a new daemon does not have one yet. With `GRPC_REQUIRE_MTLS=true` you will not
> be able to register a new dedicated server. To add a node later, temporarily set it back to
> `false`, register the daemon and turn it back on.

Daemons registered by a different panel will not work: their certificates are issued by a foreign
certificate authority.

## Daemon configuration

Connection parameters are set in the daemon configuration — `/etc/gameap-daemon/gameap-daemon.yaml`
on Linux, `C:\gameap\daemon\gameap-daemon.yaml` on Windows — in the `grpc` block:

```yaml
grpc:
  address: panel.example.com:31718
  insecure: false
  heartbeat_interval: 30s
  connect_timeout: 30s
  initial_reconnect_delay: 1s
  max_reconnect_delay: 60s
```

| Parameter                 | Default | Purpose                                                          |
|---------------------------|---------|------------------------------------------------------------------|
| `address`                 | —       | Panel address in `host:port` form                                |
| `insecure`                | `false` | Disable TLS. For debugging only                                  |
| `heartbeat_interval`      | `30s`   | Heartbeat interval. The panel may assign its own value           |
| `connect_timeout`         | `30s`   | Connection establishment timeout                                 |
| `initial_reconnect_delay` | `1s`    | Initial pause before reconnecting                                |
| `max_reconnect_delay`     | `60s`   | Maximum pause before reconnecting                                |

If `address` is not set, it is derived from the deprecated `api_host` parameter: the host name is
taken and the port is replaced with 31718. Setting `address` explicitly is more reliable.

> The daemon has no `grpc.enabled` key either. `gameapctl` writes it during migration as a marker
> that the migration has been performed; the daemon ignores this key.

The remaining daemon configuration parameters are described on the
[GameAP Daemon](/en/daemon/daemon.html) page.

### Reconnection and connection loss

If the connection is lost, the daemon reconnects on its own, with an exponential delay: it doubles
from `initial_reconnect_delay` up to `max_reconnect_delay` and is spread out with a random jitter of
±10 %, so that many daemons do not arrive at once. With the default values that is 1 s, 2 s, 4 s, 8 s
and so on up to 60 s. After a successful connection the counter is reset. On a planned shutdown the
panel may assign the daemon a pause before its next attempt itself.

While the panel is unavailable, **game servers keep running** — only their management from the panel
is interrupted. Once the connection is restored, the daemon registers again, reports the tasks it is
currently running, and receives the full current state from the panel: the list of servers, tasks,
games and modifications, server settings. That is why tasks started before the disconnect are not
lost.

> The gRPC keepalive mechanism is not used on either side: connection liveness relies solely on the
> heartbeat every 30 seconds. If there is NAT or a firewall between the daemon and the panel that
> closes idle connections sooner, reduce `heartbeat_interval`.

## Migrating from the old protocol

A daemon installed before gRPC appeared is switched to the new protocol with the command:

```bash
gameapctl daemon upgrade --switch-to-grpc
```

What the command does:

1. Checks that the migration has not been performed yet, and determines the panel address from
   `api_host` (or takes it from `--grpc-address`).
2. Checks that the configuration contains `api_key`, `ds_id` and all three certificate files.
3. **Before making any changes** it checks that the panel is reachable: it establishes a TCP
   connection and performs a real TLS handshake with the existing certificates. If the certificates
   were issued by a different panel, the command will report this and suggest a reinstall.
4. Makes a backup copy of the configuration next to the original, with a timestamp in the name.
5. Writes the gRPC address and removes the deprecated `api_host`, `listen_ip` and `listen_port`.
6. Restarts the daemon and makes sure the panel has revoked access over the old HTTP API.
7. On any failure after step 5 it rolls the configuration back from the backup copy and starts the
   daemon again.

All that is required is that port 31718 of the panel is reachable from the dedicated server. Nothing
has to be enabled on the panel side — contrary to what the command's own error message and its
description say, the `GRPC_ENABLED` variable does not exist.

## Verification

The panel answers the standard gRPC health check and reports the `SERVING` status for the
`gameap.DaemonGateway` and `gameap.FileTransferService` services.

The simplest check that the port is reachable from the dedicated server:

```bash
nc -zv panel.example.com 31718
```

The daemon's connection status is visible in the panel on the **"Administration"** →
**"Dedicated Servers"** page. The details are in the daemon log:
`/var/log/gameap-daemon/output.log` on Linux, `C:\gameap\daemon\logs\output.log` on Windows.

Typical messages in the daemon log:

| Message                        | What it means                                                               |
|--------------------------------|-----------------------------------------------------------------------------|
| `gRPC connection failed`       | The connection was not established or was broken, followed by a pause and a new attempt |
| `registration failed: ...`     | There is a connection, but the panel rejected the registration — wrong `ds_id` or `api_key` |
| Certificate verification error | The name in the panel's certificate does not match the connection address. See `GRPC_EXTERNAL_HOST` |
