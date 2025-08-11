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
	"git-stats/git"
	"git-stats/models"
	"git-stats/visualizers"
	"os"
	"time"
)

func Summarize() {
	fmt.Println("Git Repository Summary")
	fmt.Println("======================")

	// Get current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Create git repository instance
	repoConfig := git.RepositoryConfig{
		Path: workingDir,
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

	fmt.Printf("Repository: %s\n", repoInfo.Name)
	fmt.Printf("Path: %s\n", repoInfo.Path)
	fmt.Printf("Total Commits: %d\n", repoInfo.TotalCommits)

	if !repoInfo.FirstCommit.IsZero() {
		fmt.Printf("First Commit: %s\n", repoInfo.FirstCommit.Format("2006-01-02 15:04:05"))
	}
	if !repoInfo.LastCommit.IsZero() {
		fmt.Printf("Last Commit: %s\n", repoInfo.LastCommit.Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("Branches: %d\n\n", len(repoInfo.Branches))

	// If repository is empty, stop here
	if repoInfo.TotalCommits == 0 {
		fmt.Println("Repository has no commits yet.")
		return
	}

	// Get commits for analysis (last year by default)
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)

	commits, err := repo.GetCommits(startDate, endDate, "")
	if err != nil {
		fmt.Printf("Error getting commits: %v\n", err)
		return
	}

	if len(commits) == 0 {
		fmt.Println("No commits found in the last year.")
		return
	}

	// Convert git.Commit to models.Commit
	modelCommits := convertGitCommitsToModelCommits(commits)

	// Create analysis configuration
	config := models.AnalysisConfig{
		TimeRange: models.TimeRange{
			Start: startDate,
			End:   endDate,
		},
		IncludeMerges: true,
		Limit:         0, // No limit
	}

	// Analyze statistics
	statsAnalyzer := analyzers.NewStatisticsAnalyzer()
	summary, err := statsAnalyzer.AnalyzeStatistics(modelCommits, config)
	if err != nil {
		fmt.Printf("Error analyzing statistics: %v\n", err)
		return
	}

	// Get contributors
	gitContributors, err := repo.GetContributors()
	if err != nil {
		fmt.Printf("Error getting contributors: %v\n", err)
		return
	}

	// Convert git.Contributor to models.Contributor
	modelContributors := convertGitContributorsToModelContributors(gitContributors)

	// Create render configuration
	renderConfig := models.RenderConfig{
		Width:       80,
		Height:      25,
		ColorScheme: "default",
		ShowLegend:  true,
		Interactive: false,
	}

	// Create charts renderer and display summary
	chartsRenderer := visualizers.NewChartsRenderer(renderConfig)

	// Render summary statistics
	summaryOutput, err := chartsRenderer.RenderSummaryStats(summary, renderConfig)
	if err != nil {
		fmt.Printf("Error rendering summary: %v\n", err)
		return
	}

	fmt.Print(summaryOutput)

	// Render contributor statistics if we have contributors
	if len(modelContributors) > 0 {
		fmt.Println("\n")
		contributorOutput, err := chartsRenderer.RenderContributorStats(modelContributors, renderConfig)
		if err != nil {
			fmt.Printf("Error rendering contributors: %v\n", err)
			return
		}

		fmt.Print(contributorOutput)
	}

	// Render time-based analysis
	fmt.Println("\n")
	timeAnalysisOutput, err := chartsRenderer.RenderTimeBasedAnalysis(summary, renderConfig)
	if err != nil {
		fmt.Printf("Error rendering time analysis: %v\n", err)
		return
	}

	fmt.Print(timeAnalysisOutput)
}


