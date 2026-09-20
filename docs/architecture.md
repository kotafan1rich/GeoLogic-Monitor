# Архитектура

Целевая архитектура хакатонного MVP (история рейтинга, планировщик и Mini App проектируются): три независимых Go-модуля со своими go.mod, go.sum, Dockerfile, конфигурацией и DI. Общего Go-пакета и импортов соседних модулей нет.

OSRM — отдельный инфраструктурный контейнер, а не четвёртый Go-модуль. Его образ и подготовка локального графа находятся в каталоге `osrm`.

[Контейнеры](diagrams/rendered/c4-containers.svg) · [Компоненты API](diagrams/rendered/c4-components.svg) · [Мониторинг](diagrams/rendered/monitoring-sequence.svg)

## Ответственность модулей

| Модуль | Ответственность | Доступ к данным |
| --- | --- | --- |
| api | Пользователи, точки, геокодинг, справочники, инфраструктура, события и история рейтинга, плановый пересчёт, авторизация Mini App | Собственная PostgreSQL через pgxpool; OSRM для рейтинга |
| ingestion | Загрузка источников, мониторинг изменений и публикация сигналов | Постоянные данные через внутренний API; отдельная PostgreSQL для кэша 2ГИС |
| bot | MAX Long Polling для технического /start, Kafka consumer и отправка оповещений | HTTP API и Kafka; без БД |

API владеет постоянными данными. Ingestion сначала загружает InfraType, затем BusinessType/InfraObject, затем Event. Загрузка выполняется через внутренние HTTP-методы. Кэш 2ГИС остаётся локальной ответственностью ingestion и хранится в отдельной БД без внешних ключей к БД API. Каждый модуль применяет только собственные миграции.

При добавлении точки API грубо отбирает инфраструктуру через PostGIS, уточняет пешеходные расстояния через OSRM и рассчитывает рейтинг. Создание точки и расчёт атомарны: при недоступности OSRM или общей ошибке расчёта точка не сохраняется. Точка и первый рейтинг сохраняются атомарно. API добавляет новые записи истории по расписанию: по умолчанию 1-го числа месяца в 03:00 Europe/Moscow. Пересчёт независим от мониторинга ingestion. Подробности расписания и ошибок — в [расчёте рейтинга](impact-engine.md).

Polling MAX и consumer Kafka работают параллельно в одном bot-процессе. Бизнес-диалогов нет; /start регистрирует чат и предлагает открыть Mini App; Kafka-сообщения проходят отдельный обработчик отправки MAX. Остановка consumer после ошибки не должна блокировать polling; ошибка видна в логах.

## Слои и каталоги

Зависимости: `cmd → app/DI → handler → service → repository`. Слои взаимодействуют через интерфейсы, зависимости собираются вручную в app/di.go. Domain не зависит от HTTP, pgx и Kafka. DTO обработчиков — в handler/.../dto; специфичные структуры адаптеров — в repository/.../model. Генератор Go-кода пока не выбирается.

```text
api/
├── cmd/api/main.go
├── internal/
│   ├── api/server.go
│   ├── app/{app.go,di.go}
│   ├── args/args.go
│   ├── config/{config.go,database.go,logging.go,server.go}
│   ├── database/
│   │   ├── postgres.go
│   │   └── postgres/{context.go,manager.go,pool.go,tx.go}
│   ├── domain/{user,geo,tracked_location,infra_type,infra_object,business_type,event,location_features,calculated_rating,rating_history}.go
│   ├── errs/pgerrors/
│   ├── handler/
│   │   ├── router.go
│   │   └── {users,locations,geocoding,business_types,infra,events}/{handler.go,dto/}
│   ├── logger/logger.go
│   ├── repository/
│   │   ├── interfaces.go
│   │   ├── database/{users,locations,infra_types,infra,business_types,events,rating}/
│   │   ├── geocoder/{repository.go,model/}
│   │   └── osrm/{repository.go,model/}
│   ├── scheduler/
│   └── service/{interfaces.go,users/,locations/,infra/,business_types/,events/,rating/}
├── migrations/
├── Dockerfile
├── go.mod
└── go.sum

bot/
├── cmd/bot/main.go
├── internal/
│   ├── app/{app.go,di.go}
│   ├── config/
│   ├── domain/
│   ├── handler/{max,kafka}/
│   ├── service/{onboarding,notifications}/
│   ├── repository/{api,max}/
│   └── logger/
├── Dockerfile
├── go.mod
└── go.sum

ingestion/
├── cmd/ingestion/main.go
├── internal/
│   ├── app/{app.go,di.go}
│   ├── config/
│   ├── database/
│   ├── domain/{infra_type,infra_object,business_type,event,competitor,route}.go
│   ├── repository/{api,database/competitor_cache,twogis,events,osm,osrm,kafka}/
│   ├── service/{loading,monitoring}/
│   ├── errs/
│   └── logger/
├── Dockerfile
├── go.mod
└── go.sum

osrm/
├── cmd/prepare.sh
├── data/SanktPetersburg.osm.pbf
├── .env.template
├── Dockerfile
└── README.md
```

В каждом repository при необходимости создаётся model/. Каждый сервис применяет миграции своей БД при запуске.

## Локальная инфраструктура OSRM

Compose собирает OSRM из закреплённого образа `ghcr.io/project-osrm/osrm-backend:v26.9.0-debian`. При сборке образ использует локальный `osrm/data/SanktPetersburg.osm.pbf`, а при его отсутствии автоматически скачивает городскую выгрузку Санкт-Петербурга с BBBike. Исходный файл сохраняется в образе как `/source/SanktPetersburg.osm.pbf`, подготовленный MLD-граф — в volume `osrm_graph`. Для маршрутизации используется профиль `/opt/foot.lua`.

Внутри сети `geo_logic_net` сервис доступен по адресу `http://osrm:5000`. Порт хоста настраивается переменной `OSRM_HTTP_PORT` и по умолчанию равен `5000`. Инструкции по сборке, обновлению карты и проверочным запросам находятся в [README OSRM](../osrm/README.md); общая настройка окружения — в [инструкции по локальной разработке](local-development.md).

## Мониторинг и загрузка

Ingestion получает через API снимок точек и необработанные события рядом с ними. Данные инфраструктуры и событий заполняются адаптером загрузки, основанным на старом проекте; для повторных запусков используются идемпотентные PUT. OSM служит источником графа OSRM. Локальная версия прежнего ingestion содержит заготовки, поэтому готовая совместимость не предполагается: JSON адаптеров сверяется с новым OpenAPI.

Ежедневный запуск ищет конкурентов за вчера и события на завтра по Europe/Moscow. Ingestion может использовать OSRM для маршрута, показываемого в конкретном уведомлении, но не собирает признаки и не рассчитывает рейтинг. Данные 2ГИС не копируются в постоянную таблицу инфраструктуры.

Для события сначала публикуются сообщения для всех подходящих точек, затем через API устанавливается отметка обработки. Точка из снимка может быть удалена во время запуска и получить последнее сообщение.

## Контракты и доступ

[OpenAPI](openapi/geologic.yaml) описывает пользовательские и внутренние методы; [AsyncAPI](asyncapi/notifications.yaml) — задания на отправку. При реализации нового контракта обновляются локальные DTO производителей и потребителей и примеры; общий модуль моделей не вводится. Текущее обновление описывает целевой контракт без изменения Go DTO.

Mini App — браузерный интерфейс внутри MAX: добавление, просмотр и удаление точек, текущий рейтинг и история. Он обращается напрямую к пользовательским HTTP-методам API. Frontend-стек и размещение статических файлов не фиксируются; четвёртый Go-модуль не вводится. Дерево выше описывает целевую структуру, новые компоненты ещё не реализованы.

Mini App передаёт исходную строку `window.WebApp.initData` в `X-Max-Init-Data`. API проверяет подпись по [официальному алгоритму MAX](https://dev.max.ru/docs/webapps/validation), отклоняет повторяющиеся параметры и получает MAX user ID только из проверенного `user.id`. Для проверки используется серверный `MAX_BOT_TOKEN`, отличный от внутренних service token. Дополнительное правило проекта: `auth_date` не должен быть из будущего или старше `MINIAPP_INIT_DATA_MAX_AGE` (по умолчанию 1 час); при истечении срока Mini App предлагает повторное открытие. Неверные данные дают 401. Секреты и initData не логируются и не передаются в клиентский код. Пользовательские запросы идут по HTTPS; при разных origin API разрешает CORS только для настроенного origin Mini App.

Bot использует BOT_SERVICE_TOKEN и X-Max-User-Id из проверенного обновления MAX только для `PUT /api/v1/users/me`: `/start` регистрирует личный чат уведомлений. Mini App не может произвольно задать адресата уведомлений и предлагает запустить бота, если пользователь ещё не зарегистрирован. Заголовок X-Max-User-Id сам по себе не авторизует. Ingestion использует отдельный INGESTION_SERVICE_TOKEN; токены не взаимозаменяемы. API проверяет владельца точки. Секреты и сырые ответы 2ГИС не логируются.

## Kafka и ошибки

Один topic notifications, одна partition для демо, group bot-notifications, franz-go на producer и consumer. Retention — один час; удаление сегментов Kafka асинхронное. Bot подтверждает offset после успешной отправки MAX. Три попытки с паузами 1, 2, 4 секунды; после исчерпания попыток consumer останавливается с ошибкой без commit. Некорректное сообщение также требует вмешательства оператора.

Outbox, DLQ и распределённой транзакции нет: при сбоях возможны дубли и пропуски. Публикация в Kafka не означает доставку или прочтение MAX. Детали отметок и TTL находятся только в [модели данных](data-model.md).

Ошибка отдельной точки не прерывает обработку остальных. При ошибке OSRM во время мониторинга кандидат пропускается; прямое расстояние не показывается как пешее. Недоступность OSRM во время создания точки приводит к 503 и откату создания. Catch-up пропущенных дней и одновременные запуски ingestion в MVP не поддерживаются.
