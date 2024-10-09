package sources

type Metadata map[string]string

func (m *Metadata) Get(attribute string) string {
	if m == nil {
		return ""
	}

	if value, ok := (*m)[attribute]; ok {
		return value
	}

	return ""
}

func (m *Metadata) Set(attribute, value string) {
	if *m == nil { // Check if the underlying map is nil
		*m = make(Metadata) // Initialize the map
	}

	(*m)[attribute] = value
}
