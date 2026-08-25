---
title: Usuarios, roles y permisos
layout: default
lang: es
category: Administración
order: 332
---

El acceso en el panel se basa en tres conceptos:

* **Permiso** — una acción permitida concreta, por ejemplo, "iniciar un servidor de juego".
* **Rol** — un conjunto de permisos con nombre. Un rol se asigna a un usuario.
* **Permiso directo** — un permiso otorgado a un usuario sin pasar por los roles.

Los permisos pueden otorgarse de forma global o vincularse a un objeto concreto — con mayor
frecuencia a un único servidor de juego. Así es como un usuario obtiene acceso a su propio
servidor y a nada más.

## Roles predeterminados

Durante la inicialización de la base de datos se crean dos roles:

| Rol     | Nombre en la interfaz | Permisos                        |
|---------|-----------------------|---------------------------------|
| `admin` | Administrador         | `admin roles & permissions`     |
| `user`  | Usuario               | Ninguno                         |

El primer usuario creado recibe el rol `admin`. El rol `user` no contiene ningún permiso: a un
usuario normal se le da acceso no mediante el rol, sino mediante permisos sobre servidores de
juego concretos.

> El permiso **`admin roles & permissions`** otorgado de forma global es lo que convierte a
> alguien en administrador. Controla el acceso a la sección "Administración" y también determina
> a quién se aplica la autenticación de dos factores obligatoria — consulte
> [Seguridad](/es/security.html).
>
> En el panel no existe un rol especial de "superadministrador", y el usuario con `id = 1` no
> tiene privilegios por defecto: solo importa disponer de este permiso.

## Permisos de servidor de juego

| Permiso                    | Nombre en la interfaz                    |
|----------------------------|-------------------------------------------|
| `game-server-common`       | Capacidad común del servidor de juego     |
| `game-server-start`        | Iniciar servidor de juego                 |
| `game-server-stop`         | Detener servidor de juego                 |
| `game-server-restart`      | Reiniciar servidor de juego               |
| `game-server-pause`        | Pausar servidor de juego                  |
| `game-server-update`       | Actualizar servidor de juego              |
| `game-server-files`        | Acceso al gestor de archivos              |
| `game-server-tasks`        | Acceso al planificador de tareas          |
| `game-server-settings`     | Acceso a la configuración                 |
| `game-server-console-view` | Acceso de lectura a la consola del servidor |
| `game-server-console-send` | Acceso para enviar comandos a la consola  |
| `game-server-rcon-console` | Consola RCON                              |
| `game-server-rcon-players` | Gestión de jugadores RCON                 |
| `game-server-metrics`      | Acceso a las métricas del servidor        |

`game-server-common` proporciona acceso básico al servidor: verlo en la lista y abrir su
página. Sin él, los demás permisos son prácticamente inútiles — empiece por este.

Los permisos de lectura y escritura de la consola están separados: `game-server-console-view`
solo permite verla, `game-server-console-send` — enviar comandos. Lo mismo ocurre con RCON: el
acceso a la consola RCON y la gestión de jugadores se otorgan por separado.

Además de los permisos de servidor de juego, existen los permisos genéricos `create`, `view`,
`edit` y `delete` — se aplican a las secciones del panel.

## Otorgar permisos sobre un servidor concreto

**Administración** → **Usuarios** → seleccione un usuario → **Editar**.

El formulario tiene dos bloques:

* **Roles** — `Administrador` o `Usuario`.
* **Servidores de juego** — una lista de servidores, en la que se pueden marcar los permisos
  necesarios para cada uno.

Los permisos otorgados aquí están vinculados al servidor concreto: el usuario obtiene acceso a
él y a nada más.

Lo mismo se hace a través de la API:

```http
GET  /api/users/{id}/servers                              — los servidores del usuario
GET  /api/users/{id}/servers/{server}/permissions         — permisos sobre el servidor
PUT  /api/users/{id}/servers/{server}/permissions         — cambiar permisos
```

Sus propios permisos sobre un servidor pueden consultarse con
`GET /api/servers/{server}/abilities`.

## Cómo se comprueba el acceso

Los permisos de un usuario se recopilan de dos fuentes: los otorgados directamente y los que
provienen de todos sus roles. Después se aplican dos reglas.

**Una denegación es más fuerte que una concesión.** Si existe un registro de denegación para el
permiso en cuestión, el acceso no se concede, sin importar cuántas concesiones se hayan emitido
por otras vías ni en qué orden estén registradas.

**El ámbito debe coincidir.** Un permiso otorgado para un servidor no tiene efecto sobre otro.
Una denegación establecida en un servidor tampoco afecta a los demás.

Un permiso puede otorgarse en tres ámbitos:

| Ámbito                      | Significado                                                       |
|-----------------------------|-------------------------------------------------------------------|
| Global                      | Se aplica en todas partes                                         |
| Sobre un tipo de objeto     | Se aplica a todos los servidores de juego, todos los servidores dedicados, etc. |
| Sobre un objeto concreto    | Se aplica solo a este servidor                                    |

Las comprobaciones que requieren acceso global — la comprobación de administrador, por ejemplo —
tienen en cuenta **únicamente un permiso otorgado de forma global**. Una concesión vinculada a
un único servidor no convierte a nadie en administrador.

> Revocar un permiso en la interfaz no siempre significa simplemente eliminar el registro. Si
> tras la eliminación el usuario siguiera recibiendo el permiso a través de un rol, el panel
> establece además una denegación explícita — de lo contrario, el permiso volvería a través del
> rol.

### Caché

Los resultados de las comprobaciones se almacenan en caché durante el tiempo indicado en
`RBAC_CACHE_TTL` (30 segundos por defecto). Los cambios de permisos realizados a través de la
interfaz y la API vacían la caché inmediatamente. Las ediciones hechas directamente en la base
de datos no vacían la caché — esos cambios surten efecto dentro del tiempo configurado.

## Gestión de usuarios

| Método y ruta            | Propósito                       |
|--------------------------|---------------------------------|
| `GET /api/users`         | Lista de usuarios               |
| `POST /api/users`        | Crear un usuario                |
| `GET /api/users/{id}`    | Detalles del usuario            |
| `PUT /api/users/{id}`    | Actualizar un usuario           |
| `DELETE /api/users/{id}` | Eliminar un usuario             |

La lista de usuarios está disponible solo para los administradores.

La contraseña establecida al crear o actualizar un usuario pasa por la comprobación de la
política de contraseñas: al menos 12 bytes y no debe figurar en la lista de contraseñas
comprometidas. Consulte [Seguridad](/es/security.html) para más detalles.

En el panel no hay autoregistro de usuarios: las cuentas las crea un administrador.

## Tareas habituales

### Dar a un usuario acceso a un único servidor

1. Cree un usuario: **Administración** → **Usuarios** → **Crear**.
2. Asígnele el rol `Usuario`.
3. En el bloque "Servidores de juego", seleccione el servidor y marque los permisos. El
   conjunto mínimo funcional: `game-server-common`, `game-server-start`, `game-server-stop`,
   `game-server-restart`.
4. Añada según sea necesario: el gestor de archivos, la consola, RCON, las tareas.

### Convertir a un usuario en administrador

Asigne el rol `Administrador`. Contiene el permiso `admin roles & permissions`, que otorga
acceso completo.

> A partir de ese momento, el requisito de autenticación de dos factores se aplica al usuario:
> en su próximo inicio de sesión verá un recordatorio y, tras 30 días, una exigencia de
> activarla. Consulte [Seguridad](/es/security.html).

### Revocar el acceso sin eliminar la cuenta

Elimine los permisos sobre los servidores y deje el rol `Usuario`. El usuario podrá iniciar
sesión, pero no verá ningún servidor. Las sesiones activas no se terminan con esto — siguen
funcionando hasta que caducan (24 horas, o 7 días para los inicios de sesión con "recordarme"
activado).
