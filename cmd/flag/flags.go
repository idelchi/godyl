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

func flags() {
	pflag.Bool("update", false, "")
	pflag.String("update-strategy", "", "")
}

func parseFlags() (cfg Config, err error) {
	flags()

	// Parse the command-line flags
	// pflag.Parse()
	// Parse the command-line flags with suggestions enabled
	if err := flagexp.ParseWithSuggestions(os.Args[1:]); err != nil {
		return cfg, fmt.Errorf("parsing flags: %w", err)
	}

	// viper.BindPFlag("update-now", pflag.CommandLine.Lookup("update"))
	// viper.BindPFlag("update.strategy", pflag.CommandLine.Lookup("strategy"))

	// Bind pflag flags to viper
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return cfg, fmt.Errorf("binding flags: %w", err)
	}

	// Set viper to automatically read from environment variables
	viper.SetEnvPrefix("godyl")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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
