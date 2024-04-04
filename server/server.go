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
	s.router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	s.router.HandleFunc("/registration", s.RegistrationFormHandler).Methods("GET")
	s.router.HandleFunc("/registration", s.RenderingRegistrationDetails).Methods("POST")
	s.router.HandleFunc("/login", s.LoginFormHandler).Methods("GET")
	s.router.HandleFunc("/login", s.LoginHandler).Methods("POST")
	s.router.HandleFunc("/", s.HomePageHandler).Methods("GET")

	log.Println("Server started on port 8080")
	return http.ListenAndServe(":8080", s.router)
}

func (s *Server) RegistrationFormHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := GetSessionIDFromCookie(r)
	_, err := s.app.SessionService.ValidateSession(sessionID)

	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	tmpl, _ := template.ParseFiles("views/registration.html")
	err = tmpl.Execute(w, nil)
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

func (s *Server) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := GetSessionIDFromCookie(r)
	session, err := s.app.SessionService.ValidateSession(sessionID)

	var username string

	if err == nil {
		user, err := s.app.UserService.GetUserByID(session.UserId)
		if err == nil {
			username = user.Username
		}
	}
	tmpl, _ := template.ParseFiles("views/homepage.html")
	err = tmpl.Execute(w, struct {
		Username string
	}{username})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := GetSessionIDFromCookie(r)
	_, err := s.app.SessionService.ValidateSession(sessionID)

	if err != nil {
		tmpl, _ := template.ParseFiles("views/login.html")
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := s.app.UserService.ValidateUserCredentials(email, password)
	if err != nil {
		tmpl, err := template.ParseFiles("views/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Invalid email or password",
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	sessionID, err := s.app.SessionService.CreateSession(user.UserId, time.Now().Add(app.SessionDurationInHours*time.Hour))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SetSessionCookie(w, sessionID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(app.SessionDurationInHours * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func GetSessionIDFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	sessionID := cookie.Value
	return sessionID
}
