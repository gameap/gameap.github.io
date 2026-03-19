---
title: Games Import
layout: default
lang: en
category: Panel Configuration
order: 330
---

GameAP 4.1 introduced the ability to export and import game settings and templates for creating game servers.

To import games, go to **Administration** → **Games** → **Import**.

![](/images/en/gameap_configure/games_import/import_button.png)

On the import page, upload a YAML file with game settings, then click the **Import** button.

![](/images/en/gameap_configure/games_import/import_page.png)

GameAP supports importing games from other control panels. Currently, import is supported from the following panels:
- [Pterodactyl](https://pterodactyl.io/)
- [Pelican](https://pelican.io/)

On the games import page, you can switch between import formats from different panels.

## Pelican and Pterodactyl Features

Working with imported Pelican Eggs and Pterodactyl Eggs is only possible with
Docker and Podman [process managers](/en/daemon/process_managers.html).

You need to configure GameAP Daemon to work with one of these process managers.
To do this, when adding a new node, select the desired process manager in the "Advanced Settings" section.

![](/images/en/gameap_configure/games_import/daemon_process_manager.png)
