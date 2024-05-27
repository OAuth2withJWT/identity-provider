package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
)

const authorizationCodeLength = 12
const authorizationCodeExpirationTime = 10

type AuthorizationCodeInfo struct {
	ClientID    string
	RedirectURI string
	State       string
	Scopes      []string
	Expiration  int64
}

func generateAuthorizationCode() (string, error) {
	authorizationCode, err := app.GenerateRandomBytes(authorizationCodeLength)
	if err != nil {
		return "", err
	}
	return authorizationCode, nil
}

func storeAuthorizationCode(s *Server, codeInfo *AuthorizationCodeInfo) error {
	ctx := context.Background()
	codeInfoJSON, err := json.Marshal(codeInfo)
	if err != nil {
		return err
	}

	err = s.app.RedisService.Set(ctx, "authorizationCode", string(codeInfoJSON), authorizationCodeExpirationTime*time.Minute)
	if err != nil {
		return err
	}

	return nil
}
