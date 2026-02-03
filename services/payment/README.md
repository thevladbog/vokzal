# Payment Service

Микросервис для обработки платежей через различные провайдеры в системе Вокзал.ТЕХ.

## Функционал

### Поддерживаемые методы оплаты
- **Tinkoff Acquiring** — оплата банковскими картами
- **СБП (Система быстрых платежей)** — оплата по QR коду
- **Наличные** — оплата на кассе

### Возможности
- Инициализация платежей
- Проверка статуса платежей
- Обработка webhooks от провайдеров
- Автоматическая публикация событий в NATS
- История всех платежей

## API Endpoints

### Initialize Payments

```bash
# Инициализировать Tinkoff платёж
POST /v1/payments/tinkoff/init
{
  "ticket_id": "uuid",
  "amount": 1500.00,
  "description": "Билет Ростов-Казань"
}

# Ответ
{
  "data": {
    "id": "uuid",
    "amount": 1500.00,
    "status": "processing",
    "payment_url": "https://securepay.tinkoff.ru/...",
    "external_id": "tinkoff_payment_id"
  }
}

# Инициализировать СБП платёж
POST /v1/payments/sbp/init
{
  "ticket_id": "uuid",
  "amount": 1500.00,
  "description": "Билет Ростов-Казань"
}

# Ответ
{
  "data": {
    "id": "uuid",
    "amount": 1500.00,
    "status": "processing",
    "qr_code": "data:image/svg+xml;base64,...",
    "external_id": "sbp_payment_id"
  }
}

# Создать наличную оплату
POST /v1/payments/cash/init
{
  "ticket_id": "uuid",
  "amount": 1500.00
}
```

### Get Payment

```bash
# Получить платёж по ID
GET /v1/payments/:id

# Проверить статус платежа (обновить из провайдера)
GET /v1/payments/:id/status

# Получить платежи по билету
GET /v1/payments?ticket_id=uuid

# Список всех платежей
GET /v1/payments/list?limit=50
```

### Webhooks

```bash
# Webhook от Tinkoff
POST /v1/webhooks/tinkoff
```

## Workflow платежей

### Tinkoff (карты)
1. POS → `POST /payments/tinkoff/init` (amount, description)
2. Payment service → Tinkoff API → получить PaymentURL
3. Открыть PaymentURL на экране покупателя
4. Покупатель вводит данные карты
5. Tinkoff отправляет webhook → `/webhooks/tinkoff`
6. Обновить статус платежа → отправить событие `payment.confirmed`
7. Ticket service завершает продажу

### СБП (QR код)
1. POS → `POST /payments/sbp/init` (amount, description)
2. Payment service → СБП API → получить QR код
3. Показать QR код на экране покупателя
4. Покупатель сканирует QR в банковском приложении
5. Polling: `GET /payments/:id/status` (каждые 3 сек)
6. Статус "confirmed" → завершить продажу

### Наличные
1. POS → `POST /payments/cash/init` (amount)
2. Payment service создаёт запись (status: confirmed)
3. Отправить событие `payment.confirmed`
4. Ticket service завершает продажу

## NATS События

### Публикуемые события
- `payment.confirmed` — платёж подтверждён

## Конфигурация

```yaml
server:
  port: "8085"
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

tinkoff:
  terminal_key: "YOUR_TERMINAL_KEY"
  password: "YOUR_PASSWORD"
  api_url: "https://securepay.tinkoff.ru/v2"

sbp:
  merchant_id: "YOUR_MERCHANT_ID"
  api_key: "YOUR_API_KEY"
  api_url: "https://api.sbp.nspk.ru"
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
docker build -t vokzal/payment-service:latest .

# Запустить
docker run -p 8085:8085 \
  -e VOKZAL_PAYMENT_DATABASE_HOST=postgres \
  -e VOKZAL_PAYMENT_TINKOFF_TERMINAL_KEY=your_key \
  vokzal/payment-service:latest
```

## Зависимости

- Go 1.23+
- PostgreSQL 15+
- NATS 2.10+
- Tinkoff Acquiring API v2
- СБП API

## Структура БД

### payments
- `id` (UUID PK)
- `ticket_id` (UUID FK, nullable, index)
- `amount` (DECIMAL)
- `currency` (VARCHAR, default: RUB)
- `method` (VARCHAR: card, sbp, cash)
- `provider` (VARCHAR: tinkoff, sbp, manual)
- `status` (VARCHAR: pending, processing, confirmed, failed, refunded)
- `external_id` (VARCHAR, index) — ID у провайдера
- `payment_url` (VARCHAR) — ссылка для оплаты
- `qr_code` (TEXT) — QR код для СБП
- `error_msg` (TEXT)
- `confirmed_at` (TIMESTAMP)
- `refunded_at` (TIMESTAMP)
- `refund_amount` (DECIMAL)
- `metadata` (JSONB)

## Tinkoff Acquiring API

### Init Payment
```
POST https://securepay.tinkoff.ru/v2/Init
{
  "TerminalKey": "...",
  "Amount": 150000,  // в копейках
  "OrderId": "uuid",
  "Description": "Билет на автобус",
  "Token": "sha256_hash"
}
```

### Get State
```
POST https://securepay.tinkoff.ru/v2/GetState
{
  "TerminalKey": "...",
  "PaymentId": "...",
  "Token": "sha256_hash"
}
```

### Webhook
Tinkoff отправляет POST запрос на указанный URL при изменении статуса:
```json
{
  "TerminalKey": "...",
  "PaymentId": "...",
  "Status": "CONFIRMED",
  "OrderId": "uuid",
  "Amount": 150000
}
```

## СБП API

### Generate QR
```
POST https://api.sbp.nspk.ru/qr/generate
{
  "merchantId": "...",
  "amount": 1500.00,
  "currency": "RUB",
  "purpose": "Билет на автобус",
  "qrType": "dynamic"
}
```

### Check Status
```
POST https://api.sbp.nspk.ru/payment/status
{
  "merchantId": "...",
  "paymentId": "..."
}
```

## Безопасность

### Tinkoff Token
Генерация SHA-256 хеша от отсортированных параметров + Password:
```go
params := {TerminalKey, Amount, OrderId, Password}
token := sha256(join(sort(params)))
```

### Webhook Verification
Проверка подписи webhook запросов от Tinkoff.

## Health Check

```bash
GET /health
```

Ответ:
```json
{
  "status": "ok",
  "service": "payment",
  "version": "1.0.0"
}
```

---

© 2025 Вокзал.ТЕХ
