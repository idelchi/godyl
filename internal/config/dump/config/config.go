package config

import (
	"github.com/idelchi/godyl/internal/config/common"
)

type Config struct {
	common.Tracker `json:"-"   mapstructure:"-"`
	All            bool `json:"all" mapstructure:"all"`
}
