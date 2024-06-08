package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX uses the regexp.MustCompile() function to parse a regular expression pattern
// for sanity checking the format of an email address. This returns a pointer to
// a 'compiled' regexp.Regexp type, or panics in the event of an error. Parsing
// this pattern once at startup and storing the compiled *regexp.Regexp in a
// variable is more performant than re-parsing the pattern each time we need it.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	//for errors not related to any particular field. Like incorrect email or password error during login
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField adds an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank returns if string is empty or not
func NotBlank(field string) bool {
	return strings.TrimSpace(field) != ""
}

// MaxChars returns true if a value contains no more than n characters.
func MaxChars(field string, limit int) bool {
	return utf8.RuneCountInString(field) <= limit
}

func MinChars(field string, min int) bool {
	return utf8.RuneCountInString(field) >= min
}

// PermittedInt returns true if a value is in a list of permitted integers.
// replace PermittedInt with PermittedValue to make it more generic and work with different types
// of datatypes
func PermittedValue[T comparable](num T, permittedVal ...T) bool {
	for i := range permittedVal {
		if num == permittedVal[i] {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
