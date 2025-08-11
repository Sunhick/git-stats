// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for NCurses GUI components (without external dependencies)

package visualizers

import (
	"git-stats/models"
	"git-stats/visualizers"
	"testing"
	"time"
)

// TestGUIStateEnhanced tests enhanced GUI state functionality
func TestGUIStateEnhanced(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Test initial state
	if state.CurrentView != visualizers.ContributionView {
		t.Errorf("Expected initial view to be ContributionView, got %v", state.CurrentView)
	}

	// Test view switching
	state.SwitchView(visualizers.StatisticsView)
	if state.CurrentView != visualizers.StatisticsView {
		t.Errorf("Expected view to be StatisticsView, got %v", state.CurrentView)
	}

	// Test date selection with commit data
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	state.SelectDate(testDate)
	if !state.SelectedDate.Equal(testDate) {
		t.Errorf("Expected selected date to be %v, got %v", testDate, state.SelectedDate)
	}

	// Test status message contains commit count
	expectedMessage := "Selected: 2024-01-15 (5 commits)"
	if state.StatusMessage != expectedMessage {
		t.Errorf("Expected status message '%s', got '%s'", expectedMessage, state.StatusMessage)
	}
}

// TestWidgetCreationAndBasicFunctionality tests widget creation without GUI dependencies
func TestWidgetCreationAndBasicFunctionality(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Test that state contains expected data
	if state.Data == nil {
		t.Fatal("Expected state to contain data")
	}

	if state.Data.ContribGraph == nil {
		t.Error("Expected contribution graph data")
	}

	if len(state.Data.Contributors) == 0 {
		t.Error("Expected contributors data")
	}

	if state.Data.Summary == nil {
		t.Error("Expected summary data")
	}

	if state.Data.HealthMetrics == nil {
		t.Error("Expected health metrics data")
	}

	// Test view transitions
	originalView := state.CurrentView
	state.SwitchView(visualizers.HealthView)
	if state.CurrentView == originalView {
		t.Error("Expected view to change")
	}

	// Test help toggle
	originalHelp := state.ShowHelp
	state.ToggleHelp()
	if state.ShowHelp == originalHelp {
		t.Error("Expected help state to toggle")
	}
}

// TestKeyCommandStructure tests basic key command concepts
func TestKeyCommandStructure(t *testing.T) {
	// Test that we can define key command concepts
	// (KeyCommand type is only available with gui build tag)

	// Test basic key mapping concepts
	keyMappings := map[rune]string{
		'c': "[C]ontrib",
		's': "[S]tats",
		't': "[T]eam",
		'H': "[H]ealth",
	}

	if keyMappings['c'] != "[C]ontrib" {
		t.Errorf("Expected '[C]ontrib', got %s", keyMappings['c'])
	}

	if len(keyMappings) != 4 {
		t.Errorf("Expected 4 key mappings, got %d", len(keyMappings))
	}
}

// TestViewTypeString tests ViewType string representation
func TestViewTypeString(t *testing.T) {
	testCases := []struct {
		view     visualizers.ViewType
		expected string
	}{
		{visualizers.ContributionView, "Contribution"},
		{visualizers.StatisticsView, "Statistics"},
		{visualizers.ContributorsView, "Contributors"},
		{visualizers.HealthView, "Health"},
	}

	for _, tc := range testCases {
		result := tc.view.String()
		if result != tc.expected {
			t.Errorf("Expected %v.String() to be '%s', got '%s'", tc.view, tc.expected, result)
		}
	}
}

// TestGUIStateNavigation tests navigation functionality
func TestGUIStateNavigation(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	originalStart := state.ViewStartDate
	originalEnd := state.ViewEndDate

	// Test month navigation
	state.NavigateMonth(1)
	expectedStart := originalStart.AddDate(0, 1, 0)
	expectedEnd := originalEnd.AddDate(0, 1, 0)

	if !state.ViewStartDate.Equal(expectedStart) {
		t.Errorf("Expected start date %v, got %v", expectedStart, state.ViewStartDate)
	}

	if !state.ViewEndDate.Equal(expectedEnd) {
		t.Errorf("Expected end date %v, got %v", expectedEnd, state.ViewEndDate)
	}

	// Test backward navigation
	state.NavigateMonth(-2)
	expectedStart = originalStart.AddDate(0, -1, 0)
	expectedEnd = originalEnd.AddDate(0, -1, 0)

	if !state.ViewStartDate.Equal(expectedStart) {
		t.Errorf("Expected start date %v, got %v", expectedStart, state.ViewStartDate)
	}

	if !state.ViewEndDate.Equal(expectedEnd) {
		t.Errorf("Expected end date %v, got %v", expectedEnd, state.ViewEndDate)
	}
}

// Helper function to create test data
func createTestAnalysisResult() *models.AnalysisResult {
	now := time.Now()
	yearAgo := now.AddDate(-1, 0, 0)

	// Create test contribution graph
	contribGraph := &models.ContributionGraph{
		StartDate:    yearAgo,
		EndDate:      now,
		DailyCommits: make(map[string]int),
		MaxCommits:   10,
		TotalCommits: 100,
	}

	// Add some test data
	contribGraph.DailyCommits["2024-01-15"] = 5
	contribGraph.DailyCommits["2024-01-16"] = 3
	contribGraph.DailyCommits["2024-01-17"] = 8

	// Create test summary
	summary := &models.StatsSummary{
		TotalCommits:     100,
		TotalInsertions:  1000,
		TotalDeletions:   500,
		FilesChanged:     50,
		ActiveDays:       30,
		AvgCommitsPerDay: 3.33,
		CommitsByHour:    make(map[int]int),
		CommitsByWeekday: make(map[time.Weekday]int),
		TopFiles:         []models.FileStats{},
		TopFileTypes:     []models.FileTypeStats{},
	}

	// Create test contributors
	contributors := []models.Contributor{
		{
			Name:            "John Doe",
			Email:           "john@example.com",
			TotalCommits:    50,
			TotalInsertions: 500,
			TotalDeletions:  250,
			FirstCommit:     yearAgo,
			LastCommit:      now,
			ActiveDays:      20,
		},
		{
			Name:            "Jane Smith",
			Email:           "jane@example.com",
			TotalCommits:    30,
			TotalInsertions: 300,
			TotalDeletions:  150,
			FirstCommit:     yearAgo.AddDate(0, 1, 0),
			LastCommit:      now.AddDate(0, 0, -1),
			ActiveDays:      15,
		},
	}

	// Create test health metrics
	healthMetrics := &models.HealthMetrics{
		RepositoryAge:      time.Since(yearAgo),
		CommitFrequency:    3.33,
		ContributorCount:   2,
		ActiveContributors: 2,
		BranchCount:        3,
		ActivityTrend:      "stable",
		MonthlyGrowth:      []models.MonthlyStats{},
	}

	// Create test repository info
	repoInfo := &models.RepositoryInfo{
		Path:         "/test/repo",
		Name:         "test-repo",
		TotalCommits: 100,
		FirstCommit:  yearAgo,
		LastCommit:   now,
		Branches:     []string{"main", "develop", "feature"},
	}

	return &models.AnalysisResult{
		Repository:    repoInfo,
		Summary:       summary,
		Contributors:  contributors,
		ContribGraph:  contribGraph,
		HealthMetrics: healthMetrics,
		TimeRange: models.TimeRange{
			Start: yearAgo,
			End:   now,
		},
	}
}
