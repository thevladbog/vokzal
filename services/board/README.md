# Board Service

Микросервис для управления информационными табло с WebSocket real-time обновлениями в системе Вокзал.ТЕХ.

## Функционал

### Real-time обновления
- WebSocket соединения для live updates
- Автоматическая отправка изменений статуса рейсов
- Поддержка множественных клиентов
- Ping/pong для поддержания соединения

### Табло
- **Общее табло** — все рейсы дня
- **Перронное табло** — рейсы конкретного перрона
- Статистика посадки пассажиров

### Кэширование
- Redis кэш (TTL 60 сек)
- Автоматическая инвалидация при изменениях
- Оптимизация запросов к БД

## API Endpoints

### WebSocket

```bash
# Подключение к WebSocket
ws://localhost:8086/v1/board/ws

# Или с SSL
wss://api.vokzal.tech/v1/board/ws
```

Клиент получает сообщения в формате:

```json
{
  "type": "trip_update",
  "trip_id": "uuid",
  "status": "delayed",
  "delay_minutes": 15,
  "timestamp": "2026-04-15T10:30:00Z"
}
```

Типы сообщений:
- `trip_created` — новый рейс создан
- `trip_update` — статус рейса изменён

### HTTP Endpoints

```bash
# Получить данные для общего табло
GET /v1/board/public?date=2026-04-15

# Ответ
{
  "data": [
    {
      "id": "uuid",
      "date": "2026-04-15",
      "departure_time": "08:30:00",
      "route_name": "Ростов — Казань",
      "platform": "3",
      "status": "scheduled",
      "delay_minutes": 0
    }
  ]
}

# Получить данные для перронного табло
GET /v1/board/platform/3?date=2026-04-15

# Ответ включает статистику посадки
{
  "data": [
    {
      "id": "uuid",
      "route_name": "Ростов — Казань",
      "departure_time": "08:30:00",
      "status": "boarding",
      "total_tickets": 45,
      "boarded_count": 32
    }
  ]
}

# Статистика WebSocket соединений
GET /v1/board/stats
```

## WebSocket клиент (JavaScript)

```javascript
const ws = new WebSocket('ws://localhost:8086/v1/board/ws');

ws.onopen = () => {
  console.log('Connected to board updates');
};

ws.onmessage = (event) => {
  const update = JSON.parse(event.data);
  console.log('Board update:', update);
  
  if (update.type === 'trip_update') {
    updateTripOnScreen(update.trip_id, update.status, update.delay_minutes);
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected, reconnecting...');
  setTimeout(reconnect, 3000);
};
```

## NATS События

### Подписки
- `trip.created` — новый рейс создан
- `trip.status_changed` — статус рейса изменён

При получении события:
1. Инвалидировать Redis кэш
2. Отправить обновление через WebSocket всем клиентам

## Конфигурация

```yaml
server:
  port: "8086"
  mode: "debug"

database:
  host: "localhost"
  port: 5432
  user: "admin"
  password: "vokzal_secret_2026"
  dbname: "vokzal"
  sslmode: "disable"

redis:
  host: "localhost"
  port: 6379
  password: "vokzal_redis_2026"
  db: 0

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
docker build -t vokzal/board-service:latest .

# Запустить
docker run -p 8086:8086 \
  -e VOKZAL_BOARD_DATABASE_HOST=postgres \
  -e VOKZAL_BOARD_REDIS_HOST=redis \
  -e VOKZAL_BOARD_NATS_URL=nats://nats:4222 \
  vokzal/board-service:latest
```

## Зависимости

- Go 1.23+
- PostgreSQL 15+
- Redis 7+
- NATS 2.10+
- Gorilla WebSocket v1.5+

## Архитектура

### WebSocket Hub
- Управление всеми активными соединениями
- Broadcast сообщений всем клиентам
- Автоматический ping/pong
- Graceful disconnect

### Redis Cache
- TTL 60 секунд для данных табло
- Автоматическая инвалидация
- Снижение нагрузки на PostgreSQL

### NATS Integration
- Event-driven обновления
- Подписка на изменения рейсов
- Декаплинг от других сервисов

## UI интеграция

### React компонент

```typescript
import { useEffect, useState } from 'react';

function BoardDisplay() {
  const [trips, setTrips] = useState([]);
  
  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8086/v1/board/ws');
    
    ws.onmessage = (event) => {
      const update = JSON.parse(event.data);
      if (update.type === 'trip_update') {
        setTrips(prev => prev.map(trip => 
          trip.id === update.trip_id 
            ? { ...trip, ...update } 
            : trip
        ));
      }
    };
    
    return () => ws.close();
  }, []);
  
  return (
    <div className="board-display">
      {trips.map(trip => (
        <TripRow key={trip.id} trip={trip} />
      ))}
    </div>
  );
}
```

## Health Check

```bash
GET /health
```

Ответ:
```json
{
  "status": "ok",
  "service": "board",
  "version": "1.0.0"
}
```

## Производительность

- Поддержка 1000+ одновременных WebSocket соединений
- Redis кэш для снижения latency
- Эффективная broadcast модель

---

© 2025 Вокзал.ТЕХ
