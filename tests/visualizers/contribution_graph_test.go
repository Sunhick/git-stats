// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Contribution graph visualizer tests

package visualizers

import (
	"git-stats/models"
	"git-stats/visualizers"
	"strings"
	"testing"
	"time"
)

func TestNewContributionGraphRenderer(t *testing.T) {
	config := models.RenderConfig{
		Width:       80,
		Height:      20,
		ColorScheme: "default",
		ShowLegend:  true,
		Interactive: false,
	}

	renderer := visualizers.NewContributionGraphRenderer(config)
	if renderer == nil {
		t.Fatal("Expected renderer to be created, got nil")
	}
}

func TestRenderContributionGraph(t *testing.T) {
	// Create test data
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	dailyCommits := make(map[string]int)
	dailyCommits["2023-01-01"] = 0
	dailyCommits["2023-01-02"] = 1
	dailyCommits["2023-01-03"] = 5
	dailyCommits["2023-01-04"] = 10

	graph := &models.ContributionGraph{
		StartDate:    startDate,
		EndDate:      endDate,
		DailyCommits: dailyCommits,
		MaxCommits:   10,
		TotalCommits: 16,
	}

	config := models.RenderConfig{
		Width:       80,
		Height:      20,
		ShowLegend:  true,
		Interactive: false,
	}

	renderer := visualizers.NewContributionGraphRenderer(config)
	result, err := renderer.RenderContributionGraph(graph, config)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	// Check that result contains month labels
	if !strings.Contains(result, "Jan") {
		t.Error("Expected result to contain month labels")
	}

	// Check that result contains day indicators
	if !strings.Contains(result, "S") {
		t.Error("Expected result to contain day indicators")
	}

	// Check that result contains contribution cells
	if !strings.Contains(result, "░") || !strings.Contains(result, "▒") ||
	   !strings.Contains(result, "▓") || !strings.Contains(result, "█") {
		t.Error("Expected result to contain contribution cells with different intensities")
	}

	// Check that legend is included when ShowLegend is true
	if !strings.Contains(result, "Less") || !strings.Contains(result, "More") {
		t.Error("Expected result to contain legend when ShowLegend is true")
	}
}

func TestRenderContributionGraphNilInput(t *testing.T) {
	config := models.RenderConfig{ShowLegend: false}
	renderer := visualizers.NewContributionGraphRenderer(config)

	_, err := renderer.RenderContributionGraph(nil, config)
	if err == nil {
		t.Fatal("Expected error for nil graph input")
	}

	if !strings.Contains(err.Error(), "cannot be nil") {
		t.Errorf("Expected error message about nil input, got: %v", err)
	}
}

func TestRenderMonthLabels(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewContributionGraphRenderer(config)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	result := renderer.RenderMonthLabels(startDate, endDate)

	if result == "" {
		t.Fatal("Expected non-empty month labels")
	}

	// Check that it contains some month names
	expectedMonths := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	monthCount := 0
	for _, month := range expectedMonths {
		if strings.Contains(result, month) {
			monthCount++
		}
	}

	if monthCount < 6 { // Should contain at least half the months
		t.Errorf("Expected at least 6 months in labels, found %d", monthCount)
	}
}

func TestRenderDayIndicators(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewContributionGraphRenderer(config)

	result := renderer.RenderDayIndicators()

	if result == "" {
		t.Fatal("Expected non-empty day indicators")
	}

	// Check that it contains day indicators
	lines := strings.Split(result, "\n")
	if len(lines) < 7 {
		t.Errorf("Expected 7 lines for day indicators, got %d", len(lines))
	}

	// Check that some lines contain day letters (only every other day is shown)
	hasS := strings.Contains(result, "S")
	hasT := strings.Contains(result, "T")

	// Since we only show every other day, we expect S, T, T, S (positions 0, 2, 4, 6)
	if !hasS || !hasT {
		t.Errorf("Expected day indicators to contain S and T (every other day). Has S:%t, T:%t", hasS, hasT)
	}

	// Verify we have 7 lines (one for each day)
	resultLines := strings.Split(strings.TrimSpace(result), "\n")
	if len(resultLines) != 7 {
		t.Errorf("Expected exactly 7 lines for day indicators, got %d", len(resultLines))
	}
}

func TestGetDayCommits(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewContributionGraphRenderer(config)

	dailyCommits := map[string]int{
		"2023-01-01": 5,
		"2023-01-02": 0,
		"2023-01-03": 10,
	}

	graph := &models.ContributionGraph{
		StartDate:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
		DailyCommits: dailyCommits,
		MaxCommits:   10,
		TotalCommits: 15,
	}

	testCases := []struct {
		date     time.Time
		expected int
	}{
		{time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), 5},
		{time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), 0},
		{time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), 10},
		{time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), 0}, // Date not in map
	}

	for _, tc := range testCases {
		result := renderer.GetDayCommits(graph, tc.date)
		if result != tc.expected {
			t.Errorf("Expected %d commits for %v, got %d",
				tc.expected, tc.date, result)
		}
	}
}

func TestGetDateFromPosition(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewContributionGraphRenderer(config)

	// Start with a Monday (2023-01-02)
	startDate := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		week     int
		day      int
		expected time.Time
	}{
		{0, 0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}, // First Sunday
		{0, 1, time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)}, // First Monday
		{1, 0, time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC)}, // Second Sunday
		{1, 6, time.Date(2023, 1, 14, 0, 0, 0, 0, time.UTC)}, // Second Saturday
	}

	for _, tc := range testCases {
		result := renderer.GetDateFromPosition(startDate, tc.week, tc.day)
		if !result.Equal(tc.expected) {
			t.Errorf("Expected date %v for week %d, day %d, got %v",
				tc.expected, tc.week, tc.day, result)
		}
	}
}

func TestValidatePosition(t *testing.T) {
	config := models.RenderConfig{}
	renderer := visualizers.NewContributionGraphRenderer(config)

	testCases := []struct {
		week     int
		day      int
		expected bool
	}{
		{0, 0, true},     // Valid position
		{25, 3, true},    // Valid middle position
		{52, 6, true},    // Valid last position
		{-1, 0, false},   // Invalid week (negative)
		{53, 0, false},   // Invalid week (too high)
		{0, -1, false},   // Invalid day (negative)
		{0, 7, false},    // Invalid day (too high)
		{-1, -1, false},  // Both invalid
	}

	for _, tc := range testCases {
		result := renderer.ValidatePosition(tc.week, tc.day)
		if result != tc.expected {
			t.Errorf("Expected %t for week %d, day %d, got %t",
				tc.expected, tc.week, tc.day, result)
		}
	}
}

func TestRenderLegend(t *testing.T) {
	config := models.RenderConfig{ShowLegend: true}
	renderer := visualizers.NewContributionGraphRenderer(config)

	graph := &models.ContributionGraph{
		StartDate:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		DailyCommits: map[string]int{"2023-01-01": 0},
		MaxCommits:   20,
		TotalCommits: 0,
	}

	result, err := renderer.RenderContributionGraph(graph, config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that legend contains expected elements
	if !strings.Contains(result, "Less") {
		t.Error("Expected legend to contain 'Less'")
	}

	if !strings.Contains(result, "More") {
		t.Error("Expected legend to contain 'More'")
	}

	// Check that legend contains the different cell types
	if !strings.Contains(result, "░") || !strings.Contains(result, "▒") ||
	   !strings.Contains(result, "▓") || !strings.Contains(result, "█") {
		t.Error("Expected legend to contain all cell types")
	}

	// Check that legend contains numeric values
	if !strings.Contains(result, "0") {
		t.Error("Expected legend to contain numeric values")
	}
}

func TestRenderContributionGraphWithoutLegend(t *testing.T) {
	graph := &models.ContributionGraph{
		StartDate:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		DailyCommits: map[string]int{"2023-01-01": 5},
		MaxCommits:   10,
		TotalCommits: 5,
	}

	config := models.RenderConfig{ShowLegend: false}
	renderer := visualizers.NewContributionGraphRenderer(config)

	result, err := renderer.RenderContributionGraph(graph, config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that legend is NOT included when ShowLegend is false
	if strings.Contains(result, "Less") || strings.Contains(result, "More") {
		t.Error("Expected result to NOT contain legend when ShowLegend is false")
	}
}
