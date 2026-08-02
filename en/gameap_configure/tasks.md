---
title: Task Scheduler
layout: default
lang: en
category: Panel settings
order: 315
---

The scheduler runs commands on a game server on a schedule: a nightly restart, regular updates,
stopping for maintenance.

Tasks are tied to a specific game server and live on the **Task Scheduler** tab of its page.

## Creating a Task

Open the game server page → the **Task Scheduler** tab → **New Task**.

| Field               | Description                                                              |
|---------------------|---------------------------------------------------------------------------|
| **Name**            | Optional, up to 128 characters. Helps tell tasks apart in the list        |
| **Command**         | What to do: start, stop, restart, update, or reinstall                    |
| **Date**            | When to run the task for the first time                                   |
| **Timezone**        | The zone the schedule is interpreted in. UTC by default                   |
| **Repeat**          | Once, forever, or a set number of times                                   |
| **Repeat period**   | The interval between runs. At least 10 minutes                            |

The available commands are the same as the server control buttons: `start`, `stop`, `restart`,
`update`, `reinstall`.

### Timezone

If no zone is set, the schedule is computed in UTC. Specify the zone if the task must fire at a
particular local time — otherwise, after a daylight saving change the task will drift relative
to the expected hour.

## Advanced Settings

The **Advanced** block defines behavior in two abnormal situations.

### Overlap Behavior

What to do when it is time for a new run but the previous one is still executing.

| Value                       | Behavior                                              |
|-----------------------------|--------------------------------------------------------|
| **Skip the new run**        | Skip the tick. The default                             |
| **Add to queue**            | Run right after the current run finishes               |

Skipping suits tasks that do not have to run every time — regular updates, for example. The
queue makes sense when every run matters.

### Missed Run Behavior

What to do with runs that fell within a period when the daemon was unavailable.

| Value                       | Behavior                                                   |
|-----------------------------|-------------------------------------------------------------|
| **Skip missed runs**        | Discard everything that was missed. The default             |
| **Run once**                | Merge all missed slots into a single run                    |

The second option protects against an avalanche: if the daemon was down for a day and the task
runs hourly, after the connection is restored it will run once, not twenty-four times.

## Enabling and Disabling

A task has an **Active** toggle. A disabled task stays in the list with the **Paused** status
and does not run. This beats deleting when the task is only needed inactive temporarily.

## Execution History

The **Execution History** button shows what has been happening with the task.

| Status        | Meaning                                                 |
|---------------|----------------------------------------------------------|
| `running`     | Running right now                                        |
| `success`     | Completed successfully                                   |
| `failed`      | Ended with an error; the reason is in the error message  |
| `canceled`    | Canceled                                                 |
| `skipped`     | Skipped by the overlap or missed-run policy              |
| `timed_out`   | The execution time limit was exceeded                    |

For each run, the start and end time, duration, exit code, error message, and command output
are stored.

Via the API, the history is available with:

```
GET /api/servers/{server}/tasks/{id}/executions
```

## Daemon Tasks

The scheduler creates daemon tasks — the same ones visible in **Administration** →
**GDaemon tasks**. Manually created tasks run there too: starting the server with a button,
installation, updates.

A stuck or unneeded task can be canceled: open it and click **Cancel**. Via the API:

```
POST /api/gdaemon_tasks/{id}/cancel
```

Cancellation helps when a task is stuck waiting — for example, because the daemon was
unavailable when it was created. The panel itself periodically marks stuck tasks: the check
interval and the threshold are set with the `TASK_REAPER_INTERVAL` and
`TASK_REAPER_STALE_THRESHOLD` variables, see the [config.env Reference](/en/config.html).

## Permissions

Working with the scheduler requires the `game-server-tasks` permission on the server, see
[Users, Roles, and Permissions](/en/users.html).

## Plugin Tasks

Plugins can register their own periodic tasks — they are not related to game servers and do not
appear in this section. Their limits are set with the `PLUGIN_SCHEDULER_*` variables, see the
[config.env Reference](/en/config.html).
