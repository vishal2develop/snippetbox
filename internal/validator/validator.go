package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// create a validator struct to hold form field validation errors
type Validator struct {
	FieldErrors map[string]string
}

// valid() return true if no errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Utility: AddFieldError() adds an error message to the FieldErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) AddFieldError(key, msg string) {
	// init the map first, if it isnt already

	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	// add the error message to the key if it doesnt already exist
	// check if key exists; if it does exist, dont add the error. If it doesnt exist, add the error
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = msg
	}

}

// CheckField() adds an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, msg string) {
	if !ok {
		v.AddFieldError(key, msg)
	}
}

// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue() returns true if a value is in a list of specific permitted
// values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
