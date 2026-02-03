# Админ-панель Вокзал.ТЕХ

Веб-приложение для администраторов и диспетчеров автовокзала.

## Возможности

### Для Администратора
- Управление расписаниями и маршрутами
- Управление станциями, автобусами, водителями
- Отчёты: продажи, выручка, заполняемость
- Управление пользователями и ролями

### Для Диспетчера
- Управление рейсами (статусы, платформы)
- Генерация рейсов по расписанию
- Мониторинг текущих отправлений
- Управление табло

### Для Кассира
- Продажа билетов
- Возврат билетов
- Отчёты по кассе

### Для Бухгалтера
- Финансовые отчёты
- Z-отчёты
- Аудит операций

## Технологии

- React 18.3
- TypeScript 5.6
- Vite 5.4
- Fluent UI React v9
- TanStack Query (React Query)
- TanStack Table
- React Router v6
- Zustand (state management)
- React Hook Form + Zod
- Axios
- Recharts (графики)

## Установка

```bash
cd ui/admin-panel
npm install
```

## Разработка

```bash
npm run dev
```

Откроется на `http://localhost:3001`.

API проксируется на `http://localhost:80/api/v1`.

## Сборка

```bash
npm run build
```

## Структура

```
src/
├── components/        # Переиспользуемые компоненты
│   ├── ProtectedRoute.tsx
│   ├── Layout.tsx
│   └── ...
├── pages/            # Страницы приложения
│   ├── LoginPage.tsx
│   ├── DashboardPage.tsx
│   ├── SchedulesPage.tsx
│   ├── TripsPage.tsx
│   ├── TicketsPage.tsx
│   ├── ReportsPage.tsx
│   └── ...
├── services/         # API клиенты
│   ├── api.ts
│   ├── auth.ts
│   ├── schedule.ts
│   ├── ticket.ts
│   └── ...
├── stores/           # Zustand stores
│   ├── authStore.ts
│   └── ...
├── types/            # TypeScript типы
│   └── index.ts
├── utils/            # Утилиты
├── App.tsx           # Корневой компонент
└── main.tsx          # Entry point
```

## Роли и права

- **admin** — полный доступ ко всем функциям
- **dispatcher** — управление рейсами, расписаниями
- **cashier** — продажа и возврат билетов
- **accountant** — отчёты, аудит
- **controller** — проверка билетов

## API Endpoints

```
GET  /api/v1/auth/login
POST /api/v1/auth/logout
POST /api/v1/auth/refresh

GET  /api/v1/schedule/routes
POST /api/v1/schedule/routes
GET  /api/v1/schedule/schedules
POST /api/v1/schedule/schedules
GET  /api/v1/schedule/trips
PUT  /api/v1/schedule/trips/:id

GET  /api/v1/tickets
POST /api/v1/tickets/sell
POST /api/v1/tickets/:id/return
GET  /api/v1/tickets/reports/sales
```

## Environment Variables

Создайте `.env` файл:

```env
VITE_API_URL=http://localhost/api/v1
```

## Брендинг

Приложение использует фирменные цвета Вокзал.ТЕХ:

- Основной: `#0078D4` (синий)
- Дополнительный: `#FFC107` (жёлтый)
- Акцент: `#E91E63` (розовый)

---

© 2025 Вокзал.ТЕХ
