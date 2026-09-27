# Локальная разработка

Проект запускается через Docker Compose; Taskfile предоставляет короткие команды для тех же операций. Детальные настройки находятся в документации соответствующих сервисов.

## Требования

- Docker с поддержкой Compose;
- Go 1.27 для запуска Go-модулей без Docker;
- Node.js, если Mini App запускается вне контейнера;
- [Task](https://taskfile.dev/docs/installation) — необязательно.

## Подготовка окружения

```bash
cp .env.template .env
cp api/.env.template api/.env
cp ingestion/.env.template ingestion/.env
cp bot/.env.template bot/.env
```

Замените заглушки токенов и внешних ключей. Локальные `.env` не коммитятся.

Корневой `.env` задаёт опубликованные порты и параметры общей инфраструктуры. Env-файлы приложений содержат их собственные подключения и секреты.

## Запуск всей системы

```bash
docker compose up --build -d
# или
task all:up
```

Проверить состояние и журналы:

```bash
docker compose ps
docker compose logs -f
# или
task logs
```

Остановка без удаления volumes:

```bash
task all:down
```

## Запуск компонентов

```bash
task postgres:up
task osrm:up
task api:up
task ingestion:up
task kafka:up
task bot:up
```

Первичная подготовка графа OSRM может потребовать до 1 ГБ памяти и занимает больше времени обычного запуска. Подробности: [локальный OSRM](../osrm/README.md).

Настройка и диагностика сервисов описаны в документации разработки [API](../api/docs/development.md) и [Bot](../bot/docs/development.md).

## Проверки

Go-модули проверяются независимо:

```bash
cd api && go test ./...
cd ingestion && go test ./...
cd bot && go test ./...
```

Mini App:

```bash
cd miniapp
npm ci
npm run build
```

Общего `go.mod` в корне намеренно нет. Команды Go нужно запускать из каталога соответствующего модуля.

## Состояние сквозного сценария

API, Mini App, загрузочные jobs ingestion, webhook MAX и Kafka consumer bot реализованы. Producer уведомлений в ingestion пока отсутствует, поэтому полный путь «обнаружение изменения → Kafka → сообщение MAX» недоступен.
