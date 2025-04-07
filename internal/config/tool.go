package config

import (
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/path/files"
)

// Tool holds the configuration options for fetching tools.
// These are used as flags, environment variables for the corresponding CLI commands,
// and used to set the tool configuration for each tool requested, unless explicitly set by the tool itself.
type Tool struct {
	// Path to output the fetched tools to
	Output string

	// Tags to filter tools by
	Tags []string

	// Source from which to install the tools
	Source sources.Type `validate:"oneof=github url go command"`

	// Strategy to use for updating tools
	Strategy tools.Strategy `validate:"oneof=none upgrade force"`

	// Operating system to install the tools for
	OS string

	// Architecture to install the tools for
	Arch string

	// Path to tools configuration file
	Tools files.Files // Positional argument

	// Number of parallel downloads (>= 0)
	Parallel int `validate:"gte=0"`

	// Skip SSL verification
	NoVerifySSL bool `mapstructure:"no-verify-ssl"`

	// Additional hints to use for tool resolution
	Hints []string

	// Version of the tool to install
	Version string

	// Show enables output display
	Show bool

	// Viper instance
	viperable `mapstructure:"-" yaml:"-" json:"-"`
}
