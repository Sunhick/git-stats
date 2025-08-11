// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Health metrics analysis implementation

package analyzers

import (
	"git-stats/models"
	"sort"
	"strings"
	"time"
)

// HealthAnalyzerImpl implements the HealthAnalyzer interface
type HealthAnalyzerImpl struct{}

// NewHealthAnalyzer creates a new health analyzer
func NewHealthAnalyzer() *HealthAnalyzerImpl {
	return &HealthAnalyzerImpl{}
}

// AnalyzeHealth analyzes repository health metrics
func (ha *HealthAnalyzerImpl) AnalyzeHealth(commits []models.Commit, contributors []models.Contributor, config models.AnalysisConfig) (*models.HealthMetrics, error) {
	if len(commits) == 0 {
		return &models.HealthMetrics{
			RepositoryAge:      0,
			CommitFrequency:    0,
			ContributorCount:   len(contributors),
			ActiveContributors: 0,
			BranchCount:        0,
			ActivityTrend:      "stable",
			MonthlyGrowth:      []models.MonthlyStats{},
		}, nil
	}

	// Filter commits based on configuration
	filteredCommits := ha.filterCommits(commits, config)
	if len(filteredCommits) == 0 {
		return &models.HealthMetrics{
			RepositoryAge:      0,
			CommitFrequency:    0,
			ContributorCount:   len(contributors),
			ActiveContributors: 0,
			BranchCount:        0,
			ActivityTrend:      "stable",
			MonthlyGrowth:      []models.MonthlyStats{},
		}, nil
	}

	// Calculate repository age
	repositoryAge := ha.calculateRepositoryAge(filteredCommits)

	// Calculate commit frequency (commits per day)
	commitFrequency := ha.calculateCommitFrequency(filteredCommits, repositoryAge)

	// Count active contributors (active in last 3 months)
	activeContributors := ha.countActiveContributors(contributors)

	// Calculate activity trend
	activityTrend := ha.CalculateActivityTrend(filteredCommits)

	// Calculate monthly growth
	monthlyGrowth := ha.CalculateMonthlyGrowth(filteredCommits)

	return &models.HealthMetrics{
		RepositoryAge:      repositoryAge,
		CommitFrequency:    commitFrequency,
		ContributorCount:   len(contributors),
		ActiveContributors: activeContributors,
		BranchCount:        0, // This would need to be passed from repository info
		ActivityTrend:      activityTrend,
		MonthlyGrowth:      monthlyGrowth,
	}, nil
}

// CalculateActivityTrend analyzes commit patterns to determine activity trend
func (ha *HealthAnalyzerImpl) CalculateActivityTrend(commits []models.Commit) string {
	if len(commits) < 2 {
		return "stable"
	}

	// Sort commits by date
	sortedCommits := make([]models.Commit, len(commits))
	copy(sortedCommits, commits)
	sort.Slice(sortedCommits, func(i, j int) bool {
		return sortedCommits[i].AuthorDate.Before(sortedCommits[j].AuthorDate)
	})

	// Calculate monthly commit counts
	monthlyCommits := make(map[string]int)
	for _, commit := range sortedCommits {
		monthKey := commit.AuthorDate.Format("2006-01")
		monthlyCommits[monthKey]++
	}

	// Get sorted months
	months := make([]string, 0, len(monthlyCommits))
	for month := range monthlyCommits {
		months = append(months, month)
	}
	sort.Strings(months)

	if len(months) < 3 {
		return "stable"
	}

	// Compare recent months with earlier months
	recentMonths := months[len(months)-3:]
	earlierMonths := months[:len(months)-3]

	// Calculate averages
	recentAvg := 0.0
	for _, month := range recentMonths {
		recentAvg += float64(monthlyCommits[month])
	}
	recentAvg /= float64(len(recentMonths))

	earlierAvg := 0.0
	for _, month := range earlierMonths {
		earlierAvg += float64(monthlyCommits[month])
	}
	earlierAvg /= float64(len(earlierMonths))

	// Determine trend based on comparison
	threshold := 0.2 // 20% change threshold
	change := (recentAvg - earlierAvg) / earlierAvg

	if change > threshold {
		return "increasing"
	} else if change < -threshold {
		return "decreasing"
	}
	return "stable"
}

// CalculateMonthlyGrowth calculates monthly statistics for growth analysis
func (ha *HealthAnalyzerImpl) CalculateMonthlyGrowth(commits []models.Commit) []models.MonthlyStats {
	if len(commits) == 0 {
		return []models.MonthlyStats{}
	}

	// Group commits by month and track unique authors
	monthlyData := make(map[string]*monthlyGrowthData)

	for _, commit := range commits {
		monthKey := commit.AuthorDate.Format("2006-01")

		if data, exists := monthlyData[monthKey]; exists {
			data.commits++
			data.authors[commit.Author.Email] = true
		} else {
			authors := make(map[string]bool)
			authors[commit.Author.Email] = true
			monthlyData[monthKey] = &monthlyGrowthData{
				month:   monthKey,
				commits: 1,
				authors: authors,
			}
		}
	}

	// Convert to sorted slice
	var monthlyStats []models.MonthlyStats
	for monthKey, data := range monthlyData {
		monthTime, err := time.Parse("2006-01", monthKey)
		if err != nil {
			continue // Skip invalid dates
		}

		monthlyStats = append(monthlyStats, models.MonthlyStats{
			Month:   monthTime,
			Commits: data.commits,
			Authors: len(data.authors),
		})
	}

	// Sort by month
	sort.Slice(monthlyStats, func(i, j int) bool {
		return monthlyStats[i].Month.Before(monthlyStats[j].Month)
	})

	return monthlyStats
}

// GetRepositoryHealthScore calculates an overall health score (0-100)
func (ha *HealthAnalyzerImpl) GetRepositoryHealthScore(metrics *models.HealthMetrics) int {
	if metrics == nil {
		return 0
	}

	score := 0

	// Commit frequency score (0-30 points)
	// Good: >1 commit/day, Average: 0.1-1 commit/day, Poor: <0.1 commit/day
	if metrics.CommitFrequency >= 1.0 {
		score += 30
	} else if metrics.CommitFrequency >= 0.1 {
		score += int(20 + (metrics.CommitFrequency-0.1)*10/0.9)
	} else {
		score += int(metrics.CommitFrequency * 200) // Scale 0-0.1 to 0-20
	}

	// Active contributors score (0-25 points)
	if metrics.ActiveContributors >= 5 {
		score += 25
	} else if metrics.ActiveContributors >= 2 {
		score += 15 + (metrics.ActiveContributors-2)*3
	} else if metrics.ActiveContributors == 1 {
		score += 10
	}

	// Activity trend score (0-20 points)
	switch metrics.ActivityTrend {
	case "increasing":
		score += 20
	case "stable":
		score += 15
	case "decreasing":
		score += 5
	}

	// Repository age score (0-15 points)
	// Mature repositories (>1 year) get more points
	ageInDays := int(metrics.RepositoryAge.Hours() / 24)
	if ageInDays >= 365 {
		score += 15
	} else if ageInDays >= 90 {
		score += 10
	} else if ageInDays >= 30 {
		score += 5
	}

	// Monthly growth consistency score (0-10 points)
	if len(metrics.MonthlyGrowth) >= 3 {
		consistency := ha.calculateGrowthConsistency(metrics.MonthlyGrowth)
		score += int(consistency * 10)
	}

	// Ensure score is within bounds
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

// GetHealthInsights provides textual insights about repository health
func (ha *HealthAnalyzerImpl) GetHealthInsights(metrics *models.HealthMetrics) []string {
	var insights []string

	if metrics == nil {
		return []string{"No health data available"}
	}

	// Commit frequency insights
	if metrics.CommitFrequency >= 1.0 {
		insights = append(insights, "High commit frequency indicates active development")
	} else if metrics.CommitFrequency < 0.1 {
		insights = append(insights, "Low commit frequency may indicate inactive project")
	}

	// Contributor insights
	if metrics.ActiveContributors == 0 {
		insights = append(insights, "No active contributors in the last 3 months")
	} else if metrics.ActiveContributors == 1 {
		insights = append(insights, "Single active contributor - consider encouraging more participation")
	} else if metrics.ActiveContributors >= 5 {
		insights = append(insights, "Good contributor diversity with multiple active developers")
	}

	// Activity trend insights
	switch metrics.ActivityTrend {
	case "increasing":
		insights = append(insights, "Repository activity is trending upward")
	case "decreasing":
		insights = append(insights, "Repository activity is declining - may need attention")
	case "stable":
		insights = append(insights, "Repository activity is stable")
	}

	// Repository age insights
	ageInDays := int(metrics.RepositoryAge.Hours() / 24)
	if ageInDays < 30 {
		insights = append(insights, "New repository - still establishing development patterns")
	} else if ageInDays >= 365 {
		insights = append(insights, "Mature repository with established history")
	}

	// Monthly growth insights
	if len(metrics.MonthlyGrowth) >= 3 {
		recent := metrics.MonthlyGrowth[len(metrics.MonthlyGrowth)-1]
		if recent.Commits == 0 {
			insights = append(insights, "No commits in the most recent month")
		} else if recent.Authors == 1 {
			insights = append(insights, "Recent development concentrated to single contributor")
		}
	}

	return insights
}

// calculateRepositoryAge calculates the age of the repository based on commits
func (ha *HealthAnalyzerImpl) calculateRepositoryAge(commits []models.Commit) time.Duration {
	if len(commits) == 0 {
		return 0
	}

	var earliest, latest time.Time
	for i, commit := range commits {
		if i == 0 {
			earliest = commit.AuthorDate
			latest = commit.AuthorDate
		} else {
			if commit.AuthorDate.Before(earliest) {
				earliest = commit.AuthorDate
			}
			if commit.AuthorDate.After(latest) {
				latest = commit.AuthorDate
			}
		}
	}

	return latest.Sub(earliest)
}

// calculateCommitFrequency calculates commits per day
func (ha *HealthAnalyzerImpl) calculateCommitFrequency(commits []models.Commit, repositoryAge time.Duration) float64 {
	if len(commits) == 0 || repositoryAge == 0 {
		return 0
	}

	days := repositoryAge.Hours() / 24
	if days < 1 {
		days = 1 // Minimum 1 day to avoid division by zero
	}

	return float64(len(commits)) / days
}

// countActiveContributors counts contributors active in the last 3 months
func (ha *HealthAnalyzerImpl) countActiveContributors(contributors []models.Contributor) int {
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	activeCount := 0

	for _, contributor := range contributors {
		if contributor.LastCommit.After(threeMonthsAgo) {
			activeCount++
		}
	}

	return activeCount
}

// calculateGrowthConsistency calculates how consistent the monthly growth is (0-1)
func (ha *HealthAnalyzerImpl) calculateGrowthConsistency(monthlyGrowth []models.MonthlyStats) float64 {
	if len(monthlyGrowth) < 2 {
		return 0
	}

	// Calculate variance in monthly commit counts
	var total, mean float64
	for _, month := range monthlyGrowth {
		total += float64(month.Commits)
	}
	mean = total / float64(len(monthlyGrowth))

	var variance float64
	for _, month := range monthlyGrowth {
		diff := float64(month.Commits) - mean
		variance += diff * diff
	}
	variance /= float64(len(monthlyGrowth))

	// Convert variance to consistency score (lower variance = higher consistency)
	if mean == 0 {
		return 0
	}

	// Coefficient of variation (CV) = standard deviation / mean
	cv := (variance / (mean * mean))

	// Convert to consistency score (0-1, where 1 is most consistent)
	consistency := 1.0 / (1.0 + cv)

	return consistency
}

// filterCommits applies configuration filters to commits
func (ha *HealthAnalyzerImpl) filterCommits(commits []models.Commit, config models.AnalysisConfig) []models.Commit {
	var filtered []models.Commit

	for _, commit := range commits {
		// Apply time range filter
		if !config.TimeRange.Start.IsZero() && commit.AuthorDate.Before(config.TimeRange.Start) {
			continue
		}
		if !config.TimeRange.End.IsZero() && commit.AuthorDate.After(config.TimeRange.End) {
			continue
		}

		// Apply author filter
		if config.AuthorFilter != "" && !ha.matchesAuthor(commit.Author, config.AuthorFilter) {
			continue
		}

		// Apply merge commit filter
		if !config.IncludeMerges && commit.IsMergeCommit() {
			continue
		}

		filtered = append(filtered, commit)
	}

	return filtered
}

// matchesAuthor checks if a commit author matches the filter
func (ha *HealthAnalyzerImpl) matchesAuthor(author models.Author, filter string) bool {
	if filter == "" {
		return true
	}

	filterLower := strings.ToLower(filter)
	nameLower := strings.ToLower(author.Name)
	emailLower := strings.ToLower(author.Email)

	return strings.Contains(nameLower, filterLower) || strings.Contains(emailLower, filterLower)
}

// monthlyGrowthData is a helper struct for calculating monthly growth
type monthlyGrowthData struct {
	month   string
	commits int
	authors map[string]bool
}
