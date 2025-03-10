// Package install implements the install command for godyl.
// It provides functionality to install tools from a YAML configuration file.
package install

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	"github.com/idelchi/godyl/internal/core/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/utils"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/godyl/pkg/validate"
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
func NewInstallCommand(cfg *config.Config, files config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:     "install [tools.yml]",
		Aliases: []string{"i", "get"},
		Short:   "Install tools from a YAML file",
		Long:    "Install tools as specified in a YAML configuration file",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := flags.Bind(cmd, cmd.Root().Name(), &cfg.Tool); err != nil {
				return fmt.Errorf("common pre-run: %w", err)
			}

			return validate.Validate(cfg.Tool)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			if cfg.Root.Show {
				iutils.Print("yaml", cfg.Root, cfg.Tool)

				return nil
			}

			lvl, err := logger.LevelString(cfg.Root.Log)
			if err != nil {
				return fmt.Errorf("parsing log level: %w", err)
			}

			// Set the tools file if provided as an argument
			if len(args) > 0 {
				cfg.Tool.Tools = file.File(args[0])
			} else {
				cfg.Tool.Tools = "tools.yml"
			}

			log := logger.New(lvl)
			log.Info("*** ***")
			log.Info("godyl install running with:")
			log.Info("*** ***")

			// Load defaults
			toolDefaults := tools.Defaults{}
			if err := defaults.LoadDefaults(&toolDefaults, cfg.Root.Defaults.Name(), files, *cfg); err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			log.Info("platform:")
			log.Info("%s", pretty.YAML(toolDefaults.Platform))
			log.Info("*** ***")

			// Load tools
			toolsList, err := utils.LoadTools(cfg.Tool.Tools, log)
			if err != nil {
				return fmt.Errorf("loading tools: %w", err)
			}

			tags, withoutTags := utils.SplitTags(cfg.Tool.Tags)

			proc := processor.New(toolsList, toolDefaults, *cfg, log)
			if err := proc.Process(tags, withoutTags); err != nil {
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
