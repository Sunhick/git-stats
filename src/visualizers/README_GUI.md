# NCurses GUI Interface

This directory contains the NCurses GUI interface implementation for the git-stats tool.

## Build Tags

The GUI functionality uses Go build tags to manage dependencies:

- **Default build** (no tags): Uses stub implementation without external dependencies
- **GUI build** (`-tags gui`): Uses full implementation with tview/tcell dependencies

## Building

### Default Build (Stub)
```bash
go build ./visualizers
```

### GUI Build (Full Implementation)
```bash
go build -tags gui ./visualizers
```

## Testing

### Unit Tests (Stub Implementation)
```bash
go test ./visualizers
```

### Integration Tests (Full Implementation)
```bash
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
