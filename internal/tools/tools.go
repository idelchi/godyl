package tools

// Tools represents a collection of Tool configurations.
type Tools []*Tool

func NewTools(d Defaults, length int) Tools {
	collection := make(Tools, length)

	for i := range collection {
		collection[i] = NewTool(d)
	}

	return collection
}
