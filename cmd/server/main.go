package main

import (
	"log"

	"github.com/OAuth2withJWT/identity-provider/server"
)

func main() {
	s := server.New()
	log.Fatal(s.Run())
}
