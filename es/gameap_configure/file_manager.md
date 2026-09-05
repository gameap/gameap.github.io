---
title: Gestor de archivos
layout: default
lang: es
category: Configuración del panel
order: 317
---

El gestor de archivos permite trabajar con los archivos del servidor de juego directamente
desde el panel: visualización y edición, subida y descarga, permisos. Se encuentra en la
pestaña **Archivos** de la página del servidor de juego.

Todas las operaciones se realizan a través de GameAP Daemon en el servidor dedicado y están
limitadas al directorio del servidor de juego: no es posible salir de él.

El acceso requiere el permiso `game-server-files` en el servidor; consulte
[Usuarios, roles y permisos](/es/users.html).

## Subida de archivos

El panel utiliza dos métodos de subida y elige entre ellos automáticamente.

**Subida normal**: para archivos pequeños, en una sola solicitud. El tamaño máximo es de
**100 MB** y no es configurable.

**Subida por fragmentos**: para archivos grandes. El archivo se divide en fragmentos de
`FILES_UPLOAD_CHUNK_SIZE` (8 MB de forma predeterminada); los fragmentos se envían en
paralelo, de cuatro en cuatro, con tres intentos cada uno.

El tamaño máximo de archivo para una subida por fragmentos es el tamaño del fragmento
multiplicado por `FILES_UPLOAD_MAX_CHUNKS` (100 000 de forma predeterminada), lo que equivale
a unos 780 GB.

### Reanudación

Antes del envío, el navegador calcula la suma de verificación SHA-256 del archivo completo y
la pasa al panel junto con los parámetros de subida. El panel registra qué fragmentos se han
recibido.

Si la subida se interrumpe —la pestaña se cerró, la red se cayó—, al reintentarlo el panel
informa la lista de fragmentos faltantes y solo se envían esos. No es necesario volver a
enviar el archivo completo.

Las subidas incompletas se conservan durante `FILES_UPLOAD_SESSION_TTL` (24 horas de forma
predeterminada), tras lo cual son eliminadas por una limpieza en segundo plano que se ejecuta
cada `FILES_UPLOAD_JANITOR_INTERVAL` (12 horas de forma predeterminada).

> La suma de verificación se calcula en el navegador mediante WebAssembly. Si ha editado la
> política CSP y ha eliminado `'wasm-unsafe-eval'` de `script-src`, la subida de archivos
> dejará de funcionar.
> Consulte [Seguridad](/es/security.html).

### Subida de directorios

Además de archivos individuales, se puede subir un directorio completo: la estructura de
subdirectorios anidados se conserva.

### Qué archivos están permitidos

El tipo de archivo se detecta **a partir del contenido**, no de la extensión ni de lo que
informa el navegador. Permitidos:

* imágenes: PNG, JPEG, GIF, WebP, BMP, ICO;
* archivos de texto, incluidos los archivos de configuración: se detectan como `text/plain`;
* JSON, XML, CSV, YAML;
* PDF.

**No permitidos de forma predeterminada:**

| Qué                                    | Variable                      | Por qué                                                                             |
|----------------------------------------|-------------------------------|------------------------------------------------------------------------------------|
| Archivos comprimidos: zip, tar, gzip, bzip2, 7z, xz | `FILES_UPLOAD_ALLOW_ARCHIVES` | Un archivo comprimido puede contener ejecutables que se descomprimirán en el servidor dedicado |
| Archivos binarios arbitrarios          | `FILES_UPLOAD_ALLOW_BINARY`   | El tipo `application/octet-stream` no da ninguna pista sobre su contenido           |

SVG y HTML están bloqueados siempre: pueden ocultar un script.

La lista de tipos permitidos se puede ampliar, sin habilitar por completo los archivos
comprimidos, con la variable `FILES_UPLOAD_ALLOWED_MIMES`: **amplía** la lista
predeterminada en lugar de reemplazarla.

Las subidas rechazadas se registran en el registro de auditoría con el tipo de archivo
detectado y el motivo del rechazo.

## Descarga

Un archivo individual se descarga tal cual. Un directorio o varios archivos seleccionados se
empaquetan en un **ZIP** y se transmiten en una sola pasada, sin crear un archivo temporal
en el disco.

Los límites se establecen con variables:

| Variable                               | Predeterminado | Propósito                                        |
|----------------------------------------|----------------|--------------------------------------------------|
| `FILES_ARCHIVE_MAX_BYTES`              | `100G`         | Tamaño máximo del archivo comprimido             |
| `FILES_ARCHIVE_MAX_FILES`              | `500000`       | Número máximo de archivos en un archivo comprimido |
| `FILES_ARCHIVE_CONCURRENT_PER_SERVER`  | `2`            | Operaciones de empaquetado simultáneas por servidor de juego |

## Permisos de archivos

Los permisos de archivos y directorios se pueden cambiar con **chmod**, igual que en una
terminal. Esto es necesario, por ejemplo, para hacer ejecutable un script de inicio de un
servidor de juego.

## Edición

Los archivos de texto se abren en el editor integrado. Los plugins pueden añadir sus propios
editores para archivos concretos, por ejemplo, un editor hexadecimal para binarios. Los
editores de los plugins no abren archivos de más de 1 MB. Consulte [Plugins](/es/plugins/index.html).

## Si un archivo no se sube

| Síntoma                                  | Causa                                                                     |
|------------------------------------------|----------------------------------------------------------------------------|
| Rechazado inmediatamente, sin subirse    | El tipo de archivo no está permitido. La mayoría de las veces es un archivo comprimido o un binario |
| Rechazado en un archivo de más de 100 MB | Se aplicó el límite de la subida normal; los archivos grandes usan la subida por fragmentos |
| La subida se interrumpe al calcular la suma de verificación | Se ha eliminado `'wasm-unsafe-eval'` de la política CSP |
| La subida no comienza                    | Falta el permiso `game-server-files` o el daemon no está accesible |

Los detalles del rechazo se pueden ver en el registro del panel; consulte
[Solución de problemas](/es/troubleshooting.html).
