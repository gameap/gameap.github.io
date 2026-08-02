---
title: Плагины
layout: default
lang: ru
category: Плагины
order: 340
---

Плагины расширяют функциональность панели GameAP: добавляют новые страницы, вкладки на странице игрового сервера, редакторы файлов, кнопки на главной странице и интеграции с внешними сервисами.

Плагин — это один файл `.wasm` (WASM-модуль, таргет `wasm32-wasip1`). Плагины можно писать на любом языке, компилируемом в WASM: для Rust есть готовый SDK, возможна разработка на AssemblyScript. Плагин может включать фронтенд-часть на Vue 3, встроенную в тот же файл `.wasm` — страницы и компоненты плагина работают прямо в интерфейсе панели.

![](/images/ru/plugins/hex-editor.png)

*HEX-редактор файлов игрового сервера — пример работающего плагина.*

## Безопасность

Плагины выполняются в изолированном окружении (WASM-рантайм): у них нет прямого доступа к файловой системе и сети, все обращения идут через контролируемый интерфейс панели. Устанавливать и удалять плагины может только администратор панели. При установке из каталога панель проверяет SHA-256-хеш скачанного файла.

## Каталог плагинов

Официальный каталог плагинов находится по адресу [plugins.gameap.ru](https://plugins.gameap.ru/) (английская версия — [plugins.gameap.dev](https://plugins.gameap.dev/)). В каталоге могут публиковаться как плагины команды GameAP, так и плагины сторонних разработчиков — любой желающий может зарегистрироваться и опубликовать свой плагин (см. [Публикация в каталоге](/ru/plugins/publishing.html)).

Плагины из каталога устанавливаются из интерфейса панели в несколько кликов (см. [Установка и управление](/ru/plugins/management.html)).

### Официальные плагины

| Плагин | Описание | Каталог | Репозиторий |
|---|---|---|---|
| FTP (files) | Управление FTP(S)/SFTP-сервером на выделенных серверах (нодах) | [plugins.gameap.ru/plugins/files](https://plugins.gameap.ru/plugins/files) | [github.com/gameap/plugin-files](https://github.com/gameap/plugin-files) |
| HEX Редактор | Просмотр и редактирование файлов игрового сервера в шестнадцатеричном виде | [plugins.gameap.ru/plugins/hexeditor4jm2](https://plugins.gameap.ru/plugins/hexeditor4jm2) | [github.com/gameap/plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) |
| GoldSource Addons | Управление плагинами Metamod и AMX Mod X на GoldSource-серверах (Half-Life, CS 1.6 и др.) | [plugins.gameap.ru/plugins/ezvdsxmlu6fbk](https://plugins.gameap.ru/plugins/ezvdsxmlu6fbk) | [github.com/gameap/plugin-goldsrc-addons](https://github.com/gameap/plugin-goldsrc-addons) |
| Minecraft Modrinth | Поиск, установка и обновление модов и плагинов Minecraft с modrinth.com | [plugins.gameap.ru/plugins/dshdabjp2l73a](https://plugins.gameap.ru/plugins/dshdabjp2l73a) | [github.com/gameap/plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth) |

## В этом разделе

* [Установка и управление](/ru/plugins/management.html) — установка плагинов из каталога и из файла, обновление, удаление, права доступа, переменные окружения.
* [Разработка плагинов](/ru/plugins/development.html) — архитектура плагинной системы, интерфейс плагина, события, HTTP-маршруты, host-функции, сборка на Rust.
* [Фронтенд плагина](/ru/plugins/frontend.html) — встраивание интерфейса плагина в панель, манифест `PluginDefinition`, слоты, SDK, локальная отладка.
* [Публикация в каталоге](/ru/plugins/publishing.html) — регистрация в каталоге, публикация версий, модерация, публикация из CI/CD.
