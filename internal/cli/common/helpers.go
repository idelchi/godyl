package common

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/internal/config/root"
)

func FromEnvOrFile(keys ...string) (string, error) {
	for _, key := range keys {
		if value := viper.GetString(key); value != "" {
			return value, nil
		}
	}

	return "", nil
}

func ExitOnShow(show root.ShowFuncType, args ...string) bool {
	if show() != nil && len(args) == 0 {
		return true
	}

	return false
}

func SetSubcommandDefaults(cmd *cobra.Command, local any, show root.ShowFuncType) {
	var config Trackable

	if local != nil {
		local, ok := local.(Trackable)
		if !ok {
			panic("configuration may only be passed as Trackable type")
		}

		config = local
	}

	cmd.PersistentPreRunE = KCreateSubcommandPreRunE(cmd.Name(), config, show)
	SetSubcommandConfig(cmd, config)
}

func SetSubcommandConfig(cmd *cobra.Command, config Trackable) {
	ctx := context.WithValue(context.Background(), "config", config)
	cmd.SetContext(ctx)
}
