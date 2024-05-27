package app

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/OAuth2withJWT/identity-provider/app/validation"
)

type ClientService struct {
	repository ClientRepository
}

func NewClientService(cr ClientRepository) *ClientService {
	return &ClientService{
		repository: cr,
	}
}

type ClientRepository interface {
	Create(req CreateClientRequest, credentials ClientCredentials) (*Client, error)
	GetClientByID(id string) (Client, error)
}

type Client struct {
	Id          string
	Name        string
	Secret      string
	Scope       string
	RedirectURI string
	CreatedBy   int
}

type CreateClientRequest struct {
	Name        string
	Scope       string
	RedirectURI string
	CreatedBy   int
}

type ClientCredentials struct {
	Id     string
	Secret string
}

func (req *CreateClientRequest) validate() error {
	v := validation.New()
	v.IsEmpty("Client Name", req.Name)
	v.IsEmpty("Scope", req.Scope)
	v.IsEmpty("Redirect URI", req.RedirectURI)
	v.IsValidURI("Redirect URI", req.RedirectURI)
	return v.Validate()
}

func GenerateRandomBytes(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}

func (s *ClientService) Create(req CreateClientRequest) (*Client, error) {

	err := req.validate()
	if err != nil {
		return nil, err
	}

	const clientIdLength = 16

	clientID, err := GenerateRandomBytes(clientIdLength)
	if err != nil {
		return nil, err
	}

	const clientSecretLength = 32

	clientSecret, err := GenerateRandomBytes(clientSecretLength)
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

func (s *ClientService) GetClientByID(id string) (Client, error) {
	client, err := s.repository.GetClientByID(id)
	if err != nil {
		return Client{}, err
	}

	return client, nil
}
