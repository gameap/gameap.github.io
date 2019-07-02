---
title: Manual panel install
layout: default
lang: en
category: Main
order: 30
---

* This will become a table of contents (this text will be scraped).
{:toc}

Install GameAP 3.x from GitHub. These are detailed instructions for installing GameAP from GitHub, installing required dependencies and building a styles.

> Notes!
1. This tutorial actual for Debian/Ubuntu. Howewer, for Ubuntu/CentOS and other distributions it can works with some changes.
2. This tutorial not actual for a Shared hostings.

Requirements:
* OS: Debian 9
* PHP version: >= 7.2
* Web server: Nginx
* Database: MySQL
* CLI Utilities: [git](requirements.html#git), [composer](requirements.html#composer)

## Install required packages

These packages will be needed in the future:
```bash
sudo apt-get install -y wget software-properties-common ca-certificates apt-transport-https gnupg curl lsb-release
```

Install required PHP extensions:
```bash
sudo apt-get -y install \
    php7.2-common \
    php7.2-cli \
    php7.2-fpm \
    php7.2-pdo \
    php7.2-mysql \
    php7.2-redis \
    php7.2-curl \
    php7.2-bz2 \
    php7.2-zip \
    php7.2-xml \
    php7.2-mbstring \
    php7.2-bcmath \
    php7.2-gmp \
    php7.2-intl
```

Install composer:
```bash
curl -sS https://getcomposer.org/installer | sudo php -- --install-dir=/usr/local/bin --filename=composer
```

> If after the execution of commands you get the following errors:
```
Reading package lists... Done
Building dependency tree       
Reading state information... Done
E: Unable to locate package php7.2-common
E: Couldn't find any package by glob 'php7.2-common'
E: Couldn't find any package by regex 'php7.2-common'
E: Unable to locate package php7.2-cli
E: Couldn't find any package by glob 'php7.2-cli'
E: Couldn't find any package by regex 'php7.2-cli'
...
```
For Debian, run the following command:
```
wget -O - https://packages.sury.org/php/apt.gpg | sudo apt-key add -
echo "deb https://packages.sury.org/php/ $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/php.list
sudo apt-get update
```

Install database and web server. You can use MySQL Nginx:

```
sudo apt-get install mysql-server nginx
```

> Please note that for testing the panel is not required to use MySQL and Nginx, you can use SQLite and the PHP: Built-in web server.
Read more about this in [Simple Environment](#simple-environment)

## Setup GameAP

Create project directory, example `/var/www/gameap`. Clone GameAP from Git:
```bash
git clone https://github.com/et-nik/gameap /var/www/gameap
```

Go to GameAP home directory:

```bash
cd /var/www/gameap
```

The panel configuration is in the `.env` file. Create your configuration file from the example:
```bash
cp .env.example .env
```

Install required PHP dependencies:
```bash
composer install --no-dev --optimize-autoloader
```
Create database, edit `.env`, set host, login, password and database name

Generate an encryption key:
```bash
php artisan key:generate --force
```

After that execute command:
```bash
php artisan migrate --seed
```

This command will create the necessary tables in the database, add the administrator and add some data for example.
Go to the URL page with the control panel.

Administrator default data:
**Login:** admin
**Password:** fpwPOuZD

Do not forget to change your password!

## Build a styles

Install NPM:
```
curl -sL https://deb.nodesource.com/setup_10.x | bash -
apt-get install -y nodejs
```

Install required npm packages:
```
npm install
```

Build a styles:
```
npm run prod
```

Setup complete. All required dependencies are installed, all styles builded. You can use GameAP.

> Do not forget that on the dedicated server on which you plan to host game servers you need to install and configure GameAP Daemon, you can read about it here - [Install GameAP Daemon] (/en/gameap_daemon.html).

## Cron

Add cron schedule

```
* * * * * cd /patch/to/gameap && php artisan schedule:run >> /dev/null 2>&1
```

Set your path instead  `/patch/to/gameap`.

## Simple Environment

If you need a panel to test something, you can use SQLite and the PHP Build-in web server. It is not necessary to install the MySQL database (or any other) and the Nginx web server (or Apache).

Create empty database file:
```
touch /var/www/gameap/database.sqlite
```

Open `.env` and edit the options:
```
DB_CONNECTION=sqlite
DB_DATABASE=/var/www/gameap/database.sqlite
```

Migrate database:
```bash
php artisan migrate --seed
```

Run PHP Build-in Web server:
```
cd public
php -S localhost:8080
```

Go to URL `http://localhost:8080`
Enter the administrator username and password.