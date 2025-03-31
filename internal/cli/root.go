package cli

import (
	"embed"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/download"
	"github.com/idelchi/godyl/internal/cli/dump"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/cli/install"
	"github.com/idelchi/godyl/internal/cli/update"
	"github.com/idelchi/godyl/internal/cli/version"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/cobraext"
	"github.com/idelchi/godyl/pkg/file"
)

// Command encapsulates a root cobra command with its associated config and embedded files.
type Command struct {
	// Command is the root cobra.Command instance
	Command *cobra.Command
	// Config contains application configuration
	Config *config.Config
	// Files contains the embedded configuration files and templates
	Files config.Embedded
}

// Run executes the root command.
func (cmd *Command) Run() error {
	return cmd.Command.Execute()
}

// Flags adds all root-level flags to the command.
func (cmd *Command) Flags() {
	flags.Root(cmd.Command)
}

// Subcommands adds all subcommands to the root command.
func (cmd *Command) Subcommands() {
	cmd.Command.AddCommand(
		version.NewCommand(),
		dump.NewCommand(cmd.Config, cmd.Files),
		install.NewCommand(cmd.Config, cmd.Files),
		download.NewCommand(cmd.Config, cmd.Files),
		update.NewCommand(cmd.Config, cmd.Files),
	)
}

// NewRootCommand creates the root cobra command with configuration and embedded files.
func NewRootCommand(cfg *config.Config, files config.Embedded, version string) *Command {
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
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			// Bind root-level flags
			if err := flags.Bind(cmd.Root(), &cfg.Root); err != nil {
				return fmt.Errorf("binding flags: %w", err)
			}

			// Validate the root configuration
			if err := cfg.Root.Validate(); err != nil {
				return fmt.Errorf("validating config: %w", err)
			}

			// Load environment variables from .env file such that it's available for the subcommands
			if err := utils.LoadDotEnv(file.File(cfg.Root.EnvFile)); err != nil {
				if cfg.Root.IsSet("env-file") {
					return fmt.Errorf("loading .env file: %w", err)
				}
			}

			// Bind root-level flags
			// Once more to get the .env file values too
			if err := flags.Bind(cmd.Root(), &cfg.Root); err != nil {
				return fmt.Errorf("binding flags: %w", err)
			}

			return nil
		},
		RunE: cobraext.UnknownSubcommandAction,
	}

	root.CompletionOptions.DisableDefaultCmd = true
	root.Flags().SortFlags = false
	root.SetVersionTemplate("{{ .Version }}\n")

	return &Command{
		Command: root,
		Config:  cfg,
		Files:   files,
	}
}

// NewCommand creates a fully configured Command instance with embedded files and subcommands.
func NewCommand(cfg *config.Config, version string, embeds embed.FS) (*Command, error) {
	// Get the embedded files
	files, err := config.NewEmbeddedFiles(embeds)
	if err != nil {
		return nil, fmt.Errorf("creating embedded files: %w", err)
	}

	// Create the root command
	cmd := NewRootCommand(cfg, files, version)

	// Add root-level flags
	cmd.Flags()

	// Add subcommands
	cmd.Subcommands()

	return cmd, nil
}
