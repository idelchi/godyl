package config

// Dump holds the configuration for the `dump` command.
type Dump struct {
	Tools     Tools `mapstructure:"-"`
	viperable `json:"-" mapstructure:"-" yaml:"-"`
	Format    string `validate:"oneof=json yaml"`
}

// Tools holds the configuration for the `dump tools` subcommand.
type Tools struct {
	viperable `json:"-" mapstructure:"-" yaml:"-"`
	Full      bool
	Tags      []string
}
