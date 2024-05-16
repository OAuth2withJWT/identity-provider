package validation

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"unicode"
)

func (v *Validator) IsValidURI(field string, uri string) {
	parsedURI, err := url.ParseRequestURI(uri)

	if err != nil {
		v.AddError(field, fmt.Errorf("%s is not a valid URI", field))
		return
	}

	if parsedURI.Scheme != "https" {
		v.AddError(field, fmt.Errorf("%s must use https scheme", field))
		return
	}

}

func (v *Validator) IsEmpty(field string, value string) {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, fmt.Errorf("%s cannot be empty", field))
	}
}

func (v *Validator) IsEmail(field string, value string) {
	_, err := mail.ParseAddress(value)
	if err != nil {
		v.errors[field] = append(v.errors[field], fmt.Errorf("%s is not valid", field))
	}
}

func (v *Validator) IsValidPassword(field string, password string) {
	var errors []string

	rules := map[string]func(string) bool{
		"8 characters":          func(s string) bool { return len(s) >= 8 },
		"one uppercase letter":  func(s string) bool { return containsType(s, unicode.IsUpper) },
		"one lowercase letter":  func(s string) bool { return containsType(s, unicode.IsLower) },
		"one digit":             func(s string) bool { return containsType(s, unicode.IsDigit) },
		"one special character": func(s string) bool { return containsSpecialChar(s) },
	}

	for rule, isValid := range rules {
		if !isValid(password) {
			errors = append(errors, rule)
		}
	}

	if len(errors) > 0 {
		v.AddError(field, fmt.Errorf("Password must contain at least: "+strings.Join(errors, ", ")))
	}
}

func (v *Validator) AddError(field string, err error) {
	if v.errors == nil {
		v.errors = make(map[string][]error)
	}
	v.errors[field] = append(v.errors[field], err)
}

func containsType(s string, check func(rune) bool) bool {
	for _, char := range s {
		if check(char) {
			return true
		}
	}
	return false
}

func containsSpecialChar(s string) bool {
	for _, char := range s {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			return true
		}
	}
	return false
}
