// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - CLI package for input validation

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Validator interface for input validation
type Validator interface {
	ValidateConfig(config *Config) error
	ValidateDateRange(since, until *time.Time) error
	ValidateAuthor(author string) error
	ValidateFormat(format string) error
	ValidateOutputFile(filename string) error
	ValidateRepositoryPath(path string) error
	ValidateLimit(limit int) error
}

// CLIValidator implements the Validator interface
type CLIValidator struct{}

// NewCLIValidator creates a new CLI validator
func NewCLIValidator() *CLIValidator {
	return &CLIValidator{}
}

// ValidateConfig validates the entire configuration
func (v *CLIValidator) ValidateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Skip validation if help is requested
	if config.ShowHelp {
		return nil
	}

	// Validate command
	if err := v.validateCommand(config.Command); err != nil {
		return err
	}

	// Validate date range
	if err := v.ValidateDateRange(config.Since, config.Until); err != nil {
		return err
	}

	// Validate author
	if config.Author != "" {
		if err := v.ValidateAuthor(config.Author); err != nil {
			return err
		}
	}

	// Validate format
	if err := v.ValidateFormat(config.Format); err != nil {
		return err
	}

	// Validate output file
	if config.OutputFile != "" {
		if err := v.ValidateOutputFile(config.OutputFile); err != nil {
			return err
		}
	}

	// Validate repository path
	if err := v.ValidateRepositoryPath(config.RepoPath); err != nil {
		return err
	}

	// Validate limit
	if err := v.ValidateLimit(config.Limit); err != nil {
		return err
	}

	return nil
}

// validateCommand validates the command
func (v *CLIValidator) validateCommand(command string) error {
	validCommands := []string{"contrib", "summary", "contributors", "health"}

	for _, valid := range validCommands {
		if command == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid command '%s'. Valid commands: %s", command, strings.Join(validCommands, ", "))
}

// ValidateDateRange validates that the date range is logical
func (v *CLIValidator) ValidateDateRange(since, until *time.Time) error {
	if since == nil && until == nil {
		return nil // No date range specified is valid
	}

	// Check if the date range is reasonable (not too far in the future)
	now := time.Now()

	if since != nil {
		if since.After(now.AddDate(0, 0, 1)) {
			return fmt.Errorf("since date (%s) cannot be more than 1 day in the future",
				since.Format("2006-01-02"))
		}
	}

	if until != nil {
		if until.After(now.AddDate(0, 0, 1)) {
			return fmt.Errorf("until date (%s) cannot be more than 1 day in the future",
				until.Format("2006-01-02"))
		}
	}

	if since != nil && until != nil {
		if since.After(*until) {
			return fmt.Errorf("since date (%s) cannot be after until date (%s)",
				since.Format("2006-01-02"), until.Format("2006-01-02"))
		}
	}

	return nil
}

// ValidateAuthor validates the author string
func (v *CLIValidator) ValidateAuthor(author string) error {
	if author == "" {
		return fmt.Errorf("author cannot be empty")
	}

	// Trim whitespace
	author = strings.TrimSpace(author)
	if author == "" {
		return fmt.Errorf("author cannot be only whitespace")
	}

	// Check length (reasonable limits)
	if len(author) > 100 {
		return fmt.Errorf("author string too long (max 100 characters)")
	}

	// Check for potentially dangerous characters (basic security)
	if strings.ContainsAny(author, "\n\r\t") {
		return fmt.Errorf("author cannot contain newline or tab characters")
	}

	// Validate email format if it looks like an email
	if strings.Contains(author, "@") {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(author) {
			return fmt.Errorf("invalid email format: %s", author)
		}
	}

	return nil
}

// ValidateFormat validates the output format
func (v *CLIValidator) ValidateFormat(format string) error {
	if format == "" {
		return fmt.Errorf("format cannot be empty")
	}

	validFormats := []string{"terminal", "json", "csv"}
	format = strings.ToLower(strings.TrimSpace(format))

	for _, valid := range validFormats {
		if format == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid format '%s'. Valid formats: %s", format, strings.Join(validFormats, ", "))
}

// ValidateOutputFile validates the output file path
func (v *CLIValidator) ValidateOutputFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("output file cannot be empty")
	}

	// Trim whitespace
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return fmt.Errorf("output file cannot be only whitespace")
	}

	// Check for dangerous characters
	if strings.ContainsAny(filename, "\n\r\t") {
		return fmt.Errorf("output file path cannot contain newline or tab characters")
	}

	// Validate the directory exists or can be created
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("output directory does not exist: %s", dir)
		}
	}

	// Check if file already exists and is writable
	if _, err := os.Stat(filename); err == nil {
		// File exists, check if it's writable
		file, err := os.OpenFile(filename, os.O_WRONLY, 0)
		if err != nil {
			return fmt.Errorf("output file is not writable: %s", filename)
		}
		file.Close()
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking output file: %w", err)
	}

	return nil
}

// ValidateRepositoryPath validates the repository path
func (v *CLIValidator) ValidateRepositoryPath(path string) error {
	if path == "" {
		return fmt.Errorf("repository path cannot be empty")
	}

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("repository path does not exist: %s", path)
	}

	// Check if it's a directory
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error accessing repository path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("repository path is not a directory: %s", path)
	}

	// Check if it's a git repository (look for .git directory)
	gitPath := filepath.Join(path, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository (no .git directory found): %s", path)
	}

	return nil
}

// ValidateLimit validates the commit limit
func (v *CLIValidator) ValidateLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("limit must be greater than 0, got %d", limit)
	}

	if limit > 1000000 {
		return fmt.Errorf("limit too large (max 1,000,000), got %d", limit)
	}

	return nil
}

// ValidationError represents a validation error with additional context
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error for field '%s' (value: %v): %s", e.Field, e.Value, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}
