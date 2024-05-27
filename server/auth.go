package server

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
)

const authorizationCodeLength = 12
const authorizationCodeExpirationTime = 10

func (s *Server) handleAuthPage(w http.ResponseWriter, r *http.Request) {
	responseType := r.URL.Query().Get("response_type")
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")

	if clientID == "" || redirectURI == "" || responseType != "code" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	client, err := s.app.ClientService.GetClientByID(clientID)
	if err != nil {
		http.Error(w, "unauthorized_client", http.StatusUnauthorized)
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
	scopes := r.Form["scopes"]

	authSessionID := getAuthSessionIDFromCookie(r)
	authSessionData, err := getAuthSessionFromStore(authSessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURI := authSessionData.RedirectURI
	state := authSessionData.State
	clientID := authSessionData.State

	delete(authSessionStore, authSessionID)
	deleteAuthSessionCookie(w)

	authorizationCode, err := app.GenerateRandomBytes(authorizationCodeLength)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	codeInfo := map[string]interface{}{
		"client_id":    clientID,
		"redirect_uri": redirectURI,
		"state":        state,
		"expiration":   time.Now().Add(time.Minute * authorizationCodeExpirationTime).Unix(),
		"scopes":       scopes,
	}

	if err := saveAuthorizationCode(s, codeInfo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURI+"?code="+authorizationCode+"&state="+state, http.StatusFound)
}

func saveAuthorizationCode(s *Server, codeInfo map[string]interface{}) error {
	codeInfoJSON, err := json.Marshal(codeInfo)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = s.app.RedisClient.Set(ctx, "authorizationCode", string(codeInfoJSON), authorizationCodeExpirationTime).Err()
	if err != nil {
		return err
	}

	return nil
}
