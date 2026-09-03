// Package cache provides the configuration for the cache subcommand and its subcommands.
package cache

import (
	"github.com/idelchi/godyl/internal/config/shared"
)

// Cache represents the configuration for cache-related commands.
type Cache struct {
	// Tracker records which cache options were supplied explicitly.
	shared.Tracker `mapstructure:"-" yaml:"-"`
}
