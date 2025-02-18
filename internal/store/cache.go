package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expTime time.Duration) error
	Delete(ctx context.Context, key string) error
}

type RedisCache struct {
	redis *redis.Client
}

func NewRedisCache(redis *redis.Client) *RedisCache {
	return &RedisCache{
		redis: redis,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, expTime time.Duration) error {
	return r.redis.Set(ctx, key, value, 0).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}
