package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	r *mux.Router
}

func New() *Server {
	return &Server{
		r: mux.NewRouter(),
	}
}

func (server *Server) Run() error {
	server.r.HandleFunc("/registration", RegistrationFormHandler).Methods("GET")
	server.r.HandleFunc("/registration", RenderingRegistrationDetails).Methods("POST")

	log.Println("Server started on port 8080")
	return http.ListenAndServe(":8080", server.r)
}

type User struct {
	Email    string
	Username string
	Password string
}

func RegistrationFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("views/registration.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderingRegistrationDetails(w http.ResponseWriter, r *http.Request) {
	registrationDetails := User{
		Email:    r.FormValue("email"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	_ = registrationDetails

	tmpl, _ := template.ParseFiles("views/registration.html")
	err := tmpl.Execute(w, struct{ Success bool }{true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
