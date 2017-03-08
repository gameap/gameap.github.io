---
title: Manual panel install
layout: default
lang: en
category: Main
---

# Manual panel install

Install GameAP 2.x from GitHub.

Notes!
1. This tutorial actual for Debian 8. Howewer, for Ubuntu/CentOS and other distributions it can works with some changes.
2. This tutorial not actual for a Shared hostings.

I use:
* OS: Debian 8
* PHP version: 7.1
* Web server: Nginx
* Database: MySQL
* CLI Utilities: [git](requirements.html#git), [composer](requirements.html#composer)

Install required PHP extensions:
```bash
apt-get -y install php7.1-cli php7.1-fpm php7.1-pdo php7.1-mysql php7.1-redis php7.1-gd php7.1-mcrypt php7.1-curl php7.1-bz2 php7.1-xml php7.1-mbstring php7.1-bcmath
```

Create project directory, example **/var/www/gameap**. Clone GameAP from Git:
```bash
git clone https://github.com/et-nik/GameAP /var/www/gameap
```

Install required PHP dependencies:
```bash
composer install
```

Create database, edit  application/config/database.php.

Update database schema:
```bash
php sprint database refresh app
```

Create administrator:
```bash
php sprint database seed UserSeeder
```
This command create administrator with data:

**Login:** admin
**Password:** fpwPOuZD

Do not forget change password!

Add default Games and Game Types:
```bash
php sprint database seed GamesSeeder
php sprint database seed GameTypesSeeder
```
