# Notify Service

Микросервис для отправки уведомлений через различные каналы в системе Вокзал.ТЕХ.

## Функционал

### Каналы уведомлений
- **SMS** — SMS.ru API
- **Email** — SMTP (Yandex/Gmail)
- **Telegram** — Telegram Bot API
- **TTS** — голосовые оповещения через локальный агент

### Возможности
- История всех уведомлений
- Статус доставки
- Автоматические повторы при ошибках
- Поддержка шаблонов

## API Endpoints

```bash
# Отправить SMS
POST /v1/notify/sms
{
  "phone": "+79001234567",
  "message": "Ваш билет: https://vokzal.tech/t/TK123456"
}

# Отправить Email
POST /v1/notify/email
{
  "to": "user@example.com",
  "subject": "Билет на автобус",
  "body": "<html><body>Ваш билет...</body></html>"
}

# Отправить Telegram сообщение
POST /v1/notify/telegram
{
  "chat_id": 123456789,
  "message": "Ваш рейс задерживается на 15 минут"
}

# Голосовое оповещение
POST /v1/notify/tts
{
  "text": "Посадка на рейс в Казань начнётся у перрона 3",
  "language": "ru",
  "priority": "high"
}

# Получить уведомление
GET /v1/notify/:id

# Список уведомлений
GET /v1/notify/list?limit=50
```

## Интеграции

### SMS.ru API

```bash
GET https://sms.ru/sms/send?api_id=xxx&to=79001234567&msg=text&json=1
```

Параметры:
- `api_id` — API ключ
- `to` — номер телефона (без +)
- `msg` — текст сообщения
- `json=1` — ответ в JSON

### SMTP Email

Поддержка:
- Yandex (smtp.yandex.ru:587)
- Gmail (smtp.gmail.com:587)
- Mailgun, SendGrid

Требуется:
- SMTP host, port
- Username, password
- TLS/STARTTLS

### Telegram Bot API

Библиотека: `gopkg.in/telebot.v3`

```go
bot, _ := tele.NewBot(tele.Settings{
  Token: "YOUR_BOT_TOKEN",
})

bot.Send(&tele.Chat{ID: 123}, "Message")
```

### TTS через локальный агент

```bash
POST http://localhost:8081/tts/announce
{
  "text": "Посадка началась",
  "language": "ru",
  "priority": "high"
}
```

## Конфигурация

```yaml
server:
  port: "8087"
  mode: "debug"

sms:
  api_id: "YOUR_SMS_RU_API_ID"
  url: "https://sms.ru/sms/send"

email:
  smtp_host: "smtp.yandex.ru"
  smtp_port: 587
  username: "noreply@vokzal.tech"
  password: "YOUR_EMAIL_PASSWORD"
  from: "Вокзал.ТЕХ <noreply@vokzal.tech>"

telegram:
  bot_token: "YOUR_BOT_TOKEN"
  webhook_url: "https://api.vokzal.tech/v1/notify/telegram/webhook"

local_agent:
  url: "http://localhost:8081"
```

## Запуск

### Локально

```bash
go mod download
go run cmd/main.go
```

### Docker

```bash
docker build -t vokzal/notify-service:latest .

docker run -p 8087:8087 \
  -e VOKZAL_NOTIFY_SMS_API_ID=your_api_id \
  -e VOKZAL_NOTIFY_TELEGRAM_BOT_TOKEN=your_token \
  vokzal/notify-service:latest
```

## Зависимости

- Go 1.23+
- PostgreSQL 15+
- NATS 2.10+
- Telebot v3.3+

## Структура БД

### notifications
- `id` (UUID PK)
- `type` (VARCHAR: sms, email, telegram, tts)
- `recipient` (VARCHAR)
- `message` (TEXT)
- `subject` (VARCHAR, для email)
- `status` (VARCHAR: pending, sent, failed)
- `sent_at` (TIMESTAMP)
- `error_msg` (TEXT)
- `metadata` (JSONB)

## Примеры использования

### Отправка билета по SMS

```bash
curl -X POST http://localhost:8087/v1/notify/sms \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+79001234567",
    "message": "Билет Ростов-Казань, 15.04.2026 08:30, место 12. QR: vokzal.tech/t/TK123456"
  }'
```

### Уведомление о задержке по Email

```bash
curl -X POST http://localhost:8087/v1/notify/email \
  -H "Content-Type: application/json" \
  -d '{
    "to": "user@example.com",
    "subject": "Задержка рейса",
    "body": "<p>Ваш рейс в Казань задерживается на 15 минут.</p>"
  }'
```

### Голосовое оповещение на вокзале

```bash
curl -X POST http://localhost:8087/v1/notify/tts \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Внимание! Посадка на рейс в Казань начинается у перрона номер три",
    "language": "ru",
    "priority": "high"
  }'
```

## Health Check

```bash
GET /health
```

---

© 2025 Вокзал.ТЕХ
