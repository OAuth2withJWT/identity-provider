package app

import (
	"github.com/OAuth2withJWT/identity-provider/app/validation"
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
	UserId    int
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

func (req *CreateUserRequest) Validate() string {
	v := &validation.Validator{}
	v.Errors = make(map[string]error)

	v.IsEmpty(req.FirstName)
	v.IsEmpty(req.LastName)
	v.IsEmpty(req.Username)
	v.IsEmpty(req.Email)
	v.IsEmpty(req.Password)
	v.IsEmail("email", req.Email)
	v.IsValidPassword("password", req.Password)

	return v.Error()
}

func (s *UserService) Create(req CreateUserRequest) (*User, error) {
	errorMessage := req.Validate()
	if errorMessage != "" {
		return nil, &FieldError{Message: errorMessage}
	}

	user, err := s.repository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user != (User{}) {
		errorMessage = "User with that email already exists"
		return nil, &FieldError{Message: errorMessage}
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword
	newUser, err := s.repository.Create(req)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) ValidateUserCredentials(email, password string) (User, error) {
	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return User{}, err
	}

	return user, nil
}

type UserRepository interface {
	Create(CreateUserRequest) (*User, error)
	GetUserByEmail(email string) (User, error)
	GetUserByID(user_id int) (User, error)
}

func (s *UserService) GetUserByEmail(email string) (User, error) {
	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(user_id int) (User, error) {
	user, err := s.repository.GetUserByID(user_id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
