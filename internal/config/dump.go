package config

// Dump holds the configuration for the `dump` command.
type Dump struct {
	Tools  Tools  `yaml:"-"`
	Format string `validate:"oneof=json yaml"`

	trackable `json:"-" mapstructure:"-" yaml:"-"`
}

// Tools holds the configuration for the `dump tools` subcommand.
type Tools struct {
	Full bool
	Tags []string

	trackable `json:"-" mapstructure:"-" yaml:"-"`
}
