package match

// Hints represents a collection of Hint objects used to evaluate asset matches.
type Hints []Hint

// Append appends a set of hints to the current collection.
func (h *Hints) Append(hints Hints) {
	*h = append(*h, hints...)
}

// Add appends a set of hints to the current collection.
func (h *Hints) Add(hints ...Hint) {
	h.Append(hints)
}
