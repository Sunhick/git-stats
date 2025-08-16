// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Integration between CLI, config, and filtering systems

package integration

import (
	"time"

	"git-stats/cli"
	"git-stats/config"
	"git-stats/filters"
	"git-stats/models"
)

// FilteredConfigManager manages the integration between CLI config, app config, and filtering
type FilteredConfigManager struct {
	configManager *config.ConfigManager
	filterBuilder *filters.FilterBuilder
}

// NewFilteredConfigManager creates a new filtered config manager
func NewFilteredConfigManager() *FilteredConfigManager {
	configManager := config.NewConfigManager()
	filterBuilder := filters.NewFilterBuilder(configManager)

	return &FilteredConfigManager{
		configManager: configManager,
		filterBuilder: filterBuilder,
	}
}

// NewFilteredConfigManagerWithPath creates a new filtered config manager with custom config path
func NewFilteredConfigManagerWithPath(configPath string) *FilteredConfigManager {
	configManager := config.NewConfigManagerWithPath(configPath)
	filterBuilder := filters.NewFilterBuilder(configManager)

	return &FilteredConfigManager{
		configManager: configManager,
		filterBuilder: filterBuilder,
	}
}

// LoadConfig loads the application configuration
func (fcm *FilteredConfigManager) LoadConfig() error {
	return fcm.configManager.Load()
}

// SaveConfig saves the application configuration
func (fcm *FilteredConfigManager) SaveConfig() error {
	return fcm.configManager.Save()
}

// GetConfigManager returns the config manager
func (fcm *FilteredConfigManager) GetConfigManager() *config.ConfigManager {
	return fcm.configManager
}

// GetFilterBuilder returns the filter builder
func (fcm *FilteredConfigManager) GetFilterBuilder() *filters.FilterBuilder {
	return fcm.filterBuilder
}

// BuildFilterChain builds a filter chain from CLI configuration
func (fcm *FilteredConfigManager) BuildFilterChain(cliConfig *cli.Config) (*filters.FilterChain, error) {
	return fcm.filterBuilder.BuildFromCLIConfig(cliConfig)
}

// BuildDefaultFilterChain builds a filter chain from application configuration only
func (fcm *FilteredConfigManager) BuildDefaultFilterChain() (*filters.FilterChain, error) {
	return fcm.filterBuilder.BuildFromConfig()
}

// ApplyFilters applies the filter chain to commits
func (fcm *FilteredConfigManager) ApplyFilters(commits []models.Commit, cliConfig *cli.Config) ([]models.Commit, error) {
	chain, err := fcm.BuildFilterChain(cliConfig)
	if err != nil {
		return nil, err
	}

	return chain.Apply(commits), nil
}

// GetFilterSummary returns a human-readable summary of active filters
func (fcm *FilteredConfigManager) GetFilterSummary(cliConfig *cli.Config) (string, error) {
	chain, err := fcm.BuildFilterChain(cliConfig)
	if err != nil {
		return "", err
	}

	return fcm.filterBuilder.GetFilterSummary(chain), nil
}

// CreateAnalysisConfig creates an AnalysisConfig from CLI configuration
func (fcm *FilteredConfigManager) CreateAnalysisConfig(cliConfig *cli.Config) models.AnalysisConfig {
	appConfig := fcm.configManager.GetConfig()

	analysisConfig := models.AnalysisConfig{
		Limit:         cliConfig.Limit,
		IncludeMerges: appConfig.Filters.IncludeMerges,
	}

	// Set time range
	if cliConfig.Since != nil || cliConfig.Until != nil {
		analysisConfig.TimeRange = models.TimeRange{
			Start: getTimeOrZero(cliConfig.Since),
			End:   getTimeOrZero(cliConfig.Until),
		}
	}

	// Set author filter
	if cliConfig.Author != "" {
		analysisConfig.AuthorFilter = cliConfig.Author
	} else if appConfig.Filters.DefaultAuthor != "" {
		analysisConfig.AuthorFilter = appConfig.Filters.DefaultAuthor
	}

	return analysisConfig
}

// UpdateConfigFromCLI updates application config with CLI overrides
func (fcm *FilteredConfigManager) UpdateConfigFromCLI(cliConfig *cli.Config) {
	appConfig := fcm.configManager.GetConfig()

	// Update defaults based on CLI flags
	if cliConfig.Command != "" {
		appConfig.Defaults.Command = cliConfig.Command
	}

	if cliConfig.Format != "" {
		appConfig.Defaults.Format = cliConfig.Format
	}

	if cliConfig.ShowProgress {
		appConfig.Defaults.ShowProgress = cliConfig.ShowProgress
	}

	if cliConfig.RepoPath != "" && cliConfig.RepoPath != "." {
		appConfig.Defaults.RepoPath = cliConfig.RepoPath
	}

	// Update output settings
	if cliConfig.NoColor {
		appConfig.Output.ColorEnabled = false
	}

	if cliConfig.ColorTheme != "" {
		appConfig.Output.ColorTheme = cliConfig.ColorTheme
	}

	// Update performance settings
	if cliConfig.Limit > 0 && cliConfig.Limit != 10000 { // 10000 is default
		appConfig.Performance.MaxCommits = cliConfig.Limit
	}

	fcm.configManager.SetConfig(appConfig)
}

// GetEffectiveConfig returns the effective configuration combining app config and CLI overrides
func (fcm *FilteredConfigManager) GetEffectiveConfig(cliConfig *cli.Config) *EffectiveConfig {
	appConfig := fcm.configManager.GetConfig()

	return &EffectiveConfig{
		Command:      getStringOrDefault(cliConfig.Command, appConfig.Defaults.Command),
		Format:       getStringOrDefault(cliConfig.Format, appConfig.Defaults.Format),
		RepoPath:     getStringOrDefault(cliConfig.RepoPath, appConfig.Defaults.RepoPath),
		ShowProgress: cliConfig.ShowProgress || appConfig.Defaults.ShowProgress,
		Since:        cliConfig.Since,
		Until:        cliConfig.Until,
		Author:       getStringOrDefault(cliConfig.Author, appConfig.Filters.DefaultAuthor),
		OutputFile:   cliConfig.OutputFile,
		Limit:        getIntOrDefault(cliConfig.Limit, appConfig.Performance.MaxCommits),
		GUIMode:      cliConfig.GUIMode,
		ShowHelp:     cliConfig.ShowHelp,
		NoColor:      cliConfig.NoColor || !appConfig.Output.ColorEnabled,
		ColorTheme:   getStringOrDefault(cliConfig.ColorTheme, appConfig.Output.ColorTheme),

		// Additional effective settings
		IncludeMerges:     appConfig.Filters.IncludeMerges,
		CaseSensitive:     appConfig.Filters.CaseSensitive,
		AuthorMatchType:   appConfig.Filters.AuthorMatchType,
		ExcludePatterns:   appConfig.Filters.ExcludePatterns,
		IncludePatterns:   appConfig.Filters.IncludePatterns,
		PrettyPrint:       appConfig.Output.PrettyPrint,
		IncludeMetadata:   appConfig.Output.IncludeMetadata,
		DateFormat:        appConfig.Output.DateFormat,
		TimeFormat:        appConfig.Output.TimeFormat,
		ChunkSize:         appConfig.Performance.ChunkSize,
		CacheEnabled:      appConfig.Performance.CacheEnabled,
		ParallelProcessing: appConfig.Performance.ParallelProcessing,
		MaxWorkers:        appConfig.Performance.MaxWorkers,
	}
}

// EffectiveConfig represents the final effective configuration
type EffectiveConfig struct {
	// CLI-configurable options
	Command      string
	Format       string
	RepoPath     string
	ShowProgress bool
	Since        *time.Time
	Until        *time.Time
	Author       string
	OutputFile   string
	Limit        int
	GUIMode      bool
	ShowHelp     bool
	NoColor      bool
	ColorTheme   string

	// Additional config-only options
	IncludeMerges      bool
	CaseSensitive      bool
	AuthorMatchType    string
	ExcludePatterns    []string
	IncludePatterns    []string
	PrettyPrint        bool
	IncludeMetadata    bool
	DateFormat         string
	TimeFormat         string
	ChunkSize          int
	CacheEnabled       bool
	ParallelProcessing bool
	MaxWorkers         int
}

// Helper functions
func getStringOrDefault(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

func getIntOrDefault(value, defaultValue int) int {
	if value > 0 {
		return value
	}
	return defaultValue
}

func getTimeOrZero(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

// ValidateEffectiveConfig validates the effective configuration
func (fcm *FilteredConfigManager) ValidateEffectiveConfig(effectiveConfig *EffectiveConfig) error {
	// Create a temporary config for validation
	tempConfig := fcm.configManager.GetConfig()

	// Update with effective values
	tempConfig.Defaults.Command = effectiveConfig.Command
	tempConfig.Defaults.Format = effectiveConfig.Format
	tempConfig.Defaults.RepoPath = effectiveConfig.RepoPath
	tempConfig.Output.ColorTheme = effectiveConfig.ColorTheme
	tempConfig.Performance.MaxCommits = effectiveConfig.Limit
	tempConfig.Performance.ChunkSize = effectiveConfig.ChunkSize
	tempConfig.Performance.MaxWorkers = effectiveConfig.MaxWorkers

	// Create temporary manager for validation
	tempManager := config.NewConfigManager()
	tempManager.SetConfig(tempConfig)

	return tempManager.Validate()
}
