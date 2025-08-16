// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for filtering system

package filters

import (
	"testing"
	"time"

	"git-stats/filters"
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
		{
			Hash:       "ghi789",
			Message:    "Fix bug",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: baseTime.AddDate(0, 0, 2),
			Stats: models.CommitStats{
				FilesChanged: 1,
				Insertions:   3,
				Deletions:    1,
				Files: []models.FileChange{
					{Path: "main.go", Status: "M", Insertions: 3, Deletions: 1},
				},
			},
		},
		{
			Hash:       "jkl012",
			Message:    "Merge branch 'feature'",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: baseTime.AddDate(0, 0, 3),
			ParentHashes: []string{"abc123", "def456"}, // Merge commit
			Stats: models.CommitStats{
				FilesChanged: 0,
				Insertions:   0,
				Deletions:    0,
				Files:        []models.FileChange{},
			},
		},
	}
}

func TestFilterChain(t *testing.T) {
	commits := createTestCommits()
	chain := filters.NewFilterChain()

	// Test empty chain
	result := chain.Apply(commits)
	if len(result) != len(commits) {
		t.Errorf("Empty chain should return all commits, got %d, expected %d", len(result), len(commits))
	}

	// Test adding filters
	authorFilter := filters.NewAuthorFilter("John")
	chain.Add(authorFilter)

	result = chain.Apply(commits)
	expectedCount := 3 // John Doe has 3 commits
	if len(result) != expectedCount {
		t.Errorf("Author filter should return %d commits, got %d", expectedCount, len(result))
	}

	// Test clearing chain
	chain.Clear()
	result = chain.Apply(commits)
	if len(result) != len(commits) {
		t.Errorf("Cleared chain should return all commits, got %d, expected %d", len(result), len(commits))
	}
}

func TestDateRangeFilter(t *testing.T) {
	commits := createTestCommits()
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		since    *time.Time
		until    *time.Time
		expected int
	}{
		{
			name:     "No date filter",
			since:    nil,
			until:    nil,
			expected: 4,
		},
		{
			name:     "Since filter",
			since:    &baseTime,
			until:    nil,
			expected: 4,
		},
		{
			name:     "Until filter",
			since:    nil,
			until:    &baseTime,
			expected: 1,
		},
		{
			name:     "Date range filter",
			since:    &baseTime,
			until:    func() *time.Time { t := baseTime.AddDate(0, 0, 1); return &t }(),
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := filters.NewDateRangeFilter(tt.since, tt.until)
			result := filter.Apply(commits)

			if len(result) != tt.expected {
				t.Errorf("Expected %d commits, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestAuthorFilter(t *testing.T) {
	commits := createTestCommits()

	tests := []struct {
		name      string
		pattern   string
		matchType filters.AuthorMatchType
		caseSensitive bool
		expected  int
	}{
		{
			name:      "Contains match - name",
			pattern:   "John",
			matchType: filters.ContainsMatch,
			caseSensitive: false,
			expected:  3,
		},
		{
			name:      "Contains match - email",
			pattern:   "jane@",
			matchType: filters.ContainsMatch,
			caseSensitive: false,
			expected:  1,
		},
		{
			name:      "Exact match - name",
			pattern:   "John Doe",
			matchType: filters.ExactMatch,
			caseSensitive: false,
			expected:  3,
		},
		{
			name:      "Email only match",
			pattern:   "example.com",
			matchType: filters.EmailMatch,
			caseSensitive: false,
			expected:  4,
		},
		{
			name:      "Name only match",
			pattern:   "Smith",
			matchType: filters.NameMatch,
			caseSensitive: false,
			expected:  1,
		},
		{
			name:      "Case sensitive match",
			pattern:   "john",
			matchType: filters.ContainsMatch,
			caseSensitive: true,
			expected:  3, // matches "john@example.com" in email
		},
		{
			name:      "Case insensitive match",
			pattern:   "john",
			matchType: filters.ContainsMatch,
			caseSensitive: false,
			expected:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := filters.NewAuthorFilterWithOptions(tt.pattern, tt.matchType, tt.caseSensitive)
			if err != nil {
				t.Fatalf("Failed to create author filter: %v", err)
			}

			result := filter.Apply(commits)

			if len(result) != tt.expected {
				t.Errorf("Expected %d commits, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestAuthorFilterRegex(t *testing.T) {
	commits := createTestCommits()

	tests := []struct {
		name     string
		pattern  string
		expected int
		shouldErr bool
	}{
		{
			name:     "Valid regex - email domain",
			pattern:  `@example\.com$`,
			expected: 4,
			shouldErr: false,
		},
		{
			name:     "Valid regex - name pattern",
			pattern:  `^John`,
			expected: 3,
			shouldErr: false,
		},
		{
			name:     "Invalid regex",
			pattern:  `[`,
			expected: 0,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := filters.NewAuthorFilterWithOptions(tt.pattern, filters.RegexMatch, false)

			if tt.shouldErr {
				if err == nil {
					t.Error("Expected error for invalid regex, got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			result := filter.Apply(commits)

			if len(result) != tt.expected {
				t.Errorf("Expected %d commits, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestMergeCommitFilter(t *testing.T) {
	commits := createTestCommits()

	tests := []struct {
		name          string
		includeMerges bool
		expected      int
	}{
		{
			name:          "Include merge commits",
			includeMerges: true,
			expected:      4,
		},
		{
			name:          "Exclude merge commits",
			includeMerges: false,
			expected:      3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := filters.NewMergeCommitFilter(tt.includeMerges)
			result := filter.Apply(commits)

			if len(result) != tt.expected {
				t.Errorf("Expected %d commits, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestLimitFilter(t *testing.T) {
	commits := createTestCommits()

	tests := []struct {
		name     string
		limit    int
		expected int
	}{
		{
			name:     "No limit",
			limit:    0,
			expected: 4,
		},
		{
			name:     "Limit less than total",
			limit:    2,
			expected: 2,
		},
		{
			name:     "Limit greater than total",
			limit:    10,
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := filters.NewLimitFilter(tt.limit)
			result := filter.Apply(commits)

			if len(result) != tt.expected {
				t.Errorf("Expected %d commits, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestFilePathFilter(t *testing.T) {
	commits := createTestCommits()

	tests := []struct {
		name      string
		patterns  []string
		matchType filters.FileMatchType
		expected  int
	}{
		{
			name:      "Contains match - go files",
			patterns:  []string{".go"},
			matchType: filters.FileContainsMatch,
			expected:  3, // commits 1, 2, 3 all have .go files
		},
		{
			name:      "Exact match - main.go",
			patterns:  []string{"main.go"},
			matchType: filters.FileExactMatch,
			expected:  2,
		},
		{
			name:      "Glob match - go files",
			patterns:  []string{"*.go"},
			matchType: filters.FileGlobMatch,
			expected:  3, // commits 1, 2, 3 all have .go files
		},
		{
			name:      "Multiple patterns",
			patterns:  []string{"main.go", "README.md"},
			matchType: filters.FileExactMatch,
			expected:  2, // commits 1 and 3 (commit 1 has both, commit 3 has main.go)
		},
		{
			name:      "No matching files",
			patterns:  []string{"nonexistent.txt"},
			matchType: filters.FileExactMatch,
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := filters.NewFilePathFilter(tt.patterns, tt.matchType, false)
			result := filter.Apply(commits)

			if len(result) != tt.expected {
				t.Errorf("Expected %d commits, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestExcludeFilePathFilter(t *testing.T) {
	commits := createTestCommits()

	tests := []struct {
		name      string
		patterns  []string
		matchType filters.FileMatchType
		expected  int
	}{
		{
			name:      "Exclude go files",
			patterns:  []string{".go"},
			matchType: filters.FileContainsMatch,
			expected:  1, // Only merge commit (commit 4) remains
		},
		{
			name:      "Exclude main.go",
			patterns:  []string{"main.go"},
			matchType: filters.FileExactMatch,
			expected:  2, // Only commits with feature.go and merge commit
		},
		{
			name:      "Exclude multiple patterns",
			patterns:  []string{"main.go", "feature.go"},
			matchType: filters.FileExactMatch,
			expected:  1, // Only merge commit (commit 4) remains
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := filters.NewExcludeFilePathFilter(tt.patterns, tt.matchType, false)
			result := filter.Apply(commits)

			if len(result) != tt.expected {
				t.Errorf("Expected %d commits, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestFilterDescriptions(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		filter   filters.Filter
		expected string
	}{
		{
			name:     "Date range filter",
			filter:   filters.NewDateRangeFilter(&baseTime, nil),
			expected: "Since: 2024-01-01",
		},
		{
			name:     "Author filter",
			filter:   filters.NewAuthorFilter("John"),
			expected: "Author contains match: 'John'",
		},
		{
			name:     "Merge commit filter - include",
			filter:   filters.NewMergeCommitFilter(true),
			expected: "Include merge commits",
		},
		{
			name:     "Merge commit filter - exclude",
			filter:   filters.NewMergeCommitFilter(false),
			expected: "Exclude merge commits",
		},
		{
			name:     "File path filter",
			filter:   filters.NewFilePathFilter([]string{"*.go"}, filters.FileGlobMatch, false),
			expected: "File paths: *.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			description := tt.filter.Description()
			if description != tt.expected {
				t.Errorf("Expected description '%s', got '%s'", tt.expected, description)
			}
		})
	}
}

func TestComplexFilterChain(t *testing.T) {
	commits := createTestCommits()
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// Create a complex filter chain
	chain := filters.NewFilterChain()

	// Add date range filter (first 2 days)
	until := baseTime.AddDate(0, 0, 1)
	chain.Add(filters.NewDateRangeFilter(&baseTime, &until))

	// Add author filter for John
	chain.Add(filters.NewAuthorFilter("John"))

	// Add file filter for .go files
	chain.Add(filters.NewFilePathFilter([]string{".go"}, filters.FileContainsMatch, false))

	result := chain.Apply(commits)

	// Should return 1 commit: John's initial commit with main.go
	expected := 1
	if len(result) != expected {
		t.Errorf("Complex filter chain should return %d commits, got %d", expected, len(result))
	}

	// Verify it's the correct commit
	if len(result) > 0 && result[0].Hash != "abc123" {
		t.Errorf("Expected commit abc123, got %s", result[0].Hash)
	}
}

// Benchmark tests
func BenchmarkDateRangeFilter(b *testing.B) {
	commits := createTestCommits()
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	filter := filters.NewDateRangeFilter(&baseTime, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Apply(commits)
	}
}

func BenchmarkAuthorFilter(b *testing.B) {
	commits := createTestCommits()
	filter := filters.NewAuthorFilter("John")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Apply(commits)
	}
}

func BenchmarkComplexFilterChain(b *testing.B) {
	commits := createTestCommits()
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	chain := filters.NewFilterChain()
	chain.Add(filters.NewDateRangeFilter(&baseTime, nil))
	chain.Add(filters.NewAuthorFilter("John"))
	chain.Add(filters.NewFilePathFilter([]string{".go"}, filters.FileContainsMatch, false))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain.Apply(commits)
	}
}
