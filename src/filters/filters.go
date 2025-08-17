// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Filtering system for commits and data

package filters

import (
	"fmt"
	"git-stats/models"
	"regexp"
	"strings"
	"time"
)

// Filter represents a generic filter interface
type Filter interface {
	Apply(commits []models.Commit) []models.Commit
	Description() string
}

// FilterChain represents a chain of filters that can be applied sequentially
type FilterChain struct {
	filters []Filter
}

// NewFilterChain creates a new filter chain
func NewFilterChain() *FilterChain {
	return &FilterChain{
		filters: make([]Filter, 0),
	}
}

// Add adds a filter to the chain
func (fc *FilterChain) Add(filter Filter) *FilterChain {
	fc.filters = append(fc.filters, filter)
	return fc
}

// Apply applies all filters in the chain sequentially
func (fc *FilterChain) Apply(commits []models.Commit) []models.Commit {
	result := commits
	for _, filter := range fc.filters {
		result = filter.Apply(result)
	}
	return result
}

// GetFilters returns all filters in the chain
func (fc *FilterChain) GetFilters() []Filter {
	return fc.filters
}

// Clear removes all filters from the chain
func (fc *FilterChain) Clear() {
	fc.filters = fc.filters[:0]
}

// DateRangeFilter filters commits by date range
type DateRangeFilter struct {
	Since *time.Time
	Until *time.Time
}

// NewDateRangeFilter creates a new date range filter
func NewDateRangeFilter(since, until *time.Time) *DateRangeFilter {
	return &DateRangeFilter{
		Since: since,
		Until: until,
	}
}

// Apply applies the date range filter
func (drf *DateRangeFilter) Apply(commits []models.Commit) []models.Commit {
	if drf.Since == nil && drf.Until == nil {
		return commits
	}

	var filtered []models.Commit
	for _, commit := range commits {
		// Check since date
		if drf.Since != nil && commit.AuthorDate.Before(*drf.Since) {
			continue
		}

		// Check until date
		if drf.Until != nil && commit.AuthorDate.After(*drf.Until) {
			continue
		}

		filtered = append(filtered, commit)
	}

	return filtered
}

// Description returns a description of the filter
func (drf *DateRangeFilter) Description() string {
	if drf.Since != nil && drf.Until != nil {
		return "Date range: " + drf.Since.Format("2006-01-02") + " to " + drf.Until.Format("2006-01-02")
	} else if drf.Since != nil {
		return "Since: " + drf.Since.Format("2006-01-02")
	} else if drf.Until != nil {
		return "Until: " + drf.Until.Format("2006-01-02")
	}
	return "No date filter"
}

// AuthorFilter filters commits by author with advanced matching
type AuthorFilter struct {
	Pattern     string
	MatchType   AuthorMatchType
	CaseSensitive bool
	compiled    *regexp.Regexp
}

// AuthorMatchType defines how author matching should be performed
type AuthorMatchType int

const (
	// ExactMatch requires exact string match
	ExactMatch AuthorMatchType = iota
	// ContainsMatch requires the pattern to be contained in name or email
	ContainsMatch
	// RegexMatch uses regular expression matching
	RegexMatch
	// EmailMatch matches only email addresses
	EmailMatch
	// NameMatch matches only names
	NameMatch
)

// NewAuthorFilter creates a new author filter with default settings
func NewAuthorFilter(pattern string) *AuthorFilter {
	return &AuthorFilter{
		Pattern:       pattern,
		MatchType:     ContainsMatch,
		CaseSensitive: false,
	}
}

// NewAuthorFilterWithOptions creates a new author filter with custom options
func NewAuthorFilterWithOptions(pattern string, matchType AuthorMatchType, caseSensitive bool) (*AuthorFilter, error) {
	af := &AuthorFilter{
		Pattern:       pattern,
		MatchType:     matchType,
		CaseSensitive: caseSensitive,
	}

	// Compile regex if needed
	if matchType == RegexMatch {
		flags := "i" // case insensitive by default
		if caseSensitive {
			flags = ""
		}
		regexPattern := "(?" + flags + ")" + pattern
		compiled, err := regexp.Compile(regexPattern)
		if err != nil {
			return nil, err
		}
		af.compiled = compiled
	}

	return af, nil
}

// Apply applies the author filter
func (af *AuthorFilter) Apply(commits []models.Commit) []models.Commit {
	if af.Pattern == "" {
		return commits
	}

	var filtered []models.Commit
	for _, commit := range commits {
		if af.matchesAuthor(commit.Author) {
			filtered = append(filtered, commit)
		}
	}

	return filtered
}

// matchesAuthor checks if an author matches the filter criteria
func (af *AuthorFilter) matchesAuthor(author models.Author) bool {
	switch af.MatchType {
	case ExactMatch:
		return af.exactMatch(author)
	case ContainsMatch:
		return af.containsMatch(author)
	case RegexMatch:
		return af.regexMatch(author)
	case EmailMatch:
		return af.emailMatch(author)
	case NameMatch:
		return af.nameMatch(author)
	default:
		return af.containsMatch(author)
	}
}

// exactMatch performs exact string matching
func (af *AuthorFilter) exactMatch(author models.Author) bool {
	if af.CaseSensitive {
		return author.Name == af.Pattern || author.Email == af.Pattern
	}

	pattern := strings.ToLower(af.Pattern)
	name := strings.ToLower(author.Name)
	email := strings.ToLower(author.Email)

	return name == pattern || email == pattern
}

// containsMatch performs substring matching
func (af *AuthorFilter) containsMatch(author models.Author) bool {
	if af.CaseSensitive {
		return strings.Contains(author.Name, af.Pattern) || strings.Contains(author.Email, af.Pattern)
	}

	pattern := strings.ToLower(af.Pattern)
	name := strings.ToLower(author.Name)
	email := strings.ToLower(author.Email)

	return strings.Contains(name, pattern) || strings.Contains(email, pattern)
}

// regexMatch performs regular expression matching
func (af *AuthorFilter) regexMatch(author models.Author) bool {
	if af.compiled == nil {
		return false
	}

	return af.compiled.MatchString(author.Name) || af.compiled.MatchString(author.Email)
}

// emailMatch matches only email addresses
func (af *AuthorFilter) emailMatch(author models.Author) bool {
	if af.CaseSensitive {
		return strings.Contains(author.Email, af.Pattern)
	}

	pattern := strings.ToLower(af.Pattern)
	email := strings.ToLower(author.Email)

	return strings.Contains(email, pattern)
}

// nameMatch matches only names
func (af *AuthorFilter) nameMatch(author models.Author) bool {
	if af.CaseSensitive {
		return strings.Contains(author.Name, af.Pattern)
	}

	pattern := strings.ToLower(af.Pattern)
	name := strings.ToLower(author.Name)

	return strings.Contains(name, pattern)
}

// Description returns a description of the filter
func (af *AuthorFilter) Description() string {
	matchTypeStr := ""
	switch af.MatchType {
	case ExactMatch:
		matchTypeStr = "exact"
	case ContainsMatch:
		matchTypeStr = "contains"
	case RegexMatch:
		matchTypeStr = "regex"
	case EmailMatch:
		matchTypeStr = "email"
	case NameMatch:
		matchTypeStr = "name"
	}

	caseStr := ""
	if af.CaseSensitive {
		caseStr = " (case-sensitive)"
	}

	return "Author " + matchTypeStr + " match: '" + af.Pattern + "'" + caseStr
}

// MergeCommitFilter filters merge commits
type MergeCommitFilter struct {
	IncludeMerges bool
}

// NewMergeCommitFilter creates a new merge commit filter
func NewMergeCommitFilter(includeMerges bool) *MergeCommitFilter {
	return &MergeCommitFilter{
		IncludeMerges: includeMerges,
	}
}

// Apply applies the merge commit filter
func (mcf *MergeCommitFilter) Apply(commits []models.Commit) []models.Commit {
	if mcf.IncludeMerges {
		return commits // Include all commits
	}

	var filtered []models.Commit
	for _, commit := range commits {
		if !commit.IsMergeCommit() {
			filtered = append(filtered, commit)
		}
	}

	return filtered
}

// Description returns a description of the filter
func (mcf *MergeCommitFilter) Description() string {
	if mcf.IncludeMerges {
		return "Include merge commits"
	}
	return "Exclude merge commits"
}

// LimitFilter limits the number of commits
type LimitFilter struct {
	Limit int
}

// NewLimitFilter creates a new limit filter
func NewLimitFilter(limit int) *LimitFilter {
	return &LimitFilter{
		Limit: limit,
	}
}

// Apply applies the limit filter
func (lf *LimitFilter) Apply(commits []models.Commit) []models.Commit {
	if lf.Limit <= 0 || len(commits) <= lf.Limit {
		return commits
	}

	return commits[:lf.Limit]
}

// Description returns a description of the filter
func (lf *LimitFilter) Description() string {
	return "Limit: " + string(rune(lf.Limit)) + " commits"
}

// FilePathFilter filters commits that affect specific file paths
type FilePathFilter struct {
	Patterns      []string
	MatchType     FileMatchType
	CaseSensitive bool
}

// FileMatchType defines how file path matching should be performed
type FileMatchType int

const (
	// FileExactMatch requires exact path match
	FileExactMatch FileMatchType = iota
	// FileContainsMatch requires the pattern to be contained in the path
	FileContainsMatch
	// FileGlobMatch uses glob pattern matching
	FileGlobMatch
	// FileRegexMatch uses regular expression matching
	FileRegexMatch
)

// NewFilePathFilter creates a new file path filter
func NewFilePathFilter(patterns []string, matchType FileMatchType, caseSensitive bool) *FilePathFilter {
	return &FilePathFilter{
		Patterns:      patterns,
		MatchType:     matchType,
		CaseSensitive: caseSensitive,
	}
}

// Apply applies the file path filter
func (fpf *FilePathFilter) Apply(commits []models.Commit) []models.Commit {
	if len(fpf.Patterns) == 0 {
		return commits
	}

	var filtered []models.Commit
	for _, commit := range commits {
		if fpf.commitAffectsFiles(commit) {
			filtered = append(filtered, commit)
		}
	}

	return filtered
}

// commitAffectsFiles checks if a commit affects any of the specified file patterns
func (fpf *FilePathFilter) commitAffectsFiles(commit models.Commit) bool {
	for _, file := range commit.Stats.Files {
		for _, pattern := range fpf.Patterns {
			if fpf.matchesPath(file.Path, pattern) {
				return true
			}
		}
	}
	return false
}

// matchesPath checks if a file path matches a pattern
func (fpf *FilePathFilter) matchesPath(path, pattern string) bool {
	if !fpf.CaseSensitive {
		path = strings.ToLower(path)
		pattern = strings.ToLower(pattern)
	}

	switch fpf.MatchType {
	case FileExactMatch:
		return path == pattern
	case FileContainsMatch:
		return strings.Contains(path, pattern)
	case FileGlobMatch:
		// Simple glob matching (basic implementation)
		return fpf.simpleGlobMatch(path, pattern)
	case FileRegexMatch:
		matched, _ := regexp.MatchString(pattern, path)
		return matched
	default:
		return strings.Contains(path, pattern)
	}
}

// simpleGlobMatch performs basic glob matching
func (fpf *FilePathFilter) simpleGlobMatch(path, pattern string) bool {
	// Convert glob pattern to regex
	regexPattern := strings.ReplaceAll(pattern, "*", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "?", ".")
	regexPattern = "^" + regexPattern + "$"

	matched, _ := regexp.MatchString(regexPattern, path)
	return matched
}

// Description returns a description of the filter
func (fpf *FilePathFilter) Description() string {
	return "File paths: " + strings.Join(fpf.Patterns, ", ")
}

// ExcludeFilePathFilter filters out commits that affect specific file paths
type ExcludeFilePathFilter struct {
	*FilePathFilter
}

// NewExcludeFilePathFilter creates a new exclude file path filter
func NewExcludeFilePathFilter(patterns []string, matchType FileMatchType, caseSensitive bool) *ExcludeFilePathFilter {
	return &ExcludeFilePathFilter{
		FilePathFilter: NewFilePathFilter(patterns, matchType, caseSensitive),
	}
}

// Apply applies the exclude file path filter (inverse of include)
func (efpf *ExcludeFilePathFilter) Apply(commits []models.Commit) []models.Commit {
	if len(efpf.Patterns) == 0 {
		return commits
	}

	var filtered []models.Commit
	for _, commit := range commits {
		if !efpf.commitAffectsFiles(commit) {
			filtered = append(filtered, commit)
		}
	}

	return filtered
}

// Description returns a description of the filter
func (efpf *ExcludeFilePathFilter) Description() string {
	return "Exclude file paths: " + strings.Join(efpf.Patterns, ", ")
}

// BranchFilter filters commits by branch name
type BranchFilter struct {
	Branches      []string
	MatchType     BranchMatchType
	CaseSensitive bool
}

// BranchMatchType defines how branch matching should be performed
type BranchMatchType int

const (
	// BranchExactMatch requires exact branch name match
	BranchExactMatch BranchMatchType = iota
	// BranchContainsMatch requires the pattern to be contained in branch name
	BranchContainsMatch
	// BranchRegexMatch uses regular expression matching
	BranchRegexMatch
)

// NewBranchFilter creates a new branch filter
func NewBranchFilter(branches []string, matchType BranchMatchType, caseSensitive bool) *BranchFilter {
	return &BranchFilter{
		Branches:      branches,
		MatchType:     matchType,
		CaseSensitive: caseSensitive,
	}
}

// Apply applies the branch filter
func (bf *BranchFilter) Apply(commits []models.Commit) []models.Commit {
	if len(bf.Branches) == 0 {
		return commits
	}

	var filtered []models.Commit
	for _, commit := range commits {
		if bf.matchesBranch(commit) {
			filtered = append(filtered, commit)
		}
	}

	return filtered
}

// matchesBranch checks if a commit matches any of the branch patterns
func (bf *BranchFilter) matchesBranch(commit models.Commit) bool {
	// Note: This is a simplified implementation. In a real scenario,
	// you would need to determine which branch(es) a commit belongs to
	// by using git commands or maintaining branch information in the commit model

	// For now, we'll assume the commit has a Branch field (this would need to be added to the model)
	// This is a placeholder implementation
	return true // Always match for now - would need actual branch detection logic
}

// Description returns a description of the filter
func (bf *BranchFilter) Description() string {
	return "Branch filter: " + strings.Join(bf.Branches, ", ")
}

// MessageFilter filters commits by commit message content
type MessageFilter struct {
	Pattern       string
	MatchType     MessageMatchType
	CaseSensitive bool
	compiled      *regexp.Regexp
}

// MessageMatchType defines how message matching should be performed
type MessageMatchType int

const (
	// MessageContainsMatch requires the pattern to be contained in the message
	MessageContainsMatch MessageMatchType = iota
	// MessageRegexMatch uses regular expression matching
	MessageRegexMatch
	// MessageStartsWithMatch matches messages that start with the pattern
	MessageStartsWithMatch
	// MessageEndsWithMatch matches messages that end with the pattern
	MessageEndsWithMatch
)

// NewMessageFilter creates a new message filter
func NewMessageFilter(pattern string, matchType MessageMatchType, caseSensitive bool) (*MessageFilter, error) {
	mf := &MessageFilter{
		Pattern:       pattern,
		MatchType:     matchType,
		CaseSensitive: caseSensitive,
	}

	// Compile regex if needed
	if matchType == MessageRegexMatch {
		flags := "i" // case insensitive by default
		if caseSensitive {
			flags = ""
		}
		regexPattern := "(?" + flags + ")" + pattern
		compiled, err := regexp.Compile(regexPattern)
		if err != nil {
			return nil, err
		}
		mf.compiled = compiled
	}

	return mf, nil
}

// Apply applies the message filter
func (mf *MessageFilter) Apply(commits []models.Commit) []models.Commit {
	if mf.Pattern == "" {
		return commits
	}

	var filtered []models.Commit
	for _, commit := range commits {
		if mf.matchesMessage(commit.Message) {
			filtered = append(filtered, commit)
		}
	}

	return filtered
}

// matchesMessage checks if a message matches the filter criteria
func (mf *MessageFilter) matchesMessage(message string) bool {
	switch mf.MatchType {
	case MessageContainsMatch:
		return mf.containsMatch(message)
	case MessageRegexMatch:
		return mf.regexMatch(message)
	case MessageStartsWithMatch:
		return mf.startsWithMatch(message)
	case MessageEndsWithMatch:
		return mf.endsWithMatch(message)
	default:
		return mf.containsMatch(message)
	}
}

// containsMatch performs substring matching
func (mf *MessageFilter) containsMatch(message string) bool {
	if mf.CaseSensitive {
		return strings.Contains(message, mf.Pattern)
	}

	pattern := strings.ToLower(mf.Pattern)
	msg := strings.ToLower(message)
	return strings.Contains(msg, pattern)
}

// regexMatch performs regular expression matching
func (mf *MessageFilter) regexMatch(message string) bool {
	if mf.compiled == nil {
		return false
	}
	return mf.compiled.MatchString(message)
}

// startsWithMatch checks if message starts with pattern
func (mf *MessageFilter) startsWithMatch(message string) bool {
	if mf.CaseSensitive {
		return strings.HasPrefix(message, mf.Pattern)
	}

	pattern := strings.ToLower(mf.Pattern)
	msg := strings.ToLower(message)
	return strings.HasPrefix(msg, pattern)
}

// endsWithMatch checks if message ends with pattern
func (mf *MessageFilter) endsWithMatch(message string) bool {
	if mf.CaseSensitive {
		return strings.HasSuffix(message, mf.Pattern)
	}

	pattern := strings.ToLower(mf.Pattern)
	msg := strings.ToLower(message)
	return strings.HasSuffix(msg, pattern)
}

// Description returns a description of the filter
func (mf *MessageFilter) Description() string {
	matchTypeStr := ""
	switch mf.MatchType {
	case MessageContainsMatch:
		matchTypeStr = "contains"
	case MessageRegexMatch:
		matchTypeStr = "regex"
	case MessageStartsWithMatch:
		matchTypeStr = "starts with"
	case MessageEndsWithMatch:
		matchTypeStr = "ends with"
	}

	caseStr := ""
	if mf.CaseSensitive {
		caseStr = " (case-sensitive)"
	}

	return "Message " + matchTypeStr + " match: '" + mf.Pattern + "'" + caseStr
}

// FileSizeFilter filters commits based on the size of changes
type FileSizeFilter struct {
	MinInsertions int
	MaxInsertions int
	MinDeletions  int
	MaxDeletions  int
	MinFiles      int
	MaxFiles      int
}

// NewFileSizeFilter creates a new file size filter
func NewFileSizeFilter(minInsertions, maxInsertions, minDeletions, maxDeletions, minFiles, maxFiles int) *FileSizeFilter {
	return &FileSizeFilter{
		MinInsertions: minInsertions,
		MaxInsertions: maxInsertions,
		MinDeletions:  minDeletions,
		MaxDeletions:  maxDeletions,
		MinFiles:      minFiles,
		MaxFiles:      maxFiles,
	}
}

// Apply applies the file size filter
func (fsf *FileSizeFilter) Apply(commits []models.Commit) []models.Commit {
	var filtered []models.Commit
	for _, commit := range commits {
		if fsf.matchesSize(commit) {
			filtered = append(filtered, commit)
		}
	}

	return filtered
}

// matchesSize checks if a commit matches the size criteria
func (fsf *FileSizeFilter) matchesSize(commit models.Commit) bool {
	stats := commit.Stats

	// Check insertions
	if fsf.MinInsertions > 0 && stats.Insertions < fsf.MinInsertions {
		return false
	}
	if fsf.MaxInsertions > 0 && stats.Insertions > fsf.MaxInsertions {
		return false
	}

	// Check deletions
	if fsf.MinDeletions > 0 && stats.Deletions < fsf.MinDeletions {
		return false
	}
	if fsf.MaxDeletions > 0 && stats.Deletions > fsf.MaxDeletions {
		return false
	}

	// Check files changed
	if fsf.MinFiles > 0 && stats.FilesChanged < fsf.MinFiles {
		return false
	}
	if fsf.MaxFiles > 0 && stats.FilesChanged > fsf.MaxFiles {
		return false
	}

	return true
}

// Description returns a description of the filter
func (fsf *FileSizeFilter) Description() string {
	var parts []string

	if fsf.MinInsertions > 0 || fsf.MaxInsertions > 0 {
		if fsf.MinInsertions > 0 && fsf.MaxInsertions > 0 {
			parts = append(parts, fmt.Sprintf("insertions: %d-%d", fsf.MinInsertions, fsf.MaxInsertions))
		} else if fsf.MinInsertions > 0 {
			parts = append(parts, fmt.Sprintf("insertions: >=%d", fsf.MinInsertions))
		} else {
			parts = append(parts, fmt.Sprintf("insertions: <=%d", fsf.MaxInsertions))
		}
	}

	if fsf.MinDeletions > 0 || fsf.MaxDeletions > 0 {
		if fsf.MinDeletions > 0 && fsf.MaxDeletions > 0 {
			parts = append(parts, fmt.Sprintf("deletions: %d-%d", fsf.MinDeletions, fsf.MaxDeletions))
		} else if fsf.MinDeletions > 0 {
			parts = append(parts, fmt.Sprintf("deletions: >=%d", fsf.MinDeletions))
		} else {
			parts = append(parts, fmt.Sprintf("deletions: <=%d", fsf.MaxDeletions))
		}
	}

	if fsf.MinFiles > 0 || fsf.MaxFiles > 0 {
		if fsf.MinFiles > 0 && fsf.MaxFiles > 0 {
			parts = append(parts, fmt.Sprintf("files: %d-%d", fsf.MinFiles, fsf.MaxFiles))
		} else if fsf.MinFiles > 0 {
			parts = append(parts, fmt.Sprintf("files: >=%d", fsf.MinFiles))
		} else {
			parts = append(parts, fmt.Sprintf("files: <=%d", fsf.MaxFiles))
		}
	}

	if len(parts) == 0 {
		return "No size filter"
	}

	return "Size filter: " + strings.Join(parts, ", ")
}
