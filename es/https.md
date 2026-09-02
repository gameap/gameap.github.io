---
title: HTTPS y certificados
layout: default
lang: es
category: Administración
order: 333
---

El panel puede servir HTTPS por sí mismo, sin un proxy inverso. El certificado puede tomarse de
archivos, definirse directamente en la configuración u obtenerse automáticamente a través de Let's Encrypt.

Los ajustes se definen en `config.env` — `/etc/gameap/config.env` en Linux, `C:\gameap\web\config.env`
en Windows. Después de un cambio es necesario reiniciar: `gameapctl panel restart`.

## Fuentes de certificados

El panel elige la fuente por sí mismo, en este orden:

1. **ACME** — si `ACME_ENABLED=true`, `ACME_EMAIL` y `ACME_DOMAINS` están definidos.
2. **Archivos** — si `TLS_CERT_FILE` y `TLS_KEY_FILE` están definidos.
3. **Valores en la configuración** — si `TLS_CERT` y `TLS_KEY` están definidos.
4. **Sin certificado** — el panel sirve solo HTTP.

Lo que se comprueba son las **parejas**: un `TLS_CERT_FILE` solo, sin `TLS_KEY_FILE`, no cuenta
como una fuente configurada y se ignora silenciosamente.

HTTPS escucha en el puerto de `HTTPS_PORT` (`443` por defecto) y solo cuando hay un certificado
disponible. HTTP en `HTTP_PORT` (`8025` por defecto) está siempre activo.

> El certificado del panel no tiene relación con los certificados gRPC que el panel usa para
> comunicarse con los daemons. Estos se emiten automáticamente mediante una autoridad de
> certificación interna; ACME no se aplica a ellos.
> Consulte [GRPC API](/es/daemon/grpc.html) para más detalles.

## Certificado desde archivos

```dotenv
TLS_CERT_FILE=/etc/gameap/certs/panel.crt
TLS_KEY_FILE=/etc/gameap/certs/panel.key
HTTPS_PORT=443
```

El archivo del certificado debe contener la cadena completa: el certificado mismo y después los
intermedios. Sin los intermedios, algunos clientes no podrán verificar la firma.

Los archivos se leen al arrancar. Después de reemplazar el certificado, reinicie el panel — no
vigila los archivos por sí mismo para detectar cambios.

## Certificado directamente en la configuración

Es conveniente cuando la configuración se despliega mediante un sistema de gestión de secretos y
no se desean archivos adicionales en el disco.

```dotenv
TLS_CERT=LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t...
TLS_KEY=LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0t...
```

Se aceptan tanto PEM plano como PEM codificado en base64 — el panel detecta el formato por sí
mismo. Dado que el formato de `config.env` no admite valores multilínea, en la práctica se usa
base64:

```bash
base64 -w0 panel.crt
base64 -w0 panel.key
```

## Let's Encrypt

El panel tiene un cliente ACME integrado: el certificado se emite y se renueva sin herramientas
externas como certbot.

| Variable                      | Valor predeterminado | Propósito                                                  |
|-------------------------------|-----------------|-------------------------------------------------------------|
| `ACME_ENABLED`                | `false`         | Activa la emisión automática                                |
| `ACME_EMAIL`                  | `""`            | Dirección para notificaciones de caducidad. Obligatorio     |
| `ACME_DOMAINS`                | `""`            | Dominios separados por comas. Obligatorio                   |
| `ACME_CHALLENGE_TYPE`         | `http-01`       | Tipo de desafío: `http-01` o `dns-01`                       |
| `ACME_DNS_PROVIDER`           | `""`            | Proveedor DNS, solo para `dns-01`                           |
| `ACME_DIRECTORY_URL`          | ACME de producción | URL del directorio ACME                                  |
| `ACME_RENEWAL_THRESHOLD`      | `720h`          | Cuánto tiempo antes de la caducidad renovar. 30 días por defecto |
| `ACME_RENEWAL_CHECK_INTERVAL` | `12h`           | Con qué frecuencia comprobar la fecha de caducidad          |
| `ACME_PROPAGATION_TIMEOUT`    | `180s`          | Cuánto tiempo esperar la propagación del registro DNS con `dns-01` |
| `ACME_STORAGE_PATH`           | `acme`          | Directorio para los certificados y la clave de la cuenta ACME |

> ACME se activa **solo** cuando `ACME_ENABLED=true`, `ACME_EMAIL` y `ACME_DOMAINS` están todos
> definidos, y para `dns-01` también `ACME_DNS_PROVIDER`. Si falta algo, el panel arranca
> sin ACME y sin error. Compruebe el resultado a través del estado en el área de administración.

### Desafío http-01

El método por defecto. No requiere nada, salvo que el panel sea accesible desde internet.

```dotenv
ACME_ENABLED=true
ACME_EMAIL=admin@example.com
ACME_DOMAINS=panel.example.com
ACME_CHALLENGE_TYPE=http-01
TLS_FORCE_HTTPS=true
```

Qué se necesita:

* los dominios de `ACME_DOMAINS` deben resolverse a la dirección de este servidor;
* **las peticiones al puerto 80 deben llegar al panel** — ese es el puerto al que se conecta la
  autoridad de certificación;
* el panel sirve la ruta `/.well-known/acme-challenge/` por sí mismo, en su puerto HTTP.

> Por defecto el panel escucha en el puerto **8025**, mientras que la autoridad de certificación
> siempre se conecta al puerto **80**. Por sí solos, estos no coinciden. Defina `HTTP_PORT=80` o
> redirija el puerto 80 al puerto del panel mediante herramientas del sistema:
>
> ```bash
> iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 8025
> ```
>
> Esta es la causa más común de que falle la emisión con `http-01`.

El puerto 80 se necesita no solo para la primera emisión: el desafío se repite en cada renovación,
por lo que no debe cerrarse después.

El método `http-01` **no emite certificados wildcard** (`*.example.com`) — para esos se requiere
`dns-01`.

### Desafío dns-01

Es necesario cuando el panel no es accesible desde internet en el puerto 80 o cuando se requiere
un certificado wildcard.

De los proveedores integrados, solo se admite **Cloudflare**. Otros se añaden mediante plugins: en
ese caso `ACME_DNS_PROVIDER` se define como `<plugin-id>:<provider-name>`.

```dotenv
ACME_ENABLED=true
ACME_EMAIL=admin@example.com
ACME_DOMAINS=panel.example.com,*.example.com
ACME_CHALLENGE_TYPE=dns-01
ACME_DNS_PROVIDER=cloudflare
CLOUDFLARE_DNS_API_TOKEN=token_from_the_cloudflare_dashboard
```

El token se crea en Cloudflare con el permiso **Zone → DNS → Edit** para la zona en cuestión.
Además de `CLOUDFLARE_DNS_API_TOKEN`, se aceptan las variables `CF_DNS_API_TOKEN`,
`CLOUDFLARE_API_TOKEN` y `CF_API_TOKEN`, así como la pareja heredada
«global key + email»: `CLOUDFLARE_API_KEY` junto con `CLOUDFLARE_EMAIL`. Es preferible un token
con permisos restringidos.

Si los registros DNS se propagan lentamente, aumente `ACME_PROPAGATION_TIMEOUT`.

### Configuración mediante gameapctl

En lugar de editar `config.env` a mano, puede usar el asistente:

```bash
gameapctl panel letsencrypt setup
```

Este solicita los dominios, la dirección de correo electrónico y el tipo de desafío, escribe los
ajustes en `config.env` y reinicia el panel.

La misma llamada sin preguntas:

```bash
gameapctl panel letsencrypt setup --non-interactive \
  --domains=panel.example.com \
  --email=admin@example.com \
  --challenge=http-01
```

Opciones útiles:

| Opción              | Propósito                                                                     |
|---------------------|-------------------------------------------------------------------------------|
| `--challenge`       | `http-01` o `dns-01`                                                          |
| `--domains`         | Dominios separados por comas                                                  |
| `--email`           | Dirección de la cuenta ACME                                                   |
| `--dns-provider`    | Proveedor DNS para `dns-01`                                                   |
| `--env`             | Líneas `KEY=VALUE` adicionales para `config.env` — para las credenciales DNS  |
| `--staging`         | Directorio staging de Let's Encrypt                                           |
| `--non-interactive` | No hacer preguntas; fallar con un error cuando falten parámetros              |

Desactivación:

```bash
gameapctl panel letsencrypt disable
```

El comando elimina las variables `ACME_*` de `config.env` y reinicia el panel. La opción
`--purge-certs` está declarada pero aún no implementada — los certificados emitidos permanecen en
el disco.

### Depuración de la emisión

El directorio de producción de Let's Encrypt tiene límites estrictos sobre el número de intentos
por dominio, y es fácil agotarlos mientras se configura. Hasta que su configuración funcione, use
el directorio staging:

```dotenv
ACME_DIRECTORY_URL=https://acme-staging-v02.api.letsencrypt.org/directory
```

Los navegadores tratan sus certificados como no confiables, pero los límites son mucho más
flexibles. Una vez que la emisión tenga éxito, elimine esta variable, borre el contenido del
directorio `ACME_STORAGE_PATH` y reinicie el panel para obtener un certificado de producción.

### Renovación

El panel comprueba la fecha de caducidad cada `ACME_RENEWAL_CHECK_INTERVAL` (cada 12 horas por
defecto) y renueva el certificado cuando queda menos de `ACME_RENEWAL_THRESHOLD` para la
caducidad (30 días por defecto). No se necesita un planificador ni una tarea de cron aparte.

Los certificados, la clave de la cuenta ACME y los datos auxiliares se almacenan en el directorio
`ACME_STORAGE_PATH` dentro del almacenamiento de archivos del panel. Con `FILES_DRIVER=s3`
terminan en S3 — esto es lo que permite que varias instancias del panel compartan un mismo
certificado.

El almacenamiento compartido por sí solo no basta para varias instancias: el bloqueo que impide
que soliciten un certificado al mismo tiempo funciona a través de Redis y está activo solo con
`CACHE_DRIVER=redis`. Con la caché en memoria el bloqueo es local a cada instancia, y estas se
estorbarán entre sí. Consulte [Varias instancias del panel](/es/multi_instance.html).

## Estado del certificado

El estado actual está disponible para un administrador en `GET /api/admin/letsencrypt/status`:

```json
{
  "enabled": true,
  "state": "active",
  "challenge_type": "http-01",
  "domains": ["panel.example.com"],
  "not_after": "2026-10-30T12:00:00Z",
  "last_renewal_at": "2026-08-01T12:00:00Z",
  "next_renewal_check_at": "2026-08-02T00:00:00Z"
}
```

Valores posibles de `state`:

| Valor      | Significado                                           |
|------------|-------------------------------------------------------|
| `disabled` | ACME está desactivado                                 |
| `pending`  | El certificado aún no ha sido emitido                 |
| `active`   | El certificado está emitido y es válido               |
| `renewing` | La renovación está en curso                           |
| `failed`   | El último intento falló; la razón está en `last_error` |

## Redirección a HTTPS

```dotenv
TLS_FORCE_HTTPS=true
```

Todas las peticiones HTTP reciben una redirección `301`, excepto `/.well-known/acme-challenge/` —
de lo contrario, el desafío `http-01` dejaría de funcionar.

Cuando ACME está activado, una petición cuya cabecera `Host` no coincide con ninguna entrada de
`ACME_DOMAINS` se redirige al primer dominio de la lista que no sea wildcard, así que ponga el
dominio principal del panel primero. Las entradas wildcard (`*.example.com`) nunca se usan como
destino de la redirección; si la lista solo las contiene, se conserva el host solicitado.

La misma variable afecta a otros dos mecanismos: la cabecera HSTS empieza a enviarse incluso
cuando TLS termina en un proxy inverso, y el origen CORS se calcula con el esquema `https`.

## Panel detrás de un proxy inverso

Si TLS termina en nginx, Traefik u otro proxy, no hace falta configurar certificados en el panel —
deje `ACME_ENABLED=false` y no defina `TLS_*`. El panel servirá HTTP en `8025`, y el proxy se
encargará de HTTPS.

Lo que importa en esta configuración:

* El proxy debe pasar la cabecera `X-Forwarded-Proto: https`, de lo contrario el panel no sabrá
  que la conexión es segura y no enviará HSTS.
* El proxy debe **sobrescribir** la cabecera `X-Forwarded-Proto` y la cabecera de
  `AUDIT_CLIENT_IP_HEADER`, no añadir a ellas: el panel confía en ellas sin verificar el remitente.
* El puerto **31718** normalmente no pasa por el proxy — los daemons deben conectarse al panel
  directamente. Defina `GRPC_EXTERNAL_HOST` con la dirección en la que el panel es accesible para
  los daemons.
* Si la dirección pública difiere de `HTTP_HOST`, inclúyala en `HTTP_ALLOWED_ORIGINS`.

## Problemas comunes

| Síntoma                                              | Causa                                                                                      |
|------------------------------------------------------|---------------------------------------------------------------------------------------------|
| El panel arrancó, pero HTTPS no está escuchando      | Una pareja de variables no está completa, o falta una de las variables `ACME_*` requeridas |
| `state: failed` con `http-01`                        | El puerto 80 no es accesible desde fuera, el dominio no resuelve a este servidor, u otro servicio lo intercepta |
| `state: failed` con `dns-01`                         | El token DNS no tiene permisos, o el registro no se propagó a tiempo — aumente `ACME_PROPAGATION_TIMEOUT` |
| La emisión dejó de funcionar después de varios intentos | Se agotó el límite de peticiones del directorio de producción de Let's Encrypt. Cambie al directorio staging y termine la configuración allí |
| El navegador se queja de la cadena                   | `TLS_CERT_FILE` contiene solo el certificado, sin los intermedios                           |
| El certificado fue reemplazado, pero se sirve el antiguo | Los archivos se leen al arrancar — hace falta `gameapctl panel restart`                |
