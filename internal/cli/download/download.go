// Package download implements the download command for godyl.
// It provides functionality to download and extract tools from various sources.
package download

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/defaults"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/tools/tags"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/utils"
)

// Command encapsulates the download cobra command with its associated config and embedded files.
type Command struct {
	// Command is the download cobra.Command instance
	Command *cobra.Command
}

// Flags adds download-specific flags to the command.
func (cmd *Command) Flags() {
	flags.Tool(cmd.Command)
}

// NewDownloadCommand creates a Command for downloading and unpacking tools.
func NewDownloadCommand(cfg *config.Config, embedded config.Embedded) *Command {
	// Create the download command
	cmd := &cobra.Command{
		Use:     "download [tool]",
		Aliases: []string{"dl", "x"},
		Short:   "Download and unpack tools",
		Long:    "Download and unpack tools from GitHub, URLs, or Go projects",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Tool, cmd.Root().Name(), "tool")
		},
		RunE: func(_ *cobra.Command, args []string) error {
			lvl, err := logger.LevelString(cfg.Root.Log)
			if err != nil {
				return fmt.Errorf("parsing log level: %w", err)
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

			tmp, err := tmp.GodylCreateRandomDir()
			if err != nil {
				return fmt.Errorf("creating temporary directory: %w", err)
			}

			type Download struct {
				Name string
				Path string
			}

			var downloads []Download

			for _, name := range args {
				if utils.IsURL(name) {
					downloads = append(downloads, Download{
						Name: file.New(name).Base(),
						Path: name,
					})
				} else {
					downloads = append(downloads, Download{
						Name: name,
					})
				}
			}

			toolsFile := tmp.WithFile("godyl.yaml")
			defer tmp.Remove()

			if err := toolsFile.WriteYAML(downloads); err != nil {
				return fmt.Errorf("writing YAML: %w", err)
			}

			toolsList, err := iutils.LoadTools(toolsFile, defaults, cfg.Root.Default)
			if err != nil {
				return fmt.Errorf("loading tools: %w", err)
			}

			for i, tool := range toolsList {
				toolsList[i].Mode = mode.Extract
				toolsList[i].Strategy = strategy.Force
				toolsList[i].Version.Version = cfg.Tool.Version
				if tool.URL != "" {
					toolsList[i].Source.Type = sources.URL
				}
			}

			// Process tools
			proc := processor.New(toolsList, cfg, log)
			if err := proc.Process(tags.IncludeTags{}, false); err != nil {
				return fmt.Errorf("processing tools: %w", err)
			}

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance containing the download command.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the download command
	cmd := NewDownloadCommand(cfg, files)

	// Add tool-specific flags
	cmd.Flags()

	return cmd.Command
}
