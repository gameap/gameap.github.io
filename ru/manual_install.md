---
title: Ручная установка панели
layout: default
lang: ru
category: Общее
order: 30
hidden: true
---

* This will become a table of contents (this text will be scraped).
{:toc}

>! В большинстве случаев использования данный способ не рекомендуется использовать!
> Используйте этот способ, только если нет возможности установить рекомендуемыми способами:
> * [Установка на Linux](/ru/install/install_on_linux.html)
> * [Установка на Windows](/ru/install/install_on_windows.html)

Установка GameAP 3.x с GitHub. Это подробная инструкция по установке GameAP из GitHub, установке необходимых зависимостей и
сборке стилей.

> Примечания!
1. Этот мануал больше актуален для Debian. Однако, для Ubuntu/CentOS и др дистрибутивов это тоже должно работать, но
с некоторыми изменениями.
2. Этот мануал не актуален для Shared хостингов.

Нам понадобится:
* ОС: Debian 9
* Версия PHP: >= 7.3
* Веб-сервер: Nginx
* База данных: MySQL
* Консольные утилиты: [git](requirements.html#git), [composer](requirements.html#composer)

## Установка необходимых пакетов

Установите необходимые пакеты, которые в дальнейшем могут понадобится. Обычно, они уже должны быть установлены, но бывает, что их нет и
из-за этого могут вознкнуть трудности в выполнении дальнейших инструкций:
```bash
sudo apt-get -y install 
    wget \
    software-properties-common \
    ca-certificates \
    apt-transport-https \
    gnupg \
    curl \
    lsb-release
```

Установите PHP и необходимые расширения PHP:
```bash
sudo apt-get -y install \
    php7.3-common \
    php7.3-cli \
    php7.3-fpm \
    php7.3-pdo \
    php7.3-mysql \
    php7.3-redis \
    php7.3-curl \
    php7.3-bz2 \
    php7.3-zip \
    php7.3-xml \
    php7.3-mbstring \
    php7.3-bcmath \
    php7.3-gmp \
    php7.3-intl
```

Установите composer:
```bash
curl -sS https://getcomposer.org/installer | sudo php -- \
    --install-dir=/usr/local/bin \
    --filename=composer
```

> Если после выполнения команд вы получаете слудующие ошибки:
```
Reading package lists... Done
Building dependency tree       
Reading state information... Done
E: Unable to locate package php7.3-common
E: Couldn't find any package by glob 'php7.3-common'
E: Couldn't find any package by regex 'php7.3-common'
E: Unable to locate package php7.3-cli
E: Couldn't find any package by glob 'php7.3-cli'
E: Couldn't find any package by regex 'php7.3-cli'
...
```
Для Debian выполните следующие команды:
```
wget -O - https://packages.sury.org/php/apt.gpg | sudo apt-key add -
echo "deb https://packages.sury.org/php/ $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/php.list
sudo apt-get update
```

Установите базу данных и веб сервер. Можете использовать MySQL и Nginx:

```
sudo apt-get install mysql-server nginx
```

О настройке веб-сервера здесь будет опущено. Воспользуйтесь одной и тысяч инструкций по настройке в интернете =)

> Обратите внимание, что для тестирования панель не обязательно использовать с MySQL и Nginx, можно воспользоваться SQLite и встроенным веб сервером PHP. 
Более подробно об этом прочитайте в [Настройке простого окружения](#простое-окружение)

## Установка и настройка панели

После того, как установили PHP и все необходимые расширения, склонируйте gameap в рабочий каталог, например в `/var/www/gameap`

```bash
git clone https://github.com/et-nik/gameap /var/www/gameap
```
Обратите внимание, что если у вас другой каталог (не `/var/www/gameap`), то указывайте свой.

Дальше будем работаеть в этом каталоге, переходим в него:

```bash
cd /var/www/gameap
```

Основная конфигурация панели находится в файле `.env`. Создайте свой файл конфигурации из примера:
```bash
cp .env.example .env
```

Установите необходимые зависимости:
```bash
composer install --no-dev --optimize-autoloader
```
Создайте базу данных, после этого отредактируйте параметры базы данных в файле `.env`, укажите хост, логин, пароль, название базы данных и т.д.

Сгенерируйте ключ:
```bash
php artisan key:generate --force
```

Обновите схему БД:
```bash
php artisan migrate --seed
```

Эта команда создаст необходимые таблицы в базе данных, добавит администратора и добавит некоторые данные для примера.
Перейдите на URL страницу с панелью управления.

Данные администратора по умолчанию:
Логин: admin
Пароль: fpwPOuZD

Не забудьте сменить пароль!

## Сборка стилей

Установите NPM:
```
curl -sL https://deb.nodesource.com/setup_10.x | bash -
apt-get install -y nodejs
```

Установите необходимые зависимости:
```
npm install
```

Сборка стилей:
```
npm run prod
```

Панель настроена, все необходимые зависимости установлены, все стили и js библиотеки собраны, можно приступать к использованию. 

> Не забудьте, что на выделенном сервере, на котором планируется размещать игровые серверы нужно установить и настроить GameAP Daemon, об этом можете прочитать здесь -- [Установка GameAP Daemon](/ru/gameap_daemon.html).

## Cron

Добавьте cron задание

```
* * * * * cd /patch/to/gameap && php artisan schedule:run >> /dev/null 2>&1
```

Вместо `/patch/to/gameap` укажите свой путь к панели управления.

## Простое окружение

Если вам нужна панель для того, чтобы протестировать что-то, то можно воспользоваться SQLite и встроенным веб-сервером PHP. Не обязательно устанавливать базу данных MySQL (или любую другую) и веб сервер Nginx (или Apache).

Создайте файл базы данных:
```
touch /var/www/gameap/database.sqlite
```

Откройте `.env` и отредактируйте параметры базы данных:
```
DB_CONNECTION=sqlite
DB_DATABASE=/var/www/gameap/database.sqlite
```

Обновите схему базы данных:
```bash
php artisan migrate --seed
```

Запустите Web-сервер
```
cd public
php -S localhost:8080
```

Перейдите по адресу `http://localhost:8080`
Введите логин и пароль администратора.

---
Документация в процессе написания!
