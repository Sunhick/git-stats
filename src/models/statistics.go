// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Models package for statistical data structures

package models

import (
	"time"
)

// AnalysisResult contains the complete analysis results
type AnalysisResult struct {
	Repository    *RepositoryInfo
	Summary       *StatsSummary
	Contributors  []ContributorStats
	ContribGraph  *ContributionGraph
	HealthMetrics *HealthMetrics
	TimeRange     TimeRange
}

// RepositoryInfo contains metadata about the analyzed repository
type RepositoryInfo struct {
	Path         string
	Name         string
	TotalCommits int
	FirstCommit  time.Time
	LastCommit   time.Time
	Branches     []string
}

// StatsSummary contains overall repository statistics
type StatsSummary struct {
	TotalCommits     int
	TotalInsertions  int
	TotalDeletions   int
	FilesChanged     int
	ActiveDays       int
	AvgCommitsPerDay float64
	CommitsByHour    map[int]int
	CommitsByWeekday map[time.Weekday]int
	TopFiles         []FileStats
	TopFileTypes     []FileTypeStats
}

// ContributorStats is an alias for backward compatibility
// Use Contributor from contributor.go for new implementations
type ContributorStats = Contributor

// ContributionGraph represents the GitHub-style contribution graph
type ContributionGraph struct {
	StartDate    time.Time
	EndDate      time.Time
	DailyCommits map[string]int // date -> commit count
	MaxCommits   int
	TotalCommits int
}

// HealthMetrics contains repository health indicators
type HealthMetrics struct {
	RepositoryAge      time.Duration
	CommitFrequency    float64 // commits per day
	ContributorCount   int
	ActiveContributors int // contributors in last 3 months
	BranchCount        int
	ActivityTrend      string // increasing, decreasing, stable
	MonthlyGrowth      []MonthlyStats
}

// FileStats contains statistics for a specific file
type FileStats struct {
	Path         string
	Commits      int
	Insertions   int
	Deletions    int
	LastModified time.Time
}

// FileTypeStats contains statistics for a file type
type FileTypeStats struct {
	Extension string
	Files     int
	Commits   int
	Lines     int
}

// MonthlyStats contains statistics for a specific month
type MonthlyStats struct {
	Month   time.Time
	Commits int
	Authors int
}

// TimeRange represents a time period for analysis
type TimeRange struct {
	Start time.Time
	End   time.Time
}
