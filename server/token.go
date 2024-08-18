package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/golang-jwt/jwt"
)

const tokenExpirationTime = 24

var (
	privateKey       any
	publicKey        any
	resourceServer   string
	identityProvider string
)

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	CodeVerifier string `json:"code_verifier"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token,omitempty"`
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

// TODO: refactor to config function
func init() {
	privKeyData, err := os.ReadFile("keys/private_key.pem")
	if err != nil {
		panic(fmt.Sprintf("Failed to load private key: %v", err))
	}

	block, _ := pem.Decode(privKeyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		panic("Failed to decode PEM block containing private key")
	}

	privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse private key: %v", err))
	}

	pubKeyData, err := os.ReadFile("keys/public_key.pem")
	if err != nil {
		panic(fmt.Sprintf("Failed to load public key: %v", err))
	}

	block, _ = pem.Decode(pubKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		panic("Failed to decode PEM block containing public key")
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse public key: %v", err))
	}

	publicKey = pubKeyInterface.(*rsa.PublicKey)

	resourceServer = os.Getenv("RESOURCE_SERVER")
	if resourceServer == "" {
		resourceServer = "http://localhost:3000/api"
	}

	identityProvider = os.Getenv("IDENTITY_PROVIDER")
	if identityProvider == "" {
		identityProvider = "http://localhost:8080"
	}
}

func createAccessToken(clientID string, scopes []string, userID int) (string, error) {
	if clientID == "" {
		return "", fmt.Errorf("client id cannot be empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud":   resourceServer,
		"exp":   time.Now().Add(time.Hour * tokenExpirationTime).Unix(),
		"iss":   clientID,
		"sub":   strconv.Itoa(userID),
		"scope": scopes,
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createIDToken(clientID string, user app.User, atHash string) (string, error) {
	if clientID == "" {
		return "", fmt.Errorf("client ID cannot be empty")
	}

	idToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":     identityProvider,
		"sub":     strconv.Itoa(user.UserId),
		"aud":     clientID,
		"exp":     time.Now().Add(time.Hour * tokenExpirationTime).Unix(),
		"iat":     time.Now().Unix(),
		"name":    user.FirstName,
		"email":   user.Email,
		"at_hash": atHash,
	})

	idTokenString, err := idToken.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return idTokenString, nil
}

func (s *Server) validateAccessToken(tokenString string, client app.Client) ([]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
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
	if !containsAnyScope(scope, requiredScopes) {
		return nil, fmt.Errorf("invalid or missing scopes")
	}

	return scope, nil
}

func containsAnyScope(tokenScopes []interface{}, requiredScopes []string) bool {
	scopeSet := make(map[string]struct{})
	for _, s := range tokenScopes {
		if scopeStr, ok := s.(string); ok {
			scopeSet[scopeStr] = struct{}{}
		}
	}

	for _, reqScope := range requiredScopes {
		if _, found := scopeSet[reqScope]; found {
			return true
		}
	}

	return false
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
