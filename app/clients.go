package app

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/OAuth2withJWT/identity-provider/app/validation"
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

func (req *CreateClientRequest) validateClientRegistrationFields(s *ClientService) error {
	v := validation.New()
	v.IsEmpty("Client Name", req.ClientName)
	v.IsEmpty("Scope", req.Scope)
	v.IsEmpty("Redirect URI", req.RedirectURI)
	v.IsValidURI("Redirect URI", req.RedirectURI)
	if s.hasClientWithName(req.ClientName) {
		v.AddError("Client Name", fmt.Errorf("Client with that name already exists"))
	}

	if s.hasClientWithRedirectURI(req.RedirectURI) {
		v.AddError("Redirect URI", fmt.Errorf("Client with that redirect URI already exists"))
	}

	return v.Validate()
}

func (s *ClientService) hasClientWithName(name string) bool {
	client, _ := s.repository.GetClientByName(name)
	return client != (Client{})
}

func (s *ClientService) hasClientWithRedirectURI(redirectURI string) bool {
	client, _ := s.repository.GetClientByRedirectURI(redirectURI)
	return client != (Client{})
}

func (s *ClientService) Create(req CreateClientRequest) (*Client, error) {

	err := req.validateClientRegistrationFields(s)
	if err != nil {
		return nil, err
	}

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
	GetClientByName(name string) (Client, error)
	GetClientByRedirectURI(redirectURI string) (Client, error)
}

func (s *ClientService) GetClientByName(name string) (Client, error) {
	client, err := s.repository.GetClientByName(name)
	if err != nil {
		return Client{}, err
	}

	return client, nil
}

func (s *ClientService) GetClientByRedirectURI(redirectURI string) (Client, error) {
	client, err := s.repository.GetClientByName(redirectURI)
	if err != nil {
		return Client{}, err
	}

	return client, nil
}

func (s *ClientService) generateClientCredentials(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}
