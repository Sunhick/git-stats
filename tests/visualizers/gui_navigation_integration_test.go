// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Integration tests for GUI navigation workflows (without GUI dependencies)

package visualizers

import (
	"fmt"
	"git-stats/models"
	"git-stats/visualizers"
	"testing"
	"time"
)

// TestCompleteNavigationWorkflow tests the complete navigation workflow
func TestCompleteNavigationWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForNavigation()
	state := visualizers.NewGUIState(testData)

	// Test initial state
	if state.CurrentView != visualizers.ContributionView {
		t.Errorf("Expected initial view to be ContributionView, got %v", state.CurrentView)
	}

	initialDate := state.SelectedDate
	initialStartDate := state.ViewStartDate
	initialEndDate := state.ViewEndDate

	// Test day navigation workflow
	t.Run("DayNavigation", func(t *testing.T) {
		// Navigate forward one day
		newDate := initialDate.AddDate(0, 0, 1)
		state.SelectDate(newDate)

		if state.SelectedDate.Equal(initialDate) {
			t.Error("Expected selected date to change after SelectDate")
		}

		if !state.SelectedDate.Equal(newDate) {
			t.Errorf("Expected selected date to be %v, got %v", newDate, state.SelectedDate)
		}

		// Navigate backward one day
		prevDate := newDate.AddDate(0, 0, -1)
		state.SelectDate(prevDate)

		if !state.SelectedDate.Equal(prevDate) {
			t.Errorf("Expected selected date to be %v, got %v", prevDate, state.SelectedDate)
		}
	})

	// Test month navigation workflow
	t.Run("MonthNavigation", func(t *testing.T) {
		// Navigate forward one month
		state.NavigateMonth(1)

		expectedStart := initialStartDate.AddDate(0, 1, 0)
		expectedEnd := initialEndDate.AddDate(0, 1, 0)

		if !state.ViewStartDate.Equal(expectedStart) {
			t.Errorf("Expected start date %v, got %v", expectedStart, state.ViewStartDate)
		}

		if !state.ViewEndDate.Equal(expectedEnd) {
			t.Errorf("Expected end date %v, got %v", expectedEnd, state.ViewEndDate)
		}

		// Navigate backward two months (should be one month before initial)
		state.NavigateMonth(-2)

		expectedStart = initialStartDate.AddDate(0, -1, 0)
		expectedEnd = initialEndDate.AddDate(0, -1, 0)

		if !state.ViewStartDate.Equal(expectedStart) {
			t.Errorf("Expected start date %v, got %v", expectedStart, state.ViewStartDate)
		}

		if !state.ViewEndDate.Equal(expectedEnd) {
			t.Errorf("Expected end date %v, got %v", expectedEnd, state.ViewEndDate)
		}
	})

	// Test year navigation workflow
	t.Run("YearNavigation", func(t *testing.T) {
		// Reset to initial state
		state.ViewStartDate = initialStartDate
		state.ViewEndDate = initialEndDate

		// Navigate forward one year
		state.NavigateYear(1)

		expectedStart := initialStartDate.AddDate(1, 0, 0)
		expectedEnd := initialEndDate.AddDate(1, 0, 0)

		if !state.ViewStartDate.Equal(expectedStart) {
			t.Errorf("Expected start date %v, got %v", expectedStart, state.ViewStartDate)
		}

		if !state.ViewEndDate.Equal(expectedEnd) {
			t.Errorf("Expected end date %v, got %v", expectedEnd, state.ViewEndDate)
		}

		// Navigate backward two years (should be one year before initial)
		state.NavigateYear(-2)

		expectedStart = initialStartDate.AddDate(-1, 0, 0)
		expectedEnd = initialEndDate.AddDate(-1, 0, 0)

		if !state.ViewStartDate.Equal(expectedStart) {
			t.Errorf("Expected start date %v, got %v", expectedStart, state.ViewStartDate)
		}

		if !state.ViewEndDate.Equal(expectedEnd) {
			t.Errorf("Expected end date %v, got %v", expectedEnd, state.ViewEndDate)
		}
	})
}

// TestCompleteViewSwitchingWorkflow tests the complete view switching workflow
func TestCompleteViewSwitchingWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForNavigation()
	state := visualizers.NewGUIState(testData)

	// Test switching through all views in sequence
	t.Run("SequentialViewSwitching", func(t *testing.T) {
		views := []visualizers.ViewType{
			visualizers.StatisticsView,
			visualizers.ContributorsView,
			visualizers.HealthView,
			visualizers.ContributionView,
		}

		for i, view := range views {
			state.SwitchView(view)

			if state.CurrentView != view {
				t.Errorf("Step %d: Expected current view to be %v, got %v", i, view, state.CurrentView)
			}

			expectedMessage := fmt.Sprintf("Switched to %s view", view.String())
			if state.StatusMessage != expectedMessage {
				t.Errorf("Step %d: Expected status message '%s', got '%s'", i, expectedMessage, state.StatusMessage)
			}
		}
	})

	// Test rapid view switching
	t.Run("RapidViewSwitching", func(t *testing.T) {
		// Switch rapidly between views
		for i := 0; i < 10; i++ {
			view := visualizers.ViewType(i % 4)
			state.SwitchView(view)

			if state.CurrentView != view {
				t.Errorf("Rapid switch %d: Expected view %v, got %v", i, view, state.CurrentView)
			}
		}
	})
}

// TestCommitSelectionWorkflow tests the complete commit selection workflow
func TestCommitSelectionWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForNavigation()
	state := visualizers.NewGUIState(testData)

	t.Run("CommitSelectionAndRetrieval", func(t *testing.T) {
		// Test initial state - no selected commits
		if len(state.SelectedCommits) != 0 {
			t.Errorf("Expected no selected commits initially, got %d", len(state.SelectedCommits))
		}

		// Test selecting a date with commits
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		state.SelectDate(testDate)

		// Verify status message includes commit count
		expectedMessage := "Selected: 2024-01-15 (5 commits)"
		if state.StatusMessage != expectedMessage {
			t.Errorf("Expected status message '%s', got '%s'", expectedMessage, state.StatusMessage)
		}

		// Test updating selected commits
		testCommits := createTestCommits(3)
		state.UpdateSelectedCommits(testCommits)

		if len(state.SelectedCommits) != len(testCommits) {
			t.Errorf("Expected %d selected commits, got %d", len(testCommits), len(state.SelectedCommits))
		}

		// Test retrieving commits for date
		retrievedCommits := state.GetCommitsForDate(testDate)
		if len(retrievedCommits) != len(testCommits) {
			t.Errorf("Expected %d retrieved commits, got %d", len(testCommits), len(retrievedCommits))
		}

		// Verify commit data integrity
		for i, commit := range testCommits {
			if state.SelectedCommits[i].Hash != commit.Hash {
				t.Errorf("Expected commit hash %s, got %s", commit.Hash, state.SelectedCommits[i].Hash)
			}
			if state.SelectedCommits[i].Message != commit.Message {
				t.Errorf("Expected commit message %s, got %s", commit.Message, state.SelectedCommits[i].Message)
			}
		}
	})

	t.Run("EmptyDateSelection", func(t *testing.T) {
		// Test selecting a date with no commits
		emptyDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
		state.SelectDate(emptyDate)

		expectedMessage := "Selected: 2024-02-01 (0 commits)"
		if state.StatusMessage != expectedMessage {
			t.Errorf("Expected status message '%s', got '%s'", expectedMessage, state.StatusMessage)
		}
	})
}

// TestCombinedNavigationWorkflow tests navigation combined with view switching
func TestCombinedNavigationWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForNavigation()
	state := visualizers.NewGUIState(testData)

	t.Run("NavigationWithViewSwitching", func(t *testing.T) {
		// Start in contribution view
		if state.CurrentView != visualizers.ContributionView {
			t.Error("Expected to start in contribution view")
		}

		// Navigate to a specific date
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		state.SelectDate(testDate)

		// Switch to statistics view
		state.SwitchView(visualizers.StatisticsView)

		// Verify date selection is preserved
		if !state.SelectedDate.Equal(testDate) {
			t.Errorf("Expected selected date to be preserved: %v, got %v", testDate, state.SelectedDate)
		}

		// Navigate months while in statistics view
		initialStart := state.ViewStartDate
		state.NavigateMonth(2)

		expectedStart := initialStart.AddDate(0, 2, 0)
		if !state.ViewStartDate.Equal(expectedStart) {
			t.Errorf("Expected start date %v, got %v", expectedStart, state.ViewStartDate)
		}

		// Switch back to contribution view
		state.SwitchView(visualizers.ContributionView)

		// Verify navigation state is preserved
		if !state.ViewStartDate.Equal(expectedStart) {
			t.Errorf("Expected navigation state to be preserved: %v, got %v", expectedStart, state.ViewStartDate)
		}
	})
}

// TestHelpToggleWorkflow tests the help toggle functionality
func TestHelpToggleWorkflow(t *testing.T) {
	testData := createTestAnalysisResultForNavigation()
	state := visualizers.NewGUIState(testData)

	t.Run("HelpToggle", func(t *testing.T) {
		// Test initial help state
		if state.ShowHelp {
			t.Error("Expected help to be hidden initially")
		}

		// Toggle help on
		state.ToggleHelp()
		if !state.ShowHelp {
			t.Error("Expected help to be shown after toggle")
		}
		if state.StatusMessage != "Help displayed" {
			t.Errorf("Expected status message 'Help displayed', got '%s'", state.StatusMessage)
		}

		// Toggle help off
		state.ToggleHelp()
		if state.ShowHelp {
			t.Error("Expected help to be hidden after second toggle")
		}
		if state.StatusMessage != "Help hidden" {
			t.Errorf("Expected status message 'Help hidden', got '%s'", state.StatusMessage)
		}
	})
}

// Helper function to create test data for navigation tests
func createTestAnalysisResultForNavigation() *models.AnalysisResult {
	now := time.Now()
	yearAgo := now.AddDate(-1, 0, 0)

	// Create test contribution graph with more data points
	contribGraph := &models.ContributionGraph{
		StartDate:    yearAgo,
		EndDate:      now,
		DailyCommits: make(map[string]int),
		MaxCommits:   15,
		TotalCommits: 250,
	}

	// Add test data for multiple dates
	contribGraph.DailyCommits["2024-01-15"] = 5
	contribGraph.DailyCommits["2024-01-16"] = 3
	contribGraph.DailyCommits["2024-01-17"] = 8
	contribGraph.DailyCommits["2024-01-18"] = 12
	contribGraph.DailyCommits["2024-01-19"] = 7
	contribGraph.DailyCommits["2024-01-20"] = 15
	contribGraph.DailyCommits["2024-01-21"] = 2

	// Create test summary
	summary := &models.StatsSummary{
		TotalCommits:     250,
		TotalInsertions:  2500,
		TotalDeletions:   1200,
		FilesChanged:     120,
		ActiveDays:       75,
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
			TotalCommits:    150,
			TotalInsertions: 1500,
			TotalDeletions:  750,
			FirstCommit:     yearAgo,
			LastCommit:      now,
			ActiveDays:      50,
		},
		{
			Name:            "Jane Smith",
			Email:           "jane@example.com",
			TotalCommits:    100,
			TotalInsertions: 1000,
			TotalDeletions:  450,
			FirstCommit:     yearAgo.AddDate(0, 1, 0),
			LastCommit:      now.AddDate(0, 0, -1),
			ActiveDays:      25,
		},
	}

	// Create test health metrics
	healthMetrics := &models.HealthMetrics{
		RepositoryAge:      time.Since(yearAgo),
		CommitFrequency:    3.33,
		ContributorCount:   2,
		ActiveContributors: 2,
		BranchCount:        5,
		ActivityTrend:      "increasing",
		MonthlyGrowth:      []models.MonthlyStats{},
	}

	// Create test repository info
	repoInfo := &models.RepositoryInfo{
		Path:         "/test/repo",
		Name:         "test-repo",
		TotalCommits: 250,
		FirstCommit:  yearAgo,
		LastCommit:   now,
		Branches:     []string{"main", "develop", "feature", "hotfix", "release"},
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

// Helper function to create test commits
func createTestCommits(count int) []models.Commit {
	commits := make([]models.Commit, count)
	baseTime := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)

	for i := 0; i < count; i++ {
		commits[i] = models.Commit{
			Hash:    fmt.Sprintf("abc123%02d", i),
			Message: fmt.Sprintf("Test commit %d", i+1),
			Author: models.Author{
				Name:  fmt.Sprintf("Author %d", i+1),
				Email: fmt.Sprintf("author%d@example.com", i+1),
			},
			AuthorDate: baseTime.Add(time.Duration(i) * time.Hour),
			Stats: models.CommitStats{
				FilesChanged: 1 + i,
				Insertions:   10 + i*5,
				Deletions:    2 + i,
			},
		}
	}

	return commits
}
