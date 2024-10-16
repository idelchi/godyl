package tools

import (
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Aliases represents a tool's alias names, allowing the configuration to either
// be a single alias string or a slice of alias strings.
// This flexibility is useful when configuring tools that might have multiple alternative names.
type Aliases = unmarshal.SingleOrSlice[string]
