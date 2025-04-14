package tools

// Tools represents a collection of Tool configurations.
type Tools []*Tool

func NewTools(d Defaults, length int) (collection Tools, err error) {
	collection = make(Tools, length)

	for i := range collection {
		if collection[i], err = NewTool(d); err != nil {
			return nil, err
		}
	}

	return collection, nil
}
