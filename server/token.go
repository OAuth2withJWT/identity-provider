package server

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenExpirationTimeInHours = 30 * 24

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

func createToken(clientID string, clientSecret string, scopes []string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":    "http://localhost:3000/api",
		"exp":    time.Hour * tokenExpirationTimeInHours,
		"iss":    clientID,
		"sub":    clientID,
		"scopes": scopes,
	})

	tokenString, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
