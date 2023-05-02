---
title: Web API
layout: default
lang: en
category: API
order: 10
---

* This will become a table of contents (this text will be scraped).
  {:toc}

In GameAP, there is an API that allows working with many panel functions 
using HTTP REST requests. 

See the full specification on the website [openapi.gameap.io](https://openapi.gameap.io).

## Token Generation

To authenticate with the API, you need to generate a token.
To do this, go to your user profile, open the "Tokens" tab, then click "Generate New Token".
Next, select the privileges that will be granted to the token, and then click "Save".

After the token is saved, it will be displayed at the top of the page.
Copy it and use it in your HTTP requests to the API.

![Token generation in GameAP](/images/en/api/generate_token_en.gif)


## Request Examples

### Retrieving a list of servers

[Specification](https://openapi.gameap.io/#tag/servers/paths/~1api~1servers/get)

Replace `YOUR_API_KEY` with the token obtained in the previous step.
Replace `demo.gameap.ru` with the address of your panel.

#### cURL

Example request using cURL:

```bash
curl --request GET \
--url 'https://demo.gameap.ru/api/servers' \
--header 'Content-Type: application/json: ' \
--header 'Authorization: Bearer YOUR_API_KEY'
```

#### PHP

Example PHP request using curl:

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

Example PHP request using [guzzlehttp/guzzle](https://github.com/guzzle/guzzle) library:

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


### Starting a game server

[Specification](https://openapi.gameap.io/#tag/servers/paths/~1api~1servers~1%7Bserver%7D~1start/post)

Replace `YOUR_API_KEY` with the token obtained in the previous step.
Replace `demo.gameap.ru` with the address of your panel.
Replace `6` with the ID of your server.


#### cURL

Example request using cURL:

```bash
curl -X POST "https://demo.gameap.ru/api/servers/6/start" \
  -H "Content-Type: application/json" \
  -H 'Authorization: Bearer YOUR_API_KEY'
```

#### PHP

Example PHP request using curl:

```php
<?php

$url = 'https://demo.gameap.ru/api/servers/6/start';
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

Example PHP request using [guzzlehttp/guzzle](https://github.com/guzzle/guzzle) library:

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