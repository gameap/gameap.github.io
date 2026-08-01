---
title: Фронтенд плагина
layout: default
lang: ru
category: Плагины
order: 343
---

* This will become a table of contents (this text will be scraped).
{:toc}

## Как встраивается фронтенд

Плагин может содержать интерфейс — Vue 3-компоненты, собранные в JS-бандл и CSS и встроенные в тот же файл `.wasm`. Панель отдаёт склейку фронтенд-бандлов всех загруженных плагинов по адресам `/plugins.js` и `/plugins.css` (только для аутентифицированных пользователей). Загрузчик панели импортирует этот код как модуль и регистрирует каждый экспортированный объект `PluginDefinition`.

Фронтенд плагина работает в контексте основного SPA панели (это не iframe и не web component) и использует библиотеки панели через глобальные объекты:

| Глобальный объект | Содержимое |
|---|---|
| `window.Vue` | Vue 3 |
| `window.VueRouter` | Vue Router |
| `window.Pinia` | Pinia |
| `window.axios` | Настроенный axios-инстанс панели (с авторизацией) |
| `window.NaiveUI` | Naive UI |

## Манифест `PluginDefinition`

Отдельного файла-манифеста нет: манифестом фронтенда служит экспортируемый из бандла объект `PluginDefinition`:

| Поле | Обязательное | Описание |
|---|---|---|
| `id` | Да | Идентификатор плагина; должен совпадать с `id` из `PluginInfo` бэкенда |
| `name` | Да | Название плагина |
| `version` | Да | Версия (семвер) |
| `apiVersion` | Да | Версия фронтенд API, только `'1.0'` |
| `description` | Нет | Описание |
| `author` | Нет | Автор |
| `routes` | Нет | Собственные страницы плагина |
| `menuItems` | Нет | Пункты левого меню (сайдбара) |
| `slots` | Нет | Компоненты во встроенных слотах панели |
| `homeButtons` | Нет | Кнопки на главной странице |
| `fileEditors` | Нет | Редакторы файлов для файлового менеджера |
| `translations` | Нет | Словари переводов `{ en: {...}, ru: {...} }` |
| `onInit` | Нет | Хук инициализации при регистрации плагина |

## Точки интеграции

| Механизм | Где появляется |
|---|---|
| `routes` | Собственные страницы с адресами `/plugins/{id}/...` |
| Слот `server-tabs` | Вкладка на странице игрового сервера (рядом с «Консолью», «Файлами» и др.) |
| Слот `dashboard-widgets` | Виджет на главной странице панели |
| `homeButtons` | Кнопки на главной странице |
| `menuItems` | Пункты сайдбара (секции `servers`, `admin`, `custom`) |
| Слот `admin-user-info` | Блок в окне информации о пользователе (администрирование) |
| `fileEditors` | Контекстное меню файлового менеджера — открытие файла в редакторе плагина |

Слоты `sidebar-sections` и `admin-pages` объявлены в SDK, но в текущей версии панели не интегрированы.

Особенности:

* Для `server-tabs` доступна проверка прав `checkPermission: { type: 'hasServerPermissions', permissions: [...] }` — вкладка показывается, только если у пользователя есть все перечисленные права на сервер (плагинные права имеют вид `plugin:{id}:...`, например `plugin:ezvdsxmlu6fbk:manage`).
* Редакторы файлов регистрируются с match-правилами (`fileName`, `extensions`, `pathContains`, `fullPath`, `gameCode` и др.): редактор с наибольшей специфичностью становится редактором по умолчанию для файла. Файлы больше 1 МБ плагинными редакторами не открываются. Компонент редактора получает props `content`, `filePath`, `fileName`, `extension`, `pluginId`, `gameCode` и `gameName` и эмитит события `save` и `close`; сохранение файла на сервер выполняет сама панель.

## Переводы

Переводы задаются словарями `translations: { en: {...}, ru: {...} }`. В полях `label`, `name`, `text` поддерживаются ссылки на ключи перевода вида `@:ключ` — панель подставит строку на текущем языке интерфейса.

## Доступ к API

Фронтенд плагина использует `window.axios` — тот же инстанс, что и панель, с авторизацией текущего пользователя. Через него доступны:

* API панели — например, отправка RCON-команды `POST /api/servers/{id}/rcon` или работа с файлами через `/api/file-manager/...`;
* бэкенд самого плагина по адресам `/api/plugins/{id}/...` (HTTP-маршруты, зарегистрированные WASM-частью).

## SDK `@gameap/plugin-sdk`

npm-пакет `@gameap/plugin-sdk` предоставляет:

* TypeScript-типы: `PluginDefinition`, `PluginRoute`, `PluginMenuItem`, `PluginSlotComponent`, `PluginHomeButton`, `PluginFileEditor`, `PluginContext` и другие;
* хуки контекста: `usePluginContext`, `useServer`, `useServerId`, `useServerAbilities`, `useCurrentUser`, `useIsAdmin`, `useIsAuthenticated`, `usePluginRoute`, `usePluginId`;
* хуки переводов: `usePluginTrans`, `providePluginTrans`;
* UI-компоненты панели (реэкспорт из `@gameap/ui`): `GCard`, `GDataTable`, `GModal`, `GStatusBadge`, `GSwitch` и другие;
* `createPluginConfig` — готовую конфигурацию Vite: сборка в lib-режиме (ES-модуль `plugin.js`), внешние зависимости (`vue`, `vue-router`, `pinia`, `axios`, `@gameap/ui`) переписываются на глобальные объекты. Учтите: панель экспонирует `window.NaiveUI`, но не `window.gameapUI`, поэтому реальные плагины используют собственный vite-конфиг с externals на `naive-ui` (образец — `frontend/vite.config.js` в [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor)).

## Сборка фронтенда

```bash
npm run build   # → dist/plugin.js (+ CSS-файл)
```

Собранные `plugin.js` и CSS встраиваются в `.wasm` плагина: для Rust — скриптом `build.rs` с копированием в `OUT_DIR` и подключением через `include_bytes!` (см. [Разработка плагинов](/ru/plugins/development.html)), для AssemblyScript — скриптом кодогенерации (образец — `scripts/embed-frontend.mjs` в plugin-minecraft-modrinth).

## Локальная отладка

Пакет `@gameap/debug` запускает отладочное окружение с реальным фронтендом панели и моками API (MSW):

```bash
PLUGIN_PATH=./dist npx @gameap/debug
```

Окружение открывается на `http://localhost:5174`. Плавающая debug-панель позволяет переключать тип пользователя (администратор / обычный пользователь / гость), задержку сети и локаль. Плагин должен быть предварительно собран (`npm run build`).

## Пример `PluginDefinition`

Пример вкладки на странице игрового сервера (по образцу plugin-goldsrc-addons):

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
