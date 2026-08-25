---
title: Múltiples instancias del panel
layout: default
lang: es
category: Administración
order: 338
---

El panel puede ejecutarse como varias instancias detrás de un balanceador de carga — para
tolerancia a fallos o para repartir la carga. Para ello, las instancias deben compartir el
estado: la base de datos, la caché, los archivos y el intercambio de eventos.

> La mayoría de las instalaciones no necesitan esto. Un solo panel gestiona cientos de
> servidores de juego — el trabajo pesado lo hacen los daemons en los servidores dedicados,
> no el panel.

## Qué debe ser compartido

| Componente        | Configuración                              | Qué ocurre en caso contrario                                  |
|-------------------|--------------------------------------------|--------------------------------------------------------------|
| Base de datos     | PostgreSQL o MySQL                         | SQLite no está diseñado para varias instancias conectadas    |
| Caché             | `CACHE_DRIVER=redis`                       | Con `memory`, las sesiones y las claves de configuración solo son visibles para una instancia |
| Intercambio de eventos | `PUBSUB_DRIVER=redis` o `postgres`    | Con `memory`, las instancias nunca se enteran de los eventos de las demás |
| Archivos          | `FILES_DRIVER=s3`                          | Con `local`, cada instancia tiene sus propios archivos y sus propios certificados |

Ejemplo de configuración:

```dotenv
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://gameap:password@db.example.com:5432/gameap?sslmode=require

CACHE_DRIVER=redis
CACHE_REDIS_ADDR=redis.example.com:6379

PUBSUB_DRIVER=redis
PUBSUB_REDIS_ADDR=redis.example.com:6379
PUBSUB_INSTANCE_ID=panel-1

FILES_DRIVER=s3
FILES_S3_ENDPOINT=https://s3.example.com
FILES_S3_BUCKET=gameap
FILES_S3_ACCESS_KEY_ID=...
FILES_S3_SECRET_ACCESS_KEY=...
```

`AUTH_SECRET` y `ENCRYPTION_KEY` deben ser **idénticos** en todas las instancias: de lo
contrario, un token emitido por una no será aceptado por otra, y los datos cifrados serán
ilegibles.

## Identificador de instancia

```dotenv
PUBSUB_INSTANCE_ID=panel-1
```

El valor debe ser **único para cada instancia**. Las instancias lo utilizan para distinguir
sus propios mensajes de los ajenos y para dirigirse respuestas entre sí.

Si la variable no está definida, se utiliza el valor `default` — el mismo en cada instancia, y
el intercambio entre instancias no funcionará. Defínala explícitamente.

## Cómo funciona

Un daemon se conecta a una de las instancias — aquella a la que el balanceador de carga lo
haya dirigido — y mantiene una conexión persistente con ella. Ninguna otra instancia puede
alcanzar ese daemon directamente.

Así, un comando que llega a una instancia que no posee la conexión con el daemon necesario se
reenvía a través del intercambio de eventos compartido a la instancia que sí la posee; la
respuesta regresa por el mismo camino. Las operaciones con archivos, la consola y las
peticiones de métricas funcionan del mismo modo.

De esto se deduce que el intercambio de eventos compartido es obligatorio: sin él, el panel
solo controlará los daemons que hayan conectado con esa instancia en particular.

Las tareas que aparecieron en la base de datos mientras un daemon estaba desconectado, o
mientras su conexión pertenecía a otra instancia, se entregan inmediatamente al conectarse, en
lugar de esperar a la próxima reconexión.

## Balanceador de carga y puertos

Separe los dos flujos:

* **HTTP (8025)** — puede balancearse de la forma habitual; la afinidad de sesión no es
  necesaria: los tokens se verifican con la clave compartida, y el estado de la sesión reside
  en la caché compartida.
* **gRPC (31718)** — las conexiones son de larga duración. El balanceador de carga debe ser
  capaz de mantenerlas y no cortarlas por un tiempo de inactividad: el canal no tiene un
  mecanismo keepalive propio, y la disponibilidad se apoya en un heartbeat cada 30 segundos.

Establezca `GRPC_EXTERNAL_HOST` a la dirección en la que los daemons ven el panel a través del
balanceador, y `GRPC_EXTERNAL_PORT` si el puerto se publica externamente con un número
diferente.

## Certificados

**Los certificados gRPC deben ser compartidos.** Se almacenan en el almacenamiento de archivos
del panel, por lo que con `FILES_DRIVER=s3` son compartidos automáticamente por todas las
instancias.

Con `FILES_DRIVER=local`, cada instancia generará su propia autoridad de certificación, y un
daemon registrado a través de una instancia no podrá conectarse a otra.

### HTTPS

Si el certificado del panel se emite mediante ACME, hay dos cosas que debe saber.

El almacenamiento ACME también reside en el almacenamiento de archivos — con S3 es compartido.

El bloqueo que impide que varias instancias soliciten un certificado al mismo tiempo funciona
a través de Redis. Está activo solo con `CACHE_DRIVER=redis`; con la caché en memoria el
bloqueo es local, y las instancias se estorbarán entre sí.

El desafío `http-01` requiere que la petición de la autoridad de certificación llegue a la
misma instancia que inició la emisión. Es más sencillo utilizar `dns-01` o terminar TLS en el
balanceador. Consulte [HTTPS y certificados](/es/https.html).

## Limitaciones

**Los contadores de intentos de inicio de sesión** residen en la caché. Con
`CACHE_DRIVER=redis` son compartidos; con `memory` cada instancia tiene los suyos, por lo que
el límite efectivo se multiplica por el número de instancias.

**El registro de auditoría** se escribe en la salida de cada instancia por separado. Para
agregarlos se necesita un recolector de registros externo, consulte [Seguridad](/es/security.html).

**Las tareas de los plugins** se desduplican entre instancias mediante bloqueos distribuidos,
de modo que la tarea periódica de un plugin se ejecuta una sola vez, y no en cada instancia.

## Comprobación

Tras el arranque, asegúrese de que:

* un daemon registrado a través de una instancia puede ser controlado a través de otra;
* un inicio de sesión realizado en una instancia es válido en las demás;
* el gestor de archivos se abre independientemente de la instancia a la que haya llegado la
  petición.

Si algo de esto no funciona, la causa es casi siempre que uno de los cuatro componentes
compartidos ha permanecido local.
