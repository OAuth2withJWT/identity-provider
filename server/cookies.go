package server

import (
	"net/http"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
)

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

func setAuthSessionCookie(w http.ResponseWriter, authSessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_session_id",
		Value:    authSessionID,
		Expires:  time.Now().Add(authorizationCodeExpiration),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func setRedirectCookie(w http.ResponseWriter, redirectURL string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "redirect",
		Value:    redirectURL,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func deleteSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func deleteRedirectCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "redirect",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func deleteAuthSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_session_id",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		Secure:   true,
		HttpOnly: true,
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

func getAuthSessionIDFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("auth_session_id")
	if err != nil {
		return ""
	}
	authSessionID := cookie.Value
	return authSessionID
}
