package config

import "github.com/idelchi/godyl/internal/tools/strategy"

// Update holds the configuration options for self-updating the tool.
// These are used as flags, environment variables for the corresponding CLI commands.
type Update struct {
	Version string
	Pre     bool
	Check   bool
	Cleanup bool
	Force   bool

	trackable `json:"-" mapstructure:"-" yaml:"-"`
}

func (u Update) ToCommon() Common {
	s := strategy.Sync
	if u.Force {
		s = strategy.Force
	}

	return Common{
		Strategy: s,

		trackable: u.trackable,
	}
}
