package app

import (
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository UserRepository
}

func NewUserService(ur UserRepository) *UserService {
	return &UserService{
		repository: ur,
	}
}

type User struct {
	FirstName string
	LastName  string
	Email     string
	Username  string
	Password  string
}

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	Username  string
	Password  string
}

type FieldError struct {
	Message string
}

func (e *FieldError) Error() string {
	return e.Message
}

func (s *UserService) Create(req CreateUserRequest) (*User, error) {
	ErrorMessage := fieldsNotEmpty(req)
	if ErrorMessage != "" {
		return nil, &FieldError{Message: ErrorMessage}
	}

	ErrorMessage = validatePassword(req.Password)
	if ErrorMessage != "" {
		return nil, &FieldError{Message: ErrorMessage}
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword
	user, err := s.repository.Create(req)

	if err != nil {
		return nil, err
	}

	return user, nil
}

type UserRepository interface {
	Create(CreateUserRequest) (*User, error)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func fieldsNotEmpty(req CreateUserRequest) string {
	firstName := strings.TrimSpace(req.FirstName)
	lastName := strings.TrimSpace(req.LastName)
	email := strings.TrimSpace(req.Email)
	username := strings.TrimSpace(req.Username)
	password := req.Password

	if firstName == "" || lastName == "" || email == "" || username == "" || password == "" {
		return "Fields cannot be empty"
	}

	return ""
}

func validatePassword(password string) string {
	var errors []string

	rules := map[string]func(string) bool{
		"at least 8 characters": func(s string) bool { return len(s) >= 8 },
		"one uppercase letter":  func(s string) bool { return containsType(s, unicode.IsUpper) },
		"one lowercase letter":  func(s string) bool { return containsType(s, unicode.IsLower) },
		"one digit":             func(s string) bool { return containsType(s, unicode.IsDigit) },
		"one special character": func(s string) bool { return containsSpecialChar(s) },
	}

	for rule, isValid := range rules {
		if !isValid(password) {
			errors = append(errors, rule)
		}
	}

	return "Password must contain " + strings.Join(errors, ", ")
}

func containsType(s string, check func(rune) bool) bool {
	for _, char := range s {
		if check(char) {
			return true
		}
	}
	return false
}

func containsSpecialChar(s string) bool {
	for _, char := range s {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			return true
		}
	}
	return false
}
