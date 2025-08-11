// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Git package for repository operations

package git

import (
	"time"
)

// RepositoryInfo contains metadata about the git repository
type RepositoryInfo struct {
	Path         string
	Name         string
	TotalCommits int
	FirstCommit  time.Time
	LastCommit   time.Time
	Branches     []string
}

// Repository interface for git operations
type Repository interface {
	GetCommits(since, until time.Time, author string) ([]Commit, error)
	GetContributors() ([]Contributor, error)
	GetBranches() ([]string, error)
	IsValidRepository() bool
	GetRepositoryInfo() (*RepositoryInfo, error)
}

// Commit represents a git commit with enhanced statistics
type Commit struct {
	Hash          string
	Message       string
	Author        Author
	Committer     Author
	AuthorDate    time.Time
	CommitterDate time.Time
	ParentHashes  []string
	TreeHash      string
	Stats         CommitStats
}

// Author represents commit author information
type Author struct {
	Name  string
	Email string
}

// CommitStats contains statistics about changes in a commit
type CommitStats struct {
	FilesChanged int
	Insertions   int
	Deletions    int
	Files        []FileChange
}

// FileChange represents changes to a specific file
type FileChange struct {
	Path       string
	Status     string // A, M, D, R, C
	Insertions int
	Deletions  int
}

// Contributor represents a repository contributor
type Contributor struct {
	Name            string
	Email           string
	TotalCommits    int
	TotalInsertions int
	TotalDeletions  int
	FirstCommit     time.Time
	LastCommit      time.Time
	ActiveDays      int
	CommitsByDay    map[string]int
}
