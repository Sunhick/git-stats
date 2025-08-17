# Enhanced Filtering and Configuration System

This document describes the enhanced filtering and configuration system implemented for task 7.2 of the GitHub Activity Visualization feature.

## Overview

The filtering and configuration system provides comprehensive capabilities for filtering git commits and managing application settings. It supports multiple filter types, flexible configuration management, and extensive customization options.

## Features Implemented

### 1. Date Range Filtering (--since, --until)

**Requirements Addressed:** 4.1, 4.3

The system supports comprehensive date range filtering with multiple date formats:

#### Absolute Date Formats
- `YYYY-MM-DD` (e.g., `2024-01-15`)
- `YYYY-MM-DD HH:MM:SS` (e.g., `2024-01-15 14:30:00`)
- `YYYY-MM-DDTHH:MM:SSZ` (ISO format)
- `MM/DD/YYYY` (US format)
- `DD-MM-YYYY` (European format)

#### Relative Date Formats
- `today`, `yesterday`
- `this week`, `last week`
- `this month`, `last month`
- `this year`, `last year`
- `this quarter`, `last quarter`
- `N days ago`, `N weeks ago`, `N months ago`, `N years ago`
- Natural language: `a week ago`, `an hour ago`

#### Usage Examples
```bash
git-stats --since "2024-01-01" --until "2024-12-31"
git-stats --since "1 month ago"
git-stats --since "last week" --until "today"
```

### 2. Author Filtering with Partial Matching

**Requirements Addressed:** 4.2, 4.4

Advanced author filtering with multiple matching strategies:

#### Match Types
- **Contains Match**: Partial matching in name or email (default)
- **Exact Match**: Exact string matching
- **Regex Match**: Regular expression matching
- **Email Match**: Match only in email field
- **Name Match**: Match only in name field

#### Features
- Case-sensitive and case-insensitive matching
- Partial name and email matching
- Domain-based filtering (e.g., `@company.com`)
- Complex regex patterns

#### Usage Examples
```bash
git-stats --author "john"                    # Partial match
git-stats --author "john@example.com"        # Email match
git-stats --author "@company.com"            # Domain match
```

### 3. Enhanced Configuration Management

**Requirements Addressed:** All filtering requirements

Comprehensive configuration system with:

#### Configuration Categories
- **Default Settings**: Command, date range, format, etc.
- **Filter Settings**: Author matching, file patterns, size limits
- **Output Settings**: Colors, themes, formatting
- **Performance Settings**: Limits, caching, parallel processing
- **GUI Settings**: Interface preferences, key bindings

#### Configuration Features
- JSON-based configuration files
- Hierarchical configuration merging
- Validation and error handling
- Import/export capabilities
- CLI override support

### 4. Additional Filter Types

#### Branch Filtering
Filter commits by branch name with support for:
- Exact branch name matching
- Partial branch name matching
- Regex-based branch filtering

#### Message Filtering
Filter commits by commit message content:
- Contains matching
- Starts with / ends with matching
- Regex pattern matching
- Case-sensitive options

#### File Size Filtering
Filter commits based on change size:
- Minimum/maximum insertions
- Minimum/maximum deletions
- Minimum/maximum files changed
- Combined size criteria

### 5. Logical Filter Combination

**Requirements Addressed:** 4.5

All filters can be combined using logical AND operations:
- Date range + author filtering
- Multiple filter types simultaneously
- Configurable filter chains
- Filter priority and ordering

## Architecture

### Core Components

#### Filter Interface
```go
type Filter interface {
    Apply(commits []models.Commit) []models.Commit
    Description() string
}
```

#### Filter Chain
```go
type FilterChain struct {
    filters []Filter
}
```

#### Filter Builder
```go
type FilterBuilder struct {
    configManager *config.ConfigManager
}
```

### Filter Types

1. **DateRangeFilter**: Time-based filtering
2. **AuthorFilter**: Author-based filtering with advanced matching
3. **BranchFilter**: Branch-based filtering
4. **MessageFilter**: Commit message filtering
5. **FileSizeFilter**: Size-based filtering
6. **FilePathFilter**: File path filtering
7. **MergeCommitFilter**: Merge commit inclusion/exclusion
8. **LimitFilter**: Result count limiting

### Configuration Management

#### Configuration Structure
```go
type Config struct {
    Defaults    DefaultConfig     `json:"defaults"`
    Filters     FilterConfig      `json:"filters"`
    Output      OutputConfig      `json:"output"`
    Performance PerformanceConfig `json:"performance"`
    GUI         GUIConfig         `json:"gui"`
}
```

#### Configuration Manager
- Load/save configuration files
- Merge with defaults
- Validate configuration
- Handle CLI overrides

## Usage Examples

### Basic Filtering
```bash
# Date range filtering
git-stats --since "2024-01-01" --until "2024-06-30"

# Author filtering
git-stats --author "john@company.com"

# Combined filtering
git-stats --since "1 month ago" --author "john" --limit 100
```

### Advanced Configuration
```json
{
  "defaults": {
    "command": "contrib",
    "date_range": "1 year ago",
    "format": "terminal"
  },
  "filters": {
    "include_merges": true,
    "default_author": "",
    "author_match_type": "contains",
    "case_sensitive": false,
    "min_insertions": 0,
    "max_insertions": 1000
  }
}
```

### Programmatic Usage
```go
// Create filter chain
builder := filters.NewFilterBuilder(configManager)
chain, err := builder.BuildFromCLIConfig(cliConfig)

// Apply filters
filtered := chain.Apply(commits)

// Get filter summary
summary := builder.GetFilterSummary(chain)
```

## Testing

### Test Coverage

The implementation includes comprehensive test coverage:

#### Unit Tests
- **Filter Tests**: Individual filter functionality
- **Builder Tests**: Filter chain building
- **Config Tests**: Configuration management
- **Integration Tests**: End-to-end scenarios

#### Test Categories
1. **Basic Functionality**: Core filter operations
2. **Edge Cases**: Boundary conditions and error handling
3. **Integration**: Multiple filter combinations
4. **Configuration**: Config loading, validation, merging
5. **Performance**: Large dataset handling

#### Test Files
- `tests/filters/filters_test.go`: Core filter tests
- `tests/filters/builder_test.go`: Filter builder tests
- `tests/filters/comprehensive_filtering_test.go`: Comprehensive scenarios
- `tests/filters/enhanced_filtering_test.go`: Enhanced features tests
- `tests/config/config_test.go`: Configuration tests
- `tests/integration/enhanced_integration_test.go`: Integration tests

### Running Tests
```bash
# Run all filter tests
go test ./tests/filters/... -v

# Run configuration tests
go test ./tests/config/... -v

# Run integration tests
go test ./tests/integration/... -v
```

## Performance Considerations

### Optimization Features
1. **Filter Chain Optimization**: Efficient filter ordering
2. **Early Termination**: Stop processing when limits are reached
3. **Memory Management**: Streaming for large datasets
4. **Caching**: Configuration and filter result caching

### Performance Settings
```json
{
  "performance": {
    "max_commits": 10000,
    "chunk_size": 1000,
    "cache_enabled": true,
    "parallel_processing": true,
    "max_workers": 4
  }
}
```

## Error Handling

### Error Types
- **Validation Errors**: Invalid configuration or parameters
- **Parse Errors**: Date parsing, regex compilation
- **Runtime Errors**: Filter application failures

### Error Recovery
- Graceful degradation for non-critical errors
- User-friendly error messages
- Fallback to default configurations

## Future Enhancements

### Planned Features
1. **Tag Filtering**: Filter by git tags
2. **File Type Filtering**: Filter by file extensions
3. **Commit Type Filtering**: Filter by conventional commit types
4. **Performance Metrics**: Filter performance monitoring
5. **Custom Filter Plugins**: Extensible filter system

### Configuration Migration
- Automatic configuration migration
- Backward compatibility support
- Configuration versioning

## API Reference

### Key Functions

#### Filter Creation
```go
// Date range filter
filter := filters.NewDateRangeFilter(since, until)

// Author filter with options
filter, err := filters.NewAuthorFilterWithOptions(pattern, matchType, caseSensitive)

// Message filter
filter, err := filters.NewMessageFilter(pattern, matchType, caseSensitive)

// Size filter
filter := filters.NewFileSizeFilter(minIns, maxIns, minDel, maxDel, minFiles, maxFiles)
```

#### Filter Chain Management
```go
// Create chain
chain := filters.NewFilterChain()

// Add filters
chain.Add(dateFilter)
chain.Add(authorFilter)

// Apply chain
filtered := chain.Apply(commits)
```

#### Configuration Management
```go
// Create manager
manager := config.NewConfigManager()

// Load configuration
err := manager.Load()

// Update configuration
manager.UpdateFilters(newFilters)

// Save configuration
err := manager.Save()
```

## Conclusion

The enhanced filtering and configuration system provides a robust, flexible, and extensible foundation for git repository analysis. It fully addresses the requirements for date range filtering, author filtering with partial matching, and comprehensive configuration management, while adding significant additional capabilities for advanced use cases.

The system is designed with performance, usability, and maintainability in mind, providing both simple interfaces for basic use cases and powerful advanced features for complex filtering scenarios.
