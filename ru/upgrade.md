---
title: Обновление
layout: default
lang: ru
category: Установка GameAP
order: 190
---

При установке GameAP будет установлена утилита `gameapctl`, 
которая позволяет управлять окружением панели, в том числе и обновлением.

## GameAP Web/API

### Linux

Для обновления панели необходимо выполнить команду:
```shell
gameapctl panel upgrade
```

### Windows

Чтобы обновить панель на Windows можно выполнить команду, 
где установлена gameapctl:
```powershell
C:\path\to\gameapctl.exe panel upgrade
```

Либо воспользуйтесь UI. Запустите gameapctl.exe,
а в открывшемся окне в браузере нажмите **"Upgrade"** в разделе Web/API

![](/images/en/gameapctl/ui.png)

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