---
title: Process Managers
layout: default
lang: en
category: GameAP Daemon
order: 410
---

A process manager is a system utility that manages processes on the server.
It is responsible for starting, stopping, and restarting game servers, monitoring their status, collecting statistics, and limiting resources (CPU / RAM).
By default, systemd is used on Linux, and [Shawl](https://github.com/mtkennerly/shawl) on Windows.

## Configuration

The process manager is configured in the GameAP Daemon configuration file.

### Basic Structure

```yaml
process_manager:
  name: <manager_name>
  config:
    <key>: <value>
```

**Available managers:**
- Linux: `systemd` (default), `docker`, `podman`, `tmux`
- Windows: `shawl` (default), `winsw`

## Linux

### Systemd

Used by default on Linux. It is a modern process manager that provides high performance and reliability.
Isolation capabilities for this process manager are limited.

| Feature                                |                 |
|----------------------------------------|-----------------|
| Start, stop, restart servers           | ✅               |
| Statistics                             | ⚠️              |
| Resource limits (CPU / RAM)            | ✅               |
| Console reading                        | ✅               |
| Sending commands to console            | ✅               |
| Isolation                              | ⚠️ Limited      |

#### Systemd Configuration

Systemd does not require additional configuration. The Daemon automatically creates and manages unit files.

##### Configuration Example

```yaml
process_manager:
  name: systemd
```

### Docker

Docker is used to isolate game servers in containers.
Use Docker if you plan to use Pterodactyl Eggs or Pelican Eggs.

| Feature                                |    |
|----------------------------------------|----|
| Start, stop, restart servers           | ✅  |
| Statistics                             | ⚠️ |
| Resource limits (CPU / RAM)            | ✅  |
| Console reading                        | ✅  |
| Sending commands to console            | ✅  |
| Isolation                              | ✅  |

#### Docker Configuration

##### Configuration Examples

Minimal:
```yaml
process_manager:
  name: docker
```

### Podman

Podman is an alternative to Docker that provides isolation of game servers in containers.
You can also use Podman if you plan to use Pterodactyl Eggs or Pelican Eggs.

| Feature                                |    |
|----------------------------------------|----|
| Start, stop, restart servers           | ✅  |
| Statistics                             | ⚠️ |
| Resource limits (CPU / RAM)            | ✅  |
| Console reading                        | ✅  |
| Sending commands to console            | ✅  |
| Isolation                              | ✅  |

#### Podman Configuration

Podman uses the same parameters as Docker, plus a specific parameter for connecting to the API.

##### Specific Parameters

| Config Key    | Description | Default |
|---------------|-------------|---------|
| `socket_path` | Path to Podman unix socket | `unix:///run/user/{UID}/podman/podman.sock` (rootless) |

##### Configuration Examples

Rootless (default):
```yaml
process_manager:
  name: podman
```

Rootful:
```yaml
process_manager:
  name: podman
  config:
    socket_path: "unix:///run/podman/podman.sock"
```

### Tmux

[Tmux](https://github.com/tmux/tmux) is a terminal multiplexer.
It is currently a deprecated process manager in GameAP,
which is not recommended for use as it does not provide full functionality.

Tmux can be used on older systems, as well as inside containers (LXC, Docker/Podman, etc.), virtual systems,
and systems that do not have Systemd and/or cannot use Docker or Podman.

| Feature                                |    |
|----------------------------------------|----|
| Start, stop, restart servers           | ✅  |
| Statistics                             | ⚠️ |
| Resource limits (CPU / RAM)            | ❌  |
| Console reading                        | ✅  |
| Sending commands to console            | ✅  |
| Isolation                              | ❌  |

#### Tmux Configuration

Tmux does not require additional configuration.

##### Configuration Example

```yaml
process_manager:
  name: tmux
```

## Windows

### Shawl

[Shawl](https://github.com/mtkennerly/shawl) is a lightweight process manager for Windows
that provides basic functionality for managing game servers.
It is written in Rust and uses the Windows API to run applications as Windows services.

| Feature                                |    |
|----------------------------------------|----|
| Start, stop, restart servers           | ✅  |
| Statistics                             | ⚠️ |
| Resource limits (CPU / RAM)            | ❌  |
| Console reading                        | ✅  |
| Sending commands to console            | ❌  |
| Isolation                              | ❌  |

#### Shawl Configuration

Shawl does not require additional configuration.

##### Operation Details

| Parameter | Value |
|-----------|-------|
| Configuration directory | `C:\gameap\services` |
| Stop timeout | 10000 ms |
| Log rotation | Daily |
| Log retention | 7 days |

##### Configuration Example

```yaml
process_manager:
  name: shawl
```

### WinSW

[WinSW](https://github.com/winsw/winsw) (Windows Service Wrapper) is a process manager written in C#.
It allows running applications as Windows services.
In GameAP, it is a deprecated manager and has been replaced by Shawl.

| Feature                                |    |
|----------------------------------------|----|
| Start, stop, restart servers           | ✅  |
| Statistics                             | ⚠️ |
| Resource limits (CPU / RAM)            | ❌  |
| Console reading                        | ✅  |
| Sending commands to console            | ❌  |
| Isolation                              | ❌  |

#### WinSW Configuration

WinSW does not require additional configuration.

##### Operation Details

| Parameter | Value |
|-----------|-------|
| Configuration directory | `C:\gameap\services` |
| Configuration format | XML |

##### Configuration Example

```yaml
process_manager:
  name: winsw
```
