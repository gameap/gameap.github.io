---
title: Referencia de config.env
layout: default
lang: es
category: Administración
order: 334
---

El panel se configura mediante variables de entorno. Los valores se leen del archivo `config.env`:

* Linux — `/etc/gameap/config.env`
* Windows — `C:\gameap\web\config.env`

El formato del archivo es un par `NAME=value` por línea, sin comillas y sin `export`. Las líneas
que comienzan con `#` se ignoran. Después de modificar el archivo, hay que reiniciar el panel:

```bash
gameapctl panel restart
```

Las variables definidas en el entorno del proceso tienen prioridad sobre el archivo.

Solo dos variables son obligatorias: `DATABASE_URL` y `AUTH_SECRET`. El panel no arranca sin
ellas. Todo lo demás tiene valores predeterminados adecuados para una instalación típica de un solo servidor.

Notación utilizada en las tablas: *duration* — un número con sufijo `s`, `m` u `h` (`30s`, `5m`,
`12h`); *list* — valores separados por comas; *size* — un número con sufijo binario `K`, `M`, `G`, `T`
(`8M`, `100G`); un número sin sufijo se interpreta como bytes.

## HTTP

| Variable               | Tipo   | Predeterminado | Propósito                                                                       |
|------------------------|--------|----------------|----------------------------------------------------------------------------------|
| `HTTP_HOST`            | string | `0.0.0.0`      | Nombre de host del panel. Se usa para CORS y como nombre en los certificados autofirmados |
| `HTTP_BIND_IP`         | string | `""`           | Dirección IP de escucha. Vacío — todas las interfaces                            |
| `HTTP_PORT`            | number | `8025`         | Puerto de la interfaz web y de la API                                            |
| `HTTPS_PORT`           | number | `443`          | Puerto HTTPS. Solo escucha cuando hay un certificado configurado                 |
| `HTTP_ALLOWED_ORIGINS` | list   | `""`           | Orígenes autorizados a hacer solicitudes desde el navegador. Vacío — un único origen derivado de `HTTP_HOST` |

`HTTP_ALLOWED_ORIGINS` acepta orígenes completos incluyendo el esquema: `https://panel.example.com`.
El comodín `*` no está soportado.

## Base de datos

| Variable          | Tipo   | Predeterminado | Propósito                                 |
|-------------------|--------|----------------|--------------------------------------------|
| `DATABASE_DRIVER` | string | `mysql`        | `mysql`, `postgres`, `sqlite`, `inmemory`  |
| `DATABASE_URL`    | string | —              | **Obligatoria.** Cadena de conexión        |

Los nombres del driver de PostgreSQL son intercambiables: `postgres`, `postgresql`, `pgsql`, `pg`, `pgx`.

Formatos de la cadena de conexión:

```text
# PostgreSQL
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://gameap:password@localhost:5432/gameap?sslmode=disable

# MySQL / MariaDB
DATABASE_DRIVER=mysql
DATABASE_URL=gameap:password@tcp(localhost:3306)/gameap?parseTime=true

# SQLite
DATABASE_DRIVER=sqlite
DATABASE_URL=file:/var/lib/gameap/db.sqlite?_busy_timeout=5000&_journal_mode=WAL&cache=shared
```

El driver `inmemory` mantiene los datos solo en RAM y está pensado para pruebas: todo se pierde al
reiniciar.

## TLS y certificados

| Variable          | Tipo   | Predeterminado | Propósito                                          |
|-------------------|--------|----------------|-----------------------------------------------------|
| `TLS_CERT_FILE`   | string | `""`           | Ruta al archivo del certificado                     |
| `TLS_KEY_FILE`    | string | `""`           | Ruta al archivo de la clave privada                 |
| `TLS_CERT`        | string | `""`           | Certificado directamente en la variable, PEM o base64 |
| `TLS_KEY`         | string | `""`           | Clave privada directamente en la variable, PEM o base64 |
| `TLS_FORCE_HTTPS` | bool   | `false`        | Redirigir HTTP a HTTPS                              |

HTTPS se activa cuando al menos una fuente de certificados está configurada: archivos, valores en
variables o ACME. Si ninguna está configurada, el panel sirve solo HTTP.

Con `TLS_FORCE_HTTPS=true`, toda solicitud HTTP recibe una redirección `301`, excepto
`/.well-known/acme-challenge/` — de lo contrario sería imposible emitir un certificado de Let's Encrypt.

### ACME (Let's Encrypt)

| Variable                      | Tipo     | Predeterminado  | Propósito                                                      |
|-------------------------------|----------|-----------------|-----------------------------------------------------------------|
| `ACME_ENABLED`                | bool     | `false`         | Emisión automática de certificados                              |
| `ACME_EMAIL`                  | string   | `""`            | Dirección para notificaciones. **Obligatoria** cuando ACME está activado |
| `ACME_DOMAINS`                | list     | `""`            | Dominios del certificado. **Obligatoria** cuando ACME está activado |
| `ACME_CHALLENGE_TYPE`         | string   | `http-01`       | `http-01` o `dns-01`                                            |
| `ACME_DNS_PROVIDER`           | string   | `""`            | Proveedor DNS para `dns-01`, por ejemplo `cloudflare`           |
| `ACME_DIRECTORY_URL`          | string   | ACME de producción | URL del directorio ACME. Para pruebas, apúntela al staging de Let's Encrypt |
| `ACME_RENEWAL_THRESHOLD`      | duration | `720h`          | Con cuánta antelación al vencimiento renovar. 30 días por defecto |
| `ACME_RENEWAL_CHECK_INTERVAL` | duration | `12h`           | Con qué frecuencia comprobar la fecha de vencimiento            |
| `ACME_PROPAGATION_TIMEOUT`    | duration | `180s`          | Cuánto tiempo esperar la propagación del registro DNS con `dns-01` |
| `ACME_STORAGE_PATH`           | string   | `acme`          | Directorio para almacenar los certificados y la clave de la cuenta ACME |

ACME solo se activa cuando `ACME_ENABLED=true`, `ACME_EMAIL` y `ACME_DOMAINS` están definidos.
Si falta algo, el panel continúa silenciosamente sin ACME.

El valor predeterminado es el directorio de producción de Let's Encrypt, con sus estrictos límites
sobre el número de intentos. Mientras depura su configuración, use staging:
`ACME_DIRECTORY_URL=https://acme-staging-v02.api.letsencrypt.org/directory`.

## Autenticación

| Variable                        | Tipo     | Predeterminado | Propósito                                                                |
|---------------------------------|----------|----------------|---------------------------------------------------------------------------|
| `AUTH_SECRET`                   | string   | —              | **Obligatoria.** Clave de firma de tokens, exactamente 32 bytes aleatorios |
| `ENCRYPTION_KEY`                | string   | `""`           | Clave de cifrado de los secretos en la base de datos, exactamente 32 bytes aleatorios |
| `AUTH_SERVICE`                  | string   | `paseto`       | Formato de token: `paseto` o `jwt`. Cualquier otro valor — el panel no arrancará |
| `AUTH_BCRYPT_COST`              | number   | `13`           | Coste de bcrypt, de 10 a 14                                               |
| `AUTH_ALLOW_WEAK_PASSWORDS`     | bool     | `false`        | Desactiva la comprobación contra la lista de contraseñas comprometidas    |
| `AUTH_REQUIRE_MFA_FOR_ADMINS`   | bool     | `true`         | Exigir 2FA a los administradores                                          |
| `AUTH_MFA_HARD_FAIL_DAYS`       | number   | `30`           | Días hasta el bloqueo. `0` — solo recordatorio                            |
| `AUTH_MFA_ENROLLMENT_TOKEN_TTL` | duration | `15m`          | Duración de la sesión restringida emitida tras la fecha límite            |
| `AUTH_SHORT_LIVED_TOKEN_TTL`    | duration | `10s`          | Duración de los tokens de un solo uso `glst_`. Limitada efectivamente a 10 segundos |

Consulte la página [Seguridad](/es/security.html) para más detalles.

> `AUTH_SECRET` se ajusta silenciosamente a 32 bytes: un valor corto se rellena, uno largo se trunca.
> Defina exactamente 32 bytes aleatorios, por ejemplo `openssl rand -hex 16`.

### Primer administrador

Estas tres variables se leen solo durante la inicialización de una base de datos vacía:

| Variable         | Propósito                                                                      |
|------------------|---------------------------------------------------------------------------------|
| `ADMIN_LOGIN`    | Nombre de usuario del primer administrador                                      |
| `ADMIN_EMAIL`    | Dirección de correo electrónico                                                 |
| `ADMIN_PASSWORD` | Contraseña. Si no se define, se genera una aleatoria y se muestra en el log en el primer arranque |

> La contraseña de `ADMIN_PASSWORD` **no** se comprueba contra la política de contraseñas: ni la
> longitud ni la lista de contraseñas comprometidas. Elija con cuidado.

## Control de acceso y caché

| Variable                    | Tipo     | Predeterminado   | Propósito                                     |
|-----------------------------|----------|------------------|------------------------------------------------|
| `RBAC_CACHE_TTL`            | duration | `30s`            | Duración de la caché de comprobación de permisos |
| `CACHE_DRIVER`              | string   | `memory`         | `memory`, `redis`, `postgres`, `mysql`. Véase más abajo |
| `CACHE_REDIS_ADDR`          | string   | `localhost:6379` | Dirección de Redis                             |
| `CACHE_REDIS_PASSWORD`      | string   | `""`             | Contraseña de Redis                            |
| `CACHE_REDIS_DB`            | number   | `0`              | Número de la base de datos de Redis            |
| `CACHE_TTL_RBAC`            | duration | `24h`            | Duración de la caché de permisos               |
| `CACHE_TTL_GAMES`           | duration | `48h`            | Duración de la caché de juegos                 |
| `CACHE_TTL_NODES`           | duration | `24h`            | Duración de la caché de servidores dedicados   |
| `CACHE_TTL_USERS`           | duration | `6h`             | Duración de la caché de usuarios               |
| `CACHE_TTL_PERSONAL_TOKENS` | duration | `24h`            | Duración de la caché de tokens personales      |
| `CACHE_TTL_SERVER_SETTINGS` | duration | `12h`            | Duración de la caché de la configuración de servidores |

`CACHE_DRIVER` acepta `memory` (alias `inmemory`), `redis`, `mysql` (alias `database`) y
`postgres` (alias `postgresql`, `pgsql`, `pg`). Un valor desconocido hace que el panel falle al
arrancar.

La caché contiene más que datos de referencia: ahí residen la clave de configuración del daemon,
la lista de tokens revocados y los contadores de intentos de inicio de sesión. Con
`CACHE_DRIVER=memory` todo eso se pierde al reiniciar y no se comparte entre varias instancias del
panel. Para una instalación con múltiples instancias, use `redis`.

## Archivos

| Variable                     | Tipo   | Predeterminado | Propósito                              |
|------------------------------|--------|----------------|-----------------------------------------|
| `FILES_DRIVER`               | string | `local`        | `local` o `s3`                          |
| `FILES_LOCAL_BASE_PATH`      | string | `""`           | Directorio base para el driver `local`  |
| `FILES_S3_ENDPOINT`          | string | `""`           | Dirección del almacenamiento compatible con S3 |
| `FILES_S3_USE_SSL`           | bool   | `true`         | Acceder al almacenamiento por HTTPS     |
| `FILES_S3_ACCESS_KEY_ID`     | string | `""`           | ID de la clave de acceso                |
| `FILES_S3_SECRET_ACCESS_KEY` | string | `""`           | Clave de acceso secreta                 |
| `FILES_S3_BUCKET`            | string | `""`           | Nombre del bucket                       |

### Subida de archivos

| Variable                        | Tipo     | Predeterminado | Propósito                                           |
|---------------------------------|----------|----------------|------------------------------------------------------|
| `FILES_UPLOAD_CHUNK_SIZE`       | size     | `8M`           | Tamaño de los fragmentos en las subidas por partes   |
| `FILES_UPLOAD_MAX_CHUNKS`       | number   | `100000`       | Número máximo de fragmentos por archivo              |
| `FILES_UPLOAD_SESSION_TTL`      | duration | `24h`          | Cuánto tiempo vive una subida sin terminar           |
| `FILES_UPLOAD_DISPATCH_TIMEOUT` | duration | `2m`           | Tiempo de espera para entregar el archivo ensamblado al daemon |
| `FILES_UPLOAD_JANITOR_INTERVAL` | duration | `12h`          | Con qué frecuencia limpiar las subidas expiradas     |
| `FILES_UPLOAD_ALLOWED_MIMES`    | list     | `""`           | Amplía la lista de tipos permitidos, no la reemplaza |
| `FILES_UPLOAD_ALLOW_ARCHIVES`   | bool     | `false`        | Permitir archivos comprimidos: zip, tar, gzip, bzip2, 7z, xz |
| `FILES_UPLOAD_ALLOW_BINARY`     | bool     | `false`        | Permitir archivos binarios arbitrarios               |

El tamaño máximo de archivo es `FILES_UPLOAD_CHUNK_SIZE` multiplicado por `FILES_UPLOAD_MAX_CHUNKS`.
Con los valores predeterminados eso es unos 780 GB. Una subida normal de una sola solicitud tiene un
límite aparte de 100 MB; no es configurable.

### Archivos comprimidos

| Variable                               | Tipo   | Predeterminado | Propósito                                      |
|----------------------------------------|--------|----------------|-------------------------------------------------|
| `FILES_ARCHIVE_MAX_BYTES`              | size   | `100G`         | Tamaño máximo de un archivo comprimido creado   |
| `FILES_ARCHIVE_MAX_FILES`              | number | `500000`       | Número máximo de archivos en un comprimido      |
| `FILES_ARCHIVE_CONCURRENT_PER_SERVER`  | number | `2`            | Operaciones de compresión concurrentes por servidor |

## gRPC

| Variable                      | Tipo   | Predeterminado | Propósito                                          |
|-------------------------------|--------|----------------|-----------------------------------------------------|
| `GRPC_PORT`                   | number | `31718`        | Puerto del servidor gRPC para los daemons           |
| `GRPC_TLS_ENABLED`            | bool   | `true`         | Cifrado de las conexiones de los daemons            |
| `GRPC_REQUIRE_MTLS`           | bool   | `false`        | Exigir un certificado de cliente al daemon          |
| `GRPC_EXTERNAL_HOST`          | string | `""`           | Dirección del panel comunicada al daemon            |
| `GRPC_EXTERNAL_PORT`          | number | `0`            | Puerto comunicado al daemon. `0` — se usa `GRPC_PORT` |
| `GRPC_MAX_RECV_MSG_SIZE`      | number | `10485760`     | Tamaño máximo del mensaje entrante, en bytes        |
| `GRPC_MAX_SEND_MSG_SIZE`      | number | `10485760`     | Tamaño máximo del mensaje saliente, en bytes        |
| `GRPC_MAX_CONCURRENT_STREAMS` | number | `100`          | Flujos concurrentes por conexión                    |
| `GRPC_ENABLE_REFLECTION`      | bool   | `false`        | Reflexión del esquema para herramientas de depuración |
| `DAEMON_SETUP_KEY`            | string | `""`           | Clave permanente de configuración del daemon en lugar de una temporal |

Consulte la página [GRPC API](/es/daemon/grpc.html) para más detalles.

> `DAEMON_SETUP_KEY` define una clave que nunca expira. Es cómodo para despliegues automatizados,
> pero cualquiera que conozca esa clave puede registrar un nuevo servidor dedicado en el panel.
> Para una instalación normal, deje esta variable sin definir — el panel emitirá una clave temporal
> válida durante una hora.

## Seguridad

La descripción completa está en la página [Seguridad](/es/security.html).

| Variable                           | Tipo   | Predeterminado                  | Propósito                                 |
|------------------------------------|--------|---------------------------------|--------------------------------------------|
| `SECURITY_HEADERS_ENABLED`         | bool   | `true`                          | Interruptor maestro de las cabeceras de seguridad |
| `SECURITY_CONTENT_TYPE_OPTIONS`    | bool   | `true`                          | `X-Content-Type-Options: nosniff`          |
| `SECURITY_FRAME_OPTIONS`           | string | `SAMEORIGIN`                    | `X-Frame-Options`                          |
| `SECURITY_REFERRER_POLICY`         | string | `strict-origin-when-cross-origin` | `Referrer-Policy`                        |
| `SECURITY_HSTS_ENABLED`            | bool   | `true`                          | HSTS. Solo se envía por HTTPS              |
| `SECURITY_HSTS_MAX_AGE`            | number | `31536000`                      | Duración de HSTS en segundos               |
| `SECURITY_HSTS_INCLUDE_SUBDOMAINS` | bool   | `false`                         | Extender HSTS a los subdominios            |
| `SECURITY_HSTS_PRELOAD`            | bool   | `false`                         | Añadir `preload`                           |
| `SECURITY_CSP_ENABLED`             | bool   | `true`                          | Content Security Policy                    |
| `SECURITY_CSP_REPORT_ONLY`         | bool   | `false`                         | Solo informes, sin bloqueo                 |
| `SECURITY_CSP_POLICY`              | string | `""`                            | Reemplaza por completo la política generada |
| `SECURITY_CSP_REPORT_URI`          | string | `""`                            | Dirección para los informes de CSP         |
| `SECURITY_CSP_EXTRA_SCRIPT_SRC`    | list   | `""`                            | Amplía `script-src`                        |
| `SECURITY_CSP_EXTRA_STYLE_SRC`     | list   | `""`                            | Amplía `style-src`                         |
| `SECURITY_CSP_EXTRA_CONNECT_SRC`   | list   | `""`                            | Amplía `connect-src`                       |
| `SECURITY_CSP_EXTRA_IMG_SRC`       | list   | `""`                            | Amplía `img-src`                           |
| `SECURITY_CSP_EXTRA_FRAME_SRC`     | list   | `""`                            | Amplía `frame-src`                         |
| `SECURITY_CSP_EXTRA_FONT_SRC`      | list   | `""`                            | Amplía `font-src`                          |
| `SECURITY_SENSITIVE_PATH_PREFIXES` | list   | véase más abajo                 | Rutas cuyas respuestas nunca deben almacenarse en caché |

El valor predeterminado de `SECURITY_SENSITIVE_PATH_PREFIXES`:
`/api/auth/,/api/profile/,/api/users/,/api/tokens/`.

### CAPTCHA

| Variable             | Tipo   | Predeterminado | Propósito                                                      |
|----------------------|--------|----------------|-----------------------------------------------------------------|
| `CAPTCHA_PROVIDER`   | string | `""`           | `recaptcha_v2`, `recaptcha_v3` o `turnstile`. Vacío — desactivado |
| `CAPTCHA_SITE_KEY`   | string | `""`           | Clave pública                                                   |
| `CAPTCHA_SECRET_KEY` | string | `""`           | Clave secreta                                                   |
| `CAPTCHA_MIN_SCORE`  | float  | `0.5`          | Umbral, solo reCAPTCHA v3                                       |
| `CAPTCHA_FAIL_OPEN`  | bool   | `false`        | Dejar pasar los inicios de sesión cuando el servicio de verificación no está disponible |
| `CAPTCHA_VERIFY_URL` | string | `""`           | Dirección de verificación personalizada                         |

### Registro de auditoría

| Variable                 | Tipo   | Predeterminado | Propósito                             |
|--------------------------|--------|----------------|----------------------------------------|
| `AUDIT_ENABLED`          | bool   | `true`         | Registro de eventos de seguridad       |
| `AUDIT_CLIENT_IP_HEADER` | string | `""`           | Cabecera que transporta la IP real del cliente |

## Plugins

| Variable                   | Tipo   | Predeterminado                 | Propósito                                |
|----------------------------|--------|--------------------------------|-------------------------------------------|
| `PLUGINS_DISABLED`         | bool   | `false`                        | Desactivar los plugins por completo       |
| `PLUGINS_AUTOLOAD`         | list   | `""`                           | Plugins cargados en el arranque           |
| `PLUGINS_CACHE_ENABLED`    | bool   | `true`                         | Caché de los módulos WebAssembly compilados |
| `PLUGINS_CACHE_DIR`        | string | `""`                           | Directorio de esa caché                   |
| `PLUGIN_STORE_URL`         | string | `https://plugins.gameap.dev/api` | Dirección del catálogo de plugins       |
| `PLUGIN_STORE_LICENSE_KEY` | string | `""`                           | Clave de licencia para los plugins de pago |

### Solicitudes de red de los plugins

| Variable                                | Tipo   | Predeterminado | Propósito                                        |
|-----------------------------------------|--------|----------------|---------------------------------------------------|
| `PLUGIN_HTTP_BLOCK_PRIVATE_IPS`         | bool   | `true`         | Bloquear solicitudes a direcciones de redes privadas |
| `PLUGIN_HTTP_ALLOWED_SCHEMES`           | list   | `https`        | Esquemas permitidos                               |
| `PLUGIN_HTTP_ALLOWED_HOSTS`             | list   | `""`           | Hosts exentos del bloqueo de direcciones privadas |
| `PLUGIN_HTTP_MAX_TIMEOUT_SECONDS`       | number | `30`           | Tiempo máximo de una solicitud                    |
| `PLUGIN_HTTP_MAX_REDIRECTS`             | number | `5`            | Número máximo de redirecciones                    |
| `PLUGIN_HTTP_RESPONSE_HEADER_ALLOWLIST` | list   | `""`           | Amplía la lista de cabeceras que se pasan al plugin |

Las direcciones del servicio de metadatos de los proveedores de nube están siempre bloqueadas;
`PLUGIN_HTTP_ALLOWED_HOSTS` no tiene efecto sobre ellas.

Un grupo aparte de ajustes gobierna los protocolos RCON y Query personalizados implementados por los plugins:

| Variable                         | Tipo   | Predeterminado | Propósito                                               |
|----------------------------------|--------|----------------|----------------------------------------------------------|
| `PLUGIN_NET_ENABLED`             | bool   | `true`         | Permitir que los plugins se comuniquen con los servidores de juego por la red |
| `PLUGIN_NET_BLOCK_PRIVATE_IPS`   | bool   | `false`        | Bloquear direcciones privadas. Desactivado: los servidores de juego suelen estar en la red interna |
| `PLUGIN_NET_ALLOWED_HOSTS`       | list   | `""`           | Hosts exentos                                            |
| `PLUGIN_NET_MAX_TIMEOUT_SECONDS` | number | `10`           | Tiempo máximo de una operación individual                |
| `PLUGIN_NET_READ_BUFFER_BYTES`   | number | `65536`        | Tamaño máximo de una lectura individual                  |
| `PLUGIN_NET_MAX_CONNECTIONS`     | number | `8`            | Conexiones concurrentes por plugin                       |

### Planificador de plugins

| Variable                                | Tipo     | Predeterminado | Propósito                                     |
|-----------------------------------------|----------|----------------|------------------------------------------------|
| `PLUGIN_SCHEDULER_MIN_INTERVAL`         | duration | `1s`           | Intervalo de tarea mínimo permitido            |
| `PLUGIN_SCHEDULER_MAX_TASKS_PER_PLUGIN` | number   | `32`           | Número máximo de tareas por plugin             |
| `PLUGIN_SCHEDULER_CALL_TIMEOUT`         | duration | `60s`          | Tiempo de espera predeterminado de la llamada al manejador |
| `PLUGIN_SCHEDULER_MAX_CALL_TIMEOUT`     | duration | `5m`           | Tiempo de espera máximo que un plugin puede solicitar |
| `PLUGIN_SCHEDULER_MAX_RETRIES`          | number   | `10`           | Número máximo de reintentos                    |
| `PLUGIN_SCHEDULER_MAX_RETRY_DELAY`      | duration | `10m`          | Retraso máximo entre reintentos                |
| `PLUGIN_SCHEDULER_MAX_JITTER`           | duration | `30s`          | Variación aleatoria máxima del inicio          |
| `PLUGIN_SCHEDULER_REFRESH_INTERVAL`     | duration | `30s`          | Con qué frecuencia releer las tareas de la base de datos |

## Intercambio de eventos entre instancias

Necesario solo cuando se ejecutan varias instancias del panel. Con una sola instancia los valores
predeterminados son suficientes. Un valor desconocido de `PUBSUB_DRIVER` hace que el panel falle al arrancar.

| Variable                     | Tipo     | Predeterminado | Propósito                                          |
|------------------------------|----------|----------------|-----------------------------------------------------|
| `PUBSUB_DRIVER`              | string   | `memory`       | `memory`, `redis`, `postgres`                       |
| `PUBSUB_INSTANCE_ID`         | string   | `""`           | Identificador de la instancia, debe ser único       |
| `PUBSUB_REDIS_ADDR`          | string   | `""`           | Dirección de Redis                                  |
| `PUBSUB_REDIS_PASSWORD`      | string   | `""`           | Contraseña de Redis                                 |
| `PUBSUB_REDIS_DB`            | number   | `1`            | Número de la base de datos de Redis                 |
| `PUBSUB_RETRY_ENABLED`       | bool     | `true`         | Reintentar la entrega en caso de error              |
| `PUBSUB_RETRY_MAX_RETRIES`   | number   | `3`            | Número máximo de reintentos                         |
| `PUBSUB_RETRY_INITIAL_DELAY` | duration | `100ms`        | Retraso inicial antes de un reintento               |
| `PUBSUB_RETRY_MAX_DELAY`     | duration | `5s`           | Retraso máximo antes de un reintento                |
| `PUBSUB_RETRY_MULTIPLIER`    | float    | `2.0`          | Factor por el que crece el retraso                  |
| `PUBSUB_DLQ_ENABLED`         | bool     | `false`        | Enviar los eventos no entregados a una cola aparte  |
| `PUBSUB_DLQ_DRIVER`          | string   | `memory`       | Almacenamiento de esa cola                          |
| `PUBSUB_DLQ_MAX_SIZE`        | number   | `1000`         | Tamaño máximo de la cola                            |

## Varios

| Variable                      | Tipo     | Predeterminado           | Propósito                                      |
|-------------------------------|----------|--------------------------|-------------------------------------------------|
| `LOGGER_LEVEL`                | string   | `info`                   | `debug`, `info`, `warn`, `error`                |
| `LOGGER_LOG_DB_QUERIES`       | bool     | `false`                  | Registrar las consultas a la base de datos. Solo para depuración |
| `DEFAULT_LANGUAGE`            | string   | `""`                     | Idioma predeterminado de la interfaz, por ejemplo `ru` |
| `GLOBAL_API_URL`              | string   | `https://api.gameap.com` | Dirección de la API global — actualizaciones de juegos |
| `GAMES_CDN_URLS`              | list     | véase más abajo          | Fuentes del catálogo de juegos, se prueban en orden |
| `TASK_REAPER_INTERVAL`        | duration | `1m`                     | Con qué frecuencia buscar tareas atascadas      |
| `TASK_REAPER_STALE_THRESHOLD` | duration | `10m`                    | Tiempo de inactividad tras el cual una tarea se considera atascada |

El valor predeterminado de `GAMES_CDN_URLS`:
`https://cdn.gameap.ru/games.json,https://cdn.gameap.com/games.json`.

## Variables sin efecto

Estas variables son analizadas por el panel pero no afectan a su comportamiento. No confíe en ellas:

| Variable                        | Problema                                                                                 |
|---------------------------------|-------------------------------------------------------------------------------------------|
| `AUTH_SESSION_IDLE_TIMEOUT`     | No existe terminación de sesión por inactividad. Una sesión vive todo su ciclo — 24 horas o 7 días |
| `AUTH_SESSION_IDLE_UPDATE_FREQ` | Igual que la anterior                                                                     |
| `GRPC_FILE_TRANSFER_BASE_PATH`  | El valor no se pasa a ninguna parte                                                       |
| `GRPC_ENABLED`                  | Esta variable no existe en absoluto. gRPC siempre está en marcha y no se puede desactivar. Las versiones antiguas de `gameapctl` añaden esta línea a `config.env` |

Los límites de frecuencia de inicio de sesión tampoco son configurables mediante variables de
entorno — están definidos como constantes en el código. Los valores figuran en la página
[Seguridad](/es/security.html).

## Ejemplo de configuración mínima

```dotenv
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://gameap:password@localhost:5432/gameap?sslmode=disable

AUTH_SECRET=replace_with_32_random_bytes
ENCRYPTION_KEY=replace_with_32_random_bytes

HTTP_HOST=panel.example.com
HTTP_PORT=8025
```

Todo lo demás recurre a los valores predeterminados.
