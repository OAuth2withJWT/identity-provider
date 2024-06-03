package server

import (
	"encoding/json"
	"html/template"
	"log"
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
	clientID := authSessionData.ClientID

	delete(authSessionStore, authSessionID)
	deleteAuthSessionCookie(w)

	if len(scopes) == 0 {
		http.Redirect(w, r, redirectURI+"?error=access_denied&state="+state, http.StatusFound)
		return
	}

	authorizationCode, err := generateAuthorizationCode()
	if err != nil {
		http.Error(w, "Unable to generate authorization code", http.StatusInternalServerError)
		return
	}

	codeInfo := &AuthorizationCodeInfo{
		Value:       authorizationCode,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		State:       state,
		Scopes:      scopes,
		Expiration:  time.Now().Add(time.Minute * authorizationCodeExpirationTime).Unix(),
	}

	s.storeAuthorizationCode(codeInfo)

	http.Redirect(w, r, redirectURI+"?code="+authorizationCode+"&state="+state, http.StatusFound)
}

func (s *Server) handleTokenRequest(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Print("Invalid request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.GrantType != "authorization_code" {
		log.Print("Unsupported grant type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scopes, err := s.validateAuthorizationCode(req.Code, req.ClientID, req.RedirectURI)
	if err != nil {
		log.Print("Invalid authorization code")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := createToken(req.ClientID, req.ClientSecret, scopes)
	if err != nil {
		log.Print("Error generating token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   (time.Hour * time.Duration(tokenExpirationTime)).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleTokenVerification(w http.ResponseWriter, r *http.Request) {
	var req VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := VerificationErrorResponse{
			Error:            "invalid_request",
			ErrorDescription: "Request body is not as expected",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	clientID, err := extractClientID(req.Token)
	if err != nil {
		response := VerificationErrorResponse{
			Error:            "invalid_client",
			ErrorDescription: "The client authentication was invalid",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	client, err := s.app.ClientService.GetClientByID(clientID)
	if err != nil {
		response := VerificationErrorResponse{
			Error:            "invalid_client",
			ErrorDescription: "The client authentication was invalid",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	scopes, err := s.verifyToken(req.Token, client)
	if err != nil {
		response := VerificationResponse{
			Active: "false",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	scopesStr := []string{}
	for _, s := range scopes {
		if scope, ok := s.(string); ok {
			scopesStr = append(scopesStr, scope)
		}
	}

	response := VerificationResponse{
		Active: "true",
		Scope:  scopesStr,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
