package validate

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
)

// allSubCommands returns all subcommands (recursive) of the given root command.
func allSubCommands(cmd *cobra.Command) []*cobra.Command {
	var cmds []*cobra.Command

	var collect func(*cobra.Command)
	collect = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			cmds = append(cmds, sub)
			collect(sub)
		}
	}

	collect(cmd)

	return cmds
}

func run(cfg config.Config, cmd *cobra.Command, args []string) error {
	all := allSubCommands(cmd.Root())

	errs := []error{}

	for _, sub := range all {
		if sub.PersistentPreRunE == nil {
			continue
		}

		if err := sub.PersistentPreRunE(sub, nil); err != nil {
			errs = append(errs, fmt.Errorf("validating command %q: %w", cmd.CommandPath(), err))
		}
	}

	return errors.Join(errs...)
}
