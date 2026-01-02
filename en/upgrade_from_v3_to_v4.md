---
title: Upgrade from v3 to v4
layout: default
lang: en
category: Install GameAP
order: 191
---

## Using gameapctl to Automatically Upgrade

If you already have the old panel installed, you can upgrade to the new version using [`gameapctl`](https://github.com/gameap/gameapctl).
This utility is installed automatically with GameAP and allows you to manage the panel environment, including updates.

Run the following commands:

```shell
gameapctl self-update
gameapctl panel upgrade --to=v4
```

### Don't have gameapctl installed? (gameapctl: command not found) 

If you don't have `gameapctl` installed, you can download it from [GitHub releases](https://github.com/gameap/gameapctl/tags).

Or use curl to download for Linux amd64:

```shell
curl -OL https://github.com/gameap/gameapctl/releases/download/v0.21.1/gameapctl-v0.21.1-linux-amd64.tar.gz
```

Unpack to `/usr/local/bin`:
```shell
tar xvfz gameapctl-v0.21.1-linux-amd64.tar.gz -C /usr/local/bin
```

Then run the self-update command:

```shell
gameapctl self-update
```

## Running Alongside the Old Version

If you want to run GameAP v3 and GameAP v4 side by side for testing purposes before upgrading, follow these steps:

Download GameAP from [https://github.com/gameap/gameap/releases](https://github.com/gameap/gameap/releases) for your platform.

Or use curl to download for Linux amd64:

```shell
curl -OL https://github.com/gameap/gameap/releases/download/v4.0.0/gameap-v4.0.0-linux-amd64.tar.gz
```

Unpack to `/usr/bin/gameap`:

```shell
tar xvfz gameap-v4.0.0-linux-amd64.tar.gz -C /usr/bin
```

Run GameAP:

```shell
gameap --legacy-env /var/www/gameap/.env
```

GameAP v4 will run on port 8025.

### Systemd Service Configuration

You can create a separate systemd service for GameAP v4 to manage it easily.

First, create the `gameap` user and group:

```shell
useradd -r -s /usr/sbin/nologin -d /var/lib/gameap gameap
```

Create the working directory:
```shell
mkdir -p /var/lib/gameap
```

Set ownership:
```shell
chown gameap:gameap /var/lib/gameap
```

Then create the file `/etc/systemd/system/gameap.service`:

```ini
[Unit]
Description=GameAP - Game Server Control Panel
Documentation=https://docs.gameap.com
After=network.target
Wants=network-online.target
Requires=network.target

[Service]
Type=simple
User=gameap
Group=gameap

# Working directory
WorkingDirectory=/var/lib/gameap

ExecStart=/usr/bin/gameap

# Allow binding to privileged ports
AmbientCapabilities=CAP_NET_BIND_SERVICE

# Graceful stop
ExecStop=/bin/kill -TERM $MAINPID
KillMode=mixed
KillSignal=SIGTERM
TimeoutStopSec=30

# Restart policy
Restart=always
RestartSec=5
StartLimitInterval=60
StartLimitBurst=3

# Environment configuration
# EnvironmentFile=/etc/gameap/config.env

RuntimeDirectory=gameap
PIDFile=/run/gameap/gameap.pid

# Filesystem permissions
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true

ReadWritePaths=/var/lib/gameap /var/lib/gameap/files
# This is needed if you run GameAP with legacy environment
# ReadWritePaths=/var/www/gameap

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=gameap

[Install]
WantedBy=multi-user.target
```

Then enable and start the service:

```shell
systemctl daemon-reload
systemctl enable gameap
systemctl start gameap
```