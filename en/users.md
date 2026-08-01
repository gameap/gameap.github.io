---
title: Users, Roles, and Permissions
layout: default
lang: en
category: Administration
order: 332
---

* This will become a table of contents (this text will be scraped).
{:toc}

Access in the panel is built from three concepts:

* **Permission** — a single allowed action, for example "start a game server".
* **Role** — a named set of permissions. A role is assigned to a user.
* **Direct permission** — a permission granted to a user bypassing roles.

Permissions can be granted globally or tied to a specific object — most often a single game
server. That is how a user is given access to their own server and nothing else.

## Default roles

Two roles are created during the initial seeding of the database:

| Role    | Name in the interface | Permissions                     |
|---------|-----------------------|---------------------------------|
| `admin` | Administrator         | `admin roles & permissions`     |
| `user`  | User                  | None                            |

The first created user gets the `admin` role. The `user` role contains no permissions at all: a
regular user is given access not by the role but by permissions on specific game servers.

> The **`admin roles & permissions`** permission granted globally is what makes an administrator.
> It gates access to the "Administration" section, and it also determines who mandatory
> two-factor authentication applies to — see [Security](/en/security.html).
>
> There is no special "superadministrator" role in the panel, and the user with `id = 1` has no
> privileges by default: only having this permission matters.

## Game server permissions

| Permission                 | Name in the interface                |
|----------------------------|---------------------------------------|
| `game-server-common`       | Common Game Server Ability            |
| `game-server-start`        | Start Game Server                     |
| `game-server-stop`         | Stop Game Server                      |
| `game-server-restart`      | Restart Game Server                   |
| `game-server-pause`        | Pause Game Server                     |
| `game-server-update`       | Update Game Server                    |
| `game-server-files`        | Access to filemanager                 |
| `game-server-tasks`        | Access to task scheduler              |
| `game-server-settings`     | Access to settings                    |
| `game-server-console-view` | Access to read server console         |
| `game-server-console-send` | Access to send console commands       |
| `game-server-rcon-console` | RCON console                          |
| `game-server-rcon-players` | RCON players manage                   |
| `game-server-metrics`      | Access to server metrics              |

`game-server-common` gives basic access to the server: seeing it in the list and opening its
page. Without it the other permissions are practically useless — start with it.

Console read and write permissions are separate: `game-server-console-view` only allows viewing,
`game-server-console-send` — sending commands. The same goes for RCON: access to the RCON console
and player management are granted separately.

Besides game server permissions there are generic `create`, `view`, `edit`, `delete` permissions —
they apply to panel sections.

## Granting permissions on a specific server

**Administration** → **Users** → select a user → **Edit**.

The form has two blocks:

* **Roles** — `Administrator` or `User`.
* **Game servers** — a list of servers, with the needed permissions checkable for each one.

Permissions granted here are tied to the specific server: the user gets access to it and nothing
more.

The same is done via the API:

```
GET  /api/users/{id}/servers                              — the user's servers
GET  /api/users/{id}/servers/{server}/permissions         — permissions on the server
PUT  /api/users/{id}/servers/{server}/permissions         — change permissions
```

Your own permissions on a server can be viewed with `GET /api/servers/{server}/abilities`.

## How access is checked

A user's permissions are collected from two sources: those granted directly and those coming
from all of their roles. Then two rules apply.

**A denial is stronger than a grant.** If there is a deny record for the permission in question,
access is not given, no matter how many grants were issued by other paths or in what order they
are recorded.

**The scope must match.** A permission granted for one server has no effect on another. A denial
set on one server likewise does not affect the rest.

A permission can be granted in three scopes:

| Scope                       | Meaning                                                          |
|-----------------------------|-------------------------------------------------------------------|
| Globally                    | Applies everywhere                                                |
| On an object type           | Applies to all game servers, all dedicated servers, and so on     |
| On a specific object        | Applies to this server only                                       |

Checks that need global access — the administrator check, for example — count **only a globally
granted permission**. A grant tied to a single server does not make one an administrator.

> Revoking a permission in the interface does not always mean simply deleting the record. If
> after the deletion the user would still receive the permission through a role, the panel
> additionally sets an explicit denial — otherwise the permission would come back through the
> role.

### Caching

Check results are cached for the time from `RBAC_CACHE_TTL` (30 seconds by default). Permission
changes made through the interface and the API flush the cache immediately. Edits made directly
in the database do not flush the cache — those changes take effect within the configured time.

## Managing users

| Method and path          | Purpose                        |
|--------------------------|---------------------------------|
| `GET /api/users`         | User list                       |
| `POST /api/users`        | Create a user                   |
| `GET /api/users/{id}`    | User details                    |
| `PUT /api/users/{id}`    | Update a user                   |
| `DELETE /api/users/{id}` | Delete a user                   |

The user list is available to administrators only.

The password set when creating or updating a user goes through the password policy check: at
least 12 bytes and not on the compromised password list. See [Security](/en/security.html) for
details.

There is no user self-registration in the panel: accounts are created by an administrator.

## Common tasks

### Give a user access to a single server

1. Create a user: **Administration** → **Users** → **Create**.
2. Assign them the `User` role.
3. In the "Game servers" block, select the server and check the permissions. The minimal working
   set: `game-server-common`, `game-server-start`, `game-server-stop`, `game-server-restart`.
4. Add as needed: the file manager, the console, RCON, tasks.

### Make a user an administrator

Assign the `Administrator` role. It contains the `admin roles & permissions` permission, which
gives full access.

> From then on the two-factor authentication requirement applies to the user: on their next
> login they will see a reminder, and after 30 days — a demand to enable it.
> See [Security](/en/security.html).

### Revoke access without deleting the account

Remove the server permissions and leave the `User` role. The user will be able to log in but
will not see a single server. Active sessions are not terminated by this — they keep working
until they expire (24 hours, or 7 days for logins with "remember me" checked).
