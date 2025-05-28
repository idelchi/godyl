package cli

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/iutils"
	penv "github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/koanfx"
	"github.com/idelchi/godyl/pkg/logger"
	pfile "github.com/idelchi/godyl/pkg/path/file"
)

func run(cmd *cobra.Command, args []string, cfg *config.Config) error {
	// Extract some commonly used strings
	cmd = cmd.Root()

	commandPath := common.BuildCommandPath(cmd)
	envPrefix := commandPath.Env().Scoped()
	sectionPrefix := ""
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
		env.Provider(envPrefix.String(), ".", iutils.MatchEnvToFlag(envPrefix)),
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
	// Unmarshal the config into a temporary struct to get the config file path.
	// var tmpConfig config.Config
	// if err := k.Unmarshal(&tmpConfig); err != nil {
	// 	log.Fatalf("Failed to unmarshal config into struct: %v", err)
	// }

	// configFile := tmpConfig.ConfigFile
	configFile := pfile.New(k.Get("config-file").(string))
	isSet := k.IsSet("config-file")

	// We fail if:
	// The provider returns an error and
	//    the config file exists
	// 	OR
	//    the config file was set explicitly
	failureCriteria := func(err error) bool {
		return err != nil && (configFile.Exists() || isSet)
	}

	// Once the location is known, load the config using a new Koanf instance.
	// This way we can validate the full configuration file, as the environment variables and flags
	// won't match the structure (yet).
	// Load environment variables
	k = koanfx.NewWithTracker(nil)

	if err := k.Load(file.Provider(configFile.Path()), yaml.Parser()); failureCriteria(err) {
		return fmt.Errorf("loading config file %q: %w", configFile, err)
	}

	var logWarning []string

	// If the config file was set in the config file, issue a warning
	if k.Exists("config-file") {
		logWarning = append(logWarning, heredoc.Docf(`
					%q was set to %q in the loaded config file %q.
					This might be confusing as it will have no effect.
					Instead, use only:
						'GODYL_CONFIG_FILE' or '--config-file')
					`, "config-file", k.Get("config-file"), configFile))

		k.Delete("config-file") // Clear the value that was collected from the config file
	}

	// We can already validate the config file here by unmarshalling it with koanfx.WithErrorUnused()
	// This will throw an error if there are any unused fields in the config file.
	var tmpConfig config.Config
	if err := k.Unmarshal(&tmpConfig, koanfx.WithErrorUnused()); err != nil {
		return fmt.Errorf("unmarshalling config file %q: %w", configFile, err)
	}

	// Set the context for future use in subcommands. Each subcommand has to be able to extract
	// it's own subsection, but no longer need to validate it.
	// Make sure this is a copy of the Koanf instance, as it will be modified in the next steps.
	config := koanfx.NewWithTracker(nil)
	cmd.SetContext(context.WithValue(cmd.Context(), "config", config.ResetKoanf(k.Copy())))

	// 2nd Pass
	//  - Load the config file
	//	- Load environment variables
	// 	- Load flags
	// Determine the location of the `.env` files

	// Load config file
	// Create a new Koanf instance, basing off the already loaded config file but cut at the relevant section
	k = koanfx.NewWithTracker(flags).ResetKoanf(k.Cut(sectionPrefix)).TrackAll().Track()

	// Load environment variables
	if err := k.TrackAll().Load(
		env.Provider(envPrefix.String(), ".", iutils.MatchEnvToFlag(envPrefix)),
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
	if err := k.Unmarshal(&tmpConfig); err != nil {
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

	for i := len(tmpConfig.EnvFile) - 1; i >= 0; i-- {
		file := tmpConfig.EnvFile[i]
		dotenv, err := iutils.LoadDotEnv(file)

		if failureCriteria(err) {
			return err
		} else {
			dotenvs = dotenvs.MergedWith(dotenv)
		}
	}

	// If the config file or env-file was set in the .env file,
	// we remove it from the loaded env vars and issue a warning
	if dotenvs.Exists("GODYL_ENV_FILE") {
		logWarning = append(logWarning, heredoc.Docf(`
					%q was set to %q in the loaded .env file(s) %q.
					This might be confusing as it will have no effect.
					Instead, use only:
						'GODYL_ENV_FILE' or '--env-file'
				`, "env-file", dotenvs.Get("GODYL_ENV_FILE"), tmpConfig.EnvFile))

		dotenvs.Delete("GODYL_ENV_FILE")
	}

	if dotenvs.Exists("GODYL_CONFIG_FILE") {
		logWarning = append(logWarning, heredoc.Docf(`
					%q was set to %q in the loaded .env file(s) %q.
					This might be confusing as it will have no effect.
					Instead, use:
						'GODYL_CONFIG_FILE' or '--config-file'
				`, "config-file", dotenvs.Get("GODYL_CONFIG_FILE"), tmpConfig.EnvFile))

		dotenvs.Delete("GODYL_CONFIG_FILE")
	}

	cmd.SetContext(context.WithValue(cmd.Context(), "dotenv", dotenvs))

	penv := penv.FromEnv()
	penv = penv.MergedWith(dotenvs)

	if err := penv.Export(); err != nil {
		return fmt.Errorf("exporting environment variables: %w", err)
	}

	// At this point, we have env + env-file(s) loaded, with preference to env.
	// This means subsequent subcommands can just use env vars and no longer need to load the env-file(s) again.

	// Set a default value for the tools file if it was not set in the config file or env-file(s)
	if cfg.Tools == "" {
		cfg.Tools = "tools.yml"
	}

	// 3rd Pass
	//  - Load the config file
	//	- Load environment variables
	// 	- Load flags
	// Fully populate the configuration struct for this command.
	if err := common.KCreateSubcommandPreRunE(name, cfg, cfg.ShowFunc)(cmd, args); err != nil {
		return err
	}

	// Full config available here
	lvl, err := logger.LevelString(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("parsing log level: %w", err)
	}

	log, err := logger.New(lvl)
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	// Make a simple warning: If the env-file was set from the .envs, the user might be confused
	if len(logWarning) > 0 {
		log.Warn(strings.TrimSpace(strings.Join(logWarning, "\n")))
	}

	return nil
}
