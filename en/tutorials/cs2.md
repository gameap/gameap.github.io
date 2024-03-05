---
title: Counter-Strike 2
layout: default
lang: en
category: Tutorials
order: 201
---

Counter-Strike 2 is a tactical first-person shooter released in 2023, 
developed and published by Valve. It is the fifth main installment in 
the Counter-Strike series, building on the success of its predecessors 
with updated graphics, gameplay mechanics, and new features.

## Setting Up the Environment

GameAP fully supports Counter-Strike 2, including player management. 
To start working with the GameAP control panel, you need to install 
it by choosing one of the options provided:

* [Installing GameAP on Linux](/en/install/install_on_linux.html)
* [Installing GameAP on Windows](/en/install/install_on_windows.html)

### Installing GameAP Daemon

To host a game server on your virtual or dedicated server machine, 
you need to install GameAP Daemon. 
When installing GameAP, you can choose the option to perform 
a full installation along with the Daemon.

In the control panel, go to **"Administration"** -> **"Dedicated Servers"** 
-> **"Create"**. 
A window with an offer for automatic installation will appear. 
Copy the code and execute it on the dedicated server.

After this, you can proceed with the Counter-Strike server installation.

## Installing Counter-Strike in GameAP

Go to **Administration** -> **Game Servers** -> **Create**

![](/images/en/tutorials/cs2/create_form.png)

* In the "Name" field, enter any server name, 
  for example, "My Counter-Strike 2 Server".
* In the "Game" field, select Counter-Strike 2.
* In the modification field, choose the modification.
* In the "Dedicated Server" field, select the desired node 
  on which the game server will be located.
* In the IP field, select the desired address of your server, 
  then you can choose a free port of your server, or use the suggested one.

### Server Token

After installation, you need to specify a unique token `sv_setsteamaccount` 
for the server; without it, the server will not work.

To generate it, go to [https://steamcommunity.com/dev/managegameservers](https://steamcommunity.com/dev/managegameservers).

![](/images/en/tutorials/cs2/token_generation.png)

After generation, the token value will appear in the table, 
use the 32-character value:

![](/images/en/tutorials/cs2/token_table.png)

You need to copy the token value and specify it in the control panel settings. 
Go to **Servers** -> select your server -> **Management** -> **Settings**

![](/images/en/tutorials/cs2/set_token.png)

After this, you can launch your server.

## Configuring Counter-Strike 2 Server

To change the server configuration, in the file manager, 
go to the `/game/csgo/cfg` directory, where you can find many *.cfg files.

The main configuration file for Counter-Strike 2 server is `server.cfg`

![](/images/en/tutorials/cs2/server_config.png)