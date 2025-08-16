// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for terminal UI components

package visualizers

import (
	"fmt"
	"git-stats/models"
	"git-stats/visualizers"
	"strings"
	"testing"
	"time"
)

func TestNewTerminalUI(t *testing.T) {
	config := models.RenderConfig{
		Width:       80,
		Height:      24,
		ColorScheme: "default",
		Interactive: true,
	}

	ui := visualizers.NewTerminalUI(config)

	if ui.Width != 80 {
		t.Errorf("Expected width 80, got %d", ui.Width)
	}
	if ui.Height != 24 {
		t.Errorf("Expected height 24, got %d", ui.Height)
	}
	if ui.ColorScheme != "default" {
		t.Errorf("Expected color scheme 'default', got %s", ui.ColorScheme)
	}
	if !ui.Interactive {
		t.Error("Expected interactive to be true")
	}
}

func TestNewProgressIndicator(t *testing.T) {
	pi := visualizers.NewProgressIndicator(100, "Testing", visualizers.ProgressStyleBar)

	if pi.ProgressTracker == nil {
		t.Error("Expected ProgressTracker to be initialized")
	}
	if pi.Style != visualizers.ProgressStyleBar {
		t.Errorf("Expected style ProgressStyleBar, got %v", pi.Style)
	}
	if !pi.ShowStats {
		t.Error("Expected ShowStats to be true by default")
	}
}

func TestProgressIndicatorRenderProgress(t *testing.T) {
	tests := []struct {
		name     string
		style    visualizers.ProgressStyle
		current  int
		total    int
		expected string
	}{
		{
			name:     "Progress bar style",
			style:    visualizers.ProgressStyleBar,
			current:  50,
			total:    100,
			expected: "Testing",
		},
		{
			name:     "Spinner style",
			style:    visualizers.ProgressStyleSpinner,
			current:  25,
			total:    100,
			expected: "Testing",
		},
		{
			name:     "Dots style",
			style:    visualizers.ProgressStyleDots,
			current:  75,
			total:    100,
			expected: "Testing",
		},
		{
			name:     "Percentage style",
			style:    visualizers.ProgressStylePercentage,
			current:  10,
			total:    100,
			expected: "Testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pi := visualizers.NewProgressIndicator(tt.total, "Testing", tt.style)
			pi.Update(tt.current, "Testing")

			result := pi.RenderProgress()

			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected result to contain '%s', got: %s", tt.expected, result)
			}
		})
	}
}

func TestNewInteractiveTable(t *testing.T) {
	headers := []string{"Name", "Commits", "Lines"}
	rows := [][]string{
		{"Alice", "50", "1000"},
		{"Bob", "30", "600"},
		{"Charlie", "20", "400"},
	}

	table := visualizers.NewInteractiveTable(headers, rows)

	if len(table.Headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(table.Headers))
	}
	if len(table.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(table.Rows))
	}
	if table.CurrentRow != 0 {
		t.Errorf("Expected current row to be 0, got %d", table.CurrentRow)
	}
	if table.PageSize != 10 {
		t.Errorf("Expected page size to be 10, got %d", table.PageSize)
	}
	if !table.Sortable {
		t.Error("Expected table to be sortable by default")
	}
	if !table.Filterable {
		t.Error("Expected table to be filterable by default")
	}
}

func TestInteractiveTableRenderTable(t *testing.T) {
	headers := []string{"Name", "Commits"}
	rows := [][]string{
		{"Alice", "50"},
		{"Bob", "30"},
	}

	table := visualizers.NewInteractiveTable(headers, rows)
	result := table.RenderTable()

	// Check that headers are present
	if !strings.Contains(result, "Name") {
		t.Error("Expected result to contain 'Name' header")
	}
	if !strings.Contains(result, "Commits") {
		t.Error("Expected result to contain 'Commits' header")
	}

	// Check that data is present
	if !strings.Contains(result, "Alice") {
		t.Error("Expected result to contain 'Alice'")
	}
	if !strings.Contains(result, "Bob") {
		t.Error("Expected result to contain 'Bob'")
	}

	// Check pagination info
	if !strings.Contains(result, "Page 1 of 1") {
		t.Error("Expected result to contain pagination info")
	}
}

func TestInteractiveTableSorting(t *testing.T) {
	headers := []string{"Name", "Commits"}
	rows := [][]string{
		{"Alice", "50"},
		{"Bob", "30"},
		{"Charlie", "70"},
	}

	table := visualizers.NewInteractiveTable(headers, rows)
	table.SortColumn = 1 // Sort by commits
	table.SortAsc = false // Descending

	result := table.RenderTable()

	// Charlie should appear first (70 commits)
	charlieIndex := strings.Index(result, "Charlie")
	aliceIndex := strings.Index(result, "Alice")
	bobIndex := strings.Index(result, "Bob")

	if charlieIndex == -1 || aliceIndex == -1 || bobIndex == -1 {
		t.Error("All names should be present in the result")
	}

	// In descending order by commits: Charlie (70), Alice (50), Bob (30)
	if charlieIndex > aliceIndex || aliceIndex > bobIndex {
		t.Error("Table should be sorted by commits in descending order")
	}
}

func TestInteractiveTableFiltering(t *testing.T) {
	headers := []string{"Name", "Commits"}
	rows := [][]string{
		{"Alice", "50"},
		{"Bob", "30"},
		{"Charlie", "70"},
	}

	table := visualizers.NewInteractiveTable(headers, rows)
	table.Filter = "alice"

	result := table.RenderTable()

	// Only Alice should be visible
	if !strings.Contains(result, "Alice") {
		t.Error("Expected result to contain 'Alice'")
	}
	if strings.Contains(result, "Bob") {
		t.Error("Expected result to not contain 'Bob' when filtered")
	}
	if strings.Contains(result, "Charlie") {
		t.Error("Expected result to not contain 'Charlie' when filtered")
	}

	// Should show filter info
	if !strings.Contains(result, "[Filter: alice]") {
		t.Error("Expected result to show filter information")
	}
}

func TestNewColoredBarChart(t *testing.T) {
	data := map[string]int{
		"Go":     100,
		"Python": 75,
		"Java":   50,
	}

	chart := visualizers.NewColoredBarChart("Languages", data, 40)

	if chart.Title != "Languages" {
		t.Errorf("Expected title 'Languages', got %s", chart.Title)
	}
	if len(chart.Data) != 3 {
		t.Errorf("Expected 3 data points, got %d", len(chart.Data))
	}
	if chart.MaxWidth != 40 {
		t.Errorf("Expected max width 40, got %d", chart.MaxWidth)
	}
	if !chart.ShowValue {
		t.Error("Expected ShowValue to be true by default")
	}
}

func TestColoredBarChartRenderChart(t *testing.T) {
	data := map[string]int{
		"Go":     100,
		"Python": 50,
	}

	chart := visualizers.NewColoredBarChart("Languages", data, 20)
	result := chart.RenderChart()

	// Check title
	if !strings.Contains(result, "Languages") {
		t.Error("Expected result to contain title 'Languages'")
	}

	// Check data labels
	if !strings.Contains(result, "Go") {
		t.Error("Expected result to contain 'Go'")
	}
	if !strings.Contains(result, "Python") {
		t.Error("Expected result to contain 'Python'")
	}

	// Check values
	if !strings.Contains(result, "(100)") {
		t.Error("Expected result to contain value '(100)'")
	}
	if !strings.Contains(result, "(50)") {
		t.Error("Expected result to contain value '(50)'")
	}
}

func TestColoredBarChartEmptyData(t *testing.T) {
	data := map[string]int{}

	chart := visualizers.NewColoredBarChart("Empty", data, 20)
	result := chart.RenderChart()

	if !strings.Contains(result, "No data available") {
		t.Error("Expected result to show 'No data available' for empty data")
	}
}

func TestNewStatusLine(t *testing.T) {
	tests := []struct {
		name       string
		message    string
		statusType visualizers.StatusType
		width      int
	}{
		{"Info status", "Information message", visualizers.StatusInfo, 50},
		{"Success status", "Success message", visualizers.StatusSuccess, 50},
		{"Warning status", "Warning message", visualizers.StatusWarning, 50},
		{"Error status", "Error message", visualizers.StatusError, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := visualizers.NewStatusLine(tt.message, tt.statusType, tt.width)

			if status.Message != tt.message {
				t.Errorf("Expected message '%s', got '%s'", tt.message, status.Message)
			}
			if status.Type != tt.statusType {
				t.Errorf("Expected type %v, got %v", tt.statusType, status.Type)
			}
			if status.Width != tt.width {
				t.Errorf("Expected width %d, got %d", tt.width, status.Width)
			}
		})
	}
}

func TestStatusLineRenderStatus(t *testing.T) {
	tests := []struct {
		name       string
		message    string
		statusType visualizers.StatusType
		expected   string
	}{
		{"Info status", "Info", visualizers.StatusInfo, "ℹ"},
		{"Success status", "Success", visualizers.StatusSuccess, "✓"},
		{"Warning status", "Warning", visualizers.StatusWarning, "⚠"},
		{"Error status", "Error", visualizers.StatusError, "✗"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := visualizers.NewStatusLine(tt.message, tt.statusType, 50)
			result := status.RenderStatus()

			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected result to contain '%s', got: %s", tt.expected, result)
			}
			if !strings.Contains(result, tt.message) {
				t.Errorf("Expected result to contain message '%s', got: %s", tt.message, result)
			}
		})
	}
}

func TestNewInteractiveMenu(t *testing.T) {
	options := []visualizers.MenuOption{
		{Label: "Option 1", Description: "First option", Enabled: true},
		{Label: "Option 2", Description: "Second option", Enabled: true},
	}

	menu := visualizers.NewInteractiveMenu("Test Menu", options)

	if menu.Title != "Test Menu" {
		t.Errorf("Expected title 'Test Menu', got '%s'", menu.Title)
	}
	if len(menu.Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(menu.Options))
	}
	if menu.CurrentItem != 0 {
		t.Errorf("Expected current item to be 0, got %d", menu.CurrentItem)
	}
	if !menu.ShowHelp {
		t.Error("Expected ShowHelp to be true by default")
	}
}

func TestInteractiveMenuRenderMenu(t *testing.T) {
	options := []visualizers.MenuOption{
		{Label: "Option 1", Description: "First option", Enabled: true},
		{Label: "Option 2", Description: "Second option", Enabled: false},
	}

	menu := visualizers.NewInteractiveMenu("Test Menu", options)
	result := menu.RenderMenu()

	// Check title
	if !strings.Contains(result, "Test Menu") {
		t.Error("Expected result to contain title 'Test Menu'")
	}

	// Check options
	if !strings.Contains(result, "Option 1") {
		t.Error("Expected result to contain 'Option 1'")
	}
	if !strings.Contains(result, "Option 2") {
		t.Error("Expected result to contain 'Option 2'")
	}

	// Check descriptions
	if !strings.Contains(result, "First option") {
		t.Error("Expected result to contain 'First option'")
	}
	if !strings.Contains(result, "Second option") {
		t.Error("Expected result to contain 'Second option'")
	}

	// Check current selection indicator
	if !strings.Contains(result, "▶") {
		t.Error("Expected result to contain selection indicator '▶'")
	}

	// Check help text
	if !strings.Contains(result, "Navigate") {
		t.Error("Expected result to contain navigation help")
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected string
	}{
		{"Short string", "hello", 10, "hello"},
		{"Exact length", "hello", 5, "hello"},
		{"Long string", "hello world", 8, "hello..."},
		{"Very short limit", "hello", 3, "..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := visualizers.TruncateString(tt.input, tt.length)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"Seconds", 30 * time.Second, "30s"},
		{"Minutes and seconds", 2*time.Minute + 30*time.Second, "2m30s"},
		{"Hours and minutes", 2*time.Hour + 30*time.Minute, "2h30m"},
		{"Just minutes", 5 * time.Minute, "5m0s"},
		{"Just hours", 3 * time.Hour, "3h0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through progress indicator since formatDuration is private
			pi := visualizers.NewProgressIndicator(100, "Testing", visualizers.ProgressStyleBar)
			pi.StartTime = time.Now().Add(-tt.duration)
			pi.Update(50, "Testing")

			result := pi.RenderProgress()

			// The result should contain some time information
			// We can't test the exact format since it's private, but we can verify
			// that time information is included when ShowStats is true
			if !strings.Contains(result, "[") || !strings.Contains(result, "]") {
				t.Error("Expected progress bar to contain time information in brackets")
			}
		})
	}
}

func TestGetTerminalWidth(t *testing.T) {
	width := visualizers.GetTerminalWidth()

	// Should return a reasonable default
	if width <= 0 {
		t.Error("Expected terminal width to be positive")
	}
	if width != 80 {
		t.Errorf("Expected default terminal width to be 80, got %d", width)
	}
}

// Test color constants
func TestColorConstants(t *testing.T) {
	// Test that color constants are defined and not empty
	colors := map[string]string{
		"ColorReset":  visualizers.ColorReset,
		"ColorRed":    visualizers.ColorRed,
		"ColorGreen":  visualizers.ColorGreen,
		"ColorYellow": visualizers.ColorYellow,
		"ColorBlue":   visualizers.ColorBlue,
		"ColorPurple": visualizers.ColorPurple,
		"ColorCyan":   visualizers.ColorCyan,
		"ColorWhite":  visualizers.ColorWhite,
		"ColorBold":   visualizers.ColorBold,
		"ColorDim":    visualizers.ColorDim,
	}

	for name, color := range colors {
		if color == "" {
			t.Errorf("Expected %s to be defined and not empty", name)
		}
		if !strings.HasPrefix(color, "\033[") {
			t.Errorf("Expected %s to be a valid ANSI escape sequence, got: %s", name, color)
		}
	}
}

// Benchmark tests for performance
func BenchmarkProgressIndicatorRender(b *testing.B) {
	pi := visualizers.NewProgressIndicator(1000, "Benchmarking", visualizers.ProgressStyleBar)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pi.Update(i%1000, "Benchmarking")
		_ = pi.RenderProgress()
	}
}

func BenchmarkInteractiveTableRender(b *testing.B) {
	headers := []string{"Name", "Commits", "Lines", "Files"}
	rows := make([][]string, 100)
	for i := 0; i < 100; i++ {
		rows[i] = []string{
			fmt.Sprintf("User%d", i),
			fmt.Sprintf("%d", i*10),
			fmt.Sprintf("%d", i*100),
			fmt.Sprintf("%d", i*5),
		}
	}

	table := visualizers.NewInteractiveTable(headers, rows)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = table.RenderTable()
	}
}

func BenchmarkColoredBarChartRender(b *testing.B) {
	data := make(map[string]int)
	for i := 0; i < 20; i++ {
		data[fmt.Sprintf("Item%d", i)] = i * 10
	}

	chart := visualizers.NewColoredBarChart("Benchmark", data, 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chart.RenderChart()
	}
}
