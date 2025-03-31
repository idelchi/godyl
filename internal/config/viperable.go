package config

import (
	"github.com/spf13/viper"
)

// Viperable is a generic struct that holds a viper instance.
type viperable struct {
	// Viper instance
	viper *viper.Viper `mapstructure:"-" yaml:"-" json:"-"`
}

// SetViper sets the viper instance for the configuration.
func (v *viperable) SetViper(viper *viper.Viper) {
	v.viper = viper
}

// IsSet checks if a flag is set in viper,
// to avoid using it's default values unless explicitly passed.
func (v *viperable) IsSet(flag string) bool {
	if v.viper == nil {
		return false
	}

	return v.viper.IsSet(flag)
}
