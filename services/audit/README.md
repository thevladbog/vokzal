# Audit Service

Микросервис для логирования всех операций и соответствия 152-ФЗ в системе Вокзал.ТЕХ.

## Функционал

### Аудит операций
- Логирование всех действий пользователей
- История изменений сущностей
- IP-адрес и User-Agent
- Хранение old/new значений (JSONB)
- Автоматическая подписка на NATS события

### Соответствие 152-ФЗ
- Хранение логов 1+ год
- Отслеживание доступа к ПД
- Возможность предоставления логов регулятору

## API Endpoints

```bash
# Создать запись аудита
POST /v1/audit/log
{
  "entity_type": "ticket",
  "entity_id": "uuid",
  "action": "create",
  "user_id": "uuid",
  "old_value": null,
  "new_value": {"price": 1500, "status": "active"}
}

# Получить лог по ID
GET /v1/audit/:id

# История изменений сущности
GET /v1/audit/entity?entity_type=ticket&entity_id=uuid

# Действия пользователя
GET /v1/audit/user?user_id=uuid&limit=100

# Логи за период
GET /v1/audit/date-range?from=2026-04-01&to=2026-04-15

# Список всех логов
GET /v1/audit/list?limit=100
```

## NATS События

### Подписка
- `audit.log` — автоматическое создание записи аудита

Формат события:
```json
{
  "entity_type": "ticket",
  "entity_id": "uuid",
  "action": "refund",
  "user_id": "uuid",
  "old_value": {"status": "active", "price": 1500},
  "new_value": {"status": "returned", "refund_amount": 1350}
}
```

## Конфигурация

```yaml
server:
  port: "8088"
  mode: "debug"

database:
  host: "localhost"
  port: 5432
  user: "admin"
  password: "vokzal_secret_2026"
  dbname: "vokzal"
  sslmode: "disable"

nats:
  url: "nats://localhost:4222"
  user: "vokzal"
  password: "nats_secret_2026"
```

## Запуск

```bash
go mod download
go run cmd/main.go
```

## Структура БД

### audit_logs
- `id` (UUID PK)
- `entity_type` (VARCHAR: ticket, trip, user, route, etc.)
- `entity_id` (UUID, index)
- `action` (VARCHAR: create, update, delete, refund, etc.)
- `user_id` (UUID FK, index)
- `old_value` (JSONB)
- `new_value` (JSONB)
- `ip_address` (VARCHAR)
- `user_agent` (VARCHAR)
- `created_at` (TIMESTAMP, index)

**Индексы:**
- `(entity_type, entity_id)` — быстрый поиск истории
- `(user_id)` — действия пользователя
- `(created_at)` — временные запросы

## Примеры использования

### История билета

```bash
curl "http://localhost:8088/v1/audit/entity?entity_type=ticket&entity_id=uuid"
```

Ответ:
```json
{
  "data": [
    {
      "id": "uuid",
      "entity_type": "ticket",
      "entity_id": "ticket-uuid",
      "action": "refund",
      "user_id": "user-uuid",
      "old_value": {"status": "active", "price": 1500},
      "new_value": {"status": "returned", "refund_amount": 1350},
      "ip_address": "192.168.1.10",
      "created_at": "2026-04-15T10:30:00Z"
    },
    {
      "action": "create",
      "old_value": null,
      "new_value": {"status": "active", "price": 1500},
      "created_at": "2026-04-14T08:15:00Z"
    }
  ]
}
```

### Действия пользователя за день

```bash
curl "http://localhost:8088/v1/audit/user?user_id=uuid&limit=50"
```

### Все операции за период

```bash
curl "http://localhost:8088/v1/audit/date-range?from=2026-04-01&to=2026-04-15"
```

## Соответствие 152-ФЗ

### Требования
- ✅ Логирование доступа к персональным данным
- ✅ Хранение логов минимум 1 год
- ✅ IP-адрес и User-Agent запросов
- ✅ Old/new значения для аудита изменений

### Отчёты для регулятора
Экспорт всех логов в JSON/CSV:
```bash
curl "http://localhost:8088/v1/audit/date-range?from=2025-01-01&to=2025-12-31" > audit_2025.json
```

## Health Check

```bash
GET /health
```

---

© 2025 Вокзал.ТЕХ
