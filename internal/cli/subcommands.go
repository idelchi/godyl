// Package cli provides the command-line interface for the application.
package cli

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/dump"
	"github.com/idelchi/godyl/internal/cli/tool"
	"github.com/idelchi/godyl/internal/cli/version"
	"github.com/idelchi/godyl/pkg/logger"
)

// addSubcommands adds all subcommands to the root command.
func (f *CommandFactory) addSubcommands(root *cobra.Command) {
	root.AddCommand(
		version.NewCommand(root.Version),
		dump.NewCommand(f.cfg, f.files),
		tool.NewInstallCommand(f.cfg, f.files),
		tool.NewDownloadCommand(f.cfg, f.files),
		tool.NewUpdateCommand(f.cfg, f.files),
	)
}

// addRootFlags adds root-level flags to the command.
func (f *CommandFactory) addRootFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("dry", false, "Run without making any changes (dry run)")
	cmd.Flags().String("log", logger.INFO.String(), "Log level (DEBUG, INFO, WARN, ERROR, SILENT)")
	cmd.Flags().String("env-file", ".env", "Path to .env file")
	cmd.Flags().StringP("defaults", "d", "defaults.yml", "Path to defaults file")
	cmd.Flags().BoolP("show", "s", false, "Show the configuration and exit")
}
