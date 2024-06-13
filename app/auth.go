package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"
)

const (
	AuthExpiration = 10 * time.Minute
	AuthIDLength   = 30
)

type Auth struct {
	AuthID              string
	ClientID            string
	RedirectURI         string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

type AuthService struct {
	repository AuthRepository
}

func NewAuthService(rs AuthRepository) *AuthService {
	return &AuthService{
		repository: rs,
	}
}

type AuthRepository interface {
	Create(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

func (s *AuthService) Create(authInfo Auth) error {
	ctx := context.Background()
	authInfoJSON, err := json.Marshal(authInfo)
	if err != nil {
		return err
	}

	return s.repository.Create(ctx, authInfo.AuthID, string(authInfoJSON), AuthExpiration)
}

func (s *AuthService) Get(key string) (Auth, error) {
	ctx := context.Background()
	res, err := s.repository.Get(ctx, key)
	if err != nil {
		return Auth{}, err
	}

	var authInfo Auth
	err = json.Unmarshal([]byte(res), &authInfo)
	if err != nil {
		return Auth{}, err
	}
	return authInfo, nil
}

func (s *AuthService) Delete(key string) error {
	ctx := context.Background()
	return s.repository.Delete(ctx, key)
}

func GenerateAuthID() (string, error) {
	randomBytes := make([]byte, AuthIDLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
