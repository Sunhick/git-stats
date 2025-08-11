// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for git repository interfaces

package git_test

import (
	"testing"
	"time"
	"git-stats/git"
)

// MockRepository implements the Repository interface for testing
type MockRepository struct {
	commits      []git.Commit
	contributors []git.Contributor
	branches     []string
	isValid      bool
	repoInfo     *git.RepositoryInfo
}

func (m *MockRepository) GetCommits(since, until time.Time, author string) ([]git.Commit, error) {
	return m.commits, nil
}

func (m *MockRepository) GetContributors() ([]git.Contributor, error) {
	return m.contributors, nil
}

func (m *MockRepository) GetBranches() ([]string, error) {
	return m.branches, nil
}

func (m *MockRepository) IsValidRepository() bool {
	return m.isValid
}

func (m *MockRepository) GetRepositoryInfo() (*git.RepositoryInfo, error) {
	return m.repoInfo, nil
}

func TestRepositoryInterface(t *testing.T) {
	// Test that MockRepository implements Repository interface
	var _ git.Repository = &MockRepository{}

	repo := &MockRepository{
		commits: []git.Commit{
			{
				Hash:    "abc123",
				Message: "Test commit",
				Author: git.Author{
					Name:  "Test User",
					Email: "test@example.com",
				},
				AuthorDate: time.Now(),
			},
		},
		contributors: []git.Contributor{
			{
				Name:         "Test User",
				Email:        "test@example.com",
				TotalCommits: 1,
			},
		},
		branches: []string{"main", "develop"},
		isValid:  true,
		repoInfo: &git.RepositoryInfo{
			Path:         "/test/repo",
			Name:         "test-repo",
			TotalCommits: 1,
		},
	}

	// Test GetCommits
	commits, err := repo.GetCommits(time.Now().AddDate(0, 0, -7), time.Now(), "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(commits) != 1 {
		t.Errorf("Expected 1 commit, got %d", len(commits))
	}

	// Test GetContributors
	contributors, err := repo.GetContributors()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(contributors) != 1 {
		t.Errorf("Expected 1 contributor, got %d", len(contributors))
	}

	// Test GetBranches
	branches, err := repo.GetBranches()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Errorf("Expected 2 branches, got %d", len(branches))
	}

	// Test IsValidRepository
	if !repo.IsValidRepository() {
		t.Error("Expected repository to be valid")
	}

	// Test GetRepositoryInfo
	info, err := repo.GetRepositoryInfo()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if info.Name != "test-repo" {
		t.Errorf("Expected repo name 'test-repo', got '%s'", info.Name)
	}
}

func TestCommitStruct(t *testing.T) {
	commit := git.Commit{
		Hash:    "abc123def456",
		Message: "Add new feature",
		Author: git.Author{
			Name:  "John Doe",
			Email: "john@example.com",
		},
		Committer: git.Author{
			Name:  "John Doe",
			Email: "john@example.com",
		},
		AuthorDate:    time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		CommitterDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		ParentHashes:  []string{"parent123"},
		TreeHash:      "tree456",
		Stats: git.CommitStats{
			FilesChanged: 2,
			Insertions:   10,
			Deletions:    5,
			Files: []git.FileChange{
				{
					Path:       "main.go",
					Status:     "M",
					Insertions: 8,
					Deletions:  3,
				},
				{
					Path:       "README.md",
					Status:     "M",
					Insertions: 2,
					Deletions:  2,
				},
			},
		},
	}

	if commit.Hash != "abc123def456" {
		t.Errorf("Expected hash 'abc123def456', got '%s'", commit.Hash)
	}

	if commit.Author.Name != "John Doe" {
		t.Errorf("Expected author name 'John Doe', got '%s'", commit.Author.Name)
	}

	if commit.Stats.FilesChanged != 2 {
		t.Errorf("Expected 2 files changed, got %d", commit.Stats.FilesChanged)
	}

	if len(commit.Stats.Files) != 2 {
		t.Errorf("Expected 2 file changes, got %d", len(commit.Stats.Files))
	}
}

func TestContributorStruct(t *testing.T) {
	contributor := git.Contributor{
		Name:            "Jane Smith",
		Email:           "jane@example.com",
		TotalCommits:    25,
		TotalInsertions: 500,
		TotalDeletions:  200,
		FirstCommit:     time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
		LastCommit:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		ActiveDays:      45,
		CommitsByDay:    map[string]int{"2024-01-15": 3, "2024-01-14": 2},
	}

	if contributor.Name != "Jane Smith" {
		t.Errorf("Expected name 'Jane Smith', got '%s'", contributor.Name)
	}

	if contributor.TotalCommits != 25 {
		t.Errorf("Expected 25 total commits, got %d", contributor.TotalCommits)
	}

	if contributor.CommitsByDay["2024-01-15"] != 3 {
		t.Errorf("Expected 3 commits on 2024-01-15, got %d", contributor.CommitsByDay["2024-01-15"])
	}
}
