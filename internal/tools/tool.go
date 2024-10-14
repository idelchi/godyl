package tools

import (
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/unmarshal"
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
	After        sources.Commands
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
