package defaults

import (
	"fmt"

	"github.com/idelchi/godyl/internal/tools/tool"
)

// MergeWith merges all the stored defaults with the provided tools.
// The sequence is:
//
//	others[0] <-- others[1] <-- others[2]... <-- d[name]
func (d *Defaults) MergeWith(others ...*tool.Tool) error {
	for name, t := range *d {
		if t == nil {
			return fmt.Errorf("tool %q is nil", name)
		}

		if err := t.MergeInto(others...); err != nil {
			return fmt.Errorf("merging %q into %q: %w", name, t.Name, err)
		}
	}

	return nil
}

// MergeFrom merges the current tool with the provided tools,
// the sequence being.
//
//	d[name] <- others[0] <-- others[1] <-- others[2]...
func (d *Defaults) MergeFrom(others ...*tool.Tool) error {
	if len(others) == 0 {
		return fmt.Errorf("no defaults to merge")
	}

	for name, t := range *d {
		if t == nil {
			return fmt.Errorf("tool %q is nil", name)
		}

		if err := t.MergeFrom(others...); err != nil {
			return fmt.Errorf("merging %q into %q: %w", name, t.Name, err)
		}
	}

	return nil
}
