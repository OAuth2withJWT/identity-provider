package server

import (
	"html/template"
	"net/http"
	"strings"
	"time"
)

func (s *Server) handleAuthPage(w http.ResponseWriter, r *http.Request) {
	responseType := r.URL.Query().Get("response_type")
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")

	if clientID == "" || redirectURI == "" || responseType != "code" {
		http.Error(w, "Invalid request: Missing or invalid parameters", http.StatusBadRequest)
		return
	}

	client, err := s.app.ClientService.GetClientByID(clientID)
	if err != nil {
		http.Error(w, "Unauthorized client", http.StatusUnauthorized)
		return
	}

	setAuthSession(w, clientID, redirectURI, state)

	scopes := strings.Split(client.Scope, ",")
	data := struct {
		ClientName string
		Scopes     []string
	}{
		ClientName: client.Name,
		Scopes:     scopes,
	}

	tmpl, err := template.ParseFiles("views/consent_screen.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleAuthForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	scopes := r.Form["scopes"]

	authSessionID := getAuthSessionIDFromCookie(r)
	authSessionData, err := getAuthSessionFromStore(authSessionID)
	if err != nil {
		http.Error(w, "Unable to retrieve authentication session", http.StatusInternalServerError)
		return
	}

	redirectURI := authSessionData.RedirectURI
	state := authSessionData.State
	clientID := authSessionData.State

	delete(authSessionStore, authSessionID)
	deleteAuthSessionCookie(w)

	authorizationCode, err := generateAuthorizationCode()
	if err != nil {
		http.Error(w, "Unable to generate authorization code", http.StatusInternalServerError)
		return
	}

	codeInfo := &AuthorizationCodeInfo{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		State:       state,
		Scopes:      scopes,
		Expiration:  time.Now().Add(time.Minute * authorizationCodeExpirationTime).Unix(),
	}

	storeAuthorizationCode(s, codeInfo)

	http.Redirect(w, r, redirectURI+"?code="+authorizationCode+"&state="+state, http.StatusFound)
}
