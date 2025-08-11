// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Contribution analysis implementation

package analyzers

import (
	"git-stats/models"
	"sort"
	"strings"
	"time"
)

// ContributionAnalyzerImpl implements the ContributionAnalyzer interface
type ContributionAnalyzerImpl struct{}

// NewContributionAnalyzer creates a new contribution analyzer
func NewContributionAnalyzer() *ContributionAnalyzerImpl {
	return &ContributionAnalyzerImpl{}
}

// AnalyzeContributions analyzes commit data to generate a contribution graph
func (ca *ContributionAnalyzerImpl) AnalyzeContributions(commits []models.Commit, config models.AnalysisConfig) (*models.ContributionGraph, error) {
	if len(commits) == 0 {
		return &models.ContributionGraph{
			StartDate:    config.TimeRange.Start,
			EndDate:      config.TimeRange.End,
			DailyCommits: make(map[string]int),
			MaxCommits:   0,
			TotalCommits: 0,
		}, nil
	}

	// Determine time range
	startDate, endDate := ca.determineTimeRange(commits, config.TimeRange)

	// Initialize daily commits map
	dailyCommits := make(map[string]int)

	// Fill in all dates in range with zero commits
	current := startDate
	for !current.After(endDate) {
		dateKey := current.Format("2006-01-02")
		dailyCommits[dateKey] = 0
		current = current.AddDate(0, 0, 1)
	}

	// Count commits per day
	totalCommits := 0
	maxCommits := 0

	for _, commit := range commits {
		// Apply author filter if specified
		if config.AuthorFilter != "" && !ca.matchesAuthor(commit.Author, config.AuthorFilter) {
			continue
		}

		// Skip merge commits if configured
		if !config.IncludeMerges && commit.IsMergeCommit() {
			continue
		}

		// Check if commit is within time range
		if !ca.isWithinTimeRange(commit.AuthorDate, startDate, endDate) {
			continue
		}

		dateKey := commit.AuthorDate.Format("2006-01-02")
		dailyCommits[dateKey]++
		totalCommits++

		if dailyCommits[dateKey] > maxCommits {
			maxCommits = dailyCommits[dateKey]
		}
	}

	return &models.ContributionGraph{
		StartDate:    startDate,
		EndDate:      endDate,
		DailyCommits: dailyCommits,
		MaxCommits:   maxCommits,
		TotalCommits: totalCommits,
	}, nil
}

// CalculateActivityLevels calculates activity levels for each day
// Returns a map of date -> activity level (0, 1, 2, 3, 4)
func (ca *ContributionAnalyzerImpl) CalculateActivityLevels(dailyCommits map[string]int) map[string]int {
	activityLevels := make(map[string]int)

	for date, commits := range dailyCommits {
		level := ca.getActivityLevel(commits)
		activityLevels[date] = level
	}

	return activityLevels
}

// CalculateStreaks calculates current and longest commit streaks
func (ca *ContributionAnalyzerImpl) CalculateStreaks(dailyCommits map[string]int) (current int, longest int) {
	if len(dailyCommits) == 0 {
		return 0, 0
	}

	// Get sorted dates
	dates := make([]string, 0, len(dailyCommits))
	for date := range dailyCommits {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Calculate longest streak by scanning all dates
	longestStreak := 0
	tempStreak := 0

	for _, date := range dates {
		commits := dailyCommits[date]
		if commits > 0 {
			tempStreak++
			if tempStreak > longestStreak {
				longestStreak = tempStreak
			}
		} else {
			tempStreak = 0
		}
	}

	// Calculate current streak by scanning from the end (most recent dates)
	currentStreak := 0

	// Count backwards from the most recent date
	for i := len(dates) - 1; i >= 0; i-- {
		date := dates[i]
		commits := dailyCommits[date]
		if commits > 0 {
			currentStreak++
		} else {
			break // Stop at first day with no commits
		}
	}

	return currentStreak, longestStreak
}

// GetContributionSummary generates a summary of contribution statistics
func (ca *ContributionAnalyzerImpl) GetContributionSummary(graph *models.ContributionGraph) *ContributionSummary {
	if graph == nil {
		return &ContributionSummary{}
	}

	activityLevels := ca.CalculateActivityLevels(graph.DailyCommits)
	currentStreak, longestStreak := ca.CalculateStreaks(graph.DailyCommits)

	// Calculate average commits per day
	totalDays := len(graph.DailyCommits)
	avgCommitsPerDay := 0.0
	if totalDays > 0 {
		avgCommitsPerDay = float64(graph.TotalCommits) / float64(totalDays)
	}

	// Count active days
	activeDays := 0
	for _, commits := range graph.DailyCommits {
		if commits > 0 {
			activeDays++
		}
	}

	return &ContributionSummary{
		TotalCommits:     graph.TotalCommits,
		MaxCommitsPerDay: graph.MaxCommits,
		ActiveDays:       activeDays,
		TotalDays:        totalDays,
		CurrentStreak:    currentStreak,
		LongestStreak:    longestStreak,
		AvgCommitsPerDay: avgCommitsPerDay,
		ActivityLevels:   activityLevels,
	}
}

// determineTimeRange determines the appropriate time range for analysis
func (ca *ContributionAnalyzerImpl) determineTimeRange(commits []models.Commit, configRange models.TimeRange) (time.Time, time.Time) {
	// If explicit time range is provided, use it
	if !configRange.Start.IsZero() && !configRange.End.IsZero() {
		return configRange.Start, configRange.End
	}

	// Default to past year if no commits or no range specified
	now := time.Now()
	defaultStart := now.AddDate(-1, 0, 0)
	defaultEnd := now

	if len(commits) == 0 {
		return defaultStart, defaultEnd
	}

	// Find actual commit range
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

	// Use configured start/end if provided, otherwise use commit range or defaults
	startDate := defaultStart
	endDate := defaultEnd

	if !configRange.Start.IsZero() {
		startDate = configRange.Start
	} else if !earliest.IsZero() {
		// Use the earlier of default start or first commit
		if earliest.Before(defaultStart) {
			startDate = earliest
		}
	}

	if !configRange.End.IsZero() {
		endDate = configRange.End
	} else if !latest.IsZero() {
		// Use the later of default end or last commit
		if latest.After(defaultEnd) {
			endDate = latest
		}
	}

	return startDate, endDate
}

// getActivityLevel returns activity level based on commit count
// 0: no commits, 1: 1-3 commits, 2: 4-9 commits, 3: 10-19 commits, 4: 20+ commits
func (ca *ContributionAnalyzerImpl) getActivityLevel(commits int) int {
	if commits == 0 {
		return 0
	} else if commits <= 3 {
		return 1
	} else if commits <= 9 {
		return 2
	} else if commits <= 19 {
		return 3
	}
	return 4
}

// matchesAuthor checks if a commit author matches the filter
func (ca *ContributionAnalyzerImpl) matchesAuthor(author models.Author, filter string) bool {
	if filter == "" {
		return true
	}

	// Case-insensitive partial matching on name and email
	filterLower := strings.ToLower(filter)
	nameLower := strings.ToLower(author.Name)
	emailLower := strings.ToLower(author.Email)

	return strings.Contains(nameLower, filterLower) || strings.Contains(emailLower, filterLower)
}

// isWithinTimeRange checks if a date is within the specified range
func (ca *ContributionAnalyzerImpl) isWithinTimeRange(date, start, end time.Time) bool {
	return !date.Before(start) && !date.After(end)
}

// ContributionSummary provides summary statistics for contributions
type ContributionSummary struct {
	TotalCommits     int            `json:"total_commits"`
	MaxCommitsPerDay int            `json:"max_commits_per_day"`
	ActiveDays       int            `json:"active_days"`
	TotalDays        int            `json:"total_days"`
	CurrentStreak    int            `json:"current_streak"`
	LongestStreak    int            `json:"longest_streak"`
	AvgCommitsPerDay float64        `json:"avg_commits_per_day"`
	ActivityLevels   map[string]int `json:"activity_levels"`
}
