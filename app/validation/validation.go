package validation

type Validator struct {
	Errors map[string]error
}

func (v *Validator) Error() string {
	errorList := ""
	if v.Errors["emptyField"] != nil {
		errorList = errorList + v.Errors["emptyField"].Error()
		return errorList
	}
	if v.Errors["email"] != nil {
		errorList = errorList + v.Errors["email"].Error()
		return errorList
	}
	if v.Errors["password"] != nil {
		errorList = errorList + v.Errors["password"].Error()
		return errorList
	}
	return errorList
}
