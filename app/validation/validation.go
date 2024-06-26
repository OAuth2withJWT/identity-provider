package validation

type Validator struct {
	errors map[string][]error
}

func New() *Validator {
	return &Validator{
		errors: make(map[string][]error),
	}
}

type Error struct {
	Errors map[string][]error
}

func (e *Error) Error() string {
	return "validation error"
}

func (v *Validator) Validate() error {
	if len(v.errors) == 0 {
		return nil
	}
	return &Error{Errors: v.errors}
}
