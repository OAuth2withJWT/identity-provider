package app

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Application struct {
	UserService *UserService
}
type User struct {
	FirstName string
	LastName  string
	Email     string
	Username  string
	Password  string
}

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db,
	}
}

type CreateUserRequest struct {
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

func (s *UserService) Create(req CreateUserRequest) (*User, error) {

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	row := s.db.QueryRow("INSERT INTO users (first_name, last_name, email, username, password) VALUES ($1, $2, $3, $4, $5) RETURNING first_name, last_name, email, username, password",
		req.FirstName, req.LastName, req.Email, req.Username, hashedPassword)

	user := &User{}
	err = row.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type UserRepository interface {
	Create(CreateUserRequest) (*User, error)
}
