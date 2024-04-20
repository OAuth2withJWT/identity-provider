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

	err := s.app.VerificationService.Verify(user.UserId, verificationCode)
	if err != nil {
		var errorMessage string
		if fieldErr, ok := err.(*app.Error); ok {
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

func (s *Server) handleEnterEmailPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("views/enter_email.html")
	err := tmpl.Execute(w, struct {
		ErrorMessage string
	}{ErrorMessage: ""})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleEnterEmailForm(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	user, err := s.app.UserService.GetUserByEmail(email)
	if err != nil {
		data := struct {
			Email        string
			ErrorMessage string
		}{Email: email, ErrorMessage: "Invalid email"}

		tmpl, _ := template.ParseFiles("views/enter_email.html")
		err = tmpl.Execute(w, data)
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
	fmt.Println("Use this link to reset your password: http://localhost:8080/password-reset?email=" + email + "&code=" + code)

	http.Redirect(w, r, "/success-message", http.StatusFound)
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("views/message.html")
	err := tmpl.Execute(w, struct {
		Message string
	}{Message: "Check your email and follow instructions."})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
