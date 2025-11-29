---
title: Install on Linux
layout: default
lang: en
category: Install GameAP
order: 100
---

* This will become a table of contents (this text will be scraped). 
{:toc}

## Installation

Installation on Linux is performed with a single command:
```shell
bash <(curl -s https://gameap.com/install.sh)
```

During the installation process, you will be prompted to enter some information.

### Host

Specify the domain or IP address at which the panel will be accessible.

In the case of an IP address, it must be an address assigned to the network 
interface on the VDS. If your network uses NAT, do not specify 
the external IP, but rather the internal one, and then configure port 
forwarding.

Any domain can be specified, but do not forget to configure DNS.

Examples of correct values:
* 10.182.104.8
* 10.182.104.8:2080
* example.com
* http://example.com


### Database

The database where data will be stored: users, information about servers, etc. 
You can use [PostgreSQL](https://www.postgresql.org/),
[MySQL](https://www.mysql.com/) and 
[SQLite](https://www.sqlite.org/).

PostgreSQL is recommended in most cases. If the load on your server 
is expected to be low, and you do not plan to use more than 10 game servers, 
you can use SQLite.

Some distributions may have [MariaDB](https://mariadb.org/) installed.

## Completing the Installation

At the end of the installation, the access details for the panel 
will be displayed. Do not forget to save this information to access the panel.

![](/images/en/gameapctl/gameap_finished_installation.png)

## Additional Installation Options

### Non-Interactive Installation

This type of installation allows you to install the panel without 
entering any data during the process. 
You can add the `--database` flags, 
and in this case, the installer will not require any additional data from you.

```shell
bash <(curl -s https://gameap.com/install.sh) \
  --database=postgresql \
  --host=http://127.0.0.1 \
  --port=80
```

### Full Installation

To install the GameAP Daemon in addition to the panel itself, 
add the `--with-daemon` flag.

This method is recommended if you plan to host both the panel 
and game servers on the same VDS.

```Shell
bash <(curl -s https://gameap.com/install.sh) --with-daemon
```