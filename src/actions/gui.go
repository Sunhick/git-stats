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

	// Perform analysis
	fmt.Println("Analyzing repository...")

	// Get commits for analysis
	startTime := getStartTime(config.Since)
	endTime := getEndTime(config.Until)

	commits, err := repo.GetCommits(startTime, endTime, config.Author)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get commits: %v\n", err)
		os.Exit(1)
	}

	// Convert git.Commit to models.Commit
	modelCommits := convertGitCommitsToModelCommits(commits)

	// Get contributors
	gitContributors, err := repo.GetContributors()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get contributors: %v\n", err)
		os.Exit(1)
	}

	// Convert git.Contributor to models.Contributor
	modelContributors := convertGitContributorsToModelContributors(gitContributors)

	// Create analysis configuration
	analysisConfig := models.AnalysisConfig{
		TimeRange: models.TimeRange{
			Start: startTime,
			End:   endTime,
		},
		AuthorFilter:  config.Author,
		IncludeMerges: true,
		Limit:         config.Limit,
	}

	// Create analyzers
	contribAnalyzer := analyzers.NewContributionAnalyzer()
	statsAnalyzer := analyzers.NewStatisticsAnalyzer()
	healthAnalyzer := analyzers.NewHealthAnalyzer()

	// Get contribution graph
	contribGraph, err := contribAnalyzer.AnalyzeContributions(modelCommits, analysisConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to generate contribution graph: %v\n", err)
		os.Exit(1)
	}

	// Get statistics
	summary, err := statsAnalyzer.AnalyzeStatistics(modelCommits, analysisConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to generate statistics: %v\n", err)
		os.Exit(1)
	}

	// Get health metrics
	healthMetrics, err := healthAnalyzer.AnalyzeHealth(modelCommits, modelContributors, analysisConfig)
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
		Repository: &models.RepositoryInfo{
			Path:         repoInfo.Path,
			Name:         repoInfo.Name,
			TotalCommits: repoInfo.TotalCommits,
			FirstCommit:  repoInfo.FirstCommit,
			LastCommit:   repoInfo.LastCommit,
			Branches:     repoInfo.Branches,
		},
		Summary:       summary,
		Contributors:  modelContributors,
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
