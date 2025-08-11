// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - CLI package for command line parsing

package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents the configuration for the git-stats tool
type Config struct {
	Command      string     // contrib, summary, contributors, health, gui
	Since        *time.Time // --since flag
	Until        *time.Time // --until flag
	Author       string     // --author flag
	Format       string     // json, csv, terminal
	OutputFile   string     // --output flag
	RepoPath     string     // repository path
	ShowProgress bool       // --progress flag
	Limit        int        // --limit flag for large repos
	GUIMode      bool       // --gui flag for ncurses interface
	ShowHelp     bool       // --help flag
	NoColor      bool       // --no-color flag to disable colors
	ColorTheme   string     // --theme flag for color theme (github, blue, fire)
}

// Parser interface for command line parsing
type Parser interface {
	Parse(args []string) (*Config, error)
	PrintUsage()
	PrintHelp()
	PrintErrorWithSuggestion(err error)
}

// CLIParser implements the Parser interface
type CLIParser struct {
	validator Validator
}

// NewCLIParser creates a new CLI parser with the given validator
func NewCLIParser(validator Validator) *CLIParser {
	return &CLIParser{
		validator: validator,
	}
}

// Parse parses command line arguments and returns a Config
func (p *CLIParser) Parse(args []string) (*Config, error) {
	config := &Config{
		Format:   "terminal", // default format
		RepoPath: ".",        // default to current directory
		Limit:    10000,      // default limit
	}

	// Create a new flag set to avoid conflicts with global flags
	fs := flag.NewFlagSet("git-stats", flag.ContinueOnError)
	fs.Usage = func() {
		p.PrintUsage()
	}

	// Define flags
	var (
		contrib      = fs.Bool("contrib", false, "Show git contribution graph (GitHub-style)")
		summary      = fs.Bool("summary", false, "Show detailed repository statistics")
		contributors = fs.Bool("contributors", false, "Show contributor statistics")
		health       = fs.Bool("health", false, "Show repository health metrics")
		gui          = fs.Bool("gui", false, "Launch interactive ncurses GUI")
		since        = fs.String("since", "", "Show commits since date (YYYY-MM-DD or relative like '1 week ago')")
		until        = fs.String("until", "", "Show commits until date (YYYY-MM-DD or relative like '1 week ago')")
		author       = fs.String("author", "", "Filter commits by author (name or email, supports partial matching)")
		format       = fs.String("format", "terminal", "Output format: terminal, json, csv")
		output       = fs.String("output", "", "Output file path (default: stdout)")
		progress     = fs.Bool("progress", false, "Show progress indicators for long operations")
		limit        = fs.Int("limit", 10000, "Limit number of commits to process (for large repositories)")
		help         = fs.Bool("help", false, "Show help information")
		h            = fs.Bool("h", false, "Show help information (short form)")
		noColor      = fs.Bool("no-color", false, "Disable colored output")
		colorTheme   = fs.String("theme", "github", "Color theme for contribution graph: github, blue, fire")
	)

	// Parse arguments
	err := fs.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Set help flag
	config.ShowHelp = *help || *h

	// If help is requested, return early
	if config.ShowHelp {
		return config, nil
	}

	// Determine command based on flags
	commandCount := 0
	if *contrib {
		config.Command = "contrib"
		commandCount++
	}
	if *summary {
		config.Command = "summary"
		commandCount++
	}
	if *contributors {
		config.Command = "contributors"
		commandCount++
	}
	if *health {
		config.Command = "health"
		commandCount++
	}
	if *gui {
		config.GUIMode = true
		if config.Command == "" {
			config.Command = "contrib" // default command for GUI mode
		}
		commandCount++
	}

	// If no command specified, default to contrib
	if commandCount == 0 {
		config.Command = "contrib"
	} else if commandCount > 1 && !*gui {
		return nil, fmt.Errorf("only one command can be specified at a time")
	}

	// Parse date flags
	if *since != "" {
		sinceTime, err := parseDate(*since)
		if err != nil {
			return nil, fmt.Errorf("invalid since date '%s': %w", *since, err)
		}
		config.Since = sinceTime
	}

	if *until != "" {
		untilTime, err := parseDate(*until)
		if err != nil {
			return nil, fmt.Errorf("invalid until date '%s': %w", *until, err)
		}
		config.Until = untilTime
	}

	// Set other configuration values
	config.Author = strings.TrimSpace(*author)
	config.Format = strings.ToLower(strings.TrimSpace(*format))
	config.OutputFile = strings.TrimSpace(*output)
	config.ShowProgress = *progress
	config.Limit = *limit
	config.NoColor = *noColor
	config.ColorTheme = strings.ToLower(strings.TrimSpace(*colorTheme))

	// Get repository path from remaining arguments or use current directory
	remainingArgs := fs.Args()
	if len(remainingArgs) > 0 {
		config.RepoPath = remainingArgs[0]
	}

	// Validate configuration
	if p.validator != nil {
		if err := p.validator.ValidateConfig(config); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}

	return config, nil
}

// parseDate parses various date formats
func parseDate(dateStr string) (*time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)

	// Try different date formats
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"01/02/2006",
		"02-01-2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t, nil
		}
	}

	// Try relative dates
	if relativeTime, err := parseRelativeDate(dateStr); err == nil {
		return relativeTime, nil
	}

	return nil, fmt.Errorf("unable to parse date '%s'. Supported formats: YYYY-MM-DD, YYYY-MM-DD HH:MM:SS, relative dates like '1 week ago', '2 months ago'", dateStr)
}

// parseRelativeDate parses relative date strings like "1 week ago", "2 months ago"
func parseRelativeDate(dateStr string) (*time.Time, error) {
	dateStr = strings.ToLower(strings.TrimSpace(dateStr))
	now := time.Now()

	// Handle "today", "yesterday"
	switch dateStr {
	case "today":
		return &now, nil
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		return &yesterday, nil
	}

	// Parse patterns like "1 week ago", "2 months ago"
	if strings.HasSuffix(dateStr, " ago") {
		dateStr = strings.TrimSuffix(dateStr, " ago")
		parts := strings.Fields(dateStr)

		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid relative date format")
		}

		amount, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid number in relative date: %s", parts[0])
		}

		unit := parts[1]
		if strings.HasSuffix(unit, "s") {
			unit = strings.TrimSuffix(unit, "s") // Remove plural 's'
		}

		var result time.Time
		switch unit {
		case "day":
			result = now.AddDate(0, 0, -amount)
		case "week":
			result = now.AddDate(0, 0, -amount*7)
		case "month":
			result = now.AddDate(0, -amount, 0)
		case "year":
			result = now.AddDate(-amount, 0, 0)
		default:
			return nil, fmt.Errorf("unsupported time unit: %s", unit)
		}

		return &result, nil
	}

	return nil, fmt.Errorf("unsupported relative date format")
}

// PrintUsage prints the usage information
func (p *CLIParser) PrintUsage() {
	fmt.Fprintf(os.Stderr, "Git Stats - Enhanced Git Repository Analysis Tool\n\n")
	fmt.Fprintf(os.Stderr, "Usage: git-stats [options] [repository-path]\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  -contrib         Show git contribution graph (GitHub-style) [default]\n")
	fmt.Fprintf(os.Stderr, "  -summary         Show detailed repository statistics\n")
	fmt.Fprintf(os.Stderr, "  -contributors    Show contributor statistics\n")
	fmt.Fprintf(os.Stderr, "  -health          Show repository health metrics\n")
	fmt.Fprintf(os.Stderr, "  -gui             Launch interactive ncurses GUI\n\n")
	fmt.Fprintf(os.Stderr, "Filtering Options:\n")
	fmt.Fprintf(os.Stderr, "  -since <date>    Show commits since date (YYYY-MM-DD or relative)\n")
	fmt.Fprintf(os.Stderr, "  -until <date>    Show commits until date (YYYY-MM-DD or relative)\n")
	fmt.Fprintf(os.Stderr, "  -author <name>   Filter commits by author (supports partial matching)\n\n")
	fmt.Fprintf(os.Stderr, "Output Options:\n")
	fmt.Fprintf(os.Stderr, "  -format <fmt>    Output format: terminal, json, csv [default: terminal]\n")
	fmt.Fprintf(os.Stderr, "  -output <file>   Output file path [default: stdout]\n")
	fmt.Fprintf(os.Stderr, "  -progress        Show progress indicators for long operations\n\n")
	fmt.Fprintf(os.Stderr, "Performance Options:\n")
	fmt.Fprintf(os.Stderr, "  -limit <n>       Limit number of commits to process [default: 10000]\n\n")
	fmt.Fprintf(os.Stderr, "Other Options:\n")
	fmt.Fprintf(os.Stderr, "  -help, -h        Show this help information\n\n")
}

// PrintHelp prints detailed help information with examples
func (p *CLIParser) PrintHelp() {
	p.PrintUsage()
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  Basic Usage:\n")
	fmt.Fprintf(os.Stderr, "    git-stats                                    # Show contribution graph for current repo\n")
	fmt.Fprintf(os.Stderr, "    git-stats -summary                           # Show detailed statistics\n")
	fmt.Fprintf(os.Stderr, "    git-stats -contributors                      # Show contributor statistics\n")
	fmt.Fprintf(os.Stderr, "    git-stats -health                            # Show repository health metrics\n")
	fmt.Fprintf(os.Stderr, "    git-stats -gui                               # Launch interactive GUI\n\n")
	fmt.Fprintf(os.Stderr, "  Date Filtering:\n")
	fmt.Fprintf(os.Stderr, "    git-stats -contrib -since \"2024-01-01\"       # Show contributions since Jan 1, 2024\n")
	fmt.Fprintf(os.Stderr, "    git-stats -summary -since \"1 month ago\"       # Show stats for last month\n")
	fmt.Fprintf(os.Stderr, "    git-stats -health -since \"yesterday\" -until \"today\"  # Show health for yesterday\n\n")
	fmt.Fprintf(os.Stderr, "  Author Filtering:\n")
	fmt.Fprintf(os.Stderr, "    git-stats -contributors -author \"john\"        # Show stats for authors matching 'john'\n")
	fmt.Fprintf(os.Stderr, "    git-stats -contrib -author \"john@example.com\" # Filter by email\n\n")
	fmt.Fprintf(os.Stderr, "  Output Formats:\n")
	fmt.Fprintf(os.Stderr, "    git-stats -summary -format json              # Output as JSON\n")
	fmt.Fprintf(os.Stderr, "    git-stats -contributors -format csv          # Output as CSV\n")
	fmt.Fprintf(os.Stderr, "    git-stats -summary -format json -output report.json  # Save to file\n\n")
	fmt.Fprintf(os.Stderr, "  Advanced Options:\n")
	fmt.Fprintf(os.Stderr, "    git-stats -summary -progress -limit 5000     # Show progress, limit commits\n")
	fmt.Fprintf(os.Stderr, "    git-stats -contrib /path/to/repo              # Analyze specific repository\n")
	fmt.Fprintf(os.Stderr, "    git-stats -gui -since \"1 year ago\"            # GUI with date filter\n\n")
	fmt.Fprintf(os.Stderr, "Date Formats:\n")
	fmt.Fprintf(os.Stderr, "  Absolute: 2024-01-15, 2024-01-15 14:30:00, 01/15/2024, 15-01-2024\n")
	fmt.Fprintf(os.Stderr, "  Relative: today, yesterday, 1 week ago, 2 months ago, 1 year ago\n\n")
	fmt.Fprintf(os.Stderr, "Author Matching:\n")
	fmt.Fprintf(os.Stderr, "  Supports partial name matching and email matching\n")
	fmt.Fprintf(os.Stderr, "  Examples: \"john\", \"john@example.com\", \"John Doe\"\n\n")
	fmt.Fprintf(os.Stderr, "GUI Mode:\n")
	fmt.Fprintf(os.Stderr, "  The --gui flag launches an interactive ncurses interface\n")
	fmt.Fprintf(os.Stderr, "  Use arrow keys to navigate, 'q' to quit, '?' for help\n")
	fmt.Fprintf(os.Stderr, "  Switch views with 'c' (contrib), 's' (stats), 't' (team), 'h' (health)\n\n")
	fmt.Fprintf(os.Stderr, "Performance:\n")
	fmt.Fprintf(os.Stderr, "  Use --limit to process fewer commits for large repositories\n")
	fmt.Fprintf(os.Stderr, "  Use --progress to see progress indicators for long operations\n\n")
	fmt.Fprintf(os.Stderr, "Note: Make sure to run this command from within a git repository or specify the repository path.\n")
}

// PrintErrorWithSuggestion prints an error message with helpful suggestions
func (p *CLIParser) PrintErrorWithSuggestion(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)

	errorMsg := err.Error()

	// Provide contextual suggestions based on error type
	if strings.Contains(errorMsg, "not a git repository") {
		fmt.Fprintf(os.Stderr, "Suggestion: Make sure you're in a git repository or specify a valid repository path.\n")
		fmt.Fprintf(os.Stderr, "Example: git-stats /path/to/your/git/repo\n\n")
	} else if strings.Contains(errorMsg, "invalid since date") || strings.Contains(errorMsg, "invalid until date") {
		fmt.Fprintf(os.Stderr, "Suggestion: Use a valid date format. Supported formats:\n")
		fmt.Fprintf(os.Stderr, "  - Absolute: 2024-01-15, 2024-01-15 14:30:00, 01/15/2024\n")
		fmt.Fprintf(os.Stderr, "  - Relative: today, yesterday, 1 week ago, 2 months ago\n")
		fmt.Fprintf(os.Stderr, "Example: git-stats -since \"2024-01-01\" -until \"2024-12-31\"\n\n")
	} else if strings.Contains(errorMsg, "invalid format") {
		fmt.Fprintf(os.Stderr, "Suggestion: Use one of the supported output formats: terminal, json, csv\n")
		fmt.Fprintf(os.Stderr, "Example: git-stats -format json\n\n")
	} else if strings.Contains(errorMsg, "only one command can be specified") {
		fmt.Fprintf(os.Stderr, "Suggestion: Choose only one command at a time:\n")
		fmt.Fprintf(os.Stderr, "  -contrib, -summary, -contributors, or -health\n")
		fmt.Fprintf(os.Stderr, "Example: git-stats -summary (not git-stats -summary -contrib)\n\n")
	} else if strings.Contains(errorMsg, "limit must be greater than 0") {
		fmt.Fprintf(os.Stderr, "Suggestion: Use a positive number for the limit option.\n")
		fmt.Fprintf(os.Stderr, "Example: git-stats -limit 5000\n\n")
	} else if strings.Contains(errorMsg, "since date") && strings.Contains(errorMsg, "cannot be after until date") {
		fmt.Fprintf(os.Stderr, "Suggestion: Make sure the 'since' date is before the 'until' date.\n")
		fmt.Fprintf(os.Stderr, "Example: git-stats -since \"2024-01-01\" -until \"2024-12-31\"\n\n")
	} else if strings.Contains(errorMsg, "repository path does not exist") {
		fmt.Fprintf(os.Stderr, "Suggestion: Check that the repository path exists and is accessible.\n")
		fmt.Fprintf(os.Stderr, "Example: git-stats /path/to/existing/repo\n\n")
	}

	fmt.Fprintf(os.Stderr, "For more help, run: git-stats -help\n")
}
