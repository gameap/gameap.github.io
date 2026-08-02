---
title: Plugin development
layout: default
lang: en
category: Plugins
order: 342
---

## Architecture

A GameAP plugin is a WASM module for the `wasm32-wasip1` target (WASI preview 1), built as a reactor. The panel runs each plugin in an isolated [wazero](https://github.com/tetratelabs/wazero) runtime: no file system, no network, no environment variables; stdout and stderr are discarded.

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
| `required_permissions` | Declared permissions (not checked in the current version of the panel) |
| `api_version` | Plugin API version, must be `"1"` |

**Requirements for `id`.** Use a stable identifier made of base32 alphabet characters `a-z2-7`, without hyphens. The panel normalizes the id: a string with hyphens or other characters is replaced by a hash, which breaks the `/api/plugins/{id}/...` and `/plugins/{id}/...` paths. Avoid purely numeric ids as well: such an identifier is treated as a decimal numeric ID (id parsing first tries to parse the string as a number). Examples of valid ids from real plugins: `hexeditor4jm2`, `ezvdsxmlu6fbk`, `dshdabjp2l73a`.

## Events

A plugin subscribes to panel events with the `GetSubscribedEvents` method. The event types (listed without the `EVENT_TYPE_` prefix):

| Event | Cancellable | Delivery |
|---|---|---|
| `SERVER_PRE_START`, `SERVER_PRE_STOP`, `SERVER_PRE_RESTART`, `SERVER_PRE_INSTALL`, `SERVER_PRE_UPDATE`, `SERVER_PRE_REINSTALL`, `SERVER_PRE_DELETE` | Yes | Synchronous |
| `SERVER_POST_START`, `SERVER_POST_STOP`, `SERVER_POST_RESTART`, `SERVER_POST_INSTALL`, `SERVER_POST_UPDATE`, `SERVER_POST_REINSTALL`, `SERVER_POST_DELETE` | No | Asynchronous |
| `SERVER_CREATED`, `SERVER_UPDATED`, `SERVER_DELETED` | No | Asynchronous |
| `DAEMON_TASK_CREATED`, `DAEMON_TASK_COMPLETED`, `DAEMON_TASK_FAILED` | No | Asynchronous |

Pre-events are delivered synchronously and block the operation: a plugin can cancel it by returning an `EventResult` with `should_cancel = true` and a `message`, or modify the data via `modified_data`. Post-events and the other types are delivered asynchronously and do not affect the operation.

A server event contains a full snapshot of the server (`ServerEventPayload`); a task event contains the daemon task data (`TaskEventPayload`). The event handler timeout is 10 seconds; once it expires the plugin is disabled until the panel is restarted.

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

## Host functions

All calls from a plugin to the panel and the outside world go through host functions, grouped into `gameap-*` modules:

| Module | Functions | Purpose |
|---|---|---|
| `gameap-log` | `log` | Writing to the panel log |
| `gameap-cache` | `get`, `set`, `delete` | The panel's shared cache (keys prefixed with `plugin:`) |
| `gameap-crypto` | `random_uint64`, `random_string`, `argon2_hash`, `argon2_verify` | Random value generators, Argon2 hashing |
| `gameap-http` | `fetch` | Outbound HTTP requests with SSRF protection |
| `gameap-storage` | `get`, `set`, `delete`, `list` | The plugin's persistent key-value storage |
| `gameap-servercontrol` | `start_server`, `stop_server`, `restart_server`, `update_server`, `install_server`, `reinstall_server` | Game server control (returns a `task_id`) |
| `gameap-servers` | `find_servers`, `get_server`, `save_server`, `delete_server` | Game servers |
| `gameap-users` | `find_users`, `get_user` | Panel users |
| `gameap-nodes` | `find_nodes`, `get_node` | Dedicated servers (nodes) |
| `gameap-games` | `find_games`, `get_game` | Games |
| `gameap-gamemods` | `find_game_mods`, `get_game_mod` | Game mods |
| `gameap-daemontasks` | `find_daemon_tasks`, `create_daemon_task` | Daemon tasks |
| `gameap-serversettings` | `find_server_settings`, `save_server_setting` | Game server settings |
| `gameap-nodefs` | `read_dir`, `mk_dir`, `copy`, `move`, `download`, `upload`, `remove`, `get_file_info`, `chmod` | File operations on a node |
| `gameap-nodecmd` | `execute_command` | Executing a command on a node |

Details:

* `gameap-http` proxies requests through the panel with SSRF protection: by default only the `https` scheme is allowed, private and service IPs are blocked, and the response body is limited to 10 MB. The policy is configured with the `PLUGIN_HTTP_*` environment variables (see [Installation and management](/en/plugins/management.html)).
* `gameap-storage` is the plugin's persistent storage, isolated by `plugin_id`, with optional binding of records to an entity (`entity_type`, `entity_id`). Use it for plugin settings: the separate configuration mechanism (`config` in `InitializeRequest`) is not used in the current version of the panel.
* `gameap-nodefs` and `gameap-nodecmd` work with files and commands on the dedicated server (node) through GameAP Daemon.

**Important:** host functions run with the panel's own privileges, without additional checks. That is why installing plugins is entrusted to administrators only — install plugins only from sources you trust.

## Runtime limits

| Limit | Value |
|---|---|
| Concurrent calls to one plugin | 1 (calls are serialized) |
| Plugin call timeout | 30 s (module start — 60 s) |
| Event handler timeout | 10 s |
| Size of an uploaded `.wasm` file | 100 MB |
| Body of an HTTP request to a plugin | 1 MB |
| Body of a `gameap-http` response | 10 MB |

When a call timeout is exceeded, the panel disables the plugin until a restart. The file system and the network are not available from WASM — only through host functions.

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

The panel can validate a `.wasm` file without installing it — the `POST /api/admin/plugins/upload/dry-run` endpoint (multipart, `file` field). The response returns the plugin metadata, HTTP routes, server permissions, event subscriptions, whether a frontend is present, and a list of errors. The same check runs in the interface when a file is uploaded.

## Plugin examples

* [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) — a minimal Rust plugin: metadata plus an embedded frontend, no host functions.
* [plugin-goldsrc-addons](https://github.com/gameap/plugin-goldsrc-addons) — a Rust plugin with HTTP routes, file operations on a node (`nodefs`), command execution (`nodecmd`) and panel API calls (RCON) from the frontend.
* [plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth) — an AssemblyScript plugin: calls to an external API (modrinth.com) through `gameap-http`, caching, persistent storage.
