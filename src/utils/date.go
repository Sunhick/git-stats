// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Utils package for date utilities

package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DateParser interface for parsing various date formats
type DateParser interface {
	ParseDate(dateStr string) (time.Time, error)
	ParseRelativeDate(relativeStr string) (time.Time, error)
	FormatDate(date time.Time, format string) string
}

// ParseDate parses various date formats commonly used in git
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, NewGitStatsError(ErrInvalidDateFormat, "empty date string", nil)
	}

	// Try different date formats
	formats := []string{
		"2006-01-02",           // YYYY-MM-DD
		"2006-01-02 15:04:05",  // YYYY-MM-DD HH:MM:SS
		"2006/01/02",           // YYYY/MM/DD
		"01/02/2006",           // MM/DD/YYYY
		"02-01-2006",           // DD-MM-YYYY
		time.RFC3339,           // ISO 8601
		time.RFC822,            // RFC 822
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// Try relative dates
	if t, err := ParseRelativeDate(dateStr); err == nil {
		return t, nil
	}

	return time.Time{}, NewGitStatsError(ErrInvalidDateFormat,
		fmt.Sprintf("unable to parse date: %s", dateStr), nil)
}

// ParseRelativeDate parses relative date strings like "1 week ago", "2 months ago"
func ParseRelativeDate(relativeStr string) (time.Time, error) {
	now := time.Now()
	lower := strings.ToLower(strings.TrimSpace(relativeStr))

	// Handle special cases
	switch lower {
	case "today":
		return now.Truncate(24 * time.Hour), nil
	case "yesterday":
		return now.AddDate(0, 0, -1).Truncate(24 * time.Hour), nil
	case "last week":
		return now.AddDate(0, 0, -7), nil
	case "last month":
		return now.AddDate(0, -1, 0), nil
	case "last year":
		return now.AddDate(-1, 0, 0), nil
	}

	// Parse "N unit ago" format
	parts := strings.Fields(lower)
	if len(parts) >= 3 && parts[len(parts)-1] == "ago" {
		if num, err := strconv.Atoi(parts[0]); err == nil {
			unit := parts[1]
			if strings.HasSuffix(unit, "s") {
				unit = unit[:len(unit)-1] // Remove plural 's'
			}

			switch unit {
			case "second":
				return now.Add(-time.Duration(num) * time.Second), nil
			case "minute":
				return now.Add(-time.Duration(num) * time.Minute), nil
			case "hour":
				return now.Add(-time.Duration(num) * time.Hour), nil
			case "day":
				return now.AddDate(0, 0, -num), nil
			case "week":
				return now.AddDate(0, 0, -num*7), nil
			case "month":
				return now.AddDate(0, -num, 0), nil
			case "year":
				return now.AddDate(-num, 0, 0), nil
			}
		}
	}

	return time.Time{}, NewGitStatsError(ErrInvalidDateFormat,
		fmt.Sprintf("unable to parse relative date: %s", relativeStr), nil)
}

// FormatDate formats a time.Time according to the specified format
func FormatDate(date time.Time, format string) string {
	switch format {
	case "iso":
		return date.Format(time.RFC3339)
	case "short":
		return date.Format("2006-01-02")
	case "long":
		return date.Format("January 2, 2006")
	case "relative":
		return FormatRelativeDate(date)
	default:
		return date.Format(format)
	}
}

// FormatRelativeDate formats a date as a relative string (e.g., "2 days ago")
func FormatRelativeDate(date time.Time) string {
	now := time.Now()
	diff := now.Sub(date)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / (24 * 7))
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := int(diff.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

// GetDateRange returns a start and end date for common ranges
func GetDateRange(rangeStr string) (time.Time, time.Time, error) {
	now := time.Now()

	switch strings.ToLower(rangeStr) {
	case "today":
		start := now.Truncate(24 * time.Hour)
		return start, now, nil
	case "yesterday":
		start := now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
		end := start.Add(24 * time.Hour)
		return start, end, nil
	case "week", "this week":
		// Start of current week (Sunday)
		weekday := int(now.Weekday())
		start := now.AddDate(0, 0, -weekday).Truncate(24 * time.Hour)
		return start, now, nil
	case "last week":
		weekday := int(now.Weekday())
		thisWeekStart := now.AddDate(0, 0, -weekday).Truncate(24 * time.Hour)
		start := thisWeekStart.AddDate(0, 0, -7)
		end := thisWeekStart
		return start, end, nil
	case "month", "this month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return start, now, nil
	case "last month":
		start := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
		end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return start, end, nil
	case "year", "this year":
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		return start, now, nil
	case "last year":
		start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
		end := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		return start, end, nil
	default:
		return time.Time{}, time.Time{}, NewGitStatsError(ErrInvalidDateFormat,
			fmt.Sprintf("unknown date range: %s", rangeStr), nil)
	}
}
