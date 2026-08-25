---
title: Instalación en Windows
layout: default
lang: es
category: Instalación de GameAP
order: 101
---

## Versiones compatibles

| Versión     | Compatible |
|-------------|------------|
| Server 2025 | ✔          |
| Server 2022 | ✔          |
| Server 2019 | ✔          |
| 11          | ✔          |
| 10          | ✔          |


## Descarga de GameAP Control

Necesita descargar la utilidad GameAP Control (gameapctl)
para gestionar el entorno de GameAP.

Para ello, vaya a la página de versiones de gameapctl en Github:
[https://github.com/gameap/gameapctl/releases](https://github.com/gameap/gameapctl/releases)

Seleccione la última versión y haga clic en ella.

<video width="1280" height="720" controls>
  <source src="/images/en/gameapctl/download.webm" type="video/webm">
  Su navegador no soporta la etiqueta de vídeo.
</video>

Después, busque la versión adecuada para usted.
La arquitectura más común es Windows AMD64,
por lo que lo más probable es que necesite descargar este archivo:

![Elección del archivo gameapctl para Windows AMD64 en la página de versiones](/images/en/gameapctl/download_release_windows_amd64.png)

## Instalación del panel mediante GameAP Control UI

Después de descargar el archivo de gameapctl, ejecútelo.

Se abrirá una ventana del navegador, en la que debe hacer clic en "Install"
en la sección Web/API.

<video width="1280" height="720" controls>
  <source src="/images/en/gameapctl/windows-install.webm" type="video/webm">
  Su navegador no soporta la etiqueta de vídeo.
</video>

### Parámetros de instalación

Especifique los datos necesarios para la instalación.

![Formulario de parámetros de instalación del panel en la interfaz de gameapctl](/images/en/gameapctl/ui_gameap_installation.png)

#### Host

Especifique el dominio o la dirección IP donde estará accesible el panel.

En el caso de una dirección IP, debe ser la dirección asignada a
la interfaz de red del VDS. Si su red utiliza NAT,
no especifique la IP externa, sino la interna,
y luego configure el reenvío de puertos.

Se puede especificar cualquier dominio, pero no olvide configurar el DNS.

Ejemplos de valores correctos:
* 10.182.104.8
* 10.182.104.8:2080
* example.com
* http://example.com

#### Base de datos

La base de datos donde se almacenarán los datos: usuarios, información de los servidores, etc.
Puede utilizar:
* [PostgreSQL](https://www.postgresql.org/). Recomendado para proyectos grandes con muchos servidores de juego y usuarios.
* [MySQL](https://www.mysql.com/)/[MariaDB](https://mariadb.org/)
* [SQLite](https://www.sqlite.org/). Si se espera que la carga de su servidor sea baja y no planea utilizar más de 10 servidores de juego.

#### Instalación de GameAP Daemon

Además del panel web, también puede instalar
la parte de servidor GameAP Daemon, que gestiona las operaciones de los servidores de juego.

### Finalización de la instalación

Espere a que finalice la instalación.
Algunas etapas pueden tardar un tiempo considerable.

No olvide guardar los datos de acceso y la información de la base de datos
que se mostrarán al final.

![Datos de acceso al panel que se muestran al finalizar la instalación](/images/en/gameapctl/gameap_finished_installation.png)
