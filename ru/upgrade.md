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

На VDS, где установлен GameAP Daemon нужно выполнить следующие команды:
```
apt update
apt install gameap-daemon
service gameap-daemon restart
```

## Shared хостинг

Для тех, у кого шаред. Скачиваем архив http://www.gameap.ru/gameap_latest.zip
Из архива **всё кроме storage каталога** и `.env` файла копируем  с заменой.
Затем заходим в панель, переходим в "Модули" и нажимаем кнопку "Запустить миграцию"

Затем очистите кеш. Удалите файлы `storage/framework/cache`, но не сам каталог

# Обновление GDaemon

Для обновления GameAP Daemon остановите его:
```bash
service gameap-daemon stop
```

Затем установите новую версию пакета.

**Debian**
```bash
apt update
apt install gameap-daemon
```

**CentOS**
```bash
yum install gameap-daemon
```

После этого запустите GameAP Daemon вновь:
```
service gameap-daemon start
```