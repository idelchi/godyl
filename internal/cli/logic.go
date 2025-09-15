package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/internal/tokenstore"
	"github.com/idelchi/godyl/pkg/cobraext"
	penv "github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/koanfx"
	"github.com/idelchi/godyl/pkg/logger"
	pfile "github.com/idelchi/godyl/pkg/path/file"
)

// TODO(Idelchi): Some subcommands should NOT validate the config file, such as `auth`, `config`.

//nolint:maintidx,funlen,gocognit,gocyclo,cyclop // Handles multiple configuration sources and validation steps
func run(cmd *cobra.Command, cfg *root.Config, calledFrom *cobra.Command) error {
	debug.Debug("[PersistentPreRunE root] Current command: %s\n", cmd.CommandPath())
	debug.Debug("[PersistentPreRunE root] Called from: %s\n", calledFrom.CommandPath())

	// Extract some commonly used strings
	commandPath := core.BuildCommandPath(cmd)
	envPrefix := commandPath.Env().Scoped()
	sectionPrefix := commandPath.Section().String()

	// Get the current command's flags
	flags := cmd.Flags()

	// Create a new Koanf instance
	k := koanfx.New()

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
	if err := k.WithFlags(flags).TrackFlags().Load(
		posflag.Provider(flags, "", k),
		nil,
	); err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// configFile := tmproot.ConfigFile
	configFileValue, ok := k.Get("config-file").(string)
	if !ok {
		return errors.New("config-file value is not a string")
	}

	configFile := pfile.New(configFileValue)
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
	k = koanfx.New()

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
	var tmpConfig root.Config
	if err := k.Unmarshal(&tmpConfig, koanfx.WithErrorUnused()); err != nil {
		return fmt.Errorf("unmarshalling config file %q: %w", configFile, err)
	}

	// Store the parsed and validated configuration for future use in subcommands.
	// Each subcommand has to be able to extract it's own subsection, but no longer need to validate it.
	// Make sure this is a copy of the Koanf instance, as it will be modified in the next steps.
	configuration := koanfx.New().WithKoanf(k.Copy())

	// 2nd Pass
	//  - Load the config file
	//	- Load environment variables
	// 	- Load flags
	// Determine the location of the `.env` files

	// Load config file
	// Create a new Koanf instance, basing off the already loaded config file but cut at the relevant section
	k = koanfx.New().WithKoanf(k.Cut(sectionPrefix)).TrackAll().Track()

	// Load environment variables
	if err := k.TrackAll().Load(
		env.Provider(envPrefix.String(), ".", iutils.MatchEnvToFlag(envPrefix)),
		nil,
	); err != nil {
		return fmt.Errorf("loading env vars: %w", err)
	}

	// Load flags
	if err := k.WithFlags(flags).TrackFlags().Load(
		posflag.Provider(flags, "", k),
		nil,
	); err != nil {
		return fmt.Errorf("loading config: %w", err)
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
		}

		dotenvs = dotenvs.MergedWith(dotenv)
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

	env := penv.FromEnv()
	menv := env.MergedWith(dotenvs)

	if err := menv.Export(); err != nil {
		return fmt.Errorf("exporting environment variables: %w", err)
	}

	// At this point, we have env + env-file(s) loaded, with preference to env.
	// This means subsequent subcommands can just use env vars and no longer need to load the env-file(s) again.

	// Set a default value for the tools file if it was not set in the config file or env-file(s)
	if cfg.Tools == "" {
		cfg.Tools = "tools.yml"
	}

	// Store the various processed values in the global context.
	core.GlobalContext.Config = configuration
	core.GlobalContext.Env = &env
	core.GlobalContext.DotEnv = &dotenvs

	// 3rd Pass
	//  - Load the config file
	//	- Load environment variables
	// 	- Load flags
	// Fully populate the configuration struct for this command.
	if err := core.KCreateSubcommandPreRunE(cmd, cfg, root.NoShow)(cmd, []string{}); err != nil {
		return err
	}

	// 4th Pass
	// Default values for tokens are deferred such that they can be
	// set with .env files or the keyring without unnecessary checks

	// TODO(Idelchi): Allow also GITHUB_TOKEN_FILE, GITLAB_TOKEN_FILE, URL_TOKEN_FILE
	// provides valuable context for future development
	githubToken := menv.GetAny("GITHUB_TOKEN", "GH_TOKEN")
	gitlabToken := menv.GetAny("GITLAB_TOKEN", "CI_JOB_TOKEN")
	urlToken := menv.GetAny("URL_TOKEN")

	if !cfg.AllTokensSet() && cfg.Keyring {
		commandPath := calledFrom.CommandPath()
		if !strings.HasPrefix(commandPath, "godyl auth store") {
			store := tokenstore.New()

			if ok, err := store.Available(); !ok {
				return err
			}

			ghToken, _ := store.Get("github-token")
			glToken, _ := store.Get("gitlab-token")
			uToken, _ := store.Get("url-token")

			githubToken = iutils.Any(ghToken, githubToken)
			gitlabToken = iutils.Any(glToken, gitlabToken)
			urlToken = iutils.Any(uToken, urlToken)
		}
	}

	if err := cobraext.SetFlagIfNotSet(flags.Lookup("github-token"), githubToken); err != nil {
		return err
	}

	if err := cobraext.SetFlagIfNotSet(flags.Lookup("gitlab-token"), gitlabToken); err != nil {
		return err
	}

	if err := cobraext.SetFlagIfNotSet(flags.Lookup("url-token"), urlToken); err != nil {
		return err
	}

	// Parse again with the new defaults
	if err := core.KCreateSubcommandPreRunE(cmd, cfg, root.NoShow)(cmd, []string{}); err != nil {
		return err
	}

	if cfg.Tokens.GitHub == "" {
		if !k.Tracker.IsSet("parallel") {
			if err := cobraext.SetFlagIfNotSet(flags.Lookup("parallel"), "1"); err != nil {
				return err
			}

			logWarning = append(
				logWarning,
				"GitHub token is not set. Limiting parallelism to 1 to avoid hitting rate limits.",
			)
		} else {
			logWarning = append(
				logWarning,
				fmt.Sprintf("GitHub token is not set. Make sure you don't hit rate limits with current level: %v", cfg.Parallel),
			)
		}
	}

	// Parse last time
	if err := core.KCreateSubcommandPreRunE(cmd, cfg, cfg.ShowFunc)(cmd, []string{}); err != nil {
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

	if len(logWarning) > 0 && !strings.HasPrefix(calledFrom.CommandPath(), "godyl dump") {
		log.Warn(strings.TrimSpace(strings.Join(logWarning, "\n")))
	}

	return nil
}
