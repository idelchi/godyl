package config

type Cache struct {
	Dump CDump `yaml:"-"`

	trackable `json:"-" mapstructure:"-" yaml:"-"`
}

type CDump struct {
	Format string `mapstructure:"format" validate:"oneof=json yaml"`

	trackable `json:"-" mapstructure:"-" yaml:"-"`
}
