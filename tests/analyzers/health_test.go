// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Health analyzer tests

package analyzers

import (
	"fmt"
	"git-stats/analyzers"
	"git-stats/models"
	"strings"
	"testing"
	"time"
)

func TestNewHealthAnalyzer(t *testing.T) {
	analyzer := analyzers.NewHealthAnalyzer()
	if analyzer == nil {
		t.Fatal("NewHealthAnalyzer should not return nil")
	}
}

func TestAnalyzeHealth_EmptyCommits(t *testing.T) {
	analyzer := analyzers.NewHealthAnalyzer()

	config := models.AnalysisConfig{}
	contributors := []models.Contributor{}

	result, err := analyzer.AnalyzeHealth([]models.Commit{}, contributors, config)
	if err != nil {
		t.Fatalf("AnalyzeHealth failed: %v", err)
	}

	if result.RepositoryAge != 0 {
		t.Errorf("Expected 0 repository age, got %v", result.RepositoryAge)
	}

	if result.CommitFrequency != 0 {
		t.Errorf("Expected 0 commit frequency, got %f", result.CommitFrequency)
	}

	if result.ActivityTrend != "stable" {
		t.Errorf("Expected 'stable' activity trend, got %s", result.ActivityTrend)
	}

	if len(result.MonthlyGrowth) != 0 {
		t.Errorf("Expected 0 monthly growth entries, got %d", len(result.MonthlyGrowth))
	}
}

func TestAnalyzeHealth_WithCommits(t *testing.T) {
	analyzer := analyzers.NewHealthAnalyzer()

	// Create test commits spanning 30 days
	baseDate := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	commits := []models.Commit{
		{
			Hash:       "abc123",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: baseDate,
		},
		{
			Hash:       "def456",
			Author:     models.Author{Name: "Jane Smith", Email: "jane@example.com"},
			AuthorDate: baseDate.AddDate(0, 0, 15),
		},
		{
			Hash:       "ghi789",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: baseDate.AddDate(0, 0, 30),
		},
	}

	// Create test contributors
	contributors := []models.Contributor{
		{
			Name:        "John Doe",
			Email:       "john@example.com",
			LastCommit:  time.Now().AddDate(0, -1, 0), // Active (1 month ago)
		},
		{
			Name:        "Jane Smith",
			Email:       "jane@example.com",
			LastCommit:  time.Now().AddDate(0, -6, 0), // Inactive (6 months ago)
		},
	}

	config := models.AnalysisConfig{}
	result, err := analyzer.AnalyzeHealth(commits, contributors, config)
	if err != nil {
		t.Fatalf("AnalyzeHealth failed: %v", err)
	}

	// Repository age should be 30 days
	expectedAge := 30 * 24 * time.Hour
	if result.RepositoryAge != expectedAge {
		t.Errorf("Expected repository age %v, got %v", expectedAge, result.RepositoryAge)
	}

	// Commit frequency should be 3 commits / 30 days = 0.1 commits/day
	expectedFreq := 3.0 / 30.0
	if result.CommitFrequency != expectedFreq {
		t.Errorf("Expected commit frequency %f, got %f", expectedFreq, result.CommitFrequency)
	}

	if result.ContributorCount != 2 {
		t.Errorf("Expected 2 contributors, got %d", result.ContributorCount)
	}

	if result.ActiveContributors != 1 {
		t.Errorf("Expected 1 active contributor, got %d", result.ActiveContributors)
	}

	if len(result.MonthlyGrowth) == 0 {
		t.Error("Expected monthly growth data")
	}
}

func TestCalculateActivityTrend(t *testing.T) {
	analyzer := analyzers.NewHealthAnalyzer()

	tests := []struct {
		name     string
		commits  []models.Commit
		expected string
	}{
		{
			name:     "Empty commits",
			commits:  []models.Commit{},
			expected: "stable",
		},
		{
			name: "Single commit",
			commits: []models.Commit{
				{AuthorDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			expected: "stable",
		},
		{
			name: "Increasing trend",
			commits: []models.Commit{
				// Earlier months (low activity)
				{AuthorDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
				// Recent months (high activity)
				{AuthorDate: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 4, 2, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 4, 3, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 5, 2, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC)},
			},
			expected: "increasing",
		},
		{
			name: "Decreasing trend",
			commits: []models.Commit{
				// Earlier months (high activity)
				{AuthorDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 3, 3, 0, 0, 0, 0, time.UTC)},
				// Recent months (low activity)
				{AuthorDate: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)},
				{AuthorDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)},
			},
			expected: "decreasing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.CalculateActivityTrend(tt.commits)
			if result != tt.expected {
				t.Errorf("Expected activity trend %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCalculateMonthlyGrowth(t *testing.T) {
	analyzer := analyzers.NewHealthAnalyzer()

	commits := []models.Commit{
		{
			Author:     models.Author{Email: "john@example.com"},
			AuthorDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Author:     models.Author{Email: "john@example.com"},
			AuthorDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			Author:     models.Author{Email: "jane@example.com"},
			AuthorDate: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			Author:     models.Author{Email: "bob@example.com"},
			AuthorDate: time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	result := analyzer.CalculateMonthlyGrowth(commits)

	if len(result) != 2 {
		t.Fatalf("Expected 2 monthly growth entries, got %d", len(result))
	}

	// Check January stats
	jan := result[0]
	if jan.Month.Month() != time.January {
		t.Errorf("Expected January, got %v", jan.Month.Month())
	}
	if jan.Commits != 3 {
		t.Errorf("Expected 3 commits in January, got %d", jan.Commits)
	}
	if jan.Authors != 2 {
		t.Errorf("Expected 2 authors in January, got %d", jan.Authors)
	}

	// Check February stats
	feb := result[1]
	if feb.Month.Month() != time.February {
		t.Errorf("Expected February, got %v", feb.Month.Month())
	}
	if feb.Commits != 1 {
		t.Errorf("Expected 1 commit in February, got %d", feb.Commits)
	}
	if feb.Authors != 1 {
		t.Errorf("Expected 1 author in February, got %d", feb.Authors)
	}
}

func TestGetRepositoryHealthScore(t *testing.T) {
	analyzer := analyzers.NewHealthAnalyzer()

	tests := []struct {
		name     string
		metrics  *models.HealthMetrics
		minScore int
		maxScore int
	}{
		{
			name:     "Nil metrics",
			metrics:  nil,
			minScore: 0,
			maxScore: 0,
		},
		{
			name: "Healthy repository",
			metrics: &models.HealthMetrics{
				RepositoryAge:      365 * 24 * time.Hour, // 1 year
				CommitFrequency:    2.0,                  // 2 commits/day
				ContributorCount:   10,
				ActiveContributors: 5,
				ActivityTrend:      "increasing",
				MonthlyGrowth: []models.MonthlyStats{
					{Commits: 30}, {Commits: 32}, {Commits: 35},
				},
			},
			minScore: 80,
			maxScore: 100,
		},
		{
			name: "Unhealthy repository",
			metrics: &models.HealthMetrics{
				RepositoryAge:      7 * 24 * time.Hour, // 1 week
				CommitFrequency:    0.01,               // Very low
				ContributorCount:   1,
				ActiveContributors: 0,
				ActivityTrend:      "decreasing",
				MonthlyGrowth:      []models.MonthlyStats{},
			},
			minScore: 0,
			maxScore: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := analyzer.GetRepositoryHealthScore(tt.metrics)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Expected score between %d and %d, got %d", tt.minScore, tt.maxScore, score)
			}
		})
	}
}

func TestGetHealthInsights(t *testing.T) {
	analyzer := analyzers.NewHealthAnalyzer()

	metrics := &models.HealthMetrics{
		RepositoryAge:      30 * 24 * time.Hour, // 30 days
		CommitFrequency:    0.05,                // Low frequency
		ContributorCount:   3,
		ActiveContributors: 0, // No active contributors
		ActivityTrend:      "decreasing",
		MonthlyGrowth: []models.MonthlyStats{
			{Commits: 10, Authors: 2},
			{Commits: 5, Authors: 1},
			{Commits: 0, Authors: 0}, // No commits in recent month
		},
	}

	insights := analyzer.GetHealthInsights(metrics)

	if len(insights) == 0 {
		t.Fatal("Expected health insights")
	}

	// Check for specific insights
	foundLowActivity := false
	foundNoActiveContributors := false
	foundDecreasingTrend := false

	for _, insight := range insights {
		if strings.Contains(insight, "Low commit frequency") {
			foundLowActivity = true
		}
		if strings.Contains(insight, "No active contributors") {
			foundNoActiveContributors = true
		}
		if strings.Contains(insight, "declining") {
			foundDecreasingTrend = true
		}
	}

	if !foundLowActivity {
		t.Error("Expected insight about low commit frequency")
	}
	if !foundNoActiveContributors {
		t.Error("Expected insight about no active contributors")
	}
	if !foundDecreasingTrend {
		t.Error("Expected insight about decreasing trend")
	}
}

func TestAnalyzeHealth_WithFilters(t *testing.T) {
	analyzer := analyzers.NewHealthAnalyzer()

	commits := []models.Commit{
		{
			Hash:       "abc123",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			Hash:       "def456",
			Author:     models.Author{Name: "Jane Smith", Email: "jane@example.com"},
			AuthorDate: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
		},
		{
			Hash:       "ghi789",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC), // Outside range
		},
	}

	contributors := []models.Contributor{}

	// Test time range filter
	config := models.AnalysisConfig{
		TimeRange: models.TimeRange{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		},
	}

	result, err := analyzer.AnalyzeHealth(commits, contributors, config)
	if err != nil {
		t.Fatalf("AnalyzeHealth failed: %v", err)
	}

	// Should only analyze 2 commits (the ones in January 2024)
	expectedAge := 1 * 24 * time.Hour // 1 day between Jan 15 and Jan 16
	if result.RepositoryAge != expectedAge {
		t.Errorf("Expected repository age %v after filtering, got %v", expectedAge, result.RepositoryAge)
	}
}

func TestAnalyzeHealth_ExcludeMerges(t *testing.T) {
	analyzer := analyzers.NewHealthAnalyzer()

	commits := []models.Commit{
		{
			Hash:         "abc123",
			Author:       models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			ParentHashes: []string{"parent1"},
		},
		{
			Hash:         "def456",
			Author:       models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate:   time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
			ParentHashes: []string{"parent1", "parent2"}, // Merge commit
		},
	}

	contributors := []models.Contributor{}

	config := models.AnalysisConfig{
		IncludeMerges: false,
	}

	result, err := analyzer.AnalyzeHealth(commits, contributors, config)
	if err != nil {
		t.Fatalf("AnalyzeHealth failed: %v", err)
	}

	// Should only analyze 1 commit (excluding the merge)
	if result.RepositoryAge != 0 {
		t.Errorf("Expected 0 repository age with single commit, got %v", result.RepositoryAge)
	}
}

// Benchmark tests
func BenchmarkAnalyzeHealth(b *testing.B) {
	analyzer := analyzers.NewHealthAnalyzer()

	// Create a large set of test commits
	commits := make([]models.Commit, 1000)
	contributors := make([]models.Contributor, 50)
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 1000; i++ {
		commits[i] = models.Commit{
			Hash:       fmt.Sprintf("commit%d", i),
			Author:     models.Author{Name: fmt.Sprintf("User%d", i%50), Email: fmt.Sprintf("user%d@example.com", i%50)},
			AuthorDate: baseDate.AddDate(0, 0, i%365),
		}
	}

	for i := 0; i < 50; i++ {
		contributors[i] = models.Contributor{
			Name:       fmt.Sprintf("User%d", i),
			Email:      fmt.Sprintf("user%d@example.com", i),
			LastCommit: time.Now().AddDate(0, -i%12, 0), // Vary activity
		}
	}

	config := models.AnalysisConfig{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := analyzer.AnalyzeHealth(commits, contributors, config)
		if err != nil {
			b.Fatalf("AnalyzeHealth failed: %v", err)
		}
	}
}

func BenchmarkCalculateActivityTrend(b *testing.B) {
	analyzer := analyzers.NewHealthAnalyzer()

	// Create commits spanning multiple months
	commits := make([]models.Commit, 500)
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 500; i++ {
		commits[i] = models.Commit{
			AuthorDate: baseDate.AddDate(0, 0, i%180), // 6 months
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.CalculateActivityTrend(commits)
	}
}
