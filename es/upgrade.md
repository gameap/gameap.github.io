---
title: Actualización
layout: default
lang: es
category: Instalación de GameAP
order: 190
---

Al instalar GameAP, se instala la utilidad `gameapctl`,
que permite administrar el entorno del panel, incluidas las actualizaciones.

Esta página trata sobre la actualización dentro de la cuarta versión. La migración desde GameAP 3 se describe
por separado: [Actualización de v3 a v4](/es/upgrade_from_v3_to_v4.html).

Haga una copia de seguridad de la base de datos antes de actualizar: `gameapctl` no la guarda, y una actualización puede cambiar
el esquema.

> **Autenticación de dos factores obligatoria para los administradores.**
>
> En GameAP 4 este requisito está activado por defecto. Después de la actualización, los administradores sin 2FA
> verán un recordatorio, y transcurridos 30 días el inicio de sesión ya no otorgará una sesión completa hasta que
> se active 2FA. La cuenta regresiva de cada administrador comienza con su primer inicio de sesión después de la actualización.
>
> Para mantener el recordatorio pero eliminar el bloqueo, defina `AUTH_MFA_HARD_FAIL_DAYS=0` en `config.env`.
> Para desactivar el requisito por completo — `AUTH_REQUIRE_MFA_FOR_ADMINS=false`.
>
> Los detalles y qué hacer si pierde el acceso están en la
> página de [Seguridad](/es/security.html).

## GameAP Web/API

### Linux

Para actualizar el panel, ejecute el comando:
```shell
gameapctl panel upgrade
```

### Windows

Para actualizar el panel en Windows,
puede ejecutar el comando donde esté instalado `gameapctl`:
```powershell
C:\path\to\gameapctl.exe panel upgrade
```

O use la interfaz gráfica. Ejecute `gameapctl.exe` y, en la ventana del navegador que se abra,
haga clic en **"Upgrade"** en la sección Web/API.

![El botón Upgrade en la sección Web/API de la interfaz de gameapctl](/images/en/gameapctl/ui.png)

### Actualización a una versión específica

Por defecto se instala la última versión estable. Para elegir otra, especifíquela mediante su etiqueta:

```shell
gameapctl panel upgrade --version=4.2.0
```

## Actualización de GameAP Daemon

### Linux

Para actualizar el Daemon, ejecute el comando:
```shell
gameapctl daemon upgrade
```

### Windows

Para actualizar GameAP Daemon en Windows, puede ejecutar el comando
donde esté instalado `gameapctl`:
```powershell
C:\path\to\gameapctl.exe daemon upgrade
```

O use la interfaz gráfica. Ejecute `gameapctl.exe` y, en la ventana del navegador que se abra,
haga clic en **"Upgrade"** en la sección GameAP Daemon.

![El botón Upgrade en la sección GameAP Daemon de la interfaz de gameapctl](/images/en/gameapctl/ui.png)
