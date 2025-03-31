package config

// Dump holds the configuration for the `dump` command.
type Dump struct {
	// Format for outputting the configuration
	Format string `validate:"oneof=json yaml"`

	// Tools configuration
	Tools Tools `mapstructure:"-"`

	// Viper instance
	viperable `mapstructure:"-" yaml:"-" json:"-"`
}

// Tools holds the configuration for the `dump tools` subcommand.
type Tools struct {
	// Full specifies whether to output the tools in full syntax
	Full bool

	// Viper instance
	viperable `mapstructure:"-" yaml:"-" json:"-"`
}
