---
title: Web API
layout: default
lang: ru
category: API
order: 10
---

* This will become a table of contents (this text will be scraped).
{:toc}


В GameAP существует API, позволяющее работать со многими функциями панели 
используя HTTP REST запросы.

Полную спецификацию смотрите на сайте [openapi.gameap.io](https://openapi.gameap.io).

## Генерация токена

Для авторизации в API необходимо сгенерировать токен. Для этого нужно зайти в
профиль пользователя, открыть вкладку "Токены", затем нажать "Сгенерировать токен",
после этого выбрать привилегии, которые будут выданы токену и затем нажать "Сохранить".

После сохранения токена он высветиться наверху страницы, его нужно скопировать, а затем 
использовать в HTTP запросах к API.

![Генерация токена в GameAP](/images/ru/api/generate_token_ru.gif)

## Примеры запросов

### Получение списка серверов

[Спецификация](https://openapi.gameap.io/#tag/servers/paths/~1api~1servers/get)

Вместо `YOUR_API_KEY` нужно подставить токен, полученный в предыдущем пункте.
Вместо `demo.gameap.ru` нужно подставить адрес вашей панели.

#### cURL

Пример запроса с помощью cURL:

```bash
curl --request GET \
--url 'https://demo.gameap.ru/api/servers' \
--header 'Content-Type: application/json: ' \
--header 'Authorization: Bearer YOUR_API_KEY'
```

#### PHP

Пример запроса на PHP с использование curl:

```php
<?php

$url = 'https://demo.gameap.ru/api/servers';
$apiKey = 'YOUR_API_KEY';

$curl = curl_init();

curl_setopt_array($curl, [
  CURLOPT_URL => $url,
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_ENCODING => '',
  CURLOPT_MAXREDIRS => 10,
  CURLOPT_TIMEOUT => 0,
  CURLOPT_FOLLOWLOCATION => true,
  CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
  CURLOPT_CUSTOMREQUEST => 'GET',
  CURLOPT_HTTPHEADER => [
    'Content-Type: application/json',
    'Authorization: Bearer ' . $apiKey,
  ],
]);

$response = curl_exec($curl);

if (curl_errno($curl)) {
    echo 'Curl error: ' . curl_error($curl);
    exit;
}

curl_close($curl);

$result = json_decode($response, true);
print_r($result);
```

Пример запроса на PHP с использованием [guzzlehttp/guzzle](https://github.com/guzzle/guzzle):

```php
<?php

require_once __DIR__ . '/vendor/autoload.php';

use GuzzleHttp\Client;

$url = 'https://demo.gameap.ru/api/servers';
$apiKey = 'YOUR_API_KEY';

$client = new Client();

$response = $client->request('GET', $url, [
  'headers' => [
    'Content-Type' => 'application/json',
    'Authorization' => 'Bearer ' . $apiKey,
  ]
]);

$result = json_decode($response->getBody(), true);
print_r($result);
```


### Запуск игрового сервера

[Спецификация](https://openapi.gameap.io/#tag/servers/paths/~1api~1servers~1%7Bserver%7D~1start/post)

Вместо `YOUR_API_KEY` нужно подставить токен, полученный в предыдущем пункте.
Вместо `demo.gameap.ru` нужно подставить адрес вашей панели.
Вместо `6` нужно подставить ID сервера.

#### cURL

```bash
curl -X POST "https://demo.gameap.ru/api/servers/6/start" \
  -H "Content-Type: application/json" \
  -H 'Authorization: Bearer YOUR_API_KEY'
```

#### PHP

Пример запроса на PHP с использование curl:

```php
<?php

$url = 'https://demo.gameap.ru/api/servers/6/start
$apiKey = 'YOUR_API_KEY';

$curl = curl_init();

curl_setopt_array($curl, [
  CURLOPT_URL => $url,
  CURLOPT_POST => true,
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
  CURLOPT_HTTPHEADER => [
    'Content-Type: application/json',
    'Authorization: Bearer ' . $apiKey,
  ],
]);

$response = curl_exec($curl);

if (curl_errno($curl)) {
    echo 'Curl error: ' . curl_error($curl);
    exit;
}

curl_close($curl);

$result = json_decode($response, true);
print_r($result);
```

Пример запроса на PHP с использованием [guzzlehttp/guzzle](https://github.com/guzzle/guzzle):

```php
<?php

use GuzzleHttp\Client;

$url = 'https://demo.gameap.ru/api/servers/6/start';
$apiKey = 'YOUR_API_KEY';

$client = new Client();

$response = $client->request('POST', $url, [
    'headers' => [
        'Content-Type' => 'application/json',
        'Authorization' => 'Bearer ' . $apiKey,
    ]
]);

$result = json_decode($response->getBody(), true);
print_r($result);

```