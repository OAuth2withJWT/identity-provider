package app

import (
	"fmt"

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
	UserId        int
	Password      string
	ErrorPassword string
}

func (req *CreateUserRequest) validateRegistrationFields(s *UserService) error {
	v := validation.New()
	v.IsEmpty("First name", req.FirstName)
	v.IsEmpty("Last name", req.LastName)
	v.IsEmpty("Username", req.Username)
	v.IsEmpty("Email", req.Email)
	v.IsEmpty("Password", req.Password)
	v.IsEmail("Email", req.Email)
	v.IsValidPassword("Password", req.Password)
	if s.hasUserWithEmail(req.Email) {
		v.AddError("Email", fmt.Errorf("User with that email already exists"))
	}

	return v.Validate()
}

func (s *UserService) Create(req CreateUserRequest) (*User, error) {
	if err := req.validateRegistrationFields(s); err != nil {
		return nil, err
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
		return User{}, fmt.Errorf(errorMessage)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return User{}, fmt.Errorf(errorMessage)
	}

	return user, nil
}

type UserRepository interface {
	Create(req CreateUserRequest) (*User, error)
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

func (s *UserService) hasUserWithEmail(email string) bool {
	user, _ := s.repository.GetUserByEmail(email)
	return user != (User{})
}

func (s *UserService) GetUserByID(user_id int) (User, error) {
	user, err := s.repository.GetUserByID(user_id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (req *PasswordResetRequest) validateNewPassword() error {
	v := validation.New()
	v.IsEmpty("Password", req.Password)
	v.IsValidPassword("Password", req.Password)

	return v.Validate()
}

func (s *UserService) ResetPassword(req *PasswordResetRequest) error {
	err := req.validateNewPassword()
	if err != nil {
		return err
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

func (req *PasswordResetRequest) setPasswordResetErrors(errors map[string][]error) bool {
	if len(errors) == 0 {
		return true
	}

	for field, err := range errors {
		if field == "Password" {
			req.ErrorPassword = err[0].Error()
		}
	}

	return false
}
