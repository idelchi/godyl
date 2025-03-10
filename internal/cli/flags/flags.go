package flags

import (
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/spf13/cobra"
)

func Root(cmd *cobra.Command) {
	cmd.Flags().Bool("dry", false, "Run without making any changes (dry run)")
	cmd.Flags().String("log", logger.INFO.String(), "Log level (DEBUG, INFO, WARN, ERROR, SILENT)")
	cmd.Flags().String("env-file", ".env", "Path to .env file")
	cmd.Flags().StringP("defaults", "d", "defaults.yml", "Path to defaults file")
	cmd.Flags().BoolP("show", "s", false, "Show the configuration and exit")
}

func Tool(cmd *cobra.Command) {
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

func Update(cmd *cobra.Command) {
	cmd.Flags().String("github-token", "brooo", "GitHub token for authentication")
	cmd.Flags().BoolP("no-verify-ssl", "k", false, "Skip SSL verification")
}
