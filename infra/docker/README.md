# Вокзал.ТЕХ — Docker Compose инфраструктура

## Компоненты

### Базы данных
- **PostgreSQL 16** - основная БД (порт 5432)
- **Redis 7.2** - кэширование (порт 6379)

### Очереди и хранилище
- **NATS 2.10** - очереди событий (порт 4222)
- **MinIO** - объектное хранилище (порты 9000, 9001)

### Инфраструктура
- **Traefik 3.0** - API Gateway (порты 80, 443, 8080)
- **Prometheus** - метрики (порт 9090)
- **Grafana** - визуализация (порт 3000)
- **Loki** - логирование (порт 3100)

## Запуск

```bash
# Запустить все сервисы
docker-compose up -d

# Проверить статус
docker-compose ps

# Просмотр логов
docker-compose logs -f

# Остановить
docker-compose down

# Остановить с удалением volumes
docker-compose down -v
```

## Доступ к сервисам

### Базы данных
- PostgreSQL: `localhost:5432`
  - User: `admin`
  - Password: `vokzal_secret_2026`
  - Database: `vokzal`

- Redis: `localhost:6379`
  - Password: `vokzal_redis_2026`

### Мониторинг
- Grafana: http://localhost:3000
  - User: `admin`
  - Password: `grafana_secret_2026`

- Prometheus: http://localhost:9090
- Traefik Dashboard: http://localhost:8080

### Хранилище
- MinIO Console: http://localhost:9001
  - User: `vokzal`
  - Password: `minio_secret_2026`

### Очереди
- NATS Monitoring: http://localhost:8222

## Переменные окружения

Все пароли и секреты находятся в `docker-compose.yml`. 

⚠️ **Важно:** В production используй `.env` файл и Kubernetes Secrets!

## Volumes

Данные сохраняются в Docker volumes:
- `postgres-data` - данные PostgreSQL
- `redis-data` - данные Redis
- `minio-data` - объекты MinIO
- `prometheus-data` - метрики Prometheus
- `grafana-data` - дашборды Grafana
- `loki-data` - логи Loki

## Healthchecks

Все сервисы имеют healthcheck'и. Проверить:

```bash
docker-compose ps
```

Статус `healthy` означает, что сервис готов к работе.

## Troubleshooting

### PostgreSQL не стартует
```bash
# Проверить логи
docker-compose logs postgres

# Пересоздать контейнер
docker-compose up -d --force-recreate postgres
```

### Порты заняты
Измени порты в `docker-compose.yml` если они конфликтуют с локальными сервисами.

### Недостаточно места
```bash
# Очистить неиспользуемые volumes
docker volume prune

# Очистить всё
docker system prune -a --volumes
```

---

© 2025 Вокзал.ТЕХ
