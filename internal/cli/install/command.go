// Package install implements the install command for godyl.
// It provides functionality to install tools from a YAML configuration file.
package install

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/tools"
	utils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Command encapsulates the install cobra command with its associated config and embedded files.
type Command struct {
	// Command is the install cobra.Command instance
	Command *cobra.Command
}

// NewInstallCommand creates a Command for installing tools from a YAML file.
func NewInstallCommand(cfg *config.Config, embedded *config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:               "install [tools.yml...|-]",
		Aliases:           []string{"i"},
		Short:             "Install tools from one of more YAML files",
		Long:              "Install tools as specified in the YAML file(s). Use '-' to read from stdin.",
		Args:              cobra.ArbitraryArgs,
		PersistentPreRunE: common.KCreateSubcommandPreRunE("install", &cfg.Install, &cfg.Root.Show),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Root.Show {
				return nil
			}

			// Load the tools from the source as []byte
			data, err := utils.ReadFromArgs("tools.yml", args...)
			if err != nil {
				return fmt.Errorf("reading arguments %v: %w", args, err)
			}

			// The tools can now be unmarshalled into a tools.Tools instance
			var tools tools.Tools
			if err := unmarshal.Strict(data, &tools); err != nil {
				return fmt.Errorf("unmarshalling tools: %w", err)
			}

			// Generate a common configuration for the command
			cfg.SetCommon(cfg.Install.ToCommon())

			runner := common.NewHandler(cfg, embedded)
			if err := runner.SetupLogger(cfg.Root.LogLevel); err != nil {
				return fmt.Errorf("setting up logger: %w", err)
			}
			if err := runner.Resolve(&cfg.Root.Defaults, &tools); err != nil {
				return err
			}

			// At this point, all tools have been resolved and can be processed by the processor
			proc := processor.New(tools, cfg, runner.Logger())
			if err := proc.Process(utils.SplitTags(cfg.Install.Tags), false); err != nil {
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
func NewCommand(cfg *config.Config, files *config.Embedded) *cobra.Command {
	// Create the install command
	cmd := NewInstallCommand(cfg, files)

	// Add tool-specific flags
	cmd.Flags()

	cmd.Command.Flags().MarkHidden("version")

	return cmd.Command
}
