// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Git package for output parsing

package git

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parser interface for parsing git command output
type Parser interface {
	ParseCommitLog(output string) ([]Commit, error)
	ParseDiffStat(output string) (*CommitStats, error)
	ParseContributors(output string) ([]Contributor, error)
	ParseBranches(output string) ([]string, error)
}

// GitOutputParser implements the Parser interface for git command output
type GitOutputParser struct{}

// NewGitOutputParser creates a new GitOutputParser instance
func NewGitOutputParser() *GitOutputParser {
	return &GitOutputParser{}
}

// ParseCommitLog parses git log output with --numstat format
func (p *GitOutputParser) ParseCommitLog(output string) ([]Commit, error) {
	if output == "" {
		return []Commit{}, nil
	}

	var commits []Commit
	lines := strings.Split(output, "\n")

	var currentCommit *Commit

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Check if this is a commit header line (contains multiple | separators)
		if strings.Count(line, "|") >= 7 {
			// Save previous commit if exists
			if currentCommit != nil {
				commits = append(commits, *currentCommit)
			}

			commit, err := p.parseCommitHeader(line)
			if err != nil {
				// Skip malformed commit headers but continue processing
				continue
			}
			currentCommit = commit
		} else if currentCommit != nil && line != "" {
			// This should be a numstat line (skip empty lines)
			if err := p.parseNumStatLine(currentCommit, line); err != nil {
				// Skip invalid numstat lines but continue processing
				continue
			}
		} else if line == "" && currentCommit != nil {
			// Empty line - check if this is the end of the commit
			// Look ahead to see if there are more numstat lines
			hasMoreStats := false
			for j := i + 1; j < len(lines); j++ {
				nextLine := strings.TrimSpace(lines[j])
				if nextLine == "" {
					continue
				}
				// If next non-empty line is a commit header, we're done
				if strings.Count(nextLine, "|") >= 7 {
					break
				}
				// If next non-empty line looks like numstat, continue
				if strings.Contains(nextLine, "\t") {
					hasMoreStats = true
					break
				}
				break
			}

			// If no more stats and we're at the end or next commit, finalize this commit
			if !hasMoreStats {
				commits = append(commits, *currentCommit)
				currentCommit = nil
			}
		}
	}

	// Add the last commit if exists
	if currentCommit != nil {
		commits = append(commits, *currentCommit)
	}

	return commits, nil
}

// parseCommitHeader parses a commit header line in format: hash|author|email|date|committer|cemail|cdate|message|parents|tree
func (p *GitOutputParser) parseCommitHeader(line string) (*Commit, error) {
	parts := strings.Split(line, "|")
	if len(parts) < 8 {
		return nil, fmt.Errorf("invalid commit header format: %s", line)
	}

	commit := &Commit{
		Hash:    parts[0],
		Message: parts[7],
		Author: Author{
			Name:  parts[1],
			Email: parts[2],
		},
		Committer: Author{
			Name:  parts[4],
			Email: parts[5],
		},
		Stats: CommitStats{
			Files: []FileChange{},
		},
	}

	// Parse author date
	if authorDate, err := time.Parse("2006-01-02 15:04:05 -0700", parts[3]); err == nil {
		commit.AuthorDate = authorDate
	}

	// Parse committer date
	if committerDate, err := time.Parse("2006-01-02 15:04:05 -0700", parts[6]); err == nil {
		commit.CommitterDate = committerDate
	}

	// Parse parent hashes
	if len(parts) > 8 && parts[8] != "" {
		commit.ParentHashes = strings.Split(parts[8], " ")
	}

	// Parse tree hash
	if len(parts) > 9 {
		commit.TreeHash = parts[9]
	}

	return commit, nil
}

// parseNumStatLine parses a numstat line in format: insertions\tdeletions\tfilename
func (p *GitOutputParser) parseNumStatLine(commit *Commit, line string) error {
	parts := strings.Split(line, "\t")
	if len(parts) < 3 {
		return fmt.Errorf("invalid numstat format: %s", line)
	}

	fileChange := FileChange{
		Path: parts[2],
	}

	// Parse insertions (might be "-" for binary files)
	if parts[0] != "-" {
		if insertions, err := strconv.Atoi(parts[0]); err == nil {
			fileChange.Insertions = insertions
			commit.Stats.Insertions += insertions
		}
	}

	// Parse deletions (might be "-" for binary files)
	if parts[1] != "-" {
		if deletions, err := strconv.Atoi(parts[1]); err == nil {
			fileChange.Deletions = deletions
			commit.Stats.Deletions += deletions
		}
	}

	// Determine file status (simplified - would need more complex logic for renames/copies)
	if fileChange.Insertions > 0 && fileChange.Deletions == 0 {
		fileChange.Status = "A" // Added
	} else if fileChange.Insertions == 0 && fileChange.Deletions > 0 {
		fileChange.Status = "D" // Deleted
	} else {
		fileChange.Status = "M" // Modified
	}

	commit.Stats.Files = append(commit.Stats.Files, fileChange)
	commit.Stats.FilesChanged++

	return nil
}

// ParseDiffStat parses git diff --stat output
func (p *GitOutputParser) ParseDiffStat(output string) (*CommitStats, error) {
	stats := &CommitStats{
		Files: []FileChange{},
	}

	if output == "" {
		return stats, nil
	}

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip summary line (contains "files changed", "insertions", "deletions")
		if strings.Contains(line, "file") && (strings.Contains(line, "changed") ||
			strings.Contains(line, "insertion") || strings.Contains(line, "deletion")) {
			continue
		}

		// Parse individual file stats
		if fileChange, err := p.parseDiffStatLine(line); err == nil {
			stats.Files = append(stats.Files, fileChange)
			stats.FilesChanged++
			stats.Insertions += fileChange.Insertions
			stats.Deletions += fileChange.Deletions
		}
	}

	return stats, nil
}

// parseDiffStatLine parses a single line from git diff --stat output
func (p *GitOutputParser) parseDiffStatLine(line string) (FileChange, error) {
	// Example: " src/main.go | 15 +++++++++------"
	parts := strings.Split(line, "|")
	if len(parts) != 2 {
		return FileChange{}, fmt.Errorf("invalid diff stat line: %s", line)
	}

	fileChange := FileChange{
		Path:   strings.TrimSpace(parts[0]),
		Status: "M", // Default to modified
	}

	// Parse the changes part
	changesStr := strings.TrimSpace(parts[1])

	// Extract number from the beginning
	re := regexp.MustCompile(`^(\d+)`)
	matches := re.FindStringSubmatch(changesStr)
	if len(matches) > 1 {
		if totalChanges, err := strconv.Atoi(matches[1]); err == nil {
			// Count + and - characters to determine insertions/deletions
			plusCount := strings.Count(changesStr, "+")
			minusCount := strings.Count(changesStr, "-")

			if plusCount > 0 && minusCount == 0 {
				fileChange.Insertions = totalChanges
				fileChange.Status = "A"
			} else if minusCount > 0 && plusCount == 0 {
				fileChange.Deletions = totalChanges
				fileChange.Status = "D"
			} else {
				// Approximate distribution for mixed changes
				fileChange.Insertions = (totalChanges * plusCount) / (plusCount + minusCount)
				fileChange.Deletions = totalChanges - fileChange.Insertions
			}
		}
	}

	return fileChange, nil
}

// ParseContributors parses git shortlog -sne output
func (p *GitOutputParser) ParseContributors(output string) ([]Contributor, error) {
	var contributors []Contributor

	if output == "" {
		return contributors, nil
	}

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		contributor, err := p.parseContributorLine(line)
		if err != nil {
			continue // Skip invalid lines
		}

		contributors = append(contributors, contributor)
	}

	return contributors, nil
}

// parseContributorLine parses a single contributor line from shortlog output
// Format: "    42  John Doe <john@example.com>"
func (p *GitOutputParser) parseContributorLine(line string) (Contributor, error) {
	// Use regex to parse the format
	re := regexp.MustCompile(`^\s*(\d+)\s+(.+?)\s+<(.+?)>$`)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 4 {
		return Contributor{}, fmt.Errorf("invalid contributor line: %s", line)
	}

	commits, err := strconv.Atoi(matches[1])
	if err != nil {
		return Contributor{}, fmt.Errorf("invalid commit count: %s", matches[1])
	}

	contributor := Contributor{
		Name:         matches[2],
		Email:        matches[3],
		TotalCommits: commits,
		CommitsByDay: make(map[string]int),
	}

	return contributor, nil
}

// ParseBranches parses git branch -a output
func (p *GitOutputParser) ParseBranches(output string) ([]string, error) {
	var branches []string

	if output == "" {
		return branches, nil
	}

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove current branch indicator (*) and whitespace
		branch := strings.TrimSpace(strings.TrimPrefix(line, "*"))

		// Skip HEAD references
		if strings.Contains(branch, "HEAD ->") {
			continue
		}

		// Clean up remote branch names
		if strings.HasPrefix(branch, "remotes/") {
			branch = strings.TrimPrefix(branch, "remotes/")
		}

		branches = append(branches, branch)
	}

	return branches, nil
}
