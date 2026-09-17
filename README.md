# GeoLogic Monitor

Сервис в MAX для владельцев кофеен, кафе и небольших ресторанов Санкт-Петербурга: следит за изменениями вокруг торговой точки и отправляет уведомления о новых конкурентах и ближайших мероприятиях.

MVP позволяет добавить несколько точек, посмотреть список и удалить точку. После выбора адреса мониторинг включается автоматически. Рейтинг окружения на основе инфраструктуры служит вспомогательной аналитикой.

## Архитектура и стек

Три независимых Go-модуля: `api` хранит данные, `ingestion` загружает окружение и публикует сигналы в Kafka, `bot` обслуживает MAX Long Polling и отправляет сообщения из Kafka.

Go 1.27, стандартный net/http, pgx/pgxpool и ручной SQL, PostgreSQL 18 + PostGIS 3.6, franz-go, OSRM с foot.lua, log/slog, Docker Compose и Taskfile. Frontend-стек пока не определяется.

## Установка Task

Для запуска команд проекта нужен [Task](https://taskfile.dev/docs/installation). Установить его можно одним из способов:

```bash
# Linux через Snap
sudo snap install task --classic

# macOS или Linux через Homebrew
brew install go-task/tap/go-task

# При установленном Go
go install github.com/go-task/task/v3/cmd/task@latest

# Windows через WinGet
winget install Task.Task
```

Проверить установку и посмотреть доступные команды:

```bash
task --version
task --list
```

Основные команды для Docker Compose:

```bash
task all:up                 # запустить всё
task all:down               # остановить всё
task api:up                 # запустить отдельный сервис
task api:down               # остановить отдельный сервис
task ingestion:up
task ingestion:down
task bot:up
task bot:down
task api-db:up
task api-db:down
task ingestion-db:up
task ingestion-db:down
task logs                   # смотреть общие логи
```

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

Подготовлены документация хакатонного MVP, базовые Docker Compose и Taskfile. Реализация сервисов и миграций продолжается. История уведомлений, Mini App и другие города не входят в объём работ.

Организация документации основана на [референсном проекте](https://github.com/talense-tasks/backend-trainee-assignment-autumn-2026-kotafan1rich-aee3bbdd/tree/main/docs).
