---
title: Panel errors
layout: default
lang: en
category: Main
order: 40
---

* This will become a table of contents (this text will be scraped).
{:toc}

Description some possible errors and their fixing.

## Errors after installation

### Whoops, looks like something went wrong.

#### Generation the encryption key

![](/images/errors/key_generate.png)

This problem will occur if you have not generated a key. Go to the panel directory and run the command:

```bash
php artisan key:generate --force
```

#### Database migration

![](/images/errors/db_migrate.png)

This problem will occur if you have not migrate database. Go to the panel directory and run the command:

```
php artisan migrate --seed
```