# Schedule Service

Микросервис для управления маршрутами, расписаниями и рейсами в системе Вокзал.ТЕХ.

## Функционал

### Маршруты (Routes)
- Создание, чтение, обновление, удаление маршрутов
- JSONB поле `stops` для гибкого хранения промежуточных остановок
- Расчёт расстояния и времени в пути

### Расписания (Schedules)
- Привязка к маршрутам
- Настройка дней недели (JSONB: `[1,2,3,4,5]`)
- Время отправления
- Назначение перронов

### Рейсы (Trips)
- Генерация рейсов из расписания
- Статусы: `scheduled`, `delayed`, `departed`, `arrived`, `cancelled`
- Назначение автобусов и водителей
- Отслеживание задержек

## API Endpoints

### Routes

```bash
# Создать маршрут
POST /v1/routes
{
  "name": "Ростов — Казань",
  "stops": [
    {"station_id": "rostov", "order": 1, "arrival_offset_min": 0},
    {"station_id": "voronezh", "order": 2, "arrival_offset_min": 240},
    {"station_id": "kazan", "order": 3, "arrival_offset_min": 720}
  ],
  "distance_km": 1150.5,
  "duration_min": 720
}

# Список маршрутов
GET /v1/routes?active=true

# Получить маршрут
GET /v1/routes/:id

# Обновить маршрут
PATCH /v1/routes/:id
{
  "name": "Ростов — Казань (обновлённый)",
  "is_active": true
}

# Удалить маршрут
DELETE /v1/routes/:id
```

### Schedules

```bash
# Создать расписание
POST /v1/schedules
{
  "route_id": "uuid",
  "departure_time": "08:30:00",
  "days_of_week": [1, 2, 3, 4, 5],
  "platform": "3"
}

# Список расписаний по маршруту
GET /v1/schedules?route_id=uuid

# Получить расписание
GET /v1/schedules/:id

# Обновить расписание
PATCH /v1/schedules/:id
{
  "platform": "5",
  "is_active": false
}

# Удалить расписание
DELETE /v1/schedules/:id
```

### Trips

```bash
# Создать рейс
POST /v1/trips
{
  "schedule_id": "uuid",
  "date": "2026-04-15",
  "bus_id": "uuid",
  "driver_id": "uuid",
  "platform": "3"
}

# Список рейсов по дате
GET /v1/trips?date=2026-04-15

# Получить рейс
GET /v1/trips/:id

# Обновить статус рейса
PATCH /v1/trips/:id/status
{
  "status": "delayed",
  "delay_minutes": 15
}

# Сгенерировать рейсы из расписания
POST /v1/trips/generate
{
  "schedule_id": "uuid",
  "from_date": "2026-04-01",
  "to_date": "2026-04-30"
}
```

## NATS События

Сервис публикует события:
- `trip.created` — новый рейс создан
- `trip.status_changed` — статус рейса изменён

## Конфигурация

```yaml
server:
  port: "8082"
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
docker build -t vokzal/schedule-service:latest .

# Запустить
docker run -p 8082:8082 \
  -e VOKZAL_SCHEDULE_DATABASE_HOST=postgres \
  -e VOKZAL_SCHEDULE_NATS_URL=nats://nats:4222 \
  vokzal/schedule-service:latest
```

## Зависимости

- Go 1.22+
- PostgreSQL 15+
- NATS 2.10+
- Gin v1.9+
- GORM v1.25+

## Структура БД

### routes
- `id` (UUID PK)
- `name` (VARCHAR)
- `stops` (JSONB)
- `distance_km` (DECIMAL)
- `duration_min` (INTEGER)
- `is_active` (BOOLEAN)

### schedules
- `id` (UUID PK)
- `route_id` (UUID FK)
- `departure_time` (TIME)
- `days_of_week` (JSONB)
- `platform` (VARCHAR)
- `is_active` (BOOLEAN)

### trips
- `id` (UUID PK)
- `schedule_id` (UUID FK)
- `date` (DATE)
- `status` (VARCHAR)
- `delay_minutes` (INTEGER)
- `platform` (VARCHAR)
- `bus_id` (UUID FK)
- `driver_id` (UUID FK)

## Health Check

```bash
GET /health
```

Ответ:
```json
{
  "status": "ok",
  "service": "schedule",
  "version": "1.0.0"
}
```

---

© 2025 Вокзал.ТЕХ
