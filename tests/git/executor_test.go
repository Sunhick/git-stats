// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Git command executor tests

package git

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"git-stats/git"
)

func TestNewGitCommandExecutor(t *testing.T) {
	tests := []struct {
		name    string
		config  git.ExecutorConfig
		wantErr bool
	}{
		{
			name: "valid config with defaults",
			config: git.ExecutorConfig{
				WorkingDirectory: "",
			},
			wantErr: false,
		},
		{
			name: "valid config with custom values",
			config: git.ExecutorConfig{
				WorkingDirectory: "",
				DefaultTimeout:   60 * time.Second,
				MaxOutputSize:    50 * 1024 * 1024,
			},
			wantErr: false,
		},
		{
			name: "invalid working directory",
			config: git.ExecutorConfig{
				WorkingDirectory: "/nonexistent/directory",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor, err := git.NewGitCommandExecutor(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGitCommandExecutor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && executor == nil {
				t.Error("NewGitCommandExecutor() returned nil executor without error")
			}
		})
	}
}

func TestGitCommandExecutor_SetWorkingDirectory(t *testing.T) {
	executor, err := git.NewGitCommandExecutor(git.ExecutorConfig{})
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "nonexistent directory",
			path:    "/nonexistent/directory",
			wantErr: true,
		},
		{
			name:    "non-git directory",
			path:    os.TempDir(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.SetWorkingDirectory(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetWorkingDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitCommandExecutor_ValidateGitCommand(t *testing.T) {
	executor, err := git.NewGitCommandExecutor(git.ExecutorConfig{})
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{
			name:    "valid command - log",
			command: "log",
			wantErr: false,
		},
		{
			name:    "valid command - shortlog",
			command: "shortlog",
			wantErr: false,
		},
		{
			name:    "valid command - branch",
			command: "branch",
			wantErr: false,
		},
		{
			name:    "empty command",
			command: "",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			command: "log; rm -rf /",
			wantErr: true,
		},
		{
			name:    "disallowed command",
			command: "push",
			wantErr: true,
		},
		{
			name:    "command with spaces",
			command: "log --oneline",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use reflection or make validateGitCommand public for testing
			// For now, we'll test through sanitizeCommand
			err := executor.SanitizeCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGitCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitCommandExecutor_ValidateArgument(t *testing.T) {
	executor, err := git.NewGitCommandExecutor(git.ExecutorConfig{})
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	tests := []struct {
		name    string
		arg     string
		wantErr bool
	}{
		{
			name:    "valid argument",
			arg:     "--oneline",
			wantErr: false,
		},
		{
			name:    "valid date argument",
			arg:     "--since=2023-01-01",
			wantErr: false,
		},
		{
			name:    "empty argument",
			arg:     "",
			wantErr: false,
		},
		{
			name:    "argument with semicolon",
			arg:     "--format='; rm -rf /'",
			wantErr: true,
		},
		{
			name:    "argument with backtick in format",
			arg:     "--format='%H `cat /etc/passwd`'",
			wantErr: true,
		},
		{
			name:    "argument with backtick",
			arg:     "--format='`rm -rf /`'",
			wantErr: true,
		},
		{
			name:    "argument with null byte",
			arg:     "test\x00",
			wantErr: true,
		},
		{
			name:    "very long argument",
			arg:     strings.Repeat("a", 5000),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through sanitizeCommand with a valid command
			err := executor.SanitizeCommand("log", tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateArgument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitCommandExecutor_ExecuteWithTimeout(t *testing.T) {
	// Skip if git is not available
	if !git.IsGitAvailable() {
		t.Skip("git not available in PATH")
	}

	executor, err := git.NewGitCommandExecutor(git.ExecutorConfig{})
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	tests := []struct {
		name    string
		command string
		args    []string
		timeout time.Duration
		wantErr bool
	}{
		{
			name:    "valid command with sufficient timeout",
			command: "version",
			args:    []string{},
			timeout: 5 * time.Second,
			wantErr: false,
		},
		{
			name:    "command with very short timeout",
			command: "log",
			args:    []string{"--all"},
			timeout: 1 * time.Nanosecond,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := executor.ExecuteWithTimeout(tt.command, tt.timeout, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Error("ExecuteWithTimeout() returned nil result without error")
			}
		})
	}
}

func TestIsGitAvailable(t *testing.T) {
	available := git.IsGitAvailable()
	// We can't assert a specific value since it depends on the test environment
	// Just ensure the function doesn't panic
	t.Logf("Git available: %v", available)
}

func TestGetGitVersion(t *testing.T) {
	if !git.IsGitAvailable() {
		t.Skip("git not available in PATH")
	}

	version, err := git.GetGitVersion()
	if err != nil {
		t.Errorf("GetGitVersion() error = %v", err)
		return
	}

	if version == "" {
		t.Error("GetGitVersion() returned empty version")
	}

	if !strings.Contains(version, "git version") {
		t.Errorf("GetGitVersion() returned unexpected format: %s", version)
	}
}

// Helper function to create a temporary git repository for testing
func createTempGitRepo(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "git-stats-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Initialize git repository
	if git.IsGitAvailable() {
		// Create executor without working directory first
		executor, err := git.NewGitCommandExecutor(git.ExecutorConfig{})
		if err != nil {
			t.Fatalf("Failed to create executor: %v", err)
		}

		// Change to temp directory and initialize
		oldDir, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldDir)

		ctx := context.Background()
		executor.Execute(ctx, "init")
		executor.Execute(ctx, "config", "user.name", "Test User")
		executor.Execute(ctx, "config", "user.email", "test@example.com")
	}

	return tempDir
}

func TestGitCommandExecutor_Integration(t *testing.T) {
	if !git.IsGitAvailable() {
		t.Skip("git not available in PATH")
	}

	// Create a temporary git repository
	tempDir := createTempGitRepo(t)
	defer os.RemoveAll(tempDir)

	executor, err := git.NewGitCommandExecutor(git.ExecutorConfig{
		WorkingDirectory: tempDir,
	})
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Test basic git commands
	ctx := context.Background()

	// Test git status
	result, err := executor.Execute(ctx, "status", "--porcelain")
	if err != nil {
		t.Errorf("Execute(status) error = %v", err)
	}
	if result == nil {
		t.Error("Execute(status) returned nil result")
	}

	// Test git config
	result, err = executor.Execute(ctx, "config", "user.name")
	if err != nil {
		t.Errorf("Execute(config) error = %v", err)
	}
	if result != nil && !strings.Contains(result.Output, "Test User") {
		t.Errorf("Execute(config) unexpected output: %s", result.Output)
	}
}

// Benchmark tests
func BenchmarkGitCommandExecutor_Execute(b *testing.B) {
	if !git.IsGitAvailable() {
		b.Skip("git not available in PATH")
	}

	executor, err := git.NewGitCommandExecutor(git.ExecutorConfig{})
	if err != nil {
		b.Fatalf("Failed to create executor: %v", err)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := executor.Execute(ctx, "version")
		if err != nil {
			b.Errorf("Execute() error = %v", err)
		}
	}
}
