package tool

import (
	"fmt"
	"slices"

	"github.com/goccy/go-yaml"
	"github.com/idelchi/godyl/internal/detect"
	utils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// MergeFrom merges the current tool with the provided tools,
// the sequence being.
//
//	tool <-- others[0] <-- others[1] <-- others[2]...
func (t *Tool) MergeFrom(others ...*Tool) error {
	if len(others) == 0 {
		return fmt.Errorf("no defaults to merge")
	}

	for _, o := range others {
		if err := utils.Merge(t, o); err != nil {
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
		return fmt.Errorf("no defaults to merge")
	}

	// Create an empty tool first to not override the first one.
	copied := &Tool{}
	// Append the current tool into the list to have it at the end.
	others = append(others, t)

	for _, o := range others {
		if err := utils.Merge(copied, o); err != nil {
			return fmt.Errorf("merging %q into %q: %w", o.Name, copied.Name, err)
		}
	}

	*t = *copied

	return nil
}

// MergeRightToLeft merges the provided tools into a new tool,
// right to left.
func MergeRightToLeft(others ...*Tool) (*Tool, error) {
	if len(others) == 0 {
		return nil, fmt.Errorf("no defaults to merge")
	}

	// Create an empty tool
	tool := &Tool{}

	for _, o := range others {
		if err := utils.Merge(tool, o); err != nil {
			return nil, fmt.Errorf("merging %q into %q: %w", o.Name, tool.Name, err)
		}
	}

	return tool, nil
}

// MergeLeftToRight merges the provided tools into a new tool,
// left to right.
func MergeLeftToRight(others ...*Tool) (*Tool, error) {
	if len(others) == 0 {
		return nil, fmt.Errorf("no defaults to merge")
	}

	slices.Reverse(others)

	// Create an empty tool
	tool := &Tool{}

	for _, o := range others {
		if err := utils.Merge(tool, o); err != nil {
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
