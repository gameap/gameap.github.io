---
title: WebSocket and Metrics
layout: default
lang: en
category: Administration
order: 337
---

Over WebSocket the panel serves everything that updates in real time: task progress, the game
server console, and metrics. The interface works through these same connections.

## Connecting

Connections are opened at six addresses:

| Address                              | What it serves                                    |
|--------------------------------------|----------------------------------------------------|
| `/api/ws/tasks/{id}`                 | Task progress and its output                       |
| `/api/ws/servers/{server}/console`   | The game server console, both directions           |
| `/api/ws/servers/{server}/attach`    | An interactive session with the game server        |
| `/api/ws/servers/{server}/metrics`   | Game server metrics                                |
| `/api/ws/nodes/{id}/metrics`         | Dedicated server metrics                           |
| `/api/ws/nodes/metrics`              | Metrics of all dedicated servers                   |

### Authorization

The token is passed in a query parameter because headers cannot be set when opening a WebSocket
from a browser:

```text
wss://panel.example.com:8025/api/ws/servers/1/console?token=<token>
```

**Only short-lived tokens** with the `glst_` prefix **are accepted in the `token` parameter**.
A personal access token cannot be passed there — a long-lived access key will not end up in the
URL even by mistake.

> The short-lived token itself does stay in the address and can settle in the logs of a reverse
> proxy, a web server or a monitoring system. Being single-use and living for 10 seconds keeps the
> window narrow, but until it is used or expires the token still works, so if logs are kept for a
> long time, it is better to strip the `token` parameter out of them.

Getting a short-lived token:

```bash
curl -X POST https://panel.example.com:8025/api/auth/short-lived-token \
  -H "Authorization: Bearer <session token>"
```

The token is single-use and lives no longer than 10 seconds, so request it immediately before
opening the connection. See [API and Tokens](/en/api.html) for details.

## Frame Format

All messages are JSON of the same shape:

```json
{
  "type": "task.status",
  "payload": { },
  "ts": 1711234567
}
```

| Field     | Description                                  |
|-----------|----------------------------------------------|
| `type`    | Message type                                 |
| `payload` | Data, depends on the type                    |
| `ts`      | Timestamp, Unix seconds                      |

### Message Types

**Tasks** (`/api/ws/tasks/{id}`):

| Type            | When it arrives                    |
|-----------------|------------------------------------|
| `task.status`   | The task status changed            |
| `task.output`   | New command output appeared        |
| `task.complete` | The task finished                  |

**Console** (`/api/ws/servers/{server}/console`):

| Type              | Direction        | Purpose                              |
|-------------------|------------------|--------------------------------------|
| `console.history` | From the panel   | Accumulated output on connection     |
| `console.command` | To the panel     | Sending a command to the server      |

**Interactive session** (`/api/ws/servers/{server}/attach`):

| Type            | Direction   | Purpose                        |
|-----------------|-------------|--------------------------------|
| `attach.input`  | To the panel | Input into the session        |
| `attach.detach` | To the panel | Detach without stopping the server |

**Metrics**:

| Type                   | Purpose                                           |
|------------------------|---------------------------------------------------|
| `metrics.replay`       | Accumulated values for the retention period       |
| `metrics.replay.done`  | Accumulated values delivered; live values follow  |
| `metrics.error`        | Metrics collection error                          |

On connecting to metrics, the panel first serves the history for the retention period, then
sends new values as they arrive. The switch is marked by a `metrics.replay.done` frame.

> In the game server console, the RCON password is masked: the game server prints it as part of
> the startup command line.

## Metric Series

Metrics are collected by the daemon and passed to the panel. The sampling interval and
retention time are set in the daemon configuration with `metrics.collection_interval` (5
seconds by default) and `metrics.retention_duration` (10 minutes by default), see
[GameAP Daemon](/en/daemon/daemon.html#metrics-collection).

### Game Server

| Series                                      | Value                                      |
|---------------------------------------------|--------------------------------------------|
| `gameap_server_up`                          | Whether the server is running              |
| `gameap_server_cpu_usage_percent`           | CPU usage, percent                         |
| `gameap_server_memory_usage_bytes`          | Memory used, bytes                         |
| `gameap_server_memory_limit_bytes`          | Memory limit, bytes                        |
| `gameap_server_memory_usage_percent`        | Memory used, percent of the limit          |
| `gameap_server_network_receive_bytes_total` | Network bytes received, cumulative         |
| `gameap_server_network_transmit_bytes_total`| Network bytes transmitted, cumulative      |
| `gameap_server_block_io_read_bytes_total`   | Disk bytes read, cumulative                |
| `gameap_server_block_io_write_bytes_total`  | Disk bytes written, cumulative             |
| `gameap_server_process_pids`                | Number of server processes                 |

The memory limit metrics are populated only by process managers that can enforce it: systemd,
Docker, and Podman.

### Dedicated Server

| Series                                    | Value                                        |
|-------------------------------------------|----------------------------------------------|
| `gameap_node_cpu_usage_percent`           | CPU usage, percent                           |
| `gameap_node_memory_usage_bytes`          | Memory used, bytes                           |
| `gameap_node_memory_total_bytes`          | Total memory, bytes                          |
| `gameap_node_memory_usage_percent`        | Memory used, percent                         |
| `gameap_node_swap_usage_bytes`            | Swap used, bytes                             |
| `gameap_node_swap_total_bytes`            | Total swap, bytes                            |
| `gameap_node_disk_usage_bytes`            | Disk used, bytes                             |
| `gameap_node_disk_total_bytes`            | Total disk, bytes                            |
| `gameap_node_disk_usage_percent`          | Disk used, percent                           |
| `gameap_node_network_receive_bytes_total` | Network bytes received, cumulative           |
| `gameap_node_network_transmit_bytes_total`| Network bytes transmitted, cumulative        |
| `gameap_node_load1`                       | Load average over 1 minute                   |
| `gameap_node_load5`                       | Load average over 5 minutes                  |
| `gameap_node_load15`                      | Load average over 15 minutes                 |
| `gameap_node_uptime_seconds_total`        | Uptime, seconds                              |

Load average is not collected on Windows.

Which network interfaces and drives to include is set with the `if_list` and `drives_list`
parameters in the daemon configuration.

## Limitations

Metrics are kept in the daemon's RAM for no longer than `metrics.retention_duration` (10 to 60
minutes allowed) and are lost when the daemon restarts. The panel has no long-term storage or
export to monitoring systems.

Metrics collection is disabled with `metrics.enabled: false` in the daemon configuration.
