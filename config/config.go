package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"
)

type RSAConfig struct {
	PrivateKey          *rsa.PrivateKey
	PublicKey           *rsa.PublicKey
	ResourceServer      string
	IdentityProvider    string
	TokenExpirationTime time.Duration
}

func LoadRSAConfig() RSAConfig {
	privateKey, err := LoadPrivateKey("keys/private_key.pem")
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKey, err := LoadPublicKey("keys/public_key.pem")
	if err != nil {
		log.Fatalf("Failed to load public key: %v", err)
	}

	identityProvider := os.Getenv("IDENTITY_PROVIDER")
	if identityProvider == "" {
		log.Fatalf("IDENTITY_PROVIDER environment variable not set")
	}

	resourceServer := os.Getenv("RESOURCE_SERVER")
	if resourceServer == "" {
		log.Fatalf("RESOURCE_SERVER environment variable not set")
	}

	return RSAConfig{
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		ResourceServer:      resourceServer,
		IdentityProvider:    identityProvider,
		TokenExpirationTime: 24 * time.Hour,
	}
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	privKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(privKeyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not of type RSA")
	}

	return rsaPrivateKey, nil
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	pubKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(pubKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type RSA")
	}

	return rsaPublicKey, nil
}
