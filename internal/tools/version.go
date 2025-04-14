package tools

import (
	"github.com/fatih/structs"

	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// Version defines how a tool's version is determined and extracted.
type Version struct {
	// Version is the explicit version string or the parsed result.
	Version string

	// Commands is a list of shell commands that can be used to determine the version.
	// Each command's output is matched against the version patterns.
	Commands unmarshal.SingleOrSliceType[string]

	// Patterns contains regex patterns for extracting version strings.
	// Used to parse version information from command output or other sources.
	Patterns unmarshal.SingleOrSliceType[string]
}

// UnmarshalYAML implements custom YAML unmarshaling for Version configuration.
// Supports both scalar values (treated as explicit version) and map values with
// additional configuration for version detection and extraction.
func (v *Version) UnmarshalYAML(value *yaml.Node) error {
	// If the YAML value is a scalar (e.g., just the version), handle it directly by setting the Version field.
	if value.Kind == yaml.ScalarNode {
		v.Version = value.Value

		return nil
	}

	// Perform custom unmarshaling with field validation, allowing only known fields.
	type raw Version

	return unmarshal.DecodeWithOptionalKnownFields(value, (*raw)(v), true, structs.New(v).Name())
}
