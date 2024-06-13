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

func getSessionIDFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	sessionID := cookie.Value
	return sessionID
}

func setAuthCookie(w http.ResponseWriter, authID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_id",
		Value:    authID,
		Expires:  time.Now().Add(app.AuthExpiration),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func deleteAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_id",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func getAuthIDFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("auth_id")
	if err != nil {
		return ""
	}
	authSessionID := cookie.Value
	return authSessionID
}
