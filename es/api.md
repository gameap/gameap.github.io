---
title: API y tokens
layout: default
lang: es
category: Administración
order: 336
---

El panel se controla completamente a través de la API HTTP: la interfaz funciona a través de la misma API.
Una descripción completa de los métodos con los esquemas de solicitud y respuesta está disponible en
[openapi.gameap.io](https://openapi.gameap.io/).

Esta página explica cómo autenticarse en la API.

## Métodos de autenticación

| Método                          | Para qué sirve                                       | Duración            |
|---------------------------------|------------------------------------------------------|---------------------|
| Token de sesión                 | La interfaz                                          | 24 horas o 7 días   |
| Token de acceso personal (PAT)  | Scripts, integraciones, automatización               | Indefinida          |
| Token de corta duración         | WebSocket y descarga de archivos                     | 10 segundos         |

El token se pasa en una cabecera:

```http
Authorization: Bearer <token>
```

## Tokens de acceso personal

Este es el método principal para la automatización: el token no está vinculado a una sesión, no caduca
y tiene su propio conjunto de permisos.

### Creación

Un token se puede crear en el perfil o con una solicitud:

```bash
curl -X POST https://panel.example.com:8025/api/tokens \
  -H "Authorization: Bearer <session token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "ci-deploy", "abilities": ["server:list", "server:restart"]}'
```

La respuesta contiene el token completo:

```json
{"token": "12|kJ3n8sQm..."}
```

> El token se muestra **una sola vez**. En la base de datos solo se almacena su hash SHA-256, y el
> valor no se puede recuperar: si lo pierde, emita uno nuevo.

El formato del token es `{id}|{secret}`. El separador `|` es obligatorio: pase el valor completo,
exactamente como fue emitido.

### Permisos (abilities) del token

Los permisos se especifican en el momento de la creación y restringen el token independientemente de los
permisos del usuario: el token no puede hacer más de lo que se le permite, ni más de lo que se le permite
a su propietario.

| Permiso                   | Qué permite                              |
|---------------------------|------------------------------------------|
| `server:list`             | Ver la lista de servidores               |
| `server:start`            | Iniciar un servidor                      |
| `server:stop`             | Detener un servidor                      |
| `server:restart`          | Reiniciar un servidor                    |
| `server:update`           | Actualizar un servidor                   |
| `server:console`          | Leer y escribir en la consola            |
| `server:rcon-console`     | Consola RCON                             |
| `server:rcon-players`     | Gestionar jugadores a través de RCON     |
| `server:tasks-manage`     | Gestionar las tareas del servidor        |
| `server:settings-manage`  | Gestionar la configuración del servidor  |
| `admin:server:create`     | Crear servidores                         |
| `admin:gdaemon-task:read` | Leer las tareas del daemon               |

Los permisos con el prefijo `admin:` solo los puede conceder un administrador: el intento de un usuario
normal de añadirlos falla.

La lista actual está disponible con:

```http
GET /api/tokens/abilities
```

### Listado y revocación

```http
GET    /api/tokens        — lista de sus tokens
DELETE /api/tokens/{id}   — revocar un token
```

La lista muestra el nombre, los permisos y la hora del último uso, lo que resulta útil para encontrar
tokens sin utilizar.

> **Cambiar la contraseña revoca los tokens.** Todos los tokens personales creados antes del cambio de
> contraseña dejan de funcionar. Después de cambiar la contraseña, emita los tokens de nuevo.

### Ejemplo de uso

```bash
TOKEN='12|kJ3n8sQm...'

# lista de servidores
curl -H "Authorization: Bearer $TOKEN" \
  https://panel.example.com:8025/api/servers

# reiniciar un servidor
curl -X POST -H "Authorization: Bearer $TOKEN" \
  https://panel.example.com:8025/api/servers/1/restart
```

## Token de sesión

Se emite al iniciar sesión con nombre de usuario y contraseña:

```bash
curl -X POST https://panel.example.com:8025/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login": "admin", "password": "..."}'
```

El formato del token es PASETO v4.local. Una sesión normal dura 24 horas, o 7 días con la opción
«recordarme» activada. Al cerrar sesión (`POST /api/auth/logout`), el token se añade a la lista de
revocación, que se comprueba en cada solicitud.

Si el usuario tiene activada la autenticación de dos factores, el inicio de sesión devuelve
`two_factor_required` junto con un `challenge_token` en lugar de un token; el segundo factor se confirma
con `POST /api/auth/2fa/verify`. Consulte [Seguridad](/es/security.html) para más detalles.

Los tokens de sesión son incómodos para la automatización: caducan, y el inicio de sesión está protegido
por limitación de velocidad (rate limiting) y posiblemente por un CAPTCHA. Utilice tokens personales.

## Tokens de corta duración

Se necesitan donde el token tiene que pasarse en la dirección de la página: conexiones WebSocket y
descargas de archivos. Se emiten con `POST /api/auth/short-lived-token`, llevan el prefijo `glst_`,
son de un solo uso y viven no más de 10 segundos independientemente de la configuración.

Solo estos tokens se aceptan en el parámetro de consulta `?token=`; un token personal no se puede pasar
allí: esto evita que aparezca en los registros del servidor web y en el historial del navegador.

## Límites y códigos de respuesta

| Código | Motivo                                                                  |
|--------|-------------------------------------------------------------------------|
| `401`  | El token falta, no es válido o está revocado                            |
| `403`  | El token o el usuario no tienen permisos                                |
| `422`  | Error de validación de la solicitud                                     |
| `429`  | Se superó el límite de intentos de inicio de sesión                     |

La limitación de velocidad se aplica solo al inicio de sesión y a la verificación del segundo factor:
20 intentos fallidos por dirección y 5 por nombre de usuario en 15 minutos. Las solicitudes con un token
personal no están sujetas a limitación de velocidad.

## CORS

Si la API se llama desde un navegador en un origen diferente, indique los orígenes permitidos en
`HTTP_ALLOWED_ORIGINS`, en forma completa, incluido el esquema. El comodín `*` no está soportado.
Consulte la [referencia de config.env](/es/config.html).
