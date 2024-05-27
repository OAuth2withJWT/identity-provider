package app

import "github.com/redis/go-redis/v9"

type Application struct {
	UserService         *UserService
	SessionService      *SessionService
	VerificationService *VerificationService
	ClientService       *ClientService
	RedisClient         *redis.Client
}
