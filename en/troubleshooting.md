---
title: Troubleshooting
layout: default
lang: en
category: Troubleshooting
order: 400
---

This page describes some possible errors and how to fix them.

## Where to Start

For almost any game server problem, details can be found in two places.

**Tasks in the panel.** **Administration** → **GDaemon tasks**, open the latest task — it
contains the command's result and the error output.

**The daemon log** on the dedicated server:

* Linux — `/var/log/gameap-daemon/output.log`
* Windows — `C:\gameap\daemon\logs\output.log`

If the daemon runs under systemd, the same output is available with:

```shell
journalctl -u gameap-daemon -n 200 --no-pager
```

For more detail, set `log_level: debug` in the daemon configuration and restart it.

## Server Startup Errors

### Server status is displayed incorrectly

Sometimes the server starts, but the panel shows it as offline.
See [Server Status Display Errors](#server-status-display-errors).

### Empty server start command

This error occurs when the start command for the game server is empty.

Go to the game server administration page: **Administration** → **Servers** → select the game
server. Or from the main page: **Servers List** → select the server → **Control** →
**Administration**.

Find the **Game server start command** field and enter the command. For Counter-Strike 1.6 it
will be something like this:

```text
./hlds_run -game cstrike +ip {ip} +port {port} +map {default_map} +maxplayers {maxplayers} +sys_ticrate {fps}
```

See [game server configuration](/en/gameap_configure/game_servers.html#run-command) for details.
The default start command can be set
[in the mod settings](/en/gameap_configure/games.html#default-startup-commands).

### The information modal window does not change for a long time

A start usually takes less than 10 seconds. If the progress bar is frozen or the status does not
change for several minutes, check whether the daemon is running and connected to the panel:

```shell
systemctl status gameap-daemon
```

If the daemon is not running:

```shell
systemctl start gameap-daemon
```

If it is running, check the log to see whether the connection to the panel is established. The
messages `gRPC connection failed` and `Reconnecting to panel...` mean the daemon could not reach
the panel — see [The daemon does not connect to the panel](#the-daemon-does-not-connect-to-the-panel).

### Server start task is already exists

This error occurs when a start task has already been created and has not been executed yet.

Go to the **GDaemon tasks** page, find the task in the waiting status, open it with the **View**
button, and click **Cancel**. Then start the server again.

If tasks regularly get stuck waiting, the daemon is not picking them up — check that it is
running and connected to the panel.

## Server Status Display Errors

### Time mismatch

The problem is caused by a time difference between the panel server and the dedicated server.
Set up time synchronization on both.

Changing the time zone on Debian and Ubuntu:

```bash
dpkg-reconfigure tzdata
```

### The daemon is connected, but the status is not updating

The daemon reports server state every `metrics.collection_interval` (5 seconds by default) and
sends a heartbeat every 30 seconds. If the data stays stale longer than that, check the daemon
log for connection drops.

## Server Installation Errors

To find the cause, look at the result of the installation command: **Administration** →
**GDaemon tasks** → find the installation task and open it.

Common causes:

* [No source](#no-source)
* [Incorrectly formed installation archive](#incorrectly-formed-installation-archive)
* [Incorrect installation source](#incorrect-installation-source)

### No source

The panel does not know where to install the game server from. There are several ways:

* via SteamCMD — the Steam APP ID must be set in the game settings;
* from a remote repository — a link to a ZIP or TAR archive must be provided (RAR is not supported);
* from a local repository — a path to a directory with files or to a ZIP or TAR archive. The
  files must be located on the dedicated server where the daemon runs.

The source is set on the **Administration** → **Games** page → select the game → **Edit**.

### Incorrectly formed installation archive

The archive must contain the game server files at its root. The most common mistake is files
sitting in a nested directory.

An incorrectly formed archive for GTA: San Andreas Multiplayer:

![Incorrect installation archive: game server files sit in a nested directory](/images/errors/source_archive_wrong.jpg)

Correct:

![Correct installation archive: game server files sit at the archive root](/images/errors/source_archive_right.jpg)

### Incorrect installation source

For a local repository, check that the directory exists on the dedicated server where the daemon
runs.

For a remote one, check that the link actually starts a file download. Links to Yandex Disk,
Google Drive, and similar storage services are not supported: they serve a page, not a file.

### Failed to install via steamcmd

Sometimes SteamCMD aborts in the middle of a download. The panel can work with the alternative
depot downloader — replace the SteamCMD script:

```shell
cd /srv/gameap/steamcmd
mv steamcmd.sh steamcmd.sh.orig
curl -O https://raw.githubusercontent.com/gameap/steamcmd-depotdownloader/main/steamcmd.sh
chmod 755 steamcmd.sh && chown gameap:gameap steamcmd.sh
```

### Steam account login required

Some games cannot be downloaded anonymously — an account with a purchased copy is needed.
Specify it in the `steam_config` section of the daemon configuration, see
[GameAP Daemon](/en/daemon/daemon.html#steam-account).

Two-factor authentication must be disabled on this account, otherwise the daemon will not be
able to log in.

## GameAP Daemon Errors

### The daemon does not connect to the panel

In GameAP 4 the connection is established by the daemon: it connects to the panel itself over
gRPC on port **31718**. The panel does not connect to the daemon, and no inbound ports on the
dedicated server are needed for this.

The telltale sign in the daemon log is repeated `gRPC connection failed` and
`Reconnecting to panel...` messages.

What to check:

1. **Port reachability** from the dedicated server:

   ```shell
   nc -zv panel.example.com 31718
   ```

   If the connection cannot be established, the port is blocked by a firewall or the panel is
   listening on a different address.

2. **The panel address** in the daemon configuration — the `grpc.address` parameter. It must
   point to an address reachable from the dedicated server.

3. **Certificate verification error.** If the log contains a message about a certificate name
   mismatch, the daemon is connecting via an address that is not in the panel's certificate.
   Set `GRPC_EXTERNAL_HOST` on the panel and reissue the certificate — see
   [GRPC API](/en/daemon/grpc.html).

4. **`registration failed`.** The connection is established, but the panel rejected the
   registration: a wrong `ds_id` or `api_key` in the daemon configuration. The easiest fix is
   to register the daemon again.

### Console or file manager does not work

If servers start and stop but the console and file manager do not work, the problem is almost
always the connection to the panel: control commands may have completed before a disconnect,
while the console and files require a live data stream.

Check the connection as described in the section above; restart the daemon if needed:

```shell
systemctl restart gameap-daemon
```

### A file cannot be uploaded via the file manager

By default, uploading archives and arbitrary binary files is not allowed — the type is detected
from the file content, not the extension. They can be allowed with the
`FILES_UPLOAD_ALLOW_ARCHIVES` and `FILES_UPLOAD_ALLOW_BINARY` variables, see
[Security](/en/security.html).

The maximum file size for a regular upload is 100 MB.

### The daemon does not start after a server reboot

Enable autostart:

```shell
systemctl enable gameap-daemon
```

### Task complete with an error

A generic error during game server installation, start, restart, or stop. There can be many
causes.

Open **Administration** → **GDaemon tasks**, find the latest task, and look at the details.
Then check the daemon log.

#### No source to install game

No installation source is set for the game. See [No source](#no-source).

## Panel Login Errors

### An administrator cannot log in after an upgrade

In GameAP 4, two-factor authentication is mandatory for administrators: 30 days after the first
reminder, logging in stops issuing a full session. This is not an account lockout — enable 2FA
and log in again.

The steps, including what to do when the device with the codes is lost, are on the
[Security](/en/security.html) page.

### Too many login attempts

After 20 failed attempts from one address or 5 attempts for one login, logging in is blocked
for 15 minutes with a `429` response. The limit clears itself when the time passes.

### A password is rejected when changing it

The minimum password length is 12 characters, and the password is checked against the
compromised password list. See [Security](/en/security.html).

## Finding the Cause

If your problem is not listed here, look for details in the logs.

**The daemon log** — `/var/log/gameap-daemon/output.log` on Linux,
`C:\gameap\daemon\logs\output.log` on Windows. For verbose output, set `log_level: debug` in
the daemon configuration and restart it.

**The panel log** — the output of the `gameap` process. Under systemd:

```shell
journalctl -u gameap -n 200 --no-pager
```

The verbosity level is set with the `LOGGER_LEVEL` variable, see the
[config.env Reference](/en/config.html).

**Tasks** — in the panel, **Administration** → **GDaemon tasks**: the result of every start,
installation, and restart command is there.
