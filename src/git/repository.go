// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Git package for repository operations

package git

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
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

// GitRepository implements the Repository interface using git commands
type GitRepository struct {
	executor Executor
	parser   Parser
	path     string
}

// RepositoryConfig contains configuration for creating a GitRepository
type RepositoryConfig struct {
	Path     string
	Executor Executor
	Parser   Parser
}

// NewGitRepository creates a new GitRepository instance
func NewGitRepository(config RepositoryConfig) (*GitRepository, error) {
	if config.Path == "" {
		return nil, fmt.Errorf("repository path cannot be empty")
	}

	if config.Executor == nil {
		executorConfig := ExecutorConfig{
			WorkingDirectory: config.Path,
			DefaultTimeout:   30 * time.Second,
		}
		executor, err := NewGitCommandExecutor(executorConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create git executor: %w", err)
		}
		config.Executor = executor
	}

	if config.Parser == nil {
		config.Parser = NewGitOutputParser()
	}

	repo := &GitRepository{
		executor: config.Executor,
		parser:   config.Parser,
		path:     config.Path,
	}

	// Validate that this is a git repository
	if !repo.IsValidRepository() {
		return nil, fmt.Errorf("not a valid git repository: %s", config.Path)
	}

	return repo, nil
}

// GetCommits retrieves commits from the repository with optional filtering
func (r *GitRepository) GetCommits(since, until time.Time, author string) ([]Commit, error) {
	ctx := context.Background()

	// Build git log command arguments
	args := []string{
		"--pretty=format:%H|%an|%ae|%ad|%cn|%ce|%cd|%s|%P|%T",
		"--date=iso",
		"--numstat",
		"--all",
	}

	// Add date filters
	if !since.IsZero() {
		args = append(args, "--since="+since.Format("2006-01-02"))
	}
	if !until.IsZero() {
		args = append(args, "--until="+until.Format("2006-01-02"))
	}

	// Add author filter
	if author != "" {
		args = append(args, "--author="+author)
	}

	// Execute git log command
	result, err := r.executor.Execute(ctx, "log", args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute git log: %w", err)
	}

	// Parse the output
	commits, err := r.parser.ParseCommitLog(result.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commit log: %w", err)
	}

	return commits, nil
}

// GetContributors retrieves contributor information from the repository
func (r *GitRepository) GetContributors() ([]Contributor, error) {
	ctx := context.Background()

	// Get contributor summary using shortlog
	result, err := r.executor.Execute(ctx, "shortlog", "-sne", "--all")
	if err != nil {
		return nil, fmt.Errorf("failed to execute git shortlog: %w", err)
	}

	// Parse contributors
	contributors, err := r.parser.ParseContributors(result.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse contributors: %w", err)
	}

	// Enhance contributor data with detailed statistics
	for i := range contributors {
		if err := r.enhanceContributorData(&contributors[i]); err != nil {
			return nil, fmt.Errorf("failed to enhance contributor data for %s: %w", contributors[i].Name, err)
		}
	}

	return contributors, nil
}

// GetBranches retrieves all branches from the repository
func (r *GitRepository) GetBranches() ([]string, error) {
	ctx := context.Background()

	// Get all branches
	result, err := r.executor.Execute(ctx, "branch", "-a")
	if err != nil {
		return nil, fmt.Errorf("failed to execute git branch: %w", err)
	}

	// Parse branches
	branches, err := r.parser.ParseBranches(result.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse branches: %w", err)
	}

	return branches, nil
}

// IsValidRepository checks if the path contains a valid git repository
func (r *GitRepository) IsValidRepository() bool {
	ctx := context.Background()

	// Try to execute a simple git command
	_, err := r.executor.Execute(ctx, "rev-parse", "--git-dir")
	return err == nil
}

// GetRepositoryInfo retrieves comprehensive repository metadata
func (r *GitRepository) GetRepositoryInfo() (*RepositoryInfo, error) {
	ctx := context.Background()

	info := &RepositoryInfo{
		Path: r.path,
		Name: filepath.Base(r.path),
	}

	// Get total commit count (handle empty repositories)
	result, err := r.executor.Execute(ctx, "rev-list", "--count", "HEAD")
	if err != nil {
		// If HEAD doesn't exist, this is likely an empty repository
		info.TotalCommits = 0
	} else {
		totalCommits, err := strconv.Atoi(strings.TrimSpace(result.Output))
		if err != nil {
			return nil, fmt.Errorf("failed to parse commit count: %w", err)
		}
		info.TotalCommits = totalCommits
	}

	// Get first commit date
	result, err = r.executor.Execute(ctx, "log", "--reverse", "--pretty=format:%ad", "--date=iso", "-1")
	if err == nil && result.Output != "" {
		firstCommit, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(result.Output))
		if err == nil {
			info.FirstCommit = firstCommit
		}
	}

	// Get last commit date
	result, err = r.executor.Execute(ctx, "log", "--pretty=format:%ad", "--date=iso", "-1")
	if err == nil && result.Output != "" {
		lastCommit, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(result.Output))
		if err == nil {
			info.LastCommit = lastCommit
		}
	}

	// Get branches
	branches, err := r.GetBranches()
	if err == nil {
		info.Branches = branches
	}

	return info, nil
}

// enhanceContributorData adds detailed statistics to a contributor
func (r *GitRepository) enhanceContributorData(contributor *Contributor) error {
	ctx := context.Background()

	// Get detailed commit information for this contributor
	args := []string{
		"--pretty=format:%ad|%H",
		"--date=short",
		"--numstat",
		"--author=" + contributor.Email,
		"--all",
	}

	result, err := r.executor.Execute(ctx, "log", args...)
	if err != nil {
		return fmt.Errorf("failed to get contributor commits: %w", err)
	}

	// Parse and enhance contributor data
	if err := r.parseContributorCommits(contributor, result.Output); err != nil {
		return fmt.Errorf("failed to parse contributor commits: %w", err)
	}

	return nil
}

// parseContributorCommits parses commit data to enhance contributor statistics
func (r *GitRepository) parseContributorCommits(contributor *Contributor, output string) error {
	if contributor.CommitsByDay == nil {
		contributor.CommitsByDay = make(map[string]int)
	}

	lines := strings.Split(output, "\n")
	var currentDate string
	var insertions, deletions int

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this is a commit header line (date|hash)
		if strings.Contains(line, "|") && len(strings.Split(line, "|")) == 2 {
			parts := strings.Split(line, "|")
			currentDate = parts[0]

			// Update commits by day
			contributor.CommitsByDay[currentDate]++

			// Parse commit date for first/last commit tracking
			commitDate, err := time.Parse("2006-01-02", currentDate)
			if err == nil {
				if contributor.FirstCommit.IsZero() || commitDate.Before(contributor.FirstCommit) {
					contributor.FirstCommit = commitDate
				}
				if contributor.LastCommit.IsZero() || commitDate.After(contributor.LastCommit) {
					contributor.LastCommit = commitDate
				}
			}
		} else {
			// This should be a numstat line (insertions\tdeletions\tfilename)
			parts := strings.Split(line, "\t")
			if len(parts) >= 3 {
				if ins, err := strconv.Atoi(parts[0]); err == nil {
					insertions += ins
				}
				if dels, err := strconv.Atoi(parts[1]); err == nil {
					deletions += dels
				}
			}
		}
	}

	// Update totals
	contributor.TotalInsertions = insertions
	contributor.TotalDeletions = deletions
	contributor.ActiveDays = len(contributor.CommitsByDay)

	return nil
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
