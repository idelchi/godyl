// Package tools implements the tools dump subcommand for godyl.
// It displays information about the configured tools.
package tools

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/utils"
	iutils "github.com/idelchi/godyl/internal/utils"
)

// Command encapsulates the tools dump command with its associated configuration.
type Command struct {
	// Command is the tools cobra.Command instance
	Command *cobra.Command
	// Config contains application configuration
	Config *config.Config
	// Files contains the embedded configuration files and templates
	Files config.Embedded
}

// Flags adds tools-specific flags to the command.
func (cmd *Command) Flags() {
	cmd.Command.Flags().BoolP("full", "f", false, "Show the tools in full syntax")
}

// NewToolsCommand creates a Command for displaying tools information.
func NewToolsCommand(cfg *config.Config, files config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Display tools information",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Dump.Tools)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := getTools(files, cfg.Dump.Tools.Full)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return &Command{
		Command: cmd,
		Config:  cfg,
		Files:   files,
	}
}

// NewCommand creates a cobra.Command instance for the tools dump subcommand.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the tools command
	cmd := NewToolsCommand(cfg, files)

	// Add tools-specific flags
	cmd.Flags()

	return cmd.Command
}

// getTools returns the tools configuration from embedded files.
func getTools(files config.Embedded, rendered bool) (any, error) {
	if !rendered {
		return utils.PrintYAMLBytes(files.Tools), nil
	}

	tools := &tools.Tools{}

	err := yaml.Unmarshal(files.Tools, tools)
	if err != nil {
		return nil, err
	}

	return tools, nil
}
