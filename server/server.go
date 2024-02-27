package server

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	r  *mux.Router
	db *sql.DB
}

func New(db *sql.DB) *Server {
	return &Server{
		r:  mux.NewRouter(),
		db: db,
	}
}

func (server *Server) Run() error {
	server.r.HandleFunc("/registration", RegistrationFormHandler).Methods("GET")
	server.r.HandleFunc("/registration", server.RenderingRegistrationDetails).Methods("POST")

	log.Println("Server started on port 8080")
	return http.ListenAndServe(":8080", server.r)
}

type User struct {
	FirstName string
	LastName  string
	Email     string
	Username  string
	Password  string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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
	registrationDetails := User{
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
	}

	hashedPassword, err := HashPassword(r.FormValue("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	registrationDetails.Password = hashedPassword

	_, err = s.db.Query("INSERT INTO users (first_name, last_name, email, username, password) VALUES ($1, $2, $3, $4, $5)",
		registrationDetails.FirstName, registrationDetails.LastName, registrationDetails.Email, registrationDetails.Username, registrationDetails.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles("views/registration.html")
	err = tmpl.Execute(w, struct{ Success bool }{true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
