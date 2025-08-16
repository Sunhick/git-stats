// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Integration tests for CLI, config, and filtering

package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"git-stats/cli"
	"git-stats/integration"
	"git-stats/models"
)

// Helper function to create test commits
func createTestCommits() []models.Commit {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	return []models.Commit{
		{
			Hash:       "abc123",
			Message:    "Initial commit",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: baseTime,
			Stats: models.CommitStats{
				FilesChanged: 2,
				Insertions:   10,
				Deletions:    0,
				Files: []models.FileChange{
					{Path: "main.go", Status: "A", Insertions: 5, Deletions: 0},
					{Path: "README.md", Status: "A", Insertions: 5, Deletions: 0},
				},
			},
		},
		{
			Hash:       "def456",
			Message:    "Add feature",
			Author:     models.Author{Name: "Jane Smith", Email: "jane@example.com"},
			AuthorDate: baseTime.AddDate(0, 0, 1),
			Stats: models.CommitStats{
				FilesChanged: 1,
				Insertions:   15,
				Deletions:    2,
				Files: []models.FileChange{
					{Path: "feature.go", Status: "A", Insertions: 15, Deletions: 2},
				},
			},
		},
	}
}

func TestFilteredConfigManager_Basic(t *testing.T) {
	// Create temporary directory for test config
	tempDir, err := os.MkdirTemp("", "git-stats-integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")
	fcm := integration.NewFilteredConfigManagerWithPath(configPath)

	// Test loading config (should use defaults for non-existent file)
	err = fcm.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test getting config manager
	configManager := fcm.GetConfigManager()
	if configManager == nil {
		t.Fatal("Config manager should not be nil")
	}

	// Test getting filter builder
	filterBuilder := fcm.GetFilterBuilder()
	if filterBuilder == nil {
		t.Fatal("Filter builder should not be nil")
	}
}

func TestFilteredConfigManager_BuildFilterChain(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		cliConfig *cli.Config
		expectErr bool
	}{
		{
			name: "Basic CLI config",
			cliConfig: &cli.Config{
				Command: "contrib",
				Format:  "terminal",
			},
			expectErr: false,
		},
		{
			name: "CLI config with filters",
			cliConfig: &cli.Config{
				Command: "summary",
				Since:   &baseTime,
				Author:  "John Doe",
				Limit:   1000,
				Format:  "json",
			},
			expectErr: false,
		},
		{
			name: "CLI config with invalid author regex",
			cliConfig: &cli.Config{
				Command: "contrib",
				Author:  "[invalid-regex",
				Format:  "terminal",
			},
			expectErr: false, // Should not error with basic author filter
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain, err := fcm.BuildFilterChain(tt.cliConfig)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectErr && chain == nil {
				t.Error("Filter chain should not be nil")
			}
		})
	}
}

func TestFilteredConfigManager_ApplyFilters(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()
	commits := createTestCommits()

	tests := []struct {
		name      string
		cliConfig *cli.Config
		expected  int
	}{
		{
			name: "No filters",
			cliConfig: &cli.Config{
				Command: "contrib",
				Format:  "terminal",
			},
			expected: 2, // All commits
		},
		{
			name: "Author filter",
			cliConfig: &cli.Config{
				Command: "contrib",
				Author:  "John",
				Format:  "terminal",
			},
			expected: 1, // Only John's commit
		},
		{
			name: "Date filter",
			cliConfig: &cli.Config{
				Command: "contrib",
				Since:   func() *time.Time { t := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC); return &t }(),
				Until:   func() *time.Time { t := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC); return &t }(),
				Format:  "terminal",
			},
			expected: 1, // Only first commit
		},
		{
			name: "Limit filter",
			cliConfig: &cli.Config{
				Command: "contrib",
				Limit:   1,
				Format:  "terminal",
			},
			expected: 1, // Limited to 1 commit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, err := fcm.ApplyFilters(commits, tt.cliConfig)
			if err != nil {
				t.Fatalf("Failed to apply filters: %v", err)
			}

			if len(filtered) != tt.expected {
				t.Errorf("Expected %d commits, got %d", tt.expected, len(filtered))
			}
		})
	}
}

func TestFilteredConfigManager_CreateAnalysisConfig(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	cliConfig := &cli.Config{
		Command: "summary",
		Since:   &baseTime,
		Until:   func() *time.Time { t := baseTime.AddDate(0, 0, 7); return &t }(),
		Author:  "test@example.com",
		Limit:   5000,
		Format:  "json",
	}

	analysisConfig := fcm.CreateAnalysisConfig(cliConfig)

	// Verify analysis config
	if analysisConfig.Limit != 5000 {
		t.Errorf("Expected limit 5000, got %d", analysisConfig.Limit)
	}

	if analysisConfig.AuthorFilter != "test@example.com" {
		t.Errorf("Expected author filter 'test@example.com', got '%s'", analysisConfig.AuthorFilter)
	}

	if analysisConfig.TimeRange.Start != baseTime {
		t.Errorf("Expected start time %v, got %v", baseTime, analysisConfig.TimeRange.Start)
	}

	expectedEnd := baseTime.AddDate(0, 0, 7)
	if analysisConfig.TimeRange.End != expectedEnd {
		t.Errorf("Expected end time %v, got %v", expectedEnd, analysisConfig.TimeRange.End)
	}
}

func TestFilteredConfigManager_GetEffectiveConfig(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	// Load default config
	err := fcm.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	cliConfig := &cli.Config{
		Command:    "summary",
		Format:     "json",
		Author:     "cli-author",
		ShowProgress: true,
		NoColor:    true,
		Limit:      2000,
	}

	effectiveConfig := fcm.GetEffectiveConfig(cliConfig)

	// Verify CLI overrides
	if effectiveConfig.Command != "summary" {
		t.Errorf("Expected command 'summary', got '%s'", effectiveConfig.Command)
	}

	if effectiveConfig.Format != "json" {
		t.Errorf("Expected format 'json', got '%s'", effectiveConfig.Format)
	}

	if effectiveConfig.Author != "cli-author" {
		t.Errorf("Expected author 'cli-author', got '%s'", effectiveConfig.Author)
	}

	if !effectiveConfig.ShowProgress {
		t.Error("Expected show progress to be true")
	}

	if !effectiveConfig.NoColor {
		t.Error("Expected no color to be true")
	}

	if effectiveConfig.Limit != 2000 {
		t.Errorf("Expected limit 2000, got %d", effectiveConfig.Limit)
	}

	// Verify config defaults are preserved
	if effectiveConfig.RepoPath != "." {
		t.Errorf("Expected default repo path '.', got '%s'", effectiveConfig.RepoPath)
	}

	if effectiveConfig.ColorTheme != "github" {
		t.Errorf("Expected default color theme 'github', got '%s'", effectiveConfig.ColorTheme)
	}
}

func TestFilteredConfigManager_UpdateConfigFromCLI(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	// Load default config
	err := fcm.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	cliConfig := &cli.Config{
		Command:      "health",
		Format:       "csv",
		ShowProgress: true,
		NoColor:      true,
		ColorTheme:   "blue",
		Limit:        15000,
		RepoPath:     "/custom/path",
	}

	// Update config from CLI
	fcm.UpdateConfigFromCLI(cliConfig)

	// Verify updates
	appConfig := fcm.GetConfigManager().GetConfig()

	if appConfig.Defaults.Command != "health" {
		t.Errorf("Expected updated command 'health', got '%s'", appConfig.Defaults.Command)
	}

	if appConfig.Defaults.Format != "csv" {
		t.Errorf("Expected updated format 'csv', got '%s'", appConfig.Defaults.Format)
	}

	if !appConfig.Defaults.ShowProgress {
		t.Error("Expected updated show progress to be true")
	}

	if appConfig.Output.ColorEnabled {
		t.Error("Expected color enabled to be false due to NoColor")
	}

	if appConfig.Output.ColorTheme != "blue" {
		t.Errorf("Expected updated color theme 'blue', got '%s'", appConfig.Output.ColorTheme)
	}

	if appConfig.Performance.MaxCommits != 15000 {
		t.Errorf("Expected updated max commits 15000, got %d", appConfig.Performance.MaxCommits)
	}

	if appConfig.Defaults.RepoPath != "/custom/path" {
		t.Errorf("Expected updated repo path '/custom/path', got '%s'", appConfig.Defaults.RepoPath)
	}
}

func TestFilteredConfigManager_GetFilterSummary(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	tests := []struct {
		name      string
		cliConfig *cli.Config
		expectEmpty bool
	}{
		{
			name: "No filters",
			cliConfig: &cli.Config{
				Command: "contrib",
				Format:  "terminal",
			},
			expectEmpty: true,
		},
		{
			name: "With filters",
			cliConfig: &cli.Config{
				Command: "contrib",
				Author:  "John",
				Since:   func() *time.Time { t := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC); return &t }(),
				Format:  "terminal",
			},
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary, err := fcm.GetFilterSummary(tt.cliConfig)
			if err != nil {
				t.Fatalf("Failed to get filter summary: %v", err)
			}

			if tt.expectEmpty && summary != "No filters applied" && summary != "Default filters applied" {
				t.Errorf("Expected empty summary, got '%s'", summary)
			}

			if !tt.expectEmpty && (summary == "No filters applied" || summary == "Default filters applied") {
				t.Errorf("Expected non-empty summary, got '%s'", summary)
			}
		})
	}
}

func TestFilteredConfigManager_ValidateEffectiveConfig(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	tests := []struct {
		name           string
		effectiveConfig *integration.EffectiveConfig
		expectErr      bool
	}{
		{
			name: "Valid config",
			effectiveConfig: &integration.EffectiveConfig{
				Command:    "contrib",
				Format:     "terminal",
				RepoPath:   ".",
				ColorTheme: "github",
				Limit:      10000,
				ChunkSize:  1000,
				MaxWorkers: 4,
			},
			expectErr: false,
		},
		{
			name: "Invalid command",
			effectiveConfig: &integration.EffectiveConfig{
				Command:    "invalid",
				Format:     "terminal",
				RepoPath:   ".",
				ColorTheme: "github",
				Limit:      10000,
				ChunkSize:  1000,
				MaxWorkers: 4,
			},
			expectErr: true,
		},
		{
			name: "Invalid format",
			effectiveConfig: &integration.EffectiveConfig{
				Command:    "contrib",
				Format:     "invalid",
				RepoPath:   ".",
				ColorTheme: "github",
				Limit:      10000,
				ChunkSize:  1000,
				MaxWorkers: 4,
			},
			expectErr: true,
		},
		{
			name: "Invalid limit",
			effectiveConfig: &integration.EffectiveConfig{
				Command:    "contrib",
				Format:     "terminal",
				RepoPath:   ".",
				ColorTheme: "github",
				Limit:      -1,
				ChunkSize:  1000,
				MaxWorkers: 4,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fcm.ValidateEffectiveConfig(tt.effectiveConfig)

			if tt.expectErr && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}
