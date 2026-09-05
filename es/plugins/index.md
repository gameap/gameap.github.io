---
title: Plugins
layout: default
lang: es
category: Plugins
order: 340
---

Los plugins amplían la funcionalidad del panel GameAP: añaden nuevas páginas, pestañas en la página del servidor de juego, editores de archivos, botones en la página principal e integraciones con servicios externos.

Un plugin es un único archivo `.wasm` (un módulo WASM compilado para `wasm32-wasip1`). Los plugins pueden escribirse en cualquier lenguaje que compile a WASM: existe un SDK listo para Rust y es posible desarrollar en AssemblyScript. Un plugin puede incluir un frontend en Vue 3 integrado en el mismo archivo `.wasm`: las páginas y los componentes del plugin se ejecutan directamente dentro de la interfaz del panel.

![Editor hexadecimal para archivos del servidor de juego — un ejemplo de plugin de GameAP](/images/ru/plugins/hex-editor.png)

*Un editor hexadecimal para archivos del servidor de juego: un ejemplo de plugin en funcionamiento.*

## Seguridad

Los plugins se ejecutan en un entorno aislado (un runtime WASM): no tienen acceso directo al sistema de archivos ni a la red, y cada llamada pasa por una interfaz controlada del panel. Solo un administrador del panel puede instalar y eliminar plugins. Al instalar desde el catálogo, el panel verifica el hash SHA-256 del archivo descargado.

## Catálogo de plugins

El catálogo oficial de plugins está disponible en [plugins.gameap.dev](https://plugins.gameap.dev/) (versión en ruso: [plugins.gameap.ru](https://plugins.gameap.ru/)). Allí pueden publicarse tanto plugins del equipo de GameAP como plugins de desarrolladores externos: cualquiera puede registrarse y publicar su propio plugin (consulte [Publicación en el catálogo](/es/plugins/publishing.html)).

Los plugins del catálogo se instalan desde la interfaz del panel en unos pocos clics (consulte [Instalación y administración](/es/plugins/management.html)).

### Plugins oficiales

| Plugin | Descripción | Catálogo | Repositorio |
|---|---|---|---|
| FTP (files) | Administración de un servidor FTP(S)/SFTP en servidores dedicados (nodos) | [plugins.gameap.dev/plugins/files](https://plugins.gameap.dev/plugins/files) | [github.com/gameap/plugin-files](https://github.com/gameap/plugin-files) |
| HEX Editor | Visualización y edición de archivos del servidor de juego en formato hexadecimal | [plugins.gameap.dev/plugins/hexeditor4jm2](https://plugins.gameap.dev/plugins/hexeditor4jm2) | [github.com/gameap/plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) |
| GoldSource Addons | Administración de plugins de Metamod y AMX Mod X en servidores GoldSource (Half-Life, CS 1.6, etc.) | [plugins.gameap.dev/plugins/ezvdsxmlu6fbk](https://plugins.gameap.dev/plugins/ezvdsxmlu6fbk) | [github.com/gameap/plugin-goldsrc-addons](https://github.com/gameap/plugin-goldsrc-addons) |
| Minecraft Modrinth | Búsqueda, instalación y actualización de mods y plugins de Minecraft desde modrinth.com | [plugins.gameap.dev/plugins/dshdabjp2l73a](https://plugins.gameap.dev/plugins/dshdabjp2l73a) | [github.com/gameap/plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth) |

## En esta sección

* [Instalación y administración](/es/plugins/management.html) — instalación de plugins desde el catálogo y desde un archivo, actualización, eliminación, permisos, variables de entorno.
* [Desarrollo de plugins](/es/plugins/development.html) — arquitectura del sistema de plugins, la interfaz del plugin, eventos, rutas HTTP, funciones del host, compilación con Rust.
* [Frontend de plugins](/es/plugins/frontend.html) — integración de la interfaz del plugin en el panel, el manifiesto `PluginDefinition`, slots, SDK, depuración local.
* [Publicación en el catálogo](/es/plugins/publishing.html) — registro en el catálogo, publicación de versiones, moderación, publicación desde CI/CD.
