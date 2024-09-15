package dump

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/cli/dump/defaults"
	"github.com/idelchi/godyl/internal/cli/dump/env"
	"github.com/idelchi/godyl/internal/cli/dump/platform"
	"github.com/idelchi/godyl/internal/cli/dump/tools"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/cobraext"
)

// Command encapsulates the dump cobra command with its associated config and embedded files.
type Command struct {
	// Command is the dump cobra.Command instance
	Command *cobra.Command
	// Config contains application configuration
	Config *config.Config
	// Files contains the embedded configuration files and templates
	Files *config.Embedded
}

// Subcommands adds all subcommands to the dump command.
func (cmd *Command) Subcommands() {
	cmd.Command.AddCommand(
		defaults.NewCommand(cmd.Config, cmd.Files),
		env.NewCommand(cmd.Config),
		platform.NewCommand(cmd.Config),
		tools.NewCommand(cmd.Config, cmd.Files),
	)
}

// NewDumpCommand creates a Command for displaying configuration information.
func NewDumpCommand(cfg *config.Config, embedded *config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "Dump configuration information",
		Example: heredoc.Doc(`
			$ godyl dump --format json defaults
			$ godyl dump env
			$ godyl dump platform
			$ godyl dump tools
		`),
		PersistentPreRunE: common.KCreateSubcommandPreRunE("dump", &cfg.Dump, &cfg.Root.Show),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Root.Show && len(args) == 0 {
				return nil
			}

			return cobraext.UnknownSubcommandAction(cmd, args)
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
func NewCommand(cfg *config.Config, files *config.Embedded) *cobra.Command {
	// Create the dump command
	cmd := NewDumpCommand(cfg, files)

	// Add dump-specific flags
	cmd.Flags()

	// Add subcommands
	cmd.Subcommands()

	return cmd.Command
}
