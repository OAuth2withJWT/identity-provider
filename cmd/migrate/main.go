package main

import (
	"log"

	"github.com/OAuth2withJWT/identity-provider/db"
)

func main() {
	d, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer d.Close()

	err = db.Setup(d)
	if err != nil {
		log.Fatal("Failed to initialize the database: ", err)
	}
}
