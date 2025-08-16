// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Integration tests for NCurses GUI interface

//go:build gui
// +build gui

package visualizers

import (
	"fmt"
	"git-stats/models"
	"git-stats/visualizers"
	"testing"
	"time"
)

// TestGUIInterfaceWithDependencies tests the GUI interface with actual dependencies
func TestGUIInterfaceWithDependencies(t *testing.T) {
	// This test only runs when built with -tags gui
	testData := createTestAnalysisResultForIntegration()

	gui := visualizers.NewGUIInterface()
	if gui == nil {
		t.Fatal("NewGUIInterface returned nil")
	}

	err := gui.Initialize()
	if err != nil {
		t.Errorf("Expected Initialize to succeed, got error: %v", err)
	}

	// Test that Run method exists and can be called
	// Note: We don't actually run the GUI in tests as it would block
	// This just verifies the interface is properly implemented

	err = gui.Cleanup()
	if err != nil {
		t.Errorf("Expected Cleanup to succeed, got error: %v", err)
	}
}

// TestGUIViewSwitchingWorkflow tests the complete view switching workflow
func TestGUIViewSwitchingWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForIntegration()
	state := visualizers.NewGUIState(testData)

	// Test initial state
	if state.CurrentView != visualizers.ContributionView {
		t.Errorf("Expected initial view to be ContributionView, got %v", state.CurrentView)
	}

	// Test switching to each view
	views := []visualizers.ViewType{
		visualizers.StatisticsView,
		visualizers.ContributorsView,
		visualizers.HealthView,
		visualizers.ContributionView,
	}

	for _, view := range views {
		state.SwitchView(view)
		if state.CurrentView != view {
			t.Errorf("Expected current view to be %v, got %v", view, state.CurrentView)
		}

		// Verify status message is updated
		expectedMessage := fmt.Sprintf("Switched to %s view", view.String())
		if state.StatusMessage != expectedMessage {
			t.Errorf("Expected status message '%s', got '%s'", expectedMessage, state.StatusMessage)
		}
	}
}

// TestGUINavigationWorkflow tests the complete navigation workflow
func TestGUINavigationWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForIntegration()
	state := visualizers.NewGUIState(testData)

	initialDate := state.SelectedDate
	initialStartDate := state.ViewStartDate
	initialEndDate := state.ViewEndDate

	// Test day navigation
	state.SelectDate(initialDate.AddDate(0, 0, 1))
	if state.SelectedDate.Equal(initialDate) {
		t.Error("Expected selected date to change after SelectDate")
	}

	// Test month navigation
	state.NavigateMonth(1)
	if state.ViewStartDate.Equal(initialStartDate) {
		t.Error("Expected view start date to change after NavigateMonth")
	}
	if state.ViewEndDate.Equal(initialEndDate) {
		t.Error("Expected view end date to change after NavigateMonth")
	}

	// Test year navigation
	initialStartAfterMonth := state.ViewStartDate
	initialEndAfterMonth := state.ViewEndDate

	state.NavigateYear(1)
	if state.ViewStartDate.Equal(initialStartAfterMonth) {
		t.Error("Expected view start date to change after NavigateYear")
	}
	if state.ViewEndDate.Equal(initialEndAfterMonth) {
		t.Error("Expected view end date to change after NavigateYear")
	}

	// Verify year navigation moved by exactly one year
	expectedStartDate := initialStartAfterMonth.AddDate(1, 0, 0)
	expectedEndDate := initialEndAfterMonth.AddDate(1, 0, 0)

	if !state.ViewStartDate.Equal(expectedStartDate) {
		t.Errorf("Expected start date %v, got %v", expectedStartDate, state.ViewStartDate)
	}
	if !state.ViewEndDate.Equal(expectedEndDate) {
		t.Errorf("Expected end date %v, got %v", expectedEndDate, state.ViewEndDate)
	}
}

// TestGUICommitSelectionWorkflow tests the commit selection and detail display workflow
func TestGUICommitSelectionWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForIntegration()
	state := visualizers.NewGUIState(testData)

	// Test initial state - no selected commits
	if len(state.SelectedCommits) != 0 {
		t.Errorf("Expected no selected commits initially, got %d", len(state.SelectedCommits))
	}

	// Test updating selected commits
	testCommits := []models.Commit{
		{
			Hash:    "abc12345",
			Message: "Test commit 1",
			Author: models.Author{
				Name:  "Test Author",
				Email: "test@example.com",
			},
			AuthorDate: time.Now(),
			Stats: models.CommitStats{
				FilesChanged: 2,
				Insertions:   10,
				Deletions:    5,
			},
		},
		{
			Hash:    "def67890",
			Message: "Test commit 2",
			Author: models.Author{
				Name:  "Another Author",
				Email: "another@example.com",
			},
			AuthorDate: time.Now().Add(-time.Hour),
			Stats: models.CommitStats{
				FilesChanged: 1,
				Insertions:   5,
				Deletions:    2,
			},
		},
	}

	state.UpdateSelectedCommits(testCommits)

	if len(state.SelectedCommits) != len(testCommits) {
		t.Errorf("Expected %d selected commits, got %d", len(testCommits), len(state.SelectedCommits))
	}

	// Verify commits are correctly stored
	for i, commit := range testCommits {
		if state.SelectedCommits[i].Hash != commit.Hash {
			t.Errorf("Expected commit hash %s, got %s", commit.Hash, state.SelectedCommits[i].Hash)
		}
		if state.SelectedCommits[i].Message != commit.Message {
			t.Errorf("Expected commit message %s, got %s", commit.Message, state.SelectedCommits[i].Message)
		}
	}
}

// TestGUIDetailPanelContentWorkflow tests the detail panel content updates for different views
func TestGUIDetailPanelContentWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForIntegration()
	state := visualizers.NewGUIState(testData)
	detailPanel := visualizers.NewDetailPanelWidget(state, "Test Details")

	// Test content for each view
	views := []visualizers.ViewType{
		visualizers.ContributionView,
		visualizers.StatisticsView,
		visualizers.ContributorsView,
		visualizers.HealthView,
	}

	for _, view := range views {
		state.SwitchView(view)
		detailPanel.UpdateContent()

		content := detailPanel.GetText(false)
		if content == "" {
			t.Errorf("Expected non-empty content for %s view", view.String())
		}

		// Verify view-specific content
		switch view {
		case visualizers.ContributionView:
			if !contains(content, "Selected Date:") {
				t.Errorf("Expected contribution view content to contain 'Selected Date:', got: %s", content)
			}
		case visualizers.StatisticsView:
			if !contains(content, "Total Commits:") {
				t.Errorf("Expected statistics view content to contain 'Total Commits:', got: %s", content)
			}
		case visualizers.ContributorsView:
			if !contains(content, "Total Contributors:") {
				t.Errorf("Expected contributors view content to contain 'Total Contributors:', got: %s", content)
			}
		case visualizers.HealthView:
			if !contains(content, "Repository Age:") {
				t.Errorf("Expected health view content to contain 'Repository Age:', got: %s", content)
			}
		}
	}
}

// TestGUIStatusBarWorkflow tests the status bar updates for different states
func TestGUIStatusBarWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForIntegration()
	state := visualizers.NewGUIState(testData)
	statusBar := visualizers.NewStatusBarWidget(state)

	// Test status bar for each view
	views := []visualizers.ViewType{
		visualizers.ContributionView,
		visualizers.StatisticsView,
		visualizers.ContributorsView,
		visualizers.HealthView,
	}

	for _, view := range views {
		state.SwitchView(view)
		statusBar.UpdateStatus()

		content := statusBar.GetText(false)
		if content == "" {
			t.Errorf("Expected non-empty status content for %s view", view.String())
		}

		// Verify view name is displayed
		if !contains(content, view.String()) {
			t.Errorf("Expected status bar to contain view name '%s', got: %s", view.String(), content)
		}
	}

	// Test relevant commands for each view
	for _, view := range views {
		state.SwitchView(view)
		commands := statusBar.GetRelevantCommands()

		if len(commands) == 0 {
			t.Errorf("Expected non-empty commands for %s view", view.String())
		}

		// Verify base commands are always present
		hasViewSwitching := false
		for _, cmd := range commands {
			if cmd.Rune == 'c' || cmd.Rune == 's' || cmd.Rune == 't' || cmd.Rune == 'H' {
				hasViewSwitching = true
				break
			}
		}
		if !hasViewSwitching {
			t.Errorf("Expected view switching commands for %s view", view.String())
		}
	}
}

// TestGUIContributionGraphNavigationWorkflow tests the contribution graph navigation
func TestGUIContributionGraphNavigationWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForIntegration()
	state := visualizers.NewGUIState(testData)
	widget := visualizers.NewContributionGraphWidget(testData.ContribGraph, state)

	initialDate := state.SelectedDate

	// Test day navigation (simulated key events)
	// Note: We can't actually send key events in unit tests, but we can test the navigation logic

	// Test date selection
	newDate := initialDate.AddDate(0, 0, 1)
	state.SelectDate(newDate)

	if state.SelectedDate.Equal(initialDate) {
		t.Error("Expected selected date to change")
	}

	// Test that the widget has the correct data
	if widget.Data != testData.ContribGraph {
		t.Error("Expected widget to have correct contribution graph data")
	}

	if widget.State != state {
		t.Error("Expected widget to have correct state reference")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
			 s[len(s)-len(substr):] == substr ||
			 containsInMiddle(s, substr))))
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestContributionGraphWidgetWithDependencies(t *testing.T) {
	// This test only runs when built with -tags gui
	testData := createTestAnalysisResultForIntegration()
	state := visualizers.NewGUIState(testData)

	// Test widget creation with actual dependencies
	widget := visualizers.NewContributionGraphWidget(testData.ContribGraph, state)
	if widget == nil {
		t.Fatal("NewContributionGraphWidget returned nil")
	}

	if widget.Data != testData.ContribGraph {
		t.Error("Expected widget data to be set correctly")
	}

	if widget.State != state {
		t.Error("Expected widget state to be set correctly")
	}
}

func TestDetailPanelWidgetWithDependencies(t *testing.T) {
	// This test only runs when built with -tags gui
	testData := createTestAnalysisResultForIntegration()
	state := visualizers.NewGUIState(testData)

	widget := visualizers.NewDetailPanelWidget(state, "Test Panel")
	if widget == nil {
		t.Fatal("NewDetailPanelWidget returned nil")
	}

	if widget.State != state {
		t.Error("Expected widget state to be set correctly")
	}

	// Test content update
	widget.UpdateContent()

	// Verify the widget has content
	text := widget.GetText(false)
	if text == "" {
		t.Error("Expected non-empty content after UpdateContent")
	}
}

func TestStatusBarWidgetWithDependencies(t *testing.T) {
	// This test only runs when built with -tags gui
	testData := createTestAnalysisResultForIntegration()
	state := visualizers.NewGUIState(testData)

	widget := visualizers.NewStatusBarWidget(state)
	if widget == nil {
		t.Fatal("NewStatusBarWidget returned nil")
	}

	if widget.State != state {
		t.Error("Expected widget state to be set correctly")
	}

	// Test status update
	widget.UpdateStatus()

	// Verify the widget has content
	text := widget.GetText(false)
	if text == "" {
		t.Error("Expected non-empty content after UpdateStatus")
	}
}

// Helper function for integration tests
func createTestAnalysisResultForIntegration() *models.AnalysisResult {
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
	}

	// Create test health metrics
	healthMetrics := &models.HealthMetrics{
		RepositoryAge:      time.Since(yearAgo),
		CommitFrequency:    3.33,
		ContributorCount:   1,
		ActiveContributors: 1,
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
