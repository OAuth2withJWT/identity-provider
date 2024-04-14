package postgres

import (
	"database/sql"
)

type VerificationRepository struct {
	db *sql.DB
}

func NewVerificationRepository(db *sql.DB) *VerificationRepository {
	return &VerificationRepository{
		db: db,
	}
}

func (sr *VerificationRepository) CreateVerification(userId int, code string) (string, error) {
	err := sr.db.QueryRow("INSERT INTO verifications (user_id, code, verified) VALUES ($1, $2, $3) RETURNING code", userId, code, false).Scan(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (sr *VerificationRepository) UpdateVerified(userId int) error {
	query := `UPDATE verifications SET verified = true WHERE user_id = $1`
	_, err := sr.db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (sr *VerificationRepository) GetVerificationCodeByUserID(userId int) (string, error) {
	var verificationCode string
	err := sr.db.QueryRow("SELECT code FROM verifications WHERE user_id = $1", userId).Scan(&verificationCode)
	if err != nil {
		return "", err
	}
	return verificationCode, nil
}
