package server

import (
	"html/template"
	"log"
	"net/http"
)

func (s *Server) handleVerificationPage(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")

	user, _ := s.app.UserService.GetUserByEmail(email)
	err := s.app.VerificationService.Verify(user.UserId, code)
	if err != nil {
		http.Redirect(w, r, "/account-status-message?verification-error=true", http.StatusFound)
	}

	http.Redirect(w, r, "/account-status-message?verified=true", http.StatusFound)
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
	log.Print("Use this link to reset your password: http://localhost:8080/password-reset?email=" + email + "&code=" + code)

	http.Redirect(w, r, "/account-status-message", http.StatusFound)
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	verificationError := false
	errorStr := r.URL.Query().Get("verification-error")
	if errorStr == "true" {
		verificationError = true
	}

	verified := false
	verifiedStr := r.URL.Query().Get("verified")
	if verifiedStr == "true" {
		verified = true
	}

	successReset := false
	successResetStr := r.URL.Query().Get("success-reset")
	if successResetStr == "true" {
		successReset = true
	}

	tmpl, _ := template.ParseFiles("views/message.html")
	err := tmpl.Execute(w, struct {
		VerificationError bool
		SuccessReset      bool
		Verified          bool
	}{VerificationError: verificationError, SuccessReset: successReset, Verified: verified})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
