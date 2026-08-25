---
title: Counter-Strike 2
layout: default
lang: es
category: Tutoriales
order: 201
---

Counter-Strike 2 es un shooter táctico en primera persona lanzado en 2023,
desarrollado y publicado por Valve. Es la quinta entrega principal de la
serie Counter-Strike, que aprovecha el éxito de sus predecesores con
gráficos actualizados, mecánicas de juego y nuevas funciones.

## Preparación del entorno

GameAP es totalmente compatible con Counter-Strike 2, incluida la gestión
de jugadores. Para empezar a trabajar con el panel de control de GameAP,
es necesario instalarlo eligiendo una de las opciones disponibles:

* [Instalación de GameAP en Linux](/es/install/install_on_linux.html)
* [Instalación de GameAP en Windows](/es/install/install_on_windows.html)

### Instalación de GameAP Daemon

Para alojar un servidor de juego en su máquina virtual o servidor dedicado,
es necesario instalar GameAP Daemon.
Durante la instalación de GameAP, puede elegir la opción de realizar
una instalación completa junto con el Daemon.

En el panel de control, vaya a **"Administración"** → **"Servidores dedicados"**
→ **"Crear"**.
Aparecerá una ventana con una oferta de instalación automática.
Copie el código y ejecútelo en el servidor dedicado.

Después de esto, puede proceder con la instalación del servidor de Counter-Strike.

## Instalación de Counter-Strike en GameAP

Vaya a **Administración** → **Servidores de juego** → **Crear**

![Formulario de creación de servidor de juego para Counter-Strike 2](/images/en/tutorials/cs2/create_form.png)

* En el campo "Name", introduzca cualquier nombre de servidor,
  por ejemplo, "My Counter-Strike 2 Server".
* En el campo "Game", seleccione Counter-Strike 2.
* En el campo de modificación, elija la modificación.
* En el campo "Dedicated Server", seleccione el nodo deseado
  en el que se ubicará el servidor de juego.
* En el campo IP, elija la dirección deseada de su servidor;
  luego puede elegir un puerto libre de su servidor o utilizar el sugerido.

### Token del servidor

Después de la instalación, debe especificar un token único `sv_setsteamaccount`
para el servidor; sin él, el servidor no funcionará.

Para generarlo, vaya a [https://steamcommunity.com/dev/managegameservers](https://steamcommunity.com/dev/managegameservers).

![Generación de un token de servidor de juego en el sitio web de Steam](/images/en/tutorials/cs2/token_generation.png)

Tras la generación, el valor del token aparecerá en la tabla;
utilice el valor de 32 caracteres:

![Tabla con el token de servidor de juego de Steam generado](/images/en/tutorials/cs2/token_table.png)

Debe copiar el valor del token y especificarlo en la configuración del panel de control.
Vaya a **Servidores** → seleccione su servidor → **Gestión** → **Configuración**

![Campo de token de Steam en la configuración del servidor de juego](/images/en/tutorials/cs2/set_token.png)

Después de esto, puede iniciar su servidor.

## Configuración del servidor de Counter-Strike 2

Para cambiar la configuración del servidor, en el gestor de archivos,
vaya al directorio `/game/csgo/cfg`, donde encontrará muchos archivos *.cfg.

El archivo de configuración principal del servidor de Counter-Strike 2 es `server.cfg`

![El archivo server.cfg de un servidor de Counter-Strike 2 en el gestor de archivos](/images/en/tutorials/cs2/server_config.png)
