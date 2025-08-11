// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Shared utilities for actions

package actions

import (
	"git-stats/git"
	"git-stats/models"
)

// convertGitCommitsToModelCommits converts git.Commit slice to models.Commit slice
func convertGitCommitsToModelCommits(gitCommits []git.Commit) []models.Commit {
	modelCommits := make([]models.Commit, len(gitCommits))

	for i, gitCommit := range gitCommits {
		// Convert file changes
		fileChanges := make([]models.FileChange, len(gitCommit.Stats.Files))
		for j, gitFile := range gitCommit.Stats.Files {
			fileChanges[j] = models.FileChange{
				Path:       gitFile.Path,
				Status:     gitFile.Status,
				Insertions: gitFile.Insertions,
				Deletions:  gitFile.Deletions,
			}
		}

		modelCommits[i] = models.Commit{
			Hash:    gitCommit.Hash,
			Message: gitCommit.Message,
			Author: models.Author{
				Name:  gitCommit.Author.Name,
				Email: gitCommit.Author.Email,
			},
			Committer: models.Author{
				Name:  gitCommit.Committer.Name,
				Email: gitCommit.Committer.Email,
			},
			AuthorDate:    gitCommit.AuthorDate,
			CommitterDate: gitCommit.CommitterDate,
			ParentHashes:  gitCommit.ParentHashes,
			TreeHash:      gitCommit.TreeHash,
			Stats: models.CommitStats{
				FilesChanged: gitCommit.Stats.FilesChanged,
				Insertions:   gitCommit.Stats.Insertions,
				Deletions:    gitCommit.Stats.Deletions,
				Files:        fileChanges,
			},
		}
	}

	return modelCommits
}

// convertGitContributorsToModelContributors converts git.Contributor slice to models.Contributor slice
func convertGitContributorsToModelContributors(gitContributors []git.Contributor) []models.Contributor {
	modelContributors := make([]models.Contributor, len(gitContributors))

	for i, gitContrib := range gitContributors {
		modelContributors[i] = models.Contributor{
			Name:             gitContrib.Name,
			Email:            gitContrib.Email,
			TotalCommits:     gitContrib.TotalCommits,
			TotalInsertions:  gitContrib.TotalInsertions,
			TotalDeletions:   gitContrib.TotalDeletions,
			FirstCommit:      gitContrib.FirstCommit,
			LastCommit:       gitContrib.LastCommit,
			ActiveDays:       gitContrib.ActiveDays,
			CommitsByDay:     gitContrib.CommitsByDay,
			CommitsByHour:    make(map[int]int),    // Initialize empty
			CommitsByWeekday: make(map[int]int),    // Initialize empty
			FileTypes:        make(map[string]int), // Initialize empty
			TopFiles:         []string{},           // Initialize empty
		}
	}

	return modelContributors
}
