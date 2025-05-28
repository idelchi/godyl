package config

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/config/download"
	"github.com/idelchi/godyl/internal/config/dump"
	"github.com/idelchi/godyl/internal/config/install"
	"github.com/idelchi/godyl/internal/config/status"
	"github.com/idelchi/godyl/internal/config/update"
	"github.com/idelchi/godyl/internal/tools/hints"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Config holds the top level configuration for godyl.
// It is split into sub-structs for each command.
type Config struct {
	/* Subcommands */
	// Dump contains the configuration for the `godyl dump` command
	Dump dump.Dump `json:"dump" mapstructure:"dump" validate:"-"`

	// Cache contains the configuration for the `godyl cache` command (empty)
	// Cache cache.Cache `json:"-" mapstructure:"-" validate:"-"`

	// Config contains the configuration for the `godyl config` command (empty)
	// Config config.Config `json:"-" mapstructure:"-" validate:"-"`

	// Update contains the configuration for the `godyl update` command
	Update update.Update `json:"update" mapstructure:"update" validate:"-"`

	// Status contains the configuration for the `godyl status` command
	Status status.Status `json:"status" mapstructure:"status" validate:"-"`

	// Download contains the configuration for the `godyl download` command
	Download download.Download `json:"download" mapstructure:"download" validate:"-"`

	// Install contains the configuration for the `godyl install` command
	Install install.Install `json:"install" mapstructure:"install" validate:"-"`

	/* Flags */
	// Tokens store authentication tokens for various sources
	Tokens Tokens `json:"tokens" mapstructure:",squash"`

	// Inherit specifies the default scheme to inherit from when no scheme is specified
	Inherit string `json:"inherit" mapstructure:"inherit"`

	// ErrorFile specifies the file to log errors
	ErrorFile file.File `json:"error-file" mapstructure:"error-file"`

	// Tools specifies the tools file to be used
	Tools string `json:"tools" mapstructure:"tools"`

	// Defaults specifies the default file to be used
	Defaults file.File `json:"defaults" mapstructure:"defaults"`

	// ConfigFile specifies the configuration file to be used
	ConfigFile file.File `json:"config-file" mapstructure:"config-file"`

	// LogLevel specifies the logging level
	LogLevel string `json:"log-level" mapstructure:"log-level" validate:"oneof=silent debug info warn error always"`

	// EnvFile specifies the environment files to be used
	EnvFile []file.File `json:"env-file" mapstructure:"env-file"`

	// Cache holds the cache configuration options
	Cache Cache `json:"cache" mapstructure:",squash"`

	// Parallel specifies the number of parallel operations
	Parallel int `json:"parallel" mapstructure:"parallel" validate:"gte=0"`

	// Verbose specifies the verbosity level
	Verbose int `json:"verbose" mapstructure:"verbose"`

	// Show specifies the verbosity level for showing output
	Show Verbosity `json:"show" mapstructure:"show"`

	// NoVerifySSL disables SSL verification
	NoVerifySSL bool `json:"no-verify-ssl" mapstructure:"no-verify-ssl"`

	// NoProgress disables progress indicators
	NoProgress bool `json:"no-progress" mapstructure:"no-progress"`

	/* Other Options */
	// Common contains a subset of common configuration options
	Common common.Common `json:"-" mapstructure:"-"`

	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	common.Tracker `json:"-" mapstructure:"-"`
}

// Cache holds the configuration options for caching.
type Cache struct {
	// Path to cache folder
	Dir folder.Folder `json:"cache-dir" mapstructure:"cache-dir"`

	// Disabled disables cache interaction
	Disabled bool `json:"no-cache" mapstructure:"no-cache"`
}

// Tokens holds the configuration options for authentication tokens.
type Tokens struct {
	// GitHub token for authentication
	GitHub string `json:"github-token" mapstructure:"github-token" mask:"fixed"`

	// GitLab token for authentication
	GitLab string `json:"gitlab-token" mapstructure:"gitlab-token" mask:"fixed"`

	// URL token for authentication
	URL string `json:"url-token" mapstructure:"url-token" mask:"fixed"`
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
			return func(name string) bool {
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

	for _, hint := range c.Common.Hints {
		tool.Hints.Add(hints.Hint{
			Pattern: hint,
		})
	}

	return &tool
}
