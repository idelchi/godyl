package config

import (
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/tool"
)

// Config holds the top level configuration for godyl.
// It is split into sub-structs for each command.

type Config struct {
	// Root level configuration, mapping configurations on the root `godyl` command
	Root Root `mapstructure:"godyl"`

	// Install level configuration, mapping configurations on the `install` command
	Install Install

	// Download level configuration, mapping configurations on the `download` command
	Download Download

	// Update level configuration, mapping configurations on the `update` command
	Update Update

	// Dump level configuration, mapping configurations on the `dump` command
	Dump Dump

	// Cache level configuration, mapping configurations on the `cache` command
	Cache Cache

	// Status level configuration, mapping configurations on the `status` command
	Status Status

	Common Common `mapstructure:"-" yaml:"-"`
}

func (c *Config) SetCommon(common Common) {
	c.Common = common
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
		tool.Source.URL.Token.Token = c.Root.Tokens.URL.Token
	}

	if isSet(&c.Root)("url-token-header") {
		tool.Source.URL.Token.Header = c.Root.Tokens.URL.Header
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
		tool.Hints.Add(match.Hint{
			Pattern: hint,
			Weight:  "1",
		})
	}

	return &tool
}
