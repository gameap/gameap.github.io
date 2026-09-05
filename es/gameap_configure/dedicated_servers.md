---
title: Servidores dedicados
layout: default
lang: es
category: Configuración del panel
order: 300
---

## Nuevo servidor dedicado

Al comenzar a trabajar con el panel, debe agregar un servidor dedicado (VDS/VPS, contenedor, servidor
físico) en el que se ejecutará GameAP Daemon.

Vaya a la página **"Administración"** → **"Servidores dedicados"** y haga clic en el botón **"Crear"**.
Se abrirá una ventana con un comando de instalación listo para usar: uno para Linux y otro para Windows.

Los certificados los crea y los entrega al daemon el propio panel: no es necesario generarlos
ni firmarlos manualmente.

> La clave de configuración es válida durante **1 hora** y es de un solo uso: se revoca una vez que
> el daemon se ha registrado correctamente. La clave es compartida en todo el panel, por lo que al
> volver a abrir la ventana "Crear" se genera una nueva clave, y un comando copiado anteriormente
> deja de funcionar.

### Qué necesita en el servidor dedicado

* Derechos de superusuario (root en Linux, administrador en Windows).
* Acceso saliente al panel por el puerto **31718/TCP**: el daemon lo utiliza para comunicarse con
  el panel a través de gRPC. Es el único puerto del panel que necesita un daemon en ejecución.
* Acceso saliente a `github.com` y `api.github.com`: desde allí se descargan gameapctl y el propio daemon.
* Para el script de Linux: tener instalados `curl`, `tar` e `install`.
* Arquitectura: `amd64`, `arm64`, `386` o `arm`.

El puerto **8025** es la interfaz web y la API del panel. Lo necesitan el navegador del administrador
y el comando `curl` del comando de instalación de una línea, pero un daemon en ejecución no lo requiere.

El puerto **31717** es el que el daemon informa como propio durante el registro. En GameAP 4 el panel
no establece conexiones entrantes con el daemon, por lo que no es necesario abrir este puerto en el
servidor dedicado.

### Instalación en Linux

Copie el comando de la pestaña Linux y ejecútelo en el servidor dedicado **como root**:

```bash
bash <(curl -fsSL 'https://your-panel/nodes/setup/zItWHWlI4RKPl9ZsYc3y3WgdKq7mNvBx')
```

El script comprueba el entorno, instala `gameapctl` en `/usr/local/bin` (o lo actualiza si ya está
instalado) y luego ejecuta `gameapctl daemon install`. Esto instala GameAP Daemon, crea el usuario
`gameap`, instala SteamCMD, registra el daemon en el panel y lo inicia como el servicio systemd
`gameap-daemon`.

El comando debe ejecutarse como root, no a través de `sudo`: la sustitución de procesos `bash <(...)`
no sobrevive a `sudo`, y el script se detendrá con un error. Si no está trabajando como root, descargue
el script a un archivo y ejecútelo:

```bash
curl -fsSL 'https://your-panel/nodes/setup/zItWHWlI4RKPl9ZsYc3y3WgdKq7mNvBx' -o gameap-setup.sh
sudo bash gameap-setup.sh
```

Rutas de instalación:

| Elemento            | Ruta                                    |
|---------------------|-----------------------------------------|
| Daemon              | `/usr/bin/gameap-daemon`                |
| Configuración       | `/etc/gameap-daemon/gameap-daemon.yaml` |
| Certificados        | `/etc/gameap-daemon/certs`              |
| Directorio de trabajo | `/srv/gameap`                         |
| SteamCMD            | `/srv/gameap/steamcmd`                  |
| Registros           | `/var/log/gameap-daemon/output.log`     |

### Instalación en Windows

No hay un comando de PowerShell de una línea; la instalación se realiza a través de gameapctl:

1. Descargue el archivo de gameapctl para su arquitectura desde la página de
   [versiones de gameapctl](https://github.com/gameap/gameapctl/releases). El archivo se llama
   `gameapctl-<version>-windows-amd64.zip`.
2. Descomprima el archivo y ejecute `gameapctl.exe`.
3. En la sección **GameAP Daemon** haga clic en **Install**.
4. Pegue la cadena de la pestaña Windows del panel en el campo **Connect URL**.
5. Haga clic en **Install**.

Lo mismo se puede hacer con un comando en la consola:

```shell
gameapctl daemon install --connect=grpc://your-panel:31718/zItWHWlI4RKPl9ZsYc3y3WgdKq7mNvBx
```

Rutas de instalación:

| Elemento            | Ruta                                     |
|---------------------|------------------------------------------|
| Daemon              | `C:\gameap\daemon\gameap-daemon.exe`     |
| Configuración       | `C:\gameap\daemon\gameap-daemon.yaml`    |
| Certificados        | `C:\gameap\daemon\certs`                 |
| Directorio de trabajo | `C:\gameap`                            |
| SteamCMD            | `C:\gameap\steamcmd`                     |
| Registros           | `C:\gameap\daemon\logs\output.log`       |

El daemon se registra como el servicio **GameAP Daemon**.

### Configuración avanzada de la instalación

La ventana de creación tiene un bloque plegable **"Advanced Settings"** (Configuración avanzada):

* **Process manager**: qué debe gestionar los procesos de los servidores de juego. En Linux están
  disponibles `systemd`, `docker`, `podman`, `tmux` y `simple`; en Windows, `winsw`, `shawl` y
  `simple`. Se selecciona automáticamente por defecto. Consulte los detalles en la página
  [Gestores de procesos](/es/daemon/process_managers.html).
* **GitHub**: compilar el daemon desde el código fuente en lugar de usar una versión lista para usar.
* **Branch**: rama del repositorio, si el daemon se compila desde el código fuente.

El panel añade los valores seleccionados al comando de instalación:

```bash
bash <(curl -s '...') --config='process_manager.name=docker' --github --branch=master
```

### Instalación manual

Si el script automático no se ajusta a sus necesidades —una distribución no estándar, sus propias
reglas de instalación de paquetes, la instalación en una imagen preparada—, el daemon se puede
registrar manualmente.

1. Descargue el binario `gameap-daemon` para su plataforma desde la página de
   [versiones del daemon](https://github.com/gameap/daemon/releases) y colóquelo en el servidor dedicado.
2. Abra la ventana **"Crear"** en el panel y tome la URL de conexión de la forma
   `grpc://host:port/key` del comando de Windows.
3. Ejecute el registro:

```bash
gameap-daemon enroll --connect=grpc://your-panel:31718/zItWHWlI4RKPl9ZsYc3y3WgdKq7mNvBx
```

El daemon se conecta al panel, crea en él un registro de servidor dedicado, recibe los certificados
(`ca.crt`, `server.crt`, `server.key` con permisos `0600` en un directorio con permisos `0700`)
y escribe el archivo de configuración.

| Flag              | Valor por defecto en Linux              | Valor por defecto en Windows          | Propósito                                                                   |
|-------------------|-----------------------------------------|---------------------------------------|----------------------------------------------------------------------------|
| `--connect`       | —                                       | —                                     | URL de conexión. Obligatorio                                                |
| `--config-path`   | `/etc/gameap-daemon/gameap-daemon.yaml` | `C:\gameap\daemon\gameap-daemon.yaml` | Dónde escribir la configuración                                             |
| `--certs-dir`     | `/etc/gameap-daemon/certs`              | `C:\gameap\daemon\certs`              | Dónde guardar los certificados                                              |
| `--work-path`     | `/srv/gameap`                           | `C:\gameap`                           | Directorio de trabajo de los servidores de juego                            |
| `--steamcmd-path` | `/srv/gameap/steamcmd`                  | `C:\gameap\steamcmd`                  | Directorio de SteamCMD                                                      |
| `--listen-ip`     | `0.0.0.0`                               | `0.0.0.0`                             | IP que el daemon informa sobre sí mismo. Con `0.0.0.0` se detecta automáticamente |
| `--listen-port`   | `31717`                                 | `31717`                               | Puerto que el daemon informa sobre sí mismo                                 |

> El comando `enroll` sobrescribe el archivo de configuración **por completo**: sin fusionarlo con
> la configuración existente y sin copia de seguridad. Si el daemon ya ha sido configurado, guarde
> la configuración de antemano.

Después de registrar el daemon, debe registrarlo como servicio del sistema e iniciarlo usted mismo.

### Si la instalación falla

| Mensaje                                         | Causa                                                                                                                                        |
|-------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| `This script must be run as root.`              | El script no se ejecutó como root. Descárguelo a un archivo y ejecútelo a través de `sudo`, como se muestra arriba                            |
| `Error: 'curl' is required but not installed.`  | El servidor no tiene `curl`, `tar` o `install`. Instale el paquete que falta e inténtelo de nuevo                                             |
| `Unsupported architecture: ...`                 | No hay compilaciones listas para la arquitectura del servidor                                                                                 |
| `Failed to detect latest gameapctl version`     | No hay acceso a `api.github.com`, o se agotó el límite de la API de GitHub (60 solicitudes por hora desde una dirección). Espere o instale gameapctl manualmente |
| `cannot reach gRPC server at ...`               | El puerto 31718 del panel no es accesible desde el servidor dedicado. Compruebe el firewall y la dirección del panel — consulte [GRPC API](/es/daemon/grpc.html) |
| El enlace devuelve `403`                        | La clave de configuración ha caducado o ha sido regenerada. Abra de nuevo la ventana "Crear" y copie el nuevo comando                         |

## Edición de servidores dedicados

Para editar un servidor dedicado (nodo), vaya a la página **"Administración"** → **"Servidores dedicados"**,
luego seleccione el servidor dedicado que desea editar y haga clic en el botón **"Editar"**.

### Descripción de los parámetros

#### Básicos

##### Nombre

Nombre del servidor dedicado. Puede tomar cualquier valor no vacío; no afecta a ninguna función.

##### Directorio de trabajo

Este directorio contiene los scripts básicos para gestionar los procesos de los servidores de juego.
Los subdirectorios del directorio de trabajo contienen los archivos de los servidores de juego. Para
la ruta especificada se asigna el [directorio del servidor de juego](/es/gameap_configure/game_servers.html#directorio).
El valor por defecto es `/srv/gameap`.

##### Ruta a SteamCMD

Ruta al directorio de SteamCMD (allí se encuentra el script `steamcmd.sh`). El valor por defecto es `/srv/gameap/steamcmd`.

##### Lista de IP

Lista de IP o hosts donde se ejecutarán los servidores de juego.

#### Scripts

Plantillas de comandos que el panel utiliza para gestionar los servidores de juego en este servidor
dedicado. Completarlas es opcional: si un campo está vacío, se usa el comando por defecto.

#### GameAP Daemon

Datos de conexión del daemon. Se rellenan automáticamente cuando se registra el servidor dedicado
y, por lo general, no es necesario modificarlos manualmente.

En GameAP 4 la conexión la establece el daemon: se conecta él mismo al panel a través de gRPC y
mantiene una conexión persistente. El panel no se conecta al daemon.

##### Host y puerto de GameAP Daemon

La dirección y el puerto que el daemon informó sobre sí mismo durante el registro. El puerto por
defecto es `31717`. Los valores son informativos: el panel no los utiliza para conectarse, y no es
necesario abrir este puerto en el servidor dedicado.

##### Usuario y contraseña de GameAP Daemon

Campos heredados de GameAP 3, donde el panel se conectaba al daemon por sí mismo. No se utilizan en
GameAP 4 y no es necesario rellenarlos.
