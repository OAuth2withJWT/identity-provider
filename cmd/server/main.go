package main

import (
	"log"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/OAuth2withJWT/identity-provider/db"
	"github.com/OAuth2withJWT/identity-provider/server"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.Close()

	app := app.Application{
		UserService: app.NewUserService(db),
	}
	s := server.New(&app)

	log.Fatal(s.Run())
}
