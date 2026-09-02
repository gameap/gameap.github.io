---
title: Seguridad
layout: default
lang: es
category: Administración
order: 331
---

La configuración de seguridad se define mediante variables de entorno en el archivo de configuración
del panel: `/etc/gameap/config.env` en Linux, `C:\gameap\web\config.env` en Windows. Después de
modificar el archivo, hay que reiniciar el panel: `gameapctl panel restart`.

## Autenticación de dos factores

### 2FA obligatorio para administradores

**Activado por defecto.** Un administrador sin autenticación de dos factores primero ve un
recordatorio y, después de 30 días, una exigencia de activarla, sin la cual no se puede usar el
panel.

| Variable                        | Valor por defecto | Propósito                                                              |
|---------------------------------|-------------------|------------------------------------------------------------------------|
| `AUTH_REQUIRE_MFA_FOR_ADMINS`   | `true`            | Exigir 2FA a los administradores. `false` desactiva el mecanismo por completo |
| `AUTH_MFA_HARD_FAIL_DAYS`       | `30`              | Cuántos días se conceden para activarla. `0` — solo recordatorio, sin bloqueo |
| `AUTH_MFA_ENROLLMENT_TOKEN_TTL` | `15m`             | Duración de la sesión restringida que se emite tras el plazo           |

Un administrador es un usuario que tiene el permiso **`admin roles & permissions`** concedido de
forma global, ya sea a través de un rol o directamente. En el panel no existe un "rol de
administrador" separado, y el usuario con `id = 1` no recibe privilegios por defecto.

### Cómo se cuentan los 30 días

La cuenta atrás comienza **en el primer inicio de sesión exitoso del administrador sin 2FA**, no en
la instalación del panel, ni en la fecha de actualización, ni cuando se activó la opción.

La fecha del primer recordatorio se guarda en el campo `metadata` de la tabla `users` bajo la clave
`mfa_first_shown_at`. Si un administrador nunca ha iniciado sesión en el panel, su cuenta atrás aún
no ha comenzado.

De esto se deduce que, tras actualizar una instalación existente, cada administrador dispone de los
30 días completos a partir de su próximo inicio de sesión.

El botón "Recordarme más tarde" pospone el diálogo 24 horas, pero **no mueve la fecha límite**.
Cerrar el diálogo con la cruz o con la tecla Esc funciona de la misma manera.

### Qué ocurre después del plazo

El inicio de sesión sigue funcionando, pero en lugar de una sesión normal se emite un token
restringido válido durante 15 minutos. Da acceso a solo cinco rutas, las necesarias para activar
2FA:

* `GET /api/config/public`
* `POST /api/auth/logout`
* `GET /api/profile`
* `POST /api/profile/2fa/setup`
* `POST /api/profile/2fa/confirm`

Todas las demás peticiones devuelven `403` con el mensaje `session is restricted to two-factor
enrollment`. La interfaz muestra un diálogo modal sin botón de cierre.

> Esto no es un bloqueo de cuenta. Active 2FA, inicie sesión de nuevo y el trabajo continúa.
> Un token restringido no puede "ascenderse" a uno completo, por lo que una vez activado 2FA el
> panel le cierra la sesión por sí mismo y le pide que inicie sesión de nuevo.

Dos excepciones que conviene conocer de antemano:

* **Una sesión completa emitida antes del plazo sigue funcionando** hasta que caduca. Un token de
  "recordarme" obtenido el día 29 seguirá siendo válido durante otros 7 días.
* **Los tokens de acceso personal (PAT) no se restringen en absoluto.** Un token emitido con
  antelación seguirá funcionando con la API después del plazo: un recurso razonable para la
  automatización.

### Activación de 2FA

Perfil → autenticación de dos factores. El panel muestra un código QR y un secreto para introducirlo
manualmente, después de lo cual se introduce el código de la aplicación de autenticación.

Los parámetros son fijos y no configurables: **TOTP según RFC 6238, HMAC-SHA1, 6 dígitos, período de
30 segundos**, tolerancia ±1 paso (es decir, aproximadamente ±30 segundos de desviación de reloj).
Esta combinación es compatible con todas las aplicaciones habituales: Google Authenticator, Authy,
1Password y otras. En la aplicación la entrada se llamará `GameAP` y el nombre de la cuenta será su
nombre de usuario.

Un código ya utilizado no puede volver a usarse, ni siquiera dentro de su ventana de 30 segundos.

Si el reloj del servidor o del teléfono se desvía más de medio minuto, los códigos dejan de
coincidir: compruebe la sincronización horaria en ambos lados.

### Códigos de recuperación

Al activar 2FA se emiten **10 códigos de recuperación** de la forma `abcde-fghjk`. El alfabeto no
contiene vocales ni caracteres fáciles de confundir (`0`, `o`, `1`, `l`, `i`).

* Cada código es de un solo uso.
* Los códigos se muestran **exactamente una vez**: al activar 2FA. En la base de datos solo se
  guardan sus hashes, por lo que no hay forma de volver a verlos.
* Un código de recuperación puede usarse tanto para iniciar sesión como para desactivar 2FA.
* Regeneración: perfil → regenerar códigos de recuperación; se le pedirá su contraseña. Todos los
  códigos anteriores quedan invalidados.

Guarde los códigos de inmediato, y no en el mismo gestor de contraseñas que guarda la contraseña del
panel.

### Pérdida de acceso

**No** existe un comando CLI integrado para restablecer 2FA: ni `gameapctl` ni el propio panel pueden
desactivar la autenticación de dos factores de otro usuario. Las opciones, en orden:

**1. Un código de recuperación.** Introdúzcalo en lugar del código de la aplicación.

**2. El plazo ha pasado, pero 2FA aún no está activado.** Esto no es una pérdida de acceso: inicie
sesión como de costumbre y termine de activarla; las páginas necesarias están disponibles.

**3. Eliminar la exigencia por completo.** En `config.env`:

```dotenv
AUTH_REQUIRE_MFA_FOR_ADMINS=false
```

y `gameapctl panel restart`. La restricción se levanta y los tokens restringidos ya emitidos se
convierten en completos.

> Esta opción **no ayudará a quien ya activó TOTP y perdió su dispositivo**: el segundo factor se
> verifica antes de comprobar la exigencia, por lo que el panel seguirá pidiendo un código.

**4. Mantener el recordatorio, pero eliminar el bloqueo.**

```dotenv
AUTH_MFA_HARD_FAIL_DAYS=0
```

**5. Edición de la base de datos.** La única opción que queda cuando se han perdido tanto el
dispositivo como los códigos de recuperación. Detenga el panel, haga una copia de seguridad de la
base de datos y ejecute la consulta.

PostgreSQL:

```sql
UPDATE users
   SET two_factor_enabled = false,
       two_factor_secret = NULL,
       two_factor_recovery_codes = NULL,
       two_factor_last_used_step = NULL
 WHERE login = 'admin';
```

MySQL y SQLite: lo mismo, pero con `two_factor_enabled = 0`.

Para restablecer también la cuenta atrás de 30 días, elimine la clave `mfa_first_shown_at` del campo
`metadata`. En PostgreSQL, donde este campo es de tipo `JSONB`:

```sql
UPDATE users SET metadata = metadata - 'mfa_first_shown_at' WHERE login = 'admin';
```

En MySQL el campo se almacena como texto JSON, y la clave se elimina así:

```sql
UPDATE users SET metadata = JSON_REMOVE(metadata, '$.mfa_first_shown_at') WHERE login = 'admin';
```

En SQLite, a partir de la versión 3.38:

```sql
UPDATE users SET metadata = json_remove(metadata, '$.mfa_first_shown_at') WHERE login = 'admin';
```

> No limpie el campo `metadata` por completo (`SET metadata = NULL`): además de la cuenta atrás de
> 2FA puede contener otros datos del usuario, y se perderían.

Después, inicie el panel y active 2FA de nuevo.

## Contraseñas

Requisitos de contraseña: **no menos de 12 ni más de 128 bytes**. Deliberadamente no hay requisitos
de composición (letras mayúsculas, dígitos, caracteres especiales); en su lugar se utiliza una
comprobación contra una lista de contraseñas comprometidas.

> El límite se cuenta en bytes, no en caracteres. Una contraseña de 12 letras cirílicas ocupa 24
> bytes y supera la comprobación con margen de sobra.

La lista de contraseñas comunes se toma de [SecLists](https://github.com/danielmiessler/SecLists)
(`xato-net-10-million-passwords`), filtrada por longitud, y contiene unas 46 000 entradas. Está
**integrada en el binario**: el panel no contacta con nada al comprobar una contraseña ni transmite
nada sobre ella.

| Variable                    | Valor por defecto | Propósito                                                              |
|-----------------------------|-------------------|------------------------------------------------------------------------|
| `AUTH_ALLOW_WEAK_PASSWORDS` | `false`           | Desactiva solo la comprobación por lista. Los límites de longitud se mantienen |
| `AUTH_BCRYPT_COST`          | `13`              | Coste de bcrypt. Rango permitido de 10 a 14; en caso contrario el panel no arrancará |

Las contraseñas se hashean con bcrypt sobre un SHA-256 preliminar, de modo que el límite de 72 bytes
de bcrypt no trunca las contraseñas largas. Al iniciar sesión, un hash con un coste inferior al
actual se vuelve a hashear automáticamente; el coste de los hashes ya almacenados no puede
reducirse, aunque disminuya el valor de la variable.

La comprobación se aplica cuando un administrador crea un usuario, cuando se modifica un usuario y
cuando se cambia la propia contraseña. No se aplica al iniciar sesión; de lo contrario, los usuarios
con contraseñas débiles antiguas perderían el acceso.

> La contraseña del primer administrador, establecida con la variable `ADMIN_PASSWORD` durante el
> llenado inicial de la base de datos, **no** se comprueba. Elija de forma consciente.

## CAPTCHA

Desactivado por defecto. Protege **únicamente** el formulario de inicio de sesión
(`POST /api/auth/login`); la verificación del segundo factor no está cubierta por el CAPTCHA.

| Variable              | Valor por defecto | Propósito                                                            |
|-----------------------|-------------------|----------------------------------------------------------------------|
| `CAPTCHA_PROVIDER`    | `""`              | `recaptcha_v2`, `recaptcha_v3` o `turnstile`. Vacío — desactivado    |
| `CAPTCHA_SITE_KEY`    | `""`              | Clave pública, se envía al navegador                                 |
| `CAPTCHA_SECRET_KEY`  | `""`              | Clave secreta, nunca se envía al exterior                            |
| `CAPTCHA_MIN_SCORE`   | `0.5`             | Umbral solo para reCAPTCHA v3, ignorado por los demás proveedores    |
| `CAPTCHA_FAIL_OPEN`   | `false`           | Si se permite el inicio de sesión cuando el servicio de verificación no está disponible |
| `CAPTCHA_VERIFY_URL`  | `""`              | Dirección de verificación personalizada, para proxyar el tráfico saliente |

Con `CAPTCHA_FAIL_OPEN=false`, un servicio de verificación no disponible significa un `503` al
iniciar sesión en el panel.

> Si define `CAPTCHA_PROVIDER` pero no `CAPTCHA_SECRET_KEY`, el widget aparecerá en el formulario de
> inicio de sesión, pero la verificación del lado del servidor permanecerá desactivada
> silenciosamente. Defina ambas variables juntas.

El panel añade por sí mismo los dominios del proveedor seleccionado a la política CSP; no hace falta
configurar `SECURITY_CSP_EXTRA_SCRIPT_SRC` adicionalmente.

### reCAPTCHA v3

Las claves se emiten en la [consola de reCAPTCHA](https://www.google.com/recaptcha/admin): registre
el sitio, elija el tipo **reCAPTCHA v3** y especifique el dominio del panel.

```dotenv
CAPTCHA_PROVIDER=recaptcha_v3
CAPTCHA_SITE_KEY=6LcExampleSiteKeyExampleSiteKeyExam
CAPTCHA_SECRET_KEY=6LcExampleSecretKeyExampleSecretKeyEx
CAPTCHA_MIN_SCORE=0.5
```

reCAPTCHA v3 no pregunta nada al usuario: devuelve una puntuación de `0.0` a `1.0`, donde uno
significa casi con seguridad un humano. El inicio de sesión se rechaza si la puntuación es inferior
a `CAPTCHA_MIN_SCORE`.

Comience con el valor por defecto de `0.5` y ajústelo según lo requieran las circunstancias: si la
gente se queja de que no puede iniciar sesión, bájelo; si los intentos de adivinar contraseñas
continúan, súbalo. Un valor superior a `0.7` estorba notablemente a los usuarios con bloqueadores en
el navegador y en modo incógnito.

El panel debe abrirse en el dominio especificado en la configuración de la clave; de lo contrario la
verificación fallará. Para varios dominios, enumérelos todos en la consola de reCAPTCHA.

La verificación se realiza con una petición a `https://www.google.com/recaptcha/api/siteverify`;
esta dirección debe ser accesible desde el servidor del panel.

### Turnstile

Las claves se emiten en el panel de [Cloudflare](https://dash.cloudflare.com/), en la sección
**Turnstile**. Basta una cuenta gratuita y no hace falta delegar el dominio a Cloudflare.

```dotenv
CAPTCHA_PROVIDER=turnstile
CAPTCHA_SITE_KEY=0x4AAAAAAAExampleSiteKey
CAPTCHA_SECRET_KEY=0x4AAAAAAAExampleSecretKey
```

`CAPTCHA_MIN_SCORE` no se aplica a Turnstile: el proveedor devuelve solo "aprobado" o "rechazado",
por lo que no hace falta definir esta variable.

En la mayoría de los casos Turnstile se supera sin que el usuario lo note, y muestra un reto breve
cuando algo parece sospechoso. El modo del widget (**Managed**, **Non-interactive** o **Invisible**)
se elige del lado de Cloudflare al crear la clave; no es configurable desde el panel.

La verificación se realiza con una petición a
`https://challenges.cloudflare.com/turnstile/v0/siteverify`.

reCAPTCHA v2 se configura igual que v3, pero sin `CAPTCHA_MIN_SCORE`.

## Limitación de intentos de inicio de sesión

Los límites están **fijados en el código**; no hay variables de entorno para ellos:

* ventana: 15 minutos;
* no más de 20 intentos fallidos desde una misma dirección IP;
* no más de 5 intentos fallidos por nombre de usuario.

El límite se aplica en dos rutas: `POST /api/auth/login` y `POST /api/auth/2fa/verify`. El cliente
recibe un `429` y una cabecera `Retry-After: 900`. Solo las respuestas `401` incrementan el
contador; un inicio de sesión exitoso restablece el contador del nombre de usuario, pero no el de la
dirección IP.

Los contadores se almacenan en la caché del panel. Con `CACHE_DRIVER=memory` (el valor por defecto)
se restablecen al reiniciar y no se comparten entre varias instancias del panel; para una
instalación tolerante a fallos utilice `CACHE_DRIVER=redis`.

Aparte de esto, la verificación del segundo factor permite no más de 5 intentos de introducir el
código por intento de inicio de sesión, tras lo cual hay que introducir la contraseña de nuevo.

El panel no tiene bloqueo de cuentas: la fuerza bruta está limitada únicamente por lo descrito
arriba.

## Cabeceras HTTP y Content Security Policy

Las cabeceras de seguridad están activadas por defecto.

| Variable                            | Valor por defecto                 | Cabecera                                         |
|-------------------------------------|-----------------------------------|--------------------------------------------------|
| `SECURITY_HEADERS_ENABLED`          | `true`                            | Interruptor general                              |
| `SECURITY_CONTENT_TYPE_OPTIONS`     | `true`                            | `X-Content-Type-Options: nosniff`                |
| `SECURITY_FRAME_OPTIONS`            | `SAMEORIGIN`                      | `X-Frame-Options`; un valor vacío elimina la cabecera |
| `SECURITY_REFERRER_POLICY`          | `strict-origin-when-cross-origin` | `Referrer-Policy`                                |
| `SECURITY_HSTS_ENABLED`             | `true`                            | `Strict-Transport-Security`                      |
| `SECURITY_HSTS_MAX_AGE`             | `31536000`                        | Duración en segundos, un año por defecto         |
| `SECURITY_HSTS_INCLUDE_SUBDOMAINS`  | `false`                           | Añade `includeSubDomains`                        |
| `SECURITY_HSTS_PRELOAD`             | `false`                           | Añade `preload`                                  |

HSTS solo se envía cuando se accede al panel por HTTPS: por una conexión TLS real, por la cabecera
`X-Forwarded-Proto: https` o con `TLS_FORCE_HTTPS=true`. Trabajar por HTTP simple durante el
desarrollo no dejará el navegador "atascado".

### Política CSP

| Variable                            | Valor por defecto | Propósito                                                             |
|-------------------------------------|-------------------|-----------------------------------------------------------------------|
| `SECURITY_CSP_ENABLED`              | `true`            | Activa la política                                                    |
| `SECURITY_CSP_REPORT_ONLY`          | `false`           | Envía `Content-Security-Policy-Report-Only` en lugar de la política de bloqueo |
| `SECURITY_CSP_POLICY`               | `""`              | Sustituye por completo la política generada                           |
| `SECURITY_CSP_REPORT_URI`           | `""`              | Añade `report-uri`                                                    |
| `SECURITY_CSP_EXTRA_SCRIPT_SRC`     | `""`              | Amplía `script-src`, valores separados por comas                      |
| `SECURITY_CSP_EXTRA_STYLE_SRC`      | `""`              | Amplía `style-src`                                                    |
| `SECURITY_CSP_EXTRA_CONNECT_SRC`    | `""`              | Amplía `connect-src`                                                  |
| `SECURITY_CSP_EXTRA_IMG_SRC`        | `""`              | Amplía `img-src`                                                      |
| `SECURITY_CSP_EXTRA_FRAME_SRC`      | `""`              | Amplía `frame-src`                                                    |
| `SECURITY_CSP_EXTRA_FONT_SRC`       | `""`              | Amplía `font-src`                                                     |

La política generada:

```text
default-src 'self'; base-uri 'self'; object-src 'none'; frame-ancestors 'self'; form-action 'self';
script-src 'self' blob: 'wasm-unsafe-eval' <hashes of inline scripts>;
style-src 'self' 'unsafe-inline';
img-src 'self' data: blob:;
font-src 'self';
connect-src 'self';
frame-src 'self';
worker-src 'self' blob:
```

Por qué contiene relajaciones:

* `'wasm-unsafe-eval'`: al subir archivos, la suma de comprobación SHA-256 se calcula en el
  navegador mediante WebAssembly. Sin este permiso la subida de archivos dejará de funcionar.
* `blob:` en `script-src`: así se cargan las partes frontend de los plugins.
* `'unsafe-inline'` en `style-src`: los estilos de los plugins y de Vue se añaden a la página en
  línea.

Las direcciones del proveedor de CAPTCHA seleccionado se añaden a `script-src` y `frame-src`
automáticamente.

> **Un plugin que cargue scripts desde un CDN de terceros será bloqueado por la política.** Añada el
> dominio necesario a `SECURITY_CSP_EXTRA_SCRIPT_SRC`. Es mejor no usar `SECURITY_CSP_POLICY` para
> esto: sustituye toda la política, junto con la autorización automática de los dominios de CAPTCHA
> y los hashes de los scripts en línea del panel.

## Cifrado de secretos

| Variable         | Obligatoria | Propósito                                                           |
|------------------|-------------|---------------------------------------------------------------------|
| `AUTH_SECRET`    | sí          | Clave de firma de los tokens de sesión. Sin ella el panel no arrancará |
| `ENCRYPTION_KEY` | no          | Clave de cifrado de los secretos en la base de datos                |

Ambos valores deben ser aleatorios, no una frase de contraseña. Pero sus requisitos de longitud
son **distintos**: el panel los trata de forma diferente.

**`AUTH_SECRET` se usa tal cual y se ajusta exactamente a 32 bytes:** un valor más corto se rellena,
uno más largo se **trunca**, y al registro solo llega una advertencia. Por eso indique exactamente
32 caracteres:

```bash
openssl rand -base64 24
```

No use aquí `openssl rand -hex 32`: obtendrá 64 caracteres, el panel conservará solo los primeros
32, y esos 32 caracteres hex llevan apenas 16 bytes aleatorios: menos que los 24 bytes aleatorios
de `openssl rand -base64 24`.

**`ENCRYPTION_KEY` se hashea por completo con SHA-256**, su longitud no está limitada y no se pierde
nada. Aquí puede usar un valor más largo:

```bash
openssl rand -hex 32
```

El hasheo conserva la entropía del valor original, pero no la aumenta, así que la clave debe ser
aleatoria de todos modos. Una frase de contraseña no es segura aquí: puede obtenerse por fuerza
bruta si el valor cifrado se filtra.

> `AUTH_SECRET` se ajusta silenciosamente a 32 bytes: un valor más corto se rellena, uno más largo
> se trunca, y solo una advertencia llega al registro. Un `AUTH_SECRET` corto o predecible significa
> que los tokens de sesión pueden falsificarse.

`ENCRYPTION_KEY` se utiliza para cifrar los secretos TOTP y la contraseña de conexión del daemon en
la base de datos. Si no está definida, los secretos TOTP se cifran con una clave derivada de
`AUTH_SECRET`, y la contraseña del daemon se almacena en texto plano; el panel advierte de ello al
arrancar.

> **No defina `ENCRYPTION_KEY` por primera vez en una instalación en funcionamiento que ya tenga 2FA
> activado.** La clave de cifrado de los secretos TOTP cambiará de `AUTH_SECRET` a `ENCRYPTION_KEY`,
> los secretos previamente almacenados quedarán ilegibles y **todos los usuarios tendrán que activar
> 2FA de nuevo**. Los códigos de recuperación seguirán funcionando: se almacenan por separado.
>
> Lo mismo ocurrirá si `ENCRYPTION_KEY` se pierde o se cambia. Guárdela junto con la copia de
> seguridad de la base de datos: sin ella, parte de los datos de la copia no podrá restaurarse.

Las contraseñas de los usuarios se hashean con bcrypt; los tokens personales y la clave del daemon
se almacenan como SHA-256: son transformaciones irreversibles y `ENCRYPTION_KEY` no tiene nada que
ver con ellas.

## Registro de auditoría

| Variable                 | Valor por defecto | Propósito                                          |
|--------------------------|-------------------|----------------------------------------------------|
| `AUDIT_ENABLED`          | `true`            | Registro de eventos de seguridad                   |
| `AUDIT_CLIENT_IP_HEADER` | `""`              | Cabecera con la IP real del cliente, por ejemplo `X-Real-IP` |

Qué se registra: inicios de sesión exitosos y fallidos, alcances del límite de intentos, denegaciones
de acceso, activación y desactivación de 2FA, regeneración de códigos de recuperación, cambios de
usuarios y asignaciones de roles, creación y revocación de tokens, cambios y eliminación de
servidores dedicados, operaciones con archivos, instalación y eliminación de plugins.

Cada registro contiene: tipo de evento, categoría, resultado, identificador y nombre del usuario que
actúa, método de autenticación, dirección IP, User-Agent, método y ruta de la petición,
identificador de la petición.

> El registro de auditoría es un conjunto de líneas estructuradas del registro de la aplicación con
> el campo `component=audit`. **No** hay tabla separada en la base de datos, ni archivo separado, ni
> rotación, ni interfaz de consulta, ni API de lectura. Si los registros necesitan almacenarse y
> buscarse, configure la recolección del registro del panel con las herramientas estándar de su
> sistema, por ejemplo, mediante `journald` y un recolector de registros externo.
>
> `AUDIT_CLIENT_IP_HEADER` confía en la cabecera especificada procedente de **cualquier** remitente:
> el panel no tiene lista de proxies de confianza. Active esta variable solo si el proxy inverso
> sobrescribe garantizadamente la cabecera en las peticiones entrantes. De lo contrario, la
> dirección IP puede falsearse y, con ella, eludirse el límite de intentos de inicio de sesión por
> IP.

## Sesiones y tokens

Las sesiones se emiten en formato PASETO v4.local (`AUTH_SERVICE=paseto`, la alternativa es `jwt`).
Una sesión normal dura 24 horas; con la opción "recordarme", 7 días. Al cerrar sesión, el token
entra en una lista de revocación que se comprueba en cada petición.

Para los casos en que el token debe pasarse en la dirección de la página —conexiones WebSocket,
descargas de archivos— se emiten tokens de un solo uso y corta duración con el prefijo `glst_`. Su
vida útil está limitada a 10 segundos independientemente del valor de
`AUTH_SHORT_LIVED_TOKEN_TTL`.

El panel no tiene protección CSRF y no la necesita: la autenticación se realiza mediante la cabecera
`Authorization`, no con cookies.

La lista de orígenes permitidos para acceder a la API desde un navegador se define con la variable
`HTTP_ALLOWED_ORIGINS` (valores separados por comas). Si está vacía, se permite un único origen
calculado a partir de `HTTP_HOST`. El carácter `*` no está soportado.

## Subida de archivos

El tipo de un archivo subido se determina por su contenido, no por la extensión ni por la cabecera
enviada por el cliente. Por defecto se permiten imágenes, archivos de texto, JSON, XML, CSV, YAML y
PDF. SVG y HTML están prohibidos deliberadamente: pueden contener scripts.

| Variable                      | Valor por defecto | Propósito                                               |
|-------------------------------|-------------------|---------------------------------------------------------|
| `FILES_UPLOAD_ALLOWED_MIMES`  | `""`              | Amplía la lista de tipos permitidos, no la sustituye    |
| `FILES_UPLOAD_ALLOW_ARCHIVES` | `false`           | Permitir archivos comprimidos: zip, tar, gzip, bzip2, 7z, xz |
| `FILES_UPLOAD_ALLOW_BINARY`   | `false`           | Permitir archivos binarios arbitrarios                  |

> La prohibición de archivos comprimidos y binarios es la causa más común de la pregunta "¿por qué
> no se sube el archivo?". Un archivo comprimido puede contener ejecutables que se desempaquetarán
> en el servidor dedicado, por eso la subida está prohibida por defecto. Active estas opciones de
> forma consciente.

Las subidas rechazadas pasan al registro de auditoría con el tipo de archivo detectado y el motivo
del rechazo. El límite de tamaño por archivo es de 100 MB y no es configurable.

## Plugins

Los plugins se ejecutan en un sandbox de WebAssembly y no tienen acceso directo al sistema. Las
peticiones de red de los plugins se restringen por separado:

| Variable                          | Valor por defecto | Propósito                                                 |
|-----------------------------------|-------------------|-----------------------------------------------------------|
| `PLUGINS_DISABLED`                | `false`           | Desactivar por completo el mecanismo de plugins           |
| `PLUGIN_HTTP_BLOCK_PRIVATE_IPS`   | `true`            | Prohibir peticiones a direcciones de redes internas       |
| `PLUGIN_HTTP_ALLOWED_SCHEMES`     | `https`           | Esquemas permitidos                                       |
| `PLUGIN_HTTP_ALLOWED_HOSTS`       | `""`              | Lista de hosts permitidos; vacío — sin restricciones de host |
| `PLUGIN_HTTP_MAX_TIMEOUT_SECONDS` | `30`              | Límite de tiempo de la petición                           |
| `PLUGIN_HTTP_MAX_REDIRECTS`       | `5`               | Límite de redirecciones; cada una se comprueba de nuevo   |

Las direcciones de los servicios de metadatos de los proveedores de nube están siempre bloqueadas y
no pueden desbloquearse. Las cabeceras `Set-Cookie`, `Authorization`, `WWW-Authenticate` y
`Clear-Site-Data` no se transmiten al plugin.

Consulte los detalles en la página [Plugins](/es/plugins/index.html).

## Qué está fijado en el código

Algunas variables están presentes en la configuración pero no tienen efecto en el comportamiento del
panel. No confíe en ellas:

* **`AUTH_SESSION_IDLE_TIMEOUT` y `AUTH_SESSION_IDLE_UPDATE_FREQ`**: actualmente no hay cierre de
  sesión por inactividad. Una sesión dura exactamente su plazo completo: 24 horas o 7 días.
* **`GRPC_ENABLED`**: no existe tal opción; el servidor gRPC siempre está en funcionamiento. Las
  versiones antiguas de `gameapctl` añaden esta línea a `config.env`; es inofensiva, pero no hace
  nada.
* Los límites de intentos de inicio de sesión están definidos por constantes en el código; no hay
  variables `RATE_LIMIT*`.

Los porcentajes de cumplimiento de OWASP ASVS en los archivos `docs/security/ASVS.md` y
`ASVS_L2.md` del repositorio del panel están desactualizados; consulte esta página en su lugar.

## Notificación de vulnerabilidades

Las vulnerabilidades se aceptan a través de un
[GitHub Security Advisory](https://github.com/gameap/gameap/security/advisories/new) —es el canal
preferido— o por correo electrónico a `security@gameap.com`.

Plazos: acuse de recibo — 72 horas, evaluación — 14 días, divulgación pública — 90 días.
Correcciones: críticas — 14 días, altas — 30 días, medias — 60 días, bajas — en la próxima versión.

Una configuración involuntariamente insegura establecida por el propio administrador (por ejemplo,
`AUTH_ALLOW_WEAK_PASSWORDS=true` o `SECURITY_HEADERS_ENABLED=false`) no se considera una
vulnerabilidad.
