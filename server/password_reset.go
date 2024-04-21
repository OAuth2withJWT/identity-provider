package server

import (
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
)

func (s *Server) handlePasswordResetPage(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")

	user, _ := s.app.UserService.GetUserByEmail(email)
	err := s.app.VerificationService.Verify(user.UserId, code)
	if err != nil {
		http.Redirect(w, r, "/account-status-message?verification-error=true", http.StatusFound)
	}

	tmpl, _ := template.ParseFiles("views/password_reset.html")
	err = tmpl.Execute(w, struct {
		Email        string
		Code         string
		ErrorMessage string
	}{Email: email, Code: code, ErrorMessage: ""})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handlePasswordResetForm(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")
	newPassword := r.FormValue("password")

	if newPassword != r.FormValue("confirmPassword") {
		data := struct {
			Email        string
			Code         string
			ErrorMessage string
		}{Email: email, Code: code, ErrorMessage: "Passwords don't match"}

		tmpl, _ := template.ParseFiles("views/password_reset.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	user, _ := s.app.UserService.GetUserByEmail(email)
	err := s.app.UserService.ResetPassword(app.PasswordResetRequest{UserId: user.UserId, Password: newPassword})
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
			Code         string
			ErrorMessage string
		}{Email: email, Code: code, ErrorMessage: errorMessage}

		tmpl, _ := template.ParseFiles("views/password_reset.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Redirect(w, r, "/account-status-message?success-reset=true", http.StatusFound)
}
