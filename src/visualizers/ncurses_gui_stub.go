// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - NCurses GUI interface stub for testing

//go:build !gui
// +build !gui

package visualizers

import (
	"fmt"
	"git-stats/models"
	"time"
)

// ViewType represents different GUI views
type ViewType int

const (
	ContributionView ViewType = iota
	StatisticsView
	ContributorsView
	HealthView
)

// String returns the string representation of ViewType
func (vt ViewType) String() string {
	switch vt {
	case ContributionView:
		return "Contribution"
	case StatisticsView:
		return "Statistics"
	case ContributorsView:
		return "Contributors"
	case HealthView:
		return "Health"
	default:
		return "Unknown"
	}
}

// GUIState manages the state of the GUI interface
type GUIState struct {
	CurrentView     ViewType
	SelectedDate    time.Time
	ViewStartDate   time.Time
	ViewEndDate     time.Time
	SelectedCommits []models.Commit
	ShowHelp        bool
	StatusMessage   string
	Data            *models.AnalysisResult
}

// NewGUIState creates a new GUI state with default values
func NewGUIState(data *models.AnalysisResult) *GUIState {
	now := time.Now()
	yearAgo := now.AddDate(-1, 0, 0)

	return &GUIState{
		CurrentView:   ContributionView,
		SelectedDate:  now,
		ViewStartDate: yearAgo,
		ViewEndDate:   now,
		ShowHelp:      false,
		StatusMessage: "Ready",
		Data:          data,
	}
}

// SwitchView changes the current view
func (gs *GUIState) SwitchView(view ViewType) {
	gs.CurrentView = view
	gs.StatusMessage = fmt.Sprintf("Switched to %s view", view.String())
}

// SelectDate sets the selected date and updates related state
func (gs *GUIState) SelectDate(date time.Time) {
	gs.SelectedDate = date
	gs.StatusMessage = fmt.Sprintf("Selected date: %s", date.Format("2006-01-02"))

	// Update selected commits for the date if contribution graph data is available
	if gs.Data != nil && gs.Data.ContribGraph != nil {
		dateStr := date.Format("2006-01-02")
		if commits, exists := gs.Data.ContribGraph.DailyCommits[dateStr]; exists {
			gs.StatusMessage = fmt.Sprintf("Selected: %s (%d commits)", dateStr, commits)
		} else {
			gs.StatusMessage = fmt.Sprintf("Selected: %s (0 commits)", dateStr)
		}
	}
}

// NavigateMonth moves the view by the specified number of months
func (gs *GUIState) NavigateMonth(months int) {
	gs.ViewStartDate = gs.ViewStartDate.AddDate(0, months, 0)
	gs.ViewEndDate = gs.ViewEndDate.AddDate(0, months, 0)
	gs.StatusMessage = fmt.Sprintf("Viewing: %s to %s",
		gs.ViewStartDate.Format("Jan 2006"),
		gs.ViewEndDate.Format("Jan 2006"))
}

// NavigateYear moves the view by the specified number of years
func (gs *GUIState) NavigateYear(years int) {
	gs.ViewStartDate = gs.ViewStartDate.AddDate(years, 0, 0)
	gs.ViewEndDate = gs.ViewEndDate.AddDate(years, 0, 0)
	gs.StatusMessage = fmt.Sprintf("Viewing: %s to %s",
		gs.ViewStartDate.Format("Jan 2006"),
		gs.ViewEndDate.Format("Jan 2006"))
}

// GetCommitsForDate returns commits for a specific date
func (gs *GUIState) GetCommitsForDate(date time.Time) []models.Commit {
	if gs.Data == nil {
		return nil
	}

	// This would typically be populated by the analyzer
	// For now, return empty slice as the actual commit details
	// would be fetched from the repository when needed
	return gs.SelectedCommits
}

// UpdateSelectedCommits updates the selected commits for the current date
func (gs *GUIState) UpdateSelectedCommits(commits []models.Commit) {
	gs.SelectedCommits = commits
}

// ToggleHelp toggles the help display
func (gs *GUIState) ToggleHelp() {
	gs.ShowHelp = !gs.ShowHelp
	if gs.ShowHelp {
		gs.StatusMessage = "Help displayed"
	} else {
		gs.StatusMessage = "Help hidden"
	}
}

// GUIInterface implements the main GUI interface (stub version)
type GUIInterface struct {
	state *GUIState
}

// NewGUIInterface creates a new GUI interface
func NewGUIInterface() *GUIInterface {
	return &GUIInterface{}
}

// Initialize sets up the GUI interface
func (gui *GUIInterface) Initialize() error {
	return nil
}

// Run starts the GUI with the provided data (stub implementation)
func (gui *GUIInterface) Run(data *models.AnalysisResult) error {
	gui.state = NewGUIState(data)
	return fmt.Errorf("GUI mode requires building with -tags gui\n\nTo enable GUI mode:\n1. Run: go build -tags gui -o git-stats .\n2. Then use: ./git-stats -gui /path/to/repo\n\nAlternatively, use terminal mode:\n- git-stats -contrib /path/to/repo\n- git-stats -summary /path/to/repo")
}

// HandleInput processes input events (stub implementation)
func (gui *GUIInterface) HandleInput() error {
	return nil
}

// Render renders the GUI (stub implementation)
func (gui *GUIInterface) Render() error {
	return nil
}

// Cleanup performs cleanup operations
func (gui *GUIInterface) Cleanup() error {
	return nil
}
