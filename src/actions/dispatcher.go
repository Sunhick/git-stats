// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Command dispatcher for routing and execution

package actions

import (
	"fmt"
	"git-stats/cli"
	"git-stats/git"
	"os"
	"os/exec"
	"path/filepath"
)

// CommandDispatcher handles routing and execution of different analysis commands
type CommandDispatcher struct {
	validator cli.Validator
}

// NewCommandDispatcher creates a new command dispatcher
func NewCommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		validator: cli.NewCLIValidator(),
	}
}

// ExecuteCommand dispatches and executes the appropriate command based on configuration
func (d *CommandDispatcher) ExecuteCommand(config *cli.Config) error {
	// Validate configuration
	if err := d.validateConfiguration(config); err != nil {
		return NewCommandError(ErrInvalidConfiguration, fmt.Sprintf("Configuration validation failed: %v", err), err)
	}

	// Validate system requirements
	if err := d.validateSystemRequirements(); err != nil {
		return NewCommandError(ErrSystemRequirements, fmt.Sprintf("System requirements not met: %v", err), err)
	}

	// Validate repository
	if err := d.validateRepository(config.RepoPath); err != nil {
		return NewCommandError(ErrRepositoryAccess, fmt.Sprintf("Repository validation failed: %v", err), err)
	}

	// Route to appropriate command handler
	switch config.Command {
	case "contrib":
		return d.executeContribCommand(config)
	case "summary":
		return d.executeSummaryCommand(config)
	case "contributors":
		return d.executeContributorsCommand(config)
	case "health":
		return d.executeHealthCommand(config)
	default:
		return NewCommandError(ErrUnknownCommand, fmt.Sprintf("Unknown command: %s", config.Command), nil)
	}
}

// validateConfiguration validates the command configuration
func (d *CommandDispatcher) validateConfiguration(config *cli.Config) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Validate individual components but skip repository validation
	// (we'll do that separately with better error classification)
	if err := d.validator.ValidateDateRange(config.Since, config.Until); err != nil {
		return err
	}

	if config.Author != "" {
		if err := d.validator.ValidateAuthor(config.Author); err != nil {
			return err
		}
	}

	if err := d.validator.ValidateFormat(config.Format); err != nil {
		return err
	}

	if config.OutputFile != "" {
		if err := d.validator.ValidateOutputFile(config.OutputFile); err != nil {
			return err
		}
	}

	if err := d.validator.ValidateLimit(config.Limit); err != nil {
		return err
	}

	// Validate command separately
	validCommands := []string{"contrib", "summary", "contributors", "health"}
	validCommand := false
	for _, valid := range validCommands {
		if config.Command == valid {
			validCommand = true
			break
		}
	}
	if !validCommand {
		return fmt.Errorf("invalid command '%s'", config.Command)
	}

	return nil
}

// validateSystemRequirements checks if git is available and accessible
func (d *CommandDispatcher) validateSystemRequirements() error {
	// Check if git is installed and accessible
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git is not installed or not in PATH. Please install git and try again")
	}

	// Test git execution
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git is installed but not working properly: %v", err)
	}

	return nil
}

// validateRepository validates that the repository path is accessible and is a git repository
func (d *CommandDispatcher) validateRepository(repoPath string) error {
	// Check if path exists and is accessible
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("repository path does not exist: %s", repoPath)
	} else if err != nil {
		return fmt.Errorf("cannot access repository path: %v", err)
	}

	// Check if it's a directory
	info, err := os.Stat(repoPath)
	if err != nil {
		return fmt.Errorf("error accessing repository: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("repository path is not a directory: %s", repoPath)
	}

	// Check if it's a git repository by looking for .git directory or file
	gitPath := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository (no .git directory found): %s", repoPath)
	}

	// Try to create a git repository instance to validate it's working
	repoConfig := git.RepositoryConfig{
		Path: repoPath,
	}
	repo, err := git.NewGitRepository(repoConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %v", err)
	}

	// Validate repository is accessible
	if !repo.IsValidRepository() {
		return fmt.Errorf("invalid git repository: %s", repoPath)
	}

	return nil
}

// executeContribCommand executes the contribution graph command
func (d *CommandDispatcher) executeContribCommand(config *cli.Config) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Fatal error in contribution analysis: %v\n", r)
		}
	}()

	if config.GUIMode {
		LaunchGUI(config)
		return nil
	}

	// The ContribWithConfig function handles its own errors
	// and prints directly to stdout/stderr, so we just call it
	ContribWithConfig(config)
	return nil
}

// executeSummaryCommand executes the summary statistics command
func (d *CommandDispatcher) executeSummaryCommand(config *cli.Config) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Fatal error in summary analysis: %v\n", r)
		}
	}()

	if config.GUIMode {
		LaunchGUI(config)
		return nil
	}

	SummarizeWithConfig(config)
	return nil
}

// executeContributorsCommand executes the contributors analysis command
func (d *CommandDispatcher) executeContributorsCommand(config *cli.Config) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Fatal error in contributors analysis: %v\n", r)
		}
	}()

	if config.GUIMode {
		LaunchGUI(config)
		return nil
	}

	ContributorsWithConfig(config)
	return nil
}

// executeHealthCommand executes the repository health analysis command
func (d *CommandDispatcher) executeHealthCommand(config *cli.Config) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Fatal error in health analysis: %v\n", r)
		}
	}()

	if config.GUIMode {
		LaunchGUI(config)
		return nil
	}

	HealthWithConfig(config)
	return nil
}

// CommandErrorType represents different types of command errors
type CommandErrorType int

const (
	ErrUnknownCommand CommandErrorType = iota
	ErrInvalidConfiguration
	ErrSystemRequirements
	ErrRepositoryAccess
	ErrNotImplemented
	ErrExecutionFailed
)

// CommandError represents an error that occurred during command execution
type CommandError struct {
	Type    CommandErrorType
	Message string
	Cause   error
}

// Error implements the error interface
func (e *CommandError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause error
func (e *CommandError) Unwrap() error {
	return e.Cause
}

// NewCommandError creates a new command error
func NewCommandError(errorType CommandErrorType, message string, cause error) *CommandError {
	return &CommandError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// IsCommandError checks if an error is a CommandError
func IsCommandError(err error) bool {
	_, ok := err.(*CommandError)
	return ok
}

// GetErrorType returns the error type if the error is a CommandError
func GetErrorType(err error) (CommandErrorType, bool) {
	if cmdErr, ok := err.(*CommandError); ok {
		return cmdErr.Type, true
	}
	return 0, false
}

// GetUserFriendlyMessage returns a user-friendly error message with suggestions
func GetUserFriendlyMessage(err error) string {
	if cmdErr, ok := err.(*CommandError); ok {
		switch cmdErr.Type {
		case ErrUnknownCommand:
			return fmt.Sprintf("Error: %s\n\nSuggestion: Use one of the available commands: contrib, summary, contributors, health\nExample: git-stats -contrib", cmdErr.Message)
		case ErrInvalidConfiguration:
			return fmt.Sprintf("Error: %s\n\nSuggestion: Check your command line arguments and try again.\nFor help, run: git-stats -help", cmdErr.Message)
		case ErrSystemRequirements:
			return fmt.Sprintf("Error: %s\n\nSuggestion: Please install git and ensure it's available in your PATH.\nYou can download git from: https://git-scm.com/downloads", cmdErr.Message)
		case ErrRepositoryAccess:
			return fmt.Sprintf("Error: %s\n\nSuggestion: Make sure you're in a git repository directory or specify a valid repository path.\nExample: git-stats /path/to/your/git/repo", cmdErr.Message)
		case ErrNotImplemented:
			return fmt.Sprintf("Error: %s\n\nSuggestion: This feature is coming soon. Try using -contrib or -summary commands instead.", cmdErr.Message)
		case ErrExecutionFailed:
			return fmt.Sprintf("Error: %s\n\nSuggestion: Please check the repository state and try again. If the problem persists, try with a smaller date range using --since and --until flags.", cmdErr.Message)
		default:
			return fmt.Sprintf("Error: %s", cmdErr.Message)
		}
	}
	return fmt.Sprintf("Error: %v", err)
}
