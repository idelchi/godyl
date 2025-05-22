package config

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/spf13/cobra"
)

// flags for the `config` command.
func flags(cmd *cobra.Command) {
}

type Config struct {
	common.Tracker `json:"-" mapstructure:"-"`
}
