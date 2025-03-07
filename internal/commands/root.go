// Package commands provides the command-line interface for the application.
package commands

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
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

	root := cobraext.NewDefaultRootCommand(version)

	root.Use = "godyl [command]"
	root.Short = "Asset downloader for tools"
	root.Long = "godyl helps with batch-downloading and installing statically compiled binaries from GitHub releases, URLs, and Go projects."

	// Disable direct execution with tools.yml
	root.Args = cobra.NoArgs

	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Skip validation for help command
		if cmd.Name() == "help" {
			return nil
		}

		// Check for version flag
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag && cmd.Name() == root.Name() {
			cmd.Print(version + "\n")
			return cobraext.ErrExitGracefully
		}

		// Set viper to automatically read from environment variables
		viper.SetEnvPrefix("godyl")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		viper.AutomaticEnv()

		// Bind flags to viper
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return fmt.Errorf("binding flags: %w", err)
		}

		if err := loadDotEnv(file.File(viper.GetString("dot-env"))); err != nil {
			if config.IsSet("dot-env") {
				return fmt.Errorf("loading .env file: %w", err)
			}
		}

		decoderConfig := func(dc *mapstructure.DecoderConfig) {
			dc.ErrorUnused = true // Throw error on unknown fields
		}

		// Unmarshal the configuration into the Config struct
		if err := viper.Unmarshal(cfg, decoderConfig); err != nil {
			return fmt.Errorf("unmarshalling config: %w", err)
		}

		// Handle exit flags
		if err := handleExitFlags(cmd, version, cfg, emb.defaults); err != nil {
			return err
		}

		// Validate the configuration
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validating configuration: %w", err)
		}

		return nil
	}

	// Add subcommands
	root.AddCommand(
		NewDumpCommand(cfg, emb),
		NewUpdateCommand(cfg),
		NewInstallCommand(cfg),
		NewDownloadCommand(cfg),
	)

	// Set up flags
	setupRootFlags(root)

	return root
}

// setupRootFlags configures all command-line flags for the application.
func setupRootFlags(root *cobra.Command) {
	// Configuration file flags
	root.PersistentFlags().String("dot-env", ".env", "Path to .env file")
	root.PersistentFlags().StringP("defaults", "d", "defaults.yml", "Path to defaults file")

	// Application flags
	root.PersistentFlags().BoolP("show", "s", false, "Show the configuration and exit")
	root.PersistentFlags().Bool("dry", false, "Run without making any changes (dry run)")
	root.PersistentFlags().String("log", logger.INFO.String(), "Log level (debug, info, warn, error, silent)")
	root.PersistentFlags().IntP("parallel", "j", 0, "Number of parallel downloads. 0 means unlimited.")
	root.PersistentFlags().BoolP("no-verify-ssl", "k", false, "Skip SSL verification")

	// Tool flags
	root.PersistentFlags().StringP("output", "o", "", "Output path for the downloaded tools")
	root.PersistentFlags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
	root.PersistentFlags().String("source", "github", "Source from which to install the tools")
	root.PersistentFlags().String("strategy", "none", "Strategy to use for updating tools")
	root.PersistentFlags().String("github-token", "", "GitHub token for authentication")
	root.PersistentFlags().String("os", "", "Operating system to install the tools for")
	root.PersistentFlags().String("arch", "", "Architecture to install the tools for")

	root.PersistentFlags().SortFlags = false
}
