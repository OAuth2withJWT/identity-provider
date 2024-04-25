package server

type Page struct {
	// vidi da li treba:)
	Success    bool
	FormFields map[string]string
	FormErrors map[string]string
}
