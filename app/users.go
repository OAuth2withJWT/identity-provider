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

type PasswordResetRequest struct {
	UserId   int
	Password string
}

type Error struct {
	Message string
	Kind    string
}

func (e *Error) Error() string {
	return e.Message
}

func (req *CreateUserRequest) validateFields() string {
	v := &validation.Validator{}
	v.Errors = make(map[string]error)

	v.IsEmpty("First name", req.FirstName)
	v.IsEmpty("Last name", req.LastName)
	v.IsEmpty("Username", req.Username)
	v.IsEmpty("Email", req.Email)
	v.IsEmpty("Password", req.Password)
	v.IsEmail("email", req.Email)
	v.IsValidPassword("password", req.Password)

	return v.Error()
}

func (s *UserService) Create(req CreateUserRequest) (*User, error) {
	errorMessage := req.validateFields()
	if errorMessage != "" {
		return nil, &Error{Message: errorMessage}
	}

	user, _ := s.repository.GetUserByEmail(req.Email)
	if user != (User{}) {
		errorMessage = "User with that email already exists"
		return nil, &Error{Message: errorMessage}
	}

	hashedPassword, err := hashPassword(req.Password)
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
	errorMessage := "Invalid email or password"
	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return User{}, &Error{Message: errorMessage}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return User{}, &Error{Message: errorMessage}
	}

	return user, nil
}

type UserRepository interface {
	Create(CreateUserRequest) (*User, error)
	GetUserByEmail(email string) (User, error)
	GetUserByID(user_id int) (User, error)
	UpdatePassword(hashedPassword string, userId int) error
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

func (req *PasswordResetRequest) validateNewPassword() string {
	v := &validation.Validator{}
	v.Errors = make(map[string]error)
	v.IsEmpty("password", req.Password)
	v.IsValidPassword("password", req.Password)

	return v.Error()
}

func (s *UserService) ResetPassword(req PasswordResetRequest) error {
	errorMessage := req.validateNewPassword()

	if errorMessage != "" {
		return &Error{Message: errorMessage}
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	err = s.repository.UpdatePassword(hashedPassword, req.UserId)
	if err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
