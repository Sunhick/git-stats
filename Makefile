# Copyright (c) 2019 Sunil
# Enhanced git-stats tool - Makefile

TARGET = git-stats
SRC_DIR = src
TEST_DIR = tests
BUILD_DIR = build
MAIN_FILE = $(SRC_DIR)/git-stats.go

# Go build flags
GO_BUILD_FLAGS = -ldflags="-s -w"
GO_TEST_FLAGS = -v -race -coverprofile=coverage.out

.PHONY: all
all: ${TARGET}

include common.mk

# Build targets
.PHONY: build
build: ${TARGET}

.PHONY: rebuild
rebuild: clean ${TARGET}

.PHONY: debug
debug:
	${E} "Building debug version..."
	${Q} cd $(SRC_DIR) && ${GO} build -gcflags="all=-N -l" -o ../$(TARGET) .

${TARGET}: $(MAIN_FILE)
	${E} "Building $(TARGET)..."
	${Q} cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -o ../$(TARGET) .

# Test targets
.PHONY: test
test:
	${E} "Running tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./...

.PHONY: test-utils
test-utils:
	${E} "Running utility tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./utils/

.PHONY: test-models
test-models:
	${E} "Running model tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./models/

.PHONY: test-git
test-git:
	${E} "Running git tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./git/

.PHONY: test-cli
test-cli:
	${E} "Running CLI tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./cli/

.PHONY: test-coverage
test-coverage: test
	${E} "Generating coverage report..."
	${Q} cd $(TEST_DIR) && ${GO} tool cover -html=coverage.out -o coverage.html
	${E} "Coverage report generated: $(TEST_DIR)/coverage.html"

.PHONY: test-integration
test-integration:
	${E} "Running integration tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./utils/integration_test.go ./utils/date_test.go ./utils/errors_test.go ./utils/progress_test.go

# Development targets
.PHONY: run
run: ${TARGET}
	${E} "Running $(TARGET)..."
	${Q} ./${TARGET}

.PHONY: run-summary
run-summary: ${TARGET}
	${E} "Running $(TARGET) with summary..."
	${Q} ./${TARGET} -summarize

.PHONY: run-contrib
run-contrib: ${TARGET}
	${E} "Running $(TARGET) with contribution graph..."
	${Q} ./${TARGET} -contrib

.PHONY: dev
dev:
	${E} "Starting development mode..."
	${Q} cd $(SRC_DIR) && ${GO} run .

# Code quality targets
.PHONY: fmt
fmt:
	${E} "Formatting Go code..."
	${Q} cd $(SRC_DIR) && ${GO} fmt ./...
	${Q} cd $(TEST_DIR) && ${GO} fmt ./...

.PHONY: vet
vet:
	${E} "Running go vet..."
	${Q} cd $(SRC_DIR) && ${GO} vet ./...
	${Q} cd $(TEST_DIR) && ${GO} vet ./...

.PHONY: lint
lint:
	${E} "Running golint (if available)..."
	${Q} command -v golint >/dev/null 2>&1 && (cd $(SRC_DIR) && golint ./... && cd ../$(TEST_DIR) && golint ./...) || echo "golint not installed, skipping..."

.PHONY: check
check: fmt vet test
	${E} "All checks passed!"

# Dependency management
.PHONY: deps
deps:
	${E} "Downloading dependencies..."
	${Q} cd $(SRC_DIR) && ${GO} mod download
	${Q} cd $(TEST_DIR) && ${GO} mod download

.PHONY: deps-update
deps-update:
	${E} "Updating dependencies..."
	${Q} cd $(SRC_DIR) && ${GO} mod tidy
	${Q} cd $(TEST_DIR) && ${GO} mod tidy

.PHONY: deps-vendor
deps-vendor:
	${E} "Vendoring dependencies..."
	${Q} cd $(SRC_DIR) && ${GO} mod vendor

# Installation targets
.PHONY: install
install: ${TARGET}
	${E} "Installing $(TARGET)..."
	${Q} cp $(TARGET) /usr/local/bin/

.PHONY: uninstall
uninstall:
	${E} "Uninstalling $(TARGET)..."
	${Q} rm -f /usr/local/bin/$(TARGET)

# Distribution targets
.PHONY: dist
dist: clean
	${E} "Creating distribution..."
	${Q} mkdir -p $(BUILD_DIR)
	${Q} GOOS=linux GOARCH=amd64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -o ../$(BUILD_DIR)/$(TARGET)-linux-amd64 .
	${Q} GOOS=darwin GOARCH=amd64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -o ../$(BUILD_DIR)/$(TARGET)-darwin-amd64 .
	${Q} GOOS=darwin GOARCH=arm64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -o ../$(BUILD_DIR)/$(TARGET)-darwin-arm64 .
	${Q} GOOS=windows GOARCH=amd64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -o ../$(BUILD_DIR)/$(TARGET)-windows-amd64.exe .

# Clean targets
.PHONY: clean
clean: decruft
	${E} "Cleaning build artifacts..."
	${Q} rm -f ${TARGET}
	${Q} rm -rf $(BUILD_DIR)
	${Q} rm -f $(TEST_DIR)/coverage.out $(TEST_DIR)/coverage.html

.PHONY: clean-all
clean-all: clean
	${E} "Cleaning all generated files..."
	${Q} cd $(SRC_DIR) && ${GO} clean -cache -modcache -testcache

# Help target
.PHONY: help
help:
	${E} "Available targets:"
	${E} "  build          - Build the application"
	${E} "  rebuild        - Clean and build"
	${E} "  debug          - Build debug version"
	${E} "  test           - Run all tests"
	${E} "  test-utils     - Run utility tests"
	${E} "  test-models    - Run model tests"
	${E} "  test-git       - Run git tests"
	${E} "  test-cli       - Run CLI tests"
	${E} "  test-coverage  - Run tests with coverage report"
	${E} "  test-integration - Run integration tests"
	${E} "  run            - Build and run the application"
	${E} "  run-summary    - Run with summary flag"
	${E} "  run-contrib    - Run with contribution graph flag"
	${E} "  dev            - Run in development mode"
	${E} "  fmt            - Format Go code"
	${E} "  vet            - Run go vet"
	${E} "  lint           - Run golint (if available)"
	${E} "  check          - Run fmt, vet, and test"
	${E} "  deps           - Download dependencies"
	${E} "  deps-update    - Update dependencies"
	${E} "  deps-vendor    - Vendor dependencies"
	${E} "  install        - Install to /usr/local/bin"
	${E} "  uninstall      - Remove from /usr/local/bin"
	${E} "  dist           - Create distribution binaries"
	${E} "  clean          - Clean build artifacts"
	${E} "  clean-all      - Clean all generated files"
	${E} "  help           - Show this help message"
