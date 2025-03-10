package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Bind connects cobra flags to viper and unmarshals the configuration into the provided struct.
// It sets up environment variable handling with the given prefix and handles flag binding.
func Bind(cmd *cobra.Command, cfg any, prefix ...string) error {
	viper.SetEnvPrefix(strings.Join(prefix, "_"))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("binding flags: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshalling config: %w", err)
	}

	return nil
}
