package server

import (
	"html/template"
	"net/http"
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

}
