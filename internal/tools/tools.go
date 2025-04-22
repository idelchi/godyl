// Package tools provides functionality for managing tool configurations.
package tools

import (
	"fmt"

	"github.com/idelchi/godyl/internal/defaults"
	"github.com/idelchi/godyl/internal/tools/inherit"
	"github.com/idelchi/godyl/internal/tools/tool"
)

// Tools represents a collection of Tool configurations.
type Tools []*tool.Tool

// NewToolsFromDefaults creates a tool collection from a given tool default and length.
// This will enable unmarshaling the tools into the defaults, with the following default (and custom) properties:
// - defaults will populate the missing fields in tool
// - missing fields in the tool will be populated by those of the defaults
// - maps will be merged
// - slices will not be combined, and must be handled by custom logic
// Above is default `yaml.Unmarshal` behavior.
func NewToolsFromDefaults(d *defaults.Defaults, inherits []inherit.Inherit) (collection Tools, err error) {
	collection = make(Tools, len(inherits))

	for i, inherit := range inherits {
		def := d.GetDefault(string(inherit))
		if def == nil {
			return nil, fmt.Errorf("no default found for %q", inherit)
		}

		if tool, err := def.ToTool().Copy(); err != nil {
			return nil, err
		} else {
			collection[i] = &tool
		}
	}

	return collection, nil
}
