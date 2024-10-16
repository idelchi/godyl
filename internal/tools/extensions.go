package tools

import (
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Extensions represents a collection of file extensions related to the tool.
// It can either be a single string or a slice of strings, providing flexibility
// when configuring tools that may involve multiple file types.
type Extensions = unmarshal.SingleOrSlice[string]
