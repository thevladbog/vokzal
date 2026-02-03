# POS Приложение Вокзал.ТЕХ

Desktop приложение на Tauri для кассиров автовокзала.

## Возможности

### Для кассира
- **Продажа билетов** — выбор рейса, ввод данных пассажира, оплата
- **Возврат билетов** — поиск по ID/штрихкоду, расчёт штрафа
- **Печать билетов** — через локальный агент (термопринтер)
- **Печать чеков** — фискализация через ККТ
- **Экран покупателя** — отображение суммы и деталей заказа

### Интеграция
- **Backend API** — все микросервисы через Traefik API Gateway
- **Локальный агент** — печать билетов, чеков, работа с ККТ
- **Real-time** — обновление списка рейсов

## Технологии

### Frontend
- React 18.3
- TypeScript 5.6
- Vite 5.4
- Fluent UI React v9
- TanStack Query
- Zustand
- Tauri API

### Backend (Rust)
- Tauri 1.5
- Reqwest (HTTP клиент)
- Serde (JSON)
- Tokio (async runtime)

## Установка

### Предварительные требования

1. Node.js 18+
2. Rust 1.70+
3. Tauri CLI

```bash
# Установка Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Установка Tauri CLI
cargo install tauri-cli
```

### Установка зависимостей

```bash
cd ui/pos-app
npm install
```

## Разработка

```bash
npm run tauri:dev
```

Запустится:
- Vite dev server на `http://localhost:1420`
- Tauri desktop окно с hot reload

## Сборка

```bash
npm run tauri:build
```

Бинарные файлы будут в `src-tauri/target/release/`.

## Структура

```
pos-app/
├── src/                    # React frontend
│   ├── components/        # UI компоненты
│   ├── pages/            # Страницы
│   │   ├── LoginPage.tsx
│   │   ├── SalePage.tsx
│   │   └── RefundPage.tsx
│   ├── services/         # API клиенты
│   │   ├── api.ts       # HTTP клиент
│   │   └── pos.ts       # Tauri commands
│   ├── stores/          # Zustand stores
│   ├── types/           # TypeScript типы
│   ├── App.tsx
│   └── main.tsx
├── src-tauri/             # Rust backend
│   ├── src/
│   │   └── main.rs      # Tauri commands
│   ├── Cargo.toml
│   └── tauri.conf.json  # Конфигурация
├── customer-display.html  # Экран покупателя
├── package.json
└── README.md
```

## Tauri Commands

### `sell_ticket`

Продажа билета через backend API.

```typescript
import { invoke } from '@tauri-apps/api/tauri';

const ticket = await invoke('sell_ticket', {
  apiUrl: 'http://localhost/api/v1',
  token: 'jwt_token',
  request: {
    trip_id: 'uuid',
    passenger_fio: 'Иванов И.И.',
    passenger_phone: '+79001234567',
  },
});
```

### `return_ticket`

Возврат билета.

```typescript
const ticket = await invoke('return_ticket', {
  apiUrl: 'http://localhost/api/v1',
  token: 'jwt_token',
  ticketId: 'uuid',
});
```

### `print_ticket`

Печать билета через локальный агент.

```typescript
const success = await invoke('print_ticket', {
  agentUrl: 'http://localhost:8081',
  ticketData: {
    ticket_id: 'uuid',
    route: 'Москва-Тверь',
    price: 450.00,
    // ...
  },
});
```

### `print_receipt`

Печать чека на ККТ.

```typescript
const receipt = await invoke('print_receipt', {
  agentUrl: 'http://localhost:8081',
  receiptData: {
    operation: 'sell',
    items: [...],
    payment: { type: 'card', amount: 450.00 },
  },
});
```

### `open_customer_display`

Открыть окно экрана покупателя.

```typescript
await invoke('open_customer_display');
```

## Экран покупателя

Отдельное окно для отображения информации покупателю:
- Маршрут
- Цена
- Дополнительная информация

Обновляется через `window.postMessage()`.

## Environment Variables

Создайте `.env`:

```env
VITE_API_URL=http://localhost/api/v1
VITE_AGENT_URL=http://localhost:8081
```

## Роли

Доступ к POS приложению имеют только пользователи с ролями:
- `cashier`
- `admin`

## Горячие клавиши

- `F1` — Продажа билета
- `F2` — Возврат билета
- `F5` — Обновить список рейсов
- `Esc` — Отменить текущую операцию

## Производительность

- Быстрый старт (Tauri использует WebView)
- Малый размер (~10 MB)
- Низкое потребление памяти (~100 MB)

## Безопасность

- JWT аутентификация
- HTTPS для production
- Валидация на клиенте и сервере
- Audit logging всех операций

---

© 2025 Вокзал.ТЕХ
