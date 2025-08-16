// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Unit tests for configuration management

package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"git-stats/config"
)

func TestConfigManager_DefaultConfig(t *testing.T) {
	manager := config.NewConfigManager()
	cfg := manager.GetConfig()

	// Test default values
	if cfg.Defaults.Command != "contrib" {
		t.Errorf("Expected default command 'contrib', got '%s'", cfg.Defaults.Command)
	}

	if cfg.Defaults.Format != "terminal" {
		t.Errorf("Expected default format 'terminal', got '%s'", cfg.Defaults.Format)
	}

	if cfg.Defaults.DateRange != "1 year ago" {
		t.Errorf("Expected default date range '1 year ago', got '%s'", cfg.Defaults.DateRange)
	}

	if !cfg.Filters.IncludeMerges {
		t.Error("Expected default include merges to be true")
	}

	if cfg.Output.ColorTheme != "github" {
		t.Errorf("Expected default color theme 'github', got '%s'", cfg.Output.ColorTheme)
	}

	if cfg.Performance.MaxCommits != 10000 {
		t.Errorf("Expected default max commits 10000, got %d", cfg.Performance.MaxCommits)
	}
}

func TestConfigManager_SaveAndLoad(t *testing.T) {
	// Create temporary directory for test config
	tempDir, err := os.MkdirTemp("", "git-stats-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")
	manager := config.NewConfigManagerWithPath(configPath)

	// Modify config
	cfg := manager.GetConfig()
	cfg.Defaults.Command = "summary"
	cfg.Defaults.Format = "json"
	cfg.Filters.DefaultAuthor = "test@example.com"
	cfg.Output.ColorEnabled = false
	cfg.Performance.MaxCommits = 5000

	// Save config
	err = manager.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Create new manager and load config
	newManager := config.NewConfigManagerWithPath(configPath)
	err = newManager.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded config
	loadedCfg := newManager.GetConfig()
	if loadedCfg.Defaults.Command != "summary" {
		t.Errorf("Expected loaded command 'summary', got '%s'", loadedCfg.Defaults.Command)
	}

	if loadedCfg.Defaults.Format != "json" {
		t.Errorf("Expected loaded format 'json', got '%s'", loadedCfg.Defaults.Format)
	}

	if loadedCfg.Filters.DefaultAuthor != "test@example.com" {
		t.Errorf("Expected loaded author 'test@example.com', got '%s'", loadedCfg.Filters.DefaultAuthor)
	}

	if loadedCfg.Output.ColorEnabled {
		t.Error("Expected loaded color enabled to be false")
	}

	if loadedCfg.Performance.MaxCommits != 5000 {
		t.Errorf("Expected loaded max commits 5000, got %d", loadedCfg.Performance.MaxCommits)
	}
}

func TestConfigManager_LoadNonExistentFile(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "git-stats-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "nonexistent.json")
	manager := config.NewConfigManagerWithPath(configPath)

	// Load should not error for non-existent file
	err = manager.Load()
	if err != nil {
		t.Errorf("Load should not error for non-existent file: %v", err)
	}

	// Should use default config
	cfg := manager.GetConfig()
	if cfg.Defaults.Command != "contrib" {
		t.Errorf("Expected default command when file doesn't exist, got '%s'", cfg.Defaults.Command)
	}
}

func TestConfigManager_Validate(t *testing.T) {
	manager := config.NewConfigManager()

	tests := []struct {
		name      string
		modify    func(*config.Config)
		shouldErr bool
	}{
		{
			name:      "Valid default config",
			modify:    func(c *config.Config) {},
			shouldErr: false,
		},
		{
			name: "Invalid command",
			modify: func(c *config.Config) {
				c.Defaults.Command = "invalid"
			},
			shouldErr: true,
		},
		{
			name: "Invalid format",
			modify: func(c *config.Config) {
				c.Defaults.Format = "invalid"
			},
			shouldErr: true,
		},
		{
			name: "Invalid author match type",
			modify: func(c *config.Config) {
				c.Filters.AuthorMatchType = "invalid"
			},
			shouldErr: true,
		},
		{
			name: "Invalid color theme",
			modify: func(c *config.Config) {
				c.Output.ColorTheme = "invalid"
			},
			shouldErr: true,
		},
		{
			name: "Invalid max commits",
			modify: func(c *config.Config) {
				c.Performance.MaxCommits = -1
			},
			shouldErr: true,
		},
		{
			name: "Invalid chunk size",
			modify: func(c *config.Config) {
				c.Performance.ChunkSize = 0
			},
			shouldErr: true,
		},
		{
			name: "Invalid max workers",
			modify: func(c *config.Config) {
				c.Performance.MaxWorkers = -1
			},
			shouldErr: true,
		},
		{
			name: "Invalid GUI view",
			modify: func(c *config.Config) {
				c.GUI.DefaultView = "invalid"
			},
			shouldErr: true,
		},
		{
			name: "Invalid contrib graph width",
			modify: func(c *config.Config) {
				c.GUI.ContribGraphWidth = 0
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := manager.GetConfig()
			tt.modify(cfg)
			manager.SetConfig(cfg)

			err := manager.Validate()

			if tt.shouldErr && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestConfigManager_UpdateMethods(t *testing.T) {
	manager := config.NewConfigManager()

	// Test UpdateDefaults
	newDefaults := config.DefaultConfig{
		Command:      "summary",
		DateRange:    "6 months ago",
		Format:       "json",
		ShowProgress: true,
		RepoPath:     "/custom/path",
	}
	manager.UpdateDefaults(newDefaults)

	cfg := manager.GetConfig()
	if cfg.Defaults.Command != "summary" {
		t.Errorf("Expected updated command 'summary', got '%s'", cfg.Defaults.Command)
	}

	// Test UpdateFilters
	newFilters := config.FilterConfig{
		IncludeMerges:   false,
		DefaultAuthor:   "admin@example.com",
		ExcludePatterns: []string{"*.log", "*.tmp"},
		CaseSensitive:   true,
		AuthorMatchType: "exact",
	}
	manager.UpdateFilters(newFilters)

	cfg = manager.GetConfig()
	if cfg.Filters.IncludeMerges {
		t.Error("Expected updated include merges to be false")
	}
	if cfg.Filters.DefaultAuthor != "admin@example.com" {
		t.Errorf("Expected updated author 'admin@example.com', got '%s'", cfg.Filters.DefaultAuthor)
	}

	// Test UpdateOutput
	newOutput := config.OutputConfig{
		ColorEnabled:    false,
		ColorTheme:      "blue",
		PrettyPrint:     false,
		IncludeMetadata: false,
		DateFormat:      "01/02/2006",
		TimeFormat:      "03:04 PM",
	}
	manager.UpdateOutput(newOutput)

	cfg = manager.GetConfig()
	if cfg.Output.ColorEnabled {
		t.Error("Expected updated color enabled to be false")
	}
	if cfg.Output.ColorTheme != "blue" {
		t.Errorf("Expected updated color theme 'blue', got '%s'", cfg.Output.ColorTheme)
	}

	// Test UpdatePerformance
	newPerformance := config.PerformanceConfig{
		MaxCommits:         20000,
		ChunkSize:          2000,
		CacheEnabled:       true,
		CacheTTL:           time.Hour * 48,
		ParallelProcessing: false,
		MaxWorkers:         8,
	}
	manager.UpdatePerformance(newPerformance)

	cfg = manager.GetConfig()
	if cfg.Performance.MaxCommits != 20000 {
		t.Errorf("Expected updated max commits 20000, got %d", cfg.Performance.MaxCommits)
	}
	if !cfg.Performance.CacheEnabled {
		t.Error("Expected updated cache enabled to be true")
	}

	// Test UpdateGUI
	newGUI := config.GUIConfig{
		DefaultView:       "health",
		RefreshInterval:   30,
		KeyBindings:       map[string]string{"quit": "x"},
		ShowHelp:          true,
		ContribGraphWidth: 60,
	}
	manager.UpdateGUI(newGUI)

	cfg = manager.GetConfig()
	if cfg.GUI.DefaultView != "health" {
		t.Errorf("Expected updated default view 'health', got '%s'", cfg.GUI.DefaultView)
	}
	if cfg.GUI.RefreshInterval != 30 {
		t.Errorf("Expected updated refresh interval 30, got %d", cfg.GUI.RefreshInterval)
	}
}

func TestConfigManager_Reset(t *testing.T) {
	manager := config.NewConfigManager()

	// Modify config
	cfg := manager.GetConfig()
	cfg.Defaults.Command = "summary"
	cfg.Filters.DefaultAuthor = "test@example.com"
	cfg.Output.ColorEnabled = false

	// Reset to defaults
	manager.Reset()

	// Verify reset
	resetCfg := manager.GetConfig()
	if resetCfg.Defaults.Command != "contrib" {
		t.Errorf("Expected reset command 'contrib', got '%s'", resetCfg.Defaults.Command)
	}
	if resetCfg.Filters.DefaultAuthor != "" {
		t.Errorf("Expected reset author to be empty, got '%s'", resetCfg.Filters.DefaultAuthor)
	}
	if !resetCfg.Output.ColorEnabled {
		t.Error("Expected reset color enabled to be true")
	}
}

func TestConfigManager_ExportImport(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "git-stats-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manager := config.NewConfigManager()

	// Modify config
	cfg := manager.GetConfig()
	cfg.Defaults.Command = "contributors"
	cfg.Filters.DefaultAuthor = "export@example.com"
	cfg.Performance.MaxCommits = 15000

	// Export config
	exportPath := filepath.Join(tempDir, "export.json")
	err = manager.ExportConfig(exportPath)
	if err != nil {
		t.Fatalf("Failed to export config: %v", err)
	}

	// Create new manager and import
	newManager := config.NewConfigManager()
	err = newManager.ImportConfig(exportPath)
	if err != nil {
		t.Fatalf("Failed to import config: %v", err)
	}

	// Verify imported config
	importedCfg := newManager.GetConfig()
	if importedCfg.Defaults.Command != "contributors" {
		t.Errorf("Expected imported command 'contributors', got '%s'", importedCfg.Defaults.Command)
	}
	if importedCfg.Filters.DefaultAuthor != "export@example.com" {
		t.Errorf("Expected imported author 'export@example.com', got '%s'", importedCfg.Filters.DefaultAuthor)
	}
	if importedCfg.Performance.MaxCommits != 15000 {
		t.Errorf("Expected imported max commits 15000, got %d", importedCfg.Performance.MaxCommits)
	}
}

func TestConfigManager_MergeWithDefaults(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "git-stats-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create partial config file
	configPath := filepath.Join(tempDir, "partial.json")
	partialConfig := `{
		"defaults": {
			"command": "summary"
		},
		"filters": {
			"default_author": "partial@example.com"
		}
	}`

	err = os.WriteFile(configPath, []byte(partialConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write partial config: %v", err)
	}

	// Load partial config
	manager := config.NewConfigManagerWithPath(configPath)
	err = manager.Load()
	if err != nil {
		t.Fatalf("Failed to load partial config: %v", err)
	}

	// Verify merged config has both loaded and default values
	cfg := manager.GetConfig()

	// Should have loaded values
	if cfg.Defaults.Command != "summary" {
		t.Errorf("Expected loaded command 'summary', got '%s'", cfg.Defaults.Command)
	}
	if cfg.Filters.DefaultAuthor != "partial@example.com" {
		t.Errorf("Expected loaded author 'partial@example.com', got '%s'", cfg.Filters.DefaultAuthor)
	}

	// Should have default values for unspecified fields
	if cfg.Defaults.Format != "terminal" {
		t.Errorf("Expected default format 'terminal', got '%s'", cfg.Defaults.Format)
	}
	if cfg.Output.ColorTheme != "github" {
		t.Errorf("Expected default color theme 'github', got '%s'", cfg.Output.ColorTheme)
	}
	if cfg.Performance.MaxCommits != 10000 {
		t.Errorf("Expected default max commits 10000, got %d", cfg.Performance.MaxCommits)
	}
}

func TestConfigManager_GetConfigPath(t *testing.T) {
	// Test with custom path
	customPath := "/custom/config.json"
	manager := config.NewConfigManagerWithPath(customPath)

	if manager.GetConfigPath() != customPath {
		t.Errorf("Expected custom path '%s', got '%s'", customPath, manager.GetConfigPath())
	}

	// Test with default path
	defaultManager := config.NewConfigManager()
	defaultPath := defaultManager.GetConfigPath()

	if defaultPath == "" {
		t.Error("Default config path should not be empty")
	}

	// Should contain git-stats in the path
	if !contains(defaultPath, "git-stats") {
		t.Errorf("Default path should contain 'git-stats', got '%s'", defaultPath)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
