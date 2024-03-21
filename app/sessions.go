package app

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
)

type Session struct {
	UserId    int
	SessionId int
	ExpiresAt time.Time
}

type SessionService struct {
	db *sql.DB
}

func NewSessionService(db *sql.DB) *SessionService {
	return &SessionService{
		db: db,
	}
}

func (s *SessionService) SaveSession(userID int, expiresAt time.Time) (int, error) {
	var sessionID int
	err := s.db.QueryRow("INSERT INTO sessions (user_id, expires_at) VALUES ($1, $2) RETURNING session_id", userID, expiresAt).Scan(&sessionID)
	if err != nil {
		return 0, err
	}
	return sessionID, nil
}

func (s *SessionService) SetSessionCookie(w http.ResponseWriter, sessionID int) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   strconv.Itoa(sessionID),
		Expires: time.Now().Add(24 * time.Hour),
	})
}

func (s *SessionService) GetSessionIDFromCookie(r *http.Request) int {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0
	}
	sessionID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0
	}
	return sessionID
}
