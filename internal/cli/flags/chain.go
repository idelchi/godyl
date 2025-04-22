package flags

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/pkg/validator"
)

// ChainPreRun is a helper function to chain the PreRunE functions of a command and its parent.
func ChainPreRun(cmd *cobra.Command, s Viperable, prefix ...string) error {
	if err := cmd.Parent().PreRunE(cmd.Parent(), nil); err != nil {
		return err
	}

	if s == nil {
		return nil
	}

	if err := Bind(cmd, s, prefix...); err != nil {
		return err
	}

	return validator.Validate(s)
}
