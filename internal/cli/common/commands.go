package common

import (
	"fmt"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/koanfx"
	"github.com/idelchi/godyl/pkg/validator"
)

type Trackable interface {
	Store(tracker *koanfx.Tracker)
	Validate() error
}

func KCreateSubcommandPreRunE(
	cmd *cobra.Command,
	cfg Trackable,
	show root.ShowFuncType,
) func(_ *cobra.Command, _ []string) error {
	return func(_ *cobra.Command, _ []string) error {
		debug.Debug("[PersistentPreRunE] Current command: %s\n", cmd.CommandPath())

		commandPath := BuildCommandPath(cmd)
		envPrefix := commandPath.Env().Scoped()
		sectionPrefix := commandPath.Section()

		defer func() {
			if f := show(); f != nil {
				if cfg != nil || cmd.Flags().NArg() > 0 {
					fmt.Println("")
					fmt.Printf("****** %s ******\n", sectionPrefix)
				}

				if cfg != nil {
					fmt.Println("-- configuration --")
					f(excludeFields(cfg, "validate"))

					fmt.Println("")
				}

				if cmd.Flags().NArg() > 0 {
					fmt.Println("-- arguments --")
					f(cmd.Flags().Args())

					fmt.Println("")
				}
			}
		}()

		// If no config is provided, we don't need to do any parsing anyway.
		if cfg == nil {
			return nil
		}

		// Layering:
		// 0. Flags (defaults)
		// 1. Configuration file
		// 2. Environment variables
		// 3. Flags (overrides)

		// Load configuration file
		koanf := koanfx.New().WithKoanf(GlobalContext.Config.Cut(sectionPrefix.String()))

		// Load environment variables
		if err := koanf.TrackAll().Load(
			env.Provider(envPrefix.String(), ".", iutils.MatchEnvToFlag(envPrefix)),
			nil,
		); err != nil {
			return fmt.Errorf("loading env vars: %w", err)
		}

		// Load flags
		if err := koanf.WithFlags(cmd.Flags()).TrackFlags().Load(
			posflag.Provider(cmd.Flags(), "", koanf),
			nil,
		); err != nil {
			return fmt.Errorf("loading flags: %w", err)
		}

		if err := koanf.Unmarshal(cfg); err != nil {
			return fmt.Errorf("unmarshalling config into struct: %w", err)
		}

		cfg.Store(koanf.Tracker)

		// Validate the configuration
		if err := validator.Validate(cfg); err != nil {
			return fmt.Errorf("validating config for command %q: %w", sectionPrefix, err)
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validating config for command %q: %w", sectionPrefix, err)
		}

		return nil
	}
}
