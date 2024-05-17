package app

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

const SessionDurationInHours = 24

type Session struct {
	Id        int
	UserId    int
	SessionId string
	ExpiresAt time.Time
}

type SessionService struct {
	repository SessionRepository
}

func NewSessionService(sr SessionRepository) *SessionService {
	return &SessionService{
		repository: sr,
	}
}

type SessionRepository interface {
	CreateSession(sessionID string, userID int, expiresAt time.Time) (string, error)
	GetSessionByID(sessionID string) (Session, error)
	UpdateStatus(sessionID string) error
}

func (s *SessionService) UpdateStatus(sessionID string) error {
	err := s.repository.UpdateStatus(sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionService) ValidateSession(sessionID string) (Session, error) {
	session, err := s.GetSessionByID(sessionID)
	if err != nil {
		return Session{}, err
	}
	if session.SessionId == "" || session.ExpiresAt.Before(time.Now()) {
		return Session{}, err
	}
	return session, nil
}

func (s *SessionService) generateSessionID() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

func (s *SessionService) CreateSession(userID int, expiresAt time.Time) (string, error) {
	sessionID, err := s.generateSessionID()
	if err != nil {
		return "", err
	}
	sessionID, err = s.repository.CreateSession(sessionID, userID, expiresAt)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (s *SessionService) GetSessionByID(sessionID string) (Session, error) {
	session, err := s.repository.GetSessionByID(sessionID)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}
