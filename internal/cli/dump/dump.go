package dump

import (
	"github.com/spf13/cobra"

	dump "github.com/idelchi/godyl/internal/cli/dump/commands"
	"github.com/idelchi/godyl/internal/config"
)

func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dump [command]",
		Aliases: []string{"show"},
		Short:   "Dump configuration information",
		Long:    "Display various configuration settings and information about the environment",
	}

	cmd.Flags().StringP("format", "f", "yaml", "Output format (json or yaml)")

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	cmd.AddCommand(
		dump.NewConfigCommand(cfg, files),
		dump.NewDefaultsCommand(cfg, files),
		dump.NewEnvCommand(cfg, files),
		dump.NewPlatformCommand(cfg, files),
		dump.NewToolsCommand(cfg, files),
	)

	return cmd
}
