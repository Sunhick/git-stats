// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Formatters package for output formatting interfaces

package formatters

import (
	"git-stats/git"
	"git-stats/models"
)

// Formatter interface for output formatting
type Formatter interface {
	Format(data *models.AnalysisResult, config models.FormatConfig) ([]byte, error)
}

// JSONFormatter interface for JSON output formatting
type JSONFormatter interface {
	FormatJSON(data *models.AnalysisResult, config models.FormatConfig) ([]byte, error)
	FormatPrettyJSON(data *models.AnalysisResult, config models.FormatConfig) ([]byte, error)
}

// CSVFormatter interface for CSV output formatting
type CSVFormatter interface {
	FormatCSV(data *models.AnalysisResult, config models.FormatConfig) ([]byte, error)
	FormatCommitsCSV(commits []git.Commit) ([]byte, error)
	FormatContributorsCSV(contributors []models.ContributorStats) ([]byte, error)
}

// TerminalFormatter interface for terminal output formatting
type TerminalFormatter interface {
	FormatTerminal(data *models.AnalysisResult, config models.FormatConfig) ([]byte, error)
	FormatColorized(text string, color string) string
	FormatTable(headers []string, rows [][]string) string
}
