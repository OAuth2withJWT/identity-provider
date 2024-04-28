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

func (ur *ClientRepository) Create(req app.CreateClientRequest, credentials app.ClientCredentials) (*app.Client, error) {
	row := ur.db.QueryRow("INSERT INTO clients (client_name, client_id, client_secret, scope,redirect_uri) VALUES ($1, $2, $3, $4, $5) RETURNING  id, client_name, client_id, client_secret, scope,redirect_uri",
		req.ClientName, credentials.ClientID, credentials.ClientSecret, req.Scope, req.RedirectURI)

	client := &app.Client{}
	err := row.Scan(&client.Id, &client.ClientName, &client.ClientId, &client.ClientSecret, &client.Scope, &client.RedirectURI)
	if err != nil {
		return nil, err
	}

	return client, nil
}
