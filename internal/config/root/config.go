// Package root provides the root configuration structure for the application.
package root

import (
	"github.com/idelchi/godyl/internal/config/download"
	"github.com/idelchi/godyl/internal/config/dump"
	"github.com/idelchi/godyl/internal/config/install"
	"github.com/idelchi/godyl/internal/config/shared"
	"github.com/idelchi/godyl/internal/config/status"
	"github.com/idelchi/godyl/internal/config/update"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// TODO(Idelchi): Change all to be .Config instead of .Dump, .Update, etc. //nolint:godox // TODO comment provides
// valuable context for future development

// Config holds the root level configuration options.
// It is split into sub-structs for each command.
type Config struct {
	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	shared.Tracker `mapstructure:"-" yaml:"-"`

	/* Subcommands */
	// Dump contains the configuration for the `godyl dump` command
	Dump dump.Dump `mapstructure:"dump" validate:"-" yaml:"dump"`

	// Cache contains the configuration for the `godyl cache` command (empty)
	// Cache cache.Cache `yaml:"-" mapstructure:"-" validate:"-"`

	// Config contains the configuration for the `godyl config` command (empty)
	// Config root.Config `yaml:"-" mapstructure:"-" validate:"-"`

	// Update contains the configuration for the `godyl update` command
	Update update.Update `mapstructure:"update" validate:"-" yaml:"update"`

	// Status contains the configuration for the `godyl status` command
	Status status.Status `mapstructure:"status" validate:"-" yaml:"status"`

	// Download contains the configuration for the `godyl download` command
	Download download.Download `mapstructure:"download" validate:"-" yaml:"download"`

	// Install contains the configuration for the `godyl install` command
	Install install.Install `mapstructure:"install" validate:"-" yaml:"install"`

	/* Flags */
	// Tokens store authentication tokens for various sources
	Tokens Tokens `mapstructure:",squash" yaml:",inline,flatten"`

	// Inherit specifies the default scheme to inherit from when no scheme is specified
	Inherit string `mapstructure:"inherit" yaml:"inherit"`

	// ErrorFile specifies the file to log errors
	ErrorFile file.File `mapstructure:"error-file" yaml:"error-file"`

	// Tools specifies the tools file to be used
	Tools string `mapstructure:"tools" yaml:"tools"`

	// Defaults specifies the default file to be used
	Defaults file.File `mapstructure:"defaults" yaml:"defaults"`

	// ConfigFile specifies the configuration file to be used
	ConfigFile file.File `mapstructure:"config-file" yaml:"config-file"`

	// LogLevel specifies the logging level
	LogLevel string `mapstructure:"log-level" validate:"oneof=silent debug info warn error always" yaml:"log-level"`

	// EnvFile specifies the environment files to be used
	EnvFile []file.File `mapstructure:"env-file" yaml:"env-file"`

	// Cache holds the cache configuration options
	Cache Cache `mapstructure:",squash" yaml:",inline,flatten"`

	// Parallel specifies the number of parallel operations
	Parallel int `mapstructure:"parallel" validate:"gte=0" yaml:"parallel"`

	// Verbose specifies the verbosity level
	Verbose int `mapstructure:"verbose" yaml:"verbose"`

	// Show specifies the verbosity level for showing output
	Show Verbosity `mapstructure:"show" yaml:"show"`

	// NoVerifySSL disables SSL verification
	NoVerifySSL bool `mapstructure:"no-verify-ssl" yaml:"no-verify-ssl"`

	// NoProgress disables progress indicators
	NoProgress bool `mapstructure:"no-progress" yaml:"no-progress"`

	// Keyring enables the use of the keyring for retrieving tokens
	Keyring bool `mapstructure:"keyring" yaml:"keyring"`

	/* Other Options */
	// Common contains a subset of common configuration options
	Common shared.Common `mapstructure:"-" yaml:"-"`
}

// Cache holds the configuration options for caching.
type Cache struct {
	// Path to cache folder
	Dir folder.Folder `mapstructure:"cache-dir" yaml:"cache-dir"`

	// Disabled disables cache interaction
	Disabled bool `mapstructure:"no-cache" yaml:"no-cache"`
}

// Tokens holds the configuration options for authentication tokens.
type Tokens struct {
	// GitHub token for authentication
	GitHub string `mapstructure:"github-token" mask:"fixed" yaml:"github-token"`

	// GitLab token for authentication
	GitLab string `mapstructure:"gitlab-token" mask:"fixed" yaml:"gitlab-token"`

	// URL token for authentication
	URL string `mapstructure:"url-token" mask:"fixed" yaml:"url-token"`
}

// AllTokensSet checks if all of the tokens are set.
func (c *Config) AllTokensSet() bool {
	return c.IsSet("github-token") && c.IsSet("gitlab-token") && c.IsSet("url-token")
}

// ToTool converts the Config to a tool.Tool instance,
// holding either default values or values set by the user.
func (c *Config) ToTool(forced bool) *tool.Tool {
	var tool tool.Tool

	type Settable interface {
		IsSet(string) bool
	}

	isSet := func(settable Settable) func(name string) bool {
		if forced {
			return func(_ string) bool {
				return true
			}
		}

		if settable == nil {
			panic("settable is nil")
		}

		return settable.IsSet
	}

	if isSet(c)("github-token") {
		tool.Source.GitHub.Token = c.Tokens.GitHub
	}

	if isSet(c)("gitlab-token") {
		tool.Source.GitLab.Token = c.Tokens.GitLab
	}

	if isSet(c)("url-token") {
		tool.Source.URL.Token = c.Tokens.URL
	}

	if isSet(c)("no-cache") {
		tool.NoCache = c.Cache.Disabled
	}

	if isSet(c)("no-verify-ssl") {
		tool.NoVerifySSL = c.NoVerifySSL
	}

	if isSet(&c.Common)("output") {
		tool.Output = c.Common.Output
	}

	if isSet(&c.Common)("strategy") {
		tool.Strategy = c.Common.Strategy
	}

	if isSet(&c.Common)("source") {
		tool.Source.Type = c.Common.Source
	}

	if isSet(&c.Common)("os") {
		tool.Platform.OS.Name = c.Common.OS
	}

	if isSet(&c.Common)("arch") {
		tool.Platform.Architecture.Name = c.Common.Arch
	}

	if isSet(&c.Common)("pre") {
		tool.Source.GitHub.Pre = c.Common.Pre
		tool.Source.GitLab.Pre = c.Common.Pre
	}

	return &tool
}
