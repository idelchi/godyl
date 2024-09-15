// Package version provides functionality for managing tool version information.
package version

import (
	"github.com/goccy/go-yaml/ast"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Version defines the target version of a tool, as well as how it can be parsed.
type Version struct {
	// Version is the explicit version string or the parsed result.
	Version string `single:"true"`

	// Commands is a list of shell commands that can be used to determine the version.
	// Each command's output is matched against the version patterns.
	Commands *Commands

	// Patterns contains regex patterns for extracting version strings.
	// Used to parse version information from command output or other sources.
	Patterns *Patterns
}

type (
	Patterns = unmarshal.SingleOrSliceType[string]
	Commands = unmarshal.SingleOrSliceType[string]
)

func (v *Version) UnmarshalYAML(node ast.Node) error {
	type raw Version

	return unmarshal.SingleStringOrStruct(node, (*raw)(v))
}
