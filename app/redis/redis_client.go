package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redisClient *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{
		redisClient: redisClient,
	}
}

func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.redisClient.Set(ctx, key, value, expiration).Err()
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisRepository) Delete(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}
