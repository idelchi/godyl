// Package cli provides the command-line interface for the application.
package cli

import (
	"embed"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// CommandFactory creates and configures commands.
type CommandFactory struct {
	cfg     *config.Config
	version string
	files   config.Embedded
}

// NewRootCmd creates the root command with common configuration.
// It sets up environment variable binding and flag handling.
func NewRootCmd(cfg *config.Config, version string, embeds embed.FS) (*cobra.Command, error) {
	factory := &CommandFactory{
		cfg:     cfg,
		version: version,
	}

	var err error

	factory.files.Defaults, err = embeds.ReadFile("defaults.yml")
	if err != nil {
		return nil, fmt.Errorf("reading defaults file: %w", err)
	}

	factory.files.Tools, err = embeds.ReadFile("tools.yml")
	if err != nil {
		return nil, fmt.Errorf("reading tools file: %w", err)
	}

	factory.files.Template, err = embeds.ReadFile("internal/core/updater/scripts/cleanup.bat.template")
	if err != nil {
		return nil, fmt.Errorf("reading cleanup template: %w", err)
	}

	return factory.CreateRootCommand(), nil
}

// CreateRootCommand creates and configures the root command.
func (f *CommandFactory) CreateRootCommand() *cobra.Command {
	// Custom functions for the NewDefaultRootCommand
	funcs := []func(*cobra.Command, []string) error{
		f.loadDotEnvFunc(),
	}

	root := NewRootCommand(f.version, f.cfg, funcs...)

	// Add subcommands
	f.addSubcommands(root)

	// Add root-level flags
	f.addRootFlags(root)

	// Make certain flags persistent so they can be used with subcommands

	return root
}

// NewDefaultRootCommand creates a root command with default settings.
// It sets up integration with viper, with environment variable and flag binding.
// Additional functions can be passed to be executed before the command is run.
func NewRootCommand(version string, cfg *config.Config, funcs ...func(*cobra.Command, []string) error) *cobra.Command {
	root := &cobra.Command{
		Use:   "godyl [command]",
		Short: "Asset downloader for tools",
		Long: "godyl helps with batch-downloading and installing statically compiled binaries from GitHub releases, " +
			"URLs, and Go projects.",
		Version:          version,
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			viper.SetEnvPrefix(cmd.Root().Name())
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
			viper.AutomaticEnv()

			if err := viper.BindPFlags(cmd.Root().Flags()); err != nil {
				return fmt.Errorf("binding command flags: %w", err)
			}

			if err := viper.Unmarshal(&cfg.Root); err != nil {
				return fmt.Errorf("unmarshalling config: %w", err)
			}

			if err := cfg.Root.Validate(); err != nil {
				return fmt.Errorf("validating config: %w", err)
			}

			for _, f := range funcs {
				if err := f(cmd, args); err != nil {
					return err
				}
			}

			return nil
		},
		RunE: cobraext.UnknownSubcommandAction,
	}

	root.CompletionOptions.DisableDefaultCmd = true
	root.Flags().SortFlags = false

	root.SetVersionTemplate("{{ .Version }}\n")

	return root
}

// loadDotEnvFunc returns a function that loads environment variables from a .env file.
func (f *CommandFactory) loadDotEnvFunc() func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		if err := utils.LoadDotEnv(file.File(viper.GetString("env-file"))); err != nil {
			if config.IsSet("env-file") {
				return fmt.Errorf("loading .env file: %w", err)
			}
		}

		return nil
	}
}
