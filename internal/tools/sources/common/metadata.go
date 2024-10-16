package common

// Metadata represents a map of string key-value pairs used to store metadata information.
type Metadata map[string]string

// Get retrieves the value of a specific attribute from the metadata.
// If the attribute does not exist, it returns an empty string.
func (m *Metadata) Get(attribute string) string {
	if m == nil {
		return ""
	}

	if value, ok := (*m)[attribute]; ok {
		return value
	}

	return ""
}

// Set assigns a value to a specific attribute in the metadata.
// If the underlying map is nil, it initializes the map before setting the value.
func (m *Metadata) Set(attribute, value string) {
	if *m == nil { // Check if the underlying map is nil
		*m = make(Metadata) // Initialize the map
	}

	(*m)[attribute] = value
}
