package main

import (
	"log"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/OAuth2withJWT/identity-provider/app/postgres"
	"github.com/OAuth2withJWT/identity-provider/db"
	"github.com/OAuth2withJWT/identity-provider/server"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.Close()

	userRepository := postgres.NewUserRepository(db)
	sessionRepository := postgres.NewSessionRepository(db)
	verificationRepository := postgres.NewVerificationRepository(db)

	app := app.Application{
		UserService:         app.NewUserService(userRepository),
		SessionService:      app.NewSessionService(sessionRepository),
		VerificationService: app.NewVerificationService(verificationRepository),
	}
	s := server.New(&app)

	log.Fatal(s.Run())
}
