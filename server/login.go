package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
)

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Context().Value(userContextKey).(app.User); ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	tmpl, err := template.ParseFiles("views/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *Server) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := s.app.UserService.ValidateUserCredentials(email, password)

	if err != nil {
		data := struct {
			Email        string
			Password     string
			ErrorMessage string
		}{Email: email, Password: password, ErrorMessage: err.Error()}

		tmpl, err := template.ParseFiles("views/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	err = s.app.VerificationService.IsUserVerified(user.UserId)

	if err != nil {
		data := struct {
			Email        string
			Password     string
			ErrorMessage string
		}{Email: email, Password: password, ErrorMessage: err.Error()}

		tmpl, err := template.ParseFiles("views/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	sessionID, err := s.app.SessionService.CreateSession(user.UserId, time.Now().Add(app.SessionDurationInHours*time.Hour))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, sessionID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) handleLogoutForm(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionIDFromCookie(r)
	_, err := s.app.SessionService.ValidateSession(sessionID)
	if err == nil {
		err := s.app.SessionService.UpdateStatus(sessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		deleteCookie(w)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
