// Package tools contains the subcommand definition for `dump tools`.
package tools

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/dump/tools"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/ierrors"
)

// Command returns the `dump tools` command.
func Command(global *root.Config, local any, embedded *core.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools [tools.yml|-]...",
		Short: "Display tools information",
		Long: heredoc.Doc(`
			Dumps out tools configuration.

			Use with the --tags flag to filter the output and createc ustom tools configuration files,
			or pipe the output back to the tool for installation.

			The command will default to 'tools.yml' if no file is specified.

			If the --embedded flag is used, the command will dump out the embedded 'tools.yml' file.
		`),
		Example: heredoc.Doc(`
			# Dump out the embedded tools.yml file
			$ godyl dump tools --embedded

			# Filter the embedded tools.yml file to only include go and python tools
			# and dump it to dev-tools.yml
			$ godyl dump tools --embedded --tags=go,python > dev-tools.yml

			# Dump out contents of dev-tools.yml
			$ godyl dump tools dev-tools.yml

			# Pipe the contents of dev-tools.yml to godyl dump tools
			$ cat dev-tools.yml | godyl dump tools - --tags=go,python

			# Pipe the embedded tools.yml file to godyl install
			$ godyl dump tools --embedded --tags go | godyl install -
		`),
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if core.ExitOnShow(global.ShowFunc) {
				return nil
			}

			if global.Dump.Tools.Embedded && len(args) > 0 {
				return fmt.Errorf(
					"%w: cannot specify arguments together with the --embedded flag, use one or the other",
					ierrors.ErrUsage,
				)
			}

			return run(core.Input{Global: global, Cmd: cmd, Args: args, Embedded: embedded})
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	tools.Flags(cmd)

	return cmd
}
