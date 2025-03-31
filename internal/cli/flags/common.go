package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Viperable is an interface for types that can hold a viper instance
type Viperable interface {
	SetViper(v *viper.Viper)
}

// Bind connects cobra flags to viper and unmarshals the configuration into the provided struct.
// It sets up environment variable handling with the given prefix and handles flag binding.
// Omit the prefix to use the command hierarchy as the prefix.
func Bind(cmd *cobra.Command, cfg Viperable, prefix ...string) error {
	// Set up Viper with our environment prefix
	envPrefix := prefixFromCmdOrPrefixes(cmd, prefix...)

	viper := viper.New()

	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("binding flags: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshalling config for %q: %w", cmd.Name(), err)
	}

	// Set the viper instance on the config struct
	cfg.SetViper(viper)

	fmt.Printf("Processed cmd %s with prefix %s\n", cmd.Name(), envPrefix)

	return nil
}

// prefixFromCmdOrPrefixes builds an environment variable prefix string either from
// a command's hierarchy or from explicitly provided prefix parts.
// When prefixes are provided, they take precedence over the command hierarchy.
func prefixFromCmdOrPrefixes(cmd *cobra.Command, prefixes ...string) string {
	if len(prefixes) > 0 {
		// Use explicitly provided prefixes if available
		return strings.Join(prefixes, "_")
	}

	// Otherwise build prefix from command hierarchy
	var commandPathParts []string
	currentCmd := cmd

	// Traverse up the command tree to build the path
	for currentCmd != nil {
		commandPathParts = append([]string{currentCmd.Name()}, commandPathParts...)
		currentCmd = currentCmd.Parent()
	}

	return strings.Join(commandPathParts, "_")
}
