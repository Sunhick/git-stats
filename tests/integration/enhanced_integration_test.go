// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Enhanced integration tests for task 7.2

package integration

import (
	"testing"
	"time"

	"git-stats/cli"
	"git-stats/integration"
	"git-stats/models"
)

// TestEnhancedFilteredConfigManager tests the enhanced filtered config manager
func TestEnhancedFilteredConfigManager(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	// Load default configuration
	err := fcm.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test enhanced effective configuration
	t.Run("Enhanced effective configuration", func(t *testing.T) {
		cliConfig := &cli.Config{
			Command:    "contrib",
			Author:     "john@example.com",
			Format:     "json",
			Limit:      5000,
			NoColor:    true,
		}

		effectiveConfig := fcm.GetEffectiveConfig(cliConfig)

		// Test that all new configuration fields are present
		if effectiveConfig.AuthorMatchType == "" {
			t.Error("Expected author match type to be set")
		}

		if effectiveConfig.CaseSensitive {
			t.Error("Expected case sensitive to be false by default")
		}

		if len(effectiveConfig.ExcludePatterns) < 0 {
			t.Error("Expected exclude patterns to be initialized")
		}

		if len(effectiveConfig.IncludePatterns) < 0 {
			t.Error("Expected include patterns to be initialized")
		}
	})

	// Test configuration validation
	t.Run("Enhanced configuration validation", func(t *testing.T) {
		cliConfig := &cli.Config{
			Command: "contrib",
			Format:  "json",
			Limit:   1000,
		}

		effectiveConfig := fcm.GetEffectiveConfig(cliConfig)

		err := fcm.ValidateEffectiveConfig(effectiveConfig)
		if err != nil {
			t.Errorf("Valid configuration should not produce validation error: %v", err)
		}
	})

	// Test filter chain building with enhanced options
	t.Run("Enhanced filter chain building", func(t *testing.T) {
		baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

		cliConfig := &cli.Config{
			Command: "contrib",
			Since:   &baseTime,
			Until:   func() *time.Time { t := baseTime.AddDate(0, 0, 7); return &t }(),
			Author:  "john",
			Format:  "terminal",
			Limit:   1000,
		}

		chain, err := fcm.BuildFilterChain(cliConfig)
		if err != nil {
			t.Fatalf("Failed to build filter chain: %v", err)
		}

		if chain == nil {
			t.Fatal("Filter chain should not be nil")
		}

		filters := chain.GetFilters()
		if len(filters) == 0 {
			t.Error("Filter chain should contain filters")
		}

		// Test filter summary
		summary, err := fcm.GetFilterSummary(cliConfig)
		if err != nil {
			t.Errorf("Failed to get filter summary: %v", err)
		}

		if summary == "" {
			t.Error("Filter summary should not be empty")
		}
	})

	// Test applying filters with enhanced functionality
	t.Run("Enhanced filter application", func(t *testing.T) {
		baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

		commits := []models.Commit{
			{
				Hash:       "abc123",
				Message:    "feat: add new feature",
				Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
				AuthorDate: baseTime.AddDate(0, 0, -2),
				Stats:      models.CommitStats{FilesChanged: 3, Insertions: 50, Deletions: 10},
			},
			{
				Hash:       "def456",
				Message:    "fix: resolve bug",
				Author:     models.Author{Name: "Jane Smith", Email: "jane@example.com"},
				AuthorDate: baseTime.AddDate(0, 0, -1),
				Stats:      models.CommitStats{FilesChanged: 1, Insertions: 20, Deletions: 5},
			},
			{
				Hash:       "ghi789",
				Message:    "docs: update README",
				Author:     models.Author{Name: "Bob Wilson", Email: "bob@example.com"},
				AuthorDate: baseTime,
				Stats:      models.CommitStats{FilesChanged: 1, Insertions: 10, Deletions: 2},
			},
		}

		cliConfig := &cli.Config{
			Command: "contrib",
			Since:   func() *time.Time { t := baseTime.AddDate(0, 0, -3); return &t }(),
			Author:  "john",
			Format:  "terminal",
		}

		filtered, err := fcm.ApplyFilters(commits, cliConfig)
		if err != nil {
			t.Fatalf("Failed to apply filters: %v", err)
		}

		if len(filtered) != 1 {
			t.Errorf("Expected 1 filtered commit, got %d", len(filtered))
		}

		if len(filtered) > 0 && filtered[0].Hash != "abc123" {
			t.Errorf("Expected commit abc123, got %s", filtered[0].Hash)
		}
	})
}

// TestAnalysisConfigCreation tests enhanced analysis config creation
func TestAnalysisConfigCreation(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	cliConfig := &cli.Config{
		Command: "contrib",
		Since:   &baseTime,
		Until:   func() *time.Time { t := baseTime.AddDate(0, 0, 7); return &t }(),
		Author:  "john@example.com",
		Limit:   5000,
	}

	analysisConfig := fcm.CreateAnalysisConfig(cliConfig)

	// Test that analysis config is properly created
	if analysisConfig.Limit != 5000 {
		t.Errorf("Expected limit 5000, got %d", analysisConfig.Limit)
	}

	if analysisConfig.AuthorFilter != "john@example.com" {
		t.Errorf("Expected author filter 'john@example.com', got '%s'", analysisConfig.AuthorFilter)
	}

	if analysisConfig.TimeRange.Start.IsZero() {
		t.Error("Expected time range start to be set")
	}

	if analysisConfig.TimeRange.End.IsZero() {
		t.Error("Expected time range end to be set")
	}
}

// TestConfigurationUpdatesFromCLI tests updating configuration from CLI
func TestConfigurationUpdatesFromCLI(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	cliConfig := &cli.Config{
		Command:      "summary",
		Format:       "csv",
		ShowProgress: true,
		RepoPath:     "/custom/path",
		NoColor:      true,
		ColorTheme:   "blue",
		Limit:        15000,
	}

	// Update config from CLI
	fcm.UpdateConfigFromCLI(cliConfig)

	// Get updated config
	config := fcm.GetConfigManager().GetConfig()

	// Verify updates
	if config.Defaults.Command != "summary" {
		t.Errorf("Expected command 'summary', got '%s'", config.Defaults.Command)
	}

	if config.Defaults.Format != "csv" {
		t.Errorf("Expected format 'csv', got '%s'", config.Defaults.Format)
	}

	if !config.Defaults.ShowProgress {
		t.Error("Expected show progress to be true")
	}

	if config.Defaults.RepoPath != "/custom/path" {
		t.Errorf("Expected repo path '/custom/path', got '%s'", config.Defaults.RepoPath)
	}

	if config.Output.ColorEnabled {
		t.Error("Expected color enabled to be false")
	}

	if config.Output.ColorTheme != "blue" {
		t.Errorf("Expected color theme 'blue', got '%s'", config.Output.ColorTheme)
	}

	if config.Performance.MaxCommits != 15000 {
		t.Errorf("Expected max commits 15000, got %d", config.Performance.MaxCommits)
	}
}

// TestComplexFilteringScenarios tests complex filtering scenarios
func TestComplexFilteringScenarios(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	// Create a diverse set of commits for testing
	commits := []models.Commit{
		{
			Hash:       "abc123",
			Message:    "feat: add authentication system",
			Author:     models.Author{Name: "John Doe", Email: "john.doe@company.com"},
			AuthorDate: baseTime.AddDate(0, 0, -10),
			Stats:      models.CommitStats{FilesChanged: 5, Insertions: 150, Deletions: 20},
		},
		{
			Hash:       "def456",
			Message:    "fix: resolve login bug #123",
			Author:     models.Author{Name: "Jane Smith", Email: "jane.smith@company.com"},
			AuthorDate: baseTime.AddDate(0, 0, -8),
			Stats:      models.CommitStats{FilesChanged: 2, Insertions: 30, Deletions: 15},
		},
		{
			Hash:       "ghi789",
			Message:    "docs: update API documentation",
			Author:     models.Author{Name: "Bob Wilson", Email: "bob.wilson@external.org"},
			AuthorDate: baseTime.AddDate(0, 0, -5),
			Stats:      models.CommitStats{FilesChanged: 3, Insertions: 80, Deletions: 10},
		},
		{
			Hash:       "jkl012",
			Message:    "refactor: improve code structure",
			Author:     models.Author{Name: "Alice Brown", Email: "alice.brown@company.com"},
			AuthorDate: baseTime.AddDate(0, 0, -3),
			Stats:      models.CommitStats{FilesChanged: 8, Insertions: 200, Deletions: 100},
		},
		{
			Hash:       "mno345",
			Message:    "test: add unit tests for auth",
			Author:     models.Author{Name: "John Doe", Email: "john.doe@company.com"},
			AuthorDate: baseTime.AddDate(0, 0, -1),
			Stats:      models.CommitStats{FilesChanged: 4, Insertions: 120, Deletions: 5},
		},
	}

	tests := []struct {
		name          string
		cliConfig     *cli.Config
		expectedCount int
		description   string
	}{
		{
			name: "Company email domain filter",
			cliConfig: &cli.Config{
				Command: "contrib",
				Author:  "@company.com",
				Format:  "terminal",
			},
			expectedCount: 4,
			description:   "Should filter commits from company domain",
		},
		{
			name: "Date range with author filter",
			cliConfig: &cli.Config{
				Command: "contrib",
				Since:   func() *time.Time { t := baseTime.AddDate(0, 0, -7); return &t }(),
				Author:  "john",
				Format:  "terminal",
			},
			expectedCount: 1,
			description:   "Should combine date and author filters",
		},
		{
			name: "Recent commits only",
			cliConfig: &cli.Config{
				Command: "contrib",
				Since:   func() *time.Time { t := baseTime.AddDate(0, 0, -4); return &t }(),
				Format:  "terminal",
			},
			expectedCount: 2, // Only commits from -3 and -1 days (within last 4 days)
			description:   "Should filter recent commits",
		},
		{
			name: "Limited results",
			cliConfig: &cli.Config{
				Command: "contrib",
				Limit:   2,
				Format:  "terminal",
			},
			expectedCount: 2,
			description:   "Should limit number of results",
		},
		{
			name: "Complex combination",
			cliConfig: &cli.Config{
				Command: "contrib",
				Since:   func() *time.Time { t := baseTime.AddDate(0, 0, -12); return &t }(),
				Until:   func() *time.Time { t := baseTime.AddDate(0, 0, -2); return &t }(),
				Author:  "@company.com",
				Limit:   10,
				Format:  "terminal",
			},
			expectedCount: 3,
			description:   "Should apply all filters in combination",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, err := fcm.ApplyFilters(commits, tt.cliConfig)
			if err != nil {
				t.Fatalf("Failed to apply filters: %v", err)
			}

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d commits, got %d. %s", tt.expectedCount, len(filtered), tt.description)
			}
		})
	}
}
