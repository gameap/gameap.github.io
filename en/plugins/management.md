---
title: Installation and management
layout: default
lang: en
category: Plugins
order: 341
---

## The Plugins page

Plugin management is available to an administrator on the **Administration** → **Plugins** page. The page has two tabs:

* **Installed** — the list of plugins loaded into the panel. The table shows the name, category, rating, download count, version and actions. Badges show the installation source (`file` — "Local file", `store` — "Store"), the status and whether an update is available.
* **Store** — the list of plugins from the plugins.gameap.dev catalog with server-side pagination. Already installed plugins get a badge; paid ones get a subscription purchase button.

## Installing from the catalog

1. On the **Store** tab, click the plugin name or the **Install** button — the plugin card opens with its description, author, license, tags, rating and version list.
2. Select a version from the drop-down list (the latest one is offered by default) and click **Install**.
3. The panel downloads the `.wasm` file, verifies its SHA-256 hash, writes it to the `plugins/` directory, loads the plugin and reloads the page (to refresh `/plugins.js` with the plugin frontends).

For paid plugins (marked `requires_subscription`), a subscription purchase button is offered instead of installation. To let the panel download paid plugins, set the license key in the `PLUGIN_STORE_LICENSE_KEY` environment variable.

## Installing from a file

1. On the **Installed** tab, click the **Upload** button and select the plugin's `.wasm` file (no larger than 100 MB).
2. Click **Check**: the panel performs a trial module load (a dry run, without installing) and shows the plugin metadata — name, version, author, Plugin API version, whether HTTP routes and a frontend are present, and a list of validation errors.
3. If the plugin is valid, click **Install** — the panel installs the plugin and reloads the page.

## Updating and removing

When a new version of a plugin installed from the catalog is released, an update badge appears in the table and the version column shows `installed → latest`. The update is performed with the button in the list or from the plugin card: the panel downloads the new version, verifies the hash, replaces the file and reloads the plugin.

Removal is performed with the **Remove** button, with confirmation. The panel unloads the plugin from the runtime and deletes the `.wasm` file and the plugin record. Note that the plugin's key-value storage (the `plugin_storage` table) is **not** cleared on removal — on a repeat installation the plugin will see its previous settings.

After installation, updating and removal, the page reloads in order to refresh the plugin frontends (`/plugins.js`).

## Enabling and disabling

The current version of the panel has no separate "enable/disable" action for a plugin in the interface: a plugin works from the moment it is installed until it is removed.

If a plugin stops responding (exceeds the call timeout), the panel temporarily disables it — it stops receiving events and HTTP requests until the panel is restarted. The status is shown as a badge in the plugin list.

To disable all plugins completely, set the `PLUGINS_DISABLED=true` environment variable and restart the panel.

## Permissions

A plugin can register its own game server permissions of the form `plugin:{id}:{ability}`, for example `plugin:ezvdsxmlu6fbk:manage`. These permissions appear in the game server permission list alongside the built-in ones and are granted to users in the usual way (configuring a user's permissions for a server). Administrators receive all plugin permissions automatically.

Based on these permissions, a plugin hides or shows interface elements (for example, a tab on the game server page), and can also restrict access to its HTTP routes with the `requires_auth` and `admin_only` flags.

## Panel environment variables

| Variable | Default | Description |
|---|---|---|
| `PLUGINS_DISABLED` | `false` | Completely disables the plugin system: plugins are not loaded, and plugin frontends and routes are not registered |
| `PLUGINS_AUTOLOAD` | — | A comma-separated list of `.wasm` file names to register and load automatically when the panel starts |
| `PLUGIN_STORE_URL` | `https://plugins.gameap.dev/api` | Plugin catalog API address |
| `PLUGIN_STORE_LICENSE_KEY` | — | License key for downloading paid plugins from the catalog |
| `PLUGIN_HTTP_ALLOWED_SCHEMES` | `https` | Allowed schemes for outbound plugin HTTP requests |
| `PLUGIN_HTTP_BLOCK_PRIVATE_IPS` | `true` | Blocks plugin requests to private, loopback and service IPs (SSRF protection) |
| `PLUGIN_HTTP_ALLOWED_HOSTS` | — | A list of exception hosts that may be requested even with private IPs |
| `PLUGIN_HTTP_MAX_REDIRECTS` | `5` | Maximum number of redirects in outbound plugin requests |
| `PLUGIN_HTTP_MAX_TIMEOUT_SECONDS` | `30` | Maximum timeout of an outbound plugin HTTP request, in seconds |

## Manual installation through the file system

A plugin can be installed without the interface: copy the `.wasm` file into the panel's `plugins/` directory, add the file name to the `PLUGINS_AUTOLOAD` environment variable (names separated by commas) and restart the panel — the plugin will be registered and loaded at startup.
