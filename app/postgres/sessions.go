package postgres

import (
	"database/sql"
	"time"

	"github.com/OAuth2withJWT/identity-provider/app"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (sr *SessionRepository) CreateSession(sessionID string, userID int, expiresAt time.Time) (string, error) {
	err := sr.db.QueryRow("INSERT INTO sessions (session_id, user_id, expires_at) VALUES ($1, $2, $3) RETURNING session_id", sessionID, userID, expiresAt).Scan(&sessionID)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (sr *SessionRepository) UpdateFlag(sessionID string) error {
	query := `UPDATE sessions SET status = 'inactive' WHERE session_id = $1`
	_, err := sr.db.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionRepository) GetSessionByID(sessionID string) (app.Session, error) {
	var session app.Session
	err := sr.db.QueryRow("SELECT id, session_id, user_id, expires_at FROM sessions WHERE session_id = $1", sessionID).Scan(&session.Id, &session.SessionId, &session.UserId, &session.ExpiresAt)
	if err != nil {
		return app.Session{}, err
	}
	return session, nil
}
