# Board Display — Табло отправлений

Веб-приложение для отображения информации о рейсах на мониторах станции в реальном времени.

## Возможности

### Общее табло (/)
- Отображение всех рейсов
- Real-time обновления через WebSocket
- Информация: время, маршрут, направление, перрон, статус, места
- Индикация статусов: посадка, задержка, отменён, отправился
- Часы и дата

### Перронное табло (/platform?platform=1&name=Перрон 1)
- Отображение рейсов только с конкретного перрона
- Увеличенный шрифт для лучшей читаемости
- Real-time обновления
- Индикация статуса "ПОСАДКА" с анимацией

### Особенности
- Автоматическое переподключение WebSocket
- Индикатор соединения
- Автообновление каждую минуту
- Адаптивная вёрстка
- Красивые анимации для статусов

## Технологии

- React 18.3
- TypeScript 5.6
- Vite 5.4
- Fluent UI React v9
- WebSocket (native)
- Zustand (state management)
- Axios
- date-fns

## Установка

```bash
cd ui/board-display
npm install
```

## Разработка

```bash
npm run dev
```

Откроется на `http://localhost:3002`.

## Сборка

```bash
npm run build
```

## Структура

```
board-display/
├── src/
│   ├── components/
│   │   └── TripRow.tsx        # Строка с информацией о рейсе
│   ├── pages/
│   │   ├── PublicBoardPage.tsx    # Общее табло
│   │   └── PlatformBoardPage.tsx  # Перронное табло
│   ├── services/
│   │   └── api.ts             # HTTP клиент
│   ├── hooks/
│   │   └── useWebSocket.ts    # WebSocket hook
│   ├── stores/
│   │   └── boardStore.ts      # Zustand store
│   ├── types/
│   │   └── index.ts           # TypeScript типы
│   ├── App.tsx
│   └── main.tsx
├── package.json
└── README.md
```

## API Endpoints

### HTTP (REST)

```
GET  /api/v1/board/public           # Получить все рейсы
GET  /api/v1/board/platform/:id     # Рейсы с конкретного перрона
GET  /api/v1/board/stats            # Статистика
```

### WebSocket

```
WS   /api/v1/board/ws               # Real-time обновления
```

#### WebSocket сообщения

```typescript
// Обновление рейса
{
  type: 'trip_update',
  data: {
    id: 'uuid',
    route_name: 'Москва-Тверь',
    departure_datetime: '2025-02-15T14:30:00Z',
    status: 'boarding',
    platform: '2',
    // ...
  }
}

// Создание рейса
{
  type: 'trip_created',
  data: { /* ... */ }
}

// Изменение статуса
{
  type: 'status_changed',
  data: { /* ... */ }
}
```

## URL Parameters

### Общее табло

```
http://localhost:3002/
```

### Перронное табло

```
http://localhost:3002/platform?platform=1&name=Перрон%201
http://localhost:3002/platform?platform=2&name=Перрон%202
```

Parameters:
- `platform` — ID перрона
- `name` — Название перрона для отображения

## Статусы рейсов

- **scheduled** — По расписанию (серый)
- **boarding** — ПОСАДКА (зелёный, мигающий)
- **departed** — Отправился (серый, зачёркнутый)
- **delayed** — Задержка (жёлтый)
- **cancelled** — ОТМЕНЁН (красный)

## Deployment

### На TV/монитор

1. Соберите приложение:
```bash
npm run build
```

2. Разверните `dist/` на веб-сервере

3. Откройте в браузере на TV:
```
http://your-server/
http://your-server/platform?platform=1&name=Перрон%201
```

4. Включите полноэкранный режим (F11)

### Рекомендации

- Используйте Chromium/Chrome в kiosk mode
- Отключите спящий режим монитора
- Настройте автозапуск браузера при загрузке системы
- Используйте Raspberry Pi для дешёвого решения

### Kiosk mode (Linux)

```bash
chromium-browser --kiosk --noerrdialogs --disable-infobars \
  --disable-session-crashed-bubble \
  http://localhost:3002/
```

## Environment Variables

Создайте `.env`:

```env
VITE_API_URL=http://localhost/api/v1
VITE_WS_URL=ws://localhost/api/v1/board/ws
```

## Брендинг

Приложение использует фирменные цвета Вокзал.ТЕХ:
- Основной: синий `#0078D4`
- Посадка: зелёный
- Задержка: жёлтый
- Отменён: красный

## Интеграция

Табло интегрируется с:
- **Board Service** — HTTP API и WebSocket
- **Schedule Service** — данные о рейсах
- **NATS** — события в реальном времени

---

© 2025 Вокзал.ТЕХ
