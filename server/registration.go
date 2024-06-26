package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/OAuth2withJWT/identity-provider/app/validation"
)

func (s *Server) handleRegistrationPage(w http.ResponseWriter, r *http.Request) {

	if _, ok := r.Context().Value(userContextKey).(app.User); ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	page := Page{
		FormFields: map[string]string{
			"First name": "",
			"Last name":  "",
			"Email":      "",
			"Username":   "",
			"Password":   "",
		},
	}

	tmpl, err := template.ParseFiles("views/registration.html")
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

func (s *Server) handleRegistrationForm(w http.ResponseWriter, r *http.Request) {
	page := Page{
		FormFields: map[string]string{
			"First name": r.FormValue("firstName"),
			"Last name":  r.FormValue("lastName"),
			"Email":      r.FormValue("email"),
			"Username":   r.FormValue("username"),
			"Password":   r.FormValue("password"),
		},
	}

	user, err := s.app.UserService.Create(app.CreateUserRequest{
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
		Password:  r.FormValue("password"),
	})

	if err != nil {
		switch v := err.(type) {
		case *validation.Error:
			page.FormErrors = make(map[string]string)

			for field, errs := range v.Errors {
				page.FormErrors[field] = errs[0].Error()
			}

			tmpl, err := template.ParseFiles("views/registration.html")
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

	code, err := s.app.VerificationService.CreateVerification(user.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Print("Use this link to verify your account: http://localhost:8080/verification?email=" + user.Email + "&code=" + code)

	http.Redirect(w, r, "/account-message?status=email-sent", http.StatusFound)
}
