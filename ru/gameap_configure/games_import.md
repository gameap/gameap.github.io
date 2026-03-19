---
title: Импорт игр
layout: default
lang: ru
category: Настройка панели
order: 330
---

В GameAP 4.1 появилась возможность экспортировать и импортировать настройки игр, заготовки для создания игровых серверов.

Чтобы импортировать игры, перейдите на страницу **Администрирование** → **Игры** → **Импорт**.

![](/images/ru/gameap_configure/games_import/import_button.png)

На странице импорта доступна загрузите YAML файл с настройками игр, после этого нажмите на кнопку **Импортировать**.

![](/images/ru/gameap_configure/games_import/import_page.png)

GameAP поддерживает импорт игр из других панелей управления. На данный момент поддерживается импорт из следующих панелей:
- [Pterodactyl](https://pterodactyl.io/)
- [Pelican](https://pelican.io/)

На странице импорта игр доступна возможность переключения между форматами импорта из разных панелей.

## Особенности Pelican и Pterodactyl

Работа с импортированными Pelican Eggs и Pterodactyl Eggs возможна только на 
Docker и Podman [процесс менеджерах](/ru/daemon/process_manager.html).

Вам необходимо сконфигурировать GameAP Daemon для работы с одним из этих процесс менеджером.
Для этого при добавлении новой ноды, в разделе "Дополнительные настройки" выберите нужный процесс менеджер.

![](/images/ru/gameap_configure/games_import/daemon_process_manager.png)