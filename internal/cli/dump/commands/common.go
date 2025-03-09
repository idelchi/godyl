package dump

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/internal/config"
)

func commonPreRunE(cmd *cobra.Command, dump *config.Dump) error {
	parent := cmd.Parent()
	if err := viper.BindPFlags(parent.Flags()); err != nil {
		return fmt.Errorf("binding flags: %w", err)
	}

	if err := viper.Unmarshal(dump); err != nil {
		return fmt.Errorf("unmarshalling config: %w", err)
	}

	return config.Validate(dump)
}
