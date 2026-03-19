---
title: GameAP Daemon
layout: default
lang: en
category: GameAP Daemon
order: 100
---

GameAP Daemon is a background application that exchanges data with the panel and works with game servers
(installs, removes, starts, stops, etc.).

It is the main application that controls the status of game servers, restarts them when necessary or
on demand.

![](/images/en/gameap_architecture.svg)

## Installation

### Automatic Installation via Panel

In the control panel, go to **Administration** → **Dedicated Servers** → **Create**.
A window will appear on the page offering automatic installation. Copy the code and execute it on
the dedicated server.

![](/images/en/daemon/autoinstall.png)

#### Installation Parameters

During installation, you can specify some parameters. To do this, open "Advanced Settings" in the auto-install window.

##### Process Manager

A process manager is a system utility that manages processes on the server.
It is responsible for starting, stopping, and restarting game servers, as well as monitoring their status.
By default, systemd is used on Linux, and [Shawl](https://github.com/mtkennerly/shawl) on Windows.
See details on the [Process Manager](/en/daemon/process_managers.html) page.

## Configuration

The GameAP Daemon configuration is located in the file `/etc/gameap-daemon/gameap-daemon.yaml` (Linux)
or `C:\gameap\daemon\daemon.yaml` (Windows).

### Basic Parameters

| Parameter                 | Required              | Type      | Information
|---------------------------|-----------------------|-----------|------------
| ds_id                     | yes                   | integer   | Dedicated server ID
| listen_port               | no (default 31717)    | integer   | Port to listen on
| api_host                  | yes                   | string    | API host
| api_key                   | yes                   | string    | API key


### SSL/TLS

| Parameter                 | Required              | Type      | Information
|---------------------------|-----------------------|-----------|------------
| ca_certificate_file       | yes                   | string    | CA certificate
| certificate_chain_file    | yes                   | string    | Server certificate
| private_key_file          | yes                   | string    | Server private key
| private_key_password      | no                    | string    | Server key password
| dh_file                   | yes                   | string    | Diffie-Hellman certificate

[Creating certificates](https://github.com/gameap/GDaemon2#creating-certificates)

### Login and Password Authentication

| Parameter                 | Required              | Type      | Information
|---------------------------|-----------------------|-----------|------------
| password_authentication   | no                    | boolean   | Enable login and password authentication
| daemon_login              | no                    | string    | Login. On Linux, if empty or not set, PAM will be used.
| daemon_password           | no                    | string    | Password. On Linux, if empty or not set, PAM will be used.

### Using a Non-Anonymous Steam Account

Many servers cannot be downloaded via steamcmd without logging into an account with a purchased copy of the game.
In this case, you need to provide the daemon with the login and password of such a Steam account for further authorization in steamcmd.

> Note! Two-factor authentication must be disabled on your account, otherwise the daemon will not be able to authorize in steamcmd.

[Example yaml file](https://github.com/gameap/daemon/blob/master/config/gameap-daemon.yaml)

After this, you can add the Steam account configuration:
```
steam_config:
    login: "your login"
    password: "Pa$$worD"
```
