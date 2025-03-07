// Package commands provides the command-line interface for the application.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// rootEmbedded struct to hold the embedded default and tools configuration.
type rootEmbedded struct {
	// defaults holds the embedded default configuration.
	defaults []byte

	// tools holds the embedded tools configuration.
	tools []byte

	// embeds holds static template scripts.
	embeds interface{}
}

// NewRootCmd creates the root command with common configuration.
// It sets up environment variable binding and flag handling.
func NewRootCmd(cfg *config.Config, version string, defaultsFile, toolsFile []byte, embeds interface{}) *cobra.Command {
	emb := rootEmbedded{
		defaults: defaultsFile,
		tools:    toolsFile,
		embeds:   embeds,
	}

	// Custom functions for the NewDefaultRootCommand
	funcs := []func(*cobra.Command, []string) error{
		func(cmd *cobra.Command, args []string) error {
			if err := loadDotEnv(file.File(viper.GetString("env-file"))); err != nil {
				if config.IsSet("env-file") {
					return fmt.Errorf("loading .env file: %w", err)
				}
			}
			return nil
		},
	}

	root := cobraext.NewDefaultRootCommand(version, funcs...)

	root.Use = "godyl [command]"
	root.Short = "Asset downloader for tools"
	root.Long = "godyl helps with batch-downloading and installing statically compiled binaries from GitHub releases, URLs, and Go projects."

	// Add subcommands
	root.AddCommand(
		NewDumpCommand(cfg, emb),
		NewUpdateCommand(cfg, emb),
		NewInstallCommand(cfg, emb),
		NewDownloadCommand(cfg, emb),
	)

	// Add root-level flags
	addRootFlags(root)

	// Make certain flags persistent so they can be used with subcommands
	root.PersistentFlags().BoolP("show", "s", false, "Show the configuration and exit")

	return root
}

// addRootFlags adds root-level flags to the command.
func addRootFlags(cmd *cobra.Command) {
	// Root command flags
	cmd.Flags().Bool("dry", false, "Run without making any changes (dry run)")
	cmd.Flags().String("log", logger.INFO.String(), "Log level (debug, info, warn, error, silent)")
	cmd.Flags().IntP("parallel", "j", 0, "Number of parallel downloads. 0 means unlimited.")
	cmd.Flags().BoolP("no-verify-ssl", "k", false, "Skip SSL verification")
	cmd.Flags().String("env-file", ".env", "Path to .env file")
	cmd.Flags().StringP("defaults", "d", "defaults.yml", "Path to defaults file")
}

// addToolFlags adds tool-related flags to the command.
func addToolFlags(cmd *cobra.Command) {
	// Configuration file flags

	// Tool flags
	cmd.Flags().StringP("output", "o", "./bin", "Output path for the downloaded tools")
	cmd.Flags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
	cmd.Flags().String("source", "github", "Source from which to install the tools")
	cmd.Flags().String("strategy", "none", "Strategy to use for updating tools")
	cmd.Flags().String("github-token", "", "GitHub token for authentication")
	cmd.Flags().String("os", "", "Operating system to install the tools for")
	cmd.Flags().String("arch", "", "Architecture to install the tools for")
}
