# Vokzal.TECH Local Agent

Локальный агент для работы с кассовым оборудованием и устройствами POS-терминала.

## Возможности

### 1. Фискализация (ККТ)
- Печать чеков продажи и возврата
- Создание Z-отчётов
- Интеграция с АТОЛ драйвером
- Отправка данных в ОФД

### 2. Печать билетов
- Термопечать на ESC/POS принтерах
- Поддержка QR-кодов и штрихкодов
- Печать ПД-2 форм

### 3. Сканирование
- Чтение QR/штрихкодов билетов
- Поддержка HID и Serial сканеров
- Валидация при посадке

### 4. Голосовые оповещения (TTS)
- Объявления отправления рейсов
- Вызов пассажиров на посадку
- Поддержка RHVoice, eSpeak

## Архитектура

```
local-agent/
├── cmd/
│   └── main.go              # Точка входа
├── internal/
│   ├── config/             # Конфигурация
│   ├── handlers/           # HTTP handlers
│   ├── kkt/               # АТОЛ ККТ клиент
│   ├── printer/           # ESC/POS принтер
│   ├── scanner/           # Сканер QR/штрихкодов
│   └── tts/               # TTS (голосовые оповещения)
├── systemd/
│   └── vokzal-agent.service
├── config.yaml
└── go.mod
```

## API Endpoints

### Health & Status
```
GET  /health                # Health check
GET  /status                # Статус всех устройств
```

### ККТ (Фискализация)
```
POST /kkt/receipt           # Печать чека
POST /kkt/z-report          # Создать Z-отчёт
GET  /kkt/status            # Статус ККТ
```

### Принтер
```
POST /printer/ticket        # Печать билета
GET  /printer/status        # Статус принтера
```

### Сканер
```
POST /scanner/scan          # Отсканировать код
GET  /scanner/status        # Статус сканера
```

### TTS
```
POST /tts/announce          # Голосовое оповещение
GET  /tts/status            # Статус TTS
```

## Примеры запросов

### Печать чека на ККТ

```bash
curl -X POST http://localhost:8081/kkt/receipt \
  -H "Content-Type: application/json" \
  -d '{
    "operation": "sell",
    "items": [
      {
        "name": "Билет Москва-Тверь",
        "quantity": 1,
        "price": 450.00,
        "vat": "vat20"
      }
    ],
    "payment": {
      "type": "card",
      "amount": 450.00
    }
  }'
```

### Печать билета

```bash
curl -X POST http://localhost:8081/printer/ticket \
  -H "Content-Type: application/json" \
  -d '{
    "ticket_id": "550e8400-e29b-41d4-a716-446655440000",
    "route": "Москва-Тверь",
    "date": "2025-02-15",
    "time": "14:30",
    "platform": "2",
    "seat": "15",
    "price": 450.00,
    "passenger_fio": "Иванов И.И.",
    "qr_code": "TICKET-550e8400",
    "bar_code": "123456789012"
  }'
```

### Голосовое оповещение

```bash
curl -X POST http://localhost:8081/tts/announce \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Внимание! Отправление автобуса 101 по маршруту Москва-Тверь со второго перрона через 5 минут",
    "language": "ru",
    "priority": "high"
  }'
```

## Установка

### 1. Компиляция

```bash
cd agent/local-agent
go mod tidy
go build -o local-agent cmd/main.go
```

### 2. Настройка конфигурации

Отредактируйте `config.yaml`:

```yaml
kkt:
  device_path: "/dev/ttyUSB0"
  atol_driver: "http://localhost:10001"
  inn: "ВАШ_ИНН"

printer:
  device_path: "/dev/usb/lp0"

scanner:
  device_path: "/dev/input/by-id/usb-scanner-event-kbd"

tts:
  voice: "alena"
```

### 3. Установка как systemd сервис

```bash
sudo cp local-agent /opt/vokzal/agent/
sudo cp config.yaml /opt/vokzal/agent/
sudo cp systemd/vokzal-agent.service /etc/systemd/system/

sudo systemctl daemon-reload
sudo systemctl enable vokzal-agent
sudo systemctl start vokzal-agent
```

### 4. Проверка статуса

```bash
sudo systemctl status vokzal-agent
curl http://localhost:8081/health
```

## Драйвера и зависимости

### АТОЛ ККТ Драйвер

Скачайте и установите драйвер с сайта АТОЛ:

```bash
wget https://atol.ru/drivers/kkt_driver.tar.gz
tar -xzf kkt_driver.tar.gz
cd kkt_driver
sudo ./install.sh
```

### RHVoice (TTS)

```bash
sudo apt install rhvoice rhvoice-russian
```

### ESC/POS принтер

```bash
sudo apt install cups
sudo lpadmin -p thermal_printer -E -v usb://...
```

## Логи

```bash
# Системные логи
sudo journalctl -u vokzal-agent -f

# Application логи
tail -f /var/log/vokzal/agent.log
```

## Безопасность

- Агент работает от имени пользователя `vokzal`
- Доступ к устройствам ограничен через группу `dialout`
- Логи записываются в `/var/log/vokzal/`

## Интеграция с сервисами

Агент используется:

- **Fiscal Service** — для печати чеков (POST `/kkt/receipt`)
- **Ticket Service** — для печати билетов (POST `/printer/ticket`)
- **Board Service** — для голосовых оповещений (POST `/tts/announce`)
- **Ticket Service** — для сканирования при посадке (POST `/scanner/scan`)

---

© 2025 Вокзал.ТЕХ
