# SQL Миграции для Вокзал.ТЕХ

## Структура

Миграции организованы в формате `{номер}_{описание}.{up|down}.sql`:

- `001_initial_schema.up.sql` - создание начальной схемы
- `001_initial_schema.down.sql` - откат начальной схемы

## Применение миграций

### Вручную (psql)

```bash
# Применить миграцию
psql -U admin -d vokzal -f 001_initial_schema.up.sql

# Откатить миграцию
psql -U admin -d vokzal -f 001_initial_schema.down.sql
```

### С помощью golang-migrate

```bash
# Установить golang-migrate
brew install golang-migrate

# Применить все миграции
migrate -path /Users/thevladbog/PRSOME/vokzal/infra/migrations \
        -database "postgresql://admin:vokzal_secret_2026@localhost:5432/vokzal?sslmode=disable" \
        up

# Откатить последнюю миграцию
migrate -path /Users/thevladbog/PRSOME/vokzal/infra/migrations \
        -database "postgresql://admin:vokzal_secret_2026@localhost:5432/vokzal?sslmode=disable" \
        down 1

# Проверить статус
migrate -path /Users/thevladbog/PRSOME/vokzal/infra/migrations \
        -database "postgresql://admin:vokzal_secret_2026@localhost:5432/vokzal?sslmode=disable" \
        version
```

## Таблицы

### Основные сущности
- `stations` - автовокзалы
- `buses` - автобусы
- `drivers` - водители
- `routes` - маршруты
- `schedules` - расписание
- `trips` - рейсы
- `seats` - места в автобусах
- `users` - пользователи системы

### Продажа билетов
- `tickets` - билеты
- `blocking_rules` - правила блокировки мест
- `boarding_events` - события посадки

### Фискализация и аудит
- `fiscal_receipts` - фискальные чеки (54-ФЗ)
- `audit_logs` - журнал аудита (152-ФЗ)

### Уведомления
- `notifications` - уведомления (SMS, Email, Telegram)
- `announcements` - голосовые оповещения (TTS)

### Дополнительные
- `sessions` - сессии пользователей
- `document_templates` - шаблоны документов

## Пользователь по умолчанию

После применения миграции создаётся администратор:
- **Username:** `admin`
- **Password:** `admin123`

⚠️ **Важно:** Сменить пароль после первого входа!

## Индексы

Все критически важные индексы созданы:
- `idx_tickets_trip_seat` - проверка занятости мест
- `idx_trips_schedule_date` - запросы расписания
- `idx_audit_entity` - поиск в журнале аудита
- И многие другие...

## Триггеры

Автоматическое обновление `updated_at` для всех основных таблиц.

---

© 2025 Вокзал.ТЕХ
