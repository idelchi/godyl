package config

type Cache struct {
	Delete bool

	// Viper instance
	viperable `mapstructure:"-" yaml:"-" json:"-"`
}
