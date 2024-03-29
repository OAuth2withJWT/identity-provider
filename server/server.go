package server

import (
	"html/template"
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
	return &Server{
		router: mux.NewRouter(),
		app:    a,
	}
}

func (s *Server) Run() error {
	s.router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	s.router.HandleFunc("/registration", RegistrationFormHandler).Methods("GET")
	s.router.HandleFunc("/login", LoginFormHandler).Methods("GET")
	s.router.HandleFunc("/registration", s.RenderingRegistrationDetails).Methods("POST")
	s.router.HandleFunc("/", HomePageHandler).Methods("GET")

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

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("views/index.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("views/login.html")
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

	_, err := s.app.UserService.Create(req)
	if err != nil {
		var errorMessage string
		if fieldErr, ok := err.(*app.FieldError); ok {
			errorMessage = fieldErr.Message
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			ErrorMessage string
		}{ErrorMessage: errorMessage}

		tmpl, _ := template.ParseFiles("views/registration.html")
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
