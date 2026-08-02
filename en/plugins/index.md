---
title: Plugins
layout: default
lang: en
category: Plugins
order: 340
---

* This will become a table of contents (this text will be scraped).
{:toc}

Plugins extend the functionality of the GameAP panel: they add new pages, tabs on the game server page, file editors, buttons on the home page, and integrations with external services.

A plugin is a single `.wasm` file (a WASM module targeting `wasm32-wasip1`). Plugins can be written in any language that compiles to WASM: there is a ready-made SDK for Rust, and development in AssemblyScript is possible. A plugin may include a Vue 3 frontend embedded in the same `.wasm` file — the plugin's pages and components run directly inside the panel interface.

![](/images/ru/plugins/hex-editor.png)

*A hex editor for game server files — an example of a working plugin.*

## Security

Plugins run in an isolated environment (a WASM runtime): they have no direct access to the file system or the network, and every call goes through a controlled panel interface. Only a panel administrator can install and remove plugins. When installing from the catalog, the panel verifies the SHA-256 hash of the downloaded file.

## Plugin catalog

The official plugin catalog is available at [plugins.gameap.dev](https://plugins.gameap.dev/) (Russian version — [plugins.gameap.ru](https://plugins.gameap.ru/)). Both plugins by the GameAP team and plugins by third-party developers may be published there — anyone can register and publish their own plugin (see [Publishing to the catalog](/en/plugins/publishing.html)).

Plugins from the catalog are installed from the panel interface in a few clicks (see [Installation and management](/en/plugins/management.html)).

### Official plugins

| Plugin | Description | Catalog | Repository |
|---|---|---|---|
| FTP (files) | Managing an FTP(S)/SFTP server on dedicated servers (nodes) | [plugins.gameap.dev/plugins/files](https://plugins.gameap.dev/plugins/files) | [github.com/gameap/plugin-files](https://github.com/gameap/plugin-files) |
| HEX Editor | Viewing and editing game server files in hexadecimal | [plugins.gameap.dev/plugins/hexeditor4jm2](https://plugins.gameap.dev/plugins/hexeditor4jm2) | [github.com/gameap/plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) |
| GoldSource Addons | Managing Metamod and AMX Mod X plugins on GoldSource servers (Half-Life, CS 1.6, etc.) | [plugins.gameap.dev/plugins/ezvdsxmlu6fbk](https://plugins.gameap.dev/plugins/ezvdsxmlu6fbk) | [github.com/gameap/plugin-goldsrc-addons](https://github.com/gameap/plugin-goldsrc-addons) |
| Minecraft Modrinth | Searching, installing and updating Minecraft mods and plugins from modrinth.com | [plugins.gameap.dev/plugins/dshdabjp2l73a](https://plugins.gameap.dev/plugins/dshdabjp2l73a) | [github.com/gameap/plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth) |

## In this section

* [Installation and management](/en/plugins/management.html) — installing plugins from the catalog and from a file, updating, removing, permissions, environment variables.
* [Plugin development](/en/plugins/development.html) — plugin system architecture, the plugin interface, events, HTTP routes, host functions, building with Rust.
* [Plugin frontend](/en/plugins/frontend.html) — embedding the plugin interface into the panel, the `PluginDefinition` manifest, slots, SDK, local debugging.
* [Publishing to the catalog](/en/plugins/publishing.html) — registering in the catalog, publishing versions, moderation, publishing from CI/CD.
