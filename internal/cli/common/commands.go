package common

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/koanfx"
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
	Store(tracker *koanfx.Tracker)
	Validate() error
}

func excludeSubcommandsFields(s any) any {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	var (
		fields []reflect.StructField
		values []reflect.Value
	)

	for i := range t.NumField() {
		field := t.Field(i)
		if field.Tag.Get("validate") != "-" && !field.Anonymous {
			fields = append(fields, field)
			values = append(values, v.Field(i))
		}
	}

	newType := reflect.StructOf(fields)
	newStruct := reflect.New(newType).Elem()

	for i, val := range values {
		newStruct.Field(i).Set(val)
	}

	return newStruct.Interface()
}

func KCreateSubcommandPreRunE(
	name string,
	cfg Trackable,
	show root.ShowFuncType,
) func(cmd *cobra.Command, arg []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// We might be called from a subcommand - means we have to find the correct parent
		if cmd.HasParent() {
			if cmd = GetCommand(cmd, name); cmd == nil {
				return errors.New("command not found")
			}
		}

		commandPath := BuildCommandPath(cmd)
		envPrefix := commandPath.Env().Scoped()
		sectionPrefix := commandPath.Section()
		name = commandPath.Last()

		if cmd.HasParent() {
			// Strip the root command from the section prefix
			sectionPrefix = commandPath.WithoutRoot().Section()
		} else {
			// If we are at the root command, we don't need a section prefix
			sectionPrefix = ""
		}

		defer func() {
			if f := show(); f != nil {
				fmt.Println("")
				fmt.Printf("****** %s ******\n", sectionPrefix)

				if cfg != nil {
					fmt.Println("-- configuration --")
					f(excludeSubcommandsFields(cfg))

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

		flags := cmd.Flags()

		// Load in the configuration file as base
		k, ok := cmd.Root().Context().Value("config").(*koanfx.KoanfWithTracker)
		if !ok {
			return errors.New("failed to get config from context")
		}

		// Cut out the section for this (sub)command
		k = koanfx.NewWithTracker(flags).ResetKoanf(k.Cut(sectionPrefix.String())).TrackAll().Track()

		// Load environment variables
		if err := k.TrackAll().Load(
			env.Provider(envPrefix.String(), ".", iutils.MatchEnvToFlag(envPrefix)),
			nil,
		); err != nil {
			return fmt.Errorf("loading env vars: %w", err)
		}

		// Load flags
		if err := k.TrackFlags().Load(
			posflag.Provider(flags, "", k),
			nil,
		); err != nil {
			return fmt.Errorf("loading flags: %w", err)
		}

		// Unmarshal the config into the struct to get the `.env` file path.
		if err := k.Unmarshal(cfg); err != nil {
			return fmt.Errorf("unmarshalling config into struct: %w", err)
		}

		cfg.Store(k.Tracker)

		// Validate the configuration
		if err := validator.Validate(cfg); err != nil {
			return fmt.Errorf("validating config: %w", err)
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validating config: %w", err)
		}

		return nil
	}
}

type CommandPath []string

func (c CommandPath) Section() iutils.Prefix {
	return iutils.Prefix(strings.Join(c, ".")).Lower()
}

func (c CommandPath) Env() iutils.Prefix {
	return iutils.Prefix(strings.Join(c, "_")).Upper()
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
