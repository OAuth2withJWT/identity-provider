package main

import (
	"log"

	"github.com/OAuth2withJWT/identity-provider/internal/db"
	"github.com/OAuth2withJWT/identity-provider/server"
)

func main() {
	DB, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.CloseDB()

	s := server.New(DB)
	log.Fatal(s.Run())
}
