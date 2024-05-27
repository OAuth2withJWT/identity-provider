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
	s.router.Handle("/", s.withUser(s.handleHomePage)).Methods("GET")
	s.router.Handle("/registration", s.withUser(s.handleRegistrationPage)).Methods("GET")
	s.router.HandleFunc("/registration", s.handleRegistrationForm).Methods("POST")
	s.router.Handle("/login", s.withUser(s.handleLoginPage)).Methods("GET")
	s.router.HandleFunc("/login", s.handleLoginForm).Methods("POST")
	s.router.HandleFunc("/logout", s.handleLogoutForm).Methods("POST")
	s.router.HandleFunc("/verification", s.handleVerification).Methods("GET")
	s.router.HandleFunc("/password-reset", s.handlePasswordResetPage).Methods("GET")
	s.router.HandleFunc("/password-reset", s.handlePasswordResetForm).Methods("POST")
	s.router.HandleFunc("/request-password-reset", s.handleEnterEmailPage).Methods("GET")
	s.router.HandleFunc("/request-password-reset", s.handleEnterEmailForm).Methods("POST")
	s.router.HandleFunc("/account-message", s.handleMessage).Methods("GET")
	s.router.Handle("/client-registration", s.protected(http.HandlerFunc(s.handleClientRegistrationPage))).Methods("GET")
	s.router.Handle("/client-registration", s.protected(http.HandlerFunc(s.handleClientRegistrationForm))).Methods("POST")
	s.router.Handle("/oauth2/auth", s.protected(http.HandlerFunc(s.handleAuthPage))).Methods("GET")
	s.router.HandleFunc("/oauth2/auth", s.handleAuthForm).Methods("POST")
}
