---
title: Rust
layout: default
lang: en
category: Tutorials
order: 202
---

Rust is a survival game developed by Facepunch Studios. 
Players need to gather resources, craft items, 
and build shelters to survive in a hostile open-world environment.

Rust is somewhat similar to [Minecraft](/en/tutorials/minecraft.html) 
with its survival system and item crafting.

## Environment Setup

GameAP provides full support for Rust game servers, including player management. 
Installing the panel will take a little time; choose one of the following 
control panel installation options:

* [Installing GameAP on Linux](/en/install/install_on_linux.html)
* [Installing GameAP on Windows](/en/install/install_on_windows.html)

### Installing GameAP Daemon

To install the game server on a machine (VDS), 
you need to install GameAP Daemon. 
During the GameAP installation, you can choose to install it fully, 
including Daemon.

In the control panel, go to **"Administration"** → **"Dedicated Servers"** 
→ **"Create"**. 
A window with an offer for automatic installation will appear. 
Copy the code and execute it on the dedicated server.

After this, you can proceed to install the Rust server.

## Installing a Rust Server in GameAP

Go to **Administration** → **Game Servers** → **Create**

![](/images/en/tutorials/rust/create_form.png)

* In the "Name" field, enter any server name, for example, "My Rust Server".
* In the "Game" field, select Rust from the list.
* In the modification field, select the modification, by default Vanilla.
* In the "Dedicated Server" field, select the desired node where the game server will be located.
* In the IP field, choose the desired address of your server, then you can choose a free port for your server, or use the suggested one.

## Configuring the Game Server

To change the server configuration, go to the **Server** section, 
select your server, and click **Management**. Then, open the **Settings** tab.

![](/images/en/tutorials/rust/settings.png)

You need to restart the server for the settings to take effect.

### Server Hostname

This is the name of your Rust server, which will be displayed 
to all players in the game in the server search window.

### Maximum Players on Server

The maximum number of players that can play on the server simultaneously. 
By default, 32.

### Map on the Server

The map on the server. 
By default, the procedurally generated map Procedural Map. 
Possible values:

* Barren
* Craggy Island
* Hapis
* Savas Island

### Map Generation Seed

The seed for generating the Procedural Map.

### World Size

The size of the world.

### Server Save Interval

The world save interval. 
This is the period during which the server data will be saved to disk.