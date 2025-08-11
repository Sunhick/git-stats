// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for CLI parser interfaces

package cli_test

import (
	"git-stats/cli"
	"testing"
	"time"
)

// MockParser implements the Parser interface for testing
type MockParser struct {
	config *cli.Config
	err    error
}

func (m *MockParser) Parse(args []string) (*cli.Config, error) {
	return m.config, m.err
}

func (m *MockParser) PrintUsage() {
	// Mock implementation
}

func (m *MockParser) PrintHelp() {
	// Mock implementation
}

func TestParserInterface(t *testing.T) {
	// Test that MockParser implements Parser interface
	var _ cli.Parser = &MockParser{}

	testTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	parser := &MockParser{
		config: &cli.Config{
			Command:      "contrib",
			Since:        &testTime,
			Until:        &testTime,
			Author:       "test@example.com",
			Format:       "json",
			OutputFile:   "output.json",
			RepoPath:     "/test/repo",
			ShowProgress: true,
			Limit:        1000,
			GUIMode:      false,
		},
		err: nil,
	}

	config, err := parser.Parse([]string{"--contrib", "--author", "test@example.com"})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if config.Command != "contrib" {
		t.Errorf("Expected command 'contrib', got '%s'", config.Command)
	}

	if config.Author != "test@example.com" {
		t.Errorf("Expected author 'test@example.com', got '%s'", config.Author)
	}

	if config.Format != "json" {
		t.Errorf("Expected format 'json', got '%s'", config.Format)
	}

	if !config.ShowProgress {
		t.Error("Expected ShowProgress to be true")
	}

	if config.Limit != 1000 {
		t.Errorf("Expected limit 1000, got %d", config.Limit)
	}
}

func TestConfigStruct(t *testing.T) {
	since := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	config := &cli.Config{
		Command:      "summary",
		Since:        &since,
		Until:        &until,
		Author:       "john@example.com",
		Format:       "csv",
		OutputFile:   "report.csv",
		RepoPath:     "/home/user/project",
		ShowProgress: false,
		Limit:        500,
		GUIMode:      true,
	}

	if config.Command != "summary" {
		t.Errorf("Expected command 'summary', got '%s'", config.Command)
	}

	if config.Since == nil || !config.Since.Equal(since) {
		t.Errorf("Expected since date %v, got %v", since, config.Since)
	}

	if config.Until == nil || !config.Until.Equal(until) {
		t.Errorf("Expected until date %v, got %v", until, config.Until)
	}

	if config.Author != "john@example.com" {
		t.Errorf("Expected author 'john@example.com', got '%s'", config.Author)
	}

	if config.Format != "csv" {
		t.Errorf("Expected format 'csv', got '%s'", config.Format)
	}

	if config.OutputFile != "report.csv" {
		t.Errorf("Expected output file 'report.csv', got '%s'", config.OutputFile)
	}

	if config.RepoPath != "/home/user/project" {
		t.Errorf("Expected repo path '/home/user/project', got '%s'", config.RepoPath)
	}

	if config.ShowProgress {
		t.Error("Expected ShowProgress to be false")
	}

	if config.Limit != 500 {
		t.Errorf("Expected limit 500, got %d", config.Limit)
	}

	if !config.GUIMode {
		t.Error("Expected GUIMode to be true")
	}
}

// MockValidator implements the Validator interface for testing
type MockValidator struct {
	validationError error
}

func (m *MockValidator) ValidateConfig(config *cli.Config) error {
	return m.validationError
}

func (m *MockValidator) ValidateDateRange(since, until *time.Time) error {
	if since != nil && until != nil && since.After(*until) {
		return &MockValidationError{"since date cannot be after until date"}
	}
	return m.validationError
}

func (m *MockValidator) ValidateAuthor(author string) error {
	if author == "" {
		return &MockValidationError{"author cannot be empty"}
	}
	return m.validationError
}

func (m *MockValidator) ValidateFormat(format string) error {
	validFormats := []string{"json", "csv", "terminal"}
	for _, valid := range validFormats {
		if format == valid {
			return m.validationError
		}
	}
	return &MockValidationError{"invalid format"}
}

func (m *MockValidator) ValidateOutputFile(filename string) error {
	if filename == "" {
		return &MockValidationError{"output file cannot be empty"}
	}
	return m.validationError
}

type MockValidationError struct {
	message string
}

func (e *MockValidationError) Error() string {
	return e.message
}

func TestValidatorInterface(t *testing.T) {
	// Test that MockValidator implements Validator interface
	var _ cli.Validator = &MockValidator{}

	validator := &MockValidator{validationError: nil}

	// Test ValidateConfig
	config := &cli.Config{
		Command: "contrib",
		Format:  "json",
		Author:  "test@example.com",
	}

	err := validator.ValidateConfig(config)
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}

	// Test ValidateDateRange with valid range
	since := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	err = validator.ValidateDateRange(&since, &until)
	if err != nil {
		t.Errorf("Unexpected validation error for valid date range: %v", err)
	}

	// Test ValidateDateRange with invalid range
	err = validator.ValidateDateRange(&until, &since)
	if err == nil {
		t.Error("Expected validation error for invalid date range")
	}

	// Test ValidateAuthor with valid author
	err = validator.ValidateAuthor("test@example.com")
	if err != nil {
		t.Errorf("Unexpected validation error for valid author: %v", err)
	}

	// Test ValidateAuthor with empty author
	err = validator.ValidateAuthor("")
	if err == nil {
		t.Error("Expected validation error for empty author")
	}

	// Test ValidateFormat with valid format
	err = validator.ValidateFormat("json")
	if err != nil {
		t.Errorf("Unexpected validation error for valid format: %v", err)
	}

	// Test ValidateFormat with invalid format
	err = validator.ValidateFormat("xml")
	if err == nil {
		t.Error("Expected validation error for invalid format")
	}
}
