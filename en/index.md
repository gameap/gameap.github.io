---
title: Overview
description: "Open source game server control panel: installation, dedicated server setup, working with games, access control, plugins and API."
layout: default
lang: en
order: 1
---

GameAP is an open source panel for managing game servers and services.

## Features

* Game server and service management (start, stop, restart)
* Game server file management: editing, uploading and downloading, search within a folder, checksums (MD5, SHA-1, SHA-256, SHA-512, CRC32, CRC64), creating and extracting archives right on the dedicated server (zip, tar, gzip, bzip2, xz, zstd; 7z and rar are extraction-only)
* Game server management via RCON (send commands, view console, manage players)
* Game server resource limits (CPU, RAM)
* Typed game server settings — string, text, integer, decimal number, switch, select, password — with validation rules, predefined values and administrator-only variables
* Task scheduling (automatic restart, updates, etc.)
* Access control (users, roles, permissions) and granting a user access to a specific server
* Single sign-on: at the request of an external system (for example, a billing panel) the panel mints a single-use ticket that logs a specific user in
* Arbitrary key-value metadata on dedicated servers — plugins tag the nodes they provision with it, for example to correlate a node with the cloud instance it runs on
* Update notifications: the home page shows administrators the installed panel version and the latest releases of GameAP and GameAP Daemon
* Interface in four languages — English, Russian, Spanish and German (Spanish and German were added in 4.5)
* API for integration with other systems and automation. API documentation is available at [openapi.gameap.io](https://openapi.gameap.io/)

Panel functionality is extended with plugins. Plugins are available in the [plugins.gameap.dev](https://plugins.gameap.dev/) catalog.
Anyone can develop and publish their own plugin (publishing goes through moderation).
Plugins can be written in any language that compiles to WASM: there is a ready-made SDK for Rust, and there are examples in Go and AssemblyScript.
Read more: [Plugins](/en/plugins/index.html).

## Supported Games

The panel supports starting, stopping, and restarting absolutely any games and services.

| Game                                      | Query | Rcon | Notes                                                                        |
|-------------------------------------------|-------|------|------------------------------------------------------------------------------|
| [Minecraft](/en/tutorials/minecraft.html) | ✔     | ✔    | Many mods are supported                                                      |
| Half-Life                                 | ✔     | ✔    | All versions and popular mods are supported (Sven Co-op, HeadCrab Frenzy)    |
| [Counter-Strike](/en/tutorials/cs2.html)  | ✔     | ✔    | All versions are supported (1.6, Source, Global Offensive, Counter-Strike 2) |
| Team Fortress 2                           | ✔     | ✔    |                                                                              |
| Garry's Mod                               | ✔     | ✔    |                                                                              |
| [Quake](/en/tutorials/quake3.html)        | ✔     | ✔    |                                                                              |
| [Rust](/en/tutorials/rust.html)           | ✔     | ✔    |                                                                              |
| FiveM                                     | ✔     | ✘    | Grand Theft Auto V online mod                                                |
| [Hytale](/en/tutorials/hytale.html)       | ✘     | ✘    |                                                                              |
| Terraria                                  |       |      |                                                                              |
| San Andreas: MP                           |       |      |                                                                              |

and many more...

The panel supports importing games from other control panels such as Pterodactyl and Pelican.
You can import Pelican Eggs and Pterodactyl Eggs to quickly add pre-configured settings
and create game servers based on them.
Read more about this in the [Games Import](/en/gameap_configure/games_import.html) section.

## Automatic Panel Installation

Available for Linux and Windows. The panel is also distributed as a ready-made Docker image.

You need to run the script, and it will automatically install the necessary packages and the panel.
Installation takes just a few minutes, and after it's complete, you can start using the panel right away.

* [Panel installation on Linux](/en/install/install_on_linux.html)
* [Panel installation on Windows](/en/install/install_on_windows.html)
* [Panel installation in Docker](/en/install/install_docker.html) — the ready-made `gameap/gameap` image instead of a script
