package server

import (
	"html/template"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
)

func (s *Server) handleHomePage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(app.User)
	var username string

	if ok && user != (app.User{}) {
		username = user.Username
	}

	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		Username string
	}{Username: username})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
