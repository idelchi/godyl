package tool

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	"github.com/idelchi/godyl/internal/core/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/utils"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
)

func NewInstallCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install [tools.yml]",
		Aliases: []string{"i", "get"},
		Short:   "Install tools from a YAML file",
		Long:    "Install tools as specified in a YAML configuration file",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := commonPreRunE(cmd, &cfg.Tool); err != nil {
				return fmt.Errorf("common pre-run: %w", err)
			}

			return config.Validate(cfg.Tool)
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
			if err := defaults.LoadDefaults(&toolDefaults, cfg.Root.Defaults.Name(), files.Defaults, *cfg); err != nil {
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

	// Add tool-specific flags
	addToolFlags(cmd)

	return cmd
}
