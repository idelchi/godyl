package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/utils"
)

// Command encapsulates the status cobra command with its associated config and embedded files.
type Command struct {
	// Command is the status cobra.Command instance
	Command *cobra.Command
}

// Flags adds status-specific flags to the command.
func (cmd *Command) Flags() {
	flags.Status(cmd.Command)
}

// NewStatusCommand creates a Command for statusing tools from a YAML file.
func NewStatusCommand(cfg *config.Config, embedded config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:     "status [tools.yml]...",
		Aliases: []string{"diff", "s"},
		Short:   "Status of installed tools as specified in the YAML file(s).",
		Long:    "Status of installed tools as specified in the YAML file(s).",
		Args:    cobra.ArbitraryArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Status)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			toolsList, log, err := common.Common(cfg, embedded, args)
			if err != nil {
				return err
			}

			for i := range toolsList {
				toolsList[i].Strategy = strategy.Sync
			}

			proc := processor.New(toolsList, cfg, log)
			if err := proc.Process(utils.SplitTags(cfg.Tool.Tags), true); err != nil {
				return fmt.Errorf("processing tools: %w", err)
			}

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance containing the status command.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the status command
	cmd := NewStatusCommand(cfg, files)

	// Add tool-specific flags
	cmd.Flags()

	cmd.Command.Flags().MarkHidden("version")

	return cmd.Command
}
