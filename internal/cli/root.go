// Package cli provides the command-line interface for the application.
package cli

import (
	"embed"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// CommandFactory creates and configures commands.
type CommandFactory struct {
	cfg     *config.Config
	version string
	files   Embedded
}

// Embedded holds the embedded files for the application.
type Embedded struct {
	Defaults []byte
	Tools    []byte
	Template []byte
}

// NewRootCmd creates the root command with common configuration.
// It sets up environment variable binding and flag handling.
func NewRootCmd(cfg *config.Config, version string, embeds embed.FS) (*cobra.Command, error) {
	e := Embedded{}
	var err error

	e.Defaults, err = embeds.ReadFile("defaults.yml")
	if err != nil {
		return nil, fmt.Errorf("reading defaults file: %w", err)
	}

	e.Tools, err = embeds.ReadFile("tools.yml")
	if err != nil {
		return nil, fmt.Errorf("reading tools file: %w", err)
	}

	e.Template, err = embeds.ReadFile("cleanup.bat.template")
	if err != nil {
		return nil, fmt.Errorf("reading cleanup template: %w", err)
	}

	factory := &CommandFactory{
		cfg:     cfg,
		version: version,
		files:   e,
	}
	return factory.CreateRootCommand(), nil
}

// CreateRootCommand creates and configures the root command.
func (f *CommandFactory) CreateRootCommand() *cobra.Command {
	// Custom functions for the NewDefaultRootCommand
	funcs := []func(*cobra.Command, []string) error{
		f.loadDotEnvFunc(),
	}

	root := cobraext.NewDefaultRootCommand(f.version, funcs...)

	root.Use = "godyl [command]"
	root.Short = "Asset downloader for tools"
	root.Long = "godyl helps with batch-downloading and installing statically compiled binaries from GitHub releases, URLs, and Go projects."

	// Add subcommands
	f.addSubcommands(root)

	// Add root-level flags
	f.addRootFlags(root)

	// Make certain flags persistent so they can be used with subcommands
	root.PersistentFlags().BoolP("show", "s", false, "Show the configuration and exit")

	return root
}

// loadDotEnvFunc returns a function that loads environment variables from a .env file.
func (f *CommandFactory) loadDotEnvFunc() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := utils.LoadDotEnv(file.File(viper.GetString("env-file"))); err != nil {
			if config.IsSet("env-file") {
				return fmt.Errorf("loading .env file: %w", err)
			}
		}
		return nil
	}
}

// addSubcommands adds all subcommands to the root command.
func (f *CommandFactory) addSubcommands(root *cobra.Command) {
	root.AddCommand(
		NewDumpCommand(f.cfg, f.files),
		NewInstallCommand(f.cfg, f.files),
		NewDownloadCommand(f.cfg, f.files),
		NewUpdateCommand(f.cfg, f.files),
	)
}

// addRootFlags adds root-level flags to the command.
func (f *CommandFactory) addRootFlags(cmd *cobra.Command) {
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
	// Tool flags
	cmd.Flags().StringP("output", "o", "./bin", "Output path for the downloaded tools")
	cmd.Flags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
	cmd.Flags().String("source", "github", "Source from which to install the tools")
	cmd.Flags().String("strategy", "none", "Strategy to use for updating tools")
	cmd.Flags().String("github-token", "", "GitHub token for authentication")
	cmd.Flags().String("os", "", "Operating system to install the tools for")
	cmd.Flags().String("arch", "", "Architecture to install the tools for")
}
