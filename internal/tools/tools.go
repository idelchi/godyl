package tools

import (
	"fmt"

	"github.com/idelchi/godyl/internal/debug"
	defaults "github.com/idelchi/godyl/internal/defaults"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/tools/inherit"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/generic"
)

// Tools represents a collection of Tool configurations.
type Tools []*tool.Tool

// Append adds a tool to the collection.
func (ts *Tools) Append(t *tool.Tool) {
	if t == nil {
		panic("nil tool in tools collection")
	}

	*ts = append(*ts, t)
}

// MergeWith merges all the stored defaults with the provided tools.
// The sequence is:
//
//	others[0] <-- others[1] <-- others[2]... <-- tools[i]
func (ts Tools) MergeWith(others ...*tool.Tool) error {
	for _, t := range ts {
		if t == nil {
			panic("nil tool in tools collection")
		}

		if err := t.MergeInto(others...); err != nil {
			return fmt.Errorf("merging %q into %q: %w", t.Name, others[0].Name, err)
		}
	}

	return nil
}

// Get retrieves a tool by name from the collection.
func (ts Tools) Get(name string) *tool.Tool {
	for _, t := range ts {
		if t.Name == name {
			return t
		}
	}

	return nil
}

// GetFirst returns the first tool in the collection.
func (ts Tools) GetFirst() *tool.Tool {
	if len(ts) > 0 {
		return ts[0]
	}

	panic("no tools in collection")
}

// DefaultInheritance sets default inheritance for tools that don't have it specified.
func (ts Tools) DefaultInheritance(inheritance string) {
	for _, t := range ts {
		if t.Inherit == nil {
			t.Inherit = &inherit.Inherit{inheritance}
		}
	}
}

// ResolveInheritance processes tool inheritance using the provided defaults.
func (ts Tools) ResolveInheritance(d *defaults.Defaults) error {
	for _, t := range ts {
		if t == nil {
			panic("nil tool in tools collection")
		}

		if generic.IsSliceNilOrEmpty(t.Inherit) {
			debug.Debug("No inheritance for %q", t.Name)

			continue
		}

		inherits, err := d.Pick(*t.Inherit...)
		if err != nil {
			return fmt.Errorf("resolving inheritance for %q: %w", t.Name, err)
		}

		debug.Debug("Inherits for %q: %v", t.Name, *t.Inherit)

		// Construct the default from the inherits
		toolDefault, err := tool.MergeRightToLeft(inherits...)
		if err != nil {
			return fmt.Errorf("merging %q into %v: %w", t.Name, *t.Inherit, err)
		}

		// Merge the tool with the defaults
		if err := t.MergeWithOther(toolDefault); err != nil {
			return fmt.Errorf("merging %q into %v: %w", t.Name, *t.Inherit, err)
		}
	}

	return nil
}

// ResolveNilPointers ensures all tools have non-nil pointer fields.
func (ts Tools) ResolveNilPointers() error {
	for _, t := range ts {
		emptyTool := tool.NewEmptyTool()

		if err := t.MergeInto(emptyTool); err != nil {
			return fmt.Errorf("merging %q into %v: %w", t.Name, t.Inherit, err)
		}
	}

	return nil
}

// MergePlatform merges detected platform information into all tools.
func (ts Tools) MergePlatform() error {
	platform := detect.Platform{}
	if err := platform.Detect(); err != nil {
		return fmt.Errorf("detecting platform: %w", err)
	}

	for _, t := range ts {
		t.MergePlatform(platform)
	}

	return nil
}

// Copy performs deep copying of all tools in the collection.
func (ts Tools) Copy() error {
	for _, t := range ts {
		err := t.Copy()
		if err != nil {
			return fmt.Errorf("deep copying tool %q: %w", t.Name, err)
		}
	}

	return nil
}
