---
title: Quake III Arena
layout: default
lang: en
category: Tutorials
order: 204
---

Quake III Arena is a multiplayer first-person shooter developed by id Software. 
The game focuses entirely on multiplayer combat and lacks a traditional single-player campaign.

The game is similar to Unreal Tournament and other first-person shooters from the late 1990s and early 2000s.

## Environment Setup

GameAP provides support for Quake III Arena game servers, their configuration and management.

* [Installing GameAP on Linux](/en/install/install_on_linux.html)
* [Installing GameAP on Windows](/en/install/install_on_windows.html)

### Installing GameAP Daemon

GameAP Daemon is an agent responsible for managing game servers on dedicated machines. To install a game server on a machine (VDS), you need to install GameAP Daemon.

When installing GameAP, you can choose a full installation including the Daemon (using the `--with-daemon` flag).

In the control panel, go to **"Administration"** → **"Dedicated Servers"** → **"Create"**. A window will appear offering automatic installation. Copy the code and run it on the dedicated server.

After that, you can proceed with installing the Quake III Arena server.

## Installing Quake III Arena Server in GameAP

Go to **Administration** → **Game Servers** → **Create**

![](/images/en/tutorials/quake3/create_form.png)

* In the "Name" field, enter any server name, for example "My Quake III Server".
* In the "Game" field, select "Quake 3" from the dropdown list.
* In the "Modification" field, specify the desired modification; the default is `ioquake3`. The `quake3e` option is also available.
* In the "Dedicated Server" field, specify the node where the game server will be hosted. By default, the first available node is selected.
* In the "IP" field, select the desired address for your server, then specify an available port or leave the one suggested by the system.

## Configuring the Quake III Arena Server

To change the server configuration, go to the **Server** section, select your server, and click **Manage**. Then open the **Settings** tab.

![](/images/en/tutorials/quake3/settings.png)

For changes to take effect, you need to restart the server.

### Server Hostname

The name of your Quake III server that will be visible to all players in the server list when searching.

### Maximum Players on Server

The maximum number of players that can be on the server at the same time. The default value is 16.

### Default Map of the Server

The map that will automatically load when the server starts. The default map is `q3dm17`, known as "The Longest Yard".

### Timelimit in Minutes

The maximum round duration in minutes. The default value is 20 minutes.

### Frag Limit

The number of frags (kills) at which the round ends. The default value is 20 frags.