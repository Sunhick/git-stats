// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Git package for command execution

package git

import (
	"context"
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
