// Package cli provides the command-line interface for the application.
package cli

import (
	"embed"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/internal/cli/download"
	"github.com/idelchi/godyl/internal/cli/dump"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/cli/install"
	"github.com/idelchi/godyl/internal/cli/update"
	"github.com/idelchi/godyl/internal/cli/version"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// NewRootCmd creates the root command with all configuration.
func NewRootCmd(cfg *config.Config, version string, embeds embed.FS) (*cobra.Command, error) {
	files := config.Embedded{}
	var err error

	// Read embedded files
	if files.Defaults, err = embeds.ReadFile("defaults.yml"); err != nil {
		return nil, fmt.Errorf("reading defaults file: %w", err)
	}

	if files.Tools, err = embeds.ReadFile("tools.yml"); err != nil {
		return nil, fmt.Errorf("reading tools file: %w", err)
	}

	if files.Template, err = embeds.ReadFile("internal/core/updater/scripts/cleanup.bat.template"); err != nil {
		return nil, fmt.Errorf("reading cleanup template: %w", err)
	}

	// Create the root command
	root := &cobra.Command{
		Use:   "godyl [command]",
		Short: "Asset downloader for tools",
		Long: "godyl helps with batch-fetching and extracting statically compiled binaries from GitHub releases, " +
			"URLs, and Go projects.",
		Version:          version,
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Bind(cmd.Root(), cmd.Root().Name(), &cfg.Root); err != nil {
				return fmt.Errorf("binding flags: %w", err)
			}

			if err := cfg.Root.Validate(); err != nil {
				return fmt.Errorf("validating config: %w", err)
			}

			// Load environment variables from .env file
			if err := utils.LoadDotEnv(file.File(viper.GetString("env-file"))); err != nil {
				if config.IsSet("env-file") {
					return fmt.Errorf("loading .env file: %w", err)
				}
			}

			return nil
		},
		RunE: cobraext.UnknownSubcommandAction,
	}

	root.CompletionOptions.DisableDefaultCmd = true
	root.Flags().SortFlags = false
	root.SetVersionTemplate("{{ .Version }}\n")

	// Add root-level flags
	flags.Root(root)

	// Add subcommands
	subcommands(root, cfg, files)

	return root, nil
}

func subcommands(cmd *cobra.Command, cfg *config.Config, files config.Embedded) {
	cmd.AddCommand(
		version.NewCommand(cmd.Version),
		dump.NewCommand(cfg, files),
		install.NewCommand(cfg, files),
		download.NewCommand(cfg, files),
		update.NewCommand(cfg, files),
	)
}
