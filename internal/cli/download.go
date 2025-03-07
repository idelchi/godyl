package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	"github.com/idelchi/godyl/internal/core/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/godyl/pkg/utils"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// NewDownloadCommand creates the download command for downloading and unpacking tools.
func NewDownloadCommand(cfg *config.Config, emb rootEmbedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download [tool]",
		Aliases: []string{"dl", "unpack"},
		Short:   "Download and unpack tools",
		Long:    "Download and unpack tools from GitHub, URLs, or Go projects",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: func(_ *cobra.Command, args []string) error {
			return cobraext.Validate(cfg, &cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			lvl, err := logger.LevelString(cfg.Log)
			if err != nil {
				return fmt.Errorf("parsing log level: %w", err)
			}

			log := logger.New(lvl)
			log.Info("*** ***")
			log.Info("godyl download running with:")
			log.Info("*** ***")

			// Load defaults
			toolDefaults := tools.Defaults{}
			if err := defaults.LoadDefaults(&toolDefaults, cfg.Defaults.Name(), emb.defaults, *cfg); err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			log.Info("platform:")
			log.Info(pretty.YAML(toolDefaults.Platform))
			log.Info("*** ***")

			toolsList := []tools.Tool{}
			for _, name := range args {
				tool := tools.Tool{
					Name: name,
					Mode: tools.Extract,
				}
				if utils.IsURL(name) {
					tool.Name = filepath.Base(name)
					tool.Path = name
					tool.Source.Type = sources.DIRECT
				}

				toolsList = append(toolsList, tool)
			}

			// Process tools
			proc := processor.New(toolsList, toolDefaults, *cfg, log)
			if err := proc.Process(nil, nil); err != nil {
				return fmt.Errorf("processing tools: %w", err)
			}

			return nil
		},
	}

	// Add tool-specific flags
	addToolFlags(cmd)

	return cmd
}
