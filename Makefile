.PHONY: help dev-up dev-down services-build services-run ui-dev test test-unit test-services test-ui test-e2e test-load test-load-smoke lint

help:
	@echo "Вокзал.ТЕХ — Makefile команды:"
	@echo "  make dev-up            - Запустить инфраструктуру (Docker Compose)"
	@echo "  make dev-down          - Остановить инфраструктуру"
	@echo "  make services-build    - Собрать все Go микросервисы"
	@echo "  make services-run      - Запустить все микросервисы"
	@echo "  make ui-dev            - Запустить UI приложения (dev mode)"
	@echo ""
	@echo "Тестирование:"
	@echo "  make test              - Запустить все тесты (unit + e2e smoke + load smoke)"
	@echo "  make test-unit         - Запустить все unit тесты (Go + JS)"
	@echo "  make test-services     - Запустить unit тесты Go сервисов"
	@echo "  make test-ui           - Запустить unit тесты UI приложений"
	@echo "  make test-e2e          - Запустить E2E тесты (Cypress headless)"
	@echo "  make test-e2e-open     - Открыть Cypress GUI"
	@echo "  make test-load         - Запустить load тесты (k6)"
	@echo "  make test-load-smoke   - Запустить smoke load тест"
	@echo ""
	@echo "  make lint              - Запустить линтеры"

dev-up:
	cd infra/docker && docker-compose up -d

dev-down:
	cd infra/docker && docker-compose down

services-build:
	@for service in services/*; do \
		if [ -d "$$service" ] && [ -f "$$service/go.mod" ]; then \
			echo "Building $$service..."; \
			(cd $$service && go build -o bin/service cmd/main.go) || exit 1; \
		fi \
	done

services-run:
	@echo "Запуск микросервисов..."
	# Будет реализовано после создания сервисов

ui-dev:
	@echo "Запуск UI приложений в dev режиме..."
	# Будет реализовано после создания UI

# Все тесты
test: test-unit test-e2e test-load-smoke
	@echo "✅ Все тесты завершены!"

# Unit тесты
test-unit: test-services test-ui
	@echo "✅ Unit тесты завершены!"

# Go unit тесты
test-services:
	@echo "🧪 Запуск unit тестов Go сервисов..."
	@for service in services/*; do \
		if [ -d "$$service" ]; then \
			echo "Testing $$service..."; \
			cd $$service && go test -v -cover ./... && cd ../..; \
		fi \
	done
	@echo "✅ Go unit тесты завершены!"

# Go unit тесты с покрытием
test-services-coverage:
	@echo "🧪 Запуск unit тестов Go сервисов с покрытием..."
	@for service in services/*; do \
		if [ -d "$$service" ]; then \
			echo "Testing $$service with coverage..."; \
			cd $$service && \
			go test -coverprofile=coverage.out ./... && \
			go tool cover -html=coverage.out -o coverage.html && \
			cd ../..; \
		fi \
	done

# UI unit тесты
test-ui:
	@echo "🧪 Запуск unit тестов UI приложений..."
	@for app in ui/admin-panel ui/pos-app ui/board-display ui/passenger-portal ui/controller-app; do \
		if [ -d "$$app" ] && [ -f "$$app/package.json" ]; then \
			echo "Testing $$app..."; \
			cd $$app && npm test -- --run && cd ../..; \
		fi \
	done
	@echo "✅ UI unit тесты завершены!"

# E2E тесты
test-e2e:
	@echo "🧪 Запуск E2E тестов (Cypress headless)..."
	@cd tests/e2e && npm run cypress:run
	@echo "✅ E2E тесты завершены!"

test-e2e-open:
	@echo "🧪 Открываем Cypress GUI..."
	@cd tests/e2e && npm run cypress:open

test-e2e-chrome:
	@echo "🧪 Запуск E2E тестов в Chrome..."
	@cd tests/e2e && npm run cypress:run:chrome

# Load тесты
test-load: test-load-auth test-load-search
	@echo "✅ Load тесты завершены!"

test-load-smoke:
	@echo "🧪 Запуск smoke load теста..."
	@k6 run --vus 1 --duration 1m tests/load/scenarios/auth.js
	@echo "✅ Smoke load тест завершён!"

test-load-auth:
	@echo "🧪 Запуск load теста авторизации..."
	@k6 run tests/load/scenarios/auth.js

test-load-search:
	@echo "🧪 Запуск load теста поиска рейсов..."
	@k6 run tests/load/scenarios/search-trips.js

test-load-stress:
	@echo "🧪 Запуск stress load теста..."
	@k6 run --stage 1m:0,2m:100,5m:100,2m:200,3m:200,2m:0 tests/load/scenarios/auth.js

lint:
	@echo "Запуск линтеров..."
	@for service in services/*; do \
		echo "Linting $$service..."; \
		cd $$service && golangci-lint run && cd ../..; \
	done
