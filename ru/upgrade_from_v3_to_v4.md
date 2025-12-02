---
title: Обновление с v3 до v4
layout: default
lang: ru
category: Установка GameAP
order: 191
---

## Автоматическое обновление с помощью gameapctl

Если у вас уже установлена старая версия панели, вы можете обновиться до новой версии с помощью [`gameapctl`](https://github.com/gameap/gameapctl).
Эта утилита устанавливается автоматически вместе с GameAP и позволяет управлять окружением панели, включая обновления.

Выполните следующие команды:

```shell
gameapctl self-update
gameapctl panel upgrade --to=v4
```

### Не установлен gameapctl? (gameapctl: command not found)

Если у вас не установлен `gameapctl`, вы можете скачать его с [GitHub releases](https://github.com/gameap/gameapctl/tags).

Или используйте curl для загрузки для Linux amd64:

```shell
curl -OL https://github.com/gameap/gameapctl/releases/download/v0.21.1/gameapctl-v0.21.1-linux-amd64.tar.gz
tar xvfz gameapctl-v0.21.1-linux-amd64.tar.gz -C /usr/local/bin
```

Затем выполните команду self-update:

```shell
gameapctl self-update
```

## Запуск параллельно со старой версией

Если вы хотите запустить GameAP v3 и GameAP v4 параллельно для тестирования перед обновлением, выполните следующие шаги:

Скачайте GameAP с [https://github.com/gameap/gameap/releases](https://github.com/gameap/gameap/releases) для вашей платформы.

Или используйте curl для загрузки для Linux amd64:

```shell
curl -OL https://github.com/gameap/gameap/releases/download/v4.0.0/gameap-v4.0.0-linux-amd64.tar.gz
```

Распакуйте в `/usr/bin/gameap`:

```shell
tar xvfz gameap-v4.0.0-linux-amd64.tar.gz -C /usr/bin
```

Запустите GameAP:

```shell
gameap --legacy-env /var/www/gameap/.env
```

GameAP v4 будет запущен на порту 8025.

### Настройка systemd-сервиса

Вы можете создать отдельный systemd-сервис для GameAP v4 для удобного управления.

Сначала создайте пользователя и группу `gameap`:

```shell
useradd -r -s /usr/sbin/nologin -d /var/lib/gameap gameap
mkdir -p /var/lib/gameap
chown gameap:gameap /var/lib/gameap
```

Затем создайте файл `/etc/systemd/system/gameap.service`:

```ini
[Unit]
Description=GameAP - Game Server Control Panel
Documentation=https://docs.gameap.com
After=network.target
Wants=network-online.target
Requires=network.target

[Service]
Type=simple
User=gameap
Group=gameap

# Рабочая директория
WorkingDirectory=/var/lib/gameap

ExecStart=/usr/bin/gameap

# Разрешить привязку к привилегированным портам
AmbientCapabilities=CAP_NET_BIND_SERVICE

# Корректное завершение
ExecStop=/bin/kill -TERM $MAINPID
KillMode=mixed
KillSignal=SIGTERM
TimeoutStopSec=30

# Политика перезапуска
Restart=always
RestartSec=5
StartLimitInterval=60
StartLimitBurst=3

# Конфигурация окружения
# EnvironmentFile=/etc/gameap/config.env

RuntimeDirectory=gameap
PIDFile=/run/gameap/gameap.pid

# Права доступа к файловой системе
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true

ReadWritePaths=/var/lib/gameap /var/lib/gameap/files
# Это необходимо, если вы запускаете GameAP с legacy-окружением
# ReadWritePaths=/var/www/gameap

# Логирование
StandardOutput=journal
StandardError=journal
SyslogIdentifier=gameap

[Install]
WantedBy=multi-user.target
```

Затем включите и запустите сервис:

```shell
systemctl daemon-reload
systemctl enable gameap
systemctl start gameap
```