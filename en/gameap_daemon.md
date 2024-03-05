---
title: GameAP Daemon
layout: default
lang: en
category: GameAP Daemon
order: 400
---

* This will become a table of contents (this text will be scraped).
{:toc}

GameAP Daemon — backgrouond application that allows panel to manage game servers (install, delete, run, stop etc.).

This is the main application, what control game servers statuses and reload them in case of failure or by user request

## Installation

### Automatically from panel

Go to **Administration** -> **Dedicated servers** -> **Create**
That action open suggestion about autoinstall. Copy the code and run on dedicated server.

## Configuration

Gameap-daemon config is here:

| OS         | Path
|------------|------------
| Linux      | /etc/gameap-daemon/gameap-daemon.yaml
| Windows    | C:\gameap\daemon\daemon.yaml

### Base examples

| Parameter                 | Required              | Type       | Info
|---------------------------|-----------------------|------------|------------
| ds_id                     | yes                   | integer    | Dedicated server ID
| listen_port               | no (default: 31717)   | integer    | Listening port
| api_host                  | yes                   | string     | API host
| api_key                   | yes                   | string     | API key


### SSL/TLS

| Parameter                 | Required              | Type      | Info
|---------------------------|-----------------------|-----------|------------
| ca_certificate_file       | yes                   | string    | CA certificate
| certificate_chain_file    | yes                   | string    | Server certificate
| private_key_file          | yes                   | string    | Server private key
| private_key_password      | no                    | string    | Server private key password
| dh_file                   | yes                   | string    | DH certificate

[Certificate creation](https://github.com/gameap/GDaemon2#creating-certificates)

### Authentication by login and password

| Parameter                 | Required              | Type      | Info
|---------------------------|-----------------------|-----------|------------
| password_authentication   | no                    | boolean   | Enable login & password auth?
| daemon_login              | no                    | string    | Login (in Linux it will use PAM if user is not exist)
| daemon_password           | no                    | string    | Password. (in Linux it will use PAM if user is not exist)

### Statistics

| Parameter                 | Required              | Type       | Info
|---------------------------|-----------------------|-----------|------------
| if_list                   | no                    | string    | Interfaces list
| drives_list               | no                    | string    | Drives list
| stats_update_period       | no                    | integer   | Statistics update period
| stats_db_update_period    | no                    | integer   | Statistics update period in db

[Config example](https://github.com/gameap/GDaemon2#example-daemoncfg)

### Non anonymous steam account usage

Many steam game servers require log in to account that has a purchased copy of game, in that case you need to add steam login and password in to config file.

You need to transform config from .cfg to .yaml manually (ofcourse with yaml syntax requirements)

| OS         | Before            | After
|------------|-------------------|------------
| Linux      | gameap-daemon.cfg | gameap-daemon.yaml
| Windows    | daemon.cfg        | daemon.yaml

[yaml config example](https://github.com/gameap/daemon/blob/master/config/gameap-daemon.yaml)

After that you can add steam account configuration rows:

```
steam_config:
    login: "your login"
    password: "Pa$$worD"
```
