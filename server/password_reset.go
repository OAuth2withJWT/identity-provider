package server

import (
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
)

func (s *Server) handlePasswordResetPage(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")

	tmpl, _ := template.ParseFiles("views/password_reset.html")
	err := tmpl.Execute(w, struct {
		Email        string
		Code         string
		ErrorMessage string
	}{Email: email, Code: code, ErrorMessage: ""})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, _ := s.app.UserService.GetUserByEmail(email)
	err = s.app.VerificationService.Verify(user.UserId, code)
	if err != nil {
		data := struct {
			Email        string
			Code         string
			ErrorMessage string
		}{Email: email, Code: code, ErrorMessage: "Invalid url"}

		tmpl, _ := template.ParseFiles("views/password_reset.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	err := s.app.UserService.ResetPassword(user.UserId, newPassword)
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

	http.Redirect(w, r, "/login", http.StatusFound)
}
