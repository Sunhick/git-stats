// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Comprehensive filtering tests for task 7.2

package filters

import (
	"testing"
	"time"

	"git-stats/cli"
	"git-stats/config"
	"git-stats/filters"
	"git-stats/integration"
	"git-stats/models"
)

// TestComprehensiveDateRangeFiltering tests all aspects of date range filtering
func TestComprehensiveDateRangeFiltering(t *testing.T) {
	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	// Create test commits spanning different dates
	commits := []models.Commit{
		{
			Hash:       "abc123",
			Message:    "commit 1",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: baseTime.AddDate(0, 0, -10), // 10 days ago
		},
		{
			Hash:       "def456",
			Message:    "commit 2",
			Author:     models.Author{Name: "Jane Smith", Email: "jane@example.com"},
			AuthorDate: baseTime.AddDate(0, 0, -5), // 5 days ago
		},
		{
			Hash:       "ghi789",
			Message:    "commit 3",
			Author:     models.Author{Name: "Bob Wilson", Email: "bob@example.com"},
			AuthorDate: baseTime, // today
		},
		{
			Hash:       "jkl012",
			Message:    "commit 4",
			Author:     models.Author{Name: "Alice Brown", Email: "alice@example.com"},
			AuthorDate: baseTime.AddDate(0, 0, 5), // 5 days in future
		},
	}

	tests := []struct {
		name           string
		since          *time.Time
		until          *time.Time
		expectedCount  int
		description    string
	}{
		{
			name:          "No date filter",
			since:         nil,
			until:         nil,
			expectedCount: 4,
			description:   "Should return all commits when no date filter is applied",
		},
		{
			name:          "Since filter only",
			since:         func() *time.Time { t := baseTime.AddDate(0, 0, -7); return &t }(),
			until:         nil,
			expectedCount: 3,
			description:   "Should return commits from 7 days ago onwards",
		},
		{
			name:          "Until filter only",
			since:         nil,
			until:         func() *time.Time { t := baseTime.AddDate(0, 0, -1); return &t }(),
			expectedCount: 2,
			description:   "Should return commits until 1 day ago",
		},
		{
			name:          "Date range filter",
			since:         func() *time.Time { t := baseTime.AddDate(0, 0, -7); return &t }(),
			until:         func() *time.Time { t := baseTime.AddDate(0, 0, 1); return &t }(),
			expectedCount: 2,
			description:   "Should return commits within date range",
		},
		{
			name:          "Exact date match",
			since:         &baseTime,
			until:         &baseTime,
			expectedCount: 1,
			description:   "Should return commits on exact date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := filters.NewDateRangeFilter(tt.since, tt.until)
			filtered := filter.Apply(commits)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d commits, got %d. %s", tt.expectedCount, len(filtered), tt.description)
			}
		})
	}
}

// TestComprehensiveAuthorFiltering tests all aspects of author filtering
func TestComprehensiveAuthorFiltering(t *testing.T) {
	commits := []models.Commit{
		{
			Hash:    "abc123",
			Message: "commit 1",
			Author:  models.Author{Name: "John Doe", Email: "john.doe@example.com"},
		},
		{
			Hash:    "def456",
			Message: "commit 2",
			Author:  models.Author{Name: "Jane Smith", Email: "jane.smith@company.org"},
		},
		{
			Hash:    "ghi789",
			Message: "commit 3",
			Author:  models.Author{Name: "John Wilson", Email: "j.wilson@test.com"},
		},
		{
			Hash:    "jkl012",
			Message: "commit 4",
			Author:  models.Author{Name: "Alice Brown", Email: "alice@example.com"},
		},
	}

	tests := []struct {
		name          string
		pattern       string
		matchType     filters.AuthorMatchType
		caseSensitive bool
		expectedCount int
		description   string
	}{
		{
			name:          "Contains match - name",
			pattern:       "john",
			matchType:     filters.ContainsMatch,
			caseSensitive: false,
			expectedCount: 2,
			description:   "Should match both John Doe and John Wilson",
		},
		{
			name:          "Contains match - email domain",
			pattern:       "example.com",
			matchType:     filters.ContainsMatch,
			caseSensitive: false,
			expectedCount: 2,
			description:   "Should match authors with example.com domain",
		},
		{
			name:          "Exact match - name",
			pattern:       "John Doe",
			matchType:     filters.ExactMatch,
			caseSensitive: false,
			expectedCount: 1,
			description:   "Should match only exact name",
		},
		{
			name:          "Email only match",
			pattern:       "jane.smith",
			matchType:     filters.EmailMatch,
			caseSensitive: false,
			expectedCount: 1,
			description:   "Should match only in email field",
		},
		{
			name:          "Name only match",
			pattern:       "Alice",
			matchType:     filters.NameMatch,
			caseSensitive: false,
			expectedCount: 1,
			description:   "Should match only in name field",
		},
		{
			name:          "Case sensitive match",
			pattern:       "JOHN",
			matchType:     filters.ContainsMatch,
			caseSensitive: true,
			expectedCount: 0,
			description:   "Should not match with case sensitive enabled",
		},
		{
			name:          "Case insensitive match",
			pattern:       "JOHN",
			matchType:     filters.ContainsMatch,
			caseSensitive: false,
			expectedCount: 2,
			description:   "Should match with case insensitive enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := filters.NewAuthorFilterWithOptions(tt.pattern, tt.matchType, tt.caseSensitive)
			if err != nil {
				t.Fatalf("Failed to create author filter: %v", err)
			}

			filtered := filter.Apply(commits)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d commits, got %d. %s", tt.expectedCount, len(filtered), tt.description)
			}
		})
	}
}

// TestRegexAuthorFiltering tests regex-based author filtering
func TestRegexAuthorFiltering(t *testing.T) {
	commits := []models.Commit{
		{
			Hash:    "abc123",
			Message: "commit 1",
			Author:  models.Author{Name: "John Doe", Email: "john.doe@example.com"},
		},
		{
			Hash:    "def456",
			Message: "commit 2",
			Author:  models.Author{Name: "Jane Smith", Email: "jane.smith@company.org"},
		},
		{
			Hash:    "ghi789",
			Message: "commit 3",
			Author:  models.Author{Name: "Bob123", Email: "bob@test.com"},
		},
	}

	tests := []struct {
		name          string
		pattern       string
		expectedCount int
		shouldError   bool
		description   string
	}{
		{
			name:          "Valid regex - email domains",
			pattern:       `@(example\.com|company\.org)$`,
			expectedCount: 2,
			shouldError:   false,
			description:   "Should match specific email domains",
		},
		{
			name:          "Valid regex - names with numbers",
			pattern:       `\d+`,
			expectedCount: 1,
			shouldError:   false,
			description:   "Should match names containing numbers",
		},
		{
			name:          "Valid regex - names starting with J",
			pattern:       `^J`,
			expectedCount: 2,
			shouldError:   false,
			description:   "Should match names starting with J",
		},
		{
			name:          "Invalid regex",
			pattern:       `[`,
			expectedCount: 0,
			shouldError:   true,
			description:   "Should error on invalid regex pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := filters.NewAuthorFilterWithOptions(tt.pattern, filters.RegexMatch, false)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error for invalid regex but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error creating regex filter: %v", err)
			}

			filtered := filter.Apply(commits)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d commits, got %d. %s", tt.expectedCount, len(filtered), tt.description)
			}
		})
	}
}

// TestMultipleFiltersIntegration tests combining multiple filters
func TestMultipleFiltersIntegration(t *testing.T) {
	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	commits := []models.Commit{
		{
			Hash:       "abc123",
			Message:    "commit 1",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: baseTime.AddDate(0, 0, -10),
			Stats:      models.CommitStats{FilesChanged: 2, Insertions: 50, Deletions: 10},
		},
		{
			Hash:       "def456",
			Message:    "commit 2",
			Author:     models.Author{Name: "John Smith", Email: "john@company.org"},
			AuthorDate: baseTime.AddDate(0, 0, -5),
			Stats:      models.CommitStats{FilesChanged: 1, Insertions: 20, Deletions: 5},
		},
		{
			Hash:       "ghi789",
			Message:    "commit 3",
			Author:     models.Author{Name: "Jane Doe", Email: "jane@example.com"},
			AuthorDate: baseTime,
			Stats:      models.CommitStats{FilesChanged: 3, Insertions: 75, Deletions: 15},
		},
		{
			Hash:       "jkl012",
			Message:    "merge commit",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: baseTime.AddDate(0, 0, 2),
			ParentHashes: []string{"abc123", "def456"}, // This makes it a merge commit
			Stats:      models.CommitStats{FilesChanged: 0, Insertions: 0, Deletions: 0},
		},
	}

	fcm := integration.NewFilteredConfigManager()

	tests := []struct {
		name          string
		cliConfig     *cli.Config
		expectedCount int
		description   string
	}{
		{
			name: "Date and author filters",
			cliConfig: &cli.Config{
				Command: "contrib",
				Since:   func() *time.Time { t := baseTime.AddDate(0, 0, -7); return &t }(),
				Author:  "john",
				Format:  "terminal",
			},
			expectedCount: 2,
			description:   "Should filter by date range and author",
		},
		{
			name: "Date, author, and limit filters",
			cliConfig: &cli.Config{
				Command: "contrib",
				Since:   func() *time.Time { t := baseTime.AddDate(0, 0, -15); return &t }(),
				Author:  "john",
				Limit:   1,
				Format:  "terminal",
			},
			expectedCount: 1,
			description:   "Should apply all filters including limit",
		},
		{
			name: "All filters with no matches",
			cliConfig: &cli.Config{
				Command: "contrib",
				Since:   func() *time.Time { t := baseTime.AddDate(0, 0, 10); return &t }(),
				Author:  "nonexistent",
				Format:  "terminal",
			},
			expectedCount: 0,
			description:   "Should return no results when filters don't match",
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

// TestConfigurationManagement tests configuration management features
func TestConfigurationManagement(t *testing.T) {
	fcm := integration.NewFilteredConfigManager()

	// Test default configuration
	t.Run("Default configuration", func(t *testing.T) {
		config := fcm.GetConfigManager().GetConfig()

		if config.Defaults.Command != "contrib" {
			t.Errorf("Expected default command 'contrib', got '%s'", config.Defaults.Command)
		}

		if config.Defaults.Format != "terminal" {
			t.Errorf("Expected default format 'terminal', got '%s'", config.Defaults.Format)
		}

		if config.Filters.AuthorMatchType != "contains" {
			t.Errorf("Expected default author match type 'contains', got '%s'", config.Filters.AuthorMatchType)
		}

		if !config.Filters.IncludeMerges {
			t.Error("Expected default include merges to be true")
		}
	})

	// Test configuration updates
	t.Run("Configuration updates", func(t *testing.T) {
		configManager := fcm.GetConfigManager()

		// Update defaults
		newDefaults := config.DefaultConfig{
			Command:      "summary",
			DateRange:    "6 months ago",
			Format:       "json",
			ShowProgress: true,
		}
		configManager.UpdateDefaults(newDefaults)

		// Update filters
		newFilters := config.FilterConfig{
			IncludeMerges:   false,
			DefaultAuthor:   "test@example.com",
			AuthorMatchType: "exact",
			CaseSensitive:   true,
		}
		configManager.UpdateFilters(newFilters)

		// Verify updates
		updatedConfig := configManager.GetConfig()
		if updatedConfig.Defaults.Command != "summary" {
			t.Errorf("Expected updated command 'summary', got '%s'", updatedConfig.Defaults.Command)
		}

		if updatedConfig.Filters.DefaultAuthor != "test@example.com" {
			t.Errorf("Expected updated author 'test@example.com', got '%s'", updatedConfig.Filters.DefaultAuthor)
		}

		if updatedConfig.Filters.IncludeMerges {
			t.Error("Expected updated include merges to be false")
		}
	})

	// Test effective configuration
	t.Run("Effective configuration", func(t *testing.T) {
		cliConfig := &cli.Config{
			Command:    "contributors",
			Author:     "override@example.com",
			Format:     "csv",
			Limit:      5000,
			NoColor:    true,
		}

		effectiveConfig := fcm.GetEffectiveConfig(cliConfig)

		// CLI overrides should take precedence
		if effectiveConfig.Command != "contributors" {
			t.Errorf("Expected effective command 'contributors', got '%s'", effectiveConfig.Command)
		}

		if effectiveConfig.Author != "override@example.com" {
			t.Errorf("Expected effective author 'override@example.com', got '%s'", effectiveConfig.Author)
		}

		if effectiveConfig.Format != "csv" {
			t.Errorf("Expected effective format 'csv', got '%s'", effectiveConfig.Format)
		}

		if !effectiveConfig.NoColor {
			t.Error("Expected effective no color to be true")
		}
	})
}

// TestDateFormatParsing tests various date format parsing
func TestDateFormatParsing(t *testing.T) {
	configManager := config.NewConfigManager()
	builder := filters.NewFilterBuilder(configManager)

	// Test various date formats through the filter builder
	tests := []struct {
		name        string
		dateRange   string
		shouldError bool
		description string
	}{
		{
			name:        "Today",
			dateRange:   "today",
			shouldError: false,
			description: "Should parse 'today' correctly",
		},
		{
			name:        "Yesterday",
			dateRange:   "yesterday",
			shouldError: false,
			description: "Should parse 'yesterday' correctly",
		},
		{
			name:        "This week",
			dateRange:   "this week",
			shouldError: false,
			description: "Should parse 'this week' correctly",
		},
		{
			name:        "This month",
			dateRange:   "this month",
			shouldError: false,
			description: "Should parse 'this month' correctly",
		},
		{
			name:        "This year",
			dateRange:   "this year",
			shouldError: false,
			description: "Should parse 'this year' correctly",
		},
		{
			name:        "Relative - 1 day ago",
			dateRange:   "1 day ago",
			shouldError: false,
			description: "Should parse '1 day ago' correctly",
		},
		{
			name:        "Relative - 2 weeks ago",
			dateRange:   "2 weeks ago",
			shouldError: false,
			description: "Should parse '2 weeks ago' correctly",
		},
		{
			name:        "Relative - 3 months ago",
			dateRange:   "3 months ago",
			shouldError: false,
			description: "Should parse '3 months ago' correctly",
		},
		{
			name:        "Relative - 1 year ago",
			dateRange:   "1 year ago",
			shouldError: false,
			description: "Should parse '1 year ago' correctly",
		},
		{
			name:        "Relative - a week ago",
			dateRange:   "a week ago",
			shouldError: false,
			description: "Should parse 'a week ago' correctly",
		},
		{
			name:        "Invalid format",
			dateRange:   "invalid date format",
			shouldError: true,
			description: "Should error on invalid date format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the date range in config and try to build filter
			appConfig := configManager.GetConfig()
			appConfig.Defaults.DateRange = tt.dateRange
			configManager.SetConfig(appConfig)

			_, err := builder.BuildFromConfig()

			if tt.shouldError && err == nil {
				t.Errorf("Expected error for '%s' but got none. %s", tt.dateRange, tt.description)
			}

			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error for '%s': %v. %s", tt.dateRange, err, tt.description)
			}
		})
	}
}
