package config

import "github.com/idelchi/godyl/internal/config/common"

type Config struct {
	common.Tracker `mapstructure:"-" yaml:"-"`
}
