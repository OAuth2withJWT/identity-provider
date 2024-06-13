package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthorizationCodeRepository struct {
	redisClient *redis.Client
}

func NewAuthorizationCodeRepository(redisClient *redis.Client) *AuthorizationCodeRepository {
	return &AuthorizationCodeRepository{
		redisClient: redisClient,
	}
}

func (r *AuthorizationCodeRepository) Create(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.redisClient.Set(ctx, key, value, expiration).Err()
}

func (r *AuthorizationCodeRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *AuthorizationCodeRepository) Delete(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}
