---
title: FTP Manager
layout: default
lang: ru
category: Модули
order: 501
---

Модуль управления FTP аккаунтами

[https://github.com/gameap/ftp-module](https://github.com/gameap/ftp-module)

Модуль поддерживает создание, удаление, обновление аккаунтов FTP на выделенных серверах. После создания FTP аккаунта можно подключиться к нему например через FileZilla или другой клиент.


# Установка модуля

Зайдите в панель, выберите **Модули** → **Marketplace**

Найдите в списке FTP модуль и нажмите "Установить." 

# Настройка модуля

Также, нужно настроить модуль, задать команды, которыми будут управляться FTP аккаунты.

С главной страницы панели перейдите в **"FTP"** → **"Команды"**, задайте команды как это указано в примере:

Команда создания: 
`./ftp.sh add --username="{username}" --password="{password}" --directory="{dir}"`

Команда обновления: 
`./ftp.sh update --username="{username}" --password="{password}" --directory="{dir}"`

Команда удаления: 
`./ftp.sh delete --username="{username}"`

## Настройка выделенного сервера (ноды)

Скачайте `ftp.sh` [отсюда](https://github.com/gameap/scripts/tree/master/ftp) и поместите его в рабочий каталог на выделенном сервере (по умолчанию `/srv/gameap`) и установите ему права на выполнение:
```
chmod +x /srv/gameap/ftp.sh
```
