// Package config provides configuration and flags for the godyl application.
package config

import (
	"github.com/idelchi/godyl/internal/config/cache"
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/config/config"
	"github.com/idelchi/godyl/internal/config/download"
	"github.com/idelchi/godyl/internal/config/dump"
	"github.com/idelchi/godyl/internal/config/install"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/config/status"
	"github.com/idelchi/godyl/internal/config/update"
	"github.com/idelchi/godyl/internal/tools/hints"
	"github.com/idelchi/godyl/internal/tools/tool"
)

// Config holds the top level configuration for godyl.
// It is split into sub-structs for each command.
type Config struct {
	// Root contains the configuration for the `godyl` command.
	Root root.Root `json:"godyl" mapstructure:"godyl"`

	// Dump contains the configuration for the `godyl dump` command
	Dump dump.Dump `json:"dump" mapstructure:"dump"`

	// Cache contains the configuration for the `godyl cache` command (empty)
	Cache cache.Cache `json:"-" mapstructure:"-"`

	// Config contains the configuration for the `godyl config` command (empty)
	Config config.Config `json:"-" mapstructure:"-"`

	// Update contains the configuration for the `godyl update` command
	Update update.Update `json:"update" mapstructure:"update"`

	// Status contains the configuration for the `godyl status` command
	Status status.Status `json:"status" mapstructure:"status"`

	// Download contains the configuration for the `godyl download` command
	Download download.Download `json:"download" mapstructure:"download"`

	// Install contains the configuration for the `godyl install` command
	Install install.Install `json:"install" mapstructure:"install"`

	// Common contains a subset of common configuration options
	Common common.Common `json:"-" mapstructure:"-"`
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

	if isSet(&c.Root)("github-token") {
		tool.Source.GitHub.Token = c.Root.Tokens.GitHub
	}

	if isSet(&c.Root)("gitlab-token") {
		tool.Source.GitLab.Token = c.Root.Tokens.GitLab
	}

	if isSet(&c.Root)("url-token") {
		tool.Source.URL.Token = c.Root.Tokens.URL
	}

	if isSet(&c.Root)("no-cache") {
		tool.NoCache = c.Root.Cache.Disabled
	}

	if isSet(&c.Root)("no-verify-ssl") {
		tool.NoVerifySSL = c.Root.NoVerifySSL
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
