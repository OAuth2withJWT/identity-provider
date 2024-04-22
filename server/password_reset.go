package server

import (
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
)

func (s *Server) handlePasswordResetPage(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")

	user, err := s.app.UserService.GetUserByEmail(email)
	if err != nil {
		http.Redirect(w, r, "/account-message?status=verification-error", http.StatusFound)
	}

	err = s.app.VerificationService.Verify(user.UserId, code)
	if err != nil {
		http.Redirect(w, r, "/account-message?status=verification-error", http.StatusFound)
	}

	tmpl, _ := template.ParseFiles("views/password_reset.html")
	err = tmpl.Execute(w, struct {
		Email         string
		Code          string
		ErrorPassword string
		Password      string
	}{Email: email, Code: code, ErrorPassword: "", Password: ""})

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
			Email         string
			Code          string
			ErrorPassword string
			Password      string
		}{Email: email, Code: code, ErrorPassword: "Passwords don't match", Password: newPassword}

		tmpl, _ := template.ParseFiles("views/password_reset.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	user, _ := s.app.UserService.GetUserByEmail(email)
	req := app.PasswordResetRequest{
		UserId:   user.UserId,
		Password: newPassword,
	}
	err := s.app.UserService.ResetPassword(&req)
	if err != nil {
		data := struct {
			Email         string
			Code          string
			ErrorPassword string
			Password      string
		}{Email: email, Code: code, ErrorPassword: req.ErrorPassword, Password: req.Password}

		tmpl, _ := template.ParseFiles("views/password_reset.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Redirect(w, r, "/account-message?status=password-reset", http.StatusFound)
}
