---
title: Solución de problemas
layout: default
lang: es
category: Solución de problemas
order: 400
---

Esta página describe algunos posibles errores y cómo solucionarlos.

## Por dónde empezar

Para casi cualquier problema con un servidor de juego, los detalles se pueden encontrar en dos
lugares.

**Tareas en el panel.** **Administración** → **Tareas de GDaemon**, abra la última tarea: contiene
el resultado del comando y la salida del error.

**El registro del daemon** en el servidor dedicado:

* Linux — `/var/log/gameap-daemon/output.log`
* Windows — `C:\gameap\daemon\logs\output.log`

Si el daemon se ejecuta bajo systemd, la misma salida está disponible con:

```shell
journalctl -u gameap-daemon -n 200 --no-pager
```

Para más detalles, establezca `log_level: debug` en la configuración del daemon y reinícielo.

## Errores de inicio del servidor

### El estado del servidor se muestra incorrectamente

A veces el servidor se inicia, pero el panel lo muestra como desconectado.
Consulte [Errores de visualización del estado del servidor](#errores-de-visualización-del-estado-del-servidor).

### Comando de inicio del servidor vacío

Este error ocurre cuando el comando de inicio del servidor de juego está vacío.

Vaya a la página de administración del servidor de juego: **Administración** → **Servidores** →
seleccione el servidor de juego. O desde la página principal: **Lista de servidores** → seleccione
el servidor → **Control** → **Administración**.

Busque el campo **Comando de inicio del servidor de juego** e introduzca el comando. Para
Counter-Strike 1.6 será algo así:

```text
./hlds_run -game cstrike +ip {ip} +port {port} +map {default_map} +maxplayers {maxplayers} +sys_ticrate {fps}
```

Consulte la [configuración de servidores de juego](/es/gameap_configure/game_servers.html#comando-de-inicio)
para más detalles. El comando de inicio predeterminado se puede establecer
[en la configuración del mod](/es/gameap_configure/games.html#comandos-de-inicio-predeterminados).

### La ventana modal de información no cambia durante mucho tiempo

Un inicio suele tardar menos de 10 segundos. Si la barra de progreso está congelada o el estado no
cambia durante varios minutos, compruebe si el daemon está en ejecución y conectado al panel:

```shell
systemctl status gameap-daemon
```

Si el daemon no está en ejecución:

```shell
systemctl start gameap-daemon
```

Si está en ejecución, revise el registro para ver si la conexión con el panel está establecida. Los
mensajes `gRPC connection failed` y `Reconnecting to panel...` significan que el daemon no pudo
comunicarse con el panel; consulte
[El daemon no se conecta al panel](#el-daemon-no-se-conecta-al-panel).

### Server start task is already exists

Este error ocurre cuando ya se ha creado una tarea de inicio y todavía no se ha ejecutado.

Vaya a la página **Tareas de GDaemon**, busque la tarea en estado de espera, ábrala con el botón
**Ver** y haga clic en **Cancelar**. Después inicie el servidor de nuevo.

Si las tareas se quedan en espera con regularidad, el daemon no las está recogiendo: compruebe que
está en ejecución y conectado al panel.

## Errores de visualización del estado del servidor

### Desincronización de la hora

El problema se debe a una diferencia de hora entre el servidor del panel y el servidor dedicado.
Configure la sincronización de la hora en ambos.

Cambiar la zona horaria en Debian y Ubuntu:

```bash
dpkg-reconfigure tzdata
```

### El daemon está conectado, pero el estado no se actualiza

El daemon informa del estado de los servidores cada `metrics.collection_interval` (5 segundos por
defecto) y envía un heartbeat cada 30 segundos. Si los datos permanecen desactualizados durante más
tiempo, revise el registro del daemon en busca de caídas de la conexión.

## Errores de instalación del servidor

Para encontrar la causa, mire el resultado del comando de instalación: **Administración** →
**Tareas de GDaemon** → busque la tarea de instalación y ábrala.

Causas comunes:

* [No hay fuente](#no-hay-fuente)
* [Archivo de instalación formado incorrectamente](#archivo-de-instalación-formado-incorrectamente)
* [Fuente de instalación incorrecta](#fuente-de-instalación-incorrecta)

### No hay fuente

El panel no sabe desde dónde instalar el servidor de juego. Hay varias formas:

* mediante SteamCMD: el Steam APP ID debe estar configurado en los ajustes del juego;
* desde un repositorio remoto: debe proporcionarse un enlace a un archivo ZIP o TAR (RAR no está
  soportado);
* desde un repositorio local: una ruta a un directorio con archivos o a un archivo ZIP o TAR. Los
  archivos deben estar ubicados en el servidor dedicado donde se ejecuta el daemon.

La fuente se configura en la página **Administración** → **Juegos** → seleccione el juego →
**Editar**.

### Archivo de instalación formado incorrectamente

El archivo debe contener los archivos del servidor de juego en su raíz. El error más común es que
los archivos estén en un directorio anidado.

Un archivo formado incorrectamente para GTA: San Andreas Multiplayer:

![Archivo de instalación incorrecto: los archivos del servidor de juego están en un directorio anidado](/images/errors/source_archive_wrong.jpg)

Correcto:

![Archivo de instalación correcto: los archivos del servidor de juego están en la raíz del archivo](/images/errors/source_archive_right.jpg)

### Fuente de instalación incorrecta

Para un repositorio local, compruebe que el directorio existe en el servidor dedicado donde se
ejecuta el daemon.

Para uno remoto, compruebe que el enlace realmente inicia la descarga de un archivo. Los enlaces a
Yandex Disk, Google Drive y servicios de almacenamiento similares no están soportados: sirven una
página, no un archivo.

### Failed to install via steamcmd

A veces SteamCMD se interrumpe en mitad de una descarga. El panel puede trabajar con el descargador
alternativo depot downloader: sustituya el script de SteamCMD:

```shell
cd /srv/gameap/steamcmd
mv steamcmd.sh steamcmd.sh.orig
curl -O https://raw.githubusercontent.com/gameap/steamcmd-depotdownloader/main/steamcmd.sh
chmod 755 steamcmd.sh && chown gameap:gameap steamcmd.sh
```

### Se requiere iniciar sesión en una cuenta de Steam

Algunos juegos no se pueden descargar de forma anónima: se necesita una cuenta con una copia
comprada. Especifíquela en la sección `steam_config` de la configuración del daemon; consulte
[GameAP Daemon](/es/daemon/daemon.html#cuenta-de-steam).

La autenticación de dos factores debe estar desactivada en esta cuenta; de lo contrario, el daemon
no podrá iniciar sesión.

## Errores de GameAP Daemon

### El daemon no se conecta al panel

En GameAP 4 la conexión la establece el daemon: se conecta al panel por gRPC en el puerto **31718**.
El panel no se conecta al daemon, y para ello no se necesitan puertos entrantes en el servidor
dedicado.

La señal reveladora en el registro del daemon son los mensajes repetidos `gRPC connection failed` y
`Reconnecting to panel...`.

Qué comprobar:

1. **Accesibilidad del puerto** desde el servidor dedicado:

   ```shell
   nc -zv panel.example.com 31718
   ```

   Si la conexión no se puede establecer, el puerto está bloqueado por un firewall o el panel está
   escuchando en otra dirección.

2. **La dirección del panel** en la configuración del daemon: el parámetro `grpc.address`. Debe
   apuntar a una dirección accesible desde el servidor dedicado.

3. **Error de verificación del certificado.** Si el registro contiene un mensaje sobre una
   discrepancia de nombre en el certificado, el daemon se está conectando mediante una dirección que
   no está en el certificado del panel. Establezca `GRPC_EXTERNAL_HOST` en el panel y reemita el
   certificado; consulte [GRPC API](/es/daemon/grpc.html).

4. **`registration failed`.** La conexión se establece, pero el panel rechazó el registro: un
   `ds_id` o `api_key` incorrecto en la configuración del daemon. La solución más sencilla es
   registrar el daemon de nuevo.

### La consola o el gestor de archivos no funcionan

Si los servidores se inician y se detienen, pero la consola y el gestor de archivos no funcionan,
el problema es casi siempre la conexión con el panel: los comandos de control pueden haberse
completado antes de una desconexión, mientras que la consola y los archivos requieren un flujo de
datos activo.

Compruebe la conexión como se describe en la sección anterior; reinicie el daemon si es necesario:

```shell
systemctl restart gameap-daemon
```

### No se puede subir un archivo mediante el gestor de archivos

Por defecto, no se permite subir archivos comprimidos ni archivos binarios arbitrarios: el tipo se
detecta a partir del contenido del archivo, no de la extensión. Se pueden permitir con las
variables `FILES_UPLOAD_ALLOW_ARCHIVES` y `FILES_UPLOAD_ALLOW_BINARY`; consulte
[Seguridad](/es/security.html).

El tamaño máximo de archivo para una subida normal es de 100 MB.

### El daemon no se inicia tras reiniciar el servidor

Active el inicio automático:

```shell
systemctl enable gameap-daemon
```

### Task complete with an error

Un error genérico durante la instalación, el inicio, el reinicio o la detención de un servidor de
juego. Puede haber muchas causas.

Abra **Administración** → **Tareas de GDaemon**, busque la última tarea y mire los detalles.
Después revise el registro del daemon.

#### No source to install game

No se ha configurado ninguna fuente de instalación para el juego. Consulte
[No hay fuente](#no-hay-fuente).

## Errores de inicio de sesión en el panel

### Un administrador no puede iniciar sesión después de una actualización

En GameAP 4, la autenticación de dos factores es obligatoria para los administradores: 30 días
después del primer recordatorio, el inicio de sesión deja de emitir una sesión completa. Esto no es
un bloqueo de la cuenta: active 2FA e inicie sesión de nuevo.

Los pasos, incluido qué hacer cuando se pierde el dispositivo con los códigos, están en la página
[Seguridad](/es/security.html).

### Demasiados intentos de inicio de sesión

Después de 20 intentos fallidos desde una dirección o 5 intentos para un mismo usuario, el inicio
de sesión se bloquea durante 15 minutos con una respuesta `429`. El límite se elimina por sí solo
cuando pasa el tiempo.

### La contraseña se rechaza al cambiarla

La longitud mínima de la contraseña es de 12 caracteres, y la contraseña se comprueba contra la
lista de contraseñas comprometidas. Consulte [Seguridad](/es/security.html).

## Búsqueda de la causa

Si su problema no aparece aquí, busque detalles en los registros.

**El registro del daemon**: `/var/log/gameap-daemon/output.log` en Linux,
`C:\gameap\daemon\logs\output.log` en Windows. Para una salida detallada, establezca
`log_level: debug` en la configuración del daemon y reinícielo.

**El registro del panel**: la salida del proceso `gameap`. Bajo systemd:

```shell
journalctl -u gameap -n 200 --no-pager
```

El nivel de detalle se configura con la variable `LOGGER_LEVEL`; consulte la
[Referencia de config.env](/es/config.html).

**Tareas**: en el panel, **Administración** → **Tareas de GDaemon**: ahí está el resultado de cada
comando de inicio, instalación y reinicio.
