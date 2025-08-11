// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Analyzers package for statistical analysis interfaces

package analyzers

import (
	"time"
	"git-stats/git"
	"git-stats/models"
)

// Analyzer interface for statistical analysis
type Analyzer interface {
	Analyze(commits []git.Commit, config models.AnalysisConfig) (*models.AnalysisResult, error)
}

// ContributionAnalyzer interface for contribution analysis
type ContributionAnalyzer interface {
	AnalyzeContributions(commits []git.Commit, config models.AnalysisConfig) (*models.ContributionGraph, error)
	CalculateActivityLevels(dailyCommits map[string]int) map[string]int
	CalculateStreaks(dailyCommits map[string]int) (current int, longest int)
}

// StatisticsAnalyzer interface for general statistics analysis
type StatisticsAnalyzer interface {
	AnalyzeStatistics(commits []git.Commit, config models.AnalysisConfig) (*models.StatsSummary, error)
	AnalyzeCommitPatterns(commits []git.Commit) (map[int]int, map[time.Weekday]int)
	AnalyzeFileStatistics(commits []git.Commit) ([]models.FileStats, []models.FileTypeStats)
}

// HealthAnalyzer interface for repository health analysis
type HealthAnalyzer interface {
	AnalyzeHealth(commits []git.Commit, contributors []git.Contributor, config models.AnalysisConfig) (*models.HealthMetrics, error)
	CalculateActivityTrend(commits []git.Commit) string
	CalculateMonthlyGrowth(commits []git.Commit) []models.MonthlyStats
}
