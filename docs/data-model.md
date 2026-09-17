# Модель данных

[ER-диаграмма](diagrams/rendered/er-model.svg). Постоянные данные принадлежат API и создаются его миграциями; ingestion заполняет их по HTTP. Кэш конкурентов находится в отдельной БД и создаётся миграциями ingestion.

## Сохранённые модели старого API

| Модель / таблица | Поля |
| --- | --- |
| User / users | id, max_user_id (unique), max_chat_id, created_at, updated_at |
| TrackedLocation / tracked_locations | id, user_id FK, business_type_id FK, address, location, created_at, updated_at |
| InfraType / infra_types | id, slug (unique), name, weight, max_radius, created_at, updated_at |
| InfraObject / infra_objects | id, type_id FK infra_types, location, address, name nullable, created_at, updated_at |
| BusinessType / business_types | id, infra_type_id FK infra_types (unique), created_at, updated_at |
| Event / events | id, location, date, info nullable, provider, external_id, notified_at nullable, created_at, updated_at |

InfraType сохраняет положительные Weight и MaxRadius; InfraObject — GeoPoint, Address, необязательное Name и Type. BusinessType сохраняет связь с одним InfraType, определяющим прямых конкурентов. Имя типа бизнеса для кнопки берётся из InfraType.Name. Начальные записи: кофейня, кафе и ресторан; они загружаются ingestion, ID не зашиваются в бот.

Event сохраняет Date и Info, а не заменяет их на starts_at/title/category. В HTTP дата называется date, координаты — lat/lon. В прошлой заготовке ingestion у Date был ошибочный JSON-тег address; его нельзя переносить в контракт. Название и текст мероприятия в MVP содержатся в info. Отдельных концов события и категорий в модели нет.

User заменяет Telegram ID на MAX ID. У TrackedLocation добавлены тип бизнеса и адрес; user_id индексируется, но не уникален — один владелец может отслеживать несколько точек. Rent и его производные не используются.

## Типы, связи и индексы

Внутренние ID и ссылающиеся на них внешние ключи — UUID. Идентификаторы MAX остаются bigint, внешние идентификаторы провайдеров — text. `created_at` и `updated_at` хранятся в БД для диагностики и журналирования, но не возвращаются через HTTP API. Даты — timestamptz UTC, календарные границы — Europe/Moscow. GeoPoint хранится как geometry(Point,4326); PostGIS и OSRM принимают lon,lat, HTTP использует именованные lat/lon. Поиск в метрах выполняется через geography, GiST-индекс должен соответствовать выражению location::geography.

Уникальность events — (provider, external_id). Индексы: tracked_locations.user_id, infra_objects.type_id, events.date, пространственные индексы infra_objects и events. Удаление точки физическое. API не обращается к БД ingestion, поэтому кэш удалённой точки исчезает при очередной TTL-очистке. Типы с существующими объектами/точками не удаляются: загрузочное API MVP не предоставляет удаления справочников.

Upsert типов инфраструктуры выполняется по slug, объектов инфраструктуры — по стабильному id загрузчика, типов бизнеса — по infra_type_id, событий — по provider/external_id. Для infra_objects загрузчик обязан поддерживать стабильные ID; это не идентификатор карточки 2ГИС. Перезапись справочника обновляет поля, но не создаёт дубликаты.

## Отметка события

events.notified_at означает завершённую публикацию для всех подходящих точек, не доставку MAX. Ingestion выбирает ещё не отмеченные события на завтра, рассылает для всех подходящих точек и вызывает API отметки. Upsert не сбрасывает существующую отметку.

Если получателей нет, отметка не ставится. Частичный сбой может повторить часть сообщений; точка, добавленная после завершённой рассылки, не получает старое событие. Таблица доставок в MVP не вводится.

## CompetitorCache

Таблица competitor_cache в БД ingestion: tracked_location_id, external_id, name, type_id, address, location, opened_at (date), notified_at nullable, expires_at, created_at, updated_at. Primary key — (tracked_location_id, external_id), индекс — expires_at. `tracked_location_id` и `type_id` — внешние UUID без FK: связанные таблицы находятся в другой БД.

TTL — 24 часа от фактического получения; повтор не продлевает срок. Данные с истёкшим TTL не читаются. Физическая очистка выполняется ingestion перед ежедневным и ручным запуском; остановка процесса откладывает удаление. Данные не служат историей, не попадают в логи и не копируются в infra_objects.

Кэш создаётся перед Kafka publish; notified_at обновляется после успешной публикации. Наличие свежей записи подавляет повтор кандидата, поэтому сбой между записью и publish может пропустить сообщение. Это принятое упрощение без outbox.

Дата первого наблюдения не считается датой открытия. Доступность даты открытия и условия временного хранения по ключу 2ГИС проверяются при подключении.

## Расчётные модели без таблиц

- GeoPoint: lat, lon.
- InfraTypeFeatures: TypeID, Slug, Weight, MaxRadiusMeters, ObjectCount, NearestMeters, MeanMeters, DistanceSum, Influence, Count100m, Count300m, Count500m, IsCompetitor.
- LocationFeatures: BusinessTypeID, InfraTypes.
- CalculatedRating: Value.
- Route: distance_meters, duration_seconds.
- Notification: получатель, точка, сигнал, маршрут, reasons, recommendations.

Расчётные модели рейтинга принадлежат API и существуют только во время создания точки; рейтинг в таблицах не хранится. Набор признаков сохранён из старого API, но расстояния уточняются OSRM. Их подготовка и использование описаны в [рейтинге](impact-engine.md); транспортные формы — в спецификациях.
