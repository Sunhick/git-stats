// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Integration tests for NCurses GUI interface

//go:build gui
// +build gui

package visualizers

import (
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
