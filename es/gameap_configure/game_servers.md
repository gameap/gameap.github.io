---
title: Servidores de juego
layout: default
lang: es
category: Configuración del panel
order: 310
---

## Página del servidor de juego

Abra la **Lista de servidores** y seleccione un servidor. La página está dividida en pestañas; cuáles están presentes
depende de los permisos otorgados al usuario y de las capacidades del juego.

| Pestaña                | Qué contiene                                                               |
|--------------------|------------------------------------------------------------------------------|
| **Control**        | Botones de inicio, detención, reinicio, actualización y reinstalación, estado del servidor, consola y estadísticas |
| **RCON**           | Consola RCON y gestión de jugadores. Aparece si el juego admite RCON        |
| **Archivos**          | El gestor de archivos del servidor                                                    |
| **Planificador de tareas** | Tareas periódicas: reinicio, actualización, comandos arbitrarios                         |
| **Configuración**       | Valores de las variables del mod para este servidor                                  |

Los plugins pueden añadir sus propias pestañas al servidor: aparecen junto a las integradas.

### Control

La pestaña principal. Además de los botones de control, contiene:

**Estado del servidor** — en ejecución o detenido, el mapa actual y el número de jugadores, si el juego informa
de estos datos a través del protocolo Query.

**Consola** — la salida del servidor de juego y un campo para enviar comandos. No funciona a través de
RCON, sino a través del gestor de procesos en el servidor dedicado, por lo que está disponible incluso cuando el
servidor aún no se ha iniciado.

**Estadísticas** — gráficos de uso de CPU, memoria y red. Los datos son recopilados por el daemon
cada `metrics.collection_interval` (5 segundos de forma predeterminada) y se conservan durante
`metrics.retention_duration` (10 minutos de forma predeterminada), consulte
[GameAP Daemon](/es/daemon/daemon.html#metrics-collection).

### Archivos

El gestor de archivos funciona dentro del directorio del servidor de juego. Permite ver y editar archivos,
subirlos y descargarlos, crear directorios y cambiar permisos.

Las subidas están restringidas por tipo de archivo: el tipo se detecta a partir del contenido, no de la extensión,
y de forma predeterminada **los archivos comprimidos y los archivos binarios arbitrarios no están permitidos**. Esta es la razón más común
por la que se rechaza una subida. Cómo permitirlos — [Seguridad](/es/security.html).

### Planificador de tareas

Tareas periódicas para el servidor: reinicio programado, actualización, ejecución de comandos. Las tareas son
ejecutadas por el daemon en el servidor dedicado.

### Configuración

Valores de las variables declaradas en el mod del juego — por ejemplo, el mapa predeterminado o el número
de slots. Se sustituyen en el comando de inicio en lugar de los shortcodes. Qué variables están
disponibles se define [en la configuración del mod](/es/gameap_configure/games.html#variables).

## Edición de servidores de juego

Para editar un servidor de juego (nodo), vaya a la página **"Administración"** → **"Servidores de juego"**, luego
seleccione el servidor de juego que desea editar y haga clic en el botón **"Editar"**.

### Descripción de los parámetros

#### Básicos

Grupo básico de parámetros

##### UUID

UUID es un identificador único universal. Se genera automáticamente para cada servidor. Se utiliza como
identificador de los procesos del servidor de juego.

##### Nombre

Nombre del servidor. Puede ser cualquier cadena de texto. Puede introducir cualquier nombre de servidor de juego.

##### Juego

El juego al que pertenece el servidor de juego. Lea la documentación sobre la configuración de juegos en la sección correspondiente —
[Configuración de juegos](/es/gameap_configure/games.html).

Consulte también [Manual para añadir juegos que faltan](/es/tutorials/additional_games.html).

##### Mod

El mod al que pertenece el servidor de juego. El mod puede determinar configuraciones adicionales (variables)
para el servidor de juego, como el mapa predeterminado (`{default_map}`), FPS y otras. Consulte la información sobre cómo
añadir sus propias configuraciones y otros detalles al respecto en la
[página de configuración de juegos](/es/gameap_configure/games.html#variables).

##### Contraseña RCON

Contraseña para administrar el servidor de juego a través de RCON.

##### Directorio

Directorio del servidor de juego relativo al [directorio de trabajo del servidor dedicado](/es/gameap_configure/dedicated_servers.html#working-directory).

Por ejemplo, si el directorio es `servers/my_server` y el directorio de trabajo del servidor dedicado es `/srv/gameap`, entonces el
servidor de juego se ubicará en `/srv/gameap/servers/my_server`.

##### Nombre de usuario en el servidor dedicado

El usuario que ejecuta el servidor de juego en el nodo. De forma predeterminada es `gameap`.

El usuario especificado en este campo debe existir en el servidor dedicado.

#### Servidor dedicado, IP, puertos

Un grupo de parámetros relacionados con el servidor dedicado (VDS, nodo) y la conexión al servidor de juego.

##### Servidor dedicado

##### IP

IP o host del servidor de juego. Ejemplos: `127.0.0.1`, `my-server.gameap.ru`.

##### Puerto del servidor

Puerto principal del servidor. Se utiliza para conectar a los jugadores al servidor.

##### Puerto Query

Puerto del servidor para consultas. Query se utiliza para obtener datos generales del servidor: mapa actual, número actual
de jugadores en el servidor, lista de jugadores.

En GoldSource y Source, coincide con el puerto principal del servidor. En Minecraft, puede especificar cualquier puerto. En algunos
otros servidores de juego, puede ser una unidad más o menos.

##### Puerto RCON

Puerto del servidor para la administración remota.

En GoldSource y Source, coincide con el puerto principal del servidor. En Minecraft, puede especificar cualquier puerto. En algunos
otros servidores de juego, puede ser una unidad más o menos.

#### Límites de recursos

Límites sobre el consumo de recursos del servidor de juego.

##### Límite de CPU

Se establece en millicores: `1000` es un núcleo completo, `500` es medio núcleo. En la interfaz el valor
también puede introducirse como porcentaje o en núcleos; el panel lo convierte automáticamente.

##### Límite de RAM

La cantidad máxima de RAM. En la interfaz se introduce en bytes, megabytes o gigabytes.

> Los límites solo los aplican los gestores de procesos que los admiten: **systemd**, **Docker**
> y **Podman**. Con `tmux`, `winsw`, `shawl` o `simple`, los valores configurados no tienen
> efecto. Consulte [Gestores de procesos](/es/daemon/process_managers.html).

Un valor vacío o cero significa sin límite.

#### Comando de inicio

El comando de inicio es un parámetro importante y obligatorio para iniciar un servidor de juego. Es una cadena con diversas
opciones y parámetros de arranque del servidor de juego. Puede ser individual para cada juego y mod.

Los parámetros pueden incluir shortcodes que se reemplazan automáticamente por el valor de la configuración (variables)
del servidor de juego específico. Los shortcodes son palabras entre llaves `{` y `}`, por ejemplo, `{ip}`,
`{port}`, `{default_map}` y otros.

Lea cómo añadir sus propios shortcodes, que se reemplazarán automáticamente por el valor de la configuración, en
la página de [configuración de juegos](/es/gameap_configure/games.html#variables).

## Envío de comandos a la consola

Un comando introducido en la consola en la pestaña **Control** se entrega al servidor de juego de una
de dos maneras.

Si el servidor dedicado tiene configurada la plantilla **Script Send Command** (`script_send_command`),
el panel sustituye el comando en ella y lo ejecuta en el servidor dedicado.
Si la plantilla no está configurada, el comando se escribe en el archivo de entrada del servidor de juego.

> **No ponga `{command}` entre comillas.** El panel sustituye el valor ya entrecomillado: antes de la
> sustitución, el comando se escapa para el shell y se envuelve en comillas simples. Añadir sus propias
> comillas provoca un doble escape, y el servidor recibe el comando junto con las comillas.

Correcto:

```text
tmux send-keys -t {uuid} {command} Enter
```

Incorrecto:

```text
tmux send-keys -t {uuid} "{command}" Enter
```

La plantilla se configura en la página **Administración** → **Servidores dedicados** → seleccione el servidor →
**Editar** → la pestaña **Scripts**.
