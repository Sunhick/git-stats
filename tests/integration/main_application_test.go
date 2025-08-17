// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Main application integration tests

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMainApplicationIntegration tests the main application binary
func TestMainApplicationIntegration(t *testing.T) {
	// Build the main application
	srcDir := filepath.Join("..", "..", "src")
	binaryPath := filepath.Join(srcDir, "git-stats-test")

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "git-stats-test", ".")
	buildCmd.Dir = srcDir
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build main application: %v", err)
	}

	// Clean up binary after test
	defer os.Remove(binaryPath)

	// Create a test repository
	tempDir, cleanup := createTestRepository(t)
	defer cleanup()

	tests := []struct {
		name        string
		args        []string
		expectError bool
		expectText  string
	}{
		{
			name:        "Help command",
			args:        []string{"-help"},
			expectError: false,
			expectText:  "Git Stats - Enhanced Git Repository Analysis Tool",
		},
		{
			name:        "Contrib command with test repo",
			args:        []string{"-contrib", "-since", "2019-01-01", tempDir},
			expectError: false,
			expectText:  "No commits found", // This is the expected behavior for empty time range
		},
		{
			name:        "Summary command with test repo",
			args:        []string{"-summary", "-since", "2019-01-01", tempDir},
			expectError: false,
			expectText:  "No commits found", // This is the expected behavior for empty time range
		},
		{
			name:        "Contributors command with test repo",
			args:        []string{"-contributors", "-since", "2019-01-01", tempDir},
			expectError: false,
			expectText:  "No commits found", // This is the expected behavior for empty time range
		},
		{
			name:        "Health command with test repo",
			args:        []string{"-health", "-since", "2019-01-01", tempDir},
			expectError: false,
			expectText:  "No commits found", // This is the expected behavior for empty time range
		},
		{
			name:        "JSON output format",
			args:        []string{"-contrib", "-format", "json", "-since", "2019-01-01", tempDir},
			expectError: false,
			expectText:  "No commits found", // This is the expected behavior for empty time range
		},
		{
			name:        "CSV output format",
			args:        []string{"-summary", "-format", "csv", "-since", "2019-01-01", tempDir},
			expectError: false,
			expectText:  "No commits found", // This is the expected behavior for empty time range
		},
		{
			name:        "Invalid repository",
			args:        []string{"-contrib", "/nonexistent/path"},
			expectError: true,
			expectText:  "repository path does not exist",
		},
		{
			name:        "Invalid command",
			args:        []string{"-invalid"},
			expectError: true,
			expectText:  "flag provided but not defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded. Output: %s", outputStr)
				}
			} else {
				if err != nil {
					t.Errorf("Expected success but command failed with error: %v. Output: %s", err, outputStr)
				}
			}

			if tt.expectText != "" && !strings.Contains(outputStr, tt.expectText) {
				t.Errorf("Expected output to contain '%s', but got: %s", tt.expectText, outputStr)
			}
		})
	}
}

// TestMainApplicationBackwardCompatibility tests backward compatibility
func TestMainApplicationBackwardCompatibility(t *testing.T) {
	// Build the main application
	srcDir := filepath.Join("..", "..", "src")
	binaryPath := filepath.Join(srcDir, "git-stats-test")

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "git-stats-test", ".")
	buildCmd.Dir = srcDir
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build main application: %v", err)
	}

	// Clean up binary after test
	defer os.Remove(binaryPath)

	// Create a test repository
	tempDir, cleanup := createTestRepository(t)
	defer cleanup()

	// Test backward compatibility scenarios
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Default contrib command",
			args: []string{tempDir},
		},
		{
			name: "Explicit contrib command",
			args: []string{"-contrib", tempDir},
		},
		{
			name: "Summary command",
			args: []string{"-summary", tempDir},
		},
		{
			name: "With limit flag",
			args: []string{"-contrib", "-limit", "100", tempDir},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Errorf("Backward compatibility test failed: %v. Output: %s", err, string(output))
			}
		})
	}
}

// TestMainApplicationNewFeatures tests new features
func TestMainApplicationNewFeatures(t *testing.T) {
	// Build the main application
	srcDir := filepath.Join("..", "..", "src")
	binaryPath := filepath.Join(srcDir, "git-stats-test")

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "git-stats-test", ".")
	buildCmd.Dir = srcDir
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build main application: %v", err)
	}

	// Clean up binary after test
	defer os.Remove(binaryPath)

	// Create a test repository
	tempDir, cleanup := createTestRepository(t)
	defer cleanup()

	// Test new features that were previously not available
	tests := []struct {
		name        string
		args        []string
		expectText  string
	}{
		{
			name:       "Contributors command (new)",
			args:       []string{"-contributors", tempDir},
			expectText: "No commits found", // Expected behavior for empty time range
		},
		{
			name:       "Health command (new)",
			args:       []string{"-health", tempDir},
			expectText: "No commits found", // Expected behavior for empty time range
		},
		{
			name:       "JSON output (new)",
			args:       []string{"-contrib", "-format", "json", tempDir},
			expectText: "No commits found", // Expected behavior for empty time range
		},
		{
			name:       "CSV output (new)",
			args:       []string{"-summary", "-format", "csv", tempDir},
			expectText: "No commits found", // Expected behavior for empty time range
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if err != nil {
				t.Errorf("New feature test failed: %v. Output: %s", err, outputStr)
			}

			if !strings.Contains(outputStr, tt.expectText) {
				t.Errorf("Expected output to contain '%s', but got: %s", tt.expectText, outputStr)
			}
		})
	}
}
