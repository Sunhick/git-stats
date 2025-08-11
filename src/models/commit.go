// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Enhanced commit data models

package models

import (
	"time"
)

// Commit represents an enhanced git commit with comprehensive statistics
type Commit struct {
	Hash          string      `json:"hash"`
	Message       string      `json:"message"`
	Author        Author      `json:"author"`
	Committer     Author      `json:"committer"`
	AuthorDate    time.Time   `json:"author_date"`
	CommitterDate time.Time   `json:"committer_date"`
	ParentHashes  []string    `json:"parent_hashes"`
	TreeHash      string      `json:"tree_hash"`
	Stats         CommitStats `json:"stats"`
}

// Author represents commit author or committer information
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CommitStats contains detailed statistics for a commit
type CommitStats struct {
	FilesChanged int          `json:"files_changed"`
	Insertions   int          `json:"insertions"`
	Deletions    int          `json:"deletions"`
	Files        []FileChange `json:"files"`
}

// FileChange represents changes to a specific file in a commit
type FileChange struct {
	Path       string `json:"path"`
	Status     string `json:"status"` // A, M, D, R, C (Added, Modified, Deleted, Renamed, Copied)
	Insertions int    `json:"insertions"`
	Deletions  int    `json:"deletions"`
	OldPath    string `json:"old_path,omitempty"` // For renamed/copied files
}

// Validate checks if the commit data is valid
func (c *Commit) Validate() error {
	if c.Hash == "" {
		return NewValidationError("commit hash cannot be empty")
	}
	if len(c.Hash) < 7 {
		return NewValidationError("commit hash must be at least 7 characters")
	}
	if c.Author.Name == "" {
		return NewValidationError("author name cannot be empty")
	}
	if c.Author.Email == "" {
		return NewValidationError("author email cannot be empty")
	}
	if c.AuthorDate.IsZero() {
		return NewValidationError("author date cannot be zero")
	}
	return nil
}

// IsEmpty returns true if the commit is empty (no changes)
func (c *Commit) IsEmpty() bool {
	return c.Stats.FilesChanged == 0 && c.Stats.Insertions == 0 && c.Stats.Deletions == 0
}

// IsMergeCommit returns true if this is a merge commit
func (c *Commit) IsMergeCommit() bool {
	return len(c.ParentHashes) > 1
}

// GetFileExtensions returns unique file extensions modified in this commit
func (c *Commit) GetFileExtensions() []string {
	extensions := make(map[string]bool)
	for _, file := range c.Stats.Files {
		if ext := getFileExtension(file.Path); ext != "" {
			extensions[ext] = true
		}
	}

	result := make([]string, 0, len(extensions))
	for ext := range extensions {
		result = append(result, ext)
	}
	return result
}

// Validate checks if the author data is valid
func (a *Author) Validate() error {
	if a.Name == "" {
		return NewValidationError("author name cannot be empty")
	}
	if a.Email == "" {
		return NewValidationError("author email cannot be empty")
	}
	return nil
}

// Validate checks if the file change data is valid
func (fc *FileChange) Validate() error {
	if fc.Path == "" {
		return NewValidationError("file path cannot be empty")
	}
	if fc.Status == "" {
		return NewValidationError("file status cannot be empty")
	}
	validStatuses := map[string]bool{"A": true, "M": true, "D": true, "R": true, "C": true}
	if !validStatuses[fc.Status] {
		return NewValidationError("invalid file status: " + fc.Status)
	}
	if fc.Insertions < 0 || fc.Deletions < 0 {
		return NewValidationError("insertions and deletions cannot be negative")
	}
	return nil
}

// IsAdded returns true if the file was added
func (fc *FileChange) IsAdded() bool {
	return fc.Status == "A"
}

// IsModified returns true if the file was modified
func (fc *FileChange) IsModified() bool {
	return fc.Status == "M"
}

// IsDeleted returns true if the file was deleted
func (fc *FileChange) IsDeleted() bool {
	return fc.Status == "D"
}

// IsRenamed returns true if the file was renamed
func (fc *FileChange) IsRenamed() bool {
	return fc.Status == "R"
}

// IsCopied returns true if the file was copied
func (fc *FileChange) IsCopied() bool {
	return fc.Status == "C"
}

// getFileExtension extracts file extension from path
func getFileExtension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
		if path[i] == '/' {
			break
		}
	}
	return ""
}
