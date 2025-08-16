// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Terminal UI components demonstration

package visualizers

import (
	"fmt"
	"git-stats/models"
	"time"
)

// DemoTerminalUI demonstrates the terminal UI components
func DemoTerminalUI() {
	fmt.Println("=== Terminal UI Components Demo ===\n")

	// Demo progress indicators
	fmt.Println("1. Progress Indicators:")
	demoProgressIndicators()
	fmt.Println()

	// Demo interactive table
	fmt.Println("2. Interactive Table:")
	demoInteractiveTable()
	fmt.Println()

	// Demo colored bar chart
	fmt.Println("3. Colored Bar Chart:")
	demoColoredBarChart()
	fmt.Println()

	// Demo status lines
	fmt.Println("4. Status Lines:")
	demoStatusLines()
	fmt.Println()

	// Demo interactive menu
	fmt.Println("5. Interactive Menu:")
	demoInteractiveMenu()
	fmt.Println()
}

func demoProgressIndicators() {
	// Progress bar style
	pi1 := NewProgressIndicator(100, "Processing commits", ProgressStyleBar)
	pi1.Update(75, "Processing commits")
	fmt.Println(pi1.RenderProgress())

	// Spinner style
	pi2 := NewProgressIndicator(100, "Analyzing repository", ProgressStyleSpinner)
	pi2.Update(50, "Analyzing repository")
	fmt.Println(pi2.RenderProgress())

	// Dots style
	pi3 := NewProgressIndicator(100, "Generating statistics", ProgressStyleDots)
	pi3.Update(90, "Generating statistics")
	fmt.Println(pi3.RenderProgress())

	// Percentage style
	pi4 := NewProgressIndicator(100, "Creating visualizations", ProgressStylePercentage)
	pi4.Update(25, "Creating visualizations")
	fmt.Println(pi4.RenderProgress())
}

func demoInteractiveTable() {
	headers := []string{"Author", "Commits", "Lines Added", "Lines Deleted"}
	rows := [][]string{
		{"Alice Johnson", "127", "3,450", "892"},
		{"Bob Smith", "89", "2,100", "456"},
		{"Charlie Brown", "156", "4,200", "1,100"},
		{"Diana Prince", "67", "1,800", "234"},
		{"Eve Wilson", "203", "5,600", "1,450"},
	}

	table := NewInteractiveTable(headers, rows)
	table.SortColumn = 1 // Sort by commits
	table.SortAsc = false // Descending order

	fmt.Println(table.RenderTable())
}

func demoColoredBarChart() {
	data := map[string]int{
		"Go":         450,
		"JavaScript": 320,
		"Python":     280,
		"TypeScript": 190,
		"Shell":      120,
		"Dockerfile": 45,
		"YAML":       30,
	}

	chart := NewColoredBarChart("Lines of Code by Language", data, 30)
	fmt.Println(chart.RenderChart())
}

func demoStatusLines() {
	statuses := []struct {
		message string
		sType   StatusType
	}{
		{"Repository analysis completed successfully", StatusSuccess},
		{"Found 1,234 commits to process", StatusInfo},
		{"Large repository detected - this may take a while", StatusWarning},
		{"Failed to access some files due to permissions", StatusError},
	}

	for _, status := range statuses {
		statusLine := NewStatusLine(status.message, status.sType, 60)
		fmt.Print(statusLine.RenderStatus())
	}
}

func demoInteractiveMenu() {
	options := []MenuOption{
		{
			Label:       "View Contribution Graph",
			Description: "Display GitHub-style contribution graph",
			Enabled:     true,
		},
		{
			Label:       "Show Statistics",
			Description: "Display detailed repository statistics",
			Enabled:     true,
		},
		{
			Label:       "Analyze Contributors",
			Description: "Show contributor analysis and rankings",
			Enabled:     true,
		},
		{
			Label:       "Export Data",
			Description: "Export analysis results to file",
			Enabled:     false,
		},
		{
			Label:       "Settings",
			Description: "Configure analysis parameters",
			Enabled:     true,
		},
	}

	menu := NewInteractiveMenu("Git Stats Analysis Options", options)
	menu.CurrentItem = 1 // Highlight second option

	fmt.Println(menu.RenderMenu())
}

// CreateSampleAnalysisResult creates sample data for testing terminal UI components
func CreateSampleAnalysisResult() *models.AnalysisResult {
	now := time.Now()

	return &models.AnalysisResult{
		Repository: &models.RepositoryInfo{
			Path:         "/path/to/repo",
			Name:         "sample-project",
			TotalCommits: 1234,
			FirstCommit:  now.AddDate(-2, 0, 0),
			LastCommit:   now,
			Branches:     []string{"main", "develop", "feature/ui"},
		},
		Summary: &models.StatsSummary{
			TotalCommits:     1234,
			TotalInsertions:  45678,
			TotalDeletions:   12345,
			FilesChanged:     567,
			ActiveDays:       180,
			AvgCommitsPerDay: 6.86,
			CommitsByHour: map[int]int{
				9:  45,
				10: 67,
				11: 89,
				14: 78,
				15: 92,
				16: 56,
			},
			CommitsByWeekday: map[time.Weekday]int{
				time.Monday:    234,
				time.Tuesday:   198,
				time.Wednesday: 267,
				time.Thursday:  189,
				time.Friday:    201,
				time.Saturday:  89,
				time.Sunday:    56,
			},
			TopFiles: []models.FileStats{
				{Path: "src/main.go", Commits: 45, Insertions: 1200, Deletions: 300},
				{Path: "src/utils.go", Commits: 32, Insertions: 800, Deletions: 150},
				{Path: "README.md", Commits: 28, Insertions: 400, Deletions: 100},
			},
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
				FirstCommit:     now.AddDate(-1, -6, 0),
				LastCommit:      now.AddDate(0, 0, -2),
				ActiveDays:      120,
			},
			{
				Name:            "Bob Smith",
				Email:           "bob@example.com",
				TotalCommits:    234,
				TotalInsertions: 8500,
				TotalDeletions:  2100,
				FirstCommit:     now.AddDate(-1, -3, 0),
				LastCommit:      now.AddDate(0, 0, -1),
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
				"2024-01-18": 8,
				"2024-01-19": 1,
			},
			MaxCommits:   8,
			TotalCommits: 1234,
		},
		HealthMetrics: &models.HealthMetrics{
			RepositoryAge:      730 * 24 * time.Hour, // 2 years
			CommitFrequency:    6.86,
			ContributorCount:   5,
			ActiveContributors: 3,
			BranchCount:        3,
			ActivityTrend:      "increasing",
			MonthlyGrowth: []models.MonthlyStats{
				{Month: now.AddDate(0, -2, 0), Commits: 89, Authors: 3},
				{Month: now.AddDate(0, -1, 0), Commits: 124, Authors: 4},
				{Month: now, Commits: 156, Authors: 5},
			},
		},
		TimeRange: models.TimeRange{
			Start: now.AddDate(-1, 0, 0),
			End:   now,
		},
	}
}
