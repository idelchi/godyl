package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// NewInstallCommand creates the install command for installing tools from a YAML file.
func NewInstallCommand(cfg *config.Config, emb rootEmbedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [tools.yml]",
		Short: "Install tools from a YAML file",
		Long:  "Install tools as specified in a YAML configuration file",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(_ *cobra.Command, args []string) error {
			return cobraext.Validate(cfg, &cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
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
			defaults := tools.Defaults{}
			if err := loadDefaults(&defaults, cfg.Defaults.Name(), emb.defaults, *cfg); err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			log.Info("platform:")
			log.Info(pretty.YAML(defaults.Platform))
			log.Info("*** ***")

			// Load tools
			toolsList, err := loadTools(cfg.Tools, log)
			if err != nil {
				return fmt.Errorf("loading tools: %w", err)
			}

			tags, withoutTags := splitTags(cfg.Tags)

			processor := NewToolProcessor(toolsList, defaults, *cfg, log)
			if err := processor.Process(tags, withoutTags); err != nil {
				return fmt.Errorf("processing tools: %w", err)
			}

			return nil
		},
	}

	// Add tool-specific flags
	addToolFlags(cmd)

	return cmd
}
