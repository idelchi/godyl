package tools

import (
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/unmarshal"
	"github.com/idelchi/godyl/pkg/utils"
	"gopkg.in/yaml.v3"
)

// Tool represents a single tool configuration
type Tool struct {
	// Name of the tool
	Name string
	// Description of the tool
	Description string
	// Version of the tool
	Version string
	// Path to fetch the tool
	Path string
	// Checksum
	Checksum string
	// Output path for the tool
	Output string
	// Name of the executable itself
	Exe          Exe
	Platform     detect.Platform
	Aliases      Aliases
	Values       map[string]any
	Fallbacks    Fallbacks
	Hints        match.Hints
	Source       sources.Source
	Tags         Tags
	Strategy     Strategy
	Extensions   Extensions
	Skip         Skip
	Test         sources.Commands
	AllowFailure bool `yaml:"allow_failure" mapstructure:"allow_failure"`
	Post         sources.Commands
	Mode         Mode
	Settings     Settings
	Env          env.Env
}

// UnmarshalYAML implements custom unmarshaling for Tool with KnownFields check
func (t *Tool) UnmarshalYAML(value *yaml.Node) error {
	// If it's a scalar (e.g., just the name), handle it directly
	if value.Kind == yaml.ScalarNode {
		t.Name = value.Value

		return nil
	}

	type rawTool Tool

	return unmarshal.DecodeWithOptionalKnownFields(value, (*rawTool)(t), true, t)
}

func (t *Tool) ApplyDefaults(d Defaults) {
	utils.SetIfEmpty(&t.Output, d.Output)
	utils.SetIfEmpty(&t.Source.Type, d.Source.Type)
	utils.SetIfEmpty(&t.Source.Github.Token, d.Source.Github.Token)
	utils.SetIfEmpty(&t.Strategy, d.Strategy)
	utils.SetSliceIfNil(&t.Skip, Condition{Condition: "false"})
	utils.SetIfEmpty(&t.Mode, d.Mode)
	utils.SetSliceIfNil(&t.Exe.Patterns, d.Exe.Patterns...)
	utils.SetSliceIfNil(&t.Extensions, d.Extensions...)
	utils.SetMapIfNil(&t.Values, d.Values)
	utils.DeepMergeMapsWithoutOverwrite(t.Values, d.Values)
	t.Env.Merge(d.Env)

	t.Platform.Merge(d.Platform)
	t.Hints.Add(d.Hints)
}
