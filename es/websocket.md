---
title: WebSocket y métricas
layout: default
lang: es
category: Administración
order: 337
---

A través de WebSocket el panel sirve todo lo que se actualiza en tiempo real: el progreso de
las tareas, la consola del servidor de juego y las métricas. La interfaz funciona a través de
estas mismas conexiones.

## Conexión

Las conexiones se abren en seis direcciones:

| Dirección                            | Qué sirve                                          |
|--------------------------------------|----------------------------------------------------|
| `/api/ws/tasks/{id}`                 | Progreso de la tarea y su salida                   |
| `/api/ws/servers/{server}/console`   | La consola del servidor de juego, en ambos sentidos|
| `/api/ws/servers/{server}/attach`    | Una sesión interactiva con el servidor de juego    |
| `/api/ws/servers/{server}/metrics`   | Métricas del servidor de juego                     |
| `/api/ws/nodes/{id}/metrics`         | Métricas del servidor dedicado                     |
| `/api/ws/nodes/metrics`              | Métricas de todos los servidores dedicados         |

### Autorización

El token se pasa en un parámetro de consulta porque no se pueden establecer cabeceras al abrir
un WebSocket desde un navegador:

```text
wss://panel.example.com:8025/api/ws/servers/1/console?token=<token>
```

**Solo se aceptan tokens de corta duración** con el prefijo `glst_` **en el parámetro `token`**.
Un token de acceso personal no se puede pasar ahí, de modo que una clave de acceso de larga
duración no acabará en la URL ni por error.

> El propio token de corta duración sí permanece en la dirección y puede quedar registrado en los
> registros de un proxy inverso, de un servidor web o de los sistemas de observabilidad. Que sea de
> un solo uso y viva 10 segundos deja esa entrada sin valor, pero si los registros se conservan
> mucho tiempo, es mejor recortar de ellos el parámetro `token`.

Obtención de un token de corta duración:

```bash
curl -X POST https://panel.example.com:8025/api/auth/short-lived-token \
  -H "Authorization: Bearer <session token>"
```

El token es de un solo uso y no vive más de 10 segundos, por lo que debe solicitarse
inmediatamente antes de abrir la conexión. Consulte [API y tokens](/es/api.html) para más
detalles.

## Formato de los frames

Todos los mensajes son JSON con la misma estructura:

```json
{
  "type": "task.status",
  "payload": { },
  "ts": 1711234567
}
```

| Campo     | Descripción                                  |
|-----------|----------------------------------------------|
| `type`    | Tipo de mensaje                              |
| `payload` | Datos, dependen del tipo                     |
| `ts`      | Marca de tiempo, en segundos Unix            |

### Tipos de mensajes

**Tareas** (`/api/ws/tasks/{id}`):

| Tipo            | Cuándo llega                       |
|-----------------|------------------------------------|
| `task.status`   | Cambió el estado de la tarea       |
| `task.output`   | Apareció nueva salida del comando  |
| `task.complete` | La tarea finalizó                  |

**Consola** (`/api/ws/servers/{server}/console`):

| Tipo              | Dirección        | Propósito                            |
|-------------------|------------------|--------------------------------------|
| `console.history` | Desde el panel   | Salida acumulada al conectar         |
| `console.command` | Hacia el panel   | Envío de un comando al servidor      |

**Sesión interactiva** (`/api/ws/servers/{server}/attach`):

| Tipo            | Dirección     | Propósito                       |
|-----------------|---------------|---------------------------------|
| `attach.input`  | Hacia el panel | Entrada en la sesión           |
| `attach.detach` | Hacia el panel | Desconectar sin detener el servidor |

**Métricas**:

| Tipo                   | Propósito                                           |
|------------------------|-----------------------------------------------------|
| `metrics.replay`       | Valores acumulados del período de retención         |
| `metrics.replay.done`  | Valores acumulados entregados; siguen los valores en vivo |
| `metrics.error`        | Error de recolección de métricas                    |

Al conectarse a las métricas, el panel primero sirve el historial del período de retención y
luego envía los nuevos valores a medida que llegan. El cambio se marca con un frame
`metrics.replay.done`.

> En la consola del servidor de juego, la contraseña RCON está enmascarada: el servidor de
> juego la muestra como parte de la línea de comandos de inicio.

## Series de métricas

Las métricas son recolectadas por el daemon y pasadas al panel. El intervalo de muestreo y el
tiempo de retención se configuran en la configuración del daemon con
`metrics.collection_interval` (5 segundos por defecto) y `metrics.retention_duration` (10
minutos por defecto); consulte [GameAP Daemon](/es/daemon/daemon.html#recolección-de-métricas).

### Servidor de juego

| Serie                                       | Valor                                      |
|---------------------------------------------|--------------------------------------------|
| `gameap_server_up`                          | Si el servidor está en ejecución           |
| `gameap_server_cpu_usage_percent`           | Uso de CPU, porcentaje                     |
| `gameap_server_memory_usage_bytes`          | Memoria utilizada, bytes                   |
| `gameap_server_memory_limit_bytes`          | Límite de memoria, bytes                   |
| `gameap_server_memory_usage_percent`        | Memoria utilizada, porcentaje del límite   |
| `gameap_server_network_receive_bytes_total` | Bytes de red recibidos, acumulativo        |
| `gameap_server_network_transmit_bytes_total`| Bytes de red transmitidos, acumulativo     |
| `gameap_server_block_io_read_bytes_total`   | Bytes de disco leídos, acumulativo         |
| `gameap_server_block_io_write_bytes_total`  | Bytes de disco escritos, acumulativo       |
| `gameap_server_process_pids`                | Número de procesos del servidor            |

Las métricas de límite de memoria solo las proporcionan los gestores de procesos que pueden
aplicarlo: systemd, Docker y Podman.

### Servidor dedicado

| Serie                                     | Valor                                          |
|-------------------------------------------|------------------------------------------------|
| `gameap_node_cpu_usage_percent`           | Uso de CPU, porcentaje                         |
| `gameap_node_memory_usage_bytes`          | Memoria utilizada, bytes                       |
| `gameap_node_memory_total_bytes`          | Memoria total, bytes                           |
| `gameap_node_memory_usage_percent`        | Memoria utilizada, porcentaje                  |
| `gameap_node_swap_usage_bytes`            | Swap utilizado, bytes                          |
| `gameap_node_swap_total_bytes`            | Swap total, bytes                              |
| `gameap_node_disk_usage_bytes`            | Disco utilizado, bytes                         |
| `gameap_node_disk_total_bytes`            | Disco total, bytes                             |
| `gameap_node_disk_usage_percent`          | Disco utilizado, porcentaje                    |
| `gameap_node_network_receive_bytes_total` | Bytes de red recibidos, acumulativo            |
| `gameap_node_network_transmit_bytes_total`| Bytes de red transmitidos, acumulativo         |
| `gameap_node_load1`                       | Carga promedio de 1 minuto                     |
| `gameap_node_load5`                       | Carga promedio de 5 minutos                    |
| `gameap_node_load15`                      | Carga promedio de 15 minutos                   |
| `gameap_node_uptime_seconds_total`        | Tiempo de actividad, segundos                  |

La carga promedio no se recolecta en Windows.

Qué interfaces de red y unidades de disco incluir se define con los parámetros `if_list` y
`drives_list` en la configuración del daemon.

## Limitaciones

Las métricas se mantienen en la RAM del daemon no más de `metrics.retention_duration` (se
permite de 10 a 60 minutos) y se pierden cuando el daemon se reinicia. El panel no tiene
almacenamiento a largo plazo ni exportación a sistemas de monitorización.

La recolección de métricas se desactiva con `metrics.enabled: false` en la configuración del
daemon.
