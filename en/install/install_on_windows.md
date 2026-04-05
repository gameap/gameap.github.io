---
title: Install on Windows
layout: default
lang: en
category: Install GameAP
order: 101
---

## Supported Versions

| Version     | Supported |
|-------------|-----------|
| Server 2025 | ✔         |
| Server 2022 | ✔         |
| Server 2022 | ✔         |
| Server 2019 | ✔         |
| 11          | ✔         |
| 10          | ✔         |


## Downloading GameAP Control

You need to download the GameAP Control utility (gameapctl) 
to manage the GameAP environment.

For this, go to the gameapctl release page on Github: 
[https://github.com/gameap/gameapctl/releases](https://github.com/gameap/gameapctl/releases)

Select the latest release and click on it.

<video width="1280" height="720" controls>
  <source src="/images/en/gameapctl/download.webm" type="video/webm">
  Your browser does not support the video tag.
</video>

After that, find the version suitable for you. 
The most popular architecture is Windows AMD64, 
so you will most likely need to download this archive:

![](/images/en/gameapctl/download_release_windows_amd64.png)

## Installing the Panel Using GameAP Control UI

After downloading the gameapctl archive, run it.

A browser window will open, in which you need to click "Install" 
in the Web/API section.

<video width="1280" height="720" controls>
  <source src="/images/en/gameapctl/windows-install.webm" type="video/webm">
  Your browser does not support the video tag.
</video>

### Installation Parameters

Specify the necessary data for installation.

![](/images/en/gameapctl/ui_gameap_installation.png)

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

#### Database

The database where data will be stored: users, server information, etc. 
You can use:
* [PostgreSQL](https://www.postgresql.org/). Recommended for large projects with many game servers and users.
* [MySQL](https://www.mysql.com/)/[MariaDB](https://mariadb.org/)
* [SQLite](https://www.sqlite.org/). If the load on your server is expected to be low and you do not plan to use more than 10 game servers.

#### Installing GameAP Daemon

In addition to the Web panel, you can also install 
the GameAP Daemon server part, which handles game server operations.

### Completion of Installation

Wait for the installation to finish. 
Some stages may take a considerable amount of time.

Do not forget to save the login data and database information 
that will be provided at the end.

![](/images/en/gameapctl/gameap_finished_installation.png)
