// Copyright (c) 2019 Sunil
// Enhanced git-stats tool - CLI package for input validation

package cli

import (
	"time"
)

// Validator interface for input validation
type Validator interface {
	ValidateConfig(config *Config) error
	ValidateDateRange(since, until *time.Time) error
	ValidateAuthor(author string) error
	ValidateFormat(format string) error
	ValidateOutputFile(filename string) error
}
