---
title: Troubleshooting
layout: default
lang: en
category: Main
order: 40
---

* This will become a table of contents (this text will be scraped).
{:toc}

Description some possible errors and their fixing.

## Errors after installation

### Whoops, looks like something went wrong.

#### Generation the encryption key

![](/images/errors/key_generate.png)

This problem will occur if you have not generated a key. Go to the panel directory and run the command:

```bash
php artisan key:generate --force
```

#### Database migration

![](/images/errors/db_migrate.png)

This problem will occur if you have not migrate database. Go to the panel directory and run the command:

```
php artisan migrate --seed
```

## Server Startup Errors

More details you can view in logs.

Go to **Administration** -> **GDaemon tasks**, open last server startup task.

Check GameAP Daemon logs. 
Проверьте логи GameAP Daemon. Logs are stored on a dedicated server in a directory `/var/log/gameap-daemon`.

### Server status is displayed incorrectly

Sometimes it happens that the server starts, but its status in the panel is displayed as offline.

### Empty server start command

This error occurs when the launch command for the game server is empty, it must be filled.

Go to the game server administration page: **Administration** -> **Servers** -> Then select your game server.

Or from main page **Servers List** -> Select game server -> **Control** -> **Administration**

Find the field **Start command**, fill the field. For Counter-Strike 1.6 startup command will be something like this:
```
./hlds_run -game cstrike +ip {ip} +port {port} +map {default_map} +maxplayers {maxplayers} +sys_ticrate {fps}
```

Sometimes after settings changing you should restart GameAP Daemon.

### The information modal window does not change for a long time

Usually server startup takes less than 10 seconds to start. If the status bar is frozen at 10% or the status does not 
change for several minutes, then try to start/restart GameAP Daemon:
```
service gameap-daemon restart
```
