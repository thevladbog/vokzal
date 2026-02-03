# Auth Service

Сервис аутентификации и авторизации для Вокзал.ТЕХ.

## Функционал

- ✅ Аутентификация пользователей (логин/пароль)
- ✅ JWT токены (access + refresh)
- ✅ RBAC (Role-Based Access Control)
- ✅ Управление сессиями
- ✅ Graceful shutdown

## Роли

- `cashier` — кассир (продажа/возврат билетов)
- `dispatcher` — диспетчер (управление рейсами, табло)
- `controller` — контролёр (фиксация посадки)
- `admin` — администратор (полный доступ)

## API Endpoints

### POST /v1/auth/login
Вход в систему

**Request:**
```json
{
  "username": "admin",
  "password": "admin123",
  "station_id": "rostov"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "expires_in": 86400,
    "user": {
      "id": "...",
      "username": "admin",
      "full_name": "Системный администратор",
      "role": "admin",
      "station_id": null
    }
  }
}
```

### POST /v1/auth/refresh
Обновление токена

**Request:**
```json
{
  "refresh_token": "eyJhbGci..."
}
```

### POST /v1/auth/logout
Выход

**Headers:**
```
X-Refresh-Token: eyJhbGci...
```

### GET /v1/auth/me
Информация о текущем пользователе

**Headers:**
```
Authorization: Bearer eyJhbGci...
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
docker build -t vokzal/auth-service:latest .

# Запустить
docker run -p 8081:8081 \
  -e VOKZAL_AUTH_DATABASE_HOST=postgres \
  vokzal/auth-service:latest
```

## Конфигурация

См. `config.yaml` для настроек или используй переменные окружения:

- `VOKZAL_AUTH_SERVER_PORT` — порт сервера
- `VOKZAL_AUTH_DATABASE_HOST` — хост БД
- `VOKZAL_AUTH_JWT_SECRET` — секрет JWT

## Тестирование

```bash
# Unit тесты
go test ./...

# С покрытием
go test -cover ./...
```

---

© 2025 Вокзал.ТЕХ
