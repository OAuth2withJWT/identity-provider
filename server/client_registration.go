package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
)

func (s *Server) handleClientRegistrationPage(w http.ResponseWriter, r *http.Request) {
	page := Page{
		FormFields: map[string]string{
			"Client Name":  "",
			"Scope":        "",
			"Redirect URI": "",
		},
	}

	tmpl, _ := template.ParseFiles("public/html/client_registration.html")
	err := tmpl.Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleClientRegistrationForm(w http.ResponseWriter, r *http.Request) {

	client, err := s.app.ClientService.Create(app.CreateClientRequest{
		ClientName:  r.FormValue("clientName"),
		Scope:       r.FormValue("scope"),
		RedirectURI: r.FormValue("redirectUri"),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "{\"clientId\": \"%s\", \"clientSecret\": \"%s\"}", client.ClientId, client.ClientSecret)

}
