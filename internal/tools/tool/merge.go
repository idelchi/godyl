package tool

import (
	"errors"
	"fmt"
	"slices"

	"github.com/goccy/go-yaml"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// MergeFrom merges the current tool with the provided tools,
// the sequence being.
//
//	tool <-- others[0] <-- others[1] <-- others[2]...
func (t *Tool) MergeFrom(others ...*Tool) error {
	if len(others) == 0 {
		return errors.New("no defaults to merge")
	}

	for _, o := range others {
		if err := iutils.Merge(t, o); err != nil {
			return fmt.Errorf("merging %q into %q: %w", o.Name, t.Name, err)
		}
	}

	return nil
}

// MergeInto the current tool into the provided tools,
// the sequence being.
//
//	others[0] <-- others[1] <-- others[2]... <-- tool
func (t *Tool) MergeInto(others ...*Tool) error {
	if len(others) == 0 {
		return errors.New("no defaults to merge")
	}

	// Create an empty tool first to not override the first one.
	copied := &Tool{}
	// Append the current tool into the list to have it at the end.
	others = append(others, t)

	for _, o := range others {
		if err := iutils.Merge(copied, o); err != nil {
			return fmt.Errorf("merging %q into %q: %w", o.Name, copied.Name, err)
		}
	}

	*t = *copied

	return nil
}

// MergeWithOther behaves like MergeInto, with special handling for slices that need to be appended instead of replaced.
func (t *Tool) MergeWithOther(other *Tool) error {
	copied := &Tool{}

	if err := copied.MergeFrom(other, t); err != nil {
		return fmt.Errorf("merging %q into %q: %w", other.Name, copied.Name, err)
	}

	// copied will now contain `other` <- `t`, where `t` has overwritten the values of `other`.

	// Reapply the hints from `other`, only if t.Hints was not set with len(0)
	if other.Hints != nil && t.Hints.Has() {
		copied.Hints.Append(*other.Hints)
	}

	copied.Env.Merge(other.Env)
	copied.Values.Merge(other.Values)

	*t = *copied

	return nil
}

// MergeRightToLeft merges the provided tools into a new tool,
// right to left.
func MergeRightToLeft(others ...*Tool) (*Tool, error) {
	if len(others) == 0 {
		return nil, errors.New("no defaults to merge")
	}

	// Create an empty tool
	tool := &Tool{}

	for _, o := range others {
		if err := iutils.Merge(tool, o); err != nil {
			return nil, fmt.Errorf("merging %q into %q: %w", o.Name, tool.Name, err)
		}
	}

	return tool, nil
}

// MergeLeftToRight merges the provided tools into a new tool,
// left to right.
func MergeLeftToRight(others ...*Tool) (*Tool, error) {
	if len(others) == 0 {
		return nil, errors.New("no defaults to merge")
	}

	slices.Reverse(others)

	// Create an empty tool
	tool := &Tool{}

	for _, o := range others {
		if err := iutils.Merge(tool, o); err != nil {
			return nil, fmt.Errorf("merging %q into %q: %w", o.Name, tool.Name, err)
		}
	}

	return tool, nil
}

// MergePlatform merges the provided platform into the tool.
// Only empty fields are set.
func (t *Tool) MergePlatform(platform detect.Platform) {
	t.Platform.Merge(platform)
}

// Marshal marshals the tool into a YAML byte slice.
func (t *Tool) Marshal() ([]byte, error) {
	data, err := yaml.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("marshalling tool %q: %w", t.Name, err)
	}

	return data, nil
}

// UnmarshalFrom unmarshals the tool from another tool.
func (t *Tool) UnmarshalFrom(other *Tool) error {
	data, err := other.Marshal()
	if err != nil {
		return fmt.Errorf("marshalling tool %q: %w", t.Name, err)
	}

	if err := unmarshal.Strict(data, t); err != nil {
		return fmt.Errorf("unmarshalling tool %q: %w", t.Name, err)
	}

	return nil
}

// UnmarshalInto unmarshals the tool into another tool.
func (t *Tool) UnmarshalInto(other *Tool) error {
	// Avoid pointers being copied. As such, we can always "copy" both the source and destination.
	copied, err := other.Copied()
	if err != nil {
		return fmt.Errorf("copying source in preparation for merge: %w", err)
	}

	data, err := t.Marshal()
	if err != nil {
		return fmt.Errorf("marshalling tool %q: %w", t.Name, err)
	}

	if err := unmarshal.Strict(data, copied); err != nil {
		return fmt.Errorf("unmarshalling tool %q: %w", t.Name, err)
	}

	*t = *copied

	return nil
}
