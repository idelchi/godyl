// Package cache provides the configuration for the cache subcommand and its subcommands.
package cache

import (
	"github.com/idelchi/godyl/internal/config/shared"
)

// Cache represents the configuration for cache-related commands.
type Cache struct {
	shared.Tracker `mapstructure:"-" yaml:"-"`
}
