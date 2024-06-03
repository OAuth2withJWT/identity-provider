package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/golang-jwt/jwt/v5"
)

const tokenExpirationTime = 30 * 24
const resourceServer = "http://localhost:3000/api"

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

type VerificationRequest struct {
	Token string `json:"token"`
}

type VerificationErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type VerificationResponse struct {
	Active string   `json:"active"`
	Scope  []string `json:"scope,omitempty"`
}

func createToken(clientID string, clientSecret string, scopes []string) (string, error) {
	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("client id and client secret cannot be empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":   resourceServer,
		"exp":   time.Now().Add(time.Hour * tokenExpirationTime).Unix(),
		"iss":   clientID,
		"sub":   clientID,
		"scope": scopes,
	})

	tokenString, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Server) verifyToken(tokenString string, client app.Client) ([]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(client.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		scope, err := s.validateClaims(claims, client)
		if err != nil {
			return nil, err
		}
		return scope, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *Server) validateClaims(claims jwt.MapClaims, client app.Client) ([]interface{}, error) {
	if aud, ok := claims["aud"].(string); !ok || aud != resourceServer {
		return nil, fmt.Errorf("invalid audience")
	}

	if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
		return nil, fmt.Errorf("token has expired")
	}

	if iss, ok := claims["iss"].(string); ok {
		if sub, ok := claims["sub"].(string); ok {
			if iss != sub {
				return nil, fmt.Errorf("incompatible issuer and subject")
			}
		} else {
			return nil, fmt.Errorf("invalid subject")
		}
	} else {
		return nil, fmt.Errorf("invalid issuer")
	}

	scope := claims["scope"].([]interface{})

	requiredScopes := strings.Split(client.Scope, ",")
	if !containsAllScopes(scope, requiredScopes) {
		return nil, fmt.Errorf("invalid or missing scopes")
	}

	return scope, nil
}

func containsAllScopes(tokenScopes []interface{}, requiredScopes []string) bool {
	scopeSet := make(map[string]struct{})
	for _, s := range tokenScopes {
		if scopeStr, ok := s.(string); ok {
			scopeSet[scopeStr] = struct{}{}
		}
	}

	for _, reqScope := range requiredScopes {
		if _, found := scopeSet[reqScope]; !found {
			return false
		}
	}

	return true
}

func extractClientID(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if clientID, ok := claims["iss"].(string); ok {
			return clientID, nil
		}
	}

	return "", fmt.Errorf("client ID not found in token")
}
