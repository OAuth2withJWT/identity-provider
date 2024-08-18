package server

import (
	"fmt"
)

func (s *Server) validateAuthorizationCode(code string, clientID string, redirectUri string) ([]string, int, error) {
	authorizationCodeData, err := s.app.AuthorizationCodeService.Get(code)
	if err != nil {
		return []string{}, 0, fmt.Errorf("invalid authorization code")
	}
	if clientID != authorizationCodeData.ClientID || redirectUri != authorizationCodeData.RedirectURI {
		return []string{}, 0, fmt.Errorf("invalid authorization code: client ID or redirect URI mismatch")
	}
	return authorizationCodeData.Scopes, authorizationCodeData.UserId, nil
}
