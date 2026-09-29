# Локальная разработка

Рекомендуемый способ запуска — Docker Compose из корня репозитория. Go, Node.js и Task нужны только для запуска компонентов или проверок напрямую на хосте.

## Требования

Для запуска в Docker:

- Docker Engine с Docker Compose v2;
- свободные порты `80`, `443`, `5000`, `5432`, `8080`, `8083`, `8088` и `9092` (`8082` также нужен для Bot);
- `curl` для команды проверки health (необязательно);
- доступ в интернет для загрузки образов, карты OSRM и данных ingestion;
- около 1 ГБ свободной оперативной памяти для первой подготовки OSRM.

Дополнительно для разработки на хосте:

- Go 1.27;
- Node.js 22;
- [Task](https://taskfile.dev/docs/installation) — необязательно;
- [cloudflared](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/) — только для публичного туннеля MAX Bot;
- Bash и `openssl` — для генератора тестовых MAX Mini App initData.

## 1. Подготовка окружения

Из корня репозитория выполните:

```bash
cp .env.example .env
cp api/.env.example api/.env
cp ingestion/.env.example ingestion/.env
cp bot/.env.example bot/.env
```

Локальные `.env` не коммитятся. Шаблоны уже содержат согласованные адреса контейнеров, имена баз данных, пароли и сервисные токены для локальной разработки.

Перед запуском учитывайте внешние интеграции:

- `DADATA_API_KEY` в `api/.env` нужен только для прямого и обратного геокодирования;
- значение `MAX_BOT_TOKEN` из `api/.env` можно использовать с локальным генератором initData; реальный токен нужен для работы с MAX;
- `bot` требует реальный `MAX_BOT_TOKEN`, публичный `WEBHOOK_URL` и `WEBHOOK_SECRET`;
- `BOT_SERVICE_TOKEN` должен совпадать в `api/.env` и `bot/.env`;
- `INGESTION_SERVICE_TOKEN` должен совпадать в `api/.env` и `ingestion/.env`.

## 2. Запуск основного стека

Для API, Mini App, загрузки данных и расчёта рейтинга Bot не требуется:

```bash
docker compose up --build -d postgres kafka osrm api ingestion miniapp caddy kafka-ui
```

Compose автоматически создаст сеть и volumes, дождётся готовности PostgreSQL и Kafka и применит миграции при старте API и ingestion.

Первый запуск OSRM дольше обычного: образ скачивает карту Санкт-Петербурга и строит пешеходный граф. Следить за подготовкой можно командой:

```bash
docker compose logs -f osrm
```

Когда в журнале появится запуск `osrm-routed`, проверьте контейнеры и API:

```bash
docker compose ps
curl --fail http://localhost:8080/health
```

Ожидаемый ответ API:

```json
{"status":"ok"}
```

Доступные интерфейсы:

| Сервис | Адрес |
| --- | --- |
| Mini App через Caddy | <http://localhost/> |
| Mini App напрямую | <http://localhost:8083/> |
| Swagger UI | <http://localhost:8080/docs/> |
| Kafka UI | <http://localhost:8088/> |
| API health | <http://localhost:8080/health> |
| OSRM | <http://localhost:5000/> |

Если контейнер не запустился или стал `unhealthy`, сначала посмотрите его журнал:

```bash
docker compose logs --tail=200 api
docker compose logs --tail=200 ingestion
docker compose logs --tail=200 osrm
```

## 3. Запуск MAX Bot

Bot является отдельным шагом, потому что MAX должен иметь доступ к публичному HTTPS webhook.

1. Запишите реальные `MAX_BOT_TOKEN` и `WEBHOOK_SECRET` в `bot/.env`.
2. Убедитесь, что `MAX_BOT_TOKEN` и `BOT_SERVICE_TOKEN` совпадают со значениями в `api/.env`.
3. Откройте туннель к локальному порту Bot:

```bash
task bot:tunel
# без Task: cloudflared tunnel --url http://localhost:8082
```

4. Добавьте к выданному HTTPS-адресу путь `/webhook` и запишите результат в `WEBHOOK_URL` файла `bot/.env`.
5. Запустите Bot:

```bash
docker compose up --build -d bot
docker compose logs -f bot
```

После настройки всех внешних значений весь стек можно поднимать одной командой:

```bash
docker compose up --build -d
# или
task all:up
```

## Управление сервисами

Просмотр состояния и журналов:

```bash
docker compose ps
docker compose logs -f
# или
task logs
```

Остановка без удаления volumes:

```bash
docker compose down
# или
task all:down
```

Запуск отдельных компонентов:

```bash
task postgres:up
task kafka:up
task osrm:up
task api:up
task ingestion:up
task bot:up
```

## Запуск компонентов на хосте

При запуске приложения вне Docker замените контейнерные адреса в его `.env`: PostgreSQL — на `localhost:5432`, Kafka — на `localhost:9092`, API — на `http://localhost:8080`, OSRM — на `http://localhost:5000`.

Подробности находятся в документации разработки [API](../api/docs/development.md), [Bot](../bot/docs/development.md) и [OSRM](../osrm/README.md).

## Проверки

Go-модули проверяются независимо:

```bash
(cd api && go test ./...)
(cd ingestion && go test ./...)
(cd bot && go test ./...)
```

Mini App:

```bash
(cd miniapp && npm ci && npm run build)
```

Каждый Go-модуль содержит собственный `go.mod`; команды Go запускаются из каталога соответствующего модуля.

## Сквозной сценарий уведомлений

Ingestion обнаруживает новых прямых конкурентов и ближайшие события, рассчитывает расстояние через API и публикует задания в Kafka. Bot читает topic `notifications` и доставляет сообщения в зарегистрированный личный чат MAX. Порядок пользовательской проверки и ожидаемые результаты приведены в [корневом README](../README.md#как-тестировать).
