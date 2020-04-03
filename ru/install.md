---
title: Установка панели
layout: default
lang: ru
category: Общее
order: 20
---

* This will become a table of contents (this text will be scraped).
{:toc}

Установите панель одним из следующих способов

* [Автоматическая установка](/ru/auto_install.html). Выполняется за короткое время скриптом автоустановки.
* [Установка на Shared хостинг](/ru/shared_install.md). Установка на обычный хостинг копированием файлов, без доступа по SSH.
* [Ручная установка](/ru/manual_install.html). Для опытных пользователей.


## Настройка веб-сервера

### Nginx

Пример конфигурации nginx:

```
server {
    listen       80;
    server_name  _;

    #charset koi8-r;
    #access_log  /var/log/nginx/log/host.access.log  main;
    root /var/www/gameap/public;
    index index.php index.html;

    location / {
        try_files $uri $uri/ /index.php$is_args$args;

        location = /index.php
        {
            # fastcgi_pass   php:9000;
            fastcgi_pass    unix:/var/run/php/php7.2-fpm.sock;
            fastcgi_param  SCRIPT_FILENAME $document_root$fastcgi_script_name;
            include        fastcgi_params;
        }
    }

    #error_page  404              /404.html;

    # redirect server error pages to the static page /50x.html
    #
    error_page   500 502 503 504  /50x.html;
    location = /50x.html {
        root   /usr/share/nginx/html;
    }
}
```
Обратите внимание на параметры `server_name` и `root`, скорее всего у вас они другие и их нужно указать.
В `server_name` указывается хост сервера, например gameap.ru.
В `root` указывается путь к public файлам панели (там, где index.php)

### Apache

Пример конфигурации хоста Apache:

```
NameVirtualHost *:80
Listen 80
 
<VirtualHost *:80>
    ServerName example.com
    DocumentRoot /var/www/gameap/public
     
    <Directory /var/www/gameap/public>
            Options Indexes FollowSymLinks MultiViews
            AllowOverride All
            Order allow,deny
            allow from all
            Require all granted
    </Directory>
     
    LogLevel debug
    ErrorLog ${APACHE_LOG_DIR}/error.log
    CustomLog ${APACHE_LOG_DIR}/access.log combined
</VirtualHost>
```

Обратите внимание на значения `DocumentRoot` и `Directory`, они должны вести в public каталог панели, а не в её корень.
Например `/var/www/gameap/public`.

## Настройка базы данных

### MySQL

После установки MySQL создайте базу данных `gameap`, откройте `.env` файл и отредактируйте параметры:
```
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=gameap
DB_USERNAME=gameap
DB_PASSWORD=secret
```

После найстройки параметров перейдите в каталог с панелью управления `/var/www/gameap` и обновите схему базы данных:
```bash
php artisan migrate --seed
```

### SQLite

Создайте пустой файл базы данных:
```
touch /var/www/gameap/database.sqlite
```

Откройте `.env` и отредактируйте параметры базы данных:
```
DB_CONNECTION=sqlite
DB_DATABASE=/var/www/gameap/database.sqlite
```

После найстройки параметров перейдите в каталог с панелью управления `/var/www/gameap` и обновите схему базы данных:
```bash
php artisan migrate --seed
```
