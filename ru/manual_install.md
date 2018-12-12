---
title: Ручная установка панели
layout: default
lang: ru
category: Общее
order: 30
---

* This will become a table of contents (this text will be scraped).
{:toc}

Установка GameAP 3.x с GitHub. Это подробная инструкция по установке GameAP из GitHub, установке необходимых зависимостей и
сборке стилей.

> Примечания!
1. Этот мануал больше актуален для Debian 8. Однако, для Ubuntu/CentOS и др дистрибутивов это тоже должно работать, но
с некоторыми изменениями.
2. Этот мануал не актуален для Shared хостингов.

Нам понадобится:
* ОС: Debian 8
* Версия PHP: 7.1
* Веб-сервер: Nginx
* База данных: MySQL
* Консольные утилиты: [git](requirements.html#git), [composer](requirements.html#composer)

## Установка необходимых пакетов

Установите необходимые пакеты, которые в дальнейшем могут понадобится. Обычно, они уже должны быть установлены, но бывает, что их нет и
из-за этого могут вознкнуть трудности в выполнении дальнейших инструкций:
```bash
sudo apt-get install -y software-properties-common apt-transport-https gnupg curl
```

Установите PHP и необходимые расширения PHP:
```bash
sudo apt-get -y install php7.1-common php7.1-cli php7.1-fpm php7.1-pdo php7.1-mysql php7.1-redis php7.1-curl php7.1-bz2 php7.1-zip php7.1-xml php7.1-mbstring php7.1-bcmath
```

> Если после выполнения команд вы получаете слудующие ошибки:
```
Reading package lists... Done
Building dependency tree       
Reading state information... Done
E: Unable to locate package php7.1-common
E: Couldn't find any package by glob 'php7.1-common'
E: Couldn't find any package by regex 'php7.1-common'
E: Unable to locate package php7.1-cli
E: Couldn't find any package by glob 'php7.1-cli'
E: Couldn't find any package by regex 'php7.1-cli'
...
```
Тогда вам необходимо добавить новый репозиторий в менеджер пакетов. О том, как это сделать можете прочитать например здесь [https://sapeg.in/2018/06/847/](), либо
тут [https://neoblogs.net/2017/03/22/ustanovka-php7-1-na-ubuntu-nastrojka-php7-1/]().

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

Сгенерируйте ключ:
```bash
php artisan key:generate --force
```

Создайте базу данных, после этого отредактируйте параметры базы данных в файле `.env`, укажите хост, логин, пароль, название базы данных и т.д.

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

## Компиляция стилей

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