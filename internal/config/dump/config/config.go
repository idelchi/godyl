package config

import (
	"github.com/idelchi/godyl/internal/config/common"
)

// Config represents the configuration for the `dump config` command.
type Config struct {
	// All is a boolean flag that indicates whether all configuration options should be dumped.
	All bool `json:"all" mapstructure:"all"`

	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	common.Tracker `json:"-" mapstructure:"-"`
}
