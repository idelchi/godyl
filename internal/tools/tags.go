package tools

import (
	"slices"

	"github.com/idelchi/godyl/pkg/unmarshal"
	"gopkg.in/yaml.v3"
)

type Tags []string

func (t *Tags) UnmarshalYAML(value *yaml.Node) error {
	result, err := unmarshal.UnmarshalSingleOrSlice[string](value, false)
	if err != nil {
		return err
	}
	*t = result
	return nil
}

func (t *Tags) Append(tags ...string) {
	for _, tag := range tags {
		if !slices.Contains(*t, tag) {
			*t = append(*t, tag)
		}
	}
}

func (t Tags) Has(tags Tags) bool {
	if len(tags) == 0 {
		return true
	}

	for _, tag := range tags {
		if slices.Contains(t, tag) {
			return true
		}
	}

	return false
}

func (t Tags) HasNot(tags Tags) bool {
	if len(tags) == 0 {
		return true
	}
	return !t.Has(tags)
}
