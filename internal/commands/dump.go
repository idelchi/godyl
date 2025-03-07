package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// NewDumpCommand creates the show command for displaying various configurations.
func NewDumpCommand(cfg *config.Config, emb rootEmbedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump [config|defaults|env|platform|tools]",
		Short: "Dump configuration information",
		Long:  "Display various configuration settings and information about the environment",
	}

	// Add subcommands
	cmd.AddCommand(
		newDumpConfigCommand(cfg, emb.defaults),
		newDumpDefaultsCommand(cfg, emb.defaults),
		newDumpEnvCommand(),
		newDumpPlatformCommand(),
		newDumpToolsCommand(emb.tools),
	)

	return cmd
}

// newDumpConfigCommand creates a command to show the current configuration.
func newDumpConfigCommand(cfg *config.Config, defaultsData []byte) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Dump the current configuration",
		Long:  "Display the current configuration settings",
		PreRunE: func(_ *cobra.Command, args []string) error {
			return cobraext.Validate(cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			defaults := Defaults{}
			if err := defaults.Load(cfg.Defaults.Name(), defaultsData); err != nil {
				return fmt.Errorf("error loading defaults: %v", err)
			}
			if err := defaults.Merge(*cfg); err != nil {
				return fmt.Errorf("error merging defaults: %v", err)
			}

			pretty.PrintYAMLMasked(defaults)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			defaults := tools.Defaults{}
			if err := loadDefaults(&defaults, cfg.Defaults.Name(), defaultsData, *cfg); err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			pretty.PrintYAMLMasked(defaults)
			return nil
		},
	}
}

// newDumpEnvCommand creates a command to show environment variables.
func newDumpEnvCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "Dump environment variables",
		Long:  "Display environment variables that affect the application",
		RunE: func(cmd *cobra.Command, args []string) error {
			pretty.PrintYAMLMasked(env.FromEnv())
			return nil
		},
	}
}

// newDumpPlatformCommand creates a command to show platform information.
func newDumpPlatformCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "platform",
		Short: "Dump platform information",
		Long:  "Display information about the current platform (OS, architecture, etc.)",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := detect.Platform{}
			if err := p.Detect(); err != nil {
				return fmt.Errorf("detecting platform: %w", err)
			}

			pretty.PrintYAML(p)
			return nil
		},
	}
}

// newDumpToolsCommand creates a command to show available tools.
func newDumpToolsCommand(toolsData []byte) *cobra.Command {
	return &cobra.Command{
		Use:   "tools",
		Short: "Dump available tools",
		Long:  "Display information about available tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			pretty.PrintYAML(PrintYAMLBytes(toolsData))
			return nil
		},
	}
}

// PrintYAMLBytes parses YAML bytes into a generic data structure.
func PrintYAMLBytes(yamlBytes []byte) any {
	var data any
	if err := yaml.Unmarshal(yamlBytes, &data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return data
}
