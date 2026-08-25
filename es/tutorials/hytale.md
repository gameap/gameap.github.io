---
title: Hytale
layout: default
lang: es
category: Tutoriales
order: 205
---

Hytale es un juego sandbox desarrollado por Hypixel Studios. El estilo visual pixelado general y la mecánica
de juego son similares a los de Minecraft, con animaciones y efectos de personajes más avanzados.

El juego incluye numerosos biomas, mobs y objetos, además de compatibilidad con mods.

## Configuración del entorno

Con GameAP, puede crear, administrar y configurar fácilmente un servidor de juego de Hytale.

Primero, debe [instalar GameAP](/es/get_started.html#panel-installation):

* [Instalación de GameAP en Linux](/es/install/install_on_linux.html)
* [Instalación de GameAP en Windows](/es/install/install_on_windows.html)

### Instalación de GameAP Daemon

GameAP Daemon es un agente responsable de gestionar los servidores de juego en máquinas dedicadas.
Puede instalarse tanto en la misma máquina donde se ejecuta GameAP como en una máquina independiente.

Al instalar GameAP, puede optar por una instalación completa del panel junto con el Daemon utilizando el indicador `--with-daemon`.

## Instalación del servidor de Hytale en GameAP

Vaya a **Administración** → **Servidores de juego** → **Crear**

![Formulario de creación de un servidor de juego para Hytale](/images/en/tutorials/hytale/create_form.png)

* En el campo "Nombre", introduzca cualquier nombre para su servidor.
* En el campo "Juego", seleccione la opción Hytale.
* En el campo "Servidor dedicado", seleccione el nodo deseado donde se alojará el servidor de juego.
* En el campo IP, seleccione la dirección deseada para su servidor; a continuación, puede elegir un puerto disponible o utilizar el sugerido.
* Introduzca el puerto del servidor de juego; el valor predeterminado es 5520. No es necesario introducir los puertos rcon y query.

### Configuración después del primer inicio

Después del primer inicio, debe completar varios pasos de autorización.

#### Autorización del descargador de archivos del servidor de juego

Después del primer inicio, debe autorizar el dispositivo en el que se ejecuta el servidor de juego.
En la consola del servidor de juego, busque la línea con el enlace de autorización, copie el enlace y ábralo en su navegador.

![Enlace de autorización en la consola del servidor de Hytale](/images/en/tutorials/hytale/auth_download.png)

> ¡Nota! Para descargar los archivos del servidor de juego, debe ser propietario del juego.
> Si no es propietario del juego, verá el mensaje
> "error fetching manifest: could not get signed URL for manifest: could not get signed URL: HTTP status: 403 Forbidden"

Durante el proceso de autorización en el sitio web del desarrollador, deberá introducir un código que se enviará a su correo electrónico.

![Introducción del código de autorización en el sitio del desarrollador de Hytale](/images/en/tutorials/hytale/auth_enter_code.png)

Después de introducir el código, haga clic en "Verify".
A continuación, en la ventana emergente, haga clic en "Approve" para conceder acceso a su cuenta.

![Concesión de acceso con el botón Approve](/images/en/tutorials/hytale/auth_press_approve.png)

Después de esto, verá un mensaje indicando que el dispositivo ha sido autorizado.

![Mensaje que confirma la autorización del dispositivo](/images/en/tutorials/hytale/auth_device_approved.png)

Los archivos del servidor de juego comenzarán a descargarse. Esto puede tardar algún tiempo
dependiendo de la velocidad de su conexión a internet.

![Descarga de los archivos del servidor de juego de Hytale](/images/en/tutorials/hytale/files_downloading.png)

#### Autorización del servidor de juego

Una vez descargados los archivos del servidor de juego, debe autorizar el propio servidor de juego.
Para ello, introduzca el comando `/auth login device` en la consola del servidor de juego.

![Comando de autorización del servidor de juego en la consola de Hytale](/images/en/tutorials/hytale/auth_server.png)

Después de esto, aparecerá en la consola del servidor de juego una línea con un enlace de autorización,
similar al proceso de autorización del descargador de archivos.
Copie el enlace y ábralo en su navegador, luego repita los mismos pasos que en la autorización del descargador de archivos.

Si el proceso se realiza correctamente, verá en la consola del servidor de juego un mensaje indicando
que la autorización se completó con éxito (`Authentication successful! Mode: OAUTH_DEVICE`)

![Mensaje de consola que confirma la autorización del servidor de juego](/images/en/tutorials/hytale/auth_server_success.png)

#### Autorización automática en cada inicio

De forma predeterminada, tras la autorización, deberá introducir el comando `/auth login device`
para autorizar el servidor de juego en cada inicio. Verá este mensaje:
```text
WARNING: Credentials stored in memory only - they will be lost on restart!
To persist credentials, run: /auth persistence <type>
Available types: Memory, Encrypted
```

Para evitar introducir el comando de autorización cada vez, active la persistencia de credenciales ejecutando el comando:
```text
/auth persistence Encrypted
```

Después de esto, aparecerá un archivo `auth.enc` en el directorio raíz del servidor de juego,
que almacenará los datos de autorización cifrados.

![Archivo que contiene los datos de autorización cifrados de Hytale](/images/en/tutorials/hytale/auth_enc_file.png)

### Configuración del servidor de juego

#### Configuración mediante archivos

La mayoría de los ajustes del servidor de juego y de los mundos de Hytale se configuran mediante archivos de configuración.


| Ruta                | Descripción                               |
|---------------------|-------------------------------------------|
| .cache/             | Caché del servidor de juego               |
| logs/               | Registros del servidor de juego           |
| mods/               | Mods instalados                           |
| universe/           | Archivos de mundos y datos de jugadores   |
| bans.json           | Jugadores baneados                        |
| config.json         | Configuración principal del servidor de juego |
| permissions.json    | Configuración de permisos                 |
| whitelist.json      | Lista blanca                              |

Los ajustes principales se encuentran en el archivo `config.json`.

![El archivo config.json con los ajustes principales del servidor de Hytale](/images/en/tutorials/hytale/main_config.png)

Aquí puede configurar los siguientes parámetros:
* `ServerName` — el nombre de su servidor que se mostrará en la lista de servidores.
* `MaxPlayers` — el número máximo de jugadores que pueden estar en el servidor al mismo tiempo.
* `ServerPassword` — contraseña para acceder a su servidor si desea hacerlo privado.
* `MaxViewRadius` — la distancia máxima a la que los jugadores pueden ver objetos en el mundo.
   Establecer un valor demasiado alto puede afectar negativamente al rendimiento del servidor.
   La distancia de visión es el principal factor que afecta al consumo de RAM.
* `MOTD` — mensaje del día que se mostrará a los jugadores al conectarse al servidor.


Los ajustes y datos de los mundos se encuentran en el directorio `universe/worlds/`, que contiene las carpetas de los mundos.
Cada carpeta de mundo contiene un archivo `config.json` con los ajustes del mundo, por ejemplo `universe/worlds/default/config.json`:

![El archivo config.json con los ajustes del mundo de Hytale](/images/en/tutorials/hytale/world_config.png)

Aquí puede configurar la Seed (semilla de generación del mundo), varios parámetros de generación, ajustes de chunks, NPC, PVP
y otros parámetros.

### Enlaces útiles

* [Sitio web oficial de Hytale](https://hytale.com/)
* [Manual en el sitio web oficial](https://support.hytale.com/hc/en-us/articles/45326769420827-Hytale-Server-Manual)
