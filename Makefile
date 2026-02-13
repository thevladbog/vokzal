# Вокзал.ТЕХ — Makefile
# Требуется Unix-подобная оболочка (sh/bash). На Windows запускайте make из Git Bash, WSL или MSYS2.

.PHONY: help dev-up dev-down services-build services-start services-stop services-restart services-status ui-dev test test-unit test-services test-ui test-load test-load-smoke lint lint-fix

help:
	@echo "Вокзал.ТЕХ — Makefile команды:"
	@echo ""
	@echo "Инфраструктура:"
	@echo "  make dev-up            - Запустить инфраструктуру (Docker Compose)"
	@echo "  make dev-down          - Остановить инфраструктуру"
	@echo ""
	@echo "Микросервисы:"
	@echo "  make services-build    - Собрать все Go микросервисы"
	@echo "  make services-start    - Запустить все микросервисы в фоне"
	@echo "  make services-stop     - Остановить все микросервисы"
	@echo "  make services-restart  - Перезапустить все микросервисы"
	@echo "  make services-status   - Проверить статус микросервисов"
	@echo ""
	@echo "UI:"
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
	@echo "Линтинг:"
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

# PID файлы хранятся в logs/.pids/
services-start: services-build
	@echo "🚀 Запуск всех микросервисов Вокзал.ТЕХ..."
	@mkdir -p logs/.pids
	@for service in services/*; do \
		if [ -d "$$service" ] && [ -f "$$service/go.mod" ]; then \
			service_name=$$(basename $$service); \
			pid_file="$$(pwd)/logs/.pids/$$service_name.pid"; \
			log_file="$$(pwd)/logs/$$service_name.log"; \
			binary="$$service/bin/service"; \
			if [ ! -f "$$binary" ]; then \
				echo "❌ Бинарник $$binary не найден! Запустите 'make services-build'"; \
				exit 1; \
			fi; \
			if [ -f "$$pid_file" ] && kill -0 $$(cat $$pid_file) 2>/dev/null; then \
				echo "⚠️  $$service_name уже запущен (PID: $$(cat $$pid_file))"; \
			else \
				echo "▶️  Запуск $$service_name..."; \
				(cd $$service && nohup ./bin/service > $$log_file 2>&1 & echo $$! > $$pid_file); \
				sleep 0.5; \
			fi \
		fi \
	done
	@echo "✅ Все микросервисы запущены! Логи: logs/"
	@echo "💡 Используйте 'make services-status' для проверки статуса"

services-stop:
	@echo "🛑 Остановка всех микросервисов Вокзал.ТЕХ..."
	@for service in services/*; do \
		if [ -d "$$service" ] && [ -f "$$service/go.mod" ]; then \
			service_name=$$(basename $$service); \
			pid_file="$$(pwd)/logs/.pids/$$service_name.pid"; \
			if [ -f "$$pid_file" ]; then \
				pid=$$(cat $$pid_file); \
				if kill -0 $$pid 2>/dev/null; then \
					echo "⏹️  Остановка $$service_name (PID: $$pid)..."; \
					kill $$pid 2>/dev/null || true; \
					sleep 0.3; \
					kill -0 $$pid 2>/dev/null && kill -9 $$pid 2>/dev/null || true; \
				else \
					echo "⚠️  $$service_name не запущен"; \
				fi; \
				rm -f "$$pid_file"; \
			else \
				echo "⚠️  PID файл для $$service_name не найден"; \
			fi; \
		fi; \
	done
	@echo "✅ Все микросервисы остановлены!"

services-restart: services-stop
	@sleep 2
	@$(MAKE) services-start

services-status:
	@echo "📊 Статус микросервисов Вокзал.ТЕХ:"
	@echo ""
	@services_count=0; \
	running_count=0; \
	for service in services/*; do \
		if [ -d "$$service" ] && [ -f "$$service/go.mod" ]; then \
			services_count=$$((services_count + 1)); \
			service_name=$$(basename $$service); \
			pid_file="$$(pwd)/logs/.pids/$$service_name.pid"; \
			case "$$service_name" in \
				auth) port=8081;; \
				schedule) port=8082;; \
				ticket) port=8083;; \
				fiscal) port=8084;; \
				payment) port=8085;; \
				board) port=8086;; \
				notify) port=8087;; \
				audit) port=8098;; \
				document) port=8089;; \
				geo) port=8090;; \
				*) port="???";; \
			esac; \
			if [ -f "$$pid_file" ]; then \
				pid=$$(cat $$pid_file); \
				if kill -0 $$pid 2>/dev/null; then \
					running_count=$$((running_count + 1)); \
					if curl -s http://localhost:$$port/health > /dev/null 2>&1; then \
						echo "✅ $$service_name (PID: $$pid, Port: $$port) - HEALTHY"; \
					else \
						echo "⚠️  $$service_name (PID: $$pid, Port: $$port) - RUNNING (не отвечает)"; \
					fi \
				else \
					echo "❌ $$service_name - STOPPED (stale PID)"; \
				fi \
			else \
				echo "❌ $$service_name - STOPPED"; \
			fi \
		fi \
	done; \
	echo ""; \
	echo "📈 Статистика: $$running_count/$$services_count сервисов запущено"

services-run:
	@echo "⚠️  Команда устарела! Используйте:"
	@echo "  make services-start  - для запуска всех сервисов"
	@echo "  make services-stop   - для остановки всех сервисов"
	@echo "  make services-status - для проверки статуса"

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
