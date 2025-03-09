package config

// Dump contains configuration options for dumping data in different formats.
type Dump struct {
	// Format to dump the configuration in
	Format string `validate:"oneof=json yaml"`
}
