// Package defaults provides functionality for managing default values and configurations.
package defaults

import (
	"fmt"

	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Defaults represents a collection of default values for tools.
type Defaults map[string]*tool.Tool

// NewDefaults initializes a new Defaults instance from the provided data.
func NewDefaultsFromBytes(data []byte) (*Defaults, error) {
	var defaults Defaults

	if err := defaults.Load(data); err != nil {
		return nil, fmt.Errorf("loading defaults: %w", err)
	}

	return &defaults, nil
}

// NewDefaults initializes the Defaults instance from the provided file.
func (d *Defaults) Load(data []byte) error {
	if err := unmarshal.Strict(data, d); err != nil {
		return fmt.Errorf("defaults: %w", err)
	}

	return nil
}

// get returns the tool with the given name from the defaults map.
func (d Defaults) get(name string) (*tool.Tool, error) {
	t, ok := d[name]
	if !ok {
		return nil, fmt.Errorf("%q not found in defaults", name)
	}

	debug.Debug("Found %q in defaults", name)

	return t, nil
}

// Get returns the tool with the given name from the defaults map.
// If the tool is not found, it returns nil.
func (d Defaults) Get(name string) *tool.Tool {
	t, _ := d.get(name)

	return t
}

// Pick returns a slice of tools from the defaults map based on the provided names.
func (d Defaults) Pick(names ...string) (tools []*tool.Tool, err error) {
	for _, name := range names {
		t, err := d.get(name)
		if err != nil {
			return nil, err
		}

		debug.Debug("Found %q in defaults", name)

		tools = append(tools, t)

	}

	return tools, nil
}
