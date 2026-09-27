# Разработка и запуск API

## Требования

- Go 1.27;
- PostgreSQL с PostGIS;
- OSRM с пешеходным профилем для полного сценария рейтинга;
- Docker Compose для рекомендуемого локального запуска;
- Task — необязательная оболочка над командами Compose.

## Подготовка окружения

```bash
cp api/.env.template api/.env
```

Значения `change-me` подходят только как локальные заглушки. Файл `api/.env` не коммитится.

### HTTP и документация

| Переменная | По умолчанию | Назначение |
| --- | --- | --- |
| `API_PORT` | `8080` | Порт HTTP-сервера |
| `READ_TIMEOUT` | `5s` | Таймаут чтения запроса |
| `WRITE_TIMEOUT` | `10s` | Таймаут записи ответа |
| `IDLE_TIMEOUT` | `60s` | Таймаут keep-alive соединения |
| `DOCS_DIR` | `./docs` | Каталог Swagger UI и OpenAPI относительно рабочего каталога API |

### PostgreSQL

| Переменная | По умолчанию в коде | Назначение |
| --- | --- | --- |
| `POSTGRES_HOST` | `localhost` | Хост PostgreSQL; в Compose используется `postgres` |
| `POSTGRES_PORT` | `5432` | Порт PostgreSQL |
| `POSTGRES_USER` | `postgres` | Пользователь БД |
| `POSTGRES_PASSWORD` | `postgres` | Пароль БД |
| `POSTGRES_NAME` | `postgres` | Имя БД |
| `POSTGRES_SSL_MODE` | `disable` | Режим SSL pgx |
| `MIN_IDLE_CONNS` | `0` | Минимум простаивающих соединений |
| `MAX_OPEN_CONNS` | `100` | Максимум открытых соединений |
| `MAX_CONN_LIFETIME` | `5m` | Время жизни соединения |

Шаблон также содержит `POSTGRES_DB`; текущая конфигурация API читает имя базы из `POSTGRES_NAME`.

### Внешние интеграции

| Переменная | По умолчанию | Назначение |
| --- | --- | --- |
| `OSRM_BASE_URL` | `http://osrm:5000` | Адрес OSRM |
| `OSRM_TIMEOUT` | `5s` | Таймаут OSRM |
| `DADATA_BASE_URL` | URL DaData Suggest | Базовый URL геокодера |
| `DADATA_API_KEY` | — | Ключ DaData |
| `GEOCODER_TIMEOUT` | `5s` | Таймаут геокодера |

### Безопасность и расписание

| Переменная | По умолчанию | Назначение |
| --- | --- | --- |
| `MAX_BOT_TOKEN` | обязательна | Проверка подписи Mini App initData |
| `MINIAPP_INIT_DATA_MAX_AGE` | `1h` | Максимальный возраст initData |
| `BOT_SERVICE_TOKEN` | обязательна | Bearer-токен сервиса `bot` |
| `INGESTION_SERVICE_TOKEN` | обязательна | Bearer-токен сервиса `ingestion` |
| `CORS_ALLOWED_ORIGIN` | обязательна | Разрешённый origin Mini App |
| `RATING_RECALC_CRON` | `0 3 1 * *` | Пятичастное расписание в `Europe/Moscow` |
| `RATING_CONFIG_PATH` | `rating.yml` | Путь к конфигурации формулы |

Логирование настраивается через `LOG_LEVEL`, `LOG_FORMAT` (`json` или `text`) и `LOG_ADD_SOURCE`.

## Запуск с Docker Compose

```bash
docker compose up -d postgres osrm
docker compose up --build -d api
docker compose logs -f api
```

Или через Task:

```bash
task api:up
task logs
task api:down
```

Контейнер запускается из `/app`, поэтому стандартные пути `./migrations`, `./docs` и `rating.yml` уже включены в образ.

## Запуск напрямую

Укажите в `api/.env` адреса зависимостей, доступные с хоста, затем:

```bash
cd api
go run ./cmd/api
```

При каждом старте goose применяет ещё не выполненные миграции из `api/migrations`. Отдельная ручная команда миграции для штатного запуска не требуется.

## Проверка

```bash
curl http://localhost:8080/health
cd api
go test ./...
```

Swagger UI: <http://localhost:8080/docs/>. `/health` проверяет только HTTP-процесс и намеренно не опрашивает PostgreSQL, DaData или OSRM.

## Диагностика

- Ошибка до запуска HTTP-сервера: проверьте обязательные env-переменные, `rating.yml`, доступность БД и миграции.
- Создание точки возвращает ошибку провайдера: проверьте `OSRM_BASE_URL`, готовность графа и таймаут.
- Геокодирование не работает: проверьте `DADATA_API_KEY`, базовый URL и сетевой доступ.
- Mini App получает `401`: проверьте сырой `X-Max-Init-Data`, соответствие `MAX_BOT_TOKEN` и допустимый возраст подписи.
- Браузер блокирует запрос: `CORS_ALLOWED_ORIGIN` должен точно совпадать с origin Mini App.

Секреты и сырые ответы провайдеров не следует добавлять в журналы или коммитить в репозиторий.
