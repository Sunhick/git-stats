// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Utils package for progress tracking utilities

package utils

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// ProgressTracker tracks progress of long-running operations
type ProgressTracker struct {
	Total     int
	Current   int
	StartTime time.Time
	Message   string
	Width     int
	ShowETA   bool
}

// ProgressReporter interface for progress reporting
type ProgressReporter interface {
	Start(total int, message string)
	Update(current int, message string)
	Increment(message string)
	Finish(message string)
	SetTotal(total int)
	GetProgress() (current, total int, percentage float64)
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(total int, message string) *ProgressTracker {
	return &ProgressTracker{
		Total:     total,
		Current:   0,
		StartTime: time.Now(),
		Message:   message,
		Width:     50,
		ShowETA:   true,
	}
}

// Start initializes the progress tracker
func (p *ProgressTracker) Start(total int, message string) {
	p.Total = total
	p.Current = 0
	p.StartTime = time.Now()
	p.Message = message
	p.render()
}

// Update updates the current progress
func (p *ProgressTracker) Update(current int, message string) {
	p.Current = current
	if message != "" {
		p.Message = message
	}
	p.render()
}

// Increment increments the progress by 1
func (p *ProgressTracker) Increment(message string) {
	p.Current++
	if message != "" {
		p.Message = message
	}
	p.render()
}

// Finish completes the progress tracking
func (p *ProgressTracker) Finish(message string) {
	p.Current = p.Total
	if message != "" {
		p.Message = message
	}
	p.render()
	fmt.Println() // New line after completion
}

// SetTotal updates the total count
func (p *ProgressTracker) SetTotal(total int) {
	p.Total = total
}

// GetProgress returns current progress information
func (p *ProgressTracker) GetProgress() (current, total int, percentage float64) {
	percentage = 0.0
	if p.Total > 0 {
		percentage = float64(p.Current) / float64(p.Total) * 100
	}
	return p.Current, p.Total, percentage
}

// render displays the progress bar
func (p *ProgressTracker) render() {
	if p.Total <= 0 {
		fmt.Printf("\r%s... %d processed", p.Message, p.Current)
		return
	}

	percentage := float64(p.Current) / float64(p.Total)
	filled := int(percentage * float64(p.Width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.Width-filled)

	elapsed := time.Since(p.StartTime)
	eta := ""

	if p.ShowETA && p.Current > 0 && p.Current < p.Total {
		remaining := time.Duration(float64(elapsed) * (float64(p.Total)/float64(p.Current) - 1))
		eta = fmt.Sprintf(" ETA: %s", formatDuration(remaining))
	}

	fmt.Printf("\r%s [%s] %d/%d (%.1f%%)%s",
		p.Message, bar, p.Current, p.Total, percentage*100, eta)
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) - 60*minutes
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	} else {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) - 60*hours
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
}

// SimpleSpinner provides a simple spinning indicator for indeterminate progress
type SimpleSpinner struct {
	chars   []rune
	current int
	message string
	active  bool
	mu      sync.Mutex
}

// NewSimpleSpinner creates a new spinner
func NewSimpleSpinner(message string) *SimpleSpinner {
	return &SimpleSpinner{
		chars:   []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'},
		current: 0,
		message: message,
		active:  false,
	}
}

// Start begins the spinner animation
func (s *SimpleSpinner) Start() {
	s.mu.Lock()
	s.active = true
	s.mu.Unlock()

	go func() {
		for {
			s.mu.Lock()
			if !s.active {
				s.mu.Unlock()
				break
			}
			fmt.Printf("\r%c %s", s.chars[s.current], s.message)
			s.current = (s.current + 1) % len(s.chars)
			s.mu.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// Stop stops the spinner animation
func (s *SimpleSpinner) Stop() {
	s.mu.Lock()
	s.active = false
	s.mu.Unlock()
	fmt.Print("\r") // Clear the line
}

// UpdateMessage updates the spinner message
func (s *SimpleSpinner) UpdateMessage(message string) {
	s.mu.Lock()
	s.message = message
	s.mu.Unlock()
}

// BatchProgressTracker handles progress tracking for batch operations
type BatchProgressTracker struct {
	*ProgressTracker
	BatchSize        int
	ProcessedBatches int
	TotalBatches     int
	CurrentBatch     []interface{}
}

// NewBatchProgressTracker creates a new batch progress tracker
func NewBatchProgressTracker(total int, batchSize int, message string) *BatchProgressTracker {
	totalBatches := (total + batchSize - 1) / batchSize // Ceiling division
	return &BatchProgressTracker{
		ProgressTracker:  NewProgressTracker(total, message),
		BatchSize:        batchSize,
		ProcessedBatches: 0,
		TotalBatches:     totalBatches,
		CurrentBatch:     make([]interface{}, 0, batchSize),
	}
}

// AddToBatch adds an item to the current batch
func (b *BatchProgressTracker) AddToBatch(item interface{}) {
	b.CurrentBatch = append(b.CurrentBatch, item)
}

// ProcessBatch processes the current batch and updates progress
func (b *BatchProgressTracker) ProcessBatch(processor func([]interface{}) error) error {
	if len(b.CurrentBatch) == 0 {
		return nil
	}

	err := processor(b.CurrentBatch)
	if err != nil {
		return err
	}

	b.ProcessedBatches++
	b.Current += len(b.CurrentBatch)
	b.CurrentBatch = b.CurrentBatch[:0] // Clear batch

	// Update progress message
	batchMessage := fmt.Sprintf("%s (batch %d/%d)", b.Message, b.ProcessedBatches, b.TotalBatches)
	b.Update(b.Current, batchMessage)

	return nil
}

// MultiStageProgressTracker handles progress across multiple stages
type MultiStageProgressTracker struct {
	stages       []ProgressStage
	currentStage int
	totalWeight  int
	completed    int
}

// ProgressStage represents a stage in multi-stage processing
type ProgressStage struct {
	Name    string
	Weight  int // Relative weight of this stage
	Total   int // Total items in this stage
	Current int // Current progress in this stage
}

// NewMultiStageProgressTracker creates a new multi-stage progress tracker
func NewMultiStageProgressTracker(stages []ProgressStage) *MultiStageProgressTracker {
	totalWeight := 0
	for _, stage := range stages {
		totalWeight += stage.Weight
	}

	return &MultiStageProgressTracker{
		stages:       stages,
		currentStage: 0,
		totalWeight:  totalWeight,
		completed:    0,
	}
}

// UpdateStage updates progress for the current stage
func (m *MultiStageProgressTracker) UpdateStage(current int, message string) {
	if m.currentStage >= len(m.stages) {
		return
	}

	m.stages[m.currentStage].Current = current

	// Calculate overall progress
	overallProgress := 0
	for i, stage := range m.stages {
		if i < m.currentStage {
			overallProgress += stage.Weight
		} else if i == m.currentStage {
			stageProgress := 0
			if stage.Total > 0 {
				stageProgress = (stage.Current * stage.Weight) / stage.Total
			}
			overallProgress += stageProgress
		}
	}

	percentage := float64(overallProgress) / float64(m.totalWeight) * 100

	stageMessage := fmt.Sprintf("[%s] %s", m.stages[m.currentStage].Name, message)
	fmt.Printf("\r%s (%.1f%% complete)", stageMessage, percentage)
}

// NextStage moves to the next stage
func (m *MultiStageProgressTracker) NextStage() {
	if m.currentStage < len(m.stages)-1 {
		m.currentStage++
	}
}

// Finish completes all stages
func (m *MultiStageProgressTracker) Finish(message string) {
	fmt.Printf("\r%s (100%% complete)\n", message)
}

// GetStages returns the stages
func (m *MultiStageProgressTracker) GetStages() []ProgressStage {
	return m.stages
}

// GetTotalWeight returns the total weight
func (m *MultiStageProgressTracker) GetTotalWeight() int {
	return m.totalWeight
}
