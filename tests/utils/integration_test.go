// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Integration tests for utility functions

package utils_test

import (
	"errors"
	"git-stats/utils"
	"testing"
	"time"
)

func TestDateUtilitiesIntegration(t *testing.T) {
	// Test parsing and formatting workflow
	dateStr := "2024-01-15"
	parsed, err := utils.ParseDate(dateStr)
	if err != nil {
		t.Fatalf("Failed to parse date: %v", err)
	}

	formatted := utils.FormatDate(parsed, "short")
	if formatted != dateStr {
		t.Errorf("Expected formatted date '%s', got '%s'", dateStr, formatted)
	}

	// Test relative date parsing and formatting
	relativeStr := "1 week ago"
	relativeParsed, err := utils.ParseRelativeDate(relativeStr)
	if err != nil {
		t.Fatalf("Failed to parse relative date: %v", err)
	}

	relativeFormatted := utils.FormatRelativeDate(relativeParsed)
	if relativeFormatted != "1 week ago" {
		t.Logf("Relative formatted: %s (expected approximately '1 week ago')", relativeFormatted)
	}

	// Test date range validation
	start, end, err := utils.GetDateRange("last month")
	if err != nil {
		t.Fatalf("Failed to get date range: %v", err)
	}

	err = utils.ValidateDateRange(start, end)
	if err != nil {
		t.Errorf("Date range validation failed: %v", err)
	}
}

func TestProgressTrackingIntegration(t *testing.T) {
	// Test regular progress tracking
	tracker := utils.NewProgressTracker(100, "Processing items")
	tracker.Start(100, "Starting processing")

	for i := 0; i < 100; i += 10 {
		tracker.Update(i, "Processing...")
		current, total, percentage := tracker.GetProgress()

		if current != i {
			t.Errorf("Expected current %d, got %d", i, current)
		}

		if total != 100 {
			t.Errorf("Expected total 100, got %d", total)
		}

		expectedPercentage := float64(i)
		if percentage != expectedPercentage {
			t.Errorf("Expected percentage %.1f, got %.1f", expectedPercentage, percentage)
		}
	}

	tracker.Finish("Processing complete")

	// Test batch progress tracking
	batchTracker := utils.NewBatchProgressTracker(50, 10, "Batch processing")

	items := make([]string, 50)
	for i := range items {
		items[i] = "item" + string(rune('0'+i%10))
	}

	processedCount := 0
	for _, item := range items {
		batchTracker.AddToBatch(item)

		if len(batchTracker.CurrentBatch) == batchTracker.BatchSize {
			err := batchTracker.ProcessBatch(func(batch []interface{}) error {
				processedCount += len(batch)
				return nil
			})
			if err != nil {
				t.Errorf("Batch processing failed: %v", err)
			}
		}
	}

	// Process remaining items
	if len(batchTracker.CurrentBatch) > 0 {
		err := batchTracker.ProcessBatch(func(batch []interface{}) error {
			processedCount += len(batch)
			return nil
		})
		if err != nil {
			t.Errorf("Final batch processing failed: %v", err)
		}
	}

	if processedCount != 50 {
		t.Errorf("Expected 50 items processed, got %d", processedCount)
	}
}

func TestErrorHandlingIntegration(t *testing.T) {
	collector := utils.NewErrorCollector()

	// Simulate various error scenarios
	scenarios := []struct {
		errType utils.ErrorType
		message string
		context map[string]interface{}
	}{
		{utils.ErrInvalidDateFormat, "Invalid date provided", map[string]interface{}{"input": "bad-date"}},
		{utils.ErrNotGitRepository, "Not in a git repository", map[string]interface{}{"path": "/tmp"}},
		{utils.ErrRepositoryCorrupted, "Repository is corrupted", map[string]interface{}{"corruption": "missing objects"}},
	}

	for _, scenario := range scenarios {
		err := utils.NewGitStatsError(scenario.errType, scenario.message, nil)
		for key, value := range scenario.context {
			err.WithContext(key, value)
		}
		collector.AddError(err)
	}

	// Test error collection and summary
	if !collector.HasErrors() {
		t.Error("Expected collector to have errors")
	}

	if !collector.HasWarnings() {
		t.Error("Expected collector to have warnings")
	}

	summary := collector.GetSummary()
	if summary == "" {
		t.Error("Expected non-empty error summary")
	}

	// Test error wrapping
	originalErr := errors.New("underlying system error")
	wrappedErr := utils.WrapError(originalErr, utils.ErrFileNotFound, "Configuration file not found")

	if wrappedErr.Unwrap() != originalErr {
		t.Error("Expected wrapped error to unwrap to original error")
	}

	formatted := wrappedErr.FormatUserFriendlyError()
	if formatted == "" {
		t.Error("Expected non-empty formatted error message")
	}
}

func TestMultiStageProgressIntegration(t *testing.T) {
	stages := []utils.ProgressStage{
		{Name: "Initialization", Weight: 10, Total: 5},
		{Name: "Data Processing", Weight: 70, Total: 100},
		{Name: "Finalization", Weight: 20, Total: 10},
	}

	tracker := utils.NewMultiStageProgressTracker(stages)

	// Stage 1: Initialization
	for i := 1; i <= 5; i++ {
		tracker.UpdateStage(i, "Initializing...")
		time.Sleep(1 * time.Millisecond) // Brief pause to simulate work
	}
	tracker.NextStage()

	// Stage 2: Data Processing
	for i := 1; i <= 100; i += 10 {
		tracker.UpdateStage(i, "Processing data...")
		time.Sleep(1 * time.Millisecond) // Brief pause to simulate work
	}
	tracker.NextStage()

	// Stage 3: Finalization
	for i := 1; i <= 10; i++ {
		tracker.UpdateStage(i, "Finalizing...")
		time.Sleep(1 * time.Millisecond) // Brief pause to simulate work
	}

	tracker.Finish("All processing complete")
}

func TestUtilityFunctionsEdgeCases(t *testing.T) {
	// Test edge cases for date utilities
	_, err := utils.ParseDate("")
	if err == nil {
		t.Error("Expected error for empty date string")
	}

	_, err = utils.ParseRelativeDate("invalid relative date")
	if err == nil {
		t.Error("Expected error for invalid relative date")
	}

	// Test edge cases for progress tracking
	tracker := utils.NewProgressTracker(0, "Zero total")
	_, _, percentage := tracker.GetProgress()
	if percentage != 0.0 {
		t.Errorf("Expected 0%% for zero total, got %.1f%%", percentage)
	}

	// Test edge cases for error handling
	unknownErr := utils.NewGitStatsError(utils.ErrorType(999), "Unknown error type", nil)
	severity := unknownErr.GetSeverity()
	if severity != "UNKNOWN" {
		t.Errorf("Expected 'UNKNOWN' severity for unknown error type, got '%s'", severity)
	}

	suggestion := unknownErr.GetRecoverySuggestion()
	if suggestion == "" {
		t.Error("Expected non-empty suggestion even for unknown error types")
	}
}
