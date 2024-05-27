package app

import (
	"context"
	"time"
)

type RedisService struct {
	repository RedisRepository
}

func NewRedisService(rs RedisRepository) *RedisService {
	return &RedisService{
		repository: rs,
	}
}

type RedisRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

func (s *RedisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.repository.Set(ctx, key, value, expiration)
}

func (s *RedisService) Get(ctx context.Context, key string) (string, error) {
	return s.repository.Get(ctx, key)
}

func (s *RedisService) Delete(ctx context.Context, key string) error {
	return s.repository.Delete(ctx, key)
}
