# Разработка и запуск Bot

## Требования

- Go 1.27;
- доступный GeoLogic Monitor API;
- Kafka для обработки уведомлений;
- MAX bot token и публичный HTTPS URL для webhook;
- Docker Compose для рекомендуемого локального запуска;
- Task и Cloudflare Tunnel — необязательные оболочки для команд запуска и туннеля.

## Подготовка окружения

```bash
cp bot/.env.template bot/.env
```

Замените все значения-заглушки. Файл `bot/.env` не коммитится. `MAX_BOT_TOKEN` и `BOT_SERVICE_TOKEN` обязательны для загрузчика конфигурации; `WEBHOOK_URL` и `WEBHOOK_SECRET` должны быть непустыми для рабочего webhook, хотя код конфигурации отдельно это не проверяет.

### HTTP-сервер и API

| Переменная | По умолчанию в коде | Назначение |
| --- | --- | --- |
| `API_PORT` | `8080` | Порт HTTP-сервера Bot; шаблон и Compose используют `8082` |
| `READ_TIMEOUT` | `5s` | Таймаут чтения HTTP-запроса |
| `WRITE_TIMEOUT` | `10s` | Таймаут записи HTTP-ответа |
| `IDLE_TIMEOUT` | `60s` | Таймаут keep-alive соединения |
| `API_BASE_URL` | `http://localhost:8080` | Базовый URL GeoLogic Monitor API; Compose переопределяет на `http://api:8080` |
| `API_TIMEOUT` | `5s` | Таймаут исходящего запроса к API |

HTTP-сервер публикует только `POST /webhook`; health endpoint у Bot отсутствует.

### MAX и сервисная авторизация

| Переменная | По умолчанию | Назначение |
| --- | --- | --- |
| `WEBHOOK_URL` | пусто | Публичный HTTPS URL с путём `/webhook`, регистрируемый в MAX |
| `WEBHOOK_SECRET` | пусто | Общий секрет для проверки входящего webhook |
| `MAX_BOT_TOKEN` | обязательна | Токен для MAX Bot API и регистрации webhook |
| `BOT_SERVICE_TOKEN` | обязательна | Bearer-токен запроса регистрации пользователя в API |

Назначение токенов и границы доверия описаны в [документе по безопасности](security.md).

### Kafka

| Переменная | По умолчанию в коде | Назначение |
| --- | --- | --- |
| `KAFKA_BROKERS` | `kafka:19092` | Список брокеров через запятую |
| `KAFKA_TOPIC` | `notifications` | Topic заданий на уведомления |
| `KAFKA_GROUP_ID` | `bot-notifications` | Consumer group Bot |

Шаблон содержит значения Kafka для сети Compose. При запуске Bot на хосте укажите `KAFKA_BROKERS=localhost:9092`.

### Логирование

| Переменная | По умолчанию | Назначение |
| --- | --- | --- |
| `LOG_LEVEL` | `info` | Минимальный уровень сообщений |
| `LOG_FORMAT` | `json` | Формат `json` или `text` |
| `LOG_ADD_SOURCE` | `true` | Добавлять источник записи |

## Запуск с Docker Compose

Из корня репозитория запустите API и Kafka:

```bash
docker compose up -d api kafka
```

В отдельном терминале откройте туннель к опубликованному порту Bot:

```bash
task bot:tunel
# эквивалентно: cloudflared tunnel --url http://localhost:8082
```

Запишите выданный HTTPS-адрес с суффиксом `/webhook` в `WEBHOOK_URL`, затем запустите Bot:

```bash
task bot:up
docker compose logs -f bot
```

`task bot:up` собирает и запускает только контейнер Bot. Kafka не указана в его Compose `depends_on`, поэтому её нужно запускать отдельно для уведомлений.

Остановка Bot не удаляет данные Kafka:

```bash
task bot:down
```

## Запуск напрямую

Подготовьте `bot/.env`, запустите API и Kafka с адресами, доступными с хоста, и откройте HTTPS-туннель к выбранному `API_PORT`. Затем выполните:

```bash
cd bot
KAFKA_BROKERS=localhost:9092 go run ./cmd/bot
```

Загрузчик читает `.env` из текущего рабочего каталога. При успешном старте в журнале появляются записи о запуске HTTP-сервера и регистрации webhook.

## Проверка

```bash
cd bot
go test ./...
```

Автоматические тесты не требуют действующих MAX, API или Kafka. Для ручной проверки сценария `/start` нужны корректные токены, доступный API и зарегистрированный публичный webhook.

## Диагностика

- Bot завершается при старте: проверьте обязательные токены, формат `LOG_LEVEL` и `LOG_FORMAT`, а также ошибку регистрации webhook.
- MAX не доставляет обновления: проверьте HTTPS-доступность `WEBHOOK_URL`, суффикс `/webhook` и совпадение `WEBHOOK_SECRET`.
- `/start` не регистрирует чат: проверьте `API_BASE_URL`, одинаковое значение `BOT_SERVICE_TOKEN` в Bot и API и доступность API.
- Consumer не получает сообщения: проверьте `KAFKA_BROKERS`, `KAFKA_TOPIC`, запуск Kafka и журналы Bot.
- Consumer остановился после сообщения: проверьте соответствие JSON контракту AsyncAPI, поддерживаемый тип уведомления и доступность MAX.
- Сообщение пропущено с предупреждением `chat.not.found`: пользовательский чат больше недоступен в MAX; такое задание считается обработанным.

Секреты и полные тела webhook или Kafka-сообщений не следует добавлять в журналы или коммитить в репозиторий.
