// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - CSV output formatter

package formatters

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git-stats/git"
	"git-stats/models"
)

// CSVFormatterImpl implements the CSVFormatter interface
type CSVFormatterImpl struct{}

// NewCSVFormatter creates a new CSV formatter instance
func NewCSVFormatter() *CSVFormatterImpl {
	return &CSVFormatterImpl{}
}

// Format implements the Formatter interface for CSV output
func (cf *CSVFormatterImpl) Format(data *models.AnalysisResult, config models.FormatConfig) ([]byte, error) {
	if data == nil {
		return nil, NewFormatterError("analysis result cannot be nil")
	}

	return cf.FormatCSV(data, config)
}

// FormatCSV formats analysis results as CSV
func (cf *CSVFormatterImpl) FormatCSV(data *models.AnalysisResult, config models.FormatConfig) ([]byte, error) {
	var result bytes.Buffer

	// Write metadata header if requested
	if config.Metadata {
		if err := cf.writeMetadata(&result, data); err != nil {
			return nil, NewFormatterOperationError("metadata", err.Error())
		}
		result.WriteString("\n")
	}

	// Write contributors CSV
	if len(data.Contributors) > 0 {
		result.WriteString("# Contributors\n")
		contributorsCSV, err := cf.FormatContributorsCSV(data.Contributors)
		if err != nil {
			return nil, NewFormatterOperationError("contributors", err.Error())
		}
		result.Write(contributorsCSV)
		result.WriteString("\n")
	} else {
		result.WriteString("# Contributors\n# No contributors data available\n\n")
	}

	// Write summary statistics CSV
	if data.Summary != nil {
		result.WriteString("# Summary Statistics\n")
		summaryCSV, err := cf.formatSummaryCSV(data.Summary)
		if err != nil {
			return nil, NewFormatterOperationError("summary", err.Error())
		}
		result.Write(summaryCSV)
		result.WriteString("\n")
	} else {
		result.WriteString("# Summary Statistics\n# No summary data available\n\n")
	}

	// Write file statistics CSV
	if data.Summary != nil && len(data.Summary.TopFiles) > 0 {
		result.WriteString("# File Statistics\n")
		filesCSV, err := cf.formatFilesCSV(data.Summary.TopFiles)
		if err != nil {
			return nil, NewFormatterOperationError("files", err.Error())
		}
		result.Write(filesCSV)
		result.WriteString("\n")
	}

	// Write file type statistics CSV
	if data.Summary != nil && len(data.Summary.TopFileTypes) > 0 {
		result.WriteString("# File Type Statistics\n")
		fileTypesCSV, err := cf.formatFileTypesCSV(data.Summary.TopFileTypes)
		if err != nil {
			return nil, NewFormatterOperationError("file_types", err.Error())
		}
		result.Write(fileTypesCSV)
		result.WriteString("\n")
	}

	// Write contribution graph CSV
	if data.ContribGraph != nil {
		result.WriteString("# Daily Contributions\n")
		contribCSV, err := cf.formatContributionGraphCSV(data.ContribGraph)
		if err != nil {
			return nil, NewFormatterOperationError("contribution_graph", err.Error())
		}
		result.Write(contribCSV)
		result.WriteString("\n")
	}

	return result.Bytes(), nil
}

// FormatCommitsCSV formats commits as CSV
func (cf *CSVFormatterImpl) FormatCommitsCSV(commits []git.Commit) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{
		"Hash", "Message", "Author Name", "Author Email",
		"Author Date", "Committer Name", "Committer Email", "Committer Date",
		"Files Changed", "Insertions", "Deletions",
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write commit data
	for _, commit := range commits {
		record := []string{
			commit.Hash,
			commit.Message,
			commit.Author.Name,
			commit.Author.Email,
			cf.formatTimeForCSV(commit.AuthorDate),
			commit.Committer.Name,
			commit.Committer.Email,
			cf.formatTimeForCSV(commit.CommitterDate),
			strconv.Itoa(commit.Stats.FilesChanged),
			strconv.Itoa(commit.Stats.Insertions),
			strconv.Itoa(commit.Stats.Deletions),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write commit record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// FormatContributorsCSV formats contributors as CSV
func (cf *CSVFormatterImpl) FormatContributorsCSV(contributors []models.Contributor) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{
		"Name", "Email", "Total Commits", "Total Insertions", "Total Deletions",
		"First Commit", "Last Commit", "Active Days", "Activity Level",
		"Avg Commits Per Day", "Most Active Hour", "Most Active Weekday", "Top File Type",
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write contributor data
	for _, contributor := range contributors {
		record := []string{
			contributor.Name,
			contributor.Email,
			strconv.Itoa(contributor.TotalCommits),
			strconv.Itoa(contributor.TotalInsertions),
			strconv.Itoa(contributor.TotalDeletions),
			cf.formatTimeForCSV(contributor.FirstCommit),
			cf.formatTimeForCSV(contributor.LastCommit),
			strconv.Itoa(contributor.ActiveDays),
			contributor.GetActivityLevel(),
			fmt.Sprintf("%.2f", contributor.GetAverageCommitsPerDay()),
			strconv.Itoa(contributor.GetMostActiveHour()),
			contributor.GetMostActiveWeekday().String(),
			contributor.GetTopFileType(),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write contributor record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// writeMetadata writes metadata as CSV comments
func (cf *CSVFormatterImpl) writeMetadata(buf *bytes.Buffer, data *models.AnalysisResult) error {
	buf.WriteString("# Metadata\n")
	buf.WriteString(fmt.Sprintf("# Generated at: %s\n", time.Now().UTC().Format(time.RFC3339)))
	buf.WriteString("# Format: CSV\n")
	buf.WriteString("# Version: 1.0\n")

	if data.Repository != nil {
		buf.WriteString(fmt.Sprintf("# Repository: %s\n", data.Repository.Name))
		buf.WriteString(fmt.Sprintf("# Repository Path: %s\n", data.Repository.Path))
	}

	buf.WriteString(fmt.Sprintf("# Analysis Period: %s to %s\n",
		cf.formatTimeForCSV(data.TimeRange.Start),
		cf.formatTimeForCSV(data.TimeRange.End)))

	return nil
}

// formatSummaryCSV formats summary statistics as CSV
func (cf *CSVFormatterImpl) formatSummaryCSV(summary *models.StatsSummary) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{"Metric", "Value"}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write summary CSV header: %w", err)
	}

	// Write summary data
	summaryData := [][]string{
		{"Total Commits", strconv.Itoa(summary.TotalCommits)},
		{"Total Insertions", strconv.Itoa(summary.TotalInsertions)},
		{"Total Deletions", strconv.Itoa(summary.TotalDeletions)},
		{"Files Changed", strconv.Itoa(summary.FilesChanged)},
		{"Active Days", strconv.Itoa(summary.ActiveDays)},
		{"Avg Commits Per Day", fmt.Sprintf("%.2f", summary.AvgCommitsPerDay)},
	}

	for _, record := range summaryData {
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write summary record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("summary CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// formatFilesCSV formats file statistics as CSV
func (cf *CSVFormatterImpl) formatFilesCSV(files []models.FileStats) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{"Path", "Commits", "Insertions", "Deletions", "Last Modified"}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write files CSV header: %w", err)
	}

	// Write file data
	for _, file := range files {
		record := []string{
			file.Path,
			strconv.Itoa(file.Commits),
			strconv.Itoa(file.Insertions),
			strconv.Itoa(file.Deletions),
			cf.formatTimeForCSV(file.LastModified),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write file record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("files CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// formatFileTypesCSV formats file type statistics as CSV
func (cf *CSVFormatterImpl) formatFileTypesCSV(fileTypes []models.FileTypeStats) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{"Extension", "Files", "Commits", "Lines"}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write file types CSV header: %w", err)
	}

	// Write file type data
	for _, fileType := range fileTypes {
		record := []string{
			fileType.Extension,
			strconv.Itoa(fileType.Files),
			strconv.Itoa(fileType.Commits),
			strconv.Itoa(fileType.Lines),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write file type record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("file types CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// formatContributionGraphCSV formats contribution graph as CSV
func (cf *CSVFormatterImpl) formatContributionGraphCSV(graph *models.ContributionGraph) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{"Date", "Commits"}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write contribution graph CSV header: %w", err)
	}

	// Write daily contribution data
	for date, commits := range graph.DailyCommits {
		record := []string{
			date,
			strconv.Itoa(commits),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write contribution record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("contribution graph CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// escapeCSVField properly escapes CSV fields
func (cf *CSVFormatterImpl) escapeCSVField(field string) string {
	// Remove any existing quotes and escape internal quotes
	field = strings.ReplaceAll(field, "\"", "\"\"")

	// If field contains comma, newline, or quote, wrap in quotes
	if strings.Contains(field, ",") || strings.Contains(field, "\n") || strings.Contains(field, "\"") {
		return "\"" + field + "\""
	}

	return field
}

// formatTimeForCSV formats time for CSV output
func (cf *CSVFormatterImpl) formatTimeForCSV(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
