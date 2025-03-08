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
		Use:     "dump [config|defaults|env|platform|tools]",
		Aliases: []string{"show"},
		Short:   "Dump configuration information",
		Long:    "Display various configuration settings and information about the environment",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return cobraext.Validate(cfg)
		},
	}

	// Add subcommands
	cmd.AddCommand(
		newDumpEnvCommand(cfg),
		newDumpConfigCommand(cfg, files.Defaults),
		newDumpDefaultsCommand(cfg, files.Defaults),
		newDumpPlatformCommand(cfg),
		newDumpToolsCommand(cfg, files.Tools),
	)

	cmd.PersistentFlags().StringP("format", "f", "yaml", "Output format (json or yaml)")

	return cmd
}

// newDumpPlatformCommand creates a command to show platform information.
func newDumpPlatformCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "platform",
		Short: "Dump platform information",
		Long:  "Display information about the current platform (OS, architecture, etc.)",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return cobraext.Validate(cfg)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			p := detect.Platform{}
			if err := p.Detect(); err != nil {
				return fmt.Errorf("detecting platform: %w", err)
			}

			Print(p, cfg.Format)

			return nil
		},
	}
}

// newDumpConfigCommand creates a command to show the current configuration.
func newDumpConfigCommand(cfg *config.Config, defaultsData []byte) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Dump the current configuration",
		Long:  "Display the current configuration settings",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return cobraext.Validate(cfg)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			defs := defaults.Defaults{}
			if err := defs.Load(cfg.Defaults.Name(), defaultsData); err != nil {
				return fmt.Errorf("error loading defaults: %w", err)
			}

			if err := defs.Merge(*cfg); err != nil {
				return fmt.Errorf("error merging defaults: %w", err)
			}

			Print(cfg, cfg.Format)

			return nil
		},
	}
}

// newDumpDefaultsCommand creates a command to show the default configuration.
func newDumpDefaultsCommand(cfg *config.Config, defaultsData []byte) *cobra.Command {
	return &cobra.Command{
		Use:   "defaults",
		Short: "Dump the default configuration",
		Long:  "Display the default configuration settings",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return cobraext.Validate(cfg)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			toolDefaults := tools.Defaults{}
			if err := defaults.LoadDefaults(&toolDefaults, cfg.Defaults.Name(), defaultsData, *cfg); err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			Print(toolDefaults, cfg.Format)

			return nil
		},
	}
}

// newDumpEnvCommand creates a command to show environment variables.
func newDumpEnvCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "Dump environment variables",
		Long:  "Display environment variables that affect the application",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return cobraext.Validate(cfg)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			Print(env.FromEnv(), cfg.Format)

			return nil
		},
	}
}

// newDumpToolsCommand creates a command to show available tools.
func newDumpToolsCommand(cfg *config.Config, toolsData []byte) *cobra.Command {
	return &cobra.Command{
		Use:   "tools",
		Short: "Dump available tools",
		Long:  "Display information about available tools",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return cobraext.Validate(cfg)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			Print(utils.PrintYAMLBytes(toolsData), cfg.Format)

			return nil
		},
	}
}
