---
layout: default
lang: en
order: 1
---

GameAP is an open source dashboard for managing game servers and services.

* This will become a table of contents (this text will be scraped).
{:toc}

## Supported Games

GameAP supports starting, stopping and restarting any games and services.

| Game | Query | Rcon | Notes |
| ------ | ------- | ------ | ------- |
| Minecraft | ✔ | ✘ |
| Half-Life| ✔ | ✘ | Supported all versions and many popular mods (Sven Co-op, HeadCrab Frenzy) |
| Counter-Strike | ✔ | ✘ | Supported all versions (1.6, Source, Global Offencive) |
| Team Fortress 2 | ✔ | ✘ |
| Garry's Mod | ✔ | ✘ |
| Rust | ✔ | ✘ |
| Terraria | | 
| San Andreas: MP | |

and many more... 

## Installation methods

### Automatic installation

Available for Debian, Ubuntu, Windows.
Fully automatic installation. You need to run the script, it will automatically install the necessary packages and panel.

Difficulty: easy

[GameAP 3.x Automatic Installation](/en/auto_install.html)

### Via package manager

Availabel only for Linux. Panel installation via APT package manager.
You need to add a repository and execute the installation command.

Difficulty: medium

[GameAP 3.x installation](/en/install.html)

### Manual installation from GitHub

Installing from source codes. You will need to download the panel,
install the necessary package dependencies, build styles using NPM.

Difficulty: high

[GameAP 3.x manual installation](/ru/manual_install.html)

### Installation on Shared hosting

Installing the panel on a shared hosting, where there is no access to the command line.

Not recommended method

* [GameAP 3.x installation on Shared hosting](/ru/shared_install.html)