# Тестирование Вокзал.ТЕХ

Комплексная стратегия тестирования для всех компонентов системы.

## Оглавление

- [Unit тесты (Go)](#unit-тесты-go)
- [Unit тесты (Jest)](#unit-тесты-jest)
- [Load тесты (k6)](#load-тесты-k6)
- [Запуск тестов](#запуск-тестов)
- [CI/CD интеграция](#cicd-интеграция)

## Unit тесты (Go)

### Структура

Каждый микросервис содержит unit-тесты рядом с исходным кодом:

```
services/auth/
├── internal/
│   ├── service/
│   │   ├── auth.go
│   │   └── auth_test.go    # Unit тесты для service layer
│   ├── repository/
│   │   ├── repository.go
│   │   └── repository_test.go
│   └── handlers/
│       ├── auth.go
│       └── auth_test.go
```

### Инструменты

- **testing** — стандартный пакет Go
- **testify/assert** — удобные assertion
- **testify/mock** — мокирование зависимостей
- **sqlmock** — мокирование SQL запросов

### Запуск

```bash
# Все тесты в сервисе
cd services/auth
go test ./...

# С покрытием
go test -cover ./...

# С подробным выводом
go test -v ./...

# Генерация отчёта о покрытии
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Makefile команды

```bash
# Все тесты во всех сервисах
make test-services

# Конкретный сервис
make test-auth
make test-schedule
```

### Примеры

См. файлы:
- `services/auth/internal/service/auth_test.go`
- `services/ticket/internal/repository/repository_test.go`

## Unit тесты (Jest)

### Структура

Тесты для React компонентов и utilities находятся рядом с исходниками:

```
ui/admin-panel/
├── src/
│   ├── components/
│   │   ├── ProtectedRoute.tsx
│   │   └── ProtectedRoute.test.tsx
│   ├── services/
│   │   ├── api.ts
│   │   └── api.test.ts
│   └── utils/
│       ├── format.ts
│       └── format.test.ts
```

### Инструменты

- **Vitest** — быстрый test runner (аналог Jest для Vite)
- **@testing-library/react** — тестирование React компонентов
- **@testing-library/user-event** — симуляция пользовательских действий
- **msw** — мокирование HTTP запросов

### Установка (для каждого UI приложения)

```bash
cd ui/admin-panel
npm install -D vitest @testing-library/react @testing-library/jest-dom @testing-library/user-event jsdom
```

### Запуск

```bash
# Все тесты
cd ui/admin-panel
npm test

# Watch mode
npm run test:watch

# С покрытием
npm run test:coverage
```

### Makefile команды

```bash
# Все UI тесты
make test-ui

# Конкретное приложение
make test-admin-panel
make test-pos-app
```

### Примеры

См. файлы:
- `ui/admin-panel/src/utils/format.test.ts`
- `ui/passenger-portal/src/components/TripCard.test.tsx`

## Load тесты (k6)

### Структура

Load-тесты находятся в `tests/load`:

```
tests/load/
├── scenarios/
│   ├── auth.js                  # Авторизация
│   ├── search-trips.js          # Поиск рейсов
│   ├── ticket-purchase.js       # Покупка билета
│   └── boarding.js              # Посадка
├── utils/
│   └── helpers.js               # Общие утилиты
└── package.json
```

### Инструменты

- **k6** — современный load testing tool

### Установка

```bash
# macOS
brew install k6

# Linux
sudo apt-get install k6

# Windows
choco install k6
```

### Запуск

```bash
cd tests/load

# Один пользователь
k6 run scenarios/auth.js

# 100 виртуальных пользователей, 30 секунд
k6 run --vus 100 --duration 30s scenarios/search-trips.js

# Рампирование: от 0 до 100 за 2 минуты, держим 5 минут, снижаем до 0 за 2 минуты
k6 run --stage 2m:100,5m:100,2m:0 scenarios/ticket-purchase.js

# С выводом метрик в Grafana
K6_PROMETHEUS_REMOTE_WRITE_URL=http://localhost:9090/api/v1/write k6 run scenarios/auth.js
```

### Makefile команды

```bash
# Smoke test (1 VU, 1 минута)
make test-load-smoke

# Load test (50 VUs, 5 минут)
make test-load

# Stress test (100 VUs с рампированием)
make test-load-stress

# Spike test (резкий всплеск нагрузки)
make test-load-spike
```

### Метрики

k6 автоматически собирает:

- **http_req_duration** — время ответа (p95, p99)
- **http_req_failed** — % неудачных запросов
- **iterations** — количество итераций
- **vus** — количество виртуальных пользователей

### Пороговые значения (thresholds)

```javascript
export const options = {
  thresholds: {
    'http_req_duration': ['p(95)<500', 'p(99)<1000'], // 95% < 500ms, 99% < 1s
    'http_req_failed': ['rate<0.01'],                 // <1% ошибок
  },
};
```

См. файлы:
- `tests/load/scenarios/ticket-purchase.js`
- `tests/load/scenarios/search-trips.js`

## Запуск тестов

### Через Makefile

```bash
# Все тесты (Unit + Load smoke)
make test

# Только unit тесты
make test-unit

# Только load
make test-load
```

### Вручную

```bash
# Go unit тесты
cd services/auth && go test ./...

# Frontend unit тесты
cd ui/admin-panel && npm test

# Load тесты
cd tests/load && k6 run scenarios/auth.js
```

## CI/CD интеграция

### GitHub Actions

Workflows в `.github/workflows/` запускают тесты при push и pull request:

1. **Unit тесты (Go)** — на каждый commit
2. **Unit тесты (Jest)** — на каждый commit
3. **Load тесты (k6)** — по расписанию или вручную (workflow `load-tests.yml`)

### Покрытие кода

- Go services: целевое покрытие **80%**
- Frontend apps: целевое покрытие **70%**

Отчёты о покрытии публикуются в PR комментариях через GitHub Actions.

## Лучшие практики

### Go тесты

1. Используйте **Table-driven tests** для множества сценариев
2. Мокируйте внешние зависимости (БД, HTTP, NATS)
3. Именуйте тесты как `TestFunctionName_Scenario_ExpectedBehavior`
4. Используйте `t.Parallel()` для параллельных тестов
5. Не тестируйте внешние библиотеки

### Frontend тесты

1. Тестируйте **поведение**, а не реализацию
2. Используйте `data-testid` для поиска элементов
3. Мокируйте API запросы через MSW
4. Тестируйте accessibility (a11y)
5. Snapshot тесты только для критичных компонентов

### Load тесты

1. Начинайте со **smoke tests** (1 VU)
2. Постепенно увеличивайте нагрузку
3. Определите пороговые значения (thresholds)
4. Мониторьте метрики инфраструктуры (CPU, RAM, DB)
5. Тестируйте на production-like окружении

## Отладка

### Go тесты

```bash
# Запуск конкретного теста
go test -run TestFunctionName

# С логами
go test -v -run TestFunctionName

# С race detector
go test -race ./...
```

### Frontend тесты

```bash
# Debug mode
npm run test:debug

# UI mode (Vitest)
npm run test:ui
```

### Load тесты

```bash
# Подробный вывод
k6 run --verbose scenarios/auth.js

# С логированием в файл
k6 run scenarios/auth.js --log-output=file=test.log
```

## Дополнительно

- Все тесты используют тестовую БД (`vokzal_test`)
- Тестовые пользователи создаются через фикстуры
- Моки для внешних API (Yandex Maps, Tinkoff, SMS.ru)
- Автоматическая очистка тестовых данных

## Поддержка

При возникновении проблем:

1. Проверьте, что все зависимости установлены
2. Убедитесь, что инфраструктура запущена (`make dev-up`)
3. Проверьте логи тестов для деталей ошибок
4. Обратитесь к документации инструментов:
   - [Go testing](https://golang.org/pkg/testing/)
   - [Vitest](https://vitest.dev/)
   - [k6](https://k6.io/docs/)
