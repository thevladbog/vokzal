# Ticket Service

Микросервис для продажи, возврата билетов и управления посадкой в системе Вокзал.ТЕХ.

## Функционал

### Продажа билетов
- Продажа билетов с выбором места или без места
- Проверка доступности мест
- Генерация QR и ШК кодов
- Поддержка различных методов оплаты
- События в NATS для фискализации

### Возврат билетов
- Возврат с автоматическим расчётом штрафа
- Блокировка возврата после начала посадки
- Учёт времени до отправления:
  - \> 24 часа: 10% штраф
  - 12-24 часа: 20% штраф
  - < 12 часов: 30% штраф
- Аудит всех операций возврата

### Посадка
- Начало посадки (блокировка возвратов)
- Отметка посадки по QR/ШК
- Статистика посадки (всего/посажено)
- Предотвращение дублирования отметок

## API Endpoints

### Tickets

```bash
# Продать билет
POST /v1/tickets/sell
{
  "trip_id": "uuid",
  "seat_id": "uuid",
  "passenger_name": "Иванов Иван Иванович",
  "passenger_doc": "4500 123456",
  "phone": "+79001234567",
  "email": "ivan@example.com",
  "price": 1500.00,
  "payment_method": "card"
}

# Список билетов на рейс
GET /v1/tickets?trip_id=uuid

# Получить билет по ID
GET /v1/tickets/:id

# Получить билет по QR коду
GET /v1/tickets/qr?qr_code=TK12345678

# Возврат билета
POST /v1/tickets/:id/refund
```

Ответ на возврат:
```json
{
  "message": "Ticket refunded successfully",
  "data": {
    "original_amount": 1500.00,
    "penalty": 150.00,
    "refund_amount": 1350.00
  }
}
```

### Boarding

```bash
# Начать посадку
POST /v1/boarding/start
{
  "trip_id": "uuid"
}

# Отметить посадку пассажира
POST /v1/boarding/mark
{
  "ticket_id": "uuid",
  "user_id": "uuid",
  "scan_method": "qr"
}

# Статус посадки
GET /v1/boarding/status?trip_id=uuid
```

Ответ на статус:
```json
{
  "data": {
    "trip_id": "uuid",
    "boarding_active": true,
    "started_at": "2026-04-15T10:30:00Z",
    "total_tickets": 45,
    "boarded_count": 32
  }
}
```

## NATS События

### Публикуемые события
- `ticket.sold` — билет продан
- `ticket.returned` — билет возвращён
- `boarding.started` — посадка началась
- `audit.log` — запись аудита

## Конфигурация

```yaml
server:
  port: "8083"
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

business:
  refund_penalty:
    over_24_hours: 0.10
    between_12_24: 0.20
    under_12_hours: 0.30
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
docker build -t vokzal/ticket-service:latest .

# Запустить
docker run -p 8083:8083 \
  -e VOKZAL_TICKET_DATABASE_HOST=postgres \
  -e VOKZAL_TICKET_NATS_URL=nats://nats:4222 \
  vokzal/ticket-service:latest
```

## Зависимости

- Go 1.23+
- PostgreSQL 15+
- NATS 2.10+
- Gin v1.10+
- GORM v1.25+

## Структура БД

### tickets
- `id` (UUID PK)
- `trip_id` (UUID FK)
- `seat_id` (UUID FK, nullable)
- `passenger_name` (VARCHAR)
- `passenger_doc` (VARCHAR)
- `phone` (VARCHAR)
- `email` (VARCHAR)
- `price` (DECIMAL)
- `status` (VARCHAR: active, returned, cancelled)
- `payment_method` (VARCHAR)
- `qr_code` (VARCHAR, unique)
- `bar_code` (VARCHAR, unique)
- `refunded_at` (TIMESTAMP)
- `refund_amount` (DECIMAL)
- `refund_penalty` (DECIMAL)

### boarding_events
- `id` (UUID PK)
- `trip_id` (UUID FK, unique)
- `started_at` (TIMESTAMP)
- `started_by` (UUID FK)

### boarding_marks
- `id` (UUID PK)
- `ticket_id` (UUID FK)
- `marked_at` (TIMESTAMP)
- `marked_by` (UUID FK)
- `scan_method` (VARCHAR: qr, barcode, manual)

## Бизнес-логика

### Проверки при продаже
1. Доступность места (если указан seat_id)
2. Валидация данных пассажира
3. Проверка суммы (price > 0)

### Проверки при возврате
1. Билет в статусе "active"
2. Посадка НЕ начата
3. Расчёт штрафа по времени до отправления

### Проверки при посадке
1. Билет в статусе "active"
2. Посадка начата
3. Билет не отмечен ранее

## Health Check

```bash
GET /health
```

Ответ:
```json
{
  "status": "ok",
  "service": "ticket",
  "version": "1.0.0"
}
```

## Соответствие 152-ФЗ

- Все операции с ПД логируются в audit сервис
- Персональные данные (ФИО, паспорт, телефон) шифруются
- Согласие на обработку ПД при покупке

---

© 2025 Вокзал.ТЕХ
