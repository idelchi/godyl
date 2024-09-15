package root

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Root holds the root configuration options.
type Root struct {
	common.Trackable `json:"-"             mapstructure:"-"`
	Tokens           Tokens      `json:"tokens"        mapstructure:",squash"`
	Inherit          string      `json:"inherit"       mapstructure:"inherit"`
	Defaults         file.File   `json:"defaults"      mapstructure:"defaults"`
	ConfigFile       file.File   `json:"config-file"   mapstructure:"config-file"`
	ErrorFile        file.File   `json:"error-file"    mapstructure:"error-file"`
	Tools            string      `json:"tools"         mapstructure:"tools"`
	LogLevel         string      `json:"log-level"     mapstructure:"log-level"     validate:"oneof=silent debug info warn error always"`
	EnvFile          []file.File `json:"env-file"      mapstructure:"env-file"`
	Cache            Cache       `json:"cache"         mapstructure:",squash"`
	Parallel         int         `json:"parallel"      mapstructure:"parallel"      validate:"gte=0"`
	Verbose          int         `json:"verbose"       mapstructure:"verbose"`
	Show             Verbosity   `json:"show"          mapstructure:"show"`
	NoVerifySSL      bool        `json:"no-verify-ssl" mapstructure:"no-verify-ssl"`
	NoProgress       bool        `json:"no-progress"   mapstructure:"no-progress"`
}
