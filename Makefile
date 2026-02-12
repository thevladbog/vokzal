# Вокзал.ТЕХ — Makefile
# Требуется Unix-подобная оболочка (sh/bash). На Windows запускайте make из Git Bash, WSL или MSYS2.

.PHONY: help dev-up dev-down services-build services-run ui-dev test test-unit test-services test-ui test-load test-load-smoke lint lint-fix

help:
	@echo "Вокзал.ТЕХ — Makefile команды:"
	@echo "  make dev-up            - Запустить инфраструктуру (Docker Compose)"
	@echo "  make dev-down          - Остановить инфраструктуру"
	@echo "  make services-build    - Собрать все Go микросервисы"
	@echo "  make services-run      - Запустить все микросервисы"
	@echo "  make ui-dev            - Запустить UI приложения (dev mode)"
	@echo ""
	@echo "Тестирование:"
	@echo "  make test              - Запустить все тесты (unit + load smoke)"
	@echo "  make test-unit         - Запустить все unit тесты (Go + JS)"
	@echo "  make test-services     - Запустить unit тесты Go сервисов"
	@echo "  make test-ui           - Запустить unit тесты UI приложений"
	@echo "  make test-load         - Запустить load тесты (k6)"
	@echo "  make test-load-smoke   - Запустить smoke load тест"
	@echo ""
	@echo "  make lint              - Запустить линтеры"
	@echo "  make lint-fix          - Автофикс выравнивания полей (fieldalignment)"
	@echo ""
	@echo "На Windows запускайте make из Git Bash, WSL или MSYS2."

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
	@echo "Запуск микросервисов: каждый сервис — в отдельном терминале (подробнее: QUICKSTART.md)."
	@echo "  cd services/auth && go run cmd/main.go"
	@echo "  cd services/schedule && go run cmd/main.go"
	@echo "  cd services/ticket && go run cmd/main.go"
	@echo "  ... и т.д. Инфраструктура: make dev-up (infra/docker)."

ui-dev:
	@echo "Запуск UI приложений в dev режиме..."
	# Будет реализовано после создания UI

# Все тесты
test: test-unit test-load-smoke
	@echo "✅ Все тесты завершены!"

# Unit тесты
test-unit: test-services test-ui
	@echo "✅ Unit тесты завершены!"

# Go unit тесты
test-services:
	@echo "🧪 Запуск unit тестов Go сервисов..."
	@for service in services/*; do \
		if [ -d "$$service" ] && [ -f "$$service/go.mod" ]; then \
			echo "Testing $$service..."; \
			(cd $$service && go test -v -cover ./...) || exit 1; \
		fi \
	done
	@echo "✅ Go unit тесты завершены!"

# Go unit тесты с покрытием
test-services-coverage:
	@echo "🧪 Запуск unit тестов Go сервисов с покрытием..."
	@for service in services/*; do \
		if [ -d "$$service" ] && [ -f "$$service/go.mod" ]; then \
			echo "Testing $$service with coverage..."; \
			(cd $$service && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html) || exit 1; \
		fi \
	done

# UI unit тесты
test-ui:
	@echo "🧪 Запуск unit тестов UI приложений..."
	@for app in ui/admin-panel ui/pos-app ui/board-display ui/passenger-portal ui/controller-app; do \
		if [ -d "$$app" ] && [ -f "$$app/package.json" ]; then \
			echo "Testing $$app..."; \
			(cd $$app && npm test -- --run) || exit 1; \
		fi \
	done
	@echo "✅ UI unit тесты завершены!"

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
		if [ -d "$$service" ] && [ -f "$$service/go.mod" ]; then \
			echo "Linting $$service..."; \
			(cd $$service && golangci-lint run) || exit 1; \
		fi \
	done

# Автофикс выравнивания полей структур (govet fieldalignment) в internal/config
lint-fix:
	@echo "Запуск fieldalignment -fix..."
	@for service in services/*; do \
		if [ -d "$$service" ] && [ -f "$$service/go.mod" ] && [ -d "$$service/internal/config" ]; then \
			echo "Fixing $$service/internal/config..."; \
			(cd $$service && go run golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest -fix ./internal/config/) || exit 1; \
		fi \
	done
	@echo "✅ fieldalignment завершён!"
