---
title: Cómo funciona GameAP
layout: default
lang: es
category: General
order: 5
---

GameAP consta de dos aplicaciones: el **panel**, con el que trabaja el administrador, y el
**daemon**, que se ejecuta en cada servidor dedicado y gestiona los servidores de juego.

![Arquitectura de GameAP: el panel, los daemons en los servidores dedicados y los servidores de juego](/images/en/gameap_architecture.svg)

## Tres capas

El **panel** es una aplicación única con una interfaz web integrada. Almacena todos los datos:
usuarios, servidores dedicados, juegos, servidores de juego, tareas. No inicia los servidores
de juego por sí mismo ni toca sus archivos — solo envía comandos a los daemons.

**GameAP Daemon** se ejecuta en cada servidor dedicado. Inicia y detiene los servidores de
juego, supervisa su estado, los instala y actualiza, trabaja con archivos y recopila métricas.

Los **servidores de juego** son procesos que el daemon gestiona a través de un gestor de
procesos: systemd, Docker, Podman, tmux y otros.

El panel y el daemon pueden instalarse en un mismo servidor — en ese caso las tres capas se
encuentran en la misma máquina.

## Quién se conecta a quién

La conexión siempre la establece el **daemon**: se conecta al panel por sí mismo y mantiene
una única conexión persistente por la que pasa todo el tráfico.

El panel no se conecta al daemon. El servidor dedicado no necesita puertos entrantes abiertos
ni una dirección IP externa — basta con el acceso saliente al panel.

Esto difiere de GameAP 3, donde el panel se conectaba al daemon.

| Dirección                               | Puerto          | Protocolo        |
|-----------------------------------------|-----------------|------------------|
| Navegador del administrador → panel     | `8025`          | HTTP, HTTPS      |
| Daemon → panel                          | `31718`         | gRPC             |
| Panel → servidor de juego               | puerto del juego | Query, RCON      |

Los puertos `8025` y `31718` son **listeners independientes**, no un solo puerto con detección
de protocolo.

El panel realiza las solicitudes Query y RCON a los servidores de juego por sí mismo,
directamente — no pasan por el daemon.

## Qué ocurre al iniciar un servidor

1. El administrador hace clic en un botón del panel.
2. El panel crea una tarea y la guarda en la base de datos.
3. La tarea llega al daemon a través de la conexión establecida.
4. El daemon inicia el servidor de juego a través del gestor de procesos.
5. El daemon envía el progreso y la salida del comando; el panel los muestra en tiempo real a
   través de WebSocket.
6. A partir de entonces, el daemon informa periódicamente del estado del servidor y las métricas.

Si la conexión con el panel se pierde, los servidores de juego siguen funcionando — solo el
control desde el panel no está disponible. Una vez restablecida la conexión, el daemon se
reconecta y recibe el estado actual completo del panel.

## Dónde se almacena cada cosa

| Datos                                        | Dónde                                             |
|----------------------------------------------|---------------------------------------------------|
| Usuarios, servidores, juegos, tareas         | Base de datos del panel                           |
| Certificados gRPC, datos ACME                | Almacenamiento de archivos del panel              |
| Sesiones, contadores, clave de instalación   | Caché del panel                                   |
| Configuración del panel                      | `config.env`                                      |
| Archivos del servidor de juego               | El servidor dedicado, en el directorio de trabajo del daemon |
| Configuración del daemon                     | `gameap-daemon.yaml` en el servidor dedicado      |
| Métricas                                     | La RAM del daemon, durante no más de una hora     |

Los archivos de los servidores de juego no se copian al panel: el gestor de archivos trabaja
con ellos a través del daemon.

## Ampliación

Los **plugins** se ejecutan dentro del panel en un entorno aislado de WebAssembly. Añaden
páginas, pestañas e integraciones, pero no tienen acceso propio al sistema — solo a través de
la interfaz controlada del panel. Consulte [Plugins](/es/plugins/index.html).

**API** — todo lo que hace la interfaz está disponible a través de la API HTTP; la propia
interfaz funciona a través de ella. Consulte [API y tokens](/es/api.html).

## A continuación

* [Requisitos](/es/requirements.html) — lo necesario para la instalación
* [Primeros pasos](/es/get_started.html) — instalación del panel y del primer servidor dedicado
* [API gRPC](/es/daemon/grpc.html) — detalles del canal de comunicación panel–daemon
* [Múltiples instancias del panel](/es/multi_instance.html) — una configuración tolerante a fallos
