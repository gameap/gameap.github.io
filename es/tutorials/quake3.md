---
title: Quake III Arena
layout: default
lang: es
category: Tutoriales
order: 204
---

Quake III Arena es un shooter multijugador en primera persona desarrollado por id Software.
El juego se centra por completo en el combate multijugador y carece de una campaña tradicional para un jugador.

El juego es similar a Unreal Tournament y a otros shooters en primera persona de finales de los años 1990 y principios de los 2000.

## Preparación del entorno

GameAP ofrece soporte para los servidores de juego de Quake III Arena, su configuración y gestión.

* [Instalación de GameAP en Linux](/es/install/install_on_linux.html)
* [Instalación de GameAP en Windows](/es/install/install_on_windows.html)

### Instalación de GameAP Daemon

GameAP Daemon es un agente responsable de gestionar los servidores de juego en máquinas dedicadas. Para instalar un servidor de juego en una máquina (VDS), es necesario instalar GameAP Daemon.

Al instalar GameAP, puede elegir una instalación completa que incluya el Daemon (utilizando la opción `--with-daemon`).

En el panel de control, vaya a **"Administración"** → **"Servidores dedicados"** → **"Crear"**. Aparecerá una ventana con una oferta de instalación automática. Copie el código y ejecútelo en el servidor dedicado.

Después de esto, puede proceder con la instalación del servidor de Quake III Arena.

## Instalación del servidor de Quake III Arena en GameAP

Vaya a **Administración** → **Servidores de juego** → **Crear**

![Formulario de creación de servidor de juego para Quake III Arena](/images/en/tutorials/quake3/create_form.png)

* En el campo "Name", introduzca cualquier nombre de servidor, por ejemplo "My Quake III Server".
* En el campo "Game", seleccione "Quake 3" en la lista desplegable.
* En el campo "Modification", especifique la modificación deseada; el valor predeterminado es `ioquake3`. También está disponible la opción `quake3e`.
* En el campo "Dedicated Server", especifique el nodo donde se alojará el servidor de juego. De forma predeterminada, se selecciona el primer nodo disponible.
* En el campo "IP", seleccione la dirección deseada para su servidor; luego especifique un puerto disponible o deje el sugerido por el sistema.

## Configuración del servidor de Quake III Arena

Para cambiar la configuración del servidor, vaya a la sección **Servidor**, seleccione su servidor y haga clic en **Gestión**. A continuación, abra la pestaña **Configuración**.

![Pestaña de configuración de un servidor de juego de Quake III Arena](/images/en/tutorials/quake3/settings.png)

Para que los cambios surtan efecto, es necesario reiniciar el servidor.

### Server Hostname / Nombre del servidor

El nombre de su servidor de Quake III que será visible para todos los jugadores en la lista de servidores al realizar la búsqueda.

### Maximum Players on Server / Número máximo de jugadores en el servidor

El número máximo de jugadores que pueden estar en el servidor al mismo tiempo. El valor predeterminado es 16.

### Default Map of the Server / Mapa predeterminado del servidor

El mapa que se cargará automáticamente al iniciar el servidor. El mapa predeterminado es `q3dm17`, conocido como "The Longest Yard".

### Timelimit in Minutes / Límite de tiempo en minutos

La duración máxima de la ronda en minutos. El valor predeterminado es 20 minutos.

### Frag Limit / Límite de frags

El número de frags (eliminaciones) al alcanzar el cual finaliza la ronda. El valor predeterminado es 20 frags.
