---
title: Primeros pasos
layout: default
lang: es
category: General
order: 2
---

Para empezar, es recomendable disponer de dos servidores dedicados o virtuales. 
En uno se encuentra el panel de control, en el otro los servidores de juego. 
También puede instalar todo en un único servidor dedicado.

## Instalación del panel

El panel se instala en un servidor dedicado con una base de datos (PostgreSQL, MySQL, SQLite).

&nbsp;&nbsp;&nbsp;&nbsp;<img src="/images/icons/linux.svg" alt="Linux" width="20" height="20" style="vertical-align: middle"> [Instalación en Linux](/es/install/install_on_linux.html)

&nbsp;&nbsp;&nbsp;&nbsp;<img src="/images/icons/windows.svg" alt="Windows" width="20" height="20" style="vertical-align: middle"> [Instalación en Windows](/es/install/install_on_windows.html)

## La instalación más sencilla en Linux

Si tiene Linux, CURL instalado y no quiere entrar en los detalles de la instalación, ejecute el comando:
```bash
bash <(curl -s https://gameap.com/install.sh) --with-daemon
```

## Añadir un servidor dedicado

Añada un nuevo servidor dedicado (VDS) en el que después instalará los servidores de juego. 
Después de instalar el panel, inicie sesión y seleccione **"Administración"** **"Servidores dedicados"** → **"Crear"** en el menú. A continuación,
se abrirá una ventana con instrucciones, sígalas.

![Añadir un servidor dedicado en el panel de GameAP](/images/en/get_started/add_dedicated_server.gif)

Para información más detallada sobre la instalación y la configuración, consulte la página [Servidores dedicados](/es/gameap_configure/dedicated_servers.html).

## Añadir un servidor de juego

Vaya a **"Administración"** → **"Servidores de juego"** → **"Crear"**.

![Crear un servidor de juego en el panel de GameAP](/images/en/get_started/add_game_server.gif)

Para más detalles, consulte la página de [servidores de juego](/es/gameap_configure/game_servers.html).
