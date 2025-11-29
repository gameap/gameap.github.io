---
title: Requirements
layout: default
lang: en
category: Main
order: 10
---

* This will become a table of contents (this text will be scraped).
{:toc}

## System Requirements

### Web/API

* RAM: 1 GB or more
* Disk space: 200 MB or more
* GameAP is not demanding on the CPU, so one core is sufficient for basic operation.

### GameAP Daemon

GameAP Daemon manages game servers, which require resources, and some game servers can be very demanding. 
The data below represent the requirements for running the server environment without considering the game servers.

* RAM: 128 MB
* Disk space: 1 GB or more
* GameAP Daemon is not demanding on the CPU, so one core is sufficient for operation.

## Package and Operating System Requirements

In most cases, gameapctl will install the necessary packages for the panel and GameAP Daemon operation.

Basic packages required for GameAP Web/API:

* PostgreSQL 16 and above, or
* MySQL 5.6 and above, or MariaDB 10.0 and above. In the case of installing SQLite, no additional database is required.

### Curl

The server must have the curl utility installed, which will allow downloading the installation script.

#### Installation

##### Debian/Ubuntu
```shell
apt install curl
```

##### CentOS

```shell
yum install curl
```