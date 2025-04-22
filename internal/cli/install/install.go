// Package install implements the install command for godyl.
// It provides functionality to install tools from a YAML configuration file.
package install

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/utils"
	iutils "github.com/idelchi/godyl/internal/utils"
)

// Command encapsulates the install cobra command with its associated config and embedded files.
type Command struct {
	// Command is the install cobra.Command instance
	Command *cobra.Command
}

// Flags adds install-specific flags to the command.
func (cmd *Command) Flags() {
	flags.Tool(cmd.Command)
}

// NewInstallCommand creates a Command for installing tools from a YAML file.
func NewInstallCommand(cfg *config.Config, embedded config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:     "install [tools.yml]...",
		Aliases: []string{"i"},
		Short:   "Install tools from one of more YAML files",
		Long:    "Install tools as specified in the YAML file(s).",
		Args:    cobra.ArbitraryArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Tool, cmd.Root().Name(), "tool")
		},
		RunE: func(_ *cobra.Command, args []string) error {
			if cfg.Tool.Show {
				iutils.Print("yaml", cfg.Root, cfg.Tool)

				return nil
			}

			toolsList, log, err := common.Common(cfg, embedded, args)
			if err != nil {
				return err
			}

			proc := processor.New(toolsList, cfg, log)
			if err := proc.Process(utils.SplitTags(cfg.Tool.Tags), false); err != nil {
				return fmt.Errorf("processing tools: %w", err)
			}

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance containing the install command.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the install command
	cmd := NewInstallCommand(cfg, files)

	// Add tool-specific flags
	cmd.Flags()

	cmd.Command.Flags().MarkHidden("version")

	return cmd.Command
}
