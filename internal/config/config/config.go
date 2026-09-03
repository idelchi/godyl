// Package config provides the configuration for the config subcommand.
package config

import "github.com/idelchi/godyl/internal/config/shared"

// Config represents the configuration for config-related commands.
type Config struct {
	// Tracker records which config options were supplied explicitly.
	shared.Tracker `mapstructure:"-" yaml:"-"`
}
