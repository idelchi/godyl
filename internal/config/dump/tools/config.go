package tools

import (
	"github.com/idelchi/godyl/internal/config/common"
)

// Tools holds the configuration for the `dump tools` subcommand.
type Tools struct {
	common.Tracker `json:"-" mapstructure:"-"`
	Tags           []string `json:"tags"     mapstructure:"tags"`
	Full           bool     `json:"full"     mapstructure:"full"`
	Embedded       bool     `json:"embedded" mapstructure:"embedded"`
}
