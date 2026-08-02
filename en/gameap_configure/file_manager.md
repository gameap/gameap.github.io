---
title: File Manager
layout: default
lang: en
category: Panel settings
order: 317
---

The file manager works with game server files right from the panel: viewing and editing,
uploading and downloading, permissions. It lives on the **Files** tab of the game server page.

All operations are performed through GameAP Daemon on the dedicated server and are confined to
the game server directory — there is no way to escape it.

Access requires the `game-server-files` permission on the server, see
[Users, Roles, and Permissions](/en/users.html).

## Uploading Files

The panel uses two upload methods and picks between them itself.

**Regular upload** — for small files, in a single request. The maximum size is **100 MB** and
is not configurable.

**Chunked upload** — for large files. The file is split into chunks of
`FILES_UPLOAD_CHUNK_SIZE` (8 MB by default); chunks are sent in parallel, four at a time, with
three attempts each.

The maximum file size for a chunked upload is the chunk size multiplied by
`FILES_UPLOAD_MAX_CHUNKS` (100,000 by default), which is about 780 GB.

### Resuming

Before sending, the browser computes the SHA-256 checksum of the whole file and passes it to
the panel along with the upload parameters. The panel tracks which chunks have been received.

If the upload is interrupted — the tab closed, the network dropped — on retry the panel reports
the list of missing chunks, and only those are sent. The whole file does not have to be sent
again.

Unfinished uploads live for `FILES_UPLOAD_SESSION_TTL` (24 hours by default), after which they
are removed by a background cleanup that runs every `FILES_UPLOAD_JANITOR_INTERVAL` (12 hours
by default).

> The checksum is computed in the browser via WebAssembly. If you have edited the CSP policy
> and removed `'wasm-unsafe-eval'` from `script-src`, file uploads will stop working.
> See [Security](/en/security.html).

### Uploading Directories

Besides individual files, a whole directory can be uploaded — the nested subdirectory structure
is preserved.

### Which Files Are Allowed

The file type is detected **from the content**, not from the extension and not from what the
browser reports. Allowed:

* images: PNG, JPEG, GIF, WebP, BMP, ICO;
* text files, including configuration files — they are detected as `text/plain`;
* JSON, XML, CSV, YAML;
* PDF.

**Not allowed by default:**

| What                                   | Variable                      | Why                                                                                |
|----------------------------------------|-------------------------------|------------------------------------------------------------------------------------|
| Archives: zip, tar, gzip, bzip2, 7z, xz | `FILES_UPLOAD_ALLOW_ARCHIVES` | An archive may contain executables that will be unpacked on the dedicated server   |
| Arbitrary binary files                 | `FILES_UPLOAD_ALLOW_BINARY`   | The `application/octet-stream` type gives no clue what is inside                   |

SVG and HTML are always blocked: a script can be hidden in them.

The list of allowed types can be extended, without opening up archives entirely, with the
`FILES_UPLOAD_ALLOWED_MIMES` variable — it **extends** the default list rather than replacing
it.

Rejected uploads go to the audit log with the detected file type and the reason for rejection.

## Downloading

A single file is downloaded as is. A directory or several selected files are packed into a
**ZIP** and streamed in one pass, without creating a temporary archive on disk.

The limits are set with variables:

| Variable                               | Default  | Purpose                                        |
|----------------------------------------|----------|-------------------------------------------------|
| `FILES_ARCHIVE_MAX_BYTES`              | `100G`   | Maximum archive size                            |
| `FILES_ARCHIVE_MAX_FILES`              | `500000` | Maximum number of files in an archive           |
| `FILES_ARCHIVE_CONCURRENT_PER_SERVER`  | `2`      | Concurrent packing operations per game server   |

## File Permissions

File and directory permissions can be changed with **chmod** — just like in a shell. This is
needed, for example, to make a game server startup script executable.

## Editing

Text files open in the built-in editor. Plugins can add their own editors for particular
files — for example, a hex editor for binaries. Files larger than 1 MB are not opened by plugin
editors. See [Plugins](/en/plugins/index.html).

## If a File Does Not Upload

| Symptom                                  | Cause                                                                     |
|------------------------------------------|----------------------------------------------------------------------------|
| Rejected immediately, without uploading  | The file type is not allowed. Most often it is an archive or a binary file |
| Rejected on a file larger than 100 MB    | The regular upload limit kicked in; large files use chunked upload         |
| The upload breaks off at checksum computation | `'wasm-unsafe-eval'` has been removed from the CSP policy             |
| The upload does not start                | The `game-server-files` permission is missing, or the daemon is unreachable |

The rejection details are visible in the panel log, see
[Troubleshooting](/en/troubleshooting.html).
