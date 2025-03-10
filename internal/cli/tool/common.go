package tool

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/internal/config"
)

func addToolFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("output", "o", "./bin", "Output path for the downloaded tools")
	cmd.Flags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
	cmd.Flags().String("source", "github", "Source from which to install the tools")
	cmd.Flags().String("strategy", "none", "Strategy to use for updating tools")
	cmd.Flags().String("github-token", "", "GitHub token for authentication")
	cmd.Flags().String("os", "", "Operating system to install the tools for")
	cmd.Flags().String("arch", "", "Architecture to install the tools for")
	cmd.Flags().BoolP("no-verify-ssl", "k", false, "Skip SSL verification")
	cmd.Flags().IntP("parallel", "j", 0, "Number of parallel downloads. 0 means unlimited.")
	cmd.Flags().StringSlice("hints", []string{""}, "Hints to use for tool resolution")
}

func addUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().String("github-token", "", "GitHub token for authentication")
	cmd.Flags().BoolP("no-verify-ssl", "k", false, "Skip SSL verification")
}

func commonPreRunE(cmd *cobra.Command, tool *config.Tool) error {
	viper.SetEnvPrefix(cmd.Root().Name())
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("binding flags: %w", err)
	}

	if err := viper.Unmarshal(tool); err != nil {
		return fmt.Errorf("unmarshalling config: %w", err)
	}

	return nil
}
