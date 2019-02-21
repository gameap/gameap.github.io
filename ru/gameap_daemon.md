---
title: GameAP Daemon
layout: default
lang: ru
category: GameAP Daemon
order: 50
---

* This will become a table of contents (this text will be scraped).
{:toc}

GameAP Daemon — фоновое приложение, которое обменивается с панелью данными, работает с игровыми серверами
(устанавливает, удаляет, запускает, останавливает и т.д).

Это основное приложение, которое будет контролировать статус игровых серверов, перезапускать их в случае надобности или
по требованию.

## Установка

### Debian/Ubuntu

Следующая инструкция актуальна для Debian/Ubuntu

### Установка необходимых пакетов

Установите необходимые пакеты, которые потом могут понадобиться:

```
sudo apt-get update
sudo apt-get install -y wget curl gnupg ca-certificates apt-transport-https lsb-release
```

После этого необходимо добавить репозиторий с GameAP Daemon.

Скачайте и добавьте ключ:
```bash
wget -O - http://packages.gameap.ru/gameap-rep.gpg.key | apt-key add -
```

Добавьте репозиторий:
```bash
echo "deb http://packages.gameap.ru/debian/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/gameap.list
```

Обновите список пакетов и запустите установку:
```bash
apt-get update
apt-get install gameap-daemon
```

## Конфигурация

Конфигурация GameAP Daemon находится в файле daemon.cfg

### Базовые параметры

| Параметр                  | Обязателен            | Тип       | Информация
|---------------------------|-----------------------|-----------|------------
| ds_id                     | да                    | integer   | ID выделенного сервера
| listen_port               | нет (по умолчанию 31717) |     integer   | Порт, который будет слушать
| api_host                  | да                    | string    | API хост
| api_key                   | да                    | string    | API ключ


### SSL/TLS

| Параметр                  | Обязателен            | Тип       | Информация
|---------------------------|-----------------------|-----------|------------
| client_certificate_file   | да                    | string    | Сертификат клиента (панели)
| certificate_chain_file    | да                    | string    | Сертификат сервера
| private_key_file          | да                    | string    | Приватный ключ сервера
| private_key_password      | нет                   | string    | Пароль от ключа сервера
| dh_file                   | да                    | string    | Сертификат Ди́ффи — Хе́ллмана

[Создание сертификатов](https://github.com/gameap/GDaemon2#creating-certificates)

### Аутентификация по логину и паролю

| Параметр                  | Обязателен            | Тип       | Информация
|---------------------------|-----------------------|-----------|------------
| password_authentication   | нет                   | boolean   | Включить аутентификацию по логину и паролю
| daemon_login              | нет                   | string    | Логин. В Linux, если пуст или не задан, то будет использоваться PAM.
| daemon_password           | нет                   | string    | Password. В Linux, если пуст или не задан, то будет использоваться PAM.

### Статистика

| Параметр                  | Обязателен            | Тип       | Информация
|---------------------------|-----------------------|-----------|------------
| if_list                   | нет                   | string    | Список интерфейсов
| drives_list               | нет                   | string    | Список дисков
| stats_update_period       | нет                   | integer   | Период обновления статистики
| stats_db_update_period    | нет                   | integer   | Период обновления базы данных


[Пример конфигурации](https://github.com/gameap/GDaemon2#example-daemoncfg)
