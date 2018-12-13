---
title: Ошибки панели
layout: default
lang: ru
category: Общее
order: 40
---

* This will become a table of contents (this text will be scraped).
{:toc}

Описание некоторых возможных ошибок и их исправление.

## Ошибки после установки

### Whoops, looks like something went wrong.

#### Генерация ключа

![](/images/errors/key_generate.png)

Эта проблема возникнет если вы забыли сгенерировать ключ. Перейдите в каталог с панелью и выполните команду:

```
php artisan key:generate --force
```

#### Миграция базы данных

![](/images/errors/db_migrate.png)

Эта проблема возникает, если вы забыли выполнить команду миграции базы данных после установки панели. Перейдите в каталог с панелью и выполните команду:

```
php artisan migrate --seed
```