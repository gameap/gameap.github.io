---
title: Rust
layout: default
lang: ru
category: Туториалы
order: 202
---

Rust - игра на выживание, разработанная Facepunch Studios. 
Игрокам предстоит собирать ресурсы, изготавливать предметы и строить убежища, 
чтобы выжить во враждебной среде открытого мира.

Rust в некотором смысле похожа на [Minecraft](/ru/tutorials/minecraft.html) 
системой выживания и крафтом предметов.

## Настройка окружения

GameAP обеспечивает полную поддержку игрового сервера Rust, 
включая управление игроками. 
Установка панели займет немного времени, 
выберите один из следующих вариантов установки панели управления:

* [Установка GameAP на Linux](/ru/install/install_on_linux.html)
* [Установка GameAP на Windows](/ru/install/install_on_windows.html)

### Установка GameAP Daemon

Для установки игрового сервера на машине (VDS) необходимо 
установить GameAP Daemon. При установке GameAP вы можете выбрать 
полную установку, включая Daemon.

В панели управления перейдите на страницу **"Администрирование"** -> 
**"Выделенные серверы"** -> **"Создать"**. 
Появится окно с предложением автоматической установки. 
Скопируйте код и выполните его на выделенном сервере.

После этого можно приступить к установке сервера Rust.

## Установка Rust сервера в GameAP

Перейдите на страницу **Администрирование** -> **Игровые серверы** -> **Создать**

![](/images/ru/tutorials/rust/create_form.png)

* В поле "Имя" впишите любое название сервера, к примеру "Мой Rust сервер".
* В поле "Игра" выберите из списка Rust.
* В поле модификация выберите модификацию, по умолчанию Vanilla.
* В поле "Выделенный сервер" выберите желаемую ноду, на которой будет располагаться игровой сервер.
* В поле IP выберите желаемый адрес вашего сервера, затем можете выбрать свободный порт вашего сервера, либо
  использовать предложенный.

## Настройка игрового сервера

Для изменения конфигурации сервера, перейдите в раздел **Сервер**, выберите ваш сервер и нажмите **Управление**.
Далее откройте вкладку **Настройки**.

![](/images/ru/tutorials/rust/settings.png)

Чтобы настройки вступили в силу вам необходимо перезагрузить сервер.

### Server Hostname

Это название вашего сервера Rust, которое будет отображаться у всех игроков в игре
в окне поиска серверов.

### Maximum players on server

Максимальное количество игроков, которое может одновременно играть на сервере.
По умолчанию 32.

### Map on the server

Карта на сервере. По умолчанию процедурно генерируемая карта Procedural Map
Возможные значения:

* Barren
* Craggy Island
* Hapis
* Savas Island

### Map generation seed

Зерно для генерации карты Procedural Map.

### World size

Размер мира.

### Server save interval

Интервал сохранения мира.
Это период в течение которого данные сервера будут сохраняться на диск.