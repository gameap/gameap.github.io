---
title: Overview
layout: default
lang: en
order: 1
---

GameAP is an open source panel for managing game servers and services.

## Features

* Game server and service management (start, stop, restart)
* Game server file management (edit, upload, download)
* Game server management via RCON (send commands, view console, manage players)
* Game server resource limits (CPU, RAM)
* Task scheduling (automatic restart, updates, etc.)
* Access control (users, roles, permissions)
* API for integration with other systems and automation. API documentation is available at [openapi.gameap.io](https://openapi.gameap.io/)

Panel functionality is extended with plugins. Plugins are available in the [plugins.gameap.dev](https://plugins.gameap.dev/) catalog.
Anyone can develop and publish their own plugin (publishing goes through moderation).
Plugins can be written in any language that compiles to WASM: there is a ready-made SDK for Rust, and there are examples in Go and AssemblyScript.
Read more: [Plugins](/en/plugins/index.html).

## Supported Games

The panel supports starting, stopping, and restarting absolutely any games and services.

| Game                                      | Query | Rcon | Notes                                                                      |
|-------------------------------------------|-------|------|----------------------------------------------------------------------------|
| [Minecraft](/en/tutorials/minecraft.html) | ✔     | ✔    | Many mods are supported                                                    |
| Half-Life                                 | ✔     | ✔    | All versions and popular mods are supported (Sven Co-op, HeadCrab Frenzy)  |
| [Counter-Strike](/en/tutorials/cs2.html)  | ✔     | ✔    | All versions are supported (1.6, Source, Global Offensive, Counter-Strike 2) |
| Team Fortress 2                           | ✔     | ✔    |                                                                            |
| Garry's Mod                               | ✔     | ✔    |                                                                            |
| [Quake](/en/tutorials/quake3.html)        | ✔     | ✔    |                                                                            |
| [Rust](/en/tutorials/rust.html)           | ✔     | ✔    |                                                                            |
| FiveM                                     | ✔     | ✘    | Grand Theft Auto V online mod                                              |
| [Hytale](/en/tutorials/hytale.html)       | ✘     | ✘    |                                                                            |
| Terraria                                  |       |      |                                                                            |
| San Andreas: MP                           |       |      |                                                                            |

and many more...

The panel supports importing games from other control panels such as Pterodactyl and Pelican.
You can import Pelican Eggs and Pterodactyl Eggs to quickly add pre-configured settings
and create game servers based on them.
Read more about this in the [Games Import](/en/gameap_configure/games_import.html) section.

## Automatic Panel Installation

Available for Linux and Windows.

You need to run the script, and it will automatically install the necessary packages and the panel.
Installation takes just a few minutes, and after it's complete, you can start using the panel right away.

* [Panel installation on Linux](/en/install/install_on_linux.html)
* [Panel installation on Windows](/en/install/install_on_windows.html)
