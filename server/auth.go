package server

import (
	"html/template"
	"net/http"
	"strings"
)

func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) {
	responseType := r.URL.Query().Get("response_type")
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	//scope := r.URL.Query().Get("scope")
	//state := r.URL.Query().Get("state")

	if clientID == "" || redirectURI == "" || responseType != "code" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	client, err := s.app.ClientService.GetClientByID(clientID)
	if err != nil {
		http.Error(w, "unauthorized_client", http.StatusUnauthorized)
		return
	}

	scopes := strings.Split(client.Scope, ",")
	page := Page{
		FormFields: map[string]string{
			"Client Name": client.Name,
		},
		QueryParameters: make(map[string]string),
	}

	for _, scope := range scopes {
		page.QueryParameters[scope] = scope
	}

	tmpl, err := template.ParseFiles("views/consent_screen.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//example : http://localhost:8080/oauth2/auth?response_type=code&client_id=ce4de456add545acc0205b79bb7bbc59&redirect_uri=https%3A%2F%2Fexample-app.com%2Fauth&scope=photos&state=5ca75bd30
	//authorizationCode, _ := s.app.ClientService.GenerateAuthorizationCode(clientID, redirectURI, scope)

	//http.Redirect(w, r, redirectURI+"?code="+authorizationCode+"&state="+state, http.StatusFound)
}
