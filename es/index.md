---
title: Descripción general
layout: default
lang: es
order: 1
---

GameAP es un panel de código abierto para administrar servidores de juego y servicios.

## Características

* Administración de servidores de juego y servicios (iniciar, detener, reiniciar)
* Gestión de archivos del servidor de juego (editar, subir, descargar)
* Administración del servidor de juego mediante RCON (enviar comandos, ver la consola, gestionar jugadores)
* Límites de recursos del servidor de juego (CPU, RAM)
* Planificación de tareas (reinicio automático, actualizaciones, etc.)
* Control de acceso (usuarios, roles, permisos)
* API para la integración con otros sistemas y la automatización. La documentación de la API está disponible en [openapi.gameap.io](https://openapi.gameap.io/)

La funcionalidad del panel se amplía con plugins. Los plugins están disponibles en el catálogo [plugins.gameap.dev](https://plugins.gameap.dev/).
Cualquiera puede desarrollar y publicar su propio plugin (la publicación pasa por moderación).
Los plugins pueden escribirse en cualquier lenguaje que compile a WASM: existe un SDK listo para Rust, y hay ejemplos en Go y AssemblyScript.
Más información: [Plugins](/es/plugins/index.html).

## Juegos compatibles

El panel permite iniciar, detener y reiniciar absolutamente cualquier juego y servicio.

| Juego                                     | Query | Rcon | Notas                                                                      |
|-------------------------------------------|-------|------|----------------------------------------------------------------------------|
| [Minecraft](/es/tutorials/minecraft.html) | ✔     | ✔    | Se admiten muchos mods                                                     |
| Half-Life                                 | ✔     | ✔    | Se admiten todas las versiones y mods populares (Sven Co-op, HeadCrab Frenzy) |
| [Counter-Strike](/es/tutorials/cs2.html)  | ✔     | ✔    | Se admiten todas las versiones (1.6, Source, Global Offensive, Counter-Strike 2) |
| Team Fortress 2                           | ✔     | ✔    |                                                                            |
| Garry's Mod                               | ✔     | ✔    |                                                                            |
| [Quake](/es/tutorials/quake3.html)        | ✔     | ✔    |                                                                            |
| [Rust](/es/tutorials/rust.html)           | ✔     | ✔    |                                                                            |
| FiveM                                     | ✔     | ✘    | Mod en línea de Grand Theft Auto V                                         |
| [Hytale](/es/tutorials/hytale.html)       | ✘     | ✘    |                                                                            |
| Terraria                                  |       |      |                                                                            |
| San Andreas: MP                           |       |      |                                                                            |

y muchos más...

El panel admite la importación de juegos desde otros paneles de control como Pterodactyl y Pelican.
Puede importar Pelican Eggs y Pterodactyl Eggs para añadir rápidamente configuraciones predefinidas
y crear servidores de juego basados en ellas.
Más información en la sección [Importación de juegos](/es/gameap_configure/games_import.html).

## Instalación automática del panel

Disponible para Linux y Windows.

Debe ejecutar el script, y este instalará automáticamente los paquetes necesarios y el panel.
La instalación tarda solo unos minutos y, una vez finalizada, puede empezar a usar el panel de inmediato.

* [Instalación del panel en Linux](/es/install/install_on_linux.html)
* [Instalación del panel en Windows](/es/install/install_on_windows.html)
