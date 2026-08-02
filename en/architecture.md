---
title: How GameAP Works
layout: default
lang: en
category: Main
order: 5
---

GameAP consists of two applications: the **panel**, which the administrator works with, and the
**daemon**, which runs on every dedicated server and manages game servers.

![](/images/en/gameap_architecture.svg)

## Three Layers

**The panel** is a single application with a built-in web interface. It stores all the data:
users, dedicated servers, games, game servers, tasks. It does not start game servers itself and
does not touch their files — it only issues commands to daemons.

**GameAP Daemon** runs on every dedicated server. It starts and stops game servers, monitors
their state, installs and updates them, works with files, and collects metrics.

**Game servers** are processes the daemon manages through a process manager: systemd, Docker,
Podman, tmux, and others.

The panel and the daemon can be installed on one server — then all three layers end up on the
same machine.

## Who Connects to Whom

The connection is always established by the **daemon**: it connects to the panel itself and
keeps a single persistent connection that carries all the traffic.

The panel does not connect to the daemon. The dedicated server needs no open inbound ports and
no external IP address — outbound access to the panel is enough.

This differs from GameAP 3, where the panel connected to the daemon.

| Direction                             | Port          | Protocol         |
|---------------------------------------|---------------|------------------|
| Administrator's browser → panel       | `8025`        | HTTP, HTTPS      |
| Daemon → panel                        | `31718`       | gRPC             |
| Panel → game server                   | game port     | Query, RCON      |

Ports `8025` and `31718` are **separate listeners**, not one port with protocol detection.

The panel performs Query and RCON requests to game servers itself, directly — they do not go
through the daemon.

## What Happens When a Server Starts

1. The administrator clicks a button in the panel.
2. The panel creates a task and puts it in the database.
3. The task goes to the daemon over the established connection.
4. The daemon starts the game server through the process manager.
5. The daemon sends the progress and the command output; the panel shows them in real time over
   WebSocket.
6. From then on, the daemon regularly reports the server state and metrics.

If the connection to the panel is lost, game servers keep running — only control from the panel
is unavailable. Once the connection is restored, the daemon reconnects and receives the complete
current state from the panel.

## What Is Stored Where

| Data                                         | Where                                             |
|----------------------------------------------|---------------------------------------------------|
| Users, servers, games, tasks                 | Panel database                                    |
| gRPC certificates, ACME data                 | Panel file storage                                |
| Sessions, counters, setup key                | Panel cache                                       |
| Panel configuration                          | `config.env`                                      |
| Game server files                            | The dedicated server, in the daemon's working directory |
| Daemon configuration                         | `gameap-daemon.yaml` on the dedicated server      |
| Metrics                                      | The daemon's RAM, for no longer than an hour      |

Game server files are not copied to the panel: the file manager works with them through the
daemon.

## Extending

**Plugins** run inside the panel in a WebAssembly sandbox. They add pages, tabs, and
integrations, but have no system access of their own — only through the panel's controlled
interface. See [Plugins](/en/plugins/index.html).

**API** — everything the interface does is available through the HTTP API; the interface itself
works through it. See [API and Tokens](/en/api.html).

## Next

* [Requirements](/en/requirements.html) — what is needed for installation
* [Getting Started](/en/get_started.html) — installing the panel and the first dedicated server
* [GRPC API](/en/daemon/grpc.html) — details of the panel–daemon communication channel
* [Multiple Panel Instances](/en/multi_instance.html) — a fault-tolerant setup
