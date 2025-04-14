package flags

import (
	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/spf13/cobra"
)

// Root adds the root-level command flags to the provided cobra command.
// These flags apply to the root command.
func Root(cmd *cobra.Command) {
	env := env.FromEnv()

	// Only flags that are not directly translatable to a `defaults` setting should have a default value here.
	// Or if additional env variables are to be used.
	cmd.Flags().String("log", logger.INFO.String(), "Log level (silent, debug, info, warn, error, always)")
	cmd.Flags().StringP("config-file", "c", tmp.ConfigDir().WithFile("godyl.yml").Path(), "Path to config file")
	cmd.Flags().StringSliceP("env-file", "e", []string{".env"}, "Paths to .env files")
	cmd.Flags().StringP("defaults", "d", "defaults.yml", "Path to defaults file")
	cmd.Flags().String("github-token", env.GetAny("GITHUB_TOKEN", "GH_TOKEN"), "GitHub token for authentication")
	cmd.Flags().String("gitlab-token", env.Get("GITLAB_TOKEN"), "GitLab token for authentication")
	cmd.Flags().String("url-token", env.Get("URL_TOKEN"), "URL token for authentication")
	cmd.Flags().String("url-token-header", "Authorization", "URL token for authentication")
	cmd.Flags().StringP("cache-dir", "", tmp.CacheDir().Path(), "Path to cache directory")
	cmd.Flags().BoolP("no-cache", "", false, "Disable cache")
}

// Tool adds tool-related command flags to the provided cobra command.
// These flags control how tools are downloaded and installed.
func Tool(cmd *cobra.Command) {
	cmd.Flags().StringP("output", "o", "./bin", "Output path for the downloaded tools")
	cmd.Flags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
	cmd.Flags().String("source", "github", "Source from which to install the tools (github, url, go, command)")
	cmd.Flags().String("strategy", tools.None.String(), "Strategy to use for updating tools (none, upgrade, force)")
	cmd.Flags().String("os", "", "Operating system to install the tools for")
	cmd.Flags().String("arch", "", "Architecture to install the tools for")
	cmd.Flags().StringSlice("hints", []string{""}, "Hints to use for tool resolution")
	cmd.Flags().BoolP("no-verify-ssl", "k", false, "Skip SSL verification")
	cmd.Flags().IntP("parallel", "j", 0, "Number of parallel downloads. 0 means unlimited.")
	cmd.Flags().String("version", "", "Version of the tool to install. Empty means latest. Obviously not so useful when downloading multiple tools.")
	cmd.Flags().BoolP("show", "s", false, "Show the configuration and exit")
	cmd.Flags().Bool("no-progress", false, "Disable progress bar")
}

// Update adds update-related command flags to the provided cobra command.
// These flags control the self-update.
func Update(cmd *cobra.Command) {
	cmd.Flags().BoolP("no-verify-ssl", "k", false, "Skip SSL verification")
	cmd.Flags().String("version", "", "Version of the tool to install. Empty means latest.")
	cmd.Flags().Bool("pre", false, "Enable pre-release versions")
}

func Cache(cmd *cobra.Command) {
	cmd.Flags().BoolP("delete", "d", false, "Delete the cache")
}
