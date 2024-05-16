package server

import (
	"html/template"
	"net/http"
)

func (s *Server) handleHomePage(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionIDFromCookie(r)
	session, err := s.app.SessionService.ValidateSession(sessionID)

	var username string

	if err == nil {
		user, err := s.app.UserService.GetUserByID(session.UserId)
		if err == nil {
			username = user.Username
		}
	}
	tmpl, _ := template.ParseFiles("views/index.html")
	err = tmpl.Execute(w, struct {
		Username string
	}{username})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
