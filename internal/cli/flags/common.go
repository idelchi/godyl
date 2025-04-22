package flags

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Viperable is an interface for types that can hold a viper instance.
type Viperable interface {
	SetViper(v *viper.Viper)
	GetViper() *viper.Viper
}

func PrefixToYAML(prefix, root string) string {
	prefix = strings.TrimPrefix(prefix, root)
	prefix = strings.ReplaceAll(prefix, "_", ".")
	prefix = strings.TrimPrefix(prefix, ".")

	return prefix
}

// Bind connects cobra flags to viper and unmarshals the configuration into the provided struct.
// It sets up environment variable handling with the given prefix and handles flag binding.
// Omit the prefix to use the command hierarchy as the prefix.
func Bind(cmd *cobra.Command, cfg Viperable, prefix ...string) error {
	// Set up Viper with our environment prefix
	envPrefix := prefixFromCmdOrPrefixes(cmd, prefix...)

	// Reuse the same instance if already set
	if cfg.GetViper() == nil {
		cfg.SetViper(viper.New())
	}

	viper := cfg.GetViper()

	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	configFile := cmd.Root().Context().Value("config-file")
	isSet := cmd.Root().Context().Value("config-file-set")

	isConfigError := func(err error) bool {
		return err != nil && isSet != nil && isSet.(bool)
	}

	if configFile != nil {
		config := file.File(configFile.(string))

		root := cmd.Root().Name()
		if root == cmd.Name() {
			root = ""
		}

		prefix := PrefixToYAML(envPrefix, root)

		content, err := Trim(config, prefix)
		if isConfigError(err) {
			return fmt.Errorf("trimming config file: %w", err)
		} else if err == nil {
			viper.SetConfigType("yaml")

			if err := viper.ReadConfig(content); isConfigError(err) {
				return fmt.Errorf("reading config file: %w", err)
			}
		}
	}

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("binding flags: %w", err)
	}

	var md mapstructure.Metadata

	if err := viper.Unmarshal(cfg, func(config *mapstructure.DecoderConfig) {
		config.Metadata = &md
	}); err != nil {
		return fmt.Errorf("unmarshalling config for %q: %w", cmd.Name(), err)
	}

	for _, val := range md.Unused {
		if val != "help" {
			return fmt.Errorf("unrecognized config key %q in %q", val, cmd.Name())
		}
	}

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
