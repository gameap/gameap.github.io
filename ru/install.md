---
title: Установка панели
layout: default
lang: ru
category: Общее
order: 20
---

* This will become a table of contents (this text will be scraped).
{:toc}

## Debian

### Добавление репозиториев

Скачайте и добавьте GPG ключ:
```bash
wget -O - http://packages.gameap.ru/debian/gameap-rep.gpg.key | sudo apt-key add -
```

Этот ключ служит для проверки подленности пакетов в репозиториях GameAP.


Добавьте репозиторий:
```bash
sudo echo "deb http://packages.gameap.ru/debian/ stretch main" > /etc/apt/sources.list.d/gameap.list
```

### Установка панели

После того, как добавили репозиторий, обновите список пакетов и установите панель следующими командами:

```bash
sudo apt-get update
sudo apt-get install gameap
```

## Ubuntu

### Добавление репозиториев

Скачайте и добавьте GPG ключ:
```bash
wget -O - http://packages.gameap.ru/ubuntu/gameap-rep.gpg.key | sudo apt-key add -
```

Этот ключ служит для проверки подленности пакетов в репозиториях GameAP.


Добавьте репозиторий:
```bash
sudo echo "deb http://packages.gameap.ru/ubuntu/ trusty main" > /etc/apt/sources.list.d/gameap.list
```

### Установка панели

Обновите список пакетов:

```bash
sudo apt-get update
```

Установите панель:

```bash
sudo apt-get install gameap
```

## CentOS

