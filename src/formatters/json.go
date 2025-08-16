// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - JSON output formatter

package formatters

import (
	"encoding/json"
	"time"

	"git-stats/models"
)

// JSONFormatterImpl implements the JSONFormatter interface
type JSONFormatterImpl struct{}

// NewJSONFormatter creates a new JSON formatter instance
func NewJSONFormatter() *JSONFormatterImpl {
	return &JSONFormatterImpl{}
}

// Format implements the Formatter interface for JSON output
func (jf *JSONFormatterImpl) Format(data *models.AnalysisResult, config models.FormatConfig) ([]byte, error) {
	if data == nil {
		return nil, NewFormatterError("analysis result cannot be nil")
	}

	if config.Pretty {
		return jf.FormatPrettyJSON(data, config)
	}
	return jf.FormatJSON(data, config)
}

// FormatJSON formats analysis results as compact JSON
func (jf *JSONFormatterImpl) FormatJSON(data *models.AnalysisResult, config models.FormatConfig) ([]byte, error) {
	output := jf.prepareJSONOutput(data, config)
	return json.Marshal(output)
}

// FormatPrettyJSON formats analysis results as pretty-printed JSON
func (jf *JSONFormatterImpl) FormatPrettyJSON(data *models.AnalysisResult, config models.FormatConfig) ([]byte, error) {
	output := jf.prepareJSONOutput(data, config)
	return json.MarshalIndent(output, "", "  ")
}

// prepareJSONOutput prepares the data structure for JSON serialization
func (jf *JSONFormatterImpl) prepareJSONOutput(data *models.AnalysisResult, config models.FormatConfig) map[string]interface{} {
	output := make(map[string]interface{})

	// Add metadata if requested
	if config.Metadata {
		metadata := map[string]interface{}{
			"generated_at": time.Now().UTC().Format(time.RFC3339),
			"format":       "json",
			"version":      "1.0",
		}

		if data.Repository != nil {
			metadata["repository_path"] = data.Repository.Path
			metadata["repository_name"] = data.Repository.Name
		}

		output["metadata"] = metadata
	}

	// Add repository information
	if data.Repository != nil {
		output["repository"] = jf.formatRepositoryInfo(data.Repository)
	}

	// Add time range
	output["time_range"] = jf.formatTimeRange(data.TimeRange)

	// Add summary statistics
	if data.Summary != nil {
		output["summary"] = jf.formatStatsSummary(data.Summary)
	}

	// Add contributors
	if len(data.Contributors) > 0 {
		output["contributors"] = jf.formatContributors(data.Contributors)
	}

	// Add contribution graph
	if data.ContribGraph != nil {
		output["contribution_graph"] = jf.formatContributionGraph(data.ContribGraph)
	}

	// Add health metrics
	if data.HealthMetrics != nil {
		output["health_metrics"] = jf.formatHealthMetrics(data.HealthMetrics)
	}

	return output
}

// formatRepositoryInfo formats repository information for JSON
func (jf *JSONFormatterImpl) formatRepositoryInfo(repo *models.RepositoryInfo) map[string]interface{} {
	return map[string]interface{}{
		"path":          repo.Path,
		"name":          repo.Name,
		"total_commits": repo.TotalCommits,
		"first_commit":  jf.formatTime(repo.FirstCommit),
		"last_commit":   jf.formatTime(repo.LastCommit),
		"branches":      repo.Branches,
	}
}

// formatTimeRange formats time range for JSON
func (jf *JSONFormatterImpl) formatTimeRange(tr models.TimeRange) map[string]interface{} {
	return map[string]interface{}{
		"start": jf.formatTime(tr.Start),
		"end":   jf.formatTime(tr.End),
	}
}

// formatStatsSummary formats statistics summary for JSON
func (jf *JSONFormatterImpl) formatStatsSummary(summary *models.StatsSummary) map[string]interface{} {
	result := map[string]interface{}{
		"total_commits":       summary.TotalCommits,
		"total_insertions":    summary.TotalInsertions,
		"total_deletions":     summary.TotalDeletions,
		"files_changed":       summary.FilesChanged,
		"active_days":         summary.ActiveDays,
		"avg_commits_per_day": summary.AvgCommitsPerDay,
	}

	// Format commits by hour
	if len(summary.CommitsByHour) > 0 {
		result["commits_by_hour"] = summary.CommitsByHour
	}

	// Format commits by weekday
	if len(summary.CommitsByWeekday) > 0 {
		weekdayMap := make(map[string]int)
		for weekday, count := range summary.CommitsByWeekday {
			weekdayMap[weekday.String()] = count
		}
		result["commits_by_weekday"] = weekdayMap
	}

	// Format top files
	if len(summary.TopFiles) > 0 {
		topFiles := make([]map[string]interface{}, len(summary.TopFiles))
		for i, file := range summary.TopFiles {
			topFiles[i] = map[string]interface{}{
				"path":          file.Path,
				"commits":       file.Commits,
				"insertions":    file.Insertions,
				"deletions":     file.Deletions,
				"last_modified": jf.formatTime(file.LastModified),
			}
		}
		result["top_files"] = topFiles
	}

	// Format top file types
	if len(summary.TopFileTypes) > 0 {
		topFileTypes := make([]map[string]interface{}, len(summary.TopFileTypes))
		for i, fileType := range summary.TopFileTypes {
			topFileTypes[i] = map[string]interface{}{
				"extension": fileType.Extension,
				"files":     fileType.Files,
				"commits":   fileType.Commits,
				"lines":     fileType.Lines,
			}
		}
		result["top_file_types"] = topFileTypes
	}

	return result
}

// formatContributors formats contributors for JSON
func (jf *JSONFormatterImpl) formatContributors(contributors []models.Contributor) []map[string]interface{} {
	result := make([]map[string]interface{}, len(contributors))

	for i, contributor := range contributors {
		contrib := map[string]interface{}{
			"name":              contributor.Name,
			"email":             contributor.Email,
			"total_commits":     contributor.TotalCommits,
			"total_insertions":  contributor.TotalInsertions,
			"total_deletions":   contributor.TotalDeletions,
			"first_commit":      jf.formatTime(contributor.FirstCommit),
			"last_commit":       jf.formatTime(contributor.LastCommit),
			"active_days":       contributor.ActiveDays,
			"activity_level":    contributor.GetActivityLevel(),
			"avg_commits_per_day": contributor.GetAverageCommitsPerDay(),
		}

		// Add commits by day if available
		if len(contributor.CommitsByDay) > 0 {
			contrib["commits_by_day"] = contributor.CommitsByDay
		}

		// Add commits by hour if available
		if len(contributor.CommitsByHour) > 0 {
			contrib["commits_by_hour"] = contributor.CommitsByHour
		}

		// Add commits by weekday if available
		if len(contributor.CommitsByWeekday) > 0 {
			weekdayMap := make(map[string]int)
			for weekday, count := range contributor.CommitsByWeekday {
				weekdayMap[time.Weekday(weekday).String()] = count
			}
			contrib["commits_by_weekday"] = weekdayMap
		}

		// Add file types if available
		if len(contributor.FileTypes) > 0 {
			contrib["file_types"] = contributor.FileTypes
		}

		// Add top files if available
		if len(contributor.TopFiles) > 0 {
			contrib["top_files"] = contributor.TopFiles
		}

		result[i] = contrib
	}

	return result
}

// formatContributionGraph formats contribution graph for JSON
func (jf *JSONFormatterImpl) formatContributionGraph(graph *models.ContributionGraph) map[string]interface{} {
	return map[string]interface{}{
		"start_date":    jf.formatTime(graph.StartDate),
		"end_date":      jf.formatTime(graph.EndDate),
		"daily_commits": graph.DailyCommits,
		"max_commits":   graph.MaxCommits,
		"total_commits": graph.TotalCommits,
	}
}

// formatHealthMetrics formats health metrics for JSON
func (jf *JSONFormatterImpl) formatHealthMetrics(health *models.HealthMetrics) map[string]interface{} {
	result := map[string]interface{}{
		"repository_age_days":  int(health.RepositoryAge.Hours() / 24),
		"commit_frequency":     health.CommitFrequency,
		"contributor_count":    health.ContributorCount,
		"active_contributors":  health.ActiveContributors,
		"branch_count":         health.BranchCount,
		"activity_trend":       health.ActivityTrend,
	}

	// Format monthly growth if available
	if len(health.MonthlyGrowth) > 0 {
		monthlyGrowth := make([]map[string]interface{}, len(health.MonthlyGrowth))
		for i, month := range health.MonthlyGrowth {
			monthlyGrowth[i] = map[string]interface{}{
				"month":   month.Month.Format("2006-01"),
				"commits": month.Commits,
				"authors": month.Authors,
			}
		}
		result["monthly_growth"] = monthlyGrowth
	}

	return result
}

// formatTime formats time for JSON output
func (jf *JSONFormatterImpl) formatTime(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}
