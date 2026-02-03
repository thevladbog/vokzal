# Document Service

Микросервис для генерации PDF документов (билеты, ПД-2, кастомные шаблоны) в системе Вокзал.ТЕХ.

## Функционал

### Генерация документов
- **Электронный билет** — PDF с QR кодом
- **ПД-2** — проездной документ (форма)
- **Кастомные шаблоны** — настраиваемые PDF
- Хранение в MinIO S3

### Возможности
- Генерация QR/Barcode
- Брендирование (логотипы, шрифты)
- История сгенерированных документов
- Ссылки для скачивания

## API Endpoints

```bash
# Сгенерировать электронный билет
POST /v1/document/ticket
{
  "ticket_id": "uuid",
  "passenger_fio": "Иванов Иван Иванович",
  "passenger_doc": "4500 123456",
  "route": "Ростов-на-Дону — Казань",
  "date": "2026-04-15",
  "time": "08:30",
  "platform": "3",
  "seat": "12",
  "price": 1500.00,
  "qr_code": "TK123456",
  "bar_code": "1234567890123"
}

# Ответ
{
  "data": {
    "id": "uuid",
    "document_type": "ticket",
    "entity_id": "uuid",
    "file_url": "http://localhost:9000/vokzal-documents/tickets/uuid_20260415083000.pdf",
    "file_name": "tickets/uuid_20260415083000.pdf",
    "status": "generated",
    "created_at": "2026-04-15T08:30:00Z"
  }
}

# Сгенерировать ПД-2
POST /v1/document/pd2
{
  "number": "123456",
  "series": "АБ",
  "passenger_fio": "Иванов Иван Иванович",
  "passenger_doc": "4500 123456",
  "route_from": "Ростов-на-Дону",
  "route_to": "Казань",
  "date": "2026-04-15",
  "price": 1500.00,
  "issue_date": "2026-04-14",
  "issuer_name": "Петрова А.В.",
  "bus_number": "А123БВ"
}

# Получить документ по ID
GET /v1/document/:id

# Список документов
GET /v1/document/list?limit=50
```

## Структура БД

### document_templates
- `id` (UUID PK)
- `name` (VARCHAR, unique)
- `type` (VARCHAR: pd2, ticket, invoice, custom)
- `description` (TEXT)
- `content` (TEXT) — HTML шаблон
- `is_active` (BOOLEAN)
- `created_at`, `updated_at` (TIMESTAMP)

### generated_documents
- `id` (UUID PK)
- `template_id` (UUID FK, nullable)
- `document_type` (VARCHAR)
- `entity_id` (UUID, index) — ID билета, рейса и т.д.
- `file_url` (VARCHAR)
- `file_name` (VARCHAR)
- `status` (VARCHAR: generated, archived)
- `created_at` (TIMESTAMP)

## Конфигурация

```yaml
server:
  port: "8089"
  mode: "debug"

minio:
  endpoint: "localhost:9000"
  access_key: "minioadmin"
  secret_key: "minioadmin"
  bucket: "vokzal-documents"
  use_ssl: false
```

## Запуск

```bash
go mod download
go run cmd/main.go
```

## MinIO интеграция

Требуется MinIO для хранения PDF:
```bash
docker run -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  minio/minio server /data --console-address ":9001"
```

Создать bucket:
```bash
mc alias set vokzal http://localhost:9000 minioadmin minioadmin
mc mb vokzal/vokzal-documents
mc anonymous set download vokzal/vokzal-documents
```

## Зависимости

- Go 1.23+
- PostgreSQL 15+
- MinIO (S3-compatible storage)
- gofpdf v1.16+ (PDF генерация)
- go-qrcode v0.0.0+ (QR коды)

## Примеры использования

### Генерация билета

```bash
curl -X POST http://localhost:8089/v1/document/ticket \
  -H "Content-Type: application/json" \
  -d '{
    "ticket_id": "TK123456",
    "passenger_fio": "Иванов Иван Иванович",
    "passenger_doc": "4500 123456",
    "route": "Ростов — Казань",
    "date": "2026-04-15",
    "time": "08:30",
    "platform": "3",
    "seat": "12",
    "price": 1500.00,
    "qr_code": "https://vokzal.tech/t/TK123456"
  }'
```

Ответ содержит `file_url` для скачивания PDF.

### Генерация ПД-2

```bash
curl -X POST http://localhost:8089/v1/document/pd2 \
  -H "Content-Type: application/json" \
  -d '{
    "number": "123456",
    "series": "АБ",
    "passenger_fio": "Иванов Иван Иванович",
    "passenger_doc": "4500 123456",
    "route_from": "Ростов-на-Дону",
    "route_to": "Казань",
    "date": "2026-04-15",
    "price": 1500.00,
    "issue_date": "2026-04-14",
    "issuer_name": "Петрова А.В.",
    "bus_number": "А123БВ"
  }'
```

## PDF структура

### Электронный билет
- Логотип "Вокзал.ТЕХ"
- Номер билета
- ФИО и документ пассажира
- Маршрут, дата, время
- Перрон, место
- Стоимость
- QR код для контроля

### ПД-2 (проездной документ)
- Серия и номер
- ФИО пассажира
- Откуда/Куда
- Дата отправления
- Номер автобуса
- Стоимость
- Дата выдачи, кассир
- Правовая информация

## Health Check

```bash
GET /health
```

---

© 2025 Вокзал.ТЕХ
