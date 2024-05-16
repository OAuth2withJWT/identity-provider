package server

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/OAuth2withJWT/identity-provider/app/validation"
)

func (s *Server) handleClientRegistrationPage(w http.ResponseWriter, r *http.Request) {

	sessionID := getSessionIDFromCookie(r)
	_, err := s.app.SessionService.ValidateSession(sessionID)

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	page := Page{
		FormFields: map[string]string{
			"Client Name":  "",
			"Scope":        "",
			"Redirect URI": "",
		},
	}

	tmpl, _ := template.ParseFiles("templates/client_registration.html")
	err = tmpl.Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleClientRegistrationForm(w http.ResponseWriter, r *http.Request) {

	sessionID := getSessionIDFromCookie(r)
	session, err := s.app.SessionService.ValidateSession(sessionID)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	client, err := s.app.ClientService.Create(app.CreateClientRequest{
		Name:        r.FormValue("clientName"),
		Scope:       r.FormValue("scope"),
		RedirectURI: r.FormValue("redirectUri"),
		CreatedBy:   session.UserId,
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
