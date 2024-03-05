---
title: Grand Theft Auto V
layout: default
lang: en
category: Tutorials
order: 203
---

Grand Theft Auto V is a popular open-world game developed by Rockstar North.

FiveM is a modification for GTA V that allows online multiplayer play. 
In addition to FiveM, there is also a similar modification called RageMP.

You can create and manage your own GTA V (FiveM) server with the help of GameAP.

## Environment Setup

With GameAP, you can manage the main parameters of your server. 
GameAP provides full support for FiveM game servers, 
including player management.

Installing the panel will take a little time. Depending on your operating system, choose one of the following control panel installation options:

* [Installation on Linux](/en/install/install_on_linux.html)
* [Installation on Windows](/en/install/install_on_windows.html)

### Installing GameAP Daemon

To install the game server on a machine (VDS), you need to install GameAP Daemon. 
During the GameAP installation, you can choose to install it fully, 
including Daemon.

In the control panel, go to **"Administration"** -> **"Dedicated Servers"** 
-> **"Create"**. A window with an offer for automatic installation will appear. 
Copy the code and execute it on the dedicated server.

After this, you can proceed to install the FiveM server.

### Creating a FiveM Server in GameAP

Go to **Administration** -> **Game Servers** -> **Create**

![](/images/en/tutorials/fivem/create_form.png)

* In the "Name" field, enter any server name, for example, "My GTA V Server".
* In the "Game" field, select FiveM from the list.
* In the modification field, select the modification, by default Vanilla.
* In the "Dedicated Server" field, select the desired node where the game server will be located.
* In the IP field, choose the desired address of your server, then you can choose a free port for your server, or use the suggested one.

## Server Key

After installation, you need to specify a key that you must obtain from 
[keymaster.fivem.net](https://keymaster.fivem.net)

![](/images/en/tutorials/fivem/generate_key.png)

After generating, you will see a message. You need to copy the key value.

![](/images/en/tutorials/fivem/key.png)

You need to copy the key value and specify it in 
the settings in the control panel. 
Go to **Servers** -> select your FiveM server -> **Management** -> **Settings**

![](/images/en/tutorials/fivem/set_key.png)

Now you can start your FiveM server in the panel.

## Configuring the FiveM Server

The FiveM server configuration is located in the `server.cfg` file, 
which is located in the root directory. 
You can edit this file in the panel's file manager.

Go to **Servers** -> select your FiveM server -> **Management** -> **Files**

![](/images/en/tutorials/fivem/server_config.png)