# Вокзал.ТЕХ — Быстрый старт

## 🚀 Начало работы

### 1. Запуск инфраструктуры

```bash
cd infra/docker
docker-compose up -d

# Проверить статус
docker-compose ps
```

### 2. Применение миграций

```bash
# Установить golang-migrate (если нужно)
brew install golang-migrate

# Применить миграции
migrate -path infra/migrations \
        -database "postgresql://admin:vokzal_secret_2026@localhost:5432/vokzal?sslmode=disable" \
        up
```

### 3. Запуск Auth Service

```bash
cd services/auth
go mod download
go run cmd/main.go
```

Сервис запустится на порту 8081.

### 4. Тестирование Auth API

```bash
# Login
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# Получить токен и использовать для /me
curl -X GET http://localhost:8081/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 5. Запуск Schedule Service

```bash
cd services/schedule
go mod download
go run cmd/main.go
```

Сервис запустится на порту 8082.

## 📦 Доступ к компонентам

**Базы данных:**
- PostgreSQL: `localhost:5432` (admin/vokzal_secret_2026)
- Redis: `localhost:6379` (password: vokzal_redis_2026)
- NATS: `localhost:4222` (vokzal/nats_secret_2026)

**Мониторинг:**
- Grafana: http://localhost:30001 (admin/grafana_secret_2026)
- Prometheus: http://localhost:9090
- Traefik Dashboard: http://localhost:8080

**Хранилище:**
- MinIO Console: http://localhost:9001 (vokzal/minio_secret_2026)

## 🛠️ Разработка

### Создание нового сервиса

```bash
# Скопировать структуру auth сервиса
cp -r services/auth services/new-service

# Обновить go.mod
cd services/new-service
go mod init github.com/vokzal-tech/new-service
go mod tidy

# Настроить конфигурацию
vim config.yaml
```

### Правила разработки

Все правила находятся в `.cursor/rules/`:
- `vokzal-core.mdc` — общие правила проекта
- `go-microservices.mdc` — стандарты Go кода
- `react-typescript.mdc` — стандарты React
- `database.mdc` — работа с БД

## 📚 Документация

- [План реализации](/.cursor/plans/vokzal.tech_implementation_f15a4aff.plan.md)
- [Прогресс](PROGRESS.md)
- [Архитектура](docs/initial/03.md)
- [API](docs/initial/05.md)
- [База данных](docs/initial/04.md)

## 🔒 Безопасность

⚠️ **Важно:** Пароли в конфигурации только для разработки!

В production:
- Используй `.env` файлы
- Храни секреты в Kubernetes Secrets
- Используй HashiCorp Vault или AWS Secrets Manager

## 🐛 Troubleshooting

### PostgreSQL не подключается
```bash
docker-compose logs postgres
docker-compose restart postgres
```

### Порт занят
Измени порты в `docker-compose.yml`

### Миграции не применяются
```bash
# Проверить версию
migrate -database "..." version

# Откатить и применить заново
migrate -database "..." down
migrate -database "..." up
```

## 💡 Полезные команды

```bash
# Остановить всё
docker-compose down

# Удалить с данными
docker-compose down -v

# Логи всех сервисов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f postgres

# Перезапустить сервис
docker-compose restart redis
```

---

© 2025 Вокзал.ТЕХ
