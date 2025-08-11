// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for statistics models

package models_test

import (
	"testing"
	"time"
	"git-stats/models"
)

func TestAnalysisResult(t *testing.T) {
	result := &models.AnalysisResult{
		Repository: &models.RepositoryInfo{
			Path:         "/test/repo",
			Name:         "test-repo",
			TotalCommits: 100,
			FirstCommit:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			LastCommit:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Branches:     []string{"main", "develop", "feature/test"},
		},
		Summary: &models.StatsSummary{
			TotalCommits:     100,
			TotalInsertions:  2000,
			TotalDeletions:   500,
			FilesChanged:     50,
			ActiveDays:       200,
			AvgCommitsPerDay: 0.5,
			CommitsByHour:    map[int]int{9: 20, 10: 15, 14: 25},
			CommitsByWeekday: map[time.Weekday]int{time.Monday: 20, time.Tuesday: 15},
		},
		Contributors: []models.ContributorStats{
			{
				Name:            "Alice",
				Email:           "alice@example.com",
				TotalCommits:    60,
				TotalInsertions: 1200,
				TotalDeletions:  300,
				FirstCommit:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				LastCommit:      time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
				ActiveDays:      120,
			},
			{
				Name:            "Bob",
				Email:           "bob@example.com",
				TotalCommits:    40,
				TotalInsertions: 800,
				TotalDeletions:  200,
				FirstCommit:     time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				LastCommit:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				ActiveDays:      80,
			},
		},
		TimeRange: models.TimeRange{
			Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	// Test repository info
	if result.Repository.Name != "test-repo" {
		t.Errorf("Expected repo name 'test-repo', got '%s'", result.Repository.Name)
	}

	if result.Repository.TotalCommits != 100 {
		t.Errorf("Expected 100 total commits, got %d", result.Repository.TotalCommits)
	}

	if len(result.Repository.Branches) != 3 {
		t.Errorf("Expected 3 branches, got %d", len(result.Repository.Branches))
	}

	// Test summary stats
	if result.Summary.TotalCommits != 100 {
		t.Errorf("Expected 100 total commits in summary, got %d", result.Summary.TotalCommits)
	}

	if result.Summary.AvgCommitsPerDay != 0.5 {
		t.Errorf("Expected 0.5 avg commits per day, got %f", result.Summary.AvgCommitsPerDay)
	}

	// Test contributors
	if len(result.Contributors) != 2 {
		t.Errorf("Expected 2 contributors, got %d", len(result.Contributors))
	}

	alice := result.Contributors[0]
	if alice.Name != "Alice" {
		t.Errorf("Expected first contributor name 'Alice', got '%s'", alice.Name)
	}

	if alice.TotalCommits != 60 {
		t.Errorf("Expected Alice to have 60 commits, got %d", alice.TotalCommits)
	}

	// Test time range
	expectedStart := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !result.TimeRange.Start.Equal(expectedStart) {
		t.Errorf("Expected start time %v, got %v", expectedStart, result.TimeRange.Start)
	}
}

func TestContributionGraph(t *testing.T) {
	graph := &models.ContributionGraph{
		StartDate:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DailyCommits: map[string]int{
			"2024-01-01": 3,
			"2024-01-02": 1,
			"2024-01-03": 0,
			"2024-01-04": 5,
			"2024-01-05": 2,
		},
		MaxCommits:   5,
		TotalCommits: 11,
	}

	if graph.MaxCommits != 5 {
		t.Errorf("Expected max commits 5, got %d", graph.MaxCommits)
	}

	if graph.TotalCommits != 11 {
		t.Errorf("Expected total commits 11, got %d", graph.TotalCommits)
	}

	if graph.DailyCommits["2024-01-04"] != 5 {
		t.Errorf("Expected 5 commits on 2024-01-04, got %d", graph.DailyCommits["2024-01-04"])
	}

	if graph.DailyCommits["2024-01-03"] != 0 {
		t.Errorf("Expected 0 commits on 2024-01-03, got %d", graph.DailyCommits["2024-01-03"])
	}
}

func TestHealthMetrics(t *testing.T) {
	health := &models.HealthMetrics{
		RepositoryAge:      365 * 24 * time.Hour, // 1 year
		CommitFrequency:    2.5,                  // commits per day
		ContributorCount:   5,
		ActiveContributors: 3,
		BranchCount:        8,
		ActivityTrend:      "increasing",
		MonthlyGrowth: []models.MonthlyStats{
			{
				Month:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Commits: 75,
				Authors: 3,
			},
			{
				Month:   time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
				Commits: 60,
				Authors: 2,
			},
		},
	}

	expectedAge := 365 * 24 * time.Hour
	if health.RepositoryAge != expectedAge {
		t.Errorf("Expected repository age %v, got %v", expectedAge, health.RepositoryAge)
	}

	if health.CommitFrequency != 2.5 {
		t.Errorf("Expected commit frequency 2.5, got %f", health.CommitFrequency)
	}

	if health.ActivityTrend != "increasing" {
		t.Errorf("Expected activity trend 'increasing', got '%s'", health.ActivityTrend)
	}

	if len(health.MonthlyGrowth) != 2 {
		t.Errorf("Expected 2 monthly growth entries, got %d", len(health.MonthlyGrowth))
	}

	if health.MonthlyGrowth[0].Commits != 75 {
		t.Errorf("Expected 75 commits in first month, got %d", health.MonthlyGrowth[0].Commits)
	}
}

func TestFileStats(t *testing.T) {
	fileStats := []models.FileStats{
		{
			Path:         "main.go",
			Commits:      25,
			Insertions:   500,
			Deletions:    100,
			LastModified: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			Path:         "README.md",
			Commits:      10,
			Insertions:   200,
			Deletions:    50,
			LastModified: time.Date(2024, 1, 10, 14, 20, 0, 0, time.UTC),
		},
	}

	if len(fileStats) != 2 {
		t.Errorf("Expected 2 file stats, got %d", len(fileStats))
	}

	mainGo := fileStats[0]
	if mainGo.Path != "main.go" {
		t.Errorf("Expected path 'main.go', got '%s'", mainGo.Path)
	}

	if mainGo.Commits != 25 {
		t.Errorf("Expected 25 commits for main.go, got %d", mainGo.Commits)
	}
}

func TestTimeRange(t *testing.T) {
	timeRange := models.TimeRange{
		Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
	}

	if timeRange.Start.After(timeRange.End) {
		t.Error("Start time should be before end time")
	}

	duration := timeRange.End.Sub(timeRange.Start)
	expectedDuration := 30*24*time.Hour + 23*time.Hour + 59*time.Minute + 59*time.Second

	if duration != expectedDuration {
		t.Errorf("Expected duration %v, got %v", expectedDuration, duration)
	}
}
