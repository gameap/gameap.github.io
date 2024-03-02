---
title: Panel installation
layout: default
lang: en
category: Main
order: 21
hidden: true
---

* This will become a table of contents (this text will be scraped).
{:toc}

Install GameAP in one of the following ways:

* [Automatic Installation](/en/auto_install.html). Fast way. Using a bash script that will do everything itself.
* [Installation on Shared hosting](/en/shared_install.html). Installation on a shared hosting by copying files. If you do not have SSH. Or on Windows.
* [Manual Installation](/en/manual_install.html). For advanced users.


## Web-server setup

### Nginx

Example Nginx virtual host config:

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
            fastcgi_pass    unix:/var/run/php/php7.3-fpm.sock;
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
Pay attention to the `server_name` and` root` parameters, most likely you have value different and you need to specify them.
In `server_name` specifies the server host, for example gameap.ru.
The `root` is the path to the public panel files (where index.php is)


### Apache

Example Apache virtual host config:

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

## Database setup

### MySQL

After installing MySQL, create the `gameap` database, open the` .env` file and edit the parameters:
```
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=gameap
DB_USERNAME=gameap
DB_PASSWORD=secret
```

Go to the directory with the control panel `/var/www/gameap` and migrate the database schema:
```bash
php artisan migrate --seed
```

### SQLite

Create an empty database file:
```
touch /var/www/gameap/database.sqlite
```

Open `.env` and edit database parameters:
```
DB_CONNECTION=sqlite
DB_DATABASE=/var/www/gameap/database.sqlite
```

Go to the directory with the control panel `/var/www/gameap` and migrate the database schema:
```bash
php artisan migrate --seed
```
