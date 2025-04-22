package config

type Cache struct {
	viperable `json:"-" mapstructure:"-" yaml:"-"`

	Format string `validate:"oneof=json yaml"`
}
