---
title: GameAP Daemon
layout: default
lang: es
category: GameAP Daemon
order: 400
---

GameAP Daemon es una aplicación en segundo plano que se ejecuta en un servidor dedicado y gestiona los servidores
de juego: los instala, elimina, inicia y detiene, supervisa su estado y ejecuta
comandos del panel.

El propio daemon se conecta al panel a través de gRPC y mantiene una conexión persistente. El panel
no se conecta al daemon, y el daemon no abre puertos entrantes en el servidor dedicado.

![Arquitectura de GameAP: el panel, los daemons en los servidores dedicados y los servidores de juego](/images/en/gameap_architecture.svg)

## Instalación

### Automáticamente desde el panel

En el panel, vaya a **Administración** → **Servidores dedicados** → **Crear**, copie el comando
y ejecútelo en el servidor dedicado.

![Diálogo de instalación automática de GameAP Daemon con el comando para el servidor dedicado](/images/en/daemon/autoinstall.png)

La guía de instalación completa, que incluye Windows, el registro manual y la solución de problemas, está en
la página [Servidores dedicados](/es/gameap_configure/dedicated_servers.html).

### Gestor de procesos

Un gestor de procesos es una utilidad del sistema que inicia, detiene y reinicia los servidores de juego y supervisa
su estado. Se puede elegir durante la instalación, en el bloque "Configuración avanzada".

Si no se especifica ningún gestor, el daemon elige uno por sí mismo: en Linux — `systemd`, con una opción
alternativa cuando systemd no está disponible o el daemon se ejecuta en un contenedor; en Windows —
[Shawl](https://github.com/mtkennerly/shawl); en macOS — `tmux`.

Consulte [Gestores de procesos](/es/daemon/process_managers.html) para más detalles.

## Configuración

El archivo de configuración está en formato YAML:

* Linux — `/etc/gameap-daemon/gameap-daemon.yaml`
* Windows — `C:\gameap\daemon\gameap-daemon.yaml`

La ruta se puede establecer explícitamente con la opción `--config` (`-c`). Sin esta opción, el daemon
busca en una lista de rutas: primero el directorio actual, luego `/etc/gameap-daemon/`,
`/etc/gameap/`, después `/etc/gameap-daemon.yaml` y el directorio personal del usuario. Se utiliza el primer
archivo encontrado.

> La lista de búsqueda todavía contiene nombres con la extensión `.cfg`, un remanente de GameAP 3.
> Dicho archivo será encontrado, pero **solo se procesa YAML**: el daemon se cerrará con un
> error de formato no compatible. Utilice `.yaml` o `.yml`.

Las claves desconocidas en el archivo se ignoran. Las rutas relativas a los archivos de certificados se resuelven
respecto al directorio en el que se encuentra el propio archivo de configuración.

### Parámetros obligatorios

| Parámetro       | Tipo   | Descripción                                                    |
|-----------------|--------|----------------------------------------------------------------|
| `ds_id`         | number | ID del servidor dedicado en el panel                           |
| `api_key`       | string | Clave de acceso a la API del panel                             |
| `grpc.address`  | string | Dirección del panel en formato `host:port`                     |

Los tres se rellenan automáticamente cuando se registra el daemon.

Si `grpc.address` no está definido, la dirección se deriva del parámetro obsoleto `api_host`: se
toma el nombre del host y el puerto se sustituye por `31718`. Definir `grpc.address` explícitamente es
más fiable: `api_host` se conserva solo por compatibilidad.

### Conexión con el panel

| Parámetro                      | Predeterminado | Descripción                                                 |
|--------------------------------|----------------|--------------------------------------------------------------|
| `grpc.address`                 | —              | Dirección del panel `host:port`                              |
| `grpc.insecure`                | `false`        | Desactiva TLS. Solo para depuración                          |
| `grpc.heartbeat_interval`      | `30s`          | Intervalo de envío de heartbeat                              |
| `grpc.connect_timeout`         | `30s`          | Tiempo de espera para establecer la conexión                 |
| `grpc.initial_reconnect_delay` | `1s`           | Retraso inicial antes de la reconexión                       |
| `grpc.max_reconnect_delay`     | `60s`          | Retraso máximo antes de la reconexión                        |

Más información sobre el protocolo, la reconexión y la autenticación mutua: [GRPC API](/es/daemon/grpc.html).

### Certificados

La conexión con el panel está protegida por TLS, y el daemon presenta un certificado de cliente.
Los tres archivos son emitidos por el panel durante el registro.

| Parámetro                | Descripción                                                        |
|--------------------------|---------------------------------------------------------------------|
| `ca_certificate_file`    | Ruta al certificado de la autoridad de certificación del panel      |
| `certificate_chain_file` | Ruta al certificado del daemon                                      |
| `private_key_file`       | Ruta a la clave privada del daemon                                  |
| `private_key_password`   | Contraseña de la clave privada, si la clave está cifrada            |
| `ca_certificate`         | Certificado de la CA directamente en el archivo de configuración    |
| `certificate_chain`      | Certificado del daemon directamente en el archivo de configuración  |
| `private_key`            | Clave privada directamente en el archivo de configuración           |

Para cada uno de los tres certificados, defina la ruta del archivo o el contenido: el valor en el
archivo de configuración tiene prioridad. Con `grpc.insecure: true`, no se requieren certificados.

### Rutas

| Parámetro       | Predeterminado      | Descripción                                                        |
|-----------------|---------------------|--------------------------------------------------------------------|
| `work_path`     | —                   | Directorio de trabajo. Los archivos de los servidores de juego se encuentran en sus subdirectorios |
| `tools_path`    | `{work_path}/tools` | Directorio para herramientas auxiliares                            |
| `steamcmd_path` | —                   | Directorio de SteamCMD                                             |
| `path_7zip`     | —                   | Ruta a 7-Zip. Solo Windows                                         |
| `path_starter`  | —                   | Ruta al programa starter. Solo Windows                             |

Valores predeterminados establecidos durante la instalación: `work_path` — `/srv/gameap` en Linux y `C:\gameap` en
Windows, `steamcmd_path` — `/srv/gameap/steamcmd` y `C:\gameap\steamcmd` respectivamente.

### Registro (logging)

| Parámetro    | Predeterminado | Descripción                                         |
|--------------|----------------|-----------------------------------------------------|
| `log_level`  | `info`         | `debug`, `info`, `warn`, `error`                    |
| `output_log` | —              | Archivo de registro para mensajes normales          |
| `error_log`  | —              | Archivo de registro de errores                      |

Después de la instalación, el registro se escribe en `/var/log/gameap-daemon/output.log` en Linux y
`C:\gameap\daemon\logs\output.log` en Windows.

### Recolección de métricas

| Parámetro                     | Predeterminado | Descripción                                                    |
|-------------------------------|----------------|----------------------------------------------------------------|
| `metrics.enabled`             | `true`         | Recolección de métricas                                        |
| `metrics.collection_interval` | `5s`           | Intervalo de muestreo, no inferior a `1s`                      |
| `metrics.retention_duration`  | `10m`          | Tiempo de conservación de las muestras. Rango permitido de `10m` a `60m` |
| `if_list`                     | —              | Interfaces de red sobre las que se recopilan métricas          |
| `drives_list`                 | —              | Unidades sobre las que se recopilan métricas                   |

Los valores de `retention_duration` fuera del rango permitido se ajustan a sus límites.

### Ejecución de tareas

| Parámetro                      | Predeterminado | Descripción                                             |
|--------------------------------|----------------|---------------------------------------------------------|
| `task_manager.run_task_period` | `10ms`         | Intervalo de sondeo de la cola de tareas                |
| `task_manager.task_timeout`    | `2h`           | Tiempo máximo de ejecución de una sola tarea            |
| `task_manager.workers_count`   | —              | Número de tareas ejecutadas simultáneamente             |

### Gestor de procesos

| Parámetro                | Predeterminado         | Descripción                                  |
|--------------------------|------------------------|----------------------------------------------|
| `process_manager.name`   | detectado automáticamente | Nombre del gestor de procesos             |
| `process_manager.config` | —                      | Parámetros adicionales del gestor            |

El único parámetro adicional admitido es `scope` con el valor `system` o `user`, y solo
para `systemd`. Para otros gestores provoca un error al iniciar.

```yaml
process_manager:
  name: systemd
  config:
    scope: user
```

### Cuenta de Steam

Muchos servidores de juego no se pueden descargar a través de SteamCMD de forma anónima: se necesita
una cuenta con una copia comprada del juego.

```yaml
steam_config:
  login: your_login
  password: your_password
  group: gameap
```

El parámetro `group` establece un grupo compartido para el directorio de SteamCMD, de modo que los servidores
de juego que se ejecutan con sus propios usuarios puedan actualizar SteamCMD. Si no se establece, se utiliza
el grupo principal del usuario del servidor.

> La autenticación de dos factores debe estar desactivada en la cuenta de Steam utilizada; de lo contrario,
> el daemon no podrá iniciar sesión en SteamCMD.

### Sustitución de direcciones de repositorio

Permite sustituir el host desde el que se descargan los archivos de los servidores de juego, por ejemplo,
por un espejo más cercano.

```yaml
remote_repository_replacements:
  files.gameap.ru: cdn.gameap.com
  files.gameap.com:
    - cdn1.gameap.com
    - replace: cdn2.gameap.com
      priority: 10
```

La primera línea cubre el caso más común: `files.gameap.ru` es una dirección retirada, y las entradas
de juegos creadas antes de su desactivación todavía apuntan a ella. Esta sustitución las corrige sin
editar cada juego.

El valor puede ser una sola dirección o una lista. Los elementos de la lista pueden especificar una
`priority`: cuanto mayor sea el número, mayor será la prioridad.

### Parámetros de Windows

| Parámetro                   | Predeterminado | Descripción                                                    |
|-----------------------------|----------------|----------------------------------------------------------------|
| `use_network_service_user`  | `false`        | Ejecutar los servidores de juego como `NT AUTHORITY\NETWORK SERVICE` |
| `users`                     | —              | Contraseñas de los usuarios con los que se ejecutan los servidores de juego |

Con `use_network_service_user: true`, los servidores se ejecutan bajo una cuenta del sistema con
privilegios limitados: esta es la opción predeterminada al instalar mediante `gameapctl`. De lo contrario,
se utilizan las cuentas del bloque `users`:

```yaml
users:
  gameap_user1: password
  gameap_user2: base64:cGFyb2xi
```

La contraseña se puede escribir codificada en base64 con el prefijo `base64:`.

## Ejemplo de archivo de configuración

```yaml
ds_id: 1
api_key: your_key

grpc:
  address: panel.example.com:31718

ca_certificate_file: /etc/gameap-daemon/certs/ca.crt
certificate_chain_file: /etc/gameap-daemon/certs/server.crt
private_key_file: /etc/gameap-daemon/certs/server.key

work_path: /srv/gameap
steamcmd_path: /srv/gameap/steamcmd

if_list: []
drives_list: []

log_level: info
```

Este es exactamente el archivo que crea el comando de registro. Un ejemplo completo con todos
los parámetros y comentarios está disponible
[en el repositorio del daemon](https://github.com/gameap/daemon/blob/master/config/gameap-daemon.yaml).

## Gestión del servicio

### Linux

```bash
systemctl start gameap-daemon
systemctl stop gameap-daemon
systemctl restart gameap-daemon
systemctl status gameap-daemon
```

Actualización:

```bash
gameapctl daemon upgrade
```

### Windows

El daemon se ejecuta como el servicio **GameAP Daemon**. Se puede gestionar mediante `gameapctl.exe` o con
las herramientas estándar de Windows.
