package match

import (
	"github.com/idelchi/godyl/pkg/unmarshal"
	"gopkg.in/yaml.v3"
)

// Hints represents a collection of Hint objects used to evaluate asset matches.
type Hints unmarshal.SingleOrSliceType[Hint]

// UnmarshalYAML allows unmarshaling the YAML node as either a single Hint or a slice of Hints,
// while appending the original slice.
func (h *Hints) UnmarshalYAML(value *yaml.Node) error {
	result, err := unmarshal.SingleOrSlice[Hint](value, false)
	if err != nil {
		return err
	}

	h.Append(result)

	return nil
}

// Append appends a set of hints to the current collection.
func (h *Hints) Append(hints Hints) {
	*h = append(*h, hints...)
}

// Add appends a set of hints to the current collection.
func (h *Hints) Add(hints ...Hint) {
	h.Append(hints)
}

func (h *Hints) Parse() error {
	for i, hint := range *h {
		if err := hint.Parse(); err != nil {
			return err
		}

		// Update the hint in the collection with the parsed weight.
		(*h)[i] = hint
	}

	return nil
}
