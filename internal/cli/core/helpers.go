package core

import (
	"reflect"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config/root"
)

// ExitOnShow checks if the command should exit early based on the show function and arguments.
// Returns true if the command should exit without further processing.
// ExitOnShow determines if the command should exit early when show function is active and no arguments are provided.
func ExitOnShow(show root.ShowFuncType, args ...string) bool {
	if show() != nil && len(args) == 0 {
		return true
	}

	return false
}

// SetSubcommandDefaults configures default behavior for a subcommand, including setting up
// the PersistentPreRunE function with the provided configuration and show function.
// SetSubcommandDefaults configures default settings for subcommands including pre-run hooks and configuration handling.
func SetSubcommandDefaults(cmd *cobra.Command, local any, show root.ShowFuncType) {
	var config Trackable

	if local != nil {
		local, ok := local.(Trackable)
		if !ok {
			panic("configuration may only be passed as Trackable type")
		}

		config = local
	}

	cmd.PersistentPreRunE = KCreateSubcommandPreRunE(cmd, config, show)
}

// Input is the input structure for the CLI commands.
type Input struct {
	// Global contains the global configuration.
	Global *root.Config
	// Embedded contains the embedded configuration.
	Embedded *Embedded
	// Cmd is the current command being executed.
	Cmd *cobra.Command
	// Args are the arguments passed to the command.
	Args []string
}

// Unpack is a convenience method to unpack the Input structure into its components.
func (i Input) Unpack() (*root.Config, *Embedded, *Context, *cobra.Command, []string) {
	return i.Global, i.Embedded, &GlobalContext, i.Cmd, i.Args
}

func excludeFields(s any, tag string) any {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
		t = t.Elem()
	}

	var (
		fields []reflect.StructField
		values []reflect.Value
	)

	for i := range t.NumField() {
		field := t.Field(i)
		if field.Tag.Get(tag) != "-" && !field.Anonymous {
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
