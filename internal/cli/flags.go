package cli

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// FlagHandler is a function that handles a specific flag and returns an error if the application should exit.
type FlagHandler func(cmd *cobra.Command, cfg *config.Config) error

// handleExitFlags handles flags that cause the application to exit.
func handleExitFlags(cmd *cobra.Command, version string, cfg *config.Config, defaultEmbedded []byte) error {
	// Define handlers for each exit flag
	handlers := []FlagHandler{
		// Version flag handler
		func(cmd *cobra.Command, cfg *config.Config) error {
			if cfg.Version {
				cmd.Println(version)
				return cobraext.ErrExitGracefully
			}
			return nil
		},
		// Help flag handler
		func(cmd *cobra.Command, cfg *config.Config) error {
			if cfg.Help {
				cmd.Help()
				return cobraext.ErrExitGracefully
			}
			return nil
		},
		// Config dump handler
		func(_ *cobra.Command, cfg *config.Config) error {
			if cfg.Dump.Config {
				pretty.PrintYAMLMasked(*cfg)
				return cobraext.ErrExitGracefully
			}
			return nil
		},
		// Environment dump handler
		func(_ *cobra.Command, cfg *config.Config) error {
			if cfg.Dump.Env {
				pretty.PrintYAMLMasked(env.FromEnv())
				return cobraext.ErrExitGracefully
			}
			return nil
		},
	}

	// Execute each handler until one returns an error
	for _, handler := range handlers {
		if err := handler(cmd, cfg); err != nil {
			return err
		}
	}

	return nil
}
