---
title: Dedicated Servers
layout: default
lang: en
category: Panel settings
order: 300
---

* This will become a table of contents (this text will be scraped).
{:toc}

## New dedicated server

When getting started working with the panel, you should add a dedicated server (VDS/VPS, container, physical server) to run GameAP Daemon.

You can install GameAP Daemon on a dedicated server automatically by executing a script, or you can do it manually. When installing manually, you should also generate SSL certificates and sign them, as well as edit the configuration file `/etc/gameap-daemon/gameap-daemon.cfg`.

### Automatic GameAP Daemon installation on a dedicated server

Go to the "Administration" → "Dedicated Servers" page, click on the "Create" button.
A window appears prompting you to automatically install GameAP Daemon on a dedicated server. Copy the command and run it on the dedicated server as the superuser. The command looks something like this:
```
curl http://your-panel/gdaemon/setup/zItWHWlI4RKPl9ZsYc3y3Wgd | bash --
```
The command will be available for execution within 5 minutes. After the time elapses, reload the page to see the new command.

## Editing dedicated servers

To edit a dedicated server (node), go to **"Administration"** → **"Dedicated Servers"** page, then 
select the dedicated server you want to edit and click on the **"Edit"** button.

### Description of parameters

#### Basic

##### Name

Dedicated server name. It can take any non-empty value, it does not affect any features.

##### Working directory

This directory contains basic scripts for managing the processes of game servers. Subdirectories
of the working directory contain the files of game servers. For the specified path, 
[game server directory](/en/gameap_configure/game_servers.html#directory) is assigned. The default is `/srv/gameap`.

##### Path to SteamCMD

Path to the SteamCMD directory (the `steamcmd.sh` script is located there). The default is `/srv/gameap/steamcmd`.

##### IP list

List of IP or hosts where game servers will run.

#### Scripts

Scripts for executing commands. Optional. All fields can be empty.

#### GameAP Daemon

Parameters for connecting GameAP Daemon to the server. The panel interacts with a dedicated server through GameAP Daemon 
(GDaemon). GameAP Daemon is a multi-component server with a service that monitors gaming
servers, launches them as needed and executes the necessary commands on them.

These are important settings. If something is specified incorrectly, work with game servers will most likely be impossible.

##### GameAP Daemon host

Host for connecting GameAP Daemon to the server.

##### GameAP Daemon port

Port for connecting GameAP Daemon to the server. The default is `31717`. Please note that this port must be
open in the firewall settings on a dedicated server.

##### GameAP Daemon login and password

GameAP Daemon uses certificates for authentication and security. Login and password can be used to provide
extra security. Fields are optional and are empty by default.
