// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - Models package for configuration structures

package models

// No imports needed for this file currently

// AnalysisConfig contains configuration for statistical analysis
type AnalysisConfig struct {
	TimeRange    TimeRange
	AuthorFilter string
	Limit        int
	IncludeMerges bool
}

// RenderConfig contains configuration for visualization rendering
type RenderConfig struct {
	Width       int
	Height      int
	ColorScheme string
	ShowLegend  bool
	Interactive bool
}

// FormatConfig contains configuration for output formatting
type FormatConfig struct {
	Format     string // json, csv, terminal
	OutputFile string
	Pretty     bool
	Metadata   bool
}

// SystemConfig contains system-wide configuration
type SystemConfig struct {
	DefaultDateRange    string
	MaxCommitsToProcess int
	CacheEnabled        bool
	PluginsEnabled      []string
	OutputDefaults      FormatConfig
}
