package cobraext

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewDefaultRootCommand creates a root command with default settings.
// It sets up integration with viper, with environment variable and flag binding.
// Additional functions can be passed to be executed before the command is run.
func NewDefaultRootCommand(version string, funcs ...func(*cobra.Command, []string) error) *cobra.Command {
	root := &cobra.Command{
		Version:          version,
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true, // TODO(Idelchi): Breaks suggestions, see below.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			viper.SetEnvPrefix(cmd.Root().Name())
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
			viper.AutomaticEnv()

			if err := viper.BindPFlags(cmd.Root().Flags()); err != nil {
				return fmt.Errorf("binding root flags: %w", err)
			}

			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return fmt.Errorf("binding command flags: %w", err)
			}

			for _, f := range funcs {
				if err := f(cmd, args); err != nil {
					return err
				}
			}

			return nil
		},
		RunE: UnknownSubcommandAction,
	}

	root.CompletionOptions.DisableDefaultCmd = true
	root.Flags().SortFlags = false

	root.SetVersionTemplate("{{ .Version }}\n")

	return root
}

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
