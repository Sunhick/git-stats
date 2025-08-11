// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Git package for command execution

package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// CommandResult represents the result of a git command execution
type CommandResult struct {
	Output   string
	Error    string
	ExitCode int
	Duration time.Duration
}

// Executor interface for git command execution
type Executor interface {
	Execute(ctx context.Context, command string, args ...string) (*CommandResult, error)
	ExecuteWithTimeout(command string, timeout time.Duration, args ...string) (*CommandResult, error)
	SetWorkingDirectory(path string) error
	GetWorkingDirectory() string
}

// GitCommandExecutor implements the Executor interface with security and performance features
type GitCommandExecutor struct {
	workingDir     string
	defaultTimeout time.Duration
	maxOutputSize  int64
}

// ExecutorConfig contains configuration options for the GitCommandExecutor
type ExecutorConfig struct {
	WorkingDirectory string
	DefaultTimeout   time.Duration
	MaxOutputSize    int64
}

// NewGitCommandExecutor creates a new GitCommandExecutor with the given configuration
func NewGitCommandExecutor(config ExecutorConfig) (*GitCommandExecutor, error) {
	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = 30 * time.Second
	}
	if config.MaxOutputSize == 0 {
		config.MaxOutputSize = 100 * 1024 * 1024 // 100MB default
	}

	executor := &GitCommandExecutor{
		workingDir:     config.WorkingDirectory,
		defaultTimeout: config.DefaultTimeout,
		maxOutputSize:  config.MaxOutputSize,
	}

	if config.WorkingDirectory != "" {
		if err := executor.SetWorkingDirectory(config.WorkingDirectory); err != nil {
			return nil, fmt.Errorf("failed to set working directory: %w", err)
		}
	}

	return executor, nil
}

// Execute runs a git command with the given context and arguments
func (e *GitCommandExecutor) Execute(ctx context.Context, command string, args ...string) (*CommandResult, error) {
	startTime := time.Now()

	// Sanitize command and arguments
	if err := e.sanitizeCommand(command, args...); err != nil {
		return nil, fmt.Errorf("command sanitization failed: %w", err)
	}

	// Prepare the command
	cmd := exec.CommandContext(ctx, "git", append([]string{command}, args...)...)

	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}

	// Set environment variables for consistent output
	cmd.Env = append(os.Environ(),
		"LC_ALL=C",           // Consistent locale
		"TZ=UTC",             // Consistent timezone
		"GIT_PAGER=",         // Disable pager
		"GIT_EDITOR=",        // Disable editor
		"GIT_ASKPASS=echo",   // Disable password prompts
	)

	// Execute the command
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	result := &CommandResult{
		Output:   string(output),
		Duration: duration,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
			result.Error = string(output)
		} else {
			result.Error = err.Error()
			result.ExitCode = -1
		}
		return result, fmt.Errorf("git command failed: %w", err)
	}

	// Check output size limits
	if int64(len(output)) > e.maxOutputSize {
		return result, fmt.Errorf("command output exceeds maximum size limit (%d bytes)", e.maxOutputSize)
	}

	return result, nil
}

// ExecuteWithTimeout runs a git command with a specific timeout
func (e *GitCommandExecutor) ExecuteWithTimeout(command string, timeout time.Duration, args ...string) (*CommandResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return e.Execute(ctx, command, args...)
}

// SetWorkingDirectory sets the working directory for git commands
func (e *GitCommandExecutor) SetWorkingDirectory(path string) error {
	if path == "" {
		return fmt.Errorf("working directory path cannot be empty")
	}

	// Clean and validate the path
	cleanPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Check if directory exists
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", cleanPath)
	}

	// Verify it's a git repository
	gitDir := filepath.Join(cleanPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository: %s", cleanPath)
	}

	e.workingDir = cleanPath
	return nil
}

// GetWorkingDirectory returns the current working directory
func (e *GitCommandExecutor) GetWorkingDirectory() string {
	return e.workingDir
}

// sanitizeCommand validates and sanitizes git commands and arguments to prevent injection
func (e *GitCommandExecutor) sanitizeCommand(command string, args ...string) error {
	// Validate command
	if err := e.validateGitCommand(command); err != nil {
		return err
	}

	// Validate arguments
	for i, arg := range args {
		if err := e.validateArgument(arg); err != nil {
			return fmt.Errorf("invalid argument at position %d: %w", i, err)
		}
	}

	return nil
}

// validateGitCommand checks if the command is a valid git subcommand
func (e *GitCommandExecutor) validateGitCommand(command string) error {
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Allow only alphanumeric characters, hyphens, and underscores
	validCommand := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validCommand.MatchString(command) {
		return fmt.Errorf("invalid command format: %s", command)
	}

	// Whitelist of allowed git commands for security
	allowedCommands := map[string]bool{
		"log":       true,
		"show":      true,
		"rev-list":  true,
		"shortlog":  true,
		"branch":    true,
		"status":    true,
		"diff":      true,
		"ls-files":  true,
		"rev-parse": true,
		"config":    true,
		"remote":    true,
		"tag":       true,
		"version":   true,
		"init":      true,
	}

	if !allowedCommands[command] {
		return fmt.Errorf("command not allowed: %s", command)
	}

	return nil
}

// validateArgument checks if an argument is safe to use
func (e *GitCommandExecutor) validateArgument(arg string) error {
	if arg == "" {
		return nil // Empty arguments are allowed
	}

	// Check for dangerous characters that could be used for injection
	// Note: pipe (|) is allowed in format strings for git log
	dangerousChars := []string{";", "&", "`", "$", "(", ")", "<", ">", "\\"}
	for _, char := range dangerousChars {
		if strings.Contains(arg, char) {
			return fmt.Errorf("argument contains dangerous character '%s': %s", char, arg)
		}
	}

	// Check for null bytes
	if strings.Contains(arg, "\x00") {
		return fmt.Errorf("argument contains null byte")
	}

	// Limit argument length to prevent buffer overflow attacks
	if len(arg) > 4096 {
		return fmt.Errorf("argument too long (max 4096 characters): %d", len(arg))
	}

	return nil
}

// IsGitAvailable checks if git is available in the system PATH
func IsGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// GetGitVersion returns the version of git installed on the system
func GetGitVersion() (string, error) {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git version: %w", err)
	}

	version := strings.TrimSpace(string(output))
	return version, nil
}

// SanitizeCommand is a public wrapper for testing the sanitization logic
func (e *GitCommandExecutor) SanitizeCommand(command string, args ...string) error {
	return e.sanitizeCommand(command, args...)
}
