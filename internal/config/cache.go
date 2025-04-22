package config

type Cache struct {
	viperable `json:"-" mapstructure:"-" yaml:"-"`
	Delete    bool
	Sync      bool
}
