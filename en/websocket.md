---
title: WebSocket and Metrics
layout: default
lang: en
category: Administration
order: 337
---

Over WebSocket the panel serves everything that updates in real time: task progress, the game
server console, metrics, and the progress of archive operations in the file manager. The
interface works through these same connections.

## Connecting

Connections are opened at seven addresses:

| Address                                                    | What it serves                                     |
|------------------------------------------------------------|----------------------------------------------------|
| `/api/ws/tasks/{id}`                                       | Task progress and its output                       |
| `/api/ws/servers/{server}/console`                         | The game server console, both directions           |
| `/api/ws/servers/{server}/attach`                          | An interactive session with the game server        |
| `/api/ws/servers/{server}/file-manager/archive-operations` | Progress of file manager archive operations        |
| `/api/ws/servers/{server}/metrics`                         | Game server metrics                                |
| `/api/ws/nodes/{id}/metrics`                               | Dedicated server metrics                           |
| `/api/ws/nodes/metrics`                                    | Metrics of all dedicated servers                   |

### Authorization

The token is passed in a query parameter because headers cannot be set when opening a WebSocket
from a browser:

```text
wss://panel.example.com:8025/api/ws/servers/1/console?token=<token>
```

**Only short-lived tokens** with the `glst_` prefix **are accepted in the `token` parameter**.
A personal access token cannot be passed there — this keeps it out of web server logs and
browser history.

> The short-lived token itself still appears in the URL and may end up in the logs of a reverse
> proxy, web server or monitoring system. Until it is used or expires the token is a working
> credential: whoever reads it out of a log within those seconds can open the connection in your
> place. Single use and the 10-second lifetime narrow that window but do not close it, so strip
> the `token` parameter from any logs that are retained.

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

**Archive operations** (`/api/ws/servers/{server}/file-manager/archive-operations`):

| Type               | Direction      | Purpose                                               |
|--------------------|----------------|-------------------------------------------------------|
| `archive.progress` | From the panel | Progress of creating or unpacking an archive          |
| `archive.complete` | From the panel | The operation finished, successfully or with an error |

No initial state is sent on connection: open the socket **before** starting the operation and
match frames to the `operation_id` returned in the `202` response to the create or extract
request. Both payloads carry the `operation_id`, the kind of operation and the counts of
processed files and bytes; `archive.complete` adds `success` and, on failure, `error`. See
[File Manager](/en/gameap_configure/file_manager.html#archives).

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
