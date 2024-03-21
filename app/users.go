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

func (s *UserService) Create(req CreateUserRequest) (*User, error) {
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

func (s *UserService) Authenticate(username, password string) (int, error) {
	userID, err := s.repository.Authenticate(username, password)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

type UserRepository interface {
	Create(CreateUserRequest) (*User, error)
	Authenticate(username, password string) (int, error)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
