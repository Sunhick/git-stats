// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Contribution graph visualizer

package visualizers

import (
	"fmt"
	"git-stats/models"
	"strings"
	"time"
)

// ContributionGraphRenderer implements the ContributionGraphVisualizer interface
type ContributionGraphRenderer struct {
	config     models.RenderConfig
	useColors  bool
	colorTheme string
}

// NewContributionGraphRenderer creates a new contribution graph renderer
func NewContributionGraphRenderer(config models.RenderConfig) *ContributionGraphRenderer {
	return &ContributionGraphRenderer{
		config:     config,
		useColors:  true, // Enable colors by default
		colorTheme: "github", // Default to GitHub-style colors
	}
}

// SetColorOptions configures color settings for the renderer
func (cgr *ContributionGraphRenderer) SetColorOptions(useColors bool, theme string) {
	cgr.useColors = useColors
	cgr.colorTheme = theme
}

// RenderContributionGraph renders the GitHub-style contribution graph
func (cgr *ContributionGraphRenderer) RenderContributionGraph(graph *models.ContributionGraph, config models.RenderConfig) (string, error) {
	if graph == nil {
		return "", fmt.Errorf("contribution graph cannot be nil")
	}

	var result strings.Builder

	// Calculate the grid dimensions (52 weeks * 7 days)
	startDate := graph.StartDate
	endDate := graph.EndDate

	// Ensure we have a full year view
	if endDate.Sub(startDate) < 365*24*time.Hour {
		startDate = endDate.AddDate(-1, 0, 0)
	}

	// Adjust start date to Sunday for proper week alignment
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	// Render month labels
	monthLabels := cgr.RenderMonthLabels(startDate, endDate)
	result.WriteString(monthLabels)
	result.WriteString("\n")

	// Render day-of-week indicators and contribution cells
	dayLabels := cgr.RenderDayIndicators()
	contributionCells := cgr.renderContributionCells(graph, startDate, endDate)

	// Combine day labels with contribution cells
	dayLines := strings.Split(dayLabels, "\n")
	cellLines := strings.Split(contributionCells, "\n")

	for i := 0; i < len(cellLines); i++ {
		if i < len(dayLines) {
			result.WriteString(dayLines[i])
		} else {
			result.WriteString(" ")
		}
		result.WriteString(" ")
		result.WriteString(cellLines[i])
		result.WriteString("\n")
	}

	// Add legend if enabled
	if config.ShowLegend {
		result.WriteString("\n")
		result.WriteString(cgr.renderLegend(graph.MaxCommits))
	}

	return result.String(), nil
}

// RenderMonthLabels renders the month labels above the contribution graph
func (cgr *ContributionGraphRenderer) RenderMonthLabels(startDate, endDate time.Time) string {
	var result strings.Builder

	// Add spacing for day-of-week labels (2 chars: label + space)
	result.WriteString("  ")

	// Adjust start date to Sunday for proper week alignment
	current := startDate
	for current.Weekday() != time.Sunday {
		current = current.AddDate(0, 0, -1)
	}

	// Build a slice of month labels positioned at the correct week offsets
	// Each week = 1 character cell width
	type monthMark struct {
		week int
		name string
	}
	var marks []monthMark

	weekCount := 0
	lastMonth := -1

	for weekCount < 53 && current.Before(endDate.AddDate(0, 0, 7)) {
		currentMonth := int(current.Month())

		if currentMonth != lastMonth && (current.Day() <= 7 || weekCount == 0) {
			marks = append(marks, monthMark{week: weekCount, name: current.Format("Jan")})
			lastMonth = currentMonth
		}

		current = current.AddDate(0, 0, 7)
		weekCount++
	}

	// Render labels, placing them at their week positions
	// Only show a label if there's enough room before the next one
	pos := 0
	for i, m := range marks {
		// Pad to reach this mark's position
		for pos < m.week {
			result.WriteString(" ")
			pos++
		}
		// Check space available before next label
		nextPos := weekCount
		if i+1 < len(marks) {
			nextPos = marks[i+1].week
		}
		space := nextPos - m.week
		if space >= 3 {
			result.WriteString(m.name)
			pos += len(m.name)
		} else {
			// Not enough room, just skip
			result.WriteString(" ")
			pos++
		}
	}

	return result.String()
}

// RenderDayIndicators renders the day-of-week indicators (S, M, T, W, T, F, S)
func (cgr *ContributionGraphRenderer) RenderDayIndicators() string {
	dayLabels := []string{"S", "M", "T", "W", "T", "F", "S"}
	var result strings.Builder

	for i, label := range dayLabels {
		if i%2 == 0 { // Show only every other day to avoid crowding
			result.WriteString(fmt.Sprintf("%s\n", label))
		} else {
			result.WriteString(" \n")
		}
	}

	return result.String()
}

// renderContributionCells renders the actual contribution cells
func (cgr *ContributionGraphRenderer) renderContributionCells(graph *models.ContributionGraph, startDate, endDate time.Time) string {
	var lines [7]strings.Builder // One for each day of the week

	current := startDate
	weekCount := 0

	for current.Before(endDate) || current.Equal(endDate) {
		weekday := int(current.Weekday())
		dateStr := current.Format("2006-01-02")
		commits := graph.DailyCommits[dateStr]

		// Get the appropriate character for the commit count
		cell := cgr.getCommitCell(commits, graph.MaxCommits)
		lines[weekday].WriteString(cell)

		current = current.AddDate(0, 0, 1)

		// Move to next week after Saturday
		if current.Weekday() == time.Sunday {
			weekCount++
			if weekCount > 52 { // Limit to one year
				break
			}
		}
	}

	// Combine all day lines
	var result strings.Builder
	for i := 0; i < 7; i++ {
		result.WriteString(lines[i].String())
		if i < 6 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// getCommitCell returns the appropriate character/symbol for the commit count with colors
func (cgr *ContributionGraphRenderer) getCommitCell(commits, maxCommits int) string {
	// Get the base character
	char := cgr.getCommitChar(commits, maxCommits)

	// Return with or without colors based on configuration
	if !cgr.useColors {
		return char
	}

	return cgr.getColoredCell(char, commits, maxCommits)
}

// getCommitChar returns the appropriate character for the commit count
func (cgr *ContributionGraphRenderer) getCommitChar(commits, maxCommits int) string {
	if commits == 0 {
		return "░" // Light shade for no commits
	} else if commits <= maxCommits/4 {
		return "▒" // Medium shade for low activity
	} else if commits <= maxCommits/2 {
		return "▓" // Dark shade for medium activity
	} else if commits <= maxCommits*3/4 {
		return "█" // Full block for high activity
	} else {
		return "█" // Full block for very high activity
	}
}

// getColoredCell returns the character with appropriate colors
func (cgr *ContributionGraphRenderer) getColoredCell(char string, commits, maxCommits int) string {
	colors := cgr.getColorScheme()

	if commits == 0 {
		return colors.noActivity + char + colors.reset
	} else if commits <= maxCommits/4 {
		return colors.lowActivity + char + colors.reset
	} else if commits <= maxCommits/2 {
		return colors.mediumActivity + char + colors.reset
	} else if commits <= maxCommits*3/4 {
		return colors.highActivity + char + colors.reset
	} else {
		return colors.veryHighActivity + char + colors.reset
	}
}

// ColorScheme defines colors for different activity levels
type ColorScheme struct {
	reset            string
	noActivity       string
	lowActivity      string
	mediumActivity   string
	highActivity     string
	veryHighActivity string
}

// getColorScheme returns the color scheme based on the current theme
func (cgr *ContributionGraphRenderer) getColorScheme() ColorScheme {
	switch cgr.colorTheme {
	case "github":
		return ColorScheme{
			reset:            "\033[0m",
			noActivity:       "\033[90m",   // Dark gray
			lowActivity:      "\033[92m",   // Light green
			mediumActivity:   "\033[32m",   // Medium green
			highActivity:     "\033[32;1m", // Bold green
			veryHighActivity: "\033[33;1m", // Bold yellow-green
		}
	case "blue":
		return ColorScheme{
			reset:            "\033[0m",
			noActivity:       "\033[90m",   // Dark gray
			lowActivity:      "\033[94m",   // Light blue
			mediumActivity:   "\033[34m",   // Medium blue
			highActivity:     "\033[34;1m", // Bold blue
			veryHighActivity: "\033[36;1m", // Bold cyan
		}
	case "fire":
		return ColorScheme{
			reset:            "\033[0m",
			noActivity:       "\033[90m",   // Dark gray
			lowActivity:      "\033[93m",   // Light yellow
			mediumActivity:   "\033[33m",   // Yellow
			highActivity:     "\033[31;1m", // Bold red
			veryHighActivity: "\033[35;1m", // Bold magenta
		}
	default:
		// Default to GitHub theme
		return ColorScheme{
			reset:            "\033[0m",
			noActivity:       "\033[90m",   // Dark gray
			lowActivity:      "\033[92m",   // Light green
			mediumActivity:   "\033[32m",   // Medium green
			highActivity:     "\033[32;1m", // Bold green
			veryHighActivity: "\033[33;1m", // Bold yellow-green
		}
	}
}

// RenderLegend renders the legend showing commit levels with colors (public for testing)
func (cgr *ContributionGraphRenderer) RenderLegend(maxCommits int) string {
	return cgr.renderLegend(maxCommits)
}

// renderLegend renders the legend showing commit levels with colors
func (cgr *ContributionGraphRenderer) renderLegend(maxCommits int) string {
	var result strings.Builder
	colors := cgr.getColorScheme()

	result.WriteString("Less ")

	if cgr.useColors {
		result.WriteString(colors.noActivity + "░" + colors.reset + " ")
		result.WriteString(colors.lowActivity + "▒" + colors.reset + " ")
		result.WriteString(colors.mediumActivity + "▓" + colors.reset + " ")
		result.WriteString(colors.highActivity + "█" + colors.reset + " ")
		result.WriteString(colors.veryHighActivity + "█" + colors.reset)
	} else {
		result.WriteString("░ ▒ ▓ █ █")
	}

	result.WriteString(" More\n")

	// Add numeric legend
	if cgr.useColors {
		result.WriteString(fmt.Sprintf("%s0%s   %s%d%s   %s%d%s   %s%d%s   %s%d+%s",
			colors.noActivity, colors.reset,
			colors.lowActivity, maxCommits/4, colors.reset,
			colors.mediumActivity, maxCommits/2, colors.reset,
			colors.highActivity, maxCommits*3/4, colors.reset,
			colors.veryHighActivity, maxCommits, colors.reset))
	} else {
		result.WriteString(fmt.Sprintf("0   %d   %d   %d   %d+",
			maxCommits/4,
			maxCommits/2,
			maxCommits*3/4,
			maxCommits))
	}

	return result.String()
}

// GetDayCommits returns the commit count for a specific day (used for interactive selection)
func (cgr *ContributionGraphRenderer) GetDayCommits(graph *models.ContributionGraph, date time.Time) int {
	dateStr := date.Format("2006-01-02")
	return graph.DailyCommits[dateStr]
}

// GetDateFromPosition calculates the date from a position in the contribution graph
func (cgr *ContributionGraphRenderer) GetDateFromPosition(startDate time.Time, week, day int) time.Time {
	// Adjust start date to Sunday
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	// Calculate the target date
	targetDate := startDate.AddDate(0, 0, week*7+day)
	return targetDate
}

// ValidatePosition checks if a position is valid within the contribution graph
func (cgr *ContributionGraphRenderer) ValidatePosition(week, day int) bool {
	return week >= 0 && week < 53 && day >= 0 && day < 7
}
