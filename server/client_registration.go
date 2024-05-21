package server

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/OAuth2withJWT/identity-provider/app/validation"
)

func (s *Server) handleClientRegistrationPage(w http.ResponseWriter, r *http.Request) {
	page := Page{
		FormFields: map[string]string{
			"Client Name":  "",
			"Scope":        "",
			"Redirect URI": "",
		},
	}

	tmpl, err := template.ParseFiles("views/client_registration.html")
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

func (s *Server) handleClientRegistrationForm(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userContextKey).(app.User)

	client, err := s.app.ClientService.Create(app.CreateClientRequest{
		Name:        r.FormValue("clientName"),
		Scope:       r.FormValue("scope"),
		RedirectURI: r.FormValue("redirectUri"),
		CreatedBy:   user.UserId,
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch v := err.(type) {
		case *validation.Error:
			formErrors := make(map[string]string)
			for field, errs := range v.Errors {
				formErrors[field] = errs[0].Error()
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":      err.Error(),
				"formErrors": formErrors,
			})
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"clientId":     client.Id,
		"clientSecret": client.Secret,
	}
	json.NewEncoder(w).Encode(response)
}
