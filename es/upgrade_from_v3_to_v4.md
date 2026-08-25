---
title: Actualización de v3 a v4
layout: default
lang: es
category: Instalación de GameAP
order: 191
---

La actualización desde GameAP 3 conservando la base de datos existente solo es posible a
las **versiones 4.0 y 4.1**. Estas reconocen una base de datos de la tercera versión,
amplían su esquema con sus propias tablas y siguen trabajando con los mismos usuarios,
servidores y juegos.

> **A partir de la versión 4.2, no se admite la actualización in situ.** Las versiones
> posteriores esperan el esquema de la cuarta versión y no funcionarán con una base de
> datos de GameAP 3.
>
> Si necesita la última versión del panel, instálela sobre una base de datos limpia y
> migre los datos usted mismo: consulte
> [Si necesita la última versión](#if-you-need-the-latest-version).

La última versión de la línea compatible es **4.1.2**.

## Antes de actualizar

> **Haga una copia de seguridad de la base de datos manualmente.** `gameapctl` copia el
> directorio de instalación de v3 y la configuración del servidor web, pero **no guarda
> la base de datos**, y la actualización cambia su esquema de forma irreversible.

```shell
# MySQL o MariaDB
mysqldump -u root -p gameap > gameap-v3-backup.sql

# PostgreSQL
pg_dump -U gameap gameap > gameap-v3-backup.sql

# SQLite: basta con copiar el archivo
cp /var/www/gameap/database.sqlite gameap-v3-backup.sqlite
```

Compruebe que la copia de seguridad no esté vacía y solo entonces continúe.

### Qué bases de datos se pueden actualizar in situ

| SGBD en GameAP 3 | Actualización in situ |
|------------------|-----------------------|
| MySQL, MariaDB   | sí                    |
| SQLite           | sí                    |
| PostgreSQL       | **no**                |

Para MySQL y SQLite, el panel reconoce una base de datos existente de GameAP 3 y no
intenta crear las tablas de nuevo. Para PostgreSQL no existe tal reconocimiento: la
primera migración intentará crear tablas que ya existen y fallará.

Si GameAP 3 funciona con PostgreSQL, instale GameAP 4 sobre una base de datos limpia y
migre los datos usted mismo.

## Actualización a 4.1

`gameapctl` se instala junto con el panel y puede realizar la transición.

```shell
gameapctl self-update
gameapctl panel upgrade
```

La utilidad detectará que está instalada la tercera versión y realizará la actualización.

> **Compruebe qué versión instala `gameapctl`.** Al pasar desde la tercera versión,
> ignora el flag `--version` e instala la última versión estable de la línea 4.x. Si
> resulta ser 4.2 o más reciente, la base de datos de producción no debe actualizarse.
>
> Pruebe primero la transición en una copia de la base de datos: el procedimiento se
> describe en [Prueba en una copia de la base de datos](#testing-on-a-database-copy).
> Asegúrese de que se instaló 4.0 o 4.1 y solo entonces actualice la instalación de
> producción.

Qué ocurre durante la actualización:

1. Se lee el archivo `.env` de la instalación de GameAP 3 y se toman de él los
   parámetros de conexión a la base de datos.
2. El directorio de instalación de v3 y la configuración del servidor web se copian a un
   directorio temporal. La ruta de la copia se muestra en la consola: anótela.
3. Se crea la configuración de GameAP 4 `config.env`, con la misma base de datos
   conectada y nuevas claves `AUTH_SECRET` y `ENCRYPTION_KEY` generadas.
4. GameAP 4 se instala y se inicia.

En el primer inicio, el panel aplica las migraciones: detecta que la base de datos
pertenece a la tercera versión, omite la creación de tablas y añade solo lo que le falta
a la cuarta versión.

### gameapctl no está instalado

Descárguelo desde la [página de releases](https://github.com/gameap/gameapctl/releases)
o con un comando que sustituye la versión correcta:

```shell
VERSION=$(curl -sL https://api.github.com/repos/gameap/gameapctl/releases/latest \
  | grep -m1 '"tag_name"' | cut -d'"' -f4)
curl -OL "https://github.com/gameap/gameapctl/releases/download/${VERSION}/gameapctl-${VERSION}-linux-amd64.tar.gz"
tar xvfz "gameapctl-${VERSION}-linux-amd64.tar.gz" -C /usr/local/bin
```

## Qué cambia con la transición

**Ya no se necesitan un servidor web ni PHP.** GameAP 4 es un único ejecutable con una
interfaz web integrada. Tras la actualización, nginx o Apache pueden conservarse como
proxy inverso o eliminarse por completo.

**El puerto predeterminado es 8025**, no 80.

**Las contraseñas de los usuarios siguen funcionando.** Las versiones 4.0 y 4.1
verifican las contraseñas del mismo modo que la tercera versión. No es necesario pedir
a los usuarios que cambien sus contraseñas.

**Los daemons siguen funcionando con el protocolo antiguo.** En las versiones 4.0 y 4.1
el panel intercambia datos con GameAP Daemon del mismo modo que la tercera versión y se
conecta al propio daemon. La configuración de los servidores dedicados no requiere
cambios, y el puerto del daemon (`31717` por defecto) debe permanecer abierto.

> No ejecute `gameapctl daemon upgrade --switch-to-grpc` después de actualizar a 4.1.
> El intercambio de datos por gRPC solo apareció en la versión 4.2, y un panel 4.1 no
> podrá conectarse a un daemon de ese tipo.

## Reversión

Si algo no funciona después de la actualización:

1. Detenga GameAP 4: `gameapctl panel stop`.
2. Restaure la base de datos desde la copia de seguridad realizada antes de la
   actualización.
3. Recupere el directorio de instalación de v3 desde la copia temporal cuya ruta mostró
   `gameapctl`.
4. Restaure la configuración del servidor web e inícielo.

> Una reversión sin copia de seguridad de la base de datos es imposible: el esquema ha
> cambiado y GameAP 3 ya no funcionará con ella.

## Prueba en una copia de la base de datos {#testing-on-a-database-copy}

Antes de actualizar la instalación de producción, conviene ejecutar la transición en una
copia de la base de datos. Esto también muestra qué versión exacta instalará
`gameapctl`.

> No apunte la instalación de prueba a la base de datos de producción: en el primer
> inicio el panel aplicará migraciones sobre ella y GameAP 3 dejará de funcionar.

Cree una copia de la base de datos:

```shell
mysqldump -u root -p gameap > /tmp/gameap.sql
mysql -u root -p -e "CREATE DATABASE gameap_v4_test"
mysql -u root -p gameap_v4_test < /tmp/gameap.sql
```

Descargue una versión de la línea 4.1 desde la
[página de releases](https://github.com/gameap/gameap/releases) y desempaquétela:

```shell
curl -OL https://github.com/gameap/gameap/releases/download/v4.1.2/gameap-v4.1.2-linux-amd64.tar.gz
tar xvfz gameap-v4.1.2-linux-amd64.tar.gz -C /usr/bin
```

Cree el archivo de configuración `/etc/gameap/config.env`:

```dotenv
DATABASE_DRIVER=mysql
DATABASE_URL=gameap:password@tcp(127.0.0.1:3306)/gameap_v4_test?parseTime=true

AUTH_SECRET=replace_with_32_random_bytes
ENCRYPTION_KEY=replace_with_32_random_bytes

HTTP_PORT=8025
```

Las claves se generan fácilmente con `openssl rand -hex 16`.

Ejecute:

```shell
gameap --env /etc/gameap/config.env
```

El panel estará disponible en el puerto 8025. La lista completa de parámetros de
configuración está en la [referencia de config.env](/es/config.html).

Asegúrese de que los datos están intactos: los usuarios pueden iniciar sesión, y los
servidores de juego y los juegos se muestran. Después de eso puede actualizar la
instalación de producción.

### Configuración de un servicio systemd

Para una gestión cómoda de la instalación de prueba, cree un servicio aparte.

Cree un usuario y un directorio:

```shell
useradd -r -s /usr/sbin/nologin -d /var/lib/gameap gameap
mkdir -p /var/lib/gameap
chown gameap:gameap /var/lib/gameap
```

Luego el archivo `/etc/systemd/system/gameap.service`:

```ini
[Unit]
Description=GameAP - Game Server Control Panel
Documentation=https://docs.gameap.com
After=network.target
Wants=network-online.target
Requires=network.target

[Service]
Type=simple
User=gameap
Group=gameap

WorkingDirectory=/var/lib/gameap

ExecStart=/usr/bin/gameap

# Permitir el enlace a puertos privilegiados
AmbientCapabilities=CAP_NET_BIND_SERVICE

# Parada controlada
ExecStop=/bin/kill -TERM $MAINPID
KillMode=mixed
KillSignal=SIGTERM
TimeoutStopSec=30

# Política de reinicio
Restart=always
RestartSec=5
StartLimitInterval=60
StartLimitBurst=3

EnvironmentFile=/etc/gameap/config.env

RuntimeDirectory=gameap
PIDFile=/run/gameap/gameap.pid

# Permisos del sistema de archivos
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true

ReadWritePaths=/var/lib/gameap

# Registro
StandardOutput=journal
StandardError=journal
SyslogIdentifier=gameap

[Install]
WantedBy=multi-user.target
```

Habilite e inicie el servicio:

```shell
systemctl daemon-reload
systemctl enable gameap
systemctl start gameap
```

## Si necesita la última versión {#if-you-need-the-latest-version}

No se puede actualizar desde GameAP 3 directamente a 4.2 o más reciente. El
procedimiento es el siguiente:

1. Instale la versión actual de GameAP 4 sobre una base de datos **limpia**, dejando
   intacta la instalación de producción; consulte
   [Instalación en Linux](/es/install/install_on_linux.html).
2. Migre los datos de GameAP 3: vuelva a crear los usuarios, servidores dedicados,
   juegos y servidores de juego, manualmente o a través de la
   [API](https://openapi.gameap.io/).
3. Asegúrese de que todo funciona y solo entonces retire la instalación antigua.

Los archivos de los servidores de juego no necesitan moverse a ninguna parte:
permanecen en el servidor dedicado; basta con describir los servidores en el nuevo
panel con los mismos directorios y puertos.
