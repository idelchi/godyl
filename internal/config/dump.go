package config

// Dump holds the configuration for the `dump` command.
type Dump struct {
	// Format for outputting the configuration
	Format string `validate:"oneof=json yaml"`

	// Tools configuration
	Tools Tools `mapstructure:"-"`

	// Cache configuration
	Cache Cache `mapstructure:"-"`

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

// Cache holds the configuration for the `dump cache` subcommand.
type Cache struct {
	// File specifies whether to output the path to the cache file
	File bool

	// Viper instance
	viperable `mapstructure:"-" yaml:"-" json:"-"`
}
