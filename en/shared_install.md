---
title: Install on shared hosting
layout: default
lang: en
category: Main
order: 31
hidden: true
---

>! Documentation in the process of writing. This section requires more material.

>! t is not recommended to use the Shared method. Some panel functionality will be limited

* This will become a table of contents (this text will be scraped).
{:toc}

Installing GameAP on a shared hosting with no access to the command line.

* **Minimum PHP version**: 7.3
* **Recommended PHP version**: >= 7.3

## Downloading panel archive

Download the following archive:
> [gameap_latest.zip](http://www.gameap.ru/gameap_latest.zip)

Unpack archive into your directory on hosting.

## Import data into database

Import the `gameap.sql` SQL file into the database. This can be done through PHPMyAdmin, most hosting provide it.

## GameAP Setting Up

* Copy `.env.example` under the new name `.env`.
* Open the file `.env`, set `APP_KEY` value. You can use next random generated value:

<pre id="app-key"></pre>

<a href="#" onclick="generateAndPrintKey(); return false;">Generate new key</a>

<script>
function randomKey(length) {
   var result           = '';
   var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789\\/!@#$%^&*()_-+=';
   var charactersLength = characters.length;
   for ( var i = 0; i < length; i++ ) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
   }
   return result;
}

function generateAndPrintKey() {
    document.getElementById('app-key').innerHTML = "base64:" + btoa(randomKey(32));
}

generateAndPrintKey();

</script>

* Set database paramenters (host, database name, login, password).

## Setting up a web server on a hosting

You must specify in the settings so that the root directory of the site follow to the public directory `public`. For example, if the hosting files for the panel are located in `/home/your_site`, then you must specify `/home/your_site/public` in the root directory of the site.

Все запросы должны обрабатываться фронт контроллером `index.php`, запросы к несуществующим каталогам тоже (например к `example.com/login`, `example.com/home` и т.д.). Фронт контроллер не должен обрабатывать запросы загрузки файлов css, js.

All requests should be processed by the front controller `index.php`, requests to non-existing directories too (for example, to `example.com/login`, `example.com/home`, etc.). The front controller should not process requests to download files css, js.

If you specify the directory above public (`/www/your_site`), then the panel will work, but this is not recommended.

It is not recommended that the root directory of the site follow to the root directory with the panel (where the files are `.env`,` composer.json`, directories `app`, `vendor`). The site root directory should follow to `public` for various reasons:

* Security. If the web server is configured incorrectly, it may happen that someone reads important files, such as the `.env` file, which contains the login and password for accessing the database.
* High speed.

Many hostings do not allow this to be configured in their control panel. You need to contact technical support for manual configuration.

For an example of web server configuration, see [Installation/Web Server Setup](/ru/install.html#настройка-веб-сервера)

## Login Details

After setting up the panel, importing data and setting up a web server, your panel is ready for work. Go to the main page of your site, you should see a login form.

This admin to log in:

Login: admin
Password: fpwPOuZD

> ! Do not forget to change your password after logging in.
