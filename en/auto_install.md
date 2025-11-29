---
title: Automatic panel installation
layout: default
lang: en
category: Main
order: 20
hidden: true
---

> ! This page is outdated, please use one of the following up-to-date guides:
* [Installation on Linux](/en/install/install_on_linux.html)
* [Installation on Windows](/en/install/install_on_windows.html)

---

* This will become a table of contents (this text will be scraped).
{:toc}

Fully automatic installation. 
You need to run the script, it will automatically install the packages 
dependencies and panel. 
During the installation, you will need to enter and select some parameters:

- Panel host. Domain or IP
- Select database: PostgreSQL, MySQL, SQLite

## Dependencies installation

To run the script you need CURL.

### Debian/Ubuntu

```bash
sudo apt-get update
sudo apt-get install curl
```

### Centos

```bash
sudo dnf update
sudo dnf install curl
```

Or, for older versions:

```bash
sudo yum update
sudo yum install curl
```

## Running auto installation script

Download and run script:
```bash
curl -sLO https://gameap.com/install.sh
bash ./install.sh
```

## Running with parameters

For silent installation no questions asked set parameters:

- `--host` Panel host, IP or domain name.
- `--database` Possible values: `postgres`, `mysql`, `sqlite`, `none`

### Examples

#### Panel installation

The following example will automatically install the panel, install and configure database. 
The panel should be available at the address `http://example.com` specified in `--host`:

```bash
./installer.sh \
    --host=example.com \
    --database=mysql \
    --github
```

The following example will automatically install the panel with SQLite database.

```bash
./installer.sh \
    --host=localhost \
    --database=sqlite
```

#### Panel upgrading

The followind example will upgrade panel to last available version:
```bash
./installer.sh --upgrade
```

## Other packages that the script installs

In addition to the panel itself, the script installs the necessary packages and their dependencies. Most likely, all or most of them will already be installed on your system.

### Packages

The script automatically installs the following packages.:

- `software-properties-common` APT Repository Management package.
- `apt-transport-https` https support for APT
- `gnupg` Package for working with digital signatures and keys. Required to authenticate packages and add repository keys.