package cli

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/pkg/cobraext"
)

// Command returns the root `godyl` command.
func Command(files *core.Embedded, version string) *cobra.Command {
	cfg := &root.Config{}

	cobra.EnableTraverseRunHooks = true
	cobra.EnableCommandSorting = false

	// Create the root command
	cmd := &cobra.Command{
		Use:   "godyl",
		Short: "Asset downloader for tools",
		Long: heredoc.Doc(`godyl helps with batch-fetching and extracting statically compiled binaries from:
			- GitHub releases
			- GitLab release
			- URLs
			- Go projects.

			as well as providing custom commands.
			`),
		Example: heredoc.Doc(`
			$ godyl install tools.yml
			$ godyl download goreleaser/goreleaser --output /usr/local/bin
			`),
		Version:          version,
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	return run(cmd, cfg)
		// },
		RunE: func(cmd *cobra.Command, args []string) error {
			if core.ExitOnShow(cfg.ShowFunc, args...) {
				return nil
			}

			return cobraext.UnknownSubcommandAction(cmd, args)
		},
	}

	cmd.PersistentPreRunE = func(calledFrom *cobra.Command, _ []string) error {
		return run(cmd, cfg, calledFrom)
	}

	cmd.CompletionOptions.DisableDefaultCmd = false
	cmd.SetVersionTemplate("{{ .Version }}\n")

	root.Flags(cmd)
	subcommands(cmd, cfg, files)

	return cmd
}
