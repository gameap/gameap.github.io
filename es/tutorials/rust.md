---
title: Rust
layout: default
lang: es
category: Tutoriales
order: 202
---

Rust es un juego de supervivencia desarrollado por Facepunch Studios. 
Los jugadores deben recolectar recursos, fabricar objetos 
y construir refugios para sobrevivir en un entorno hostil de mundo abierto.

Rust es algo similar a [Minecraft](/es/tutorials/minecraft.html) 
por su sistema de supervivencia y fabricación de objetos.

## Preparación del entorno

GameAP ofrece soporte completo para servidores de juego de Rust, incluida la gestión de jugadores. 
La instalación del panel tomará poco tiempo; elige una de las siguientes 
opciones de instalación del panel de control:

* [Instalación de GameAP en Linux](/es/install/install_on_linux.html)
* [Instalación de GameAP en Windows](/es/install/install_on_windows.html)

### Instalación de GameAP Daemon

Para instalar el servidor de juego en una máquina (VDS), 
es necesario instalar GameAP Daemon. 
Durante la instalación de GameAP puedes elegir instalarlo por completo, 
incluido Daemon.

En el panel de control, ve a **"Administración"** → **"Servidores dedicados"** 
→ **"Crear"**. 
Aparecerá una ventana con una oferta de instalación automática. 
Copia el código y ejecútalo en el servidor dedicado.

Después de esto, puedes proceder a instalar el servidor de Rust.

## Instalación de un servidor de Rust en GameAP

Ve a **Administración** → **Servidores de juego** → **Crear**

![Formulario de creación de servidor de juego para Rust](/images/en/tutorials/rust/create_form.png)

* En el campo "Nombre", introduce cualquier nombre de servidor, por ejemplo, "My Rust Server".
* En el campo "Juego", selecciona Rust de la lista.
* En el campo de modificación, selecciona la modificación; por defecto, Vanilla.
* En el campo "Servidor dedicado", selecciona el nodo deseado donde se ubicará el servidor de juego.
* En el campo IP, elige la dirección deseada de tu servidor; luego puedes elegir un puerto libre para tu servidor o usar el sugerido.

## Configuración del servidor de juego

Para cambiar la configuración del servidor, ve a la sección **Servidor**, 
selecciona tu servidor y haz clic en **Administración**. Luego, abre la pestaña **Configuración**.

![Pestaña de configuración de un servidor de juego de Rust](/images/en/tutorials/rust/settings.png)

Es necesario reiniciar el servidor para que la configuración surta efecto.

### Nombre del servidor

Este es el nombre de tu servidor de Rust, que se mostrará 
a todos los jugadores en el juego en la ventana de búsqueda de servidores.

### Número máximo de jugadores en el servidor

El número máximo de jugadores que pueden jugar en el servidor simultáneamente. 
Por defecto, 32.

### Mapa en el servidor

El mapa en el servidor. 
Por defecto, el mapa generado proceduralmente Procedural Map. 
Valores posibles:

* Procedural Map
* Barren
* Craggy Island
* Hapis
* Savas Island

### Semilla de generación del mapa

La semilla para generar el Procedural Map.

### Tamaño del mundo

El tamaño del mundo.

### Intervalo de guardado del servidor

El intervalo de guardado del mundo. 
Este es el período durante el cual los datos del servidor se guardarán en el disco.
