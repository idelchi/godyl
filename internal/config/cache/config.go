// Package cache provides the configuration for the cache subcommand and its subcommands.
package cache

import (
	"github.com/idelchi/godyl/internal/config/common"
)

// Cache represents the configuration for cache-related commands.
type Cache struct {
	common.Tracker `mapstructure:"-" yaml:"-"`
}
