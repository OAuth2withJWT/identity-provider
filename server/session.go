package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
)

type contextKey string

const userContextKey contextKey = "user"

func generateSessionID() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

type AuthSessionData struct {
	ClientID            string
	RedirectURI         string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

var authSessionStore = make(map[string]AuthSessionData)

func (s *Server) withUser(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := getSessionIDFromCookie(r)
		session, err := s.app.SessionService.ValidateSession(sessionID)

		if err == nil {
			user, nil := s.app.UserService.GetUserByID(session.UserId)
			if err != nil {
				http.Error(w, "unauthorized_user", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userContextKey, user)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) protected(next http.Handler) http.Handler {
	return s.withUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(userContextKey).(app.User)

		if !ok {
			originURL := r.URL.RequestURI()
			setRedirectCookie(w, originURL)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	}))
}

func setAuthSession(w http.ResponseWriter, clientID string, redirectURI string, state string, codeChallenge string, codeChallengeMethod string) {
	sessionID, err := generateSessionID()

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	authSessionData := AuthSessionData{
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		State:               state,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}

	authSessionStore[sessionID] = authSessionData

	setAuthCookie(w, sessionID)
}

func getAuthSessionFromStore(authSessionID string) (AuthSessionData, error) {
	session, ok := authSessionStore[authSessionID]
	if !ok {
		return AuthSessionData{}, errors.New("session not found")
	}
	return session, nil
}
