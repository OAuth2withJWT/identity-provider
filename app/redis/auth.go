package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	redisClient *redis.Client
}

func NewAuthRepository(redisClient *redis.Client) *AuthRepository {
	return &AuthRepository{
		redisClient: redisClient,
	}
}

func (r *AuthRepository) Create(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.redisClient.Set(ctx, key, value, expiration).Err()
}

func (r *AuthRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *AuthRepository) Delete(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}
