package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/flagexp"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func IsSet(flag string) bool {
	return viper.IsSet(flag)
}

func flags() {
	pflag.String("output", "", "Output path for the downloaded tools")
	pflag.String("tools", "", "Path to tools configuration file")
	pflag.StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by")
	pflag.StringP("defaults", "d", "defaults.yml", "Path to defaults file")

	// Update flags
	pflag.Bool("update", false, "Update the tools")
	pflag.String("strategy", string(tools.None), "Strategy to use for updating tools")

	pflag.Bool("dry", false, "Run without making any changes (dry run)")
	pflag.Bool("detect", false, "Detect the platform and exit")
	pflag.String("log", string(logger.INFO), "Log level (DEBUG, INFO, WARN, ERROR)")

	// Tokens flags
	pflag.String("github-token", "", "GitHub token for authentication")

	pflag.String("source", "", "Source from which to install the tools")

	pflag.String("dot-env", "", "Path to .env file")

	pflag.BoolP("help", "h", false, "Show help message and exit")
	pflag.Bool("show-config", false, "Show the parsed configuration and exit")
	pflag.Bool("show-defaults", false, "Show the parsed default configuration and exit")
	pflag.Bool("show-env", false, "Show the parsed environment variables and exit")
	pflag.Bool("version", false, "Show version information and exit")

	pflag.IntP("parallel", "j", 0, "Number of parallel downloads")

	pflag.CommandLine.SortFlags = false
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [tools]\n\n", "godyl")
		fmt.Fprintf(os.Stderr, "Tool manager that installs tools as specified in a YAML file.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		pflag.PrintDefaults()
	}
}

// parseFlags parses the application configuration (in order of precedence) from:
//   - command-line flags
//   - environment variables
func parseFlags() (cfg Config, err error) {
	flags()

	// Parse the command-line flags with suggestions enabled
	if err := flagexp.ParseWithSuggestions(os.Args[1:]); err != nil {
		return cfg, fmt.Errorf("parsing flags: %w", err)
	}

	// Bind pflag flags to viper
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return cfg, fmt.Errorf("binding flags: %w", err)
	}

	// Set viper to automatically read from environment variables
	viper.SetEnvPrefix("godyl")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if IsSet("dot-env") {
		if err := loadDotEnv(file.File(viper.GetString("dot-env"))); err != nil {
			return cfg, fmt.Errorf("loading .env file: %w", err)
		}
	}

	decoderConfig := func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnused = true // Throw error on unknown fields
	}

	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(&cfg, decoderConfig); err != nil {
		return cfg, fmt.Errorf("unmarshalling config: %w", err)
	}

	// Validate the input
	if err := validateInput(&cfg); err != nil {
		return cfg, fmt.Errorf("validating input: %w", err)
	}

	// Handle the commandline flags that exit the application
	handleExitFlags(cfg)

	return cfg, nil
}

func validateInput(cfg *Config) error {
	switch pflag.NArg() {
	case 0:
		cfg.Tools = "tools.yml"
	case 1:
		cfg.Tools = pflag.Arg(0)
	default:
		return fmt.Errorf("too many arguments: %d", pflag.NArg())
	}

	return nil
}

//nolint:forbidigo // Function will print & exit for various help messages.
func handleExitFlags(cfg Config) {
	// Check if the version flag was provided
	if cfg.Version {
		fmt.Println(version)

		os.Exit(0)
	}

	// Check if the help flag was provided
	if cfg.Help {
		pflag.Usage()

		os.Exit(0)
	}

	if cfg.Show.Config {
		pretty.PrintYAML(cfg)

		os.Exit(0)
	}

	if cfg.Show.Env {
		pretty.PrintYAML(env.FromEnv())

		os.Exit(0)
	}

	if cfg.Show.Defaults {
		defaults := Defaults{}
		if err := defaults.Load(cfg.Defaults.Name()); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading defaults: %v\n", err)

			os.Exit(1)
		}

		defaults.Merge(cfg)

		pretty.PrintYAML(defaults)

		os.Exit(0)
	}
}

// PrintJSON returns a pretty-printed JSON representation of the provided object.
func PrintJSON(obj any) string {
	bytes, err := json.MarshalIndent(obj, "  ", "    ")
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}

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
