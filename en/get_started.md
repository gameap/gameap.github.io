---
title: Get Started
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
curl -sLO http://packages.gameap.ru/installer.sh && bash installer.sh
```

## Panel Installation

The panel is installed on a dedicated server with PHP, a web server (Nginx, Apache, etc.) and a database (MySQL, PgSQL, SQLite, ...)

Installation ways:

* [Install on Linux with SSH (auto, easy)](/en/auto_install.html)
* [Install on Linux with SSH (manual, difficult)](/en/manual_install.html)
* [Install on shared hosting (not recommended)](/en/shared_install.html)

### Download

You can download and unpack archive on your hosting. Some installation methods don't required archive downloading.

* [GitHub](https://github.com/et-nik/gameap)
* [gameap_latest.zip](http://www.gameap.ru/gameap_latest.zip)

## Dedicated server adding

Add a new dedicated server (VDS, VPS, Physical server). On a dedicated server, game servers will then start.
Go to **"Administration"** -> **"Dedicated servers"** -> **"Create"**.
A window appears prompting you to automatically install GameAP Daemon on a dedicated server.
Follow the instructions.

![](/images/en/get_started/add_dedicated_server.gif)

For details see [dedicated servers](/en/gameap_configure/dedicated_servers.html) page.

## Game server adding

Go to **"Administration"** -> **"Game servers"** -> **"Create"**.

![](/images/en/get_started/add_game_server.gif)

For details see [game servers](/en/gameap_configure/game_servers.html) page.

## Demo

Look GameAP features on the demo site.

[https://demo.gameap.ru](https://demo.gameap.ru)

* Login: demo
* Password: demo
