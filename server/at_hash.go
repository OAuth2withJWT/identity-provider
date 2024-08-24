package server

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func (s *Server) ComputeAtHash(accessToken string) (string, error) {
	if accessToken == "" {
		return "", errors.New("access token is empty")
	}

	hash := sha256.New()
	_, err := hash.Write([]byte(accessToken))
	if err != nil {
		return "", err
	}
	hashBytes := hash.Sum(nil)

	halfHash := hashBytes[:len(hashBytes)/2]
	atHash := base64.RawURLEncoding.EncodeToString(halfHash)
	return atHash, nil
}
