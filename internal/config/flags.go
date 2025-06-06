package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/logger"
)

// Flags adds the flags for the `godyl` root command to the provided Cobra command.
func Flags(cmd *cobra.Command) {
	env := env.FromEnv()

	cmd.Flags().StringP("log-level", "l", logger.INFO.String(), fmt.Sprintf("log level (%v)", logger.LevelValues()))
	cmd.Flags().IntP("parallel", "j", 0, "parallelism, 0 means unlimited.")
	cmd.Flags().StringP("cache-dir", "", tmp.CacheDir().Path(), "path to cache directory")
	cmd.Flags().BoolP("no-cache", "", false, "disable cache")
	cmd.Flags().BoolP("no-verify-ssl", "k", false, "skip SSL verification")
	cmd.Flags().Bool("no-progress", false, "disable progress bar")
	cmd.PersistentFlags().CountP("show", "s", "show the parsed flags and exit, repeat for unmasking tokens.")

	cmd.Flags().StringP("config-file", "c", tmp.ConfigFile().Path(), "path to config file")
	cmd.Flags().StringSliceP("env-file", "e", []string{".env"}, "paths to .env files")
	cmd.Flags().StringP("defaults", "d", "defaults.yml", "path to defaults file")
	cmd.Flags().StringP("inherit", "i", "default", "default to inherit from when unset in the tool spec")

	cmd.Flags().
		String("github-token", env.GetAny("GITHUB_TOKEN", "GH_TOKEN"), "github token for authentication")
	cmd.Flags().
		String("gitlab-token", env.GetAny("GODYL_GITLAB_TOKEN", "GITLAB_TOKEN"), "gitlab token for authentication")
	cmd.Flags().String("url-token", env.GetAny("GODYL_URL_TOKEN", "URL_TOKEN"), "url token for authentication")

	cmd.Flags().StringP("error-file", "", "", "path to error log file, empty means stdout.")
	cmd.Flags().CountP("verbose", "v", "increase verbosity (can be used multiple times)")
}
