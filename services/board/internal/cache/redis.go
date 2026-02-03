package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache обёртка для Redis
type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// SetTrips кэширует список рейсов
func (r *RedisCache) SetTrips(ctx context.Context, date string, trips interface{}) error {
	data, err := json.Marshal(trips)
	if err != nil {
		return fmt.Errorf("failed to marshal trips: %w", err)
	}

	key := fmt.Sprintf("board:trips:%s", date)
	return r.client.Set(ctx, key, data, 60*time.Second).Err()
}

// GetTrips получает кэшированные рейсы
func (r *RedisCache) GetTrips(ctx context.Context, date string) ([]byte, error) {
	key := fmt.Sprintf("board:trips:%s", date)
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return data, err
}

// InvalidateTrips удаляет кэш рейсов
func (r *RedisCache) InvalidateTrips(ctx context.Context, date string) error {
	key := fmt.Sprintf("board:trips:%s", date)
	return r.client.Del(ctx, key).Err()
}
