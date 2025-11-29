package validator

import (
	"regexp"
	"strconv"
	"strings"
)

// Validator provides validation utilities.
type Validator struct {
	errors []string
}

// New creates a new Validator.
func New() *Validator {
	return &Validator{
		errors: []string{},
	}
}

// Required validates that a string is not empty.
func (v *Validator) Required(value, field string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.addError(field + " is required")
	}
	return v
}

// MinLength validates that a string has a minimum length.
func (v *Validator) MinLength(value, field string, min int) *Validator {
	if len(value) < min {
		v.addError(field + " must be at least " + strconv.Itoa(min) + " characters")
	}
	return v
}

// MaxLength validates that a string has a maximum length.
func (v *Validator) MaxLength(value, field string, max int) *Validator {
	if len(value) > max {
		v.addError(field + " must be at most " + strconv.Itoa(max) + " characters")
	}
	return v
}

// Email validates that a string is a valid email address.
func (v *Validator) Email(value, field string) *Validator {
	if value == "" {
		return v
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.addError(field + " must be a valid email address")
	}
	return v
}

// Min validates that an integer is at least a minimum value.
func (v *Validator) Min(value int, field string, min int) *Validator {
	if value < min {
		v.addError(field + " must be at least " + strconv.Itoa(min))
	}
	return v
}

// Max validates that an integer is at most a maximum value.
func (v *Validator) Max(value int, field string, max int) *Validator {
	if value > max {
		v.addError(field + " must be at most " + strconv.Itoa(max))
	}
	return v
}

// Positive validates that an integer is positive.
func (v *Validator) Positive(value int, field string) *Validator {
	if value <= 0 {
		v.addError(field + " must be positive")
	}
	return v
}

// NonNegative validates that an integer is non-negative.
func (v *Validator) NonNegative(value int, field string) *Validator {
	if value < 0 {
		v.addError(field + " must be non-negative")
	}
	return v
}

// OneOf validates that a string is one of the allowed values.
func (v *Validator) OneOf(value, field string, allowed []string) *Validator {
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.addError(field + " must be one of: " + strings.Join(allowed, ", "))
	return v
}

// addError adds an error message to the validator.
func (v *Validator) addError(message string) {
	v.errors = append(v.errors, message)
}

// Valid returns true if there are no validation errors.
func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

// Errors returns all validation errors.
func (v *Validator) Errors() []string {
	return v.errors
}

// Error returns the first validation error or nil if valid.
func (v *Validator) Error() error {
	if v.Valid() {
		return nil
	}
	return &ValidationError{errors: v.errors}
}

// ValidationError represents a validation error.
type ValidationError struct {
	errors []string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return strings.Join(e.errors, "; ")
}

// Errors returns all validation errors.
func (e *ValidationError) Errors() []string {
	return e.errors
}
