// Package env implements the env dump subcommand for godyl.
// It displays information about the environment variables.
package env

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/env"
)

// Command encapsulates the env dump command with its associated configuration.
type Command struct {
	// Command is the env cobra.Command instance
	Command *cobra.Command
}

// NewEnvCommand creates a Command for displaying environment information.
func NewEnvCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:               "env",
		Short:             "Display environment information",
		Args:              cobra.NoArgs,
		PersistentPreRunE: common.KCreateSubcommandPreRunE("env", &cfg.Dump.Tools, &cfg.Root.Show),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cfg.Root.Show {
				return nil
			}

			c, err := getEnv()
			if err != nil {
				return err
			}

			format := "env"
			if cfg.Dump.IsSet("format") {
				format = cfg.Dump.Format
			}

			iutils.Print(format, c)

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance for the env dump subcommand.
func NewCommand(cfg *config.Config) *cobra.Command {
	// Create the env command
	cmd := NewEnvCommand(cfg)

	// Add env-specific flags
	cmd.Flags()

	return cmd.Command
}

// getEnv retrieves the current environment variables.
func getEnv() (env.Env, error) {
	return env.FromEnv(), nil
}
