package cli

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/pkg/cobraext"
)

func Command(cfg *config.Config, files *common.Embedded, version string) *cobra.Command {
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args, cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if common.ExitOnShow(cfg.ShowFunc, args...) {
				return nil
			}

			return cobraext.UnknownSubcommandAction(cmd, args)
		},
	}

	cmd.CompletionOptions.DisableDefaultCmd = false
	cmd.Flags().SortFlags = false
	cmd.SetVersionTemplate("{{ .Version }}\n")

	root.Flags(cmd)
	subcommands(cmd, cfg, files)

	return cmd
}
