package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/idelchi/godyl/internal/flagexp"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var config = ConfigFile{
	"config.yml",
}

type ConfigFile struct {
	Default string
}

func (c ConfigFile) Get() string {
	env := c.Env()
	if env != "" {
		return env
	}

	return c.Default
}

func (ConfigFile) Env() string {
	return os.Getenv("GODYL_CONFIG")
}

// configIsSet indicates whether the configuration flag is set,
// either via the command-line or environment variable.
func (c ConfigFile) IsSet() bool {
	return pflag.CommandLine.Changed("config") || c.Env() != ""
}

func flags() {
	// General flags
	pflag.Bool("version", false, "Show the version information and exit")
	pflag.BoolP("help", "h", false, "Show the help information and exit")
	pflag.BoolP("show", "s", false, "Show the configuration and exit")
	pflag.StringP("config", "c", config.Get(), "Path to configuration file")
	pflag.IntP("parallel", "j", 0, "Number of parallel downloads")

	// Selected custom flags
	pflag.String("defaults.source.github.token", "", "GitHub token for API requests")
	pflag.StringSliceP("defaults.extensions", "e", nil, "Extensions to filter tools by")
	pflag.String("defaults.strategy", "none", "")
	pflag.String("defaults.output", "~/.local/bin", "")

	pflag.StringSliceP("tags", "t", nil, "Tags to filter tools by")

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
//   - configuration file
func parseFlags() (cfg Config, err error) {
	flags()

	// Parse the command-line flags
	// pflag.Parse()
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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString("config"))
	viper.SetConfigType("yaml")
	// Only return errors if the config flag is set or the environment variable is set
	// Otherwise, either read the configuration or try to read the default value
	if err := viper.ReadInConfig(); err != nil && config.IsSet() {
		fmt.Println("Error reading config file:", err)

		return cfg, err
	} else if err != nil {
		fmt.Println("Setting defaults to config file")

		cfg.Default()
	}

	decoderConfig := func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnused = true // Throw error on unknown fields
	}

	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(&cfg, decoderConfig); err != nil {
		return cfg, fmt.Errorf("unmarshalling config: %w", err)
	}

	// Handle the commandline flags that exit the application
	handleExitFlags(cfg)

	// Validate the input
	if err := validateInput(&cfg); err != nil {
		return cfg, fmt.Errorf("validating input: %w", err)
	}

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
	if viper.GetBool("version") {
		fmt.Println(version)
		os.Exit(0)
	}

	// Check if the help flag was provided
	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}

	if viper.GetBool("show") {
		fmt.Println(PrintJSON(cfg))

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
