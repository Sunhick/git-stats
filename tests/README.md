# Git-Stats Unit Tests

This directory contains comprehensive unit tests for the git-stats project.

## Structure

```
tests/
├── cli/           # Tests for CLI parsing and validation
├── git/           # Tests for git repository operations
├── models/        # Tests for data models and structures
├── utils/         # Tests for utility functions
├── go.mod         # Go module file for tests
└── README.md      # This file
```

## Running Tests

To run all tests:
```bash
cd tests
go test ./...
```

To run tests for a specific package:
```bash
cd tests
go test ./utils
go test ./git
go test ./models
go test ./cli
```

To run tests with verbose output:
```bash
cd tests
go test -v ./...
```

To run tests with coverage:
```bash
cd tests
go test -cover ./...
```

## Test Coverage

The tests cover:

### Utils Package
- Error handling and recovery suggestions
- Date parsing (absolute and relative formats)
- Progress tracking and spinner functionality

### Git Package
- Repository interface implementations
- Commit and contributor data structures
- Mock implementations for testing

### Models Package
- Statistical analysis result structures
- Configuration models
- Time range and health metrics

### CLI Package
- Command line parsing interfaces
- Input validation
- Configuration structures

## Writing New Tests

When adding new functionality to the main project, please add corresponding tests:

1. Create test files with `_test.go` suffix
2. Use the same package structure as the main project
3. Include both positive and negative test cases
4. Use mock implementations for interfaces
5. Test edge cases and error conditions

## Test Conventions

- Test functions should start with `Test`
- Use descriptive test names that explain what is being tested
- Include both success and failure scenarios
- Use table-driven tests for multiple similar test cases
- Mock external dependencies and interfaces
