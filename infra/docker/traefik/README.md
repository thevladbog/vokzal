# Traefik Configuration for Вокзал.ТЕХ

Конфигурация Traefik v3.0 как API Gateway для всех микросервисов.

## Структура

```
traefik/
├── traefik.toml              # Статическая конфигурация
└── config/
    └── dynamic/
        └── services.toml      # Динамическая конфигурация (роутеры, сервисы)
```

## Роутинг

### API Gateway Entry Point

Все микросервисы доступны через `api.vokzal.tech`:

```
http://api.vokzal.tech/v1/auth/*       → Auth Service (8081)
http://api.vokzal.tech/v1/schedule/*   → Schedule Service (8082)
http://api.vokzal.tech/v1/tickets/*    → Ticket Service (8083)
http://api.vokzal.tech/v1/receipts/*   → Fiscal Service (8084)
http://api.vokzal.tech/v1/payment/*    → Payment Service (8085)
http://api.vokzal.tech/v1/board/*      → Board Service (8086)
http://api.vokzal.tech/v1/notify/*     → Notify Service (8087)
http://api.vokzal.tech/v1/audit/*      → Audit Service (8088)
http://api.vokzal.tech/v1/document/*   → Document Service (8089)
http://api.vokzal.tech/v1/geo/*        → Geo Service (8090)
```

### Локальная разработка

Для локальной разработки используйте `localhost`:

```
http://localhost/v1/auth/login
http://localhost/v1/tickets/sell
http://localhost/v1/board/ws (WebSocket)
```

## Middlewares

### CORS Headers
- Разрешены origins: `localhost:5173`, `localhost:3000`, `*.vokzal.tech`
- Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
- Credentials: true

### Rate Limiting
- Average: 100 req/min
- Burst: 50 requests
- Period: 1 минута

## SSL/TLS (Production)

### Let's Encrypt
- Автоматическое получение сертификатов
- Email: admin@vokzal.tech
- HTTP Challenge через entrypoint `web`

### Редирект HTTP → HTTPS

Для production добавить в `services.toml`:

```toml
[http.routers.http-redirect]
  rule = "HostRegexp(`{host:.+}`)"
  entryPoints = ["web"]
  middlewares = ["redirect-to-https"]
  service = "noop@internal"

[http.middlewares.redirect-to-https.redirectScheme]
  scheme = "https"
  permanent = true
```

## Dashboard

Traefik Dashboard доступен на:

```
http://localhost:8080/dashboard/
```

Показывает:
- Роутеры и сервисы
- Middlewares
- Метрики запросов
- Health checks

## Prometheus Metrics

Метрики доступны на:

```
http://localhost:8082/metrics
```

Экспортируемые метрики:
- `traefik_entrypoint_requests_total`
- `traefik_entrypoint_request_duration_seconds`
- `traefik_service_requests_total`

## Запуск

### С Docker Compose

```bash
cd infra/docker
docker-compose up -d traefik
```

### Логи

```bash
docker logs -f vokzal-traefik
```

## Добавление нового сервиса

1. Создать новый сервис в `docker-compose.yml`
2. Добавить роутер в `traefik/config/dynamic/services.toml`:

```toml
[http.routers.new-service]
  rule = "Host(`api.vokzal.tech`) && PathPrefix(`/v1/new`)"
  service = "new-service"
  middlewares = ["cors-headers"]
  entryPoints = ["web", "websecure"]

[http.services.new-service.loadBalancer]
  [[http.services.new-service.loadBalancer.servers]]
    url = "http://new-service:8091"
```

3. Перезапустить Traefik (или подождать автообновления)

## Health Checks

Traefik автоматически проверяет health каждого сервиса:

```
GET http://service:port/health
```

Если сервис недоступен, Traefik:
- Исключает его из балансировки
- Показывает 503 Service Unavailable

## Производительность

- Load balancing между репликами
- Keep-Alive соединения
- HTTP/2 поддержка
- Сжатие ответов (gzip, brotli)

---

© 2025 Вокзал.ТЕХ
