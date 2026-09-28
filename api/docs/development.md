# Разработка и запуск API

## Требования

- Go 1.27;
- PostgreSQL с PostGIS;
- OSRM с пешеходным профилем для полного сценария рейтинга;
- Docker Compose для рекомендуемого локального запуска;
- Task — необязательная оболочка над командами Compose.

## Подготовка окружения

```bash
cp api/.env.example api/.env
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

### Внешние интеграции

| Переменная | По умолчанию | Назначение |
| --- | --- | --- |
| `OSRM_BASE_URL` | `http://osrm:5000` | Адрес OSRM |
| `OSRM_TIMEOUT` | `5s` | Таймаут OSRM |
| `DADATA_BASE_URL` | `https://suggestions.dadata.ru/suggestions/api/4_1/rs` | Базовый URL DaData Suggest и GeoLocate |
| `DADATA_API_KEY` | — | Ключ для заголовка `Authorization: Token <key>` внешних запросов DaData |
| `GEOCODER_TIMEOUT` | `5s` | Общий таймаут HTTP-клиента DaData |

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

### Интеграция DaData

Публичные и внутренние маршруты GeoLogic остаются GET-методами, описанными в OpenAPI. API преобразует их во внешние POST-запросы DaData:

- `GET /api/v1/geocoding/suggestions` и внутренний аналог вызывают `${DADATA_BASE_URL}/suggest/address` с `count: 5`, `locations: [{"city":"Санкт-Петербург"}]` и `restrict_value: true`; параметр GeoLogic `limit` обрезает полученный список до 1–5 элементов;
- `GET /api/v1/geocoding/address` и внутренний аналог вызывают `${DADATA_BASE_URL}/geolocate/address` с `lat`, `lon` и `count: 1`;
- внешние запросы используют JSON и заголовок `Authorization: Token <DADATA_API_KEY>`;
- GeoLogic возвращает только нормализованные поля `address`, `lat` и `lon`; ошибка сети, неуспешный статус DaData или пустой результат обратного геокодирования преобразуются в `503 provider_unavailable`.

Подсказки принудительно ограничены Санкт-Петербургом. Обратное геокодирование принимает полный диапазон координат WGS 84 и отдельную проверку города не выполняет.

## Диагностика

- Ошибка до запуска HTTP-сервера: проверьте обязательные env-переменные, `rating.yml`, доступность БД и миграции.
- Создание точки возвращает ошибку провайдера: проверьте `OSRM_BASE_URL`, готовность графа и таймаут.
- Геокодирование не работает: проверьте `DADATA_API_KEY`, базовый URL и сетевой доступ.
- Mini App получает `401`: проверьте сырой `X-Max-Init-Data`, соответствие `MAX_BOT_TOKEN` и допустимый возраст подписи.
- Браузер блокирует запрос: `CORS_ALLOWED_ORIGIN` должен точно совпадать с origin Mini App.

Секреты и сырые ответы провайдеров не следует добавлять в журналы или коммитить в репозиторий.
