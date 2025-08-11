// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Git repository integration tests

package git

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"git-stats/git"
)

func TestGitRepository_BasicFunctionality(t *testing.T) {
	// Skip if git is not available
	if !git.IsGitAvailable() {
		t.Skip("git not available in PATH")
	}

	// Create a temporary git repository for testing
	tempDir := createSimpleGitRepo(t)
	defer cleanupTempRepo(tempDir)

	// Create repository instance
	repo, err := git.NewGitRepository(git.RepositoryConfig{
		Path: tempDir,
	})
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Test IsValidRepository
	if !repo.IsValidRepository() {
		t.Error("IsValidRepository() returned false for valid repository")
	}

	// Test GetRepositoryInfo
	info, err := repo.GetRepositoryInfo()
	if err != nil {
		t.Errorf("GetRepositoryInfo() error = %v", err)
	} else {
		if info.Path != tempDir {
			t.Errorf("GetRepositoryInfo() path = %v, want %v", info.Path, tempDir)
		}
		if info.TotalCommits < 0 {
			t.Errorf("GetRepositoryInfo() total commits = %v, want >= 0", info.TotalCommits)
		}
	}

	// Test GetBranches
	branches, err := repo.GetBranches()
	if err != nil {
		t.Errorf("GetBranches() error = %v", err)
	} else if len(branches) == 0 {
		t.Log("GetBranches() returned no branches - this might be expected for empty repo")
	}

	// Test GetContributors
	contributors, err := repo.GetContributors()
	if err != nil {
		t.Errorf("GetContributors() error = %v", err)
	} else if len(contributors) == 0 {
		t.Log("GetContributors() returned no contributors - this might be expected for empty repo")
	}

	// Test GetCommits
	commits, err := repo.GetCommits(time.Time{}, time.Time{}, "")
	if err != nil {
		t.Errorf("GetCommits() error = %v", err)
	} else if len(commits) == 0 {
		t.Log("GetCommits() returned no commits - this might be expected for empty repo")
	}
}

func TestNewGitRepository_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		config  git.RepositoryConfig
		wantErr bool
	}{
		{
			name: "empty path",
			config: git.RepositoryConfig{
				Path: "",
			},
			wantErr: true,
		},
		{
			name: "non-existent path",
			config: git.RepositoryConfig{
				Path: "/nonexistent/path",
			},
			wantErr: true,
		},
		{
			name: "non-git directory",
			config: git.RepositoryConfig{
				Path: os.TempDir(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := git.NewGitRepository(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGitRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && repo == nil {
				t.Error("NewGitRepository() returned nil repository without error")
			}
		})
	}
}

// Helper functions for testing

func createSimpleGitRepo(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "git-stats-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Initialize git repository using system git command
	if git.IsGitAvailable() {
		// Create executor without working directory first
		executor, err := git.NewGitCommandExecutor(git.ExecutorConfig{})
		if err != nil {
			t.Fatalf("Failed to create executor: %v", err)
		}

		// Change to temp directory and initialize
		oldDir, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldDir)

		ctx := context.Background()

		// Initialize repository
		executor.Execute(ctx, "init")
		executor.Execute(ctx, "config", "user.name", "Test User")
		executor.Execute(ctx, "config", "user.email", "test@example.com")

		// Create initial commit if possible
		testFile := filepath.Join(tempDir, "README.md")
		if err := os.WriteFile(testFile, []byte("# Test Repository\n"), 0644); err == nil {
			executor.Execute(ctx, "add", "README.md")
			executor.Execute(ctx, "commit", "-m", "Initial commit")
		}
	}

	return tempDir
}

func cleanupTempRepo(path string) {
	os.RemoveAll(path)
}
