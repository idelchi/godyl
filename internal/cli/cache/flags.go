package cache

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/spf13/cobra"
)

// flags for the `cache` command.
func flags(cmd *cobra.Command) {
}

type Cache struct {
	common.Tracker `json:"-" mapstructure:"-"`
}
