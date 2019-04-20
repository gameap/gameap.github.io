---
title: Automatic panel installation
layout: default
lang: en
category: Main
order: 20
---

* This will become a table of contents (this text will be scraped).
{:toc}

Fully automatic installation. You need to run the script, it will automatically install the packages dependencies and panel. During the installation, you will need to enter and select some parameters:

- Installation path
- Pahel host. Domain or IP
- Select database: MySQL, PostgreSQL, SQLite
- Select web server: Nginx, Apache

## Dependencies installation

To run the script you need CURL.

### Debian/Ubuntu

```bash
sudo apt-get update
sudo apt-get install curl
```

### Centos

```bash
sudo yum update
sudo yum install curl
```

## Running auto installation script

Download script and set execute permission:
```bash
curl -sLO http://packages.gameap.ru/installer.sh
chmod +x ./installer.sh
```

Run script:
```
./installer.sh
```

## Running with parameters

For silent installation no questions asked set parameters:

- `--path` Path to panel files.
- `--host` Panel host, IP or domain name.
- `--web-server` Possible values: `nginx`, `apache`, `none`
- `--database` Possible values: `mysql`, `pgsql`, `sqlite`, `none`
- `--github` Script will install and build panel from GitHub'а. Recommended options.
- `--upgrade` Use this option to upgrade panel

### Examples

#### Panel installation

Next example au

The following example will automatically install the panel in the directory `/var/www/gameap`, install and configure the web server, database. The panel should be available at the address `http://your-gameap.ru` specified in `--host`:

```bash
./installer.sh \
    --path=/var/www/gameap \
    --host=your-gameap.ru \
    --web-server=nginx \
    --database=mysql \
    --github
```

The following example will automatically install the panel in the directory `/var/www/gameap`. Web server will not be installed. The php-sqlite package will be installed for working with the SQLite database.

```bash
./installer.sh \
    --path=/var/www/gameap \
    --host=localhost \
    --web-server=none \
    --database=sqlite
```

#### Panel upgrading

The followind example will upgrade panel to last available version:
```bash
./installer.sh --upgrade
```

Upgrading panel using github:
```bash
./installer.sh --upgrade --github
```

## Other packages that the script installs

In addition to the panel itself, the script installs the necessary packages and their dependencies. Most likely, all or most of them will already be installed on your system.

### Packages

The script automatically installs the following packages.:

- `software-properties-common` APT Repository Management package.
- `apt-transport-https` https support for APT
- `gnupg` Package for working with digital signatures and keys. Required to authenticate packages and add repository keys.

If you use `--github` options, the script will also install the following packages:
- `git` Git Tool. To download and update panel from GitHub
- `composer` PHP Package Manager. To install packages dependencies.
- `npm` NodeJS Package Manager. To build GameAP Styles.

PHP extensions: php-cli, php-fpm, php-pdo, php-mysql, php-redis, php-curl 
php-bz2, php-zip, php-xml, php-mbstring, php-bcmath

### Repositories

The script can add multiple repositories to APT. For example, in Debian Stretch, the default is PHP 7.0, and the minimum PHP version for panel 7.1. In this case, the script will check the ability to install the required version of PHP, if it is not, then add the necessary repositories.

- http://packages.gameap.ru/
GameAP repository. Delete file `/etc/apt/sources.list.d/gameap.list` if you want to remove repository 

- https://packages.sury.org/php/
This repository will be added to Debian Stretch and Jessie if the ability to install PHP >= 7.1 is not possible.
Delete file `/etc/apt/sources.list.d/php.list` if you want to remove repository 

- ppa:ondrej/php
Repository with the latest versions of PHP for Ubuntu. Will be added to Ubuntu Trusty and below.
To delete, run the command: `sudo ppa-purge ppa:ondrej/php`

- http://nginx.org/packages/
Nginx official repository. It will be added if Nginx was selected as a web server.
Delete file `/etc/apt/sources.list.d/nginx.list` if you want to remove repository 

- ppa:chris-lea/node.js
NodeJS Repository. 
Репозиторий для установки NodeJS менеджера пакетов (NPM). It will be added only when installing the panel from GitHub.
To delete, run the command: `sudo ppa-purge ppa:chris-lea/node.js`