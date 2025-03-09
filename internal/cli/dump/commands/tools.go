package dump

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/utils"
	iutils "github.com/idelchi/godyl/internal/utils"
)

func NewToolsCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Display tools information",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return commonPreRunE(cmd, &cfg.Dump)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := getTools(cfg, files)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return cmd
}

func getTools(cfg *config.Config, files config.Embedded) (any, error) {
	return utils.PrintYAMLBytes(files.Tools), nil
}
