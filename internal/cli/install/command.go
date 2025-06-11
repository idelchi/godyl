// Package install contains the subcommand definition for `install`.
package install

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/install"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `install` command.
func Command(global *root.Config, local any, embedded *common.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install [tools.yml|-]...",
		Aliases: []string{"i"},
		Short:   "Install tools from one of more YAML files",
		Long: heredoc.Doc(`
		Install tools as specified in the YAML file(s).
		Use '-' to read from stdin.
		Can be combined with reading from files.
		`),
		Example: heredoc.Doc(`
			# Install tools from 'tools.yml' (if existing in the current directory)
			$ godyl install

			# Install tools from 'tools1.yml' and 'tools2.yml'
			$ godyl install tools1.yml - tools2.yml

			# Install tools from stdin
			$ cat tools.yml | godyl install -

			# Install tools from stdin and files
			$ cat tools.yml | godyl install - tools1.yml tools2.yml
		`),
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(common.Input{Global: global, Cmd: cmd, Args: args, Embedded: embedded})
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	install.Flags(cmd)

	return cmd
}
