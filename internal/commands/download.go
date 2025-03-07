package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
)

// NewDownloadCommand creates the download command for downloading and unpacking tools.
func NewDownloadCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download [tool]",
		Short: "Download and unpack tools",
		Long:  "Download and unpack tools from GitHub, URLs, or Go projects",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the tool to download if provided as an argument
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
			log.Info("godyl download running with:")
			log.Info("*** ***")

			// Load defaults
			defaults := tools.Defaults{}
			if err := loadDefaults(&defaults, cfg.Defaults.Name(), nil, *cfg); err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			log.Info("platform:")
			log.Info(pretty.YAML(defaults.Platform))
			log.Info("*** ***")

			// If the tool is not a YAML file, treat it as a single tool
			var toolsList tools.Tools

			if isYAMLFile(cfg.Tools) {
				var loadErr error
				toolsList, loadErr = loadTools(cfg.Tools, log)
				if loadErr != nil {
					return fmt.Errorf("loading tools: %w", loadErr)
				}
			} else {
				// Create a single tool from the argument
				tool := tools.Tool{
					Name: cfg.Tools,
					Mode: "extract", // Default to extract mode for single tool
				}
				// Set the source type
				tool.Source.Type = cfg.Source
				toolsList = append(toolsList, tool)
				log.Info("downloading single tool: %s", cfg.Tools)
			}

			tags, withoutTags := splitTags(cfg.Tags)

			// Process tools
			processor := NewToolProcessor(toolsList, defaults, *cfg, log)
			if err := processor.Process(tags, withoutTags); err != nil {
				return fmt.Errorf("processing tools: %w", err)
			}

			return nil
		},
	}

	return cmd
}

// isYAMLFile checks if a file is a YAML file based on its extension.
func isYAMLFile(path string) bool {
	return path[len(path)-4:] == ".yml" || path[len(path)-5:] == ".yaml"
}
