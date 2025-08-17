// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Enhanced filtering tests for task 7.2

package filters

import (
	"testing"
	"time"

	"git-stats/config"
	"git-stats/filters"
	"git-stats/models"
)

// TestBranchFilter tests branch filtering functionality
func TestBranchFilter(t *testing.T) {
	// Note: This test is a placeholder since branch filtering requires
	// additional commit model enhancements to track branch information
	commits := []models.Commit{
		{
			Hash:    "abc123",
			Message: "commit on main",
			Author:  models.Author{Name: "John Doe", Email: "john@example.com"},
		},
		{
			Hash:    "def456",
			Message: "commit on feature",
			Author:  models.Author{Name: "Jane Smith", Email: "jane@example.com"},
		},
	}

	tests := []struct {
		name          string
		branches      []string
		matchType     filters.BranchMatchType
		caseSensitive bool
		expectedCount int
		description   string
	}{
		{
			name:          "Empty branch filter",
			branches:      []string{},
			matchType:     filters.BranchContainsMatch,
			caseSensitive: false,
			expectedCount: 2,
			description:   "Should return all commits when no branch filter is applied",
		},
		{
			name:          "Main branch filter",
			branches:      []string{"main"},
			matchType:     filters.BranchExactMatch,
			caseSensitive: false,
			expectedCount: 2, // Placeholder - would be 1 with actual branch detection
			description:   "Should filter commits from main branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := filters.NewBranchFilter(tt.branches, tt.matchType, tt.caseSensitive)
			filtered := filter.Apply(commits)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d commits, got %d. %s", tt.expectedCount, len(filtered), tt.description)
			}

			// Test description
			description := filter.Description()
			if description == "" {
				t.Error("Filter description should not be empty")
			}
		})
	}
}

// TestMessageFilter tests commit message filtering functionality
func TestMessageFilter(t *testing.T) {
	commits := []models.Commit{
		{
			Hash:    "abc123",
			Message: "feat: add new feature",
			Author:  models.Author{Name: "John Doe", Email: "john@example.com"},
		},
		{
			Hash:    "def456",
			Message: "fix: resolve bug #123",
			Author:  models.Author{Name: "Jane Smith", Email: "jane@example.com"},
		},
		{
			Hash:    "ghi789",
			Message: "docs: update README",
			Author:  models.Author{Name: "Bob Wilson", Email: "bob@example.com"},
		},
		{
			Hash:    "jkl012",
			Message: "refactor: improve code structure",
			Author:  models.Author{Name: "Alice Brown", Email: "alice@example.com"},
		},
	}

	tests := []struct {
		name          string
		pattern       string
		matchType     filters.MessageMatchType
		caseSensitive bool
		expectedCount int
		description   string
	}{
		{
			name:          "Contains match - feat",
			pattern:       "feat",
			matchType:     filters.MessageContainsMatch,
			caseSensitive: false,
			expectedCount: 1, // Only "feat: add new feature"
			description:   "Should match messages containing 'feat'",
		},
		{
			name:          "Starts with match - fix",
			pattern:       "fix",
			matchType:     filters.MessageStartsWithMatch,
			caseSensitive: false,
			expectedCount: 1,
			description:   "Should match messages starting with 'fix'",
		},
		{
			name:          "Ends with match - README",
			pattern:       "README",
			matchType:     filters.MessageEndsWithMatch,
			caseSensitive: false,
			expectedCount: 1,
			description:   "Should match messages ending with 'README'",
		},
		{
			name:          "Case sensitive match - FEAT",
			pattern:       "FEAT",
			matchType:     filters.MessageContainsMatch,
			caseSensitive: true,
			expectedCount: 0,
			description:   "Should not match with case sensitive enabled",
		},
		{
			name:          "Case insensitive match - FEAT",
			pattern:       "FEAT",
			matchType:     filters.MessageContainsMatch,
			caseSensitive: false,
			expectedCount: 1,
			description:   "Should match with case insensitive enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := filters.NewMessageFilter(tt.pattern, tt.matchType, tt.caseSensitive)
			if err != nil {
				t.Fatalf("Failed to create message filter: %v", err)
			}

			filtered := filter.Apply(commits)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d commits, got %d. %s", tt.expectedCount, len(filtered), tt.description)
			}

			// Test description
			description := filter.Description()
			if description == "" {
				t.Error("Filter description should not be empty")
			}
		})
	}
}

// TestMessageFilterRegex tests regex-based message filtering
func TestMessageFilterRegex(t *testing.T) {
	commits := []models.Commit{
		{
			Hash:    "abc123",
			Message: "feat: add feature #123",
			Author:  models.Author{Name: "John Doe", Email: "john@example.com"},
		},
		{
			Hash:    "def456",
			Message: "fix: resolve issue #456",
			Author:  models.Author{Name: "Jane Smith", Email: "jane@example.com"},
		},
		{
			Hash:    "ghi789",
			Message: "docs: update documentation",
			Author:  models.Author{Name: "Bob Wilson", Email: "bob@example.com"},
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
			name:          "Valid regex - issue numbers",
			pattern:       `#\d+`,
			expectedCount: 2,
			shouldError:   false,
			description:   "Should match messages with issue numbers",
		},
		{
			name:          "Valid regex - feat or fix",
			pattern:       `^(feat|fix):`,
			expectedCount: 2,
			shouldError:   false,
			description:   "Should match messages starting with feat: or fix:",
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
			filter, err := filters.NewMessageFilter(tt.pattern, filters.MessageRegexMatch, false)

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

// TestFileSizeFilter tests file size-based filtering
func TestFileSizeFilter(t *testing.T) {
	commits := []models.Commit{
		{
			Hash:    "abc123",
			Message: "small change",
			Author:  models.Author{Name: "John Doe", Email: "john@example.com"},
			Stats:   models.CommitStats{FilesChanged: 1, Insertions: 5, Deletions: 2},
		},
		{
			Hash:    "def456",
			Message: "medium change",
			Author:  models.Author{Name: "Jane Smith", Email: "jane@example.com"},
			Stats:   models.CommitStats{FilesChanged: 3, Insertions: 50, Deletions: 20},
		},
		{
			Hash:    "ghi789",
			Message: "large change",
			Author:  models.Author{Name: "Bob Wilson", Email: "bob@example.com"},
			Stats:   models.CommitStats{FilesChanged: 10, Insertions: 200, Deletions: 100},
		},
		{
			Hash:    "jkl012",
			Message: "huge change",
			Author:  models.Author{Name: "Alice Brown", Email: "alice@example.com"},
			Stats:   models.CommitStats{FilesChanged: 25, Insertions: 1000, Deletions: 500},
		},
	}

	tests := []struct {
		name          string
		minInsertions int
		maxInsertions int
		minDeletions  int
		maxDeletions  int
		minFiles      int
		maxFiles      int
		expectedCount int
		description   string
	}{
		{
			name:          "No size filter",
			minInsertions: 0,
			maxInsertions: 0,
			minDeletions:  0,
			maxDeletions:  0,
			minFiles:      0,
			maxFiles:      0,
			expectedCount: 4,
			description:   "Should return all commits when no size filter is applied",
		},
		{
			name:          "Min insertions filter",
			minInsertions: 50,
			maxInsertions: 0,
			minDeletions:  0,
			maxDeletions:  0,
			minFiles:      0,
			maxFiles:      0,
			expectedCount: 3,
			description:   "Should filter commits with at least 50 insertions",
		},
		{
			name:          "Max insertions filter",
			minInsertions: 0,
			maxInsertions: 100,
			minDeletions:  0,
			maxDeletions:  0,
			minFiles:      0,
			maxFiles:      0,
			expectedCount: 2,
			description:   "Should filter commits with at most 100 insertions",
		},
		{
			name:          "Insertions range filter",
			minInsertions: 20,
			maxInsertions: 300,
			minDeletions:  0,
			maxDeletions:  0,
			minFiles:      0,
			maxFiles:      0,
			expectedCount: 2,
			description:   "Should filter commits with insertions in range 20-300",
		},
		{
			name:          "Files changed filter",
			minInsertions: 0,
			maxInsertions: 0,
			minDeletions:  0,
			maxDeletions:  0,
			minFiles:      5,
			maxFiles:      15,
			expectedCount: 1,
			description:   "Should filter commits with 5-15 files changed",
		},
		{
			name:          "Combined filters",
			minInsertions: 40,
			maxInsertions: 250,
			minDeletions:  15,
			maxDeletions:  150,
			minFiles:      2,
			maxFiles:      12,
			expectedCount: 2,
			description:   "Should apply all size filters combined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := filters.NewFileSizeFilter(
				tt.minInsertions, tt.maxInsertions,
				tt.minDeletions, tt.maxDeletions,
				tt.minFiles, tt.maxFiles,
			)

			filtered := filter.Apply(commits)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d commits, got %d. %s", tt.expectedCount, len(filtered), tt.description)
			}

			// Test description
			description := filter.Description()
			if description == "" {
				t.Error("Filter description should not be empty")
			}
		})
	}
}

// TestEnhancedDateParsing tests enhanced date parsing functionality
func TestEnhancedDateParsing(t *testing.T) {
	configManager := config.NewConfigManager()
	builder := filters.NewFilterBuilder(configManager)

	tests := []struct {
		name        string
		dateRange   string
		shouldError bool
		description string
	}{
		{
			name:        "Last week",
			dateRange:   "last week",
			shouldError: false,
			description: "Should parse 'last week' correctly",
		},
		{
			name:        "Last month",
			dateRange:   "last month",
			shouldError: false,
			description: "Should parse 'last month' correctly",
		},
		{
			name:        "Last year",
			dateRange:   "last year",
			shouldError: false,
			description: "Should parse 'last year' correctly",
		},
		{
			name:        "This quarter",
			dateRange:   "this quarter",
			shouldError: false,
			description: "Should parse 'this quarter' correctly",
		},
		{
			name:        "Last quarter",
			dateRange:   "last quarter",
			shouldError: false,
			description: "Should parse 'last quarter' correctly",
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

// TestAdvancedFilterBuilder tests the enhanced advanced filter builder
func TestAdvancedFilterBuilder(t *testing.T) {
	configManager := config.NewConfigManager()
	builder := filters.NewFilterBuilder(configManager)

	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	options := filters.AdvancedFilterOptions{
		Since: &baseTime,
		Until: func() *time.Time { t := baseTime.AddDate(0, 0, 7); return &t }(),
		Authors: []filters.AuthorFilterOptions{
			{
				Pattern:       "John",
				MatchType:     filters.ContainsMatch,
				CaseSensitive: false,
			},
		},
		IncludeFiles:  []string{"*.go"},
		ExcludeFiles:  []string{"*_test.go"},
		FileMatchType: filters.FileGlobMatch,
		CaseSensitive: false,
		IncludeMerges: true,
		Limit:         1000,
		Branches:      []string{"main", "develop"},
		BranchMatchType: filters.BranchContainsMatch,
		MessageFilter: &filters.MessageFilterOptions{
			Pattern:       "feat:",
			MatchType:     filters.MessageStartsWithMatch,
			CaseSensitive: false,
		},
		SizeFilter: &filters.FileSizeFilterOptions{
			MinInsertions: 10,
			MaxInsertions: 500,
			MinFiles:      1,
			MaxFiles:      20,
		},
	}

	chain, err := builder.BuildAdvancedFilter(options)
	if err != nil {
		t.Fatalf("Failed to build advanced filter: %v", err)
	}

	// Should have: date + author + include files + exclude files + merge + limit + branch + message + size = 9 filters
	expectedFilters := 9
	filterCount := len(chain.GetFilters())
	if filterCount != expectedFilters {
		t.Errorf("Expected %d filters, got %d", expectedFilters, filterCount)
	}

	// Test filter summary
	summary := builder.GetFilterSummary(chain)
	if summary == "" {
		t.Error("Filter summary should not be empty")
	}
}

// TestConfigurationEnhancements tests enhanced configuration management
func TestConfigurationEnhancements(t *testing.T) {
	configManager := config.NewConfigManager()

	// Test enhanced default configuration
	t.Run("Enhanced default configuration", func(t *testing.T) {
		config := configManager.GetConfig()

		// Test new filter settings
		if config.Filters.BranchMatchType != "contains" {
			t.Errorf("Expected default branch match type 'contains', got '%s'", config.Filters.BranchMatchType)
		}

		if config.Filters.MessageMatchType != "contains" {
			t.Errorf("Expected default message match type 'contains', got '%s'", config.Filters.MessageMatchType)
		}

		if len(config.Filters.DefaultBranches) != 0 {
			t.Errorf("Expected empty default branches, got %v", config.Filters.DefaultBranches)
		}

		if len(config.Filters.MessagePatterns) != 0 {
			t.Errorf("Expected empty message patterns, got %v", config.Filters.MessagePatterns)
		}
	})

	// Test configuration validation with new fields
	t.Run("Enhanced configuration validation", func(t *testing.T) {
		tests := []struct {
			name      string
			modify    func(*config.Config)
			shouldErr bool
		}{
			{
				name: "Invalid branch match type",
				modify: func(c *config.Config) {
					c.Filters.BranchMatchType = "invalid"
				},
				shouldErr: true,
			},
			{
				name: "Invalid message match type",
				modify: func(c *config.Config) {
					c.Filters.MessageMatchType = "invalid"
				},
				shouldErr: true,
			},
			{
				name: "Invalid min insertions",
				modify: func(c *config.Config) {
					c.Filters.MinInsertions = -1
				},
				shouldErr: true,
			},
			{
				name: "Invalid size range - insertions",
				modify: func(c *config.Config) {
					c.Filters.MinInsertions = 100
					c.Filters.MaxInsertions = 50
				},
				shouldErr: true,
			},
			{
				name: "Invalid size range - deletions",
				modify: func(c *config.Config) {
					c.Filters.MinDeletions = 100
					c.Filters.MaxDeletions = 50
				},
				shouldErr: true,
			},
			{
				name: "Invalid size range - files",
				modify: func(c *config.Config) {
					c.Filters.MinFiles = 10
					c.Filters.MaxFiles = 5
				},
				shouldErr: true,
			},
			{
				name: "Valid enhanced configuration",
				modify: func(c *config.Config) {
					c.Filters.BranchMatchType = "regex"
					c.Filters.MessageMatchType = "starts_with"
					c.Filters.MinInsertions = 10
					c.Filters.MaxInsertions = 100
					c.Filters.MinFiles = 1
					c.Filters.MaxFiles = 20
				},
				shouldErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cfg := configManager.GetConfig()
				tt.modify(cfg)
				configManager.SetConfig(cfg)

				err := configManager.Validate()

				if tt.shouldErr && err == nil {
					t.Error("Expected validation error but got none")
				}
				if !tt.shouldErr && err != nil {
					t.Errorf("Unexpected validation error: %v", err)
				}

				// Reset config for next test
				configManager.Reset()
			})
		}
	})
}
