---
title: Minecraft
layout: default
lang: es
category: Tutoriales
order: 200
---

Minecraft es un popular juego sandbox desarrollado
por Mojang Studios, una filial de Xbox Game Studios.
Permite a los jugadores construir, explorar y sobrevivir en
un mundo generado proceduralmente y compuesto por bloques.
El juego cuenta con varios modos: supervivencia, creativo
y aventura. Minecraft tiene una gran comunidad de jugadores
que crean y comparten diversas creaciones, mods y skins.

## Preparación del entorno

Para comenzar, es necesario [instalar GameAP](/es/get_started.html#panel-installation), lo que tomará unos minutos:

* [Instalación de GameAP en Linux](/es/install/install_on_linux.html)
* [Instalación de GameAP en Windows](/es/install/install_on_windows.html)

### Instalación de GameAP Daemon

En el nodo donde se alojará el servidor de juego, es necesario
instalar GameAP Daemon. Durante el proceso de instalación de GameAP,
puede elegir la opción de instalarlo junto con el Daemon.

En el panel de control, vaya a **Administración** → **Servidores dedicados**
→ **Crear**.
Aparecerá una ventana con una oferta de instalación automática.
Copie el código y ejecútelo en el servidor dedicado.

Después de esto, puede proceder a instalar el servidor de Minecraft.

## Instalación de Minecraft en GameAP

Vaya a **Administración** → **Servidores de juego** → **Crear**

![Formulario de creación de servidor de juego para Minecraft](/images/en/tutorials/minecraft/create_form.png)

* En el campo "Name", introduzca cualquier nombre de servidor.
* En el campo "Game", seleccione Minecraft.
* En el campo de modificación, recomendamos elegir "Multicore", que permite gestionar fácilmente las versiones del servidor de Minecraft.
* En el campo "Dedicated Server", seleccione el nodo deseado donde se ubicará el servidor de juego.
* En el campo IP, elija la dirección deseada de su servidor; luego puede elegir un puerto libre para su servidor o utilizar el sugerido.

Vea un breve vídeo sobre el proceso de instalación de un servidor de Minecraft en GameAP:

<video controls width="600">
  <source src="/media/en/tutorials/minecraft/installation.webm" type="video/webm" />
  <source src="/media/en/tutorials/minecraft/installation.mp4" type="video/mp4" />
</video>

## Gestión del servidor de Minecraft

Para gestionar el servidor de Minecraft, vaya a la sección **Servidores**, luego seleccione su servidor y haga clic en **Gestión**. En esta página podrá iniciar y detener su servidor, ver su consola y enviar comandos.

![Página de gestión del servidor de Minecraft con la consola](/images/en/tutorials/minecraft/server_management.png)

### Cambio de la versión del servidor

Si eligió la modificación "Multicore", puede cambiar la versión de su servidor. Para ello, vaya a la página principal de gestión de su servidor y luego seleccione **Configuración**.

![Configuración del servidor de juego de Minecraft con el selector de versión](/images/en/tutorials/minecraft/server_settings.png)

#### Ejemplos de configuración

##### Predeterminada

* Versión de Minecraft: 1.14.3
* Núcleo: vanilla

##### Forge 1.19

En este ejemplo, utilizaremos un servidor de Minecraft versión 1.19.4 con Forge API para soporte de mods versión 45.0.50

* Versión de Minecraft: 1.19.4
* Núcleo: forge
* Versión del mod del núcleo: 45.0.50

##### Spigot

En este ejemplo, utilizaremos un servidor de Minecraft versión 1.19.4 con Spigot API para soporte de plugins.

* Versión de Minecraft: 1.19.4
* Núcleo: spigot

### server.properties

`server.properties` es el archivo de configuración principal del servidor de Minecraft. Puede editarlo mediante el gestor de archivos de GameAP.

![El archivo server.properties en el gestor de archivos de GameAP](/images/en/tutorials/minecraft/server_properties.png)
