package errors

import (
	"errors"
	"fmt"
)

// AppError represents an application error with additional context.
type AppError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError.
func New(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with additional context.
func Wrap(err error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is checks if the error matches the target error.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Common error codes
const (
	CodeValidation     = "VALIDATION_ERROR"
	CodeNotFound       = "NOT_FOUND"
	CodeConflict       = "CONFLICT"
	CodeInternal       = "INTERNAL_ERROR"
	CodeUnauthorized   = "UNAUTHORIZED"
	CodeForbidden      = "FORBIDDEN"
	CodeBadRequest     = "BAD_REQUEST"
)

// Common errors
var (
	ErrValidation   = New(CodeValidation, "Validation failed")
	ErrNotFound     = New(CodeNotFound, "Resource not found")
	ErrConflict     = New(CodeConflict, "Resource conflict")
	ErrInternal     = New(CodeInternal, "Internal server error")
	ErrUnauthorized = New(CodeUnauthorized, "Unauthorized")
	ErrForbidden    = New(CodeForbidden, "Forbidden")
	ErrBadRequest   = New(CodeBadRequest, "Bad request")
)
