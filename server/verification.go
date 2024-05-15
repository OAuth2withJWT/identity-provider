package server

import (
	"html/template"
	"log"
	"net/http"
)

func (s *Server) handleVerification(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")

	user, _ := s.app.UserService.GetUserByEmail(email)
	err := s.app.VerificationService.Verify(user.UserId, code)
	if err != nil {
		http.Redirect(w, r, "/account-message?status=verification-error", http.StatusFound)
	}

	http.Redirect(w, r, "/account-message?status=verified", http.StatusFound)
}

func (s *Server) handleEnterEmailPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/enter_email.html")
	err := tmpl.Execute(w, Page{
		FormFields: map[string]string{
			"Email": "",
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleEnterEmailForm(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	page := Page{
		FormFields: map[string]string{
			"Email": email,
		},
	}

	user, err := s.app.UserService.GetUserByEmail(email)
	if err != nil {
		page.FormErrors = make(map[string]string)

		page.FormErrors["Email"] = "Invalid email"

		tmpl, _ := template.ParseFiles("templates/enter_email.html")
		err = tmpl.Execute(w, page)
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

	http.Redirect(w, r, "/account-message?status=email-sent", http.StatusFound)
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	verificationError := false
	verified := false
	successReset := false

	status := r.URL.Query().Get("status")
	if status == "verification-error" {
		verificationError = true
	} else if status == "verified" {
		verified = true
	} else if status == "password-reset" {
		successReset = true
	}
	tmpl, _ := template.ParseFiles("templates/message.html")
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
