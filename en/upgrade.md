---
title: Upgrade
layout: default
lang: en
category: Install GameAP
order: 190
---

When installing GameAP, the `gameapctl` utility will be installed, 
which allows you to manage the panel environment, including updates.

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

## Shared Hosting

Download the archive https://packages.gameap.com/gameap/gameap-3.1-shared.zip. 
From the archive, copy **everything except the storage directory** 
and the `.env` file, replacing them in the directory where your panel files 
are located. Then, go to the panel, navigate to "Modules" and click 
the "Run migration" button.

After that, clear the cache. Delete the files inside `storage/framework/cache`, 
but not the directory itself.

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