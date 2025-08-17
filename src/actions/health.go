// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Health analysis action

package actions

import (
	"fmt"
	"git-stats/analyzers"
	"git-stats/cli"
	"git-stats/git"
	"git-stats/models"
	"os"
	"time"
)

// Health executes the health analysis with default configuration
func Health() {
	HealthWithConfig(nil)
}

// HealthWithConfig executes the health analysis with the given configuration
func HealthWithConfig(config *cli.Config) {
	// Use default config if none provided
	if config == nil {
		config = &cli.Config{
			Command:  "health",
			RepoPath: ".",
			Format:   "terminal",
			Limit:    10000,
		}
	}

	// Get repository path
	repoPath := config.RepoPath
	if repoPath == "" {
		var err error
		repoPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			return
		}
	}

	// Create git repository instance
	repoConfig := git.RepositoryConfig{
		Path: repoPath,
	}

	repo, err := git.NewGitRepository(repoConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Make sure you're in a git repository directory.")
		return
	}

	// Get repository info
	repoInfo, err := repo.GetRepositoryInfo()
	if err != nil {
		fmt.Printf("Error getting repository info: %v\n", err)
		return
	}

	// If repository is empty, stop here
	if repoInfo.TotalCommits == 0 {
		fmt.Println("Repository has no commits yet.")
		return
	}

	// Determine time range
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)

	if config.Since != nil {
		startDate = *config.Since
	}
	if config.Until != nil {
		endDate = *config.Until
	}

	// Get commits
	commits, err := repo.GetCommits(startDate, endDate, config.Author)
	if err != nil {
		fmt.Printf("Error getting commits: %v\n", err)
		return
	}

	if len(commits) == 0 {
		fmt.Println("No commits found in the specified time range.")
		return
	}

	// Convert git.Commit to models.Commit
	modelCommits := convertGitCommitsToModelCommits(commits)

	// Apply limit if specified
	if config.Limit > 0 && len(modelCommits) > config.Limit {
		modelCommits = modelCommits[:config.Limit]
	}

	// Get contributors
	gitContributors, err := repo.GetContributors()
	if err != nil {
		fmt.Printf("Error getting contributors: %v\n", err)
		return
	}
	modelContributors := convertGitContributorsToModelContributors(gitContributors)

	// Create analysis configuration
	analysisConfig := models.AnalysisConfig{
		TimeRange: models.TimeRange{
			Start: startDate,
			End:   endDate,
		},
		AuthorFilter:  config.Author,
		IncludeMerges: true,
		Limit:         config.Limit,
	}

	// Analyze health metrics
	healthAnalyzer := analyzers.NewHealthAnalyzer()
	healthMetrics, err := healthAnalyzer.AnalyzeHealth(modelCommits, modelContributors, analysisConfig)
	if err != nil {
		fmt.Printf("Error analyzing health: %v\n", err)
		return
	}

	// Analyze statistics for additional context
	statsAnalyzer := analyzers.NewStatisticsAnalyzer()
	summary, err := statsAnalyzer.AnalyzeStatistics(modelCommits, analysisConfig)
	if err != nil {
		fmt.Printf("Error analyzing statistics: %v\n", err)
		return
	}

	// Analyze contributions for additional context
	contribAnalyzer := analyzers.NewContributionAnalyzer()
	contribGraph, err := contribAnalyzer.AnalyzeContributions(modelCommits, analysisConfig)
	if err != nil {
		fmt.Printf("Error analyzing contributions: %v\n", err)
		return
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
			Start: startDate,
			End:   endDate,
		},
	}

	// Handle different output formats
	switch config.Format {
	case "json":
		err = outputJSON(analysisResult, config)
	case "csv":
		err = outputCSV(analysisResult, config)
	default:
		err = outputTerminal(analysisResult, config, "health")
	}

	if err != nil {
		fmt.Printf("Error generating output: %v\n", err)
		return
	}
}
