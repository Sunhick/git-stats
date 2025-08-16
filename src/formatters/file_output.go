// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - File output handler

package formatters

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"git-stats/models"
)

// FileOutputHandler handles writing formatted output to files
type FileOutputHandler struct {
	backupEnabled bool
	overwriteMode OverwriteMode
}

// OverwriteMode defines how to handle existing files
type OverwriteMode int

const (
	OverwriteModeError OverwriteMode = iota // Error if file exists
	OverwriteModeBackup                     // Create backup before overwriting
	OverwriteModeReplace                    // Replace without backup
	OverwriteModeAppend                     // Append to existing file
)

// FileOutputConfig contains configuration for file output
type FileOutputConfig struct {
	OutputPath    string
	BackupEnabled bool
	OverwriteMode OverwriteMode
	CreateDirs    bool
	FileMode      fs.FileMode
}

// NewFileOutputHandler creates a new file output handler
func NewFileOutputHandler(backupEnabled bool, overwriteMode OverwriteMode) *FileOutputHandler {
	return &FileOutputHandler{
		backupEnabled: backupEnabled,
		overwriteMode: overwriteMode,
	}
}

// WriteToFile writes formatted data to a file with proper error handling
func (foh *FileOutputHandler) WriteToFile(data []byte, config FileOutputConfig) error {
	if config.OutputPath == "" {
		return NewFormatterError("output path cannot be empty")
	}

	// Create directories if requested
	if config.CreateDirs {
		dir := filepath.Dir(config.OutputPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return NewFormatterOperationError("create_directories", fmt.Sprintf("failed to create directories: %v", err))
		}
	}

	// Check if file exists
	fileExists, err := foh.fileExists(config.OutputPath)
	if err != nil {
		return NewFormatterOperationError("file_check", fmt.Sprintf("failed to check file existence: %v", err))
	}

	// Handle existing file based on overwrite mode
	if fileExists {
		switch config.OverwriteMode {
		case OverwriteModeError:
			return NewFormatterOperationError("file_exists", fmt.Sprintf("file already exists: %s", config.OutputPath))
		case OverwriteModeBackup:
			if err := foh.createBackup(config.OutputPath); err != nil {
				return NewFormatterOperationError("backup", fmt.Sprintf("failed to create backup: %v", err))
			}
		case OverwriteModeAppend:
			return foh.appendToFile(data, config)
		case OverwriteModeReplace:
			// Continue with normal write operation
		}
	}

	// Write data to file
	return foh.writeFile(data, config)
}

// WriteFormattedOutput writes analysis results using the specified formatter
func (foh *FileOutputHandler) WriteFormattedOutput(data *models.AnalysisResult, formatter Formatter, formatConfig models.FormatConfig, outputConfig FileOutputConfig) error {
	if data == nil {
		return NewFormatterError("analysis result cannot be nil")
	}
	if formatter == nil {
		return NewFormatterError("formatter cannot be nil")
	}

	// Format the data
	formattedData, err := formatter.Format(data, formatConfig)
	if err != nil {
		return NewFormatterOperationError("format", fmt.Sprintf("failed to format data: %v", err))
	}

	// Write to file
	return foh.WriteToFile(formattedData, outputConfig)
}

// WriteMultipleFormats writes the same data in multiple formats to different files
func (foh *FileOutputHandler) WriteMultipleFormats(data *models.AnalysisResult, formats map[string]FormatterConfig) error {
	if data == nil {
		return NewFormatterError("analysis result cannot be nil")
	}

	for formatName, config := range formats {
		if err := foh.WriteFormattedOutput(data, config.Formatter, config.FormatConfig, config.OutputConfig); err != nil {
			return NewFormatterOperationError("multi_format", fmt.Sprintf("failed to write %s format: %v", formatName, err))
		}
	}

	return nil
}

// FormatterConfig combines formatter, format config, and output config
type FormatterConfig struct {
	Formatter    Formatter
	FormatConfig models.FormatConfig
	OutputConfig FileOutputConfig
}

// fileExists checks if a file exists
func (foh *FileOutputHandler) fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// createBackup creates a backup of the existing file
func (foh *FileOutputHandler) createBackup(path string) error {
	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup_%s", path, timestamp)

	// Read original file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}

	// Write backup file
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// appendToFile appends data to an existing file
func (foh *FileOutputHandler) appendToFile(data []byte, config FileOutputConfig) error {
	file, err := os.OpenFile(config.OutputPath, os.O_APPEND|os.O_WRONLY, config.FileMode)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %w", err)
	}
	defer file.Close()

	// Add separator if file is not empty
	if stat, err := file.Stat(); err == nil && stat.Size() > 0 {
		if _, err := file.WriteString("\n\n# --- Appended at " + time.Now().Format(time.RFC3339) + " ---\n\n"); err != nil {
			return fmt.Errorf("failed to write separator: %w", err)
		}
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to append data: %w", err)
	}

	return nil
}

// writeFile writes data to a file
func (foh *FileOutputHandler) writeFile(data []byte, config FileOutputConfig) error {
	fileMode := config.FileMode
	if fileMode == 0 {
		fileMode = 0644 // Default file mode
	}

	if err := os.WriteFile(config.OutputPath, data, fileMode); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetSafeOutputPath generates a safe output path by checking for conflicts
func (foh *FileOutputHandler) GetSafeOutputPath(basePath string) (string, error) {
	if basePath == "" {
		return "", NewFormatterError("base path cannot be empty")
	}

	// If file doesn't exist, use the base path
	exists, err := foh.fileExists(basePath)
	if err != nil {
		return "", fmt.Errorf("failed to check file existence: %w", err)
	}
	if !exists {
		return basePath, nil
	}

	// Generate alternative paths
	ext := filepath.Ext(basePath)
	name := basePath[:len(basePath)-len(ext)]

	for i := 1; i <= 999; i++ {
		altPath := fmt.Sprintf("%s_%d%s", name, i, ext)
		exists, err := foh.fileExists(altPath)
		if err != nil {
			return "", fmt.Errorf("failed to check alternative path: %w", err)
		}
		if !exists {
			return altPath, nil
		}
	}

	return "", NewFormatterError("unable to generate safe output path after 999 attempts")
}

// ValidateOutputPath validates that the output path is writable
func (foh *FileOutputHandler) ValidateOutputPath(path string) error {
	if path == "" {
		return NewFormatterError("output path cannot be empty")
	}

	// Check if directory exists and is writable
	dir := filepath.Dir(path)
	if stat, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return NewFormatterOperationError("validate_path", fmt.Sprintf("directory does not exist: %s", dir))
		}
		return NewFormatterOperationError("validate_path", fmt.Sprintf("failed to stat directory: %v", err))
	} else if !stat.IsDir() {
		return NewFormatterOperationError("validate_path", fmt.Sprintf("path is not a directory: %s", dir))
	}

	// Test write permissions by creating a temporary file
	tempFile := filepath.Join(dir, ".git-stats-write-test")
	if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
		return NewFormatterOperationError("validate_path", fmt.Sprintf("directory is not writable: %s", dir))
	}
	os.Remove(tempFile) // Clean up

	return nil
}

// GetOutputStats returns statistics about the output operation
func (foh *FileOutputHandler) GetOutputStats(path string) (*OutputStats, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	return &OutputStats{
		Path:         path,
		Size:         stat.Size(),
		ModTime:      stat.ModTime(),
		Mode:         stat.Mode(),
		IsDirectory:  stat.IsDir(),
		AbsolutePath: filepath.Clean(path),
	}, nil
}

// OutputStats contains statistics about output files
type OutputStats struct {
	Path         string
	Size         int64
	ModTime      time.Time
	Mode         fs.FileMode
	IsDirectory  bool
	AbsolutePath string
}
