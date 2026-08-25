---
title: Instalación y gestión
layout: default
lang: es
category: Plugins
order: 341
---

## La página de Plugins

La gestión de plugins está disponible para el administrador en la página **Administración** → **Plugins**. La página tiene dos pestañas:

* **Instalados** — la lista de plugins cargados en el panel. La tabla muestra el nombre, la categoría, la valoración, el número de descargas, la versión y las acciones. Las insignias muestran el origen de la instalación (`file` — «Archivo local», `store` — «Tienda»), el estado y si hay una actualización disponible.
* **Tienda** — la lista de plugins del catálogo plugins.gameap.dev con paginación en el lado del servidor. Los plugins ya instalados reciben una insignia; los de pago reciben un botón de compra de suscripción.

## Instalación desde el catálogo

1. En la pestaña **Tienda**, haga clic en el nombre del plugin o en el botón **Instalar**: se abrirá la ficha del plugin con su descripción, autor, licencia, etiquetas, valoración y lista de versiones.
2. Seleccione una versión de la lista desplegable (por defecto se ofrece la más reciente) y haga clic en **Instalar**.
3. El panel descarga el archivo `.wasm`, verifica su hash SHA-256, lo escribe en el directorio `plugins/`, carga el plugin y recarga la página (para actualizar `/plugins.js` con los frontends de los plugins).

Para los plugins de pago (marcados como `requires_subscription`), se ofrece un botón de compra de suscripción en lugar de la instalación. Para permitir que el panel descargue plugins de pago, defina la clave de licencia en la variable de entorno `PLUGIN_STORE_LICENSE_KEY`.

## Instalación desde un archivo

1. En la pestaña **Instalados**, haga clic en el botón **Subir** y seleccione el archivo `.wasm` del plugin (de no más de 100 MB).
2. Haga clic en **Comprobar**: el panel realiza una carga de prueba del módulo (una ejecución en seco, sin instalar) y muestra los metadatos del plugin: nombre, versión, autor, versión de la Plugin API, si hay rutas HTTP y un frontend, y una lista de errores de validación.
3. Si el plugin es válido, haga clic en **Instalar**: el panel instala el plugin y recarga la página.

## Actualización y eliminación

Cuando se publica una nueva versión de un plugin instalado desde el catálogo, aparece una insignia de actualización en la tabla y la columna de versión muestra `instalada → última`. La actualización se realiza con el botón de la lista o desde la ficha del plugin: el panel descarga la nueva versión, verifica el hash, sustituye el archivo y recarga el plugin.

La eliminación se realiza con el botón **Eliminar**, con confirmación. El panel descarga el plugin del runtime y elimina el archivo `.wasm` y el registro del plugin. Tenga en cuenta que el almacenamiento clave-valor del plugin (la tabla `plugin_storage`) **no** se borra al eliminarlo: en una reinstalación el plugin verá su configuración anterior.

Después de la instalación, la actualización y la eliminación, la página se recarga para actualizar los frontends de los plugins (`/plugins.js`).

## Activación y desactivación

La versión actual del panel no tiene una acción separada de «activar/desactivar» para un plugin en la interfaz: un plugin funciona desde el momento en que se instala hasta que se elimina.

Si un plugin deja de responder (supera el tiempo de espera de la llamada), el panel lo desactiva temporalmente: deja de recibir eventos y solicitudes HTTP hasta que se reinicie el panel. El estado se muestra como una insignia en la lista de plugins.

Para desactivar todos los plugins por completo, defina la variable de entorno `PLUGINS_DISABLED=true` y reinicie el panel.

## Permisos

Un plugin puede registrar sus propios permisos de servidor de juego con el formato `plugin:{id}:{ability}`, por ejemplo `plugin:ezvdsxmlu6fbk:manage`. Estos permisos aparecen en la lista de permisos del servidor de juego junto con los integrados y se otorgan a los usuarios de la forma habitual (configurando los permisos de un usuario para un servidor). Los administradores reciben todos los permisos de los plugins automáticamente.

Basándose en estos permisos, un plugin oculta o muestra elementos de la interfaz (por ejemplo, una pestaña en la página del servidor de juego), y también puede restringir el acceso a sus rutas HTTP con las opciones `requires_auth` y `admin_only`.

## Variables de entorno del panel

| Variable | Valor por defecto | Descripción |
|---|---|---|
| `PLUGINS_DISABLED` | `false` | Desactiva por completo el sistema de plugins: los plugins no se cargan, y los frontends y las rutas de los plugins no se registran |
| `PLUGINS_AUTOLOAD` | — | Una lista de nombres de archivos `.wasm` separados por comas para registrar y cargar automáticamente al iniciar el panel |
| `PLUGIN_STORE_URL` | `https://plugins.gameap.dev/api` | Dirección de la API del catálogo de plugins |
| `PLUGIN_STORE_LICENSE_KEY` | — | Clave de licencia para descargar plugins de pago del catálogo |
| `PLUGIN_HTTP_ALLOWED_SCHEMES` | `https` | Esquemas permitidos para las solicitudes HTTP salientes de los plugins |
| `PLUGIN_HTTP_BLOCK_PRIVATE_IPS` | `true` | Bloquea las solicitudes de los plugins a IPs privadas, de loopback y de servicio (protección contra SSRF) |
| `PLUGIN_HTTP_ALLOWED_HOSTS` | — | Una lista de hosts de excepción que pueden solicitarse incluso con IPs privadas |
| `PLUGIN_HTTP_MAX_REDIRECTS` | `5` | Número máximo de redirecciones en las solicitudes salientes de los plugins |
| `PLUGIN_HTTP_MAX_TIMEOUT_SECONDS` | `30` | Tiempo de espera máximo de una solicitud HTTP saliente de un plugin, en segundos |

## Instalación manual a través del sistema de archivos

Un plugin puede instalarse sin la interfaz: copie el archivo `.wasm` en el directorio `plugins/` del panel, añada el nombre del archivo a la variable de entorno `PLUGINS_AUTOLOAD` (nombres separados por comas) y reinicie el panel: el plugin se registrará y cargará en el inicio.
