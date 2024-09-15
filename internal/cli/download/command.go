// Package download implements the download command for godyl.
// It provides functionality to download and extract tools from various sources.
package download

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/internal/tools/version"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/utils"
)

// Command encapsulates the download cobra command with its associated config and embedded files.
type Command struct {
	// Command is the download cobra.Command instance
	Command *cobra.Command
}

type Source struct {
	Type sources.Type `validate:"oneof=github gitlab url none go"`
}

type stubbedTool struct {
	Name     string
	URL      string
	Mode     mode.Mode
	Strategy strategy.Strategy
	Version  version.Version
	Source   Source
}

// NewDownloadCommand creates a Command for downloading and unpacking tools.
func NewDownloadCommand(cfg *config.Config, embedded *config.Embedded) *Command {
	// Create the download command
	cmd := &cobra.Command{
		Use:               "download [tool]",
		Aliases:           []string{"dl", "x"},
		Short:             "Download and extract tools",
		Args:              cobra.MinimumNArgs(1),
		PersistentPreRunE: common.KCreateSubcommandPreRunE("download", &cfg.Download, &cfg.Root.Show),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Root.Show {
				return nil
			}

			tools := tools.Tools{}

			for _, name := range args {
				tool := &tool.Tool{
					Mode:     mode.Extract,
					Strategy: strategy.Force,
					Version: version.Version{
						Version: cfg.Download.Version,
					},
				}

				if utils.IsURL(name) {
					tool.Name = file.New(name).Base()
					tool.URL = name
					tool.Source.Type = sources.URL
				} else {
					tool.Name = name
					tool.Source.Type = sources.GITHUB
				}

				tools.Append(tool)

			}

			// Generate a common configuration for the command
			cfg.SetCommon(cfg.Download.ToCommon())

			runner := common.NewHandler(cfg, embedded)
			if err := runner.SetupLogger(cfg.Root.LogLevel); err != nil {
				return fmt.Errorf("setting up logger: %w", err)
			}
			if err := runner.Resolve(&cfg.Root.Defaults, &tools); err != nil {
				return err
			}

			// Process tools
			proc := processor.New(tools, cfg, runner.Logger())
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
func NewCommand(cfg *config.Config, files *config.Embedded) *cobra.Command {
	// Create the download command
	cmd := NewDownloadCommand(cfg, files)

	// Add tool-specific flags
	cmd.Flags()

	return cmd.Command
}
