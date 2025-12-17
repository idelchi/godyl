// Package shared provides shared configuration structures and utilities used across different commands.
package shared

import (
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Common contains configuration fields shared across multiple commands.
// Common represents a shared configuration structure that provides
// command-line arguments, show functionality, and validation.
type Common struct {
	Tracker `mapstructure:"-" yaml:"-"`

	Output   string
	Strategy strategy.Strategy
	Source   sources.Type
	OS       string
	Arch     string
	Hints    []string
	Pre      bool
}
