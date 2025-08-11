// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for error utilities

package utils_test

import (
	"errors"
	"testing"
	"git-stats/utils"
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
