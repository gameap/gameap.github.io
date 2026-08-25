---
title: Install on Linux
layout: default
lang: en
category: Install GameAP
order: 100
---

## Installation

Installation on Linux is performed with a single command:
```shell
bash <(curl -s https://gameap.com/install.sh)
```

During the installation process, you will be prompted to enter some information.

### Host

Specify the domain or IP address at which the panel will be accessible.

In the case of an IP address, it must be an address assigned to the network 
interface on the VDS. If your network uses NAT, do not specify 
the external IP, but rather the internal one, and then configure port 
forwarding.

Any domain can be specified, but do not forget to configure DNS.

Examples of correct values:
* 10.182.104.8
* 10.182.104.8:2080
* example.com
* http://example.com

### Database

The database where data will be stored: users, information about servers, etc. 
You can use [PostgreSQL](https://www.postgresql.org/),
[MySQL](https://www.mysql.com/) and 
[SQLite](https://www.sqlite.org/).

PostgreSQL is recommended in most cases. If the load on your server 
is expected to be low, and you do not plan to use more than 10 game servers, 
you can use SQLite.

Some distributions may have [MariaDB](https://mariadb.org/) installed.

## Completing the Installation

At the end of the installation, the access details for the panel 
will be displayed. Do not forget to save this information to access the panel.

![Panel login details shown when the installation finishes](/images/en/gameapctl/gameap_finished_installation.png)

## Additional Installation Options

### Develop Version

You can install the version currently in development by passing the extra
flags `--github --branch=develop` to the installer.
The installation will take noticeably longer in this case, since it is
performed from source.

```shell
bash <(curl -s https://gameap.com/install.sh) \
  --github \
  --branch=develop
```

### Non-Interactive Installation

This type of installation lets you install the panel without entering any data
during the process. Pass the flags, and the installer will not need any
additional input from you.

```shell
bash <(curl -s https://gameap.com/install.sh) \
  --non-interactive \
  --host=panel.example.com \
  --port=8025 \
  --database=sqlite
```

> Put the address the panel is reached at into `--host` — a domain name or the server's external
> IP. It becomes `HTTP_HOST`, and the panel derives the listening address from it: with
> `--host=127.0.0.1` the HTTP and gRPC listeners come up on loopback only, and neither remote
> administrators nor daemons on other machines will be able to connect.

Main flags:

| Flag                  | Purpose                                                             |
|-----------------------|---------------------------------------------------------------------|
| `--non-interactive`   | Ask no questions                                                    |
| `--host`              | Address at which the panel will be accessible                       |
| `--port`              | Panel port, `8025` by default                                       |
| `--grpc-port`         | gRPC port for daemons, `31718` by default                           |
| `--database`          | `sqlite`, `mysql`, or `postgres`                                    |
| `--database-host`     | Database host                                                       |
| `--database-port`     | Database port                                                       |
| `--database-name`     | Database name                                                       |
| `--database-username` | Database user                                                       |
| `--database-password` | Database user password                                              |
| `--with-daemon`       | Also install GameAP Daemon                                          |
| `--version`           | Specific panel version                                              |

SQLite needs no connection parameters — the database file is created automatically.

> The `--path` and `--web-server` flags are leftovers from GameAP 3 and have no effect when
> installing GameAP 4: the panel is a single executable with a built-in web interface and does
> not need a separate web server.

### Full Installation

To install the GameAP Daemon in addition to the panel itself, 
add the `--with-daemon` flag.

This method is recommended if you plan to host both the panel 
and game servers on the same VDS.

```Shell
bash <(curl -s https://gameap.com/install.sh) --with-daemon
```
