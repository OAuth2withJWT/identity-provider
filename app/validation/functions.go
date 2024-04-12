package validation

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode"
)

func (v *Validator) IsEmpty(value string) {
	if strings.TrimSpace(value) == "" && v.Errors["emptyField"] == nil {
		v.Errors["emptyField"] = fmt.Errorf("Fields cannot be empty")
	}
}

func (v *Validator) IsEmail(field string, value string) {
	_, err := mail.ParseAddress(value)
	if err != nil {
		v.Errors[field] = fmt.Errorf("%s is not valid", field)
	}
}

func (v *Validator) IsValidPassword(field string, password string) {
	var errors []string

	rules := map[string]func(string) bool{
		"at least 8 characters": func(s string) bool { return len(s) >= 8 },
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
		v.Errors[field] = fmt.Errorf("Password must contain " + strings.Join(errors, ", "))
	}
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
