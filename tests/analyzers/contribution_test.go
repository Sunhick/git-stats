// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Contribution analyzer tests

package analyzers

import (
	"fmt"
	"git-stats/analyzers"
	"git-stats/models"
	"testing"
	"time"
)

func TestNewContributionAnalyzer(t *testing.T) {
	analyzer := analyzers.NewContributionAnalyzer()
	if analyzer == nil {
		t.Fatal("NewContributionAnalyzer should not return nil")
	}
}

func TestAnalyzeContributions_EmptyCommits(t *testing.T) {
	analyzer := analyzers.NewContributionAnalyzer()

	config := models.AnalysisConfig{
		TimeRange: models.TimeRange{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	result, err := analyzer.AnalyzeContributions([]models.Commit{}, config)
	if err != nil {
		t.Fatalf("AnalyzeContributions failed: %v", err)
	}

	if result.TotalCommits != 0 {
		t.Errorf("Expected 0 total commits, got %d", result.TotalCommits)
	}

	if result.MaxCommits != 0 {
		t.Errorf("Expected 0 max commits, got %d", result.MaxCommits)
	}

	if len(result.DailyCommits) != 0 {
		t.Errorf("Expected empty daily commits map, got %d entries", len(result.DailyCommits))
	}
}

func TestAnalyzeContributions_WithCommits(t *testing.T) {
	analyzer := analyzers.NewContributionAnalyzer()

	// Create test commits
	commits := []models.Commit{
		{
			Hash:       "abc123",
			Message:    "First commit",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Stats:      models.CommitStats{FilesChanged: 1, Insertions: 10, Deletions: 0},
		},
		{
			Hash:       "def456",
			Message:    "Second commit",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
			Stats:      models.CommitStats{FilesChanged: 2, Insertions: 5, Deletions: 2},
		},
		{
			Hash:       "ghi789",
			Message:    "Third commit",
			Author:     models.Author{Name: "Jane Smith", Email: "jane@example.com"},
			AuthorDate: time.Date(2024, 1, 16, 9, 0, 0, 0, time.UTC),
			Stats:      models.CommitStats{FilesChanged: 1, Insertions: 3, Deletions: 1},
		},
	}

	config := models.AnalysisConfig{
		TimeRange: models.TimeRange{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	result, err := analyzer.AnalyzeContributions(commits, config)
	if err != nil {
		t.Fatalf("AnalyzeContributions failed: %v", err)
	}

	if result.TotalCommits != 3 {
		t.Errorf("Expected 3 total commits, got %d", result.TotalCommits)
	}

	if result.MaxCommits != 2 {
		t.Errorf("Expected 2 max commits per day, got %d", result.MaxCommits)
	}

	// Check specific dates
	jan15 := "2024-01-15"
	jan16 := "2024-01-16"

	if result.DailyCommits[jan15] != 2 {
		t.Errorf("Expected 2 commits on %s, got %d", jan15, result.DailyCommits[jan15])
	}

	if result.DailyCommits[jan16] != 1 {
		t.Errorf("Expected 1 commit on %s, got %d", jan16, result.DailyCommits[jan16])
	}
}

func TestAnalyzeContributions_AuthorFilter(t *testing.T) {
	analyzer := analyzers.NewContributionAnalyzer()

	commits := []models.Commit{
		{
			Hash:       "abc123",
			Message:    "John's commit",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			Hash:       "def456",
			Message:    "Jane's commit",
			Author:     models.Author{Name: "Jane Smith", Email: "jane@example.com"},
			AuthorDate: time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
		},
	}

	config := models.AnalysisConfig{
		TimeRange: models.TimeRange{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		},
		AuthorFilter: "John",
	}

	result, err := analyzer.AnalyzeContributions(commits, config)
	if err != nil {
		t.Fatalf("AnalyzeContributions failed: %v", err)
	}

	if result.TotalCommits != 1 {
		t.Errorf("Expected 1 total commit after filtering, got %d", result.TotalCommits)
	}

	jan15 := "2024-01-15"
	if result.DailyCommits[jan15] != 1 {
		t.Errorf("Expected 1 commit on %s after filtering, got %d", jan15, result.DailyCommits[jan15])
	}
}

func TestAnalyzeContributions_ExcludeMerges(t *testing.T) {
	analyzer := analyzers.NewContributionAnalyzer()

	commits := []models.Commit{
		{
			Hash:         "abc123",
			Message:      "Regular commit",
			Author:       models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			ParentHashes: []string{"parent1"},
		},
		{
			Hash:         "def456",
			Message:      "Merge commit",
			Author:       models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate:   time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
			ParentHashes: []string{"parent1", "parent2"}, // Merge commit has multiple parents
		},
	}

	config := models.AnalysisConfig{
		TimeRange: models.TimeRange{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		},
		IncludeMerges: false,
	}

	result, err := analyzer.AnalyzeContributions(commits, config)
	if err != nil {
		t.Fatalf("AnalyzeContributions failed: %v", err)
	}

	if result.TotalCommits != 1 {
		t.Errorf("Expected 1 total commit after excluding merges, got %d", result.TotalCommits)
	}
}

func TestCalculateActivityLevels(t *testing.T) {
	analyzer := analyzers.NewContributionAnalyzer()

	dailyCommits := map[string]int{
		"2024-01-01": 0,  // Level 0
		"2024-01-02": 1,  // Level 1
		"2024-01-03": 3,  // Level 1
		"2024-01-04": 4,  // Level 2
		"2024-01-05": 9,  // Level 2
		"2024-01-06": 10, // Level 3
		"2024-01-07": 19, // Level 3
		"2024-01-08": 20, // Level 4
		"2024-01-09": 50, // Level 4
	}

	levels := analyzer.CalculateActivityLevels(dailyCommits)

	expectedLevels := map[string]int{
		"2024-01-01": 0,
		"2024-01-02": 1,
		"2024-01-03": 1,
		"2024-01-04": 2,
		"2024-01-05": 2,
		"2024-01-06": 3,
		"2024-01-07": 3,
		"2024-01-08": 4,
		"2024-01-09": 4,
	}

	for date, expectedLevel := range expectedLevels {
		if levels[date] != expectedLevel {
			t.Errorf("Expected activity level %d for %s, got %d", expectedLevel, date, levels[date])
		}
	}
}

func TestCalculateStreaks(t *testing.T) {
	analyzer := analyzers.NewContributionAnalyzer()

	tests := []struct {
		name           string
		dailyCommits   map[string]int
		expectedCurrent int
		expectedLongest int
	}{
		{
			name:           "Empty commits",
			dailyCommits:   map[string]int{},
			expectedCurrent: 0,
			expectedLongest: 0,
		},
		{
			name: "No commits",
			dailyCommits: map[string]int{
				"2024-01-01": 0,
				"2024-01-02": 0,
				"2024-01-03": 0,
			},
			expectedCurrent: 0,
			expectedLongest: 0,
		},
		{
			name: "Simple streak",
			dailyCommits: map[string]int{
				"2024-01-01": 1,
				"2024-01-02": 2,
				"2024-01-03": 1,
			},
			expectedCurrent: 3,
			expectedLongest: 3,
		},
		{
			name: "Broken streak",
			dailyCommits: map[string]int{
				"2024-01-01": 1,
				"2024-01-02": 1,
				"2024-01-03": 0,
				"2024-01-04": 1,
				"2024-01-05": 1,
			},
			expectedCurrent: 2,
			expectedLongest: 2,
		},
		{
			name: "Long streak with gap",
			dailyCommits: map[string]int{
				"2024-01-01": 1,
				"2024-01-02": 1,
				"2024-01-03": 1,
				"2024-01-04": 1,
				"2024-01-05": 0,
				"2024-01-06": 1,
				"2024-01-07": 1,
			},
			expectedCurrent: 2,
			expectedLongest: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current, longest := analyzer.CalculateStreaks(tt.dailyCommits)

			if current != tt.expectedCurrent {
				t.Errorf("Expected current streak %d, got %d", tt.expectedCurrent, current)
			}

			if longest != tt.expectedLongest {
				t.Errorf("Expected longest streak %d, got %d", tt.expectedLongest, longest)
			}
		})
	}
}

func TestGetContributionSummary(t *testing.T) {
	analyzer := analyzers.NewContributionAnalyzer()

	graph := &models.ContributionGraph{
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		DailyCommits: map[string]int{
			"2024-01-01": 0,
			"2024-01-02": 2,
			"2024-01-03": 1,
			"2024-01-04": 0,
			"2024-01-05": 3,
		},
		MaxCommits:   3,
		TotalCommits: 6,
	}

	summary := analyzer.GetContributionSummary(graph)

	if summary.TotalCommits != 6 {
		t.Errorf("Expected 6 total commits, got %d", summary.TotalCommits)
	}

	if summary.MaxCommitsPerDay != 3 {
		t.Errorf("Expected 3 max commits per day, got %d", summary.MaxCommitsPerDay)
	}

	if summary.ActiveDays != 3 {
		t.Errorf("Expected 3 active days, got %d", summary.ActiveDays)
	}

	if summary.TotalDays != 5 {
		t.Errorf("Expected 5 total days, got %d", summary.TotalDays)
	}

	expectedAvg := 6.0 / 5.0
	if summary.AvgCommitsPerDay != expectedAvg {
		t.Errorf("Expected average %.2f commits per day, got %.2f", expectedAvg, summary.AvgCommitsPerDay)
	}

	if len(summary.ActivityLevels) != 5 {
		t.Errorf("Expected 5 activity levels, got %d", len(summary.ActivityLevels))
	}
}

func TestGetContributionSummary_NilGraph(t *testing.T) {
	analyzer := analyzers.NewContributionAnalyzer()

	summary := analyzer.GetContributionSummary(nil)

	if summary.TotalCommits != 0 {
		t.Errorf("Expected 0 total commits for nil graph, got %d", summary.TotalCommits)
	}

	if summary.ActiveDays != 0 {
		t.Errorf("Expected 0 active days for nil graph, got %d", summary.ActiveDays)
	}
}

// Benchmark tests
func BenchmarkAnalyzeContributions(b *testing.B) {
	analyzer := analyzers.NewContributionAnalyzer()

	// Create a large set of test commits
	commits := make([]models.Commit, 1000)
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 1000; i++ {
		commits[i] = models.Commit{
			Hash:       fmt.Sprintf("commit%d", i),
			Message:    fmt.Sprintf("Commit %d", i),
			Author:     models.Author{Name: "Test User", Email: "test@example.com"},
			AuthorDate: baseDate.AddDate(0, 0, i%365), // Spread over a year
			Stats:      models.CommitStats{FilesChanged: 1, Insertions: 10, Deletions: 5},
		}
	}

	config := models.AnalysisConfig{
		TimeRange: models.TimeRange{
			Start: baseDate,
			End:   baseDate.AddDate(1, 0, 0),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := analyzer.AnalyzeContributions(commits, config)
		if err != nil {
			b.Fatalf("AnalyzeContributions failed: %v", err)
		}
	}
}

func BenchmarkCalculateActivityLevels(b *testing.B) {
	analyzer := analyzers.NewContributionAnalyzer()

	// Create a large daily commits map
	dailyCommits := make(map[string]int)
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 365; i++ {
		date := baseDate.AddDate(0, 0, i).Format("2006-01-02")
		dailyCommits[date] = i % 25 // Vary commit counts
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.CalculateActivityLevels(dailyCommits)
	}
}
