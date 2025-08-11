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

func Contrib() {
	fmt.Println("Git Contribution Graph")
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
	fmt.Printf("Total Commits: %d\n\n", repoInfo.TotalCommits)

	// If repository is empty, stop here
	if repoInfo.TotalCommits == 0 {
		fmt.Println("Repository has no commits yet.")
		return
	}

	// Get commits for the last year
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

	// Analyze contributions
	contribAnalyzer := analyzers.NewContributionAnalyzer()
	contribGraph, err := contribAnalyzer.AnalyzeContributions(modelCommits, config)
	if err != nil {
		fmt.Printf("Error analyzing contributions: %v\n", err)
		return
	}

	// Create render configuration
	renderConfig := models.RenderConfig{
		Width:       80,
		Height:      25,
		ColorScheme: "default",
		ShowLegend:  true,
		Interactive: false,
	}

	// Create contribution graph renderer and display
	contribRenderer := visualizers.NewContributionGraphRenderer(renderConfig)

	// Render contribution graph
	graphOutput, err := contribRenderer.RenderContributionGraph(contribGraph, renderConfig)
	if err != nil {
		fmt.Printf("Error rendering contribution graph: %v\n", err)
		return
	}

	fmt.Print(graphOutput)

	// Display contribution summary
	contribSummary := contribAnalyzer.GetContributionSummary(contribGraph)

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

	fmt.Printf("\nNote: Showing activity for the last 365 days\n")
}


