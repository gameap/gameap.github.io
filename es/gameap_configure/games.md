---
title: Juegos
layout: default
lang: es
category: Configuración del panel
order: 320
---

## Adición de nuevos juegos

El panel de control permite iniciar y realizar el control básico de cualquier servidor de juego y aplicación. Si
el panel no contiene el juego que necesita, puede agregarlo.

Para agregar un nuevo juego, vaya a **"Administración"** → **"Juegos"** y seleccione 
**"Agregar juego"**.

Después de agregar el juego, agregue el primer mod de este juego y especifique los parámetros de inicio del servidor de juego.

### Campos {#fields}

#### Código

El código del juego es un valor único. Normalmente es una abreviatura del nombre del juego, por ejemplo para 
"7 Day To Die" el código es `7d2d`. El código debe ser único. Puede introducir cualquiera, pero debe ser único.

#### Código de inicio

El código de inicio se suele especificar en los parámetros de inicio. Este valor puede ser cualquiera, pero con mayor frecuencia coincide con el código 
del juego. Puede coincidir con el de otros juegos.

#### Nombre del juego

Simplemente introduzca el nombre completo del juego.

#### Motor del juego

El motor del juego es el sistema sobre el que está escrito el juego. Half-Life y Counter-Strike están escritos en el motor GoldSource.

Half-Life 2 y Counter-Strike Source están escritos en Source. A veces no conoce el motor del juego;
en este caso, puede introducir el código del juego o alguna abreviatura del juego.

Para los juegos escritos en GoldSource y Source, introduzca el nombre exacto del motor. Para los juegos escritos en Unity, es mejor introducir
alguna abreviatura del juego, por ejemplo Rust está escrito en Unity, pero el nombre del motor es `rust`.
Minecraft está escrito sin usar ningún motor, por lo que el nombre del motor también es "Minecraft".

#### Versión del motor

Valor numérico o de cadena de la versión del motor. Puede especificar cualquier valor semántico, por ejemplo, "legacy", "beta",
etc.

#### Steam App ID

ID del servidor de juego en Steam. Se utiliza para instalar el servidor a través de SteamCMD.
Puede encontrar el SteamID en la [wiki oficial de Steam](https://developer.valvesoftware.com/wiki/Dedicated_Servers_List), o
en [SteamDB](https://steamdb.info/).

#### Steam App Set Config

Opciones adicionales para instalar el servidor a través de SteamCMD. Puede encontrar algunos valores en la 
[wiki oficial de Steam](https://developer.valvesoftware.com/wiki/Dedicated_Servers_List)

#### Repositorio local {#local-repository}

Ruta en el servidor dedicado al archivo o al directorio con la build del juego que se usa como plantilla.
Esta ruta debe existir en el servidor dedicado donde se instala un nuevo servidor de juego. 
Al agregar un nuevo servidor de juego, el archivo se descomprimirá en el directorio de trabajo del servidor de juego.

Se admiten los siguientes archivos: Zip, 7z, Tar, XZ, Bzip, GZip. 

Los archivos RAR no se admiten.

##### Ejemplos

Directorio de trabajo de ejemplo del servidor de juego `/srv/gameap/servers/example-server`

| Valor del campo "Repositorio local" | Validez del valor | Resultado de la instalación
| ------ | ------- | ------ |
| `/srv/gameap/repo/cs16_gungame.zip` | Válido, si el archivo existe en el servidor dedicado | El contenido del archivo `cs16_gungame.zip` se descomprimirá en `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_public` | Válido, si el directorio existe en el servidor dedicado | El contenido del directorio se copiará a `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_gungame.rar` | No válido. Los archivos RAR no se admiten | El método de instalación desde el repositorio local se omitirá o el servidor no se instalará
| `https://cdn.gameap.com/cstrike-1.6/rehlds-amxx-reunion.tar.xz` | No válido. Se ha especificado un valor para el repositorio remoto | El método de instalación desde el repositorio local se omitirá o el servidor no se instalará


#### Repositorio remoto {#remote-repository}

Enlace a una fuente remota. Debe ser una URL a un recurso HTTP o FTP. El archivo debe estar accesible mediante un enlace directo,
sin páginas intermedias que requieran espera o acciones adicionales. No se admiten enlaces a Yandex Disk, Google Drive, 
etc.

Las builds listas para usar se encuentran en el repositorio de GameAP: `cdn.gameap.com` para todo el mundo y `cdn.gameap.ru`
para Rusia. No es posible explorar la lista de archivos: el listado de directorios está desactivado y solo funcionan los enlaces
directos. En lugar de buscar estas URL manualmente, use el botón **Actualizar juegos**: las configuraciones de juegos
incluidas ya contienen las URL correctas, consulte
[Importación de juegos](/es/gameap_configure/games_import.html#actualización-de-juegos-desde-el-catálogo-de-gameap).

##### Ejemplos

Directorio de trabajo de ejemplo del servidor de juego `/srv/gameap/servers/example-server`

| Valor del campo "Repositorio remoto" | Validez del valor | Resultado de la instalación
| ------ | ------- | ------ |
| `https://cdn.gameap.com/cstrike-1.6/rehlds-amxx-reunion.tar.xz` | Válido | `rehlds-amxx-reunion.tar.xz` se descargará y se descomprimirá en `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_gungame.zip` | No válido. Debe especificarse un recurso http o ftp | El método de instalación se omitirá o el servidor de juego no se instalará.

## Adición de nuevos mods

Cada juego puede tener muchos mods; cada uno tiene sus propias características, configuraciones, 
parámetros de inicio, archivos de configuración, etc.

Para agregar un nuevo juego, vaya a **"Administración"** → **"Juegos"** y seleccione 
**"Agregar mod"**

### Campos

##### Juego

El juego al que pertenece el mod.

##### Nombre

Nombre del mod. Puede ser el nombre del addon, de la build, del kernel, de cualquier característica, etc. Introduzca el nombre completo
a su discreción.

#### Repositorios

El mod puede tener su propio conjunto de archivos y configuraciones que se escriben sobre los archivos de la build principal.
En el archivo puede incluir cualquier plugin adicional, contenido adicional, sonidos, música, etc.
Los campos de repositorio local son opcionales.

Se admiten los siguientes archivos: Zip, 7z, Tar, XZ, Bzip, GZip. 

Los archivos RAR no se admiten.

##### Repositorio local {#local-repository-1}

Ruta en el servidor dedicado al archivo o al directorio con los archivos del mod del juego que se usan como plantilla.
Esta ruta debe existir en el servidor dedicado donde se instala un nuevo servidor de juego. 
Al instalar el servidor de juego, el archivo se descomprimirá sobre la build principal en el directorio de trabajo 
del servidor de juego.

Se admiten los siguientes archivos: Zip, 7z, Tar, XZ, Bzip, GZip. 

Los archivos RAR no se admiten.

###### Ejemplos

Directorio de trabajo de ejemplo del servidor de juego `/srv/gameap/servers/example-server`

| Valor del campo "Repositorio local" | Validez del valor | Resultado de la instalación
| ------ | ------- | ------ |
| `/srv/gameap/repo/cs16_gungame.zip` | Válido, si el archivo existe en el servidor dedicado | El contenido del archivo `cs16_gungame.zip` se descomprimirá en `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_cs16_gungame` | Válido, si el directorio existe en el servidor dedicado | El contenido del directorio se copiará a `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_gungame.rar` | No válido. Los archivos RAR no se admiten | La descompresión del archivo se omitirá
| `https://cdn.gameap.com/cstrike-1.6/rehlds-amxx-reunion.tar.xz` | No válido. Se ha especificado un valor de repositorio remoto | La instalación del mod se omitirá


##### Repositorio remoto {#remote-repository-1}

Enlace a una fuente remota. Debe ser una URL a un recurso HTTP o FTP. El archivo debe estar accesible mediante un enlace directo,
sin páginas intermedias que requieran espera o acciones adicionales. No se admiten enlaces a Yandex Disk, Google Drive, 
etc.

Las builds listas para usar se encuentran en el repositorio de GameAP: `cdn.gameap.com` para todo el mundo y `cdn.gameap.ru`
para Rusia. No es posible explorar la lista de archivos: el listado de directorios está desactivado y solo funcionan los enlaces
directos. En lugar de buscar estas URL manualmente, use el botón **Actualizar juegos**: las configuraciones de juegos
incluidas ya contienen las URL correctas, consulte
[Importación de juegos](/es/gameap_configure/games_import.html#actualización-de-juegos-desde-el-catálogo-de-gameap).

###### Ejemplos

Directorio de trabajo de ejemplo del servidor de juego `/srv/gameap/servers/example-server`

| Valor del campo "Repositorio remoto" | Validez del valor | Resultado de la instalación
| ------ | ------- | ------ |
| `https://cdn.gameap.com/cstrike-1.6/rehlds-amxx-reunion.tar.xz` | Válido | `rehlds-amxx-reunion.tar.xz` se descargará y se descomprimirá en `/srv/gameap/servers/example-server`
| `/srv/gameap/repo/cs16_gungame.zip` | No válido. Debe especificarse un recurso http o ftp | La instalación del mod se omitirá

## Edición de mods

Después de crear un mod, puede configurarlo más a fondo especificando parámetros adicionales como 
los "comandos de inicio predeterminados", variables de inicio y diversos comandos RCON.

### Configuración principal

La configuración principal incluye parámetros como el nombre del mod, los repositorios y los comandos de inicio predeterminados. Preste especial atención
a los comandos de inicio.

#### Nombre del mod

#### Repositorios

Lea más en [Adición de nuevos mods, repositorios](#repositorios)

#### Comandos de inicio predeterminados

Al agregar un nuevo servidor de juego para el mod seleccionado, este comando se le asignará automáticamente. Si
el comando está vacío, no se asignará nada y deberá especificar el comando de inicio para cada servidor de juego manualmente.
El servidor de juego no puede iniciarse sin un comando de inicio.

Puede usar shortcodes en los comandos de inicio; luego se reemplazarán por los valores de las variables del servidor. Los shortcodes
son palabras sin espacios entre llaves `{`, `}`; por ejemplo, `{ip}`, `{port}`, `{maxplayers}`, etc.

##### Shortcodes básicos

Estos shortcodes están siempre disponibles, no requieren agregar variables adicionales en la configuración del mod.

| Shortcode | Descripción
| ------ | -------
| {ip} | IP del servidor de juego
| {port} | Puerto principal del servidor de juego. A veces llamado puerto de conexión
| {query_port} | Puerto de consulta
| {rcon_port} | Puerto de comunicación del servidor (puerto RCON)
| {rcon_password} | Contraseña RCON
| {uuid} | UUID del servidor
| {uuid_short} | UUID corto del servidor

##### Shortcodes definidos por el usuario

Puede definir estos shortcodes usted mismo para cada mod de juego concreto. Dependiendo de la 
configuración individual del servidor, estos shortcodes se reemplazarán por los valores de los parámetros del servidor de juego. Lea más 
en [Variables](#variables).

### Variables

Puede agregar configuraciones individuales para cada servidor de juego. Luego, estas configuraciones pueden ser editadas
por el administrador o por un usuario normal en la página de configuración (**"Lista de servidores"** → **"Administración"** → 
**"Configuración"**).

| Campo | Descripción
| ------ | -------
| Variable | Nombre de la variable. Sin llaves.
| Predeterminado | Valor de la variable por defecto. Este valor se usará si no se ha establecido un valor individual para el servidor de juego.
| Descripción | Descripción del servidor de juego en la página de configuración (**"Lista de servidores"** → **"Administración"** → **"Configuración"**)
| Variable de administrador | Si está marcado, solo el administrador podrá editar esta configuración para los servidores de juego.

#### Ejemplos

| Valores | Descripción
| ------ | -------
| **Variable:** default_map <br><br>**Predeterminado:** de_dust <br><br>**Descripción:** Mapa al inicio | Para cada servidor de juego de este mod, aparecerá el shortcode `{default_map}` y un nuevo parámetro en la configuración llamado "Mapa predeterminado" con el valor predeterminado "de_dust". <br><br>De forma individual para cada servidor, este parámetro se puede editar en **"Lista de servidores"** → **"Administración"** → **"Configuración"**


### Comandos RCON

Estos comandos permiten una administración más avanzada del servidor de juego. Si el juego admite trabajar con RCON o 
si se admite la consola, puede hacer lo siguiente: expulsar jugadores del servidor, banear jugadores, cambiar el mapa del servidor,
enviar mensajes de texto al chat común, establecer una contraseña.

Algunas funciones pueden estar limitadas por el propio juego o por el mod. Por ejemplo, no todos los servidores de juego admiten
el acceso al servidor con contraseña.

#### Comando de expulsión

Puede establecer el comando RCON para expulsar a un jugador del servidor.

Puede establecer shortcodes para el comando que se reemplazarán con los datos de un jugador concreto

| Shortcode | Descripción
| ------ | -------
| {id} | ID del jugador en el servidor
| {name} | Nombre del jugador en el servidor

Para muchos juegos GoldSource/Source, este es el comando: 
```text
kick #{id}
```

#### Comando de ban

Puede establecer un comando que se usará para el ban temporal o permanente de un jugador en el servidor.

Puede establecer shortcodes para el comando que se reemplazarán con los datos de un jugador concreto

| Shortcode | Descripción
| ------ | -------
| {id} | ID del jugador en el servidor
| {name} | Nombre del jugador en el servidor
| {time} | Tiempo del ban
| {reason} | Razón del ban

Para muchos juegos GoldSource (Half-Life, Counter-Strike 1.6, etc.) que ejecutan AMX Mod X, este es el comando: 
```text
amx_ban "{name}" {time} "{reason}"
```

#### Comando de cambio de nombre (nick)

Con este comando puede cambiar el nick del jugador seleccionado.

Puede establecer shortcodes para el comando que se reemplazarán con los datos de un jugador concreto

| Shortcode | Descripción
| ------ | -------
| {id} | ID del jugador en el servidor
| {name} | Nombre actual del jugador en el servidor
| {new_name} | Nuevo nombre del jugador en el servidor
| {reason} | Razón del cambio

Para muchos juegos GoldSource (Half-Life, Counter-Strike 1.6, etc.) que ejecutan AMX Mod X, este es el comando: 
```text
amx_nick #{id} {new_name}
```

#### Comando de reinicio

Puede establecer el comando para un reinicio suave del servidor, sin reiniciar el proceso del servidor de juego. Normalmente,
este comando reinicia el mapa o la ronda del juego. No es compatible con muchos juegos.

Para muchos juegos GoldSource/Source, este es el comando: 
```text
restart
```

#### Comando de cambio de mapa

Con este comando puede cambiar el mapa del servidor de juego.

| Shortcode | Descripción
| ------ | -------
| {map} | Nombre del mapa

Para muchos juegos GoldSource/Source, este es el comando: 
```text
changelevel {map}
```

#### Comando de envío de mensajes

Con este comando puede enviar un mensaje de texto al chat de todos los jugadores del servidor.

| Shortcode | Descripción
| ------ | -------
| {msg} | Mensaje que se enviará al servidor

Para muchos juegos GoldSource (Half-Life, Counter-Strike 1.6, etc.) que ejecutan AMX Mod X, este es el comando: 
```text
amx_say "{msg}"
```

#### Comando para establecer/cambiar la contraseña

Puede establecer una contraseña para el servidor de juego, de modo que solo los jugadores que conozcan
esta contraseña puedan entrar.

| Shortcode | Descripción
| ------ | -------
| {password} | Contraseña del servidor

Para muchos juegos GoldSource/Source, este es el comando: 
```text
password {password}
```

### Comandos FastRCON

Puede especificar sus propios comandos RCON opcionales. Por ejemplo, un comando de estado del servidor, obtención de estadísticas, obtención de la
lista de jugadores desconectados recientemente, etc.
