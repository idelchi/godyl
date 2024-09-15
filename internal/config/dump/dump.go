package dump

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/iutils"
)

// Dump holds the configuration for the `dump` command.
type Dump struct {
	Tools  Tools         `json:"tools"  mapstructure:"tools"       validate:"-"`
	Config Config        `json:"config" mapstructure:"config"      validate:"-"`
	Format iutils.Format `json:"format" validate:"oneof=json yaml"`

	common.Trackable `json:"-" mapstructure:"-"`
}

type Config struct {
	common.Trackable `json:"-"   mapstructure:"-"`
	All              bool `json:"all" mapstructure:"all"`
}

// Tools holds the configuration for the `dump tools` subcommand.
type Tools struct {
	Full     bool     `json:"full"     mapstructure:"full"`
	Tags     []string `json:"tags"     mapstructure:"tags"`
	Embedded bool     `json:"embedded" mapstructure:"embedded"`

	common.Trackable `json:"-" mapstructure:"-"`
}
