package server

import (
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
)

func (s *Server) handleRegistrationPage(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionIDFromCookie(r)
	_, err := s.app.SessionService.ValidateSession(sessionID)

	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	tmpl, _ := template.ParseFiles("views/registration.html")
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRegistrationForm(w http.ResponseWriter, r *http.Request) {
	req := app.CreateUserRequest{
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
		Password:  r.FormValue("password"),
	}

	_, err := s.app.UserService.Create(req)
	if err != nil {
		var errorMessage string
		if fieldErr, ok := err.(*app.FieldError); ok {
			errorMessage = fieldErr.Message
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			ErrorMessage string
		}{ErrorMessage: errorMessage}

		tmpl, _ := template.ParseFiles("views/registration.html")
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
