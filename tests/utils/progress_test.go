// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for progress utilities

package utils_test

import (
	"git-stats/utils"
	"testing"
	"time"
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

func TestNewBatchProgressTracker(t *testing.T) {
	tracker := utils.NewBatchProgressTracker(100, 10, "Processing batches")

	if tracker.BatchSize != 10 {
		t.Errorf("Expected batch size 10, got %d", tracker.BatchSize)
	}

	if tracker.TotalBatches != 10 {
		t.Errorf("Expected total batches 10, got %d", tracker.TotalBatches)
	}

	if tracker.ProcessedBatches != 0 {
		t.Errorf("Expected processed batches 0, got %d", tracker.ProcessedBatches)
	}
}

func TestBatchProgressTrackerAddToBatch(t *testing.T) {
	tracker := utils.NewBatchProgressTracker(100, 10, "Processing")

	tracker.AddToBatch("item1")
	tracker.AddToBatch("item2")

	if len(tracker.CurrentBatch) != 2 {
		t.Errorf("Expected current batch size 2, got %d", len(tracker.CurrentBatch))
	}
}

func TestBatchProgressTrackerProcessBatch(t *testing.T) {
	tracker := utils.NewBatchProgressTracker(100, 10, "Processing")

	// Add items to batch
	tracker.AddToBatch("item1")
	tracker.AddToBatch("item2")
	tracker.AddToBatch("item3")

	// Process batch
	processed := 0
	err := tracker.ProcessBatch(func(items []interface{}) error {
		processed = len(items)
		return nil
	})

	if err != nil {
		t.Errorf("Unexpected error processing batch: %v", err)
	}

	if processed != 3 {
		t.Errorf("Expected 3 items processed, got %d", processed)
	}

	if tracker.ProcessedBatches != 1 {
		t.Errorf("Expected 1 processed batch, got %d", tracker.ProcessedBatches)
	}

	if tracker.Current != 3 {
		t.Errorf("Expected current progress 3, got %d", tracker.Current)
	}

	if len(tracker.CurrentBatch) != 0 {
		t.Errorf("Expected current batch to be cleared, got %d items", len(tracker.CurrentBatch))
	}
}

func TestNewMultiStageProgressTracker(t *testing.T) {
	stages := []utils.ProgressStage{
		{Name: "Stage 1", Weight: 30, Total: 100},
		{Name: "Stage 2", Weight: 50, Total: 200},
		{Name: "Stage 3", Weight: 20, Total: 50},
	}

	tracker := utils.NewMultiStageProgressTracker(stages)

	if len(tracker.GetStages()) != 3 {
		t.Errorf("Expected 3 stages, got %d", len(tracker.GetStages()))
	}

	if tracker.GetTotalWeight() != 100 {
		t.Errorf("Expected total weight 100, got %d", tracker.GetTotalWeight())
	}
}

func TestMultiStageProgressTrackerUpdateStage(t *testing.T) {
	stages := []utils.ProgressStage{
		{Name: "Stage 1", Weight: 50, Total: 100},
		{Name: "Stage 2", Weight: 50, Total: 100},
	}

	tracker := utils.NewMultiStageProgressTracker(stages)

	// Update first stage to 50% complete
	tracker.UpdateStage(50, "Processing stage 1")

	// Move to next stage
	tracker.NextStage()

	// Update second stage to 25% complete
	tracker.UpdateStage(25, "Processing stage 2")

	// Test that we can finish
	tracker.Finish("All stages complete")
}
