---
title: Установка на Shared хостинг
layout: default
lang: ru
category: Общее
order: 31
---

>! Не рекомендуется использовать Shared способ. Некоторый функционал панели будет ограничен. При наличии возможности не использовать Shared хостинг -- воспользуйтесь этой возможностью.

* This will become a table of contents (this text will be scraped).
{:toc}

Установка GameAP на Shared хостинг, на котором нет доступа к командной строке.

* **Рекомендуемая версия PHP** >= 7.2

## Загрузка архива с панелью

Для начала вам нужно загрузить архив с панелью.
Скачайте следующий архив:
> [gameap_latest.zip](http://www.gameap.ru/gameap_latest.zip)

Распакуйте его в вашу директорию на хостинге.

## Импорт данных в базу

Имортируйте SQL файл `gameap.sql` в базу данных. Сделать это можно через PHPMyAdmin, большинство хостингов предоставляют его.

## Настройка панели управления

Скопируйте файл `.env.example` под новым именем `.env`. 
Откройте файл `.env`, задайте значения для параметра `APP_KEY`. Укажите параметры базы данных.

## Настройка веб сервера на хостинге

Вы должны указать в настройках так, чтобы корневой каталог сайта вёл в публичную директорию `public`. Например, если на хостинге файлы панели расположены в `/home/your_site`, то корневой каталог сайта вы должны указать `/home/your_site/public`.

Все запросы должны обрабатываться фронт контроллером `index.php`, запросы к несуществующим каталогам тоже (например к `example.com/login`, `example.com/home` и т.д.). Фронт контроллер не должен обрабатывать запросы загрузки файлов css, js.

Добавьте новый домен в панели управления, каталог укажите как `/path/to/gameap/public`. Для ISPManager это будет примерно так:

![](/images/install/ispmanager_add_domain.png)

Если указывать каталог выше public (`/www/your_site`), то панель будет работать, однако так делать не рекомендуется.

Не рекомендуется, чтобы корневой каталог сайта вёл в корневой каталог с панелью (туда, где файлы `.env`, `composer.json`, каталоги `app`, `vendor`). Корневой каталог сайта должен вести в `public` по разным причинам:

* Безопасность. В случае неправильной настройки веб сервера может произойти такое, что кто-то прочитает важные файлы, например файл `.env`, который содержит логин и пароль для доступа к базе данных.
* Быстродействие.

Некоторые хостинги не позволяют это настроить в своей панели управления. Вам нужно обратиться в техническую поддержку для ручной настройки.

Пример конфигурации веб серверов смотрите в разделе [Установка/Настройка веб сервера](/ru/install.html#настройка-веб-сервера)

## Данные для входа

После настройки панели, импорта данных и настройки веб-сервера ваша панель готова для работы. Зайдите на главную страницу вашего сайта, перед вами должна появиться форма для входа.

Данный администратора для входа:

Логин: admin
Пароль: fpwPOuZD

> ! Не забудьте сменить пароль после входа