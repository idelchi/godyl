package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/data"
	"github.com/idelchi/godyl/pkg/logger"
)

// Flags adds the flags for the `godyl` root command to the provided Cobra command.
func Flags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.Flags().StringP("config-file", "c", data.ConfigFile().Path(), "path to config file")
	cmd.Flags().StringSliceP("env-file", "e", []string{".env"}, "paths to .env files")
	cmd.Flags().StringP("defaults", "d", "defaults.yml", "path to defaults file")
	cmd.Flags().StringP("inherit", "i", "default", "default to inherit from when unset in the tool spec")
	cmd.PersistentFlags().CountP("show", "s", "show the parsed configuration and exit, repeat for unmasking tokens.")

	cmd.Flags().StringP("cache-dir", "", data.UserDataDir().Path(), "path to cache directory")

	cmd.Flags().String("github-token", "", "github api token, defaulting to keyring, GITHUB_TOKEN or GH_TOKEN")
	cmd.Flags().String("gitlab-token", "", "gitlab api token, default to keyring, GITLAB_TOKEN or CI_JOB_TOKEN")
	cmd.Flags().String("url-token", "", "url api token, defaulting to keyring or URL_TOKEN")

	cmd.Flags().Bool("keyring", false, "enable token retrieval from keyring")
	cmd.Flags().IntP("parallel", "j", 0, "parallelism, 0 means unlimited.")
	cmd.Flags().StringP("log-level", "l", logger.INFO.String(), fmt.Sprintf("log level (%v)", logger.LevelValues()))

	cmd.Flags().BoolP("no-cache", "", false, "disable cache")
	cmd.Flags().BoolP("no-verify-ssl", "k", false, "skip SSL verification")
	cmd.Flags().Bool("no-progress", false, "disable progress bar")
	cmd.Flags().BoolP("no-verify-checksum", "C", false, "skip checksum verification")

	cmd.Flags().StringP("error-file", "", "", "path to error log file, empty means stdout.")
	cmd.Flags().CountP("verbose", "v", "increase verbosity (can be used multiple times)")
}
