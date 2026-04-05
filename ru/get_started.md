---
title: Начало работы
layout: default
lang: ru
category: Общее
order: 2
---

Для начала работы желательно два выделенных или виртуальных сервера. 
На одном размещается панель управления, на другом игровые серверы. 
Можно всё устанавливать на одном выделенном сервере.

## Установка панели

Панель устанавливается на выделенный сервер с базой данных (PostgreSQL, MySQL, SQLite).

&nbsp;&nbsp;&nbsp;&nbsp;<img src="/images/icons/linux.svg" alt="Linux" width="20" height="20" style="vertical-align: middle"> [Установка на Linux](/ru/install/install_on_linux.html)

&nbsp;&nbsp;&nbsp;&nbsp;<img src="/images/icons/windows.svg" alt="Windows" width="20" height="20" style="vertical-align: middle"> [Установка на Windows](/ru/install/install_on_windows.html)

## Простейшая установка на Linux

Если у вас Linux, установлен CURL и вы не хотите разбираться с подробностями установки, то выполните команду:
```bash
bash <(curl -s https://gameap.com/install.sh) --with-daemon
```

## Добавление выделенного сервера

Добавьте новый выделенный сервер (VDS), на который затем будете устанавливать игровые серверы. 
После установки панели зайдите в неё и в меню выберите **"Администрирование"** **"Выделенные серверы"** → **"Создать"**. После
чего откроется окошко с инструкцией, следуйте ей.

![](/images/ru/get_started/add_dedicated_server.gif)

Более подробно об установке и настройке читайте на странице [Выделенные серверы](/ru/gameap_configure/dedicated_servers.html).

## Добавление игрового сервера

Перейдите в **"Администрирование"** → **"Игровые серверы"** → **"Создать"**.

![](/images/ru/get_started/add_game_server.gif)

Более подробно о параметрах настройки игровых серверов читайте на странице [Игровые серверы](/ru/gameap_configure/game_servers.html)