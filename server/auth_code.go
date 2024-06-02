package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
)

const authorizationCodeLength = 12
const authorizationCodeExpirationTime = 10

type AuthorizationCodeInfo struct {
	Value       string
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

func (s *Server) storeAuthorizationCode(codeInfo *AuthorizationCodeInfo) error {
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

func (s *Server) validateAuthorizationCode(code string, clientID string, redirectUri string) ([]string, error) {
	ctx := context.Background()
	res, err := s.app.RedisService.Get(ctx, "authorizationCode")
	if err != nil {
		return []string{}, err
	}

	var codeInfo AuthorizationCodeInfo
	err = json.Unmarshal([]byte(res), &codeInfo)
	if err != nil {
		return []string{}, err
	}

	if code != codeInfo.Value || clientID != codeInfo.ClientID || redirectUri != codeInfo.RedirectURI {
		return []string{}, fmt.Errorf("invalid authorization code")
	}

	return codeInfo.Scopes, nil
}
