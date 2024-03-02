---
title: Обновление
layout: default
lang: ru
category: Установка GameAP
order: 190
---

При установке GameAP будет установлена утилита gameapctl, которая позволяет управлять окружением панели, в том
числе и обновлением.

## GameAP Web/API

### Linux

Для обновления панели необходимо выполнить команду:
```shell
gameapctl panel upgrade
```

### Windows

Чтобы обновить панель на Windows можно выполнить команду, где установлена gameapctl:
```powershell
C:\path\to\gameapctl.exe panel upgrade
```

Либо воспользуйтесь UI. Запустите gameapctl.exe, а в открывшемся окне в браузере нажмите **"Upgrade"**
в разделе Web/API

![](/images/en/gameapctl/ui.png)

## Shared хостинг

Скачайте архив https://packages.gameap.com/gameap/gameap-3.1-shared.zip
Из архива **всё кроме storage каталога** и `.env` файла скопируйте с заменой в каталог, где у вас находятся файлы панели.
Затем зайдите в панель, перейдите в "Модули" и нажмите кнопку "Запустить миграцию"

Затем очистите кеш. Удалите файлы внутри `storage/framework/cache`, но не сам каталог

# Обновление GameAP Daemon

## Linux

Для обновления Daemon необходимо выполнить команду:
```shell
gameapctl daemon upgrade
```

## Windows

Чтобы обновить GameAP Daemon на Windows можно выполнить команду, где установлена gameapctl:
```powershell
C:\path\to\gameapctl.exe daemon upgrade
```

Либо воспользуйтесь UI. Запустите gameapctl.exe, а в открывшемся окне в браузере нажмите **"Upgrade"**
в разделе GameAP Daemon

![](/images/en/gameapctl/ui.png)