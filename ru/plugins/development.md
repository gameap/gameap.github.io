---
title: Разработка плагинов
layout: default
lang: ru
category: Плагины
order: 342
---

## Архитектура

Плагин GameAP — это WASM-модуль для таргета `wasm32-wasip1` (WASI preview 1), собранный как reactor. Панель исполняет каждый плагин в изолированном рантайме [wazero](https://github.com/tetratelabs/wazero): без файловой системы, сети и переменных окружения; stdout и stderr гостя попадают в лог панели (см. [Ограничения рантайма](#ограничения-рантайма)).

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
| `GetAssets` | Статические файлы плагина: переводы, отдаваемые по `/lang/`, и файлы фронтенда, отдаваемые из корня SPA (опционально) |

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
| `required_permissions` | Права, которые нужны плагину; при установке выдаются ровно они (см. [Права плагина](#права-плагина)) |
| `api_version` | Версия Plugin API, обязательно `"1"` |

**Требования к `id`.** Используйте стабильный идентификатор из символов base32-алфавита `a-z2-7` без дефисов. Панель нормализует id: строка с дефисами или другими символами будет заменена хешем, из-за чего сломаются пути `/api/plugins/{id}/...` и `/plugins/{id}/...`. Избегайте и чисто цифровых id: такой идентификатор трактуется как десятичный числовой ID (разбор id сначала пытается распарсить строку как число). Примеры корректных id реальных плагинов: `hexeditor4jm2`, `ezvdsxmlu6fbk`, `dshdabjp2l73a`.

## События

Плагин подписывается на события панели методом `GetSubscribedEvents`. Типы событий (перечислены без префикса `EVENT_TYPE_`):

| Событие | Отменяемое | Доставка |
|---|---|---|
| `SERVER_PRE_START`, `SERVER_PRE_STOP`, `SERVER_PRE_RESTART`, `SERVER_PRE_INSTALL`, `SERVER_PRE_UPDATE`, `SERVER_PRE_REINSTALL`, `SERVER_PRE_DELETE` | Да | Синхронно |
| `SERVER_POST_START`, `SERVER_POST_STOP`, `SERVER_POST_RESTART`, `SERVER_POST_INSTALL`, `SERVER_POST_UPDATE`, `SERVER_POST_REINSTALL`, `SERVER_POST_DELETE` | Нет | Асинхронно |
| `SERVER_CREATED`, `SERVER_UPDATED`, `SERVER_DELETED`, `SERVER_SETTINGS_CHANGED` | Нет | Асинхронно |
| `USER_PRE_DELETE` | Да | Синхронно |
| `USER_CREATED`, `USER_UPDATED`, `USER_DELETED` | Нет | Асинхронно |
| `NODE_PRE_DELETE` | Да | Синхронно |
| `NODE_CREATED`, `NODE_UPDATED`, `NODE_DELETED`, `NODE_ONLINE`, `NODE_OFFLINE` | Нет | Асинхронно |
| `DAEMON_TASK_CREATED`, `DAEMON_TASK_STARTED`, `DAEMON_TASK_COMPLETED`, `DAEMON_TASK_FAILED` | Нет | Асинхронно |
| `PLUGIN_LOADED`, `PLUGIN_UNLOADED`, `PLUGIN_ERROR` | Нет | Асинхронно |

Pre-события доставляются синхронно и блокируют операцию: плагин может отменить её, вернув `EventResult` с `should_cancel = true` и сообщением `message`, либо изменить данные через `modified_data`. Post-события и остальные типы доставляются асинхронно и не влияют на операцию.

Полезная нагрузка зависит от типа события: `server_event` (`ServerEventPayload`, полный снапшот сервера) для `SERVER_*`, `server_settings_event` (`ServerSettingsEventPayload`: id сервера и сохранённые настройки) для `SERVER_SETTINGS_CHANGED`, `task_event` (`TaskEventPayload`, данные задачи демона) для `DAEMON_TASK_*`, `user_event` (`UserEventPayload`) для `USER_*`, `node_event` (`NodeEventPayload`) для `NODE_*` и `plugin_event` (`PluginEventPayload`: компактный id плагина, название, версия, статус и ошибка) для `PLUGIN_*`. События `PLUGIN_*` никогда не доставляются тому плагину, о котором они сообщают.

Таймаут обработчика события — 10 секунд; по его истечении рантайм закрывает модуль, плагин отключается и затем автоматически перезагружается (см. [Ограничения рантайма](#ограничения-рантайма)). Асинхронные события доставляются в фоне с общим бюджетом 60 секунд на событие и не более чем 64 одновременными доставками.

При `PLUGINS_PERMISSIONS_ENFORCE=true` для получения событий нужно право `listen_events`: без него ответ `GetSubscribedEvents` игнорируется с предупреждением в логе панели, а право повторно проверяется перед каждой доставкой. Пока проверка прав выключена (по умолчанию в 4.5.0), подписки всех плагинов учитываются.

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

### Отдача файла с ноды в браузер

Вместо `body` маршрут может ответить полем `file` — сообщением `FileRef` (`node_id`, `path`, `filename`), указывающим на файл на ноде. Панель сама передаёт этот файл клиенту потоком, и байты не проходят через плагин; так и следует отдавать большие файлы в браузер. `path` — абсолютный путь на ноде в том же виде, что принимает `gameap-nodefs` (сегменты `..` отклоняются); `filename` — имя, предлагаемое клиенту, пустое значение означает базовое имя `path`.

Заголовки `Content-Length` и `Content-Disposition` (всегда attachment) выставляет панель; из `headers` плагина до клиента доходят только `Content-Type`, `Content-Language`, `Cache-Control`, `Expires`, `Pragma`, `Last-Modified`, `ETag`, `Vary` и заголовки `X-Plugin`/`X-Plugin-*`. Клиент должен быть аутентифицирован (иначе `401` независимо от `requires_auth`), плагину нужно право `files_read` (без него `403`), а путь проверяется [политикой путей на ноде](#политика-путей-на-ноде) (`403` при отказе). Панели, не знающие этого поля, игнорируют его и отправляют пустое тело.

## Host-функции

Все обращения плагина к панели и внешнему миру идут через host-функции, сгруппированные в 22 модуля `gameap-*`:

| Модуль | Функции | Назначение |
|---|---|---|
| `gameap-log` | `log` | Запись в лог панели |
| `gameap-cache` | `get`, `set`, `delete` | Общий кеш панели (ключи с префиксом `plugin:`, у каждого плагина своё пространство имён) |
| `gameap-crypto` | `random_uint64`, `random_string`, `argon2_hash`, `argon2_verify` | Генераторы случайных значений, хеширование Argon2 |
| `gameap-http` | `fetch` | Исходящие HTTP-запросы с защитой от SSRF |
| `gameap-storage` | `get`, `set`, `delete`, `list` | Персистентное key-value хранилище плагина |
| `gameap-secrets` | `get`, `set`, `delete`, `list_keys` | Зашифрованное хранилище учётных данных плагина |
| `gameap-servercontrol` | `start_server`, `stop_server`, `restart_server`, `update_server`, `install_server`, `reinstall_server` | Управление игровыми серверами (возвращает `task_id`) |
| `gameap-servers` | `find_servers`, `get_server`, `save_server`, `delete_server` | Игровые серверы |
| `gameap-users` | `find_users`, `get_user` | Пользователи панели |
| `gameap-nodes` | `find_nodes`, `get_node`, `update_node`, `delete_node`, `create_setup_key`, `get_setup_key`, `revoke_setup_key` | Выделенные серверы (ноды) и ключи установки демона |
| `gameap-games` | `find_games`, `get_game` | Игры |
| `gameap-gamemods` | `find_game_mods`, `get_game_mod` | Моды игр |
| `gameap-daemontasks` | `find_daemon_tasks`, `create_daemon_task` | Задачи демона |
| `gameap-serversettings` | `find_server_settings`, `save_server_setting` | Настройки игровых серверов |
| `gameap-nodefs` | `read_dir`, `mk_dir`, `copy`, `move`, `download`, `upload`, `remove`, `get_file_info`, `chmod`, `hash`, `create_archive`, `extract_archive`, `start_create_archive`, `start_extract_archive`, `cancel_archive`, `get_archive_operation` | Файловые операции на ноде |
| `gameap-nodecmd` | `execute_command` | Выполнение команды на ноде |
| `gameap-ssh` | `generate_key_pair`, `connect`, `disconnect`, `exec`, `start_exec`, `get_exec_operation`, `cancel_exec`, `write_file`, `read_file` | SSH-подключения к хостам, которые плагин указывает сам |
| `gameap-net` | `send`, `recv`, `close` | Ввод-вывод по соединениям, открытым панелью (расширения протоколов RCON/Query) |
| `gameap-scheduler` | `add_task`, `remove_task`, `list_tasks` | Задания плагина по расписанию |
| `gameap-authz` | `can`, `can_one_of`, `can_for_entity`, `can_any_for_entity`, `get_user_roles` | Проверка прав и ролей пользователей |
| `gameap-rbac` | `set_user_roles`, `allow_user_abilities_for_entity`, `revoke_or_forbid_user_abilities_for_entity`, `get_roles`, `save_role`, `delete_role`, `get_permissions`, `get_roles_for_entity`, `assign_roles_for_entity`, `clear_roles_for_entity`, `allow`, `forbid`, `revoke` | Управление ролями и правами |
| `gameap-host` | `get_grants`, `get_host_info` | Интроспекция: выданные плагину права, версия панели, версия Plugin API, id экземпляра и подключённые для плагина host-модули |

Особенности:

* `gameap-http` проксирует запросы через панель с защитой от SSRF: по умолчанию разрешена только схема `https`, приватные и служебные IP блокируются, тело ответа ограничено 10 МБ. Политика настраивается переменными окружения `PLUGINS_HTTP_*` (см. [Установка и управление](/ru/plugins/management.html)).
* `gameap-storage` — персистентное хранилище плагина, изолированное по `plugin_id`, с опциональной привязкой записей к сущности (`entity_type`, `entity_id`). Используйте его для настроек плагина: отдельный механизм конфигурации (`config` в `InitializeRequest`) в текущей версии панели не задействован. Хранилище ограничено квотами на плагин: `PLUGINS_STORAGE_MAX_KEYS_PER_PLUGIN` (10000), `PLUGINS_STORAGE_MAX_VALUE` (1M) и `PLUGINS_STORAGE_MAX_TOTAL` (64M). Данные хранятся в открытом виде — учётным данным место в `gameap-secrets`.
* Значения в `gameap-cache` ограничены `PLUGINS_CACHE_MAX_VALUE` (1M); записи истекают по TTL и не удаляются при удалении плагина.
* `gameap-nodefs` и `gameap-nodecmd` — работа с файлами и командами на выделенном сервере (ноде) через GameAP Daemon; каждый путь проверяется [политикой путей на ноде](#политика-путей-на-ноде).
* Изменяющие вызовы `gameap-nodes` (`update_node`, `delete_node`, ключи установки) требуют права `manage_nodes`, которое модуль проверяет сам; чтение открыто всем плагинам.
* `gameap-ssh` регистрируется только при `PLUGINS_SSH_ENABLED=true` (по умолчанию `false`), `gameap-net` — только при `PLUGINS_NET_ENABLED=true` (по умолчанию `true`). Плагин, импортирующий незарегистрированный панелью модуль, не загружается (см. [Обратная совместимость](#обратная-совместимость)).

### Чтение файлов с ноды по частям

`download` принимает окно `offset`/`length`. Без окна это чтение файла целиком, ограниченное `PLUGINS_NODEFS_MAX_INLINE` (32M): файл большего размера отклоняется с ошибкой, называющей оба размера. С окном в лимит должно уместиться только само окно (`length` больше лимита отклоняется, а не урезается), поэтому большой файл читается по частям. В ответе возвращаются `offset` и `total_size`; `length = 0` при заданном `offset` возвращает столько, сколько позволяет лимит, поэтому сдвигайтесь на `len(content)`, а не на запрошенную длину. Чтение продолжается, пока `offset + len(content) < total_size`; чтение на конце файла или за ним возвращает пустое содержимое, а не ошибку. Данные `upload` ограничены той же переменной.

### `gameap-secrets`

Зашифрованное хранилище для собственных учётных данных плагина (API-ключи, токены ботов), для которых открытое `gameap-storage` не подходит. Каждая функция требует права `secrets`; без него `set`/`delete` отвечают `success = false`, `get` — `found = false`, `list_keys` — пустым списком, а в `error` указано недостающее право; модуль при этом остаётся импортируемым.

* Значения шифруются ключом панели `ENCRYPTION_KEY` (AES-256-GCM); шифртекст привязан к id плагина-владельца и ключу, поэтому запись, скопированная в другое место базы, уже не расшифровывается.
* Без `ENCRYPTION_KEY` запись отклоняется, а не сохраняется в открытом виде; отключить это можно переменной `PLUGINS_SECRETS_REQUIRE_ENCRYPTION=false`.
* Ключи должны соответствовать `^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,63}$`; квоты: `PLUGINS_SECRETS_MAX_KEYS_PER_PLUGIN` (64) и `PLUGINS_SECRETS_MAX_VALUE` (8K).
* Секреты приватны для записавшего их плагина и удаляются вместе с его записями `gameap-storage` при удалении плагина.

### `gameap-ssh`

Открывает SSH-подключения к хостам, которые плагин указывает сам, и выполняет на них команды — единственное, чего не умеет `gameap-nodecmd`: тот работает через GameAP Daemon, а на подготавливаемой машине демона ещё нет. Модулю нужны право `ssh` **и** `PLUGINS_SSH_ENABLED=true` на панели (по умолчанию `false`): при выключенном переключателе модуль вообще не регистрируется, и импортирующий его плагин не загружается.

* `generate_key_pair` создаёт пару ключей запрошенного типа `KeyType`: `KEY_TYPE_ED25519` (он же подразумевается, если тип не задан), `KEY_TYPE_RSA_4096` или `KEY_TYPE_ECDSA_P256`. Право требуется для каждой функции модуля, включая генерацию ключей.
* В `connect` обязательна политика проверки ключа хоста: либо `accept_any`, либо список закреплённых ключей (отпечатки SHA-256 или публичные ключи), но не оба сразу — такое сочетание отклоняется. Наблюдаемый ключ хоста возвращается всегда, чтобы при первом подключении его можно было закрепить. Оператор может запретить `accept_any` для всей панели переменной `PLUGINS_SSH_ALLOW_ACCEPT_ANY_HOST_KEY=false` (по умолчанию `true`).
* `exec` ждёт завершения команды в пределах дедлайна гостевого вызова; когда бюджет исчерпан, возвращает `completed = false` и `operation_id`, а плагин подписывается на обратный вызов о завершении. `start_exec` сразу возвращает `operation_id`; завершение доставляется в плагин, если тот экспортирует сервис `SSHExecEventsHandler`, а `get_exec_operation` опрашивает статус и вывод; `cancel_exec` отменяет выполняющуюся команду.
* `write_file`/`read_file` передают данные через удалённый `cat`, поэтому SFTP на свежеустановленной машине не требуется.
* Подключения и операции живут в памяти того экземпляра панели, который их создал: они не переживают перезагрузку плагина, а на многоэкземплярной панели `get_exec_operation` на другом экземпляре ответит `found = false`. `disconnect` отменяет все операции, ещё выполняющиеся на подключении; при выгрузке плагина его подключения освобождаются.
* Адреса назначения проходят ту же политику, что и в `gameap-http`: `PLUGINS_SSH_BLOCK_PRIVATE_IPS` (по умолчанию `true`) и `PLUGINS_SSH_ALLOWED_HOSTS`. Вывод каждого потока ограничен `PLUGINS_SSH_MAX_OUTPUT_BYTES` (1 МБ) — сохраняется начало, а поток помечается как обрезанный.

## Права плагина

Привилегированные host-функции доступны только по выданным плагину правам — тем, которые разрешил администратор. Допустимые имена:

| Право | Что открывает |
|---|---|
| `manage_servers` | `gameap-servercontrol` (все функции), `gameap-daemontasks.create_daemon_task`, `gameap-servers.save_server`/`delete_server`, `gameap-serversettings.save_server_setting` |
| `node_commands` | `gameap-nodecmd.execute_command`; задачи демона `cmdexec` дополнительно к `manage_servers` |
| `files_read` | Чтение в `gameap-nodefs`: `read_dir`, `download`, `get_file_info`, `hash`, `get_archive_operation`; `HTTPResponse.file` |
| `files` | Все функции `gameap-nodefs`, включая запись и операции с архивами; включает `files_read` |
| `listen_events` | Подписки на события |
| `manage_rbac` | Все функции `gameap-rbac` |
| `secrets` | Все функции `gameap-secrets` |
| `ssh` | Все функции `gameap-ssh`, включая генерацию ключей; на панели также должно быть `PLUGINS_SSH_ENABLED=true` |
| `manage_nodes` | Изменяющие вызовы `gameap-nodes`: `update_node`, `delete_node`, `create_setup_key`, `get_setup_key`, `revoke_setup_key` |
| `manage_games`, `manage_game_mods`, `manage_users` | Объявлены, но пока не требуются ни одной host-функцией |

Модули только для чтения (`gameap-users`, `gameap-games`, `gameap-gamemods`, `gameap-authz`, `gameap-servers.find_servers`/`get_server`, `gameap-nodes.find_nodes`/`get_node`, …), а также `gameap-http`, `gameap-cache`, `gameap-storage`, `gameap-crypto`, `gameap-log`, `gameap-scheduler`, `gameap-net` и `gameap-host` прав не требуют.

Плагин объявляет нужные права в `PluginInfo.required_permissions`. Установка — загрузкой файла, из каталога или через `PLUGINS_AUTOLOAD` — выдаёт ровно объявленные права; неизвестные имена отбрасываются. Обновление права не расширяет: новое объявление записывается в `required_permissions`, а `allowed_permissions` остаётся прежним, поэтому у сборки, которой стало нужно больше прав, такие вызовы отклоняются, пока администратор не выдаст их в разделе **Администрирование** → **Плагины** (действие **Права** в строке плагина) или через `PUT /api/admin/plugins/{id}/permissions`.

Каждый привилегированный вызов проходит через guard: сначала проверяется право, затем лимит частоты вызовов плагина (token bucket по классам, см. [Ограничения рантайма](#ограничения-рантайма)), а пути на ноде проверяются политикой путей. Отклонённый вызов получает ответ в поле `error` — `plugin permission <name> required`, `rate limited: gameap-<module> allows N calls/s (burst M)` или `path policy: <reason>: <path>` — и записывается в журнал аудита как `access.denied` (превышение лимита — как `plugin.hostcall.ratelimited`) с плагином в роли актора. Модуль остаётся загруженным; за отклонённый вызов плагин никогда не отключается.

Проверка прав вводится поэтапно: при `PLUGINS_PERMISSIONS_ENFORCE=false` (по умолчанию в 4.5.0) проверки прав в guard проходят, а сами права по-прежнему записываются и показываются в интерфейсе. Лимиты частоты, политика путей на ноде, переключатель `PLUGINS_SSH_ENABLED` и проверка `manage_nodes` внутри `gameap-nodes` действуют всегда. Объявляйте и выдавайте нужные плагину права уже сейчас, чтобы ничего не сломалось при включении проверки.

> Права — механизм ограничения, но плагин всё равно выполняется внутри панели и имеет доступ к её данным. Поэтому установка плагинов доверяется только администраторам — устанавливайте плагины только из источников, которым доверяете.

### Политика путей на ноде

Каждый путь, который плагин передаёт в `gameap-nodefs`, `work_dir` в `gameap-nodecmd` и `HTTPResponse.file` проверяется на панели до отправки демону. В любом режиме отклоняется путь с сегментом `..` или NUL-байтом. Что разрешено помимо этого, задаёт `PLUGINS_NODEFS_PATH_POLICY`:

| Режим | Разрешённые пути |
|---|---|
| `unrestricted` (по умолчанию) | Всё, что разрешает сам демон |
| `node_workpath` | Внутри рабочего каталога ноды; относительные пути разрешаются относительно него |
| `server_dirs` | Внутри каталога игрового сервера на этой ноде (`<work_path>/<server dir>`) |

В обоих ограниченных режимах `work_dir` в `gameap-nodecmd` должен быть абсолютным путём: относительный путь демон разрешил бы относительно собственного рабочего каталога. Команда без `work_dir` выполняется в каталоге демона по умолчанию и не проверяется.

Оба ограниченных режима оставляют открытым один каталог — служебный каталог плагина `<work_path>/.plugins/<компактный id плагина>`, вычисляемый от рабочего каталога каждой ноды и открытый только тому плагину, чей id в нём указан. Держите там рабочие файлы на ноде (подготовленные запросы, результаты долгих команд, собранные для скачивания архивы) — и самый строгий режим продолжит работать. `PLUGINS_NODEFS_ALLOWED_PATHS` (абсолютные корни через запятую) расширяет ограниченные режимы на всех нодах. Отклонённый вызов получает `path policy: <reason>: <path>` (`403` для ссылок на файлы) и попадает в аудит; плагин за это не отключается.

## Ограничения рантайма

| Ограничение | Значение |
|---|---|
| Одновременные вызовы одного плагина | 1 (вызовы сериализуются) |
| Таймаут вызова плагина | 30 с (старт модуля — 60 с) |
| Таймаут обработчика события | 10 с |
| Доставка асинхронных событий | 60 с на событие, не более 64 одновременных доставок |
| Линейная память модуля | `PLUGINS_RUNTIME_MAX_MEMORY`, 256M (объявленный максимум больше лимита урезается; модуль, чья начальная память превышает лимит, не загружается) |
| Размер `.wasm`-файла | `PLUGINS_RUNTIME_MAX_MODULE_SIZE`, 128M; загрузка по HTTP дополнительно ограничена 100 МБ |
| Тело HTTP-запроса к плагину | 1 МБ |
| Тело ответа `gameap-http` | 10 МБ |
| Встроенные `upload`/`download` в `gameap-nodefs` | `PLUGINS_NODEFS_MAX_INLINE`, 32M |
| `gameap-storage` | 10000 ключей, 1M на значение, 64M суммарно на плагин |
| Значение `gameap-cache` | `PLUGINS_CACHE_MAX_VALUE`, 1M |
| `gameap-secrets` | 64 ключа, 8K на значение |

Размеры принимают суффикс единицы (`512K`, `64M`, `1G`); число без суффикса — байты.

Затратные host-функции ограничены по частоте вызовов для каждого плагина (token bucket) по классам:

| Класс | Функции | По умолчанию |
|---|---|---|
| `nodecmd` | `gameap-nodecmd.execute_command` | 5/с, burst 20 |
| `servercontrol` | `gameap-servercontrol.*`, `create_daemon_task`, `save_server`, `delete_server`, `save_server_setting` | 5/с, burst 20 |
| `nodefs` | Все функции `gameap-nodefs` | 50/с, burst 200 |
| `http` | `gameap-http.fetch` | 20/с, burst 50 |
| `rbac` | Все функции `gameap-rbac` | 10/с, burst 50 |
| `ssh` | Все функции `gameap-ssh` | 20/с, burst 60 |

Лимиты настраиваются переменными `PLUGINS_RATELIMIT_<CLASS>_RPS` / `_BURST`; RPS `0` отключает класс. Счётчики ведутся отдельно на каждом экземпляре панели.

При превышении таймаута вызова (обработчик события, HTTP-маршрут, задание по расписанию, обратный вызов) или когда гость сам завершает модуль (паника, заканчивающаяся `proc_exit`), рантайм закрывает модуль: плагин получает статус `error` с причиной в `last_error`, а панель автоматически перезагружает его через `PLUGINS_RECOVERY_INITIAL_DELAY` (30s), удваивая паузу при каждом следующем сбое вплоть до `PLUGINS_RECOVERY_MAX_DELAY` (10m). После `PLUGINS_RECOVERY_MAX_ATTEMPTS` (5) перезагрузок подряд плагин остаётся в статусе `error`, пока администратор не перезагрузит его (**Перезагрузить** в разделе **Администрирование** → **Плагины** или `POST /api/admin/plugins/{id}/reload`) или не перезапустится панель. `PLUGINS_RECOVERY_ENABLED=false` полностью выключает автоматическую перезагрузку: отключение всё равно фиксируется (статус `error`, причина, запись в аудите), но плагин вернётся в работу только после той же ручной перезагрузки или перезапуска панели. Пишите плагин с учётом этого: держите состояние в `gameap-storage`, а не в глобальных переменных модуля, и делайте `Initialize` идемпотентным — перезагрузка создаёт новый экземпляр модуля.

Stdout гостя попадает в лог панели с уровнем debug, stderr — с уровнем warn, поэтому сообщение о панике видно; строки обрезаются до 4 КиБ, а каждый поток ограничен 200 строками за 10 секунд (отброшенные строки подсчитываются и отмечаются в логе). Из WASM недоступны файловая система и сеть — только через host-функции.

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

Панель умеет валидировать `.wasm`-файл без установки — эндпоинт `POST /api/admin/plugins/upload/dry-run` (multipart, поле `file`). В ответ возвращаются метаданные плагина, HTTP-маршруты, серверные права, подписки на события, признаки фронтенд-части (`has_frontend_bundle`, `frontend_bundle_size`, `has_frontend_styles`) и список ошибок. Также возвращаются права: `required_permissions` (объявленные в манифесте), `used_permissions` (вычисленные по импортам host-функций модуля и его подпискам на события) и `undeclared_permissions` — используемые, но не объявленные. При установке выдаются только объявленные права, поэтому при `PLUGINS_PERMISSIONS_ENFORCE=true` такие вызовы будут отклоняться; со значением по умолчанию в 4.5.0 проверки проходят, а панель лишь пишет предупреждение в лог при загрузке плагина. Если плагин с тем же id уже установлен, об этом сообщают поля `installed`, `installed_version` и `installed_source_type`: загрузка файла заменит его. Та же проверка выполняется в интерфейсе при загрузке файла.

## Обратная совместимость

Совместимость плагинов гарантируется в тех пределах, в каких её проверяет CI панели: матрица `pkg/plugin/compatrust` загружает фиксированный набор тестовых Rust-плагинов через `plugin.Manager` на версиях `HEAD`, `v4.4.1` и `v4.3.5` и требует, чтобы одни и те же сборки плагинов загружались и отвечали на вызовы на каждой из этих версий панели; успешный прогон обязателен для любого изменения контрактов плагинов (`pkg/plugin/proto`, `pkg/plugin/sdk`). Версии панели вне этой матрицы ею не охвачены.

И наоборот: плагин, импортирующий host-модуль, которого панель не предоставляет, отклоняется при загрузке с ошибкой, называющей модуль, а не падает при вызове. На панели 4.5.0 обычная причина — `gameap-ssh` при `PLUGINS_SSH_ENABLED=false` или `gameap-net` при `PLUGINS_NET_ENABLED=false`; если загрузка не удалась из-за неизвестного импорта, сначала проверьте конфигурацию панели. Во время работы плагин может узнать своё окружение через `gameap-host.get_host_info` (версия панели, версия Plugin API, id экземпляра и подключённые для него host-модули) и `get_grants` (его текущие права).

## Примеры плагинов

* [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) — минимальный Rust-плагин: метаданные + встроенный фронтенд, без host-функций.
* [plugin-goldsrc-addons](https://github.com/gameap/plugin-goldsrc-addons) — Rust-плагин с HTTP-маршрутами, работой с файлами на ноде (`nodefs`), выполнением команд (`nodecmd`) и вызовами API панели (RCON) с фронтенда.
* [plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth) — плагин на AssemblyScript: обращения к внешнему API (modrinth.com) через `gameap-http`, кеширование, персистентное хранилище.
