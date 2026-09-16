# GeoLogic Monitor

Сервис в MAX для владельцев кофеен, кафе и небольших ресторанов Санкт-Петербурга: следит за изменениями вокруг торговой точки и отправляет уведомления о новых конкурентах и ближайших мероприятиях.

MVP позволяет добавить несколько точек, посмотреть список и удалить точку. После выбора адреса мониторинг включается автоматически. Рейтинг окружения на основе инфраструктуры служит вспомогательной аналитикой.

## Архитектура и стек

Три независимых Go-модуля: `api` хранит данные, `ingestion` загружает окружение и публикует сигналы в Kafka, `bot` обслуживает MAX Long Polling и отправляет сообщения из Kafka.

Go 1.27, стандартный net/http, pgx/pgxpool и ручной SQL, PostgreSQL 18 + PostGIS 3.6, franz-go, OSRM с foot.lua, log/slog, Docker Compose и Taskfile. Frontend-стек пока не определяется.

## Документация

- [Правила работы с репозиторием](CONTRIBUTING.md)
- [Продукт и пользовательский сценарий](docs/product.md)
- [Архитектура и структура каталогов](docs/architecture.md)
- [Модель данных](docs/data-model.md)
- [Рейтинг окружения](docs/impact-engine.md)
- [OpenAPI](docs/openapi/geologic.yaml)
- [AsyncAPI](docs/asyncapi/notifications.yaml)
- [Локальная разработка и демо](docs/local-development.md)
- [Этапы реализации](docs/implementation-plan.md)
- [Диаграммы: PlantUML и SVG](docs/diagrams/README.md)

## Статус

Подготовлена документация хакатонного MVP; сервисы, миграции, Docker Compose и Taskfile пока не реализованы. История уведомлений, Mini App и другие города не входят в объём работ.

Организация документации основана на [референсном проекте](https://github.com/talense-tasks/backend-trainee-assignment-autumn-2026-kotafan1rich-aee3bbdd/tree/main/docs).
