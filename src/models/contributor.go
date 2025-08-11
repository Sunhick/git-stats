// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Contributor data models

package models

import (
	"time"
)

// Contributor represents a comprehensive contributor profile
type Contributor struct {
	Name            string         `json:"name"`
	Email           string         `json:"email"`
	TotalCommits    int            `json:"total_commits"`
	TotalInsertions int            `json:"total_insertions"`
	TotalDeletions  int            `json:"total_deletions"`
	FirstCommit     time.Time      `json:"first_commit"`
	LastCommit      time.Time      `json:"last_commit"`
	ActiveDays      int            `json:"active_days"`
	CommitsByDay    map[string]int `json:"commits_by_day"` // date -> commit count
	CommitsByHour   map[int]int    `json:"commits_by_hour"` // hour -> commit count
	CommitsByWeekday map[int]int   `json:"commits_by_weekday"` // weekday -> commit count
	FileTypes       map[string]int `json:"file_types"` // extension -> commit count
	TopFiles        []string       `json:"top_files"` // most frequently modified files
}

// ContributorSummary provides a lightweight summary of contributor data
type ContributorSummary struct {
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Commits      int       `json:"commits"`
	Percentage   float64   `json:"percentage"`
	FirstCommit  time.Time `json:"first_commit"`
	LastCommit   time.Time `json:"last_commit"`
	IsActive     bool      `json:"is_active"` // active in last 3 months
}

// Validate checks if the contributor data is valid
func (c *Contributor) Validate() error {
	if c.Name == "" {
		return NewValidationError("contributor name cannot be empty")
	}
	if c.Email == "" {
		return NewValidationError("contributor email cannot be empty")
	}
	if c.TotalCommits < 0 {
		return NewValidationError("total commits cannot be negative")
	}
	if c.TotalInsertions < 0 {
		return NewValidationError("total insertions cannot be negative")
	}
	if c.TotalDeletions < 0 {
		return NewValidationError("total deletions cannot be negative")
	}
	if c.ActiveDays < 0 {
		return NewValidationError("active days cannot be negative")
	}
	if !c.FirstCommit.IsZero() && !c.LastCommit.IsZero() && c.FirstCommit.After(c.LastCommit) {
		return NewValidationError("first commit cannot be after last commit")
	}
	return nil
}

// GetActivityLevel returns the activity level based on commit count
func (c *Contributor) GetActivityLevel() string {
	if c.TotalCommits == 0 {
		return "inactive"
	} else if c.TotalCommits < 10 {
		return "low"
	} else if c.TotalCommits < 50 {
		return "medium"
	} else if c.TotalCommits < 200 {
		return "high"
	}
	return "very_high"
}

// GetContributionPeriod returns the duration of contribution activity
func (c *Contributor) GetContributionPeriod() time.Duration {
	if c.FirstCommit.IsZero() || c.LastCommit.IsZero() {
		return 0
	}
	return c.LastCommit.Sub(c.FirstCommit)
}

// GetAverageCommitsPerDay calculates average commits per active day
func (c *Contributor) GetAverageCommitsPerDay() float64 {
	if c.ActiveDays == 0 {
		return 0
	}
	return float64(c.TotalCommits) / float64(c.ActiveDays)
}

// IsActiveInPeriod checks if contributor was active in the given time range
func (c *Contributor) IsActiveInPeriod(start, end time.Time) bool {
	return !c.LastCommit.Before(start) && !c.FirstCommit.After(end)
}

// GetMostActiveHour returns the hour with most commits
func (c *Contributor) GetMostActiveHour() int {
	maxHour := 0
	maxCommits := 0
	for hour, commits := range c.CommitsByHour {
		if commits > maxCommits {
			maxCommits = commits
			maxHour = hour
		}
	}
	return maxHour
}

// GetMostActiveWeekday returns the weekday with most commits
func (c *Contributor) GetMostActiveWeekday() time.Weekday {
	maxWeekday := time.Sunday
	maxCommits := 0
	for weekday, commits := range c.CommitsByWeekday {
		if commits > maxCommits {
			maxCommits = commits
			maxWeekday = time.Weekday(weekday)
		}
	}
	return maxWeekday
}

// GetTopFileType returns the most frequently modified file type
func (c *Contributor) GetTopFileType() string {
	maxType := ""
	maxCommits := 0
	for fileType, commits := range c.FileTypes {
		if commits > maxCommits {
			maxCommits = commits
			maxType = fileType
		}
	}
	return maxType
}

// ToSummary converts a full Contributor to a ContributorSummary
func (c *Contributor) ToSummary(totalCommits int) ContributorSummary {
	percentage := 0.0
	if totalCommits > 0 {
		percentage = float64(c.TotalCommits) / float64(totalCommits) * 100
	}

	// Consider active if last commit was within 3 months
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	isActive := c.LastCommit.After(threeMonthsAgo)

	return ContributorSummary{
		Name:        c.Name,
		Email:       c.Email,
		Commits:     c.TotalCommits,
		Percentage:  percentage,
		FirstCommit: c.FirstCommit,
		LastCommit:  c.LastCommit,
		IsActive:    isActive,
	}
}

// Validate checks if the contributor summary data is valid
func (cs *ContributorSummary) Validate() error {
	if cs.Name == "" {
		return NewValidationError("contributor name cannot be empty")
	}
	if cs.Email == "" {
		return NewValidationError("contributor email cannot be empty")
	}
	if cs.Commits < 0 {
		return NewValidationError("commits cannot be negative")
	}
	if cs.Percentage < 0 || cs.Percentage > 100 {
		return NewValidationError("percentage must be between 0 and 100")
	}
	return nil
}
