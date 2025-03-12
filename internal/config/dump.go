package config

// Dump holds the configuration for the `dump` command.
type Dump struct {
	// Format for outputting the configuration
	Format string `validate:"oneof=json yaml"`

	// Tools configuration
	Tools Tools `mapstructure:"-" json:"-" yaml:"-"`
}

// Tools holds the configuration for the `dump tools` subcommand.
type Tools struct {
	// Rendered specifies whether to render the tools in full
	Rendered bool
}
