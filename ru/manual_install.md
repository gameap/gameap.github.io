---
title: Ручная установка панели
layout: default
lang: ru
category: Общее
order: 30
---

# Ручная установка панели

Установка GameAP 2.x с GitHub.

Примечания!
1. Этот мануал больше актуален для Debian 8. Однако, для Ubuntu/CentOS и др дистрибутивов это тоже должно работать, но
с некоторыми изменениями.
2. Этот мануал не актуален для Shared хостингов.

Нам понадобится:
* ОС: Debian 8
* Версия PHP: 7.1
* Веб-сервер: Nginx
* База данных: MySQL
* Консольные утилиты: [git](requirements.html#git), [composer](requirements.html#composer)

Установите необходимые расширения PHP:
```bash
apt-get -y install php7.1-cli php7.1-fpm php7.1-pdo php7.1-mysql php7.1-redis php7.1-gd php7.1-mcrypt php7.1-curl php7.1-bz2 php7.1-xml php7.1-mbstring php7.1-bcmath
```

Создайте каталог с gameap, например /var/www/gameap

```bash
git clone https://github.com/et-nik/GameAP /var/www/gameap
```
Обратите внимание, что если у вас другой каталог (не /var/www/gameap), то указывайте свой.

Установите необходимые зависимости:
```bash
composer install
```

Создаём базу данных, редактируем application/config/database.php,
указываем там свои данные БД.

Обновляем схему БД:
```bash
php sprint database refresh app
```

Создаём администратора:
```bash
php sprint database seed UserSeeder
```
Эта команда создаст администратора со следующими данными:

Логин: admin
Пароль: fpwPOuZD

Не забудьте сменить пароль!

Добавляем стандартные игры и модификации:
```bash
php sprint database seed GamesSeeder
php sprint database seed GameTypesSeeder
```

---
Документация в процессе написания!