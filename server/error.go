package server

type Page struct {
	FormFields      map[string]string
	FormErrors      map[string]string
	QueryParameters map[string]string
}
