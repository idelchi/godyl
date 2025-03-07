package commands

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/pflag"

	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/logger"
)

func flags() {
	// General flags
	pflag.BoolP("help", "h", false, "Show help message and exit")
	pflag.Bool("version", false, "Show version information and exit")

	// Configuration file flags
	pflag.String("dot-env", ".env", "Path to .env file")
	pflag.StringP("defaults", "d", "defaults.yml", "Path to defaults file")

	// Show flags
	pflag.Bool("show-config", false, "Show the parsed configuration and exit")
	pflag.Bool("show-defaults", false, "Show the parsed default configuration and exit")
	pflag.Bool("show-env", false, "Show the parsed environment variables and exit")
	pflag.Bool("show-platform", false, "Detect the platform and exit")

	// Application flags
	pflag.Bool("update", false, "Update the tools")
	pflag.Bool("dump-tools", false, "Dump out default tools.yml as stdout")
	pflag.Bool("dry", false, "Run without making any changes (dry run)")
	pflag.String("log", logger.INFO.String(), "Log level (debug, info, warn, error, silent)")
	pflag.IntP("parallel", "j", runtime.NumCPU(), "Number of parallel downloads. 0 means unlimited.")
	pflag.BoolP("no-verify-ssl", "k", false, "Skip SSL verification")

	// Tool flags
	pflag.String("output", "", "Output path for the downloaded tools")
	pflag.StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
	pflag.String("source", string(sources.GITHUB), "Source from which to install the tools")
	pflag.String("strategy", "none", "Strategy to use for updating tools")
	pflag.String("github-token", "", "GitHub token for authentication")
	pflag.String("os", "", "Operating system to install the tools for")
	pflag.String("arch", "", "Architecture to install the tools for")

	pflag.CommandLine.SortFlags = false
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [tools]\n\n", "godyl")
		fmt.Fprintf(os.Stderr, "Tool manager that installs tools as specified in a YAML file.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		pflag.PrintDefaults()
	}
}

// This file is deprecated. Use helpers.go instead.

// Deprecated functions - use helpers.go instead
