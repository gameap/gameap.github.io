---
title: Instalación en Linux
layout: default
lang: es
category: Instalación de GameAP
order: 100
---

## Instalación

La instalación en Linux se realiza con un único comando:
```shell
bash <(curl -s https://gameap.com/install.sh)
```

Durante el proceso de instalación, se le pedirá que introduzca cierta información.

### Host

Especifique el dominio o la dirección IP en la que el panel estará accesible.

En el caso de una dirección IP, debe ser una dirección asignada a la interfaz 
de red del VDS. Si su red utiliza NAT, no especifique 
la IP externa, sino la interna, y luego configure el reenvío 
de puertos.

Se puede especificar cualquier dominio, pero no olvide configurar el DNS.

Ejemplos de valores correctos:
* 10.182.104.8
* 10.182.104.8:2080
* example.com
* http://example.com

### Base de datos

La base de datos donde se almacenarán los datos: usuarios, información sobre los servidores, etc. 
Puede usar [PostgreSQL](https://www.postgresql.org/),
[MySQL](https://www.mysql.com/) y 
[SQLite](https://www.sqlite.org/).

PostgreSQL se recomienda en la mayoría de los casos. Si se espera que la carga de su servidor 
sea baja y no planea usar más de 10 servidores de juego, 
puede usar SQLite.

Algunas distribuciones pueden tener [MariaDB](https://mariadb.org/) instalado.

## Finalización de la instalación

Al final de la instalación, se mostrarán los datos de acceso al panel. 
No olvide guardar esta información para acceder al panel.

![Datos de acceso al panel que se muestran al finalizar la instalación](/images/en/gameapctl/gameap_finished_installation.png)

## Opciones adicionales de instalación

### Versión en desarrollo

Puede instalar la versión actualmente en desarrollo pasando los parámetros 
adicionales `--github --branch=develop` al instalador.
En este caso, la instalación tardará notablemente más, ya que se
realiza desde el código fuente.

```shell
bash <(curl -s https://gameap.com/install.sh) \
  --github \
  --branch=develop
```

### Instalación no interactiva

Este tipo de instalación le permite instalar el panel sin introducir ningún dato
durante el proceso. Pase los parámetros y el instalador no necesitará
ninguna entrada adicional por su parte.

```shell
bash <(curl -s https://gameap.com/install.sh) \
  --non-interactive \
  --host=panel.example.com \
  --port=8025 \
  --database=sqlite
```

> En `--host` indique la dirección por la que el panel estará accesible: un nombre de dominio o la
> IP externa del servidor. El valor pasa a `HTTP_HOST`, y de él el panel deduce la dirección de
> escucha: con `--host=127.0.0.1` los listeners HTTP y gRPC se levantan solo en la interfaz de
> loopback, y no podrán conectarse ni los administradores remotos ni los daemons de otras máquinas.

Parámetros principales:

| Parámetro             | Propósito                                                           |
|-----------------------|---------------------------------------------------------------------|
| `--non-interactive`   | No hacer preguntas                                                  |
| `--host`              | Dirección en la que el panel estará accesible                       |
| `--port`              | Puerto del panel, `8025` por defecto                                |
| `--grpc-port`         | Puerto gRPC para los daemons, `31718` por defecto                   |
| `--database`          | `sqlite`, `mysql` o `postgres`                                      |
| `--database-host`     | Host de la base de datos                                            |
| `--database-port`     | Puerto de la base de datos                                          |
| `--database-name`     | Nombre de la base de datos                                          |
| `--database-username` | Usuario de la base de datos                                         |
| `--database-password` | Contraseña del usuario de la base de datos                          |
| `--with-daemon`       | Instalar también GameAP Daemon                                      |
| `--version`           | Versión específica del panel                                        |

SQLite no necesita parámetros de conexión: el archivo de la base de datos se crea automáticamente.

> Los parámetros `--path` y `--web-server` son restos de GameAP 3 y no tienen efecto al
> instalar GameAP 4: el panel es un único ejecutable con una interfaz web integrada y no
> necesita un servidor web separado.

### Instalación completa

Para instalar GameAP Daemon además del propio panel, 
añada el parámetro `--with-daemon`.

Este método se recomienda si planea alojar tanto el panel 
como los servidores de juego en el mismo VDS.

```Shell
bash <(curl -s https://gameap.com/install.sh) --with-daemon
```
