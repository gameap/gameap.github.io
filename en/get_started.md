---
title: Getting Started
layout: default
lang: en
category: Main
order: 2
---

* This will become a table of contents (this text will be scraped).
{:toc}

To get started, it is desirable to have two dedicated or virtual servers. On one is the control panel, on the other game servers. You can install all on single dedicated server, but it is not recomended.

## The easiest installation

If you have Linux, have CURL packages and you don't want to figure installation details out. Execute command:
```bash
bash <(curl -s https://gameap.com/install.sh)
```

## Panel Installation

The panel is installed on a dedicated server with PHP, 
a web server (Nginx, Apache, etc.), and a database (MySQL, PgSQL, SQLite, ...)

* [Installation on Linux](/en/install/install_on_linux.html)
* [Installation on Windows](/en/install/install_on_windows.html)

### Outdated Installation Methods

The methods below are more complicated to install but are still relevant. Use them if there is a specific need, but remember that in most cases, it is better to use the methods described above.

* [Install on Linux with SSH (auto, easy)](/en/auto_install.html)
* [Install on Linux with SSH (manual, difficult)](/en/manual_install.html)
* [Install on shared hosting (not recommended)](/en/shared_install.html)

### Download

You can download and unpack archive on your hosting. Some installation methods don't required archive downloading.

* [gameap-3.1-shared.zip](https://packages.gameap.com/gameap/gameap-3.1-shared.zip)
* [GitHub](https://github.com/et-nik/gameap)


## Adding a Dedicated Server

Add a new dedicated server (VDS, VPS, Physical server). On a dedicated server, game servers will then start.
Go to **"Administration"** -> **"Dedicated servers"** -> **"Create"**.
A window appears prompting you to automatically install GameAP Daemon on a dedicated server.
Follow the instructions.

![](/images/en/get_started/add_dedicated_server.gif)

For more detailed information on installation and configuration, read the [Dedicated servers](/en/gameap_configure/dedicated_servers.html) page.

## Game server adding

Go to **"Administration"** -> **"Game servers"** -> **"Create"**.

![](/images/en/get_started/add_game_server.gif)

For details see [game servers](/en/gameap_configure/game_servers.html) page.