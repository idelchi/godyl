package config

import (
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/pkg/path/files"
)

// Tool holds the configuration options for fetching tools.
// These are used as flags, environment variables for the corresponding CLI commands,
// and used to set the tool configuration for each tool requested, unless explicitly set by the tool itself.
type Tool struct {
	viperable   `json:"-" mapstructure:"-" yaml:"-"`
	Version     string
	Source      sources.Type      `validate:"oneof=github gitlab url go command"`
	Strategy    strategy.Strategy `validate:"oneof=none sync force"`
	OS          string
	Arch        string
	Output      string
	Tools       files.Files
	Hints       []string
	Tags        []string
	Parallel    int  `validate:"gte=0"`
	NoVerifySSL bool `mapstructure:"no-verify-ssl"`
	Show        bool
	NoCache     bool `mapstructure:"no-cache"`
}
