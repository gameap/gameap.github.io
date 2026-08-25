---
title: API gRPC
layout: default
lang: es
category: GameAP Daemon
order: 405
---

A partir de GameAP 4.2 y GameAP Daemon 4.0, el panel y el daemon intercambian datos a través de gRPC
mediante un flujo bidireccional (bidirectional streaming). Este método ha sustituido al antiguo
intercambio a través de BINN y la API REST.

La conexión la establece el daemon: se conecta al panel por sí mismo y mantiene un único flujo de
larga duración que lo transporta todo: el registro, el heartbeat, las métricas, las tareas, los
comandos, la consola y las operaciones con archivos. El panel no se conecta al daemon y no requiere
puertos entrantes en el servidor dedicado.

## Configuración del panel

No existe un ajuste separado para activar gRPC: el servidor siempre se inicia. Solo se pueden
configurar la dirección, el cifrado y los límites.

| Variable                      | Predeterminado | Propósito                                                                 |
|-------------------------------|----------------|---------------------------------------------------------------------------|
| `GRPC_PORT`                   | `31718`        | Puerto en el que escucha el servidor gRPC del panel                       |
| `GRPC_TLS_ENABLED`            | `true`         | Cifrado de la conexión                                                    |
| `GRPC_REQUIRE_MTLS`           | `false`        | Exigir un certificado de cliente al daemon                                |
| `GRPC_EXTERNAL_HOST`          | `""`           | Dirección del panel que se comunica al daemon. Vacío: se determina a partir de la petición |
| `GRPC_EXTERNAL_PORT`          | `0`            | Puerto que se comunica al daemon. `0`: se usa `GRPC_PORT`                 |
| `GRPC_MAX_RECV_MSG_SIZE`      | `10485760`     | Tamaño máximo de un mensaje entrante, en bytes                            |
| `GRPC_MAX_SEND_MSG_SIZE`      | `10485760`     | Tamaño máximo de un mensaje saliente, en bytes                            |
| `GRPC_MAX_CONCURRENT_STREAMS` | `100`          | Número de flujos simultáneos por conexión                                 |
| `GRPC_ENABLE_REFLECTION`      | `false`        | Reflexión del esquema para herramientas de depuración como `grpcurl`. No lo active en producción |

> La variable `GRPC_ENABLED` **no existe**: el panel no la lee. Las versiones antiguas de
> `gameapctl` añaden una línea `GRPC_ENABLED=true` a `config.env`: es inofensiva, pero no tiene
> ningún efecto. El servidor gRPC no se puede desactivar.

### Puertos

gRPC funciona en un **puerto separado, 31718**, mientras que la interfaz web y la API funcionan en
`HTTP_PORT` (`8025` por defecto). Son dos listeners diferentes, no un único puerto con detección de
protocolo.

El puerto 31718 debe ser accesible desde cada servidor dedicado. Es el único puerto del panel que
necesita un daemon en funcionamiento.

> En el `Dockerfile` y el `docker-compose.yml` del panel solo se publica el puerto 8025. Al
> desplegar con Docker, debe redirigir el puerto 31718 usted mismo.

### Dirección del panel para el daemon

El panel sustituye su propia dirección en el comando de instalación del daemon y en una URL de
conexión de la forma `grpc://host:port/key`. Por defecto, el host se toma de la cabecera de la
petición con la que el administrador abrió la página de creación del servidor dedicado, es decir,
de la dirección en la barra de direcciones del navegador.

Establezca `GRPC_EXTERNAL_HOST` si esa dirección difiere de aquella a la que deben conectarse los
daemons:

* el panel está detrás de un proxy inverso, gRPC no pasa a través del proxy y los daemons tienen
  que conectarse directamente;
* el panel está detrás de NAT y tiene direcciones diferentes dentro y fuera;
* el panel se ejecuta en Docker, donde la cabecera acaba conteniendo `localhost` o el nombre del
  contenedor.

`GRPC_EXTERNAL_PORT` es necesario cuando el puerto 31718 se publica externamente con un número
diferente.

Si las variables no están establecidas, el propio servidor gRPC funciona correctamente; solo la
dirección del comando de instalación generado resulta incorrecta, y el daemon no podrá conectarse.

> `GRPC_EXTERNAL_HOST` se incluye en la lista de nombres alternativos del sujeto (SAN) del
> certificado gRPC autofirmado. El certificado se crea una sola vez, en el primer arranque, por lo
> que la variable debe establecerse **antes** de iniciar el panel por primera vez. Si la establece
> más tarde, el daemon rechazará la conexión debido a la discrepancia del nombre en el certificado:
> elimine `certs/server/api-server.crt` y `certs/server/api-server.key` y reinicie el panel para
> que el certificado se genere de nuevo.

### Cifrado y certificados

Con `GRPC_TLS_ENABLED=true` (el valor por defecto) el panel utiliza un certificado autofirmado
emitido por su propia autoridad certificadora interna: `certs/root.crt` y `certs/root.key`. El
certificado del servidor es `certs/server/api-server.crt`. Las claves son RSA de 2048 bits, válidas
durante 10 años, y todo se crea automáticamente en el primer uso.

Estos certificados no tienen nada que ver con el certificado HTTPS propio del panel: ACME y Let's
Encrypt no se aplican a gRPC, y no es necesario configurarlos por separado.

Durante el registro, el daemon recibe `ca.crt`, `server.crt` y `server.key` del panel y los coloca
en su directorio de certificados. A partir de entonces verifica el certificado del panel contra la
CA recibida y presenta su propio certificado de cliente.

Además, cada petición se autentica con la clave API del nodo, que se almacena en la base de datos
del panel como un hash SHA-256 y se compara en tiempo constante.

### Autenticación mutua (mTLS)

`GRPC_REQUIRE_MTLS=true` hace que el panel exija al cliente un certificado emitido por su propia
autoridad certificadora y rechace las peticiones sin él.

No es necesario configurar el daemon de forma especial para esto: un daemon registrado siempre
presenta su certificado de todos modos.

> **Active mTLS solo después de que todos los daemons hayan sido registrados.** El registro en sí
> ocurre sin un certificado de cliente: un daemon nuevo aún no lo tiene. Con
> `GRPC_REQUIRE_MTLS=true` no podrá registrar un nuevo servidor dedicado. Para añadir un nodo más
> tarde, vuelva a establecerlo temporalmente en `false`, registre el daemon y actívelo de nuevo.

Los daemons registrados por otro panel no funcionarán: sus certificados están emitidos por una
autoridad certificadora ajena.

## Configuración del daemon

Los parámetros de conexión se establecen en la configuración del daemon —
`/etc/gameap-daemon/gameap-daemon.yaml` en Linux, `C:\gameap\daemon\gameap-daemon.yaml` en
Windows — en el bloque `grpc`:

```yaml
grpc:
  address: panel.example.com:31718
  insecure: false
  heartbeat_interval: 30s
  connect_timeout: 30s
  initial_reconnect_delay: 1s
  max_reconnect_delay: 60s
```

| Parámetro                 | Predeterminado | Propósito                                                       |
|---------------------------|----------------|------------------------------------------------------------------|
| `address`                 | —              | Dirección del panel en formato `host:port`                       |
| `insecure`                | `false`        | Desactivar TLS. Solo para depuración                             |
| `heartbeat_interval`      | `30s`          | Intervalo de heartbeat. El panel puede asignar su propio valor   |
| `connect_timeout`         | `30s`          | Tiempo de espera para establecer la conexión                     |
| `initial_reconnect_delay` | `1s`           | Pausa inicial antes de la reconexión                             |
| `max_reconnect_delay`     | `60s`          | Pausa máxima antes de la reconexión                              |

Si `address` no está establecido, se deriva del parámetro obsoleto `api_host`: se toma el nombre
del host y el puerto se sustituye por 31718. Establecer `address` explícitamente es más fiable.

> El daemon tampoco tiene la clave `grpc.enabled`. `gameapctl` la escribe durante la migración como
> un marcador de que la migración se ha realizado; el daemon ignora esta clave.

Los demás parámetros de configuración del daemon se describen en la página
[GameAP Daemon](/es/daemon/daemon.html).

### Reconexión y pérdida de conexión

Si se pierde la conexión, el daemon se reconecta por sí mismo, con un retardo exponencial: se
duplica desde `initial_reconnect_delay` hasta `max_reconnect_delay` y se dispersa con un jitter
aleatorio de ±10 %, para que muchos daemons no lleguen a la vez. Con los valores por defecto es
1 s, 2 s, 4 s, 8 s y así sucesivamente hasta 60 s. Tras una conexión exitosa, el contador se
reinicia. En un apagado planificado, el panel puede asignar al daemon una pausa antes de su próximo
intento.

Mientras el panel no está disponible, **los servidores de juego siguen funcionando**: solo se
interrumpe su gestión desde el panel. Una vez restaurada la conexión, el daemon se registra de
nuevo, informa de las tareas que está ejecutando y recibe el estado actual completo del panel: la
lista de servidores, tareas, juegos y modificaciones, y la configuración de los servidores. Por eso
las tareas iniciadas antes de la desconexión no se pierden.

> El mecanismo keepalive de gRPC no se utiliza en ninguno de los dos lados: la vitalidad de la
> conexión depende únicamente del heartbeat cada 30 segundos. Si hay NAT o un firewall entre el
> daemon y el panel que cierra las conexiones inactivas antes, reduzca `heartbeat_interval`.

## Migración desde el protocolo antiguo

Un daemon instalado antes de la aparición de gRPC se cambia al nuevo protocolo con el comando:

```bash
gameapctl daemon upgrade --switch-to-grpc
```

Qué hace el comando:

1. Comprueba que la migración aún no se ha realizado y determina la dirección del panel a partir de
   `api_host` (o la toma de `--grpc-address`).
2. Comprueba que la configuración contiene `api_key`, `ds_id` y los tres archivos de certificados.
3. **Antes de hacer cualquier cambio**, comprueba que el panel es accesible: establece una conexión
   TCP y realiza un handshake TLS real con los certificados existentes. Si los certificados fueron
   emitidos por otro panel, el comando lo informará y sugerirá una reinstalación.
4. Crea una copia de seguridad de la configuración junto al original, con una marca de tiempo en el
   nombre.
5. Escribe la dirección gRPC y elimina los obsoletos `api_host`, `listen_ip` y `listen_port`.
6. Reinicia el daemon y se asegura de que el panel ha revocado el acceso a través de la antigua API
   HTTP.
7. Ante cualquier fallo después del paso 5, restaura la configuración desde la copia de seguridad y
   arranca el daemon de nuevo.

Todo lo que se requiere es que el puerto 31718 del panel sea accesible desde el servidor dedicado.
No hay que activar nada en el lado del panel; contrariamente a lo que dicen el propio mensaje de
error del comando y su descripción, la variable `GRPC_ENABLED` no existe.

## Verificación

El panel responde al health check estándar de gRPC e informa del estado `SERVING` para los
servicios `gameap.DaemonGateway` y `gameap.FileTransferService`.

La comprobación más sencilla de que el puerto es accesible desde el servidor dedicado:

```bash
nc -zv panel.example.com 31718
```

El estado de conexión del daemon es visible en el panel en la página **«Administración»** →
**«Servidores dedicados»**. Los detalles están en el log del daemon:
`/var/log/gameap-daemon/output.log` en Linux, `C:\gameap\daemon\logs\output.log` en Windows.

Mensajes típicos en el log del daemon:

| Mensaje                        | Qué significa                                                                 |
|--------------------------------|-------------------------------------------------------------------------------|
| `gRPC connection failed`       | La conexión no se estableció o se interrumpió; le siguen una pausa y un nuevo intento |
| `registration failed: ...`     | Hay conexión, pero el panel rechazó el registro: `ds_id` o `api_key` incorrectos |
| Error de verificación del certificado | El nombre del certificado del panel no coincide con la dirección de conexión. Véase `GRPC_EXTERNAL_HOST` |
