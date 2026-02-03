# Вокзал.ТЕХ - Мониторинг и наблюдаемость

Полный стек мониторинга для системы Вокзал.ТЕХ: Prometheus, Grafana, Loki.

## Компоненты

### Prometheus 2.45+
- **Сбор метрик** всех Go микросервисов (через `/metrics`)
- **Scrape interval**: 15 секунд
- **Retention**: 15 дней
- **URL**: http://localhost:9090

### Grafana 10.1+
- **Визуализация** метрик из Prometheus и логов из Loki
- **3 предустановленных дашборда**:
  - Обзор сервисов (Request Rate, Response Time, Errors)
  - База данных (PostgreSQL, Redis)
  - Бизнес метрики (продажи, платежи, посадка)
- **URL**: http://localhost:3000
- **Логин**: admin / admin (по умолчанию)

### Loki 2.9+
- **Централизованное хранение логов** всех сервисов
- **Интеграция с Grafana** для поиска и визуализации
- **URL**: http://localhost:3100

## Запуск

```bash
# Из корня проекта
make dev-up

# Или напрямую
cd infra/docker
docker-compose up -d prometheus grafana loki
```

Проверка статуса:

```bash
docker-compose ps
# Должны работать: prometheus, grafana, loki
```

## Доступ к интерфейсам

### Grafana
1. Откройте http://localhost:3000
2. Войдите (admin / admin)
3. Перейдите в "Dashboards" → выберите один из 3 дашбордов

### Prometheus
1. Откройте http://localhost:9090
2. Перейдите в "Graph" для запросов PromQL
3. Или "Targets" для просмотра состояния targets

### Loki (через Grafana)
1. В Grafana перейдите в "Explore"
2. Выберите datasource "Loki"
3. Введите LogQL запрос, например: `{service="auth"}`

## Дашборды Grafana

### 1. Обзор сервисов (vokzal-services-overview)

**Что показывает:**
- Request Rate (RPS) для всех сервисов
- Success Rate (% 2xx ответов)
- Response Time (p95, p99 перцентили)
- Client Errors (4xx)
- Server Errors (5xx)

**Когда использовать:**
- Проверка производительности API
- Обнаружение аномалий (внезапный рост ошибок)
- Мониторинг SLA (время ответа, доступность)

### 2. База данных (vokzal-database)

**Что показывает:**
- PostgreSQL: активные соединения, транзакции, операции с данными (INSERT/UPDATE/DELETE)
- Redis: соединения, использование памяти, команды, cache hit/miss

**Когда использовать:**
- Оптимизация запросов (высокая нагрузка на БД)
- Проверка кэширования (cache hit rate)
- Мониторинг соединений (connection pooling)

### 3. Бизнес метрики (vokzal-business-metrics)

**Что показывает:**
- Продано билетов (за час, по времени)
- Возвратов билетов
- Успешных платежей (Tinkoff, СБП, наличные)
- Созданных рейсов
- Посадка пассажиров

**Когда использовать:**
- Отчёты для менеджмента
- Анализ пиковых нагрузок (продажи в часы пик)
- Проверка работы платёжных систем

## Настройка алертов

### Создание алерта в Grafana

1. Откройте дашборд
2. Нажмите на панель → "Edit"
3. Перейдите на вкладку "Alert"
4. Настройте условия (например, "5xx errors > 10 за 5 минут")
5. Добавьте notification channel (Email, Telegram, Slack)
6. Сохраните

### Пример: алерт на 5xx ошибки

```yaml
# Условие:
WHEN max() OF query(A, 5m, now) IS ABOVE 10

# Query A:
sum(rate(gin_requests_total{status=~"5.."}[5m]))

# Notification:
- Email: dev@vokzal.tech
- Message: "Критическое количество 5xx ошибок: {{ $value }}"
```

## Метрики микросервисов

Все Go сервисы экспортируют метрики в формате Prometheus на `/metrics`:

### HTTP метрики (Gin)
- `gin_requests_total{service, method, path, status}` - всего запросов
- `gin_request_duration_seconds{service, method, path}` - время ответа (histogram)

### Go Runtime метрики
- `go_goroutines` - количество горутин
- `go_memstats_alloc_bytes` - использование памяти
- `go_gc_duration_seconds` - время GC

### Custom метрики (примеры)
- `ticket_sales_total` - всего продано билетов
- `ticket_returns_total` - всего возвратов
- `payment_amount_total` - сумма платежей
- `trips_active` - активных рейсов

## Логирование

### Формат логов

Все сервисы используют структурированное логирование (Zap):

```json
{
  "level": "info",
  "ts": "2026-02-03T15:30:00.123Z",
  "caller": "handlers/auth.go:45",
  "msg": "user logged in",
  "service": "auth",
  "user_id": "uuid-123",
  "ip": "192.168.1.100"
}
```

### Поиск логов в Loki (через Grafana Explore)

#### Базовые запросы

```logql
# Все логи сервиса auth
{service="auth"}

# Только ошибки (level=error)
{service="auth"} |= "level\":\"error"

# Логи конкретного пользователя
{service="auth"} |= "user_id\":\"uuid-123"

# Логи за последние 15 минут с фильтрацией
{service="ticket"} |= "ticket sold" [15m]
```

#### Продвинутые запросы

```logql
# Количество ошибок по сервисам
sum by (service) (count_over_time({level="error"} [5m]))

# Rate ошибок 5xx (из логов Traefik)
rate({service="traefik"} |= "status=5" [1m])

# Топ-10 медленных запросов
topk(10, 
  sum by (path) (
    rate({service=~".*"} |= "duration" | json | unwrap duration [5m])
  )
)
```

### Интеграция с Sentry (опционально)

Для критических ошибок (panic, 500 errors) рекомендуется использовать Sentry:

```bash
# Добавить в .env сервисов
SENTRY_DSN=https://<key>@sentry.io/<project>
```

Sentry автоматически отправляет:
- Stack traces
- Request context (headers, body)
- Environment info
- User info

## Troubleshooting

### Prometheus не собирает метрики

1. Проверьте targets:
   ```bash
   curl http://localhost:9090/api/v1/targets
   ```
2. Убедитесь, что сервисы экспортируют `/metrics`:
   ```bash
   curl http://localhost:8001/metrics  # auth service
   ```
3. Проверьте `prometheus.yml` (должны быть все targets)

### Grafana не показывает данные

1. Проверьте datasource (Settings → Data Sources):
   - Prometheus: http://prometheus:9090
   - Loki: http://loki:3100
2. Проверьте временной диапазон (top-right corner)
3. Проверьте query (должен быть валидный PromQL/LogQL)

### Loki не получает логи

1. Проверьте docker logs:
   ```bash
   docker-compose logs loki
   ```
2. Убедитесь, что сервисы логируют в stdout/stderr
3. Проверьте `loki/local-config.yaml` (должен быть правильный ingester)

## Best Practices

### 1. Настройка retention

По умолчанию:
- **Prometheus**: 15 дней
- **Loki**: 7 дней

Для production рекомендуется:
- Prometheus: 30 дней (краткосрочные метрики)
- Loki: 14 дней (краткосрочные логи)
- Long-term storage: S3/MinIO (для audit logs, 152-ФЗ: 3-5 лет)

### 2. Оптимизация запросов

- Используйте `rate()` вместо `increase()` для скорости
- Добавляйте labels для фильтрации (service, method, status)
- Избегайте `count()` без `rate()` (результат постоянно растёт)

### 3. Алерты

Критические алерты:
- 5xx errors > 1% от всех запросов (5 минут)
- Response time p99 > 1 секунда (5 минут)
- PostgreSQL connections > 80% от max_connections
- Redis memory > 90%

Предупреждающие алерты:
- 4xx errors > 10% от всех запросов (15 минут)
- Disk usage > 80%

### 4. Безопасность

- Защитите Grafana аутентификацией (смените admin пароль!)
- Ограничьте доступ к Prometheus/Loki (только внутренняя сеть)
- Используйте HTTPS для Grafana (Let's Encrypt через Traefik)

## Kubernetes (будущее)

При миграции на Kubernetes:

1. Используйте **Prometheus Operator** для автоматического обнаружения targets
2. Настройте **Persistent Volumes** для Prometheus/Loki
3. Добавьте **Horizontal Pod Autoscaling** на основе метрик
4. Используйте **Grafana Loki** с S3 storage backend

## Дополнительные ресурсы

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Loki Documentation](https://grafana.com/docs/loki/)
- [PromQL Cheat Sheet](https://promlabs.com/promql-cheat-sheet/)
- [LogQL Documentation](https://grafana.com/docs/loki/latest/logql/)

## Контакты

По вопросам мониторинга: dev@vokzal.tech

---

© 2026 Вокзал.ТЕХ
