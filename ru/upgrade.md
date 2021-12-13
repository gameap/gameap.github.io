---
title: Обновление
layout: default
lang: ru
category: Общее
order: 32
---

Перед обновлением панели не забудьте остановить GameAP Daemon на VDS, где работают игровые серверы.
Сами игровые серверы при этом остановлены не будут.

## VDS

Для обновления панели нужно выполнить команду на VDS где установлена панель:
```
curl -sLO http://packages.gameap.ru/installer.sh && bash installer.sh --upgrade
```

На VDS, где установлен GameAP Daemon нужно сделать следующее:

* Скачать последнюю версию `gameap-daemon` [отсюда](https://github.com/gameap/daemon/releases)
* Распаковать файл gameap-daemon в `/usr/bin/`

```bash
curl -qL "https://packages.gameap.ru/gameap-daemon/download-release?os=linux&arch=$(arch)" -o gameap-daemon.tar.gz
tar -xvf gameap-daemon.tar.gz
cp gameap-daemon /usr/bin/gameap-daemon
```

## Shared хостинг

Для тех, у кого шаред. Скачиваем архив http://www.gameap.ru/gameap_latest.zip
Из архива **всё кроме storage каталога** и `.env` файла копируем  с заменой.
Затем заходим в панель, переходим в "Модули" и нажимаем кнопку "Запустить миграцию"

Затем очистите кеш. Удалите файлы `storage/framework/cache`, но не сам каталог

# Обновление GDaemon

## Linux

Для обновления GameAP Daemon остановите его:
```bash
service gameap-daemon stop
```

```bash
curl -qL "https://packages.gameap.ru/gameap-daemon/download-release?os=linux&arch=$(arch)" -o gameap-daemon.tar.gz
tar -xvf gameap-daemon.tar.gz
cp gameap-daemon /usr/bin/gameap-daemon
```

После этого запустите GameAP Daemon вновь:
```
service gameap-daemon start
```


## Windows

* Остановите службу GameAP
* Скачайте последнюю версию с https://github.com/gameap/daemon/releases
* Замените gameap-daemon.exe
* Запустите службу GameAP
