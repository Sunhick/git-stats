// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Tests for NCurses GUI interface

//go:build gui
// +build gui

package visualizers

import (
	"git-stats/models"
	"git-stats/visualizers"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
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

// Test enhanced widget functionality
func TestDetailPanelWidgetEnhanced(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Create detail panel widget
	detailPanel := visualizers.NewDetailPanelWidget(state, "Test Details")

	// Test initial state
	if !detailPanel.ShowDetails {
		t.Error("Expected ShowDetails to be true initially")
	}

	if detailPanel.SelectedCommitIndex != 0 {
		t.Error("Expected SelectedCommitIndex to be 0 initially")
	}

	// Test content update
	detailPanel.UpdateContent()
	content := detailPanel.GetText(true)
	if content == "" {
		t.Error("Expected content to be generated")
	}

	// Test details toggle
	detailPanel.ShowDetails = false
	detailPanel.UpdateContent()
	contentWithoutDetails := detailPanel.GetText(true)
	if contentWithoutDetails == content {
		t.Error("Expected content to change when ShowDetails is toggled")
	}
}

func TestStatusBarWidgetEnhanced(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Create status bar widget
	statusBar := visualizers.NewStatusBarWidget(state)

	// Test initial state
	if !statusBar.ShowShortcuts {
		t.Error("Expected ShowShortcuts to be true initially")
	}

	if statusBar.HelpText == "" {
		t.Error("Expected HelpText to be set")
	}

	// Test status update
	statusBar.UpdateStatus()
	content := statusBar.GetText(true)
	if content == "" {
		t.Error("Expected status content to be generated")
	}

	// Test shortcuts toggle
	statusBar.ToggleShortcuts()
	if statusBar.ShowShortcuts {
		t.Error("Expected ShowShortcuts to be false after toggle")
	}

	statusBar.UpdateStatus()
	contentWithoutShortcuts := statusBar.GetText(true)
	if contentWithoutShortcuts == content {
		t.Error("Expected content to change when shortcuts are toggled")
	}
}

func TestKeyboardInputHandling(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Create widgets
	detailPanel := visualizers.NewDetailPanelWidget(state, "Test Details")
	statusBar := visualizers.NewStatusBarWidget(state)

	// Add some test commits to state
	testCommits := []models.Commit{
		{
			Hash:    "abc123",
			Message: "Test commit 1",
			Author:  models.Author{Name: "Test User", Email: "test@example.com"},
			AuthorDate: time.Now(),
			Stats:   models.CommitStats{FilesChanged: 1, Insertions: 10, Deletions: 5},
		},
		{
			Hash:    "def456",
			Message: "Test commit 2",
			Author:  models.Author{Name: "Test User", Email: "test@example.com"},
			AuthorDate: time.Now().Add(-time.Hour),
			Stats:   models.CommitStats{FilesChanged: 2, Insertions: 20, Deletions: 10},
		},
	}
	state.SelectedCommits = testCommits

	// Test detail panel input handling
	testCases := []struct {
		name     string
		key      tcell.Key
		rune     rune
		expected int
	}{
		{"Down arrow", tcell.KeyDown, 0, 1},
		{"Up arrow", tcell.KeyUp, 0, 0},
		{"j key", tcell.KeyRune, 'j', 1},
		{"k key", tcell.KeyRune, 'k', 0},
		{"Page down", tcell.KeyPageDown, 0, 1},
		{"Page up", tcell.KeyPageUp, 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := tcell.NewEventKey(tc.key, tc.rune, tcell.ModNone)
			detailPanel.HandleInput(event)

			if detailPanel.SelectedCommitIndex != tc.expected {
				t.Errorf("Expected SelectedCommitIndex to be %d, got %d", tc.expected, detailPanel.SelectedCommitIndex)
			}
		})
	}

	// Test details toggle
	initialShowDetails := detailPanel.ShowDetails
	event := tcell.NewEventKey(tcell.KeyRune, 'd', tcell.ModNone)
	detailPanel.HandleInput(event)

	if detailPanel.ShowDetails == initialShowDetails {
		t.Error("Expected ShowDetails to toggle")
	}
}

func TestWidgetInteractions(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)

	// Create widgets
	detailPanel := visualizers.NewDetailPanelWidget(state, "Test Details")
	statusBar := visualizers.NewStatusBarWidget(state)

	// Test state changes affect widgets
	originalView := state.CurrentView
	state.SwitchView(visualizers.StatisticsView)

	// Update widgets
	detailPanel.UpdateContent()
	statusBar.UpdateStatus()

	// Verify content changed
	detailContent := detailPanel.GetText(true)
	statusContent := statusBar.GetText(true)

	if detailContent == "" {
		t.Error("Expected detail panel content to be updated")
	}

	if statusContent == "" {
		t.Error("Expected status bar content to be updated")
	}

	// Verify status bar shows current view
	if !strings.Contains(statusContent, "Statistics") {
		t.Error("Expected status bar to show current view")
	}

	// Reset state
	state.SwitchView(originalView)
}

func TestKeyBindings(t *testing.T) {
	testData := createTestAnalysisResult()
	state := visualizers.NewGUIState(testData)
	statusBar := visualizers.NewStatusBarWidget(state)

	// Test that key commands are properly defined
	if len(statusBar.Commands) == 0 {
		t.Error("Expected key commands to be defined")
	}

	// Test relevant commands for different views
	state.SwitchView(visualizers.ContributionView)
	relevantCommands := statusBar.GetRelevantCommands()

	if len(relevantCommands) == 0 {
		t.Error("Expected relevant commands for contribution view")
	}

	// Check for navigation commands in contribution view
	hasNavigation := false
	for _, cmd := range relevantCommands {
		if strings.Contains(cmd.Description, "Days") || strings.Contains(cmd.Description, "Weeks") {
			hasNavigation = true
			break
		}
	}

	if !hasNavigation {
		t.Error("Expected navigation commands for contribution view")
	}

	// Test other views
	views := []visualizers.ViewType{
		visualizers.StatisticsView,
		visualizers.ContributorsView,
		visualizers.HealthView,
	}

	for _, view := range views {
		state.SwitchView(view)
		relevantCommands := statusBar.GetRelevantCommands()

		if len(relevantCommands) == 0 {
			t.Errorf("Expected relevant commands for %s view", view.String())
		}
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
