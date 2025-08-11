// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Visualizers package for rendering interfaces

package visualizers

import (
	"git-stats/models"
	"time"
)

// Visualizer interface for rendering output
type Visualizer interface {
	Render(data *models.AnalysisResult, config models.RenderConfig) (string, error)
}

// ContributionGraphVisualizer interface for contribution graph rendering
type ContributionGraphVisualizer interface {
	RenderContributionGraph(graph *models.ContributionGraph, config models.RenderConfig) (string, error)
	RenderMonthLabels(startDate, endDate time.Time) string
	RenderDayIndicators() string
}

// ChartsVisualizer interface for charts and tables rendering
type ChartsVisualizer interface {
	RenderBarChart(data map[string]int, title string, config models.RenderConfig) (string, error)
	RenderTable(headers []string, rows [][]string, config models.RenderConfig) (string, error)
	RenderSummaryStats(summary *models.StatsSummary, config models.RenderConfig) (string, error)
}

// GUIVisualizer interface for ncurses GUI rendering
type GUIVisualizer interface {
	Initialize() error
	Run(data *models.AnalysisResult) error
	Cleanup() error
	HandleInput() error
	Render() error
}
