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
all: build

include common.mk

# Build targets
.PHONY: build
build: build-gui

.PHONY: build-terminal
build-terminal: ${TARGET}

.PHONY: rebuild
rebuild: clean build

.PHONY: debug
debug:
	${E} "Building debug version..."
	${Q} cd $(SRC_DIR) && ${GO} build -gcflags="all=-N -l" -o ../$(TARGET) .

.PHONY: build-gui
build-gui: deps-gui ${TARGET}-gui

.PHONY: build-gui-offline
build-gui-offline: deps-gui-offline ${TARGET}-gui-offline

.PHONY: rebuild-gui
rebuild-gui: clean deps-gui ${TARGET}-gui

.PHONY: rebuild-gui-offline
rebuild-gui-offline: clean deps-gui-offline ${TARGET}-gui-offline

${TARGET}: $(MAIN_FILE)
	${E} "Building $(TARGET)..."
	${Q} cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -o ../$(TARGET) .

${TARGET}-gui: $(MAIN_FILE)
	@echo "Building $(TARGET) with GUI support..."
	@cd $(SRC_DIR) && (${GO} build $(GO_BUILD_FLAGS) -tags gui -o ../$(TARGET)-gui . 2>/dev/null && echo "✓ GUI build successful") || (echo "⚠ GUI dependencies missing. Building with stub implementation..."; ${GO} build $(GO_BUILD_FLAGS) -o ../$(TARGET)-gui .)

${TARGET}-gui-offline: $(MAIN_FILE)
	${E} "Building $(TARGET) with GUI support (offline mode)..."
	${Q} cd $(SRC_DIR) && (${GO} build $(GO_BUILD_FLAGS) -tags gui -o ../$(TARGET)-gui . 2>/dev/null && echo "GUI build successful") || (echo "GUI dependencies not available. Building with stub implementation..."; ${GO} build $(GO_BUILD_FLAGS) -o ../$(TARGET)-gui .)

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

.PHONY: test-analyzers
test-analyzers:
	${E} "Running analyzer tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./analyzers/

.PHONY: test-filters
test-filters:
	${E} "Running filter tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./filters/

.PHONY: test-config
test-config:
	${E} "Running configuration tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./config/

.PHONY: test-integration
test-integration:
	${E} "Running integration tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./integration/

.PHONY: test-coverage
test-coverage: test
	${E} "Generating coverage report..."
	${Q} cd $(TEST_DIR) && ${GO} tool cover -html=coverage.out -o coverage.html
	${E} "Coverage report generated: $(TEST_DIR)/coverage.html"

.PHONY: test-utils-integration
test-utils-integration:
	${E} "Running utility integration tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./utils/integration_test.go ./utils/date_test.go ./utils/errors_test.go ./utils/progress_test.go

# GUI test targets
.PHONY: test-gui
test-gui:
	${E} "Running GUI tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) -tags gui ./visualizers/

.PHONY: test-gui-unit
test-gui-unit:
	${E} "Running GUI unit tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./visualizers/ncurses_gui_unit_test.go

.PHONY: test-gui-integration
test-gui-integration:
	${E} "Running GUI integration tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./visualizers/gui_navigation_integration_test.go

.PHONY: test-gui-all
test-gui-all:
	${E} "Running all GUI tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./visualizers/ncurses_gui_unit_test.go ./visualizers/gui_navigation_integration_test.go

# Individual analyzer test targets
.PHONY: test-contribution
test-contribution:
	${E} "Running contribution analyzer tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./analyzers/contribution_test.go

.PHONY: test-statistics
test-statistics:
	${E} "Running statistics analyzer tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./analyzers/statistics_test.go

.PHONY: test-health
test-health:
	${E} "Running health analyzer tests..."
	${Q} cd $(TEST_DIR) && ${GO} test $(GO_TEST_FLAGS) ./analyzers/health_test.go

# Benchmark targets
.PHONY: bench
bench:
	${E} "Running benchmarks..."
	${Q} cd $(TEST_DIR) && ${GO} test -bench=. -benchmem ./analyzers/ ./filters/

.PHONY: bench-analyzers
bench-analyzers:
	${E} "Running analyzer benchmarks..."
	${Q} cd $(TEST_DIR) && ${GO} test -bench=Benchmark -benchmem ./analyzers/

.PHONY: bench-filters
bench-filters:
	${E} "Running filter benchmarks..."
	${Q} cd $(TEST_DIR) && ${GO} test -bench=Benchmark -benchmem ./filters/

# Development targets
.PHONY: run
run: build
	@echo "Running $(TARGET)-gui..."
	@echo "Usage: ./${TARGET}-gui [options] [repository-path]"
	@echo "Examples:"
	@echo "  ./${TARGET}-gui -contrib"
	@echo "  ./${TARGET}-gui -summary /path/to/repo"
	@echo "  ./${TARGET}-gui -gui /path/to/repo"
	@echo "  ./${TARGET}-gui -help"

.PHONY: run-summary
run-summary: build
	${E} "Running $(TARGET)-gui with summary..."
	${Q} ./${TARGET}-gui -summary

.PHONY: run-contrib
run-contrib: build
	${E} "Running $(TARGET)-gui with contribution graph..."
	${Q} ./${TARGET}-gui -contrib

.PHONY: run-health
run-health: build
	${E} "Running $(TARGET)-gui with health analysis..."
	${Q} ./${TARGET}-gui -health

.PHONY: run-detailed
run-detailed: build
	${E} "Running $(TARGET)-gui with detailed statistics..."
	${Q} ./${TARGET}-gui -detailed

.PHONY: run-files
run-files: build
	${E} "Running $(TARGET)-gui with file statistics..."
	${Q} ./${TARGET}-gui -files

# Advanced filtering examples
.PHONY: run-filter-date
run-filter-date: build
	${E} "Running $(TARGET)-gui with date filtering..."
	${Q} ./${TARGET}-gui -contrib -since "1 month ago"

.PHONY: run-filter-author
run-filter-author: build
	${E} "Running $(TARGET)-gui with author filtering..."
	${Q} ./${TARGET}-gui -contributors -author "$(shell git config user.name)"

.PHONY: run-filter-combined
run-filter-combined: build
	${E} "Running $(TARGET)-gui with combined filters..."
	${Q} ./${TARGET}-gui -summary -since "3 months ago" -author "$(shell git config user.name)" -format json

.PHONY: run-config-demo
run-config-demo: build
	${E} "Running $(TARGET)-gui configuration demo..."
	${Q} ./${TARGET}-gui --show-config || ./${TARGET}-gui -help

.PHONY: dev
dev:
	${E} "Starting development mode with GUI..."
	${Q} cd $(SRC_DIR) && ${GO} run -tags gui .

# GUI targets
.PHONY: gui
gui: ${TARGET}-gui
	@echo "Launching GUI mode..."
	@echo "Usage: ./${TARGET}-gui -gui [repository-path]"
	@echo "Example: ./${TARGET}-gui -gui /path/to/your/repo"
	@echo "Note: GUI requires building with -tags gui for full functionality"

.PHONY: gui-offline
gui-offline: ${TARGET}-gui-offline
	${E} "Launching GUI mode (offline build)..."
	${Q} ./${TARGET}-gui -gui

.PHONY: run-gui
run-gui: gui

.PHONY: run-gui-contrib
run-gui-contrib: ${TARGET}-gui
	${E} "Launching GUI mode with contribution graph..."
	${Q} ./${TARGET}-gui -gui -contrib

.PHONY: run-gui-summary
run-gui-summary: ${TARGET}-gui
	${E} "Launching GUI mode with summary..."
	${Q} ./${TARGET}-gui -gui -summary

.PHONY: run-gui-contributors
run-gui-contributors: ${TARGET}-gui
	${E} "Launching GUI mode with contributors..."
	${Q} ./${TARGET}-gui -gui -contributors

.PHONY: run-gui-health
run-gui-health: ${TARGET}-gui
	${E} "Launching GUI mode with health analysis..."
	${Q} ./${TARGET}-gui -gui -health

.PHONY: dev-gui
dev-gui:
	${E} "Starting development mode with GUI..."
	${Q} cd $(SRC_DIR) && ${GO} run -tags gui . -gui

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

.PHONY: deps-gui
deps-gui:
	@echo "Installing GUI dependencies..."
	@echo "Attempting to download tcell and tview packages..."
	@cd $(SRC_DIR) && (${GO} get github.com/gdamore/tcell/v2@v2.6.0 && echo "✓ tcell downloaded successfully") || echo "⚠ Failed to download tcell - network issue or proxy required"
	@cd $(SRC_DIR) && (${GO} get github.com/rivo/tview@v0.0.0-20230826224341-9754ab44dc1c && echo "✓ tview downloaded successfully") || echo "⚠ Failed to download tview - network issue or proxy required"
	@cd $(SRC_DIR) && ${GO} mod tidy || true
	@echo "Note: If dependencies failed to download, GUI will use stub implementation"

.PHONY: deps-gui-offline
deps-gui-offline:
	${E} "Checking for existing GUI dependencies..."
	${Q} cd $(SRC_DIR) && ${GO} mod tidy || true
	${E} "Note: GUI dependencies must be manually installed or available in module cache"
	${E} "If dependencies are missing, GUI build will use stub implementation"

.PHONY: check-gui-deps
check-gui-deps:
	@echo "Checking GUI dependencies availability..."
	@cd $(SRC_DIR) && (${GO} list -m github.com/gdamore/tcell/v2 >/dev/null 2>&1 && echo "✓ tcell dependency available") || echo "✗ tcell dependency missing"
	@cd $(SRC_DIR) && (${GO} list -m github.com/rivo/tview >/dev/null 2>&1 && echo "✓ tview dependency available") || echo "✗ tview dependency missing"
	@echo ""
	@echo "If dependencies are missing:"
	@echo "  1. Run 'make deps-gui' to download them (requires network)"
	@echo "  2. Or use 'make build-gui-offline' for stub implementation"

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
install: build
	${E} "Installing $(TARGET) and $(TARGET)-gui..."
	${Q} if [ -f $(TARGET)-gui ]; then \
		cp $(TARGET)-gui /usr/local/bin/; \
		echo "$(TARGET)-gui installed successfully"; \
	fi
	${Q} if [ -f $(TARGET) ]; then \
		cp $(TARGET) /usr/local/bin/; \
		echo "$(TARGET) installed successfully"; \
	fi

.PHONY: install-gui
install-gui: ${TARGET}-gui
	${E} "Installing $(TARGET)-gui..."
	${Q} cp $(TARGET)-gui /usr/local/bin/

.PHONY: install-all
install-all: ${TARGET} ${TARGET}-gui
	${E} "Installing both $(TARGET) and $(TARGET)-gui..."
	${Q} cp $(TARGET) /usr/local/bin/
	${Q} cp $(TARGET)-gui /usr/local/bin/

.PHONY: uninstall
uninstall:
	${E} "Uninstalling $(TARGET)..."
	${Q} rm -f /usr/local/bin/$(TARGET)
	${Q} rm -f /usr/local/bin/$(TARGET)-gui

# Distribution targets
.PHONY: dist
dist: clean
	${E} "Creating distribution..."
	${Q} mkdir -p $(BUILD_DIR)
	${Q} GOOS=linux GOARCH=amd64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -o ../$(BUILD_DIR)/$(TARGET)-linux-amd64 .
	${Q} GOOS=darwin GOARCH=amd64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -o ../$(BUILD_DIR)/$(TARGET)-darwin-amd64 .
	${Q} GOOS=darwin GOARCH=arm64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -o ../$(BUILD_DIR)/$(TARGET)-darwin-arm64 .
	${Q} GOOS=windows GOARCH=amd64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -o ../$(BUILD_DIR)/$(TARGET)-windows-amd64.exe .

.PHONY: dist-gui
dist-gui: clean
	${E} "Creating GUI distribution..."
	${Q} mkdir -p $(BUILD_DIR)
	${Q} GOOS=linux GOARCH=amd64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -tags gui -o ../$(BUILD_DIR)/$(TARGET)-gui-linux-amd64 .
	${Q} GOOS=darwin GOARCH=amd64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -tags gui -o ../$(BUILD_DIR)/$(TARGET)-gui-darwin-amd64 .
	${Q} GOOS=darwin GOARCH=arm64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -tags gui -o ../$(BUILD_DIR)/$(TARGET)-gui-darwin-arm64 .
	${Q} GOOS=windows GOARCH=amd64 cd $(SRC_DIR) && ${GO} build $(GO_BUILD_FLAGS) -tags gui -o ../$(BUILD_DIR)/$(TARGET)-gui-windows-amd64.exe .

.PHONY: dist-all
dist-all: dist dist-gui

# Clean targets
.PHONY: clean
clean: decruft
	${E} "Cleaning build artifacts..."
	${Q} rm -f ${TARGET} ${TARGET}-gui
	${Q} rm -rf $(BUILD_DIR)
	${Q} rm -f $(TEST_DIR)/coverage.out $(TEST_DIR)/coverage.html

.PHONY: clean-all
clean-all: clean
	${E} "Cleaning all generated files..."
	${Q} cd $(SRC_DIR) && ${GO} clean -cache -modcache -testcache

# Help target
.PHONY: help
help:
	@echo "Git Stats - Enhanced Git Repository Analysis Tool"
	@echo "================================================="
	@echo ""
	@echo "Quick Start:"
	@echo "  make build                    # Build the application with GUI support"
	@echo "  ./git-stats-gui -help         # Show application help"
	@echo "  ./git-stats-gui -contrib      # Show contribution graph"
	@echo "  ./git-stats-gui -gui          # Launch GUI mode"
	@echo ""
	@echo "GUI Mode:"
	@echo "  make check-gui-deps           # Check if GUI dependencies are available"
	@echo "  make deps-gui                 # Download GUI dependencies (requires network)"
	@echo "  make build-gui                # Build with GUI support"
	@echo "  make build-gui-offline        # Build GUI (offline mode, uses stub if deps missing)"
	@echo ""
	@echo "Available targets:"
	@echo ""
	@echo "Build Targets:"
	@echo "  build          - Build the application with GUI support (default)"
	@echo "  build-terminal - Build terminal-only version"
	@echo "  build-gui      - Build the application with GUI support"
	@echo "  build-gui-offline - Build GUI with existing dependencies only"
	@echo "  rebuild        - Clean and build with GUI support"
	@echo "  rebuild-gui    - Clean and build with GUI support"
	@echo "  rebuild-gui-offline - Clean and build GUI offline"
	@echo "  debug          - Build debug version"
	${E} ""
	${E} "Test Targets:"
	${E} "  test           - Run all tests"
	${E} "  test-utils     - Run utility tests"
	${E} "  test-models    - Run model tests"
	${E} "  test-git       - Run git tests"
	${E} "  test-cli       - Run CLI tests"
	${E} "  test-analyzers - Run analyzer tests"
	${E} "  test-filters   - Run filter tests"
	${E} "  test-config    - Run configuration tests"
	${E} "  test-integration - Run integration tests"
	${E} "  test-contribution - Run contribution analyzer tests"
	${E} "  test-statistics - Run statistics analyzer tests"
	${E} "  test-health    - Run health analyzer tests"
	${E} "  test-coverage  - Run tests with coverage report"
	${E} "  test-utils-integration - Run utility integration tests"
	${E} "  test-gui       - Run GUI tests (requires GUI build tags)"
	${E} "  test-gui-unit  - Run GUI unit tests"
	${E} "  test-gui-integration - Run GUI integration tests"
	${E} "  test-gui-all   - Run all GUI tests"
	${E} ""
	${E} "Run Targets:"
	${E} "  run            - Build and run the application"
	${E} "  run-summary    - Run with summary flag"
	${E} "  run-contrib    - Run with contribution graph flag"
	${E} "  run-health     - Run with health analysis flag"
	${E} "  run-detailed   - Run with detailed statistics flag"
	${E} "  run-files      - Run with file statistics flag"
	${E} "  run-filter-date - Run with date filtering example"
	${E} "  run-filter-author - Run with author filtering example"
	${E} "  run-filter-combined - Run with combined filters example"
	${E} "  run-config-demo - Run configuration demo"
	${E} "  dev            - Run in development mode"
	${E} ""
	${E} "GUI Targets:"
	${E} "  gui            - Build and launch GUI mode"
	${E} "  gui-offline    - Build and launch GUI mode (offline)"
	${E} "  run-gui        - Launch GUI mode (alias for gui)"
	${E} "  run-gui-contrib - Launch GUI with contribution graph"
	${E} "  run-gui-summary - Launch GUI with summary"
	${E} "  run-gui-contributors - Launch GUI with contributors"
	${E} "  run-gui-health - Launch GUI with health analysis"
	${E} "  dev-gui        - Run GUI in development mode"
	${E} ""
	${E} "Benchmark Targets:"
	${E} "  bench          - Run all benchmarks"
	${E} "  bench-analyzers - Run analyzer benchmarks"
	${E} "  bench-filters  - Run filter benchmarks"
	${E} ""
	${E} "Code Quality Targets:"
	${E} "  fmt            - Format Go code"
	${E} "  vet            - Run go vet"
	${E} "  lint           - Run golint (if available)"
	${E} "  check          - Run fmt, vet, and test"
	${E} ""
	${E} "Dependency Targets:"
	${E} "  deps           - Download dependencies"
	${E} "  deps-gui       - Install GUI dependencies (tcell, tview)"
	${E} "  deps-gui-offline - Check existing GUI dependencies"
	${E} "  check-gui-deps - Check GUI dependencies availability"
	${E} "  deps-update    - Update dependencies"
	${E} "  deps-vendor    - Vendor dependencies"
	${E} ""
	${E} "Installation Targets:"
	${E} "  install        - Install git-stats (and git-stats-gui if available)"
	${E} "  install-gui    - Install git-stats-gui only"
	${E} "  install-all    - Install both git-stats and git-stats-gui"
	${E} "  uninstall      - Remove both binaries from /usr/local/bin"
	${E} ""
	${E} "Distribution Targets:"
	${E} "  dist           - Create distribution binaries"
	${E} "  dist-gui       - Create GUI distribution binaries"
	${E} "  dist-all       - Create both regular and GUI distributions"
	${E} ""
	${E} "Clean Targets:"
	${E} "  clean          - Clean build artifacts"
	${E} "  clean-all      - Clean all generated files"
	${E} ""
	${E} "Other:"
	${E} "  help           - Show this help message"
