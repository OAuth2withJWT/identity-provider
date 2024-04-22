package validation

type Validator struct {
	Errors map[string][]error
}

func (v *Validator) Error() map[string][]error {
	return v.Errors
}
