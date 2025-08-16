// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Configuration management system

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the complete application configuration
type Config struct {
	// Default settings
	Defaults DefaultConfig `json:"defaults"`

	// Filter settings
	Filters FilterConfig `json:"filters"`

	// Output settings
	Output OutputConfig `json:"output"`

	// Performance settings
	Performance PerformanceConfig `json:"performance"`

	// GUI settings
	GUI GUIConfig `json:"gui"`
}

// DefaultConfig contains default application settings
type DefaultConfig struct {
	Command      string `json:"command"`       // Default command (contrib, summary, etc.)
	DateRange    string `json:"date_range"`    // Default date range (e.g., "1 year ago")
	Format       string `json:"format"`        // Default output format
	ShowProgress bool   `json:"show_progress"` // Show progress by default
	RepoPath     string `json:"repo_path"`     // Default repository path
}

// FilterConfig contains default filter settings
type FilterConfig struct {
	IncludeMerges     bool     `json:"include_merges"`      // Include merge commits by default
	DefaultAuthor     string   `json:"default_author"`      // Default author filter
	ExcludePatterns   []string `json:"exclude_patterns"`    // Default file patterns to exclude
	IncludePatterns   []string `json:"include_patterns"`    // Default file patterns to include
	CaseSensitive     bool     `json:"case_sensitive"`      // Case sensitive filtering
	AuthorMatchType   string   `json:"author_match_type"`   // Default author matching type
}

// OutputConfig contains output formatting settings
type OutputConfig struct {
	ColorEnabled  bool   `json:"color_enabled"`  // Enable colored output
	ColorTheme    string `json:"color_theme"`    // Default color theme
	PrettyPrint   bool   `json:"pretty_print"`   // Pretty print JSON output
	IncludeMetadata bool `json:"include_metadata"` // Include metadata in output
	DateFormat    string `json:"date_format"`    // Default date format
	TimeFormat    string `json:"time_format"`    // Default time format
}

// PerformanceConfig contains performance-related settings
type PerformanceConfig struct {
	MaxCommits        int           `json:"max_commits"`         // Maximum commits to process
	ChunkSize         int           `json:"chunk_size"`          // Processing chunk size
	CacheEnabled      bool          `json:"cache_enabled"`       // Enable caching
	CacheTTL          time.Duration `json:"cache_ttl"`           // Cache time-to-live
	ParallelProcessing bool         `json:"parallel_processing"` // Enable parallel processing
	MaxWorkers        int           `json:"max_workers"`         // Maximum worker goroutines
}

// GUIConfig contains GUI-specific settings
type GUIConfig struct {
	DefaultView       string `json:"default_view"`        // Default GUI view
	RefreshInterval   int    `json:"refresh_interval"`    // Auto-refresh interval (seconds)
	KeyBindings       map[string]string `json:"key_bindings"` // Custom key bindings
	ShowHelp          bool   `json:"show_help"`           // Show help on startup
	ContribGraphWidth int    `json:"contrib_graph_width"` // Contribution graph width
}

// ConfigManager manages application configuration
type ConfigManager struct {
	config     *Config
	configPath string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		config: getDefaultConfig(),
	}
}

// NewConfigManagerWithPath creates a new configuration manager with a specific config path
func NewConfigManagerWithPath(configPath string) *ConfigManager {
	return &ConfigManager{
		config:     getDefaultConfig(),
		configPath: configPath,
	}
}

// getDefaultConfig returns the default configuration
func getDefaultConfig() *Config {
	return &Config{
		Defaults: DefaultConfig{
			Command:      "contrib",
			DateRange:    "1 year ago",
			Format:       "terminal",
			ShowProgress: false,
			RepoPath:     ".",
		},
		Filters: FilterConfig{
			IncludeMerges:   true,
			DefaultAuthor:   "",
			ExcludePatterns: []string{},
			IncludePatterns: []string{},
			CaseSensitive:   false,
			AuthorMatchType: "contains",
		},
		Output: OutputConfig{
			ColorEnabled:    true,
			ColorTheme:      "github",
			PrettyPrint:     true,
			IncludeMetadata: true,
			DateFormat:      "2006-01-02",
			TimeFormat:      "15:04:05",
		},
		Performance: PerformanceConfig{
			MaxCommits:         10000,
			ChunkSize:          1000,
			CacheEnabled:       false,
			CacheTTL:           time.Hour * 24,
			ParallelProcessing: true,
			MaxWorkers:         4,
		},
		GUI: GUIConfig{
			DefaultView:       "contrib",
			RefreshInterval:   0, // No auto-refresh by default
			KeyBindings:       getDefaultKeyBindings(),
			ShowHelp:          false,
			ContribGraphWidth: 53, // Standard GitHub width
		},
	}
}

// getDefaultKeyBindings returns default key bindings for GUI
func getDefaultKeyBindings() map[string]string {
	return map[string]string{
		"quit":         "q",
		"help":         "?",
		"contrib_view": "c",
		"stats_view":   "s",
		"team_view":    "t",
		"health_view":  "h",
		"refresh":      "r",
		"export":       "e",
	}
}

// Load loads configuration from file
func (cm *ConfigManager) Load() error {
	configPath := cm.getConfigPath()

	// If config file doesn't exist, use defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil // Use default config
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Merge with defaults to ensure all fields are set
	cm.config = cm.mergeWithDefaults(&config)

	return nil
}

// Save saves configuration to file
func (cm *ConfigManager) Save() error {
	configPath := cm.getConfigPath()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// SetConfig sets the configuration
func (cm *ConfigManager) SetConfig(config *Config) {
	cm.config = config
}

// UpdateDefaults updates default settings
func (cm *ConfigManager) UpdateDefaults(defaults DefaultConfig) {
	cm.config.Defaults = defaults
}

// UpdateFilters updates filter settings
func (cm *ConfigManager) UpdateFilters(filters FilterConfig) {
	cm.config.Filters = filters
}

// UpdateOutput updates output settings
func (cm *ConfigManager) UpdateOutput(output OutputConfig) {
	cm.config.Output = output
}

// UpdatePerformance updates performance settings
func (cm *ConfigManager) UpdatePerformance(performance PerformanceConfig) {
	cm.config.Performance = performance
}

// UpdateGUI updates GUI settings
func (cm *ConfigManager) UpdateGUI(gui GUIConfig) {
	cm.config.GUI = gui
}

// getConfigPath returns the configuration file path
func (cm *ConfigManager) getConfigPath() string {
	if cm.configPath != "" {
		return cm.configPath
	}

	// Try to get user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fallback to current directory
			return "git-stats-config.json"
		}
		return filepath.Join(homeDir, ".git-stats", "config.json")
	}

	return filepath.Join(configDir, "git-stats", "config.json")
}

// mergeWithDefaults merges loaded config with defaults
func (cm *ConfigManager) mergeWithDefaults(loaded *Config) *Config {
	defaults := getDefaultConfig()

	// Merge defaults
	if loaded.Defaults.Command == "" {
		loaded.Defaults.Command = defaults.Defaults.Command
	}
	if loaded.Defaults.DateRange == "" {
		loaded.Defaults.DateRange = defaults.Defaults.DateRange
	}
	if loaded.Defaults.Format == "" {
		loaded.Defaults.Format = defaults.Defaults.Format
	}
	if loaded.Defaults.RepoPath == "" {
		loaded.Defaults.RepoPath = defaults.Defaults.RepoPath
	}

	// Merge filter settings
	if loaded.Filters.AuthorMatchType == "" {
		loaded.Filters.AuthorMatchType = defaults.Filters.AuthorMatchType
	}
	if loaded.Filters.ExcludePatterns == nil {
		loaded.Filters.ExcludePatterns = defaults.Filters.ExcludePatterns
	}
	if loaded.Filters.IncludePatterns == nil {
		loaded.Filters.IncludePatterns = defaults.Filters.IncludePatterns
	}

	// Merge output settings
	if loaded.Output.ColorTheme == "" {
		loaded.Output.ColorTheme = defaults.Output.ColorTheme
	}
	if loaded.Output.DateFormat == "" {
		loaded.Output.DateFormat = defaults.Output.DateFormat
	}
	if loaded.Output.TimeFormat == "" {
		loaded.Output.TimeFormat = defaults.Output.TimeFormat
	}

	// Merge performance settings
	if loaded.Performance.MaxCommits == 0 {
		loaded.Performance.MaxCommits = defaults.Performance.MaxCommits
	}
	if loaded.Performance.ChunkSize == 0 {
		loaded.Performance.ChunkSize = defaults.Performance.ChunkSize
	}
	if loaded.Performance.CacheTTL == 0 {
		loaded.Performance.CacheTTL = defaults.Performance.CacheTTL
	}
	if loaded.Performance.MaxWorkers == 0 {
		loaded.Performance.MaxWorkers = defaults.Performance.MaxWorkers
	}

	// Merge GUI settings
	if loaded.GUI.DefaultView == "" {
		loaded.GUI.DefaultView = defaults.GUI.DefaultView
	}
	if loaded.GUI.KeyBindings == nil {
		loaded.GUI.KeyBindings = defaults.GUI.KeyBindings
	}
	if loaded.GUI.ContribGraphWidth == 0 {
		loaded.GUI.ContribGraphWidth = defaults.GUI.ContribGraphWidth
	}

	return loaded
}

// Validate validates the configuration
func (cm *ConfigManager) Validate() error {
	config := cm.config

	// Validate defaults
	validCommands := []string{"contrib", "summary", "contributors", "health"}
	if !contains(validCommands, config.Defaults.Command) {
		return fmt.Errorf("invalid default command: %s", config.Defaults.Command)
	}

	validFormats := []string{"terminal", "json", "csv"}
	if !contains(validFormats, config.Defaults.Format) {
		return fmt.Errorf("invalid default format: %s", config.Defaults.Format)
	}

	// Validate filter settings
	validMatchTypes := []string{"exact", "contains", "regex", "email", "name"}
	if !contains(validMatchTypes, config.Filters.AuthorMatchType) {
		return fmt.Errorf("invalid author match type: %s", config.Filters.AuthorMatchType)
	}

	// Validate output settings
	validThemes := []string{"github", "blue", "fire", "mono"}
	if !contains(validThemes, config.Output.ColorTheme) {
		return fmt.Errorf("invalid color theme: %s", config.Output.ColorTheme)
	}

	// Validate performance settings
	if config.Performance.MaxCommits <= 0 {
		return fmt.Errorf("max commits must be positive: %d", config.Performance.MaxCommits)
	}
	if config.Performance.ChunkSize <= 0 {
		return fmt.Errorf("chunk size must be positive: %d", config.Performance.ChunkSize)
	}
	if config.Performance.MaxWorkers <= 0 {
		return fmt.Errorf("max workers must be positive: %d", config.Performance.MaxWorkers)
	}

	// Validate GUI settings
	validViews := []string{"contrib", "summary", "contributors", "health"}
	if !contains(validViews, config.GUI.DefaultView) {
		return fmt.Errorf("invalid default GUI view: %s", config.GUI.DefaultView)
	}

	if config.GUI.ContribGraphWidth <= 0 {
		return fmt.Errorf("contribution graph width must be positive: %d", config.GUI.ContribGraphWidth)
	}

	return nil
}

// Reset resets configuration to defaults
func (cm *ConfigManager) Reset() {
	cm.config = getDefaultConfig()
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetConfigPath returns the current config file path
func (cm *ConfigManager) GetConfigPath() string {
	return cm.getConfigPath()
}

// ExportConfig exports configuration to a specific file
func (cm *ConfigManager) ExportConfig(path string) error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ImportConfig imports configuration from a specific file
func (cm *ConfigManager) ImportConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Merge with defaults and validate
	cm.config = cm.mergeWithDefaults(&config)

	if err := cm.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return nil
}
