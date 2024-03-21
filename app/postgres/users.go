package postgres

import (
	"database/sql"

	"github.com/OAuth2withJWT/identity-provider/app"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Create(req app.CreateUserRequest) (*app.User, error) {
	row := ur.db.QueryRow("INSERT INTO users (first_name, last_name, email, username, password) VALUES ($1, $2, $3, $4, $5) RETURNING first_name, last_name, email, username, password",
		req.FirstName, req.LastName, req.Email, req.Username, req.Password)

	user := &app.User{}
	err := row.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) Authenticate(username, password string) (int, error) {
	var userID int
	err := ur.db.QueryRow("SELECT id FROM users WHERE username = $1 AND password = $2", username, password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
