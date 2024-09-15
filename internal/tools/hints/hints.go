package hints

import (
	"fmt"
	"slices"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Hints represents a collection of Hint objects used to evaluate asset matches.
type Hints unmarshal.SingleOrSliceType[Hint]

func (h *Hints) Has() bool {
	return h != nil && len(*h) > 0
}

func (h *Hints) UnmarshalYAML(node ast.Node) (err error) {
	*h, err = unmarshal.SingleOrSlice[Hint](node)
	if err != nil {
		return fmt.Errorf("unmarshaling hints: %w", err)
	}

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

// Reduced removes any Hint with an empty Pattern.
func (h *Hints) Reduced() *Hints {
	if h == nil {
		return nil
	}

	reduced := slices.DeleteFunc(*h, func(hint Hint) bool {
		return hint.Pattern == ""
	})

	result := Hints(reduced)

	return &result
}
