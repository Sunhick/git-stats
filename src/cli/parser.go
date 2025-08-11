// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - CLI package for command line parsing

package cli

import (
	"time"
)

// Config represents the configuration for the git-stats tool
type Config struct {
	Command      string     // contrib, summary, contributors, health, gui
	Since        *time.Time // --since flag
	Until        *time.Time // --until flag
	Author       string     // --author flag
	Format       string     // json, csv, terminal
	OutputFile   string     // --output flag
	RepoPath     string     // repository path
	ShowProgress bool       // --progress flag
	Limit        int        // --limit flag for large repos
	GUIMode      bool       // --gui flag for ncurses interface
}

// Parser interface for command line parsing
type Parser interface {
	Parse(args []string) (*Config, error)
	PrintUsage()
	PrintHelp()
}
