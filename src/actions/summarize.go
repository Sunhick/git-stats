// Copyright (c) 2019 Sunil

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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

func Summarize() {
	SummarizeWithConfig(nil)
}

func SummarizeWithConfig(config *cli.Config) {
	// Use default config if none provided
	if config == nil {
		config = &cli.Config{
			Command:  "summary",
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

	// Analyze statistics
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

	// Analyze health metrics
	healthAnalyzer := analyzers.NewHealthAnalyzer()
	healthMetrics, err := healthAnalyzer.AnalyzeHealth(modelCommits, modelContributors, analysisConfig)
	if err != nil {
		fmt.Printf("Error analyzing health: %v\n", err)
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
		err = outputTerminal(analysisResult, config, "summary")
	}

	if err != nil {
		fmt.Printf("Error generating output: %v\n", err)
		return
	}
}


