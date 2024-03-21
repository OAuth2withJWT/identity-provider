package server

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	app    *app.Application
}

func New(a *app.Application) *Server {
	return &Server{
		router: mux.NewRouter(),
		app:    a,
	}
}

func (s *Server) Run() error {
	s.router.HandleFunc("/registration", RegistrationFormHandler).Methods("GET")
	s.router.HandleFunc("/registration", s.RenderingRegistrationDetails).Methods("POST")
	s.router.HandleFunc("/login", s.LoginFormHandler).Methods("GET")
	s.router.HandleFunc("/login", s.LoginHandler).Methods("POST")
	s.router.HandleFunc("/homepage", HomePageHandler).Methods("GET")

	log.Println("Server started on port 8080")
	return http.ListenAndServe(":8080", s.router)
}

func RegistrationFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("views/registration.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) RenderingRegistrationDetails(w http.ResponseWriter, r *http.Request) {
	req := app.CreateUserRequest{
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
		Password:  r.FormValue("password"),
	}
	user, err := s.app.UserService.Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles("views/registration.html")
	err = tmpl.Execute(w, struct {
		Success  bool
		Username string
	}{true, user.Username})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("views/homepage.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := s.app.SessionService.GetSessionIDFromCookie(r)

	if sessionID != 0 {
		http.Redirect(w, r, "/homepage", http.StatusFound)
	}

	tmpl, _ := template.ParseFiles("views/login.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	userID, err := s.app.UserService.Authenticate(username, password)
	if err != nil {
		tmpl, err := template.ParseFiles("views/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Invalid username or password",
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	sessionID, err := s.app.SessionService.SaveSession(userID, time.Now().Add(24*time.Hour))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.app.SessionService.SetSessionCookie(w, sessionID)

	http.Redirect(w, r, "/homepage", http.StatusFound)
}
