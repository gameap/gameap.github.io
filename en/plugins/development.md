---
title: Plugin development
layout: default
lang: en
category: Plugins
order: 342
---

## Architecture

A GameAP plugin is a WASM module for the `wasm32-wasip1` target (WASI preview 1), built as a reactor. The panel runs each plugin in an isolated [wazero](https://github.com/tetratelabs/wazero) runtime: no file system, no network, no environment variables; the guest's stdout and stderr are forwarded to the panel log (see [Runtime limits](#runtime-limits)).

The panel and the plugin exchange protobuf messages through WASM linear memory: the panel writes the request into guest memory using the `malloc` exported by the plugin, calls the function and reads the response. The Plugin API version is `1`.

Message types and SDKs are published from the [github.com/gameap/gameap-proto](https://github.com/gameap/gameap-proto) repository:

* `gameap-plugin-sdk` (crates.io) — the full Rust SDK: the ABI layer, the `Plugin` trait, the `register_plugin!` macro, typed host-function clients;
* `@gameap/proto-as` (npm) — AssemblyScript message types (without the ABI layer);
* `@gameap/proto` (npm) — TypeScript types (for tooling; building a WASM plugin in plain JS is not currently possible).

## The plugin interface

A plugin implements the `PluginService` service from `plugin.proto`. The service methods:

| Method | Purpose |
|---|---|
| `GetInfo` | Returns the plugin metadata (`PluginInfo`) |
| `Initialize` | Initialization when the plugin is loaded |
| `Shutdown` | Shutdown when the plugin is unloaded |
| `HandleEvent` | Handling a panel event |
| `GetSubscribedEvents` | The list of event types the plugin subscribes to |
| `GetHTTPRoutes` | The list of the plugin's HTTP routes |
| `HandleHTTPRequest` | Handling an HTTP request on a plugin route |
| `GetFrontendBundle` | The embedded frontend: JS bundle and CSS (optional) |
| `GetServerAbilities` | Registering the plugin's server permissions (optional) |
| `GetAssets` | Static files the plugin contributes: translations served at `/lang/` and frontend files served from the SPA root (optional) |

In the Rust SDK the interface is represented by the `Plugin` trait with neutral default implementations — in practice only `get_info` is mandatory. The `register_plugin!` macro generates all the required WASM exports, including `malloc`/`free` and the API version check.

## `PluginInfo` metadata

| Field | Description |
|---|---|
| `id` | The plugin's string identifier (see the requirements below) |
| `name` | Plugin name |
| `version` | Version (semver) |
| `description` | Short description |
| `author` | Author |
| `license` | License (for example, `MIT`) |
| `homepage` | Link to the plugin page |
| `required_permissions` | The permissions the plugin needs; installing grants exactly these (see [Plugin permissions](#plugin-permissions)) |
| `api_version` | Plugin API version, must be `"1"` |

**Requirements for `id`.** Use a stable identifier made of base32 alphabet characters `a-z2-7`, without hyphens. The panel normalizes the id: a string with hyphens or other characters is replaced by a hash, which breaks the `/api/plugins/{id}/...` and `/plugins/{id}/...` paths. Avoid purely numeric ids as well: such an identifier is treated as a decimal numeric ID (id parsing first tries to parse the string as a number). Examples of valid ids from real plugins: `hexeditor4jm2`, `ezvdsxmlu6fbk`, `dshdabjp2l73a`.

## Events

A plugin subscribes to panel events with the `GetSubscribedEvents` method. The event types (listed without the `EVENT_TYPE_` prefix):

| Event | Cancellable | Delivery |
|---|---|---|
| `SERVER_PRE_START`, `SERVER_PRE_STOP`, `SERVER_PRE_RESTART`, `SERVER_PRE_INSTALL`, `SERVER_PRE_UPDATE`, `SERVER_PRE_REINSTALL`, `SERVER_PRE_DELETE` | Yes | Synchronous |
| `SERVER_POST_START`, `SERVER_POST_STOP`, `SERVER_POST_RESTART`, `SERVER_POST_INSTALL`, `SERVER_POST_UPDATE`, `SERVER_POST_REINSTALL`, `SERVER_POST_DELETE` | No | Asynchronous |
| `SERVER_CREATED`, `SERVER_UPDATED`, `SERVER_DELETED`, `SERVER_SETTINGS_CHANGED` | No | Asynchronous |
| `USER_PRE_DELETE` | Yes | Synchronous |
| `USER_CREATED`, `USER_UPDATED`, `USER_DELETED` | No | Asynchronous |
| `NODE_PRE_DELETE` | Yes | Synchronous |
| `NODE_CREATED`, `NODE_UPDATED`, `NODE_DELETED`, `NODE_ONLINE`, `NODE_OFFLINE` | No | Asynchronous |
| `DAEMON_TASK_CREATED`, `DAEMON_TASK_STARTED`, `DAEMON_TASK_COMPLETED`, `DAEMON_TASK_FAILED` | No | Asynchronous |
| `PLUGIN_LOADED`, `PLUGIN_UNLOADED`, `PLUGIN_ERROR` | No | Asynchronous |

Pre-events are delivered synchronously and block the operation: a plugin can cancel it by returning an `EventResult` with `should_cancel = true` and a `message`, or modify the data via `modified_data`. Post-events and the other types are delivered asynchronously and do not affect the operation.

The payload depends on the event type: `server_event` (`ServerEventPayload`, a full snapshot of the server) for `SERVER_*`, `server_settings_event` (`ServerSettingsEventPayload`: the server id and the saved settings) for `SERVER_SETTINGS_CHANGED`, `task_event` (`TaskEventPayload`, the daemon task data) for `DAEMON_TASK_*`, `user_event` (`UserEventPayload`) for `USER_*`, `node_event` (`NodeEventPayload`) for `NODE_*` and `plugin_event` (`PluginEventPayload`: the compact plugin id, name, version, status and error) for `PLUGIN_*`. `PLUGIN_*` events are never delivered to the plugin they describe.

The event handler timeout is 10 seconds; once it expires the runtime closes the module and the plugin is disabled, then reloaded automatically (see [Runtime limits](#runtime-limits)). Asynchronous events are delivered in the background with a total budget of 60 seconds per event and at most 64 concurrent deliveries.

Receiving events requires the `listen_events` permission when `PLUGINS_PERMISSIONS_ENFORCE=true`: without the grant the `GetSubscribedEvents` answer is ignored with a warning in the panel log, and the grant is re-checked before every delivery. With enforcement off (the default in 4.5.0) every plugin's subscriptions are honoured.

## Plugin HTTP routes

A plugin registers HTTP routes with the `GetHTTPRoutes` method; each route is an `HTTPRoute` message:

| Field | Description |
|---|---|
| `path` | The path, starts with `/`; path parameters of the form `{name}` are supported |
| `methods` | Methods: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, `OPTIONS` |
| `requires_auth` | Require an authenticated user |
| `admin_only` | Require an administrator |
| `description` | Route description |

Routes are available at `/api/plugins/{plugin_id}/...`. The request passed to the plugin contains the method, the relative path, the headers, path and query parameters, the body (no larger than 1 MB) and the user session (if the request is authenticated). The plugin returns an `HTTPResponse` with the `status_code`, `headers` and `body` fields.

### Serving a node file to the browser

Instead of `body`, a route can answer with the `file` field — a `FileRef` message (`node_id`, `path`, `filename`) pointing at a file on a node. The panel streams that file to the client itself, so the bytes never pass through the plugin; this is the way to hand large files to the browser. `path` is an absolute node path in the same form `gameap-nodefs` accepts (`..` segments are rejected); `filename` is the name offered to the client, empty means the base name of `path`.

The panel sets `Content-Length` and `Content-Disposition` (always an attachment); of the plugin's `headers` only `Content-Type`, `Content-Language`, `Cache-Control`, `Expires`, `Pragma`, `Last-Modified`, `ETag`, `Vary` and the `X-Plugin`/`X-Plugin-*` headers reach the client. The client must be authenticated (`401` otherwise, whatever `requires_auth` says), the plugin needs the `files_read` grant (`403` without it), and the path is checked against the [node path policy](#node-path-policy) (`403` when refused). Panels that predate this field ignore it and send an empty body.

## Host functions

All calls from a plugin to the panel and the outside world go through host functions, grouped into 22 `gameap-*` modules:

| Module | Functions | Purpose |
|---|---|---|
| `gameap-log` | `log` | Writing to the panel log |
| `gameap-cache` | `get`, `set`, `delete` | The panel's shared cache (keys prefixed with `plugin:`, each plugin gets its own namespace) |
| `gameap-crypto` | `random_uint64`, `random_string`, `argon2_hash`, `argon2_verify` | Random value generators, Argon2 hashing |
| `gameap-http` | `fetch` | Outbound HTTP requests with SSRF protection |
| `gameap-storage` | `get`, `set`, `delete`, `list` | The plugin's persistent key-value storage |
| `gameap-secrets` | `get`, `set`, `delete`, `list_keys` | The plugin's encrypted credential store |
| `gameap-servercontrol` | `start_server`, `stop_server`, `restart_server`, `update_server`, `install_server`, `reinstall_server` | Game server control (returns a `task_id`) |
| `gameap-servers` | `find_servers`, `get_server`, `save_server`, `delete_server` | Game servers |
| `gameap-users` | `find_users`, `get_user` | Panel users |
| `gameap-nodes` | `find_nodes`, `get_node`, `update_node`, `delete_node`, `create_setup_key`, `get_setup_key`, `revoke_setup_key` | Dedicated servers (nodes) and daemon enrollment setup keys |
| `gameap-games` | `find_games`, `get_game` | Games |
| `gameap-gamemods` | `find_game_mods`, `get_game_mod` | Game mods |
| `gameap-daemontasks` | `find_daemon_tasks`, `create_daemon_task` | Daemon tasks |
| `gameap-serversettings` | `find_server_settings`, `save_server_setting` | Game server settings |
| `gameap-nodefs` | `read_dir`, `mk_dir`, `copy`, `move`, `download`, `upload`, `remove`, `get_file_info`, `chmod`, `hash`, `create_archive`, `extract_archive`, `start_create_archive`, `start_extract_archive`, `cancel_archive`, `get_archive_operation` | File operations on a node |
| `gameap-nodecmd` | `execute_command` | Executing a command on a node |
| `gameap-ssh` | `generate_key_pair`, `connect`, `disconnect`, `exec`, `start_exec`, `get_exec_operation`, `cancel_exec`, `write_file`, `read_file` | SSH connections to hosts the plugin names itself |
| `gameap-net` | `send`, `recv`, `close` | I/O on connections opened by the panel (RCON/Query protocol extensions) |
| `gameap-scheduler` | `add_task`, `remove_task`, `list_tasks` | Scheduled plugin tasks |
| `gameap-authz` | `can`, `can_one_of`, `can_for_entity`, `can_any_for_entity`, `get_user_roles` | Checking users' abilities and roles |
| `gameap-rbac` | `set_user_roles`, `allow_user_abilities_for_entity`, `revoke_or_forbid_user_abilities_for_entity`, `get_roles`, `save_role`, `delete_role`, `get_permissions`, `get_roles_for_entity`, `assign_roles_for_entity`, `clear_roles_for_entity`, `allow`, `forbid`, `revoke` | Managing roles and abilities |
| `gameap-host` | `get_grants`, `get_host_info` | Introspection: the plugin's grants, the panel version, the Plugin API version, the instance id and the host modules instantiated for the plugin |

Details:

* `gameap-http` proxies requests through the panel with SSRF protection: by default only the `https` scheme is allowed, private and service IPs are blocked, and the response body is limited to 10 MB. The policy is configured with the `PLUGINS_HTTP_*` environment variables (see [Installation and management](/en/plugins/management.html)).
* `gameap-storage` is the plugin's persistent storage, isolated by `plugin_id`, with optional binding of records to an entity (`entity_type`, `entity_id`). Use it for plugin settings: the separate configuration mechanism (`config` in `InitializeRequest`) is not used in the current version of the panel. Storage is bounded per plugin by `PLUGINS_STORAGE_MAX_KEYS_PER_PLUGIN` (10000), `PLUGINS_STORAGE_MAX_VALUE` (1M) and `PLUGINS_STORAGE_MAX_TOTAL` (64M). Payloads are stored in plaintext — credentials belong in `gameap-secrets`.
* `gameap-cache` values are capped by `PLUGINS_CACHE_MAX_VALUE` (1M); entries expire by TTL and are not deleted when the plugin is uninstalled.
* `gameap-nodefs` and `gameap-nodecmd` work with files and commands on the dedicated server (node) through GameAP Daemon; every path is checked against the [node path policy](#node-path-policy).
* `gameap-nodes` writes (`update_node`, `delete_node`, the setup keys) require the `manage_nodes` grant, checked by the module itself; reads are open to every plugin.
* `gameap-ssh` is registered only with `PLUGINS_SSH_ENABLED=true` (default `false`), `gameap-net` only with `PLUGINS_NET_ENABLED=true` (default `true`). A plugin that imports a module the panel did not register fails to load (see [Backward compatibility](#backward-compatibility)).

### Reading node files in chunks

`download` accepts an `offset`/`length` window. Without a window it is a whole-file read capped by `PLUGINS_NODEFS_MAX_INLINE` (32M): a larger file is refused with an error naming both sizes. With a window only the window has to fit the cap (a `length` above it is refused, not clamped), so a large file is read a piece at a time. The response echoes `offset` and carries `total_size`; `length = 0` with an `offset` returns as much as the cap allows, so advance by `len(content)`, not by the requested length. A caller pages while `offset + len(content) < total_size`; reading at or past the end of the file answers empty content rather than an error. `upload` payloads are capped by the same variable.

### `gameap-secrets`

An encrypted store for the plugin's own credentials (API keys, bot tokens), which the plaintext `gameap-storage` is not suited for. Every function requires the `secrets` grant; without it `set`/`delete` answer `success = false`, `get` answers `found = false` and `list_keys` an empty list, with the missing permission in `error` — the module stays importable.

* Values are encrypted with the panel's `ENCRYPTION_KEY` (AES-256-GCM); the ciphertext is bound to the owning plugin id and key, so a row copied elsewhere in the database no longer decrypts.
* Without `ENCRYPTION_KEY` a write is refused rather than stored in plaintext; `PLUGINS_SECRETS_REQUIRE_ENCRYPTION=false` opts out.
* Keys must match `^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,63}$`; quotas: `PLUGINS_SECRETS_MAX_KEYS_PER_PLUGIN` (64) and `PLUGINS_SECRETS_MAX_VALUE` (8K).
* Secrets are private to the plugin that wrote them and are deleted together with its `gameap-storage` entries when the plugin is uninstalled.

### `gameap-ssh`

Opens SSH connections to hosts the plugin names itself and runs commands there — the one thing `gameap-nodecmd` cannot do, because it works through GameAP Daemon and a machine being provisioned does not have one yet. The module requires the `ssh` grant **and** `PLUGINS_SSH_ENABLED=true` on the panel (default `false`): with the switch off the module is not registered at all, and a plugin that imports it fails to load.

* `generate_key_pair` creates a key pair of the requested `KeyType`: `KEY_TYPE_ED25519` (also what an unset type means), `KEY_TYPE_RSA_4096` or `KEY_TYPE_ECDSA_P256`. Every function of the module, key generation included, needs the grant.
* A host key policy is mandatory in `connect`: either `accept_any` or a pin list (SHA-256 fingerprints or public keys), never both — the combination is rejected. The observed host key is always returned, so a first contact can pin it. An operator can forbid `accept_any` panel-wide with `PLUGINS_SSH_ALLOW_ACCEPT_ANY_HOST_KEY=false` (default `true`).
* `exec` blocks within the guest call deadline; when the budget runs out it answers with `completed = false` and an `operation_id`, and the plugin is subscribed to the completion callback. `start_exec` returns the `operation_id` immediately; completion is pushed into the plugin when it exports the `SSHExecEventsHandler` service, and `get_exec_operation` polls the status and output; `cancel_exec` cancels a running command.
* `write_file`/`read_file` stream through a remote `cat`, so no SFTP subsystem is needed on a freshly installed machine.
* Connections and operations live in the memory of the panel instance that created them: they do not survive a plugin reload, and on a multi-instance panel `get_exec_operation` on another instance answers `found = false`. `disconnect` cancels every operation still running on the connection; the plugin's connections are released when it is unloaded.
* Targets go through the same address policy as `gameap-http`: `PLUGINS_SSH_BLOCK_PRIVATE_IPS` (default `true`) and `PLUGINS_SSH_ALLOWED_HOSTS`. Output per stream is capped by `PLUGINS_SSH_MAX_OUTPUT_BYTES` (1 MB) — the head is kept and the stream is flagged truncated.

## Plugin permissions

Privileged host functions are gated on the plugin's grants — the permissions an administrator allowed it. The valid names:

| Permission | Gates |
|---|---|
| `manage_servers` | `gameap-servercontrol` (every function), `gameap-daemontasks.create_daemon_task`, `gameap-servers.save_server`/`delete_server`, `gameap-serversettings.save_server_setting` |
| `node_commands` | `gameap-nodecmd.execute_command`; `cmdexec` daemon tasks additionally to `manage_servers` |
| `files_read` | `gameap-nodefs` reads: `read_dir`, `download`, `get_file_info`, `hash`, `get_archive_operation`; `HTTPResponse.file` |
| `files` | Every `gameap-nodefs` function, including writes and archive operations; includes `files_read` |
| `listen_events` | Event subscriptions |
| `manage_rbac` | Every `gameap-rbac` function |
| `secrets` | Every `gameap-secrets` function |
| `ssh` | Every `gameap-ssh` function, key generation included; the panel must also have `PLUGINS_SSH_ENABLED=true` |
| `manage_nodes` | `gameap-nodes` writes: `update_node`, `delete_node`, `create_setup_key`, `get_setup_key`, `revoke_setup_key` |
| `manage_games`, `manage_game_mods`, `manage_users` | Declared, but not required by any host function yet |

Read-only modules (`gameap-users`, `gameap-games`, `gameap-gamemods`, `gameap-authz`, `gameap-servers.find_servers`/`get_server`, `gameap-nodes.find_nodes`/`get_node`, …) and `gameap-http`, `gameap-cache`, `gameap-storage`, `gameap-crypto`, `gameap-log`, `gameap-scheduler`, `gameap-net`, `gameap-host` need no grant.

A plugin declares what it needs in `PluginInfo.required_permissions`. Installing — by upload, from the store or via `PLUGINS_AUTOLOAD` — grants exactly the declared permissions; unknown names are dropped. Updating never widens the grants: the new declaration is recorded in `required_permissions` while `allowed_permissions` is left as it was, so a build that starts requiring more is refused those calls until an administrator grants them in **Administration** → **Plugins** (the **Permissions** action of the plugin row) or with `PUT /api/admin/plugins/{id}/permissions`.

Every privileged call goes through the guard: the grant is checked first, then the per-plugin rate limit (a token bucket per class, see [Runtime limits](#runtime-limits)), and node paths are checked against the path policy. A refused call answers in the response's `error` field — `plugin permission <name> required`, `rate limited: gameap-<module> allows N calls/s (burst M)` or `path policy: <reason>: <path>` — and is recorded in the audit log as `access.denied` (rate limits as `plugin.hostcall.ratelimited`) with the plugin as the actor. The module stays loaded; a plugin is never disabled for a refused call.

Enforcement is transitional: with `PLUGINS_PERMISSIONS_ENFORCE=false` (the default in 4.5.0) the guard's permission checks pass, while the grants are still recorded and shown in the interface. The rate limits, the node path policy, the `PLUGINS_SSH_ENABLED` switch and the `manage_nodes` check inside `gameap-nodes` always apply. Declare and grant the permissions your plugin needs now, so nothing breaks once enforcement is switched on.

> Grants are the enforcement mechanism, but a plugin still runs inside the panel with access to its data. That is why installing plugins is entrusted to administrators only — install plugins only from sources you trust.

### Node path policy

Every path a plugin hands to `gameap-nodefs`, the `work_dir` of `gameap-nodecmd` and `HTTPResponse.file` is checked on the panel before it reaches the daemon. In every mode a path with a `..` segment or a NUL byte is refused. `PLUGINS_NODEFS_PATH_POLICY` selects what else is allowed:

| Mode | Allowed paths |
|---|---|
| `unrestricted` (default) | Anything the daemon itself permits |
| `node_workpath` | Inside the node's work path; relative paths are resolved under it |
| `server_dirs` | Inside the directory of a game server on that node (`<work_path>/<server dir>`) |

In both restricted modes the `work_dir` of `gameap-nodecmd` must be an absolute path — a relative one would be resolved by the daemon against its own working directory. A command that names no `work_dir` runs in the daemon's default directory and is not checked.

Both restricted modes keep one directory open: the plugin's own service directory `<work_path>/.plugins/<compact plugin id>`, resolved from each node's work path and opened only to the plugin whose id names it. Put node-side scratch files (staged requests, results of long commands, archives built for download) there and the strictest mode still works. `PLUGINS_NODEFS_ALLOWED_PATHS` (comma-separated absolute roots) widens the restricted modes on every node. A refused call answers `path policy: <reason>: <path>` (`403` for file references) and is audited; the plugin is never disabled for it.

## Runtime limits

| Limit | Value |
|---|---|
| Concurrent calls to one plugin | 1 (calls are serialized) |
| Plugin call timeout | 30 s (module start — 60 s) |
| Event handler timeout | 10 s |
| Asynchronous event delivery | 60 s per event, at most 64 concurrent deliveries |
| Module linear memory | `PLUGINS_RUNTIME_MAX_MEMORY`, 256M (a larger declared maximum is clamped; a module whose initial memory exceeds the cap fails to load) |
| Size of a `.wasm` file | `PLUGINS_RUNTIME_MAX_MODULE_SIZE`, 128M; an HTTP upload is additionally capped at 100 MB |
| Body of an HTTP request to a plugin | 1 MB |
| Body of a `gameap-http` response | 10 MB |
| `gameap-nodefs` inline `upload`/`download` | `PLUGINS_NODEFS_MAX_INLINE`, 32M |
| `gameap-storage` | 10000 keys, 1M per value, 64M in total per plugin |
| `gameap-cache` value | `PLUGINS_CACHE_MAX_VALUE`, 1M |
| `gameap-secrets` | 64 keys, 8K per value |

Sizes accept a unit suffix (`512K`, `64M`, `1G`); a plain number is bytes.

The expensive host functions are rate limited per plugin with a token bucket, by class:

| Class | Functions | Default |
|---|---|---|
| `nodecmd` | `gameap-nodecmd.execute_command` | 5/s, burst 20 |
| `servercontrol` | `gameap-servercontrol.*`, `create_daemon_task`, `save_server`, `delete_server`, `save_server_setting` | 5/s, burst 20 |
| `nodefs` | Every `gameap-nodefs` function | 50/s, burst 200 |
| `http` | `gameap-http.fetch` | 20/s, burst 50 |
| `rbac` | Every `gameap-rbac` function | 10/s, burst 50 |
| `ssh` | Every `gameap-ssh` function | 20/s, burst 60 |

`PLUGINS_RATELIMIT_<CLASS>_RPS` / `_BURST` tune them; RPS `0` disables a class. Buckets are per panel instance.

When a call timeout is exceeded (an event handler, an HTTP route, a scheduled task, a callback) or the guest terminates its own module (a panic ending in `proc_exit`), the runtime closes the module: the plugin gets status `error` with the reason in `last_error`, and the panel reloads it automatically after `PLUGINS_RECOVERY_INITIAL_DELAY` (30s), doubling the wait on every further failure up to `PLUGINS_RECOVERY_MAX_DELAY` (10m). After `PLUGINS_RECOVERY_MAX_ATTEMPTS` (5) consecutive reloads the plugin stays in status `error` until an administrator reloads it (**Reload** in **Administration** → **Plugins**, or `POST /api/admin/plugins/{id}/reload`) or the panel restarts; `PLUGINS_RECOVERY_ENABLED=false` keeps the disable permanent. Write the plugin accordingly: keep state in `gameap-storage`, not in module globals, and make `Initialize` idempotent — a reload starts a fresh module instance.

Guest stdout is forwarded to the panel log at debug level and stderr at warn level, so a panic message is visible; lines are cut at 4 KiB and each stream is limited to 200 lines per 10 seconds (drops are counted and reported). The file system and the network are not available from WASM — only through host functions.

## Building a plugin with Rust

The structure of a minimal project (following [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor)):

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

A minimal `src/lib.rs` — a plugin that only returns metadata and an embedded frontend:

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

Building:

```bash
cargo build --target wasm32-wasip1 --release
# optional — size reduction (binaryen):
wasm-opt -Oz target/wasm32-wasip1/release/my_plugin.wasm -o my-plugin.wasm
```

The resulting `.wasm` file is installed through the panel interface or copied into the `plugins/` directory (see [Installation and management](/en/plugins/management.html)). The frontend is built separately and embedded into the `.wasm` (see [Plugin frontend](/en/plugins/frontend.html)).

## Development in AssemblyScript

There is no ready-made SDK for AssemblyScript yet: the `@gameap/proto-as` package provides only message types, and the ABI layer (packing pointers, error handling, exporting `malloc`/`free`) has to be written by hand. Because of AssemblyScript's garbage collector, buffers passed to the host must be pinned (`__pin`/`__unpin`), and a JSON parser and basic utilities have to be implemented yourself or pulled in from third-party libraries. A working AssemblyScript plugin example is [plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth).

## Validating a plugin before installation

The panel can validate a `.wasm` file without installing it — the `POST /api/admin/plugins/upload/dry-run` endpoint (multipart, `file` field). The response returns the plugin metadata, HTTP routes, server permissions, event subscriptions, whether a frontend is present (`has_frontend_bundle`, `frontend_bundle_size`, `has_frontend_styles`) and a list of errors. It also reports the permissions: `required_permissions` (declared in the manifest), `used_permissions` (derived from the module's host imports and its event subscriptions) and `undeclared_permissions` — used but not declared, so the install will refuse those calls. If a plugin with the same id is already installed, `installed`, `installed_version` and `installed_source_type` say so: uploading the file would replace it. The same check runs in the interface when a file is uploaded.

## Backward compatibility

A plugin compiled against an older panel keeps loading and working on newer panels: the panel's CI runs the `pkg/plugin/compatrust` fixture matrix against `HEAD`, `v4.4.1` and `v4.3.5`, and a green run is mandatory for any change to the plugin contracts (`pkg/plugin/proto`, `pkg/plugin/sdk`).

Conversely, a plugin that imports a host module the panel does not provide is rejected at load time with an error naming the module, rather than failing at call time. On a 4.5.0 panel the usual cause is `gameap-ssh` with `PLUGINS_SSH_ENABLED=false` or `gameap-net` with `PLUGINS_NET_ENABLED=false` — check the panel configuration first when a load fails with an unknown import. At runtime a plugin can ask the panel what it got with `gameap-host.get_host_info` (the panel version, the Plugin API version, the instance id and the host modules instantiated for it) and `get_grants` (its current permissions).

## Plugin examples

* [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) — a minimal Rust plugin: metadata plus an embedded frontend, no host functions.
* [plugin-goldsrc-addons](https://github.com/gameap/plugin-goldsrc-addons) — a Rust plugin with HTTP routes, file operations on a node (`nodefs`), command execution (`nodecmd`) and panel API calls (RCON) from the frontend.
* [plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth) — an AssemblyScript plugin: calls to an external API (modrinth.com) through `gameap-http`, caching, persistent storage.
