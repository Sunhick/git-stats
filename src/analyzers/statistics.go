// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Statistics analysis implementation

package analyzers

import (
	"fmt"
	"git-stats/models"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// StatisticsAnalyzerImpl implements the StatisticsAnalyzer interface
type StatisticsAnalyzerImpl struct{}

// NewStatisticsAnalyzer creates a new statistics analyzer
func NewStatisticsAnalyzer() *StatisticsAnalyzerImpl {
	return &StatisticsAnalyzerImpl{}
}

// AnalyzeStatistics analyzes commit data to generate comprehensive statistics
func (sa *StatisticsAnalyzerImpl) AnalyzeStatistics(commits []models.Commit, config models.AnalysisConfig) (*models.StatsSummary, error) {
	if len(commits) == 0 {
		return &models.StatsSummary{
			CommitsByHour:    make(map[int]int),
			CommitsByWeekday: make(map[time.Weekday]int),
			TopFiles:         []models.FileStats{},
			TopFileTypes:     []models.FileTypeStats{},
		}, nil
	}

	// Filter commits based on configuration
	filteredCommits := sa.filterCommits(commits, config)

	// Initialize summary
	summary := &models.StatsSummary{
		CommitsByHour:    make(map[int]int),
		CommitsByWeekday: make(map[time.Weekday]int),
	}

	// Calculate basic statistics
	summary.TotalCommits = len(filteredCommits)
	activeDaysMap := make(map[string]bool)

	for _, commit := range filteredCommits {
		// Count insertions and deletions
		summary.TotalInsertions += commit.Stats.Insertions
		summary.TotalDeletions += commit.Stats.Deletions
		summary.FilesChanged += commit.Stats.FilesChanged

		// Track active days
		dateKey := commit.AuthorDate.Format("2006-01-02")
		activeDaysMap[dateKey] = true
	}

	summary.ActiveDays = len(activeDaysMap)

	// Calculate average commits per day
	if summary.ActiveDays > 0 {
		summary.AvgCommitsPerDay = float64(summary.TotalCommits) / float64(summary.ActiveDays)
	}

	// Analyze commit patterns (hours and weekdays)
	summary.CommitsByHour, summary.CommitsByWeekday = sa.AnalyzeCommitPatterns(filteredCommits)

	// Analyze file statistics
	summary.TopFiles, summary.TopFileTypes = sa.AnalyzeFileStatistics(filteredCommits)

	return summary, nil
}

// AnalyzeCommitPatterns analyzes temporal patterns in commits
func (sa *StatisticsAnalyzerImpl) AnalyzeCommitPatterns(commits []models.Commit) (map[int]int, map[time.Weekday]int) {
	hourCounts := make(map[int]int)
	weekdayCounts := make(map[time.Weekday]int)

	for _, commit := range commits {
		// Count by hour
		hour := commit.AuthorDate.Hour()
		hourCounts[hour]++

		// Count by weekday
		weekday := commit.AuthorDate.Weekday()
		weekdayCounts[weekday]++
	}

	return hourCounts, weekdayCounts
}

// AnalyzeFileStatistics analyzes file and file type statistics
func (sa *StatisticsAnalyzerImpl) AnalyzeFileStatistics(commits []models.Commit) ([]models.FileStats, []models.FileTypeStats) {
	fileStatsMap := make(map[string]*models.FileStats)
	fileTypeStatsMap := make(map[string]*models.FileTypeStats)
	fileTypeFilesMap := make(map[string]map[string]bool) // extension -> set of file paths

	for _, commit := range commits {
		for _, fileChange := range commit.Stats.Files {
			// Update file statistics
			if fileStats, exists := fileStatsMap[fileChange.Path]; exists {
				fileStats.Commits++
				fileStats.Insertions += fileChange.Insertions
				fileStats.Deletions += fileChange.Deletions
				if commit.AuthorDate.After(fileStats.LastModified) {
					fileStats.LastModified = commit.AuthorDate
				}
			} else {
				fileStatsMap[fileChange.Path] = &models.FileStats{
					Path:         fileChange.Path,
					Commits:      1,
					Insertions:   fileChange.Insertions,
					Deletions:    fileChange.Deletions,
					LastModified: commit.AuthorDate,
				}
			}

			// Update file type statistics
			extension := sa.getFileExtension(fileChange.Path)
			if extension == "" {
				extension = "no-extension"
			}

			// Initialize file set for this extension if needed
			if fileTypeFilesMap[extension] == nil {
				fileTypeFilesMap[extension] = make(map[string]bool)
			}

			if typeStats, exists := fileTypeStatsMap[extension]; exists {
				typeStats.Commits++
				typeStats.Lines += fileChange.Insertions + fileChange.Deletions
				// Track unique files for this type
				if !fileTypeFilesMap[extension][fileChange.Path] {
					fileTypeFilesMap[extension][fileChange.Path] = true
					typeStats.Files++
				}
			} else {
				fileTypeStatsMap[extension] = &models.FileTypeStats{
					Extension: extension,
					Files:     1,
					Commits:   1,
					Lines:     fileChange.Insertions + fileChange.Deletions,
				}
				fileTypeFilesMap[extension][fileChange.Path] = true
			}
		}
	}

	// Convert maps to sorted slices
	topFiles := sa.sortFileStats(fileStatsMap)
	topFileTypes := sa.sortFileTypeStats(fileTypeStatsMap)

	return topFiles, topFileTypes
}

// GetCommitFrequencyAnalysis analyzes commit frequency over different time periods
func (sa *StatisticsAnalyzerImpl) GetCommitFrequencyAnalysis(commits []models.Commit) *CommitFrequencyAnalysis {
	if len(commits) == 0 {
		return &CommitFrequencyAnalysis{
			Daily:   make(map[string]int),
			Weekly:  make(map[string]int),
			Monthly: make(map[string]int),
		}
	}

	analysis := &CommitFrequencyAnalysis{
		Daily:   make(map[string]int),
		Weekly:  make(map[string]int),
		Monthly: make(map[string]int),
	}

	for _, commit := range commits {
		// Daily frequency
		dailyKey := commit.AuthorDate.Format("2006-01-02")
		analysis.Daily[dailyKey]++

		// Weekly frequency (ISO week)
		year, week := commit.AuthorDate.ISOWeek()
		weeklyKey := fmt.Sprintf("%d-W%02d", year, week)
		analysis.Weekly[weeklyKey]++

		// Monthly frequency
		monthlyKey := commit.AuthorDate.Format("2006-01")
		analysis.Monthly[monthlyKey]++
	}

	return analysis
}

// GetTimeBasedPatterns analyzes patterns based on time of day and day of week
func (sa *StatisticsAnalyzerImpl) GetTimeBasedPatterns(commits []models.Commit) *TimeBasedPatterns {
	patterns := &TimeBasedPatterns{
		HourlyDistribution:   make(map[int]int),
		WeekdayDistribution: make(map[time.Weekday]int),
		HourlyAverage:       make(map[int]float64),
		WeekdayAverage:      make(map[time.Weekday]float64),
	}

	if len(commits) == 0 {
		return patterns
	}

	// Count commits by hour and weekday
	hourCounts, weekdayCounts := sa.AnalyzeCommitPatterns(commits)
	patterns.HourlyDistribution = hourCounts
	patterns.WeekdayDistribution = weekdayCounts

	// Calculate averages (commits per hour/weekday over the total period)
	if len(commits) > 0 {
		// Find date range
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

		// Calculate total days in range
		totalDays := int(latest.Sub(earliest).Hours()/24) + 1
		if totalDays > 0 {
			// Calculate hourly averages
			for hour := 0; hour < 24; hour++ {
				patterns.HourlyAverage[hour] = float64(hourCounts[hour]) / float64(totalDays)
			}

			// Calculate weekday averages (number of each weekday in the period)
			totalWeeks := float64(totalDays) / 7.0
			for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
				patterns.WeekdayAverage[weekday] = float64(weekdayCounts[weekday]) / totalWeeks
			}
		}
	}

	return patterns
}

// filterCommits applies configuration filters to commits
func (sa *StatisticsAnalyzerImpl) filterCommits(commits []models.Commit, config models.AnalysisConfig) []models.Commit {
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
		if config.AuthorFilter != "" && !sa.matchesAuthor(commit.Author, config.AuthorFilter) {
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
func (sa *StatisticsAnalyzerImpl) matchesAuthor(author models.Author, filter string) bool {
	if filter == "" {
		return true
	}

	filterLower := strings.ToLower(filter)
	nameLower := strings.ToLower(author.Name)
	emailLower := strings.ToLower(author.Email)

	return strings.Contains(nameLower, filterLower) || strings.Contains(emailLower, filterLower)
}

// getFileExtension extracts file extension from path
func (sa *StatisticsAnalyzerImpl) getFileExtension(path string) string {
	ext := filepath.Ext(path)
	if ext != "" && len(ext) > 1 {
		return ext[1:] // Remove the dot
	}
	return ""
}



// sortFileStats sorts file statistics by commit count (descending)
func (sa *StatisticsAnalyzerImpl) sortFileStats(fileStatsMap map[string]*models.FileStats) []models.FileStats {
	var fileStats []models.FileStats
	for _, stats := range fileStatsMap {
		fileStats = append(fileStats, *stats)
	}

	sort.Slice(fileStats, func(i, j int) bool {
		return fileStats[i].Commits > fileStats[j].Commits
	})

	// Return top 20 files
	if len(fileStats) > 20 {
		fileStats = fileStats[:20]
	}

	return fileStats
}

// sortFileTypeStats sorts file type statistics by commit count (descending)
func (sa *StatisticsAnalyzerImpl) sortFileTypeStats(fileTypeStatsMap map[string]*models.FileTypeStats) []models.FileTypeStats {
	var fileTypeStats []models.FileTypeStats
	for _, stats := range fileTypeStatsMap {
		fileTypeStats = append(fileTypeStats, *stats)
	}

	sort.Slice(fileTypeStats, func(i, j int) bool {
		return fileTypeStats[i].Commits > fileTypeStats[j].Commits
	})

	// Return top 15 file types
	if len(fileTypeStats) > 15 {
		fileTypeStats = fileTypeStats[:15]
	}

	return fileTypeStats
}

// CommitFrequencyAnalysis contains frequency analysis results
type CommitFrequencyAnalysis struct {
	Daily   map[string]int `json:"daily"`   // date -> commit count
	Weekly  map[string]int `json:"weekly"`  // week -> commit count
	Monthly map[string]int `json:"monthly"` // month -> commit count
}

// TimeBasedPatterns contains time-based pattern analysis
type TimeBasedPatterns struct {
	HourlyDistribution   map[int]int             `json:"hourly_distribution"`   // hour -> commit count
	WeekdayDistribution  map[time.Weekday]int    `json:"weekday_distribution"`  // weekday -> commit count
	HourlyAverage        map[int]float64         `json:"hourly_average"`        // hour -> average commits per day
	WeekdayAverage       map[time.Weekday]float64 `json:"weekday_average"`       // weekday -> average commits per week
}
