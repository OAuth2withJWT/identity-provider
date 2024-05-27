package db

import (
	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
