---
title: Game Servers
layout: default
lang: en
category: Panel settings
order: 20
---

* This will become a table of contents (this text will be scraped).
{:toc}

## Editing game servers

To edit a game server (node), go to **"Administration"** -> **"Game Servers"** page, then 
select the game server you want to edit and click on the **"Edit"** button.

### Description of parameters

#### Basic

Basic group of parameters

##### UUID

UUID is a universal unique identifier. It is generated automatically for each server. It is used as 
an identifier for game server processes.

##### Name

Server name. It can be any string value. You can enter any game server name.

##### Game

The game which the game server belongs to. Read the documentation for setting up games in the corresponding section — 
[Game settings](/en/gameap_configure/games.html).

See also [Manual for adding missing games](/en/tutorials/additional_games.html).

##### Mod

The mod which the game server belongs to. The mod may determine additional settings (variables) 
for the game server, such as the default map (`{default_map}`), FPS, and others. See the information on how
add your settings and other detailed info about this on the 
[game settings page](/en/gameap_configure/games.html#variables).

##### RCON password

Password to manage the game server via RCON.

##### Directory

Game server directory relative to [dedicated server working directory](/en/gameap_configure/dedicated_servers.html#working-directory).

For example, if the directory is `servers/my_server`, and the dedicated server working directory  is `/srv/gameap`, then the game
the server will be located in `/srv/gameap/servers/my_server`.

##### Username on a dedicated server

The user running the game server on the node. The default is `gameap`.

The user specified in this field must exist on the dedicated server.

#### Dedicated server, IP, ports

A group of parameters that relates to a dedicated server (VDS, node) and connection to the game server.

##### Dedicated server

##### IP

Game server IP or host. Examples are `127.0.0.1`,`my-server.gameap.ru`.

##### Server port

Main server port. Used to connect players to the server.

##### Query port

Server port for queries. Query is used to get general server data: current map, current number 
of players on the server, list of players.

In GoldSource and Source, it matches the main server port. In Minecraft, you can specify any port. On some
other game servers, it may be one unit more or less.

##### RCON port

Server port for remote administration.

In GoldSource and Source, it matches the main server port. In Minecraft, you can specify any port. In some
other game servers, it may be one unit more or less.

#### Run command

The run command is an important and required parameter for starting a game server. It is a string with various
options and game server starting options. It can be individual for each game and mod.

Parameters may include shortcodes which are automatically replaced by the settings value (variables)
of specific game server. Shortcodes are words in braces `{` and `}`, for example, `{ip}`, 
`{port}`, `{default_map}` and others. 

Read on how to add your own shortcodes which will be automatically replaced by the settings value, on
the [game settings](/en/gameap_configure/games.html#variables) page.