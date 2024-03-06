package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/OAuth2withJWT/identity-provider/app"
	"github.com/gorilla/mux"
)

type Server struct {
	router      *mux.Router
	userService *app.UserService
}

func New(s *app.Application) *Server {
	return &Server{
		router:      mux.NewRouter(),
		userService: s.UserService,
	}
}

func (s *Server) Run() error {
	s.router.HandleFunc("/registration", RegistrationFormHandler).Methods("GET")
	s.router.HandleFunc("/registration", s.RenderingRegistrationDetails).Methods("POST")

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
	}
	user, err := s.userService.Create(req)
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
