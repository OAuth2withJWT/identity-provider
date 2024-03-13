package app

import "golang.org/x/crypto/bcrypt"

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

func (s *UserService) Create(req CreateUserRequest) (*User, error) {

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
