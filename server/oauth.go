package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
)

type ScopeData struct {
	Scope       string
	Description string
}

var scopeToDescription = map[string]string{
	"openid":            "associate you with your personal info",
	"cards:read":        "read your cards info",
	"transactions:read": "read your transactions info",
}

func GetScopeData(scopes []string) []ScopeData {
	var scopeData []ScopeData
	for _, scope := range scopes {
		description, exists := scopeToDescription[scope]
		if exists {
			scopeData = append(scopeData, ScopeData{
				Scope:       scope,
				Description: description,
			})
		}
	}
	return scopeData
}

func (s *Server) handleAuthPage(w http.ResponseWriter, r *http.Request) {
	responseType := r.URL.Query().Get("response_type")
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")
	codeChallenge := r.URL.Query().Get("code_challenge")
	codeChallengeMethod := r.URL.Query().Get("code_challenge_method")

	if clientID == "" || redirectURI == "" || responseType != "code" || codeChallenge == "" || codeChallengeMethod == "" {
		http.Error(w, "Invalid request: Missing or invalid parameters", http.StatusBadRequest)
		return
	}

	if codeChallengeMethod != "S256" {
		http.Error(w, "Invalid code challenge method", http.StatusBadRequest)
		return
	}

	client, err := s.app.ClientService.GetClientByID(clientID)
	if err != nil {
		http.Error(w, "Unauthorized client", http.StatusUnauthorized)
		return
	}

	authID, err := app.GenerateAuthID()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, authID)

	s.app.AuthService.Create(app.Auth{
		AuthID:              authID,
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		State:               state,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	})

	sessionID := getSessionIDFromCookie(r)
	session, _ := s.app.SessionService.ValidateSession(sessionID)
	user, err := s.app.UserService.GetUserByID(session.UserId)
	if err != nil {
		http.Error(w, "unauthorized_user", http.StatusUnauthorized)
		return
	}

	scopes := strings.Split(client.Scope, ",")
	scopeData := GetScopeData(scopes)

	data := struct {
		ClientName string
		Scopes     []ScopeData
		Email      string
	}{
		ClientName: client.Name,
		Scopes:     scopeData,
		Email:      user.Email,
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
	authID := getAuthIDFromCookie(r)

	authData, err := s.app.AuthService.Get(authID)
	if err != nil {
		http.Error(w, "Unable to retrieve auth data", http.StatusInternalServerError)
		return
	}

	err = s.app.AuthService.Delete(authID)
	if err != nil {
		http.Error(w, "Unable to delete auth data", http.StatusInternalServerError)
		return
	}

	redirectURI := authData.RedirectURI
	state := authData.State
	clientID := authData.ClientID
	codeChallenge := authData.CodeChallenge
	codeChallengeMethod := authData.CodeChallengeMethod

	deleteAuthCookie(w)

	if len(scopes) == 0 {
		parsedURI, err := url.Parse(redirectURI)
		if err != nil {
			http.Error(w, "Invalid redirect URI", http.StatusInternalServerError)
			return
		}

		rootURI := parsedURI.Scheme + "://" + parsedURI.Host

		http.Redirect(w, r, rootURI, http.StatusFound)
		return
	}

	authorizationCode, err := app.GenerateAuthorizationCode()
	if err != nil {
		http.Error(w, "Unable to generate authorization code", http.StatusInternalServerError)
		return
	}

	sessionID := getSessionIDFromCookie(r)
	session, _ := s.app.SessionService.ValidateSession(sessionID)

	s.app.AuthorizationCodeService.Create(app.AuthorizationCode{
		Value:               authorizationCode,
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		State:               state,
		Scopes:              scopes,
		Expiration:          time.Now().Add(app.AuthorizationCodeExpiration).Unix(),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		UserId:              session.UserId,
	})

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

	valid, err := s.validateCodeChallenge(req.CodeVerifier, req.Code)
	if err != nil {
		log.Printf("Error validating code challenge: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !valid {
		log.Print("Invalid code verifier")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	scopes, userId, err := s.validateAuthorizationCode(req.Code, req.ClientID, req.RedirectURI)
	if err != nil {
		log.Print("Invalid authorization code")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken, err := s.createAccessToken(scopes, userId)
	if err != nil {
		log.Print("Error generating token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	atHash, err := s.ComputeAtHash(accessToken)
	if err != nil {
		log.Print("Error computing at_hash")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.app.AuthorizationCodeService.Delete(req.Code)
	if err != nil {
		http.Error(w, "Unable to delete authorization code data", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   (time.Duration(s.RSAConfig.TokenExpirationTime)).String(),
	}

	user, err := s.app.UserService.GetUserByID(userId)
	if err != nil {
		http.Error(w, "unauthorized_user", http.StatusUnauthorized)
		return
	}

	if ContainsScope(scopes, "openid") {
		idToken, err := s.createIDToken(req.ClientID, user, atHash)
		if err != nil {
			log.Print("Error generating ID token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.IDToken = idToken
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

	scopes, err := s.validateAccessToken(req.Token, client)
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
