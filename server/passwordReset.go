package server

import (
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/gorilla/mux"
)

func (s *Server) handlePasswordResetPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	code := vars["code"]

	tmpl, _ := template.ParseFiles("views/passwordReset.html")
	err := tmpl.Execute(w, struct {
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
	vars := mux.Vars(r)
	email := vars["email"]
	code := vars["code"]
	newPassword := r.FormValue("password")

	if newPassword != r.FormValue("confirmPassword") {
		data := struct {
			Email        string
			Code         string
			ErrorMessage string
		}{Email: email, Code: code, ErrorMessage: "Passwords don't match"}

		tmpl, _ := template.ParseFiles("views/passwordReset.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	user, _ := s.app.UserService.GetUserByEmail(email)
	err := s.app.VerificationService.ValidateCode(user.UserId, code)
	if err != nil {
		data := struct {
			Email        string
			Code         string
			ErrorMessage string
		}{Email: email, Code: code, ErrorMessage: "Invalid url"}

		tmpl, _ := template.ParseFiles("views/passwordReset.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	err = s.app.UserService.ResetPassword(user.UserId, newPassword)
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
			Code         string
			ErrorMessage string
		}{Email: email, Code: code, ErrorMessage: errorMessage}

		tmpl, _ := template.ParseFiles("views/passwordReset.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
