package common

import (
	"fmt"
	"log"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/pkg/koanfx"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/godyl/pkg/utils"
	"github.com/idelchi/godyl/pkg/validator"
)

func GetCommand(cmd *cobra.Command, command string) *cobra.Command {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Name() == command {
			return c
		}
	}

	return nil
}

type Trackable interface {
	StoreTracker(tracker *koanfx.Tracker)
	Validate() error
}

func KCreateSubcommandPreRunE(name string, cfg Trackable, show *bool) func(cmd *cobra.Command, arg []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// If no config is provided, we don't need to do any parsing anyway.
		if cfg == nil {
			return nil
		}

		// We might be called from a subcommand - means we have to find the correct parent
		if cmd.HasParent() {
			if cmd = GetCommand(cmd, name); cmd == nil {
				return fmt.Errorf("command not found")
			}
		}

		commandPath := BuildCommandPath(cmd)
		envPrefix := commandPath.Env().Scoped()
		sectionPrefix := commandPath.Section()
		name = commandPath.Last()

		if cmd.HasParent() {
			// Strip the root command from the section prefix
			sectionPrefix = commandPath.WithoutRoot().Section()
		}

		flags := cmd.Flags()

		// Load in the configuration file as base
		k, ok := cmd.Root().Context().Value("config").(*koanfx.KoanfWithTracker)
		if !ok {
			return fmt.Errorf("failed to get config from context")
		}

		// Cut out the section for this (sub)command
		k = koanfx.NewWithTracker(flags).ResetK(k.Cut(sectionPrefix.String())).TrackAll().Track()

		// Load environment variables
		if err := k.TrackAll().Load(
			env.Provider(envPrefix.String(), ".", utils.MatchEnvToFlag(envPrefix)),
			nil,
		); err != nil {
			return fmt.Errorf("loading env vars: %w", err)
		}

		// Load flags
		if err := k.TrackFlags().Load(
			posflag.Provider(flags, "", k),
			nil,
		); err != nil {
			log.Fatalf("error loading config: %v", err)
		}

		// Unmarshal the config into the struct to get the `.env` file path.
		if err := k.Unmarshal(cfg); err != nil {
			return fmt.Errorf("unmarshalling config into struct: %w", err)
		}

		cfg.StoreTracker(k.Tracker)

		// Validate the configuration
		if err := validator.Validate(cfg); err != nil {
			return fmt.Errorf("validating config: %w", err)
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validating config: %w", err)
		}

		if *show {
			fmt.Println("")
			fmt.Printf("****** %s ******\n", name)
			fmt.Println("-- configuration --")
			pretty.PrintYAMLMasked(cfg)

			if len(args) > 0 {
				fmt.Println("-- arguments --")
				pretty.PrintYAMLMasked(args)
			}
		}

		return nil
	}
}

type CommandPath []string

func (c CommandPath) Section() utils.Prefix {
	return utils.Prefix(strings.Join(c, ".")).Lower()
}

func (c CommandPath) Env() utils.Prefix {
	return utils.Prefix(strings.Join(c, "_")).Upper()
}

func (c CommandPath) WithoutRoot() CommandPath {
	return c[1:]
}

func (c CommandPath) Last() string {
	return c[len(c)-1]
}

func BuildCommandPath(cmd *cobra.Command) CommandPath {
	var parts CommandPath

	// Start with the current command
	current := cmd
	parts = append(parts, current.Name())

	// Add all ancestors
	for current.Parent() != nil {
		current = current.Parent()
		parts = append(CommandPath{current.Name()}, parts...)
	}

	return parts
}
