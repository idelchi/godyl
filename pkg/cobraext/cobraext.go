package cobraext

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// UnknownSubcommandAction handles unknown cobra subcommands.
// Implements cobra.Command.RunE to provide helpful error messages
// and suggestions when an unknown subcommand is used. Required
// when TraverseChildren is true, as this disables cobra's built-in
// suggestion system. See:
// - https://github.com/spf13/cobra/issues/981
// - https://github.com/containerd/nerdctl/blob/242e6fc6e861b61b878bd7df8bf25e95674c036d/cmd/nerdctl/main.go#L401-L418
func UnknownSubcommandAction(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help() //nolint: wrapcheck
	}

	err := fmt.Sprintf("unknown subcommand %q for %q", args[0], cmd.Name())

	if suggestions := cmd.SuggestionsFor(args[0]); len(suggestions) > 0 {
		err += "\n\nDid you mean this?\n"
		for _, s := range suggestions {
			err += fmt.Sprintf("\t%v\n", s)
		}
	}

	return errors.New(err) //nolint: err113
}

// IsSet checks if a flag was explicitly set.
// Uses viper to determine if a flag was provided by the user,
// avoiding the use of default values unless specifically passed.
func IsSet(flag string) bool {
	return viper.IsSet(flag)
}
