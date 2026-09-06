---
title: Plugin frontend
layout: default
lang: en
category: Plugins
order: 343
---

## How the frontend is embedded

A plugin may contain an interface — Vue 3 components bundled into JS and CSS and embedded into the same `.wasm` file. The panel serves the concatenated frontend bundles of all loaded plugins at `/plugins.js` and `/plugins.css` (for authenticated users only). The panel loader imports this code as a module and registers every exported `PluginDefinition` object.

> **Trust boundary.** The WASM part of a plugin runs in a sandbox; the frontend does not. It is ordinary
> JavaScript running in the same origin as the panel: it has access to the whole page and to
> `window.axios` with the permissions of the current user. A plugin with a frontend effectively gets the
> same capabilities as the panel interface itself, so install only plugins whose source you trust.

A plugin frontend runs in the context of the panel's main SPA (it is not an iframe and not a web component) and uses the panel's libraries through global objects:

| Global object           | Contents                                                                    |
|-------------------------|-----------------------------------------------------------------------------|
| `window.Vue`            | Vue 3                                                                       |
| `window.VueRouter`      | Vue Router                                                                  |
| `window.Pinia`          | Pinia                                                                       |
| `window.axios`          | The panel's configured axios instance (with authorization)                  |
| `window.NaiveUI`        | Naive UI                                                                    |
| `window.naive`          | Alias of `window.NaiveUI`, kept for the debug harness                       |
| `window.gameapUI`       | The `@gameap/ui` component library (the panel's own UI components)          |
| `window.__gameap_vue_h` | `Vue.h`, kept for backwards compatibility with render-function plugins      |

## The `PluginDefinition` manifest

There is no separate manifest file: the frontend manifest is the `PluginDefinition` object exported from the bundle:

| Field          | Required | Description                                                                 |
|----------------|----------|-----------------------------------------------------------------------------|
| `id`           | Yes      | The plugin identifier; must match the `id` from the backend's `PluginInfo`  |
| `name`         | Yes      | Plugin name                                                                 |
| `version`      | Yes      | Version (semver)                                                            |
| `apiVersion`   | Yes      | Frontend API version, only `'1.0'`                                          |
| `description`  | No       | Description                                                                 |
| `author`       | No       | Author                                                                      |
| `routes`       | No       | The plugin's own pages                                                      |
| `menuItems`    | No       | Items in the left menu (sidebar)                                            |
| `slots`        | No       | Components in the panel's built-in slots                                    |
| `homeButtons`  | No       | Buttons on the home page                                                    |
| `fileEditors`  | No       | File editors for the file manager                                           |
| `translations` | No       | Translation dictionaries keyed by language code, `{ en: {...}, ru: {...} }` |
| `onInit`       | No       | Initialization hook (sync or returning a `Promise`), called when the plugin is registered |
| `onDestroy`    | No       | Cleanup hook with the same signature; declared for forward compatibility     |

* `onInit` may be asynchronous: the loader awaits it during registration. An exception thrown inside it is logged to the browser console (`Plugin <id> onInit failed:`) and does not abort the registration.
* `onDestroy` is declared by the SDK but the panel does not call it in 4.5.0 — do not rely on it for teardown.
* A definition with an `apiVersion` other than `'1.0'` is rejected with `Unsupported API version`; the error is logged and the plugin is not registered.

## Integration points

| Mechanism     | Where it appears                                                                                   |
|---------------|----------------------------------------------------------------------------------------------------|
| `routes`      | The plugin's own pages at `/plugins/{id}/...`                                                      |
| `menuItems`   | Sidebar items (the `servers`, `admin` and `custom` sections)                                       |
| `homeButtons` | Buttons on the home page, next to **Servers** and **Nodes**                                        |
| `slots`       | Components in the 20 built-in slots of the panel (see below)                                       |
| `fileEditors` | The file manager: the context menu and the double-click default — opening a file in the plugin's editor |

### Slots

| Slot                       | Where it renders                                                                        | Props                    |
|----------------------------|-----------------------------------------------------------------------------------------|--------------------------|
| `server-tabs`              | An extra tab on the game server page, next to **Control**, **Files** and the others     | `ServerTabProps`         |
| `server-control-buttons`   | The **Start** / **Stop** / **Restart** button row on the server page                    | `ServerControlProps`     |
| `server-control-blocks`    | The **Control** tab of the server page, between the status card and the console         | `ServerControlProps`     |
| `servers-list-actions`     | The **Commands** column of the server list                                              | `ServersListActionProps` |
| `dashboard-widgets`        | The bottom of the home page, below the buttons and the information blocks               | `DashboardWidgetProps`   |
| `home-buttons`             | The home page button row (a fully custom button, see [Home buttons](#home-buttons))     | registered `props` + `pluginId` |
| `navbar-items`             | The top bar, left of the theme switch                                                   | `ChromeSlotProps` (`routeName` is not passed here) |
| `sidebar-sections`         | The sidebar, below the menu sections, in whichever sidebar variant is on screen         | `SidebarSectionProps`    |
| `global-banners`           | Above the content of every page                                                         | `ChromeSlotProps`        |
| `admin-pages`              | Above the content of `/admin/*` pages, administrators only                              | `ChromeSlotProps`        |
| `profile-info-rows`        | Extra rows in the profile table                                                         | `ProfileSlotProps`       |
| `profile-blocks`           | The profile page, below the **Two-Factor Authentication** card                          | `ProfileSlotProps`       |
| `admin-user-info-above`    | The user information dialog (administration), above the details table                   | `AdminUserInfoProps`     |
| `admin-user-info-rows`     | Extra rows in the user information dialog table                                         | `AdminUserInfoProps`     |
| `admin-user-info`          | The user information dialog, below the details table                                    | `AdminUserInfoProps`     |
| `admin-user-edit-blocks`   | The user edit page, below the **Servers** card                                          | `AdminUserEditBlockProps` |
| `admin-node-edit-blocks`   | The node edit page, at the end of the **Main** tab                                      | `AdminNodeEditBlockProps` |
| `admin-server-edit-blocks` | The server edit page, below the last card                                               | `AdminServerEditBlockProps` |
| `admin-game-edit-blocks`   | The game edit page, below the tabs                                                      | `AdminGameEditBlockProps` |
| `admin-mod-edit-blocks`    | The mod edit page, below the tabs                                                       | `AdminModEditBlockProps` |

### Registering a slot component

Each entry in `slots[slotName]` is a `PluginSlotComponent`:

| Field             | Description                                                                                                   |
|-------------------|---------------------------------------------------------------------------------------------------------------|
| `component`       | The Vue component to render (required)                                                                        |
| `order`           | Sort order within the slot, ascending (default `0`)                                                           |
| `label`           | Display label (a tab caption); supports `@:key`                                                               |
| `icon`            | Icon name from the `@gameap/ui` icon registry (for example `'plug'`, `'metrics'`); a legacy Font Awesome class still renders but logs a deprecation warning |
| `name`            | Unique name within the slot; for `server-tabs` it becomes part of the tab key `plugin-{pluginId}-{name}`      |
| `props`           | Default props for the component. They are merged under the host's context (the host wins on key collisions) in every slot except `server-tabs`, which passes only `serverId`, `server` and `pluginId` and drops registered `props` |
| `checkPermission` | Permission condition, see below                                                                               |
| `checkGame`       | Game condition, see below                                                                                     |

Registering into an unknown slot name does nothing except an `Unknown slot: <name>` warning in the console.

### `checkPermission` and `checkGame`

* `checkPermission: { type: 'hasServerPermissions', permissions: [...] }` — the component is rendered only if the user has **all** the listed permissions for the server (plugin permissions have the form `plugin:{id}:...`, for example `plugin:ezvdsxmlu6fbk:manage`).
* `checkGame: { engines?: string[], codes?: string[] }` — the component is rendered if the server's game matches at least one listed engine (case-insensitive comparison with `game.engine`) **or** one listed code (exact match with `game.code`). An empty `checkGame` matches every game; while the server's game data has not loaded yet, nothing is rendered. Example: `checkGame: { engines: ['GoldSource'], codes: ['cstrike', 'valve'] }`.
* Both conditions are independent; when both are set, both must pass.
* A failed check means "not rendered", not "disabled".

> These conditions are evaluated by exactly four slots — `server-tabs`, `server-control-buttons`,
> `server-control-blocks` and `servers-list-actions` — and, for `checkPermission` only, by file editors
> (editors restrict games through the `gameCode` / `gameName` match rules instead).
> Every other slot renders each registered component unconditionally, so never use these fields to hide
> sensitive content in, for example, `admin-user-info` or the `admin-*-edit-blocks` slots.

### Slot props

| Interface                   | Slots                                                          | Props                                                                      |
|-----------------------------|----------------------------------------------------------------|----------------------------------------------------------------------------|
| `ServerTabProps`            | `server-tabs`                                                  | `serverId`, `server`, `pluginId`                                           |
| `ServerControlProps`        | `server-control-buttons`, `server-control-blocks`              | `serverId`, `server`, `abilities: Record<string, boolean>`, `pluginId`     |
| `ServersListActionProps`    | `servers-list-actions`                                         | `serverId`, `server`, `pluginId`                                           |
| `DashboardWidgetProps`      | `dashboard-widgets`                                            | `isAdmin`, `pluginId`                                                      |
| `ChromeSlotProps`           | `navbar-items`, `global-banners`, `admin-pages`                | `routeName`, `isAdmin`, `pluginId`                                         |
| `SidebarSectionProps`       | `sidebar-sections`                                             | `minimized`, `isAdmin`, `pluginId`                                         |
| `ProfileSlotProps`          | `profile-info-rows`, `profile-blocks`                          | `userId`, `user`, `pluginId`                                               |
| `AdminUserInfoProps`        | `admin-user-info-above`, `admin-user-info-rows`, `admin-user-info` | `userId`, `user`, `pluginId`                                           |
| `AdminUserEditBlockProps`   | `admin-user-edit-blocks`                                       | `userId`, `user`, `form: UserEditFormData`, `pluginId`                     |
| `AdminNodeEditBlockProps`   | `admin-node-edit-blocks`                                       | `nodeId`, `form`, `pluginId`                                               |
| `AdminServerEditBlockProps` | `admin-server-edit-blocks`                                     | `serverId`, `server` (`AdminServerSavedData` or `null`), `form`, `pluginId` |
| `AdminGameEditBlockProps`   | `admin-game-edit-blocks`                                       | `gameCode`, `form`, `pluginId`                                             |
| `AdminModEditBlockProps`    | `admin-mod-edit-blocks`                                        | `modId`, `form`, `pluginId`                                                |

### Edit-page blocks

The `form` object handed to the `admin-*-edit-blocks` slots is a snapshot of the **unsaved** form and never carries secrets. The panel does not save plugin data together with the form — a plugin persists its own state through its own API.

| Slot                       | `form` snapshot                                                                                                  |
|----------------------------|------------------------------------------------------------------------------------------------------------------|
| `admin-user-edit-blocks`   | `login`, `name`, `email`, `roles`, `servers` (`{ id, name }[]`); never the password fields                        |
| `admin-node-edit-blocks`   | `name`, `enabled`, `os`, `location`, `provider`, `workPath`, `steamcmdPath`, `ip`; never the daemon credentials, certificates or control scripts |
| `admin-server-edit-blocks` | An allowlist of the form: `name`, `enabled`, `blocked`, `ip`, `serverPort`, `queryPort`, `rconPort`, `dir`, `user`, `startCommand`, `nodeId`, `game`, `gameMod`, `metadata` — never the RCON password; `server` is an allowlist of the saved record (`id`, `uuid`, `uuid_short`, `name`, `enabled`, `installed`, `blocked`, `online`, `ds_id`, `game_id`, `game_mod_id`) that excludes the RCON password too |
| `admin-game-edit-blocks`   | The whole form                                                                                                   |
| `admin-mod-edit-blocks`    | The whole form                                                                                                   |

### Rendering conventions and per-slot notes

* `profile-info-rows` and `admin-user-info-rows` components are rendered directly inside a `<tbody>`: the root element **must** be a `<tr>` with two `<td>` cells, otherwise the table markup is invalid.
* Block slots render components as-is, without a wrapper. A plugin that wants to look like a built-in section renders its own `GCard` / `n-card`.
* `servers-list-actions` renders nothing for disabled or blocked servers (the whole commands column is empty for them); its components are created with `h()` rather than through `PluginSlot`, and the conditions are evaluated per row.
* `sidebar-sections` is rendered in whichever sidebar variant is on screen: the collapsed one (`minimized: true`, icons only) or the expanded one (`minimized: false`). The two variants are separate branches of the markup, so collapsing or expanding the sidebar destroys the component and creates it again — any state it holds is lost. The sidebar itself is not rendered at all below the `sm` breakpoint. For plain navigation links prefer `menuItems`.
* `navbar-items` sits in a 4rem-tall bar shared with the help and profile dropdowns — keep items small.
* `admin-pages` is rendered only for administrators and only on `/admin/*` paths; `global-banners` on every page for every user.

### Routes and menu items

A `PluginRoute` has `path`, `name`, `component` and an optional `meta`:

* `path` is relative to `/plugins/{pluginId}/` and must start with `/` (the panel concatenates it as `/plugins/{pluginId}{path}`);
* `name` is prefixed to `plugin.{pluginId}.{name}`, so `route: { name: 'index' }` in a menu item or home button resolves to `plugin.my-plugin.index`; a route with no `name` becomes `plugin.{pluginId}.index`;
* `meta` accepts `title` (the panel sets it as the document title), `requiresAuth` (defaults to `true`; the navigation guard requires authentication for every non-guest route anyway), `requiresAdmin` and arbitrary keys. The guard does not act on `requiresAdmin` in 4.5.0 — a page meant for administrators has to check `useIsAdmin()` itself;
* `children` is declared by the type, but the panel registers only the top-level route and drops the nested ones — declare every page as a route of its own.

A `PluginMenuItem` has `section` (`'servers'`, `'admin'` or `'custom'`; an unknown section falls back to `custom`), `icon` (registry name, default `'puzzle-piece'`), `text` (supports `@:key`), `route` (`{ name }` or `{ path }`, normalized the same way) and `order` (default `100`). The `adminOnly` field is declared in the type but **not** honoured in 4.5.0 — put admin-only items into `section: 'admin'`, which is rendered only for administrators.

### Home buttons

A `PluginHomeButton` has:

| Field       | Description                                                                                                       |
|-------------|-------------------------------------------------------------------------------------------------------------------|
| `name`      | Button caption (required); supports `@:key`                                                                       |
| `icon`      | Icon name from the `@gameap/ui` registry, for example `'metrics'`; when omitted the panel uses `'puzzle-piece'`. A name unknown to the registry is rendered as a bare `<i class="...">` |
| `component` | A custom Vue component rendered instead of the default button; it receives the normalized `route` and `pluginId` as props |
| `route`     | `{ name }` or `{ path }`, normalized like menu item routes; defaults to the plugin's index route `plugin.{pluginId}.index` |
| `order`     | Sort order (lower first)                                                                                          |

The same result can be achieved by registering a component into the `home-buttons` slot directly.

## File editors

A plugin registers file editors with `fileEditors: PluginFileEditor[]`. The editor appears in the file manager context menu and, unless it is context-menu-only, becomes the double-click default for matching files.

| Field             | Required | Description                                                                                          |
|-------------------|----------|------------------------------------------------------------------------------------------------------|
| `id`              | Yes      | Unique identifier within the plugin                                                                  |
| `name`            | Yes      | Display name — the modal title and the default menu caption `Edit with <name>`; supports `@:key`     |
| `component`       | Yes      | The Vue component that renders the editor                                                            |
| `match`           | Yes      | Match rules (see below)                                                                              |
| `contentType`     | No       | `'text'` (default), `'binary'` or `'none'`                                                           |
| `readOnly`        | No       | Read-only editor: the modal shows no save button                                                     |
| `icon`            | No       | Icon name from the `@gameap/ui` registry, for example `'file-archive'`                               |
| `contextMenuOnly` | No       | Offer the editor in the context menu only, never as the double-click default. **Mandatory** for an editor matching `allFiles`, otherwise it takes over every image, video, PDF and text preview |
| `menuLabel`       | No       | Caption of the context menu item instead of `Edit with <name>`; supports `@:key`                     |
| `menuGroup`       | No       | Which block of the context menu the item joins: `'top'` (default), `'open'`, `'modify'`, `'danger'` or `'info'` |
| `checkPermission` | No       | The same `hasServerPermissions` condition as for slots; the item is hidden when it fails             |
| `width`           | No       | Modal width as a CSS length; the default is `1000px`, no `max-width` is applied — use a viewport-relative value such as `'min(1400px, 95vw)'` |
| `hideFooter`      | No       | Remove the modal footer; the editor draws its own actions, closing stays on the header cross and Escape |
| `keepOpenOnSave`  | No       | Leave the modal open after a successful save                                                         |

### Match rules and specificity

`match` may combine `allFiles`, `fileName` (exact), `pathContains`, `fullPath` (exact), `extensions` (array, case-insensitive), `fileNameRegexp`, `gameCode` and `gameName`. All specified rules must match. When several editors match a file, all of them are listed in the context menu and the most specific one is marked `(default)` and opens on double click. Specificity: `fullPath` > `pathContains` > `fileName` > `fileNameRegexp` > `extensions` > `allFiles`; `gameCode` and `gameName` add to the score. An editor with `contextMenuOnly: true` is never the default.

### The size cap and `contentType: 'none'`

The modal downloads the file before mounting the editor, so files larger than 1 MB (1 048 576 bytes) are not opened by ordinary plugin editors — the context menu item stays visible but disabled. An editor that declares `contentType: 'none'` receives no `content` prop at all, so the cap does not apply to it: it opens on files of any size and decides itself what to read, using `fileSize` and `fileMtime` from the directory listing. Saving still works for such an editor.

### Where the item lands in the context menu

| `menuGroup` | What the block holds                                    | Where the item lands            |
|-------------|---------------------------------------------------------|---------------------------------|
| `top`       | Nothing of the file manager's — a block of its own above the rest | in plugin order       |
| `open`      | Open, Play, View, Edit, Select, Download, Download as ZIP, Zip, Unzip | after the block's own items |
| `modify`    | Copy, Cut, Rename, Permissions, Paste                   | after the block's own items     |
| `danger`    | Delete                                                  | after the block's own items     |
| `info`      | Checksums, Properties                                   | **before** the block's own items, so that Properties stays last |

An unknown `menuGroup` value falls back to `top`.

### The editor component

The component receives the `FileEditorProps`:

| Prop        | Description                                                                              |
|-------------|------------------------------------------------------------------------------------------|
| `content`   | File content — a string for `text`, an `ArrayBuffer` for `binary`; absent for `none`     |
| `filePath`  | Full file path                                                                           |
| `fileName`  | File name with extension                                                                 |
| `extension` | Extension without the dot                                                                |
| `gameCode`  | The server's game code, if available                                                     |
| `gameName`  | The server's game name, if available                                                     |
| `fileSize`  | Size in bytes, as reported by the directory listing                                      |
| `fileMtime` | Modification time in Unix seconds, as reported by the directory listing                  |
| `disk`      | The file manager disk the file lives on                                                  |
| `pluginId`  | The plugin that registered the editor                                                    |

The server the file belongs to is deliberately not a prop — read it with `useServerId()` / `useServer()` from the SDK (both throw outside a plugin context).

The component emits `save` (payload: the new content, a string or an `ArrayBuffer`) and `close`; the panel writes the file to the server itself. The footer's save button calls the component's exposed `save()` method, so the editor should `defineExpose({ save })` and emit `save` from it. After a successful write the panel calls the editor's exposed `onSaved()` if it exists and, unless `keepOpenOnSave` is set, closes the modal.

The modal caps the editor body at `--gameap-plugin-editor-height` (`calc(100vh - 250px)` in 4.5.0) and scrolls anything taller; there is no `height` field. An editor with its own scroller should size it from that variable, keeping a fallback for older panels, otherwise the modal ends up with two scrollbars and the editor's buttons below the fold:

```css
.my-editor-pane {
    height: calc(var(--gameap-plugin-editor-height, calc(100vh - 250px)) - 10rem);
    overflow: auto;
}
```

## Translations

Translations are defined with the `translations` dictionaries keyed by language code, for example `{ en: {...}, ru: {...} }`; any code is accepted — the panel itself ships `de`, `en`, `es` and `ru`. Lookup order: the current interface language → `en` → the key itself.

* `trans(key, params)` substitutes `:paramName` placeholders: `trans('greeting', { name: 'World' })` with `'Hello, :name!'`.
* `usePluginTrans()` is for plugin route pages, where the translation context is provided automatically.
* `providePluginTrans(props.pluginId)` must be called in the root component of a slot component (tab, widget, block); its children can then use `usePluginTrans()`.
* References of the form `@:key` are supported in `label` (slots), `text` (menu items), `name` (home buttons) and a file editor's `name` and `menuLabel`. When the plugin has no such key, the panel falls back to its own translation table.

## API access

A plugin frontend uses `window.axios` — the same instance the panel uses, with the current user's authorization. Through it you can reach:

* the panel API — for example, sending an RCON command with `POST /api/servers/{id}/rcon`, or working with files through `/api/file-manager/...`;
* the plugin's own backend at `/api/plugins/{id}/...` (the HTTP routes registered by the WASM part).

## The `@gameap/plugin-sdk` SDK

The `@gameap/plugin-sdk` npm package (version `0.3.3` in the panel repository at 4.5.0; peer dependencies `vue ^3.5` and `@gameap/ui ^1.3`) provides:

* TypeScript types: `PluginDefinition`, `PluginRoute`, `PluginMenuItem`, `PluginSlotComponent`, `PluginHomeButton`, `PluginFileEditor`, `SlotName`, `PluginContext`, `PluginRouteInfo`, `ServerData`, `UserData`, `PluginI18nContext`; the slot props interfaces listed above plus `UserEditFormData` and `AdminServerSavedData`; the editor types `EditorContentType`, `EditorMenuGroup`, `EditorMatchRules`, `FileEditorProps`;
* context hooks: `usePluginContext`, `useServer`, `useServerId`, `useServerAbilities`, `useCurrentUser`, `useIsAdmin`, `useIsAuthenticated`, `usePluginRoute`, `usePluginId`;
* translation hooks: `usePluginTrans`, `providePluginTrans`;
* Vue helpers re-exported for convenience: `defineComponent`, `ref`, `computed`, `watch`, `onMounted`, `onUnmounted`;
* the panel's UI components (re-exported from `@gameap/ui`): `GBreadcrumbs`, `GCard`, `GDataTable`, `GDeletableList`, `GDivider`, `GEmpty`, `GGameIcon`, `GIcon`, `GInput`, `GMenu`, `GMenuButton`, `GMenuItem`, `GMenuItems`, `GModal`, `GStatusBadge`, `GSwitch`, `GTable`, `Loading`, `Progressbar`, and the icon registry helpers `registerIcons`, `iconRegistry`, `getIcon`, `hasIcon`, `defaultIconMap`;
* `createPluginConfig` (imported from `@gameap/plugin-sdk/vite`) — a ready-made Vite configuration: a library-mode build (the `plugin.js` ES module) where the external dependencies `vue`, `vue-router`, `pinia`, `axios` and `@gameap/ui` are rewritten to the global objects listed above (`@gameap/ui` → `window.gameapUI`). A plugin needs its own Vite config only if it imports `naive-ui` directly — then externalize `naive-ui` to `window.NaiveUI` (see `frontend/vite.config.js` in [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) for an example).

## Theming

The panel's colours are CSS custom properties with the `--gameap-` prefix, defined in `/theme.css` (the `@gameap/ui` `theme.css` file, served verbatim under a stable name). The variable table is in the `@gameap/ui` README ("Theming").

Load order: `/theme.css` (a `<link>` in the document head) → the panel CSS (`main-*.css`) → the plugin CSS (`<style id="gameap-plugin-styles">`, injected before the app mounts) → lazy route chunks such as the file manager, loaded whenever their route opens. Consequently an override of a `--gameap-*` variable always wins, while an override of a panel selector is unreliable — a chunk loaded later out-cascades it at equal specificity. **Override variables, not selectors.**

Ship regular CSS in the plugin bundle; dark values must be declared under `html.dark` (the panel toggles the `dark` class on `<html>`), a plain `:root` declaration loses to the panel's own `html.dark` rules:

```css
:root {
    --gameap-primary: #e11d48;
    --gameap-primary-hover: #be123c;
}

html.dark {
    --gameap-surface: #1e1b4b;
}
```

Limitations:

* `/plugins.css` is served behind authentication, so plugin theme overrides do not apply to the login screen.
* Naive UI components are CSS-in-JS: the panel rebuilds their theme from the variables on load and on every light/dark switch, but only for the variables it maps — `--gameap-{primary,success,warning,danger}` and their `-hover` variants, `--gameap-table-header`, `--gameap-surface-overlay`, `--gameap-surface-raised`, `--gameap-surface-hover`, `--gameap-text-muted`, `--gameap-tab-accent`.

In plugin templates prefer the safelisted semantic utility classes (`bg-surface`, `text-muted`, `bg-primary`, `text-danger`, `border-strong` and the others listed in the `@gameap/ui` README) — they follow the active theme and need no `dark:` variants. In plugin CSS use `var(--gameap-*)` directly. Other Tailwind classes exist only if the panel's own markup happens to use them, so do not rely on arbitrary utilities.

## Building the frontend

```bash
npm run build   # → dist/plugin.js (+ a CSS file)
```

The built `plugin.js` and CSS are embedded into the plugin's `.wasm`: for Rust, with a `build.rs` script that copies them into `OUT_DIR` and includes them via `include_bytes!` (see [Plugin development](/en/plugins/development.html)); for AssemblyScript, with a code generation script (see `scripts/embed-frontend.mjs` in plugin-minecraft-modrinth for an example).

## Local debugging

The `@gameap/debug` package starts a debugging environment with the real panel frontend and mocked API (MSW):

```bash
PLUGINS_PATH=./dist npx @gameap/debug
```

The environment opens at `http://localhost:5174`. A floating debug panel lets you switch the user type (administrator / regular user / guest), the network delay and the locale. The plugin has to be built beforehand (`npm run build`).

## A `PluginDefinition` example

A tab on the game server page shown only for GoldSource games and only to users with the plugin permission, a button in the server control row, and a context-menu-only file inspector that loads its own content (following plugin-goldsrc-addons):

```ts
export const myPlugin: PluginDefinition = {
    id: 'myplugin2j7d',
    name: 'My Plugin',
    version: '0.1.0',
    apiVersion: '1.0',
    description: 'My first GameAP plugin',
    author: 'Me',
    translations: {
        en: { tab_label: 'My Plugin', versions_title: 'File versions', versions_menu: 'Show versions' },
        ru: { tab_label: 'Мой плагин', versions_title: 'Версии файла', versions_menu: 'Показать версии' },
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
                checkGame: {
                    engines: ['GoldSource'],
                    codes: ['cstrike', 'valve'],
                },
            },
        ],
        'server-control-buttons': [
            {
                component: BackupButton,
                order: 10,
                checkPermission: {
                    type: 'hasServerPermissions',
                    permissions: ['plugin:myplugin2j7d:manage'],
                },
            },
        ],
    },
    fileEditors: [
        {
            id: 'file-versions',
            name: '@:versions_title',
            menuLabel: '@:versions_menu',
            menuGroup: 'info',
            component: FileVersions,
            match: { allFiles: true },
            contentType: 'none',
            contextMenuOnly: true,
            readOnly: true,
            checkPermission: {
                type: 'hasServerPermissions',
                permissions: ['plugin:myplugin2j7d:manage'],
            },
        },
    ],
};
```
