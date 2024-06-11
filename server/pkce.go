package server

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func (s *Server) validateCodeChallenge(verifier string, authorizationCode string) (bool, error) {
	authorizationCodeData, err := s.app.AuthorizationCodeService.Get(authorizationCode)
	if err != nil {
		return false, fmt.Errorf("error retrieving authorization code data: %v", err)
	}
	sha256Bytes := sha256.Sum256([]byte(verifier))
	computedChallenge := base64.RawURLEncoding.EncodeToString(sha256Bytes[:])
	return computedChallenge == authorizationCodeData.CodeChallenge, nil
}
