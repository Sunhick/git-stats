// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - JSON formatter unit tests

package formatters

import (
	"encoding/json"
	"testing"
	"time"

	"git-stats/formatters"
	"git-stats/models"
)

func TestNewJSONFormatter(t *testing.T) {
	formatter := formatters.NewJSONFormatter()
	if formatter == nil {
		t.Fatal("NewJSONFormatter should return a non-nil formatter")
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	formatter := formatters.NewJSONFormatter()
	testData := createTestAnalysisResult()

	tests := []struct {
		name   string
		config models.FormatConfig
		pretty bool
	}{
		{
			name: "compact JSON format",
			config: models.FormatConfig{
				Format:   "json",
				Pretty:   false,
				Metadata: true,
			},
			pretty: false,
		},
		{
			name: "pretty JSON format",
			config: models.FormatConfig{
				Format:   "json",
				Pretty:   true,
				Metadata: true,
			},
			pretty: true,
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

			// Verify it's valid JSON
			var jsonData map[string]interface{}
			if err := json.Unmarshal(result, &jsonData); err != nil {
				t.Fatalf("Format() returned invalid JSON: %v", err)
			}

			// Verify metadata is included when requested
			if tt.config.Metadata {
				if _, exists := jsonData["metadata"]; !exists {
					t.Error("Expected metadata in JSON output when Metadata=true")
				}
			}

			// Verify required sections exist
			expectedSections := []string{"repository", "time_range", "summary", "contributors"}
			for _, section := range expectedSections {
				if _, exists := jsonData[section]; !exists {
					t.Errorf("Expected section '%s' in JSON output", section)
				}
			}
		})
	}
}

func TestJSONFormatter_FormatJSON(t *testing.T) {
	formatter := formatters.NewJSONFormatter()
	testData := createTestAnalysisResult()
	config := models.FormatConfig{
		Format:   "json",
		Pretty:   false,
		Metadata: false,
	}

	result, err := formatter.FormatJSON(testData, config)
	if err != nil {
		t.Fatalf("FormatJSON() error = %v", err)
	}

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(result, &jsonData); err != nil {
		t.Fatalf("FormatJSON() returned invalid JSON: %v", err)
	}

	// Verify repository data
	if repo, exists := jsonData["repository"].(map[string]interface{}); exists {
		if repo["name"] != "test-repo" {
			t.Errorf("Expected repository name 'test-repo', got %v", repo["name"])
		}
		if repo["total_commits"] != float64(100) {
			t.Errorf("Expected total_commits 100, got %v", repo["total_commits"])
		}
	} else {
		t.Error("Expected repository data in JSON output")
	}
}

func TestJSONFormatter_FormatPrettyJSON(t *testing.T) {
	formatter := formatters.NewJSONFormatter()
	testData := createTestAnalysisResult()
	config := models.FormatConfig{
		Format:   "json",
		Pretty:   true,
		Metadata: true,
	}

	result, err := formatter.FormatPrettyJSON(testData, config)
	if err != nil {
		t.Fatalf("FormatPrettyJSON() error = %v", err)
	}

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(result, &jsonData); err != nil {
		t.Fatalf("FormatPrettyJSON() returned invalid JSON: %v", err)
	}

	// Verify pretty formatting (should contain newlines and indentation)
	resultStr := string(result)
	if !containsIndentation(resultStr) {
		t.Error("FormatPrettyJSON() should return indented JSON")
	}

	// Verify metadata exists
	if metadata, exists := jsonData["metadata"].(map[string]interface{}); exists {
		if metadata["format"] != "json" {
			t.Errorf("Expected metadata format 'json', got %v", metadata["format"])
		}
		if metadata["version"] != "1.0" {
			t.Errorf("Expected metadata version '1.0', got %v", metadata["version"])
		}
	} else {
		t.Error("Expected metadata in pretty JSON output")
	}
}

func TestJSONFormatter_FormatContributors(t *testing.T) {
	formatter := formatters.NewJSONFormatter()
	contributors := createTestContributors()

	// Use reflection to access private method for testing
	// In a real scenario, we'd test this through the public Format method
	testData := &models.AnalysisResult{
		Contributors: contributors,
		Repository:   createTestRepository(),
		TimeRange:    models.TimeRange{Start: time.Now().AddDate(-1, 0, 0), End: time.Now()},
	}

	config := models.FormatConfig{Format: "json", Metadata: false}
	result, err := formatter.Format(testData, config)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(result, &jsonData); err != nil {
		t.Fatalf("Format() returned invalid JSON: %v", err)
	}

	contributorsData, exists := jsonData["contributors"].([]interface{})
	if !exists {
		t.Fatal("Expected contributors array in JSON output")
	}

	if len(contributorsData) != len(contributors) {
		t.Errorf("Expected %d contributors, got %d", len(contributors), len(contributorsData))
	}

	// Verify first contributor data
	if len(contributorsData) > 0 {
		contrib := contributorsData[0].(map[string]interface{})
		if contrib["name"] != "John Doe" {
			t.Errorf("Expected contributor name 'John Doe', got %v", contrib["name"])
		}
		if contrib["total_commits"] != float64(50) {
			t.Errorf("Expected total_commits 50, got %v", contrib["total_commits"])
		}
	}
}

func TestJSONFormatter_ErrorHandling(t *testing.T) {
	formatter := formatters.NewJSONFormatter()

	// Test with nil data
	config := models.FormatConfig{Format: "json"}
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

	// Verify it's still valid JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(result, &jsonData); err != nil {
		t.Fatalf("Format() returned invalid JSON for empty data: %v", err)
	}
}

func TestJSONFormatter_TimeFormatting(t *testing.T) {
	formatter := formatters.NewJSONFormatter()

	// Test with specific time
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	testData := &models.AnalysisResult{
		Repository: &models.RepositoryInfo{
			Name:        "test-repo",
			Path:        "/test/path",
			FirstCommit: testTime,
			LastCommit:  testTime.Add(24 * time.Hour),
		},
		TimeRange: models.TimeRange{Start: testTime, End: testTime.Add(24 * time.Hour)},
	}

	config := models.FormatConfig{Format: "json"}
	result, err := formatter.Format(testData, config)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(result, &jsonData); err != nil {
		t.Fatalf("Format() returned invalid JSON: %v", err)
	}

	// Verify time formatting in repository
	repo := jsonData["repository"].(map[string]interface{})
	firstCommit := repo["first_commit"].(string)
	expectedTime := testTime.Format(time.RFC3339)
	if firstCommit != expectedTime {
		t.Errorf("Expected first_commit '%s', got '%s'", expectedTime, firstCommit)
	}
}

// Helper functions for testing

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

func containsIndentation(s string) bool {
	// Check if the string contains proper JSON indentation
	return len(s) > 0 && (s[0] == '{' || s[0] == '[') &&
		   (contains(s, "\n  ") || contains(s, "\n    "))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
