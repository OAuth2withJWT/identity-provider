package app

import (
	"context"
	"encoding/json"
	"time"
)

const (
	AuthorizationCodeLength     = 12
	AuthorizationCodeExpiration = 10 * time.Minute
)

type AuthorizationCode struct {
	Value               string
	ClientID            string
	RedirectURI         string
	State               string
	Scopes              []string
	Expiration          int64
	CodeChallenge       string
	CodeChallengeMethod string
	UserId              int
}

type AuthorizationCodeService struct {
	repository AuthorizationCodeRepository
}

func NewAuthorizationCodeService(rs AuthorizationCodeRepository) *AuthorizationCodeService {
	return &AuthorizationCodeService{
		repository: rs,
	}
}

type AuthorizationCodeRepository interface {
	Create(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

func (s *AuthorizationCodeService) Create(authorizationCodeInfo AuthorizationCode) error {
	ctx := context.Background()
	authorizationCodeInfoJSON, err := json.Marshal(authorizationCodeInfo)
	if err != nil {
		return err
	}

	return s.repository.Create(ctx, authorizationCodeInfo.Value, string(authorizationCodeInfoJSON), AuthExpiration)
}

func (s *AuthorizationCodeService) Get(key string) (AuthorizationCode, error) {
	ctx := context.Background()
	res, err := s.repository.Get(ctx, key)
	if err != nil {
		return AuthorizationCode{}, err
	}

	var authorizationCodeInfo AuthorizationCode
	err = json.Unmarshal([]byte(res), &authorizationCodeInfo)
	if err != nil {
		return AuthorizationCode{}, err
	}
	return authorizationCodeInfo, nil
}

func (s *AuthorizationCodeService) Delete(key string) error {
	ctx := context.Background()
	return s.repository.Delete(ctx, key)
}

func GenerateAuthorizationCode() (string, error) {
	authorizationCode, err := GenerateRandomBytes(AuthorizationCodeLength)
	if err != nil {
		return "", err
	}
	return authorizationCode, nil
}
