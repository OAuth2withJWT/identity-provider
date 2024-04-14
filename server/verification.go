package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/gorilla/mux"
)

func (s *Server) handleVerificationPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	tmpl, _ := template.ParseFiles("views/verification.html")
	err := tmpl.Execute(w, struct {
		Email        string
		ErrorMessage string
	}{Email: email, ErrorMessage: ""})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleVerificationForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	var verificationCode string
	for i := 1; i <= 6; i++ {
		digit := r.FormValue(fmt.Sprintf("digit%d", i))
		verificationCode += digit
	}

	user, _ := s.app.UserService.GetUserByEmail(email)

	err := s.app.VerificationService.ValidateCode(user.UserId, verificationCode)
	if err != nil {
		var errorMessage string
		if fieldErr, ok := err.(*app.FieldError); ok {
			errorMessage = fieldErr.Message
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			Email        string
			ErrorMessage string
		}{Email: email, ErrorMessage: errorMessage}

		tmpl, _ := template.ParseFiles("views/verification.html")
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	fmt.Println("Verification code for ", email, ": ", verificationCode)

	http.Redirect(w, r, "/login", http.StatusFound)
}
