package server

import (
	"log"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	app    *app.Application
}

func New(a *app.Application) *Server {
	s := &Server{
		router: mux.NewRouter(),
		app:    a,
	}
	s.setupRoutes()
	return s
}

func (s *Server) Run() error {
	log.Println("Server started on port 8080")
	return http.ListenAndServe(":8080", s.router)
}

func (s *Server) setupRoutes() {
	s.router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	s.router.HandleFunc("/", s.handleHomePage).Methods("GET")
	s.router.HandleFunc("/registration", s.handleRegistrationPage).Methods("GET")
	s.router.HandleFunc("/registration", s.handleRegistrationForm).Methods("POST")
	s.router.HandleFunc("/login", s.handleLoginPage).Methods("GET")
	s.router.HandleFunc("/login", s.handleLoginForm).Methods("POST")
	s.router.HandleFunc("/logout", s.handleLogoutForm).Methods("POST")
	s.router.HandleFunc("/verification/{email}", s.handleVerificationPage).Methods("GET")
	s.router.HandleFunc("/verification/{email}", s.handleVerificationForm).Methods("POST")
}
