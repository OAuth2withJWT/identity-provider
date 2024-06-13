package server

import (
	"context"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
)

type contextKey string

const userContextKey contextKey = "user"

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
