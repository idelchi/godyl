package cache

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/cache/clean"
	show "github.com/idelchi/godyl/internal/cli/cache/dump"
	"github.com/idelchi/godyl/internal/cli/cache/path"
	"github.com/idelchi/godyl/internal/cli/cache/remove"
	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/cobraext"
)

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
		path.NewCommand(cmd.Config),
		remove.NewCommand(cmd.Config),
		show.NewCommand(cmd.Config),
		clean.NewCommand(cmd.Config, cmd.Files),
	)
}

func NewCacheCommand(cfg *config.Config, embedded *config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:               "cache [command]",
		Short:             "Interact with the cache",
		Long:              "Interact with the cache",
		Aliases:           []string{"c"},
		Args:              cobra.NoArgs,
		PersistentPreRunE: common.KCreateSubcommandPreRunE("cache", nil, &cfg.Root.Show),
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
	cmd := NewCacheCommand(cfg, files)

	// Add dump-specific flags
	cmd.Flags()

	// Add subcommands
	cmd.Subcommands()

	return cmd.Command
}
