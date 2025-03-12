// Package env implements the env dump subcommand for godyl.
// It displays information about the environment variables.
package env

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/env"
)

// Command encapsulates the env dump command with its associated configuration.
type Command struct {
	// Command is the env cobra.Command instance
	Command *cobra.Command
	// Config contains application configuration
	Config *config.Config
	// Files contains the embedded configuration files and templates
	Files config.Embedded
}

// Flags adds env-specific flags to the command.
func (cmd *Command) Flags() {
	// No specific flags for this command
}

// NewEnvCommand creates a Command for displaying environment information.
func NewEnvCommand(cfg *config.Config, files config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Display environment information",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.Bind(cmd.Parent(), &cfg.Dump, cmd.Root().Name(), cmd.Parent().Name())
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := getEnv()
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return &Command{
		Command: cmd,
		Config:  cfg,
		Files:   files,
	}
}

// NewCommand creates a cobra.Command instance for the env dump subcommand.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the env command
	cmd := NewEnvCommand(cfg, files)

	// Add env-specific flags
	cmd.Flags()

	return cmd.Command
}

// getEnv retrieves the current environment variables.
func getEnv() (env.Env, error) {
	return env.FromEnv(), nil
}
