---
title: Grand Theft Auto V
layout: default
lang: es
category: Tutoriales
order: 203
---

Grand Theft Auto V es un popular juego de mundo abierto desarrollado por Rockstar North.

FiveM es una modificación para GTA V que permite jugar en línea en modo multijugador. 
Además de FiveM, también existe una modificación similar llamada RageMP.

Puede crear y administrar su propio servidor de GTA V (FiveM) con la ayuda de GameAP.

## Preparación del entorno

Con GameAP, puede administrar los parámetros principales de su servidor. 
GameAP ofrece soporte completo para servidores de juego FiveM, 
incluida la gestión de jugadores.

La instalación del panel llevará poco tiempo. Dependiendo de su sistema operativo, elija una de las siguientes opciones de instalación del panel de control:

* [Instalación en Linux](/es/install/install_on_linux.html)
* [Instalación en Windows](/es/install/install_on_windows.html)

### Instalación de GameAP Daemon

Para instalar el servidor de juego en una máquina (VDS), debe instalar GameAP Daemon. 
Durante la instalación de GameAP, puede optar por instalarlo por completo, 
incluido Daemon.

En el panel de control, vaya a **"Administración"** → **"Servidores dedicados"** 
→ **"Crear"**. Aparecerá una ventana con una oferta de instalación automática. 
Copie el código y ejecútelo en el servidor dedicado.

Después de esto, puede proceder a instalar el servidor FiveM.

### Creación de un servidor FiveM en GameAP

Vaya a **Administración** → **Servidores de juego** → **Crear**

![Formulario de creación de servidor de juego para FiveM](/images/en/tutorials/fivem/create_form.png)

* En el campo "Nombre", introduzca cualquier nombre de servidor, por ejemplo, "My GTA V Server".
* En el campo "Juego", seleccione FiveM de la lista.
* En el campo de modificación, seleccione la modificación; por defecto es Vanilla.
* En el campo "Servidor dedicado", seleccione el nodo deseado donde estará ubicado el servidor de juego.
* En el campo IP, elija la dirección deseada de su servidor; luego puede elegir un puerto libre para su servidor o utilizar el sugerido.

## Clave del servidor

Después de la instalación, debe especificar una clave que debe obtener de 
[keymaster.fivem.net](https://keymaster.fivem.net)

![Creación de una clave de servidor en keymaster.fivem.net](/images/en/tutorials/fivem/generate_key.png)

Después de la generación, verá un mensaje. Debe copiar el valor de la clave.

![La clave del servidor FiveM generada](/images/en/tutorials/fivem/key.png)

Debe copiar el valor de la clave y especificarlo en 
la configuración del panel de control. 
Vaya a **Servidores** → seleccione su servidor FiveM → **Administración** → **Configuración**

![Campo de la clave FiveM en la configuración del servidor de juego](/images/en/tutorials/fivem/set_key.png)

Ahora puede iniciar su servidor FiveM en el panel.

## Configuración del servidor FiveM

La configuración del servidor FiveM se encuentra en el archivo `server.cfg`, 
que está ubicado en el directorio raíz. 
Puede editar este archivo en el gestor de archivos del panel.

Vaya a **Servidores** → seleccione su servidor FiveM → **Administración** → **Archivos**

![Archivo de configuración del servidor FiveM en el gestor de archivos](/images/en/tutorials/fivem/server_config.png)
