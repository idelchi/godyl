package common

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/iutils"
)

type CommandPath []string

func (c CommandPath) Section() iutils.Prefix {
	if len(c) <= 1 {
		return iutils.Prefix("")
	}

	return iutils.Prefix(strings.Join(c.WithoutRoot(), ".")).Lower()
}

func (c CommandPath) Env() iutils.Prefix {
	return iutils.Prefix(strings.Join(c, "_")).Upper()
}

func (c CommandPath) WithoutRoot() CommandPath {
	return c[1:]
}

func BuildCommandPath(cmd *cobra.Command) CommandPath {
	return strings.Split(cmd.CommandPath(), " ")
}
