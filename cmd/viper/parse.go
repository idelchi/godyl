package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/idelchi/godyl/pkg/flagexp"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// configIsSet indicates whether the configuration flag is set,
// either via the command-line or environment variable.
func IsSet(flag string) bool {
	// return pflag.CommandLine.Changed(flag)
	return viper.IsSet(flag)
}

func flags() {
	pflag.StringP("config", "c", "config.yml", "Path to configuration file")
	pflag.String("dot-env", ".env", "Path to .env file")

	pflag.CommandLine.SortFlags = false
}

// parseFlags parses the application configuration (in order of precedence) from:
//   - command-line flags
//   - environment variables
//   - configuration file
func parseFlags() (cfg Flags, err error) {
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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	decoderConfig := func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnused = true // Throw error on unknown fields
	}

	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(&cfg, decoderConfig); err != nil {
		return cfg, fmt.Errorf("unmarshalling config: %w", err)
	}

	return cfg, nil
}
