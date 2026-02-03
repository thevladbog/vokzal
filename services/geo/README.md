# Geo Service

Микросервис для геокодирования и работы с Yandex Maps API в системе Вокзал.ТЕХ.

## Функционал

### Возможности
- **Геокодирование** — адрес в координаты
- **Обратное геокодирование** — координаты в адрес
- **Расстояние** — вычисление между точками (Haversine)
- **Кэширование** — Redis (TTL 24 часа)

### Yandex Maps API
- Geocoder API
- Точные адреса для России
- Форматированные результаты

## API Endpoints

```bash
# Геокодировать адрес
GET /v1/geo/geocode?address=Ростов-на-Дону, ул. Большая Садовая, 1

# Ответ
{
  "data": {
    "address": "Россия, Ростов-на-Дону, улица Большая Садовая, 1",
    "latitude": 47.222078,
    "longitude": 39.720349
  }
}

# Обратное геокодирование
GET /v1/geo/reverse?lat=47.222078&lon=39.720349

# Ответ
{
  "data": {
    "address": "Россия, Ростов-на-Дону, улица Большая Садовая, 1",
    "latitude": 47.222078,
    "longitude": 39.720349
  }
}

# Расстояние между точками
GET /v1/geo/distance?lat1=47.222078&lon1=39.720349&lat2=55.751244&lon2=37.618423

# Ответ (Ростов — Москва)
{
  "distance_km": 953.24
}
```

## Конфигурация

```yaml
server:
  port: "8090"
  mode: "debug"

redis:
  host: "localhost"
  port: 6379
  password: "vokzal_redis_2026"
  db: 0

yandex_maps:
  api_key: "YOUR_YANDEX_MAPS_API_KEY"
  base_url: "https://geocode-maps.yandex.ru/1.x/"
```

## Запуск

```bash
go mod download
go run cmd/main.go
```

## Yandex Maps API Key

Получить API ключ:
1. https://developer.tech.yandex.ru/
2. Создать проект
3. Включить Geocoder API
4. Скопировать API ключ

Лимиты (бесплатно):
- 25,000 запросов/день
- 5 запросов/секунду

## Кэширование

Redis кэш с TTL 24 часа:
- `geocode:адрес` — координаты
- `reverse:lat,lon` — адрес

## Примеры использования

### Геокодирование остановок

```bash
curl "http://localhost:8090/v1/geo/geocode?address=Ростов-на-Дону, Автовокзал"
```

### Поиск ближайших станций

```bash
# 1. Получить координаты пользователя
curl "http://localhost:8090/v1/geo/reverse?lat=47.222078&lon=39.720349"

# 2. Вычислить расстояния до станций
curl "http://localhost:8090/v1/geo/distance?lat1=47.222078&lon1=39.720349&lat2=47.25&lon2=39.74"
```

### Интеграция с Schedule service

```go
// Получить координаты остановки
resp, _ := http.Get("http://geo:8090/v1/geo/geocode?address=" + url.QueryEscape(stopAddress))

// Сохранить в БД
stop.Latitude = result.Latitude
stop.Longitude = result.Longitude
```

## Зависимости

- Go 1.23+
- Redis 7+
- Yandex Geocoder API

## Формат Yandex Maps API

### Запрос геокодирования
```
GET https://geocode-maps.yandex.ru/1.x/?apikey=xxx&geocode=адрес&format=json&results=1
```

### Ответ
```json
{
  "response": {
    "GeoObjectCollection": {
      "featureMember": [
        {
          "GeoObject": {
            "Point": {
              "pos": "39.720349 47.222078"
            },
            "metaDataProperty": {
              "GeocoderMetaData": {
                "Address": {
                  "formatted": "Россия, Ростов-на-Дону, улица Большая Садовая, 1"
                }
              }
            }
          }
        }
      ]
    }
  }
}
```

## Health Check

```bash
GET /health
```

---

© 2025 Вокзал.ТЕХ
