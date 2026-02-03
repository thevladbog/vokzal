# Вокзал.ТЕХ

**Система автоматизации автовокзалов**

[![Services CI](https://github.com/thevladbog/vokzal/actions/workflows/services-ci.yml/badge.svg)](https://github.com/vokzal-tech/vokzal/actions/workflows/services-ci.yml)
[![UI CI](https://github.com/thevladbog/vokzal/actions/workflows/ui-ci.yml/badge.svg)](https://github.com/vokzal-tech/vokzal/actions/workflows/ui-ci.yml)
[![E2E Tests](https://github.com/thevladbog/vokzal/actions/workflows/e2e-tests.yml/badge.svg)](https://github.com/vokzal-tech/vokzal/actions/workflows/e2e-tests.yml)
[![codecov](https://codecov.io/gh/thevladbog/vokzal/branch/main/graph/badge.svg)](https://codecov.io/gh/vokzal-tech/vokzal)
[![Go Report Card](https://goreportcard.com/badge/github.com/thevladbog/vokzal)](https://goreportcard.com/report/github.com/thevladbog/vokzal)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Описание

Вокзал.ТЕХ — комплексная система управления автовокзалами с микросервисной архитектурой, включающая продажу билетов, управление расписанием, фискализацию, интеграции с платёжными системами и многое другое.

## Структура проекта

```
vokzal-tech/
├── services/          # Go микросервисы
│   ├── auth/         # Аутентификация и авторизация
│   ├── schedule/     # Управление расписанием
│   ├── ticket/       # Продажа и возврат билетов
│   ├── fiscal/       # Фискализация (ККТ)
│   ├── payment/      # Платёжные интеграции
│   ├── board/        # Управление табло
│   ├── geo/          # Географические данные
│   ├── notify/       # Уведомления (SMS, Email, Telegram)
│   ├── audit/        # Логирование операций
│   └── document/     # Генерация документов
├── ui/                  # React приложения
│   ├── admin-panel/     # Панель администратора
│   ├── passenger-portal/# Портал пассажира
│   ├── pos-app/         # POS-приложение (Tauri)
│   ├── board-display/   # Табло отправлений
│   ├── controller-app/  # Приложение контролёра
│   └── shared/          # Общие ресурсы (логотип, brand-colors)
├── agent/               # Локальный агент
│   └── local-agent/     # Работа с ККТ/принтерами
├── infra/
│   ├── docker/          # Docker Compose
│   └── migrations/      # SQL миграции
├── docs/             # Документация
└── shared/           # Общие библиотеки
    ├── go-common/    # Общий Go код
    └── ts-common/    # Общий TS код
```

## Технологии

### Backend
- **Go 1.22+** — микросервисы
- **PostgreSQL 16** — основная БД
- **Redis 7.2** — кэширование
- **NATS 2.10+** — очереди событий
- **MinIO** — хранение документов

### Frontend
- **React 18.3** — UI библиотека
- **TypeScript 5.3+** — типизация
- **Vite** — сборка
- **Tauri** — десктопные приложения
- **Tailwind CSS + Fluent UI** — стили

### Инфраструктура
- **Docker** — контейнеризация
- **Traefik** — API Gateway
- **Prometheus + Grafana** — мониторинг

## Быстрый старт

### Требования

- Docker 24+
- Go 1.22+
- Node.js 20 LTS
- Make

### Локальная разработка

```bash
# Клонировать репозиторий
git clone https://github.com/your-org/vokzal-tech.git
cd vokzal-tech

# Запустить инфраструктуру
cd infra/docker
docker-compose up -d

# Запустить сервис (пример: auth)
cd services/auth
go run cmd/main.go

# Запустить UI (пример: admin)
cd ui/admin
npm install
npm run dev
```

## Документация

- [План реализации](/.cursor/plans/vokzal.tech_implementation_f15a4aff.plan.md)
- [Архитектура](docs/initial/03.md)
- [API](docs/initial/05.md)
- [Модель данных](docs/initial/04.md)

## Соответствие законодательству

- **152-ФЗ** — защита персональных данных
- **54-ФЗ** — применение контрольно-кассовой техники

## Лицензия

Proprietary — все права защищены

## Контакты

- **Email:** dev@vokzal.tech
- **Документация:** https://docs.vokzal.tech

---

© 2025 Вокзал.ТЕХ · "Вокзалик знает, когда уезжать!"
