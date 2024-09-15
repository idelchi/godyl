// Package tools implements the tools dump subcommand for godyl.
// It displays information about the configured tools.
package tools

import (
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/tags"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/utils"
)

// Command encapsulates the tools dump command with its associated configuration.
type Command struct {
	// Command is the tools cobra.Command instance
	Command *cobra.Command
}

// NewToolsCommand creates a Command for displaying tools information.
func NewToolsCommand(cfg *config.Config, embedded *config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:               "tools",
		Short:             "Display tools information",
		Args:              cobra.NoArgs,
		PersistentPreRunE: common.KCreateSubcommandPreRunE("tools", &cfg.Dump.Tools, &cfg.Root.Show),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cfg.Root.Show {
				return nil
			}

			tags := iutils.SplitTags(cfg.Dump.Tools.Tags)

			c, err := getTools(embedded, cfg.Dump.Tools.Full, tags)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance for the tools dump subcommand.
func NewCommand(cfg *config.Config, files *config.Embedded) *cobra.Command {
	// Create the tools command
	cmd := NewToolsCommand(cfg, files)

	// Add tools-specific flags
	cmd.Flags()

	return cmd.Command
}

// getTools returns the tools configuration from embedded files.
func getTools(files *config.Embedded, rendered bool, tags tags.IncludeTags) (any, error) {
	tools := tools.Tools{}

	err := yaml.Unmarshal(files.Tools, &tools)
	if err != nil {
		return nil, err
	}

	var included []int

	for i, tool := range tools {
		tool.Tags.Append(tool.Name)

		if !tool.Tags.Include(tags.Include) || !tool.Tags.Exclude(tags.Exclude) {
			continue
		}

		included = append(included, i)
	}

	if !rendered {
		var tools []any

		err := yaml.Unmarshal(files.Tools, &tools)
		if err != nil {
			return nil, err
		}

		return utils.PickByIndices(tools, included), nil
	}

	return utils.PickByIndices(tools, included), nil
}
