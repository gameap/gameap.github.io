---
title: FastDL Manager
layout: default
lang: ru
category: Модули
order: 101
---

Модуль управления FastDL для Source и GoldSource серверов (Counter-Strike 1.6, Half-Life, Team Fortress 2, 
Counter-Strike: Global Offensive и др.)

[https://github.com/gameap/fastdl-module](https://github.com/gameap/fastdl-module)

Модуль поддерживает создание и удаление FastDL аккаунтов для игровых серверов, а также
автоматическую установку необходимых зависимостей и их настройку (Nginx, Rsync).

# Установка модуля

## При наличии доступа к консоли (на VDS)

Перейдите в каталог с панелью (например `/var/www/gameap`, это путь по умолчанию):
```
cd /var/www/gameap
```

После этого установите модуль и выполните миграцию:
```
composer require --update-no-dev gameap/fastdl-module
php artisan module:migrate Fastdl
```

## При отсутствии доступа к консоли (на Shared хостинг)

Скопируйте содержимое архива с модулем в каталог `modules`. Архив можно скачать [здесь](https://github.com/gameap/fastdl-module/archive/master.zip)

Перейдите в панель управления, там выберите в верхнем меню **"GameAP"** -> **"Модули"**, затем кликните **"Запустить миграцию"**

# Настройка модуля

Для работы FastDL у вас должен быть уже установлен хоть один выделенный сервер, а 
на нём должен иметься хоть один игровой сервер GoldSource/Source.

С главной страницы панели перейдите в **"FastDL"**, найдите в списке выделенный сервер, на котором вы хотите установить
FastDL окружение, кликните **"Редактировать"**, перед вами откроется форма, которую нужно заполнить.

## Метод

Это метод, который будет использоваться для создания FastDL аккаунтов

| Метод  | Описание                                  |
| ------ | ----------------------------------------- |
| rsync  | Рекомендуется. Файлы сервера будут синхронизированы с файлами в веб каталоге используя rsync
| copy   | Файлы игрового сервера будут скопированы в веб каталог
| link   | Будет создана ссылка на игровой сервер в веб-каталоге
| mount  | Каталог игрового сервера будет смонтирован (`mount --bind`) в веб-каталог 

## Хост

IP адрес или домен по которому будет доступен FastDL

## Порт

Порт веб сервера для FastDL. По умолчанию 80

## Автоиндекс

При переходе по адресу FastDL вы увидите список всех доступных файлов. 
Если отключить опцию, то вместо списка файлов вы увидите 403 ошибку, однако файлы будут доступны по прямой ссылке.

## Дополнительные опции

Вы можете добавить дополнительные опции. С их помощью можно более гибко настроить FastDL. 
Например можно настроить удалённый FastDL.

| Опция         | Описание                                  |
| ------------- | ----------------------------------------- |
| web-path      | Укажите свой путь к web-каталогу
| rsync-remote  | Удалённый FastDL