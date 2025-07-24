// Package config provides configuration management functionality.
// Package config provides the configuration for the config subcommand.
package config

import "github.com/idelchi/godyl/internal/config/common"

// Config represents the application configuration structure.
// Config represents the configuration for config-related commands.
type Config struct {
	common.Tracker `mapstructure:"-" yaml:"-"`
}
