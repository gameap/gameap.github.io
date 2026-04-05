---
title: Установка на Windows
layout: default
lang: ru
category: Установка GameAP
order: 101
---

## Поддерживаемые версии

| Версия      | Поддерживается | 
|-------------|----------------|
| Server 2025 | ✔              |
| Server 2022 | ✔              |
| Server 2022 | ✔              |
| Server 2019 | ✔              |
| 11          | ✔              |
| 10          | ✔              |


## Загрузка GameAP Control

Необходимо скачать утилиту GameAP Control (gameapctl) для управления окружением GameAP.

Для этого необходимо перейти на страницу релизов gameapctl в Github: 
[https://github.com/gameap/gameapctl/releases](https://github.com/gameap/gameapctl/releases)

Выберите там последний релиз и кликните по нему

<video width="1280" height="720" controls>
  <source src="/images/en/gameapctl/download.webm" type="video/webm">
  Your browser does not support the video tag.
</video>

После этого найдите версию подходящую вам. Наиболее популярной является архитектура
Windows AMD64, поэтому скорее всего вам необходимо скачать именно этот архив:

![](/images/en/gameapctl/download_release_windows_amd64.png)

## Установка панели с использованием GameAP Control UI

После того, как вы скачали архив с утилитой gameapctl, запустите его.

Откроется окно браузера, в котором в разделе Web/API необходимо нажать "Install".

<video width="1280" height="720" controls>
  <source src="/images/en/gameapctl/windows-install.webm" type="video/webm">
  Your browser does not support the video tag.
</video>

### Параметры установки 

Укажите необходимые данные для установки

![](/images/en/gameapctl/ui_gameap_installation.png)

#### Хост

Укажите домен или IP адрес, по которому будет доступна панель.

В случае IP адреса, должен быть указан адрес, который назначен сетевому интерфейсу на VDS.
Если ваша сеть использует NAT, то не указывайте внешний IP,
а указывайте внутренний, затем сконфигурируйте перенаправление портов.

Домен может быть указан любой, но не забудьте настроить DNS.

Пример корректных значений:
* 10.182.104.8
* 10.182.104.8:2080
* example.com
* http://example.com

#### База данных

База данных, где будут храниться данные: пользователи, сведения о серверах и тд.
Можно использовать:
* [PostgreSQL](https://www.postgresql.org/). Рекомендуется для больших проектов, с большим количеством игровых серверов 
  и пользователей.
* [MySQL](https://www.mysql.com/)/[MariaDB](https://mariadb.org/)
* [SQLite](https://www.sqlite.org/). Если нагрузка на ваш сервер планируется небольшой и вы не планируете использовать более 10 игровых серверов.

#### Установка GameAP Daemon

Помимо самой Web панели вы можете ещё установить серверную часть GameAP Daemon, которая
отвечает за работу с игровыми серверами.

### Окончание установки

Дождитесь окончания установки. Некоторые этапы 
могут занимать продолжительное время.

После окончания не забудьте сохранить данные для входа 
и данные от базы, которые будут указаны.

![](/images/en/gameapctl/gameap_finished_installation.png)