package dump

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/config/dump/config"
	"github.com/idelchi/godyl/internal/config/dump/tools"
	"github.com/idelchi/godyl/internal/iutils"
)

// Dump holds the configuration for the `dump` command.
type Dump struct {
	Tools  tools.Tools   `json:"tools"  mapstructure:"tools"       validate:"-"`
	Config config.Config `json:"config" mapstructure:"config"      validate:"-"`
	Format iutils.Format `json:"format" validate:"oneof=json yaml"`

	common.Tracker `json:"-" mapstructure:"-"`
}
