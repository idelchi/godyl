package core

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/iutils"
)

// CommandPath represents a hierarchical command path as a slice of strings,
// used for building configuration sections and environment variable prefixes.
type CommandPath []string

// Section returns the configuration section prefix for the command path,
// excluding the root command and using dot notation.
func (c CommandPath) Section() iutils.Prefix {
	if len(c) <= 1 {
		return iutils.Prefix("")
	}

	return iutils.Prefix(strings.Join(c.WithoutRoot(), ".")).Lower()
}

// Env returns the environment variable prefix for the command path,
// using underscore notation and uppercase letters.
func (c CommandPath) Env() iutils.Prefix {
	return iutils.Prefix(strings.Join(c, "_")).Upper()
}

// WithoutRoot returns the command path excluding the root command.
func (c CommandPath) WithoutRoot() CommandPath {
	return c[1:]
}

// BuildCommandPath constructs a CommandPath from a Cobra command's full path.
func BuildCommandPath(cmd *cobra.Command) CommandPath {
	return strings.Split(cmd.CommandPath(), " ")
}
