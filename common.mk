# Copyright (c) 2019 Sunil
# Enhanced git-stats tool - Common Makefile definitions

# Go configuration
GO = go
GOFMT = gofmt
GOVET = go vet

# Build configuration
Q = @
E = @echo
RM = rm -f

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags with version info
LDFLAGS = -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

# File patterns
CRUFT = $(wildcard *~) $(wildcard *.tmp) $(wildcard .DS_Store)
GO_FILES = $(shell find . -name "*.go" -not -path "./vendor/*" 2>/dev/null)

# Colors for output (if terminal supports it)
ifneq (,$(findstring xterm,${TERM}))
	RED     := $(shell tput -Txterm setaf 1)
	GREEN   := $(shell tput -Txterm setaf 2)
	YELLOW  := $(shell tput -Txterm setaf 3)
	BLUE    := $(shell tput -Txterm setaf 4)
	MAGENTA := $(shell tput -Txterm setaf 5)
	CYAN    := $(shell tput -Txterm setaf 6)
	WHITE   := $(shell tput -Txterm setaf 7)
	RESET   := $(shell tput -Txterm sgr0)
else
	RED     :=
	GREEN   :=
	YELLOW  :=
	BLUE    :=
	MAGENTA :=
	CYAN    :=
	WHITE   :=
	RESET   :=
endif

# Enhanced echo with colors
SUCCESS = ${E} "${GREEN}✓${RESET}"
ERROR   = ${E} "${RED}✗${RESET}"
WARN    = ${E} "${YELLOW}⚠${RESET}"
INFO    = ${E} "${BLUE}ℹ${RESET}"

.PHONY: decruft
decruft:
	${Q} if [ -n "$(CRUFT)" ]; then \
		${E} "Removing cruft files: $(CRUFT)"; \
		${RM} -- ${CRUFT}; \
	fi
