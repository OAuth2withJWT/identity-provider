package server

import (
	"html/template"
	"log"
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
	req := app.RegistrationRequest{
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
		Password:  r.FormValue("password"),
	}

	user := s.app.UserService.Create(&req)
	if user == nil {
		data := req

		tmpl, _ := template.ParseFiles("views/registration.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	code, err := s.app.VerificationService.CreateVerification(user.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Print("Use this link to verify your account: http://localhost:8080/verification?email=" + user.Email + "&code=" + code)

	http.Redirect(w, r, "/account-message?status=email-sent", http.StatusFound)
}
