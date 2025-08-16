# git-stats

Enhanced terminal git statistics utility with comprehensive analysis and visualization capabilities.

## Features

- **Contribution Analysis**: GitHub-style contribution graphs with activity levels and streak tracking
- **Statistical Analysis**: Comprehensive commit frequency, file statistics, and temporal pattern analysis
- **Health Metrics**: Repository health scoring with activity trends and growth analysis
- **Multiple Output Formats**: Terminal, JSON, CSV, and interactive visualizations
- **Advanced Filtering**: Time range, author, and merge commit filtering
- **Interactive GUI**: Full-screen ncurses interface with keyboard navigation
- **Comprehensive CLI**: Robust command-line interface with intelligent error handling
- **Flexible Date Parsing**: Support for absolute and relative date formats
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

### Date and Author Filtering

```shell
# Filter by date range
$ git-stats -contrib -since "2024-01-01" -until "2024-12-31"

# Use relative dates
$ git-stats -summary -since "1 month ago"
$ git-stats -health -since "yesterday" -until "today"

# Filter by author
$ git-stats -contributors -author "john"
$ git-stats -contrib -author "john@example.com"
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
| Flag             | Description                    |
| ---------------- | ------------------------------ |
| `-since <date>`  | Show commits since date        |
| `-until <date>`  | Show commits until date        |
| `-author <name>` | Filter by author name or email |

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

### Other Options
| Flag        | Description           |
| ----------- | --------------------- |
| `-help, -h` | Show help information |

## Date Formats

The tool supports both absolute and relative date formats:

**Absolute Dates:**
- `2024-01-15`
- `2024-01-15 14:30:00`
- `01/15/2024`
- `15-01-2024`

**Relative Dates:**
- `today`
- `yesterday`
- `1 week ago`
- `2 months ago`
- `1 year ago`

## Author Matching

Author filtering supports:
- Partial name matching: `"john"`
- Email matching: `"john@example.com"`
- Full name matching: `"John Doe"`

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
