// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Charts visualizer tests

package visualizers

import (
	"git-stats/models"
	"git-stats/visualizers"
	"strings"
	"testing"
	"time"
)

func TestNewChartsRenderer(t *testing.T) {
	config := models.RenderConfig{
		Width:       80,
		Height:      20,
		ColorScheme: "default",
		ShowLegend:  true,
		Interactive: false,
	}

	renderer := visualizers.NewChartsRenderer(config)
	if renderer == nil {
		t.Fatal("Expected renderer to be created, got nil")
	}
}

func TestRenderBarChart(t *testing.T) {
	config := models.RenderConfig{Width: 50}
	renderer := visualizers.NewChartsRenderer(config)

	data := map[string]int{
		"Go":         100,
		"JavaScript": 75,
		"Python":     50,
		"Java":       25,
	}

	result, err := renderer.RenderBarChart(data, "Programming Languages", config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	// Check that title is included
	if !strings.Contains(result, "Programming Languages") {
		t.Error("Expected result to contain title")
	}

	// Check that all data keys are present
	for key := range data {
		if !strings.Contains(result, key) {
			t.Errorf("Expected result to contain key: %s", key)
		}
	}

	// Check that bars are rendered (contains bar character)
	if !strings.Contains(result, "█") {
		t.Error("Expected result to contain bar characters")
	}

	// Check that numeric values are displayed
	if !strings.Contains(result, "100") || !strings.Contains(result, "75") {
		t.Error("Expected result to contain numeric values")
	}
}

func TestRenderBarChartNilData(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewChartsRenderer(config)

	_, err := renderer.RenderBarChart(nil, "Test", config)
	if err == nil {
		t.Fatal("Expected error for nil data")
	}

	if !strings.Contains(err.Error(), "cannot be nil") {
		t.Errorf("Expected error message about nil data, got: %v", err)
	}
}

func TestRenderBarChartEmptyData(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewChartsRenderer(config)

	data := make(map[string]int)
	_, err := renderer.RenderBarChart(data, "Test", config)
	if err == nil {
		t.Fatal("Expected error for empty data")
	}

	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Expected error message about empty data, got: %v", err)
	}
}

func TestRenderTable(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewChartsRenderer(config)

	headers := []string{"Name", "Age", "City"}
	rows := [][]string{
		{"Alice", "30", "New York"},
		{"Bob", "25", "San Francisco"},
		{"Charlie", "35", "Chicago"},
	}

	result, err := renderer.RenderTable(headers, rows, config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	// Check that headers are present
	for _, header := range headers {
		if !strings.Contains(result, header) {
			t.Errorf("Expected result to contain header: %s", header)
		}
	}

	// Check that row data is present
	for _, row := range rows {
		for _, cell := range row {
			if !strings.Contains(result, cell) {
				t.Errorf("Expected result to contain cell data: %s", cell)
			}
		}
	}

	// Check that table borders are present
	if !strings.Contains(result, "│") || !strings.Contains(result, "─") {
		t.Error("Expected result to contain table borders")
	}
}

func TestRenderTableNilHeaders(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewChartsRenderer(config)

	_, err := renderer.RenderTable(nil, [][]string{}, config)
	if err == nil {
		t.Fatal("Expected error for nil headers")
	}

	if !strings.Contains(err.Error(), "cannot be nil") {
		t.Errorf("Expected error message about nil headers, got: %v", err)
	}
}

func TestRenderTableEmptyHeaders(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewChartsRenderer(config)

	_, err := renderer.RenderTable([]string{}, [][]string{}, config)
	if err == nil {
		t.Fatal("Expected error for empty headers")
	}

	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Expected error message about empty headers, got: %v", err)
	}
}

func TestRenderSummaryStats(t *testing.T) {
	config := models.RenderConfig{Width: 50}
	renderer := visualizers.NewChartsRenderer(config)

	summary := &models.StatsSummary{
		TotalCommits:     100,
		TotalInsertions:  5000,
		TotalDeletions:   2000,
		FilesChanged:     50,
		ActiveDays:       30,
		AvgCommitsPerDay: 3.33,
		CommitsByHour: map[int]int{
			9:  10,
			14: 20,
			18: 15,
		},
		CommitsByWeekday: map[time.Weekday]int{
			time.Monday:    20,
			time.Tuesday:   15,
			time.Wednesday: 25,
			time.Thursday:  20,
			time.Friday:    20,
		},
		TopFiles: []models.FileStats{
			{
				Path:         "main.go",
				Commits:      10,
				Insertions:   500,
				Deletions:    100,
				LastModified: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		TopFileTypes: []models.FileTypeStats{
			{
				Extension: ".go",
				Files:     20,
				Commits:   80,
				Lines:     4000,
			},
		},
	}

	result, err := renderer.RenderSummaryStats(summary, config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	// Check that basic statistics are present
	if !strings.Contains(result, "Repository Statistics") {
		t.Error("Expected result to contain title")
	}

	if !strings.Contains(result, "100") { // Total commits
		t.Error("Expected result to contain total commits")
	}

	if !strings.Contains(result, "5000") { // Total insertions
		t.Error("Expected result to contain total insertions")
	}

	// Check that charts are rendered
	if !strings.Contains(result, "Commits by Hour") {
		t.Error("Expected result to contain hour chart")
	}

	if !strings.Contains(result, "Commits by Weekday") {
		t.Error("Expected result to contain weekday chart")
	}

	// Check that file information is present
	if !strings.Contains(result, "main.go") {
		t.Error("Expected result to contain file information")
	}

	if !strings.Contains(result, ".go") {
		t.Error("Expected result to contain file type information")
	}
}

func TestRenderSummaryStatsNil(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewChartsRenderer(config)

	_, err := renderer.RenderSummaryStats(nil, config)
	if err == nil {
		t.Fatal("Expected error for nil summary")
	}

	if !strings.Contains(err.Error(), "cannot be nil") {
		t.Errorf("Expected error message about nil summary, got: %v", err)
	}
}

func TestRenderContributorStats(t *testing.T) {
	config := models.RenderConfig{Width: 50}
	renderer := visualizers.NewChartsRenderer(config)

	contributors := []models.Contributor{
		{
			Name:             "Alice",
			Email:            "alice@example.com",
			TotalCommits:     50,
			TotalInsertions:  2500,
			TotalDeletions:   1000,
			ActiveDays:       20,
			FirstCommit:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			LastCommit:       time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:             "Bob",
			Email:            "bob@example.com",
			TotalCommits:     30,
			TotalInsertions:  1500,
			TotalDeletions:   500,
			ActiveDays:       15,
			FirstCommit:      time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			LastCommit:       time.Date(2023, 11, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	result, err := renderer.RenderContributorStats(contributors, config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	// Check that title is present
	if !strings.Contains(result, "Contributor Statistics") {
		t.Error("Expected result to contain title")
	}

	// Check that contributor data is present
	for _, contributor := range contributors {
		if !strings.Contains(result, contributor.Name) {
			t.Errorf("Expected result to contain contributor name: %s", contributor.Name)
		}

		if !strings.Contains(result, contributor.Email) {
			t.Errorf("Expected result to contain contributor email: %s", contributor.Email)
		}
	}

	// Check that table structure is present
	if !strings.Contains(result, "│") {
		t.Error("Expected result to contain table borders")
	}

	// Check that commits distribution chart is present
	if !strings.Contains(result, "Commits Distribution") {
		t.Error("Expected result to contain commits distribution chart")
	}
}

func TestRenderContributorStatsNilOrEmpty(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewChartsRenderer(config)

	// Test nil contributors
	_, err := renderer.RenderContributorStats(nil, config)
	if err == nil {
		t.Fatal("Expected error for nil contributors")
	}

	// Test empty contributors
	_, err = renderer.RenderContributorStats([]models.Contributor{}, config)
	if err == nil {
		t.Fatal("Expected error for empty contributors")
	}
}

func TestRenderHealthMetrics(t *testing.T) {
	config := models.RenderConfig{Width: 50}
	renderer := visualizers.NewChartsRenderer(config)

	health := &models.HealthMetrics{
		RepositoryAge:      365 * 24 * time.Hour, // 1 year
		CommitFrequency:    2.5,
		ContributorCount:   5,
		ActiveContributors: 3,
		BranchCount:        10,
		ActivityTrend:      "increasing",
		MonthlyGrowth: []models.MonthlyStats{
			{
				Month:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				Commits: 20,
				Authors: 2,
			},
			{
				Month:   time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				Commits: 25,
				Authors: 3,
			},
		},
	}

	result, err := renderer.RenderHealthMetrics(health, config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	// Check that title is present
	if !strings.Contains(result, "Repository Health Metrics") {
		t.Error("Expected result to contain title")
	}

	// Check that health metrics are present
	if !strings.Contains(result, "2.5") { // Commit frequency
		t.Error("Expected result to contain commit frequency")
	}

	if !strings.Contains(result, "increasing") { // Activity trend
		t.Error("Expected result to contain activity trend")
	}

	// Check that monthly growth chart is present
	if !strings.Contains(result, "Monthly Activity") {
		t.Error("Expected result to contain monthly activity chart")
	}
}

func TestRenderHealthMetricsNil(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewChartsRenderer(config)

	_, err := renderer.RenderHealthMetrics(nil, config)
	if err == nil {
		t.Fatal("Expected error for nil health metrics")
	}

	if !strings.Contains(err.Error(), "cannot be nil") {
		t.Errorf("Expected error message about nil health metrics, got: %v", err)
	}
}

func TestRenderTimeBasedAnalysis(t *testing.T) {
	config := models.RenderConfig{Width: 50}
	renderer := visualizers.NewChartsRenderer(config)

	summary := &models.StatsSummary{
		CommitsByHour: map[int]int{
			9:  10,
			14: 20,
			18: 15,
		},
		CommitsByWeekday: map[time.Weekday]int{
			time.Monday:    20,
			time.Tuesday:   15,
			time.Wednesday: 25,
			time.Thursday:  20,
			time.Friday:    20,
		},
	}

	result, err := renderer.RenderTimeBasedAnalysis(summary, config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	// Check that title is present
	if !strings.Contains(result, "Time-Based Analysis") {
		t.Error("Expected result to contain title")
	}

	// Check that hour analysis is present
	if !strings.Contains(result, "Activity by Hour") {
		t.Error("Expected result to contain hour analysis")
	}

	// Check that weekday analysis is present
	if !strings.Contains(result, "Activity by Day of Week") {
		t.Error("Expected result to contain weekday analysis")
	}

	// Check that time slots are present
	if !strings.Contains(result, "Morning") || !strings.Contains(result, "Afternoon") {
		t.Error("Expected result to contain time slot descriptions")
	}

	// Check that percentages are calculated
	if !strings.Contains(result, "%") {
		t.Error("Expected result to contain percentage calculations")
	}
}

func TestFormatDuration(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewChartsRenderer(config)

	testCases := []struct {
		duration time.Duration
		expected string
	}{
		{24 * time.Hour, "1 days"},
		{15 * 24 * time.Hour, "15 days"},
		{60 * 24 * time.Hour, "2 months"},
		{400 * 24 * time.Hour, "1 years, 1 months"},
		{730 * 24 * time.Hour, "2 years"},
	}

	for _, tc := range testCases {
		// This would test the private formatDuration method
		// In a real implementation, you might make this method public for testing
		// or test it through the public interface

		health := &models.HealthMetrics{
			RepositoryAge: tc.duration,
		}

		result, err := renderer.RenderHealthMetrics(health, config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Check that some duration formatting is present
		if !strings.Contains(result, "days") && !strings.Contains(result, "months") && !strings.Contains(result, "years") {
			t.Error("Expected result to contain formatted duration")
		}
	}
}

func TestGetTimeSlot(t *testing.T) {
	// This would test the private getTimeSlot method
	// Since it's private, we test it through the public interface

	config := models.RenderConfig{Width: 50}
	renderer := visualizers.NewChartsRenderer(config)

	summary := &models.StatsSummary{
		CommitsByHour: map[int]int{
			8:  10, // Morning
			14: 20, // Afternoon
			19: 15, // Evening
			2:  5,  // Night
		},
	}

	result, err := renderer.RenderTimeBasedAnalysis(summary, config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that time slot descriptions are present
	expectedSlots := []string{"Morning", "Afternoon", "Evening", "Night"}
	for _, slot := range expectedSlots {
		if !strings.Contains(result, slot) {
			t.Errorf("Expected result to contain time slot: %s", slot)
		}
	}
}
