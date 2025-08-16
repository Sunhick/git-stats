// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - CSV formatter unit tests

package formatters

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"git-stats/formatters"
	"git-stats/git"
	"git-stats/models"
)

func TestNewCSVFormatter(t *testing.T) {
	formatter := formatters.NewCSVFormatter()
	if formatter == nil {
		t.Fatal("NewCSVFormatter should return a non-nil formatter")
	}
}

func TestCSVFormatter_Format(t *testing.T) {
	formatter := formatters.NewCSVFormatter()
	testData := createTestAnalysisResult()

	tests := []struct {
		name   string
		config models.FormatConfig
	}{
		{
			name: "CSV format with metadata",
			config: models.FormatConfig{
				Format:   "csv",
				Metadata: true,
			},
		},
		{
			name: "CSV format without metadata",
			config: models.FormatConfig{
				Format:   "csv",
				Metadata: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatter.Format(testData, tt.config)
			if err != nil {
				t.Fatalf("Format() error = %v", err)
			}

			if len(result) == 0 {
				t.Fatal("Format() returned empty result")
			}

			resultStr := string(result)

			// Verify metadata is included when requested
			if tt.config.Metadata {
				if !strings.Contains(resultStr, "# Metadata") {
					t.Error("Expected metadata in CSV output when Metadata=true")
				}
				if !strings.Contains(resultStr, "# Generated at:") {
					t.Error("Expected generation timestamp in metadata")
				}
			}

			// Verify required sections exist
			expectedSections := []string{"# Contributors", "# Summary Statistics"}
			for _, section := range expectedSections {
				if !strings.Contains(resultStr, section) {
					t.Errorf("Expected section '%s' in CSV output", section)
				}
			}
		})
	}
}

func TestCSVFormatter_FormatContributorsCSV(t *testing.T) {
	formatter := formatters.NewCSVFormatter()
	contributors := createTestContributors()

	result, err := formatter.FormatContributorsCSV(contributors)
	if err != nil {
		t.Fatalf("FormatContributorsCSV() error = %v", err)
	}

	if len(result) == 0 {
		t.Fatal("FormatContributorsCSV() returned empty result")
	}

	// Parse CSV to verify structure
	reader := csv.NewReader(strings.NewReader(string(result)))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}

	// Verify header
	if len(records) == 0 {
		t.Fatal("CSV should have at least header row")
	}

	expectedHeaders := []string{
		"Name", "Email", "Total Commits", "Total Insertions", "Total Deletions",
		"First Commit", "Last Commit", "Active Days", "Activity Level",
		"Avg Commits Per Day", "Most Active Hour", "Most Active Weekday", "Top File Type",
	}

	header := records[0]
	if len(header) != len(expectedHeaders) {
		t.Errorf("Expected %d headers, got %d", len(expectedHeaders), len(header))
	}

	for i, expectedHeader := range expectedHeaders {
		if i < len(header) && header[i] != expectedHeader {
			t.Errorf("Expected header[%d] = '%s', got '%s'", i, expectedHeader, header[i])
		}
	}

	// Verify data rows
	expectedRows := len(contributors) + 1 // +1 for header
	if len(records) != expectedRows {
		t.Errorf("Expected %d rows (including header), got %d", expectedRows, len(records))
	}

	// Verify first contributor data
	if len(records) > 1 {
		firstContrib := records[1]
		if firstContrib[0] != "John Doe" {
			t.Errorf("Expected first contributor name 'John Doe', got '%s'", firstContrib[0])
		}
		if firstContrib[2] != "50" {
			t.Errorf("Expected first contributor commits '50', got '%s'", firstContrib[2])
		}
	}
}

func TestCSVFormatter_FormatCommitsCSV(t *testing.T) {
	formatter := formatters.NewCSVFormatter()
	commits := createTestCommits()

	result, err := formatter.FormatCommitsCSV(commits)
	if err != nil {
		t.Fatalf("FormatCommitsCSV() error = %v", err)
	}

	if len(result) == 0 {
		t.Fatal("FormatCommitsCSV() returned empty result")
	}

	// Parse CSV to verify structure
	reader := csv.NewReader(strings.NewReader(string(result)))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}

	// Verify header
	if len(records) == 0 {
		t.Fatal("CSV should have at least header row")
	}

	expectedHeaders := []string{
		"Hash", "Message", "Author Name", "Author Email",
		"Author Date", "Committer Name", "Committer Email", "Committer Date",
		"Files Changed", "Insertions", "Deletions",
	}

	header := records[0]
	if len(header) != len(expectedHeaders) {
		t.Errorf("Expected %d headers, got %d", len(expectedHeaders), len(header))
	}

	// Verify data rows
	expectedRows := len(commits) + 1 // +1 for header
	if len(records) != expectedRows {
		t.Errorf("Expected %d rows (including header), got %d", expectedRows, len(records))
	}

	// Verify first commit data
	if len(records) > 1 {
		firstCommit := records[1]
		if firstCommit[0] != "abc1234" {
			t.Errorf("Expected first commit hash 'abc1234', got '%s'", firstCommit[0])
		}
		if firstCommit[2] != "John Doe" {
			t.Errorf("Expected first commit author 'John Doe', got '%s'", firstCommit[2])
		}
	}
}

func TestCSVFormatter_ErrorHandling(t *testing.T) {
	formatter := formatters.NewCSVFormatter()

	// Test with nil data
	config := models.FormatConfig{Format: "csv"}
	_, err := formatter.Format(nil, config)
	if err == nil {
		t.Error("Expected error when formatting nil data")
	}

	// Test with empty analysis result
	emptyData := &models.AnalysisResult{}
	result, err := formatter.Format(emptyData, config)
	if err != nil {
		t.Fatalf("Format() should handle empty data gracefully: %v", err)
	}

	// Should still produce some output (at least newlines for empty sections)
	if len(result) == 0 {
		t.Error("Format() should produce some output even for empty data")
	}
}

func TestCSVFormatter_CSVEscaping(t *testing.T) {
	formatter := formatters.NewCSVFormatter()

	// Test contributors with special characters
	contributors := []models.Contributor{
		{
			Name:            "John \"Johnny\" Doe",
			Email:           "john,doe@example.com",
			TotalCommits:    10,
			TotalInsertions: 100,
			TotalDeletions:  50,
			FirstCommit:     time.Now().AddDate(-1, 0, 0),
			LastCommit:      time.Now(),
			ActiveDays:      5,
			CommitsByDay:    map[string]int{"2023-12-01": 5},
			CommitsByHour:   map[int]int{9: 5},
			CommitsByWeekday: map[int]int{1: 10},
			FileTypes:       map[string]int{"go": 10},
			TopFiles:        []string{"main.go"},
		},
	}

	result, err := formatter.FormatContributorsCSV(contributors)
	if err != nil {
		t.Fatalf("FormatContributorsCSV() error = %v", err)
	}

	// Parse CSV to verify proper escaping
	reader := csv.NewReader(strings.NewReader(string(result)))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to parse CSV with special characters: %v", err)
	}

	if len(records) < 2 {
		t.Fatal("Expected at least header and one data row")
	}

	// Verify that special characters are properly handled (CSV parser unescapes them)
	dataRow := records[1]
	if dataRow[0] != "John \"Johnny\" Doe" {
		t.Errorf("Expected name with quotes to be preserved after CSV parsing, got '%s'", dataRow[0])
	}
	if dataRow[1] != "john,doe@example.com" {
		t.Errorf("Expected email with comma to be preserved after CSV parsing, got '%s'", dataRow[1])
	}
}

func TestCSVFormatter_TimeFormatting(t *testing.T) {
	formatter := formatters.NewCSVFormatter()

	// Test with specific time
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	contributors := []models.Contributor{
		{
			Name:            "Test User",
			Email:           "test@example.com",
			TotalCommits:    1,
			TotalInsertions: 10,
			TotalDeletions:  5,
			FirstCommit:     testTime,
			LastCommit:      testTime.Add(24 * time.Hour),
			ActiveDays:      1,
			CommitsByDay:    map[string]int{"2023-12-25": 1},
			CommitsByHour:   map[int]int{15: 1},
			CommitsByWeekday: map[int]int{1: 1},
			FileTypes:       map[string]int{"go": 1},
			TopFiles:        []string{"test.go"},
		},
	}

	result, err := formatter.FormatContributorsCSV(contributors)
	if err != nil {
		t.Fatalf("FormatContributorsCSV() error = %v", err)
	}

	// Parse CSV to verify time formatting
	reader := csv.NewReader(strings.NewReader(string(result)))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}

	if len(records) < 2 {
		t.Fatal("Expected at least header and one data row")
	}

	dataRow := records[1]
	expectedTime := testTime.Format(time.RFC3339)
	if dataRow[5] != expectedTime {
		t.Errorf("Expected first commit time '%s', got '%s'", expectedTime, dataRow[5])
	}
}

// Helper functions for testing

func createTestCommits() []git.Commit {
	return []git.Commit{
		{
			Hash:    "abc1234",
			Message: "Initial commit",
			Author: git.Author{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			Committer: git.Author{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			AuthorDate:    time.Now().AddDate(0, 0, -2),
			CommitterDate: time.Now().AddDate(0, 0, -2),
			ParentHashes:  []string{},
			TreeHash:      "tree123",
			Stats: git.CommitStats{
				FilesChanged: 3,
				Insertions:   100,
				Deletions:    0,
				Files: []git.FileChange{
					{Path: "main.go", Status: "A", Insertions: 50, Deletions: 0},
					{Path: "utils.go", Status: "A", Insertions: 30, Deletions: 0},
					{Path: "README.md", Status: "A", Insertions: 20, Deletions: 0},
				},
			},
		},
		{
			Hash:    "def5678",
			Message: "Add feature X",
			Author: git.Author{
				Name:  "Jane Smith",
				Email: "jane@example.com",
			},
			Committer: git.Author{
				Name:  "Jane Smith",
				Email: "jane@example.com",
			},
			AuthorDate:    time.Now().AddDate(0, 0, -1),
			CommitterDate: time.Now().AddDate(0, 0, -1),
			ParentHashes:  []string{"abc1234"},
			TreeHash:      "tree456",
			Stats: git.CommitStats{
				FilesChanged: 2,
				Insertions:   75,
				Deletions:    25,
				Files: []git.FileChange{
					{Path: "feature.go", Status: "A", Insertions: 60, Deletions: 0},
					{Path: "main.go", Status: "M", Insertions: 15, Deletions: 25},
				},
			},
		},
	}
}

// Helper functions for testing (shared with json_test.go)

func createTestAnalysisResult() *models.AnalysisResult {
	return &models.AnalysisResult{
		Repository:    createTestRepository(),
		Summary:       createTestSummary(),
		Contributors:  createTestContributors(),
		ContribGraph:  createTestContributionGraph(),
		HealthMetrics: createTestHealthMetrics(),
		TimeRange:     models.TimeRange{Start: time.Now().AddDate(-1, 0, 0), End: time.Now()},
	}
}

func createTestRepository() *models.RepositoryInfo {
	return &models.RepositoryInfo{
		Path:         "/test/repo",
		Name:         "test-repo",
		TotalCommits: 100,
		FirstCommit:  time.Now().AddDate(-1, 0, 0),
		LastCommit:   time.Now(),
		Branches:     []string{"main", "develop", "feature/test"},
	}
}

func createTestSummary() *models.StatsSummary {
	return &models.StatsSummary{
		TotalCommits:     100,
		TotalInsertions:  5000,
		TotalDeletions:   2000,
		FilesChanged:     50,
		ActiveDays:       30,
		AvgCommitsPerDay: 3.33,
		CommitsByHour:    map[int]int{9: 10, 14: 15, 18: 8},
		CommitsByWeekday: map[time.Weekday]int{time.Monday: 20, time.Friday: 25},
		TopFiles: []models.FileStats{
			{Path: "main.go", Commits: 15, Insertions: 500, Deletions: 100},
			{Path: "utils.go", Commits: 10, Insertions: 300, Deletions: 50},
		},
		TopFileTypes: []models.FileTypeStats{
			{Extension: "go", Files: 10, Commits: 80, Lines: 4000},
			{Extension: "md", Files: 3, Commits: 20, Lines: 1000},
		},
	}
}

func createTestContributors() []models.Contributor {
	return []models.Contributor{
		{
			Name:            "John Doe",
			Email:           "john@example.com",
			TotalCommits:    50,
			TotalInsertions: 2500,
			TotalDeletions:  1000,
			FirstCommit:     time.Now().AddDate(-1, 0, 0),
			LastCommit:      time.Now().AddDate(0, 0, -1),
			ActiveDays:      20,
			CommitsByDay:    map[string]int{"2023-12-01": 5, "2023-12-02": 3},
			CommitsByHour:   map[int]int{9: 5, 14: 8},
			CommitsByWeekday: map[int]int{1: 10, 5: 15},
			FileTypes:       map[string]int{"go": 40, "md": 10},
			TopFiles:        []string{"main.go", "utils.go"},
		},
		{
			Name:            "Jane Smith",
			Email:           "jane@example.com",
			TotalCommits:    50,
			TotalInsertions: 2500,
			TotalDeletions:  1000,
			FirstCommit:     time.Now().AddDate(-1, 0, 0),
			LastCommit:      time.Now(),
			ActiveDays:      25,
			CommitsByDay:    map[string]int{"2023-12-01": 3, "2023-12-02": 4},
			CommitsByHour:   map[int]int{10: 6, 15: 9},
			CommitsByWeekday: map[int]int{2: 12, 4: 18},
			FileTypes:       map[string]int{"go": 35, "js": 15},
			TopFiles:        []string{"server.go", "client.js"},
		},
	}
}

func createTestContributionGraph() *models.ContributionGraph {
	return &models.ContributionGraph{
		StartDate:    time.Now().AddDate(-1, 0, 0),
		EndDate:      time.Now(),
		DailyCommits: map[string]int{"2023-12-01": 5, "2023-12-02": 3, "2023-12-03": 7},
		MaxCommits:   7,
		TotalCommits: 15,
	}
}

func createTestHealthMetrics() *models.HealthMetrics {
	return &models.HealthMetrics{
		RepositoryAge:      365 * 24 * time.Hour,
		CommitFrequency:    3.33,
		ContributorCount:   2,
		ActiveContributors: 2,
		BranchCount:        3,
		ActivityTrend:      "stable",
		MonthlyGrowth: []models.MonthlyStats{
			{Month: time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC), Commits: 45, Authors: 2},
			{Month: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC), Commits: 55, Authors: 2},
		},
	}
}
