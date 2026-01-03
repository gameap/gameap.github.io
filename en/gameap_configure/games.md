---
title: Games
layout: default
lang: en
category: Panel settings
order: 320
---

* This will become a table of contents (this text will be scraped).
{:toc}

## Adding new games

The control panel supports the launch and basic control of any game servers and applications. If
the panel does not contain a game you need, you can add it.

To add a new game, go to **"Administration"** → **"Games"**, then select 
**"Add Game"**.

After adding the game, add the first mod of this game and specify the the game server start parameters.

### Fields 

#### Code

Game code is a unique value. Usually, this is a game name abbreviation, for example for 
"7 Day To Die", the code is `7d2d`. The code must be unique. You can enter any, but it must be unique.

#### Start code

Start code is usually specified in the start parameters. This value can be anything, but more often it matches the game 
code. It may match with other games.

#### Game name

Just enter the full game name.

#### Game engine

Game engine is the system where the game is written. Half-Life, Counter-Strike are written in the GoldSource engine.

Half-Life 2, Counter-Strike Source are written in Source. Sometimes you do not know the game engine,
in this case, you can enter the game code or some game abbreviation.

For games written in GoldSource and Source, enter the exact engine name. For games written in Unity, it's better to enter
some game abbreviation, for example Rust is written in Unity, but the engine name is `rust`.
Minecraft is written without using any engine, so the engine name is also "Minecraft".

#### Engine version

Numeric or string value of the engine version. You can specify any semantic value, for example, "legacy", "beta",
and so on.

#### Steam App ID

Game server ID on Steam. Used to install the server through SteamCMD.
You can find SteamID on [Steam official wiki] (https://developer.valvesoftware.com/wiki/Dedicated_Servers_List), or
in the [SteamDB](https://steamdb.info/).

#### Steam App Set Config

Additional options for installing the server through SteamCMD. You can find some values on 
[Steam official wiki](https://developer.valvesoftware.com/wiki/Dedicated_Servers_List)

#### Local repository

Path on the dedicated server to the archive or to the directory with the game build used as a template.
This path must exist on the dedicated server where a new game server 
is installed. When adding a new game server, the archive will be unpacked into the game server working directory.

The following archives are supported: Zip, 7z, Tar, XZ, Bzip, GZip. 

RAR archives are not supported.

##### Examples

Example game server working directory `/srv/gameap/servers/example-server`

| Local Repository field value | Value validity | Installation result
| ------ | ------- | ------ |
| `/srv/gameap/repo/cs16_gungame.zip` | Valid, if the archive exists on the dedicated server | The `cs16_gungame.zip` archive contents will be unzipped into `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_public` | Valid, if the directory exists on the dedicated server | The directory contents will be copied to `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_gungame.rar` | Invalid. RAR archives not supported | Installation method from the local repository will be skipped or the server will not be installed
| `http://files.gameap.ru/cstrike-1.6/rehlds-amxx-reunion.tar.xz` | Invalid. Value for remote repository is specified | Installation method from the local repository will be skipped or the server will not be installed


#### Remote repository

Link to a remote source. This must be a URL to an HTTP or FTP resource. Archive must be accessible via direct link
without any intermediate pages that require waiting or additional action. Links to Yandex Disk, Google Drive, 
etc. are not supported.

You can find some builds on the [official GameAP repository] (http://files.gameap.ru/).

##### Examples

Example game server working directory `/srv/gameap/servers/example-server`

| Remote Repository field value | Value validity | Installation result
| ------ | ------- | ------ |
| `http://files.gameap.ru/cstrike-1.6/rehlds-amxx-reunion.tar.xz` | Valid | The `rehlds-amxx-reunion.tar.xz` will be loaded and unzipped to `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_gungame.zip` | Invalid. The http or ftp resource must be specified | Installation method will be skipped or the game server will not be installed.

## Adding new mods

Each game can have many mods, each has its own features, settings, 
startup parameters, configuration files, etc.

To add a new game, go to **"Administration"** → **"Games"**, then select 
**"Add Mod"**

### Fields

##### Game

The game which the mod belongs to.

##### Name

Mod name. This may be the name of the addon, build, kernel, any feature, etc. Enter the full name at
your discretion.

#### Repositories

The mod may have its own set of files and settings written over the files of the main build.
You can include any additional plugins, additional content, sounds, music, etc., in the file archive.
Local repository fields are optional.

The following archives are supported: Zip, 7z, Tar, XZ, Bzip, GZip. 

RAR archives are not supported.

##### Local repository

Path on the dedicated server to the archive or to the directory with the game mod files used as a template.
This path must exist on the dedicated server where a new game server 
is installed. When installing the game server, the archive will be unzipped over the main build into the game server working 
directory.

The following archives are supported: Zip, 7z, Tar, XZ, Bzip, GZip. 

RAR archives are not supported.

###### Examples

Example game server working directory `/srv/gameap/servers/example-server`

| Local Repository field value | Correctness of value | Installation result
| ------ | ------- | ------ |
| `/srv/gameap/repo/cs16_gungame.zip` | Valid, if the archive exists on the dedicated server | The `cs16_gungame.zip` archive contents will be unzipped into `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_cs16_gungame` | Valid, if the directory exists on the dedicated server | The directory contents will be copied to `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_gungame.rar` | Invalid. RAR archives not supported | Archive unzipping will be skipped
| `http://files.gameap.ru/cstrike-1.6/rehlds-amxx-reunion.tar.xz` | Invalid. Remote repository value is specified | Mod installation will be skipped


##### Remote repository

Link to a remote source. This must be a URL to an HTTP or FTP resource. Archive must be accessible via direct link
without any intermediate pages that require waiting or additional action. Links to Yandex Disk, Google Drive, 
etc. are not supported.

You can find some builds on the [official GameAP repository](http://files.gameap.ru/).

###### Examples

Example game server working directory `/srv/gameap/servers/example-server`

| Remote Repository field value | Value validity | Installation result
| ------ | ------- | ------ |
| `http://files.gameap.ru/cstrike-1.6/rehlds-amxx-reunion.tar.xz` | Valid | The `rehlds-amxx-reunion.tar.xz` will be loaded and unzipped to `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_gungame.zip` | Invalid. Http or ftp resource must be specified | Mod installation will be skipped

## Editing mods

After making a mod, you can further configure it by specifying additional parameters such as 
"Default startup commands", startup variables, various RCON commands.

### Main settings

Main settings include parameters such as mod name, repositories, and default startup commands. Pay special attention
to the startup commands.

#### Mod name

#### Repositories

Read more in [Adding new mod, repositories](#repositories)

#### Default startup commands

When adding a new game server for the selected mod, this command will be automatically assigned to it. If
the command is empty, nothing will be assigned and you should specify the startup command for each game server manually.
The game server cannot be started without a startup command.

You can use shortcodes in startup commands, they will then be replaced with the values of server variables. Shortcodes
are words without spaces in braces `{`, `}`, for example, `{ip}`, `{port}`, `{maxplayers}`, etc.

##### Basis shortcodes

These shortcodes are always available, they do not require adding additional variables in the mod settings.

| Shortcode | Description
| ------ | -------
| {ip} | Game server IP
| {port} | Game server main port. Sometimes called a connect port
| {query_port} | Query port
| {rcon_port} | Server communication port (RCON port)
| {rcon_password} | RCON password
| {uuid} | Server UUID
| {uuid_short} | Short server UUID

##### User-defined shortcodes

You can define these shortcodes for each specific game mod yourself. Depending on 
individual server settings, these shortcodes will be replaced with values of the game server parameters. Read more 
in [Variables](#variables).

### Variables

You can add individual settings for each game server. Then these settings can be edited
by administrator or regular user on the settings page (**"Server List"** → **"Administration"** → 
**"Settings"**).

| Field | Description
| ------ | -------
| Variable | Variable name. No curly brackets.
| Default | Variable value by default. This value will be used if an individual value is not set for the game server.
| Description | Game server description on the settings page (**"Server List"** → **"Administration"** → **"Settings"**)
| Admin variable | If checked, only the administrator will be able to edit this setting for game servers.

#### Examples

| Values | Description
| ------ | -------
| **Variable:** default_map <br><br>**Default:** de_dust <br><br>**Description:** Map at startup | For each game server of this mod, the shortcode `{default_map}` and a new parameter in settings called "Default Map" with the default value "de_dust" will appear. <br><br>Individually for each server, this parameter can be edited at **"Server List"** → **"Administration"** → **"Settings"**


### RCON commands

These commands allow more advanced game server administration. If the game supports working with RCON or 
if the console is supported, you can do the following: kick players from the server, ban players, change the server map,
send text messages to the common chat, set a password.

Some features may be limited by the game itself or by the mod. For example, not all game servers support
server login with a password.

#### Kick command

You can set the RCON command to kick a player from the server.

You can set shortcodes for the command that will be replaced with the data of a certain player

| Shortcode | Description
| ------ | -------
| {id} | Server player ID
| {name} | Server player name

For many GoldSource/Source games, this is the command: 
```
kick #{id}
```

#### Ban command

You can set a command that will be used for a temporary or permanent ban of a player on the server.

You can set shortcodes for the command that will be replaced with the data of a certain player

| Shortcode | Description
| ------ | -------
| {id} | Server player ID
| {name} | Server player name
| {time} | Ban time
| {reason} | Ban reason

For many GoldSource games (Half-Life, Counter-Strike 1.6, etc.) running AMX Mod X, this is the command: 
```
amx_ban "{name}" {time} "{reason}"
```

#### Name (nickname) change command

With this command, you can change the selected player's nickname.

You can set shortcodes for the command that will be replaced with the data of a certain player

| Shortcode | Description
| ------ | -------
| {id} | Server player ID
| {name} | Current server player name
| {new_name} | New server player name
| {reason} | Reason for change

For many GoldSource games (Half-Life, Counter-Strike 1.6, etc.) running AMX Mod X, this is the command: 
```
amx_nick #{id} {new_name}
```

#### Restart command

You can set the command to soft restart the server, without restarting the game server process. Usually,
this command restarts the game map or round. Not supported by many games.

For many GoldSource/Source games, this is the command: 
```
restart
```

#### Map change command

With this command, you can change the game server map.

| Shortcode | Description
| ------ | -------
| {map} | Map name

For many GoldSource/Source games, this is the command: 
```
changelevel {map}
```

#### Send message command

With this command, you can send a text chat message to all players on the server.

| Shortcode | Description
| ------ | -------
| {msg} | Message to be sent to the server

For many GoldSource games (Half-Life, Counter-Strike 1.6, etc.) running AMX Mod X, this is the command: 
```
amx_say "{msg}"
```

#### Set/change password command

You can set a password for the game server, so only players who know
this password, can log in.

| Shortcode | Description
| ------ | -------
| {password} | Server password

For many GoldSource/Source games, this is the command: 
```
password {password}
```

### FastRCON commands

You can specify your optional RCON commands. For example, server status command, receiving statistics, receiving
list of recently disconnected players, etc.
