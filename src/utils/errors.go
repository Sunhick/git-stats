// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Utils package for error handling utilities

package utils

import (
	"fmt"
	"strings"
)

// ErrorType represents different types of errors that can occur
type ErrorType int

const (
	ErrNotGitRepository ErrorType = iota
	ErrGitNotFound
	ErrInvalidDateFormat
	ErrInvalidAuthor
	ErrRepositoryCorrupted
	ErrPermissionDenied
	ErrInvalidArguments
	ErrCommandTimeout
	ErrInvalidFormat
	ErrFileNotFound
	ErrInsufficientMemory
)

// GitStatsError represents a structured error with type and context
type GitStatsError struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *GitStatsError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *GitStatsError) Unwrap() error {
	return e.Cause
}

// NewGitStatsError creates a new GitStatsError
func NewGitStatsError(errType ErrorType, message string, cause error) *GitStatsError {
	return &GitStatsError{
		Type:    errType,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context information to the error
func (e *GitStatsError) WithContext(key string, value interface{}) *GitStatsError {
	e.Context[key] = value
	return e
}

// GetRecoverySuggestion returns a user-friendly recovery suggestion based on error type
func (e *GitStatsError) GetRecoverySuggestion() string {
	switch e.Type {
	case ErrNotGitRepository:
		return "Make sure you're running this command inside a git repository. Use 'git init' to initialize a new repository."
	case ErrGitNotFound:
		return "Git is not installed or not found in PATH. Please install git and ensure it's accessible from the command line."
	case ErrInvalidDateFormat:
		return "Use date format YYYY-MM-DD or relative dates like '1 week ago', '2024-01-01', etc."
	case ErrInvalidAuthor:
		return "Specify author name or email. Use partial matching like 'john' or 'john@example.com'."
	case ErrRepositoryCorrupted:
		return "The git repository appears to be corrupted. Try running 'git fsck' to check repository integrity."
	case ErrPermissionDenied:
		return "Permission denied. Check file permissions and ensure you have read access to the repository."
	case ErrInvalidArguments:
		return "Invalid command line arguments. Use --help to see available options and usage examples."
	case ErrCommandTimeout:
		return "Command timed out. Try using --limit flag to reduce the scope of analysis for large repositories."
	case ErrInvalidFormat:
		return "Invalid output format. Supported formats are: json, csv, terminal."
	case ErrFileNotFound:
		return "File or directory not found. Check the path and ensure the file exists."
	case ErrInsufficientMemory:
		return "Insufficient memory to process the repository. Try using --limit flag or increase available memory."
	default:
		return "An unexpected error occurred. Please check the error message for more details."
	}
}

// IsRecoverable returns true if the error is recoverable
func (e *GitStatsError) IsRecoverable() bool {
	switch e.Type {
	case ErrNotGitRepository, ErrGitNotFound, ErrInvalidDateFormat, ErrInvalidAuthor, ErrInvalidArguments, ErrInvalidFormat:
		return true
	case ErrRepositoryCorrupted, ErrPermissionDenied, ErrCommandTimeout, ErrFileNotFound, ErrInsufficientMemory:
		return false
	default:
		return false
	}
}

// GetSeverity returns the severity level of the error
func (e *GitStatsError) GetSeverity() string {
	switch e.Type {
	case ErrNotGitRepository, ErrGitNotFound, ErrInvalidArguments:
		return "ERROR"
	case ErrInvalidDateFormat, ErrInvalidAuthor, ErrInvalidFormat:
		return "WARNING"
	case ErrRepositoryCorrupted, ErrPermissionDenied, ErrInsufficientMemory:
		return "CRITICAL"
	case ErrCommandTimeout, ErrFileNotFound:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// FormatUserFriendlyError formats the error for display to users
func (e *GitStatsError) FormatUserFriendlyError() string {
	severity := e.GetSeverity()
	suggestion := e.GetRecoverySuggestion()

	var output strings.Builder
	output.WriteString(fmt.Sprintf("[%s] %s\n", severity, e.Message))

	if suggestion != "" {
		output.WriteString(fmt.Sprintf("💡 Suggestion: %s\n", suggestion))
	}

	// Add context information if available
	if len(e.Context) > 0 {
		output.WriteString("📋 Details:\n")
		for key, value := range e.Context {
			output.WriteString(fmt.Sprintf("  • %s: %v\n", key, value))
		}
	}

	return output.String()
}

// ErrorCollector collects multiple errors and provides summary
type ErrorCollector struct {
	errors   []error
	warnings []error
}

// NewErrorCollector creates a new error collector
func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors:   make([]error, 0),
		warnings: make([]error, 0),
	}
}

// AddError adds an error to the collector
func (ec *ErrorCollector) AddError(err error) {
	if gitErr, ok := err.(*GitStatsError); ok {
		if gitErr.GetSeverity() == "WARNING" {
			ec.warnings = append(ec.warnings, err)
		} else {
			ec.errors = append(ec.errors, err)
		}
	} else {
		ec.errors = append(ec.errors, err)
	}
}

// HasErrors returns true if there are any errors
func (ec *ErrorCollector) HasErrors() bool {
	return len(ec.errors) > 0
}

// HasWarnings returns true if there are any warnings
func (ec *ErrorCollector) HasWarnings() bool {
	return len(ec.warnings) > 0
}

// GetErrorCount returns the number of errors
func (ec *ErrorCollector) GetErrorCount() int {
	return len(ec.errors)
}

// GetWarningCount returns the number of warnings
func (ec *ErrorCollector) GetWarningCount() int {
	return len(ec.warnings)
}

// GetSummary returns a summary of all collected errors and warnings
func (ec *ErrorCollector) GetSummary() string {
	var output strings.Builder

	if len(ec.errors) > 0 {
		output.WriteString(fmt.Sprintf("❌ %d error(s) occurred:\n", len(ec.errors)))
		for i, err := range ec.errors {
			if gitErr, ok := err.(*GitStatsError); ok {
				output.WriteString(fmt.Sprintf("%d. %s\n", i+1, gitErr.FormatUserFriendlyError()))
			} else {
				output.WriteString(fmt.Sprintf("%d. %s\n", i+1, err.Error()))
			}
		}
	}

	if len(ec.warnings) > 0 {
		output.WriteString(fmt.Sprintf("⚠️  %d warning(s):\n", len(ec.warnings)))
		for i, warning := range ec.warnings {
			if gitErr, ok := warning.(*GitStatsError); ok {
				output.WriteString(fmt.Sprintf("%d. %s\n", i+1, gitErr.FormatUserFriendlyError()))
			} else {
				output.WriteString(fmt.Sprintf("%d. %s\n", i+1, warning.Error()))
			}
		}
	}

	return output.String()
}

// WrapError wraps a standard error as a GitStatsError
func WrapError(err error, errType ErrorType, message string) *GitStatsError {
	return NewGitStatsError(errType, message, err)
}
