---
title: Frontend de plugins
layout: default
lang: es
category: Plugins
order: 343
---

## Cómo se integra el frontend

Un plugin puede contener una interfaz: componentes de Vue 3 empaquetados en JS y CSS e integrados en el mismo archivo `.wasm`. El panel sirve los bundles de frontend concatenados de todos los plugins cargados en `/plugins.js` y `/plugins.css` (solo para usuarios autenticados). El cargador del panel importa este código como un módulo y registra cada objeto `PluginDefinition` exportado.

El frontend de un plugin se ejecuta en el contexto de la SPA principal del panel (no es un iframe ni un web component) y utiliza las bibliotecas del panel a través de objetos globales:

| Objeto global | Contenido |
|---|---|
| `window.Vue` | Vue 3 |
| `window.VueRouter` | Vue Router |
| `window.Pinia` | Pinia |
| `window.axios` | La instancia de axios configurada del panel (con autorización) |
| `window.NaiveUI` | Naive UI |

## El manifiesto `PluginDefinition`

No hay un archivo de manifiesto separado: el manifiesto del frontend es el objeto `PluginDefinition` exportado desde el bundle:

| Campo | Obligatorio | Descripción |
|---|---|---|
| `id` | Sí | El identificador del plugin; debe coincidir con el `id` del `PluginInfo` del backend |
| `name` | Sí | Nombre del plugin |
| `version` | Sí | Versión (semver) |
| `apiVersion` | Sí | Versión de la API del frontend, solo `'1.0'` |
| `description` | No | Descripción |
| `author` | No | Autor |
| `routes` | No | Páginas propias del plugin |
| `menuItems` | No | Elementos del menú izquierdo (sidebar) |
| `slots` | No | Componentes en los slots integrados del panel |
| `homeButtons` | No | Botones en la página de inicio |
| `fileEditors` | No | Editores de archivos para el gestor de archivos |
| `translations` | No | Diccionarios de traducción `{ en: {...}, ru: {...} }` |
| `onInit` | No | Hook de inicialización que se llama cuando se registra el plugin |

## Puntos de integración

| Mecanismo | Dónde aparece |
|---|---|
| `routes` | Las páginas propias del plugin en `/plugins/{id}/...` |
| El slot `server-tabs` | Una pestaña en la página del servidor de juego (junto a "Console", "Files", etc.) |
| El slot `dashboard-widgets` | Un widget en la página de inicio del panel |
| `homeButtons` | Botones en la página de inicio |
| `menuItems` | Elementos del sidebar (las secciones `servers`, `admin` y `custom`) |
| El slot `admin-user-info` | Un bloque en el diálogo de información del usuario (administración) |
| `fileEditors` | El menú contextual del gestor de archivos: apertura de un archivo en el editor del plugin |

Los slots `sidebar-sections` y `admin-pages` están declarados en el SDK, pero no están integrados en la versión actual del panel.

Detalles:

* Para `server-tabs` está disponible una comprobación de permisos: `checkPermission: { type: 'hasServerPermissions', permissions: [...] }` — la pestaña se muestra solo si el usuario tiene todos los permisos indicados para el servidor (los permisos de plugin tienen la forma `plugin:{id}:...`, por ejemplo `plugin:ezvdsxmlu6fbk:manage`).
* Los editores de archivos se registran con reglas de coincidencia (`fileName`, `extensions`, `pathContains`, `fullPath`, `gameCode` y otras): el editor más específico entre los que coinciden se convierte en el editor predeterminado para el archivo y se abre con un doble clic. Un editor marcado con `contextMenuOnly: true` nunca es el predeterminado y solo se ofrece en el menú contextual, que es lo que necesita un editor que coincide con todos los archivos (`allFiles: true`) para no apropiarse de todas las vistas previas. `menuLabel` sustituye el texto «Edit with …» del elemento del menú, y `checkPermission: { type: 'hasServerPermissions', permissions: [...] }` oculta el elemento a los usuarios que no tengan los permisos de servidor indicados.
* El componente del editor recibe las props `content`, `filePath`, `fileName`, `extension`, `fileSize`, `fileMtime`, `disk`, `pluginId`, `gameCode` y `gameName`, y emite los eventos `save` y `close`; el guardado del archivo en el servidor lo realiza el propio panel.
* El panel descarga el archivo antes de montar el editor, y no lo hace con archivos de más de 1 MB. A un editor que declara `contentType: 'none'` no se le pasa la prop `content` en absoluto: carga por su cuenta lo que necesita, por lo que el límite de tamaño no se le aplica y así se puede abrir un archivo de cualquier tamaño; `fileSize` y `fileMtime` provienen del listado del directorio, de modo que un editor así puede describir el archivo sin descargarlo. Con `contentType: 'text'` (el valor predeterminado) el contenido llega como cadena, y con `'binary'` como `ArrayBuffer`.

## Traducciones

Las traducciones se definen con los diccionarios `translations: { en: {...}, ru: {...} }`. En los campos `label`, `name` y `text` se admiten referencias a claves de traducción de la forma `@:key` — el panel sustituye la cadena correspondiente al idioma actual de la interfaz.

## Acceso a la API

El frontend de un plugin utiliza `window.axios` — la misma instancia que usa el panel, con la autorización del usuario actual. A través de ella se puede acceder a:

* la API del panel — por ejemplo, enviar un comando RCON con `POST /api/servers/{id}/rcon`, o trabajar con archivos a través de `/api/file-manager/...`;
* el backend propio del plugin en `/api/plugins/{id}/...` (las rutas HTTP registradas por la parte WASM).

## El SDK `@gameap/plugin-sdk`

El paquete npm `@gameap/plugin-sdk` proporciona:

* tipos de TypeScript: `PluginDefinition`, `PluginRoute`, `PluginMenuItem`, `PluginSlotComponent`, `PluginHomeButton`, `PluginFileEditor`, `PluginContext` y otros;
* hooks de contexto: `usePluginContext`, `useServer`, `useServerId`, `useServerAbilities`, `useCurrentUser`, `useIsAdmin`, `useIsAuthenticated`, `usePluginRoute`, `usePluginId`;
* hooks de traducción: `usePluginTrans`, `providePluginTrans`;
* los componentes de UI del panel (reexportados desde `@gameap/ui`): `GCard`, `GDataTable`, `GModal`, `GStatusBadge`, `GSwitch` y otros;
* `createPluginConfig` — una configuración de Vite lista para usar: una compilación en modo biblioteca (el módulo ES `plugin.js`) donde las dependencias externas (`vue`, `vue-router`, `pinia`, `axios`, `@gameap/ui`) se reescriben a objetos globales. Nota: el panel expone `window.NaiveUI` pero no `window.gameapUI`, por lo que los plugins reales usan su propia configuración de Vite con externals que apuntan a `naive-ui` (véase `frontend/vite.config.js` en [plugin-hex-editor](https://github.com/gameap/plugin-hex-editor) como ejemplo).

## Compilación del frontend

```bash
npm run build   # → dist/plugin.js (+ un archivo CSS)
```

El `plugin.js` y el CSS compilados se integran en el `.wasm` del plugin: para Rust, con un script `build.rs` que los copia en `OUT_DIR` y los incluye mediante `include_bytes!` (véase [Desarrollo de plugins](/es/plugins/development.html)); para AssemblyScript, con un script de generación de código (véase `scripts/embed-frontend.mjs` en plugin-minecraft-modrinth como ejemplo).

## Depuración local

El paquete `@gameap/debug` inicia un entorno de depuración con el frontend real del panel y una API simulada (MSW):

```bash
PLUGIN_PATH=./dist npx @gameap/debug
```

El entorno se abre en `http://localhost:5174`. Un panel de depuración flotante permite cambiar el tipo de usuario (administrador / usuario normal / invitado), el retardo de red y el idioma. El plugin debe estar compilado previamente (`npm run build`).

## Un ejemplo de `PluginDefinition`

Un ejemplo de una pestaña en la página del servidor de juego (siguiendo plugin-goldsrc-addons):

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
