package tools

import "slices"

type Tags []string

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
