// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - NCurses GUI interface for interactive visualization

//go:build gui
// +build gui

package visualizers

import (
	"fmt"
	"git-stats/models"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

// ToggleHelp toggles the help display
func (gs *GUIState) ToggleHelp() {
	gs.ShowHelp = !gs.ShowHelp
	if gs.ShowHelp {
		gs.StatusMessage = "Help displayed"
	} else {
		gs.StatusMessage = "Help hidden"
	}
}

// KeyCommand represents a keyboard command
type KeyCommand struct {
	Key         tcell.Key
	Rune        rune
	Description string
	Action      func(*GUIState) error
}

// ContributionGraphWidget handles the contribution graph display
type ContributionGraphWidget struct {
	*tview.Box
	Data         *models.ContributionGraph
	State        *GUIState
	SelectedDay  time.Time
	ViewOffset   int
	CellWidth    int
	CellHeight   int
}

// NewContributionGraphWidget creates a new contribution graph widget
func NewContributionGraphWidget(data *models.ContributionGraph, state *GUIState) *ContributionGraphWidget {
	widget := &ContributionGraphWidget{
		Box:         tview.NewBox(),
		Data:        data,
		State:       state,
		SelectedDay: time.Now(),
		ViewOffset:  0,
		CellWidth:   2,
		CellHeight:  1,
	}

	widget.SetBorder(true).SetTitle("Contribution Graph")
	return widget
}

// Draw renders the contribution graph widget
func (cgw *ContributionGraphWidget) Draw(screen tcell.Screen) {
	cgw.Box.DrawForSubclass(screen, cgw)
	x, y, width, height := cgw.GetInnerRect()

	if cgw.Data == nil {
		tview.Print(screen, "No contribution data available", x, y, width, tcell.AlignLeft, tcell.ColorWhite)
		return
	}

	// Draw month labels
	cgw.drawMonthLabels(screen, x, y, width)

	// Draw day indicators
	cgw.drawDayIndicators(screen, x, y+2, height-2)

	// Draw contribution cells
	cgw.drawContributionCells(screen, x+4, y+3, width-4, height-3)
}

// drawMonthLabels draws the month labels at the top
func (cgw *ContributionGraphWidget) drawMonthLabels(screen tcell.Screen, x, y, width int) {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
					   "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	startMonth := cgw.State.ViewStartDate.Month()
	for i := 0; i < 12 && i*8 < width-4; i++ {
		monthIndex := (int(startMonth) - 1 + i) % 12
		tview.Print(screen, months[monthIndex], x+4+i*8, y, 3, tcell.AlignLeft, tcell.ColorYellow)
	}
}

// drawDayIndicators draws the day-of-week indicators
func (cgw *ContributionGraphWidget) drawDayIndicators(screen tcell.Screen, x, y, height int) {
	days := []string{"S", "M", "T", "W", "T", "F", "S"}
	for i, day := range days {
		if i < height {
			tview.Print(screen, day, x, y+i, 1, tcell.AlignLeft, tcell.ColorCyan)
		}
	}
}

// drawContributionCells draws the actual contribution cells
func (cgw *ContributionGraphWidget) drawContributionCells(screen tcell.Screen, x, y, width, height int) {
	if height < 7 {
		return // Not enough space for a week
	}

	// Calculate the starting date (should be a Sunday)
	startDate := cgw.State.ViewStartDate
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	// Draw cells for each week
	currentDate := startDate
	weekOffset := 0

	for weekOffset*3 < width && currentDate.Before(cgw.State.ViewEndDate.AddDate(0, 0, 7)) {
		for day := 0; day < 7 && day < height; day++ {
			cellX := x + weekOffset*3
			cellY := y + day

			if cellX >= x+width-2 {
				break
			}

			// Get commit count for this date
			dateStr := currentDate.Format("2006-01-02")
			commits := 0
			if cgw.Data.DailyCommits != nil {
				commits = cgw.Data.DailyCommits[dateStr]
			}

			// Determine cell style based on commit count
			style := cgw.getCellStyle(commits)
			char := cgw.getCellChar(commits)

			// Highlight selected date
			if currentDate.Format("2006-01-02") == cgw.State.SelectedDate.Format("2006-01-02") {
				style = style.Background(tcell.ColorBlue)
			}

			screen.SetContent(cellX, cellY, char, nil, style)
			screen.SetContent(cellX+1, cellY, char, nil, style)

			currentDate = currentDate.AddDate(0, 0, 1)
		}
		weekOffset++
	}
}

// getCellStyle returns the appropriate style for a cell based on commit count with enhanced colors
func (cgw *ContributionGraphWidget) getCellStyle(commits int) tcell.Style {
	base := tcell.StyleDefault

	// Calculate activity level based on max commits for better scaling
	maxCommits := 1
	if cgw.Data != nil {
		maxCommits = cgw.Data.MaxCommits
		if maxCommits == 0 {
			maxCommits = 1
		}
	}

	switch {
	case commits == 0:
		// No activity - dark gray
		return base.Foreground(tcell.ColorDarkGray).Background(tcell.ColorBlack)
	case commits <= maxCommits/4:
		// Low activity - light green
		return base.Foreground(tcell.ColorLightGreen).Background(tcell.ColorDarkGreen)
	case commits <= maxCommits/2:
		// Medium activity - medium green
		return base.Foreground(tcell.ColorGreen).Background(tcell.ColorDarkGreen)
	case commits <= maxCommits*3/4:
		// High activity - bright green
		return base.Foreground(tcell.ColorLime).Background(tcell.ColorGreen)
	default:
		// Very high activity - yellow-green
		return base.Foreground(tcell.ColorYellow).Background(tcell.ColorGreen)
	}
}

// getCellChar returns the appropriate character for a cell based on commit count with better scaling
func (cgw *ContributionGraphWidget) getCellChar(commits int) rune {
	// Calculate activity level based on max commits for better scaling
	maxCommits := 1
	if cgw.Data != nil {
		maxCommits = cgw.Data.MaxCommits
		if maxCommits == 0 {
			maxCommits = 1
		}
	}

	switch {
	case commits == 0:
		return '░' // Light shade for no commits
	case commits <= maxCommits/4:
		return '▒' // Medium shade for low activity
	case commits <= maxCommits/2:
		return '▓' // Dark shade for medium activity
	case commits <= maxCommits*3/4:
		return '█' // Full block for high activity
	default:
		return '█' // Full block for very high activity (different color)
	}
}

// HandleInput processes keyboard input for the contribution graph
func (cgw *ContributionGraphWidget) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyLeft:
		cgw.State.SelectDate(cgw.State.SelectedDate.AddDate(0, 0, -1))
		return nil
	case tcell.KeyRight:
		cgw.State.SelectDate(cgw.State.SelectedDate.AddDate(0, 0, 1))
		return nil
	case tcell.KeyUp:
		cgw.State.SelectDate(cgw.State.SelectedDate.AddDate(0, 0, -7))
		return nil
	case tcell.KeyDown:
		cgw.State.SelectDate(cgw.State.SelectedDate.AddDate(0, 0, 7))
		return nil
	}

	switch event.Rune() {
	case 'h', 'H':
		cgw.State.NavigateMonth(-1)
		return nil
	case 'l', 'L':
		cgw.State.NavigateMonth(1)
		return nil
	}

	return event
}

// DetailPanelWidget displays detailed information with enhanced functionality
type DetailPanelWidget struct {
	*tview.TextView
	State       *GUIState
	Title       string
	MaxLines    int
	ScrollPos   int
	Content     []string
	ShowDetails bool
	SelectedCommitIndex int
}

// NewDetailPanelWidget creates a new detail panel widget with enhanced functionality
func NewDetailPanelWidget(state *GUIState, title string) *DetailPanelWidget {
	widget := &DetailPanelWidget{
		TextView:            tview.NewTextView(),
		State:               state,
		Title:               title,
		MaxLines:            20,
		Content:             make([]string, 0),
		ShowDetails:         true,
		SelectedCommitIndex: 0,
	}

	widget.SetBorder(true).SetTitle(title)
	widget.SetDynamicColors(true)
	widget.SetScrollable(true)
	widget.SetWrap(true)
	widget.SetWordWrap(true)

	return widget
}

// UpdateContent updates the content of the detail panel
func (dpw *DetailPanelWidget) UpdateContent() {
	dpw.Clear()

	if dpw.State.Data == nil {
		dpw.SetText("No data available")
		return
	}

	var content strings.Builder

	switch dpw.State.CurrentView {
	case ContributionView:
		dpw.updateContributionDetails(&content)
	case StatisticsView:
		dpw.updateStatisticsDetails(&content)
	case ContributorsView:
		dpw.updateContributorsDetails(&content)
	case HealthView:
		dpw.updateHealthDetails(&content)
	}

	dpw.SetText(content.String())
}

// updateContributionDetails updates content for contribution view with enhanced commit details
func (dpw *DetailPanelWidget) updateContributionDetails(content *strings.Builder) {
	selectedDate := dpw.State.SelectedDate.Format("2006-01-02")
	commits := 0

	if dpw.State.Data.ContribGraph != nil && dpw.State.Data.ContribGraph.DailyCommits != nil {
		commits = dpw.State.Data.ContribGraph.DailyCommits[selectedDate]
	}

	content.WriteString(fmt.Sprintf("[yellow]Selected Date:[white] %s\n", selectedDate))
	content.WriteString(fmt.Sprintf("[yellow]Commits:[white] %d\n\n", commits))

	if dpw.State.Data.ContribGraph != nil {
		content.WriteString(fmt.Sprintf("[yellow]Total Contributions:[white] %d\n", dpw.State.Data.ContribGraph.TotalCommits))
		content.WriteString(fmt.Sprintf("[yellow]Max Daily Commits:[white] %d\n", dpw.State.Data.ContribGraph.MaxCommits))
		content.WriteString(fmt.Sprintf("[yellow]Period:[white] %s to %s\n\n",
			dpw.State.Data.ContribGraph.StartDate.Format("2006-01-02"),
			dpw.State.Data.ContribGraph.EndDate.Format("2006-01-02")))
	}

	// Show detailed commit information if available
	if dpw.ShowDetails && len(dpw.State.SelectedCommits) > 0 {
		content.WriteString("[cyan]Commit Details:[white]\n")
		for i, commit := range dpw.State.SelectedCommits {
			if i >= dpw.MaxLines-10 { // Leave space for other info
				content.WriteString(fmt.Sprintf("... and %d more commits\n", len(dpw.State.SelectedCommits)-i))
				break
			}

			// Highlight selected commit
			prefix := "  "
			if i == dpw.SelectedCommitIndex {
				prefix = "[blue]>[white] "
			}

			content.WriteString(fmt.Sprintf("%s[green]%s[white]\n", prefix, commit.Hash[:8]))
			content.WriteString(fmt.Sprintf("%s  %s\n", prefix, truncateString(commit.Message, 40)))
			content.WriteString(fmt.Sprintf("%s  [gray]%s - %s[white]\n",
				prefix, commit.Author.Name, commit.AuthorDate.Format("15:04")))
			content.WriteString(fmt.Sprintf("%s  [gray]+%d/-%d files: %d[white]\n\n",
				prefix, commit.Stats.Insertions, commit.Stats.Deletions, commit.Stats.FilesChanged))
		}
	}
}

// updateStatisticsDetails updates content for statistics view
func (dpw *DetailPanelWidget) updateStatisticsDetails(content *strings.Builder) {
	if dpw.State.Data.Summary == nil {
		content.WriteString("No statistics available")
		return
	}

	summary := dpw.State.Data.Summary
	content.WriteString(fmt.Sprintf("[yellow]Total Commits:[white] %d\n", summary.TotalCommits))
	content.WriteString(fmt.Sprintf("[yellow]Lines Added:[white] %d\n", summary.TotalInsertions))
	content.WriteString(fmt.Sprintf("[yellow]Lines Deleted:[white] %d\n", summary.TotalDeletions))
	content.WriteString(fmt.Sprintf("[yellow]Files Changed:[white] %d\n", summary.FilesChanged))
	content.WriteString(fmt.Sprintf("[yellow]Active Days:[white] %d\n", summary.ActiveDays))
	content.WriteString(fmt.Sprintf("[yellow]Avg Commits/Day:[white] %.2f\n", summary.AvgCommitsPerDay))
}

// updateContributorsDetails updates content for contributors view
func (dpw *DetailPanelWidget) updateContributorsDetails(content *strings.Builder) {
	if len(dpw.State.Data.Contributors) == 0 {
		content.WriteString("No contributors data available")
		return
	}

	content.WriteString(fmt.Sprintf("[yellow]Total Contributors:[white] %d\n\n", len(dpw.State.Data.Contributors)))

	for i, contributor := range dpw.State.Data.Contributors {
		if i >= 10 { // Limit to top 10 for display
			break
		}
		content.WriteString(fmt.Sprintf("[cyan]%s[white]\n", contributor.Name))
		content.WriteString(fmt.Sprintf("  Commits: %d\n", contributor.TotalCommits))
		content.WriteString(fmt.Sprintf("  Lines: +%d/-%d\n", contributor.TotalInsertions, contributor.TotalDeletions))
		content.WriteString("\n")
	}
}

// updateHealthDetails updates content for health view
func (dpw *DetailPanelWidget) updateHealthDetails(content *strings.Builder) {
	if dpw.State.Data.HealthMetrics == nil {
		content.WriteString("No health metrics available")
		return
	}

	health := dpw.State.Data.HealthMetrics
	content.WriteString(fmt.Sprintf("[yellow]Repository Age:[white] %s\n", health.RepositoryAge.String()))
	content.WriteString(fmt.Sprintf("[yellow]Commit Frequency:[white] %.2f/day\n", health.CommitFrequency))
	content.WriteString(fmt.Sprintf("[yellow]Total Contributors:[white] %d\n", health.ContributorCount))
	content.WriteString(fmt.Sprintf("[yellow]Active Contributors:[white] %d\n", health.ActiveContributors))
	content.WriteString(fmt.Sprintf("[yellow]Branch Count:[white] %d\n", health.BranchCount))
	content.WriteString(fmt.Sprintf("[yellow]Activity Trend:[white] %s\n", health.ActivityTrend))
}

// StatusBarWidget displays status information and keyboard shortcuts with enhanced functionality
type StatusBarWidget struct {
	*tview.TextView
	State         *GUIState
	Commands      []KeyCommand
	ShowShortcuts bool
	HelpText      string
}

// NewStatusBarWidget creates a new status bar widget with enhanced functionality
func NewStatusBarWidget(state *GUIState) *StatusBarWidget {
	widget := &StatusBarWidget{
		TextView:      tview.NewTextView(),
		State:         state,
		ShowShortcuts: true,
		HelpText:      "Press ? for help",
	}

	widget.SetDynamicColors(true)
	widget.SetTextAlign(tview.AlignLeft)

	// Define comprehensive key commands
	widget.Commands = []KeyCommand{
		{Key: tcell.KeyRune, Rune: 'c', Description: "[C]ontrib"},
		{Key: tcell.KeyRune, Rune: 's', Description: "[S]tats"},
		{Key: tcell.KeyRune, Rune: 't', Description: "[T]eam"},
		{Key: tcell.KeyRune, Rune: 'H', Description: "[H]ealth"},
		{Key: tcell.KeyLeft, Description: "← Day"},
		{Key: tcell.KeyRight, Description: "→ Day"},
		{Key: tcell.KeyUp, Description: "↑ Week"},
		{Key: tcell.KeyDown, Description: "↓ Week"},
		{Key: tcell.KeyRune, Rune: 'h', Description: "h/l Month"},
		{Key: tcell.KeyRune, Rune: 'd', Description: "[D]etails"},
		{Key: tcell.KeyRune, Rune: 'j', Description: "j/k Scroll"},
		{Key: tcell.KeyRune, Rune: 'q', Description: "[Q]uit"},
		{Key: tcell.KeyRune, Rune: '?', Description: "[?] Help"},
	}

	return widget
}

// GetText returns the text content (for testing)
func (dpw *DetailPanelWidget) GetText(stripTags bool) string {
	return dpw.TextView.GetText(stripTags)
}

// GetText returns the text content (for testing)
func (sbw *StatusBarWidget) GetText(stripTags bool) string {
	return sbw.TextView.GetText(stripTags)
}

// UpdateStatus updates the status bar content with enhanced information
func (sbw *StatusBarWidget) UpdateStatus() {
	var content strings.Builder

	// Add current status message
	content.WriteString(fmt.Sprintf("[yellow]%s[white] | ", sbw.State.StatusMessage))

	// Add current view
	content.WriteString(fmt.Sprintf("View: [cyan]%s[white]", sbw.State.CurrentView.String()))

	// Add view-specific information
	switch sbw.State.CurrentView {
	case ContributionView:
		if sbw.State.Data != nil && sbw.State.Data.ContribGraph != nil {
			selectedDate := sbw.State.SelectedDate.Format("2006-01-02")
			commits := sbw.State.Data.ContribGraph.DailyCommits[selectedDate]
			content.WriteString(fmt.Sprintf(" | [green]%s: %d commits[white]", selectedDate, commits))
		}
	case ContributorsView:
		if sbw.State.Data != nil {
			content.WriteString(fmt.Sprintf(" | [green]%d contributors[white]", len(sbw.State.Data.Contributors)))
		}
	case StatisticsView:
		if sbw.State.Data != nil && sbw.State.Data.Summary != nil {
			content.WriteString(fmt.Sprintf(" | [green]%d total commits[white]", sbw.State.Data.Summary.TotalCommits))
		}
	case HealthView:
		if sbw.State.Data != nil && sbw.State.Data.HealthMetrics != nil {
			content.WriteString(fmt.Sprintf(" | [green]%s trend[white]", sbw.State.Data.HealthMetrics.ActivityTrend))
		}
	}

	// Add keyboard shortcuts if enabled
	if sbw.ShowShortcuts {
		content.WriteString(" | ")
		relevantCommands := sbw.getRelevantCommands()
		for i, cmd := range relevantCommands {
			if i > 0 {
				content.WriteString(" ")
			}
			content.WriteString(fmt.Sprintf("[gray]%s[white]", cmd.Description))
		}
	} else {
		content.WriteString(fmt.Sprintf(" | [gray]%s[white]", sbw.HelpText))
	}

	sbw.SetText(content.String())
}

// getRelevantCommands returns commands relevant to the current view
func (sbw *StatusBarWidget) getRelevantCommands() []KeyCommand {
	baseCommands := []KeyCommand{
		{Key: tcell.KeyRune, Rune: 'c', Description: "[C]ontrib"},
		{Key: tcell.KeyRune, Rune: 's', Description: "[S]tats"},
		{Key: tcell.KeyRune, Rune: 't', Description: "[T]eam"},
		{Key: tcell.KeyRune, Rune: 'H', Description: "[H]ealth"},
	}

	switch sbw.State.CurrentView {
	case ContributionView:
		return append(baseCommands, []KeyCommand{
			{Key: tcell.KeyLeft, Description: "←→ Days"},
			{Key: tcell.KeyUp, Description: "↑↓ Weeks"},
			{Key: tcell.KeyRune, Rune: 'h', Description: "h/l Months"},
			{Key: tcell.KeyRune, Rune: 'd', Description: "[D]etails"},
		}...)
	case StatisticsView, ContributorsView, HealthView:
		return append(baseCommands, []KeyCommand{
			{Key: tcell.KeyRune, Rune: 'j', Description: "j/k Scroll"},
		}...)
	}

	return append(baseCommands, []KeyCommand{
		{Key: tcell.KeyRune, Rune: 'q', Description: "[Q]uit"},
		{Key: tcell.KeyRune, Rune: '?', Description: "[?] Help"},
	}...)
}

// ToggleShortcuts toggles the display of keyboard shortcuts
func (sbw *StatusBarWidget) ToggleShortcuts() {
	sbw.ShowShortcuts = !sbw.ShowShortcuts
	sbw.UpdateStatus()
}

// GetRelevantCommands returns commands relevant to the current view (public for testing)
func (sbw *StatusBarWidget) GetRelevantCommands() []KeyCommand {
	return sbw.getRelevantCommands()
}

// GUIInterface implements the main GUI interface
type GUIInterface struct {
	app              *tview.Application
	state            *GUIState
	layout           *tview.Flex
	contributionGraph *ContributionGraphWidget
	detailPanel      *DetailPanelWidget
	statusBar        *StatusBarWidget
	helpModal        *tview.Modal
}

// NewGUIInterface creates a new GUI interface
func NewGUIInterface() *GUIInterface {
	return &GUIInterface{
		app: tview.NewApplication(),
	}
}

// Initialize sets up the GUI interface
func (gui *GUIInterface) Initialize() error {
	// This will be called with data in Run method
	return nil
}

// Run starts the GUI with the provided data
func (gui *GUIInterface) Run(data *models.AnalysisResult) error {
	// Initialize state
	gui.state = NewGUIState(data)

	// Create widgets
	gui.contributionGraph = NewContributionGraphWidget(data.ContribGraph, gui.state)
	gui.detailPanel = NewDetailPanelWidget(gui.state, "Details")
	gui.statusBar = NewStatusBarWidget(gui.state)

	// Create help modal
	gui.helpModal = tview.NewModal().
		SetText("Git Stats - Keyboard Shortcuts\n\n" +
			"Navigation:\n" +
			"  ←→ : Navigate days\n" +
			"  ↑↓ : Navigate weeks\n" +
			"  h/l : Navigate months\n\n" +
			"Views:\n" +
			"  c : Contribution view\n" +
			"  s : Statistics view\n" +
			"  t : Team/Contributors view\n" +
			"  H : Health metrics view\n\n" +
			"Other:\n" +
			"  ? : Toggle this help\n" +
			"  q/ESC : Quit").
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			gui.state.ToggleHelp()
			gui.app.SetRoot(gui.layout, true)
		})

	// Create main layout
	gui.layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(gui.contributionGraph, 0, 2, true).
			AddItem(gui.detailPanel, 0, 1, false), 0, 1, true).
		AddItem(gui.statusBar, 1, 0, false)

	// Set up input handling
	gui.app.SetInputCapture(gui.handleGlobalInput)

	// Update initial content
	gui.updateDisplay()

	// Set root and run
	gui.app.SetRoot(gui.layout, true)
	return gui.app.Run()
}

// handleGlobalInput handles global keyboard input with enhanced navigation
func (gui *GUIInterface) handleGlobalInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		gui.app.Stop()
		return nil
	case tcell.KeyCtrlC:
		gui.app.Stop()
		return nil
	}

	switch event.Rune() {
	case 'q', 'Q':
		gui.app.Stop()
		return nil
	case '?':
		gui.state.ToggleHelp()
		if gui.state.ShowHelp {
			gui.app.SetRoot(gui.helpModal, true)
		} else {
			gui.app.SetRoot(gui.layout, true)
		}
		return nil
	case 'c', 'C':
		gui.state.SwitchView(ContributionView)
		gui.updateDisplay()
		return nil
	case 's', 'S':
		gui.state.SwitchView(StatisticsView)
		gui.updateDisplay()
		return nil
	case 't', 'T':
		gui.state.SwitchView(ContributorsView)
		gui.updateDisplay()
		return nil
	case 'H':
		gui.state.SwitchView(HealthView)
		gui.updateDisplay()
		return nil
	case 'd', 'D':
		// Toggle detail panel visibility
		if gui.detailPanel != nil {
			gui.detailPanel.ShowDetails = !gui.detailPanel.ShowDetails
			gui.updateDisplay()
		}
		return nil
	case 'j':
		// Scroll down in detail panel
		if gui.detailPanel != nil {
			gui.scrollDetailPanel(1)
		}
		return nil
	case 'k':
		// Scroll up in detail panel
		if gui.detailPanel != nil {
			gui.scrollDetailPanel(-1)
		}
		return nil
	case 'r', 'R':
		// Refresh display
		gui.updateDisplay()
		gui.state.StatusMessage = "Display refreshed"
		return nil
	}

	// Handle view-specific input
	switch gui.state.CurrentView {
	case ContributionView:
		if gui.contributionGraph != nil {
			return gui.contributionGraph.HandleInput(event)
		}
	case StatisticsView, ContributorsView, HealthView:
		// Handle scrolling for text-based views
		return gui.handleTextViewInput(event)
	}

	return event
}

// handleTextViewInput handles input for text-based views
func (gui *GUIInterface) handleTextViewInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		if gui.detailPanel != nil {
			gui.scrollDetailPanel(-1)
		}
		return nil
	case tcell.KeyDown:
		if gui.detailPanel != nil {
			gui.scrollDetailPanel(1)
		}
		return nil
	case tcell.KeyPageUp:
		if gui.detailPanel != nil {
			gui.scrollDetailPanel(-5)
		}
		return nil
	case tcell.KeyPageDown:
		if gui.detailPanel != nil {
			gui.scrollDetailPanel(5)
		}
		return nil
	case tcell.KeyHome:
		if gui.detailPanel != nil {
			gui.detailPanel.ScrollPos = 0
			gui.updateDisplay()
		}
		return nil
	case tcell.KeyEnd:
		if gui.detailPanel != nil {
			gui.detailPanel.ScrollPos = len(gui.detailPanel.Content)
			gui.updateDisplay()
		}
		return nil
	}
	return event
}

// scrollDetailPanel scrolls the detail panel by the specified amount
func (gui *GUIInterface) scrollDetailPanel(delta int) {
	if gui.detailPanel == nil {
		return
	}

	newPos := gui.detailPanel.ScrollPos + delta
	maxPos := len(gui.detailPanel.Content) - gui.detailPanel.MaxLines
	if maxPos < 0 {
		maxPos = 0
	}

	if newPos < 0 {
		newPos = 0
	} else if newPos > maxPos {
		newPos = maxPos
	}

	gui.detailPanel.ScrollPos = newPos
	gui.updateDisplay()
}

// updateDisplay updates all display components
func (gui *GUIInterface) updateDisplay() {
	if gui.detailPanel != nil {
		gui.detailPanel.UpdateContent()
	}
	if gui.statusBar != nil {
		gui.statusBar.UpdateStatus()
	}
	if gui.app != nil {
		gui.app.Draw()
	}
}

// HandleInput processes input events (part of GUIVisualizer interface)
func (gui *GUIInterface) HandleInput() error {
	// Input handling is managed by tview application
	return nil
}

// Render renders the GUI (part of GUIVisualizer interface)
func (gui *GUIInterface) Render() error {
	gui.updateDisplay()
	return nil
}

// Cleanup performs cleanup operations
func (gui *GUIInterface) Cleanup() error {
	if gui.app != nil {
		gui.app.Stop()
	}
	return nil
}

// truncateString truncates a string to the specified length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// HandleInput processes keyboard input for the detail panel
func (dpw *DetailPanelWidget) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		if dpw.SelectedCommitIndex > 0 {
			dpw.SelectedCommitIndex--
			dpw.UpdateContent()
		}
		return nil
	case tcell.KeyDown:
		if dpw.SelectedCommitIndex < len(dpw.State.SelectedCommits)-1 {
			dpw.SelectedCommitIndex++
			dpw.UpdateContent()
		}
		return nil
	case tcell.KeyPageUp:
		dpw.SelectedCommitIndex = 0
		dpw.UpdateContent()
		return nil
	case tcell.KeyPageDown:
		if len(dpw.State.SelectedCommits) > 0 {
			dpw.SelectedCommitIndex = len(dpw.State.SelectedCommits) - 1
			dpw.UpdateContent()
		}
		return nil
	}

	switch event.Rune() {
	case 'd', 'D':
		dpw.ShowDetails = !dpw.ShowDetails
		dpw.UpdateContent()
		return nil
	case 'j':
		if dpw.SelectedCommitIndex < len(dpw.State.SelectedCommits)-1 {
			dpw.SelectedCommitIndex++
			dpw.UpdateContent()
		}
		return nil
	case 'k':
		if dpw.SelectedCommitIndex > 0 {
			dpw.SelectedCommitIndex--
			dpw.UpdateContent()
		}
		return nil
	}

	return event
}

// HandleInput processes keyboard input for the status bar
func (sbw *StatusBarWidget) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'h':
		// Toggle shortcuts display
		sbw.ToggleShortcuts()
		return nil
	}
	return event
}
