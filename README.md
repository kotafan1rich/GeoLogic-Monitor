# GeoLogic Monitor

Сервис в MAX для владельцев кофеен, кафе и небольших ресторанов Санкт-Петербурга: следит за изменениями вокруг торговой точки и отправляет уведомления о новых конкурентах и ближайших мероприятиях.

Целевой MVP позволяет через MAX Mini App добавить несколько точек, посмотреть и удалить точку, увидеть текущий рейтинг и его историю за выбранный период (по умолчанию три месяца). API уже сохраняет рейтинг при создании точки и ведёт историю; автоматический пересчёт по расписанию ещё не реализован. Бот должен отправлять уведомления; технический `/start` регистрирует чат и открывает Mini App.

Пользовательские и внутренние HTTP-методы API, авторизация Mini App по MAX initData, расчёт при создании точки и история рейтинга реализованы. Планировщик пересчёта и frontend Mini App остаются целевыми возможностями.

## Архитектура и стек

Три независимых Go-модуля: `api` хранит постоянные данные и рассчитывает рейтинг, `ingestion` со своей БД загружает окружение и публикует сигналы в Kafka, `bot` обрабатывает технический `/start` через MAX Long Polling и отправляет сообщения из Kafka. Mini App обращается напрямую к API.

Go 1.27, стандартный net/http, pgx/pgxpool и ручной SQL, PostgreSQL 18 + PostGIS 3.6, franz-go, OSRM с foot.lua, log/slog, Docker Compose и Taskfile. Frontend-стек пока не определяется.

## Текущее состояние

В `api` реализованы миграции постоянных данных, методы OpenAPI, сервисная и Mini App-авторизация, геокодинг, расчёт рейтинга через OSRM и его история. Не реализован плановый пересчёт. Полный пользовательский сценарий также зависит от ещё не завершённых `bot`, `ingestion` и frontend Mini App.

Локальная инфраструктура OSRM подготовлена отдельно от Go-сервисов: контейнер автоматически строит пешеходный граф, проверяет его trial-запуском и сохраняет в Docker volume. Инструкции находятся в [README OSRM](osrm/README.md). API и OSRM можно собрать из Compose; запуск всего стека остаётся недоступен до появления исполняемых файлов и env-конфигурации `bot` и `ingestion`.

## Локальный запуск

Подробная настройка окружения и текущее состояние сервисов описаны в [инструкции по локальной разработке](docs/local-development.md).

### Установка Task

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

Определённые в проекте команды Docker Compose:

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
task db:up                  # запустить обе базы данных
task db:down                # остановить обе базы данных
task osrm:up                # запустить OSRM
task osrm:down              # остановить OSRM
task logs                   # смотреть общие логи
```

Для первой сборки OSRM или после изменения его Dockerfile и entrypoint используйте:

```bash
docker compose up --build -d osrm
```

## Документация

- [Правила работы с репозиторием](CONTRIBUTING.md)
- [Продукт и пользовательский сценарий](docs/product.md)
- [Архитектура и структура каталогов](docs/architecture.md)
- [Модель данных](docs/data-model.md)
- [Рейтинг окружения](docs/impact-engine.md)
- [OpenAPI](api/docs/openapi.yaml)
- [AsyncAPI](docs/asyncapi/notifications.yaml)
- [Локальная разработка](docs/local-development.md)
- [Локальный OSRM](osrm/README.md)
- [Диаграммы: PlantUML и SVG](docs/diagrams/README.md)
