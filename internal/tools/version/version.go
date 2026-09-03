// Package version provides functionality for managing tool version information.
package version

import (
	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Version defines the target version of a tool, as well as how it can be parsed.
type Version struct {
	// Commands contains commands used to discover an installed version.
	Commands *Commands `mapstructure:"commands" yaml:"commands"`
	// Patterns contains expressions used to extract a version from command output.
	Patterns *Patterns `mapstructure:"patterns" yaml:"patterns"`
	// Version is the desired version expression.
	Version string `mapstructure:"version"  single:"true"   yaml:"version"`
}

type (
	// Patterns represents version pattern matching rules.
	Patterns = unmarshal.SingleOrSliceType[string]
	// Commands represents version extraction commands.
	Commands = unmarshal.SingleOrSliceType[string]
)

// UnmarshalYAML implements the yaml.Unmarshaler interface for Version.
func (v *Version) UnmarshalYAML(node ast.Node) error {
	type raw Version

	return unmarshal.SingleStringOrStruct(node, (*raw)(v))
}
