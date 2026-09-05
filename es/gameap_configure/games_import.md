---
title: Importación de juegos
layout: default
lang: es
category: Configuración del panel
order: 330
---

La configuración de juegos y mods se puede trasladar entre instalaciones de GameAP y tomarse
de otros paneles de control. Esta funcionalidad apareció en GameAP 4.1.

Todo esto se encuentra en la página **Administración** → **Juegos**.

![Botones de importación y exportación en la página Juegos de la administración](/images/en/gameap_configure/games_import/import_button.png)

## Formato propio de GameAP

El formato de GameAP contiene el juego completo: el juego en sí, todos sus mods, variables,
comandos de ejecución y comandos RCON.

### Exportación

En la página **Administración** → **Juegos**, seleccione un juego y haga clic en **Exportar**.
El panel servirá un archivo `<game-code>.gameap.yaml`.

A través de la API:

```http
GET /api/games/{code}/export
```

### Importación

**Administración** → **Juegos** → **Importar YAML de GameAP**, suba el archivo y haga clic
en **Importar**.

![Página de importación de juegos con el selector del archivo de configuración YAML](/images/en/gameap_configure/games_import/import_page.png)

A través de la API:

```http
POST /api/games/import/gameap
```

La importación crea el juego si aún no existe y **actualiza el existente** si ya hay un juego
con ese código. Los mods se importan junto con el juego: después de la subida, el panel
informa cuántos se importaron.

La versión del formato se indica en la parte superior del archivo:

```yaml
schema_version: "1.0"
```

El panel no aceptará un archivo con una versión de formato diferente. También conviene
revisar este campo cuando una transferencia falla.

## Importación desde otros paneles

La importación de plantillas es compatible con:

* [Pterodactyl](https://pterodactyl.io/)
* [Pelican](https://pelican.dev/)

**Administración** → **Juegos** → **Importar Egg de Pelican**, suba el archivo egg en formato
JSON o YAML.

A través de la API:

```text
POST /api/games/import/pelican-egg
```

Ambos paneles comparten el mismo formato egg, por lo que los archivos de Pterodactyl y
Pelican se suben de la misma manera.

## Actualización de juegos desde el catálogo de GameAP

El panel incluye un catálogo de configuraciones listas para usar de juegos y mods. Para
obtener los cambios del catálogo, haga clic en **Actualizar juegos** en la página
**Administración** → **Juegos**.

A través de la API:

```text
POST /api/games/upgrade
```

> **Sus cambios pueden sobrescribirse.** La actualización lleva la configuración de juegos y
> mods al estado del catálogo. Si ha editado comandos de ejecución o variables de los juegos
> estándar, guárdelos antes de actualizar; por ejemplo, exportando el juego a un archivo.

Un mod no se actualiza si el panel encuentra varios mods con el mismo nombre para el juego:
estos casos se omiten para evitar sobrescribir el incorrecto.

Las direcciones del catálogo se configuran con la variable `GAMES_CDN_URLS`; se prueban en
orden hasta que una responda. Consulte la [referencia de config.env](/es/config.html).

## Funciones de Pelican y Pterodactyl

Trabajar con Eggs de Pelican y Eggs de Pterodactyl importados solo es posible con los
[gestores de procesos](/es/daemon/process_managers.html) Docker y Podman.

Debe configurar GameAP Daemon para trabajar con uno de estos gestores de procesos.
Para ello, al agregar un nuevo nodo, seleccione el gestor de procesos deseado en la sección
"Configuración avanzada".

![Elección del gestor de procesos Docker al agregar un servidor dedicado](/images/en/gameap_configure/games_import/daemon_process_manager.png)
