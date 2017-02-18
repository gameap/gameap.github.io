---
title: Установка GameAP Daemon
layout: default
lang: ru
category: Общее
---

# Ручная установка GameAP Daemon

## Debian

Следующая инструкция актуальна для Debian Jessie и Debian Stretch

## Добавление репозиториев

Скачайте и добавьте ключ:
```bash
wget -O - http://packages.gameap.ru/debian/gameap-rep.gpg.key | apt-key add -
```

Добавьте новый репозиторий:
```bash
echo "deb http://packages.gameap.ru/debian/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/gameap.list
```

Обновите список пакетов и запустите установку:
```bash
apt-get update
apt-get install gameap-daemon
```