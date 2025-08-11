// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Charts and statistics visualizer

package visualizers

import (
	"fmt"
	"git-stats/models"
	"sort"
	"strings"
	"time"
)

// ChartsRenderer implements the ChartsVisualizer interface
type ChartsRenderer struct {
	config models.RenderConfig
}

// NewChartsRenderer creates a new charts renderer
func NewChartsRenderer(config models.RenderConfig) *ChartsRenderer {
	return &ChartsRenderer{
		config: config,
	}
}

// RenderBarChart renders a horizontal bar chart for the given data
func (cr *ChartsRenderer) RenderBarChart(data map[string]int, title string, config models.RenderConfig) (string, error) {
	if data == nil || len(data) == 0 {
		return "", fmt.Errorf("data cannot be nil or empty")
	}

	var result strings.Builder

	// Add title
	if title != "" {
		result.WriteString(fmt.Sprintf("%s\n", title))
		result.WriteString(strings.Repeat("=", len(title)) + "\n\n")
	}

	// Convert map to sorted slice for consistent output
	type item struct {
		key   string
		value int
	}

	var items []item
	maxValue := 0
	for k, v := range data {
		items = append(items, item{k, v})
		if v > maxValue {
			maxValue = v
		}
	}

	// Sort by value (descending)
	sort.Slice(items, func(i, j int) bool {
		return items[i].value > items[j].value
	})

	// Calculate bar width based on config or default
	maxBarWidth := config.Width
	if maxBarWidth <= 0 {
		maxBarWidth = 50
	}

	// Find the longest key for alignment
	maxKeyLength := 0
	for _, item := range items {
		if len(item.key) > maxKeyLength {
			maxKeyLength = len(item.key)
		}
	}

	// Render each bar
	for _, item := range items {
		barLength := 0
		if maxValue > 0 {
			barLength = (item.value * maxBarWidth) / maxValue
		}

		// Ensure minimum bar length for visibility
		if item.value > 0 && barLength == 0 {
			barLength = 1
		}

		bar := strings.Repeat("█", barLength)

		result.WriteString(fmt.Sprintf("%-*s │ %s %d\n",
			maxKeyLength, item.key, bar, item.value))
	}

	return result.String(), nil
}

// RenderTable renders a formatted table with headers and rows
func (cr *ChartsRenderer) RenderTable(headers []string, rows [][]string, config models.RenderConfig) (string, error) {
	if headers == nil || len(headers) == 0 {
		return "", fmt.Errorf("headers cannot be nil or empty")
	}
	if rows == nil {
		return "", fmt.Errorf("rows cannot be nil")
	}

	var result strings.Builder

	// Calculate column widths
	colWidths := make([]int, len(headers))

	// Initialize with header lengths
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	// Check row data for maximum widths
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Add padding
	for i := range colWidths {
		colWidths[i] += 2
	}

	// Render header
	result.WriteString("│")
	for i, header := range headers {
		result.WriteString(fmt.Sprintf(" %-*s │", colWidths[i]-2, header))
	}
	result.WriteString("\n")

	// Render header separator
	result.WriteString("├")
	for i, width := range colWidths {
		result.WriteString(strings.Repeat("─", width))
		if i < len(colWidths)-1 {
			result.WriteString("┼")
		}
	}
	result.WriteString("┤\n")

	// Render rows
	for _, row := range rows {
		result.WriteString("│")
		for i, cell := range row {
			if i < len(colWidths) {
				result.WriteString(fmt.Sprintf(" %-*s │", colWidths[i]-2, cell))
			}
		}
		result.WriteString("\n")
	}

	// Render bottom border
	result.WriteString("└")
	for i, width := range colWidths {
		result.WriteString(strings.Repeat("─", width))
		if i < len(colWidths)-1 {
			result.WriteString("┴")
		}
	}
	result.WriteString("┘\n")

	return result.String(), nil
}

// RenderSummaryStats renders a comprehensive summary of repository statistics
func (cr *ChartsRenderer) RenderSummaryStats(summary *models.StatsSummary, config models.RenderConfig) (string, error) {
	if summary == nil {
		return "", fmt.Errorf("summary cannot be nil")
	}

	var result strings.Builder

	// Overall Statistics
	result.WriteString("Repository Statistics\n")
	result.WriteString("====================\n\n")

	result.WriteString(fmt.Sprintf("Total Commits:       %d\n", summary.TotalCommits))
	result.WriteString(fmt.Sprintf("Total Insertions:    %d\n", summary.TotalInsertions))
	result.WriteString(fmt.Sprintf("Total Deletions:     %d\n", summary.TotalDeletions))
	result.WriteString(fmt.Sprintf("Files Changed:       %d\n", summary.FilesChanged))
	result.WriteString(fmt.Sprintf("Active Days:         %d\n", summary.ActiveDays))
	result.WriteString(fmt.Sprintf("Avg Commits/Day:     %.2f\n\n", summary.AvgCommitsPerDay))

	// Commits by Hour
	if len(summary.CommitsByHour) > 0 {
		result.WriteString("Commits by Hour\n")
		result.WriteString("---------------\n")

		hourData := make(map[string]int)
		for hour, count := range summary.CommitsByHour {
			hourData[fmt.Sprintf("%02d:00", hour)] = count
		}

		hourChart, err := cr.RenderBarChart(hourData, "", config)
		if err == nil {
			result.WriteString(hourChart)
		}
		result.WriteString("\n")
	}

	// Commits by Weekday
	if len(summary.CommitsByWeekday) > 0 {
		result.WriteString("Commits by Weekday\n")
		result.WriteString("------------------\n")

		weekdayData := make(map[string]int)
		weekdays := []time.Weekday{
			time.Sunday, time.Monday, time.Tuesday, time.Wednesday,
			time.Thursday, time.Friday, time.Saturday,
		}

		for _, weekday := range weekdays {
			if count, exists := summary.CommitsByWeekday[weekday]; exists {
				weekdayData[weekday.String()] = count
			}
		}

		weekdayChart, err := cr.RenderBarChart(weekdayData, "", config)
		if err == nil {
			result.WriteString(weekdayChart)
		}
		result.WriteString("\n")
	}

	// Top Files
	if len(summary.TopFiles) > 0 {
		result.WriteString("Most Modified Files\n")
		result.WriteString("-------------------\n")

		headers := []string{"File", "Commits", "Insertions", "Deletions", "Last Modified"}
		var rows [][]string

		for _, file := range summary.TopFiles {
			row := []string{
				file.Path,
				fmt.Sprintf("%d", file.Commits),
				fmt.Sprintf("%d", file.Insertions),
				fmt.Sprintf("%d", file.Deletions),
				file.LastModified.Format("2006-01-02"),
			}
			rows = append(rows, row)
		}

		fileTable, err := cr.RenderTable(headers, rows, config)
		if err == nil {
			result.WriteString(fileTable)
		}
		result.WriteString("\n")
	}

	// Top File Types
	if len(summary.TopFileTypes) > 0 {
		result.WriteString("File Types\n")
		result.WriteString("----------\n")

		headers := []string{"Extension", "Files", "Commits", "Lines"}
		var rows [][]string

		for _, fileType := range summary.TopFileTypes {
			row := []string{
				fileType.Extension,
				fmt.Sprintf("%d", fileType.Files),
				fmt.Sprintf("%d", fileType.Commits),
				fmt.Sprintf("%d", fileType.Lines),
			}
			rows = append(rows, row)
		}

		typeTable, err := cr.RenderTable(headers, rows, config)
		if err == nil {
			result.WriteString(typeTable)
		}
	}

	return result.String(), nil
}

// RenderContributorStats renders contributor statistics in tabular format
func (cr *ChartsRenderer) RenderContributorStats(contributors []models.Contributor, config models.RenderConfig) (string, error) {
	if contributors == nil || len(contributors) == 0 {
		return "", fmt.Errorf("contributors cannot be nil or empty")
	}

	var result strings.Builder

	result.WriteString("Contributor Statistics\n")
	result.WriteString("======================\n\n")

	// Calculate total commits for percentage calculation
	totalCommits := 0
	for _, contributor := range contributors {
		totalCommits += contributor.TotalCommits
	}

	// Summary table with percentages (Requirement 3.1)
	headers := []string{"Name", "Email", "Commits", "Percentage", "Insertions", "Deletions", "Active Days", "First Commit", "Last Commit"}
	var rows [][]string

	for _, contributor := range contributors {
		percentage := 0.0
		if totalCommits > 0 {
			percentage = float64(contributor.TotalCommits) / float64(totalCommits) * 100
		}

		row := []string{
			contributor.Name,
			contributor.Email,
			fmt.Sprintf("%d", contributor.TotalCommits),
			fmt.Sprintf("%.1f%%", percentage),
			fmt.Sprintf("%d", contributor.TotalInsertions),
			fmt.Sprintf("%d", contributor.TotalDeletions),
			fmt.Sprintf("%d", contributor.ActiveDays),
			contributor.FirstCommit.Format("2006-01-02"),
			contributor.LastCommit.Format("2006-01-02"),
		}
		rows = append(rows, row)
	}

	table, err := cr.RenderTable(headers, rows, config)
	if err != nil {
		return "", fmt.Errorf("failed to render contributor table: %v", err)
	}

	result.WriteString(table)
	result.WriteString("\n")

	// Commits distribution chart
	result.WriteString("Commits Distribution\n")
	result.WriteString("--------------------\n")

	commitData := make(map[string]int)
	for _, contributor := range contributors {
		commitData[contributor.Name] = contributor.TotalCommits
	}

	commitChart, err := cr.RenderBarChart(commitData, "", config)
	if err == nil {
		result.WriteString(commitChart)
	}
	result.WriteString("\n")

	// Lines of code distribution chart
	result.WriteString("Lines Added Distribution\n")
	result.WriteString("------------------------\n")

	linesData := make(map[string]int)
	for _, contributor := range contributors {
		linesData[contributor.Name] = contributor.TotalInsertions
	}

	linesChart, err := cr.RenderBarChart(linesData, "", config)
	if err == nil {
		result.WriteString(linesChart)
	}

	return result.String(), nil
}

// RenderHealthMetrics renders repository health metrics
func (cr *ChartsRenderer) RenderHealthMetrics(health *models.HealthMetrics, config models.RenderConfig) (string, error) {
	if health == nil {
		return "", fmt.Errorf("health metrics cannot be nil")
	}

	var result strings.Builder

	result.WriteString("Repository Health Metrics\n")
	result.WriteString("=========================\n\n")

	// Basic metrics
	result.WriteString(fmt.Sprintf("Repository Age:       %s\n", cr.formatDuration(health.RepositoryAge)))
	result.WriteString(fmt.Sprintf("Commit Frequency:     %.2f commits/day\n", health.CommitFrequency))
	result.WriteString(fmt.Sprintf("Total Contributors:   %d\n", health.ContributorCount))
	result.WriteString(fmt.Sprintf("Active Contributors:  %d\n", health.ActiveContributors))
	result.WriteString(fmt.Sprintf("Branch Count:         %d\n", health.BranchCount))
	result.WriteString(fmt.Sprintf("Activity Trend:       %s\n\n", health.ActivityTrend))

	// Monthly growth chart
	if len(health.MonthlyGrowth) > 0 {
		result.WriteString("Monthly Activity\n")
		result.WriteString("----------------\n")

		monthData := make(map[string]int)
		for _, month := range health.MonthlyGrowth {
			monthKey := month.Month.Format("2006-01")
			monthData[monthKey] = month.Commits
		}

		monthChart, err := cr.RenderBarChart(monthData, "", config)
		if err == nil {
			result.WriteString(monthChart)
		}
	}

	return result.String(), nil
}

// formatDuration formats a duration in a human-readable way
func (cr *ChartsRenderer) formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	if days < 30 {
		return fmt.Sprintf("%d days", days)
	} else if days < 365 {
		months := days / 30
		return fmt.Sprintf("%d months", months)
	} else {
		years := days / 365
		months := (days % 365) / 30
		if months > 0 {
			return fmt.Sprintf("%d years, %d months", years, months)
		}
		return fmt.Sprintf("%d years", years)
	}
}

// RenderTimeBasedAnalysis renders time-based analysis charts
func (cr *ChartsRenderer) RenderTimeBasedAnalysis(summary *models.StatsSummary, config models.RenderConfig) (string, error) {
	if summary == nil {
		return "", fmt.Errorf("summary cannot be nil")
	}

	var result strings.Builder

	result.WriteString("Time-Based Analysis\n")
	result.WriteString("===================\n\n")

	// Hour distribution
	if len(summary.CommitsByHour) > 0 {
		result.WriteString("Activity by Hour of Day\n")
		result.WriteString("-----------------------\n")

		// Create a more detailed hour analysis
		hourData := make(map[string]int)
		for hour := 0; hour < 24; hour++ {
			count := summary.CommitsByHour[hour]
			timeSlot := cr.getTimeSlot(hour)
			hourData[fmt.Sprintf("%02d:00 (%s)", hour, timeSlot)] = count
		}

		hourChart, err := cr.RenderBarChart(hourData, "", config)
		if err == nil {
			result.WriteString(hourChart)
		}
		result.WriteString("\n")
	}

	// Weekday distribution with percentages
	if len(summary.CommitsByWeekday) > 0 {
		result.WriteString("Activity by Day of Week\n")
		result.WriteString("-----------------------\n")

		totalCommits := 0
		for _, count := range summary.CommitsByWeekday {
			totalCommits += count
		}

		headers := []string{"Day", "Commits", "Percentage"}
		var rows [][]string

		weekdays := []time.Weekday{
			time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
			time.Friday, time.Saturday, time.Sunday,
		}

		for _, weekday := range weekdays {
			count := summary.CommitsByWeekday[weekday]
			percentage := 0.0
			if totalCommits > 0 {
				percentage = float64(count) / float64(totalCommits) * 100
			}

			row := []string{
				weekday.String(),
				fmt.Sprintf("%d", count),
				fmt.Sprintf("%.1f%%", percentage),
			}
			rows = append(rows, row)
		}

		weekdayTable, err := cr.RenderTable(headers, rows, config)
		if err == nil {
			result.WriteString(weekdayTable)
		}
	}

	return result.String(), nil
}

// getTimeSlot returns a descriptive time slot for the given hour
func (cr *ChartsRenderer) getTimeSlot(hour int) string {
	switch {
	case hour >= 6 && hour < 12:
		return "Morning"
	case hour >= 12 && hour < 18:
		return "Afternoon"
	case hour >= 18 && hour < 22:
		return "Evening"
	default:
		return "Night"
	}
}

// RenderCollaborationPatterns renders collaboration patterns between contributors (Requirement 3.3)
func (cr *ChartsRenderer) RenderCollaborationPatterns(contributors []models.Contributor, config models.RenderConfig) (string, error) {
	if contributors == nil || len(contributors) == 0 {
		return "", fmt.Errorf("contributors cannot be nil or empty")
	}

	var result strings.Builder

	result.WriteString("Collaboration Patterns\n")
	result.WriteString("======================\n\n")

	// File overlap analysis - contributors working on similar files
	result.WriteString("File Overlap Analysis\n")
	result.WriteString("---------------------\n")

	// Create a map of files to contributors
	fileContributors := make(map[string][]string)
	for _, contributor := range contributors {
		for _, file := range contributor.TopFiles {
			fileContributors[file] = append(fileContributors[file], contributor.Name)
		}
	}

	// Find files with multiple contributors
	collaborativeFiles := make(map[string]int)
	for file, contribs := range fileContributors {
		if len(contribs) > 1 {
			collaborativeFiles[file] = len(contribs)
		}
	}

	if len(collaborativeFiles) > 0 {
		headers := []string{"File", "Contributors", "Collaboration Level"}
		var rows [][]string

		// Sort by collaboration level
		type fileCollab struct {
			file  string
			count int
		}
		var sortedFiles []fileCollab
		for file, count := range collaborativeFiles {
			sortedFiles = append(sortedFiles, fileCollab{file, count})
		}

		// Sort by count descending
		for i := 0; i < len(sortedFiles)-1; i++ {
			for j := i + 1; j < len(sortedFiles); j++ {
				if sortedFiles[i].count < sortedFiles[j].count {
					sortedFiles[i], sortedFiles[j] = sortedFiles[j], sortedFiles[i]
				}
			}
		}

		for _, fc := range sortedFiles {
			level := "Low"
			if fc.count >= 5 {
				level = "High"
			} else if fc.count >= 3 {
				level = "Medium"
			}

			contributorNames := strings.Join(fileContributors[fc.file], ", ")
			if len(contributorNames) > 50 {
				contributorNames = contributorNames[:47] + "..."
			}

			row := []string{
				fc.file,
				contributorNames,
				level,
			}
			rows = append(rows, row)
		}

		table, err := cr.RenderTable(headers, rows, config)
		if err == nil {
			result.WriteString(table)
		}
		result.WriteString("\n")
	} else {
		result.WriteString("No collaborative files found.\n\n")
	}

	// Activity overlap analysis - contributors active in similar time periods
	result.WriteString("Activity Overlap\n")
	result.WriteString("----------------\n")

	overlapData := make(map[string]int)
	for i, contrib1 := range contributors {
		for j, contrib2 := range contributors {
			if i >= j {
				continue
			}

			// Check if their active periods overlap
			overlap := cr.calculateTimeOverlap(contrib1, contrib2)
			if overlap > 0 {
				pairName := fmt.Sprintf("%s & %s", contrib1.Name, contrib2.Name)
				overlapData[pairName] = overlap
			}
		}
	}

	if len(overlapData) > 0 {
		overlapChart, err := cr.RenderBarChart(overlapData, "", config)
		if err == nil {
			result.WriteString(overlapChart)
		}
	} else {
		result.WriteString("No significant activity overlap found.\n")
	}

	return result.String(), nil
}

// calculateTimeOverlap calculates the overlap in active days between two contributors
func (cr *ChartsRenderer) calculateTimeOverlap(contrib1, contrib2 models.Contributor) int {
	// Simple overlap calculation based on active periods
	start1, end1 := contrib1.FirstCommit, contrib1.LastCommit
	start2, end2 := contrib2.FirstCommit, contrib2.LastCommit

	// Find overlap period
	overlapStart := start1
	if start2.After(start1) {
		overlapStart = start2
	}

	overlapEnd := end1
	if end2.Before(end1) {
		overlapEnd = end2
	}

	if overlapStart.Before(overlapEnd) {
		return int(overlapEnd.Sub(overlapStart).Hours() / 24)
	}

	return 0
}

// RenderFrequencyAnalysis renders commit frequency analysis (Requirement 2.3)
func (cr *ChartsRenderer) RenderFrequencyAnalysis(summary *models.StatsSummary, config models.RenderConfig) (string, error) {
	if summary == nil {
		return "", fmt.Errorf("summary cannot be nil")
	}

	var result strings.Builder

	result.WriteString("Commit Frequency Analysis\n")
	result.WriteString("=========================\n\n")

	// Basic frequency metrics
	result.WriteString("Frequency Metrics\n")
	result.WriteString("-----------------\n")
	result.WriteString(fmt.Sprintf("Average Commits per Day:    %.2f\n", summary.AvgCommitsPerDay))
	result.WriteString(fmt.Sprintf("Average Commits per Week:   %.2f\n", summary.AvgCommitsPerDay*7))
	result.WriteString(fmt.Sprintf("Average Commits per Month:  %.2f\n", summary.AvgCommitsPerDay*30))
	result.WriteString(fmt.Sprintf("Total Active Days:          %d\n", summary.ActiveDays))

	if summary.TotalCommits > 0 && summary.ActiveDays > 0 {
		result.WriteString(fmt.Sprintf("Activity Ratio:             %.1f%% (active days / total period)\n",
			float64(summary.ActiveDays)/float64(summary.TotalCommits)*100))
	}
	result.WriteString("\n")

	// Peak activity analysis
	result.WriteString("Peak Activity Analysis\n")
	result.WriteString("----------------------\n")

	// Find peak hour
	maxHourCommits := 0
	peakHour := 0
	for hour, commits := range summary.CommitsByHour {
		if commits > maxHourCommits {
			maxHourCommits = commits
			peakHour = hour
		}
	}

	// Find peak weekday
	maxWeekdayCommits := 0
	var peakWeekday time.Weekday
	for weekday, commits := range summary.CommitsByWeekday {
		if commits > maxWeekdayCommits {
			maxWeekdayCommits = commits
			peakWeekday = weekday
		}
	}

	result.WriteString(fmt.Sprintf("Peak Hour:                  %02d:00 (%d commits)\n", peakHour, maxHourCommits))
	result.WriteString(fmt.Sprintf("Peak Weekday:               %s (%d commits)\n", peakWeekday.String(), maxWeekdayCommits))
	result.WriteString(fmt.Sprintf("Peak Hour Time Slot:        %s\n", cr.getTimeSlot(peakHour)))

	return result.String(), nil
}
// RenderFileStatistics renders detailed file and file type statistics (Requirement 2.6)
func (cr *ChartsRenderer) RenderFileStatistics(summary *models.StatsSummary, config models.RenderConfig) (string, error) {
	if summary == nil {
		return "", fmt.Errorf("summary cannot be nil")
	}

	var result strings.Builder

	result.WriteString("File Statistics\n")
	result.WriteString("===============\n\n")

	// Most frequently modified files
	if len(summary.TopFiles) > 0 {
		result.WriteString("Most Frequently Modified Files\n")
		result.WriteString("------------------------------\n")

		headers := []string{"Rank", "File", "Commits", "Insertions", "Deletions", "Net Changes", "Last Modified"}
		var rows [][]string

		for i, file := range summary.TopFiles {
			netChanges := file.Insertions - file.Deletions
			netChangeStr := fmt.Sprintf("%+d", netChanges)

			row := []string{
				fmt.Sprintf("#%d", i+1),
				file.Path,
				fmt.Sprintf("%d", file.Commits),
				fmt.Sprintf("%d", file.Insertions),
				fmt.Sprintf("%d", file.Deletions),
				netChangeStr,
				file.LastModified.Format("2006-01-02"),
			}
			rows = append(rows, row)
		}

		table, err := cr.RenderTable(headers, rows, config)
		if err == nil {
			result.WriteString(table)
		}
		result.WriteString("\n")

		// File modification frequency chart
		result.WriteString("File Modification Frequency\n")
		result.WriteString("---------------------------\n")

		fileData := make(map[string]int)
		for _, file := range summary.TopFiles {
			// Truncate long file names for better display
			fileName := file.Path
			if len(fileName) > 30 {
				fileName = "..." + fileName[len(fileName)-27:]
			}
			fileData[fileName] = file.Commits
		}

		fileChart, err := cr.RenderBarChart(fileData, "", config)
		if err == nil {
			result.WriteString(fileChart)
		}
		result.WriteString("\n")
	}

	// File type statistics
	if len(summary.TopFileTypes) > 0 {
		result.WriteString("File Type Statistics\n")
		result.WriteString("--------------------\n")

		headers := []string{"Rank", "Extension", "Files", "Commits", "Total Lines", "Avg Lines/File", "Commit %"}
		var rows [][]string

		totalCommits := 0
		for _, fileType := range summary.TopFileTypes {
			totalCommits += fileType.Commits
		}

		for i, fileType := range summary.TopFileTypes {
			avgLines := 0
			if fileType.Files > 0 {
				avgLines = fileType.Lines / fileType.Files
			}

			commitPercentage := 0.0
			if totalCommits > 0 {
				commitPercentage = float64(fileType.Commits) / float64(totalCommits) * 100
			}

			extension := fileType.Extension
			if extension == "" {
				extension = "(no extension)"
			}

			row := []string{
				fmt.Sprintf("#%d", i+1),
				extension,
				fmt.Sprintf("%d", fileType.Files),
				fmt.Sprintf("%d", fileType.Commits),
				fmt.Sprintf("%d", fileType.Lines),
				fmt.Sprintf("%d", avgLines),
				fmt.Sprintf("%.1f%%", commitPercentage),
			}
			rows = append(rows, row)
		}

		table, err := cr.RenderTable(headers, rows, config)
		if err == nil {
			result.WriteString(table)
		}
		result.WriteString("\n")

		// File type distribution charts
		result.WriteString("File Type Distribution\n")
		result.WriteString("----------------------\n")

		// Commits by file type
		result.WriteString("Commits by File Type:\n")
		typeCommitData := make(map[string]int)
		for _, fileType := range summary.TopFileTypes {
			extension := fileType.Extension
			if extension == "" {
				extension = "no-ext"
			}
			typeCommitData[extension] = fileType.Commits
		}

		typeCommitChart, err := cr.RenderBarChart(typeCommitData, "", config)
		if err == nil {
			result.WriteString(typeCommitChart)
		}
		result.WriteString("\n")

		// Lines by file type
		result.WriteString("Lines by File Type:\n")
		typeLinesData := make(map[string]int)
		for _, fileType := range summary.TopFileTypes {
			extension := fileType.Extension
			if extension == "" {
				extension = "no-ext"
			}
			typeLinesData[extension] = fileType.Lines
		}

		typeLinesChart, err := cr.RenderBarChart(typeLinesData, "", config)
		if err == nil {
			result.WriteString(typeLinesChart)
		}
	}

	return result.String(), nil
}
