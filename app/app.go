package app

type Application struct {
	UserService              *UserService
	SessionService           *SessionService
	VerificationService      *VerificationService
	ClientService            *ClientService
	AuthService              *AuthService
	AuthorizationCodeService *AuthorizationCodeService
}
