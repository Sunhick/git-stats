// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Shared utilities for actions

package actions

import (
	"fmt"
	"git-stats/analyzers"
	"git-stats/cli"
	"git-stats/formatters"
	"git-stats/git"
	"git-stats/models"
	"git-stats/visualizers"
	"os"
)

// convertGitCommitsToModelCommits converts git.Commit slice to models.Commit slice
func convertGitCommitsToModelCommits(gitCommits []git.Commit) []models.Commit {
	modelCommits := make([]models.Commit, len(gitCommits))

	for i, gitCommit := range gitCommits {
		// Convert file changes
		fileChanges := make([]models.FileChange, len(gitCommit.Stats.Files))
		for j, gitFile := range gitCommit.Stats.Files {
			fileChanges[j] = models.FileChange{
				Path:       gitFile.Path,
				Status:     gitFile.Status,
				Insertions: gitFile.Insertions,
				Deletions:  gitFile.Deletions,
			}
		}

		modelCommits[i] = models.Commit{
			Hash:    gitCommit.Hash,
			Message: gitCommit.Message,
			Author: models.Author{
				Name:  gitCommit.Author.Name,
				Email: gitCommit.Author.Email,
			},
			Committer: models.Author{
				Name:  gitCommit.Committer.Name,
				Email: gitCommit.Committer.Email,
			},
			AuthorDate:    gitCommit.AuthorDate,
			CommitterDate: gitCommit.CommitterDate,
			ParentHashes:  gitCommit.ParentHashes,
			TreeHash:      gitCommit.TreeHash,
			Stats: models.CommitStats{
				FilesChanged: gitCommit.Stats.FilesChanged,
				Insertions:   gitCommit.Stats.Insertions,
				Deletions:    gitCommit.Stats.Deletions,
				Files:        fileChanges,
			},
		}
	}

	return modelCommits
}

// convertGitContributorsToModelContributors converts git.Contributor slice to models.Contributor slice
func convertGitContributorsToModelContributors(gitContributors []git.Contributor) []models.Contributor {
	modelContributors := make([]models.Contributor, len(gitContributors))

	for i, gitContrib := range gitContributors {
		modelContributors[i] = models.Contributor{
			Name:             gitContrib.Name,
			Email:            gitContrib.Email,
			TotalCommits:     gitContrib.TotalCommits,
			TotalInsertions:  gitContrib.TotalInsertions,
			TotalDeletions:   gitContrib.TotalDeletions,
			FirstCommit:      gitContrib.FirstCommit,
			LastCommit:       gitContrib.LastCommit,
			ActiveDays:       gitContrib.ActiveDays,
			CommitsByDay:     gitContrib.CommitsByDay,
			CommitsByHour:    make(map[int]int),    // Initialize empty
			CommitsByWeekday: make(map[int]int),    // Initialize empty
			FileTypes:        make(map[string]int), // Initialize empty
			TopFiles:         []string{},           // Initialize empty
		}
	}

	return modelContributors
}

// outputJSON outputs analysis results in JSON format
func outputJSON(data *models.AnalysisResult, config *cli.Config) error {
	formatter := formatters.NewJSONFormatter()

	formatConfig := models.FormatConfig{
		Format:   "json",
		Pretty:   true,
		Metadata: true,
	}

	output, err := formatter.Format(data, formatConfig)
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	return writeOutput(output, config.OutputFile)
}

// outputCSV outputs analysis results in CSV format
func outputCSV(data *models.AnalysisResult, config *cli.Config) error {
	formatter := formatters.NewCSVFormatter()

	formatConfig := models.FormatConfig{
		Format:   "csv",
		Metadata: true,
	}

	output, err := formatter.Format(data, formatConfig)
	if err != nil {
		return fmt.Errorf("failed to format CSV: %w", err)
	}

	return writeOutput(output, config.OutputFile)
}

// outputTerminal outputs analysis results in terminal format
func outputTerminal(data *models.AnalysisResult, config *cli.Config, command string) error {
	// Create render configuration
	renderConfig := models.RenderConfig{
		Width:       80,
		Height:      25,
		ColorScheme: "default",
		ShowLegend:  true,
		Interactive: false,
	}

	// Configure colors based on CLI options
	useColors := !config.NoColor
	colorTheme := config.ColorTheme
	if colorTheme == "" {
		colorTheme = "github"
	}

	switch command {
	case "contrib":
		return outputContribTerminal(data, renderConfig, useColors, colorTheme)
	case "summary":
		return outputSummaryTerminal(data, renderConfig, useColors, colorTheme)
	case "contributors":
		return outputContributorsTerminal(data, renderConfig, useColors, colorTheme)
	case "health":
		return outputHealthTerminal(data, renderConfig, useColors, colorTheme)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// outputContribTerminal outputs contribution graph in terminal format
func outputContribTerminal(data *models.AnalysisResult, renderConfig models.RenderConfig, useColors bool, colorTheme string) error {
	fmt.Println("Git Contribution Graph")
	fmt.Println("======================")

	if data.Repository != nil {
		fmt.Printf("Repository: %s\n", data.Repository.Name)
		fmt.Printf("Total Commits: %d\n\n", data.Repository.TotalCommits)
	}

	if data.ContribGraph == nil {
		fmt.Println("No contribution data available.")
		return nil
	}

	// Create contribution graph renderer
	contribRenderer := visualizers.NewContributionGraphRenderer(renderConfig)
	contribRenderer.SetColorOptions(useColors, colorTheme)

	// Render contribution graph
	graphOutput, err := contribRenderer.RenderContributionGraph(data.ContribGraph, renderConfig)
	if err != nil {
		return fmt.Errorf("error rendering contribution graph: %w", err)
	}

	fmt.Print(graphOutput)

	// Display contribution summary
	contribAnalyzer := analyzers.NewContributionAnalyzer()
	contribSummary := contribAnalyzer.GetContributionSummary(data.ContribGraph)

	fmt.Printf("\nContribution Summary:\n")
	fmt.Printf("====================\n")
	fmt.Printf("Total Commits: %d\n", contribSummary.TotalCommits)
	fmt.Printf("Active Days: %d out of %d days\n", contribSummary.ActiveDays, contribSummary.TotalDays)
	fmt.Printf("Average Commits/Day: %.2f\n", contribSummary.AvgCommitsPerDay)
	fmt.Printf("Max Commits in a Day: %d\n", contribSummary.MaxCommitsPerDay)
	fmt.Printf("Current Streak: %d days\n", contribSummary.CurrentStreak)
	fmt.Printf("Longest Streak: %d days\n", contribSummary.LongestStreak)

	// Show activity level distribution
	fmt.Printf("\nActivity Level Distribution:\n")
	fmt.Printf("============================\n")

	levelCounts := make(map[int]int)
	for _, level := range contribSummary.ActivityLevels {
		levelCounts[level]++
	}

	levelNames := []string{"No activity", "Low activity (1-3)", "Medium activity (4-9)", "High activity (10-19)", "Very high activity (20+)"}
	for level := 0; level <= 4; level++ {
		count := levelCounts[level]
		percentage := float64(count) / float64(contribSummary.TotalDays) * 100
		fmt.Printf("%s: %d days (%.1f%%)\n", levelNames[level], count, percentage)
	}

	fmt.Printf("\nNote: Showing activity for the specified time range\n")
	return nil
}

// outputSummaryTerminal outputs summary statistics in terminal format
func outputSummaryTerminal(data *models.AnalysisResult, renderConfig models.RenderConfig, useColors bool, colorTheme string) error {
	fmt.Println("Git Repository Summary")
	fmt.Println("======================")

	if data.Repository != nil {
		fmt.Printf("Repository: %s\n", data.Repository.Name)
		fmt.Printf("Path: %s\n", data.Repository.Path)
		fmt.Printf("Total Commits: %d\n", data.Repository.TotalCommits)

		if !data.Repository.FirstCommit.IsZero() {
			fmt.Printf("First Commit: %s\n", data.Repository.FirstCommit.Format("2006-01-02 15:04:05"))
		}
		if !data.Repository.LastCommit.IsZero() {
			fmt.Printf("Last Commit: %s\n", data.Repository.LastCommit.Format("2006-01-02 15:04:05"))
		}

		fmt.Printf("Branches: %d\n\n", len(data.Repository.Branches))
	}

	if data.Summary == nil {
		fmt.Println("No summary data available.")
		return nil
	}

	// Create charts renderer
	chartsRenderer := visualizers.NewChartsRenderer(renderConfig)

	// Render summary statistics
	summaryOutput, err := chartsRenderer.RenderSummaryStats(data.Summary, renderConfig)
	if err != nil {
		return fmt.Errorf("error rendering summary: %w", err)
	}

	fmt.Print(summaryOutput)

	// Render contributor statistics if available
	if len(data.Contributors) > 0 {
		fmt.Println("\n")
		contributorOutput, err := chartsRenderer.RenderContributorStats(data.Contributors, renderConfig)
		if err != nil {
			return fmt.Errorf("error rendering contributors: %w", err)
		}

		fmt.Print(contributorOutput)
	}

	// Render time-based analysis
	fmt.Println("\n")
	timeAnalysisOutput, err := chartsRenderer.RenderTimeBasedAnalysis(data.Summary, renderConfig)
	if err != nil {
		return fmt.Errorf("error rendering time analysis: %w", err)
	}

	fmt.Print(timeAnalysisOutput)
	return nil
}

// outputContributorsTerminal outputs contributor statistics in terminal format
func outputContributorsTerminal(data *models.AnalysisResult, renderConfig models.RenderConfig, useColors bool, colorTheme string) error {
	fmt.Println("Git Contributors Analysis")
	fmt.Println("=========================")

	if data.Repository != nil {
		fmt.Printf("Repository: %s\n", data.Repository.Name)
		fmt.Printf("Total Contributors: %d\n\n", len(data.Contributors))
	}

	if len(data.Contributors) == 0 {
		fmt.Println("No contributors found.")
		return nil
	}

	// Create charts renderer
	chartsRenderer := visualizers.NewChartsRenderer(renderConfig)

	// Render contributor statistics
	contributorOutput, err := chartsRenderer.RenderContributorStats(data.Contributors, renderConfig)
	if err != nil {
		return fmt.Errorf("error rendering contributors: %w", err)
	}

	fmt.Print(contributorOutput)
	return nil
}

// outputHealthTerminal outputs health metrics in terminal format
func outputHealthTerminal(data *models.AnalysisResult, renderConfig models.RenderConfig, useColors bool, colorTheme string) error {
	fmt.Println("Repository Health Analysis")
	fmt.Println("==========================")

	if data.Repository != nil {
		fmt.Printf("Repository: %s\n", data.Repository.Name)
		fmt.Printf("Total Commits: %d\n\n", data.Repository.TotalCommits)
	}

	if data.HealthMetrics == nil {
		fmt.Println("No health metrics available.")
		return nil
	}

	// Create charts renderer
	chartsRenderer := visualizers.NewChartsRenderer(renderConfig)

	// Render health metrics
	healthOutput, err := chartsRenderer.RenderHealthMetrics(data.HealthMetrics, renderConfig)
	if err != nil {
		return fmt.Errorf("error rendering health metrics: %w", err)
	}

	fmt.Print(healthOutput)
	return nil
}

// writeOutput writes data to file or stdout
func writeOutput(data []byte, outputFile string) error {
	if outputFile == "" {
		fmt.Print(string(data))
		return nil
	}

	return os.WriteFile(outputFile, data, 0644)
}
