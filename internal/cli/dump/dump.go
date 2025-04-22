package dump

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/dump/configuration"
	"github.com/idelchi/godyl/internal/cli/dump/defaults"
	"github.com/idelchi/godyl/internal/cli/dump/env"
	"github.com/idelchi/godyl/internal/cli/dump/platform"
	"github.com/idelchi/godyl/internal/cli/dump/tools"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
)

// Command encapsulates the dump cobra command with its associated config and embedded files.
type Command struct {
	// Command is the dump cobra.Command instance
	Command *cobra.Command
	// Config contains application configuration
	Config *config.Config
	// Files contains the embedded configuration files and templates
	Files config.Embedded
}

// Flags adds dump-specific flags to the command.
func (cmd *Command) Flags() {
	cmd.Command.Flags().StringP("format", "f", "yaml", "Output format (json or yaml)")
}

// Subcommands adds all subcommands to the dump command.
func (cmd *Command) Subcommands() {
	cmd.Command.AddCommand(
		configuration.NewCommand(cmd.Config),
		defaults.NewCommand(cmd.Config, cmd.Files),
		env.NewCommand(cmd.Config),
		platform.NewCommand(cmd.Config),
		tools.NewCommand(cmd.Config, cmd.Files),
	)
}

// NewDumpCommand creates a Command for displaying configuration information.
func NewDumpCommand(cfg *config.Config, embedded config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:     "dump [command]",
		Aliases: []string{"d"},
		Short:   "Dump configuration information",
		Long:    "Display various configuration settings and information about the environment",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Dump)
		},
	}

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	return &Command{
		Command: cmd,
		Config:  cfg,
		Files:   embedded,
	}
}

// NewCommand creates a cobra.Command instance containing the dump command and its subcommands.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the dump command
	cmd := NewDumpCommand(cfg, files)

	// Add dump-specific flags
	cmd.Flags()

	// Add subcommands
	cmd.Subcommands()

	return cmd.Command
}
