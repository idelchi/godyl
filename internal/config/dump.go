package config

// Dump holds the configuration for the `dump` command.
type Dump struct {
	// Format for outputting the configuration
	Format string `validate:"oneof=json yaml"`
}
