package cli

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/idelchi/godyl/internal/cli/cache"
	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/cli/download"
	"github.com/idelchi/godyl/internal/cli/dump"
	"github.com/idelchi/godyl/internal/cli/install"
	"github.com/idelchi/godyl/internal/cli/status"
	"github.com/idelchi/godyl/internal/cli/update"
	"github.com/idelchi/godyl/internal/cli/version"
	"github.com/idelchi/godyl/internal/config"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/cobraext"
	penv "github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/koanfx"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/utils"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/cobra"
)

// Command encapsulates a root cobra command with its associated config and embedded files.
type Command struct {
	// Command is the root cobra.Command instance
	Command *cobra.Command
	// Config contains application configuration
	Config *config.Config
	// Files contains the embedded configuration files and templates
	Files *config.Embedded
}

// Run executes the root command.
func (cmd *Command) Run() error {
	return cmd.Command.Execute()
}

// Subcommands adds all subcommands to the root command.
func (cmd *Command) Subcommands() {
	cmd.Command.AddCommand(
		install.NewCommand(cmd.Config, cmd.Files),
		download.NewCommand(cmd.Config, cmd.Files),
		status.NewCommand(cmd.Config, cmd.Files),
		dump.NewCommand(cmd.Config, cmd.Files),
		update.NewCommand(cmd.Config, cmd.Files),
		cache.NewCommand(cmd.Config, cmd.Files),
		version.NewCommand(cmd.Config),
	)
}

// NewRootCommand creates the root cobra command with configuration and embedded files.
func NewRootCommand(cfg *config.Config, files *config.Embedded, version string) *Command {
	cobra.EnableTraverseRunHooks = true
	cobra.EnableCommandSorting = false

	// Create the root command
	root := &cobra.Command{
		Use:   "godyl",
		Short: "Asset downloader for tools",
		Long: heredoc.Doc(`godyl helps with batch-fetching and extracting statically compiled binaries from:
			- GitHub releases
			- GitLab release
			- URLs
			- Go projects.

			as well as providing custom commands.
			`),
		Example: heredoc.Doc(`
			$ godyl install tools.yml
			$ godyl download goreleaser/goreleaser --output /usr/local/bin
			`),
		Version:          version,
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Extract some commonly used strings
			cmd = cmd.Root()

			commandPath := common.BuildCommandPath(cmd)
			envPrefix := commandPath.Env().Scoped()
			sectionPrefix := commandPath.Section().String()
			name := commandPath.Last()

			// Get the current command's flags
			flags := cmd.Flags()

			// Create a new Koanf instance
			k := koanfx.NewWithTracker(flags)

			// 1st Pass
			//	- Load environment variables
			// 	- Load flags
			// Determine the location of the config file

			// Load environment variables
			if err := k.TrackAll().Load(
				env.Provider(envPrefix.String(), ".", utils.MatchEnvToFlag(envPrefix)),
				nil,
			); err != nil {
				return fmt.Errorf("loading env vars: %w", err)
			}

			// Load flags
			if err := k.TrackFlags().Load(
				posflag.Provider(flags, "", k),
				nil,
			); err != nil {
				log.Fatalf("error loading config: %v", err)
			}

			// At this point, the location of the config file can be determined through either
			// the environment variables or the flags.
			// Unmarshal the config into the struct to get the config file path.
			if err := k.Unmarshal(&cfg.Root); err != nil {
				log.Fatalf("Failed to unmarshal config into struct: %v", err)
			}

			// We fail if:
			// The provider returns an error and
			//    the config file exists
			// 	OR
			//    the config file was set explicitly
			failureCriteria := func(err error) bool {
				return err != nil && (cfg.Root.ConfigFile.Exists() || k.IsSet("config-file"))
			}

			// Once the location is known, load the config using a new Koanf instance.
			// This way we can validate the full configuration file, as the environment variables and flags
			// won't match the structure (yet).
			// Load environment variables
			k = koanfx.NewWithTracker(nil)

			if err := k.Load(file.Provider(cfg.Root.ConfigFile.Path()), yaml.Parser()); failureCriteria(err) {
				return fmt.Errorf("loading config file %q: %w", cfg.Root.ConfigFile.Path(), err)
			}

			// We can already validate the config file here by unmarshalling it with koanfx.WithErrorUnused()
			// This will throw an error if there are any unused fields in the config file.
			if err := k.Unmarshal(&cfg, koanfx.WithErrorUnused()); err != nil {
				return fmt.Errorf("unmarshalling config file %q: %w", cfg.Root.ConfigFile.Path(), err)
			}

			// Set the context for future use in subcommands. Each subcommand has to be able to extract
			// it's own subsection, but no longer need to validate it.
			// Make sure this is a copy of the Koanf instance, as it will be modified in the next steps.
			config := koanfx.NewWithTracker(nil)
			cmd.SetContext(context.WithValue(cmd.Context(), "config", config.ResetK(k.Copy())))

			// 2nd Pass
			//  - Load the config file
			//	- Load environment variables
			// 	- Load flags
			// Determine the location of the `.env` files

			// Load config file
			// Create a new Koanf instance, basing off the already loaded config file but cut at the relevant section
			k = koanfx.NewWithTracker(flags).ResetK(k.Cut(sectionPrefix)).TrackAll().Track()

			// Load environment variables
			if err := k.TrackAll().Load(
				env.Provider(envPrefix.String(), ".", utils.MatchEnvToFlag(envPrefix)),
				nil,
			); err != nil {
				return fmt.Errorf("loading env vars: %w", err)
			}

			// Load flags
			if err := k.TrackFlags().Load(
				posflag.Provider(flags, "", k),
				nil,
			); err != nil {
				log.Fatalf("error loading config: %v", err)
			}

			// At this point, the location of the `.env` file(s) can be determined through either
			// the config file, environment variables or the flags.
			// Unmarshal the config into the struct to get the `.env` file path.
			if err := k.Unmarshal(&cfg.Root); err != nil {
				return fmt.Errorf("unmarshalling config into struct: %w", err)
			}

			// We fail if:
			// The provider returns an error and
			//    the `.env` file was set explicitly
			failureCriteria = func(err error) bool {
				return err != nil && k.IsSet("env-file")
			}

			// Load environment variables from .env files such that it's available for the subcommands
			// Precedence is:
			// - Existing env vars
			// - Env vars from env-file (in order from right to left overwriting)
			dotenvs := penv.Env{}
			penv := penv.FromEnv()

			for i := len(cfg.Root.EnvFile) - 1; i >= 0; i-- {
				file := cfg.Root.EnvFile[i]
				dotenv, err := iutils.LoadDotEnv(file)
				if failureCriteria(err) {
					return fmt.Errorf("loading .env file %q: %w", file, err)
				} else {
					dotenvs = dotenvs.MergedWith(dotenv)
				}
			}

			var logWarning bool
			// If the config file or env-file was set in the .env file, we remove it from the loaded env vars and issue a warning
			if dotenvs.Has("GODYL_ENV_FILE") {
				dotenvs.Delete("GODYL_ENV_FILE")
				logWarning = true
			}
			if dotenvs.Has("GODYL_CONFIG_FILE") {
				dotenvs.Delete("GODYL_CONFIG_FILE")
				logWarning = true
			}

			penv = penv.MergedWith(dotenvs)

			if err := penv.ToEnv(); err != nil {
				return fmt.Errorf("setting environment variables: %w", err)
			}

			// At this point, we have env + env-file(s) loaded, with preference to env.
			// This means subsequent subcommands can just use env vars and no longer need to load the env-file(s) again.

			// 3rd Pass
			//  - Load the config file
			//	- Load environment variables
			// 	- Load flags
			// Fully populate the configuration struct for this command.
			if err := common.KCreateSubcommandPreRunE(name, &cfg.Root, &cfg.Root.Show)(cmd, args); err != nil {
				return err
			}

			// Full config available here
			lvl, err := logger.LevelString(cfg.Root.LogLevel)
			if err != nil {
				return fmt.Errorf("parsing log level: %w", err)
			}
			log := logger.New(lvl)

			// Make a simple warning: If the env-file was set from the .envs, the user might be confused
			if logWarning {
				log.Warn(heredoc.Doc(`
					'config-file' and/or 'env-file' was set in the loaded '.env' file(s)
					This might be confusing as it will have no effect.
					Instead, use:
						'GODYL_CONFIG_FILE' or '--config-file'
						'GODYL_ENV_FILE' or '--env-file'
				`))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Root.Show && len(args) == 0 {
				return nil
			}

			return cobraext.UnknownSubcommandAction(cmd, args)
		},
	}

	root.CompletionOptions.DisableDefaultCmd = false
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
