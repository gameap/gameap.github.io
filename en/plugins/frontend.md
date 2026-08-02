---
title: Plugin frontend
layout: default
lang: en
category: Plugins
order: 343
---

* This will become a table of contents (this text will be scraped).
{:toc}

## How the frontend is embedded

A plugin may contain an interface — Vue 3 components bundled into JS and CSS and embedded into the same `.wasm` file. The panel serves the concatenated frontend bundles of all loaded plugins at `/plugins.js` and `/plugins.css` (for authenticated users only). The panel loader imports this code as a module and registers every exported `PluginDefinition` object.

A plugin frontend runs in the context of the panel's main SPA (it is not an iframe and not a web component) and uses the panel's libraries through global objects:

| Global object | Contents |
|---|---|
| `window.Vue` | Vue 3 |
| `window.VueRouter` | Vue Router |
| `window.Pinia` | Pinia |
| `window.axios` | The panel's configured axios instance (with authorization) |
| `window.NaiveUI` | Naive UI |

## The `PluginDefinition` manifest

There is no separate manifest file: the frontend manifest is the `PluginDefinition` object exported from the bundle:

| Field | Required | Description |
|---|---|---|
| `id` | Yes | The plugin identifier; must match the `id` from the backend's `PluginInfo` |
| `name` | Yes | Plugin name |
| `version` | Yes | Version (semver) |
| `apiVersion` | Yes | Frontend API version, only `'1.0'` |
| `description` | No | Description |
| `author` | No | Author |
| `routes` | No | The plugin's own pages |
| `menuItems` | No | Items in the left menu (sidebar) |
| `slots` | No | Components in the panel's built-in slots |
| `homeButtons` | No | Buttons on the home page |
| `fileEditors` | No | File editors for the file manager |
| `translations` | No | Translation dictionaries `{ en: {...}, ru: {...} }` |
| `onInit` | No | Initialization hook called when the plugin is registered |

## Integration points

| Mechanism | Where it appears |
|---|---|
| `routes` | The plugin's own pages at `/plugins/{id}/...` |
| The `server-tabs` slot | A tab on the game server page (next to "Console", "Files", etc.) |
| The `dashboard-widgets` slot | A widget on the panel home page |
| `homeButtons` | Buttons on the home page |
| `menuItems` | Sidebar items (the `servers`, `admin` and `custom` sections) |
| The `admin-user-info` slot | A block in the user information dialog (administration) |
| `fileEditors` | The file manager context menu — opening a file in the plugin's editor |

The `sidebar-sections` and `admin-pages` slots are declared in the SDK but are not integrated in the current version of the panel.

Details:

* For `server-tabs`, a permission check is available: `checkPermission: { type: 'hasServerPermissions', permissions: [...] }` — the tab is shown only if the user has all the listed permissions for the server (plugin permissions have the form `plugin:{id}:...`, for example `plugin:ezvdsxmlu6fbk:manage`).
* File editors are registered with match rules (`fileName`, `extensions`, `pathContains`, `fullPath`, `gameCode` and others): the editor with the highest specificity becomes the default editor for the file. Files larger than 1 MB are not opened by plugin editors. The editor component receives the `content`, `filePath`, `fileName`, `extension`, `pluginId`, `gameCode` and `gameName` props and emits the `save` and `close` events; saving the file to the server is done by the panel itself.

## Translations

Translations are defined with the `translations: { en: {...}, ru: {...} }` dictionaries. In the `label`, `name` and `text` fields, references to translation keys of the form `@:key` are supported — the panel substitutes the string for the current interface language.

## API access

A plugin frontend uses `window.axios` — the same instance the panel uses, with the current user's authorization. Through it you can reach:

* the panel API — for example, sending an RCON command with `POST /api/servers/{id}/rcon`, or working with files through `/api/file-manager/...`;
* the plugin's own backend at `/api/plugins/{id}/...` (the HTTP routes registered by the WASM part).

## The `@gameap/plugin-sdk` SDK

The `@gameap/plugin-sdk` npm package provides:

* TypeScript types: `PluginDefinition`, `PluginRoute`, `PluginMenuItem`, `PluginSlotComponent`, `PluginHomeButton`, `PluginFileEditor`, `PluginContext` and others;
* context hooks: `usePluginContext`, `useServer`, `useServerId`, `useServerAbilities`, `useCurrentUser`, `useIsAdmin`, `useIsAuthenticated`, `usePluginRoute`, `usePluginId`;
* translation hooks: `usePluginTrans`, `providePluginTrans`;
* the panel's UI components (re-exported from `@gameap/ui`): `GCard`, `GDataTable`, `GModal`, `GStatusBadge`, `GSwitch` and others;
* `createPluginConfig` — a ready-made Vite configuration: a library-mode build (the `plugin.js` ES module) where external dependencies (`vue`, `vue-router`, `pinia`, `axios`, `@gameap/ui`) are rewritten to global objects. Note: the panel exposes `window.NaiveUI` but not `window.gameapUI`, so real plugins use their own Vite config with externals pointing at `naive-ui` (see `frontend/vite.config.js` in [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) for an example).

## Building the frontend

```bash
npm run build   # → dist/plugin.js (+ a CSS file)
```

The built `plugin.js` and CSS are embedded into the plugin's `.wasm`: for Rust, with a `build.rs` script that copies them into `OUT_DIR` and includes them via `include_bytes!` (see [Plugin development](/en/plugins/development.html)); for AssemblyScript, with a code generation script (see `scripts/embed-frontend.mjs` in plugin-minecraft-modrinth for an example).

## Local debugging

The `@gameap/debug` package starts a debugging environment with the real panel frontend and mocked API (MSW):

```bash
PLUGIN_PATH=./dist npx @gameap/debug
```

The environment opens at `http://localhost:5174`. A floating debug panel lets you switch the user type (administrator / regular user / guest), the network delay and the locale. The plugin has to be built beforehand (`npm run build`).

## A `PluginDefinition` example

An example of a tab on the game server page (following plugin-goldsrc-addons):

```ts
export const myPlugin: PluginDefinition = {
    id: 'myplugin2j7d',
    name: 'My Plugin',
    version: '0.1.0',
    apiVersion: '1.0',
    description: 'My first GameAP plugin',
    author: 'Me',
    translations: {
        en: { tab_label: 'My Plugin' },
        ru: { tab_label: 'Мой плагин' },
    },
    slots: {
        'server-tabs': [
            {
                component: MyTab,
                order: 100,
                label: '@:tab_label',
                icon: 'plug',
                name: 'my-tab',
                checkPermission: {
                    type: 'hasServerPermissions',
                    permissions: ['plugin:myplugin2j7d:manage'],
                },
            },
        ],
    },
};
```
