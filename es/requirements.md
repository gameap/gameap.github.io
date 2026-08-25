---
title: Requisitos
layout: default
lang: es
category: General
order: 10
---

GameAP 4 es un único ejecutable con una interfaz web integrada. No se necesita ningún servidor
web, PHP, Composer ni Node.js para ejecutar el panel.

## Requisitos del sistema

### Panel

* RAM: 512 MB o más
* Disco: 200 MB o más
* El panel no es exigente con la CPU, un núcleo es suficiente

### GameAP Daemon

Los valores siguientes corresponden al propio daemon, **sin contar los servidores de juego**.
Reserve recursos para los servidores de juego por separado: algunos juegos tienen requisitos
elevados.

* RAM: 128 MB
* Disco: 1 GB o más, más espacio para los archivos de los servidores de juego
* El daemon no es exigente con la CPU, un núcleo es suficiente

El panel y el daemon se pueden instalar en el mismo servidor.

## Sistemas operativos

Las compilaciones se publican para:

| Sistema | Arquitecturas            |
|---------|--------------------------|
| Linux   | `amd64`, `arm64`, `386`  |
| Windows | `amd64`, `arm64`, `386`  |
| macOS   | `amd64`, `arm64`         |

La instalación mediante `gameapctl` se describe en las páginas
[Instalación en Linux](/es/install/install_on_linux.html) e
[Instalación en Windows](/es/install/install_on_windows.html).

Las versiones de Windows compatibles se indican en la página de instalación.

## Red

| Puerto  | Quién lo necesita                                          | Obligatorio |
|---------|------------------------------------------------------------|----------|
| `8025`  | Interfaz web y API. Lo necesita el navegador del administrador | sí      |
| `31718` | gRPC. Los daemons se conectan al panel a través de él      | sí, si hay servidores dedicados remotos |
| `443`   | HTTPS, cuando el propio panel lo sirve                     | no       |
| `80`    | Desafío `http-01` de Let's Encrypt                         | no       |

Los puertos 8025 y 443 se cambian con las variables `HTTP_PORT` y `HTTPS_PORT`, y el puerto gRPC
con `GRPC_PORT`. Consulte la [referencia de config.env](/es/config.html).

El daemon necesita acceso saliente al panel por el puerto 31718, y ese es el único puerto del
panel que requiere. El panel no se conecta al daemon; no es necesario abrir ningún puerto de
entrada en el servidor dedicado.

## Base de datos

Se necesita una de las siguientes:

| DBMS             | Nota                                                              |
|------------------|-------------------------------------------------------------------|
| PostgreSQL       | Recomendado para instalaciones con varias instancias del panel    |
| MySQL / MariaDB  |                                                                   |
| SQLite           | No se necesita un servidor de base de datos separado, el archivo se crea automáticamente |

Para una instalación pequeña, SQLite es suficiente: no requiere ni un servicio separado ni
ninguna configuración.

La cadena de conexión se define con las variables `DATABASE_DRIVER` y `DATABASE_URL`; los
formatos se indican en la [referencia de config.env](/es/config.html).

Al actualizar desde GameAP 3, el panel se conecta a la base de datos MySQL existente y la
migra directamente; consulte [Actualización de v3 a v4](/es/upgrade_from_v3_to_v4.html).

## Componentes opcionales

Necesarios solo en escenarios concretos, no en una instalación habitual.

| Componente       | Cuándo se necesita                                                        |
|------------------|-------------------------------------------------------------------------------|
| Redis            | Caché compartida e intercambio de eventos entre varias instancias del panel   |
| Almacenamiento S3 | Almacenamiento compartido para archivos y certificados ACME con varias instancias del panel |
| Proxy inverso    | Cuando HTTPS lo sirve nginx o Traefik en lugar del propio panel               |

## Requisitos del servidor dedicado

La instalación automática del daemon necesita:

* privilegios de root (Linux) o de administrador (Windows);
* las utilidades `curl`, `tar` e `install`, para el script de instalación en Linux;
* acceso saliente a `github.com` y `api.github.com`: desde ahí se descargan gameapctl y el
  daemon;
* acceso saliente al panel por el puerto 31718.

Los paquetes que necesitan el daemon y los servidores de juego —incluidos SteamCMD y un gestor
de procesos— los instala el propio `gameapctl`. Consulte
[Servidores dedicados](/es/gameap_configure/dedicated_servers.html) para más detalles.

### Curl

Si `curl` no está presente en el sistema, instálelo:

```shell
# Debian, Ubuntu
apt install curl

# CentOS, RHEL, Fedora
yum install curl
```
