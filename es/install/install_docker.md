---
title: Instalación en Docker
layout: default
lang: es
category: Instalación de GameAP
order: 102
---

El panel se distribuye como una imagen lista para usar `gameap/gameap`. El daemon no se ejecuta en Docker:
debe gestionar los procesos de los servidores de juego directamente en el servidor dedicado.

## Inicio rápido

La opción más sencilla es SQLite, sin una base de datos separada:

```bash
docker run -d \
  --name gameap \
  -p 8025:8025 \
  -p 31718:31718 \
  -e DATABASE_DRIVER=sqlite \
  -e 'DATABASE_URL=file:/var/lib/gameap/db.sqlite?_busy_timeout=5000&_journal_mode=WAL&cache=shared' \
  -e AUTH_SECRET=$(openssl rand -hex 16) \
  -e ENCRYPTION_KEY=$(openssl rand -hex 16) \
  -e GRPC_EXTERNAL_HOST=panel.example.com \
  -v gameap-data:/var/lib/gameap \
  gameap/gameap:latest
```

El panel estará disponible en `http://localhost:8025`. El usuario y la contraseña del primer
administrador aparecerán en el registro del contenedor:

```bash
docker logs gameap
```

> Dos parámetros del ejemplo son esenciales, y ambos son fáciles de pasar por alto porque los
> archivos `docker-compose.yml` listos para usar no los incluyen. Consulte [Puerto 31718](#puerto-31718) y
> [Dirección del panel para los daemons](#dirección-del-panel-para-los-daemons).

## Puerto 31718

Los daemons se conectan al panel a través de este puerto. El `Dockerfile` solo declara `8025`, y el
`docker-compose.yml` de ejemplo del repositorio también publica únicamente ese.

**Debe publicar el puerto 31718 usted mismo** — de lo contrario el panel funcionará, pero ningún
servidor dedicado podrá conectarse a él.

```yaml
ports:
  - "8025:8025"
  - "31718:31718"
```

## Dirección del panel para los daemons

El panel sustituye su propia dirección en el comando de instalación del daemon. Por defecto la toma
de la cabecera de la solicitud, y dentro de un contenedor suele ser `localhost` o el nombre del
servicio — un comando así no funcionará en un servidor dedicado remoto.

Configure la dirección en la que el panel es accesible desde el exterior:

```yaml
environment:
  GRPC_EXTERNAL_HOST: panel.example.com
```

Si el puerto 31718 se publica externamente con un número diferente, configúrelo también:

```yaml
environment:
  GRPC_EXTERNAL_PORT: 41718
```

> `GRPC_EXTERNAL_HOST` se incluye en la lista de nombres del certificado gRPC autofirmado, y el
> certificado se crea en el primer arranque. Configure la variable **antes** del primer arranque. Si
> la configura más tarde, elimine `certs/server/api-server.crt` y `certs/server/api-server.key` en el
> directorio de archivos del panel y reinicie el contenedor. Consulte [GRPC API](/es/daemon/grpc.html)
> para más detalles.

## Docker Compose

Hay un `docker-compose.yml` listo para usar disponible en el
[repositorio del panel](https://github.com/gameap/gameap). Levanta el panel junto con
PostgreSQL y Redis:

```bash
git clone https://github.com/gameap/gameap.git
cd gameap
docker compose up -d
```

Debe añadirle la publicación del puerto 31718 y `GRPC_EXTERNAL_HOST`, como se describe arriba.

Las contraseñas y las claves se configuran mediante un archivo `.env` junto a `docker-compose.yml`:

```dotenv
AUTH_SECRET=replace_with_32_random_bytes
ENCRYPTION_KEY=replace_with_32_random_bytes
POSTGRES_PASSWORD=replace_with_a_strong_password
LOGGER_LEVEL=info
```

> En el `docker-compose.yml` de ejemplo, los valores predeterminados de `AUTH_SECRET` y
> `ENCRYPTION_KEY` son `change-me-in-production`. Con tales claves, los tokens de sesión son
> triviales de falsificar. Configure sus propios valores antes del primer arranque.
>
> Tenga en cuenta: cambiar `ENCRYPTION_KEY` en una instalación en funcionamiento rompe la
> autenticación de dos factores para todos los usuarios. Consulte [Seguridad](/es/security.html).

## Datos y volúmenes

Todo el estado del panel se almacena en `/var/lib/gameap`:

| Qué                  | Dónde                                                      |
|----------------------|------------------------------------------------------------|
| Archivos del panel   | El directorio de `FILES_LOCAL_BASE_PATH`, `/var/lib/gameap/files` en el ejemplo de compose |
| Certificados gRPC    | El subdirectorio `certs/` dentro de él                     |
| Certificados ACME    | El subdirectorio de `ACME_STORAGE_PATH`, `acme/` por defecto |
| Base de datos SQLite | Donde apunte `DATABASE_URL`                                |

Los certificados y los datos de ACME no se almacenan en la raíz del volumen, sino **dentro del
directorio de archivos del panel**. Con la configuración del `docker-compose.yml` de ejemplo, la ruta
completa al certificado gRPC es `/var/lib/gameap/files/certs/server/api-server.crt`.

Monte este directorio como un volumen; de lo contrario, al recrear el contenedor se perderán los
certificados — y todos los daemons dejarán de conectarse.

El contenedor se ejecuta como el usuario sin privilegios `gameap`, por lo que el volumen debe tener
los permisos adecuados.

## Variables de entorno

La imagen se configura con las mismas variables que una instalación normal — la lista completa está
en la [referencia de config.env](/es/config.html). No se necesita ningún archivo `config.env` dentro
del contenedor.

El mínimo indispensable:

| Variable             | Propósito                                                  |
|----------------------|------------------------------------------------------------|
| `DATABASE_DRIVER`    | `postgres`, `mysql` o `sqlite`                             |
| `DATABASE_URL`       | Cadena de conexión a la base de datos                      |
| `AUTH_SECRET`        | Clave de firma de tokens, 32 bytes aleatorios              |
| `ENCRYPTION_KEY`     | Clave de cifrado de secretos, 32 bytes aleatorios          |
| `GRPC_EXTERNAL_HOST` | Dirección del panel que usan los daemons para conectarse   |

## Comprobación de estado

La imagen incluye una comprobación de estado — una solicitud a `/api/health`:

```bash
docker inspect --format='{{.State.Health.Status}}' gameap
```

El registro del panel:

```bash
docker logs -f gameap
```

## Detrás de un proxy inverso

Si el HTTPS lo sirve nginx o Traefik, no es necesario configurar certificados en el
panel. El proxy debe enviar la cabecera `X-Forwarded-Proto: https` y **sobrescribirla**, no
añadirla.

El puerto 31718 normalmente no pasa por un proxy inverso: los daemons deben conectarse al panel
directamente. Consulte [HTTPS y certificados](/es/https.html) para más detalles.

## Actualización

```bash
docker compose pull
docker compose up -d
```

Antes de actualizar, haga una copia de seguridad de la base de datos y del volumen de datos: las
migraciones se aplican al arrancar y son irreversibles.

```bash
docker compose exec postgres pg_dump -U gameap gameap > gameap-backup.sql
```
