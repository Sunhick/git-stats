# NCurses GUI Interface

This directory contains the NCurses GUI interface implementation for the git-stats tool.

## Build Tags

The GUI functionality uses Go build tags to manage dependencies:

- **Default build** (no tags): Uses stub implementation without external dependencies
- **GUI build** (`-tags gui`): Uses full implementation with tview/tcell dependencies

## Building

### Using Makefile (Recommended)
```bash
# Build with GUI support (installs dependencies automatically)
make build-gui

# Launch GUI mode
make gui

# Launch with specific views
make run-gui-contrib      # Start with contribution graph
make run-gui-summary      # Start with summary view
make run-gui-contributors # Start with contributors view
make run-gui-health       # Start with health metrics

# Development mode
make dev-gui
```

### Manual Build
```bash
# Default Build (Stub)
go build ./visualizers

# GUI Build (Full Implementation)
go build -tags gui ./visualizers

# Install GUI dependencies first
go get github.com/gdamore/tcell/v2@v2.6.0
go get github.com/rivo/tview@v0.0.0-20230826224341-9754ab44dc1c
```

## Testing

### Using Makefile (Recommended)
```bash
# Run all GUI tests
make test-gui-all

# Run specific test suites
make test-gui-unit         # Unit tests for GUI components
make test-gui-integration  # Integration tests for navigation workflows
```

### Manual Testing
```bash
# Unit Tests (Stub Implementation)
go test ./visualizers

# Integration Tests (Full Implementation)
go test -tags gui ./visualizers
```

## Dependencies

The GUI implementation requires:
- `github.com/gdamore/tcell/v2` - Terminal cell library
- `github.com/rivo/tview` - Terminal UI framework

These dependencies are only required when building with the `gui` tag.

## Usage

The GUI interface provides:

1. **GUIState Management**: Handles view switching, date selection, and navigation
2. **ContributionGraphWidget**: Interactive contribution graph display
3. **DetailPanelWidget**: Shows detailed information for selected items
4. **StatusBarWidget**: Displays status and keyboard shortcuts
5. **GUIInterface**: Main interface coordinator

## Key Features

- Interactive contribution graph similar to GitHub
- Keyboard navigation (arrow keys, shortcuts)
- Multiple views (contribution, statistics, contributors, health)
- Help system with keyboard shortcuts
- Real-time status updates

## Keyboard Shortcuts

- `←→`: Navigate days
- `↑↓`: Navigate weeks
- `h/l`: Navigate months
- `c`: Contribution view
- `s`: Statistics view
- `t`: Team/Contributors view
- `H`: Health metrics view
- `?`: Toggle help
- `q/ESC`: Quit

## Architecture

The GUI uses a widget-based architecture:

```
GUIInterface
├── GUIState (state management)
├── ContributionGraphWidget (main graph display)
├── DetailPanelWidget (information panel)
└── StatusBarWidget (status and shortcuts)
```

Each widget is responsible for its own rendering and input handling, coordinated through the shared GUIState.
