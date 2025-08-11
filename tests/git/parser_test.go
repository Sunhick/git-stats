// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Git output parser tests

package git

import (
	"testing"
	"time"

	"git-stats/git"
)

func TestGitOutputParser_ParseCommitLog(t *testing.T) {
	parser := git.NewGitOutputParser()

	tests := []struct {
		name     string
		output   string
		expected []git.Commit
		wantErr  bool
	}{
		{
			name:     "empty output",
			output:   "",
			expected: []git.Commit{},
			wantErr:  false,
		},
		{
			name: "single commit with stats",
			output: `abc123|John Doe|john@example.com|2024-01-15 10:30:00 -0800|John Doe|john@example.com|2024-01-15 10:30:00 -0800|Initial commit|parent123|tree456

5	2	README.md
10	0	src/main.go

`,
			expected: []git.Commit{
				{
					Hash:    "abc123",
					Message: "Initial commit",
					Author: git.Author{
						Name:  "John Doe",
						Email: "john@example.com",
					},
					Committer: git.Author{
						Name:  "John Doe",
						Email: "john@example.com",
					},
					AuthorDate:    time.Date(2024, 1, 15, 10, 30, 0, 0, time.FixedZone("PST", -8*3600)),
					CommitterDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.FixedZone("PST", -8*3600)),
					ParentHashes:  []string{"parent123"},
					TreeHash:      "tree456",
					Stats: git.CommitStats{
						FilesChanged: 2,
						Insertions:   15,
						Deletions:    2,
						Files: []git.FileChange{
							{
								Path:       "README.md",
								Status:     "M",
								Insertions: 5,
								Deletions:  2,
							},
							{
								Path:       "src/main.go",
								Status:     "A",
								Insertions: 10,
								Deletions:  0,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple commits",
			output: `def456|Jane Smith|jane@example.com|2024-01-16 14:20:00 -0800|Jane Smith|jane@example.com|2024-01-16 14:20:00 -0800|Add feature|abc123|tree789

3	1	feature.go

ghi789|Bob Wilson|bob@example.com|2024-01-17 09:15:00 -0800|Bob Wilson|bob@example.com|2024-01-17 09:15:00 -0800|Fix bug|def456|tree012

0	5	bug.go

`,
			expected: []git.Commit{
				{
					Hash:    "def456",
					Message: "Add feature",
					Author: git.Author{
						Name:  "Jane Smith",
						Email: "jane@example.com",
					},
					Committer: git.Author{
						Name:  "Jane Smith",
						Email: "jane@example.com",
					},
					AuthorDate:    time.Date(2024, 1, 16, 14, 20, 0, 0, time.FixedZone("PST", -8*3600)),
					CommitterDate: time.Date(2024, 1, 16, 14, 20, 0, 0, time.FixedZone("PST", -8*3600)),
					ParentHashes:  []string{"abc123"},
					TreeHash:      "tree789",
					Stats: git.CommitStats{
						FilesChanged: 1,
						Insertions:   3,
						Deletions:    1,
						Files: []git.FileChange{
							{
								Path:       "feature.go",
								Status:     "M",
								Insertions: 3,
								Deletions:  1,
							},
						},
					},
				},
				{
					Hash:    "ghi789",
					Message: "Fix bug",
					Author: git.Author{
						Name:  "Bob Wilson",
						Email: "bob@example.com",
					},
					Committer: git.Author{
						Name:  "Bob Wilson",
						Email: "bob@example.com",
					},
					AuthorDate:    time.Date(2024, 1, 17, 9, 15, 0, 0, time.FixedZone("PST", -8*3600)),
					CommitterDate: time.Date(2024, 1, 17, 9, 15, 0, 0, time.FixedZone("PST", -8*3600)),
					ParentHashes:  []string{"def456"},
					TreeHash:      "tree012",
					Stats: git.CommitStats{
						FilesChanged: 1,
						Insertions:   0,
						Deletions:    5,
						Files: []git.FileChange{
							{
								Path:       "bug.go",
								Status:     "D",
								Insertions: 0,
								Deletions:  5,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "commit with binary files",
			output: `xyz123|Alice Brown|alice@example.com|2024-01-18 16:45:00 -0800|Alice Brown|alice@example.com|2024-01-18 16:45:00 -0800|Add binary file||tree345

-	-	image.png
2	0	README.md

`,
			expected: []git.Commit{
				{
					Hash:    "xyz123",
					Message: "Add binary file",
					Author: git.Author{
						Name:  "Alice Brown",
						Email: "alice@example.com",
					},
					Committer: git.Author{
						Name:  "Alice Brown",
						Email: "alice@example.com",
					},
					AuthorDate:    time.Date(2024, 1, 18, 16, 45, 0, 0, time.FixedZone("PST", -8*3600)),
					CommitterDate: time.Date(2024, 1, 18, 16, 45, 0, 0, time.FixedZone("PST", -8*3600)),
					ParentHashes:  []string{},
					TreeHash:      "tree345",
					Stats: git.CommitStats{
						FilesChanged: 2,
						Insertions:   2,
						Deletions:    0,
						Files: []git.FileChange{
							{
								Path:       "image.png",
								Status:     "M",
								Insertions: 0,
								Deletions:  0,
							},
							{
								Path:       "README.md",
								Status:     "A",
								Insertions: 2,
								Deletions:  0,
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commits, err := parser.ParseCommitLog(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCommitLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(commits) != len(tt.expected) {
				t.Errorf("ParseCommitLog() returned %d commits, expected %d", len(commits), len(tt.expected))
				return
			}

			for i, commit := range commits {
				expected := tt.expected[i]

				if commit.Hash != expected.Hash {
					t.Errorf("Commit %d hash = %v, expected %v", i, commit.Hash, expected.Hash)
				}
				if commit.Message != expected.Message {
					t.Errorf("Commit %d message = %v, expected %v", i, commit.Message, expected.Message)
				}
				if commit.Author.Name != expected.Author.Name {
					t.Errorf("Commit %d author name = %v, expected %v", i, commit.Author.Name, expected.Author.Name)
				}
				if commit.Stats.FilesChanged != expected.Stats.FilesChanged {
					t.Errorf("Commit %d files changed = %v, expected %v", i, commit.Stats.FilesChanged, expected.Stats.FilesChanged)
				}
				if commit.Stats.Insertions != expected.Stats.Insertions {
					t.Errorf("Commit %d insertions = %v, expected %v", i, commit.Stats.Insertions, expected.Stats.Insertions)
				}
				if commit.Stats.Deletions != expected.Stats.Deletions {
					t.Errorf("Commit %d deletions = %v, expected %v", i, commit.Stats.Deletions, expected.Stats.Deletions)
				}
			}
		})
	}
}

func TestGitOutputParser_ParseDiffStat(t *testing.T) {
	parser := git.NewGitOutputParser()

	tests := []struct {
		name     string
		output   string
		expected *git.CommitStats
		wantErr  bool
	}{
		{
			name:   "empty output",
			output: "",
			expected: &git.CommitStats{
				Files: []git.FileChange{},
			},
			wantErr: false,
		},
		{
			name: "single file change",
			output: ` src/main.go | 15 +++++++++------
 1 file changed, 9 insertions(+), 6 deletions(-)`,
			expected: &git.CommitStats{
				FilesChanged: 1,
				Insertions:   9,
				Deletions:    6,
				Files: []git.FileChange{
					{
						Path:       "src/main.go",
						Status:     "M",
						Insertions: 9,
						Deletions:  6,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple file changes",
			output: ` README.md    | 5 +++++
 src/main.go  | 10 ++++++++++
 test.go      | 3 ---
 3 files changed, 15 insertions(+), 3 deletions(-)`,
			expected: &git.CommitStats{
				FilesChanged: 3,
				Insertions:   15,
				Deletions:    3,
				Files: []git.FileChange{
					{
						Path:       "README.md",
						Status:     "A",
						Insertions: 5,
						Deletions:  0,
					},
					{
						Path:       "src/main.go",
						Status:     "A",
						Insertions: 10,
						Deletions:  0,
					},
					{
						Path:       "test.go",
						Status:     "D",
						Insertions: 0,
						Deletions:  3,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := parser.ParseDiffStat(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDiffStat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if stats.FilesChanged != tt.expected.FilesChanged {
				t.Errorf("ParseDiffStat() files changed = %v, expected %v", stats.FilesChanged, tt.expected.FilesChanged)
			}
			if stats.Insertions != tt.expected.Insertions {
				t.Errorf("ParseDiffStat() insertions = %v, expected %v", stats.Insertions, tt.expected.Insertions)
			}
			if stats.Deletions != tt.expected.Deletions {
				t.Errorf("ParseDiffStat() deletions = %v, expected %v", stats.Deletions, tt.expected.Deletions)
			}
		})
	}
}

func TestGitOutputParser_ParseContributors(t *testing.T) {
	parser := git.NewGitOutputParser()

	tests := []struct {
		name     string
		output   string
		expected []git.Contributor
		wantErr  bool
	}{
		{
			name:     "empty output",
			output:   "",
			expected: []git.Contributor{},
			wantErr:  false,
		},
		{
			name: "single contributor",
			output: `    42	John Doe <john@example.com>`,
			expected: []git.Contributor{
				{
					Name:         "John Doe",
					Email:        "john@example.com",
					TotalCommits: 42,
					CommitsByDay: make(map[string]int),
				},
			},
			wantErr: false,
		},
		{
			name: "multiple contributors",
			output: `    25	Alice Smith <alice@example.com>
    18	Bob Wilson <bob@example.com>
     7	Charlie Brown <charlie@example.com>`,
			expected: []git.Contributor{
				{
					Name:         "Alice Smith",
					Email:        "alice@example.com",
					TotalCommits: 25,
					CommitsByDay: make(map[string]int),
				},
				{
					Name:         "Bob Wilson",
					Email:        "bob@example.com",
					TotalCommits: 18,
					CommitsByDay: make(map[string]int),
				},
				{
					Name:         "Charlie Brown",
					Email:        "charlie@example.com",
					TotalCommits: 7,
					CommitsByDay: make(map[string]int),
				},
			},
			wantErr: false,
		},
		{
			name: "contributor with special characters in name",
			output: `    15	José García-López <jose@example.com>`,
			expected: []git.Contributor{
				{
					Name:         "José García-López",
					Email:        "jose@example.com",
					TotalCommits: 15,
					CommitsByDay: make(map[string]int),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contributors, err := parser.ParseContributors(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseContributors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(contributors) != len(tt.expected) {
				t.Errorf("ParseContributors() returned %d contributors, expected %d", len(contributors), len(tt.expected))
				return
			}

			for i, contributor := range contributors {
				expected := tt.expected[i]

				if contributor.Name != expected.Name {
					t.Errorf("Contributor %d name = %v, expected %v", i, contributor.Name, expected.Name)
				}
				if contributor.Email != expected.Email {
					t.Errorf("Contributor %d email = %v, expected %v", i, contributor.Email, expected.Email)
				}
				if contributor.TotalCommits != expected.TotalCommits {
					t.Errorf("Contributor %d total commits = %v, expected %v", i, contributor.TotalCommits, expected.TotalCommits)
				}
			}
		})
	}
}

func TestGitOutputParser_ParseBranches(t *testing.T) {
	parser := git.NewGitOutputParser()

	tests := []struct {
		name     string
		output   string
		expected []string
		wantErr  bool
	}{
		{
			name:     "empty output",
			output:   "",
			expected: []string{},
			wantErr:  false,
		},
		{
			name: "single branch",
			output: `* main`,
			expected: []string{"main"},
			wantErr:  false,
		},
		{
			name: "multiple local branches",
			output: `  develop
* main
  feature/new-ui`,
			expected: []string{"develop", "main", "feature/new-ui"},
			wantErr:  false,
		},
		{
			name: "local and remote branches",
			output: `  develop
* main
  feature/new-ui
  remotes/origin/develop
  remotes/origin/main
  remotes/origin/HEAD -> origin/main`,
			expected: []string{"develop", "main", "feature/new-ui", "origin/develop", "origin/main"},
			wantErr:  false,
		},
		{
			name: "branches with remote tracking",
			output: `* main
  develop
  remotes/origin/main
  remotes/origin/develop
  remotes/upstream/main`,
			expected: []string{"main", "develop", "origin/main", "origin/develop", "upstream/main"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branches, err := parser.ParseBranches(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBranches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(branches) != len(tt.expected) {
				t.Errorf("ParseBranches() returned %d branches, expected %d", len(branches), len(tt.expected))
			} else {
				for i, branch := range branches {
					if branch != tt.expected[i] {
						t.Errorf("ParseBranches() branch %d = %v, expected %v", i, branch, tt.expected[i])
					}
				}
			}
		})
	}
}

// Test edge cases and error conditions
func TestGitOutputParser_EdgeCases(t *testing.T) {
	parser := git.NewGitOutputParser()

	t.Run("malformed commit header", func(t *testing.T) {
		output := `invalid|header|format`
		commits, err := parser.ParseCommitLog(output)
		if err != nil {
			t.Errorf("ParseCommitLog() should handle malformed headers gracefully, got error: %v", err)
		}
		// Should return empty slice or skip malformed entries
		if len(commits) > 0 {
			t.Errorf("ParseCommitLog() should not return commits for malformed input")
		}
	})

	t.Run("invalid contributor format", func(t *testing.T) {
		output := `invalid contributor line without proper format`
		contributors, err := parser.ParseContributors(output)
		if err != nil {
			t.Errorf("ParseContributors() should handle invalid format gracefully, got error: %v", err)
		}
		// Should return empty slice for invalid format
		if len(contributors) > 0 {
			t.Errorf("ParseContributors() should not return contributors for invalid input")
		}
	})

	t.Run("mixed valid and invalid lines", func(t *testing.T) {
		output := `    25	Alice Smith <alice@example.com>
invalid line
    18	Bob Wilson <bob@example.com>`
		contributors, err := parser.ParseContributors(output)
		if err != nil {
			t.Errorf("ParseContributors() should handle mixed input gracefully, got error: %v", err)
		}
		// Should return only valid contributors
		if len(contributors) != 2 {
			t.Errorf("ParseContributors() should return 2 valid contributors, got %d", len(contributors))
		}
	})
}

// Benchmark tests
func BenchmarkGitOutputParser_ParseCommitLog(b *testing.B) {
	parser := git.NewGitOutputParser()
	output := `abc123|John Doe|john@example.com|2024-01-15 10:30:00 -0800|John Doe|john@example.com|2024-01-15 10:30:00 -0800|Initial commit|parent123|tree456

5	2	README.md
10	0	src/main.go

def456|Jane Smith|jane@example.com|2024-01-16 14:20:00 -0800|Jane Smith|jane@example.com|2024-01-16 14:20:00 -0800|Add feature|abc123|tree789

3	1	feature.go

`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseCommitLog(output)
		if err != nil {
			b.Errorf("ParseCommitLog() error = %v", err)
		}
	}
}

func BenchmarkGitOutputParser_ParseContributors(b *testing.B) {
	parser := git.NewGitOutputParser()
	output := `    25	Alice Smith <alice@example.com>
    18	Bob Wilson <bob@example.com>
     7	Charlie Brown <charlie@example.com>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseContributors(output)
		if err != nil {
			b.Errorf("ParseContributors() error = %v", err)
		}
	}
}
