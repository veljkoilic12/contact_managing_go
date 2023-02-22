package validator

import "regexp"

var (
	PhoneRX = regexp.MustCompile("^(\\+)(3816)([0-9]){6,9}$")
)

// Validator contains a map of validation errors
type Validator struct {
	Errors map[string]string
}

// Returns a new Validator instance with an empty errors map
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Returns true if the errors map does not contain any entries
func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map, as long as the key does not exist
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Matches returns true if a string value matches a specific regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
