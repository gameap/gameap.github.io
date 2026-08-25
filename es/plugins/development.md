---
title: Desarrollo de plugins
layout: default
lang: es
category: Plugins
order: 342
---

## Arquitectura

Un plugin de GameAP es un módulo WASM para el target `wasm32-wasip1` (WASI preview 1), compilado como reactor. El panel ejecuta cada plugin en un entorno aislado de [wazero](https://github.com/tetratelabs/wazero): sin sistema de archivos, sin red, sin variables de entorno; stdout y stderr se descartan.

El panel y el plugin intercambian mensajes protobuf a través de la memoria lineal de WASM: el panel escribe la solicitud en la memoria del invitado usando el `malloc` exportado por el plugin, llama a la función y lee la respuesta. La versión de la API de plugins es `1`.

Los tipos de mensajes y los SDK se publican desde el repositorio [github.com/gameap/gameap-proto](https://github.com/gameap/gameap-proto):

* `gameap-plugin-sdk` (crates.io) — el SDK completo de Rust: la capa ABI, el trait `Plugin`, la macro `register_plugin!`, clientes tipados de funciones del host;
* `@gameap/proto-as` (npm) — tipos de mensajes de AssemblyScript (sin la capa ABI);
* `@gameap/proto` (npm) — tipos de TypeScript (para herramientas; actualmente no es posible compilar un plugin WASM en JS puro).

## La interfaz del plugin

Un plugin implementa el servicio `PluginService` de `plugin.proto`. Los métodos del servicio:

| Método | Propósito |
|---|---|
| `GetInfo` | Devuelve los metadatos del plugin (`PluginInfo`) |
| `Initialize` | Inicialización al cargar el plugin |
| `Shutdown` | Apagado al descargar el plugin |
| `HandleEvent` | Manejo de un evento del panel |
| `GetSubscribedEvents` | La lista de tipos de eventos a los que el plugin se suscribe |
| `GetHTTPRoutes` | La lista de rutas HTTP del plugin |
| `HandleHTTPRequest` | Manejo de una solicitud HTTP en una ruta del plugin |
| `GetFrontendBundle` | El frontend integrado: bundle JS y CSS (opcional) |
| `GetServerAbilities` | Registro de los permisos de servidor del plugin (opcional) |

En el SDK de Rust la interfaz está representada por el trait `Plugin` con implementaciones por defecto neutras — en la práctica solo `get_info` es obligatorio. La macro `register_plugin!` genera todas las exportaciones WASM necesarias, incluidas `malloc`/`free` y la comprobación de la versión de la API.

## Metadatos de `PluginInfo`

| Campo | Descripción |
|---|---|
| `id` | El identificador de cadena del plugin (véanse los requisitos a continuación) |
| `name` | Nombre del plugin |
| `version` | Versión (semver) |
| `description` | Descripción breve |
| `author` | Autor |
| `license` | Licencia (por ejemplo, `MIT`) |
| `homepage` | Enlace a la página del plugin |
| `required_permissions` | Permisos que declara el plugin: se registran al instalar y se comprueban en cada llamada a una función del host (véase más abajo) |
| `api_version` | Versión de la API de plugins, debe ser `"1"` |

**Requisitos para `id`.** Utilice un identificador estable compuesto por caracteres del alfabeto base32 `a-z2-7`, sin guiones. El panel normaliza el id: una cadena con guiones u otros caracteres se reemplaza por un hash, lo que rompe las rutas `/api/plugins/{id}/...` y `/plugins/{id}/...`. Evite también los ids puramente numéricos: dicho identificador se trata como un ID numérico decimal (el análisis del id primero intenta interpretar la cadena como un número). Ejemplos de ids válidos de plugins reales: `hexeditor4jm2`, `ezvdsxmlu6fbk`, `dshdabjp2l73a`.

## Eventos

Un plugin se suscribe a los eventos del panel con el método `GetSubscribedEvents`. Los tipos de eventos (listados sin el prefijo `EVENT_TYPE_`):

| Evento | Cancelable | Entrega |
|---|---|---|
| `SERVER_PRE_START`, `SERVER_PRE_STOP`, `SERVER_PRE_RESTART`, `SERVER_PRE_INSTALL`, `SERVER_PRE_UPDATE`, `SERVER_PRE_REINSTALL`, `SERVER_PRE_DELETE` | Sí | Síncrona |
| `SERVER_POST_START`, `SERVER_POST_STOP`, `SERVER_POST_RESTART`, `SERVER_POST_INSTALL`, `SERVER_POST_UPDATE`, `SERVER_POST_REINSTALL`, `SERVER_POST_DELETE` | No | Asíncrona |
| `SERVER_CREATED`, `SERVER_UPDATED`, `SERVER_DELETED` | No | Asíncrona |
| `DAEMON_TASK_CREATED`, `DAEMON_TASK_COMPLETED`, `DAEMON_TASK_FAILED` | No | Asíncrona |

Los pre-eventos se entregan de forma síncrona y bloquean la operación: un plugin puede cancelarla devolviendo un `EventResult` con `should_cancel = true` y un `message`, o modificar los datos mediante `modified_data`. Los post-eventos y los demás tipos se entregan de forma asíncrona y no afectan a la operación.

Un evento de servidor contiene una instantánea completa del servidor (`ServerEventPayload`); un evento de tarea contiene los datos de la tarea del daemon (`TaskEventPayload`). El tiempo de espera del manejador de eventos es de 10 segundos; una vez expirado, el plugin se desactiva hasta que se reinicie el panel.

## Rutas HTTP del plugin

Un plugin registra rutas HTTP con el método `GetHTTPRoutes`; cada ruta es un mensaje `HTTPRoute`:

| Campo | Descripción |
|---|---|
| `path` | La ruta, comienza con `/`; se admiten parámetros de ruta de la forma `{name}` |
| `methods` | Métodos: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, `OPTIONS` |
| `requires_auth` | Requerir un usuario autenticado |
| `admin_only` | Requerir un administrador |
| `description` | Descripción de la ruta |

Las rutas están disponibles en `/api/plugins/{plugin_id}/...`. La solicitud pasada al plugin contiene el método, la ruta relativa, las cabeceras, los parámetros de ruta y de consulta, el cuerpo (no mayor de 1 MB) y la sesión del usuario (si la solicitud está autenticada). El plugin devuelve un `HTTPResponse` con los campos `status_code`, `headers` y `body`.

## Funciones del host

Todas las llamadas desde un plugin al panel y al mundo exterior pasan por funciones del host, agrupadas en módulos `gameap-*`:

| Módulo | Funciones | Propósito |
|---|---|---|
| `gameap-log` | `log` | Escritura en el registro del panel |
| `gameap-cache` | `get`, `set`, `delete` | La caché compartida del panel (claves con prefijo `plugin:`) |
| `gameap-crypto` | `random_uint64`, `random_string`, `argon2_hash`, `argon2_verify` | Generadores de valores aleatorios, hash Argon2 |
| `gameap-http` | `fetch` | Solicitudes HTTP salientes con protección SSRF |
| `gameap-storage` | `get`, `set`, `delete`, `list` | El almacenamiento clave-valor persistente del plugin |
| `gameap-servercontrol` | `start_server`, `stop_server`, `restart_server`, `update_server`, `install_server`, `reinstall_server` | Control de servidores de juego (devuelve un `task_id`) |
| `gameap-servers` | `find_servers`, `get_server`, `save_server`, `delete_server` | Servidores de juego |
| `gameap-users` | `find_users`, `get_user` | Usuarios del panel |
| `gameap-nodes` | `find_nodes`, `get_node` | Servidores dedicados (nodos) |
| `gameap-games` | `find_games`, `get_game` | Juegos |
| `gameap-gamemods` | `find_game_mods`, `get_game_mod` | Mods de juegos |
| `gameap-daemontasks` | `find_daemon_tasks`, `create_daemon_task` | Tareas del daemon |
| `gameap-serversettings` | `find_server_settings`, `save_server_setting` | Configuración de servidores de juego |
| `gameap-nodefs` | `read_dir`, `mk_dir`, `copy`, `move`, `download`, `upload`, `remove`, `get_file_info`, `chmod`, `hash`, `create_archive`, `extract_archive`, `start_create_archive`, `start_extract_archive`, `cancel_archive`, `get_archive_operation` | Operaciones con archivos en un nodo |
| `gameap-nodecmd` | `execute_command` | Ejecución de un comando en un nodo |

Detalles:

* `gameap-http` actúa como proxy de las solicitudes a través del panel con protección SSRF: por defecto solo se permite el esquema `https`, las IPs privadas y de servicio están bloqueadas, y el cuerpo de la respuesta está limitado a 10 MB. La política se configura con las variables de entorno `PLUGIN_HTTP_*` (véase [Instalación y administración](/es/plugins/management.html)).
* `gameap-storage` es el almacenamiento persistente del plugin, aislado por `plugin_id`, con vinculación opcional de registros a una entidad (`entity_type`, `entity_id`). Utilícelo para la configuración del plugin: el mecanismo de configuración separado (`config` en `InitializeRequest`) no se usa en la versión actual del panel.
* `gameap-nodefs` y `gameap-nodecmd` trabajan con archivos y comandos en el servidor dedicado (nodo) a través de GameAP Daemon. Las lecturas (`read_dir`, `download`, `get_file_info`, `hash`, `get_archive_operation`) requieren el permiso `files_read`; todo lo que escribe — incluidos `chmod` y las funciones de archivado — requiere `files`, que incluye `files_read`.
* `gameap-nodefs.download` sin `offset`/`length` devuelve el archivo entero en un solo mensaje y está limitado por `PLUGIN_NODEFS_MAX_INLINE`; un archivo mayor se rechaza con un error que indica ambos tamaños. Indicar una ventana `offset`/`length` permite leer ese archivo por partes: solo la ventana tiene que caber en el límite, y un `length` mayor que el límite se rechaza en lugar de recortarse en silencio. La respuesta repite `offset` e incluye `total_size`, y `length: 0` junto con un `offset` lee tanto como permita el límite — por eso un bucle de paginación avanza según la longitud de `content` hasta alcanzar `total_size`. Leer al final del archivo o más allá no es un error: la respuesta trae `content` vacío.

**Comprobaciones en una llamada al host.** Las funciones del host se ejecutan con los propios privilegios del panel, pero la llamada no pasa sin comprobaciones:

* **Permisos.** Los módulos privilegiados están protegidos por los permisos registrados para el plugin: `files_read`, `files`, `manage_servers`, `node_commands`, `listen_events`, `manage_rbac`, `secrets`, `ssh`, `manage_nodes` (`manage_games`, `manage_game_mods` y `manage_users` están reservados para operaciones de escritura que todavía no existen). Un permiso más amplio satisface a uno más estrecho: `files` incluye `files_read`. La aplicación es transitoria: `PLUGIN_PERMISSIONS_ENFORCE` vale `false` por defecto, así que los permisos se registran, se muestran y se pueden editar mientras todas las comprobaciones siguen pasando; una versión futura lo pondrá en `true` por defecto, así que declare ya los permisos que su plugin necesita. Las llamadas de modificación de `gameap-nodes` comprueban su permiso independientemente de ese ajuste. Una llamada rechazada responde `plugin permission <name> required`.
* **Rutas.** Toda ruta entregada a `gameap-nodefs`, el `work_dir` de `gameap-nodecmd` y un archivo del nodo referenciado por una respuesta HTTP se comprueban en el panel antes de llegar al daemon. Una ruta con un segmento `..` o un byte NUL se rechaza siempre. `PLUGIN_NODEFS_PATH_POLICY` puede además confinar las rutas al `work_path` del nodo o a los directorios de los servidores de juego que haya en él; en todos los modos restringidos queda abierto el directorio de servicio del propio plugin, `<work_path>/.plugins/<id del plugin>`, y ahí es donde deben ir los archivos de trabajo en el nodo. Una llamada rechazada responde `path policy: <motivo>: <ruta>`.
* **Límites de frecuencia.** Cada plugin tiene su propio cubo de tokens por clase: `gameap-nodefs` 50 llamadas/s (ráfaga 200), `gameap-http` 20/s (50), `gameap-ssh` 20/s (60), `gameap-rbac` 10/s (50), `gameap-nodecmd` y el control de servidores 5/s (20). Una llamada rechazada responde `rate limited: ...`; el plugin nunca se desactiva por ello.
* **Registro de auditoría.** Las operaciones privilegiadas — control de servidores, comandos en el nodo, escritura de archivos, SSH, cambios de RBAC — se registran con el plugin como actor, incluidos los rechazos.

Aun así, la instalación de plugins está encomendada únicamente a los administradores — instale plugins solo de fuentes de confianza.

## Límites en tiempo de ejecución

| Límite | Valor |
|---|---|
| Llamadas concurrentes a un plugin | 1 (las llamadas se serializan) |
| Tiempo de espera de una llamada al plugin | 30 s (inicio del módulo — 60 s) |
| Tiempo de espera del manejador de eventos | 10 s |
| Tamaño de un archivo `.wasm` subido | 100 MB |
| Cuerpo de una solicitud HTTP a un plugin | 1 MB |
| Cuerpo de una respuesta de `gameap-http` | 10 MB |
| Un mensaje `download` / `upload` de `gameap-nodefs` | 32 MB (`PLUGIN_NODEFS_MAX_INLINE`); un `download` con ventana se mide por la ventana |

Cuando se supera el tiempo de espera de una llamada, el panel desactiva el plugin hasta un reinicio. El sistema de archivos y la red no están disponibles desde WASM — solo a través de las funciones del host.

## Compilación de un plugin con Rust

La estructura de un proyecto mínimo (siguiendo [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor)):

`rust-toolchain.toml`:
```toml
[toolchain]
channel = "1.94.0"
targets = ["wasm32-wasip1"]
```

`Cargo.toml`:
```toml
[lib]
crate-type = ["cdylib"]

[dependencies]
gameap-plugin-sdk = "0.1"

[profile.release]
opt-level = "z"
lto = true
strip = true
```

Un `src/lib.rs` mínimo — un plugin que solo devuelve metadatos y un frontend integrado:

```rust
#![cfg(target_arch = "wasm32")]

use gameap_plugin_sdk::proto::gameap::plugin as pb;
use gameap_plugin_sdk::{Plugin, PluginError, register_plugin};

const FRONTEND_JS: &[u8] = include_bytes!(concat!(env!("OUT_DIR"), "/plugin.js"));
const FRONTEND_CSS: &[u8] = include_bytes!(concat!(env!("OUT_DIR"), "/plugin.css"));

#[derive(Default)]
struct MyPlugin;

impl Plugin for MyPlugin {
    fn get_info(&mut self, _req: pb::GetInfoRequest) -> Result<pb::PluginInfo, PluginError> {
        Ok(pb::PluginInfo {
            id: "myplugin2j7d".into(),
            name: "My Plugin".into(),
            version: env!("CARGO_PKG_VERSION").into(),
            description: "My first GameAP plugin".into(),
            author: "Me".into(),
            api_version: "1".into(),
            ..Default::default()
        })
    }

    fn get_frontend_bundle(
        &mut self,
        _req: pb::GetFrontendBundleRequest,
    ) -> Result<pb::GetFrontendBundleResponse, PluginError> {
        Ok(pb::GetFrontendBundleResponse {
            bundle: FRONTEND_JS.to_vec(),
            has_bundle: !FRONTEND_JS.is_empty(),
            styles: FRONTEND_CSS.to_vec(),
            has_styles: !FRONTEND_CSS.is_empty(),
        })
    }
}

register_plugin!(MyPlugin);
```

Compilación:

```bash
cargo build --target wasm32-wasip1 --release
# opcional — reducción de tamaño (binaryen):
wasm-opt -Oz target/wasm32-wasip1/release/my_plugin.wasm -o my-plugin.wasm
```

El archivo `.wasm` resultante se instala a través de la interfaz del panel o se copia en el directorio `plugins/` (véase [Instalación y administración](/es/plugins/management.html)). El frontend se compila por separado y se integra en el `.wasm` (véase [Frontend del plugin](/es/plugins/frontend.html)).

## Desarrollo en AssemblyScript

Todavía no existe un SDK listo para AssemblyScript: el paquete `@gameap/proto-as` proporciona solo los tipos de mensajes, y la capa ABI (empaquetado de punteros, manejo de errores, exportación de `malloc`/`free`) hay que escribirla a mano. Debido al recolector de basura de AssemblyScript, los búferes pasados al host deben fijarse (`__pin`/`__unpin`), y un analizador JSON y utilidades básicas hay que implementarlos uno mismo o traerlos de bibliotecas de terceros. Un ejemplo funcional de plugin en AssemblyScript es [plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth).

## Validación de un plugin antes de la instalación

El panel puede validar un archivo `.wasm` sin instalarlo — el endpoint `POST /api/admin/plugins/upload/dry-run` (multipart, campo `file`). La respuesta devuelve los metadatos del plugin, las rutas HTTP, los permisos de servidor, las suscripciones a eventos, si hay un frontend presente y una lista de errores. La misma comprobación se ejecuta en la interfaz al subir un archivo.

## Ejemplos de plugins

* [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) — un plugin mínimo de Rust: metadatos más un frontend integrado, sin funciones del host.
* [plugin-goldsrc-addons](https://github.com/gameap/plugin-goldsrc-addons) — un plugin de Rust con rutas HTTP, operaciones con archivos en un nodo (`nodefs`), ejecución de comandos (`nodecmd`) y llamadas a la API del panel (RCON) desde el frontend.
* [plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth) — un plugin en AssemblyScript: llamadas a una API externa (modrinth.com) a través de `gameap-http`, caché, almacenamiento persistente.
