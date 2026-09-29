# GeoLogic Monitor

GeoLogic Monitor помогает владельцам локального офлайн-бизнеса Санкт-Петербурга оценивать окружение торговых точек и своевременно узнавать о новых конкурентах и городских событиях рядом. Пользователь работает через MAX: бот регистрирует личный чат и открывает Mini App, а Mini App позволяет управлять точками, смотреть Smart Score и историю его изменения.

![Контейнерная схема GeoLogic Monitor](docs/diagrams/rendered/c4-containers.svg)

## Возможности

- добавление, просмотр и удаление торговых точек;
- поиск адресов и обратное геокодирование через DaData;
- расчёт Smart Score по инфраструктуре и пешеходным расстояниям;
- хранение истории успешных расчётов и плановый пересчёт рейтинга;
- загрузка городской инфраструктуры и событий из внешних источников;
- обнаружение новых прямых конкурентов и событий рядом с точкой;
- доставка уведомлений через Kafka и MAX.

## Основной пользовательский сценарий

1. Предприниматель запускает бота в MAX командой `/start` и открывает Mini App.
2. Выбирает тип бизнеса, вводит название точки и находит адрес в Санкт-Петербурге.
3. API определяет инфраструктуру рядом, уточняет пешеходные расстояния через OSRM, рассчитывает Smart Score и атомарно сохраняет точку с первым значением рейтинга.
4. Mini App показывает текущую оценку, время расчёта и историю изменения рейтинга.
5. Ingestion обновляет данные об окружении и событиях, обнаруживает значимые изменения и публикует задания в Kafka.
6. Bot доставляет владельцу точки уведомление о новом прямом конкуренте или ближайшем событии в MAX.

Подробное описание продукта: [продукт и пользовательский сценарий](docs/product.md).

## Архитектура

| Компонент | Ответственность | Документация |
| --- | --- | --- |
| `api` | Постоянные данные, HTTP-контракт, геокодирование, Smart Score и планировщик | [README API](api/README.md) |
| `ingestion` | Загрузка инфраструктуры и событий, мониторинг изменений, публикация уведомлений | [README Ingestion](ingestion/README.md) |
| `bot` | Webhook MAX, `/start`, Kafka consumer и доставка уведомлений | [README Bot](bot/README.md) |
| `miniapp` | Пользовательский интерфейс точек, рейтинга и истории | [README Mini App](miniapp/README.md) |
| `osrm` | Локальная пешеходная маршрутизация | [README OSRM](osrm/README.md) |

`api`, `ingestion` и `bot` являются независимыми Go-модулями, а `miniapp` — приложением React/TypeScript. Постоянными продуктовыми данными владеет API; остальные сервисы используют его HTTP-контракт. Kafka связывает ingestion и bot, а отдельная база ingestion хранит только временный кэш мониторинга.

Подробная схема: [архитектура системы](docs/architecture.md).

## Технологии и зависимости

Go 1.27, `net/http`, PostgreSQL 18, PostGIS 3.6, React 19, TypeScript, Vite, Kafka, franz-go, OSRM, Nginx, Caddy, Docker Compose и Taskfile. Версии Go-зависимостей зафиксированы в `go.mod` и `go.sum`, frontend-зависимостей — в `package-lock.json`, базовых образов — в Dockerfile и Compose.

## Быстрый запуск

Нужны Docker Engine и Docker Compose v2. Все команды выполняются из корня репозитория.

1. Создайте локальные env-файлы из безопасных примеров:

```bash
cp .env.example .env
cp api/.env.example api/.env
cp ingestion/.env.example ingestion/.env
cp bot/.env.example bot/.env
```

2. Заполните реальные значения `DADATA_API_KEY`, `TWOGIS_API_KEY`, `MAX_BOT_TOKEN`, `WEBHOOK_URL` и `WEBHOOK_SECRET`. `BOT_SERVICE_TOKEN` должен совпадать в `bot` и API; `INGESTION_SERVICE_TOKEN` — в `ingestion` и API. Между собой служебные токены различаются.

3. Запустите все локальные компоненты одной командой:

```bash
docker compose up --build -d
```

Первая сборка OSRM скачивает выгрузку OpenStreetMap Санкт-Петербурга и строит пешеходный граф. Повторные сборки используют Docker cache.

4. Проверьте состояние:

```bash
docker compose ps
curl --fail http://localhost:8080/health
```

Ожидаемый ответ API:

```json
{"status":"ok"}
```

### Порты

| Компонент | Адрес с хоста | Назначение |
| --- | --- | --- |
| Caddy / Mini App | `http://localhost/` | Единая точка входа Mini App и API |
| API | `http://localhost:8080` | HTTP API и Swagger UI `/docs/` |
| Bot | `http://localhost:8082/webhook` | Локальный webhook, публикуемый через HTTPS-туннель |
| Mini App | `http://localhost:8083` | Прямой доступ к Nginx-контейнеру |
| Kafka UI | `http://localhost:8088` | Диагностика Kafka |
| OSRM | `http://localhost:5000` | Route и Table API |
| PostgreSQL | `localhost:5432` | API DB и временный кэш ingestion |
| Kafka | `localhost:9092` | Локальный broker |

Публикуемые порты PostgreSQL, Kafka и OSRM можно изменить в корневом `.env`.

### Переменные окружения

| Файл | Основные переменные |
| --- | --- |
| `.env` | `POSTGRES_ADMIN_USER`, `POSTGRES_ADMIN_PASSWORD`, `API_POSTGRES_PASSWORD`, `INGESTION_POSTGRES_PASSWORD`, `POSTGRES_PORT`, `KAFKA_HOST_PORT`, `KAFKA_RETENTION_MS`, `OSRM_HTTP_PORT` |
| `api/.env` | параметры PostgreSQL, `BOT_SERVICE_TOKEN`, `INGESTION_SERVICE_TOKEN`, `MAX_BOT_TOKEN`, `MINIAPP_INIT_DATA_MAX_AGE`, `DADATA_BASE_URL`, `DADATA_API_KEY`, `GEOCODER_TIMEOUT`, `OSRM_BASE_URL`, `RATING_RECALC_CRON`, параметры HTTP и логирования |
| `ingestion/.env` | параметры PostgreSQL и Kafka, `GEO_API_URL`, `INGESTION_SERVICE_TOKEN`, `TWOGIS_API_KEY`, расписания загрузки и мониторинга, параметры логирования |
| `bot/.env` | `MAX_BOT_TOKEN`, `WEBHOOK_URL`, `WEBHOOK_SECRET`, `BOT_SERVICE_TOKEN`, адреса API и Kafka, параметры HTTP и логирования |

Полный набор с безопасными значениями находится в соответствующих `.env.example`. Реальные токены, ключи, пароли и MAX `initData` не коммитятся.

## Внешние сервисы и данные

- DaData Suggest и GeoLocate — поиск адресов и обратное геокодирование;
- городские API Санкт-Петербурга `spb-classif`, `egs` и `yazzh` — инфраструктура и события;
- OpenStreetMap Overpass — коммерческая и сервисная инфраструктура;
- 2ГИС Catalog API — обнаружение новых прямых конкурентов;
- подготовленный справочник метро — станции и входы метрополитена;
- OpenStreetMap PBF — пешеходный граф OSRM;
- MAX Bot API и Web App SDK — пользовательский интерфейс и доставка сообщений.

Источники, назначение данных, порядок обновления и поведение при ошибках подробно описаны в [документации ingestion](ingestion/README.md). Рабочий MVP использует реальные внешние интеграции. Значения в OpenAPI, DATA-API и разделе проверки ниже являются только воспроизводимыми тестовыми примерами.

## Как тестировать

### Проверка основного сценария в MAX

1. Откройте `https://max.ru/t499_hakaton_max_bot`, нажмите «Начать» или отправьте `/start`.
   Ожидаемый результат: бот подтверждает регистрацию и показывает кнопку открытия GeoLogic Monitor.
2. Откройте Mini App.
   Ожидаемый результат: отображаются список точек и действие «Добавить точку»; ошибок авторизации нет.
3. Создайте точку с произвольным названием, любым типом бизнеса из списка и адресом `Санкт-Петербург, Невский проспект, 28`. Можно попробовать и любой другой адрес в Санкт-Петербурге. Выберите адрес из подсказок и подтвердите данные.
   Ожидаемый результат: точка появляется в списке, API возвращает первый Smart Score и время расчёта.
4. Откройте созданную точку.
   Ожидаемый результат: видны адрес, текущий Smart Score и график истории; после повторных расчётов новые значения добавляются в хронологическом порядке.
5. Повторно откройте Mini App.
   Ожидаемый результат: созданная точка и её история сохранены; основной сценарий можно пройти повторно без сброса приложения.

### Проверка уведомления о событии

Для проверки нужен запущенный локальный стек и зарегистрированный в MAX чат. Сначала пройдите сценарий выше и создайте точку по адресу `Санкт-Петербург, Невский проспект, 28`: тестовое событие будет добавлено рядом с ней.

1. Из корня репозитория загрузите сервисный токен из локального env-файла, сформируйте уникальный идентификатор и дату события через 24 часа:

```bash
INGESTION_SERVICE_TOKEN=$(sed -n 's/^INGESTION_SERVICE_TOKEN=//p' ingestion/.env)
EVENT_EXTERNAL_ID="readme-event-$(date +%s)"
EVENT_DATE=$(date -u -d '+24 hours' '+%Y-%m-%dT%H:%M:%SZ')
```

2. Добавьте событие через внутренний API:

```bash
curl --fail-with-body --request PUT \
  --url http://localhost:8080/internal/v1/events \
  --header "Authorization: Bearer ${INGESTION_SERVICE_TOKEN}" \
  --header 'Content-Type: application/json' \
  --data @- <<JSON
{
  "provider": "readme",
  "external_id": "${EVENT_EXTERNAL_ID}",
  "lat": 59.9358,
  "lon": 30.3259,
  "date": "${EVENT_DATE}",
  "info": "Тестовое событие GeoLogic Monitor"
}
JSON
```

Ожидаемый результат: API возвращает созданное событие с `notified_at: null`.

3. Перезапустите ingestion, чтобы job уведомлений выполнился сразу, и откройте журналы сервисов:

```bash
docker compose restart ingestion
docker compose logs --follow ingestion bot
```

Ожидаемый результат: ingestion выводит запись `event notification published`, Bot — `processing notification` с типом `event_upcoming`, а в зарегистрированный личный чат MAX приходит уведомление о тестовом событии. Остановить просмотр журналов можно сочетанием `Ctrl+C`; контейнеры продолжат работу.

Для повторной проверки выполните команды заново: новый `EVENT_EXTERNAL_ID` обязателен, поскольку после успешной публикации событие получает `notified_at` и повторно не обрабатывается.

Отдельный воспроизводимый сценарий уведомления о конкуренте не приводится. По политике использования данных 2ГИС записи конкурентов не сохраняются в основной продуктовой БД; ingestion использует только краткоживущий технический кэш для дедупликации результатов мониторинга.

### Проверка API

OpenAPI доступен в [`api/docs/openapi.yaml`](api/docs/openapi.yaml), а обязательный воспроизводимый сценарий API находится в [`DATA-API.yaml`](DATA-API.yaml). Для ручной проверки пользовательских методов сформируйте актуальную MAX initData:

```bash
INIT_DATA=$(./scripts/generate-max-init-data.sh)
```

Скрипт требует Bash и `openssl`, скрыто запрашивает токен MAX-бота и ID тестового пользователя. InitData ограничена настройкой `MINIAPP_INIT_DATA_MAX_AGE`, поэтому её формируют непосредственно перед проверкой.

### Очистка тестовых данных

- удалите созданную тестовую точку в Mini App — вместе с точкой API удалит её историю рейтинга;
- остановите локальный стек без удаления данных командой `docker compose down`;
- повторно запустите его командой `docker compose up -d`;
- для полного сброса только локального окружения используйте `docker compose down -v`, учитывая, что команда удаляет локальные PostgreSQL и Kafka volumes.

## Границы решения

- продукт и источники данных настроены для Санкт-Петербурга;
- Smart Score оценивает окружение точки и не является прогнозом выручки;
- для MAX необходимы реальные токен бота, публичный HTTPS webhook и корректно подключённый Mini App;
- внешние источники, DaData и 2ГИС требуют сетевого доступа и действующих ключей;
- Kafka обеспечивает доставку как минимум один раз, поэтому при сбое между отправкой сообщения и подтверждением offset возможен повтор уведомления.

## Документация

- [Архитектура системы](docs/architecture.md)
- [Продукт и пользовательский сценарий](docs/product.md)
- [Владение данными](docs/data-model.md)
- [Smart Score](docs/impact-engine.md)
- [Документация API](api/README.md)
- [Документация Ingestion](ingestion/README.md)
- [Документация Bot](bot/README.md)
- [Документация Mini App](miniapp/README.md)
- [Документация OSRM](osrm/README.md)
- [HTTP OpenAPI](api/docs/openapi.yaml)
- [Kafka AsyncAPI](docs/asyncapi/notifications.yaml)
- [Локальная разработка](docs/local-development.md)
- [Исходники и изображения диаграмм](docs/diagrams/README.md)
- [Правила работы с репозиторием](CONTRIBUTING.md)
