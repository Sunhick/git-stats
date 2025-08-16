// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - File output handler unit tests

package formatters

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git-stats/formatters"
	"git-stats/models"
)

func TestNewFileOutputHandler(t *testing.T) {
	handler := formatters.NewFileOutputHandler(true, formatters.OverwriteModeBackup)
	if handler == nil {
		t.Fatal("NewFileOutputHandler should return a non-nil handler")
	}
}

func TestFileOutputHandler_WriteToFile(t *testing.T) {
	handler := formatters.NewFileOutputHandler(false, formatters.OverwriteModeReplace)

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-stats-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name       string
		data       []byte
		config     formatters.FileOutputConfig
		expectErr  bool
		setupFile  bool
		fileContent string
	}{
		{
			name: "write new file",
			data: []byte("test data"),
			config: formatters.FileOutputConfig{
				OutputPath:    filepath.Join(tempDir, "test.txt"),
				CreateDirs:    true,
				OverwriteMode: formatters.OverwriteModeReplace,
				FileMode:      0644,
			},
			expectErr: false,
		},
		{
			name: "create directories",
			data: []byte("test data"),
			config: formatters.FileOutputConfig{
				OutputPath:    filepath.Join(tempDir, "subdir", "test.txt"),
				CreateDirs:    true,
				OverwriteMode: formatters.OverwriteModeReplace,
				FileMode:      0644,
			},
			expectErr: false,
		},
		{
			name: "error on existing file",
			data: []byte("new data"),
			config: formatters.FileOutputConfig{
				OutputPath:    filepath.Join(tempDir, "existing.txt"),
				OverwriteMode: formatters.OverwriteModeError,
				FileMode:      0644,
			},
			expectErr:   true,
			setupFile:   true,
			fileContent: "existing content",
		},
		{
			name: "replace existing file",
			data: []byte("new data"),
			config: formatters.FileOutputConfig{
				OutputPath:    filepath.Join(tempDir, "replace.txt"),
				OverwriteMode: formatters.OverwriteModeReplace,
				FileMode:      0644,
			},
			expectErr:   false,
			setupFile:   true,
			fileContent: "old content",
		},
		{
			name: "empty output path",
			data: []byte("test data"),
			config: formatters.FileOutputConfig{
				OutputPath: "",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup existing file if needed
			if tt.setupFile {
				if err := os.WriteFile(tt.config.OutputPath, []byte(tt.fileContent), 0644); err != nil {
					t.Fatalf("Failed to setup test file: %v", err)
				}
			}

			err := handler.WriteToFile(tt.data, tt.config)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("WriteToFile() error = %v", err)
			}

			// Verify file was written correctly
			if tt.config.OutputPath != "" {
				content, err := os.ReadFile(tt.config.OutputPath)
				if err != nil {
					t.Fatalf("Failed to read written file: %v", err)
				}

				if string(content) != string(tt.data) {
					t.Errorf("Expected file content '%s', got '%s'", string(tt.data), string(content))
				}
			}
		})
	}
}

func TestFileOutputHandler_WriteFormattedOutput(t *testing.T) {
	handler := formatters.NewFileOutputHandler(false, formatters.OverwriteModeReplace)
	formatter := formatters.NewJSONFormatter()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-stats-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testData := createTestAnalysisResult()
	formatConfig := models.FormatConfig{
		Format:   "json",
		Pretty:   true,
		Metadata: true,
	}
	outputConfig := formatters.FileOutputConfig{
		OutputPath:    filepath.Join(tempDir, "output.json"),
		CreateDirs:    true,
		OverwriteMode: formatters.OverwriteModeReplace,
		FileMode:      0644,
	}

	err = handler.WriteFormattedOutput(testData, formatter, formatConfig, outputConfig)
	if err != nil {
		t.Fatalf("WriteFormattedOutput() error = %v", err)
	}

	// Verify file was created and contains JSON
	content, err := os.ReadFile(outputConfig.OutputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Output file is empty")
	}

	// Verify it's JSON by checking for basic structure
	contentStr := string(content)
	if !strings.Contains(contentStr, "{") || !strings.Contains(contentStr, "}") {
		t.Error("Output doesn't appear to be JSON")
	}
}

func TestFileOutputHandler_WriteMultipleFormats(t *testing.T) {
	handler := formatters.NewFileOutputHandler(false, formatters.OverwriteModeReplace)

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-stats-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testData := createTestAnalysisResult()

	formats := map[string]formatters.FormatterConfig{
		"json": {
			Formatter: formatters.NewJSONFormatter(),
			FormatConfig: models.FormatConfig{
				Format:   "json",
				Pretty:   true,
				Metadata: true,
			},
			OutputConfig: formatters.FileOutputConfig{
				OutputPath:    filepath.Join(tempDir, "output.json"),
				CreateDirs:    true,
				OverwriteMode: formatters.OverwriteModeReplace,
				FileMode:      0644,
			},
		},
		"csv": {
			Formatter: formatters.NewCSVFormatter(),
			FormatConfig: models.FormatConfig{
				Format:   "csv",
				Metadata: true,
			},
			OutputConfig: formatters.FileOutputConfig{
				OutputPath:    filepath.Join(tempDir, "output.csv"),
				CreateDirs:    true,
				OverwriteMode: formatters.OverwriteModeReplace,
				FileMode:      0644,
			},
		},
	}

	err = handler.WriteMultipleFormats(testData, formats)
	if err != nil {
		t.Fatalf("WriteMultipleFormats() error = %v", err)
	}

	// Verify both files were created
	for formatName, config := range formats {
		if _, err := os.Stat(config.OutputConfig.OutputPath); os.IsNotExist(err) {
			t.Errorf("Output file for %s format was not created", formatName)
		}
	}
}

func TestFileOutputHandler_BackupMode(t *testing.T) {
	handler := formatters.NewFileOutputHandler(true, formatters.OverwriteModeBackup)

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-stats-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	outputPath := filepath.Join(tempDir, "test.txt")
	originalContent := "original content"
	newContent := "new content"

	// Create original file
	if err := os.WriteFile(outputPath, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create original file: %v", err)
	}

	config := formatters.FileOutputConfig{
		OutputPath:    outputPath,
		OverwriteMode: formatters.OverwriteModeBackup,
		FileMode:      0644,
	}

	// Write new content with backup
	err = handler.WriteToFile([]byte(newContent), config)
	if err != nil {
		t.Fatalf("WriteToFile() error = %v", err)
	}

	// Verify original file was overwritten
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != newContent {
		t.Errorf("Expected new content '%s', got '%s'", newContent, string(content))
	}

	// Verify backup was created
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	backupFound := false
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "test.txt.backup_") {
			backupFound = true
			// Verify backup content
			backupContent, err := os.ReadFile(filepath.Join(tempDir, file.Name()))
			if err != nil {
				t.Fatalf("Failed to read backup file: %v", err)
			}
			if string(backupContent) != originalContent {
				t.Errorf("Expected backup content '%s', got '%s'", originalContent, string(backupContent))
			}
			break
		}
	}

	if !backupFound {
		t.Error("Backup file was not created")
	}
}

func TestFileOutputHandler_AppendMode(t *testing.T) {
	handler := formatters.NewFileOutputHandler(false, formatters.OverwriteModeAppend)

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-stats-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	outputPath := filepath.Join(tempDir, "test.txt")
	originalContent := "original content"
	appendContent := "appended content"

	// Create original file
	if err := os.WriteFile(outputPath, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create original file: %v", err)
	}

	config := formatters.FileOutputConfig{
		OutputPath:    outputPath,
		OverwriteMode: formatters.OverwriteModeAppend,
		FileMode:      0644,
	}

	// Append content
	err = handler.WriteToFile([]byte(appendContent), config)
	if err != nil {
		t.Fatalf("WriteToFile() error = %v", err)
	}

	// Verify content was appended
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, originalContent) {
		t.Error("Original content not found in file")
	}
	if !strings.Contains(contentStr, appendContent) {
		t.Error("Appended content not found in file")
	}
	if !strings.Contains(contentStr, "Appended at") {
		t.Error("Append separator not found in file")
	}
}

func TestFileOutputHandler_GetSafeOutputPath(t *testing.T) {
	handler := formatters.NewFileOutputHandler(false, formatters.OverwriteModeReplace)

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-stats-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	basePath := filepath.Join(tempDir, "test.txt")

	// Test with non-existing file
	safePath, err := handler.GetSafeOutputPath(basePath)
	if err != nil {
		t.Fatalf("GetSafeOutputPath() error = %v", err)
	}
	if safePath != basePath {
		t.Errorf("Expected safe path '%s', got '%s'", basePath, safePath)
	}

	// Create the base file
	if err := os.WriteFile(basePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with existing file
	safePath, err = handler.GetSafeOutputPath(basePath)
	if err != nil {
		t.Fatalf("GetSafeOutputPath() error = %v", err)
	}

	expectedPath := filepath.Join(tempDir, "test_1.txt")
	if safePath != expectedPath {
		t.Errorf("Expected safe path '%s', got '%s'", expectedPath, safePath)
	}
}

func TestFileOutputHandler_ValidateOutputPath(t *testing.T) {
	handler := formatters.NewFileOutputHandler(false, formatters.OverwriteModeReplace)

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-stats-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name      string
		path      string
		expectErr bool
	}{
		{
			name:      "valid path",
			path:      filepath.Join(tempDir, "test.txt"),
			expectErr: false,
		},
		{
			name:      "empty path",
			path:      "",
			expectErr: true,
		},
		{
			name:      "non-existent directory",
			path:      filepath.Join(tempDir, "nonexistent", "test.txt"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateOutputPath(tt.path)
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateOutputPath() error = %v", err)
				}
			}
		})
	}
}

func TestFileOutputHandler_GetOutputStats(t *testing.T) {
	handler := formatters.NewFileOutputHandler(false, formatters.OverwriteModeReplace)

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-stats-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test content"

	// Create test file
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	stats, err := handler.GetOutputStats(testFile)
	if err != nil {
		t.Fatalf("GetOutputStats() error = %v", err)
	}

	if stats.Path != testFile {
		t.Errorf("Expected path '%s', got '%s'", testFile, stats.Path)
	}
	if stats.Size != int64(len(testContent)) {
		t.Errorf("Expected size %d, got %d", len(testContent), stats.Size)
	}
	if stats.IsDirectory {
		t.Error("Expected IsDirectory to be false")
	}
	if time.Since(stats.ModTime) > time.Minute {
		t.Error("ModTime seems too old")
	}
}

func TestFileOutputHandler_ErrorHandling(t *testing.T) {
	handler := formatters.NewFileOutputHandler(false, formatters.OverwriteModeReplace)

	// Test WriteFormattedOutput with nil data
	err := handler.WriteFormattedOutput(nil, formatters.NewJSONFormatter(), models.FormatConfig{}, formatters.FileOutputConfig{})
	if err == nil {
		t.Error("Expected error when data is nil")
	}

	// Test WriteFormattedOutput with nil formatter
	testData := createTestAnalysisResult()
	err = handler.WriteFormattedOutput(testData, nil, models.FormatConfig{}, formatters.FileOutputConfig{})
	if err == nil {
		t.Error("Expected error when formatter is nil")
	}

	// Test WriteMultipleFormats with nil data
	err = handler.WriteMultipleFormats(nil, map[string]formatters.FormatterConfig{})
	if err == nil {
		t.Error("Expected error when data is nil")
	}
}
// Helper functions for testing (shared with other test files)

func createTestAnalysisResult() *models.AnalysisResult {
	return &models.AnalysisResult{
		Repository:    createTestRepository(),
		Summary:       createTestSummary(),
		Contributors:  createTestContributors(),
		ContribGraph:  createTestContributionGraph(),
		HealthMetrics: createTestHealthMetrics(),
		TimeRange:     models.TimeRange{Start: time.Now().AddDate(-1, 0, 0), End: time.Now()},
	}
}

func createTestRepository() *models.RepositoryInfo {
	return &models.RepositoryInfo{
		Path:         "/test/repo",
		Name:         "test-repo",
		TotalCommits: 100,
		FirstCommit:  time.Now().AddDate(-1, 0, 0),
		LastCommit:   time.Now(),
		Branches:     []string{"main", "develop", "feature/test"},
	}
}

func createTestSummary() *models.StatsSummary {
	return &models.StatsSummary{
		TotalCommits:     100,
		TotalInsertions:  5000,
		TotalDeletions:   2000,
		FilesChanged:     50,
		ActiveDays:       30,
		AvgCommitsPerDay: 3.33,
		CommitsByHour:    map[int]int{9: 10, 14: 15, 18: 8},
		CommitsByWeekday: map[time.Weekday]int{time.Monday: 20, time.Friday: 25},
		TopFiles: []models.FileStats{
			{Path: "main.go", Commits: 15, Insertions: 500, Deletions: 100},
			{Path: "utils.go", Commits: 10, Insertions: 300, Deletions: 50},
		},
		TopFileTypes: []models.FileTypeStats{
			{Extension: "go", Files: 10, Commits: 80, Lines: 4000},
			{Extension: "md", Files: 3, Commits: 20, Lines: 1000},
		},
	}
}

func createTestContributors() []models.Contributor {
	return []models.Contributor{
		{
			Name:            "John Doe",
			Email:           "john@example.com",
			TotalCommits:    50,
			TotalInsertions: 2500,
			TotalDeletions:  1000,
			FirstCommit:     time.Now().AddDate(-1, 0, 0),
			LastCommit:      time.Now().AddDate(0, 0, -1),
			ActiveDays:      20,
			CommitsByDay:    map[string]int{"2023-12-01": 5, "2023-12-02": 3},
			CommitsByHour:   map[int]int{9: 5, 14: 8},
			CommitsByWeekday: map[int]int{1: 10, 5: 15},
			FileTypes:       map[string]int{"go": 40, "md": 10},
			TopFiles:        []string{"main.go", "utils.go"},
		},
		{
			Name:            "Jane Smith",
			Email:           "jane@example.com",
			TotalCommits:    50,
			TotalInsertions: 2500,
			TotalDeletions:  1000,
			FirstCommit:     time.Now().AddDate(-1, 0, 0),
			LastCommit:      time.Now(),
			ActiveDays:      25,
			CommitsByDay:    map[string]int{"2023-12-01": 3, "2023-12-02": 4},
			CommitsByHour:   map[int]int{10: 6, 15: 9},
			CommitsByWeekday: map[int]int{2: 12, 4: 18},
			FileTypes:       map[string]int{"go": 35, "js": 15},
			TopFiles:        []string{"server.go", "client.js"},
		},
	}
}

func createTestContributionGraph() *models.ContributionGraph {
	return &models.ContributionGraph{
		StartDate:    time.Now().AddDate(-1, 0, 0),
		EndDate:      time.Now(),
		DailyCommits: map[string]int{"2023-12-01": 5, "2023-12-02": 3, "2023-12-03": 7},
		MaxCommits:   7,
		TotalCommits: 15,
	}
}

func createTestHealthMetrics() *models.HealthMetrics {
	return &models.HealthMetrics{
		RepositoryAge:      365 * 24 * time.Hour,
		CommitFrequency:    3.33,
		ContributorCount:   2,
		ActiveContributors: 2,
		BranchCount:        3,
		ActivityTrend:      "stable",
		MonthlyGrowth: []models.MonthlyStats{
			{Month: time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC), Commits: 45, Authors: 2},
			{Month: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC), Commits: 55, Authors: 2},
		},
	}
}
