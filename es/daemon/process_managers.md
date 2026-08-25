---
title: Gestores de procesos
layout: default
lang: es
category: GameAP Daemon
order: 410
---

Un gestor de procesos es una utilidad del sistema que administra los procesos en el servidor.
Se encarga de iniciar, detener y reiniciar los servidores de juego, supervisar su estado, recopilar estadísticas y limitar los recursos (CPU / RAM).
De forma predeterminada, se utiliza systemd en Linux y [Shawl](https://github.com/mtkennerly/shawl) en Windows.

## Configuración

El gestor de procesos se configura en el archivo de configuración de GameAP Daemon.

### Estructura básica

```yaml
process_manager:
  name: <manager_name>
  config:
    <key>: <value>
```

**Gestores disponibles:**
- Linux: `systemd` (predeterminado), `docker`, `podman`, `tmux`
- Windows: `shawl` (predeterminado), `winsw`

## Linux

### Systemd

Se utiliza de forma predeterminada en Linux. Es un gestor de procesos moderno que ofrece alto rendimiento y fiabilidad.
Las capacidades de aislamiento de este gestor de procesos son limitadas.

| Función                                |                 |
|----------------------------------------|-----------------|
| Inicio, detención y reinicio de servidores | ✅               |
| Estadísticas                           | ⚠️              |
| Límites de recursos (CPU / RAM)        | ✅               |
| Lectura de consola                     | ✅               |
| Envío de comandos a la consola         | ✅               |
| Aislamiento                            | ⚠️ Limitado     |

#### Configuración de Systemd

Systemd no requiere configuración adicional. El Daemon crea y administra automáticamente los archivos de unidad.

##### Ejemplo de configuración

```yaml
process_manager:
  name: systemd
```

### Docker

Docker se utiliza para aislar los servidores de juego en contenedores.
Utilice Docker si planea usar Pterodactyl Eggs o Pelican Eggs.

| Función                                |    |
|----------------------------------------|----|
| Inicio, detención y reinicio de servidores | ✅  |
| Estadísticas                           | ⚠️ |
| Límites de recursos (CPU / RAM)        | ✅  |
| Lectura de consola                     | ✅  |
| Envío de comandos a la consola         | ✅  |
| Aislamiento                            | ✅  |

#### Configuración de Docker

##### Ejemplos de configuración

Mínima:
```yaml
process_manager:
  name: docker
```

### Podman

Podman es una alternativa a Docker que proporciona aislamiento de los servidores de juego en contenedores.
También puede utilizar Podman si planea usar Pterodactyl Eggs o Pelican Eggs.

| Función                                |    |
|----------------------------------------|----|
| Inicio, detención y reinicio de servidores | ✅  |
| Estadísticas                           | ⚠️ |
| Límites de recursos (CPU / RAM)        | ✅  |
| Lectura de consola                     | ✅  |
| Envío de comandos a la consola         | ✅  |
| Aislamiento                            | ✅  |

#### Configuración de Podman

Podman utiliza los mismos parámetros que Docker, más un parámetro específico para conectarse a la API.

##### Parámetros específicos

| Clave de configuración | Descripción | Valor predeterminado |
|---------------|-------------|---------|
| `socket_path` | Ruta al socket unix de Podman | `unix:///run/user/{UID}/podman/podman.sock` (rootless) |

##### Ejemplos de configuración

Rootless (predeterminado):
```yaml
process_manager:
  name: podman
```

Rootful:
```yaml
process_manager:
  name: podman
  config:
    socket_path: "unix:///run/podman/podman.sock"
```

### Tmux

[Tmux](https://github.com/tmux/tmux) es un multiplexor de terminal.
Actualmente es un gestor de procesos obsoleto en GameAP,
cuyo uso no se recomienda, ya que no proporciona funcionalidad completa.

Tmux puede utilizarse en sistemas antiguos, así como dentro de contenedores (LXC, Docker/Podman, etc.), sistemas virtuales
y sistemas que no disponen de Systemd o no pueden usar Docker o Podman.

| Función                                |    |
|----------------------------------------|----|
| Inicio, detención y reinicio de servidores | ✅  |
| Estadísticas                           | ⚠️ |
| Límites de recursos (CPU / RAM)        | ❌  |
| Lectura de consola                     | ✅  |
| Envío de comandos a la consola         | ✅  |
| Aislamiento                            | ❌  |

#### Configuración de Tmux

Tmux no requiere configuración adicional.

##### Ejemplo de configuración

```yaml
process_manager:
  name: tmux
```

## Windows

### Shawl

[Shawl](https://github.com/mtkennerly/shawl) es un gestor de procesos ligero para Windows
que proporciona funcionalidad básica para administrar servidores de juego.
Está escrito en Rust y utiliza la API de Windows para ejecutar aplicaciones como servicios de Windows.

| Función                                |    |
|----------------------------------------|----|
| Inicio, detención y reinicio de servidores | ✅  |
| Estadísticas                           | ⚠️ |
| Límites de recursos (CPU / RAM)        | ❌  |
| Lectura de consola                     | ✅  |
| Envío de comandos a la consola         | ❌  |
| Aislamiento                            | ❌  |

#### Configuración de Shawl

Shawl no requiere configuración adicional.

##### Detalles de funcionamiento

| Parámetro | Valor |
|-----------|-------|
| Directorio de configuración | `C:\gameap\services` |
| Tiempo de espera de detención | 10000 ms |
| Rotación de registros | Diaria |
| Retención de registros | 7 días |

##### Ejemplo de configuración

```yaml
process_manager:
  name: shawl
```

### WinSW

[WinSW](https://github.com/winsw/winsw) (Windows Service Wrapper) es un gestor de procesos escrito en C#.
Permite ejecutar aplicaciones como servicios de Windows.
En GameAP es un gestor obsoleto y ha sido reemplazado por Shawl.

| Función                                |    |
|----------------------------------------|----|
| Inicio, detención y reinicio de servidores | ✅  |
| Estadísticas                           | ⚠️ |
| Límites de recursos (CPU / RAM)        | ❌  |
| Lectura de consola                     | ✅  |
| Envío de comandos a la consola         | ❌  |
| Aislamiento                            | ❌  |

#### Configuración de WinSW

WinSW no requiere configuración adicional.

##### Detalles de funcionamiento

| Parámetro | Valor |
|-----------|-------|
| Directorio de configuración | `C:\gameap\services` |
| Formato de configuración | XML |

##### Ejemplo de configuración

```yaml
process_manager:
  name: winsw
```
