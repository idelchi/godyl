// Package install implements the install command for godyl.
// It provides functionality to install tools from a YAML configuration file.
package install

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	"github.com/idelchi/godyl/internal/core/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/utils"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/files"
	"github.com/idelchi/godyl/pkg/pretty"
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
		Aliases: []string{"i", "get"},
		Short:   "Install tools from one of more YAML files",
		Long:    "Install tools as specified in the YAML file(s).",
		Args:    cobra.ArbitraryArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Tool, cmd.Root().Name(), "tool")
		},
		RunE: func(_ *cobra.Command, args []string) error {
			lvl, err := logger.LevelString(cfg.Root.Log)
			if err != nil {
				return fmt.Errorf("parsing log level: %w", err)
			}

			// Set the tools file if provided as an argument
			if len(args) > 0 {
				cfg.Tool.Tools = files.New("", args...)
			} else {
				cfg.Tool.Tools = files.New(".", "tools.yml")
			}

			log := logger.New(lvl)

			// Load defaults
			defaults, err := defaults.Load(cfg.Root.Defaults, embedded, *cfg)
			if err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			if cfg.Tool.Show {
				iutils.Print("yaml", cfg.Root, cfg.Tool)

				return nil
			}

			toolsList := tools.Tools{}

			// Load tools
			for _, file := range cfg.Tool.Tools {
				tools, err := utils.LoadTools(file, defaults)
				if err != nil {
					return fmt.Errorf("loading tools: %w", err)
				}

				log.Info("loaded %d tools from %q", len(tools), file)

				toolsList = append(toolsList, tools...)
			}

			tags, withoutTags := utils.SplitTags(cfg.Tool.Tags)

			pretty.PrintYAML(toolsList)

			os.Exit(0)

			proc := processor.New(toolsList, defaults, *cfg, log)
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
