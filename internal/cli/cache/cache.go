package cache

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/cache/path"
	"github.com/idelchi/godyl/internal/cli/cache/remove"
	"github.com/idelchi/godyl/internal/cli/cache/show"
	"github.com/idelchi/godyl/internal/cli/cache/sync"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
)

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
}

// Subcommands adds all subcommands to the dump command.
func (cmd *Command) Subcommands() {
	cmd.Command.AddCommand(
		path.NewCommand(cmd.Config),
		remove.NewCommand(cmd.Config),
		show.NewCommand(cmd.Config),
		sync.NewCommand(cmd.Config, cmd.Files),
	)
}

// NewDumpCommand creates a Command for displaying configuration information.
func NewDumpCommand(cfg *config.Config, embedded config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:     "cache [command]",
		Short:   "Interact with the cache",
		Long:    "Interact with the cache",
		Aliases: []string{"c"},
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, nil)
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
