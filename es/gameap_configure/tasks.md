---
title: Planificador de tareas
layout: default
lang: es
category: Configuración del panel
order: 315
---

El planificador ejecuta comandos en un servidor de juego según una programación: un reinicio
nocturno, actualizaciones periódicas, una detención por mantenimiento.

Las tareas están vinculadas a un servidor de juego concreto y se encuentran en la pestaña
**Planificador de tareas** de su página.

## Crear una tarea

Abra la página del servidor de juego → la pestaña **Planificador de tareas** → **Nueva tarea**.

| Campo                 | Descripción                                                               |
|-----------------------|---------------------------------------------------------------------------|
| **Nombre**            | Opcional, hasta 128 caracteres. Ayuda a distinguir las tareas en la lista |
| **Comando**           | Qué hacer: iniciar, detener, reiniciar, actualizar o reinstalar           |
| **Fecha**             | Cuándo ejecutar la tarea por primera vez                                  |
| **Zona horaria**      | La zona en la que se interpreta la programación. UTC por defecto          |
| **Repetir**           | Una vez, para siempre o un número determinado de veces                    |
| **Período de repetición** | El intervalo entre ejecuciones. Al menos 10 minutos                   |

Los comandos disponibles son los mismos que los botones de control del servidor: `start`, `stop`,
`restart`, `update`, `reinstall`.

### Zona horaria

Si no se establece ninguna zona, la programación se calcula en UTC. Especifique la zona si la
tarea debe ejecutarse a una hora local concreta; de lo contrario, tras un cambio de horario de
verano la tarea se desviará respecto a la hora esperada.

## Configuración avanzada

El bloque **Avanzado** define el comportamiento en dos situaciones anómalas.

### Comportamiento ante solapamientos

Qué hacer cuando llega el momento de una nueva ejecución pero la anterior sigue en curso.

| Valor                         | Comportamiento                                         |
|-------------------------------|--------------------------------------------------------|
| **Omitir la nueva ejecución** | Omitir la ejecución. Valor por defecto                 |
| **Añadir a la cola**          | Ejecutar justo después de que termine la ejecución actual |

La omisión conviene a tareas que no tienen que ejecutarse cada vez; las actualizaciones
periódicas, por ejemplo. La cola tiene sentido cuando cada ejecución importa.

### Comportamiento ante ejecuciones perdidas

Qué hacer con las ejecuciones que cayeron en un período en el que el daemon no estaba disponible.

| Valor                             | Comportamiento                                            |
|-----------------------------------|------------------------------------------------------------|
| **Omitir las ejecuciones perdidas** | Descartar todo lo que se perdió. Valor por defecto       |
| **Ejecutar una vez**              | Fusionar todas las ejecuciones perdidas en una sola        |

La segunda opción protege contra una avalancha: si el daemon estuvo caído un día y la tarea se
ejecuta cada hora, tras restablecerse la conexión se ejecutará una vez, no veinticuatro.

## Activación y desactivación

Una tarea tiene un interruptor **Activa**. Una tarea desactivada permanece en la lista con el
estado **En pausa** y no se ejecuta. Esto es preferible a eliminarla cuando la tarea solo
necesita estar inactiva temporalmente.

## Historial de ejecución

El botón **Historial de ejecución** muestra lo que ha estado ocurriendo con la tarea.

| Estado        | Significado                                                   |
|---------------|----------------------------------------------------------------|
| `running`     | Se está ejecutando ahora mismo                                 |
| `success`     | Completada correctamente                                       |
| `failed`      | Terminó con un error; la razón está en el mensaje de error     |
| `canceled`    | Cancelada                                                      |
| `skipped`     | Omitida por la política de solapamiento o de ejecuciones perdidas |
| `timed_out`   | Se superó el límite de tiempo de ejecución                     |

Para cada ejecución se almacenan la hora de inicio y de fin, la duración, el código de salida,
el mensaje de error y la salida del comando.

A través de la API, el historial está disponible con:

```http
GET /api/servers/{server}/tasks/{id}/executions
```

## Tareas del daemon

El planificador crea tareas del daemon: las mismas que se ven en **Administración** →
**Tareas de GDaemon**. Allí también se ejecutan las tareas creadas manualmente: iniciar el
servidor con un botón, la instalación, las actualizaciones.

Una tarea atascada o innecesaria se puede cancelar: ábrala y haga clic en **Cancelar**. A
través de la API:

```http
POST /api/gdaemon_tasks/{id}/cancel
```

La cancelación ayuda cuando una tarea está atascada en espera, por ejemplo, porque el daemon no
estaba disponible cuando se creó. El propio panel marca periódicamente las tareas atascadas: el
intervalo de comprobación y el umbral se configuran con las variables `TASK_REAPER_INTERVAL` y
`TASK_REAPER_STALE_THRESHOLD`; consulte la [referencia de config.env](/es/config.html).

## Permisos

Para trabajar con el planificador se requiere el permiso `game-server-tasks` en el servidor;
consulte [Usuarios, roles y permisos](/es/users.html).

## Tareas de plugins

Los plugins pueden registrar sus propias tareas periódicas; no están relacionadas con servidores
de juego y no aparecen en esta sección. Sus límites se configuran con las variables
`PLUGIN_SCHEDULER_*`; consulte la [referencia de config.env](/es/config.html).
