---
title: Adding Missing Games
layout: default
lang: en
category: Tutorials
order: 299
---

The control panel supports the launch and basic control of any game servers and applications. 
This manual covers making a new game by the example of Sven Co-op. Each step will include explanations.

## Game adding

First, go to the add game page. Go to **"Administration"** menu, then select 
**"Games"**.

![](/images/en/tutorials/additional_games/game_menu.png)

Next, find the **"Add Game"** button on the top of the page and click on it.

![](/images/en/tutorials/additional_games/add_game_menu.png)

You will be redirected to the new game add page. Here you specify some details of your new game.

![](/images/en/tutorials/additional_games/example_add_svencoop.png)

You must fill in the following fields:
* **Code**. Enter the abbreviated game name
* **Start code**. You can also enter the abbreviated game name.
* **Game name**
* **Game engine**. If the game is written in Unity or without using the engine, enter the abbreviated
game name
* **Version**. Enter the version number or semantic value, for example “legacy”, “beta”, etc.

Note! To enable automatic installation, you must fill in one of the following fields: 
* [**Steam APP ID**](/en/gameap_configure/games.html#steam-app-set-config). Find out the value for the game you are interested in at 
[Steam official wiki](https://developer.valvesoftware.com/wiki/Dedicated_Servers_List), or in the database 
[SteamDB](https://steamdb.info/)
* [**Remote repository**](/en/gameap_configure/games.html#remote-repository). You can find some archives in the [GameAP repository](http://files.gameap.ru/)
* [**Local repository**](/en/gameap_configure/games.html#local-repository).

All these fields are optional, but one of them must have a value.

Read more about the meaning of fields at [Game Settings, Fields Description](/en/gameap_configure/games.html#fields) page.

## Mod adding

Each game must have at least one mod. 

Mod is a strong GameAP tool, you can enable additional plugins, configuration or any content for 
extending the basic server capabilities. The archive with files you specify for the mod will be unzipped over 
the base server build.

To add a new mod for a specific game, select the game in the list and click **"Add the first mod"**.

![](/images/en/tutorials/additional_games/example_menu_add_mod_svencoop.png)

If the game already has at least one mod, then at the very top of the game list page, select 
**"Add mod"**.

On the mod adding page, specify the mod name depending on the 
game mode features (GunGame, Jail, etc.), or availability of any modules (AMXX, ReAMXX for Counter-Strike,
 IndustrialCraft, BuildCraft for Minecraft, etc.).
 
![](/images/en/tutorials/additional_games/example_add_svencoop_mod.png)

If you have an archive with additional plugins to be written over the base build, then specify the path to
it in the local or remote repository fields. 

In the **local repository** field, specify the path to the archive or directory on a dedicated server 
running GameAP Daemon; path example `/srv/gameap/repo/svencoop_op4_maps.tar.xz`. 
See details on [Game Settings](/en/gameap_configure/games.html#local-repository-1) page.

In the **remote repository** field, specify the URL to the HTTP or FTP archive. 
Path example `http://files.gameap.ru/svencoop/svencoop_op4_maps.tar.xz`.
See details on [Game Settings](/en/gameap_configure/games.html#remote-repository-1) page. 
You can see examples of archives in the [GameAP repository](http://files.gameap.ru/)

## Mod configuring

After making a mod, you can further configure it by specifying additional parameters such as 
"Default startup commands", startup variables, various RCON commands.

Default startup commands should be specified. If you do not specify them, then the startup command will be empty
when creating a new game server, but it must be specified, otherwise the server will not start.

![](/images/en/tutorials/additional_games/game_mods_edit_basic.png)

Examples of default startup commands for some Linux games:
* Sven Co-op: 
```shell
./svends_run +ip {ip} +port {port} +maxplayers {maxplayers} +log on +map {default_map}
```
* Half-Life:
```shell
./hlds_run -game valve +ip {ip} +port {port} +map {default_map} +maxplayers {maxplayers} +sys_ticrate {fps}
```

* GTA: Multi Theft Auto
```shell
./mta-server64 -t -n --ip {ip} --port {port} --maxplayers {maxplayers}
```

Examples of default startup commands for some Windows games:

* Sven Co-op:
```shell 
SvenDS +ip  {ip} +port {port} +maxplayers {maxplayers} +log on +map {default_map}
```

* 7 Day To Die
```shell
startdedicated.bat
```

Pay attention to the values in braces `{` and `}`, such as `{ip}`, `{port}`, `{maxplayers}`, `{default_map}`,
`{fps}` and others. In GameAP, they are called shortcodes; they are replaced with server variable values.
All game servers have so-called basic variables, such as IP, ports, ID, UUID. They also have additional
variables specified in the mod settings, these include the maximum number of games, default map, 
FPS and others.

Some parameters can be changed only by administrators, and some are available for change to ordinary users. 
The list of variables and their names are specified in the mod settings, on the "Variables" tab.

![](/images/en/tutorials/additional_games/game_mods_edit_vars.png)

The variables specified in the mod for each game server can then be changed individually in the settings.

The next tab in the mod settings is "RCON commands". You can specify RCON commands for player kick, ban, 
map change, and other RCON commands, they are used for advanced game server administration.

![](/images/en/tutorials/additional_games/game_mods_edit_commands.png)

You can specify your optional RCON commands on the Fast Rcon tab. For example, the server status command or 
statistics.

![](/images/en/tutorials/additional_games/game_mods_edit_fast_rcon.png)
