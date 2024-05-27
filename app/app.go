package app

type Application struct {
	UserService         *UserService
	SessionService      *SessionService
	VerificationService *VerificationService
	ClientService       *ClientService
	RedisService        *RedisService
}
