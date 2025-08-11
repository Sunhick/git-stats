// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for date utilities

package utils_test

import (
	"testing"
	"time"
	"git-stats/utils"
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
