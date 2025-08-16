// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Filter builder for creating filter chains from configuration

package filters

import (
	"fmt"
	"git-stats/cli"
	"git-stats/config"
	"strings"
	"time"
)

// FilterBuilder builds filter chains from CLI configuration and config files
type FilterBuilder struct {
	configManager *config.ConfigManager
}

// NewFilterBuilder creates a new filter builder
func NewFilterBuilder(configManager *config.ConfigManager) *FilterBuilder {
	return &FilterBuilder{
		configManager: configManager,
	}
}

// BuildFromCLIConfig builds a filter chain from CLI configuration
func (fb *FilterBuilder) BuildFromCLIConfig(cliConfig *cli.Config) (*FilterChain, error) {
	chain := NewFilterChain()

	// Add date range filter
	if cliConfig.Since != nil || cliConfig.Until != nil {
		dateFilter := NewDateRangeFilter(cliConfig.Since, cliConfig.Until)
		chain.Add(dateFilter)
	}

	// Add author filter
	if cliConfig.Author != "" {
		authorFilter, err := fb.buildAuthorFilter(cliConfig.Author)
		if err != nil {
			return nil, fmt.Errorf("failed to build author filter: %w", err)
		}
		chain.Add(authorFilter)
	}

	// Add merge commit filter based on config
	appConfig := fb.configManager.GetConfig()
	mergeFilter := NewMergeCommitFilter(appConfig.Filters.IncludeMerges)
	chain.Add(mergeFilter)

	// Add limit filter
	if cliConfig.Limit > 0 {
		limitFilter := NewLimitFilter(cliConfig.Limit)
		chain.Add(limitFilter)
	}

	// Add file path filters from config
	if len(appConfig.Filters.IncludePatterns) > 0 {
		includeFilter := NewFilePathFilter(
			appConfig.Filters.IncludePatterns,
			FileContainsMatch,
			appConfig.Filters.CaseSensitive,
		)
		chain.Add(includeFilter)
	}

	if len(appConfig.Filters.ExcludePatterns) > 0 {
		excludeFilter := NewExcludeFilePathFilter(
			appConfig.Filters.ExcludePatterns,
			FileContainsMatch,
			appConfig.Filters.CaseSensitive,
		)
		chain.Add(excludeFilter)
	}

	return chain, nil
}

// BuildFromConfig builds a filter chain from application configuration only
func (fb *FilterBuilder) BuildFromConfig() (*FilterChain, error) {
	chain := NewFilterChain()
	appConfig := fb.configManager.GetConfig()

	// Add default date range if specified
	if appConfig.Defaults.DateRange != "" && appConfig.Defaults.DateRange != "all" {
		since, err := fb.parseDateRange(appConfig.Defaults.DateRange)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default date range: %w", err)
		}
		dateFilter := NewDateRangeFilter(since, nil)
		chain.Add(dateFilter)
	}

	// Add default author filter
	if appConfig.Filters.DefaultAuthor != "" {
		authorFilter, err := fb.buildAuthorFilter(appConfig.Filters.DefaultAuthor)
		if err != nil {
			return nil, fmt.Errorf("failed to build default author filter: %w", err)
		}
		chain.Add(authorFilter)
	}

	// Add merge commit filter
	mergeFilter := NewMergeCommitFilter(appConfig.Filters.IncludeMerges)
	chain.Add(mergeFilter)

	// Add performance limit
	if appConfig.Performance.MaxCommits > 0 {
		limitFilter := NewLimitFilter(appConfig.Performance.MaxCommits)
		chain.Add(limitFilter)
	}

	return chain, nil
}

// buildAuthorFilter builds an author filter based on configuration
func (fb *FilterBuilder) buildAuthorFilter(authorPattern string) (*AuthorFilter, error) {
	appConfig := fb.configManager.GetConfig()

	// Determine match type from config
	var matchType AuthorMatchType
	switch strings.ToLower(appConfig.Filters.AuthorMatchType) {
	case "exact":
		matchType = ExactMatch
	case "contains":
		matchType = ContainsMatch
	case "regex":
		matchType = RegexMatch
	case "email":
		matchType = EmailMatch
	case "name":
		matchType = NameMatch
	default:
		matchType = ContainsMatch
	}

	return NewAuthorFilterWithOptions(
		authorPattern,
		matchType,
		appConfig.Filters.CaseSensitive,
	)
}

// parseDateRange parses relative date ranges like "1 year ago", "6 months ago"
func (fb *FilterBuilder) parseDateRange(dateRange string) (*time.Time, error) {
	dateRange = strings.ToLower(strings.TrimSpace(dateRange))
	now := time.Now()

	// Handle special cases
	switch dateRange {
	case "today":
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return &today, nil
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		yesterday = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
		return &yesterday, nil
	case "this week":
		// Start of current week (Sunday)
		weekday := int(now.Weekday())
		startOfWeek := now.AddDate(0, 0, -weekday)
		startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
		return &startOfWeek, nil
	case "this month":
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return &startOfMonth, nil
	case "this year":
		startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		return &startOfYear, nil
	}

	// Parse relative dates like "1 year ago", "6 months ago"
	if strings.HasSuffix(dateRange, " ago") {
		return fb.parseRelativeDate(dateRange)
	}

	// Try to parse as absolute date
	return fb.parseAbsoluteDate(dateRange)
}

// parseRelativeDate parses relative date strings
func (fb *FilterBuilder) parseRelativeDate(dateStr string) (*time.Time, error) {
	dateStr = strings.TrimSuffix(dateStr, " ago")
	parts := strings.Fields(dateStr)

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid relative date format: %s", dateStr)
	}

	var amount int
	var unit string

	// Parse amount
	switch parts[0] {
	case "a", "an", "1":
		amount = 1
	case "2", "two":
		amount = 2
	case "3", "three":
		amount = 3
	case "4", "four":
		amount = 4
	case "5", "five":
		amount = 5
	case "6", "six":
		amount = 6
	case "7", "seven":
		amount = 7
	case "8", "eight":
		amount = 8
	case "9", "nine":
		amount = 9
	case "10", "ten":
		amount = 10
	default:
		// Try to parse as number
		if _, err := fmt.Sscanf(parts[0], "%d", &amount); err != nil {
			return nil, fmt.Errorf("invalid amount in relative date: %s", parts[0])
		}
	}

	// Parse unit
	unit = strings.ToLower(parts[1])
	if strings.HasSuffix(unit, "s") {
		unit = strings.TrimSuffix(unit, "s")
	}

	now := time.Now()
	var result time.Time

	switch unit {
	case "day":
		result = now.AddDate(0, 0, -amount)
	case "week":
		result = now.AddDate(0, 0, -amount*7)
	case "month":
		result = now.AddDate(0, -amount, 0)
	case "year":
		result = now.AddDate(-amount, 0, 0)
	default:
		return nil, fmt.Errorf("unsupported time unit: %s", unit)
	}

	return &result, nil
}

// parseAbsoluteDate parses absolute date strings
func (fb *FilterBuilder) parseAbsoluteDate(dateStr string) (*time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"01/02/2006",
		"02-01-2006",
		"January 2, 2006",
		"Jan 2, 2006",
		"2 January 2006",
		"2 Jan 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("unable to parse date: %s", dateStr)
}

// BuildAdvancedFilter builds a filter with advanced options
func (fb *FilterBuilder) BuildAdvancedFilter(options AdvancedFilterOptions) (*FilterChain, error) {
	chain := NewFilterChain()

	// Add date range filter
	if options.Since != nil || options.Until != nil {
		dateFilter := NewDateRangeFilter(options.Since, options.Until)
		chain.Add(dateFilter)
	}

	// Add author filters
	for _, authorOpts := range options.Authors {
		authorFilter, err := NewAuthorFilterWithOptions(
			authorOpts.Pattern,
			authorOpts.MatchType,
			authorOpts.CaseSensitive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create author filter: %w", err)
		}
		chain.Add(authorFilter)
	}

	// Add file path filters
	if len(options.IncludeFiles) > 0 {
		includeFilter := NewFilePathFilter(
			options.IncludeFiles,
			options.FileMatchType,
			options.CaseSensitive,
		)
		chain.Add(includeFilter)
	}

	if len(options.ExcludeFiles) > 0 {
		excludeFilter := NewExcludeFilePathFilter(
			options.ExcludeFiles,
			options.FileMatchType,
			options.CaseSensitive,
		)
		chain.Add(excludeFilter)
	}

	// Add merge commit filter
	mergeFilter := NewMergeCommitFilter(options.IncludeMerges)
	chain.Add(mergeFilter)

	// Add limit filter
	if options.Limit > 0 {
		limitFilter := NewLimitFilter(options.Limit)
		chain.Add(limitFilter)
	}

	return chain, nil
}

// AdvancedFilterOptions contains options for building advanced filters
type AdvancedFilterOptions struct {
	Since         *time.Time
	Until         *time.Time
	Authors       []AuthorFilterOptions
	IncludeFiles  []string
	ExcludeFiles  []string
	FileMatchType FileMatchType
	CaseSensitive bool
	IncludeMerges bool
	Limit         int
}

// AuthorFilterOptions contains options for author filtering
type AuthorFilterOptions struct {
	Pattern       string
	MatchType     AuthorMatchType
	CaseSensitive bool
}

// GetFilterSummary returns a human-readable summary of active filters
func (fb *FilterBuilder) GetFilterSummary(chain *FilterChain) string {
	if chain == nil || len(chain.GetFilters()) == 0 {
		return "No filters applied"
	}

	var descriptions []string
	for _, filter := range chain.GetFilters() {
		desc := filter.Description()
		if desc != "No date filter" && desc != "Include merge commits" {
			descriptions = append(descriptions, desc)
		}
	}

	if len(descriptions) == 0 {
		return "Default filters applied"
	}

	return "Active filters: " + strings.Join(descriptions, "; ")
}
