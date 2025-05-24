package tools

import (
	"github.com/idelchi/godyl/internal/config/common"
)

// Tools holds the configuration for the `dump tools` subcommand.
type Tools struct {
	Full     bool     `json:"full"     mapstructure:"full"`
	Tags     []string `json:"tags"     mapstructure:"tags"`
	Embedded bool     `json:"embedded" mapstructure:"embedded"`

	common.Tracker `json:"-" mapstructure:"-"`
}
