package tools

import (
	"bytes"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
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
	Exe          string
	Platform     detect.Platform
	Aliases      Aliases
	Values       map[string]any
	Fallbacks    []string
	Hints        match.Hints
	Source       sources.Source
	Tags         Tags
	Strategy     Strategy
	Extensions   []string
	SkipTemplate string `yaml:"skip" mapstructure:"skip"`
	Skip         bool   `yaml:"-" mapstructure:"-"`
	Test         sources.Commands
}

// UnmarshalYAML implements custom unmarshaling for Tool with KnownFields check
func (t *Tool) UnmarshalYAML(value *yaml.Node) error {
	// If it's a scalar (e.g., just the name), handle it directly
	if value.Kind == yaml.ScalarNode {
		t.Name = value.Value
		return nil
	}

	// Re-encode the yaml.Node to bytes to leverage yaml.NewDecoder
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return err
	}
	enc.Close()

	// Now decode from the buffer with KnownFields enabled
	decoder := yaml.NewDecoder(&buf)
	decoder.KnownFields(true)

	// Decode the Tool
	type rawTool Tool
	if err := decoder.Decode((*rawTool)(t)); err != nil {
		return err
	}

	return nil
}
