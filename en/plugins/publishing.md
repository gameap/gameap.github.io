---
title: Publishing to the catalog
layout: default
lang: en
category: Plugins
order: 344
---

* This will become a table of contents (this text will be scraped).
{:toc}

## The plugin catalog

The GameAP plugin catalog is available at [plugins.gameap.dev](https://plugins.gameap.dev/) (English version) and [plugins.gameap.ru](https://plugins.gameap.ru/) (Russian version). The public part of the catalog includes:

* the plugin list with sorting (popular, highly rated, new, recently updated);
* filtering by categories and tags;
* the plugin page: overview, reviews, per-version changelog, screenshots, license information and a link to the repository.

There is no "Install" button on the site: plugins are installed from the panel interface (see [Installation and management](/en/plugins/management.html)).

The developer dashboard is at [plugins.gameap.dev/dashboard](https://plugins.gameap.dev/dashboard).

## Registration

Registration is by email, with a username and a password (at least 8 characters) and email confirmation, or through OAuth (GitHub, Google).

## Creating a plugin in the dashboard

In the developer dashboard, create a plugin and fill in the fields:

* name;
* short description (no more than 500 characters);
* description;
* category and tags;
* license;
* link to the repository;
* minimum GameAP and Plugin API versions.

The plugin icon is uploaded after creation, on the edit page.

## Publishing a version

On the plugin page, add a new version:

* version — in semantic versioning format (for example, `1.0.0`);
* the plugin file — `.wasm` only;
* a GPG signature of the file (`.asc`, optional);
* a changelog ("what's new in this version");
* the "stable release" flag — marks the version as recommended for use;
* screenshots of the version.

The description, changelog and screenshots can be translated into other languages.

## Moderation

A new plugin goes through moderation: "draft" → "under review" → "published" or "rejected". When a plugin or a version is rejected, the moderator states the reason and the developer receives a notification.

## Publishing from CI/CD

Deploy tokens are used to publish versions automatically: they are created on the plugin page in the dashboard and are shown only once. Publishing is done with a request to the `https://plugins.gameap.dev/api/ci/plugins/{PLUGIN_ID}/versions` endpoint:

```http
POST /api/ci/plugins/{PLUGIN_ID}/versions
Authorization: Bearer <deploy token>
```

Multipart fields: `version`, `file` (the `.wasm` file), `signature` (the `.asc` file, optional), `changelog`, `is_stable`. A successful response is HTTP 201.

A GPG signature of the file is created with:

```bash
gpg --detach-sign --armor -o my-plugin.wasm.asc my-plugin.wasm
```

For a ready-made example of a release-tag publishing workflow (build, sign, upload to the catalog), see the `release.yml` file in the [plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth) and [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) repositories.

## Plugin requirements

* A single `.wasm` file (the `wasm32-wasip1` target); the frontend is embedded into the same file.
* Valid `PluginInfo` metadata: `api_version` set to `"1"` and a stable `id` are mandatory (see the id requirements in [Plugin development](/en/plugins/development.html)).
* Versions in semantic versioning format.

When a plugin is installed from the catalog, the panel verifies the SHA-256 hash of the downloaded file. The GPG signature is not verified by the panel — it is there so that users can check the file's authenticity themselves.
