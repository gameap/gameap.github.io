---
title: Upgrade
layout: default
lang: en
category: Install GameAP
order: 190
---

When installing GameAP, the `gameapctl` utility will be installed, 
which allows you to manage the panel environment, including updates.

> **Before upgrading: two-factor authentication is mandatory for administrators.**
>
> In GameAP 4 the requirement is enabled by default. After the upgrade, administrators without 2FA
> will see a reminder, and after 30 days logging in will no longer issue a full session until 2FA is
> enabled. The countdown for each administrator starts at their first login after the upgrade.
>
> To keep the reminder but remove the lockout, set `AUTH_MFA_HARD_FAIL_DAYS=0` in `config.env`.
> To disable the requirement entirely — `AUTH_REQUIRE_MFA_FOR_ADMINS=false`.
>
> Details and what to do if you lose access are on the
> [Security](/en/security.html) page.

## GameAP Web/API

### Linux

To update the panel, execute the command:
```shell
gameapctl panel upgrade
```

### Windows

To update the panel on Windows, 
you can execute the command where `gameapctl` is installed:
```powershell
C:\path\to\gameapctl.exe panel upgrade
```

Or use the UI. Run `gameapctl.exe`, and in the browser window that opens,
click **"Upgrade"** in the Web/API section.

![](/images/en/gameapctl/ui.png)

# Updating GameAP Daemon

## Linux

To update the Daemon, execute the command:
```shell
gameapctl daemon upgrade
```

## Windows

To update the GameAP Daemon on Windows, you can execute the command 
where `gameapctl` is installed:
```powershell
C:\path\to\gameapctl.exe daemon upgrade
```

Or use the UI. Run `gameapctl.exe`, and in the browser window that opens, 
click **"Upgrade"** in the GameAP Daemon section.

![](/images/en/gameapctl/ui.png)