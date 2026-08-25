---
title: Añadir juegos que faltan
layout: default
lang: es
category: Tutoriales
order: 299
---

El panel de control admite el arranque y el control básico de cualquier servidor de juego y aplicación.
Este manual explica cómo añadir un nuevo juego tomando como ejemplo Sven Co-op. Cada paso incluye explicaciones.

## Añadir un juego

Primero, vaya a la página de añadir juego. Abra el menú **"Administración"** y luego seleccione
**"Juegos"**.

![El elemento Juegos en el menú de administración](/images/en/tutorials/additional_games/game_menu.png)

A continuación, busque el botón **"Añadir juego"** en la parte superior de la página y haga clic en él.

![El botón Añadir juego en la página de la lista de juegos](/images/en/tutorials/additional_games/add_game_menu.png)

Será redirigido a la página para añadir un nuevo juego. Aquí debe especificar algunos datos de su nuevo juego.

![Formulario para añadir un juego usando Sven Co-op como ejemplo](/images/en/tutorials/additional_games/example_add_svencoop.png)

Debe rellenar los siguientes campos:
* **Código**. Introduzca el nombre abreviado del juego.
* **Código de inicio**. También puede introducir el nombre abreviado del juego.
* **Nombre del juego**.
* **Motor del juego**. Si el juego está escrito en Unity o sin usar un motor, introduzca el nombre abreviado
del juego.
* **Versión**. Introduzca el número de versión o un valor semántico, por ejemplo “legacy”, “beta”, etc.

¡Nota! Para habilitar la instalación automática, debe rellenar uno de los siguientes campos:
* [**Steam APP ID**](/es/gameap_configure/games.html#steam-app-id). Consulte el valor del juego que le interesa en la
[wiki oficial de Steam](https://developer.valvesoftware.com/wiki/Dedicated_Servers_List) o en la base de datos
[SteamDB](https://steamdb.info/)
* [**Repositorio remoto**](/es/gameap_configure/games.html#remote-repository). Un enlace a un archivo por HTTP o FTP.
* [**Repositorio local**](/es/gameap_configure/games.html#local-repository).

Todos estos campos son opcionales, pero uno de ellos debe tener un valor.

Lea más sobre el significado de los campos en la página [Configuración de juegos, descripción de los campos](/es/gameap_configure/games.html#fields).

## Añadir un mod

Cada juego debe tener al menos un mod.

El mod es una potente herramienta de GameAP: permite habilitar plugins adicionales, configuración o cualquier
contenido para ampliar las capacidades básicas del servidor. El archivo con los ficheros que especifique para el mod
se descomprimirá sobre la compilación base del servidor.

Para añadir un nuevo mod a un juego concreto, seleccione el juego en la lista y haga clic en **"Añadir el primer mod"**.

![El botón Añadir el primer mod para el juego Sven Co-op](/images/en/tutorials/additional_games/example_menu_add_mod_svencoop.png)

Si el juego ya tiene al menos un mod, en la parte superior de la página de la lista de juegos seleccione
**"Añadir mod"**.

En la página para añadir un mod, especifique el nombre del mod según las características del
modo de juego (GunGame, Jail, etc.) o la disponibilidad de módulos (AMXX, ReAMXX para Counter-Strike,
IndustrialCraft, BuildCraft para Minecraft, etc.).

![Formulario de creación de mod para el juego Sven Co-op](/images/en/tutorials/additional_games/example_add_svencoop_mod.png)

Si dispone de un archivo con plugins adicionales que deben escribirse sobre la compilación base, especifique la ruta
al mismo en los campos de repositorio local o remoto.

En el campo **repositorio local**, especifique la ruta al archivo o directorio en un servidor dedicado
que ejecute GameAP Daemon; ejemplo de ruta: `/srv/gameap/repo/svencoop_op4_maps.tar.xz`.
Consulte los detalles en la página [Configuración de juegos](/es/gameap_configure/games.html#local-repository-1).

En el campo **repositorio remoto**, especifique la URL del archivo HTTP o FTP.
Ejemplo de ruta: `https://cdn.gameap.com/svencoop/svencoop_op4_maps.tar.xz`.
Consulte los detalles en la página [Configuración de juegos](/es/gameap_configure/games.html#remote-repository-1).
Hay archivos listos para usar para muchos juegos en el repositorio de GameAP (`cdn.gameap.com`,
`cdn.gameap.ru`), pero la lista de ficheros no se puede explorar: el listado de directorios está deshabilitado
y solo funcionan los enlaces directos. Normalmente es más sencillo tomar la configuración de juegos incluida con el
botón **Actualizar juegos**: las URL ya están definidas ahí.

## Configuración del mod

Después de crear un mod, puede configurarlo especificando parámetros adicionales como
"Comandos de inicio predeterminados", variables de inicio y diversos comandos RCON.

Los comandos de inicio predeterminados deben especificarse. Si no los especifica, el comando de inicio estará vacío
al crear un nuevo servidor de juego, pero es obligatorio; de lo contrario, el servidor no arrancará.

![Configuración principal del mod con el comando de inicio predeterminado](/images/en/tutorials/additional_games/game_mods_edit_basic.png)

Ejemplos de comandos de inicio predeterminados para algunos juegos en Linux:
* Sven Co-op:
```shell
./svends_run +ip {ip} +port {port} +maxplayers {maxplayers} +log on +map {default_map}
```
* Half-Life:
```shell
./hlds_run -game valve +ip {ip} +port {port} +map {default_map} +maxplayers {maxplayers} +sys_ticrate {fps}
```

* GTA: Multi Theft Auto
```shell
./mta-server64 -t -n --ip {ip} --port {port} --maxplayers {maxplayers}
```

Ejemplos de comandos de inicio predeterminados para algunos juegos en Windows:

* Sven Co-op:
```shell
SvenDS +ip  {ip} +port {port} +maxplayers {maxplayers} +log on +map {default_map}
```

* 7 Day To Die
```shell
startdedicated.bat
```

Preste atención a los valores entre llaves `{` y `}`, como `{ip}`, `{port}`, `{maxplayers}`, `{default_map}`,
`{fps}` y otros. En GameAP se llaman shortcodes; se sustituyen por los valores de las variables del servidor.
Todos los servidores de juego tienen las llamadas variables básicas, como IP, puertos, ID y UUID. También tienen
variables adicionales que se especifican en la configuración del mod, entre ellas el número máximo de jugadores,
el mapa predeterminado, los FPS y otras.

Algunos parámetros solo pueden modificarlos los administradores, y otros están disponibles para los usuarios
normales. La lista de variables y sus nombres se especifica en la configuración del mod, en la pestaña "Variables".

![La pestaña Variables en la configuración del mod](/images/en/tutorials/additional_games/game_mods_edit_vars.png)

Las variables especificadas en el mod para cada servidor de juego pueden modificarse después individualmente en la configuración.

La siguiente pestaña en la configuración del mod es "Comandos RCON". Puede especificar comandos RCON para expulsar
jugadores, banearlos, cambiar el mapa y otros comandos RCON; se utilizan para la administración avanzada del
servidor de juego.

![La pestaña Comandos RCON en la configuración del mod](/images/en/tutorials/additional_games/game_mods_edit_commands.png)

Puede especificar sus propios comandos RCON opcionales en la pestaña Fast RCON. Por ejemplo, el comando de estado
del servidor o de estadísticas.

![La pestaña Fast RCON con comandos definidos por el usuario](/images/en/tutorials/additional_games/game_mods_edit_fast_rcon.png)
