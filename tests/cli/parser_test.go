// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for CLI parser and validator

package cli_test

import (
	"fmt"
	"git-stats/cli"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Test CLIParser implementation
func TestCLIParser_Parse_BasicCommands(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"default command", []string{}, "contrib"},
		{"contrib command", []string{"-contrib"}, "contrib"},
		{"summary command", []string{"-summary"}, "summary"},
		{"contributors command", []string{"-contributors"}, "contributors"},
		{"health command", []string{"-health"}, "health"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary git repository for testing
			tempDir := createTempGitRepo(t)
			defer os.RemoveAll(tempDir)

			args := append(tt.args, tempDir)
			config, err := parser.Parse(args)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if config.Command != tt.expected {
				t.Errorf("Expected command '%s', got '%s'", tt.expected, config.Command)
			}
		})
	}
}

func TestCLIParser_Parse_GUIMode(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	config, err := parser.Parse([]string{"-gui", tempDir})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !config.GUIMode {
		t.Error("Expected GUIMode to be true")
	}

	if config.Command != "contrib" {
		t.Errorf("Expected default command 'contrib' for GUI mode, got '%s'", config.Command)
	}
}

func TestCLIParser_Parse_DateFlags(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name      string
		args      []string
		checkFunc func(*testing.T, *cli.Config)
	}{
		{
			name: "absolute dates",
			args: []string{"-since", "2024-01-01", "-until", "2024-01-31", tempDir},
			checkFunc: func(t *testing.T, config *cli.Config) {
				if config.Since == nil {
					t.Error("Expected since date to be set")
				} else {
					expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
					if !config.Since.Equal(expected) {
						t.Errorf("Expected since date %v, got %v", expected, *config.Since)
					}
				}

				if config.Until == nil {
					t.Error("Expected until date to be set")
				} else {
					expected := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
					if !config.Until.Equal(expected) {
						t.Errorf("Expected until date %v, got %v", expected, *config.Until)
					}
				}
			},
		},
		{
			name: "relative dates",
			args: []string{"-since", "1 week ago", tempDir},
			checkFunc: func(t *testing.T, config *cli.Config) {
				if config.Since == nil {
					t.Error("Expected since date to be set")
				} else {
					// Check that it's approximately 1 week ago
					expected := time.Now().AddDate(0, 0, -7)
					diff := config.Since.Sub(expected)
					if diff > time.Hour || diff < -time.Hour {
						t.Errorf("Expected since date to be approximately 1 week ago, got %v", *config.Since)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := parser.Parse(tt.args)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			tt.checkFunc(t, config)
		})
	}
}

func TestCLIParser_Parse_AuthorFlag(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	config, err := parser.Parse([]string{"-author", "john@example.com", tempDir})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if config.Author != "john@example.com" {
		t.Errorf("Expected author 'john@example.com', got '%s'", config.Author)
	}
}

func TestCLIParser_Parse_FormatAndOutput(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	// Create a temporary output file
	outputFile := filepath.Join(tempDir, "output.json")

	config, err := parser.Parse([]string{"-format", "json", "-output", outputFile, tempDir})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if config.Format != "json" {
		t.Errorf("Expected format 'json', got '%s'", config.Format)
	}

	if config.OutputFile != outputFile {
		t.Errorf("Expected output file '%s', got '%s'", outputFile, config.OutputFile)
	}
}

func TestCLIParser_Parse_ProgressAndLimit(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	config, err := parser.Parse([]string{"-progress", "-limit", "5000", tempDir})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !config.ShowProgress {
		t.Error("Expected ShowProgress to be true")
	}

	if config.Limit != 5000 {
		t.Errorf("Expected limit 5000, got %d", config.Limit)
	}
}

func TestCLIParser_Parse_Help(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tests := []struct {
		name string
		args []string
	}{
		{"help flag", []string{"-help"}},
		{"h flag", []string{"-h"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := parser.Parse(tt.args)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !config.ShowHelp {
				t.Error("Expected ShowHelp to be true")
			}
		})
	}
}

func TestCLIParser_Parse_MultipleCommands(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	_, err := parser.Parse([]string{"-contrib", "-summary", tempDir})

	if err == nil {
		t.Error("Expected error for multiple commands")
	}

	if !strings.Contains(err.Error(), "only one command can be specified") {
		t.Errorf("Expected error about multiple commands, got: %v", err)
	}
}

func TestCLIParser_Parse_InvalidDate(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	_, err := parser.Parse([]string{"-since", "invalid-date", tempDir})

	if err == nil {
		t.Error("Expected error for invalid date")
	}

	if !strings.Contains(err.Error(), "invalid since date") {
		t.Errorf("Expected error about invalid since date, got: %v", err)
	}
}

// Test CLIValidator implementation
func TestCLIValidator_ValidateConfig(t *testing.T) {
	validator := cli.NewCLIValidator()

	tests := []struct {
		name      string
		config    *cli.Config
		expectErr bool
	}{
		{
			name: "valid config",
			config: func() *cli.Config {
				tempDir := createTempGitRepo(t)
				return &cli.Config{
					Command:  "contrib",
					Format:   "terminal",
					RepoPath: tempDir,
					Limit:    1000,
				}
			}(),
			expectErr: false,
		},
		{
			name:      "nil config",
			config:    nil,
			expectErr: true,
		},
		{
			name: "help config",
			config: &cli.Config{
				ShowHelp: true,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up temp directory if it was created for this test
			if tt.config != nil && strings.Contains(tt.config.RepoPath, "git-stats-test-") {
				defer os.RemoveAll(tt.config.RepoPath)
			}

			err := validator.ValidateConfig(tt.config)

			if tt.expectErr && err == nil {
				t.Error("Expected validation error")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestCLIValidator_ValidateDateRange(t *testing.T) {
	validator := cli.NewCLIValidator()

	tests := []struct {
		name      string
		since     *time.Time
		until     *time.Time
		expectErr bool
	}{
		{
			name:      "no dates",
			since:     nil,
			until:     nil,
			expectErr: false,
		},
		{
			name:      "valid range",
			since:     timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			until:     timePtr(time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)),
			expectErr: false,
		},
		{
			name:      "invalid range",
			since:     timePtr(time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)),
			until:     timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			expectErr: true,
		},
		{
			name:      "future date",
			since:     timePtr(time.Now().AddDate(0, 0, 2)),
			until:     nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDateRange(tt.since, tt.until)

			if tt.expectErr && err == nil {
				t.Error("Expected validation error")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestCLIValidator_ValidateAuthor(t *testing.T) {
	validator := cli.NewCLIValidator()

	tests := []struct {
		name      string
		author    string
		expectErr bool
	}{
		{"valid name", "John Doe", false},
		{"valid email", "john@example.com", false},
		{"empty author", "", true},
		{"whitespace only", "   ", true},
		{"too long", strings.Repeat("a", 101), true},
		{"invalid email", "invalid@", true},
		{"newline character", "john\ndoe", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAuthor(tt.author)

			if tt.expectErr && err == nil {
				t.Error("Expected validation error")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestCLIValidator_ValidateFormat(t *testing.T) {
	validator := cli.NewCLIValidator()

	tests := []struct {
		name      string
		format    string
		expectErr bool
	}{
		{"terminal format", "terminal", false},
		{"json format", "json", false},
		{"csv format", "csv", false},
		{"uppercase format", "JSON", false},
		{"empty format", "", true},
		{"invalid format", "xml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateFormat(tt.format)

			if tt.expectErr && err == nil {
				t.Error("Expected validation error")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestCLIValidator_ValidateLimit(t *testing.T) {
	validator := cli.NewCLIValidator()

	tests := []struct {
		name      string
		limit     int
		expectErr bool
	}{
		{"valid limit", 1000, false},
		{"minimum limit", 1, false},
		{"maximum limit", 1000000, false},
		{"zero limit", 0, true},
		{"negative limit", -1, true},
		{"too large limit", 1000001, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateLimit(tt.limit)

			if tt.expectErr && err == nil {
				t.Error("Expected validation error")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestCLIValidator_ValidateRepositoryPath(t *testing.T) {
	validator := cli.NewCLIValidator()

	// Create a temporary git repository
	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	// Create a non-git directory
	nonGitDir := filepath.Join(os.TempDir(), "non-git-dir")
	os.MkdirAll(nonGitDir, 0755)
	defer os.RemoveAll(nonGitDir)

	tests := []struct {
		name      string
		path      string
		expectErr bool
	}{
		{"valid git repo", tempDir, false},
		{"empty path", "", true},
		{"non-existent path", "/non/existent/path", true},
		{"non-git directory", nonGitDir, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateRepositoryPath(tt.path)

			if tt.expectErr && err == nil {
				t.Error("Expected validation error")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

// Test comprehensive CLI combinations
func TestCLIParser_Parse_ComprehensiveCombinations(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		args     []string
		validate func(*testing.T, *cli.Config, error)
	}{
		{
			name: "GUI with contrib command",
			args: []string{"-gui", "-contrib", tempDir},
			validate: func(t *testing.T, config *cli.Config, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !config.GUIMode {
					t.Error("Expected GUIMode to be true")
				}
				if config.Command != "contrib" {
					t.Errorf("Expected command 'contrib', got '%s'", config.Command)
				}
			},
		},
		{
			name: "Complex filtering combination",
			args: []string{"-summary", "-since", "2024-01-01", "-until", "2024-12-31", "-author", "john@example.com", "-format", "json", "-progress", tempDir},
			validate: func(t *testing.T, config *cli.Config, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if config.Command != "summary" {
					t.Errorf("Expected command 'summary', got '%s'", config.Command)
				}
				if config.Since == nil {
					t.Error("Expected since date to be set")
				}
				if config.Until == nil {
					t.Error("Expected until date to be set")
				}
				if config.Author != "john@example.com" {
					t.Errorf("Expected author 'john@example.com', got '%s'", config.Author)
				}
				if config.Format != "json" {
					t.Errorf("Expected format 'json', got '%s'", config.Format)
				}
				if !config.ShowProgress {
					t.Error("Expected ShowProgress to be true")
				}
			},
		},
		{
			name: "All flags combination",
			args: []string{"-contributors", "-since", "1 month ago", "-author", "jane", "-format", "csv", "-output", "/tmp/test.csv", "-limit", "5000", "-progress", tempDir},
			validate: func(t *testing.T, config *cli.Config, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if config.Command != "contributors" {
					t.Errorf("Expected command 'contributors', got '%s'", config.Command)
				}
				if config.Since == nil {
					t.Error("Expected since date to be set")
				}
				if config.Author != "jane" {
					t.Errorf("Expected author 'jane', got '%s'", config.Author)
				}
				if config.Format != "csv" {
					t.Errorf("Expected format 'csv', got '%s'", config.Format)
				}
				if config.OutputFile != "/tmp/test.csv" {
					t.Errorf("Expected output file '/tmp/test.csv', got '%s'", config.OutputFile)
				}
				if config.Limit != 5000 {
					t.Errorf("Expected limit 5000, got %d", config.Limit)
				}
				if !config.ShowProgress {
					t.Error("Expected ShowProgress to be true")
				}
			},
		},
		{
			name: "Health command with relative dates",
			args: []string{"-health", "-since", "yesterday", "-until", "today", tempDir},
			validate: func(t *testing.T, config *cli.Config, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if config.Command != "health" {
					t.Errorf("Expected command 'health', got '%s'", config.Command)
				}
				if config.Since == nil {
					t.Error("Expected since date to be set")
				}
				if config.Until == nil {
					t.Error("Expected until date to be set")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := parser.Parse(tt.args)
			tt.validate(t, config, err)
		})
	}
}

// Test edge cases and error conditions
func TestCLIParser_Parse_EdgeCases(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorCheck  func(error) bool
	}{
		{
			name:        "Invalid date format",
			args:        []string{"-since", "not-a-date", tempDir},
			expectError: true,
			errorCheck:  func(err error) bool { return strings.Contains(err.Error(), "invalid since date") },
		},
		{
			name:        "Invalid format",
			args:        []string{"-format", "xml", tempDir},
			expectError: true,
			errorCheck:  func(err error) bool { return strings.Contains(err.Error(), "invalid format") },
		},
		{
			name:        "Invalid limit",
			args:        []string{"-limit", "0", tempDir},
			expectError: true,
			errorCheck:  func(err error) bool { return strings.Contains(err.Error(), "limit must be greater than 0") },
		},
		{
			name:        "Date range backwards",
			args:        []string{"-since", "2024-12-31", "-until", "2024-01-01", tempDir},
			expectError: true,
			errorCheck:  func(err error) bool { return strings.Contains(err.Error(), "since date") && strings.Contains(err.Error(), "cannot be after until date") },
		},

		{
			name:        "Non-existent repository",
			args:        []string{"/non/existent/repo"},
			expectError: true,
			errorCheck:  func(err error) bool { return strings.Contains(err.Error(), "repository path does not exist") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.args)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectError && err != nil && tt.errorCheck != nil {
				if !tt.errorCheck(err) {
					t.Errorf("Error message doesn't match expected pattern: %v", err)
				}
			}
		})
	}
}

// Test help system functionality
func TestCLIParser_HelpSystem(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	// Test that help can be called without validation errors
	config, err := parser.Parse([]string{"-help"})
	if err != nil {
		t.Errorf("Help should not produce errors: %v", err)
	}
	if !config.ShowHelp {
		t.Error("Expected ShowHelp to be true")
	}

	// Test short help flag
	config, err = parser.Parse([]string{"-h"})
	if err != nil {
		t.Errorf("Short help should not produce errors: %v", err)
	}
	if !config.ShowHelp {
		t.Error("Expected ShowHelp to be true for short flag")
	}
}

// Test error suggestion functionality
func TestCLIParser_PrintErrorWithSuggestion(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	// Test that the method exists and can be called
	// We can't easily test the output without capturing stderr,
	// but we can ensure it doesn't panic
	testErrors := []error{
		fmt.Errorf("not a git repository"),
		fmt.Errorf("invalid since date '2024-13-01'"),
		fmt.Errorf("invalid format 'xml'"),
		fmt.Errorf("only one command can be specified at a time"),
		fmt.Errorf("limit must be greater than 0"),
		fmt.Errorf("since date (2024-12-31) cannot be after until date (2024-01-01)"),
		fmt.Errorf("repository path does not exist: /non/existent"),
		fmt.Errorf("some other error"),
	}

	for _, err := range testErrors {
		// This should not panic
		parser.PrintErrorWithSuggestion(err)
	}
}

// Test date parsing edge cases
func TestCLIParser_DateParsing(t *testing.T) {
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		dateStr  string
		expectOk bool
	}{
		{"ISO date", "2024-01-15", true},
		{"ISO datetime", "2024-01-15 14:30:00", true},
		{"ISO with timezone", "2024-01-15T14:30:00Z", true},
		{"US format", "01/15/2024", true},
		{"EU format", "15-01-2024", true},
		{"Today", "today", true},
		{"Yesterday", "yesterday", true},
		{"1 week ago", "1 week ago", true},
		{"2 months ago", "2 months ago", true},
		{"1 year ago", "1 year ago", true},
		{"Invalid format", "not-a-date", false},
		{"Invalid relative", "5 centuries ago", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse([]string{"-since", tt.dateStr, tempDir})

			if tt.expectOk && err != nil {
				t.Errorf("Expected date '%s' to be valid, got error: %v", tt.dateStr, err)
			}

			if !tt.expectOk && err == nil {
				t.Errorf("Expected date '%s' to be invalid, but got no error", tt.dateStr)
			}
		})
	}
}

// Helper functions
func timePtr(t time.Time) *time.Time {
	return &t
}

func createTempGitRepo(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "git-stats-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create .git directory to simulate a git repository
	gitDir := filepath.Join(tempDir, ".git")
	err = os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	return tempDir
}
