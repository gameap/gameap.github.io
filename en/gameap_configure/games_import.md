---
title: Games Import
layout: default
lang: en
category: Panel settings
order: 330
---

Game and mod settings can be moved between GameAP installations and taken from other control
panels. This capability appeared in GameAP 4.1.

All of this lives on the **Administration** → **Games** page.

![](/images/en/gameap_configure/games_import/import_button.png)

## GameAP's Own Format

The GameAP format carries the game in full: the game itself, all of its mods, variables, run
commands, and RCON commands.

### Export

On the **Administration** → **Games** page, select a game and click **Export**. The panel will
serve a `<game-code>.gameap.yaml` file.

Via the API:

```
GET /api/games/{code}/export
```

### Import

**Administration** → **Games** → **Import GameAP YAML**, upload the file and click **Import**.

![](/images/en/gameap_configure/games_import/import_page.png)

Via the API:

```
POST /api/games/import/gameap
```

The import creates the game if it does not exist yet and **updates the existing one** if a game
with that code is already present. Mods are imported together with the game — after the upload
the panel reports how many were imported.

The format version is given at the top of the file:

```yaml
schema_version: "1.0"
```

The panel will not accept a file with a different format version. This field is also worth
checking when a transfer fails.

## Importing from Other Panels

Importing templates is supported from:

* [Pterodactyl](https://pterodactyl.io/)
* [Pelican](https://pelican.dev/)

**Administration** → **Games** → **Import Pelican Egg**, upload the egg file in JSON or YAML
format.

Via the API:

```
POST /api/games/import/pelican-egg
```

Both panels share the same egg format, so Pterodactyl and Pelican files are uploaded the same
way.

## Upgrading Games from the GameAP Catalog

The panel ships with a catalog of ready-made game and mod settings. To pull changes from it,
click **Upgrade Games** on the **Administration** → **Games** page.

Via the API:

```
POST /api/games/upgrade
```

> **Your changes may be overwritten.** The upgrade brings game and mod settings to the catalog
> state. If you have edited run commands or variables of the standard games, save them before
> upgrading — for example, by exporting the game to a file.

A mod is not updated if the panel finds several mods with the same name for the game — such
cases are skipped to avoid overwriting the wrong one.

The catalog addresses are set with the `GAMES_CDN_URLS` variable; they are tried in order until
one responds. See the [config.env Reference](/en/config.html).

## Pelican and Pterodactyl Features

Working with imported Pelican Eggs and Pterodactyl Eggs is only possible with
Docker and Podman [process managers](/en/daemon/process_managers.html).

You need to configure GameAP Daemon to work with one of these process managers.
To do this, when adding a new node, select the desired process manager in the "Advanced Settings" section.

![](/images/en/gameap_configure/games_import/daemon_process_manager.png)
