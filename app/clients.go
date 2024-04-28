package app

import (
	"crypto/rand"
	"encoding/hex"
)

type ClientService struct {
	repository ClientRepository
}

type Client struct {
	Id           int
	ClientName   string
	ClientId     string
	ClientSecret string
	Scope        string
	RedirectURI  string
}

type CreateClientRequest struct {
	ClientName  string
	Scope       string
	RedirectURI string
}

type ClientCredentials struct {
	ClientID     string
	ClientSecret string
}

func (s *ClientService) Create(req CreateClientRequest) (*Client, error) {
	clientID, err := s.generateClientCredentials(16)
	if err != nil {
		return nil, err
	}

	clientSecret, err := s.generateClientCredentials(32)
	if err != nil {
		return nil, err
	}

	credentials := ClientCredentials{clientID, clientSecret}

	newClient, err := s.repository.Create(req, credentials)

	if err != nil {
		return nil, err
	}

	return newClient, nil
}

func NewClientService(cr ClientRepository) *ClientService {
	return &ClientService{
		repository: cr,
	}
}

type ClientRepository interface {
	Create(req CreateClientRequest, credentials ClientCredentials) (*Client, error)
}

func (s *ClientService) generateClientCredentials(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}
