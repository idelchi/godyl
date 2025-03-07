package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// validateInput validates the command-line arguments.
func validateInput(cfg *config.Config, args []string) error {
	switch len(args) {
	case 0:
		cfg.Tools = "tools.yml"
	case 1:
		cfg.Tools = args[0]
	}

	return nil
}

// handleExitFlags handles flags that cause the application to exit.
func handleExitFlags(cmd *cobra.Command, version string, cfg *config.Config, defaultEmbedded []byte) error {
	// Check if the version flag was provided
	if cfg.Version {
		fmt.Println(version)
		return cobraext.ErrExitGracefully
	}

	// Check if the help flag was provided
	if cfg.Help {
		cmd.Help()
		return cobraext.ErrExitGracefully
	}

	if cfg.Dump.Config {
		pretty.PrintYAMLMasked(*cfg)
		return cobraext.ErrExitGracefully
	}

	if cfg.Dump.Env {
		pretty.PrintYAMLMasked(env.FromEnv())
		return cobraext.ErrExitGracefully
	}

	if cfg.Dump.Defaults {
		defaults := tools.Defaults{}
		if err := loadDefaults(&defaults, cfg.Defaults.Name(), defaultEmbedded, *cfg); err != nil {
			return fmt.Errorf("loading defaults: %w", err)
		}

		pretty.PrintYAMLMasked(defaults)
		return cobraext.ErrExitGracefully
	}

	if cfg.Dump.Platform {
		p := detect.Platform{}
		if err := p.Detect(); err != nil {
			return fmt.Errorf("detecting platform: %w", err)
		}

		pretty.PrintYAML(p)
		return cobraext.ErrExitGracefully
	}

	if cfg.Dump.Tools {
		// This is handled by the dump-tools command
		return cobraext.ErrExitGracefully
	}

	return nil
}

// loadDotEnv loads environment variables from a .env file.
func loadDotEnv(path file.File) error {
	dotEnv, err := env.FromDotEnv(path.Name())
	if err != nil {
		return fmt.Errorf("loading environment variables from %q: %w", path.Name(), err)
	}

	env := env.FromEnv().Normalized().Merged(dotEnv.Normalized())

	if err := env.ToEnv(); err != nil {
		return fmt.Errorf("setting environment variables: %w", err)
	}

	return nil
}

// loadDefaults loads the default configuration.
func loadDefaults(defaults *tools.Defaults, path string, defaultEmbedded []byte, cfg config.Config) error {
	if config.IsSet("defaults") {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading defaults file %q: %w", path, err)
		}

		if err := yaml.Unmarshal(data, defaults); err != nil {
			return fmt.Errorf("unmarshalling defaults: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(defaultEmbedded, defaults); err != nil {
			return fmt.Errorf("unmarshalling embedded defaults: %w", err)
		}
	}

	if err := defaults.Initialize(); err != nil {
		return fmt.Errorf("initializing defaults: %w", err)
	}

	// Apply configuration overrides
	if config.IsSet("output") {
		defaults.Output = cfg.Output
	}

	if config.IsSet("source") {
		defaults.Source.Type = cfg.Source
	}

	if config.IsSet("strategy") {
		defaults.Strategy = cfg.Strategy
	}

	if config.IsSet("github-token") {
		defaults.Source.Github.Token = cfg.Tokens.GitHub
	}

	if config.IsSet("os") {
		if err := defaults.Platform.OS.Parse(cfg.OS); err != nil {
			return fmt.Errorf("parsing OS: %w", err)
		}
		defaults.Platform.Extension = defaults.Platform.Extension.Default(defaults.Platform.OS)
		defaults.Platform.Library = defaults.Platform.Library.Default(defaults.Platform.OS, defaults.Platform.Distribution)
	}

	if config.IsSet("arch") {
		if err := defaults.Platform.Architecture.Parse(cfg.Arch); err != nil {
			return fmt.Errorf("parsing architecture: %w", err)
		}
	}

	return nil
}

// loadTools loads the tools configuration.
func loadTools(path string, log *logger.Logger) (tools.Tools, error) {
	var toolsList tools.Tools

	if err := toolsList.Load(path); err != nil {
		return toolsList, fmt.Errorf("loading tools from %q: %w", path, err)
	}

	log.Info("loaded %d tools from %q", len(toolsList), path)

	return toolsList, nil
}

// splitTags splits tags into include and exclude lists.
func splitTags(tags []string) ([]string, []string) {
	var withTags, withoutTags []string

	for _, tag := range tags {
		if strings.HasPrefix(tag, "!") {
			withoutTags = append(withoutTags, tag[1:])
		} else {
			withTags = append(withTags, tag)
		}
	}

	return withTags, withoutTags
}
