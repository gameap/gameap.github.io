---
title: Install on Windows
layout: default
lang: en
category: Install GameAP
order: 101
---

## Supported Versions

| Version      | Supported | Note                                           |
|--------------|-----------|------------------------------------------------|
| Server 2022 | ✔         | Last manual test (v0.9.1): March 2, 2024       |
| Server 2019 | ✔         | Last manual test (v0.9.3): March 2, 2024       |

## Downloading GameAP Control

You need to download the GameAP Control utility (gameapctl) 
to manage the GameAP environment.

For this, go to the gameapctl release page on Github: 
[https://github.com/gameap/gameapctl/releases](https://github.com/gameap/gameapctl/releases)

Select the latest release and click on it.

![](/images/en/gameapctl/download.png)

After that, find the version suitable for you. 
The most popular architecture is Windows AMD64, 
so you will most likely need to download this archive:

![](/images/en/gameapctl/download_release_windows_amd64.png)

## Installing the Panel Using GameAP Control UI

After downloading the gameapctl archive, run it.

![](/images/en/gameapctl/exe_in_archive.png)

A browser window will open, in which you need to click "Install" 
in the Web/API section.

![](/images/en/gameapctl/ui_gameap_install_button.png)

### Installation Parameters

Specify the necessary data for installation.

![](/images/en/gameapctl/ui_gameap_installation.png)

#### Installation Path

This is the path where the main panel files will be stored.
By default, it is `C:\gameap\web`, and this path cannot currently be changed.

#### Host

Specify the domain or IP address where the panel will be accessible.

In the case of an IP address, it must be the address assigned to 
the network interface on the VDS. If your network uses NAT, 
do not specify the external IP, but the internal one, 
then configure port forwarding.

Any domain can be specified, but do not forget to configure DNS.

Correct value examples:
* 10.182.104.8
* 10.182.104.8:2080
* example.com
* http://example.com

#### Web Server

The HTTP server that will accept and process incoming requests. 
[Nginx](https://www.nginx.com/) is recommended.

#### Database

The database where data will be stored: users, server information, etc. 
You can use [MySQL](https://www.mysql.com/)/[MariaDB](https://mariadb.org/) 
and [SQLite](https://www.sqlite.org/).

If the load on your server is expected to be low, 
and you do not plan to use more than 10 game servers, 
you can use SQLite.

#### Installing GameAP Daemon

In addition to the Web panel, you can also install 
the GameAP Daemon server part, which handles game server operations.

### Completion of Installation

Wait for the installation to finish. 
Some stages may take a considerable amount of time.

Do not forget to save the login data and database information that will be provided at the end.

![](/images/en/gameapctl/gameap_finished_installation.png)