---
title: GRPC API
layout: default
lang: ru
category: GameAP Daemon
order: 405
---

Начиная с версии GameAP 4.2 и GameAP Daemon 4.0 добавлена поддержка GRPC Bidirectional Streaming API для обмена данными 
между панелью и демоном. 
Новый способ обмена данными пришёл на смену старому способу с использованием BINN и REST API протоколам.

GRPC Bidi API обеспечивает более эффективную и надежную realtime коммуникацию между панелью и демоном.

## Базовая настройка

### GameAP

Для использования GRPC API необходимо указать следующие параметры в конфигурации GameAP (/etc/gameap/config.env):

```
GRPC_ENABLED=true
GRPC_PORT=31718
GRPC_TLS_ENABLED=true
```

### GameAP Daemon

