// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Git package for output parsing

package git

// Parser interface for parsing git command output
type Parser interface {
	ParseCommitLog(output string) ([]Commit, error)
	ParseDiffStat(output string) (*CommitStats, error)
	ParseContributors(output string) ([]Contributor, error)
	ParseBranches(output string) ([]string, error)
}
