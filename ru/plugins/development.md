---
title: Разработка плагинов
layout: default
lang: ru
category: Плагины
order: 342
---

* This will become a table of contents (this text will be scraped).
{:toc}

## Архитектура

Плагин GameAP — это WASM-модуль для таргета `wasm32-wasip1` (WASI preview 1), собранный как reactor. Панель исполняет каждый плагин в изолированном рантайме [wazero](https://github.com/tetratelabs/wazero): без файловой системы, сети и переменных окружения, stdout и stderr отбрасываются.

Обмен между панелью и плагином происходит по протоколу protobuf через линейную память WASM: панель пишет запрос в память гостя через экспортируемый плагином `malloc`, вызывает функцию и читает ответ. Версия Plugin API — `1`.

Типы сообщений и SDK публикуются из репозитория [github.com/gameap/gameap-proto](https://github.com/gameap/gameap-proto):

* `gameap-plugin-sdk` (crates.io) — полный Rust SDK: ABI-прослойка, трейт `Plugin`, макрос `register_plugin!`, типизированные клиенты host-функций;
* `@gameap/proto-as` (npm) — AssemblyScript-типы сообщений (без ABI-прослойки);
* `@gameap/proto` (npm) — TypeScript-типы (для инструментария, собрать WASM-плагин на чистом JS сейчас нельзя).

## Интерфейс плагина

Плагин реализует сервис `PluginService` из `plugin.proto`. Методы сервиса:

| Метод | Назначение |
|---|---|
| `GetInfo` | Возвращает метаданные плагина (`PluginInfo`) |
| `Initialize` | Инициализация при загрузке плагина |
| `Shutdown` | Завершение работы при выгрузке плагина |
| `HandleEvent` | Обработка события панели |
| `GetSubscribedEvents` | Список типов событий, на которые подписан плагин |
| `GetHTTPRoutes` | Список HTTP-маршрутов плагина |
| `HandleHTTPRequest` | Обработка HTTP-запроса по маршруту плагина |
| `GetFrontendBundle` | Встроенный фронтенд: JS-бандл и CSS (опционально) |
| `GetServerAbilities` | Регистрация серверных прав плагина (опционально) |

В Rust SDK интерфейс представлен трейтом `Plugin` с нейтральными реализациями по умолчанию — фактически обязателен только `get_info`. Макрос `register_plugin!` генерирует все необходимые WASM-экспорты, включая `malloc`/`free` и проверку версии API.

## Метаданные `PluginInfo`

| Поле | Описание |
|---|---|
| `id` | Строковый идентификатор плагина (см. требования ниже) |
| `name` | Название плагина |
| `version` | Версия (семвер) |
| `description` | Краткое описание |
| `author` | Автор |
| `license` | Лицензия (например, `MIT`) |
| `homepage` | Ссылка на страницу плагина |
| `required_permissions` | Заявляемые разрешения (в текущей версии панели не проверяются) |
| `api_version` | Версия Plugin API, обязательно `"1"` |

**Требования к `id`.** Используйте стабильный идентификатор из символов base32-алфавита `a-z2-7` без дефисов. Панель нормализует id: строка с дефисами или другими символами будет заменена хешем, из-за чего сломаются пути `/api/plugins/{id}/...` и `/plugins/{id}/...`. Избегайте и чисто цифровых id: такой идентификатор трактуется как десятичный числовой ID (разбор id сначала пытается распарсить строку как число). Примеры корректных id реальных плагинов: `hexeditor4jm2`, `ezvdsxmlu6fbk`, `dshdabjp2l73a`.

## События

Плагин подписывается на события панели методом `GetSubscribedEvents`. Типы событий (перечислены без префикса `EVENT_TYPE_`):

| Событие | Отменяемое | Доставка |
|---|---|---|
| `SERVER_PRE_START`, `SERVER_PRE_STOP`, `SERVER_PRE_RESTART`, `SERVER_PRE_INSTALL`, `SERVER_PRE_UPDATE`, `SERVER_PRE_REINSTALL`, `SERVER_PRE_DELETE` | Да | Синхронно |
| `SERVER_POST_START`, `SERVER_POST_STOP`, `SERVER_POST_RESTART`, `SERVER_POST_INSTALL`, `SERVER_POST_UPDATE`, `SERVER_POST_REINSTALL`, `SERVER_POST_DELETE` | Нет | Асинхронно |
| `SERVER_CREATED`, `SERVER_UPDATED`, `SERVER_DELETED` | Нет | Асинхронно |
| `DAEMON_TASK_CREATED`, `DAEMON_TASK_COMPLETED`, `DAEMON_TASK_FAILED` | Нет | Асинхронно |

Pre-события доставляются синхронно и блокируют операцию: плагин может отменить её, вернув `EventResult` с `should_cancel = true` и сообщением `message`, либо изменить данные через `modified_data`. Post-события и остальные типы доставляются асинхронно и не влияют на операцию.

Событие сервера содержит полный снапшот сервера (`ServerEventPayload`), событие задачи — данные задачи демона (`TaskEventPayload`). Таймаут обработчика события — 10 секунд; по его истечении плагин отключается до перезагрузки панели.

## HTTP-маршруты плагина

Плагин регистрирует HTTP-маршруты методом `GetHTTPRoutes`, каждый маршрут — сообщение `HTTPRoute`:

| Поле | Описание |
|---|---|
| `path` | Путь, начинается с `/`; поддерживаются path-параметры вида `{name}` |
| `methods` | Методы: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, `OPTIONS` |
| `requires_auth` | Требовать аутентифицированного пользователя |
| `admin_only` | Требовать администратора |
| `description` | Описание маршрута |

Маршруты доступны по адресам `/api/plugins/{plugin_id}/...`. Запрос к плагину содержит метод, относительный путь, заголовки, path- и query-параметры, тело (не более 1 МБ) и сессию пользователя (если запрос аутентифицирован). Плагин возвращает `HTTPResponse` с полями `status_code`, `headers`, `body`.

## Host-функции

Все обращения плагина к панели и внешнему миру идут через host-функции, сгруппированные в модули `gameap-*`:

| Модуль | Функции | Назначение |
|---|---|---|
| `gameap-log` | `log` | Запись в лог панели |
| `gameap-cache` | `get`, `set`, `delete` | Общий кеш панели (ключи с префиксом `plugin:`) |
| `gameap-crypto` | `random_uint64`, `random_string`, `argon2_hash`, `argon2_verify` | Генераторы случайных значений, хеширование Argon2 |
| `gameap-http` | `fetch` | Исходящие HTTP-запросы с защитой от SSRF |
| `gameap-storage` | `get`, `set`, `delete`, `list` | Персистентное key-value хранилище плагина |
| `gameap-servercontrol` | `start_server`, `stop_server`, `restart_server`, `update_server`, `install_server`, `reinstall_server` | Управление игровыми серверами (возвращает `task_id`) |
| `gameap-servers` | `find_servers`, `get_server`, `save_server`, `delete_server` | Игровые серверы |
| `gameap-users` | `find_users`, `get_user` | Пользователи панели |
| `gameap-nodes` | `find_nodes`, `get_node` | Выделенные серверы (ноды) |
| `gameap-games` | `find_games`, `get_game` | Игры |
| `gameap-gamemods` | `find_game_mods`, `get_game_mod` | Моды игр |
| `gameap-daemontasks` | `find_daemon_tasks`, `create_daemon_task` | Задачи демона |
| `gameap-serversettings` | `find_server_settings`, `save_server_setting` | Настройки игровых серверов |
| `gameap-nodefs` | `read_dir`, `mk_dir`, `copy`, `move`, `download`, `upload`, `remove`, `get_file_info`, `chmod` | Файловые операции на ноде |
| `gameap-nodecmd` | `execute_command` | Выполнение команды на ноде |

Особенности:

* `gameap-http` проксирует запросы через панель с защитой от SSRF: по умолчанию разрешена только схема `https`, приватные и служебные IP блокируются, тело ответа ограничено 10 МБ. Политика настраивается переменными окружения `PLUGIN_HTTP_*` (см. [Установка и управление](/ru/plugins/management.html)).
* `gameap-storage` — персистентное хранилище плагина, изолированное по `plugin_id`, с опциональной привязкой записей к сущности (`entity_type`, `entity_id`). Используйте его для настроек плагина: отдельный механизм конфигурации (`config` в `InitializeRequest`) в текущей версии панели не задействован.
* `gameap-nodefs` и `gameap-nodecmd` — работа с файлами и командами на выделенном сервере (ноде) через GameAP Daemon.

**Важно:** host-функции выполняются с правами самой панели, без дополнительных проверок. Поэтому установка плагинов доверяется только администраторам — устанавливайте плагины только из источников, которым доверяете.

## Ограничения рантайма

| Ограничение | Значение |
|---|---|
| Одновременные вызовы одного плагина | 1 (вызовы сериализуются) |
| Таймаут вызова плагина | 30 с (старт модуля — 60 с) |
| Таймаут обработчика события | 10 с |
| Размер загружаемого `.wasm`-файла | 100 МБ |
| Тело HTTP-запроса к плагину | 1 МБ |
| Тело ответа `gameap-http` | 10 МБ |

При превышении таймаута вызова панель отключает плагин до перезагрузки. Из WASM недоступны файловая система и сеть — только через host-функции.

## Сборка плагина на Rust

Структура минимального проекта (по образцу [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor)):

`rust-toolchain.toml`:
```toml
[toolchain]
channel = "1.94.0"
targets = ["wasm32-wasip1"]
```

`Cargo.toml`:
```toml
[lib]
crate-type = ["cdylib"]

[dependencies]
gameap-plugin-sdk = "0.1"

[profile.release]
opt-level = "z"
lto = true
strip = true
```

Минимальный `src/lib.rs` — плагин, отдающий только метаданные и встроенный фронтенд:

```rust
#![cfg(target_arch = "wasm32")]

use gameap_plugin_sdk::proto::gameap::plugin as pb;
use gameap_plugin_sdk::{Plugin, PluginError, register_plugin};

const FRONTEND_JS: &[u8] = include_bytes!(concat!(env!("OUT_DIR"), "/plugin.js"));
const FRONTEND_CSS: &[u8] = include_bytes!(concat!(env!("OUT_DIR"), "/plugin.css"));

#[derive(Default)]
struct MyPlugin;

impl Plugin for MyPlugin {
    fn get_info(&mut self, _req: pb::GetInfoRequest) -> Result<pb::PluginInfo, PluginError> {
        Ok(pb::PluginInfo {
            id: "myplugin2j7d".into(),
            name: "My Plugin".into(),
            version: env!("CARGO_PKG_VERSION").into(),
            description: "My first GameAP plugin".into(),
            author: "Me".into(),
            api_version: "1".into(),
            ..Default::default()
        })
    }

    fn get_frontend_bundle(
        &mut self,
        _req: pb::GetFrontendBundleRequest,
    ) -> Result<pb::GetFrontendBundleResponse, PluginError> {
        Ok(pb::GetFrontendBundleResponse {
            bundle: FRONTEND_JS.to_vec(),
            has_bundle: !FRONTEND_JS.is_empty(),
            styles: FRONTEND_CSS.to_vec(),
            has_styles: !FRONTEND_CSS.is_empty(),
        })
    }
}

register_plugin!(MyPlugin);
```

Сборка:

```bash
cargo build --target wasm32-wasip1 --release
# опционально — уменьшение размера (binaryen):
wasm-opt -Oz target/wasm32-wasip1/release/my_plugin.wasm -o my-plugin.wasm
```

Полученный `.wasm`-файл устанавливается через интерфейс панели или копируется в каталог `plugins/` (см. [Установка и управление](/ru/plugins/management.html)). Фронтенд-часть собирается отдельно и встраивается в `.wasm` (см. [Фронтенд плагина](/ru/plugins/frontend.html)).

## Разработка на AssemblyScript

Готового SDK для AssemblyScript пока нет: пакет `@gameap/proto-as` предоставляет только типы сообщений, а ABI-прослойку (упаковку указателей, обработку ошибок, экспорт `malloc`/`free`) нужно писать вручную. Из-за сборщика мусора AssemblyScript буферы, передаваемые хосту, необходимо закреплять (`__pin`/`__unpin`), а JSON-парсер и базовые утилиты придётся реализовать самостоятельно или подключить сторонние библиотеки. Рабочий образец плагина на AssemblyScript — [plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth).

## Проверка плагина перед установкой

Панель умеет валидировать `.wasm`-файл без установки — эндпоинт `POST /api/admin/plugins/upload/dry-run` (multipart, поле `file`). В ответ возвращаются метаданные плагина, HTTP-маршруты, серверные права, подписки на события, признак фронтенд-части и список ошибок. Та же проверка выполняется в интерфейсе при загрузке файла.

## Примеры плагинов

* [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) — минимальный Rust-плагин: метаданные + встроенный фронтенд, без host-функций.
* [plugin-goldsrc-addons](https://github.com/gameap/plugin-goldsrc-addons) — Rust-плагин с HTTP-маршрутами, работой с файлами на ноде (`nodefs`), выполнением команд (`nodecmd`) и вызовами API панели (RCON) с фронтенда.
* [plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth) — плагин на AssemblyScript: обращения к внешнему API (modrinth.com) через `gameap-http`, кеширование, персистентное хранилище.
