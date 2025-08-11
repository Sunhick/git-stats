// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Statistics analyzer tests

package analyzers

import (
	"fmt"
	"git-stats/analyzers"
	"git-stats/models"
	"testing"
	"time"
)

func TestNewStatisticsAnalyzer(t *testing.T) {
	analyzer := analyzers.NewStatisticsAnalyzer()
	if analyzer == nil {
		t.Fatal("NewStatisticsAnalyzer should not return nil")
	}
}

func TestAnalyzeStatistics_EmptyCommits(t *testing.T) {
	analyzer := analyzers.NewStatisticsAnalyzer()

	config := models.AnalysisConfig{}
	result, err := analyzer.AnalyzeStatistics([]models.Commit{}, config)
	if err != nil {
		t.Fatalf("AnalyzeStatistics failed: %v", err)
	}

	if result.TotalCommits != 0 {
		t.Errorf("Expected 0 total commits, got %d", result.TotalCommits)
	}

	if result.TotalInsertions != 0 {
		t.Errorf("Expected 0 total insertions, got %d", result.TotalInsertions)
	}

	if result.TotalDeletions != 0 {
		t.Errorf("Expected 0 total deletions, got %d", result.TotalDeletions)
	}

	if len(result.TopFiles) != 0 {
		t.Errorf("Expected 0 top files, got %d", len(result.TopFiles))
	}

	if len(result.TopFileTypes) != 0 {
		t.Errorf("Expected 0 top file types, got %d", len(result.TopFileTypes))
	}
}

func TestAnalyzeStatistics_WithCommits(t *testing.T) {
	analyzer := analyzers.NewStatisticsAnalyzer()

	commits := []models.Commit{
		{
			Hash:       "abc123",
			Message:    "First commit",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			Stats: models.CommitStats{
				FilesChanged: 2,
				Insertions:   15,
				Deletions:    5,
				Files: []models.FileChange{
					{Path: "main.go", Status: "M", Insertions: 10, Deletions: 2},
					{Path: "README.md", Status: "A", Insertions: 5, Deletions: 3},
				},
			},
		},
		{
			Hash:       "def456",
			Message:    "Second commit",
			Author:     models.Author{Name: "Jane Smith", Email: "jane@example.com"},
			AuthorDate: time.Date(2024, 1, 16, 14, 45, 0, 0, time.UTC),
			Stats: models.CommitStats{
				FilesChanged: 1,
				Insertions:   8,
				Deletions:    3,
				Files: []models.FileChange{
					{Path: "main.go", Status: "M", Insertions: 8, Deletions: 3},
				},
			},
		},
	}

	config := models.AnalysisConfig{}
	result, err := analyzer.AnalyzeStatistics(commits, config)
	if err != nil {
		t.Fatalf("AnalyzeStatistics failed: %v", err)
	}

	if result.TotalCommits != 2 {
		t.Errorf("Expected 2 total commits, got %d", result.TotalCommits)
	}

	if result.TotalInsertions != 23 {
		t.Errorf("Expected 23 total insertions, got %d", result.TotalInsertions)
	}

	if result.TotalDeletions != 8 {
		t.Errorf("Expected 8 total deletions, got %d", result.TotalDeletions)
	}

	if result.FilesChanged != 3 {
		t.Errorf("Expected 3 files changed, got %d", result.FilesChanged)
	}

	if result.ActiveDays != 2 {
		t.Errorf("Expected 2 active days, got %d", result.ActiveDays)
	}

	expectedAvg := 2.0 / 2.0
	if result.AvgCommitsPerDay != expectedAvg {
		t.Errorf("Expected average %.2f commits per day, got %.2f", expectedAvg, result.AvgCommitsPerDay)
	}
}

func TestAnalyzeCommitPatterns(t *testing.T) {
	analyzer := analyzers.NewStatisticsAnalyzer()

	commits := []models.Commit{
		{
			AuthorDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), // Monday, 10 AM
		},
		{
			AuthorDate: time.Date(2024, 1, 15, 14, 45, 0, 0, time.UTC), // Monday, 2 PM
		},
		{
			AuthorDate: time.Date(2024, 1, 16, 9, 15, 0, 0, time.UTC), // Tuesday, 9 AM
		},
	}

	hourCounts, weekdayCounts := analyzer.AnalyzeCommitPatterns(commits)

	// Check hour counts
	if hourCounts[10] != 1 {
		t.Errorf("Expected 1 commit at hour 10, got %d", hourCounts[10])
	}
	if hourCounts[14] != 1 {
		t.Errorf("Expected 1 commit at hour 14, got %d", hourCounts[14])
	}
	if hourCounts[9] != 1 {
		t.Errorf("Expected 1 commit at hour 9, got %d", hourCounts[9])
	}

	// Check weekday counts (Monday = 1, Tuesday = 2)
	if weekdayCounts[time.Monday] != 2 {
		t.Errorf("Expected 2 commits on Monday, got %d", weekdayCounts[time.Monday])
	}
	if weekdayCounts[time.Tuesday] != 1 {
		t.Errorf("Expected 1 commit on Tuesday, got %d", weekdayCounts[time.Tuesday])
	}
}

func TestAnalyzeFileStatistics(t *testing.T) {
	analyzer := analyzers.NewStatisticsAnalyzer()

	commits := []models.Commit{
		{
			AuthorDate: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Stats: models.CommitStats{
				Files: []models.FileChange{
					{Path: "main.go", Status: "M", Insertions: 10, Deletions: 2},
					{Path: "utils.py", Status: "A", Insertions: 5, Deletions: 0},
				},
			},
		},
		{
			AuthorDate: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
			Stats: models.CommitStats{
				Files: []models.FileChange{
					{Path: "main.go", Status: "M", Insertions: 8, Deletions: 3},
					{Path: "test.js", Status: "A", Insertions: 15, Deletions: 0},
				},
			},
		},
	}

	topFiles, topFileTypes := analyzer.AnalyzeFileStatistics(commits)

	// Check top files
	if len(topFiles) < 1 {
		t.Fatal("Expected at least 1 top file")
	}

	// main.go should be the top file (2 commits)
	if topFiles[0].Path != "main.go" {
		t.Errorf("Expected main.go to be top file, got %s", topFiles[0].Path)
	}
	if topFiles[0].Commits != 2 {
		t.Errorf("Expected main.go to have 2 commits, got %d", topFiles[0].Commits)
	}
	if topFiles[0].Insertions != 18 {
		t.Errorf("Expected main.go to have 18 insertions, got %d", topFiles[0].Insertions)
	}
	if topFiles[0].Deletions != 5 {
		t.Errorf("Expected main.go to have 5 deletions, got %d", topFiles[0].Deletions)
	}

	// Check top file types
	if len(topFileTypes) < 1 {
		t.Fatal("Expected at least 1 top file type")
	}

	// Find go file type
	var goType *models.FileTypeStats
	for _, ft := range topFileTypes {
		if ft.Extension == "go" {
			goType = &ft
			break
		}
	}

	if goType == nil {
		t.Fatal("Expected to find 'go' file type in results")
	}

	if goType.Files != 1 {
		t.Errorf("Expected go type to have 1 file, got %d", goType.Files)
	}
	if goType.Commits != 2 {
		t.Errorf("Expected go type to have 2 commits, got %d", goType.Commits)
	}
}

func TestGetCommitFrequencyAnalysis(t *testing.T) {
	analyzer := analyzers.NewStatisticsAnalyzer()

	commits := []models.Commit{
		{AuthorDate: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)},
		{AuthorDate: time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)},
		{AuthorDate: time.Date(2024, 1, 22, 9, 0, 0, 0, time.UTC)},
		{AuthorDate: time.Date(2024, 2, 1, 11, 0, 0, 0, time.UTC)},
	}

	analysis := analyzer.GetCommitFrequencyAnalysis(commits)

	// Check daily frequency
	if analysis.Daily["2024-01-15"] != 2 {
		t.Errorf("Expected 2 commits on 2024-01-15, got %d", analysis.Daily["2024-01-15"])
	}
	if analysis.Daily["2024-01-22"] != 1 {
		t.Errorf("Expected 1 commit on 2024-01-22, got %d", analysis.Daily["2024-01-22"])
	}
	if analysis.Daily["2024-02-01"] != 1 {
		t.Errorf("Expected 1 commit on 2024-02-01, got %d", analysis.Daily["2024-02-01"])
	}

	// Check monthly frequency
	if analysis.Monthly["2024-01"] != 3 {
		t.Errorf("Expected 3 commits in 2024-01, got %d", analysis.Monthly["2024-01"])
	}
	if analysis.Monthly["2024-02"] != 1 {
		t.Errorf("Expected 1 commit in 2024-02, got %d", analysis.Monthly["2024-02"])
	}

	// Check weekly frequency
	if len(analysis.Weekly) == 0 {
		t.Error("Expected weekly frequency data")
	}
}

func TestGetTimeBasedPatterns(t *testing.T) {
	analyzer := analyzers.NewStatisticsAnalyzer()

	commits := []models.Commit{
		{AuthorDate: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)}, // Monday
		{AuthorDate: time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)}, // Monday
		{AuthorDate: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC)}, // Tuesday
	}

	patterns := analyzer.GetTimeBasedPatterns(commits)

	// Check hourly distribution
	if patterns.HourlyDistribution[10] != 2 {
		t.Errorf("Expected 2 commits at hour 10, got %d", patterns.HourlyDistribution[10])
	}
	if patterns.HourlyDistribution[14] != 1 {
		t.Errorf("Expected 1 commit at hour 14, got %d", patterns.HourlyDistribution[14])
	}

	// Check weekday distribution
	if patterns.WeekdayDistribution[time.Monday] != 2 {
		t.Errorf("Expected 2 commits on Monday, got %d", patterns.WeekdayDistribution[time.Monday])
	}
	if patterns.WeekdayDistribution[time.Tuesday] != 1 {
		t.Errorf("Expected 1 commit on Tuesday, got %d", patterns.WeekdayDistribution[time.Tuesday])
	}

	// Check that averages are calculated
	if len(patterns.HourlyAverage) == 0 {
		t.Error("Expected hourly averages to be calculated")
	}
	if len(patterns.WeekdayAverage) == 0 {
		t.Error("Expected weekday averages to be calculated")
	}
}

func TestAnalyzeStatistics_WithFilters(t *testing.T) {
	analyzer := analyzers.NewStatisticsAnalyzer()

	commits := []models.Commit{
		{
			Hash:       "abc123",
			Message:    "John's commit",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Stats:      models.CommitStats{FilesChanged: 1, Insertions: 10, Deletions: 5},
		},
		{
			Hash:       "def456",
			Message:    "Jane's commit",
			Author:     models.Author{Name: "Jane Smith", Email: "jane@example.com"},
			AuthorDate: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
			Stats:      models.CommitStats{FilesChanged: 2, Insertions: 8, Deletions: 3},
		},
		{
			Hash:       "ghi789",
			Message:    "Old commit",
			Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate: time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
			Stats:      models.CommitStats{FilesChanged: 1, Insertions: 5, Deletions: 2},
		},
	}

	// Test author filter
	config := models.AnalysisConfig{
		AuthorFilter: "John",
	}

	result, err := analyzer.AnalyzeStatistics(commits, config)
	if err != nil {
		t.Fatalf("AnalyzeStatistics failed: %v", err)
	}

	if result.TotalCommits != 2 {
		t.Errorf("Expected 2 commits after author filter, got %d", result.TotalCommits)
	}

	// Test time range filter
	config = models.AnalysisConfig{
		TimeRange: models.TimeRange{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		},
	}

	result, err = analyzer.AnalyzeStatistics(commits, config)
	if err != nil {
		t.Fatalf("AnalyzeStatistics failed: %v", err)
	}

	if result.TotalCommits != 2 {
		t.Errorf("Expected 2 commits after time range filter, got %d", result.TotalCommits)
	}
}

func TestAnalyzeStatistics_ExcludeMerges(t *testing.T) {
	analyzer := analyzers.NewStatisticsAnalyzer()

	commits := []models.Commit{
		{
			Hash:         "abc123",
			Message:      "Regular commit",
			Author:       models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			ParentHashes: []string{"parent1"},
			Stats:        models.CommitStats{FilesChanged: 1, Insertions: 10, Deletions: 5},
		},
		{
			Hash:         "def456",
			Message:      "Merge commit",
			Author:       models.Author{Name: "John Doe", Email: "john@example.com"},
			AuthorDate:   time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
			ParentHashes: []string{"parent1", "parent2"}, // Merge commit
			Stats:        models.CommitStats{FilesChanged: 2, Insertions: 8, Deletions: 3},
		},
	}

	config := models.AnalysisConfig{
		IncludeMerges: false,
	}

	result, err := analyzer.AnalyzeStatistics(commits, config)
	if err != nil {
		t.Fatalf("AnalyzeStatistics failed: %v", err)
	}

	if result.TotalCommits != 1 {
		t.Errorf("Expected 1 commit after excluding merges, got %d", result.TotalCommits)
	}

	if result.TotalInsertions != 10 {
		t.Errorf("Expected 10 insertions after excluding merges, got %d", result.TotalInsertions)
	}
}

// Benchmark tests
func BenchmarkAnalyzeStatistics(b *testing.B) {
	analyzer := analyzers.NewStatisticsAnalyzer()

	// Create a large set of test commits
	commits := make([]models.Commit, 1000)
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 1000; i++ {
		commits[i] = models.Commit{
			Hash:       fmt.Sprintf("commit%d", i),
			Message:    fmt.Sprintf("Commit %d", i),
			Author:     models.Author{Name: "Test User", Email: "test@example.com"},
			AuthorDate: baseDate.AddDate(0, 0, i%365),
			Stats: models.CommitStats{
				FilesChanged: 1 + i%5,
				Insertions:   10 + i%50,
				Deletions:    5 + i%20,
				Files: []models.FileChange{
					{
						Path:       fmt.Sprintf("file%d.go", i%100),
						Status:     "M",
						Insertions: 10 + i%50,
						Deletions:  5 + i%20,
					},
				},
			},
		}
	}

	config := models.AnalysisConfig{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := analyzer.AnalyzeStatistics(commits, config)
		if err != nil {
			b.Fatalf("AnalyzeStatistics failed: %v", err)
		}
	}
}

func BenchmarkAnalyzeFileStatistics(b *testing.B) {
	analyzer := analyzers.NewStatisticsAnalyzer()

	// Create commits with many file changes
	commits := make([]models.Commit, 100)
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 100; i++ {
		files := make([]models.FileChange, 10)
		for j := 0; j < 10; j++ {
			files[j] = models.FileChange{
				Path:       fmt.Sprintf("dir%d/file%d.go", j%5, (i*10+j)%50),
				Status:     "M",
				Insertions: 10 + j%20,
				Deletions:  5 + j%10,
			}
		}

		commits[i] = models.Commit{
			Hash:       fmt.Sprintf("commit%d", i),
			AuthorDate: baseDate.AddDate(0, 0, i),
			Stats: models.CommitStats{
				Files: files,
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.AnalyzeFileStatistics(commits)
	}
}
