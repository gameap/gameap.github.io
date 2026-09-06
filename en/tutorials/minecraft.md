---
title: Minecraft
layout: default
lang: en
category: Tutorials
order: 200
---

Minecraft is a popular sandbox game developed 
by Mojang Studios, a subsidiary of Xbox Game Studios. 
It allows players to build, explore, and survive in 
a procedurally generated world made of blocks. 
The game features various modes: survival, creative, 
and adventure. Minecraft has a large community of players 
who create and share various creations, mods, and skins.

## Environment Setup

To begin, you need to [install GameAP](/en/get_started.html#panel-installation), which will take a few minutes:

* [Installing GameAP on Linux](/en/install/install_on_linux.html)
* [Installing GameAP on Windows](/en/install/install_on_windows.html)

### Installing GameAP Daemon

On the node where the game server will be hosted, you need 
to install GameAP Daemon. During the GameAP installation process, 
you can choose the option to install it together with the Daemon.

In the control panel, go to **Administration** → **Dedicated servers** 
→ **Create**. 
A window with an offer for automatic installation will appear. 
Copy the code and execute it on the dedicated server.

After this, you can proceed to install the Minecraft server.

## Installing Minecraft in GameAP

Go to **Administration** → **Game servers** → **Create**

![Create game server form for Minecraft](/images/en/tutorials/minecraft/create_form.png)

* In the "Name" field, enter any server name.
* In the "Game" field, select Minecraft.
* In the "Game Mod" field, we recommend choosing "Multicore", which allows for easy management of Minecraft server versions.
* In the "Dedicated Server" field, select the desired node where the game server will be located.
* In the IP field, choose the desired address of your server, then you can choose a free port for your server, or use the suggested one.

Watch a short video on the process of installing a Minecraft server in GameAP:

<video controls width="600">
  <source src="/media/en/tutorials/minecraft/installation.webm" type="video/webm" />
  <source src="/media/en/tutorials/minecraft/installation.mp4" type="video/mp4" />
</video>

## Managing the Minecraft Server

To manage the Minecraft server, go to the **Servers** section, then select your server and click **Control**. On this page, you will be able to start, stop your server, view its console, and send commands.

![Minecraft server management page with the console](/images/en/tutorials/minecraft/server_management.png)

### Changing the Server Version

If you chose the "Multicore" modification, you can change the version of your server. To do this, go to the main page of managing your server, then select **Settings**.

![Minecraft game server settings with the version selector](/images/en/tutorials/minecraft/server_settings.png)

The "Multicore" modification provides the following settings:

| Setting                             | Type                          | Default | Description                                                                                                  |
|-------------------------------------|-------------------------------|---------|--------------------------------------------------------------------------------------------------------------|
| Minecraft version                   | Select (custom value allowed) | 1.20.4  | Server version. Pick one from the list or type any other version                                             |
| Mod                                 | Select                        | Vanilla | Server core: Vanilla, Paper, Forge, Fabric, Spigot, CraftBukkit, Cauldron, Waterfall, Velocity or BungeeCord |
| Mod version                         | String                        | —       | Version of the selected core, for example the Forge build                                                    |
| Memory. Max heap size (Xmx)         | String                        | 1G      | Maximum Java heap size                                                                                       |
| Min memory. Initial heap size (Xms) | String                        | 512M    | Initial Java heap size                                                                                       |

All of these settings are marked **Admin only**, so they are shown to administrators only.

#### Examples of Settings

##### Default

* Minecraft version: 1.20.4
* Mod: Vanilla

##### Forge 1.19

In this example, we will use a Minecraft server version 1.19.4 with Forge API for mod support version 45.0.50

* Minecraft version: 1.19.4
* Mod: Forge
* Mod version: 45.0.50

##### Spigot

In this example, we will use a Minecraft server version 1.19.4 with Spigot API for plugin support.

* Minecraft version: 1.19.4
* Mod: Spigot

### server.properties

`server.properties` is the main configuration file for the Minecraft server. You can edit it using the file manager in GameAP.

![The server.properties file in the GameAP file manager](/images/en/tutorials/minecraft/server_properties.png)