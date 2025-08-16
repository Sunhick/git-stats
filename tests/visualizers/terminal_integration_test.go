// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Integration tests for terminal UI components

package visualizers

import (
	"git-stats/models"
	"git-stats/visualizers"
	"strings"
	"testing"
	"time"
)

func TestTerminalUIIntegration(t *testing.T) {
	// Create sample analysis result
	analysisResult := createTestAnalysisResult()

	// Test terminal UI creation
	config := models.RenderConfig{
		Width:       80,
		Height:      24,
		ColorScheme: "default",
		Interactive: true,
	}

	ui := visualizers.NewTerminalUI(config)
	if ui == nil {
		t.Fatal("Failed to create terminal UI")
	}

	// Test progress indicator with real data
	pi := visualizers.NewProgressIndicator(
		analysisResult.Repository.TotalCommits,
		"Processing commits",
		visualizers.ProgressStyleBar,
	)

	// Simulate processing progress
	for i := 0; i <= analysisResult.Repository.TotalCommits; i += 100 {
		pi.Update(i, "Processing commits")
		result := pi.RenderProgress()

		if !strings.Contains(result, "Processing commits") {
			t.Error("Progress indicator should contain the message")
		}

		if i == analysisResult.Repository.TotalCommits {
			if !strings.Contains(result, "100.0%") {
				t.Error("Progress indicator should show 100% when complete")
			}
		}
	}
}

func TestInteractiveTableWithRealData(t *testing.T) {
	analysisResult := createTestAnalysisResult()

	// Create table from contributor data
	headers := []string{"Name", "Commits", "Insertions", "Deletions", "Active Days"}
	var rows [][]string

	for _, contributor := range analysisResult.Contributors {
		row := []string{
			contributor.Name,
			string(rune(contributor.TotalCommits + '0')), // Simple conversion for test
			string(rune(contributor.TotalInsertions/1000 + '0')),
			string(rune(contributor.TotalDeletions/1000 + '0')),
			string(rune(contributor.ActiveDays + '0')),
		}
		rows = append(rows, row)
	}

	table := visualizers.NewInteractiveTable(headers, rows)
	result := table.RenderTable()

	// Verify table contains contributor data
	for _, contributor := range analysisResult.Contributors {
		if !strings.Contains(result, contributor.Name) {
			t.Errorf("Table should contain contributor name: %s", contributor.Name)
		}
	}

	// Verify headers are present
	for _, header := range headers {
		if !strings.Contains(result, header) {
			t.Errorf("Table should contain header: %s", header)
		}
	}
}

func TestColoredBarChartWithFileTypes(t *testing.T) {
	analysisResult := createTestAnalysisResult()

	// Create chart from file type data
	data := make(map[string]int)
	for _, fileType := range analysisResult.Summary.TopFileTypes {
		data[fileType.Extension] = fileType.Lines
	}

	chart := visualizers.NewColoredBarChart("Lines by File Type", data, 40)
	result := chart.RenderChart()

	// Verify chart contains file type data
	for _, fileType := range analysisResult.Summary.TopFileTypes {
		if !strings.Contains(result, fileType.Extension) {
			t.Errorf("Chart should contain file type: %s", fileType.Extension)
		}
	}

	// Verify title is present
	if !strings.Contains(result, "Lines by File Type") {
		t.Error("Chart should contain the title")
	}
}

func TestStatusLinesForDifferentScenarios(t *testing.T) {
	analysisResult := createTestAnalysisResult()

	scenarios := []struct {
		condition   bool
		message     string
		statusType  visualizers.StatusType
		description string
	}{
		{
			condition:   analysisResult.Repository.TotalCommits > 1000,
			message:     "Large repository detected",
			statusType:  visualizers.StatusWarning,
			description: "Large repo warning",
		},
		{
			condition:   len(analysisResult.Contributors) > 1,
			message:     "Multiple contributors found",
			statusType:  visualizers.StatusInfo,
			description: "Multiple contributors info",
		},
		{
			condition:   analysisResult.HealthMetrics.ActivityTrend == "increasing",
			message:     "Repository activity is increasing",
			statusType:  visualizers.StatusSuccess,
			description: "Increasing activity success",
		},
		{
			condition:   analysisResult.HealthMetrics.ActiveContributors < 2,
			message:     "Low contributor activity",
			statusType:  visualizers.StatusError,
			description: "Low activity error",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			if scenario.condition {
				statusLine := visualizers.NewStatusLine(scenario.message, scenario.statusType, 60)
				result := statusLine.RenderStatus()

				if !strings.Contains(result, scenario.message) {
					t.Errorf("Status line should contain message: %s", scenario.message)
				}

				// Verify appropriate status icon is present
				switch scenario.statusType {
				case visualizers.StatusInfo:
					if !strings.Contains(result, "ℹ") {
						t.Error("Info status should contain info icon")
					}
				case visualizers.StatusSuccess:
					if !strings.Contains(result, "✓") {
						t.Error("Success status should contain success icon")
					}
				case visualizers.StatusWarning:
					if !strings.Contains(result, "⚠") {
						t.Error("Warning status should contain warning icon")
					}
				case visualizers.StatusError:
					if !strings.Contains(result, "✗") {
						t.Error("Error status should contain error icon")
					}
				}
			}
		})
	}
}

func TestInteractiveMenuForAnalysisOptions(t *testing.T) {
	analysisResult := createTestAnalysisResult()

	// Create menu options based on available data
	var options []visualizers.MenuOption

	if analysisResult.ContribGraph != nil {
		options = append(options, visualizers.MenuOption{
			Label:       "View Contribution Graph",
			Description: "Display GitHub-style contribution graph",
			Enabled:     true,
		})
	}

	if analysisResult.Summary != nil {
		options = append(options, visualizers.MenuOption{
			Label:       "Show Statistics",
			Description: "Display detailed repository statistics",
			Enabled:     true,
		})
	}

	if len(analysisResult.Contributors) > 0 {
		options = append(options, visualizers.MenuOption{
			Label:       "Analyze Contributors",
			Description: "Show contributor analysis and rankings",
			Enabled:     true,
		})
	}

	if analysisResult.HealthMetrics != nil {
		options = append(options, visualizers.MenuOption{
			Label:       "Health Metrics",
			Description: "Display repository health indicators",
			Enabled:     true,
		})
	}

	menu := visualizers.NewInteractiveMenu("Analysis Options", options)
	result := menu.RenderMenu()

	// Verify menu contains expected options
	if !strings.Contains(result, "Analysis Options") {
		t.Error("Menu should contain the title")
	}

	for _, option := range options {
		if !strings.Contains(result, option.Label) {
			t.Errorf("Menu should contain option: %s", option.Label)
		}
		if !strings.Contains(result, option.Description) {
			t.Errorf("Menu should contain description: %s", option.Description)
		}
	}

	// Verify navigation help is present
	if !strings.Contains(result, "Navigate") {
		t.Error("Menu should contain navigation help")
	}
}

func TestProgressIndicatorStyles(t *testing.T) {
	styles := []visualizers.ProgressStyle{
		visualizers.ProgressStyleBar,
		visualizers.ProgressStyleSpinner,
		visualizers.ProgressStyleDots,
		visualizers.ProgressStylePercentage,
	}

	for _, style := range styles {
		t.Run(string(rune(int(style)+'0')), func(t *testing.T) {
			pi := visualizers.NewProgressIndicator(100, "Testing style", style)
			pi.Update(50, "Testing style")

			result := pi.RenderProgress()

			if !strings.Contains(result, "Testing style") {
				t.Error("Progress indicator should contain the message")
			}

			// Each style should produce different output
			if len(result) == 0 {
				t.Error("Progress indicator should produce output")
			}
		})
	}
}

func TestTerminalUIColorCoding(t *testing.T) {
	// Test that color codes are properly applied
	tests := []struct {
		name     string
		function func() string
		expected []string
	}{
		{
			name: "Progress bar colors",
			function: func() string {
				pi := visualizers.NewProgressIndicator(100, "Testing", visualizers.ProgressStyleBar)
				pi.Update(50, "Testing")
				return pi.RenderProgress()
			},
			expected: []string{"\033[32m", "\033[0m"}, // Green and reset
		},
		{
			name: "Status line colors",
			function: func() string {
				status := visualizers.NewStatusLine("Test", visualizers.StatusSuccess, 50)
				return status.RenderStatus()
			},
			expected: []string{"\033[32m", "\033[0m"}, // Green and reset
		},
		{
			name: "Bar chart colors",
			function: func() string {
				data := map[string]int{"Test": 100}
				chart := visualizers.NewColoredBarChart("Test", data, 20)
				return chart.RenderChart()
			},
			expected: []string{"\033[34m", "\033[0m"}, // Blue and reset
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function()

			for _, expectedColor := range tt.expected {
				if !strings.Contains(result, expectedColor) {
					t.Errorf("Expected result to contain color code: %s", expectedColor)
				}
			}
		})
	}
}

// Helper function to create test analysis result
func createTestAnalysisResult() *models.AnalysisResult {
	now := time.Now()

	return &models.AnalysisResult{
		Repository: &models.RepositoryInfo{
			Path:         "/test/repo",
			Name:         "test-project",
			TotalCommits: 1234,
			FirstCommit:  now.AddDate(-1, 0, 0),
			LastCommit:   now,
			Branches:     []string{"main", "develop"},
		},
		Summary: &models.StatsSummary{
			TotalCommits:     1234,
			TotalInsertions:  45678,
			TotalDeletions:   12345,
			FilesChanged:     567,
			ActiveDays:       180,
			AvgCommitsPerDay: 6.86,
			TopFileTypes: []models.FileTypeStats{
				{Extension: ".go", Files: 25, Commits: 567, Lines: 12000},
				{Extension: ".js", Files: 18, Commits: 234, Lines: 8500},
				{Extension: ".md", Files: 8, Commits: 89, Lines: 1200},
			},
		},
		Contributors: []models.Contributor{
			{
				Name:            "Alice Johnson",
				Email:           "alice@example.com",
				TotalCommits:    456,
				TotalInsertions: 12000,
				TotalDeletions:  3000,
				ActiveDays:      120,
			},
			{
				Name:            "Bob Smith",
				Email:           "bob@example.com",
				TotalCommits:    234,
				TotalInsertions: 8500,
				TotalDeletions:  2100,
				ActiveDays:      89,
			},
		},
		ContribGraph: &models.ContributionGraph{
			StartDate:    now.AddDate(-1, 0, 0),
			EndDate:      now,
			DailyCommits: map[string]int{
				"2024-01-15": 3,
				"2024-01-16": 5,
				"2024-01-17": 2,
			},
			MaxCommits:   5,
			TotalCommits: 1234,
		},
		HealthMetrics: &models.HealthMetrics{
			RepositoryAge:      365 * 24 * time.Hour,
			CommitFrequency:    6.86,
			ContributorCount:   2,
			ActiveContributors: 2,
			BranchCount:        2,
			ActivityTrend:      "increasing",
		},
		TimeRange: models.TimeRange{
			Start: now.AddDate(-1, 0, 0),
			End:   now,
		},
	}
}
