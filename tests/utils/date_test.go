// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for date utilities

package utils_test

import (
	"git-stats/utils"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	testCases := []struct {
		input    string
		expected string // Expected date in RFC3339 format for comparison
		hasError bool
	}{
		{"2024-01-15", "2024-01-15T00:00:00Z", false},
		{"2024-01-15 14:30:00", "2024-01-15T14:30:00Z", false},
		{"2024/01/15", "2024-01-15T00:00:00Z", false},
		{"01/15/2024", "2024-01-15T00:00:00Z", false},
		{"15-01-2024", "2024-01-15T00:00:00Z", false},
		{"invalid-date", "", true},
		{"", "", true},
	}

	for _, tc := range testCases {
		result, err := utils.ParseDate(tc.input)

		if tc.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
				continue
			}

			expected, _ := time.Parse(time.RFC3339, tc.expected)
			if !result.Equal(expected) {
				t.Errorf("For input '%s', expected %v, got %v", tc.input, expected, result)
			}
		}
	}
}

func TestParseRelativeDate(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		input    string
		expected time.Duration // Duration from now
		hasError bool
	}{
		{"1 day ago", -24 * time.Hour, false},
		{"2 days ago", -48 * time.Hour, false},
		{"1 week ago", -7 * 24 * time.Hour, false},
		{"1 hour ago", -time.Hour, false},
		{"30 minutes ago", -30 * time.Minute, false},
		{"invalid relative", 0, true},
		{"", 0, true},
	}

	for _, tc := range testCases {
		result, err := utils.ParseRelativeDate(tc.input)

		if tc.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
				continue
			}

			expected := now.Add(tc.expected)
			// Allow for small time differences due to test execution time
			diff := result.Sub(expected)
			if diff > time.Second || diff < -time.Second {
				t.Errorf("For input '%s', expected around %v, got %v (diff: %v)", tc.input, expected, result, diff)
			}
		}
	}
}

func TestFormatDate(t *testing.T) {
	testDate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	testCases := []struct {
		format   string
		expected string
	}{
		{"iso", "2024-01-15T14:30:00Z"},
		{"short", "2024-01-15"},
		{"long", "January 15, 2024"},
		{"2006-01-02 15:04", "2024-01-15 14:30"},
	}

	for _, tc := range testCases {
		result := utils.FormatDate(testDate, tc.format)
		if result != tc.expected {
			t.Errorf("For format '%s', expected '%s', got '%s'", tc.format, tc.expected, result)
		}
	}
}

func TestGetDateRange(t *testing.T) {
	testCases := []struct {
		input    string
		hasError bool
	}{
		{"today", false},
		{"yesterday", false},
		{"week", false},
		{"this week", false},
		{"last week", false},
		{"month", false},
		{"this month", false},
		{"last month", false},
		{"year", false},
		{"this year", false},
		{"last year", false},
		{"invalid range", true},
	}

	for _, tc := range testCases {
		start, end, err := utils.GetDateRange(tc.input)

		if tc.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
				continue
			}

			if start.After(end) {
				t.Errorf("For input '%s', start date %v should be before end date %v", tc.input, start, end)
			}
		}
	}
}

func TestValidateDateRange(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	// Valid range
	err := utils.ValidateDateRange(start, end)
	if err != nil {
		t.Errorf("Expected no error for valid date range, got: %v", err)
	}

	// Invalid range (start after end)
	err = utils.ValidateDateRange(end, start)
	if err == nil {
		t.Error("Expected error for invalid date range (start after end)")
	}
}

func TestGetContributionGraphDateRange(t *testing.T) {
	endDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC) // A Saturday

	start, end := utils.GetContributionGraphDateRange(endDate)

	// Should be approximately 12 months ago
	expectedStart := endDate.AddDate(-1, 0, 0)
	if start.After(expectedStart.AddDate(0, 0, 7)) || start.Before(expectedStart.AddDate(0, 0, -7)) {
		t.Errorf("Start date %v should be approximately 12 months before %v", start, endDate)
	}

	// Start should be a Sunday
	if start.Weekday() != time.Sunday {
		t.Errorf("Start date %v should be a Sunday, got %v", start, start.Weekday())
	}

	// End should be a Saturday
	if end.Weekday() != time.Saturday {
		t.Errorf("End date %v should be a Saturday, got %v", end, end.Weekday())
	}
}

func TestGetWeekdayName(t *testing.T) {
	testCases := []struct {
		weekday  time.Weekday
		expected string
	}{
		{time.Sunday, "Sunday"},
		{time.Monday, "Monday"},
		{time.Tuesday, "Tuesday"},
		{time.Wednesday, "Wednesday"},
		{time.Thursday, "Thursday"},
		{time.Friday, "Friday"},
		{time.Saturday, "Saturday"},
	}

	for _, tc := range testCases {
		result := utils.GetWeekdayName(tc.weekday)
		if result != tc.expected {
			t.Errorf("For weekday %v, expected '%s', got '%s'", tc.weekday, tc.expected, result)
		}
	}
}

func TestIsWeekend(t *testing.T) {
	testCases := []struct {
		date     time.Time
		expected bool
	}{
		{time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), true},   // Saturday
		{time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), true},   // Sunday
		{time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), false},  // Monday
		{time.Date(2024, 1, 9, 0, 0, 0, 0, time.UTC), false},  // Tuesday
		{time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC), false}, // Wednesday
		{time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC), false}, // Thursday
		{time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC), false}, // Friday
	}

	for _, tc := range testCases {
		result := utils.IsWeekend(tc.date)
		if result != tc.expected {
			t.Errorf("For date %v, expected %v, got %v", tc.date, tc.expected, result)
		}
	}
}

func TestGetQuarter(t *testing.T) {
	testCases := []struct {
		date     time.Time
		expected int
	}{
		{time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), 1},  // January
		{time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC), 1},  // March
		{time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC), 2},  // April
		{time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), 2},  // June
		{time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), 3},  // July
		{time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), 3},  // September
		{time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC), 4}, // October
		{time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC), 4}, // December
	}

	for _, tc := range testCases {
		result := utils.GetQuarter(tc.date)
		if result != tc.expected {
			t.Errorf("For date %v, expected quarter %d, got %d", tc.date, tc.expected, result)
		}
	}
}
