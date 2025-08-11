// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for error utilities

package utils_test

import (
	"errors"
	"git-stats/utils"
	"strings"
	"testing"
)

func TestNewGitStatsError(t *testing.T) {
	cause := errors.New("underlying error")
	err := utils.NewGitStatsError(utils.ErrNotGitRepository, "not a git repository", cause)

	if err.Type != utils.ErrNotGitRepository {
		t.Errorf("Expected error type %v, got %v", utils.ErrNotGitRepository, err.Type)
	}

	if err.Message != "not a git repository" {
		t.Errorf("Expected message 'not a git repository', got '%s'", err.Message)
	}

	if err.Cause != cause {
		t.Errorf("Expected cause to be set correctly")
	}
}

func TestGitStatsErrorWithContext(t *testing.T) {
	err := utils.NewGitStatsError(utils.ErrInvalidDateFormat, "invalid date", nil)
	err.WithContext("date", "invalid-date-string")

	if err.Context["date"] != "invalid-date-string" {
		t.Errorf("Expected context to contain date key with value 'invalid-date-string'")
	}
}

func TestGetRecoverySuggestion(t *testing.T) {
	testCases := []struct {
		errorType utils.ErrorType
		expected  string
	}{
		{utils.ErrNotGitRepository, "Make sure you're running this command inside a git repository. Use 'git init' to initialize a new repository."},
		{utils.ErrGitNotFound, "Git is not installed or not found in PATH. Please install git and ensure it's accessible from the command line."},
		{utils.ErrInvalidDateFormat, "Use date format YYYY-MM-DD or relative dates like '1 week ago', '2024-01-01', etc."},
	}

	for _, tc := range testCases {
		err := utils.NewGitStatsError(tc.errorType, "test message", nil)
		suggestion := err.GetRecoverySuggestion()

		if suggestion != tc.expected {
			t.Errorf("For error type %v, expected suggestion '%s', got '%s'", tc.errorType, tc.expected, suggestion)
		}
	}
}

func TestIsRecoverable(t *testing.T) {
	recoverableErrors := []utils.ErrorType{
		utils.ErrNotGitRepository,
		utils.ErrGitNotFound,
		utils.ErrInvalidDateFormat,
		utils.ErrInvalidAuthor,
		utils.ErrInvalidArguments,
		utils.ErrInvalidFormat,
	}

	nonRecoverableErrors := []utils.ErrorType{
		utils.ErrRepositoryCorrupted,
		utils.ErrPermissionDenied,
		utils.ErrCommandTimeout,
		utils.ErrFileNotFound,
		utils.ErrInsufficientMemory,
	}

	for _, errType := range recoverableErrors {
		err := utils.NewGitStatsError(errType, "test", nil)
		if !err.IsRecoverable() {
			t.Errorf("Error type %v should be recoverable", errType)
		}
	}

	for _, errType := range nonRecoverableErrors {
		err := utils.NewGitStatsError(errType, "test", nil)
		if err.IsRecoverable() {
			t.Errorf("Error type %v should not be recoverable", errType)
		}
	}
}

func TestGetSeverity(t *testing.T) {
	testCases := []struct {
		errorType utils.ErrorType
		expected  string
	}{
		{utils.ErrNotGitRepository, "ERROR"},
		{utils.ErrGitNotFound, "ERROR"},
		{utils.ErrInvalidDateFormat, "WARNING"},
		{utils.ErrInvalidAuthor, "WARNING"},
		{utils.ErrRepositoryCorrupted, "CRITICAL"},
		{utils.ErrPermissionDenied, "CRITICAL"},
		{utils.ErrCommandTimeout, "ERROR"},
	}

	for _, tc := range testCases {
		err := utils.NewGitStatsError(tc.errorType, "test message", nil)
		severity := err.GetSeverity()

		if severity != tc.expected {
			t.Errorf("For error type %v, expected severity '%s', got '%s'", tc.errorType, tc.expected, severity)
		}
	}
}

func TestFormatUserFriendlyError(t *testing.T) {
	err := utils.NewGitStatsError(utils.ErrInvalidDateFormat, "Invalid date format", nil)
	err.WithContext("input", "invalid-date")

	formatted := err.FormatUserFriendlyError()

	if !strings.Contains(formatted, "[WARNING]") {
		t.Error("Expected formatted error to contain severity level")
	}

	if !strings.Contains(formatted, "Invalid date format") {
		t.Error("Expected formatted error to contain error message")
	}

	if !strings.Contains(formatted, "💡 Suggestion:") {
		t.Error("Expected formatted error to contain suggestion")
	}

	if !strings.Contains(formatted, "📋 Details:") {
		t.Error("Expected formatted error to contain context details")
	}

	if !strings.Contains(formatted, "input: invalid-date") {
		t.Error("Expected formatted error to contain context information")
	}
}

func TestErrorCollector(t *testing.T) {
	collector := utils.NewErrorCollector()

	// Add some errors
	err1 := utils.NewGitStatsError(utils.ErrNotGitRepository, "Not a git repo", nil)
	err2 := utils.NewGitStatsError(utils.ErrInvalidDateFormat, "Invalid date", nil) // Warning
	err3 := utils.NewGitStatsError(utils.ErrRepositoryCorrupted, "Corrupted repo", nil)

	collector.AddError(err1)
	collector.AddError(err2)
	collector.AddError(err3)

	if !collector.HasErrors() {
		t.Error("Expected collector to have errors")
	}

	if !collector.HasWarnings() {
		t.Error("Expected collector to have warnings")
	}

	if collector.GetErrorCount() != 2 {
		t.Errorf("Expected 2 errors, got %d", collector.GetErrorCount())
	}

	if collector.GetWarningCount() != 1 {
		t.Errorf("Expected 1 warning, got %d", collector.GetWarningCount())
	}

	summary := collector.GetSummary()
	if !strings.Contains(summary, "❌ 2 error(s) occurred:") {
		t.Error("Expected summary to contain error count")
	}

	if !strings.Contains(summary, "⚠️  1 warning(s):") {
		t.Error("Expected summary to contain warning count")
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := utils.WrapError(originalErr, utils.ErrFileNotFound, "File not found")

	if wrappedErr.Type != utils.ErrFileNotFound {
		t.Errorf("Expected error type %v, got %v", utils.ErrFileNotFound, wrappedErr.Type)
	}

	if wrappedErr.Message != "File not found" {
		t.Errorf("Expected message 'File not found', got '%s'", wrappedErr.Message)
	}

	if wrappedErr.Cause != originalErr {
		t.Error("Expected wrapped error to contain original error as cause")
	}
}
