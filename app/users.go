package app

import (
	"regexp"
	"strings"

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
	isValid := true
	errorMessage := "Password must contain at least "
	if len(password) < 8 {
		if !isValid {
			errorMessage += ", "
		} else {
			isValid = false
		}
		errorMessage += "8 characters"
	}

	hasUppercase := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if !hasUppercase {
		if !isValid {
			errorMessage += ", "
		} else {
			isValid = false
		}
		errorMessage += "one uppercase letter"
	}

	hasLowercase := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	if !hasLowercase {
		if !isValid {
			errorMessage += ", "
		} else {
			isValid = false
		}
		errorMessage += "one lowercase letter"
	}

	hasDigit, _ := regexp.MatchString(`\d`, password)
	if !hasDigit {
		if !isValid {
			errorMessage += ", "
		} else {
			isValid = false
		}
		errorMessage += "one digit"
	}

	hasSpecialChar, _ := regexp.MatchString(`[\W_]`, password)
	if !hasSpecialChar {
		if !isValid {
			errorMessage += ", "
		} else {
			isValid = false
		}
		errorMessage += "one special character"
	}

	if isValid {
		return ""
	}
	return errorMessage
}
