# GeoLogic Monitor

Сервис в MAX для владельцев кофеен, кафе и небольших ресторанов Санкт-Петербурга: следит за изменениями вокруг торговой точки и отправляет уведомления о новых конкурентах и ближайших мероприятиях.

Целевой MVP позволяет через MAX Mini App добавить несколько точек, посмотреть и удалить точку, увидеть текущий рейтинг и его историю за выбранный период (по умолчанию три месяца). API сохраняет рейтинг при создании точки, ведёт историю и автоматически пересчитывает рейтинги по расписанию. Бот должен отправлять уведомления; технический `/start` регистрирует чат и открывает Mini App.

Пользовательские и внутренние HTTP-методы API, авторизация Mini App по MAX initData, расчёт при создании точки, история рейтинга, плановый пересчёт и frontend Mini App реализованы.

## Архитектура и стек

Три независимых Go-модуля: `api` хранит постоянные данные и рассчитывает рейтинг, `ingestion` загружает окружение, `bot` принимает webhook MAX и обрабатывает технический `/start`. Целевая цепочка уведомлений связывает ingestion и bot через Kafka; producer и consumer пока не реализованы. Mini App обращается напрямую к API.

Go 1.27, стандартный net/http, pgx/pgxpool и ручной SQL, PostgreSQL 18 + PostGIS 3.6, OSRM с foot.lua, log/slog, Docker Compose и Taskfile; Mini App использует React, TypeScript и Vite. Для целевой Kafka-интеграции выбран franz-go.

## Текущее состояние

В `api` реализованы миграции постоянных данных, методы OpenAPI, сервисная и Mini App-авторизация, геокодинг, расчёт рейтинга через OSRM, его история и плановый пересчёт. Bot реализует webhook MAX и `/start`, ingestion — загрузку данных по расписанию, frontend Mini App — работу с точками и рейтингом. Не завершена доставка уведомлений через Kafka.

Локальная инфраструктура OSRM подготовлена отдельно от Go-сервисов: контейнер автоматически строит пешеходный граф, проверяет его trial-запуском и сохраняет в Docker volume. Инструкции находятся в [README OSRM](osrm/README.md). Все приложения имеют исполняемые файлы или frontend-сборку и конфигурацию Compose; сквозной сценарий уведомлений остаётся недоступен до реализации Kafka producer/consumer.

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
