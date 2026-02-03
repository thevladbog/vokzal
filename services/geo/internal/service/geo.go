// Package service — бизнес-логика Geo Service.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vokzal-tech/geo-service/internal/yandex"
	"go.uber.org/zap"
)

// GeoService — интерфейс сервиса геокодирования.
type GeoService interface {
	Geocode(ctx context.Context, address string) (*yandex.GeocodeResult, error)
	ReverseGeocode(ctx context.Context, lat, lon float64) (*yandex.GeocodeResult, error)
	GetDistance(ctx context.Context, lat1, lon1, lat2, lon2 float64) float64
}

type geoService struct {
	yandexClient *yandex.Client
	redis        *redis.Client
	logger       *zap.Logger
}

// NewGeoService создаёт новый GeoService.
func NewGeoService(yandexClient *yandex.Client, redisClient *redis.Client, logger *zap.Logger) GeoService {
	return &geoService{
		yandexClient: yandexClient,
		redis:        redisClient,
		logger:       logger,
	}
}

func (s *geoService) Geocode(ctx context.Context, address string) (*yandex.GeocodeResult, error) {
	// Проверить кэш
	cacheKey := fmt.Sprintf("geocode:%s", address)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var result yandex.GeocodeResult
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			s.logger.Debug("Geocode cache hit", zap.String("address", address))
			return &result, nil
		}
	}

	// Запросить Yandex Maps
	result, err := s.yandexClient.Geocode(address)
	if err != nil {
		return nil, err
	}

	// Кэшировать результат (TTL 24 часа)
	data, _ := json.Marshal(result)
	s.redis.Set(ctx, cacheKey, data, 24*time.Hour)

	return result, nil
}

func (s *geoService) ReverseGeocode(ctx context.Context, lat, lon float64) (*yandex.GeocodeResult, error) {
	// Проверить кэш
	cacheKey := fmt.Sprintf("reverse:%f,%f", lat, lon)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var result yandex.GeocodeResult
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			s.logger.Debug("Reverse geocode cache hit", zap.Float64("lat", lat), zap.Float64("lon", lon))
			return &result, nil
		}
	}

	// Запросить Yandex Maps
	result, err := s.yandexClient.ReverseGeocode(lat, lon)
	if err != nil {
		return nil, err
	}

	// Кэшировать результат (TTL 24 часа)
	data, _ := json.Marshal(result)
	s.redis.Set(ctx, cacheKey, data, 24*time.Hour)

	return result, nil
}

// GetDistance вычисляет расстояние между двумя точками (формула Haversine).
func (s *geoService) GetDistance(_ context.Context, lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // km

	dLat := (lat2 - lat1) * (3.141592653589793 / 180.0)
	dLon := (lon2 - lon1) * (3.141592653589793 / 180.0)

	lat1Rad := lat1 * (3.141592653589793 / 180.0)
	lat2Rad := lat2 * (3.141592653589793 / 180.0)

	a := (dLat/2)*(dLat/2) + (dLon/2)*(dLon/2)*cosine(lat1Rad)*cosine(lat2Rad)
	c := 2 * atan2(sqrt(a), sqrt(1-a))

	return earthRadius * c
}

func cosine(x float64) float64 {
	// Простая аппроксимация cos
	return 1 - (x*x)/2 + (x*x*x*x)/24
}

func sqrt(x float64) float64 {
	z := 1.0
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

func atan2(y, x float64) float64 {
	// Простая аппроксимация atan2
	if x > 0 {
		return atan(y / x)
	}
	if y > 0 {
		return 1.5707963267948966 - atan(x/y)
	}
	return -1.5707963267948966 - atan(x/y)
}

func atan(x float64) float64 {
	// Аппроксимация atan
	return x - (x*x*x)/3 + (x*x*x*x*x)/5
}
