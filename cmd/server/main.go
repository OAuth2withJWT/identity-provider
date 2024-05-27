package main

import (
	"log"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/OAuth2withJWT/identity-provider/app/postgres"
	"github.com/OAuth2withJWT/identity-provider/app/redis"
	"github.com/OAuth2withJWT/identity-provider/db"
	"github.com/OAuth2withJWT/identity-provider/server"
)

func main() {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer dbConn.Close()

	redisClient := db.ConnectRedis()

	userRepository := postgres.NewUserRepository(dbConn)
	sessionRepository := postgres.NewSessionRepository(dbConn)
	verificationRepository := postgres.NewVerificationRepository(dbConn)
	clientRepository := postgres.NewClientRepository(dbConn)
	redisRepository := redis.NewRedisRepository(redisClient)

	app := app.Application{
		UserService:         app.NewUserService(userRepository),
		SessionService:      app.NewSessionService(sessionRepository),
		VerificationService: app.NewVerificationService(verificationRepository),
		ClientService:       app.NewClientService(clientRepository),
		RedisService:        app.NewRedisService(redisRepository),
	}
	s := server.New(&app)

	log.Fatal(s.Run())
}
