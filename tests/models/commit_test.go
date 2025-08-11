// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Commit model tests

package models

import (
	"testing"
	"time"

	"../src/models"
)

func TestCommit_Validate(t *testing.T) {
	tests := []struct {
		name    string
		commit  models.Commit
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid commit",
			commit: models.Commit{
				Hash:       "abc123def456",
				Message:    "test commit",
				Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
				Committer:  models.Author{Name: "John Doe", Email: "john@example.com"},
				AuthorDate: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty hash",
			commit: models.Commit{
				Hash:       "",
				Message:    "test commit",
				Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
				AuthorDate: time.Now(),
			},
			wantErr: true,
			errMsg:  "commit hash cannot be empty",
		},
		{
			name: "short hash",
			commit: models.Commit{
				Hash:       "abc123",
				Message:    "test commit",
				Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
				AuthorDate: time.Now(),
			},
			wantErr: true,
			errMsg:  "commit hash must be at least 7 characters",
		},
		{
			name: "empty author name",
			commit: models.Commit{
				Hash:       "abc123def456",
				Message:    "test commit",
				Author:     models.Author{Name: "", Email: "john@example.com"},
				AuthorDate: time.Now(),
			},
			wantErr: true,
			errMsg:  "author name cannot be empty",
		},
		{
			name: "empty author email",
			commit: models.Commit{
				Hash:       "abc123def456",
				Message:    "test commit",
				Author:     models.Author{Name: "John Doe", Email: ""},
				AuthorDate: time.Now(),
			},
			wantErr: true,
			errMsg:  "author email cannot be empty",
		},
		{
			name: "zero author date",
			commit: models.Commit{
				Hash:       "abc123def456",
				Message:    "test commit",
				Author:     models.Author{Name: "John Doe", Email: "john@example.com"},
				AuthorDate: time.Time{},
			},
			wantErr: true,
			errMsg:  "author date cannot be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.commit.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Commit.Validate() expected error but got none")
					return
				}
				if err.Error() != "validation error: "+tt.errMsg {
					t.Errorf("Commit.Validate() error = %v, want %v", err.Error(), "validation error: "+tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Commit.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestCommit_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		commit models.Commit
		want   bool
	}{
		{
			name: "empty commit",
			commit: models.Commit{
				Stats: models.CommitStats{
					FilesChanged: 0,
					Insertions:   0,
					Deletions:    0,
				},
			},
			want: true,
		},
		{
			name: "non-empty commit",
			commit: models.Commit{
				Stats: models.CommitStats{
					FilesChanged: 1,
					Insertions:   10,
					Deletions:    5,
				},
			},
			want: false,
		},
		{
			name: "only files changed",
			commit: models.Commit{
				Stats: models.CommitStats{
					FilesChanged: 1,
					Insertions:   0,
					Deletions:    0,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.commit.IsEmpty(); got != tt.want {
				t.Errorf("Commit.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommit_IsMergeCommit(t *testing.T) {
	tests := []struct {
		name   string
		commit models.Commit
		want   bool
	}{
		{
			name: "regular commit",
			commit: models.Commit{
				ParentHashes: []string{"abc123"},
			},
			want: false,
		},
		{
			name: "merge commit",
			commit: models.Commit{
				ParentHashes: []string{"abc123", "def456"},
			},
			want: true,
		},
		{
			name: "no parents",
			commit: models.Commit{
				ParentHashes: []string{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.commit.IsMergeCommit(); got != tt.want {
				t.Errorf("Commit.IsMergeCommit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommit_GetFileExtensions(t *testing.T) {
	commit := models.Commit{
		Stats: models.CommitStats{
			Files: []models.FileChange{
				{Path: "main.go"},
				{Path: "test.go"},
				{Path: "README.md"},
				{Path: "config.json"},
				{Path: "script.sh"},
				{Path: "another.go"},
				{Path: "noextension"},
			},
		},
	}

	extensions := commit.GetFileExtensions()

	expectedExtensions := map[string]bool{
		"go":   true,
		"md":   true,
		"json": true,
		"sh":   true,
	}

	if len(extensions) != len(expectedExtensions) {
		t.Errorf("GetFileExtensions() returned %d extensions, want %d", len(extensions), len(expectedExtensions))
	}

	for _, ext := range extensions {
		if !expectedExtensions[ext] {
			t.Errorf("GetFileExtensions() returned unexpected extension: %s", ext)
		}
	}
}

func TestAuthor_Validate(t *testing.T) {
	tests := []struct {
		name    string
		author  models.Author
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid author",
			author:  models.Author{Name: "John Doe", Email: "john@example.com"},
			wantErr: false,
		},
		{
			name:    "empty name",
			author:  models.Author{Name: "", Email: "john@example.com"},
			wantErr: true,
			errMsg:  "author name cannot be empty",
		},
		{
			name:    "empty email",
			author:  models.Author{Name: "John Doe", Email: ""},
			wantErr: true,
			errMsg:  "author email cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.author.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Author.Validate() expected error but got none")
					return
				}
				if err.Error() != "validation error: "+tt.errMsg {
					t.Errorf("Author.Validate() error = %v, want %v", err.Error(), "validation error: "+tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Author.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestFileChange_Validate(t *testing.T) {
	tests := []struct {
		name       string
		fileChange models.FileChange
		wantErr    bool
		errMsg     string
	}{
		{
			name: "valid file change",
			fileChange: models.FileChange{
				Path:       "main.go",
				Status:     "M",
				Insertions: 10,
				Deletions:  5,
			},
			wantErr: false,
		},
		{
			name: "empty path",
			fileChange: models.FileChange{
				Path:   "",
				Status: "M",
			},
			wantErr: true,
			errMsg:  "file path cannot be empty",
		},
		{
			name: "empty status",
			fileChange: models.FileChange{
				Path:   "main.go",
				Status: "",
			},
			wantErr: true,
			errMsg:  "file status cannot be empty",
		},
		{
			name: "invalid status",
			fileChange: models.FileChange{
				Path:   "main.go",
				Status: "X",
			},
			wantErr: true,
			errMsg:  "invalid file status: X",
		},
		{
			name: "negative insertions",
			fileChange: models.FileChange{
				Path:       "main.go",
				Status:     "M",
				Insertions: -1,
			},
			wantErr: true,
			errMsg:  "insertions and deletions cannot be negative",
		},
		{
			name: "negative deletions",
			fileChange: models.FileChange{
				Path:      "main.go",
				Status:    "M",
				Deletions: -1,
			},
			wantErr: true,
			errMsg:  "insertions and deletions cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fileChange.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("FileChange.Validate() expected error but got none")
					return
				}
				if err.Error() != "validation error: "+tt.errMsg {
					t.Errorf("FileChange.Validate() error = %v, want %v", err.Error(), "validation error: "+tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("FileChange.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestFileChange_StatusMethods(t *testing.T) {
	tests := []struct {
		status     string
		isAdded    bool
		isModified bool
		isDeleted  bool
		isRenamed  bool
		isCopied   bool
	}{
		{"A", true, false, false, false, false},
		{"M", false, true, false, false, false},
		{"D", false, false, true, false, false},
		{"R", false, false, false, true, false},
		{"C", false, false, false, false, true},
	}

	for _, tt := range tests {
		t.Run("status_"+tt.status, func(t *testing.T) {
			fc := models.FileChange{Status: tt.status}

			if got := fc.IsAdded(); got != tt.isAdded {
				t.Errorf("FileChange.IsAdded() = %v, want %v", got, tt.isAdded)
			}
			if got := fc.IsModified(); got != tt.isModified {
				t.Errorf("FileChange.IsModified() = %v, want %v", got, tt.isModified)
			}
			if got := fc.IsDeleted(); got != tt.isDeleted {
				t.Errorf("FileChange.IsDeleted() = %v, want %v", got, tt.isDeleted)
			}
			if got := fc.IsRenamed(); got != tt.isRenamed {
				t.Errorf("FileChange.IsRenamed() = %v, want %v", got, tt.isRenamed)
			}
			if got := fc.IsCopied(); got != tt.isCopied {
				t.Errorf("FileChange.IsCopied() = %v, want %v", got, tt.isCopied)
			}
		})
	}
}
