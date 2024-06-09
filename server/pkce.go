package server

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

func (s *Server) validateCodeChallenge(verifier string) (bool, error) {
	ctx := context.Background()
	res, err := s.app.RedisService.Get(ctx, "authorizationCode")
	if err != nil {
		return false, err
	}

	var codeInfo AuthorizationCodeInfo
	err = json.Unmarshal([]byte(res), &codeInfo)
	if err != nil {
		return false, err
	}

	sha256Bytes := sha256.Sum256([]byte(verifier))
	computedChallenge := base64.RawURLEncoding.EncodeToString(sha256Bytes[:])
	return computedChallenge == codeInfo.CodeChallenge, nil
}
