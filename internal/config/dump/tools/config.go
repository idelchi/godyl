package tools

import (
	"github.com/idelchi/godyl/internal/config/common"
)

// Tools holds the configuration for the `dump tools` subcommand.
type Tools struct {
	// Tags are the tags to consider when dumping tools
	Tags []string `json:"tags" mapstructure:"tags"`

	// Full indicates whether to dump full tool information
	Full bool `json:"full" mapstructure:"full"`

	// Embedded indicates whether to dump the embedded tools
	Embedded bool `json:"embedded" mapstructure:"embedded"`

	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	common.Tracker `json:"-" mapstructure:"-"`
}
