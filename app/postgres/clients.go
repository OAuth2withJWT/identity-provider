package postgres

import (
	"database/sql"

	"github.com/OAuth2withJWT/identity-provider/app"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (cr *ClientRepository) Create(req app.CreateClientRequest, credentials app.ClientCredentials) (*app.Client, error) {
	row := cr.db.QueryRow("INSERT INTO clients (id, name, secret, scope, redirect_uri, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING  id, name, secret, scope, redirect_uri, created_by",
		credentials.Id, req.Name, credentials.Secret, req.Scope, req.RedirectURI, req.CreatedBy)

	client := &app.Client{}
	err := row.Scan(&client.Id, &client.Name, &client.Secret, &client.Scope, &client.RedirectURI, &client.CreatedBy)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (cr *ClientRepository) GetClientByID(id string) (app.Client, error) {
	var client app.Client
	err := cr.db.QueryRow("SELECT * FROM clients WHERE id = $1", id).Scan(&client.Id, &client.Name, &client.Secret, &client.Scope, &client.RedirectURI, &client.CreatedBy)
	if err != nil {
		return app.Client{}, err
	}
	return client, nil
}
