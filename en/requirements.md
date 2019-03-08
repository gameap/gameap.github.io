---
title: Requirements
layout: default
lang: en
category: Main
order: 10
---

* This will become a table of contents (this text will be scraped).
{:toc}

## Server part

* GameAP Daemon 2.1
* GameAP Starter

## Web part

* Any web server (nginx, lighttpd, apache etc.)
* PHP >= 7.1
* PHP Extensions: GD, OpenSSL, Curl, GMP, Intl

### Composer

Optional. Required for manual install.

#### Install

##### Linux
```bash
wget https://getcomposer.org/download/1.8.0/composer.phar
cp composer.phar /usr/bin/composer
chmod +x /usr/bin/composer
```

### Git

Optional. Required for manual install.

#### Install

##### Debian/Ubuntu

```bash
apt-get install git
```

##### CentOS
```bash
yum install git
```
