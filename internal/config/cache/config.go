// Package cache provides the configuration for the cache subcommand and its subcommands.
package cache

import (
	"github.com/idelchi/godyl/internal/config/common"
)

type Cache struct {
	common.Tracker `json:"-" mapstructure:"-"`
}
