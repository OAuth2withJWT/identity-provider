package app

type Application struct {
	UserService    *UserService
	SessionService *SessionService
}
