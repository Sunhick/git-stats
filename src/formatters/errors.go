// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Formatter error types

package formatters

import "fmt"

// FormatterError represents a formatting error
type FormatterError struct {
	Operation string
	Message   string
}

// Error implements the error interface
func (e *FormatterError) Error() string {
	if e.Operation != "" {
		return fmt.Sprintf("formatter error in %s: %s", e.Operation, e.Message)
	}
	return fmt.Sprintf("formatter error: %s", e.Message)
}

// NewFormatterError creates a new formatter error
func NewFormatterError(message string) *FormatterError {
	return &FormatterError{
		Message: message,
	}
}

// NewFormatterOperationError creates a new formatter error for a specific operation
func NewFormatterOperationError(operation, message string) *FormatterError {
	return &FormatterError{
		Operation: operation,
		Message:   message,
	}
}

// IsFormatterError checks if an error is a formatter error
func IsFormatterError(err error) bool {
	_, ok := err.(*FormatterError)
	return ok
}
