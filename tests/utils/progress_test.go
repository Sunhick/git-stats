// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for progress utilities

package utils_test

import (
	"testing"
	"time"
	"git-stats/utils"
)

func TestNewProgressTracker(t *testing.T) {
	tracker := utils.NewProgressTracker(100, "Testing progress")

	if tracker.Total != 100 {
		t.Errorf("Expected total 100, got %d", tracker.Total)
	}

	if tracker.Current != 0 {
		t.Errorf("Expected current 0, got %d", tracker.Current)
	}

	if tracker.Message != "Testing progress" {
		t.Errorf("Expected message 'Testing progress', got '%s'", tracker.Message)
	}
}

func TestProgressTrackerStart(t *testing.T) {
	tracker := utils.NewProgressTracker(0, "")
	tracker.Start(50, "Starting test")

	if tracker.Total != 50 {
		t.Errorf("Expected total 50, got %d", tracker.Total)
	}

	if tracker.Current != 0 {
		t.Errorf("Expected current 0, got %d", tracker.Current)
	}

	if tracker.Message != "Starting test" {
		t.Errorf("Expected message 'Starting test', got '%s'", tracker.Message)
	}
}

func TestProgressTrackerUpdate(t *testing.T) {
	tracker := utils.NewProgressTracker(100, "Test")
	tracker.Update(25, "Updated message")

	if tracker.Current != 25 {
		t.Errorf("Expected current 25, got %d", tracker.Current)
	}

	if tracker.Message != "Updated message" {
		t.Errorf("Expected message 'Updated message', got '%s'", tracker.Message)
	}
}

func TestProgressTrackerIncrement(t *testing.T) {
	tracker := utils.NewProgressTracker(100, "Test")

	tracker.Increment("First increment")
	if tracker.Current != 1 {
		t.Errorf("Expected current 1, got %d", tracker.Current)
	}

	tracker.Increment("Second increment")
	if tracker.Current != 2 {
		t.Errorf("Expected current 2, got %d", tracker.Current)
	}
}

func TestProgressTrackerGetProgress(t *testing.T) {
	tracker := utils.NewProgressTracker(100, "Test")
	tracker.Update(25, "")

	current, total, percentage := tracker.GetProgress()

	if current != 25 {
		t.Errorf("Expected current 25, got %d", current)
	}

	if total != 100 {
		t.Errorf("Expected total 100, got %d", total)
	}

	if percentage != 25.0 {
		t.Errorf("Expected percentage 25.0, got %f", percentage)
	}
}

func TestProgressTrackerGetProgressZeroTotal(t *testing.T) {
	tracker := utils.NewProgressTracker(0, "Test")
	tracker.Update(5, "")

	current, total, percentage := tracker.GetProgress()

	if current != 5 {
		t.Errorf("Expected current 5, got %d", current)
	}

	if total != 0 {
		t.Errorf("Expected total 0, got %d", total)
	}

	if percentage != 0.0 {
		t.Errorf("Expected percentage 0.0, got %f", percentage)
	}
}

func TestNewSimpleSpinner(t *testing.T) {
	spinner := utils.NewSimpleSpinner("Loading...")

	if spinner == nil {
		t.Error("Expected spinner to be created, got nil")
	}

	// Test that spinner can be started and stopped without panic
	spinner.Start()
	time.Sleep(10 * time.Millisecond) // Brief pause to let spinner run
	spinner.Stop()

	// Test updating message
	spinner.UpdateMessage("New message")
}
