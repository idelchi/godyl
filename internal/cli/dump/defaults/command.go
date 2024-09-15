// Package defaults implements the defaults dump subcommand for godyl.
// It displays the application's default configuration settings.
package defaults

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	iutils "github.com/idelchi/godyl/internal/utils"
)

// Command encapsulates the defaults dump command with its associated configuration.
type Command struct {
	// Command is the defaults cobra.Command instance
	Command *cobra.Command
}

// NewDefaultsCommand creates a Command for displaying default configuration settings.
func NewDefaultsCommand(cfg *config.Config, embedded *config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:               "defaults [name...]",
		Short:             "Display default configuration settings",
		Args:              cobra.ArbitraryArgs,
		PersistentPreRunE: common.KCreateSubcommandPreRunE("defaults", &cfg.Dump.Tools, &cfg.Root.Show),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Root.Show {
				return nil
			}

			c, err := getDefaults(embedded, args)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance for the defaults dump subcommand.
func NewCommand(cfg *config.Config, files *config.Embedded) *cobra.Command {
	// Create the defaults command
	cmd := NewDefaultsCommand(cfg, files)

	// Add defaults-specific flags
	cmd.Flags()

	return cmd.Command
}

// getDefaults loads and returns the application's default settings.
func getDefaults(files *config.Embedded, defaultNames []string) (any, error) {
	var defaults map[string]any

	err := yaml.Unmarshal(files.Defaults, &defaults)
	if err != nil {
		return nil, err
	}

	// If defaultNames is provided, filter the defaults
	if len(defaultNames) > 0 {
		filteredDefaults := make(map[string]any)
		for _, name := range defaultNames {
			if value, exists := defaults[name]; exists {
				filteredDefaults[name] = value
			} else {
				return nil, fmt.Errorf("default %q not found", name)
			}
		}
		defaults = filteredDefaults
	}

	return defaults, nil
}
