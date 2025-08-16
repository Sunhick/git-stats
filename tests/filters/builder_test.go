// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for filter builder

package filters

import (
	"testing"
	"time"

	"git-stats/cli"
	"git-stats/config"
	"git-stats/filters"
)

func TestFilterBuilder_BuildFromCLIConfig(t *testing.T) {
	configManager := config.NewConfigManager()
	builder := filters.NewFilterBuilder(configManager)

	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		cliConfig      *cli.Config
		expectedFilters int
	}{
		{
			name: "Empty CLI config",
			cliConfig: &cli.Config{
				Command: "contrib",
				Format:  "terminal",
			},
			expectedFilters: 1, // only merge filter (limit filter only added if limit > 0)
		},
		{
			name: "CLI config with date range",
			cliConfig: &cli.Config{
				Command: "contrib",
				Since:   &baseTime,
				Until:   func() *time.Time { t := baseTime.AddDate(0, 0, 7); return &t }(),
				Format:  "terminal",
			},
			expectedFilters: 2, // date filter + merge filter
		},
		{
			name: "CLI config with author",
			cliConfig: &cli.Config{
				Command: "contrib",
				Author:  "John Doe",
				Format:  "terminal",
			},
			expectedFilters: 2, // author filter + merge filter
		},
		{
			name: "CLI config with all filters",
			cliConfig: &cli.Config{
				Command: "contrib",
				Since:   &baseTime,
				Until:   func() *time.Time { t := baseTime.AddDate(0, 0, 7); return &t }(),
				Author:  "John Doe",
				Limit:   1000,
				Format:  "terminal",
			},
			expectedFilters: 4, // date + author + merge + limit filters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain, err := builder.BuildFromCLIConfig(tt.cliConfig)
			if err != nil {
				t.Fatalf("Failed to build filter chain: %v", err)
			}

			filterCount := len(chain.GetFilters())
			if filterCount != tt.expectedFilters {
				t.Errorf("Expected %d filters, got %d", tt.expectedFilters, filterCount)
			}
		})
	}
}

func TestFilterBuilder_BuildFromConfig(t *testing.T) {
	configManager := config.NewConfigManager()

	// Customize config
	appConfig := configManager.GetConfig()
	appConfig.Defaults.DateRange = "1 month ago"
	appConfig.Filters.DefaultAuthor = "test@example.com"
	appConfig.Filters.IncludeMerges = false
	appConfig.Performance.MaxCommits = 5000
	configManager.SetConfig(appConfig)

	builder := filters.NewFilterBuilder(configManager)

	chain, err := builder.BuildFromConfig()
	if err != nil {
		t.Fatalf("Failed to build filter chain from config: %v", err)
	}

	// Should have: date filter + author filter + merge filter + limit filter
	expectedFilters := 4
	filterCount := len(chain.GetFilters())
	if filterCount != expectedFilters {
		t.Errorf("Expected %d filters, got %d", expectedFilters, filterCount)
	}
}

func TestFilterBuilder_ParseDateRange(t *testing.T) {
	configManager := config.NewConfigManager()
	builder := filters.NewFilterBuilder(configManager)

	tests := []struct {
		name      string
		dateRange string
		shouldErr bool
	}{
		{
			name:      "Today",
			dateRange: "today",
			shouldErr: false,
		},
		{
			name:      "Yesterday",
			dateRange: "yesterday",
			shouldErr: false,
		},
		{
			name:      "This week",
			dateRange: "this week",
			shouldErr: false,
		},
		{
			name:      "This month",
			dateRange: "this month",
			shouldErr: false,
		},
		{
			name:      "This year",
			dateRange: "this year",
			shouldErr: false,
		},
		{
			name:      "1 day ago",
			dateRange: "1 day ago",
			shouldErr: false,
		},
		{
			name:      "2 weeks ago",
			dateRange: "2 weeks ago",
			shouldErr: false,
		},
		{
			name:      "3 months ago",
			dateRange: "3 months ago",
			shouldErr: false,
		},
		{
			name:      "1 year ago",
			dateRange: "1 year ago",
			shouldErr: false,
		},
		{
			name:      "a week ago",
			dateRange: "a week ago",
			shouldErr: false,
		},
		{
			name:      "an hour ago",
			dateRange: "an hour ago",
			shouldErr: true, // hours not supported
		},
		{
			name:      "invalid format",
			dateRange: "invalid date",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use reflection to access private method for testing
			// In a real implementation, you might want to make this method public for testing
			// or create a separate utility function

			// For now, we'll test through BuildFromConfig with the date range
			appConfig := configManager.GetConfig()
			appConfig.Defaults.DateRange = tt.dateRange
			configManager.SetConfig(appConfig)

			_, err := builder.BuildFromConfig()

			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestFilterBuilder_BuildAdvancedFilter(t *testing.T) {
	configManager := config.NewConfigManager()
	builder := filters.NewFilterBuilder(configManager)

	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	options := filters.AdvancedFilterOptions{
		Since: &baseTime,
		Until: func() *time.Time { t := baseTime.AddDate(0, 0, 7); return &t }(),
		Authors: []filters.AuthorFilterOptions{
			{
				Pattern:       "John",
				MatchType:     filters.ContainsMatch,
				CaseSensitive: false,
			},
			{
				Pattern:       "jane@example.com",
				MatchType:     filters.EmailMatch,
				CaseSensitive: false,
			},
		},
		IncludeFiles:  []string{"*.go", "*.md"},
		ExcludeFiles:  []string{"test_*"},
		FileMatchType: filters.FileGlobMatch,
		CaseSensitive: false,
		IncludeMerges: true,
		Limit:         1000,
	}

	chain, err := builder.BuildAdvancedFilter(options)
	if err != nil {
		t.Fatalf("Failed to build advanced filter: %v", err)
	}

	// Should have: date + 2 authors + include files + exclude files + merge + limit = 7 filters
	expectedFilters := 7
	filterCount := len(chain.GetFilters())
	if filterCount != expectedFilters {
		t.Errorf("Expected %d filters, got %d", expectedFilters, filterCount)
	}
}

func TestFilterBuilder_GetFilterSummary(t *testing.T) {
	configManager := config.NewConfigManager()
	builder := filters.NewFilterBuilder(configManager)

	tests := []struct {
		name     string
		chain    *filters.FilterChain
		expected string
	}{
		{
			name:     "Empty chain",
			chain:    filters.NewFilterChain(),
			expected: "No filters applied",
		},
		{
			name:     "Nil chain",
			chain:    nil,
			expected: "No filters applied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := builder.GetFilterSummary(tt.chain)
			if summary != tt.expected {
				t.Errorf("Expected summary '%s', got '%s'", tt.expected, summary)
			}
		})
	}

	// Test with actual filters
	chain := filters.NewFilterChain()
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	chain.Add(filters.NewDateRangeFilter(&baseTime, nil))
	chain.Add(filters.NewAuthorFilter("John"))

	summary := builder.GetFilterSummary(chain)
	if !contains(summary, "Since: 2024-01-01") {
		t.Error("Summary should contain date filter description")
	}
	if !contains(summary, "Author contains match: 'John'") {
		t.Error("Summary should contain author filter description")
	}
}

func TestFilterBuilder_AuthorMatchTypes(t *testing.T) {
	configManager := config.NewConfigManager()
	builder := filters.NewFilterBuilder(configManager)

	tests := []struct {
		name      string
		matchType string
		expected  filters.AuthorMatchType
	}{
		{
			name:      "exact",
			matchType: "exact",
			expected:  filters.ExactMatch,
		},
		{
			name:      "contains",
			matchType: "contains",
			expected:  filters.ContainsMatch,
		},
		{
			name:      "regex",
			matchType: "regex",
			expected:  filters.RegexMatch,
		},
		{
			name:      "email",
			matchType: "email",
			expected:  filters.EmailMatch,
		},
		{
			name:      "name",
			matchType: "name",
			expected:  filters.NameMatch,
		},
		{
			name:      "invalid",
			matchType: "invalid",
			expected:  filters.ContainsMatch, // should default to contains
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the match type in config
			appConfig := configManager.GetConfig()
			appConfig.Filters.AuthorMatchType = tt.matchType
			configManager.SetConfig(appConfig)

			// Build CLI config with author
			cliConfig := &cli.Config{
				Command: "contrib",
				Author:  "test",
				Format:  "terminal",
			}

			chain, err := builder.BuildFromCLIConfig(cliConfig)
			if err != nil {
				t.Fatalf("Failed to build filter chain: %v", err)
			}

			// Find the author filter and check its match type
			found := false
			for _, filter := range chain.GetFilters() {
				if authorFilter, ok := filter.(*filters.AuthorFilter); ok {
					// We can't directly access the match type, but we can test the behavior
					// by checking the description
					description := authorFilter.Description()
					expectedDesc := getExpectedDescription(tt.expected)
					if !contains(description, expectedDesc) {
						t.Errorf("Expected description to contain '%s', got '%s'", expectedDesc, description)
					}
					found = true
					break
				}
			}

			if !found {
				t.Error("Author filter not found in chain")
			}
		})
	}
}

// Helper function to get expected description based on match type
func getExpectedDescription(matchType filters.AuthorMatchType) string {
	switch matchType {
	case filters.ExactMatch:
		return "exact"
	case filters.ContainsMatch:
		return "contains"
	case filters.RegexMatch:
		return "regex"
	case filters.EmailMatch:
		return "email"
	case filters.NameMatch:
		return "name"
	default:
		return "contains"
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
