package server

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
)

type contextKey string

const userContextKey contextKey = "user"

func (s *Server) withUser(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := getSessionIDFromCookie(r)
		session, err := s.app.SessionService.ValidateSession(sessionID)

		if err == nil {
			user, _ := s.app.UserService.GetUserByID(session.UserId)
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
			loginURL := "/login?redirect=" + url.QueryEscape(originURL)
			http.Redirect(w, r, loginURL, http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	}))
}

func setSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(app.SessionDurationInHours * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func deleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func getSessionIDFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	sessionID := cookie.Value
	return sessionID
}
