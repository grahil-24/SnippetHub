package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField() adds an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// returns if string is empty or not
func NotBlank(message string) bool {
	return strings.TrimSpace(message) != ""
}

// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(message string, limit int) bool {
	return utf8.RuneCountInString(message) <= limit
}

// PermittedInt() returns true if a value is in a list of permitted integers.
func PermittedInt(num int, permittedVal ...int) bool {
	for i := range permittedVal {
		if num == permittedVal[i] {
			return true
		}
	}
	return false
}
