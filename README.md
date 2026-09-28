# GeoLogic Monitor

GeoLogic Monitor помогает владельцам кофеен, кафе и небольших ресторанов Санкт-Петербурга оценивать окружение торговых точек и узнавать о важных изменениях рядом. Пользователь работает через MAX Mini App: добавляет точки, смотрит текущий рейтинг и его историю. Бот регистрирует чат и доставляет уведомления.

![Контейнерная схема GeoLogic Monitor](docs/diagrams/rendered/c4-containers.svg)

## Что умеет система

- добавлять, показывать и удалять точки пользователя;
- находить адреса через DaData Suggest;
- рассчитывать рейтинг окружения по инфраструктуре и пешеходным расстояниям;
- хранить историю успешных расчётов и пересчитывать рейтинг по расписанию;
- загружать инфраструктуру и городские события;
- принимать webhook MAX и отправлять сообщения из Kafka.

Сквозной сценарий уведомлений пока не завершён: Kafka consumer в `bot` реализован, producer уведомлений в `ingestion` отсутствует. Остальные возможности выше представлены в текущем коде.

## Пользовательский сценарий

1. Предприниматель запускает бота в MAX и открывает Mini App.
2. Выбирает тип бизнеса и адрес новой точки.
3. API находит объекты рядом, уточняет расстояния через OSRM и сохраняет первый рейтинг.
4. Mini App показывает точку, текущую оценку и график истории.
5. API обновляет рейтинги по расписанию, а ingestion собирает изменения окружения для уведомлений.

Подробное описание: [продукт и пользовательский сценарий](docs/product.md).

## Сервисы

| Компонент | Ответственность | Документация |
| --- | --- | --- |
| `api` | Постоянные данные, HTTP-контракт, геокодинг, рейтинг и планировщик | [README API](api/README.md) |
| `ingestion` | Загрузка инфраструктуры и событий, мониторинг изменений | Документация будет добавлена отдельно |
| `bot` | Webhook MAX, `/start`, Kafka consumer и отправка уведомлений | [README Bot](bot/README.md) |
| `miniapp` | Пользовательский интерфейс точек, рейтинга и истории | [README Mini App](miniapp/README.md) |
| `osrm` | Локальная пешеходная маршрутизация | [README OSRM](osrm/README.md) |

Сервисы независимы: `api`, `ingestion` и `bot` являются отдельными Go-модулями, а `miniapp` — приложением React/TypeScript. Постоянными данными владеет API; остальные сервисы обращаются к нему по HTTP.

## Технологии

Go 1.27, `net/http`, PostgreSQL 18 и PostGIS 3.6, React 19, TypeScript, Vite, Kafka, franz-go, OSRM, Docker Compose и Taskfile.

## Быстрый запуск

Нужны Docker и Docker Compose; `curl` используется только для проверки API. Все команды выполняются из корня репозитория.

1. Создайте локальные env-файлы:

```bash
cp .env.template .env
cp api/.env.template api/.env
cp ingestion/.env.template ingestion/.env
cp bot/.env.template bot/.env
```

Шаблоны уже согласованы для локальной сети Compose. Для запуска API и Mini App их можно использовать без дополнительных изменений. Реальные `DADATA_API_KEY` и `MAX_BOT_TOKEN` нужны только для соответствующих внешних интеграций.

2. Запустите основной стек без MAX Bot:

```bash
docker compose up --build -d postgres kafka osrm api ingestion miniapp caddy kafka-ui
```

Первая сборка OSRM скачивает карту Санкт-Петербурга и подготавливает граф, поэтому занимает больше времени последующих запусков.

3. Проверьте состояние:

```bash
docker compose ps
curl --fail http://localhost:8080/health
```

После запуска доступны:

- Mini App: <http://localhost/>;
- Swagger UI: <http://localhost:8080/docs/>;
- Kafka UI: <http://localhost:8088/>;
- API напрямую: <http://localhost:8080/>.

4. Для запуска MAX Bot укажите в `bot/.env` настоящий `MAX_BOT_TOKEN`, публичный `WEBHOOK_URL` и `WEBHOOK_SECRET`. Значения `MAX_BOT_TOKEN` и `BOT_SERVICE_TOKEN` должны совпадать с `api/.env`. Затем запустите:

```bash
docker compose up --build -d bot
```

Остановить стек без удаления данных:

```bash
docker compose down
```

Подробная настройка, диагностика и запуск отдельных компонентов описаны в [инструкции по локальной разработке](docs/local-development.md).

## Проверка API

В корне репозитория находится [`DATA-API.yaml`](DATA-API.yaml) со сценарием проверки API.

Для получения MAX Mini App initData используйте генератор `scripts/generate-max-init-data.sh`. Скрипт требует Bash и `openssl`, скрыто запрашивает токен MAX-бота и затем ID пользователя:

```bash
INIT_DATA=$(./scripts/generate-max-init-data.sh)
```

Скрипт выводит готовую строку для заголовка `X-Max-Init-Data`. Сгенерированные initData ограничены настройкой `MINIAPP_INIT_DATA_MAX_AGE` (по умолчанию один час), поэтому перед проверкой следует формировать новое значение. Токен и полученные initData не должны попадать в Git.

## Документация

- [Архитектура системы](docs/architecture.md)
- [Продукт и пользовательский сценарий](docs/product.md)
- [Обзор владения данными](docs/data-model.md)
- [Рейтинг окружения](docs/impact-engine.md)
- [Документация API](api/README.md)
- [Документация Bot](bot/README.md)
- [Документация Mini App](miniapp/README.md)
- [HTTP OpenAPI](api/docs/openapi.yaml)
- [Kafka AsyncAPI](docs/asyncapi/notifications.yaml)
- [Локальная разработка](docs/local-development.md)
- [Исходники и изображения диаграмм](docs/diagrams/README.md)
- [Правила работы с репозиторием](CONTRIBUTING.md)
