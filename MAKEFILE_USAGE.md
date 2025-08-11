# Makefile Usage Guide

This document describes the enhanced Makefile for the git-stats project.

## Quick Start

```bash
# Build the application
make build

# Run all tests
make test

# Format code and run checks
make check

# Show all available targets
make help
```

## Build Targets

- `make build` - Build the application
- `make rebuild` - Clean and build
- `make debug` - Build debug version with debugging symbols
- `make dev` - Run in development mode (go run)

## Test Targets

- `make test` - Run all tests with race detection and coverage
- `make test-utils` - Run utility tests only
- `make test-models` - Run model tests only
- `make test-git` - Run git tests only
- `make test-cli` - Run CLI tests only
- `make test-coverage` - Run tests and generate HTML coverage report
- `make test-integration` - Run integration tests

## Development Targets

- `make run` - Build and run the application
- `make run-summary` - Run with summary flag
- `make run-contrib` - Run with contribution graph flag
- `make fmt` - Format Go code
- `make vet` - Run go vet
- `make lint` - Run golint (if available)
- `make check` - Run fmt, vet, and test

## Dependency Management

- `make deps` - Download dependencies
- `make deps-update` - Update and tidy dependencies
- `make deps-vendor` - Vendor dependencies

## Distribution

- `make dist` - Create distribution binaries for multiple platforms
- `make install` - Install to /usr/local/bin
- `make uninstall` - Remove from /usr/local/bin

## Cleanup

- `make clean` - Clean build artifacts
- `make clean-all` - Clean all generated files including caches

## Features

### Enhanced Build Process
- Optimized build flags with size reduction (`-ldflags="-s -w"`)
- Debug builds with debugging symbols
- Cross-platform distribution builds

### Comprehensive Testing
- Race detection enabled by default
- Coverage reporting with HTML output
- Separate test targets for different modules
- Integration test support

### Code Quality
- Automatic code formatting
- Static analysis with go vet
- Optional linting support
- Combined check target for CI/CD

### Development Workflow
- Development mode with hot reloading
- Multiple run configurations
- Dependency management
- Clean separation of concerns

### Version Information
The Makefile includes version information in builds:
- Git commit hash
- Build timestamp
- Version tag (if available)

## Examples

```bash
# Complete development workflow
make clean && make check && make build

# Run specific tests during development
make test-utils

# Generate coverage report
make test-coverage
open tests/coverage.html

# Create distribution
make dist
ls build/

# Install for system-wide use
sudo make install
```

## Configuration

The Makefile uses `common.mk` for shared configuration:
- Build tools and flags
- Color output support
- Version information extraction
- Common file patterns

## Troubleshooting

### Build Issues
- Ensure Go is installed and in PATH
- Run `make deps` to download dependencies
- Check `make vet` for static analysis issues

### Test Issues
- Use `make test-utils` to run specific test suites
- Check for race conditions with the `-race` flag
- Review coverage reports for missing test coverage

### Import Issues
- Ensure all imports use the module path `git-stats/...`
- Avoid relative imports like `../src/...`
- Run `make fmt` to fix formatting issues
