---
title: Process Managers
layout: default
lang: en
category: GameAP Daemon
order: 410
---

A process manager is the component GameAP Daemon uses to run game servers on a node.
It is responsible for starting, stopping, and restarting game servers, monitoring their status, collecting statistics, and limiting resources (CPU / RAM).

By default the daemon picks the manager itself. On Linux it uses `systemd`; when the daemon itself runs inside a
container, or `systemctl` is missing, or systemd is not PID 1, it falls back to `tmux`, and to `simple` when tmux is not
installed either. On Windows the default is [Shawl](https://github.com/mtkennerly/shawl), on macOS — `tmux`.

## Configuration

The process manager is configured in the `process_manager` block of the
[GameAP Daemon configuration file](/en/daemon/daemon.html). When a node is added through the panel, the manager can be
chosen in advance: expand **Additional settings** and pick it in the **Process Manager** field. The Linux tab offers
**Auto (Recommended)**, **Simple (Custom scripts controlled)**, SystemD, Docker, Podman and Tmux; the Windows tab offers
only **Auto (Recommended)**, **Simple (Custom scripts controlled)**, WinSW and Shawl, so Docker on Windows has to be set
in the configuration file by hand.

### Basic Structure

```yaml
process_manager:
  name: <manager_name>
  config:
    <key>: <value>
```

**Available managers:**
- Linux: `systemd` (default), `docker`, `podman`, `tmux`, `simple`
- macOS: `tmux` (default), `docker`, `podman`, `simple`
- Windows: `shawl` (default), `winsw`, `simple`, `docker` (Podman is not supported on Windows)

The contents of `config` depend on the manager:

| Manager                            | Keys in `config`                                                                                                 |
|------------------------------------|------------------------------------------------------------------------------------------------------------------|
| `systemd`                          | `scope` — `system` or `user`, see [Systemd Configuration](#systemd-configuration)                                |
| `docker`                           | `host`, `cert_path`, `api_version`, plus any [container parameter](#container-parameters) as a node-wide default |
| `podman`                           | `socket_path`, plus any [container parameter](#container-parameters) as a node-wide default                      |
| `tmux`, `simple`, `shawl`, `winsw` | none                                                                                                             |

> `scope` is accepted only by `systemd`. With any other manager the daemon refuses to start with the error
> `process_manager.config.scope is only valid for process_manager.name=systemd`.

### Restart After a Crash

Who brings a crashed game server back differs between managers. With `systemd` it is the unit itself (`Restart=always`),
with Shawl — the service (`--restart`), with WinSW — the service's `onfailure` actions. With `tmux`, `simple`, Docker and
Podman the daemon restarts the server itself: every 5 seconds it checks the state of all servers and starts the ones
that are down. Servers that have been idle for a long time are checked less often, so such a restart can be delayed by
up to about two minutes. In every case the server's **Autostart on crash** setting is respected: a server that was
stopped on purpose is never started again.

### Statistics

`systemd`, Docker and Podman report the full set of metrics: CPU, memory usage and limit, network traffic, disk I/O and
process count. `tmux`, `simple`, Shawl and WinSW report only whether the server is running. Metric collection is described
in [GameAP Daemon](/en/daemon/daemon.html#metrics-collection).

## Linux and macOS

### Systemd

Platforms: Linux.

Used by default on Linux. It is a modern process manager that provides high performance and reliability.
Isolation capabilities for this process manager are limited.

| Feature                      |                               |
|------------------------------|-------------------------------|
| Start, stop, restart servers | ✅                            |
| Restart after a crash        | ✅ The unit, `Restart=always` |
| Statistics                   | ✅                            |
| Resource limits (CPU / RAM)  | ✅                            |
| Console reading              | ✅                            |
| Sending commands to console  | ✅                            |
| Isolation                    | ⚠️ Limited                    |

> Units generated before the daemon started writing the `CPUAccounting=yes`, `MemoryAccounting=yes`, `IOAccounting=yes`,
> `IPAccounting=yes` and `TasksAccounting=yes` directives report zero statistics until the server is restarted and the
> unit file is regenerated.

#### Systemd Configuration

Systemd does not require additional configuration. The Daemon automatically creates and manages unit files.
The only parameter is `scope`:

| Config Key | Description                                                | Default  |
|------------|------------------------------------------------------------|----------|
| `scope`    | Where the units live and who runs them: `system` or `user` | `system` |

* `system` — units are written to `/etc/systemd/system`; the daemon needs the rights to manage system units (normally it
  runs as root). Every unit gets `User=` / `Group=` from the game server's user (the **Su User** field), so servers on one
  node can run under different users.
* `user` — units are written to `~/.config/systemd/user/` and managed with `systemctl --user`. The daemon runs as a regular
  user and **all game servers run as that same user**. Lingering must be enabled for that user
  (`sudo loginctl enable-linger <user>`), otherwise the daemon warns that user services will be killed at logout;
  `XDG_RUNTIME_DIR` (usually `/run/user/<uid>`) must be accessible. A server whose **Su User** field names a different
  user fails to start with `server requests user "..." but daemon runs as "..."`.

##### Configuration Examples

Default (system scope):
```yaml
process_manager:
  name: systemd
```

User scope:
```yaml
process_manager:
  name: systemd
  config:
    scope: user
```

### Docker

Platforms: Linux, macOS, Windows.

Docker runs every game server in its own container, isolating it from the node and from the other servers.

> Docker or Podman is **required** for games imported from Pelican or Pterodactyl eggs. An imported egg carries no GameAP
> installation rules: its files are installed by a script that runs in the egg's own container image
> (`docker_installation_image` + `docker_installation_script`). Under `systemd`, `tmux`, `simple`, `shawl` or `winsw`
> such a server cannot be installed — the installation fails with
> `could not determine the rules for installing the game`.
> See [Games Import](/en/gameap_configure/games_import.html#pelican-and-pterodactyl-features).

| Feature                      |                  |
|------------------------------|------------------|
| Start, stop, restart servers | ✅               |
| Restart after a crash        | ✅ GameAP Daemon |
| Statistics                   | ✅               |
| Resource limits (CPU / RAM)  | ✅               |
| Console reading              | ✅               |
| Sending commands to console  | ✅               |
| Isolation                    | ✅               |

#### Docker Configuration

##### Connection Parameters

| Config Key    | Description                                                                                      | Default |
|---------------|--------------------------------------------------------------------------------------------------|---------|
| `host`        | Address of the Docker daemon, for example `unix:///var/run/docker.sock` or `tcp://10.0.0.5:2376` | —       |
| `cert_path`   | Directory with `ca.pem`, `cert.pem` and `key.pem` for a TLS connection                           | —       |
| `api_version` | Docker API version to use                                                                        | —       |

If none of the three keys is set, the connection is configured from the standard Docker environment variables
(`DOCKER_HOST`, `DOCKER_CERT_PATH`, `DOCKER_TLS_VERIFY`, `DOCKER_API_VERSION`) — with none of them set, the local
Docker daemon is used.

Any [container parameter](#container-parameters) can also be placed in `config` as a default for the whole node, with
or without the `docker_` prefix (`image` and `docker_image` are equivalent). Do not set `container_name` node-wide:
every server on the node would get the same container name.

##### Configuration Examples

Minimal:
```yaml
process_manager:
  name: docker
```

Remote Docker daemon over TLS:
```yaml
process_manager:
  name: docker
  config:
    host: tcp://10.0.0.5:2376
    cert_path: /etc/gameap-daemon/docker-certs
```

Node-wide container defaults:
```yaml
process_manager:
  name: docker
  config:
    image: docker.io/gameap/debian:latest
    workdir: /server
    dns: 8.8.8.8,1.1.1.1
```

#### Container Parameters

Docker and Podman read the container settings from the same keys (Podman keeps the `docker_` prefix). The keys are
entered in the panel as metadata: in **Administration → Games**, both the game form and the mod form have a
**Metadata** tab. The **?** button next to the key field opens the **Metadata keys** reference with a description and an
example for every key and inserts the chosen key into the field.

| Key                              | Description                                                                                                                                        | Default                      |
|----------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------|
| `docker_image`                   | Image the game server runs in. `:latest` is added when the tag is omitted                                                                          | `docker.io/gameap/debian`    |
| `docker_container_name`          | Container name                                                                                                                                     | server XID                   |
| `docker_workdir`                 | Working directory inside the container; the server directory is mounted there                                                                      | `/server`                    |
| `docker_network_mode`            | Network mode. Only the value `host` has an effect and enables host networking                                                                      | bridge                       |
| `docker_dns`                     | DNS servers for the container, comma-separated. IP addresses only                                                                                  | —                            |
| `docker_volumes`                 | Additional mounts in the form `source:target[:ro]`, comma-separated or as a JSON array. A relative source is resolved against the server directory | —                            |
| `docker_capabilities`            | Additional Linux capabilities, comma-separated without spaces                                                                                      | —                            |
| `docker_privileged`              | Privileged mode. Only the exact string `true` enables it                                                                                           | —                            |
| `docker_installation_image`      | Image for the installation step. Works only together with `docker_installation_script`                                                             | —                            |
| `docker_installation_script`     | Installation script body. Runs inside the installation container in `/mnt/server`, where the server directory is mounted                           | —                            |
| `docker_installation_entrypoint` | Interpreter for the installation script                                                                                                            | shebang, otherwise `/bin/sh` |
| `docker_installation_user`       | User the installation runs as                                                                                                                      | `root`                       |

**Resolution order.** The daemon looks a key up in this order and takes the first non-empty string value:

1. game server variables — the mod's variable defaults, then the values set for the server itself;
2. game mod metadata;
3. game metadata;
4. `process_manager.config` in the daemon configuration file, where the key may be written with or without the
   `docker_` prefix.

Only string values are used; a number or a boolean in metadata is ignored.

> The **Metadata** block of the game server itself is not read by GameAP Daemon 4.1.2, even though the panel's reference
> offers the container keys there. To override a key for a single server, add a mod variable with the same name (for
> example `docker_image`) and set its value in the server's **Settings**.

> Games imported from Pelican or Pterodactyl eggs get the following keys in the metadata of the **Default** mod:
> `docker_image`, `docker_installation_image`, `docker_installation_script`, `docker_installation_entrypoint`,
> `docker_installation_user: root` and `docker_workdir: /home/container`. The import also stores `docker_startup_done`
> (the egg's startup-done marker) in the mod metadata and `pelican_egg` (the original egg JSON) in both the game and the
> mod metadata. Neither of these two keys is read by GameAP Daemon; they are kept for reference.

#### How It Works

Applies to Docker and Podman alike unless stated otherwise.

* The container is named after the server XID (`docker_container_name` overrides it). Containers created by older
  daemon versions and named after the server UUID are still found.
* The server directory is bind-mounted into the container at `docker_workdir` (default `/server`), which is also the
  working directory. Only this directory and the additional `docker_volumes` mounts survive a restart: the container
  itself is removed on every stop.
* The game process runs as the uid:gid of the game server's user (the **Su User** field). That user must exist on the
  node even though the server runs in a container, otherwise the start fails with `failed to lookup user`. When the
  field is empty, the daemon's own user is used.
* Ports are published on the server's IP address: the connect port over TCP and UDP, the query port over UDP and the
  RCON port over TCP when they differ from the connect port.
* CPU and RAM limits are taken from the server's
  [resource limits](/en/gameap_configure/game_servers.html#resource-limits) in the panel (CPU in millicores, RAM in
  bytes). There are no metadata keys for limits.
* Neither Docker nor Podman restarts the container on its own: the Docker container is created with the restart policy
  disabled, and the Podman container is created without a restart policy at all (Podman's own default is `no`). A
  crashed container is restarted by the daemon according to the server's **Autostart on crash** setting.
* Stopping waits up to 30 seconds for the process to exit, then the container is removed.
* The console shows the last 500 lines of the container log. A command from the console is delivered differently:
  Docker attaches to the container's stdin and writes the command there; Podman runs the command as an `exec` inside the
  container, so it does not reach the game process's stdin.

#### Installation in a Container

If both `docker_installation_image` and `docker_installation_script` are set, the installation runs in a separate
container:

1. the installation image is pulled;
2. the script is written to `.gameap_install.sh` in the server directory (the file is removed afterwards);
3. a temporary container `gameap-install-<XID>` is created with the server directory mounted at `/mnt/server` and the
   server's environment variables;
4. the script is run by the interpreter from `docker_installation_entrypoint` (a bare name such as `ash` is looked up
   in `/bin/`), otherwise by the interpreter from the script's shebang, otherwise by `/bin/sh`;
5. the log is streamed live with Docker; with Podman it is collected after the container exits;
6. after a successful run the files are chowned to the server's user, and the container is removed.

The installation runs as `root` unless `docker_installation_user` says otherwise. A non-zero exit code fails the
installation.

The container step runs after the regular installation from the game's repositories or Steam, if the game has any. If
the game has neither installation rules nor both container keys, the installation fails with
`could not determine the rules for installing the game`. If only one of the two keys is set, the container step is
skipped and the daemon only pulls the runtime image.

### Podman

Platforms: Linux, macOS.

Podman is an alternative to Docker that provides isolation of game servers in containers. It reads the same
[container parameters](#container-parameters) as Docker (the keys keep the `docker_` prefix) and behaves the same way at
run time, but connects differently: it talks to the libpod REST API (v4.0.0) over a Unix socket, and of the connection
parameters only `socket_path` is used — `host`, `cert_path` and `api_version` are ignored.

| Feature                      |                                      |
|------------------------------|--------------------------------------|
| Start, stop, restart servers | ✅                                   |
| Restart after a crash        | ✅ GameAP Daemon                     |
| Statistics                   | ✅                                   |
| Resource limits (CPU / RAM)  | ✅                                   |
| Console reading              | ✅                                   |
| Sending commands to console  | ⚠️ See [How It Works](#how-it-works) |
| Isolation                    | ✅                                   |

#### Podman Configuration

The Podman API socket must be running. Rootless:

```bash
systemctl --user enable --now podman.socket
```

Root:

```bash
sudo systemctl enable --now podman.socket
```

##### Specific Parameters

| Config Key    | Description                    | Default                                                                                                                                                |
|---------------|--------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| `socket_path` | Path to the Podman Unix socket | `unix:///run/user/<uid>/podman/podman.sock` when the daemon runs as a non-root user and that socket exists, otherwise `unix:///run/podman/podman.sock` |

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

Platforms: Linux, macOS.

[Tmux](https://github.com/tmux/tmux) is a terminal multiplexer.
It is currently a deprecated process manager in GameAP,
which is not recommended for use as it does not provide full functionality.

Tmux can be used on older systems, as well as inside containers (LXC, Docker/Podman, etc.), virtual systems,
and systems that do not have Systemd and/or cannot use Docker or Podman.

| Feature                      |                                       |
|------------------------------|---------------------------------------|
| Start, stop, restart servers | ✅                                    |
| Restart after a crash        | ✅ GameAP Daemon                      |
| Statistics                   | ⚠️ Only whether the server is running |
| Resource limits (CPU / RAM)  | ❌                                    |
| Console reading              | ✅                                    |
| Sending commands to console  | ✅                                    |
| Isolation                    | ❌                                    |

#### Tmux Configuration

Tmux does not require additional configuration.

##### Configuration Example

```yaml
process_manager:
  name: tmux
```

### Simple

Platforms: Linux, macOS, Windows. In the panel it is listed as **Simple (Custom scripts controlled)**.

`simple` has no built-in way of talking to a game server: it only runs the commands from the `scripts` block of the
daemon configuration, substituting the server's start, stop or restart command for `{command}`. The status and
console-reading scripts are called without a server command, so with the default `{command}` template they expand to
nothing and every such call fails with `empty command`. Without a working status script the server's state cannot be
determined: it is never reported to the panel and is not restarted after a crash.

On Linux the daemon falls back to `simple` when neither systemd nor tmux is available.

| Feature                      |                                            |
|------------------------------|--------------------------------------------|
| Start, stop, restart servers | ✅ Via the configured scripts              |
| Restart after a crash        | ✅ GameAP Daemon, requires a status script |
| Statistics                   | ⚠️ Only whether the server is running      |
| Resource limits (CPU / RAM)  | ❌                                         |
| Console reading              | ⚠️ Requires a console-reading script       |
| Sending commands to console  | ⚠️ Requires a command script               |
| Isolation                    | ❌                                         |

#### Simple Configuration

```yaml
process_manager:
  name: simple
```

## Windows

### Shawl

Platforms: Windows.

[Shawl](https://github.com/mtkennerly/shawl) is a lightweight process manager for Windows
that provides basic functionality for managing game servers.
It is written in Rust and uses the Windows API to run applications as Windows services.

| Feature                      |                                       |
|------------------------------|---------------------------------------|
| Start, stop, restart servers | ✅                                    |
| Restart after a crash        | ✅ Shawl, `--restart`                 |
| Statistics                   | ⚠️ Only whether the server is running |
| Resource limits (CPU / RAM)  | ❌                                    |
| Console reading              | ✅                                    |
| Sending commands to console  | ❌                                    |
| Isolation                    | ❌                                    |

#### Shawl Configuration

Shawl does not require additional configuration.

##### Operation Details

| Parameter               | Value                |
|-------------------------|----------------------|
| Configuration directory | `C:\gameap\services` |
| Stop timeout            | 10000 ms             |
| Log rotation            | Daily                |
| Log retention           | 7 days               |

##### Configuration Example

```yaml
process_manager:
  name: shawl
```

### WinSW

Platforms: Windows.

[WinSW](https://github.com/winsw/winsw) (Windows Service Wrapper) is a process manager written in C#.
It allows running applications as Windows services.
In GameAP, it is a deprecated manager and has been replaced by Shawl.

| Feature                      |                                       |
|------------------------------|---------------------------------------|
| Start, stop, restart servers | ✅                                    |
| Restart after a crash        | ✅ WinSW, `onfailure` actions         |
| Statistics                   | ⚠️ Only whether the server is running |
| Resource limits (CPU / RAM)  | ❌                                    |
| Console reading              | ✅                                    |
| Sending commands to console  | ❌                                    |
| Isolation                    | ❌                                    |

#### WinSW Configuration

WinSW does not require additional configuration.

##### Operation Details

| Parameter               | Value                |
|-------------------------|----------------------|
| Configuration directory | `C:\gameap\services` |
| Configuration format    | XML                  |

##### Configuration Example

```yaml
process_manager:
  name: winsw
```

## Troubleshooting

### Docker

**`failed to connect to docker daemon`** — the daemon cannot reach Docker. Check that Docker is running
(`systemctl status docker`) and that the user the daemon runs as may access `/var/run/docker.sock` — usually by being
a member of the `docker` group (`sudo usermod -aG docker <user>`, then log in again). With a remote Docker daemon check
`host` and `cert_path`.

**`failed to lookup user <name>`** — the game server's **Su User** does not exist on the node. Create the OS user or
change the field.

### Podman

**Connection errors mentioning the socket path** — the Podman API socket is not running or lives elsewhere. Rootless:

```bash
systemctl --user enable --now podman.socket
ls -la /run/user/$(id -u)/podman/podman.sock
curl --unix-socket /run/user/$(id -u)/podman/podman.sock http://d/v4.0.0/libpod/info
```

If the daemon runs as root, the socket is `/run/podman/podman.sock` (`sudo systemctl enable --now podman.socket`).
When the socket is at a non-standard path, set `socket_path`.

### Docker and Podman

**`could not determine the rules for installing the game`** when installing a server from an imported Pelican or
Pterodactyl egg — the node is running a process manager other than Docker or Podman. Switch the manager in the
daemon configuration and restart the daemon.
