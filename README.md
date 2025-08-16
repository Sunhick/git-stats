# git-stats

Enhanced terminal git statistics utility with comprehensive analysis and visualization capabilities.

## Features

- **Contribution Analysis**: GitHub-style contribution graphs with activity levels and streak tracking
- **Statistical Analysis**: Comprehensive commit frequency, file statistics, and temporal pattern analysis
- **Health Metrics**: Repository health scoring with activity trends and growth analysis
- **Multiple Output Formats**: Terminal, JSON, CSV, and interactive visualizations
- **Advanced Filtering**: Comprehensive filtering system with date ranges, author matching, file patterns, and configurable options
- **Interactive GUI**: Full-screen ncurses interface with keyboard navigation
- **Comprehensive CLI**: Robust command-line interface with intelligent error handling
- **Flexible Date Parsing**: Support for absolute and relative date formats
- **Configuration Management**: Persistent configuration with user and workspace-level settings
- **Performance Optimized**: Efficient algorithms for large repositories with progress indicators

## Architecture

The application is built with a modular architecture:

- **Models**: Core data structures for commits, contributors, and statistics
- **Git Package**: Repository operations and commit parsing
- **Analyzers**: Statistical analysis engines
  - Contribution Analyzer: Activity graphs and streak calculations
  - Statistics Analyzer: Commit patterns and file statistics
  - Health Analyzer: Repository health metrics and trends
- **CLI**: Comprehensive command-line interface with validation and error handling
  - Parser: Robust argument parsing with comprehensive flag support
  - Validator: Input validation for dates, authors, formats, and repository paths
  - Error Handling: Contextual error messages with helpful suggestions
- **Filtering System**: Advanced filtering with multiple filter types, chaining, and configuration
- **Configuration System**: JSON-based configuration management with validation and defaults
- **Visualization**: Terminal-based charts and interactive ncurses GUI
- **Formatters**: Multiple output format support (terminal, JSON, CSV)

## Quick Start

```shell
# Clone and build
$ git clone <repository-url>
$ cd git-stats
$ make build

# Show help
$ ./git-stats -help

# Basic usage (run from any git repository)
$ ./git-stats                    # Show contribution graph
$ ./git-stats -summary           # Show repository summary
$ ./git-stats -gui               # Launch interactive GUI
```

## Build & Install

```shell
# Build the application
$ make build

# Build with GUI support
$ make build-gui

# Install to /usr/local/bin
$ make install

# Or build and run directly
$ make run

# Run specific commands
$ make run-contrib    # Test contribution graph
$ make run-summary    # Test repository summary
```

## GUI Mode

The application includes a full-featured interactive ncurses GUI with keyboard navigation.

### Building and Running GUI

```shell
# Install GUI dependencies and build
$ make build-gui

# Launch GUI mode
$ make gui

# Or launch with specific views
$ make run-gui-contrib      # Start with contribution graph
$ make run-gui-summary      # Start with summary view
$ make run-gui-contributors # Start with contributors view
$ make run-gui-health       # Start with health metrics

# Development mode with GUI
$ make dev-gui
```

### Offline/Network-Restricted Environments

If you encounter network connectivity issues when building the GUI:

```shell
# Check if GUI dependencies are available
$ make check-gui-deps

# Build GUI in offline mode (uses existing dependencies or falls back to stub)
$ make build-gui-offline

# Launch GUI in offline mode
$ make gui-offline

# Alternative: Build without GUI dependencies (uses stub implementation)
$ make build
$ ./git-stats -gui  # Will show message about GUI mode requiring build tags
```

**Note**: The offline build will automatically fall back to a stub implementation if GUI dependencies (tcell and tview) are not available. The stub provides the same API but displays a message that GUI mode requires building with `-tags gui`.

### GUI Features

- **Multiple Views**: Switch between contribution graph, statistics, contributors, and health metrics
- **Interactive Navigation**: Navigate through dates, months, and years with keyboard shortcuts
- **Detailed Commit Information**: Select dates to view detailed commit information
- **Real-time Updates**: Dynamic content updates as you navigate
- **Comprehensive Help**: Built-in help system with keyboard shortcuts

### GUI Keyboard Shortcuts

#### View Switching
- `c` / `1` / `F1`: Contribution view
- `s` / `2` / `F2`: Statistics view
- `t` / `3` / `F3`: Contributors view
- `H` / `4` / `F4`: Health metrics view
- `Tab`: Cycle views forward
- `Shift+Tab`: Cycle views backward

#### Navigation (Contribution View)
- `←` / `→`: Navigate days
- `↑` / `↓`: Navigate weeks
- `j` / `k`: Navigate weeks
- `h` / `l`: Navigate months
- `H` / `L`: Navigate years
- `Ctrl+←` / `Ctrl+→`: Navigate months
- `Ctrl+↑` / `Ctrl+↓`: Navigate years
- `g`: Go to today
- `G`: Go to first commit

#### Other Controls
- `d`: Toggle detailed commit information
- `r`: Refresh display
- `?`: Toggle help
- `q` / `ESC`: Quit

### GUI Testing

```shell
# Run GUI-specific tests
$ make test-gui-unit         # Unit tests for GUI components
$ make test-gui-integration  # Integration tests for navigation workflows
$ make test-gui-all          # All GUI tests
```

### Troubleshooting GUI Build

**Network/Proxy Issues:**
```shell
# If you get "dial tcp: lookup proxy.golang.org: i/o timeout"
$ make build-gui-offline     # Uses offline build with fallback

# Check dependency availability
$ make check-gui-deps

# Manual dependency installation (if you have network access)
$ cd src && go get github.com/gdamore/tcell/v2@v2.6.0
$ cd src && go get github.com/rivo/tview@v0.0.0-20230826224341-9754ab44dc1c
```

**Corporate/Restricted Networks:**
- Use `make build-gui-offline` for builds without internet access
- The offline build automatically falls back to stub implementation
- GUI tests work with both full and stub implementations

**Dependency Issues:**
- GUI dependencies are only required when building with `-tags gui`
- Without GUI dependencies, the application builds with stub implementation
- All core functionality works without GUI dependencies

## Development

```shell
# Run all tests
$ make test

# Run specific test suites
$ make test-cli          # Test CLI parser and validator
$ make test-analyzers    # Test analysis engines
$ make test-filters      # Test filtering system
$ make test-config       # Test configuration management
$ make test-integration  # Test system integration
$ make test-models       # Test data models
$ make test-git          # Test git operations

# Run with coverage
$ make test-coverage

# Format and check code
$ make check

# Development mode (auto-rebuild)
$ make dev

# Run specific commands for testing
$ make run-contrib       # Test contribution graph
$ make run-summary       # Test repository summary
$ make run-health        # Test health analysis

# Test filtering capabilities
$ make run-filter-date   # Test date filtering
$ make run-filter-author # Test author filtering
$ make run-filter-combined # Test combined filters
$ make run-config-demo   # Test configuration system
```

## Usage

After building, you can add the folder containing `git-stats` to your PATH variable. Then `git stats` will invoke the utility as a git subcommand.

### Basic Usage

```shell
# Show contribution graph (default)
$ git-stats

# Show detailed repository statistics
$ git-stats -summary

# Show contributor statistics
$ git-stats -contributors

# Show repository health metrics
$ git-stats -health

# Launch interactive GUI
$ git-stats -gui
```

### Advanced Filtering System

The tool includes a comprehensive filtering system with multiple filter types and configuration options.

#### Date Range Filtering

```shell
# Absolute dates
$ git-stats -contrib -since "2024-01-01" -until "2024-12-31"
$ git-stats -summary -since "2024-01-15 14:30:00"

# Relative dates
$ git-stats -summary -since "1 month ago"
$ git-stats -health -since "yesterday" -until "today"
$ git-stats -contrib -since "3 weeks ago"
$ git-stats -summary -since "this week"
$ git-stats -health -since "this month"
```

#### Advanced Author Filtering

```shell
# Basic author filtering
$ git-stats -contributors -author "john"
$ git-stats -contrib -author "john@example.com"

# Advanced author matching (configured via config file)
# - Exact matching: "John Doe" matches only "John Doe"
# - Contains matching: "john" matches "John Doe" and "john@example.com"
# - Email-only matching: "example.com" matches only email addresses
# - Name-only matching: "John" matches only names, not emails
# - Regex matching: "^John.*" for pattern-based matching
# - Case-sensitive/insensitive options
```

#### Configuration Management

```shell
# Configuration is automatically managed via JSON files
# User-level config: ~/.config/git-stats/config.json (or ~/.git-stats/config.json)
# Workspace-level config: .git-stats/config.json

# View current configuration
$ git-stats --show-config

# Reset to defaults
$ git-stats --reset-config

# Export configuration
$ git-stats --export-config my-config.json

# Import configuration
$ git-stats --import-config my-config.json
```

#### Filter Combinations

```shell
# Combine multiple filters
$ git-stats -summary -since "1 year ago" -author "john" -format json

# Complex filtering with configuration
# Set default filters in config file for consistent behavior
# CLI options override configuration defaults
```

### Output Formats

```shell
# Export to JSON
$ git-stats -summary -format json

# Export to CSV
$ git-stats -contributors -format csv

# Save to file
$ git-stats -summary -format json -output report.json
```

### Advanced Options

```shell
# Show progress for large repositories
$ git-stats -summary -progress

# Limit commits processed
$ git-stats -contrib -limit 5000

# Analyze specific repository
$ git-stats -summary /path/to/repo

# Combine multiple options
$ git-stats -summary -since "1 year ago" -author "john" -format json -progress
```

## Command Line Options

### Commands
| Flag            | Description                         |
| --------------- | ----------------------------------- |
| `-contrib`      | Show contribution graph (default)   |
| `-summary`      | Show detailed repository statistics |
| `-contributors` | Show contributor statistics         |
| `-health`       | Repository health analysis          |
| `-gui`          | Launch interactive ncurses GUI      |

### Filtering Options
| Flag             | Description                                                |
| ---------------- | ---------------------------------------------------------- |
| `-since <date>`  | Show commits since date (absolute or relative)             |
| `-until <date>`  | Show commits until date (absolute or relative)             |
| `-author <name>` | Filter by author name or email (supports partial matching) |

### Output Options
| Flag             | Description                         |
| ---------------- | ----------------------------------- |
| `-format <fmt>`  | Output format (terminal, json, csv) |
| `-output <file>` | Output file path                    |
| `-progress`      | Show progress indicators            |

### Performance Options
| Flag         | Description                        |
| ------------ | ---------------------------------- |
| `-limit <n>` | Limit number of commits to process |

### Configuration Options
| Flag                     | Description                     |
| ------------------------ | ------------------------------- |
| `--show-config`          | Display current configuration   |
| `--reset-config`         | Reset configuration to defaults |
| `--export-config <file>` | Export configuration to file    |
| `--import-config <file>` | Import configuration from file  |

### Other Options
| Flag        | Description           |
| ----------- | --------------------- |
| `-help, -h` | Show help information |

## Configuration System

The tool uses a hierarchical JSON-based configuration system with automatic defaults merging.

### Configuration Locations

1. **User-level config**: `~/.config/git-stats/config.json` (or `~/.git-stats/config.json`)
2. **Workspace-level config**: `.git-stats/config.json` (workspace-specific settings)

Workspace-level settings override user-level settings, and CLI options override both.

### Configuration Structure

```json
{
  "defaults": {
    "command": "contrib",
    "date_range": "1 year ago",
    "format": "terminal",
    "show_progress": false,
    "repo_path": "."
  },
  "filters": {
    "include_merges": true,
    "default_author": "",
    "exclude_patterns": ["*.log", "*.tmp"],
    "include_patterns": [],
    "case_sensitive": false,
    "author_match_type": "contains"
  },
  "output": {
    "color_enabled": true,
    "color_theme": "github",
    "pretty_print": true,
    "include_metadata": true,
    "date_format": "2006-01-02",
    "time_format": "15:04:05"
  },
  "performance": {
    "max_commits": 10000,
    "chunk_size": 1000,
    "cache_enabled": false,
    "parallel_processing": true,
    "max_workers": 4
  },
  "gui": {
    "default_view": "contrib",
    "refresh_interval": 0,
    "show_help": false,
    "contrib_graph_width": 53
  }
}
```

### Filter Configuration Options

#### Author Match Types
- `"exact"`: Exact string matching
- `"contains"`: Substring matching (default)
- `"regex"`: Regular expression matching
- `"email"`: Email-only matching
- `"name"`: Name-only matching

#### File Pattern Filtering
- `exclude_patterns`: File patterns to exclude from analysis
- `include_patterns`: File patterns to include in analysis
- Supports glob patterns: `*.go`, `src/**/*.js`, `test_*`

#### Performance Settings
- `max_commits`: Maximum commits to process
- `chunk_size`: Processing chunk size for large repositories
- `parallel_processing`: Enable parallel processing
- `max_workers`: Maximum worker goroutines

## Date Formats

The tool supports both absolute and relative date formats:

**Absolute Dates:**
- `2024-01-15`
- `2024-01-15 14:30:00`
- `2024-01-15T14:30:00Z`
- `01/15/2024`
- `15-01-2024`
- `January 15, 2024`
- `Jan 15, 2024`

**Relative Dates:**
- `today`
- `yesterday`
- `this week`
- `this month`
- `this year`
- `1 day ago`
- `2 weeks ago`
- `3 months ago`
- `1 year ago`
- `a week ago`
- `an hour ago` (not supported - use days/weeks/months/years)

## Advanced Author Matching

Author filtering supports multiple matching modes:

**Basic Matching:**
- Partial name matching: `"john"` matches "John Doe"
- Email matching: `"john@example.com"` matches exact email
- Domain matching: `"@example.com"` matches all emails from domain

**Advanced Matching (via configuration):**
- **Exact Match**: `"John Doe"` matches only exactly "John Doe"
- **Contains Match**: `"john"` matches "John Doe" and "john@example.com"
- **Email-Only Match**: `"example.com"` matches only email addresses
- **Name-Only Match**: `"John"` matches only names, ignores emails
- **Regex Match**: `"^John.*"` for pattern-based matching
- **Case Sensitivity**: Configurable case-sensitive or insensitive matching

## Output Examples

### Contribution Graph
```bash
$ git-stats -contrib
Git Contribution Graph
======================
Repository: my-project
Total Commits: 247

Contributions in the last year:
    Jan  Feb  Mar  Apr  May  Jun  Jul  Aug  Sep  Oct  Nov  Dec
Mon ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜
Tue ⬜🟩🟩⬜⬜⬜⬜⬜⬜⬜⬜⬜
Wed ⬜⬜🟩🟩🟩⬜⬜⬜⬜⬜⬜⬜
...

Current streak: 5 days
Longest streak: 23 days
Total contributions: 247
```

### Repository Summary
```bash
$ git-stats -summary
Git Repository Summary
======================
Repository: my-project
Path: /path/to/my-project
Total Commits: 247
First Commit: 2024-01-15 09:30:00
Last Commit: 2024-08-10 16:45:00
Branches: 3

Active Contributors: 5
Most Active: john@example.com (89 commits)
```

### Help System
```bash
$ git-stats -help
Git Stats - Enhanced Git Repository Analysis Tool

Usage: git-stats [options] [repository-path]

Commands:
  -contrib         Show git contribution graph (GitHub-style) [default]
  -summary         Show detailed repository statistics
  -contributors    Show contributor statistics
  -health          Show repository health metrics
  -gui             Launch interactive ncurses GUI

[... detailed help with examples ...]
```

### Error Handling with Suggestions
```bash
$ git-stats -since "invalid-date"
Error: invalid since date 'invalid-date': unable to parse date...

Suggestion: Use a valid date format. Supported formats:
  - Absolute: 2024-01-15, 2024-01-15 14:30:00, 01/15/2024
  - Relative: today, yesterday, 1 week ago, 2 months ago
Example: git-stats -since "2024-01-01" -until "2024-12-31"

For more help, run: git-stats -help
```

## New Features in Latest Version

### Advanced Filtering System
- **Multiple Filter Types**: Date range, author, file path, merge commit, and limit filters
- **Filter Chaining**: Combine multiple filters with AND logic
- **Advanced Author Matching**: Exact, contains, regex, email-only, and name-only matching
- **File Pattern Filtering**: Include/exclude files using glob patterns
- **Performance Optimized**: Efficient filter implementations with benchmarking

### Configuration Management
- **JSON-Based Configuration**: Hierarchical configuration with automatic defaults
- **User and Workspace Configs**: Support for both user-level and workspace-level settings
- **Import/Export**: Configuration portability and backup
- **Validation**: Comprehensive validation with helpful error messages
- **CLI Integration**: Seamless integration between CLI options and configuration

### Enhanced Date Parsing
- **Extended Relative Dates**: Support for "this week", "this month", "this year"
- **Multiple Absolute Formats**: ISO dates, US dates, European dates, natural language
- **Timezone Support**: Full timezone and ISO format support

### Developer Experience
- **Comprehensive Testing**: 100% test coverage for new features
- **Benchmarking**: Performance testing for all filter operations
- **Integration Testing**: End-to-end testing of complete workflows
- **Documentation**: Extensive documentation with examples

### Performance Improvements
- **Configurable Limits**: Adjustable processing limits for large repositories
- **Parallel Processing**: Configurable parallel processing with worker pools
- **Memory Optimization**: Efficient memory usage for large datasets
- **Caching Support**: Optional caching for repeated operations
