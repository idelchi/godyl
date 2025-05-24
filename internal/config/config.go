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
	Dump     dump.Dump         `json:"dump"     mapstructure:"dump"`
	Cache    cache.Cache       `json:"cache"    mapstructure:"cache"`
	Config   config.Config     `json:"config"   mapstructure:"config"`
	Common   common.Common     `json:"-"        mapstructure:"-"`
	Update   update.Update     `json:"update"   mapstructure:"update"`
	Status   status.Status     `json:"status"   mapstructure:"status"`
	Download download.Download `json:"download" mapstructure:"download"`
	Install  install.Install   `json:"install"  mapstructure:"install"`
	Root     root.Root         `json:"godyl"    mapstructure:"godyl"`
}

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
