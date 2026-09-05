---
title: Publicación en el catálogo
layout: default
lang: es
category: Plugins
order: 344
---

## El catálogo de plugins

El catálogo de plugins de GameAP está disponible en [plugins.gameap.dev](https://plugins.gameap.dev/) (versión en inglés) y [plugins.gameap.ru](https://plugins.gameap.ru/) (versión en ruso). La parte pública del catálogo incluye:

* la lista de plugins con ordenación (populares, mejor valorados, nuevos, actualizados recientemente);
* filtrado por categorías y etiquetas;
* la página del plugin: descripción general, reseñas, registro de cambios por versión, capturas de pantalla, información de licencia y un enlace al repositorio.

No hay botón «Install» en el sitio: los plugins se instalan desde la interfaz del panel (consulte [Instalación y gestión](/es/plugins/management.html)).

El panel del desarrollador está en [plugins.gameap.dev/dashboard](https://plugins.gameap.dev/dashboard).

## Registro

El registro se realiza por correo electrónico, con un nombre de usuario y una contraseña (de al menos 8 caracteres) y confirmación por correo electrónico, o mediante OAuth (GitHub, Google).

## Creación de un plugin en el panel del desarrollador

En el panel del desarrollador, cree un plugin y complete los campos:

* nombre;
* descripción corta (no más de 500 caracteres);
* descripción;
* categoría y etiquetas;
* licencia;
* enlace al repositorio;
* versiones mínimas de GameAP y de la Plugin API.

El icono del plugin se sube después de la creación, en la página de edición.

## Publicación de una versión

En la página del plugin, añada una nueva versión:

* versión — en formato de versionado semántico (por ejemplo, `1.0.0`);
* el archivo del plugin — solo `.wasm`;
* una firma GPG del archivo (`.asc`, opcional);
* un registro de cambios («qué hay de nuevo en esta versión»);
* el indicador «stable release» — marca la versión como recomendada para su uso;
* capturas de pantalla de la versión.

La descripción, el registro de cambios y las capturas de pantalla se pueden traducir a otros idiomas.

## Moderación

Un plugin nuevo pasa por moderación: «draft» → «under review» → «published» o «rejected». Cuando un plugin o una versión es rechazado, el moderador indica el motivo y el desarrollador recibe una notificación.

## Publicación desde CI/CD

Para publicar versiones automáticamente se utilizan deploy tokens: se crean en la página del plugin en el panel del desarrollador y se muestran solo una vez. La publicación se realiza con una petición al endpoint `https://plugins.gameap.dev/api/ci/plugins/{PLUGIN_ID}/versions`:

```http
POST /api/ci/plugins/{PLUGIN_ID}/versions
Authorization: Bearer <deploy token>
```

Campos multipart: `version`, `file` (el archivo `.wasm`), `signature` (el archivo `.asc`, opcional), `changelog`, `is_stable`. Una respuesta exitosa es HTTP 201.

La firma GPG del archivo se crea con:

```bash
gpg --detach-sign --armor -o my-plugin.wasm.asc my-plugin.wasm
```

Para un ejemplo listo de un flujo de trabajo de publicación por etiquetas de release (compilación, firma, subida al catálogo), consulte el archivo `release.yml` en los repositorios [plugin-minecraft-modrinth](https://github.com/gameap/plugin-minecraft-modrinth) y [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor).

## Requisitos del plugin

* Un único archivo `.wasm` (target `wasm32-wasip1`); el frontend está integrado en el mismo archivo.
* Metadatos `PluginInfo` válidos: `api_version` con valor `"1"` y un `id` estable son obligatorios (consulte los requisitos del id en [Desarrollo de plugins](/es/plugins/development.html)).
* Versiones en formato de versionado semántico.

Cuando un plugin se instala desde el catálogo, el panel verifica el hash SHA-256 del archivo descargado. La firma GPG no es verificada por el panel — está ahí para que los usuarios puedan comprobar la autenticidad del archivo por sí mismos.
