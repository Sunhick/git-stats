// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Terminal UI components for interactive data exploration

package visualizers

import (
	"fmt"
	"git-stats/models"
	"git-stats/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Color constants for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

// TerminalUI provides terminal-based UI components
type TerminalUI struct {
	Width       int
	Height      int
	ColorScheme string
	Interactive bool
}

// NewTerminalUI creates a new terminal UI instance
func NewTerminalUI(config models.RenderConfig) *TerminalUI {
	return &TerminalUI{
		Width:       config.Width,
		Height:      config.Height,
		ColorScheme: config.ColorScheme,
		Interactive: config.Interactive,
	}
}

// ProgressIndicator provides visual progress feedback
type ProgressIndicator struct {
	*utils.ProgressTracker
	Style     ProgressStyle
	ShowStats bool
}

// ProgressStyle defines the visual style of progress indicators
type ProgressStyle int

const (
	ProgressStyleBar ProgressStyle = iota
	ProgressStyleSpinner
	ProgressStyleDots
	ProgressStylePercentage
)

// NewProgressIndicator creates a new progress indicator
func NewProgressIndicator(total int, message string, style ProgressStyle) *ProgressIndicator {
	return &ProgressIndicator{
		ProgressTracker: utils.NewProgressTracker(total, message),
		Style:           style,
		ShowStats:       true,
	}
}

// RenderProgress renders a progress indicator with the specified style
func (pi *ProgressIndicator) RenderProgress() string {
	current, total, percentage := pi.GetProgress()

	switch pi.Style {
	case ProgressStyleBar:
		return pi.renderProgressBar(current, total, percentage)
	case ProgressStyleSpinner:
		return pi.renderSpinner(current, total, percentage)
	case ProgressStyleDots:
		return pi.renderDots(current, total, percentage)
	case ProgressStylePercentage:
		return pi.renderPercentage(current, total, percentage)
	default:
		return pi.renderProgressBar(current, total, percentage)
	}
}

// renderProgressBar renders a traditional progress bar
func (pi *ProgressIndicator) renderProgressBar(current, total int, percentage float64) string {
	width := 40
	filled := int(percentage / 100.0 * float64(width))

	bar := ColorGreen + strings.Repeat("█", filled) + ColorReset +
		ColorDim + strings.Repeat("░", width-filled) + ColorReset

	stats := ""
	if pi.ShowStats {
		elapsed := time.Since(pi.StartTime)
		stats = fmt.Sprintf(" %s[%s]%s", ColorDim, formatDuration(elapsed), ColorReset)
	}

	return fmt.Sprintf("%s%s%s [%s] %d/%d (%.1f%%)%s",
		ColorBold, pi.Message, ColorReset, bar, current, total, percentage, stats)
}

// renderSpinner renders a spinning indicator
func (pi *ProgressIndicator) renderSpinner(current, total int, percentage float64) string {
	spinChars := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	spinIndex := int(time.Since(pi.StartTime).Milliseconds()/100) % len(spinChars)

	return fmt.Sprintf("%s%c%s %s %s(%d/%d - %.1f%%)%s",
		ColorCyan, spinChars[spinIndex], ColorReset,
		pi.Message, ColorDim, current, total, percentage, ColorReset)
}

// renderDots renders a dot-based progress indicator
func (pi *ProgressIndicator) renderDots(current, total int, percentage float64) string {
	dots := int(percentage / 10) // 10 dots max
	progress := ColorGreen + strings.Repeat("●", dots) + ColorReset +
		ColorDim + strings.Repeat("○", 10-dots) + ColorReset

	return fmt.Sprintf("%s%s%s %s %.1f%%",
		ColorBold, pi.Message, ColorReset, progress, percentage)
}

// renderPercentage renders a simple percentage indicator
func (pi *ProgressIndicator) renderPercentage(current, total int, percentage float64) string {
	color := ColorGreen
	if percentage < 50 {
		color = ColorYellow
	}
	if percentage < 25 {
		color = ColorRed
	}

	return fmt.Sprintf("%s%s%s %s%.1f%%%s (%d/%d)",
		ColorBold, pi.Message, ColorReset, color, percentage, ColorReset, current, total)
}

// InteractiveTable provides an interactive table for data exploration
type InteractiveTable struct {
	Headers     []string
	Rows        [][]string
	CurrentRow  int
	PageSize    int
	CurrentPage int
	Sortable    bool
	SortColumn  int
	SortAsc     bool
	Filterable  bool
	Filter      string
}

// NewInteractiveTable creates a new interactive table
func NewInteractiveTable(headers []string, rows [][]string) *InteractiveTable {
	return &InteractiveTable{
		Headers:     headers,
		Rows:        rows,
		CurrentRow:  0,
		PageSize:    10,
		CurrentPage: 0,
		Sortable:    true,
		SortColumn:  0,
		SortAsc:     true,
		Filterable:  true,
		Filter:      "",
	}
}

// RenderTable renders the interactive table
func (it *InteractiveTable) RenderTable() string {
	var result strings.Builder

	// Filter rows if filter is active
	filteredRows := it.getFilteredRows()

	// Sort rows if sorting is enabled
	if it.Sortable {
		it.sortRows(filteredRows)
	}

	// Calculate pagination
	totalPages := (len(filteredRows) + it.PageSize - 1) / it.PageSize
	startIdx := it.CurrentPage * it.PageSize
	endIdx := startIdx + it.PageSize
	if endIdx > len(filteredRows) {
		endIdx = len(filteredRows)
	}

	// Render header
	result.WriteString(it.renderTableHeader())
	result.WriteString("\n")

	// Render separator
	result.WriteString(it.renderTableSeparator())
	result.WriteString("\n")

	// Render rows
	for i := startIdx; i < endIdx; i++ {
		rowStyle := ""
		if i == it.CurrentRow {
			rowStyle = ColorBold + ColorCyan
		}
		result.WriteString(it.renderTableRow(filteredRows[i], rowStyle))
		result.WriteString("\n")
	}

	// Render footer with pagination info
	result.WriteString(it.renderTableFooter(it.CurrentPage+1, totalPages, len(filteredRows)))

	return result.String()
}

// renderTableHeader renders the table header with sorting indicators
func (it *InteractiveTable) renderTableHeader() string {
	var header strings.Builder
	header.WriteString(ColorBold)

	for i, h := range it.Headers {
		if i > 0 {
			header.WriteString(" │ ")
		}

		// Add sorting indicator
		sortIndicator := ""
		if it.Sortable && i == it.SortColumn {
			if it.SortAsc {
				sortIndicator = " ↑"
			} else {
				sortIndicator = " ↓"
			}
		}

		header.WriteString(fmt.Sprintf("%-15s%s", h, sortIndicator))
	}

	header.WriteString(ColorReset)
	return header.String()
}

// renderTableSeparator renders the table separator line
func (it *InteractiveTable) renderTableSeparator() string {
	width := len(it.Headers)*17 + (len(it.Headers)-1)*3 // Approximate width
	return ColorDim + strings.Repeat("─", width) + ColorReset
}

// renderTableRow renders a single table row
func (it *InteractiveTable) renderTableRow(row []string, style string) string {
	var result strings.Builder
	result.WriteString(style)

	for i, cell := range row {
		if i > 0 {
			result.WriteString(" │ ")
		}
		result.WriteString(fmt.Sprintf("%-15s", TruncateString(cell, 15)))
	}

	result.WriteString(ColorReset)
	return result.String()
}

// renderTableFooter renders pagination and filter information
func (it *InteractiveTable) renderTableFooter(currentPage, totalPages, totalRows int) string {
	var footer strings.Builder

	// Pagination info
	footer.WriteString(fmt.Sprintf("%sPage %d of %d (%d rows)%s",
		ColorDim, currentPage, totalPages, totalRows, ColorReset))

	// Filter info
	if it.Filter != "" {
		footer.WriteString(fmt.Sprintf(" %s[Filter: %s]%s", ColorYellow, it.Filter, ColorReset))
	}

	// Navigation help
	if it.Sortable || it.Filterable {
		footer.WriteString("\n")
		footer.WriteString(ColorDim + "Navigation: ↑↓ Select, ←→ Page, s Sort, f Filter, q Quit" + ColorReset)
	}

	return footer.String()
}

// getFilteredRows returns rows that match the current filter
func (it *InteractiveTable) getFilteredRows() [][]string {
	if it.Filter == "" {
		return it.Rows
	}

	var filtered [][]string
	filterLower := strings.ToLower(it.Filter)

	for _, row := range it.Rows {
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell), filterLower) {
				filtered = append(filtered, row)
				break
			}
		}
	}

	return filtered
}

// sortRows sorts the rows based on the current sort column and direction
func (it *InteractiveTable) sortRows(rows [][]string) {
	if it.SortColumn >= len(it.Headers) {
		return
	}

	sort.Slice(rows, func(i, j int) bool {
		if it.SortColumn >= len(rows[i]) || it.SortColumn >= len(rows[j]) {
			return false
		}

		a := rows[i][it.SortColumn]
		b := rows[j][it.SortColumn]

		// Try to parse as numbers first
		if numA, errA := strconv.ParseFloat(a, 64); errA == nil {
			if numB, errB := strconv.ParseFloat(b, 64); errB == nil {
				if it.SortAsc {
					return numA < numB
				}
				return numA > numB
			}
		}

		// Fall back to string comparison
		if it.SortAsc {
			return strings.ToLower(a) < strings.ToLower(b)
		}
		return strings.ToLower(a) > strings.ToLower(b)
	})
}

// ColoredBarChart renders a colored bar chart
type ColoredBarChart struct {
	Title     string
	Data      map[string]int
	MaxWidth  int
	ShowValue bool
	Colors    []string
}

// NewColoredBarChart creates a new colored bar chart
func NewColoredBarChart(title string, data map[string]int, maxWidth int) *ColoredBarChart {
	return &ColoredBarChart{
		Title:     title,
		Data:      data,
		MaxWidth:  maxWidth,
		ShowValue: true,
		Colors:    []string{ColorBlue, ColorGreen, ColorYellow, ColorPurple, ColorCyan},
	}
}

// RenderChart renders the colored bar chart
func (cbc *ColoredBarChart) RenderChart() string {
	var result strings.Builder

	// Title
	result.WriteString(fmt.Sprintf("%s%s%s\n", ColorBold, cbc.Title, ColorReset))
	result.WriteString(strings.Repeat("─", len(cbc.Title)) + "\n")

	// Find max value for scaling
	maxValue := 0
	for _, value := range cbc.Data {
		if value > maxValue {
			maxValue = value
		}
	}

	if maxValue == 0 {
		result.WriteString(ColorDim + "No data available" + ColorReset + "\n")
		return result.String()
	}

	// Sort keys for consistent display
	var keys []string
	for key := range cbc.Data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Render bars
	colorIndex := 0
	for _, key := range keys {
		value := cbc.Data[key]
		barWidth := int(float64(value) / float64(maxValue) * float64(cbc.MaxWidth))

		color := cbc.Colors[colorIndex%len(cbc.Colors)]
		colorIndex++

		bar := color + strings.Repeat("█", barWidth) + ColorReset

		valueStr := ""
		if cbc.ShowValue {
			valueStr = fmt.Sprintf(" %s(%d)%s", ColorDim, value, ColorReset)
		}

		result.WriteString(fmt.Sprintf("%-15s %s%s\n", TruncateString(key, 15), bar, valueStr))
	}

	return result.String()
}

// StatusLine provides a colored status line for displaying information
type StatusLine struct {
	Message string
	Type    StatusType
	Width   int
}

// StatusType defines the type of status message
type StatusType int

const (
	StatusInfo StatusType = iota
	StatusSuccess
	StatusWarning
	StatusError
)

// NewStatusLine creates a new status line
func NewStatusLine(message string, statusType StatusType, width int) *StatusLine {
	return &StatusLine{
		Message: message,
		Type:    statusType,
		Width:   width,
	}
}

// RenderStatus renders the status line with appropriate colors
func (sl *StatusLine) RenderStatus() string {
	var color string
	var prefix string

	switch sl.Type {
	case StatusInfo:
		color = ColorBlue
		prefix = "ℹ"
	case StatusSuccess:
		color = ColorGreen
		prefix = "✓"
	case StatusWarning:
		color = ColorYellow
		prefix = "⚠"
	case StatusError:
		color = ColorRed
		prefix = "✗"
	}

	message := fmt.Sprintf("%s %s", prefix, sl.Message)
	padding := sl.Width - len(message)
	if padding < 0 {
		padding = 0
		message = TruncateString(message, sl.Width-3) + "..."
	}

	return fmt.Sprintf("%s%s%s%s%s",
		color, message, strings.Repeat(" ", padding), ColorReset, "\n")
}

// InteractiveMenu provides a navigable menu interface
type InteractiveMenu struct {
	Title       string
	Options     []MenuOption
	CurrentItem int
	ShowHelp    bool
}

// MenuOption represents a menu item
type MenuOption struct {
	Label       string
	Description string
	Action      func() error
	Enabled     bool
}

// NewInteractiveMenu creates a new interactive menu
func NewInteractiveMenu(title string, options []MenuOption) *InteractiveMenu {
	return &InteractiveMenu{
		Title:       title,
		Options:     options,
		CurrentItem: 0,
		ShowHelp:    true,
	}
}

// RenderMenu renders the interactive menu
func (im *InteractiveMenu) RenderMenu() string {
	var result strings.Builder

	// Title
	result.WriteString(fmt.Sprintf("%s%s%s\n", ColorBold, im.Title, ColorReset))
	result.WriteString(strings.Repeat("═", len(im.Title)) + "\n\n")

	// Menu options
	for i, option := range im.Options {
		prefix := "  "
		color := ""

		if i == im.CurrentItem {
			prefix = "▶ "
			color = ColorBold + ColorCyan
		}

		if !option.Enabled {
			color = ColorDim
		}

		result.WriteString(fmt.Sprintf("%s%s%s%s%s\n",
			color, prefix, option.Label, ColorReset,
			func() string {
				if option.Description != "" {
					return fmt.Sprintf(" %s- %s%s", ColorDim, option.Description, ColorReset)
				}
				return ""
			}()))
	}

	// Help text
	if im.ShowHelp {
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("%s↑↓ Navigate, Enter Select, q Quit%s\n",
			ColorDim, ColorReset))
	}

	return result.String()
}

// Helper functions

// TruncateString truncates a string to the specified length
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) - 60*minutes
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	} else {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) - 60*hours
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
}

// GetTerminalWidth attempts to get the terminal width
func GetTerminalWidth() int {
	// Default width if unable to determine
	return 80
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// MoveCursor moves the cursor to the specified position
func MoveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

// HideCursor hides the terminal cursor
func HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor
func ShowCursor() {
	fmt.Print("\033[?25h")
}
