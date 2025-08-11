// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Model validation errors

package models

import (
	"fmt"
)

// ValidationError represents a data validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

// NewFieldValidationError creates a new validation error for a specific field
func NewFieldValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}
