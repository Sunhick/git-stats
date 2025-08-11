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
	config models.RenderConfig
}

// NewContributionGraphRenderer creates a new contribution graph renderer
func NewContributionGraphRenderer(config models.RenderConfig) *ContributionGraphRenderer {
	return &ContributionGraphRenderer{
		config: config,
	}
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

	for i := 0; i < len(dayLines) && i < len(cellLines); i++ {
		if strings.TrimSpace(dayLines[i]) != "" {
			result.WriteString(dayLines[i])
			result.WriteString(" ")
			result.WriteString(cellLines[i])
			result.WriteString("\n")
		}
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

	// Add spacing for day-of-week labels
	result.WriteString("    ")

	// Adjust start date to Sunday for proper week alignment
	current := startDate
	for current.Weekday() != time.Sunday {
		current = current.AddDate(0, 0, -1)
	}

	weekCount := 0
	lastMonth := -1

	for weekCount < 53 && current.Before(endDate.AddDate(0, 0, 7)) {
		currentMonth := int(current.Month())

		// Add month label at the beginning of each month
		if currentMonth != lastMonth && (current.Day() <= 7 || weekCount == 0) {
			monthName := current.Format("Jan")
			result.WriteString(fmt.Sprintf("%-4s", monthName))
			lastMonth = currentMonth
		} else {
			result.WriteString("    ") // Add spacing
		}

		current = current.AddDate(0, 0, 7) // Move to next week
		weekCount++
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

// getCommitCell returns the appropriate character/symbol for the commit count
func (cgr *ContributionGraphRenderer) getCommitCell(commits, maxCommits int) string {
	if commits == 0 {
		return "░" // Light shade for no commits
	} else if commits <= maxCommits/4 {
		return "▒" // Medium shade for low activity
	} else if commits <= maxCommits/2 {
		return "▓" // Dark shade for medium activity
	} else {
		return "█" // Full block for high activity
	}
}

// renderLegend renders the legend showing commit levels
func (cgr *ContributionGraphRenderer) renderLegend(maxCommits int) string {
	var result strings.Builder

	result.WriteString("Less ")
	result.WriteString("░ ▒ ▓ █")
	result.WriteString(" More\n")

	// Add numeric legend
	result.WriteString(fmt.Sprintf("0   %d   %d   %d+",
		maxCommits/4,
		maxCommits/2,
		maxCommits*3/4))

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
