package app

import (
	"math/rand"
	"strconv"
	"time"
)

type VerificationService struct {
	repository VerificationRepository
}

func NewVerificationService(vr VerificationRepository) *VerificationService {
	return &VerificationService{
		repository: vr,
	}
}

type VerificationRepository interface {
	CreateVerification(userId int, code string) (string, error)
	UpdateVerified(userId int) error
	GetVerificationCodeByUserID(userId int) (string, error)
}

type Verification struct {
	Id       int
	UserId   int
	Code     string
	Verified bool
}

func (v *VerificationService) SendCode(userId int) (string, error) {
	code := generateVerificationCode()
	println("Verification code: " + code)

	code, err := v.repository.CreateVerification(userId, code)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (v *VerificationService) ValidateCode(userId int, code string) error {
	actualCode, err := v.repository.GetVerificationCodeByUserID(userId)
	if err != nil {
		return &FieldError{Message: "User doesn't exist"}
	}
	if code != actualCode {
		return &FieldError{Message: "Invalid code! Try again."}
	}

	err = v.repository.UpdateVerified(userId)
	return err
}

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(900000) + 100000)
}
