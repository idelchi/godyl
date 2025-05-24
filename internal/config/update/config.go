package update

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Update holds the configuration options for self-updating the tool.
// These are used as flags, environment variables for the corresponding CLI commands.
type Update struct {
	common.Tracker `json:"-"       mapstructure:"-"`
	Version        string `json:"version" mapstructure:"version"`
	Pre            bool   `json:"pre"     mapstructure:"pre"`
	Check          bool   `json:"check"   mapstructure:"check"`
	Cleanup        bool   `json:"cleanup" mapstructure:"cleanup"`
	Force          bool   `json:"force"   mapstructure:"force"`
}

func (u Update) ToCommon() common.Common {
	s := strategy.Sync
	if u.Force {
		s = strategy.Force
	}

	return common.Common{
		Strategy: s,

		Tracker: u.Tracker,
	}
}
