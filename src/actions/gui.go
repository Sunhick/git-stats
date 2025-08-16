// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - GUI action implementation

package actions

import (
	"fmt"
	"git-stats/analyzers"
	"git-stats/cli"
	"git-stats/git"
	"git-stats/models"
	"git-stats/visualizers"
	"os"
	"time"
)

// LaunchGUI launches the GUI interface with the specified configuration
func LaunchGUI(config *cli.Config) {
	// Initialize git repository
	repo, err := git.NewRepository(config.RepoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open repository at %s: %v\n", config.RepoPath, err)
		os.Exit(1)
	}

	// Create analyzers
	contribAnalyzer := analyzers.NewContributionAnalyzer(repo)
	statsAnalyzer := analyzers.NewStatisticsAnalyzer(repo)
	healthAnalyzer := analyzers.NewHealthAnalyzer(repo)

	// Perform analysis
	fmt.Println("Analyzing repository...")

	// Get contribution graph
	contribGraph, err := contribAnalyzer.GenerateContributionGraph(config.Since, config.Until)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to generate contribution graph: %v\n", err)
		os.Exit(1)
	}

	// Get statistics
	summary, err := statsAnalyzer.GenerateStatsSummary(config.Since, config.Until, config.Author)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to generate statistics: %v\n", err)
		os.Exit(1)
	}

	// Get contributors
	contributors, err := statsAnalyzer.GetContributors(config.Since, config.Until, config.Author)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get contributors: %v\n", err)
		os.Exit(1)
	}

	// Get health metrics
	healthMetrics, err := healthAnalyzer.AnalyzeHealth(config.Since, config.Until)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to analyze health: %v\n", err)
		os.Exit(1)
	}

	// Get repository info
	repoInfo, err := repo.GetRepositoryInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get repository info: %v\n", err)
		os.Exit(1)
	}

	// Create analysis result
	analysisResult := &models.AnalysisResult{
		Repository:    repoInfo,
		Summary:       summary,
		Contributors:  contributors,
		ContribGraph:  contribGraph,
		HealthMetrics: healthMetrics,
		TimeRange: models.TimeRange{
			Start: getStartTime(config.Since),
			End:   getEndTime(config.Until),
		},
	}

	// Launch GUI
	fmt.Println("Launching GUI interface...")
	gui := visualizers.NewGUIInterface()

	err = gui.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize GUI: %v\n", err)
		os.Exit(1)
	}

	// Run the GUI
	err = gui.Run(analysisResult)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: GUI execution failed: %v\n", err)
		os.Exit(1)
	}

	// Cleanup
	err = gui.Cleanup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: GUI cleanup failed: %v\n", err)
	}
}

// Helper functions to handle time ranges
func getStartTime(since *time.Time) time.Time {
	if since != nil {
		return *since
	}
	// Default to one year ago
	return time.Now().AddDate(-1, 0, 0)
}

func getEndTime(until *time.Time) time.Time {
	if until != nil {
		return *until
	}
	// Default to now
	return time.Now()
}
