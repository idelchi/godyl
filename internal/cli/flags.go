package cli

import (
	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/logger"
)

func (cmd *Command) Flags() {
	env := env.FromEnv()

	cmd.Command.Flags().StringP("log-level", "l", logger.INFO.String(), "Log level (silent, debug, info, warn, error, always)")
	cmd.Command.Flags().IntP("parallel", "j", 0, "Parallelism. 0 means unlimited.")
	cmd.Command.Flags().StringP("cache-dir", "", tmp.CacheDir().Path(), "Path to cache directory")
	cmd.Command.Flags().BoolP("no-cache", "", false, "Disable cache")
	cmd.Command.Flags().BoolP("no-verify-ssl", "k", false, "Skip SSL verification")
	cmd.Command.Flags().Bool("no-progress", false, "Disable progress bar")
	cmd.Command.PersistentFlags().BoolP("show", "s", false, "Show the parsed flags and exit")

	cmd.Command.Flags().StringP("config-file", "c", tmp.ConfigFile().Path(), "Path to config file")
	cmd.Command.Flags().StringSliceP("env-file", "e", []string{".env"}, "Paths to .env files")
	cmd.Command.Flags().StringP("defaults", "d", "defaults.yml", "Path to defaults file")
	cmd.Command.Flags().StringP("inherit", "i", "default", "Default to inherit from when unset in the tool spec")

	cmd.Command.Flags().String("github-token", env.GetAny("GITHUB_TOKEN", "GH_TOKEN"), "GitHub token for authentication")
	cmd.Command.Flags().String("gitlab-token", env.GetAny("GODYL_GITLAB_TOKEN", "GITLAB_TOKEN"), "GitLab token for authentication")
	cmd.Command.Flags().String("url-token", env.GetAny("GODYL_URL_TOKEN", "URL_TOKEN"), "URL token for authentication")
	cmd.Command.Flags().String("url-token-header", "Authorization", "URL token for authentication")

	cmd.Command.Flags().StringP("error-file", "", "", "Path to error log file. Empty means stdout.")
	cmd.Command.Flags().BoolP("detailed", "", false, "Show detailed messages")
	cmd.Command.Flags().CountP("verbose", "v", "increase verbosity (can be used multiple times)")
}
