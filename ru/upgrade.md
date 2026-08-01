---
title: Обновление
layout: default
lang: ru
category: Установка GameAP
order: 190
---

При установке GameAP будет установлена утилита `gameapctl`, 
которая позволяет управлять окружением панели, в том числе и обновлением.

> **Перед обновлением: обязательная двухфакторная аутентификация для администраторов.**
>
> В GameAP 4 требование включено по умолчанию. После обновления администраторы без 2FA увидят
> напоминание, а через 30 дней вход перестанет выдавать полноценную сессию, пока 2FA не будет
> подключена. Отсчёт для каждого администратора начинается с его первого входа после обновления.
>
> Чтобы оставить напоминание, но убрать блокировку, задайте в `config.env`
> `AUTH_MFA_HARD_FAIL_DAYS=0`. Чтобы отключить требование целиком —
> `AUTH_REQUIRE_MFA_FOR_ADMINS=false`.
>
> Подробности и порядок действий при потере доступа — на странице
> [Безопасность](/ru/security.html).

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