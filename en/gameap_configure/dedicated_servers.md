---
title: Dedicated Servers
layout: default
lang: en
category: Panel settings
order: 300
---

## New dedicated server

When getting started with the panel, you should add a dedicated server (VDS/VPS, container, physical
server) to run GameAP Daemon on.

Go to the **"Administration"** → **"Dedicated Servers"** page and click the **"Create"** button.
A window opens with a ready-made installation command — separate ones for Linux and for Windows.

Certificates are created and passed to the daemon by the panel itself: there is no need to generate
and sign them manually.

> The setup key is valid for **1 hour** and is single-use — it is revoked once the daemon has
> successfully registered. The key is shared across the whole panel, so reopening the "Create" window
> issues a new key, and a previously copied command stops working.

### What you need on the dedicated server

* Superuser rights (root on Linux, administrator on Windows).
* Outgoing access to the panel on port **31718/TCP** — the daemon uses it to talk to the panel over
  gRPC. This is the only panel port a running daemon needs.
* Outgoing access to `github.com` and `api.github.com` — gameapctl and the daemon itself are
  downloaded from there.
* For the Linux script: `curl`, `tar` and `install` installed.
* Architecture: `amd64`, `arm64`, `386` or `arm`.

Port **8025** is the panel web interface and API. It is needed by the administrator's browser and by
the `curl` command in the installation one-liner, but a running daemon does not require it.

Port **31717** is what the daemon reports as its own during registration. In GameAP 4 the panel does
not establish incoming connections to the daemon, so there is no need to open this port on the
dedicated server.

### Installation on Linux

Copy the command from the Linux tab and run it on the dedicated server **as root**:

```bash
bash <(curl -fsSL 'https://your-panel/nodes/setup/zItWHWlI4RKPl9ZsYc3y3WgdKq7mNvBx')
```

The script checks the environment, installs `gameapctl` into `/usr/local/bin` (or updates it if it is
already installed), and then runs `gameapctl daemon install`. That installs GameAP Daemon, creates the
`gameap` user, installs SteamCMD, registers the daemon in the panel and starts it as the
`gameap-daemon` systemd service.

The command must be run as root, not through `sudo`: the `bash <(...)` process substitution does not
survive `sudo`, and the script will stop with an error. If you are not working as root, download the
script to a file and run it:

```bash
curl -fsSL 'https://your-panel/nodes/setup/zItWHWlI4RKPl9ZsYc3y3WgdKq7mNvBx' -o gameap-setup.sh
sudo bash gameap-setup.sh
```

Installation paths:

| Item              | Path                                    |
|-------------------|-----------------------------------------|
| Daemon            | `/usr/bin/gameap-daemon`                |
| Configuration     | `/etc/gameap-daemon/gameap-daemon.yaml` |
| Certificates      | `/etc/gameap-daemon/certs`              |
| Working directory | `/srv/gameap`                           |
| SteamCMD          | `/srv/gameap/steamcmd`                  |
| Logs              | `/var/log/gameap-daemon/output.log`     |

### Installation on Windows

There is no PowerShell one-liner, installation is done through gameapctl:

1. Download the gameapctl archive for your architecture from the
   [gameapctl releases](https://github.com/gameap/gameapctl/releases) page. The archive is named
   `gameapctl-<version>-windows-amd64.zip`.
2. Unpack the archive and run `gameapctl.exe`.
3. In the **GameAP Daemon** section click **Install**.
4. Paste the string from the Windows tab in the panel into the **Connect URL** field.
5. Click **Install**.

The same can be done with a command in the console:

```shell
gameapctl daemon install --connect=grpc://your-panel:31718/zItWHWlI4RKPl9ZsYc3y3WgdKq7mNvBx
```

Installation paths:

| Item              | Path                                     |
|-------------------|------------------------------------------|
| Daemon            | `C:\gameap\daemon\gameap-daemon.exe`     |
| Configuration     | `C:\gameap\daemon\gameap-daemon.yaml`    |
| Certificates      | `C:\gameap\daemon\certs`                 |
| Working directory | `C:\gameap`                              |
| SteamCMD          | `C:\gameap\steamcmd`                     |
| Logs              | `C:\gameap\daemon\logs\output.log`       |

The daemon is registered as the **GameAP Daemon** service.

### Advanced installation settings

The creation window has a collapsible **"Advanced Settings"** block:

* **Process manager** — what should manage game server processes. On Linux `systemd`, `docker`,
  `podman`, `tmux` and `simple` are available, on Windows — `winsw`, `shawl` and `simple`. Selected
  automatically by default. See details on the [Process Managers](/en/daemon/process_managers.html) page.
* **GitHub** — build the daemon from source instead of using a ready-made release.
* **Branch** — repository branch, if the daemon is built from source.

The panel appends the selected values to the installation command:

```bash
bash <(curl -s '...') --config='process_manager.name=docker' --github --branch=master
```

### Manual installation

If the automatic script does not fit — a non-standard distribution, your own package installation
rules, installation into a prepared image — the daemon can be registered manually.

1. Download the `gameap-daemon` binary for your platform from the
   [daemon releases](https://github.com/gameap/daemon/releases) page and place it on the dedicated server.
2. Open the **"Create"** window in the panel and take the connect URL of the form
   `grpc://host:port/key` from the Windows command.
3. Run the registration:

```bash
gameap-daemon enroll --connect=grpc://your-panel:31718/zItWHWlI4RKPl9ZsYc3y3WgdKq7mNvBx
```

The daemon connects to the panel, creates a dedicated server record in it, receives the certificates
(`ca.crt`, `server.crt`, `server.key` with `0600` permissions in a directory with `0700` permissions)
and writes the configuration file.

| Flag              | Linux default                           | Windows default                       | Purpose                                                                    |
|-------------------|-----------------------------------------|---------------------------------------|----------------------------------------------------------------------------|
| `--connect`       | —                                       | —                                     | Connect URL. Required                                                      |
| `--config-path`   | `/etc/gameap-daemon/gameap-daemon.yaml` | `C:\gameap\daemon\gameap-daemon.yaml` | Where to write the configuration                                           |
| `--certs-dir`     | `/etc/gameap-daemon/certs`              | `C:\gameap\daemon\certs`              | Where to save the certificates                                             |
| `--work-path`     | `/srv/gameap`                           | `C:\gameap`                           | Game servers working directory                                             |
| `--steamcmd-path` | `/srv/gameap/steamcmd`                  | `C:\gameap\steamcmd`                  | SteamCMD directory                                                         |
| `--listen-ip`     | `0.0.0.0`                               | `0.0.0.0`                             | IP the daemon reports about itself. With `0.0.0.0` it is detected automatically |
| `--listen-port`   | `31717`                                 | `31717`                               | Port the daemon reports about itself                                       |

> The `enroll` command overwrites the configuration file **completely** — without merging with
> existing settings and without a backup. If the daemon has already been configured, save the
> configuration beforehand.

After registering the daemon you should register it as a system service and start it yourself.

### If the installation fails

| Message                                         | Cause                                                                                                                                        |
|-------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| `This script must be run as root.`              | The script was not run as root. Download it to a file and run it through `sudo`, as shown above                                               |
| `Error: 'curl' is required but not installed.`  | The server does not have `curl`, `tar` or `install`. Install the missing package and try again                                                |
| `Unsupported architecture: ...`                 | There are no ready-made builds for the server architecture                                                                                    |
| `Failed to detect latest gameapctl version`     | No access to `api.github.com`, or the GitHub API limit is exhausted (60 requests per hour from one address). Wait or install gameapctl manually |
| `cannot reach gRPC server at ...`               | Port 31718 of the panel is not reachable from the dedicated server. Check the firewall and the panel address — see [GRPC API](/en/daemon/grpc.html) |
| The link returns `403`                          | The setup key has expired or has been reissued. Open the "Create" window again and copy the new command                                       |

## Editing dedicated servers

To edit a dedicated server (node), go to **"Administration"** → **"Dedicated Servers"** page, then 
select the dedicated server you want to edit and click on the **"Edit"** button.

### Description of parameters

#### Basic

##### Name

Dedicated server name. It can take any non-empty value, it does not affect any features.

##### Working directory

This directory contains the basic scripts for managing the processes of game servers. Subdirectories
of the working directory contain the files of game servers. For the specified path,
[game server directory](/en/gameap_configure/game_servers.html#directory) is assigned. The default is `/srv/gameap`.

##### Path to SteamCMD

Path to the SteamCMD directory (the `steamcmd.sh` script is located there). The default is `/srv/gameap/steamcmd`.

##### IP list

List of IP or hosts where game servers will run.

#### Scripts

Command templates the panel uses to manage game servers on this dedicated server. Filling them in is
optional: if a field is empty, the default command is used.

#### GameAP Daemon

Daemon connection details. They are filled in automatically when the dedicated server is registered,
and usually do not need to be changed manually.

In GameAP 4 the connection is established by the daemon: it connects to the panel over gRPC itself and
keeps a persistent connection. The panel does not connect to the daemon.

##### GameAP Daemon host and port

The address and port the daemon reported about itself during registration. The default port is `31717`.
The values are informational — the panel does not use them to connect, and there is no need to open
this port on the dedicated server.

##### GameAP Daemon login and password

Legacy fields from GameAP 3, where the panel connected to the daemon itself. They are not used in
GameAP 4 and do not need to be filled in.
