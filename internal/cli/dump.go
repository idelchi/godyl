package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// Print displays the configuration in the specified format.
func Print(cfg any, format string) {
	switch format {
	case "json":
		pretty.PrintJSONMasked(cfg)
	case "yaml":
		pretty.PrintYAMLMasked(cfg)
	default:
		fmt.Printf("unsupported output format: %s\n", format)
	}
}

// NewDumpCommand creates the show command for displaying various configurations.
func NewDumpCommand(cfg *config.Config, files Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dump [flags] [config|defaults|env|platform|tools] [flags]",
		Aliases: []string{"show"},
		Short:   "Dump configuration information",
		Long:    "Display various configuration settings and information about the environment",
		Args:    cobra.MaximumNArgs(1),

		PreRunE: func(_ *cobra.Command, args []string) error {
			return cobraext.Validate(cfg, &cfg.Dump)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()

				return nil
			}

			cfg.Dump.Type = args[0]

			var data any
			var err error

			switch cfg.Dump.Type {
			case "config":
				if data, err = getConfig(cfg, files); err != nil {
					return fmt.Errorf("getting config: %w", err)
				}

			case "defaults":
				if data, err = getDefaults(cfg, files); err != nil {
					return fmt.Errorf("getting defaults: %w", err)
				}
			case "env":
				if data, err = getEnv(); err != nil {
					return fmt.Errorf("getting env: %w", err)
				}
			case "platform":
				if data, err = getPlatform(); err != nil {
					return fmt.Errorf("getting platform: %w", err)
				}
			case "tools":
				if data, err = getTools(cfg, files); err != nil {
					return fmt.Errorf("getting tools: %w", err)
				}
			default:
				return fmt.Errorf("unknown subcommand: %s", cfg.Dump.Type)
			}

			Print(data, cfg.Dump.Format)

			return nil
		},
	}

	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		command.Flags().MarkHidden("show")
		command.Parent().HelpFunc()(command, strings)
	})

	cmd.PersistentFlags().StringP("format", "f", "yaml", "Output format (json or yaml)")

	return cmd
}

func getConfig(cfg *config.Config, files Embedded) (*config.Config, error) {
	defs := &defaults.Defaults{}
	if err := defs.Load(cfg.Defaults.Name(), files.Defaults); err != nil {
		return nil, fmt.Errorf("error loading defaults: %w", err)
	}

	if err := defs.Merge(*cfg); err != nil {
		return nil, fmt.Errorf("error merging defaults: %w", err)
	}

	return cfg, nil
}

func getDefaults(cfg *config.Config, files Embedded) (*tools.Defaults, error) {
	tools := &tools.Defaults{}
	if err := defaults.LoadDefaults(tools, cfg.Defaults.Name(), files.Defaults, *cfg); err != nil {
		return nil, fmt.Errorf("loading defaults: %w", err)
	}

	return tools, nil
}

func getEnv() (env.Env, error) {
	return env.FromEnv(), nil
}

func getPlatform() (*detect.Platform, error) {
	platform := &detect.Platform{}
	if err := platform.Detect(); err != nil {
		return nil, fmt.Errorf("detecting platform: %w", err)
	}

	return platform, nil
}

func getTools(cfg *config.Config, files Embedded) (any, error) {
	return utils.PrintYAMLBytes(files.Tools), nil
}
