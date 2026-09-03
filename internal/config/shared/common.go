// Package shared provides shared configuration structures and utilities used across different commands.
package shared

import (
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Common contains installation and download options shared across commands.
type Common struct {
	// Tracker records which common options were supplied explicitly.
	Tracker `mapstructure:"-" yaml:"-"`

	// Output specifies the destination path.
	Output string
	// Strategy defines how existing installations are handled.
	Strategy strategy.Strategy
	// Source selects the release provider.
	Source sources.Type
	// OS selects the target operating system.
	OS string
	// Arch selects the target architecture.
	Arch string
	// Hints help rank matching release assets.
	Hints []string
	// Pre allows prerelease versions.
	Pre bool
}
