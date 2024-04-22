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

type RegistrationRequest struct {
	ErrorFirstName string
	ErrorLastName  string
	ErrorEmail     string
	ErrorUsername  string
	ErrorPassword  string
	FirstName      string
	LastName       string
	Email          string
	Username       string
	Password       string
}

type PasswordResetRequest struct {
	UserId        int
	Password      string
	ErrorPassword string
}

func (s *UserService) validateRegistrationFields(req *RegistrationRequest) bool {
	v := &validation.Validator{}
	v.Errors = make(map[string][]error)
	v.IsEmpty("First name", req.FirstName)
	v.IsEmpty("Last name", req.LastName)
	v.IsEmpty("Username", req.Username)
	v.IsEmpty("Email", req.Email)
	v.IsEmpty("Password", req.Password)
	v.IsEmail("Email", req.Email)
	v.IsValidPassword("Password", req.Password)
	s.isEmailUsed(v, req.Email)

	return req.setRegistrationFieldErrors(v.Errors)
}

func (s *UserService) Create(req *RegistrationRequest) *User {
	if !s.validateRegistrationFields(req) {
		return nil
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil
	}

	req.Password = hashedPassword

	newUser, err := s.repository.Create(req)
	if err != nil {
		return nil
	}

	return newUser
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
	Create(req *RegistrationRequest) (*User, error)
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

func (s *UserService) isEmailUsed(v *validation.Validator, email string) bool {
	user, _ := s.repository.GetUserByEmail(email)
	if user != (User{}) {
		v.Errors["Email"] = append(v.Errors["Email"], fmt.Errorf("User with that email already exists"))
		return false
	}

	return true
}

func (s *UserService) GetUserByID(user_id int) (User, error) {
	user, err := s.repository.GetUserByID(user_id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (req *PasswordResetRequest) validateNewPassword() bool {
	v := &validation.Validator{}
	v.Errors = make(map[string][]error)
	v.IsEmpty("Password", req.Password)
	v.IsValidPassword("Password", req.Password)

	return req.setPasswordResetErrors(v.Errors)
}

func (s *UserService) ResetPassword(req *PasswordResetRequest) error {
	if !req.validateNewPassword() {
		return fmt.Errorf("")
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

func (req *RegistrationRequest) setRegistrationFieldErrors(errors map[string][]error) bool {
	if len(errors) == 0 {
		return true
	}

	for field, err := range errors {
		if field == "First name" {
			req.ErrorFirstName = err[0].Error()
		}
		if field == "Last name" {
			req.ErrorLastName = err[0].Error()
		}
		if field == "Username" {
			req.ErrorUsername = err[0].Error()
		}
		if field == "Email" {
			req.ErrorEmail = err[0].Error()
		}
		if field == "Password" {
			req.ErrorPassword = err[0].Error()
		}
	}

	return false
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
