package server

import (
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/OAuth2withJWT/identity-provider/app/validation"
)

func (s *Server) handlePasswordResetPage(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")

	page := Page{
		FormFields: map[string]string{
			"Password": "",
		},
		QueryParameters: map[string]string{
			"Email": email,
			"Code":  code,
		},
	}

	user, err := s.app.UserService.GetUserByEmail(email)
	if err != nil {
		http.Redirect(w, r, "/account-message?status=verification-error", http.StatusFound)
	}

	err = s.app.VerificationService.Verify(user.UserId, code)
	if err != nil {
		http.Redirect(w, r, "/account-message?status=verification-error", http.StatusFound)
	}

	tmpl, err := template.ParseFiles("views/password_reset.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handlePasswordResetForm(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")
	newPassword := r.FormValue("password")

	page := Page{
		FormFields: map[string]string{
			"Password": newPassword,
		},
		QueryParameters: map[string]string{
			"Email": email,
			"Code":  code,
		},
	}

	if newPassword != r.FormValue("confirmPassword") {
		page.FormErrors = make(map[string]string)

		page.FormErrors["Password"] = "Passwords don't match"

		tmpl, err := template.ParseFiles("views/password_reset.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, page)
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
		switch v := err.(type) {
		case *validation.Error:
			page.FormErrors = make(map[string]string)

			page.FormErrors["Password"] = v.Errors["Password"][0].Error()

			tmpl, err := template.ParseFiles("views/password_reset.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/account-message?status=password-reset", http.StatusFound)
}
