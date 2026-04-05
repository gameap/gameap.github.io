---
title: Getting Started
layout: default
lang: en
category: Main
order: 2
---

To get started, it is desirable to have two dedicated or virtual servers. 
On one is the control panel, on the other game servers. 
You can install everything on a single dedicated server.

## Panel Installation

The panel is installed on a dedicated server with a database (PostgreSQL, MySQL, SQLite).

&nbsp;&nbsp;&nbsp;&nbsp;<img src="/images/icons/linux.svg" alt="Linux" width="20" height="20" style="vertical-align: middle"> [Installation on Linux](/en/install/install_on_linux.html)

&nbsp;&nbsp;&nbsp;&nbsp;<img src="/images/icons/windows.svg" alt="Windows" width="20" height="20" style="vertical-align: middle"> [Installation on Windows](/en/install/install_on_windows.html)

## The Easiest Installation on Linux

If you have Linux, have CURL installed and you don't want to figure installation details out, execute command:
```bash
bash <(curl -s https://gameap.com/install.sh) --with-daemon
```

## Adding a Dedicated Server

Add a new dedicated server (VDS) on which you will then install game servers. 
After installing the panel, log in and select **"Administration"** **"Dedicated servers"** → **"Create"** from the menu. After
that, a window with instructions will open, follow them.

![](/images/en/get_started/add_dedicated_server.gif)

For more detailed information on installation and configuration, read the [Dedicated servers](/en/gameap_configure/dedicated_servers.html) page.

## Adding a Game Server

Go to **"Administration"** → **"Game servers"** → **"Create"**.

![](/images/en/get_started/add_game_server.gif)

For details see [game servers](/en/gameap_configure/game_servers.html) page.
