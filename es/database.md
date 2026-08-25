---
title: Base de datos
layout: default
lang: es
category: Administración
order: 335
---

El panel funciona con PostgreSQL, MySQL o MariaDB, y SQLite. La elección se define con dos
variables en `config.env`:

```dotenv
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://gameap:password@localhost:5432/gameap?sslmode=disable
```

## Cuál elegir

| SGBD               | Cuándo conviene                                                                                                                                            |
|--------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **SQLite**         | Un solo panel, hasta unas pocas docenas de servidores de juego. No requiere un servicio aparte                                                             |
| **PostgreSQL**     | Cuando se necesita alto rendimiento. Recomendado para ejecutar varias instancias del panel                                                                 |
| **MySQL, MariaDB** | La opción familiar, incluso al actualizar desde GameAP 3. Ligeramente más lento que PostgreSQL, pero aún adecuado para varias instancias                     |

Para una instalación típica en un solo servidor, SQLite es suficiente: no requiere un servicio
aparte, ni configuración, ni una copia de seguridad más compleja que copiar un archivo.

## Cadena de conexión

### PostgreSQL

```dotenv
DATABASE_DRIVER=postgres
DATABASE_URL=postgres://user:password@host:5432/database?sslmode=disable
```

Los nombres del driver son intercambiables: `postgres`, `postgresql`, `pgsql`, `pg`, `pgx`.

El parámetro `sslmode` define el modo de cifrado de la conexión: `disable` para una base de
datos local, `require` o `verify-full` para una remota.

### MySQL y MariaDB

```dotenv
DATABASE_DRIVER=mysql
DATABASE_URL=user:password@tcp(host:3306)/database?parseTime=true
```

El formato difiere de la URL habitual con esquema — es el formato del driver de Go. El
parámetro `parseTime=true` es obligatorio.

Conexión a través de un socket:

```dotenv
DATABASE_URL=user:password@unix(/var/run/mysqld/mysqld.sock)/database?parseTime=true
```

### SQLite

```dotenv
DATABASE_DRIVER=sqlite
DATABASE_URL=file:/var/lib/gameap/db.sqlite?_busy_timeout=5000&_journal_mode=WAL&cache=shared
```

El archivo se crea automáticamente. Los parámetros del ejemplo activan el journaling WAL y la
espera de bloqueos — consérvelos, mejoran notablemente el comportamiento con peticiones
concurrentes.

El directorio con el archivo de la base de datos debe tener permisos de escritura para el
usuario con el que se ejecuta el panel.

### inmemory

```dotenv
DATABASE_DRIVER=inmemory
```

Los datos se mantienen solo en la RAM y se pierden al reiniciar. Está pensado para pruebas; no
es adecuado para una instalación en producción.

## Migraciones

El panel aplica las migraciones por sí mismo al arrancar — no existe un comando aparte. La
versión del esquema se guarda en una tabla de servicio dentro de la misma base de datos.

De esto se derivan dos reglas prácticas:

* **Una actualización del panel cambia el esquema en el primer arranque.** Haga la copia de
  seguridad antes.
* **Volver a una versión anterior del panel sin restaurar la base de datos no funcionará** —
  el esquema ya ha cambiado.

Las migraciones se aplican incluso con huecos en la numeración, así que saltarse una versión
intermedia del panel durante una actualización no es problema.

## Copia de seguridad

El panel no hace copias de seguridad — ni por sí mismo ni mediante `gameapctl`. Configure las
copias de seguridad con las herramientas del SGBD.

### PostgreSQL

```bash
pg_dump -U gameap gameap > gameap-$(date +%F).sql
```

Restauración:

```bash
psql -U gameap gameap < gameap-2026-08-02.sql
```

### MySQL y MariaDB

```bash
mysqldump -u gameap -p gameap > gameap-$(date +%F).sql
```

Restauración:

```bash
mysql -u gameap -p gameap < gameap-2026-08-02.sql
```

### SQLite

No es necesario detener el panel si utiliza el comando incorporado:

```bash
sqlite3 /var/lib/gameap/db.sqlite ".backup '/backup/gameap-$(date +%F).sqlite'"
```

Copiar simplemente el archivo mientras el panel está en ejecución puede producir una copia
corrupta debido al journal WAL.

### Qué más guardar

La base de datos no es suficiente. Junto con ella, guarde:

* **`config.env`** — contiene `AUTH_SECRET` y `ENCRYPTION_KEY`. Sin `ENCRYPTION_KEY`,
  algunos datos de la copia de seguridad no podrán restaurarse, y la autenticación de dos
  factores dejará de funcionar para todos los usuarios;
* **el directorio de archivos del panel** — contiene los certificados gRPC con los que se
  conectan los daemons y los datos de ACME. La ruta se define con `FILES_LOCAL_BASE_PATH`.

Los archivos de los servidores de juego residen en los servidores dedicados y no forman parte
de la copia de seguridad del panel.

## Cambiar a otro SGBD

No existe una transferencia de datos incorporada entre SGBD: los esquemas de PostgreSQL, MySQL
y SQLite se crean de forma independiente, y el panel no proporciona ninguna herramienta de
migración.

El procedimiento:

1. Haga una copia de seguridad de la base de datos actual.
2. Prepare la nueva base de datos y apunte `DATABASE_DRIVER` y `DATABASE_URL` hacia ella.
3. Inicie el panel — creará el esquema desde cero.
4. Migre los datos: vuelva a crear los usuarios, servidores dedicados, juegos y servidores de
   juego manualmente o mediante la [API](/es/api.html).

Los servidores de juego en sí no se ven afectados: sus archivos permanecen en los servidores
dedicados; basta con describir los servidores en el panel con los mismos directorios y puertos.

Cargar un volcado de un SGBD directamente en otro no funcionará — los tipos de columna y la
sintaxis difieren.

## Varias instancias del panel

Para una instalación con varias instancias, utilice PostgreSQL o MySQL. SQLite no sirve: no
está diseñado para que varios clientes trabajen con el archivo de la base de datos al mismo
tiempo.

Además de la base de datos, necesitará una caché compartida, un intercambio de eventos
compartido y un almacenamiento de archivos compartido. Consulte
[Varias instancias del panel](/es/multi_instance.html) para más detalles.

## Comprobación

Si el panel no arranca por culpa de la base de datos, el log contendrá un mensaje de conexión
o de migración:

```bash
journalctl -u gameap -n 50 --no-pager
```

Causas habituales: un formato incorrecto de `DATABASE_URL` (especialmente en MySQL — necesita
el formato del driver de Go, no una URL con esquema), falta de permisos sobre el directorio del
archivo SQLite, la base de datos no creada o el usuario sin acceso concedido a ella.
