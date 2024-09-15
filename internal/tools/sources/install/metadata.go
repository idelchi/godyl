package install

// Metadata stores arbitrary key-value pairs as string attributes.
// It provides a safe interface for storing and retrieving metadata information.
type Metadata map[string]string

// Get retrieves a metadata attribute's value by its key.
// Returns an empty string if the metadata is nil or the key doesn't exist.
func (m *Metadata) Get(attribute string) string {
	if m == nil {
		return ""
	}

	if value, ok := (*m)[attribute]; ok {
		return value
	}

	return ""
}

// Set stores a key-value pair in the metadata.
// Initializes the underlying map if it is nil before setting the value.
func (m *Metadata) Set(attribute, value string) {
	if *m == nil { // Check if the underlying map is nil
		*m = make(Metadata) // Initialize the map
	}

	(*m)[attribute] = value
}
