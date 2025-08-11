// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Utils package for progress tracking utilities

package utils

import (
	"fmt"
	"strings"
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
		remaining := time.Duration(float64(elapsed) * (float64(p.Total) / float64(p.Current) - 1))
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
	s.active = true
	go func() {
		for s.active {
			fmt.Printf("\r%c %s", s.chars[s.current], s.message)
			s.current = (s.current + 1) % len(s.chars)
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// Stop stops the spinner animation
func (s *SimpleSpinner) Stop() {
	s.active = false
	fmt.Print("\r") // Clear the line
}

// UpdateMessage updates the spinner message
func (s *SimpleSpinner) UpdateMessage(message string) {
	s.message = message
}
