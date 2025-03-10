package cobraext

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// UnknownSubcommandAction is a cobra.Command.RunE function that prints an error message for unknown subcommands.
// Necessary when using `TraverseChildren: true`, because it seems to disable suggestions for unknown subcommands.
// See:
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

// IsSet checks if a flag is set in viper,
// to avoid using it's default values unless explicitly passed.
func IsSet(flag string) bool {
	return viper.IsSet(flag)
}
