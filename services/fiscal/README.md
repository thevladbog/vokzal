# Fiscal Service

Микросервис для фискализации операций и соответствия 54-ФЗ в системе Вокзал.ТЕХ.

## Функционал

### Фискализация
- Автоматическая обработка продаж билетов
- Автоматическая обработка возвратов
- Отправка чеков в ОФД через АТОЛ ККТ
- Хранение фискальных чеков (5 лет по 54-ФЗ)

### Z-отчёты
- Ежедневные Z-отчёты
- Автоматический запуск в 00:00
- Хранение истории отчётов
- Статистика продаж/возвратов

### Интеграция с ККТ
- Работа через локальный агент
- АТОЛ драйвер
- Проверка статуса ККТ
- Обработка ошибок фискализации

## API Endpoints

### Receipts

```bash
# Получить чек по ID
GET /v1/receipts/:id

# Получить чеки по билету
GET /v1/receipts?ticket_id=uuid
```

### Z-Reports

```bash
# Создать Z-отчёт
POST /v1/z-reports
{
  "date": "2026-04-15"
}

# Получить Z-отчёт по дате
GET /v1/z-reports/date?date=2026-04-15

# Список Z-отчётов
GET /v1/z-reports?limit=30
```

Ответ:
```json
{
  "data": {
    "id": "uuid",
    "date": "2026-04-15",
    "kkt_serial": "00001234567890",
    "shift_number": 42,
    "total_sales": 125000.00,
    "total_refunds": 5000.00,
    "sales_count": 85,
    "refunds_count": 3,
    "status": "completed",
    "fiscal_sign": "1234567890"
  }
}
```

### KKT Status

```bash
# Статус ККТ
GET /v1/kkt/status
```

## NATS События

### Подписки
- `ticket.sold` — обработка продажи билета
- `ticket.returned` — обработка возврата билета

### Обработка событий
1. Получение события из NATS
2. Создание записи в БД (status: pending)
3. Отправка на ККТ через локальный агент
4. Обновление статуса (confirmed/failed)
5. Сохранение OFD URL и фискального признака

## Конфигурация

```yaml
server:
  port: "8084"
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

logger:
  level: "debug"

atol:
  company_inn: "1234567890"
  company_name: "ООО «Вокзал.ТЕХ»"
  tax_system: "osn"

local_agent:
  url: "http://localhost:8081"
```

## Запуск

### Локально

```bash
# Установить зависимости
go mod download

# Запустить
go run cmd/main.go
```

### Docker

```bash
# Собрать образ
docker build -t vokzal/fiscal-service:latest .

# Запустить
docker run -p 8084:8084 \
  -e VOKZAL_FISCAL_DATABASE_HOST=postgres \
  -e VOKZAL_FISCAL_LOCAL_AGENT_URL=http://agent:8081 \
  vokzal/fiscal-service:latest
```

## Зависимости

- Go 1.23+
- PostgreSQL 15+
- NATS 2.10+
- Локальный агент (порт 8081)
- АТОЛ ККТ

## Структура БД

### fiscal_receipts
- `id` (UUID PK)
- `ticket_id` (UUID FK, index)
- `type` (VARCHAR: sale, refund)
- `amount` (DECIMAL)
- `ofd_url` (VARCHAR)
- `kkt_serial` (VARCHAR)
- `fiscal_sign` (VARCHAR)
- `status` (VARCHAR: pending, sent, confirmed, failed)
- `error_msg` (TEXT)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### z_reports
- `id` (UUID PK)
- `date` (DATE, index)
- `kkt_serial` (VARCHAR)
- `shift_number` (INTEGER)
- `total_sales` (DECIMAL)
- `total_refunds` (DECIMAL)
- `sales_count` (INTEGER)
- `refunds_count` (INTEGER)
- `status` (VARCHAR: pending, completed, failed)
- `fiscal_sign` (VARCHAR)

## Локальный агент API

Fiscal service взаимодействует с локальным агентом:

```bash
# Печать чека
POST http://localhost:8081/kkt/receipt
{
  "operation": "sell",
  "items": [...],
  "payment": {...},
  "company": {...}
}

# Z-отчёт
POST http://localhost:8081/kkt/z-report

# Статус ККТ
GET http://localhost:8081/kkt/status
```

## Соответствие 54-ФЗ

### Требования
- ✅ Все продажи через ККТ
- ✅ Отправка чеков в ОФД в течение 30 сек
- ✅ Z-отчёты ежедневно
- ✅ Хранение чеков 5 лет
- ✅ Фискальный накопитель

### Обработка ошибок
- Повторная отправка при сбое связи
- Логирование всех ошибок
- Алерты для администраторов
- Ручное исправление через API

## Автоматизация

### Ежедневные задачи
- Z-отчёт в 00:00 (cron)
- Проверка статуса ККТ
- Отправка отчётов в ОФД

## Health Check

```bash
GET /health
```

Ответ:
```json
{
  "status": "ok",
  "service": "fiscal",
  "version": "1.0.0"
}
```

---

© 2025 Вокзал.ТЕХ
