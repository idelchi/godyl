package cli

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// handleExitFlags handles flags that cause the application to exit.
func handleExitFlags(cmd *cobra.Command, version string, cfg *config.Config, defaultEmbedded []byte) error {
	// Check if the version flag was provided
	if cfg.Version {
		cmd.Println(version)
		return cobraext.ErrExitGracefully
	}

	// Check if the help flag was provided
	if cfg.Help {
		cmd.Help()
		return cobraext.ErrExitGracefully
	}

	if cfg.Dump.Config {
		pretty.PrintYAMLMasked(*cfg)
		return cobraext.ErrExitGracefully
	}

	if cfg.Dump.Env {
		pretty.PrintYAMLMasked(env.FromEnv())
		return cobraext.ErrExitGracefully
	}

	return nil
}
