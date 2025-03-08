package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	"github.com/idelchi/godyl/internal/core/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// NewInstallCommand creates the install command for installing tools from a YAML file.
func NewInstallCommand(cfg *config.Config, files Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install [tools.yml]",
		Aliases: []string{"i", "get"},
		Short:   "Install tools from a YAML file",
		Long:    "Install tools as specified in a YAML configuration file",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return cobraext.Validate(cfg, &cfg)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			// Set the tools file if provided as an argument
			if len(args) > 0 {
				cfg.Tools = args[0]
			} else {
				cfg.Tools = "tools.yml"
			}

			lvl, err := logger.LevelString(cfg.Log)
			if err != nil {
				return fmt.Errorf("parsing log level: %w", err)
			}

			log := logger.New(lvl)
			log.Info("*** ***")
			log.Info("godyl install running with:")
			log.Info("*** ***")

			// Load defaults
			toolDefaults := tools.Defaults{}
			if err := defaults.LoadDefaults(&toolDefaults, cfg.Defaults.Name(), files.Defaults, *cfg); err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			log.Info("platform:")
			log.Info("%s", pretty.YAML(toolDefaults.Platform))
			log.Info("*** ***")

			// Load tools
			toolsList, err := utils.LoadTools(cfg.Tools, log)
			if err != nil {
				return fmt.Errorf("loading tools: %w", err)
			}

			tags, withoutTags := utils.SplitTags(cfg.Tags)

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
