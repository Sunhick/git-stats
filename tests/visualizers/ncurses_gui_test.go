// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Tests for NCurses GUI interface

package visualizers

import (
	"git-stats/models"
	"git-stats/visualizers"
	"testing"
	"time"
)

// TestGUIState tests the GUI state management
func TestGUIState(t *testing.T) {
	// Create test data
	testData := createTestAnalysisResult()

	// Test NewGUIState
	state := visualizers.NewGUIState(testData)
	if state == nil {
		t.Fatal("NewGUIState returned nil")
	}

	// Test initial state
	if state.CurrentView != visualizers.ContributionView {
		t.Errorf("Expected initial view to be ContributionView, got %v", state.CurrentView)
	}

	if state.ShowHelp {
		t.Error("Expected ShowHelp to be false initially")
	}

	if state.StatusMessage != "Ready" {
		t.Errorf("Expected initial status message to be 'Ready', got %s", state.StatusMessage)
	}

	if state.Data != testData {
		t.Error("Expected Data to be set to testData")
	}
}

func TestGUIStateSwitchView(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Test switching to different views
	testCases := []visualizers.ViewType{
		visualizers.StatisticsView,
		visualizers.ContributorsView,
		visualizers.HealthView,
		visualizers.ContributionView,
	}

	for _, view := range testCases {
		state.SwitchView(view)
		if state.CurrentView != view {
			t.Errorf("Expected current view to be %v, got %v", view, state.CurrentView)
		}

		expectedMessage := "Switched to " + view.String() + " view"
		if state.StatusMessage != expectedMessage {
			t.Errorf("Expected status message '%s', got '%s'", expectedMessage, state.StatusMessage)
		}
	}
}

func TestGUIStateSelectDate(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Test selecting a date with commits
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	state.SelectDate(testDate)

	if !state.SelectedDate.Equal(testDate) {
		t.Errorf("Expected selected date to be %v, got %v", testDate, state.SelectedDate)
	}

	expectedMessage := "Selected: 2024-01-15 (5 commits)"
	if state.StatusMessage != expectedMessage {
		t.Errorf("Expected status message '%s', got '%s'", expectedMessage, state.StatusMessage)
	}

	// Test selecting a date without commits
	emptyDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	state.SelectDate(emptyDate)

	expectedMessage = "Selected: 2024-01-01 (0 commits)"
	if state.StatusMessage != expectedMessage {
		t.Errorf("Expected status message '%s', got '%s'", expectedMessage, state.StatusMessage)
	}
}

func TestGUIStateNavigateMonth(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	originalStart := state.ViewStartDate
	originalEnd := state.ViewEndDate

	// Test navigating forward
	state.NavigateMonth(1)

	expectedStart := originalStart.AddDate(0, 1, 0)
	expectedEnd := originalEnd.AddDate(0, 1, 0)

	if !state.ViewStartDate.Equal(expectedStart) {
		t.Errorf("Expected start date to be %v, got %v", expectedStart, state.ViewStartDate)
	}

	if !state.ViewEndDate.Equal(expectedEnd) {
		t.Errorf("Expected end date to be %v, got %v", expectedEnd, state.ViewEndDate)
	}

	// Test navigating backward
	state.NavigateMonth(-2)

	expectedStart = originalStart.AddDate(0, -1, 0)
	expectedEnd = originalEnd.AddDate(0, -1, 0)

	if !state.ViewStartDate.Equal(expectedStart) {
		t.Errorf("Expected start date to be %v, got %v", expectedStart, state.ViewStartDate)
	}

	if !state.ViewEndDate.Equal(expectedEnd) {
		t.Errorf("Expected end date to be %v, got %v", expectedEnd, state.ViewEndDate)
	}
}

func TestGUIStateToggleHelp(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Initially help should be false
	if state.ShowHelp {
		t.Error("Expected ShowHelp to be false initially")
	}

	// Toggle help on
	state.ToggleHelp()
	if !state.ShowHelp {
		t.Error("Expected ShowHelp to be true after first toggle")
	}
	if state.StatusMessage != "Help displayed" {
		t.Errorf("Expected status message 'Help displayed', got '%s'", state.StatusMessage)
	}

	// Toggle help off
	state.ToggleHelp()
	if state.ShowHelp {
		t.Error("Expected ShowHelp to be false after second toggle")
	}
	if state.StatusMessage != "Help hidden" {
		t.Errorf("Expected status message 'Help hidden', got '%s'", state.StatusMessage)
	}
}

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

func TestGUIInterface(t *testing.T) {
	// Test interface creation
	gui := visualizers.NewGUIInterface()
	if gui == nil {
		t.Fatal("NewGUIInterface returned nil")
	}

	// Test initialization
	err := gui.Initialize()
	if err != nil {
		t.Errorf("Expected Initialize to succeed, got error: %v", err)
	}

	// Test cleanup
	err = gui.Cleanup()
	if err != nil {
		t.Errorf("Expected Cleanup to succeed, got error: %v", err)
	}
}

// Test widget creation without dependencies on external libraries
func TestWidgetCreationLogic(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Test that we can create widgets (structure validation)
	// Note: Actual widget functionality requires tview/tcell dependencies
	// which are tested in integration tests

	if state.Data != testData {
		t.Error("Expected state to contain test data")
	}

	if state.Data.ContribGraph == nil {
		t.Error("Expected contribution graph data to be available")
	}

	if len(state.Data.Contributors) == 0 {
		t.Error("Expected contributors data to be available")
	}

	if state.Data.Summary == nil {
		t.Error("Expected summary data to be available")
	}

	if state.Data.HealthMetrics == nil {
		t.Error("Expected health metrics data to be available")
	}
}

// Test state transitions and data consistency
func TestStateTransitions(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Test view transitions
	originalView := state.CurrentView
	state.SwitchView(visualizers.StatisticsView)
	if state.CurrentView == originalView {
		t.Error("Expected view to change after SwitchView")
	}

	// Test date navigation
	originalDate := state.SelectedDate
	state.SelectDate(originalDate.AddDate(0, 0, 1))
	if state.SelectedDate.Equal(originalDate) {
		t.Error("Expected selected date to change after SelectDate")
	}

	// Test month navigation
	originalStart := state.ViewStartDate
	state.NavigateMonth(1)
	if state.ViewStartDate.Equal(originalStart) {
		t.Error("Expected view start date to change after NavigateMonth")
	}

	// Test help toggle
	originalHelp := state.ShowHelp
	state.ToggleHelp()
	if state.ShowHelp == originalHelp {
		t.Error("Expected help state to change after ToggleHelp")
	}
}

// Helper functions

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
