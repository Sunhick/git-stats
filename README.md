# git-stats

Enhanced terminal git statistics utility with comprehensive analysis and visualization capabilities.

## Features

- **Contribution Analysis**: GitHub-style contribution graphs with activity levels and streak tracking
- **Statistical Analysis**: Comprehensive commit frequency, file statistics, and temporal pattern analysis
- **Health Metrics**: Repository health scoring with activity trends and growth analysis
- **Multiple Output Formats**: Terminal, JSON, CSV, and interactive visualizations
- **Advanced Filtering**: Time range, author, and merge commit filtering
- **Performance Optimized**: Efficient algorithms for large repositories

## Architecture

The application is built with a modular architecture:

- **Models**: Core data structures for commits, contributors, and statistics
- **Git Package**: Repository operations and commit parsing
- **Analyzers**: Statistical analysis engines
  - Contribution Analyzer: Activity graphs and streak calculations
  - Statistics Analyzer: Commit patterns and file statistics
  - Health Analyzer: Repository health metrics and trends
- **CLI**: Command-line interface and output formatting
- **Visualization**: Terminal-based charts and graphs

## Build & Install

```shell
# Build the application
$ make build

# Install to /usr/local/bin
$ make install

# Or build and run directly
$ make run
```

## Development

```shell
# Run tests
$ make test

# Run specific test suites
$ make test-analyzers
$ make test-models
$ make test-git

# Run with coverage
$ make test-coverage

# Format and check code
$ make check

# Development mode (auto-rebuild)
$ make dev
```

## Usage

After building, you can add the folder containing `git-stats` to your PATH variable. Then `git stats` will invoke the utility as a git subcommand.

### Basic Usage

```shell
# Show repository summary
$ git stats

# Show contribution graph
$ git stats --contrib

# Show detailed statistics
$ git stats --detailed

# Filter by author
$ git stats --author="John Doe"

# Filter by time range
$ git stats --since="2024-01-01" --until="2024-12-31"

# Export to JSON
$ git stats --format=json --output=stats.json
```

### Advanced Analysis

```shell
# Repository health analysis
$ git stats --health

# File statistics
$ git stats --files

# Time-based patterns
$ git stats --patterns

# Exclude merge commits
$ git stats --no-merges

# Interactive mode
$ git stats --interactive
```

## Command Line Options

| Flag              | Description                         |
| ----------------- | ----------------------------------- |
| `--summary`       | Show repository summary (default)   |
| `--contrib`       | Display contribution graph          |
| `--detailed`      | Show detailed statistics            |
| `--health`        | Repository health analysis          |
| `--files`         | File and file type statistics       |
| `--patterns`      | Time-based commit patterns          |
| `--author=NAME`   | Filter by author name or email      |
| `--since=DATE`    | Start date for analysis             |
| `--until=DATE`    | End date for analysis               |
| `--no-merges`     | Exclude merge commits               |
| `--format=FORMAT` | Output format (terminal, json, csv) |
| `--output=FILE`   | Output file path                    |
| `--interactive`   | Interactive mode                    |

## Output Examples

### Contribution Graph
```
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

### Repository Health
```
Repository Health Score: 85/100

✅ High commit frequency (1.2 commits/day)
✅ Good contributor diversity (5 active contributors)
✅ Activity trending upward
⚠️  Repository age: 6 months (still maturing)
```

### File Statistics
```
Top Files by Commits:
1. src/main.go          (45 commits, 1,234 lines)
2. README.md            (23 commits, 456 lines)
3. src/analyzer.go      (19 commits, 789 lines)

Top File Types:
1. .go     (156 commits, 12,345 lines, 23 files)
2. .md     (34 commits, 2,456 lines, 8 files)
3. .json   (12 commits, 567 lines, 4 files)
```
