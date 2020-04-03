---
title: FTP Manager
layout: default
lang: ru
category: Модули
order: 100
---

Модуль управления FTP аккаунтами

[https://github.com/gameap/ftp-module](https://github.com/gameap/ftp-module)

Модуль поддерживает создание, удаление, обновление аккаунтов FTP на выделенных серверах. После создания FTP аккаунта можно подключиться к нему например через FileZilla или другой клиент.


# Установка модуля

## При наличии доступа к консоли (на VDS)

### Используя Composer

Перейдите в каталог с панелью (например /var/www/gameap):
```
cd /var/www/gameap
```

После этого установите модуль и выполните миграцию:
```
composer require --update-no-dev gameap/ftp-module
php artisan module:migrate Ftp
```

#### Если возникла ошибка "bash: composer: command not found"

Ошибка возникает если [composer](https://getcomposer.org/) не установлен. Установите его следующей командой:

```bash
curl -sS https://getcomposer.org/installer | sudo php -- \
    --install-dir=/usr/local/bin \
    --filename=composer
```

После установки попробуйте выполнить команду установки вновь.

### Используя Git

Перейдите в каталог с панелью (по умолчанию `/var/www/gameap`):
```
cd /var/www/gameap
```

Выполните команду:
```
git clone https://github.com/gameap/fastdl-module modules/Ftp
```

Выполните миграцию базы данных:
```
php artisan module:migrate Ftp
```

#### Если возникла ошибка "bash: git: command not found"

Необходимо установить утилиту git. В зависимости от вашего дистрибутива, выполните команду.

Debian/Ubuntu:
``` 
apt install git
```

CentOS:
```bash
yum install git
```

## При отсутствии доступа к консоли (на Shared хостинг)

Скопируйте содержимое архива с модулем в каталог `modules`. Архив можно скачать [здесь](https://github.com/gameap/ftp-module/archive/master.zip)

Перейдите в панель управления, там выберите в верхнем меню **"GameAP"** -> **"Модули"**, затем кликните **"Запустить миграцию"**

# Настройка модуля

Также, нужно настроить модуль, задать команды, которыми будут управляться FTP аккаунты.

С главной страницы панели перейдите в **"FTP"** -> **"Команды"**, задайте команды как это указано в примере:

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