# Locky password keeper

## Overview

- Supports following commands:
  - `/get <service>`
  - `/set <service> <login> <password>`
  - `/del <service>`
- Can support webhooks
- Worker pool for requests
- PostgreSQL as underlying storage

## Deploy

```shell
# checkout sources
git clone https://github.com/nikmy/locky.git

# go to source dir
cd locky

# build and run the bot
docker-compose up --build db app
```

## Environment

In `docker-compose.yml`, service `app` you can find environment
variables that control bot behaviour. The most important are:
- `TOKEN`: telegram API token from [@BotFather](https://t.me/BotFather)
- `WEBHOOK`: if set, is used for webhook instead of polling