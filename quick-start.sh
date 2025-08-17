#!/bin/bash
# Git Stats Quick Start Script
# Copyright (c) 2019 Sunil

set -e

echo "Git Stats - Quick Start Setup"
echo "============================="
echo

# Check if we're in the right directory
if [ ! -f "Makefile" ] || [ ! -d "src" ]; then
    echo "❌ Error: Please run this script from the git-stats project root directory"
    echo "   (The directory containing Makefile and src/ folder)"
    exit 1
fi

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed or not in PATH"
    echo "   Please install Go 1.19+ from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
echo "✓ Go version: $GO_VERSION"

# Check Git installation
if ! command -v git &> /dev/null; then
    echo "❌ Error: Git is not installed or not in PATH"
    echo "   Please install Git from https://git-scm.com/"
    exit 1
fi

GIT_VERSION=$(git --version | cut -d' ' -f3)
echo "✓ Git version: $GIT_VERSION"
echo

# Build the application
echo "Building git-stats..."
if make build; then
    echo "✓ Build successful!"
else
    echo "❌ Build failed. Please check the error messages above."
    exit 1
fi

echo

# Check if we're in a git repository for testing
if [ -d ".git" ]; then
    echo "Testing git-stats in current repository..."
    echo
    echo "Running: ./git-stats -contrib"
    ./git-stats -contrib | head -10
    echo "..."
    echo
    echo "✓ git-stats is working!"
    echo
    echo "Try these commands:"
    echo "  ./git-stats -help                    # Show help"
    echo "  ./git-stats -summary                 # Repository summary"
    echo "  ./git-stats -contributors            # Contributor analysis"
    echo "  ./git-stats -health                  # Repository health"
    echo
else
    echo "Current directory is not a git repository."
    echo "✓ git-stats built successfully!"
    echo
    echo "To test git-stats, run it in a git repository:"
    echo "  ./git-stats -contrib /path/to/git/repo"
    echo "  ./git-stats -summary /path/to/git/repo"
    echo "  ./git-stats -help"
    echo
fi

# GUI setup information
echo "GUI Mode Setup:"
echo "==============="
echo "To enable GUI mode, you need additional dependencies."
echo
echo "Check GUI dependencies:"
echo "  make check-gui-deps"
echo
echo "Install GUI dependencies (requires network):"
echo "  make deps-gui"
echo "  make build-gui"
echo
echo "Or build GUI in offline mode:"
echo "  make build-gui-offline"
echo
echo "Then run GUI mode:"
echo "  ./git-stats-gui -gui /path/to/repo"
echo

echo "Setup complete! 🎉"
echo
echo "For more information:"
echo "  ./git-stats -help        # Application help"
echo "  make help                # Build system help"
echo "  cat README.md            # Full documentation"
